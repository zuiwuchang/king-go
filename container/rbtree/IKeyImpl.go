package rbtree

//爲 int 實現的 IKey 接口
type IKeyInt int

func (i IKeyInt) Less(k IKey) bool {
	v := k.(IKeyInt)
	return i < v
}
func (i IKeyInt) Equal(k IKey) bool {
	v := k.(IKeyInt)
	return i == v
}

//爲 int8 實現的 IKey 接口
type IKeyInt8 int8

func (i IKeyInt8) Less(k IKey) bool {
	v := k.(IKeyInt8)
	return i < v
}
func (i IKeyInt8) Equal(k IKey) bool {
	v := k.(IKeyInt8)
	return i == v
}

//爲 int16 實現的 IKey 接口
type IKeyInt16 int16

func (i IKeyInt16) Less(k IKey) bool {
	v := k.(IKeyInt16)
	return i < v
}
func (i IKeyInt16) Equal(k IKey) bool {
	v := k.(IKeyInt16)
	return i == v
}

//爲 int32 實現的 IKey 接口
type IKeyInt32 int16

func (i IKeyInt32) Less(k IKey) bool {
	v := k.(IKeyInt32)
	return i < v
}
func (i IKeyInt32) Equal(k IKey) bool {
	v := k.(IKeyInt32)
	return i == v
}

//爲 int64 實現的 IKey 接口
type IKeyInt64 int16

func (i IKeyInt64) Less(k IKey) bool {
	v := k.(IKeyInt64)
	return i < v
}
func (i IKeyInt64) Equal(k IKey) bool {
	v := k.(IKeyInt64)
	return i == v
}

//爲 uint 實現的 IKey 接口
type IKeyUInt uint

func (i IKeyUInt) Less(k IKey) bool {
	v := k.(IKeyUInt)
	return i < v
}
func (i IKeyUInt) Equal(k IKey) bool {
	v := k.(IKeyUInt)
	return i == v
}

//爲 uint8 實現的 IKey 接口
type IKeyUInt8 uint8

func (i IKeyUInt8) Less(k IKey) bool {
	v := k.(IKeyUInt8)
	return i < v
}
func (i IKeyUInt8) Equal(k IKey) bool {
	v := k.(IKeyUInt8)
	return i == v
}

//爲 uint16 實現的 IKey 接口
type IKeyUInt16 uint16

func (i IKeyUInt16) Less(k IKey) bool {
	v := k.(IKeyUInt16)
	return i < v
}
func (i IKeyUInt16) Equal(k IKey) bool {
	v := k.(IKeyUInt16)
	return i == v
}

//爲 uint32 實現的 IKey 接口
type IKeyUInt32 uint32

func (i IKeyUInt32) Less(k IKey) bool {
	v := k.(IKeyUInt32)
	return i < v
}
func (i IKeyUInt32) Equal(k IKey) bool {
	v := k.(IKeyUInt32)
	return i == v
}

//爲 uint64 實現的 IKey 接口
type IKeyUInt64 uint64

func (i IKeyUInt64) Less(k IKey) bool {
	v := k.(IKeyUInt64)
	return i < v
}
func (i IKeyUInt64) Equal(k IKey) bool {
	v := k.(IKeyUInt64)
	return i == v
}

//爲 float32 實現的 IKey 接口
type IKeyFloat32 float32

func (i IKeyFloat32) Less(k IKey) bool {
	v := k.(IKeyFloat32)
	return i < v
}
func (i IKeyFloat32) Equal(k IKey) bool {
	v := k.(IKeyFloat32)
	return i == v
}

//爲 float64 實現的 IKey 接口
type IKeyFloat64 float64

func (i IKeyFloat64) Less(k IKey) bool {
	v := k.(IKeyFloat64)
	return i < v
}
func (i IKeyFloat64) Equal(k IKey) bool {
	v := k.(IKeyFloat64)
	return i == v
}

//爲 byte 實現的 IKey 接口
type IKeyByte byte

func (i IKeyByte) Less(k IKey) bool {
	v := k.(IKeyByte)
	return i < v
}
func (i IKeyByte) Equal(k IKey) bool {
	v := k.(IKeyByte)
	return i == v
}

//爲 string 實現的 IKey 接口
type IKeyString string

func (i IKeyString) Less(k IKey) bool {
	v := k.(IKeyString)
	return i < v
}
func (i IKeyString) Equal(k IKey) bool {
	v := k.(IKeyString)
	return i == v
}
