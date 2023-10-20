package utils

import (
	"errors"
	"os"
	"strings"
)

const SUFFIX = "_SOCKTOPUS_SECRET"

func GetSecret(name string) (string, error) {
	secret := os.Getenv(strings.ToUpper(name) + SUFFIX)
	if len(secret) < 1 {
		return "", errors.New("Secret doesn't exist")
	}
	return secret, nil
}
