package models

var DefaultSocketManager = NewSocketManager()

type SocketManager struct {
	connections map[string]*Connection
	joining     chan *Connection
	leaving     chan string
}

func NewSocketManager() *SocketManager {
	sm := &SocketManager{
		connections: make(map[string]*Connection),
		joining:     make(chan *Connection),
		leaving:     make(chan string),
	}

	go func() {
		for {
			select {
			case connection := <-sm.joining:
				sm.join(connection)
			case key := <-sm.leaving:
				sm.leave(key)
			}
		}
	}()

	return sm
}

func (sm *SocketManager) join(connection *Connection) {
	key := connection.Key()
	sm.leave(key)
	sm.connections[key] = connection
}

func (sm *SocketManager) leave(key string) {
	connection, ok := sm.connections[key]
	if !ok {
		return
	}
	connection.Close()
	delete(sm.connections, key)
}

func (sm *SocketManager) Join(connection *Connection) {
	sm.joining <- connection
}

func (sm *SocketManager) Leave(key string) {
	sm.leaving <- key
}

func (sm *SocketManager) Count() int {
	return len(sm.connections)
}
