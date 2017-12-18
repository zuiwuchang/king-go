//tcp 複用
package reuse

import (
	"errors"
	"net"
)

//描述 如何進行 複用
type IReuse interface {
	//定義 tcp 每次調用 block socket recv 函數時 的 緩衝區 大小
	Recv() int

	//定義 對復用的 tcp 模擬的 socket write buffer 大小
	Buffer() int
}

//創建一個 默認的 IReuse
func NewReuse() IReuse {
	return reuseImpl{
		recv:   1024 * 8,
		buffer: 1024 * 8 * 10,
	}
}

var ErrorListenerClosed error = errors.New("listener is already closed")

//創建一個 可以複用的 tcp 服務器
func Listen(network, address string) (net.Listener, error) {
	return newListenerImpl(network, address)
}

//創建一個 可以複用的 tcp 服務器
func Listen2(network, address string, reuse IReuse) (net.Listener, error) {
	return newListenerImpl(network, address)
}
