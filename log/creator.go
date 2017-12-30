package log

import (
	"github.com/fatih/color"
	"io"
	"log"
)

//日誌 創建器
type Creator struct {
	//標記 需要輸出的 內容
	Flags int

	//是否要 爲日誌 着色
	Color bool

	//日誌 前綴
	Tag string
}

//創建一個 默認的 Creator
func NewCreator() *Creator {
	return &Creator{
		Flags: log.LstdFlags | log.Lshortfile,
		Color: true,
		Tag:   "",
	}
}

//創建 日誌
func (c *Creator) New(out io.Writer, tag string) *log.Logger {
	return log.New(
		out,
		tag,
		c.Flags,
	)
}

//創建 日誌
func (c *Creator) NewColor(out io.Writer, tag string, cr *color.Color) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  cr,
			Writer: out,
		},
		tag,
		c.Flags,
	)
}
func (c *Creator) newAuto(out io.Writer, tag string, cr *color.Color) *log.Logger {
	if c.Color && cr != nil {
		return c.NewColor(out, tag, cr)
	}
	return c.New(out, tag)
}

//創建 trace 日誌
func (c *Creator) NewTrace(out io.Writer) *log.Logger {
	return c.newAuto(
		out,
		"[Trace] "+c.Tag,
		color.New(color.FgCyan),
	)
}

//創建 debug 日誌
func (c *Creator) NewDebug(out io.Writer) *log.Logger {
	return c.newAuto(
		out,
		"[Debug] "+c.Tag,
		color.New(color.FgBlue),
	)
}

//創建 info 日誌
func (c *Creator) NewInfo(out io.Writer) *log.Logger {
	return c.newAuto(
		out,
		"[Info] "+c.Tag,
		color.New(color.FgGreen),
	)
}

//創建 warn 日誌
func (c *Creator) NewWarn(out io.Writer) *log.Logger {
	return c.newAuto(
		out,
		"[Warn] "+c.Tag,
		color.New(color.FgYellow),
	)
}

//創建 error 日誌
func (c *Creator) NewError(out io.Writer) *log.Logger {
	return c.newAuto(
		out,
		"[Error] "+c.Tag,
		color.New(color.FgMagenta),
	)

}

//創建 fault 日誌
func (c *Creator) NewFault(out io.Writer) *log.Logger {
	return c.newAuto(
		out,
		"[Fault] "+c.Tag,
		color.New(color.FgRed),
	)
}
