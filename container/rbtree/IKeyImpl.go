package rbtree

//爲 int 實現的 IKey 接口
type IKeyInt int

func (i IKeyInt) Compare(k IKey) int {
	v := k.(IKeyInt)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 int8 實現的 IKey 接口
type IKeyInt8 int8

func (i IKeyInt8) Compare(k IKey) int {
	v := k.(IKeyInt8)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 int16 實現的 IKey 接口
type IKeyInt16 int16

func (i IKeyInt16) Compare(k IKey) int {
	v := k.(IKeyInt16)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 int32 實現的 IKey 接口
type IKeyInt32 int16

func (i IKeyInt32) Compare(k IKey) int {
	v := k.(IKeyInt32)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 int64 實現的 IKey 接口
type IKeyInt64 int16

func (i IKeyInt64) Compare(k IKey) int {
	v := k.(IKeyInt64)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 uint 實現的 IKey 接口
type IKeyUInt uint

func (i IKeyUInt) Compare(k IKey) int {
	v := k.(IKeyUInt)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 uint8 實現的 IKey 接口
type IKeyUInt8 uint8

func (i IKeyUInt8) Compare(k IKey) int {
	v := k.(IKeyUInt8)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 uint16 實現的 IKey 接口
type IKeyUInt16 uint16

func (i IKeyUInt16) Compare(k IKey) int {
	v := k.(IKeyUInt16)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 uint32 實現的 IKey 接口
type IKeyUInt32 uint32

func (i IKeyUInt32) Compare(k IKey) int {
	v := k.(IKeyUInt32)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 uint64 實現的 IKey 接口
type IKeyUInt64 uint64

func (i IKeyUInt64) Compare(k IKey) int {
	v := k.(IKeyUInt64)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 float32 實現的 IKey 接口
type IKeyFloat32 float32

func (i IKeyFloat32) Compare(k IKey) int {
	v := k.(IKeyFloat32)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 float64 實現的 IKey 接口
type IKeyFloat64 float64

func (i IKeyFloat64) Compare(k IKey) int {
	v := k.(IKeyFloat64)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 byte 實現的 IKey 接口
type IKeyByte byte

func (i IKeyByte) Compare(k IKey) int {
	v := k.(IKeyByte)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}

//爲 string 實現的 IKey 接口
type IKeyString string

func (i IKeyString) Compare(k IKey) int {
	v := k.(IKeyString)
	if v == i {
		return 0
	}
	if i < v {
		return -1
	}
	return 1
}
