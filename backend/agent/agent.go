package agent

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"

	"github.com/nadiramlha/relay/backend/models"
)

// Agent handles AI processing using OpenAI
type Agent struct {
	client    *openai.Client
	model     string
	systemMsg string
}

// New creates a new AI agent
func New() (*Agent, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.GPT4oMini
	}

	systemMsg := os.Getenv("AGENT_SYSTEM_MESSAGE")
	if systemMsg == "" {
		systemMsg = "You are Relay, a helpful and friendly AI assistant. Answer questions clearly and concisely. Be conversational but informative."
	}

	return &Agent{
		client:    openai.NewClient(apiKey),
		model:     model,
		systemMsg: systemMsg,
	}, nil
}

// Chat sends a conversation history to the AI and returns the assistant's reply
func (a *Agent) Chat(ctx context.Context, history []models.Message, userMessage string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: a.systemMsg,
		},
	}

	// Add conversation history (last 20 messages to stay within token limits)
	start := 0
	if len(history) > 20 {
		start = len(history) - 20
	}
	for _, msg := range history[start:] {
		role := openai.ChatMessageRoleUser
		if msg.Role == "assistant" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add the new user message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       a.model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
