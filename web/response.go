package web

import (
	"encoding/json"
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

func (w *ResponseWriter) SendError(statusCode int, errorMessage string) error {
	return w.SendJson(statusCode, H{
		"code":    statusCode,
		"message": errorMessage,
	})
}
