package gollama

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func getRequest(location string) ([]byte, error) {
	resp, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type PromptResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type ModelInfo struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
	Details    string `json:"details"`
}

type ModelList struct {
	Models []ModelInfo `json:"models"`
}

type Model struct {
	Name     string
	Addr     string
	Prompts  []string
	Response PromptResponse
}

func (model *Model) ListModels() ([]string, error) {
	var modelsResponse ModelList
	var avaiableModels []string
	resp, err := getRequest(model.Addr + "/api/tags")
	if err != nil {
		return avaiableModels, err
	}
	if err := json.Unmarshal(resp, &modelsResponse); err != nil {
		return avaiableModels, err
	}
	for _, model := range modelsResponse.Models {
		avaiableModels = append(avaiableModels, model.Name)
	}
	return avaiableModels, nil
}

func (model *Model) Generate(prompt string) (string, error) {
	var modelResponse string

	return modelResponse, nil
}

func CreateModel(name string, addr string) (*Model, error) {
	if addr == "" {
		addr = "localhost:11434"
	}
	newModel := &Model{Name: name, Addr: addr}
	models, err := newModel.ListModels()
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, errors.New("no models found")
	}
	check := false
	for _, model := range models {
		if model == name {
			check = true
		}
	}
	if !check {
		return nil, errors.New("model not found")
	}
	return newModel, nil
}
