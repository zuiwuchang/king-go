package timer

import (
	"testing"
)

func TestTimer(t *testing.T) {
	str := "2Week 6Day 23Hour 59Minute 59Second 999Millisecond 999Microsecond"
	duration := Week*2 +
		6*Day +
		23*Hour +
		59*Minute +
		59*Second +
		999*Millisecond +
		999*Microsecond
	if str != ToString(duration) {
		t.Fatal("ToString not work")
	}

	duration1, e := ToDuration(str)
	if e != nil {
		t.Fatal(e)
	}
	if duration1 != duration {
		t.Fatal("ToDuration not work")
	}

	duration1, e = ToDuration("1Week6Day1Week23Hour59Minute59Second999Millisecond999Microsecond")
	if e != nil {
		t.Fatal(e)
	}
	if duration1 != duration {

		t.Fatal("ToDuration not work")
	}

	duration1, e = ToDuration("1Week	6Day 1Week 23Hour 59Minute 59Second 999Millisecond 999Microsecond")
	if e != nil {
		t.Fatal(e)
	}
	if duration1 != duration {
		t.Fatal("ToDuration not work")
	}

	duration1, e = ToDuration("1Wek6Day1Week23Hour59Minute59Second999Millisecond999Microsecond")
	if e == nil {
		t.Fatal("ToDuration not work")
	}
}
