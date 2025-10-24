package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 非阻塞的方式处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
}

