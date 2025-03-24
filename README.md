# AI-Powered Personal Knowledge Base (PKB)

A local-first knowledge base that uses AI embeddings to organize notes, code snippets, bookmarks, and documents.

## Features

- Store and organize notes, code snippets, bookmarks, and documents
- AI-powered search using embeddings
- Content summarization using OpenAI/Anthropic APIs
- Local storage with SQLite
- Simple web UI for access

## Tech Stack

- Backend: Go
- Frontend: Svelte
- Database: SQLite with vector search extension
- AI: OpenAI/Anthropic APIs for embeddings and summarization

## Getting Started

### Prerequisites

- Go 1.18+
- Node.js 14+
- npm or yarn

### Setup

1. Clone the repository
2. Set up your API keys in `.env` file
3. Run the backend: `go run cmd/server/main.go`
4. Run the frontend: `cd web && npm install && npm run dev`

## Usage

1. Add content through the web UI
2. Search using natural language
3. View and manage your knowledge base
