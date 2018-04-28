package set

import (
	"bytes"
	"fmt"
)

// Uint16 .
type Uint16 map[uint16]bool

// NewUint16 創建 一個 Uint16 的 set
func NewUint16() Uint16 {
	return make(Uint16)
}

// Insert 插入節點
func (s Uint16) Insert(key uint16) {
	s[key] = true
}

// Remove 移除節點
func (s Uint16) Remove(key uint16) {
	delete(s, key)
}

// Ok 返回節點是否 存在
func (s Uint16) Ok(key uint16) bool {
	_, ok := s[key]
	return ok
}

// Len 返回節點 數量
func (s Uint16) Len() int {
	return len(s)
}
func (s Uint16) String() string {
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
