package manage

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/marekor555/gollama/structs"
	"github.com/marekor555/gollama/utils"
	"io"
	"net/http"
)

type Manager struct {
	Addr string
}

func (manager *Manager) ListModels() []string {
	var modelsResponse structs.ModelList
	var availableModels []string
	resp, err := utils.GetRequest(manager.Addr + "/api/tags")
	if err != nil {
		return availableModels
	}
	if err := json.Unmarshal(resp, &modelsResponse); err != nil {
		return availableModels
	}
	for _, model := range modelsResponse.Models {
		availableModels = append(availableModels, model.Name)
	}
	return availableModels
}

func (manager *Manager) Install(modelName string) error {
	pullRequest := struct {
		Name string `json:"name"`
	}{
		Name: modelName,
	}

	requestData, err := json.Marshal(pullRequest)
	if err != nil {
		return err
	}

	resp, err := http.Post(manager.Addr+"/api/pull", "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to pull model")
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var response struct {
			Status    string  `json:"status"`
			Completed float64 `json:"completed"`
			Total     float64 `json:"total"`
			Error     string  `json:"error,omitempty"`
		}

		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("failed to decode response")
		}

		if response.Error != "" {
			return errors.New(response.Error)
		}
	}

	return nil
}

func (manager *Manager) Remove(modelName string) error {
	deleteRequest := struct {
		Name string `json:"name"`
	}{
		Name: modelName,
	}

	requestData, err := json.Marshal(deleteRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, manager.Addr+"/api/delete", bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to delete model")
	}

	return nil
}

func NewManager(addr string) *Manager {
	if addr == "" {
		addr = "http://localhost:11434"
	}
	return &Manager{addr}
}
