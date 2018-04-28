package easy

import (
	"encoding/binary"
	"io"
)

type analyze struct {
	headerSize  int
	analyzeFunc func(header []byte) (int, error)
}

func (a *analyze) Header() int {
	return a.headerSize
}
func (a *analyze) Analyze(header []byte) (int, error) {
	return a.analyzeFunc(header)
}

type defaultAnalyze struct {
}

func (defaultAnalyze) Header() int {
	return DefaultHeaderLen
}
func (defaultAnalyze) Analyze(header []byte) (int, error) {
	//header
	if len(header) != DefaultHeaderLen {

		return 0, ErrorHeaderSize
	}

	//flag
	flag := binary.LittleEndian.Uint16(header[DefaultFlagPos:])
	if DefaultFlag != flag {

		return 0, ErrorHeaderFlag
	}

	//cmd
	if header[DefaultCommandPos] == 0 && header[DefaultCommandPos+1] == 0 {

		return 0, ErrorHeaderCommand
	}

	//len
	n := int(binary.LittleEndian.Uint16(header[DefaultLenPos:]))
	if n < 6 {
		return 0, ErrorMessageSize
	}

	return n, nil
}

// WriteMessage 使用 默認協議 發送一條消息
func WriteMessage(w io.Writer, cmd uint16, body []byte) (msg []byte, n int, e error) {
	msg = make([]byte, DefaultHeaderLen+len(body))
	copy(msg[DefaultHeaderLen:], body)
	// flag
	binary.LittleEndian.PutUint16(msg[DefaultFlagPos:], DefaultFlag)
	// cmd
	binary.LittleEndian.PutUint16(msg[DefaultCommandPos:], cmd)
	// size
	binary.LittleEndian.PutUint16(msg[DefaultLenPos:], (uint16)(len(msg)))

	n, e = w.Write(msg)
	return
}

// FormatMessage 使用默認 協議 填充 消息頭
func FormatMessage(cmd uint16, msg []byte) {
	// flag
	binary.LittleEndian.PutUint16(msg[DefaultFlagPos:], DefaultFlag)
	// cmd
	binary.LittleEndian.PutUint16(msg[DefaultCommandPos:], cmd)
	// size
	binary.LittleEndian.PutUint16(msg[DefaultLenPos:], (uint16)(len(msg)))
}

// GetMessageCommand 返回 默認協議的 消息 命令
func GetMessageCommand(msg []byte) uint16 {
	return binary.LittleEndian.Uint16(msg[DefaultCommandPos:])
}
