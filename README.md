# Relay 🚀

A modern AI-powered webchat application built with **Go**, **React**, **PostgreSQL** (via **Supabase**), and deployed on **Vercel**.

---

## Features

- 💬 **Real-time webchat** with an AI agent powered by OpenAI (GPT-4o-mini by default)
- 🧠 **Conversation history** — AI remembers context within a session
- 🗄️ **PostgreSQL persistence** via Supabase — all messages and conversations are stored
- ⚡ **Go backend** with Gin, graceful shutdown, and connection pooling
- ⚛️ **React + Vite frontend** with Tailwind CSS, markdown rendering, and responsive design
- 🚀 **Vercel deployment** for the frontend; containerised backend for any cloud

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│  User Browser                                        │
│  React (Vite) ─── deployed on Vercel                 │
└──────────────────┬───────────────────────────────────┘
                   │ HTTPS / REST API
┌──────────────────▼───────────────────────────────────┐
│  Go Backend (Gin)                                    │
│  ├── POST /api/chat           – send message         │
│  ├── GET  /api/conversations  – list conversations   │
│  ├── GET  /api/conversations/:id/messages            │
│  ├── POST /api/conversations  – create conversation  │
│  ├── DELETE /api/conversations/:id                   │
│  └── GET  /api/health                               │
└──────────┬──────────────────────┬────────────────────┘
           │                      │
┌──────────▼──────────┐  ┌────────▼───────────────────┐
│  Supabase           │  │  OpenAI API                 │
│  PostgreSQL DB      │  │  (GPT-4o-mini / GPT-4o)    │
└─────────────────────┘  └────────────────────────────┘
```

---

## Quick Start (Local)

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- A [Supabase](https://supabase.com/) project **or** any PostgreSQL 15+ instance
- An [OpenAI API key](https://platform.openai.com/api-keys)

### 1. Clone the repository

```bash
git clone https://github.com/nadiramlha/Relay.git
cd Relay
```

### 2. Set up the database

Run the migration against your Supabase project (SQL editor) or any PostgreSQL instance:

```bash
psql "$DATABASE_URL" -f supabase/migrations/001_initial_schema.sql
```

Or paste the contents of `supabase/migrations/001_initial_schema.sql` into the Supabase SQL editor.

### 3. Configure environment variables

```bash
cp .env.example .env
# Edit .env with your DATABASE_URL, OPENAI_API_KEY, etc.
```

### 4a. Run with Docker Compose (recommended)

```bash
docker-compose up --build
```

Frontend: <http://localhost:3000>  
Backend API: <http://localhost:8080>

### 4b. Run manually

**Backend:**
```bash
cd backend
go mod download
go run .
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

---

## Deployment

### Frontend → Vercel

1. Push your repository to GitHub.
2. Import the project in [Vercel](https://vercel.com/new).
3. Set **Root Directory** to `frontend`.
4. Add the environment variable:
   - `VITE_API_URL` → your deployed backend URL (e.g. `https://api.your-domain.com`)
5. Update `frontend/vercel.json` – replace `https://your-backend-url.com` with your actual backend URL.
6. Deploy.

### Backend → Railway / Render / Fly.io

The backend is a standard Go HTTP server. Deploy using the provided `Dockerfile`:

```bash
# Example: deploy to Railway
railway init
railway up
```

Set the following environment variables in your hosting platform:

| Variable | Description |
|---|---|
| `DATABASE_URL` | Supabase / PostgreSQL connection string |
| `OPENAI_API_KEY` | OpenAI API key |
| `OPENAI_MODEL` | Model to use (default: `gpt-4o-mini`) |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed frontend origins |
| `PORT` | Port to listen on (default: `8080`) |
| `GIN_MODE` | Set to `release` in production |

---

## Environment Variables Reference

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | ✅ | — | PostgreSQL connection string |
| `OPENAI_API_KEY` | ✅ | — | OpenAI API key |
| `OPENAI_MODEL` | ❌ | `gpt-4o-mini` | OpenAI model |
| `AGENT_SYSTEM_MESSAGE` | ❌ | built-in | Custom AI system prompt |
| `PORT` | ❌ | `8080` | Backend server port |
| `GIN_MODE` | ❌ | `debug` | `debug` or `release` |
| `ALLOWED_ORIGINS` | ❌ | `*` | CORS allowed origins |
| `VITE_API_URL` | ❌ | (empty) | Frontend: backend API base URL |

---

## API Reference

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/health` | Health check |
| `POST` | `/api/conversations` | Create a conversation |
| `GET` | `/api/conversations` | List all conversations |
| `DELETE` | `/api/conversations/:id` | Delete a conversation |
| `GET` | `/api/conversations/:id/messages` | Get messages for a conversation |
| `POST` | `/api/chat` | Send a message and get AI response |

**POST /api/chat** request body:
```json
{
  "conversation_id": "optional-uuid",
  "message": "Hello, what can you do?"
}
```

**POST /api/chat** response:
```json
{
  "conversation_id": "uuid",
  "user_message": { "id": "...", "role": "user", "content": "...", "created_at": "..." },
  "ai_message":   { "id": "...", "role": "assistant", "content": "...", "created_at": "..." }
}
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 18, Vite, Tailwind CSS |
| Backend | Go 1.24, Gin |
| Database | PostgreSQL 15 via Supabase |
| AI Agent | OpenAI GPT-4o-mini |
| Frontend Hosting | Vercel |
| Container | Docker |

---

## License

MIT