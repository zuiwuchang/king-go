//提供了 標準庫 沒有的 字符串 Sub 操作
package strings

//返回 從左向右 的n個 字符
func LeftRune(runes []rune, n int) []rune {
	if n < 1 {
		return nil
	}
	size := len(runes)
	if n > size {
		return runes
	}
	return runes[:n]
}
func Left(str string, n int) (substr string) {
	if n < 1 {
		return
	}

	runes := []rune(str)
	runes = LeftRune(runes, n)
	if runes == nil {
		return
	}

	substr = string(runes)
	return
}

//返回 從右向左 的n個 字符
func RightRune(runes []rune, n int) []rune {
	if n < 1 {
		return nil
	}

	size := len(runes)
	start := size - n
	if start < 0 {
		start = 0
	}
	return runes[start:size]
}
func Right(str string, n int) (substr string) {
	if n < 1 {
		return
	}

	runes := []rune(str)
	runes = RightRune(runes, n)
	if runes == nil {
		return
	}

	substr = string(runes)
	return
}

//使用 索引返回子串	(str,start=0,n=end)
func SubRune(runes []rune, args ...int) []rune {
	start := 0
	n := -1
	lenArgs := len(args)
	if lenArgs < 1 {
		return runes
	}
	start = args[0]
	if start < 0 {
		start = 0
	}

	if lenArgs > 1 {
		n = args[1]
		if n == 0 {
			return nil
		}
	}

	size := len(runes)
	if start >= size {
		return nil
	}
	if n < 1 {
		return runes[start:]
	}

	end := start + n
	if end > size {
		return runes[start:]
	}
	return runes[start:end]
}
func Sub(str string, args ...int) (substr string) {
	runes := []rune(str)

	runes = SubRune(runes, args...)
	if runes == nil {
		return
	}

	substr = string(runes)
	return
}
