package models

import "time"

// Conversation represents a chat session
type Conversation struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Message represents a single chat message
type Message struct {
	ID             string    `json:"id" db:"id"`
	ConversationID string    `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"` // "user" | "assistant" | "system"
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ChatRequest is the request payload for sending a message
type ChatRequest struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message" binding:"required"`
}

// ChatResponse is the response after processing a chat message
type ChatResponse struct {
	ConversationID string  `json:"conversation_id"`
	UserMessage    Message `json:"user_message"`
	AIMessage      Message `json:"ai_message"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
