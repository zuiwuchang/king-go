package easy

import (
	"net"
)

type _Listener struct {
	net.Listener
	closed bool
}

// NewListener .
func NewListener(l net.Listener) IListener {
	return &_Listener{l, false}
}
func (l *_Listener) Closed() bool {
	return l.closed
}
func (l *_Listener) Close() (e error) {
	l.closed = true
	e = l.Listener.Close()
	return
}
