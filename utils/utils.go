package utils

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
)

func GetRequest(location string) ([]byte, error) {
	resp, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func LoadAndEncode(path string) (string, error) {
	imageData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	return encodedImage, nil
}
