package models

import (
	"errors"
	"log"
)

var DefaultSocketManager = &SocketManager{}

type SocketManager struct {
	connections map[string]*Connection
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
	sm.join(connection)
}

func (sm *SocketManager) Leave(key string) {
	sm.leave(key)
}

func (sm *SocketManager) Count() int {
	return len(sm.connections)
}

func (sm *SocketManager) Send(key string, message string) error {
	connection, ok := sm.connections[key]
	if !ok {
		return errors.New("Client not connected")
	}
	return connection.Send([]byte("DATA:" + message))
}

func (sm *SocketManager) Get(key string) (*Connection, error) {
	connection, ok := sm.connections[key]
	if !ok {
		return nil, errors.New("Connection doesn't exist")
	}
	return connection, nil
}
