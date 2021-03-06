package set

import (
	"bytes"
	"fmt"
)

// StringSetEmpty .
var StringSetEmpty = "[]"

// Int64 .
type Int64 map[int64]bool

// NewInt64 創建 一個 Int64 的 set
func NewInt64() Int64 {
	return make(Int64)
}

// Insert 插入節點
func (s Int64) Insert(key int64) {
	s[key] = true
}

// Remove 移除節點
func (s Int64) Remove(key int64) {
	delete(s, key)
}

// Ok 返回節點是否 存在
func (s Int64) Ok(key int64) bool {
	_, ok := s[key]
	return ok
}

// Len 返回節點 數量
func (s Int64) Len() int {
	return len(s)
}
func (s Int64) String() string {
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
