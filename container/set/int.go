package set

import (
	"bytes"
	"fmt"
)

// Int .
type Int map[int]bool

// NewInt 創建 一個 Int 的 set
func NewInt() Int {
	return make(Int)
}

// Insert 插入節點
func (s Int) Insert(key int) {
	s[key] = true
}

// Remove 移除節點
func (s Int) Remove(key int) {
	delete(s, key)
}

// Ok 返回節點是否 存在
func (s Int) Ok(key int) bool {
	_, ok := s[key]
	return ok
}

// Len 返回節點 數量
func (s Int) Len() int {
	return len(s)
}
func (s Int) String() string {
	if len(s) == 0 {
		return StringSetEmpty
	} else if len(s) == 1 {
		for key := range s {
			return fmt.Sprintf("[%v]", key)
		}
	}

	buf := bytes.NewBufferString("[")
	first := true
	for key := range s {
		if first {
			first = false
		} else {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprint(key))

	}
	buf.WriteString("]")
	return buf.String()
}
