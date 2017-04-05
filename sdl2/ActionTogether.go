package sdl2

import (
	"time"
)

const (
	//未工作
	_ACTION_TOGETHER_NONE = 0
	//工作 已完成
	_ACTION_TOGETHER_OK = 1
	//進行 重複工作
	_ACTION_TOGETHER_MORE = 2
)

type actionTogetherNode struct {
	//原 動作結束回調
	callback ActionCallBack
	params   interface{}

	action Action

	//工作狀態
	status int
}

//將多個 action 並行執行
//ActionTogether 的 callback 會在 所有 action 完成後被 調用
type ActionTogether struct {
	ActionBase

	//並行 action 集合
	actions []*actionTogetherNode
}

//創建並行動作
func NewActionTogether(as ...Action) *ActionTogether {
	rs := &ActionTogether{}

	size := len(as)
	if size < 1 {
		return rs
	}

	actions := make([]*actionTogetherNode, size, size)
	for i := 0; i < size; i++ {
		callback, params := as[i].GetCallBack()
		as[i].SetCallBack(
			actionTogetherOk,
			&actionTogetherOKParams{
				a:     rs,
				index: i,
			},
		)
		node := &actionTogetherNode{
			action:   as[i],
			callback: callback,
			params:   params,
			status:   _ACTION_TOGETHER_NONE,
		}
		actions[i] = node
	}
	rs.actions = actions
	return rs
}

type actionTogetherOKParams struct {
	a     *ActionTogether
	index int
}

//action 完成一個 週期
func actionTogetherOk(node Object, _a Action, params interface{}) {
	ps := params.(*actionTogetherOKParams)
	at := ps.a
	an := at.actions[ps.index]
	if an.callback != nil {
		an.callback(node, an.action, an.params)
	}

	if an.action.GetLoop() {
		an.status = _ACTION_TOGETHER_MORE
	} else {
		an.status = _ACTION_TOGETHER_OK
	}
}

//執行動作
func (a *ActionTogether) DoAction(node Object, duration time.Duration) {
	ok := true
	for _, togetherNode := range a.actions {
		//fmt.Println(i, togetherNode.status)
		if togetherNode.status == _ACTION_TOGETHER_NONE {
			togetherNode.action.DoAction(node, duration)
			if ok {
				ok = false
			}
		} else if togetherNode.status == _ACTION_TOGETHER_MORE {
			togetherNode.action.DoAction(node, duration)
		}
	}

	if ok {

		//所有 action 都執行完 至少一個 週期 回調之
		if !a.loop {
			//不循環 移除 動作
			node.RemoveAction(a)
		}

		if a.callback != nil {
			a.callback(node, a, a.params)
		}

		if a.loop {
			for _, togetherNode := range a.actions {
				togetherNode.status = _ACTION_TOGETHER_NONE
			}
		}
	}
}

//釋放 動作
func (a *ActionTogether) Destory() {
	for _, node := range a.actions {
		ac := node.action
		if ac.GetAutoDestroy() {
			ac.Destroy()
		}
	}

	*a = ActionTogether{}
}

//返回一個動作副本
func (a *ActionTogether) Clone() Action {
	size := len(a.actions)
	as := make([]Action, size, size)
	for i, node := range a.actions {
		as[i] = node.action.Clone()
		as[i].SetCallBack(node.callback, node.params)
	}
	return NewActionLinear(as...)
}
