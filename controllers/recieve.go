package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/zuma206/socktopus/cli"
	"github.com/zuma206/socktopus/utils"
	"github.com/zuma206/socktopus/web"
)

func MalformedToken(w web.ResponseWriter) error {
	return w.SendError(400, "Malformed token")
}

func ExpiredToken(w web.ResponseWriter) error {
	return w.SendError(401, "Expired token")
}

func SecretNotFound(w web.ResponseWriter) error {
	return w.SendError(404, "Secret not found")
}

func InvalidSignature(w web.ResponseWriter) error {
	return w.SendError(401, "Invalid signature")
}

const TOKEN_LIFESPAN = 10_000

func HandleRecieve(w web.ResponseWriter, r web.Request) error {
	query := r.URL.Query()
	token := query.Get("token")
	parts := strings.Split(token, cli.SEPERATOR)

	if len(parts) != 4 {
		return MalformedToken(w)
	}

	timestamp, err := hex.DecodeString(parts[0])
	if err != nil {
		return MalformedToken(w)
	}
	timestampStr := string(timestamp)

	timestampInt, err := strconv.Atoi(timestampStr)
	if err != nil {
		return MalformedToken(w)
	}

	if time.Now().UnixMilli()-int64(timestampInt) > TOKEN_LIFESPAN {
		return ExpiredToken(w)
	}

	secretName, err := hex.DecodeString(parts[1])
	if err != nil {
		return MalformedToken(w)
	}

	secret, err := utils.GetSecret(string(secretName))
	if err != nil {
		return SecretNotFound(w)
	}

	connectionId, err := hex.DecodeString(parts[2])
	if err != nil {
		return MalformedToken(w)
	}
	connectionIdStr := string(connectionId)

	signature, err := hex.DecodeString(parts[3])
	if err != nil {
		return err
	}

	verification := hmac.New(sha256.New, []byte(secret))
	if _, err := verification.Write([]byte(timestampStr + connectionIdStr)); err != nil {
		return err
	}

	if ok := hmac.Equal(signature, verification.Sum(nil)); !ok {
		return InvalidSignature(w)
	}

	return w.SendJson(200, web.H{
		"secretName":   string(secretName),
		"connectionId": connectionIdStr,
	})
}
