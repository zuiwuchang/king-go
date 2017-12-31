package knet

import (
	"errors"
)

var ErrorHeaderSize error = errors.New("mesaage.header size error")
var ErrorHeaderFlag error = errors.New("mesaage.header flag not match")
var ErrorHeaderCommand error = errors.New("mesaage command not support")
var ErrorMessageSize error = errors.New("mesaage size error")

var ErrorBadReadStatus error = errors.New("socket read channel bad , please close socket .")
var ErrorBadIConnWriteTo error = errors.New("IConn Bad WriteTo , this message already miss .")
