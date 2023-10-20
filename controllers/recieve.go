package controllers

import (
	"github.com/zuma206/socktopus/models"
	"github.com/zuma206/socktopus/utils"
	"github.com/zuma206/socktopus/web"
)

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

	return w.SendJson(200, web.H{
		"secretName":   connection.SecretName,
		"connectionId": connection.ConnectionId,
	})
}
