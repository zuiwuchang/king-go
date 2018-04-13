package strings

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// GMailUserMinLen .
	GMailUserMinLen = 6
	// GMailUserMaxLen .
	GMailUserMaxLen = 30
)

// ErrMatchGMailSplitLess .
var ErrMatchGMailSplitLess = errors.New(`not match @`)

// ErrMatchGMailSplitLarge .
var ErrMatchGMailSplitLarge = errors.New(`@ large than 1`)

// ErrMatchGMailSplit .
var ErrMatchGMailSplit = errors.New(`@ at first or last`)

// ErrMatchGMailUserLess .
var ErrMatchGMailUserLess = fmt.Errorf(`user name requires at least %v characters`, GMailUserMinLen)

// ErrMatchGMailUserLarge .
var ErrMatchGMailUserLarge = fmt.Errorf(`user name up to %v characters`, GMailUserMaxLen)

// ErrMatchGMailUserBadBegin .
var ErrMatchGMailUserBadBegin = errors.New(`user name must start with [a-zA-Z0-9]`)

// ErrMatchGMailUserBadEnd .
var ErrMatchGMailUserBadEnd = errors.New(`user name must end with [a-zA-Z0-9\.]`)

// ErrMatchGMailUserPointLink .
var ErrMatchGMailUserPointLink = errors.New(`usernames can't have consecutive periods`)

// ErrMatchGMailBadHost .
var ErrMatchGMailBadHost = errors.New("bad host name")

// MatchGMail 匹配 字符串是否符合 gmail 格式 要求
/*	用戶名@主機名
	用戶名
		字符[a-zA-Z0-9]長度爲 [6,30]
		以 [a-zA-Z0-9] 開始 結束
		中間只能是 [a-zA-Z0-9\.] 且 . 和 . 不能相連
	域名
		以 .[a-zA-Z]{2,} 結尾
		不能以 . 開頭 且 . 不能相連
		由 [a-zA-Z0-9\-_\.]組成
*/
func MatchGMail(str string) (e error) {
	//matchSplite
	var strs = strings.Split(str, "@")
	var len = len(strs)
	if len < 2 {
		e = ErrMatchGMailSplitLess
		return
	} else if len > 2 {
		e = ErrMatchGMailSplitLarge
		return
	} else if strs[0] == "" || strs[1] == "" {
		e = ErrMatchGMailSplit
		return
	}

	e = matchGMailUser(strs[0])
	if e != nil {
		return
	}
	return matchGMailHost(strs[1])
}

var matchGMailRegexBegin, _ = regexp.Compile(`^[a-zA-Z0-9]`)
var matchGMailRegexEnd, _ = regexp.Compile(`[a-zA-Z0-9]$`)

func matchGMailUser(user string) (e error) {
	if len(user) < GMailUserMinLen {
		e = ErrMatchGMailUserLess
		return
	} else if len(user) > GMailUserMaxLen {
		e = ErrMatchGMailUserLarge
		return
	}
	if !matchGMailRegexBegin.MatchString(user) {
		e = ErrMatchGMailUserBadBegin
		return
	}
	if !matchGMailRegexEnd.MatchString(user) {
		e = ErrMatchGMailUserBadEnd
		return
	}
	if strings.Index(user, "..") != -1 {
		e = ErrMatchGMailUserPointLink
		return
	}
	if len(user)-strings.Count(user, ".") < GMailUserMinLen {
		e = ErrMatchGMailUserLess
	}
	return
}

var matchGMailRegexHostName, _ = regexp.Compile(`^[a-zA-Z0-9\-_][a-zA-Z0-9\-_\.]*\.[a-zA-Z0-9]{2,}$`)

func matchGMailHost(name string) (e error) {
	if !matchGMailRegexHostName.MatchString(name) {
		e = ErrMatchGMailBadHost
		return
	}
	if strings.Index(name, "..") != -1 {
		e = ErrMatchGMailBadHost
		return
	}
	return
}
