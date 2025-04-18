package gollama

import (
	"bytes"
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

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

type Prompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system"`
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
}

type ModelList struct {
	Models []ModelInfo `json:"models"`
}

type Model struct {
	Name   string
	Addr   string
	System string
}

func (model *Model) ListModels() ([]string, error) {
	var modelsResponse ModelList
	var avaiableModels []string
	resp, err := getRequest(model.Addr + "/api/tags")
	if err != nil {
		return avaiableModels, err
	}
	//println(string(resp))
	if err := json.Unmarshal(resp, &modelsResponse); err != nil {
		return avaiableModels, err
	}
	for _, model := range modelsResponse.Models {
		avaiableModels = append(avaiableModels, model.Name)
	}
	return avaiableModels, nil
}

func (model *Model) Generate(msg string) (string, error) {
	var modelResponse string
	var promptResponse PromptResponse
	prompt := Prompt{model.Name, msg, false, model.System}
	promptStr, err := json.Marshal(prompt)
	if err != nil {
		return modelResponse, err
	}
	resp, err := http.Post(model.Addr+"/api/generate", "application/json", bytes.NewBuffer(promptStr))
	if err != nil {
		return modelResponse, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&promptResponse); err != nil {
		return modelResponse, err
	}

	modelResponse = promptResponse.Response
	return modelResponse, nil
}

func (model *Model) changeModel(modelName string) error {
	models, err := model.ListModels()
	if err != nil {
		return err
	}
	if !find(models, modelName) {
		return errors.New("model not found")
	}
	model.Name = modelName
	return nil
}

func CreateModel(name string, addr string, system string) (*Model, error) {
	if addr == "" {
		addr = "http://localhost:11434"
	}
	newModel := &Model{name, addr, system}
	models, err := newModel.ListModels()
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, errors.New("no models found")
	}
	if !find(models, name) {
		return nil, errors.New("model not found")
	}
	return newModel, nil
}
