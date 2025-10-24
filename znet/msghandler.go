package znet

import (
	"fmt"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("api msgId = %d is not FOUND\n", request.GetMsgId())
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断msgId是否绑定了方法
	// 方法不能重复绑定
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api, msgId = " + fmt.Sprintf("%d", msgId))
	}
	m.Apis[msgId] = router
	fmt.Println("add api msgId = ", msgId)
}
