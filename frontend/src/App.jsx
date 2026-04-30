import { useState, useEffect, useRef } from 'react'
import Sidebar from './components/Sidebar.jsx'
import ChatWindow from './components/ChatWindow.jsx'
import { getConversations, sendMessage, getMessages } from './api/chat.js'

export default function App() {
  const [conversations, setConversations] = useState([])
  const [activeConvId, setActiveConvId] = useState(null)
  const [messages, setMessages] = useState([])
  const [loading, setLoading] = useState(false)
  const [sidebarOpen, setSidebarOpen] = useState(true)

  // Load conversations on mount
  useEffect(() => {
    loadConversations()
  }, [])

  // Load messages when active conversation changes
  useEffect(() => {
    if (activeConvId) {
      loadMessages(activeConvId)
    }
  }, [activeConvId])

  async function loadConversations() {
    try {
      const data = await getConversations()
      setConversations(data || [])
    } catch (err) {
      console.error('Failed to load conversations:', err)
    }
  }

  async function loadMessages(convId) {
    try {
      const data = await getMessages(convId)
      setMessages(data || [])
    } catch (err) {
      console.error('Failed to load messages:', err)
      setMessages([])
    }
  }

  async function handleSend(text) {
    if (!text.trim() || loading) return

    // Optimistically add user message
    const tempUserMsg = {
      id: `temp-user-${Date.now()}`,
      role: 'user',
      content: text,
      created_at: new Date().toISOString(),
    }
    setMessages((prev) => [...prev, tempUserMsg])
    setLoading(true)

    try {
      const data = await sendMessage(activeConvId, text)

      // Switch to the (possibly new) conversation
      if (!activeConvId || data.conversation_id !== activeConvId) {
        setActiveConvId(data.conversation_id)
      }

      // Replace temp message + add AI response
      setMessages((prev) => {
        const withoutTemp = prev.filter((m) => m.id !== tempUserMsg.id)
        return [...withoutTemp, data.user_message, data.ai_message]
      })

      // Refresh conversation list (title may have changed)
      await loadConversations()
    } catch (err) {
      console.error('Chat error:', err)
      // Remove temp message on error
      setMessages((prev) => prev.filter((m) => m.id !== tempUserMsg.id))
    } finally {
      setLoading(false)
    }
  }

  function handleNewChat() {
    setActiveConvId(null)
    setMessages([])
  }

  function handleSelectConversation(id) {
    setActiveConvId(id)
  }

  async function handleDeleteConversation(id) {
    setConversations((prev) => prev.filter((c) => c.id !== id))
    if (activeConvId === id) {
      setActiveConvId(null)
      setMessages([])
    }
  }

  return (
    <div className="flex h-screen bg-gray-50 overflow-hidden">
      {/* Sidebar */}
      <Sidebar
        open={sidebarOpen}
        conversations={conversations}
        activeConvId={activeConvId}
        onNewChat={handleNewChat}
        onSelect={handleSelectConversation}
        onDelete={handleDeleteConversation}
        onToggle={() => setSidebarOpen((v) => !v)}
      />

      {/* Chat area */}
      <ChatWindow
        messages={messages}
        loading={loading}
        onSend={handleSend}
        onToggleSidebar={() => setSidebarOpen((v) => !v)}
        sidebarOpen={sidebarOpen}
      />
    </div>
  )
}
