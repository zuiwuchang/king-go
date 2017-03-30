package sdl2

import (
	"log"
	"os"
)

//系統日誌
var g_log *log.Logger

func init() {
	g_log = log.New(os.Stdout, "[king]\t", log.LstdFlags|log.Lshortfile)
}
