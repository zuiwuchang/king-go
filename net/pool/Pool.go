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
	Get() (net.Conn, error)
	//將 連接 返回給 連接池
	Put(net.Conn)

	//釋放 連接池 所有資源 此後 連接池不能再被使用
	Close()
}

//創建一個 連接池
func NewPool(t IPoolTemplate) (IPool, error) {
	min := t.MinConn()
	l := list.New()
	if min > 0 {
		var e error
		var c net.Conn
		for i := 0; i < min; i++ {
			c, e = t.Conect()
			if e != nil {
				break
			}
			l.PushBack(newConn(c))
		}
		if e != nil {
			for iter := l.Front(); iter != nil; iter = iter.Next() {
				iter.Value.(*Conn).free(t)
			}
			return nil, e
		}
	}

	//啓動 定時 縮容
	//****
	return &poolImpl{
			run:      true,
			t:        t,
			l:        l,
			lastFree: time.Now(),
		},
		nil
}
