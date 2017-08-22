//定時器 相關
package timer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = Hour * 24
	Week        = Day * 7
)

//將時間字符串 轉換到 Duration
//1Hour2Minute
func ToDuration(str string) (time.Duration, error) {
	var duration time.Duration

	reg, e := regexp.Compile(`[0-9]+[a-zA-Z]+`)
	if e != nil {
		return 0, e
	}
	strs := reg.FindAllString(str, -1)
	for _, str := range strs {
		d, e := getDuration(str)
		if e != nil {
			return 0, e
		}
		duration += d
	}

	return duration, nil
}

//將 形如 1week 2HOUR(1Week2hour) 的字符串 轉爲 Duration
func getDuration(str string) (time.Duration, error) {
	pos := strings.IndexFunc(str, func(c rune) bool {
		return c < '0' || c > '9'
	})

	name := strings.ToLower(str[pos:])
	number, _ := strconv.ParseInt(str[:pos], 10, 64)
	switch name {
	case "microsecond":
		return Microsecond * time.Duration(number), nil
	case "millisecond":
		return Millisecond * time.Duration(number), nil
	case "second":
		return Second * time.Duration(number), nil
	case "minute":
		return Minute * time.Duration(number), nil
	case "hour":
		return Hour * time.Duration(number), nil
	case "day":
		return Day * time.Duration(number), nil
	case "week":
		return Week * time.Duration(number), nil
	}
	return 0, fmt.Errorf("unknow item (%s)", str)
}

//將Duration 轉換到 時間字符串
func ToString(duration time.Duration) string {
	first := true
	var buf bytes.Buffer
	if duration >= Week {
		buf.WriteString(fmt.Sprint(int(duration / Week)))
		buf.WriteString("Week")
		duration %= Week

		first = false
	}

	if duration >= Day {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(int(duration / Day)))
		buf.WriteString("Day")
		duration %= Day
	}

	if duration >= Hour {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(int(duration / Hour)))
		buf.WriteString("Hour")
		duration %= Hour
	}

	if duration >= Minute {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(int(duration / Minute)))
		buf.WriteString("Minute")
		duration %= Minute
	}

	if duration >= Second {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(int(duration / Second)))
		buf.WriteString("Second")
		duration %= Second
	}

	if duration >= Millisecond {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(int(duration / Millisecond)))
		buf.WriteString("Millisecond")
		duration %= Millisecond
	}

	if duration >= Microsecond {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprint(int(duration / Microsecond)))
		buf.WriteString("Microsecond")
		duration %= Microsecond
	}

	return buf.String()
}
