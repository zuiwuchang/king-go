package fileperm

import (
	"os"
)

const (
	//other rwx
	OX os.FileMode = 1 << (iota)
	OW
	OR

	//group rwx
	GX
	GW
	GR

	//user rwx
	UX
	UW
	UR
)

//all
const (
	AX = OX | GX | UX
	AW = OW | GW | UW
	AR = OR | GR | UR
)
const (
	ORW = OR | OW
	ORX = OR | OX
	OWX = OW | OX
	OA  = OR | OW | OX

	GRW = GR | GW
	GRX = GR | GX
	GWX = GW | GX
	GA  = GR | GW | GX

	URW = UR | UW
	URX = UR | UX
	UWX = UW | UX
	UA  = UR | UW | UX

	ARW = AR | AW
	ARX = AR | AX
	AWX = AW | AX
	AA  = AR | AW | AX
)

const (
	//drwxr-xr-x
	Directory = os.ModeDir | UA | GRX | ORX

	//-rw-r--r--
	File = URW | GR | OR

	//-rwxr-xr-x
	Execute = UA | GRX | ORX

	//-rwxr-xr-x
	ShellScript = File | AX
)

func SetPerm(args ...os.FileMode) (rs os.FileMode) {
	for _, flag := range args {
		rs |= flag
	}
	return rs
}
func UndoPerm(flags os.FileMode, args ...os.FileMode) os.FileMode {
	if len(args) == 0 {
		return flags
	}

	return flags & (^SetPerm(args...))
}
