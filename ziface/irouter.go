package ziface

type IRouter interface {
	PreHandle(request IRequest)  // 前置处理
	Handle(request IRequest)     // 处理
	PostHandle(request IRequest) // 后置处理
}
