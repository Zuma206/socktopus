package models

import (
	"errors"
	"log"
	"sync"
)

var DefaultSocketManager = &SocketManager{connections: make(map[string]*Connection)}

type SocketManager struct {
	connections map[string]*Connection
	lock        sync.RWMutex
}

func (sm *SocketManager) join(connection *Connection) {
	key := connection.Key()
	sm.leave(key)
	sm.connections[key] = connection
	log.Printf("[JOINED] %s", key)
}

func (sm *SocketManager) leave(key string) {
	connection, ok := sm.connections[key]
	if !ok {
		return
	}
	connection.Close()
	delete(sm.connections, key)
	log.Printf("[LEFT] %s", key)
}

func (sm *SocketManager) Join(connection *Connection) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.join(connection)
}

func (sm *SocketManager) Leave(key string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.leave(key)
}

func (sm *SocketManager) Count() int {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return len(sm.connections)
}

func (sm *SocketManager) Send(key string, message string) error {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	connection, ok := sm.connections[key]
	if !ok {
		return errors.New("Client not connected")
	}
	return connection.Send([]byte("DATA:" + message))
}

func (sm *SocketManager) Get(key string) (*Connection, error) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	connection, ok := sm.connections[key]
	if !ok {
		return nil, errors.New("Connection doesn't exist")
	}
	return connection, nil
}
