import { useEffect, useRef } from 'react'
import MessageItem from './MessageItem.jsx'
import MessageInput from './MessageInput.jsx'
import TypingIndicator from './TypingIndicator.jsx'

export default function ChatWindow({ messages, loading, onSend, onToggleSidebar, sidebarOpen }) {
  const bottomRef = useRef(null)

  // Auto-scroll to latest message
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, loading])

  const isEmpty = messages.length === 0 && !loading

  return (
    <div className="flex flex-col flex-1 min-w-0 h-screen">
      {/* Top bar */}
      <header className="flex items-center gap-3 px-4 py-3 border-b border-gray-200 bg-white shadow-sm">
        {!sidebarOpen && (
          <button
            onClick={onToggleSidebar}
            className="text-gray-400 hover:text-gray-600 transition-colors p-1 rounded"
            title="Open sidebar"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" />
            </svg>
          </button>
        )}
        <div className="w-8 h-8 rounded-full bg-relay-500 flex items-center justify-center text-white font-bold text-sm">
          R
        </div>
        <div>
          <h1 className="font-semibold text-gray-800 text-sm">Relay AI</h1>
          <p className="text-xs text-green-500">Online</p>
        </div>
      </header>

      {/* Message area */}
      <div className="flex-1 overflow-y-auto px-4 py-6">
        {isEmpty ? (
          <WelcomeScreen />
        ) : (
          <div className="max-w-3xl mx-auto space-y-4">
            {messages.map((msg) => (
              <MessageItem key={msg.id} message={msg} />
            ))}
            {loading && <TypingIndicator />}
            <div ref={bottomRef} />
          </div>
        )}
        {!isEmpty && <div ref={bottomRef} />}
      </div>

      {/* Input */}
      <div className="border-t border-gray-200 bg-white px-4 py-3">
        <div className="max-w-3xl mx-auto">
          <MessageInput onSend={onSend} disabled={loading} />
          <p className="text-xs text-gray-400 text-center mt-2">
            Relay can make mistakes. Consider checking important information.
          </p>
        </div>
      </div>
    </div>
  )
}

function WelcomeScreen() {
  return (
    <div className="flex flex-col items-center justify-center h-full text-center px-4 min-h-[60vh]">
      <div className="w-16 h-16 rounded-full bg-relay-500 flex items-center justify-center text-white font-bold text-3xl mb-4 shadow-lg">
        R
      </div>
      <h2 className="text-2xl font-bold text-gray-800 mb-2">How can I help you today?</h2>
      <p className="text-gray-500 max-w-sm text-sm">
        I&apos;m Relay, your AI assistant. Ask me anything — I&apos;m here to help with questions,
        writing, coding, analysis, and more.
      </p>
      <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 gap-3 max-w-lg w-full">
        {SUGGESTIONS.map((s) => (
          <SuggestionCard key={s.title} {...s} />
        ))}
      </div>
    </div>
  )
}

function SuggestionCard({ icon, title, description }) {
  return (
    <div className="border border-gray-200 rounded-xl p-4 text-left hover:border-relay-300 hover:bg-relay-50 transition-colors cursor-default">
      <div className="text-2xl mb-1">{icon}</div>
      <div className="font-medium text-sm text-gray-800">{title}</div>
      <div className="text-xs text-gray-500 mt-0.5">{description}</div>
    </div>
  )
}

const SUGGESTIONS = [
  { icon: '💡', title: 'Brainstorm ideas', description: 'Generate creative ideas for any project' },
  { icon: '📝', title: 'Help with writing', description: 'Draft, edit, or improve any text' },
  { icon: '🖥️', title: 'Explain code', description: 'Understand or debug code snippets' },
  { icon: '🔍', title: 'Answer questions', description: 'Get clear, concise answers' },
]
