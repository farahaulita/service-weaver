package main

import (
	"context"
	"fmt"
	 openai "github.com/sashabaranov/go-openai"
	"github.com/ServiceWeaver/weaver"
	
)

type ChatGPT interface {
	Complete(ctx context.Context, prompt string) (string,error)
}

type chatGPT struct{
	weaver.Implements[ChatGPT]
	weaver.WithConfig[config]
}

type config struct {
	APIKey string `toml:"api_key"`
}

func (gpt *chatGPT) Complete(ctx context.Context, prompt string) (string,error){
	// check if key ada
	if gpt.Config().APIKey == "" {
		return "", fmt.Errorf("ChatGPT api_key not provided")
	}

	client := openai.NewClient(gpt.Config().APIKey)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ChatGPT completion error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}