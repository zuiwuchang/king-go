package pool

import (
	"container/list"
	"errors"
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
	} else if n < 0 {
		//縮容
		p.resizing = true
		go p.resizeLess(-n)
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
	} else if n < 0 {
		//縮容
		p.resizing = true
		go p.resizeLess(-n)
	}

	//重置 執行 時間
	p.timer.Reset(p.t.Interval())
}

//擴容
func (p *poolImpl) resizeMore(n int) {
	for ; n > 0; n-- {
		c, e := p.t.Conect()
		if e != nil {
			p.resizing = false
			return
		}
		if nil != p.put(c) {
			p.resizing = false
			return
		}
	}
	p.lastResize = time.Now()
	p.resizing = false
}
func (p *poolImpl) put(c0 net.Conn) error {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.run {
		return errors.New("pool not run")
	}

	c := newConn(c0, _StatusIdle)

	//創建 ping
	d := p.t.PingInterval()
	if d > 0 {
		p.executePing(c, d)
	}
	//創建 timeout
	d = p.t.Timeout()
	if d > 0 {
		p.executeTimeout(c, d)
	}

	//加入 連接池
	p.l.PushBack(c)
	return nil
}

//縮容
func (p *poolImpl) resizeLess(n int) {
	p.m.Lock()
	defer p.m.Unlock()

	var next *list.Element
	for pos := p.l.Front(); pos != nil && n > 0; pos = next {
		next = pos.Next()

		c := pos.Value.(*Conn)
		if c.status != _StatusIdle {
			continue
		}

		//移除節點
		p.l.Remove(pos)
		c.free(p.t)
		n--
	}

	p.lastResize = time.Now()
	p.resizing = false
}

//爲連接 執行 超時 關閉 操作
func (p *poolImpl) executeTimeout(c *Conn, d time.Duration) {
	if c.timeout == nil {
		c.timeout = time.AfterFunc(d, func() {
			p.m.Lock()
			defer p.m.Unlock()

			if c.status == _StatusIdle {
				//超時 關閉
				c.Conn.Close()
				for iter := p.l.Front(); iter != nil; iter = iter.Next() {
					c0 := iter.Value.(*Conn)
					if c0 == c {
						p.l.Remove(iter)
						p.resize()
						break
					}
				}
			} else if c.status == _StatusPing {
				//正在ping中 延遲 執行 timeout
				c.timeout.Reset(time.Minute)
			}
		})
	} else {
		//已有 timer 直接 執行之
		c.timeout.Reset(d)
	}
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
func (p *poolImpl) Put(c *Conn) {
	p.m.Lock()
	defer p.m.Unlock()

	//減少 使用 計數
	p.use--

	if c.status == _StatusClose {
		//釋放連接
		c.free(p.t)
		return
	}

	if p.run {
		//創建 ping
		d := p.t.PingInterval()
		if d > 0 {
			p.executePing(c, d)
		}
		//創建 timeout
		d = p.t.Timeout()
		if d > 0 {
			p.executeTimeout(c, d)
		}

		//加入 連接池
		c.lastPut = time.Now()
		c.status = _StatusIdle
		p.l.PushBack(c)

		//調整 容量
		p.resize()
	} else {
		//工作 已停止 直接釋放 連接
		c.free(p.t)
	}

}

//從 連接池 中 獲取一個連接
func (p *poolImpl) Get() (*Conn, error) {
	p.m.Lock()
	defer p.m.Unlock()

	//查找 已有連接
	ele := p.l.Back()
	var c *Conn
	for ; ele != nil; ele = ele.Prev() {
		c0 := ele.Value.(*Conn)
		if c0.status == _StatusIdle {
			c = c0
			break
		}
	}
	if c != nil { //找到 可用 連接 返回之
		//從列表 移除
		p.l.Remove(ele)

		//停止 timer
		if c.timer != nil {
			c.timer.Stop()
		}
		if c.timeout != nil {
			c.timeout.Stop()
		}
		//返回 連接
		c.status = _StatusGet

		//增加 計數
		p.use++

		//調整 容量
		p.resize()
		return c, nil
	}

	//沒有 空閒 連接 創建 新連接
	c0, e := p.t.Conect()
	if e != nil {
		return nil, e
	}
	//增加 計數
	p.use++

	//執行 自動 擴容
	p.resize()
	return newConn(c0, _StatusGet), nil
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

	p.run = false

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
