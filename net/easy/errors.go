package easy

import (
	"errors"
)

// ErrorHeaderSize 解析包 長度 錯誤
var ErrorHeaderSize = errors.New("mesaage.header size error")

// ErrorHeaderFlag 解析包 標記 不匹配
var ErrorHeaderFlag = errors.New("mesaage.header flag not match")

// ErrorHeaderCommand 解析包 指令 錯誤
var ErrorHeaderCommand = errors.New("mesaage command not support")

// ErrorMessageSize 解析 包長度 小於包頭
var ErrorMessageSize = errors.New("mesaage size error")

// ErrorReadTimeout 讀取消息 超時
var ErrorReadTimeout = errors.New("tcp read timeout")

// ErrorReadChannelClosed read goroutine 已經關閉
var ErrorReadChannelClosed = errors.New("tcp read channel closed")

// ErrorWriteTimeout 寫入消息 超時
var ErrorWriteTimeout = errors.New("tcp write timeout")
