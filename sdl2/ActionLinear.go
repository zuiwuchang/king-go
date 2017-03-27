package sdl2

import (
	"time"
)

type actionLinearNode struct {
	//原 動作結束回調
	callback ActionCallBack
	params   interface{}

	action Action
}

//將多個 action 線性執行
//對於 ActionLinear 中的 Action GetLoop() 的返回值 會被忽略直接以 false 處理
type ActionLinear struct {
	//動作結束回調
	callback ActionCallBack
	params   interface{}

	//線性動作 集合
	actions []actionLinearNode
	//當前動作
	pos int

	loop bool
}

//進入下個 action
func actionLinearNext(node Object, _a Action, params interface{}) {
	al := params.(*ActionLinear)
	a := al.actions[al.pos]
	if a.callback != nil {
		a.callback(node, a.action, a.params)
	}
	al.pos++

	//已經完成 所有 action
	size := len(al.actions)
	if al.pos >= size {
		if al.loop {
			al.pos = 0
		} else {
			//不循環 移除 動作
			node.RemoveAction(al)
		}
		if al.callback != nil {
			al.callback(node, al, al.params)
		}
	}
}

//執行動作
func (a *ActionLinear) DoAction(node Object, duration time.Duration) {
	//已經完成 所有 action size == 0
	size := len(a.actions)
	if a.pos >= size {
		if !a.loop {
			//不循環 移除 動作
			node.RemoveAction(a)
		}

		if a.callback != nil {
			a.callback(node, a, a.params)
		}
		return
	}

	a.actions[a.pos].action.DoAction(node, duration)
}

//創建線性動作
func NewActionLinear(as ...Action) *ActionLinear {
	rs := &ActionLinear{}

	size := len(as)
	if size < 1 {
		return rs
	}

	actions := make([]actionLinearNode, size, size)
	for i := 0; i < size; i++ {
		callback, params := as[i].GetCallBack()
		as[i].SetCallBack(actionLinearNext, rs)
		node := actionLinearNode{
			action:   as[i],
			callback: callback,
			params:   params,
		}
		actions[i] = node
	}
	rs.actions = actions
	return rs
}

//設置 action 完成 通知
func (a *ActionLinear) SetCallBack(callback ActionCallBack, params interface{}) Action {
	a.callback = callback
	a.params = params
	return a
}

//返回 action 完成 通知
func (a *ActionLinear) GetCallBack() (ActionCallBack, interface{}) {
	return a.callback, a.params
}

//釋放 動作
func (a *ActionLinear) Destory() {
	*a = ActionLinear{}
}

//是否自動 釋放
func (a *ActionLinear) Auto() bool {
	return false
}

//返回一個動作副本
func (a *ActionLinear) Clone() Action {
	action := *a
	action.pos = 0
	return &action
}

//返回 是否 循環執行
func (a *ActionLinear) GetLoop() bool {
	return a.loop
}

//設置 是否 循環執行
func (a *ActionLinear) SetLoop(yes bool) Action {
	a.loop = yes
	return a
}
