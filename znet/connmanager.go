package znet

import (
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[conn.GetConnID()] = conn
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	panic("TODO: Implement")
}

func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	panic("TODO: Implement")
}

func (c *ConnManager) Len() int {
	panic("TODO: Implement")
}

func (c *ConnManager) ClearConn() {
	panic("TODO: Implement")
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
