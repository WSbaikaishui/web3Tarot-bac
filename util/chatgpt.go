package util

//
//import (
//	"context"
//	"fmt"
//	"strconv"
//	"strings"
//	"sync"
//
//	"github.com/neo-ngd/nchat-backend/config"
//	tokenizer "github.com/samber/go-gpt-3-encoder"
//	gpt "github.com/sashabaranov/go-openai"
//)
//
//const totalMaxToken = 4096
//
//const sysPrompt = `You are ChatGPT, a large language model trained by OpenAI. Respond conversationally. Current datetime: %s`
//const egPrompt = "Hello"
//const egResponse = "Hello! How can I help you today?"
//
//var (
//	RequestMap = sync.Map{}
//)
//
//func CreateChatGPTResponse(userId uint, prompt string, hs []*History) (string, error) {
//	// set to true to prevent user from starting next request
//	RequestMap.Store(userId, true)
//	// set to false so user can start next request
//	defer RequestMap.Store(userId, false)
//
//	c := gpt.NewClient(config.AppCfg.ChatGPT.ApiKey)
//	ctx := context.Background()
//
//	// moderation check
//	m := gpt.ModerationRequest{
//		Input: prompt,
//	}
//	r, err := c.Moderations(ctx, m)
//	if err != nil {
//		return "", err
//	}
//	if len(r.Results) > 0 && r.Results[0].Flagged {
//		return "", fmt.Errorf("inappropriate content")
//	}
//
//	// decorate prompt
//	userName := "User_" + strconv.Itoa(int(userId)) // address consumes too many tokens
//	messages, err := decoratePrompt(userName, prompt, hs)
//	if err != nil {
//		return "", err
//	}
//
//	// make completion request
//	req := gpt.ChatCompletionRequest{
//		Model:       config.AppCfg.ChatGPT.Model,
//		Messages:    messages,
//		MaxTokens:   config.AppCfg.ChatGPT.MaxReplyTokens,
//		Temperature: config.AppCfg.ChatGPT.Temperature,
//		User:        userName,
//	}
//	resp, err := c.CreateChatCompletion(ctx, req)
//	if err != nil {
//		return "", err
//	}
//	reply := strings.TrimSpace(resp.Choices[0].Message.Content)
//	return reply, nil
//}
//
//type History struct {
//	Q string
//	A string
//}
//
//func decoratePrompt(userName string, prompt string, hs []*History) ([]gpt.ChatCompletionMessage, error) {
//	for i := 0; i <= len(hs); i++ {
//		msgArray := new([]gpt.ChatCompletionMessage)
//		// system message
//		addMessageToArray(msgArray, gpt.ChatCompletionMessage{Role: gpt.ChatMessageRoleSystem, Content: fmt.Sprintf(sysPrompt, FormatUtcNow())})
//		// example user message
//		addMessageToArray(msgArray, gpt.ChatCompletionMessage{Role: gpt.ChatMessageRoleSystem, Name: "example_user", Content: egPrompt})
//		// example assistant message
//		addMessageToArray(msgArray, gpt.ChatCompletionMessage{Role: gpt.ChatMessageRoleSystem, Name: "example_assistant", Content: egResponse})
//		// history messages
//		concatHistory(userName, msgArray, hs[i:])
//		// new user message
//		addMessageToArray(msgArray, gpt.ChatCompletionMessage{Role: gpt.ChatMessageRoleUser, Name: userName, Content: prompt})
//
//		count, err := getTokenCount(msgArray)
//		if err != nil {
//			return nil, fmt.Errorf("encode messages error: %v", err)
//		}
//
//		if count+config.AppCfg.ChatGPT.MaxReplyTokens < totalMaxToken {
//			return *msgArray, nil
//		}
//	}
//	return nil, fmt.Errorf("prompt too long")
//}
//
//func concatHistory(userName string, array *[]gpt.ChatCompletionMessage, hs []*History) {
//	for _, h := range hs {
//		addMessageToArray(array, gpt.ChatCompletionMessage{Role: gpt.ChatMessageRoleUser, Name: userName, Content: h.Q})
//		addMessageToArray(array, gpt.ChatCompletionMessage{Role: gpt.ChatMessageRoleAssistant, Content: h.A})
//	}
//}
//
//func addMessageToArray(array *[]gpt.ChatCompletionMessage, msg gpt.ChatCompletionMessage) {
//	*array = append(*array, msg)
//}
//
//func getTokenCount(array *[]gpt.ChatCompletionMessage) (int, error) {
//	encoder, err := tokenizer.NewEncoder()
//	if err != nil {
//		return 0, err
//	}
//	count := 0
//	for _, m := range *array {
//		if len(m.Name) == 0 {
//			count += 1 // role
//		} else {
//			t, err := encoder.Encode(m.Name) // name
//			if err != nil {
//				return 0, err
//			}
//			count += len(t)
//		}
//		t, err := encoder.Encode(m.Content) // content
//		if err != nil {
//			return 0, err
//		}
//		count += len(t)
//	}
//	count += 4 // metadata
//	return count, nil
//}
