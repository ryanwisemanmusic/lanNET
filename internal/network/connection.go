package network

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionPool struct {
	mu          sync.RWMutex
	connections map[int]net.Conn
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connections: make(map[int]net.Conn),
	}
}

func (cp *ConnectionPool) Remove(id int) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	delete(cp.connections, id)
}

func (cp *ConnectionPool) Get(id int) (net.Conn, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	conn, ok := cp.connections[id]
	if !ok {
		return nil, fmt.Errorf("player %d not connected", id)
	}
	return conn, nil
}

func (cp *ConnectionPool) SendToAll(data []byte, excludeID int) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	for id, conn := range cp.connections {
		if id != excludeID {
			conn.Write(data)
		}
	}
}

func (cp *ConnectionPool) Count() int {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return len(cp.connections)
}
