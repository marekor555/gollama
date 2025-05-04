package gollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/marekor555/gollama/structs"
	"github.com/marekor555/gollama/utils"
	"io"
	"net/http"
)

type Model struct {
	Name   string
	Addr   string
	System string
}

func (model *Model) ListModels() ([]string, error) {
	var modelsResponse structs.ModelList
	var availableModels []string
	resp, err := utils.GetRequest(model.Addr + "/api/tags")
	if err != nil {
		return availableModels, err
	}
	//println(string(resp))
	if err := json.Unmarshal(resp, &modelsResponse); err != nil {
		return availableModels, err
	}
	for _, model := range modelsResponse.Models {
		availableModels = append(availableModels, model.Name)
	}
	return availableModels, nil
}

func (model *Model) Generate(msg string, images ...string) (string, error) {
	var modelResponse string
	var promptResponse structs.PromptResponse
	prompt := structs.Prompt{Model: model.Name, Prompt: msg, System: model.System}
	if images != nil {
		for _, imagePath := range images {
			encoded, err := utils.LoadAndEncode(imagePath)
			if err != nil {
				return "", err
			}
			prompt.Images = append(prompt.Images, encoded)
		}
	}
	promptStr, err := json.Marshal(prompt)
	if err != nil {
		return modelResponse, err
	}
	resp, err := http.Post(model.Addr+"/api/generate", "application/json", bytes.NewBuffer(promptStr))
	if err != nil {
		return modelResponse, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

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
	if !utils.Find(models, modelName) {
		return errors.New("model not found")
	}
	model.Name = modelName
	return nil
}

func NewModel(name string, addr string, system string) (*Model, error) {
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
	if !utils.Find(models, name) {
		return nil, errors.New("model not found")
	}
	return newModel, nil
}
