package signal

import "C"

import (
	"os"
	"sync"
	"syscall"
)

const (
	SIGINT   = 2  /* interrupt */
	SIGILL   = 4  /* illegal instruction - invalid function image */
	SIGFPE   = 8  /* floating point exception */
	SIGSEGV  = 11 /* segment violation */
	SIGTERM  = 15 /* Software termination signal from kill */
	SIGBREAK = 21 /* Ctrl-Break sequence */
	SIGABRT  = 22
)

type handlers struct {
	sync.Mutex
	m map[chan<- os.Signal]int
}

var g_handlers handlers

//export _king_go_os_signal_goHandler
func _king_go_os_signal_goHandler(sig int) {
	g_handlers.Lock()
	defer g_handlers.Unlock()

	switch sig {
	case 0: //CTRL_C_EVENT
		sig = SIGINT
	case 1: //CTRL_BREAK_EVENT
		sig = SIGINT
	case 2: //CTRL_CLOSE_EVENT
		sig = SIGINT
	default:
		return
	}
	for ch, _ := range g_handlers.m {
		ch <- syscall.Signal(sig)
	}
}
