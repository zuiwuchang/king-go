package sdl2

const (
	//鼠標 單擊
	UI_EVT_LBTN_CLICK = "k-lbtn_click"
	UI_EVT_RBTN_CLICK = "k-rbtn_click"
	UI_EVT_MBTN_CLICK = "k-mbtn_click"
	//鼠標 雙擊
	UI_EVT_LBTN_DCLICK = "k-lbtn_dclick"
	UI_EVT_RBTN_DCLICK = "k-rbtn_dclick"
	UI_EVT_MBTN_DCLICK = "k-mbtn_dclick"

	//鍵盤 按鍵
	UI_EVT_KEY_DOWM = "k-key-down"
	UI_EVT_KEY_UP   = "k-key-up"

	//焦點
	UI_EVT_FOCUS_IN  = "k-foucs-in"
	UI_EVT_FOCUS_OUT = "k-foucs-out"
)

const (
	R_MASK = 0x000000ff
	G_MASK = 0x0000ff00
	B_MASK = 0x00ff0000
	A_MASK = 0xff000000
)

//ui事件 回調
type UiCallBack func(ui UiObject, info interface{})

type UiObject interface {
	Object
	//設置 事件回調
	SetEventCallback(evt string /*ui事件名*/, callback UiCallBack)
	//返回 事件回調 函數
	GetEventCallback(evt string) UiCallBack
}
type UiBase struct {
	Node
	evtCallbacks map[string]UiCallBack
}

func posInRect(posX, posY, x, y, w, h int) bool {
	return posX >= x && posX <= x+w &&
		posY >= y && posY <= y+h
}
func (u *UiBase) posInRect(posX, posY int32) bool {
	x, y := u.GetDrawPos()
	w, h := u.GetDrawSize()
	return posInRect(int(posX), int(posY), int(x), int(y), int(w), int(h))
}

//設置 事件回調
func (u *UiBase) SetEventCallback(evt string /*ui事件名*/, callback UiCallBack) {
	if u.evtCallbacks == nil {
		u.evtCallbacks = make(map[string]UiCallBack)
	}
	u.evtCallbacks[evt] = callback
}

//返回 事件回調 函數
func (u *UiBase) GetEventCallback(evt string) UiCallBack {
	if u.evtCallbacks == nil {
		return nil
	}
	callback, _ := u.evtCallbacks[evt]
	return callback
}
