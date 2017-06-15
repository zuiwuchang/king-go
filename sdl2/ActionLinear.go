package sdl2

import (
	"time"
)

type actionLinearNode struct {
	//原 動作結束回調
	callback ActionCallBack
	params   interface{}

	action IAction
}

//將多個 action 線性執行
//對於 ActionLinear 中的 Action GetLoop() 的返回值 會被忽略直接以 false 處理
type ActionLinear struct {
	ActionBase

	//線性動作 集合
	actions []actionLinearNode
	//當前動作
	pos int
}

//進入下個 action
func actionLinearNext(node IObject, _a IAction, params interface{}) {
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
func (a *ActionLinear) DoAction(node IObject, duration time.Duration) {
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
func NewActionLinear(as ...IAction) *ActionLinear {
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

//釋放 動作
func (a *ActionLinear) Destory() {
	for _, node := range a.actions {
		ac := node.action
		if ac.GetAutoDestroy() {
			ac.Destroy()
		}
	}

	*a = ActionLinear{}
}

//返回一個動作副本
func (a *ActionLinear) Clone() IAction {
	size := len(a.actions)
	as := make([]IAction, size, size)
	for i, node := range a.actions {
		as[i] = node.action.Clone()
		as[i].SetCallBack(node.callback, node.params)
	}
	return NewActionLinear(as...)
}
