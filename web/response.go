package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
}

type Headers map[string]string

func (w *ResponseWriter) Headers(headers Headers) {
	header := w.Header()
	for key, value := range headers {
		header.Set(key, value)
	}
}

func (w *ResponseWriter) Send(statusCode int, body []byte) error {
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	return err
}

func (w *ResponseWriter) SendString(statusCode int, body string) error {
	return w.Send(statusCode, []byte(body))
}

func (w *ResponseWriter) SendJson(statusCode int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Headers(Headers{
		"Content-Type": "application/json",
	})
	return w.Send(statusCode, body)
}

type H map[string]interface{}

func (w *ResponseWriter) SendError(err error, statusCode int, errorMessage string) error {
	log.Println("[ERROR]", err)
	return w.SendJson(statusCode, H{
		"code":    statusCode,
		"message": errorMessage,
	})
}
