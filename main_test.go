package gollama

import (
	"fmt"
	"github.com/marekor555/gollama/manage"
	"testing"
)

func Test(t *testing.T) {
	modelName := "llava:latest"
	systemPrompt := "Dont add anything to the image"
	prompt := "Describe the image"
	manager := manage.NewManager("")
	println(manager.ListModels())
	model, err := NewModel(modelName, "", systemPrompt)
	if err != nil {
		t.Error(err)
	}

	models, err := model.ListModels()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(models)

	resp, err := model.Generate(prompt, "random.jpg")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)

}

func Test_chat(t *testing.T) {
	modelName := "llava:latest"
	systemPrompt := "Speak shortly and clearly."
	prompt1 := "Who are you?"
	prompt2 := "Describe the image"
	chat, err := NewChat(modelName, "", systemPrompt)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("USR:", prompt1)
	err = chat.Send(prompt1)
	if err != nil {
		t.Error(err)
	}

	resp, err := chat.Receive()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)

	fmt.Println("USR:", prompt2)
	resp, err = chat.SendAndReceive(prompt2, "random.jpg")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}

//func Test_manager(t *testing.T) {
//	manager := manage.NewManager("")
//	models := manager.ListModels()
//	fmt.Println(models)
//
//	err := manager.Install("nomic-embed-text")
//	if err != nil {
//		t.Error(err)
//	}
//	models = manager.ListModels()
//	fmt.Println(models)
//
//	err = manager.Remove("nomic-embed-text")
//	if err != nil {
//		t.Error(err)
//	}
//	models = manager.ListModels()
//}
