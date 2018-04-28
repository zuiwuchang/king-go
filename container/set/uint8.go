package set

import (
	"bytes"
	"fmt"
)

// Uint8 .
type Uint8 map[uint8]bool

// NewUint8 創建 一個 Uint8 的 set
func NewUint8() Uint8 {
	return make(Uint8)
}

// Insert 插入節點
func (s Uint8) Insert(key uint8) {
	s[key] = true
}

// Remove 移除節點
func (s Uint8) Remove(key uint8) {
	delete(s, key)
}

// Ok 返回節點是否 存在
func (s Uint8) Ok(key uint8) bool {
	_, ok := s[key]
	return ok
}

// Len 返回節點 數量
func (s Uint8) Len() int {
	return len(s)
}
func (s Uint8) String() string {
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
