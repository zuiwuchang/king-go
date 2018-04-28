package set

import (
	"bytes"
	"fmt"
)

// String .
type String map[string]bool

// NewString 創建 一個 string 的 set
func NewString() String {
	return make(String)
}

// Insert 插入節點
func (s String) Insert(key string) {
	s[key] = true
}

// Remove 移除節點
func (s String) Remove(key string) {
	delete(s, key)
}

// Ok 返回節點是否 存在
func (s String) Ok(key string) bool {
	_, ok := s[key]
	return ok
}

// Len 返回節點 數量
func (s String) Len() int {
	return len(s)
}
func (s String) String() string {
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
