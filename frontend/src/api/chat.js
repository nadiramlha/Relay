import axios from 'axios'

const apiBaseURL = import.meta.env.VITE_API_URL || ''

const api = axios.create({
  baseURL: apiBaseURL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 60000,
})

/** Create a new conversation */
export async function createConversation(title = '') {
  const res = await api.post('/api/conversations', { title })
  return res.data
}

/** List all conversations */
export async function getConversations() {
  const res = await api.get('/api/conversations')
  return res.data
}

/** Delete a conversation */
export async function deleteConversation(id) {
  const res = await api.delete(`/api/conversations/${id}`)
  return res.data
}

/** Get messages for a conversation */
export async function getMessages(conversationId) {
  const res = await api.get(`/api/conversations/${conversationId}/messages`)
  return res.data
}

/**
 * Send a chat message; creates a new conversation when conversationId is empty/null.
 * Returns { conversation_id, user_message, ai_message }
 */
export async function sendMessage(conversationId, message) {
  const res = await api.post('/api/chat', {
    conversation_id: conversationId || '',
    message,
  })
  return res.data
}
