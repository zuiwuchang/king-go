package sdl2

//獲取到焦點
type FocusInEvent struct {
	//得到焦點元素
	Obj Object
}

//失去焦點
type FocusOutEvent struct {
	//失去焦點元素
	Obj Object
}
