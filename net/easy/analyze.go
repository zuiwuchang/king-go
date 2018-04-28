package easy

import (
	"encoding/binary"
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
