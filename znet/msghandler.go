package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("workerID = ", workerID, " is started")
	for request := range taskQueue {
		m.DoMsgHandler(request)
	}
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("add connID = ", request.GetConnection().GetConnID(), " request msgID = ", request.GetMsgId(), " to workerID = ", workerId)
	m.TaskQueue[workerId] <- request
}
