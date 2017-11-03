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

	//上次 縮擴容 時間
	lastResize time.Time
	//擴容 縮容 定時器
	timer *time.Timer
	//是否正在 調整容量中
	resizing bool

	//已經 被 get 使用 的連接數量
	use int
}

//執行 自動 縮擴容操作
func (p *poolImpl) resizeTimer() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.resizing { //正在調整 容量中
		p.timer.Reset(p.t.Interval())
		return
	}

	//當前 時間
	now := time.Now()
	//應該執行的 時間
	pt := p.lastResize.Add(p.t.MinInterval())
	if now.Before(pt) {
		//重置 執行 時間
		p.timer.Reset(pt.Sub(now))
		return
	}

	/***	執行 縮擴容	***/
	//返回 是否需要 縮擴容
	n := p.t.Resize(p.use, p.l.Len())
	if n > 0 {
		//擴容
		p.resizing = true
		go p.resizeMore(n)
	} else {
		//縮容
		p.resizing = true
		go p.resizeLess(n)
	}

	//重置 執行 時間
	p.timer.Reset(p.t.Interval())
}

//自動 計算 縮擴容
func (p *poolImpl) resize() {
	if p.resizing { //正在調整 容量中
		return
	}

	//當前 時間
	now := time.Now()
	//應該執行的 時間
	pt := p.lastResize.Add(p.t.MinInterval())
	if now.Before(pt) { //不需要執行
		return
	}

	/***	執行 縮擴容	***/
	//返回 是否需要 縮擴容
	n := p.t.Resize(p.use, p.l.Len())
	if n > 0 {
		//擴容
		p.resizing = true
		go p.resizeMore(n)
	} else {
		//縮容
		p.resizing = true
		go p.resizeLess(n)
	}

	//重置 執行 時間
	p.timer.Reset(p.t.Interval())
}

//擴容
func (p *poolImpl) resizeMore(n int) {
	p.m.Lock()
	defer p.m.Unlock()

	p.lastResize = time.Now()
	p.resizing = false
}

//縮容
func (p *poolImpl) resizeLess(n int) {
	p.m.Lock()
	defer p.m.Unlock()

	p.lastResize = time.Now()
	p.resizing = false
}

//從 連接池 中 獲取一個連接
func (p *poolImpl) Get() (net.Conn, error) {
	p.m.Lock()
	defer p.m.Unlock()

	//查找 已有連接
	element := p.l.Back()
	var c *Conn
	for ; element != nil; element = element.Prev() {
		c0 := element.Value.(*Conn)
		if c0.status == _StatusIdle {
			c = c0
			break
		}
	}
	if c != nil { //找到 可用 連接 返回之
		//從列表 移除
		p.l.Remove(element)

		//停止 timer
		if c.timer != nil {
			c.timer.Stop()
		}
		//返回 連接
		c.status = _StatusGet

		//自動 縮擴容
		p.resize()
		return c, nil
	}

	//沒有 空閒 連接 創建 新連接
	/*if p.l.Len() == 0 {
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
	*/
	return c, nil
}

//爲連接 執行 ping 操作
func (p *poolImpl) executePing(c *Conn, d time.Duration) {
	if c.timer == nil {
		//沒有 timer 創建之
		c.timer = time.AfterFunc(d, func() {
			//驗證 狀態 是否 空閒
			p.m.Lock()
			if c.status != _StatusIdle { //已經分配工作 不需要 ping
				p.m.Unlock()
				return
			}
			c.status = _StatusPing
			p.m.Unlock()

			//執行 ping
			e := p.t.Ping(c)
			p.m.Lock()
			defer p.m.Unlock()
			if e == nil {
				//重啓 ping
				c.timer.Reset(d)
			} else {
				//節點已失效 移除
				for ele := p.l.Back(); ele != nil; ele = ele.Prev() {
					if ele.Value.(*Conn) == c {
						p.l.Remove(ele)
						break
					}
				}

				p.resize()
			}
		})
	} else {
		//已有 timer 直接 執行之
		c.timer.Reset(d)
	}
}

//將 連接 返回給 連接池
func (p *poolImpl) Put(c net.Conn) {
	/*	p.m.Lock()
		defer p.m.Unlock()

		c0 := c.(*Conn)
		if c0.status == _StatusClose {
			//釋放連接
			c0.free(p.t)
			return
		}

		if p.run {
			//創建 ping
			duration := p.t.PingInterval()
			if duration > 0 {
				p.executePing(c0, duration)
			}
			//加入 連接池
			c0.lastPut = time.Now()
			c0.status = _StatusIdle
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
	*/
}

//釋放 超時資源
func (p *poolImpl) free(n int) {
	/*p.lastFree = time.Now()
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
	*/
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
