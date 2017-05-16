package timer

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//類似 linux crontab 組件 的定時器
/*
#分鐘		小時		號數		月份		周星期
#[0,59]		[0,23]		[1,31]		[1~12]		[0,7] 0==7==星期天
0		13		*		*		0
*/
const (
	//不限制 項目
	CrontabAny     = -1
	CrontabAnyName = "*"

	//錯誤描述
	crontabEmsgMinute = "minute must at range [0,59]"
	crontabEmsgHour   = "hour must at range [0,23]"
	crontabEmsgDay    = "day must at range [1,31]"
	crontabEmsgMonth  = "month must at range [1,12]"
	crontabEmsgWeek   = "week must at range [0,7]"
)

//Crontab 定時器 回調 函數 返回 true 將 停止定時器
type CrontabCallback func(c *Crontab)
type Crontab struct {
	//分鐘 [0,59]
	Minute int

	//小時 [0,23]
	Hour int

	//號數 [1,31]
	Day int

	//月份 [1,12]
	Month int

	//星期 [0,7] 0 和 7 都是星期天
	Week int

	//結束標記
	chExit chan bool
	//wait 標記
	chWait chan bool
}

//創建一個 Crontab 定時器
func NewCrontab(minute, hour, day, month, week int) (*Crontab, error) {
	cb := Crontab{Minute: minute,
		Hour:  hour,
		Day:   day,
		Month: month,
		Week:  week,
	}
	e := cb.IsOk()
	if e != nil {
		return nil, e
	}
	return &cb, nil
}

//返回 對應的 配置 字符串
func (c *Crontab) String() (string, error) {
	e := c.IsOk()
	if e != nil {
		return "", e
	}

	var buf bytes.Buffer
	arrs := []int{c.Minute,
		c.Hour,
		c.Day,
		c.Month,
		c.Week,
	}
	for i, v := range arrs {
		if i != 0 {
			buf.WriteString("	")
		}
		if v == CrontabAny {
			buf.WriteString("*")
		} else {
			buf.WriteString(fmt.Sprint(v))
		}
	}
	return buf.String(), nil
}

//由配置 字符串 初始化
func (c *Crontab) FromString(str string) error {
	str = strings.TrimSpace(str)
	reg, e := regexp.Compile(`([0-9]{1,2})|(\*)`)
	if e != nil {
		return e
	}
	strs := reg.FindAllString(str, -1)
	if len(strs) < 5 {
		return errors.New("need 5 item for Minute,Hour,Date,Month,Week")
	}

	var minute, hour, day, month, week int
	//minute
	if strs[0] == CrontabAnyName {
		minute = CrontabAny
	} else {
		v, e := strconv.ParseInt(strs[0], 10, 32)
		if e != nil || v < 0 || v > 59 {
			return errors.New(crontabEmsgMinute)
		}
		minute = int(v)
	}

	//hour
	if strs[1] == CrontabAnyName {
		hour = CrontabAny
	} else {
		v, e := strconv.ParseInt(strs[1], 10, 32)
		if e != nil || v < 0 || v > 23 {
			return errors.New(crontabEmsgHour)
		}
		hour = int(v)
	}

	//day
	if strs[2] == CrontabAnyName {
		day = CrontabAny
	} else {
		v, e := strconv.ParseInt(strs[2], 10, 32)
		if e != nil || v < 1 || v > 31 {
			return errors.New(crontabEmsgDay)
		}
		day = int(v)
	}
	//month
	if strs[3] == CrontabAnyName {
		month = CrontabAny
	} else {
		v, e := strconv.ParseInt(strs[3], 10, 32)
		if e != nil || v < 1 || v > 12 {
			return errors.New(crontabEmsgMonth)
		}
		month = int(v)
	}
	//week
	if strs[4] == CrontabAnyName {
		week = CrontabAny
	} else {
		v, e := strconv.ParseInt(strs[4], 10, 32)
		if e != nil || v < 0 || v > 7 {
			return errors.New(crontabEmsgWeek)
		}
		week = int(v)
	}

	c.Minute = minute
	c.Hour = hour
	c.Day = day
	c.Month = month
	c.Week = week
	return nil
}

//返回 配置 是否合法
func (c *Crontab) IsOk() error {
	if c.Minute != CrontabAny &&
		(c.Minute < 0 || c.Minute > 59) {
		return errors.New(crontabEmsgMinute)
	}
	if c.Hour != CrontabAny &&
		(c.Hour < 0 || c.Hour > 23) {
		return errors.New(crontabEmsgHour)
	}
	if c.Day != CrontabAny &&
		(c.Day < 1 || c.Day > 31) {
		return errors.New(crontabEmsgDay)
	}
	if c.Month != CrontabAny &&
		(c.Month < 1 || c.Month > 12) {
		return errors.New(crontabEmsgMonth)
	}
	if c.Week != CrontabAny &&
		(c.Week < 0 || c.Week > 7) {
		return errors.New(crontabEmsgWeek)
	}
	return nil
}

//運行 定時器
//Crontab 定時器 回調 函數 返回 true 將 停止定時器
func (c *Crontab) Run(callback CrontabCallback) {
	chWait := make(chan bool)
	chExit := make(chan bool, 1)
	c.chWait = chWait
	c.chExit = chExit
	go func() {
		ch := make(chan bool)
		go func() {
			for {
				ch <- true
				time.Sleep(time.Second * 30)
			}
		}()
		c.doSelect(callback, ch, chExit)
		chWait <- true
	}()
}
func (c *Crontab) doSelect(callback CrontabCallback, chWork, chExit chan bool) {
	last := ""
	layout := "2006-01-02 15:04"
	for {
		select {
		case <-chWork:
			if c.isShouldRun() {
				if last != time.Now().Format(layout) {
					callback(c)
					last = time.Now().Format(layout)
				}
			}
		case <-chExit:
			return
		}
	}
}

//停止 定時器
func (c *Crontab) Close() {
	ch := c.chExit
	if ch != nil {
		ch <- true
	}
}

//等待 定時器 結束
func (c *Crontab) Wait() {
	ch := c.chWait
	if ch != nil {
		<-ch
	}
}

//返回是否需要 運行
func (c *Crontab) isShouldRun() bool {
	now := time.Now()
	minute := now.Minute()
	if c.Minute != CrontabAny && c.Minute != minute {
		return false
	}

	hour := now.Hour()
	if c.Hour != CrontabAny && c.Hour != hour {
		return false
	}

	day := now.Day()
	if c.Day != CrontabAny && c.Day != day {
		return false
	}
	month := int(now.Month())
	if c.Month != CrontabAny && c.Month != month {
		return false
	}
	week := int(now.Weekday())
	if c.Week != CrontabAny {
		if c.Week == 7 && week != 0 {
			return false
		}
		if c.Week != week {
			return false
		}
	}
	return true
}
