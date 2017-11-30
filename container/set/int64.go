package set

import (
	"bytes"
	"fmt"
)

var StringSetEmpty string = "[]"

type Int64 map[int64]bool

//創建 一個 Int64 的 set
func NewInt64() Int64 {
	return make(Int64)
}

//插入節點
func (this Int64) Insert(key int64) {
	this[key] = true
}

//移除節點
func (this Int64) Remove(key int64) {
	delete(this, key)
}

//返回節點是否 存在
func (this Int64) Ok(key int64) bool {
	_, ok := this[key]
	return ok
}

//返回節點 數量
func (this Int64) Len() int {
	return len(this)
}
func (this Int64) String() string {
	if len(this) == 0 {
		return StringSetEmpty
	} else if len(this) == 1 {
		for key, _ := range this {
			return fmt.Sprintf("[%v]", key)
		}
	}

	buf := bytes.NewBufferString("[")
	first := true
	for key, _ := range this {
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
