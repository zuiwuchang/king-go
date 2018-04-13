package strings

import (
	"testing"
)

func TestMatchGMail(t *testing.T) {

	var str = "king.zuiwuchang@gmail.com"
	e := MatchGMail(str)
	if e != nil {
		t.Fatal(str, e)
	}

	str = "zuiwuchang@--kl._.1.com"
	e = MatchGMail(str)
	if e != nil {
		t.Fatal(str, e)
	}

	str = "zuiwuchang@"
	e = MatchGMail(str)
	if e != ErrMatchGMailSplit {
		t.Fatal(str, e)
	}

	str = "zuiwuch@ang@"
	e = MatchGMail(str)
	if e != ErrMatchGMailSplitLarge {
		t.Fatal(str, e)
	}

	str = ".zuiwuch@ang"
	e = MatchGMail(str)
	if e != ErrMatchGMailUserBadBegin {
		t.Fatal(str, e)
	}

	str = "1234567890123456789012345678901@ang"
	e = MatchGMail(str)
	if e != ErrMatchGMailUserLarge {
		t.Fatal(str, e)
	}

	str = "king@ang"
	e = MatchGMail(str)
	if e != ErrMatchGMailUserLess {
		t.Fatal(str, e)
	}

	str = "k.i.n.g@ang"
	e = MatchGMail(str)
	if e != ErrMatchGMailUserLess {
		t.Fatal(str, e)
	}

	str = "zuiwuch.@ang"
	e = MatchGMail(str)
	if e != ErrMatchGMailUserBadEnd {
		t.Fatal(str, e)
	}

	str = "zuiwu..ch@ang"
	e = MatchGMail(str)
	if e != ErrMatchGMailUserPointLink {
		t.Fatal(str, e)
	}

	str = "zuiwuchang@kl..1.com"
	e = MatchGMail(str)
	if e != ErrMatchGMailBadHost {
		t.Fatal(str, e)
	}

	str = "zuiwuchang@.kl.1.com"
	e = MatchGMail(str)
	if e != ErrMatchGMailBadHost {
		t.Fatal(str, e)
	}
}
