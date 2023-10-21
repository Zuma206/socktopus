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

	"github.com/gorilla/websocket"
	"github.com/zuma206/socktopus/cli"
)

const CLOSE = "CLOSE"
const PING = "PING"
const PONG = "PONG"

type Connection struct {
	ConnectionId string
	SecretName   string
	Timestamp    int
	Signature    []byte
	Socket       *websocket.Conn
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
		Socket:       nil,
	}, nil
}

const TOKEN_LIFESPAN = 60_000

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

func ConnectionKey(secretName string, connectionId string) string {
	return secretName + ";;" + connectionId
}

func (c *Connection) Key() string {
	return ConnectionKey(c.SecretName, c.ConnectionId)
}

func (c *Connection) Close() {
	if c.Socket != nil {
		c.Send([]byte(CLOSE))
		c.Socket.Close()
	}
}

const READ_DEADLINE = time.Duration(10 * time.Second)

func (c *Connection) ApplyDeadline() error {
	return c.Socket.SetReadDeadline(time.Now().Add(READ_DEADLINE))
}

func (c *Connection) Send(message []byte) error {
	if err := c.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}
