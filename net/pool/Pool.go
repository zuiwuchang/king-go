//一個 tcp 連接池
package pool

import (
	"container/list"
	"net"
	"time"
)

//連接池 接口
type IPool interface {
	//從 連接池 中 獲取一個連接
	Get() (*Conn, error)
	//將 連接 返回給 連接池
	Put(*Conn)

	//釋放 連接池 所有資源 此後 連接池不能再被使用
	Close()
}

//創建一個 連接池
//sum 初始連接 數量
func NewPool(t IPoolTemplate, sum int) (IPool, error) {
	l := list.New()
	if sum > 0 {
		//創建 連接
		var e error
		var c net.Conn
		for i := 0; i < sum; i++ {
			c, e = t.Conect()
			if e != nil {
				break
			}
			l.PushBack(newConn(c, _StatusIdle))
		}

		if e != nil {
			//出錯 釋放 連接
			for iter := l.Front(); iter != nil; iter = iter.Next() {
				iter.Value.(*Conn).free(t)
			}
			return nil, e
		}
	}

	impl := &poolImpl{
		run:        true,
		t:          t,
		l:          l,
		lastResize: time.Now(),
		use:        0,
		resizing:   false,
	}

	//啓動 定時 縮容
	d := t.Interval()
	if d > 0 {
		impl.timer = time.AfterFunc(0, func() {
			impl.resizeTimer()
		})
	}

	//啓用 ping
	d = t.PingInterval()
	if d > 0 {
		for iter := l.Front(); iter != nil; iter = iter.Next() {
			c := iter.Value.(*Conn)
			impl.executePing(c, d)
		}
	}

	//啓用 超時
	d = t.Timeout()
	if d > 0 {
		for iter := l.Front(); iter != nil; iter = iter.Next() {
			c := iter.Value.(*Conn)
			impl.executeTimeout(c, d)
		}
	}
	return impl, nil
}
