package reuse

type Description struct {
	RecvWindow uint32
	RecvBuffer uint32
}

func (d *Description) Format() {
	if d.RecvWindow < 1024*4 {
		if kLog.Warn != nil {
			kLog.Warn.Printf("RecvWindow(%v) < 1024*64 use default(%v)",
				d.RecvWindow,
				DefaultRecvWindow,
			)
		}

		d.RecvWindow = DefaultRecvWindow
	}
	if d.RecvBuffer < 1024 {
		if kLog.Warn != nil {
			kLog.Warn.Printf("RecvBuffer(%v) < 1024*16 use default(%v)",
				d.RecvBuffer,
				DefaultRecvBuffer,
			)
		}

		d.RecvBuffer = DefaultRecvBuffer
	}
}
