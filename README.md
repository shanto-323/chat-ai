# AI Chat Backend API

## Overview

This project is a basic backend API for authentication and AI-powered chat. It is designed as a clean, simple service that demonstrates common backend concepts such as authentication, request validation, persistence, and structured APIs.

The system allows users to:

* Register and log in using email and password
* Receive a JWT access token valid for a fixed duration
* Chat with an AI model using text and/or images
* Retrieve previous chat history with pagination

This project is intentionally kept at a basic and readable level.

---

## Tech Stack

* Language: Go (Golang)
* Database: PostgreSQL
* Authentication: JWT (JSON Web Token)
* Containerization: Docker
* API Style: REST

---

## Base URL

All APIs are versioned and prefixed with:

```
/api/v1
```

---

## Authentication Flow

1. User registers or logs in using email and password
2. Server validates credentials
3. Server returns a JWT access token
4. Token is valid for **1 hour**
5. Token must be sent in subsequent requests using the `Authorization` header

```
Authorization: Bearer <token>
```

---

## Authentication APIs

### POST /api/v1/auth/register

Registers a new user.

Request Body:

```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

Response (200 OK):

```json
{
  "accessToken": "jwt_token_here"
}
```

---

### POST /api/v1/auth/login

Logs in an existing user.

Request Body:

```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

Response (200 OK):

```json
{
  "accessToken": "jwt_token_here"
}
```

---

## Chat APIs

All chat APIs require authentication.

---

### POST /api/v1/chat

Sends a message to the AI system. A request must include **either** a text message or at least one image.

Request Body:

```json
{
  "message_query": "Describe this image",
  "image_s": [
    {
      "url": "https://example.com/image.png"
    }
  ],
  "model_config": {
    "llm_model": "llama_70b",
    "vlm_model": "mock"
  }
}
```

#### ImageData

```json
{
  "url": "string (optional)",
  "base64": "string (optional)",
  "type": "image/jpeg | image/png | image/webp"
}
```

Rules:

* Either `message_query` or `image_s` must be provided
* `model_config` is optional
* If `model_config` is not provided, defaults are applied

Default Model Configuration:

```json
{
  "llm_model": "llama_70b",
  "vlm_model": "mock"
}
```

Response (200 OK):

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "text_query": "Describe this image",
  "image_urls": ["/uploads/abc.png"],
  "response_text": "AI response text",
  "llm_model_name": "llama_70b",
  "vlm_model_name": "mock",
  "timestamp": "Time data created_at"
}
```

---

### GET /api/v1/chat/history

Returns previous chat history for the authenticated user.

Query Parameters:

* `page` (optional, minimum 1, default 1)
* `limit` (optional, min 1, max 100, default 10)
* `order` (optional: asc or desc, default desc)

Example:

```
/api/v1/chat/history?page=1&limit=10&order=desc
```

Response (200 OK):

```json
{
  "data": [
    {
      "text_query": "Hello",
      "image_url": [],
      "response_text": "Hi there"
    }
  ],
  "limit": 10,
  "total": 1,
  "page": 1,
  "total_pages": 3
}
```
