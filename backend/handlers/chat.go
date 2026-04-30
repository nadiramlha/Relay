package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nadiramlha/relay/backend/agent"
	"github.com/nadiramlha/relay/backend/db"
	"github.com/nadiramlha/relay/backend/models"
)

// Handler holds the application dependencies
type Handler struct {
	DB    *db.DB
	Agent *agent.Agent
}

// New creates a new Handler
func New(database *db.DB, aiAgent *agent.Agent) *Handler {
	return &Handler{DB: database, Agent: aiAgent}
}

// CreateConversation creates a new conversation
func (h *Handler) CreateConversation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var req struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Title = "New Conversation"
	}
	if req.Title == "" {
		req.Title = "New Conversation"
	}

	conv := models.Conversation{}
	err := h.DB.Pool.QueryRow(ctx,
		`INSERT INTO conversations (title) VALUES ($1)
		 RETURNING id, title, created_at, updated_at`,
		req.Title,
	).Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create conversation"})
		return
	}

	c.JSON(http.StatusCreated, conv)
}

// GetConversations lists all conversations ordered by most recent
func (h *Handler) GetConversations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.DB.Pool.Query(ctx,
		`SELECT id, title, created_at, updated_at FROM conversations ORDER BY updated_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch conversations"})
		return
	}
	defer rows.Close()

	conversations := []models.Conversation{}
	for rows.Next() {
		var conv models.Conversation
		if err := rows.Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			continue
		}
		conversations = append(conversations, conv)
	}

	c.JSON(http.StatusOK, conversations)
}

// GetMessages returns all messages for a conversation
func (h *Handler) GetMessages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	convID := c.Param("id")
	if _, err := uuid.Parse(convID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid conversation id"})
		return
	}

	rows, err := h.DB.Pool.Query(ctx,
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`,
		convID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch messages"})
		return
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	c.JSON(http.StatusOK, messages)
}

// Chat handles a new user message, processes it through the AI agent, and stores the exchange
func (h *Handler) Chat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "message is required"})
		return
	}

	// Create or retrieve the conversation
	convID := req.ConversationID
	if convID == "" {
		var conv models.Conversation
		err := h.DB.Pool.QueryRow(ctx,
			`INSERT INTO conversations (title) VALUES ($1)
			 RETURNING id, title, created_at, updated_at`,
			"New Conversation",
		).Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create conversation"})
			return
		}
		convID = conv.ID
	} else if _, err := uuid.Parse(convID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid conversation_id"})
		return
	}

	// Fetch conversation history
	rows, err := h.DB.Pool.Query(ctx,
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`,
		convID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch history"})
		return
	}
	history := []models.Message{}
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			continue
		}
		history = append(history, msg)
	}
	rows.Close()

	// Persist user message
	var userMsg models.Message
	err = h.DB.Pool.QueryRow(ctx,
		`INSERT INTO messages (conversation_id, role, content) VALUES ($1, 'user', $2)
		 RETURNING id, conversation_id, role, content, created_at`,
		convID, req.Message,
	).Scan(&userMsg.ID, &userMsg.ConversationID, &userMsg.Role, &userMsg.Content, &userMsg.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to save user message"})
		return
	}

	// Update conversation title if this is the first user message
	if len(history) == 0 {
		title := req.Message
		if len(title) > 60 {
			title = title[:60] + "…"
		}
		_, _ = h.DB.Pool.Exec(ctx,
			`UPDATE conversations SET title = $1, updated_at = NOW() WHERE id = $2`,
			title, convID,
		)
	} else {
		_, _ = h.DB.Pool.Exec(ctx,
			`UPDATE conversations SET updated_at = NOW() WHERE id = $1`, convID,
		)
	}

	// Call the AI agent
	aiContent, err := h.Agent.Chat(ctx, history, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "AI agent error: " + err.Error()})
		return
	}

	// Persist AI response
	var aiMsg models.Message
	err = h.DB.Pool.QueryRow(ctx,
		`INSERT INTO messages (conversation_id, role, content) VALUES ($1, 'assistant', $2)
		 RETURNING id, conversation_id, role, content, created_at`,
		convID, aiContent,
	).Scan(&aiMsg.ID, &aiMsg.ConversationID, &aiMsg.Role, &aiMsg.Content, &aiMsg.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to save AI message"})
		return
	}

	c.JSON(http.StatusOK, models.ChatResponse{
		ConversationID: convID,
		UserMessage:    userMsg,
		AIMessage:      aiMsg,
	})
}

// DeleteConversation removes a conversation and all its messages
func (h *Handler) DeleteConversation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	convID := c.Param("id")
	if _, err := uuid.Parse(convID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid conversation id"})
		return
	}

	result, err := h.DB.Pool.Exec(ctx,
		`DELETE FROM conversations WHERE id = $1`, convID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to delete conversation"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "conversation deleted"})
}

// HealthCheck returns a simple health status
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.DB.Pool.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
