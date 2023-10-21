package controllers

import (
	"encoding/json"

	"github.com/zuma206/socktopus/models"
	"github.com/zuma206/socktopus/utils"
	"github.com/zuma206/socktopus/web"
)

type Message struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

type SendRequest struct {
	SecretName string    `json:"secretName"`
	Secret     string    `json:"secret"`
	Messages   []Message `json:"messages"`
}

type Response struct {
	Succeeded []Message `json:"succeeded"`
	Failed    []Message `json:"failed"`
}

func HandleSend(w web.ResponseWriter, r web.Request) error {
	decoder := json.NewDecoder(r.Body)
	request := new(SendRequest)
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

	response := Response{
		Failed:    make([]Message, 0),
		Succeeded: make([]Message, 0),
	}

	for _, message := range request.Messages {
		key := models.ConnectionKey(request.SecretName, message.Recipient)
		err := models.DefaultSocketManager.Send(key, message.Content)
		if err == nil {
			response.Succeeded = append(response.Succeeded, message)
		} else {
			response.Failed = append(response.Failed, message)
		}
	}

	return w.SendJson(200, response)
}
