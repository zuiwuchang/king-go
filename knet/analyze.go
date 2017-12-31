package knet

import (
	"encoding/binary"
	kio "github.com/zuiwuchang/king-go/io"
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
		if kLog.Error != nil {
			kLog.Error.Println(ErrorHeaderSize)
		}
		return 0, ErrorHeaderSize
	}

	//flag
	flag := binary.LittleEndian.Uint16(header[DefaultFlagPos:])
	if DefaultFlag != flag {
		if kLog.Error != nil {
			kLog.Error.Println(ErrorHeaderFlag, flag, "!=", DefaultFlag)
		}

		return 0, ErrorHeaderFlag
	}

	//cmd
	if header[DefaultCommandPos] == 0 && header[DefaultCommandPos+1] == 0 {
		if kLog.Error != nil {
			kLog.Error.Println(ErrorHeaderCommand, "0")
		}

		return 0, ErrorHeaderCommand
	}

	//len
	n := int(binary.LittleEndian.Uint16(header[DefaultLenPos:]))
	if n < 6 {
		if kLog.Error != nil {
			kLog.Error.Println(ErrorMessageSize, n)
		}
		return 0, ErrorMessageSize
	}

	return n, nil
}

//使用 默認協議 發送一條消息
func WriteMessage(w io.Writer, b []byte) (n int, e error) {
	header := make([]byte, 2)
	binary.LittleEndian.PutUint16(header, uint16(len(b)+2))
	n, e = kio.WriteAllEx(w, header)
	if e != nil {
		if kLog.Warn != nil {
			kLog.Warn.Println(e)
		}
		return
	}
	var sum int
	sum, e = kio.WriteAllEx(w, b)
	if sum != 0 {
		n += sum
	}
	if e != nil && kLog.Warn != nil {
		kLog.Warn.Println(e)
	}
	return

}
