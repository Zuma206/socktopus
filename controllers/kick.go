package controllers

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/zuma206/socktopus/models"
	"github.com/zuma206/socktopus/utils"
	"github.com/zuma206/socktopus/web"
)

type KickRequest struct {
	SecretName  string   `json:"secretName"`
	Secret      string   `json:"secret"`
	Connections []string `json:"connections"`
}

func HandleKick(w web.ResponseWriter, r web.Request) error {
	decoder := json.NewDecoder(r.Body)
	request := new(KickRequest)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(request); err != nil {
		return w.SendError(400, "Malformed request")
	}

	secret, err := utils.GetSecret(request.SecretName)
	if err != nil {
		return w.SendError(404, "Secret not found")
	}

	if secret != request.Secret {
		return w.SendError(401, "Invalid secret")
	}

	wg := new(sync.WaitGroup)
	for _, connectionId := range request.Connections {
		key := models.ConnectionKey(request.SecretName, connectionId)
		connection, err := models.DefaultSocketManager.Get(key)
		if err != nil {
			continue
		}
		expiresAt := int64(connection.Timestamp + utils.TOKEN_LIFESPAN)
		diff := (expiresAt - time.Now().UnixMilli()) * int64(time.Millisecond)

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-time.After(time.Duration(diff))
			models.DefaultSocketManager.Leave(key)
		}()
	}
	wg.Wait()

	return w.SendJson(200, 200)
}
