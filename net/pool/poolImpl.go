package pool

import (
	"container/list"
	"net"
	"sync"
	"time"
)

//連接池 實現
type poolImpl struct {
	//連接池是否 開始工作
	run bool

	//連接池 模板
	t IPoolTemplate

	//空閒的 連接
	l *list.List

	//同步 對象
	m sync.Mutex

	//上次 釋放時間
	lastFree time.Time
}

//從 連接池 中 獲取一個連接
func (p *poolImpl) Get() (net.Conn, error) {
	p.m.Lock()
	defer p.m.Unlock()

	//沒有 空閒 連接 創建 新連接
	if p.l.Len() == 0 {
		//執行 自動 擴容
		//****

		c, e := p.t.Conect()
		if e != nil {
			return nil, e
		}
		return newConn(c), nil
	}

	//返回 節點
	element := p.l.Back()
	c := element.Value.(*Conn)
	p.l.Remove(element)

	if p.l.Len() == 0 {
		//執行 自動 擴容
		//****
	}

	//停止 timer
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	//返回 連接
	return c, nil
}

//將 連接 返回給 連接池
func (p *poolImpl) Put(c net.Conn) {
	p.m.Lock()
	defer p.m.Unlock()

	c0 := c.(*Conn)
	if !c0.ok {
		//連接已經 失效 釋放連接
		c0.free(p.t)
		return
	}

	if p.run {
		//創建 ping
		duration := p.t.PingInterval()
		if duration > 0 {
			c0.timer = time.AfterFunc(duration, func() {
				//*****
			})
		}
		//加入 連接池
		c0.lastPut = time.Now()
		p.l.PushBack(c0)

		//需要 釋放 資源
		n := p.l.Len()
		if n > 1 && n > p.t.MinConn() {
			if p.lastFree.Add(p.t.MinFreeInterval()).After(time.Now()) { //超過了最小執行 週期
				if p.t.MinConn() == 0 {
					p.free(n - 1)
				} else {
					p.free(n - p.t.MinConn())
				}
			}
		}
	} else {
		//工作 已停止 直接釋放 連接
		c0.free(p.t)
	}
}

//釋放 超時資源
func (p *poolImpl) free(n int) {
	p.lastFree = time.Now()
	beginTime := p.lastFree.Add(-p.t.Timeout())

	//遍歷 節點 跳過 最後一個 節點
	pos := p.l.Front()
	next := pos.Next()
	for next != nil && n > 0 {
		c := pos.Value.(*Conn)
		//超時 刪除 此節點
		if c.lastPut.Before(beginTime) {
			//停止 定時器
			if c.timer != nil {
				c.timer.Stop()
				c.timer = nil
			}

			//刪除 節點
			p.l.Remove(pos)

			//減少 計數
			n--
		} else {
			//後面的 都不可能 超時 直接 跳出
			break
		}

		pos = next
		next = pos.Next()
	}
}

//釋放 連接池 所有資源 此後 連接池不能再被使用
func (p *poolImpl) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.run { //沒有工作 直接 返回
		return
	}

	//關閉 所有 連接
	for iter := p.l.Front(); iter != nil; iter = iter.Next() {
		c := iter.Value.(*Conn)
		//停止 timer
		if c.timer != nil {
			c.timer.Stop()
			c.timer = nil
		}
		//釋放 連接
		c.free(p.t)
	}
}
