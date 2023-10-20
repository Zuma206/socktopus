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

	conn, err := upgrader.Upgrade(w.ResponseWriter, r, nil)
	if err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, World!")); err != nil {
		return err
	}
	conn.Close()

	return nil
}
