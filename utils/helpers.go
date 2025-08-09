package utils

import (
	"encoding/base64"
	"os"
)

func DecodeBase64(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func GetGithubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}
