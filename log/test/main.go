package main

import (
	"github.com/zuiwuchang/king-go/log"
	slog "log"
	"os"
)

func main() {
	log.InitDebugLoggers()
	log.Trace.Println("Trace")
	log.Debug.Println("Debug")
	log.Info.Println("Info")
	log.Warn.Println("Warn")
	log.Error.Println("Error")
	log.Fault.Println("Fault")

	out := os.Stdout
	c0 := &log.Creator{
		Flags: slog.Lshortfile,
		Color: false,         //Color false 會使 NewXXX 忽略 Color
		Tag:   "[no color] ", //NewXXX 會在 tag 後加上 Tag
	}
	l0 := log.Loggers{
		Trace: c0.NewTrace(out),
		Debug: c0.NewDebug(out),
		Info:  c0.NewInfo(out),
		Warn:  c0.NewWarn(out),
		Error: c0.NewError(out),
		Fault: c0.NewFault(out),
	}
	l0.Trace.Println("Trace")
	l0.Debug.Println("Debug")
	l0.Info.Println("Info")
	l0.Warn.Println("Warn")
	l0.Error.Println("Error")
	l0.Fault.Println("Fault")

	l1 := log.NewDebugLoggers()
	l1.Trace.Println("Trace")
	l1.Debug.Println("Debug")
	l1.Info.Println("Info")
	l1.Warn.Println("Warn")
	l1.Error.Println("Error")
	l1.Fault.Println("Fault")

	c0 = &log.Creator{
		Flags: slog.LstdFlags,
		Color: true,       //New 會忽略 Color 屬性
		Tag:   "[color] ", //New 會忽略 Tag 屬性 而只使用 New 傳入的 tag
	}
	l0 = log.Loggers{
		Trace: c0.New(out, "[tag0] "),
		Debug: c0.New(out, "[tag1] "),
		Info:  c0.New(out, "[tag2] "),
		Warn:  c0.New(out, "[tag3] "),
		Error: c0.New(out, "[tag4] "),
		Fault: c0.New(out, "[tag5] "),
	}
	l0.Trace.Println("Trace")
	l0.Debug.Println("Debug")
	l0.Info.Println("Info")
	l0.Warn.Println("Warn")
	l0.Error.Println("Error")
	l0.Fault.Println("Fault")

	l1 = log.NewDebugLoggers2("l2 ")
	l1.Trace.Println("Trace")
	l1.Debug.Println("Debug")
	l1.Info.Println("Info")
	l1.Warn.Println("Warn")
	l1.Error.Println("Error")
	l1.Fault.Println("Fault")
}
