package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zuma206/socktopus/cli"
)

type Connection struct {
	ConnectionId string
	SecretName   string
	Timestamp    int
	Signature    []byte
}

func NewConnection(token string) (*Connection, error) {
	parts := strings.Split(token, cli.SEPERATOR)
	if len(parts) != 4 {
		return nil, errors.New("Malformed token")
	}

	timestamp, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	timestampInt, err := strconv.Atoi(string(timestamp))
	if err != nil {
		return nil, err
	}

	secretName, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	connectionId, err := hex.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}

	signature, err := hex.DecodeString(parts[3])
	if err != nil {
		return nil, err
	}

	return &Connection{
		ConnectionId: string(connectionId),
		SecretName:   string(secretName),
		Timestamp:    timestampInt,
		Signature:    signature,
	}, nil
}

const TOKEN_LIFESPAN = 10_000

func (c *Connection) IsExpired() bool {
	if time.Now().UnixMilli()-int64(c.Timestamp) > TOKEN_LIFESPAN {
		return true
	}
	return false
}

func (c *Connection) IsSigned(secret string) bool {
	data := fmt.Sprint(c.Timestamp) + c.ConnectionId
	verification := hmac.New(sha256.New, []byte(secret))
	if _, err := verification.Write([]byte(data)); err != nil {
		return false
	}
	return hmac.Equal(c.Signature, verification.Sum(nil))
}
