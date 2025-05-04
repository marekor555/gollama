package gollama

import (
	"bytes"
	"encoding/json"
	"github.com/marekor555/gollama/structs"
	"github.com/marekor555/gollama/utils"
	"io"
	"net/http"
)

type Chat struct {
	Model    Model
	Messages []structs.Message
}

func (chat *Chat) Send(msg string, images ...string) error {
	message := structs.Message{Role: "user", Content: msg}
	if images != nil {
		for _, imagePath := range images {
			encoded, err := utils.LoadAndEncode(imagePath)
			if err != nil {
				return err
			}
			message.Images = append(message.Images, encoded)
		}
	}
	chat.Messages = append(chat.Messages, message)
	return nil
}

func (chat *Chat) Receive() (string, error) {
	var modelResponse string
	var promptResponse structs.ChatResponse
	prompt := structs.ChatPrompt{Model: chat.Model.Name, Messages: chat.Messages, System: chat.Model.System}
	promptStr, err := json.Marshal(prompt)
	if err != nil {
		return modelResponse, err
	}
	resp, err := http.Post(chat.Model.Addr+"/api/chat", "application/json", bytes.NewBuffer(promptStr))
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

	modelResponse = promptResponse.Response.Content
	//log.Println(modelResponse)
	chat.Messages = append(chat.Messages, structs.Message{Role: "assistant", Content: modelResponse})
	return modelResponse, nil
}

func (chat *Chat) SendAndReceive(msg string, images ...string) (string, error) {
	err := chat.Send(msg, images...)
	if err != nil {
		return "", err
	}
	return chat.Receive()
}

func NewChat(modelName string, addr string, system string) (*Chat, error) {
	model, err := NewModel(modelName, addr, system)
	if err != nil {
		return nil, err
	}
	return &Chat{*model, []structs.Message{}}, nil
}
