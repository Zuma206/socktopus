package cli

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const SEPERATOR = "h"

// timestamp;secretName;connectionId;signature;

func GenerateToken(secretName string, secret string, connectionId string) (string, error) {
	timestamp := fmt.Sprint(time.Now().UnixMilli())
	signature := hmac.New(sha256.New, []byte(secret))
	_, err := signature.Write([]byte(timestamp + connectionId))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString([]byte(timestamp)) + SEPERATOR +
		hex.EncodeToString([]byte(secretName)) + SEPERATOR +
		hex.EncodeToString([]byte(connectionId)) + SEPERATOR +
		hex.EncodeToString(signature.Sum(nil)), nil
}
