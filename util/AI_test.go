package util

import (
	"fmt"
	"testing"
)

func TestCreateChatGPTResponse(t *testing.T) {
	question := "hello"
	card := "the fool"
	res, err := CreateChatGPTResponse(card, question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
