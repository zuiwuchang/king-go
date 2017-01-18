//提供了 常用的 字符串 處理函數
package strings

import (
	"regexp"
	"strconv"
)

func isIPv4Node(str string) bool {
	n := len(str)
	if n == 2 && str[0] == '0' {
		return false
	} else if n == 3 {
		if str[0] == '0' {
			return false
		}
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false
		}
		if v > 255 {
			return false
		}
	}
	return true
}

//返回 是否是 IP v4 字符串
func IsIPv4(ip string) bool {
	match, _ := regexp.Compile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
	strs := match.FindStringSubmatch(ip)
	if strs == nil {
		return false
	}
	return isIPv4Node(strs[1]) && isIPv4Node(strs[2]) && isIPv4Node(strs[3]) && isIPv4Node(strs[4])
}
