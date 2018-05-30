package command

import (
	"fmt"
	"testing"
)

const (
	testErrorVal = 64
)

type testContext struct {
	val64 int64
	val   int
}

func (t *testContext) DoneInt(v int) (e error) {
	if v == testErrorVal {
		e = fmt.Errorf("%v", testErrorVal)
		return
	}
	t.val += v
	return
}
func (t *testContext) Int64(v int64) (e error) {
	if v == testErrorVal {
		e = fmt.Errorf("%v", testErrorVal)
		return
	}
	t.val64 += v
	return
}
func (t *testContext) DoneIntp64(v *int64) (e error) {
	if *v == testErrorVal {
		e = fmt.Errorf("%v", testErrorVal)
		return
	}
	t.val64 += *v
	*v++
	return
}
func TestCommander(t *testing.T) {
	var e error
	c := New()
	hander := &testContext{}
	RegisterCommander(c, hander, "Done")

	if e = c.Execute(int8(8)); e == nil {
		t.Fatal("not Register int8 but execute")
	} else if !IsUnknow(e) {
		t.Fatal(e)
	}
	if e = c.Execute(int64(64)); e == nil {
		t.Fatal("not Register int64 but execute")
	} else if !IsUnknow(e) {
		t.Fatal(e)
	}

	if e = c.Execute(1); e != nil {
		t.Fatal(e)
	}
	if hander.val != 1 {
		t.Fatal("bad ptr hander")
	}
	if e = c.Execute(testErrorVal); e.Error() != fmt.Sprintf("%v", testErrorVal) {
		t.Fatal("testErrorVal", e)
	}

	v64 := int64(testErrorVal)
	if e = c.Execute(&v64); e.Error() != fmt.Sprintf("%v", testErrorVal) {
		t.Fatal("testErrorVal", e)
	}
	v64 = 32
	if e = c.Execute(&v64); e != nil {
		t.Fatal(e)
	}
	if hander.val64 != 32 {
		t.Fatal("bad ptr hander")
	}
	if v64 != 33 {
		t.Fatal("bad ptr param")
	}
}
