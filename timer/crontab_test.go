package timer

import (
	"testing"
	"time"
)

func TestCrontabIsOk(t *testing.T) {
	cb := Crontab{Minute: CrontabAny,
		Hour:  CrontabAny,
		Day:   CrontabAny,
		Month: CrontabAny,
		Week:  CrontabAny,
	}
	e := cb.IsOk()
	if e != nil {
		t.Fatal(e)
	}

	//Minute
	cb.Minute = -2
	e = cb.IsOk()
	if e.Error() != crontabEmsgMinute {
		t.Fatal("Minute small than 0 not work")
	}
	cb.Minute = 60
	e = cb.IsOk()
	if e.Error() != crontabEmsgMinute {
		t.Fatal("Minute large than 59 not work")
	}
	cb.Minute = CrontabAny

	//Hour
	cb.Hour = -2
	e = cb.IsOk()
	if e.Error() != crontabEmsgHour {
		t.Fatal("Hour small than 0 not work")
	}
	cb.Hour = 24
	e = cb.IsOk()
	if e.Error() != crontabEmsgHour {
		t.Fatal("Hour large than 23 not work")
	}
	cb.Hour = CrontabAny

	//Day
	cb.Day = 0
	e = cb.IsOk()
	if e.Error() != crontabEmsgDay {
		t.Fatal("Day small than 1 not work")
	}
	cb.Day = 32
	e = cb.IsOk()
	if e.Error() != crontabEmsgDay {
		t.Fatal("Day large than 31 not work")
	}
	cb.Day = CrontabAny

	//Month
	cb.Month = 0
	e = cb.IsOk()
	if e.Error() != crontabEmsgMonth {
		t.Fatal("Month small than 1 not work")
	}
	cb.Month = 13
	e = cb.IsOk()
	if e.Error() != crontabEmsgMonth {
		t.Fatal("Month large than 12 not work")
	}
	cb.Month = CrontabAny

	//Week
	cb.Week = -2
	e = cb.IsOk()
	if e.Error() != crontabEmsgWeek {
		t.Fatal("Week small than 0 not work")
	}
	cb.Week = 8
	e = cb.IsOk()
	if e.Error() != crontabEmsgWeek {
		t.Fatal("Week large than 7 not work")
	}
	cb.Week = CrontabAny

	//min
	cb.Minute = 0
	cb.Hour = 0
	cb.Day = 1
	cb.Month = 1
	cb.Week = 0
	e = cb.IsOk()
	if e != nil {
		t.Fatal(e)
	}

	//max
	cb.Minute = 59
	cb.Hour = 23
	cb.Day = 31
	cb.Month = 12
	cb.Week = 7
	e = cb.IsOk()
	if e != nil {
		t.Fatal(e)
	}

}
func TestCrontabString(t *testing.T) {
	cb := Crontab{Minute: 59,
		Hour:  23,
		Day:   31,
		Month: CrontabAny,
		Week:  7,
	}

	//String
	cb.Month = CrontabAny
	str, _ := cb.String()
	if str != "59	23	31	*	7" {
		t.Fatal("String() not work")
	}

	var cb1 Crontab
	e := cb1.FromString(str)
	if e != nil {
		t.Fatal(e)
	}

	str1, _ := cb1.String()
	if str1 != str {
		t.Fatal("FromString not work")
	}
}
func TestCrontabShouldRun(t *testing.T) {

	cb := Crontab{Minute: CrontabAny,
		Hour:  CrontabAny,
		Day:   CrontabAny,
		Month: CrontabAny,
		Week:  CrontabAny,
	}
	if !cb.isShouldRun() {
		t.Fatal("isShouldRun not work CrontabAny")
	}

	now := time.Now()
	cb.Minute = now.Minute()
	cb.Hour = now.Hour()
	cb.Day = now.Day()
	cb.Month = int(now.Month())
	cb.Week = int(now.Weekday())
	if !cb.isShouldRun() {
		t.Fatal("isShouldRun not work")
	}

	now = time.Now()
	cb.Minute = now.Minute()
	cb.Hour = now.Hour() - 1
	cb.Day = now.Day()
	cb.Month = int(now.Month())
	cb.Week = int(now.Weekday())
	if cb.isShouldRun() {
		t.Fatal("isShouldRun not work")
	}
}
func TestCrontab(t *testing.T) {
	cb, e := NewCrontab(CrontabAny, CrontabAny, CrontabAny, CrontabAny, CrontabAny)
	if e != nil {
		t.Fatal(e)
	}

	v := 0
	cb.Run(func(c *Crontab) {
		v++
		//fmt.Println(v, time.Now())
		c.Close()
	})
	cb.Wait()

	if v != 1 {
		t.Fatal("error")
	}

}
