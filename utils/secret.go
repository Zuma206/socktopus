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

func GetSecretNames() []string {
	services := make([]string, 0)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if strings.HasSuffix(parts[0], SUFFIX) {
			name := strings.SplitN(parts[0], SUFFIX, 2)
			services = append(services, strings.ToLower(name[0]))
		}
	}
	return services
}
