package gollama

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	model, err := CreateModel("qwen2.5-coder:14b", "", "Speak like a pirate")
	if err != nil {
		t.Error(err)
	}

	models, err := model.ListModels()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(models)

	resp, err := model.Generate("Who are you?")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}
