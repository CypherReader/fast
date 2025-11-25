package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fastinghero/internal/core/ports"
	"fmt"
	"io"
	"net/http"
)

type DeepSeekAdapter struct {
	apiKey string
	client *http.Client
}

func NewDeepSeekAdapter(apiKey string) ports.LLMProvider {
	return &DeepSeekAdapter{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

func (a *DeepSeekAdapter) GenerateResponse(ctx context.Context, prompt string, systemPrompt string) (string, error) {
	reqBody := chatRequest{
		Model: "deepseek-chat",
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("deepseek api error: %s - %s", resp.Status, string(bodyBytes))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("no choices returned from deepseek")
	}

	return chatResp.Choices[0].Message.Content, nil
}
