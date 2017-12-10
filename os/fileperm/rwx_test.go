package fileperm

import (
	"fmt"
	"testing"
)

func TestFunc(t *testing.T) {
	if 0 != SetPerm() {
		t.Fatal("nil")
	} else if AX != SetPerm(UX, GX, OX) {
		t.Fatal("a+x")
	}

	if AX != UndoPerm(AX) {
		t.Fatal("nil")
	} else if ARX != UndoPerm(AA, UW, GW, UW, GW, OW) {
		t.Fatal("UndoPerm AW")
	}
}
func TestDefault(t *testing.T) {
	if fmt.Sprint(Directory) != "drwxr-xr-x" {
		t.Fatal("bad directory")
	} else if fmt.Sprint(File) != "-rw-r--r--" {
		t.Fatal("bad file")
	} else if fmt.Sprint(Execute) != "-rwxr-xr-x" {
		t.Fatal("bad execute")
	} else if fmt.Sprint(ShellScript) != "-rwxr-xr-x" {
		t.Fatal("bad shell script")
	}
}
func TestRWX(t *testing.T) {
	if fmt.Sprint(OX) != "---------x" {
		t.Fatal("bad o+x")
	} else if fmt.Sprint(OW) != "--------w-" {
		t.Fatal("bad o+w")
	} else if fmt.Sprint(OR) != "-------r--" {
		t.Fatal("bad o+r")
	} else if fmt.Sprint(GX) != "------x---" {
		t.Fatal("bad g+x")
	} else if fmt.Sprint(GW) != "-----w----" {
		t.Fatal("bad g+w")
	} else if fmt.Sprint(GR) != "----r-----" {
		t.Fatal("bad g+r")
	} else if fmt.Sprint(UX) != "---x------" {
		t.Fatal("bad u+x")
	} else if fmt.Sprint(UW) != "--w-------" {
		t.Fatal("bad u+w")
	} else if fmt.Sprint(UR) != "-r--------" {
		t.Fatal("bad u+r")
	}

	if fmt.Sprint(AX) != "---x--x--x" {
		t.Fatal("bad a+x")
	} else if fmt.Sprint(AW) != "--w--w--w-" {
		t.Fatal("bad a+w")
	} else if fmt.Sprint(AR) != "-r--r--r--" {
		t.Fatal("bad a+r")
	}

	if fmt.Sprint(ORW) != "-------rw-" {
		t.Fatal("bad o+rw")
	} else if fmt.Sprint(ORX) != "-------r-x" {
		t.Fatal("bad o+rx")
	} else if fmt.Sprint(OWX) != "--------wx" {
		t.Fatal("bad o+wx")
	} else if fmt.Sprint(OA) != "-------rwx" {
		t.Fatal("bad o+rwx")
	}
	if fmt.Sprint(GRW) != "----rw----" {
		t.Fatal("bad g+rw")
	} else if fmt.Sprint(GRX) != "----r-x---" {
		t.Fatal("bad g+rx")
	} else if fmt.Sprint(GWX) != "-----wx---" {
		t.Fatal("bad g+wx")
	} else if fmt.Sprint(GA) != "----rwx---" {
		t.Fatal("bad g+rwx")
	}
	if fmt.Sprint(URW) != "-rw-------" {
		t.Fatal("bad u+rw")
	} else if fmt.Sprint(URX) != "-r-x------" {
		t.Fatal("bad u+rx")
	} else if fmt.Sprint(UWX) != "--wx------" {
		t.Fatal("bad u+wx")
	} else if fmt.Sprint(UA) != "-rwx------" {
		t.Fatal("bad u+rwx")
	}

	if fmt.Sprint(ARW) != "-rw-rw-rw-" {
		t.Fatal("bad a+rw")
	} else if fmt.Sprint(ARX) != "-r-xr-xr-x" {
		t.Fatal("bad a+rx")
	} else if fmt.Sprint(AWX) != "--wx-wx-wx" {
		t.Fatal("bad a+wx")
	} else if fmt.Sprint(AA) != "-rwxrwxrwx" {
		t.Fatal("bad a+rwx")
	}
}
