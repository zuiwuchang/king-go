//go 標準庫的 os.signal 無法接收 win32 的 raise 產生的 signal 此庫完成此功能
package signal

//#cgo LDFLAGS:
//#include <signal.h>
//extern void _king_go_os_signal_goHandler(int);
/*
#if defined(_WIN32) || defined(WIN32) || defined(_WIN64) || defined(WIN64) || defined(__WIN32__) || defined(__TOS_WIN__) || defined(__WINDOWS__)
void _king_go_os_signal_handler(int sig)
{
	_king_go_os_signal_goHandler(sig);
	signal(sig,_king_go_os_signal_handler);
}
#endif
void _king_go_os_signal_signals()
{
#if defined(_WIN32) || defined(WIN32) || defined(_WIN64) || defined(WIN64) || defined(__WIN32__) || defined(__TOS_WIN__) || defined(__WINDOWS__)
	//註冊 信號處理函數
	signal(SIGINT,_king_go_os_signal_handler);		//通常是 ctrl + c
	signal(SIGILL,_king_go_os_signal_handler);
	signal(SIGFPE,_king_go_os_signal_handler);
	signal(SIGSEGV,_king_go_os_signal_handler);
	signal(SIGTERM,_king_go_os_signal_handler);	//軟體退出 kill
	signal(SIGBREAK,_king_go_os_signal_handler);	//通常是 ctrl + break/pause /點擊了控制檯窗口的 關閉按鈕
	signal(SIGABRT,_king_go_os_signal_handler);	//調用了 abort 函數
#endif
}
*/
import "C"

import (
	"os"
	osSignal "os/signal"
	"runtime"
)

func doInit() {
	if g_handlers.m != nil {
		return
	}

	g_handlers.m = make(map[chan<- os.Signal]int)

	C._king_go_os_signal_signals()
}

func Notify(c chan<- os.Signal, sig ...os.Signal) {
	if runtime.GOOS == "windows" {
		g_handlers.Lock()
		doInit()
		g_handlers.m[c] = 1
		g_handlers.Unlock()
	} else {
		osSignal.Notify(c, sig...)
	}
}
func Stop(c chan<- os.Signal) {
	if runtime.GOOS == "windows" {
		g_handlers.Lock()
		delete(g_handlers.m, c)
		g_handlers.Unlock()
	} else {
		osSignal.Stop(c)
	}
}
