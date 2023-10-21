package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/zuma206/socktopus/models"
	"github.com/zuma206/socktopus/utils"
	"github.com/zuma206/socktopus/web"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleRecieve(w web.ResponseWriter, r web.Request) error {
	query := r.URL.Query()
	token := query.Get("token")

	connection, err := models.NewConnection(token)
	if err != nil {
		return w.SendError(400, "Malformed token")
	}

	if connection.IsExpired() {
		return w.SendError(401, "Token expired")
	}

	secret, err := utils.GetSecret(connection.SecretName)
	if err != nil {
		return w.SendError(404, "Secret not found")
	}

	if !connection.IsSigned(secret) {
		return w.SendError(401, "Invalid signature")
	}

	connection.Socket, err = upgrader.Upgrade(w.ResponseWriter, r, nil)
	if err != nil {
		return err
	}

	models.DefaultSocketManager.Join(connection)
	defer models.DefaultSocketManager.Leave(connection.Key())

	err = connection.ApplyDeadline()
	if err != nil {
		return err
	}

	for {
		messageType, message, err := connection.Socket.ReadMessage()
		msg := string(message)
		if err != nil || messageType != websocket.TextMessage || msg == models.CLOSE {
			models.DefaultSocketManager.Leave(connection.Key())
			break
		} else if msg == models.PING {
			connection.ApplyDeadline()
			connection.Send([]byte(models.PONG))
		}
	}

	return nil
}
