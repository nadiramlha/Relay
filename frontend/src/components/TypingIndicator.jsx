export default function TypingIndicator() {
  return (
    <div className="flex gap-3">
      {/* Avatar */}
      <div className="w-8 h-8 rounded-full bg-relay-500 flex items-center justify-center text-white font-bold text-sm flex-shrink-0 mt-1">
        R
      </div>
      {/* Dots */}
      <div className="chat-bubble-ai flex items-center gap-1 py-3">
        <span className="typing-dot w-2 h-2 rounded-full bg-gray-400 inline-block" />
        <span className="typing-dot w-2 h-2 rounded-full bg-gray-400 inline-block" />
        <span className="typing-dot w-2 h-2 rounded-full bg-gray-400 inline-block" />
      </div>
    </div>
  )
}
