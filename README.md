# Chat-AI - Backend API

Backend service for authentication and multi-modal AI chat. Built with Go, PostgreSQL, and Docker.

## Quick Links

- **API Docs:** `http://localhost:8080/docs`
- **Base URL:** `/api/v1`
- **Auth:** JWT tokens (1-hour expiry)

## Tech Stack

Go • PostgreSQL • JWT • Docker

## Getting Started

### Prerequisites
- Go 1.20+
- PostgreSQL 12+
- Docker 

### Run Locally

```bash
go run ./cmd/main
```

Visit `http://localhost:8080/docs` for interactive API documentation.


## Authentication

All endpoints except `/auth/register` and `/auth/login` require a JWT token.

**Include token in requests:**
```
Authorization:<token>
```

---

## API Endpoints

### Register
**POST** `/api/v1/auth/register`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "accessToken": "jwt_token_here"
}
```

### Login
**POST** `/api/v1/auth/login`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "accessToken": "jwt_token_here"
}
```

### Chat
**POST** `/api/v1/chat` (requires auth)

Send text and/or images. At least one must be provided.

```json
{
  "message_query": "What's in this image?",
  "image_s": [
    {
      "url": "https://example.com/image.jpg"
    }
  ],
  "model_config": {
    "llm_model": "llama_70b",
    "vlm_model": "mock"
  }
}
```

**Image Format:**
```json
{
  "url": "string (optional)",
  "base64": "string (optional)",
  "type": "image/jpeg | image/png | image/webp"
}
```

Response:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "query": "What's in this image?",
  "response_text": "AI response here",
  "image_urls": [],
  "llm_model_name": "llama_70b",
  "vlm_model_name": "mock",
  "timestamp": "2025-12-14T10:30:00Z"
}
```

**Defaults:**
- `llm_model`: `llama_70b`
- `vlm_model`: `mock`

### Chat History
**GET** `/api/v1/chat/history` (requires auth)

Query Parameters:
- `page` (default: 1, min: 1)
- `limit` (default: 10, min: 1, max: 100)
- `order` (default: desc, options: asc, desc)

Example: `/api/v1/chat/history?page=1&limit=20&order=desc`

Response:
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "query": "What's in this image?",
      "response_text": "AI response here",
      "image_urls": [],
      "llm_model_name": "llama_70b",
      "vlm_model_name": "mock",
      "timestamp": "2025-12-14T10:30:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 45,
  "total_pages": 3
}
```

## Error Response

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

Common Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad request
- `401` - Unauthorized (invalid/expired token)
- `500` - Server error

## Project Structure

```
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config_test.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   ├── database.go
│   │   ├── migrations
│   │   │   └── 001_user.sql
│   │   ├── migrator.go
│   │   ├── mock
│   │   │   ├── conversation_log.go
│   │   │   ├── mock.go
│   │   │   └── user.go
│   │   └── postgres
│   │       ├── conversation_log.go
│   │       ├── postgres.go
│   │       └── user.go
│   ├── server
│   │   ├── errs
│   │   │   ├── http.go
│   │   │   └── types.go
│   │   ├── handler
│   │   │   ├── auth.go
│   │   │   ├── auth_test.go
│   │   │   ├── base.go
│   │   │   ├── chat.go
│   │   │   ├── chat_test.go
│   │   │   ├── hendlers.go
│   │   │   └── openapi.go
│   │   ├── manager
│   │   │   ├── llm
│   │   │   │   ├── llm_manager.go
│   │   │   │   └── openrouter
│   │   │   │       ├── openrouter.go
│   │   │   │       └── openrouter_test.go
│   │   │   ├── manager.go
│   │   │   └── vlm
│   │   │       ├── mock
│   │   │       │   └── mock.go
│   │   │       └── vlm_manager.go
│   │   ├── middleware
│   │   │   ├── auth.go
│   │   │   ├── context.go
│   │   │   ├── global.go
│   │   │   ├── middlewares.go
│   │   │   ├── rate_limit.go
│   │   │   └── request_id.go
│   │   ├── router
│   │   │   ├── router.go
│   │   │   ├── system.go
│   │   │   └── v1
│   │   │       ├── auth.go
│   │   │       ├── chat.go
│   │   │       └── v1.go
│   │   ├── server.go
│   │   └── validation
│   │       └── utils.go
│   └── service
│       ├── auth.go
│       ├── chat.go
│       ├── image
│       │   └── image.go
│       └── services.go
├── model
│   ├── base.go
│   ├── base_lv.go
│   ├── dto
│   │   ├── chat.go
│   │   ├── conversation_log.go
│   │   ├── llm.go
│   │   ├── user.go
│   │   └── vlm.go
│   ├── entity
│   │   ├── conversation_log.go
│   │   └── user.go
│   ├── paginate.go
│   └── vlm_result.go
├── pkg
│   ├── hash.go
│   ├── jwt.go
│   └── logger
│       └── logger.go
├── scripts
│   └── Dockerfile
├── sqlerr
│   ├── error.go
│   └── handler.go
└── static
    ├── openapi.html
    └── openapi.json
```

## Notes

- Token expires after 1 hour - users must login again
- Images can be sent via URL or Base64
- Chat history is user-specific and paginated
- All timestamps are in ISO 8601 format