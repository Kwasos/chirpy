# Chirpy

A lightweight RESTful social media API built in Go.

## Features

- User authentication using JWTs and refresh tokens
- CRUD endpoints for chirps (posts)
- Profanity filtering middleware
- Webhook integration for Chirpy Red premium status
- PostgreSQL database interactions with SQLC

## Prerequisites

- Go 1.22+
- PostgreSQL

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_secret_key
POLKA_KEY=your_polka_api_key
PLATFORM=dev # (for reset endpoint)
```

## Getting Started

Clone the repository and install dependencies:

```bash
git clone https://github.com/your-username/chirpy.git
cd chirpy
go mod download
```

Run database migrations:

```bash
goose postgres "$DB_URL" up
```

Build and run the server:

```bash
go run .
```

The server will be available at http://localhost:8080.
