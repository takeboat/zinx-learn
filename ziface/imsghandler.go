package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 非阻塞的方式处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
	StartWorkerPool()                        // 启动worker 工作池
	SendMsgToTaskQueue(request IRequest)    // 消息交给TaskQueue，由worker处理
}
