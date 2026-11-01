# Poolworld Pattaya Backend

Backend API for Poolworld Pattaya CRM system built with Go, Echo framework, and MongoDB.

## Features

- JWT Authentication
- Customer Management (CRUD)
- Clean Architecture
- MongoDB Database
- Docker Support

## Prerequisites

- Go 1.21 or higher
- MongoDB 7.0 or higher
- Docker and Docker Compose (optional)

## Getting Started

### Using Docker Compose (Recommended)

1. Copy the environment file:

```bash
cp .env.example .env
```

2. Start the services:

```bash
docker compose up -d
```

The API will be available at `http://localhost:8080`

### Manual Setup

1. Install dependencies:

```bash
go mod download
```

2. Set up environment variables:

```bash
cp .env.example .env
```

3. Make sure MongoDB is running locally or update the `MONGODB_URI` in `.env`

4. Initialize admin user:

```bash
go run cmd/init-admin/main.go
```

5. Run the server:

```bash
go run cmd/server/main.go
```

## API Endpoints

### Authentication

- `POST /api/login` - Login with username and password

### Customers (Protected)

- `GET /api/customers` - Get all customers
- `GET /api/customers/:id` - Get customer by ID
- `POST /api/customers` - Create a new customer
- `PUT /api/customers/:id` - Update customer
- `DELETE /api/customers/:id` - Delete customer

## Default Credentials

- Username: `poolworld`
- Password: `poolworldadmin`

## Project Structure

```
.
├── cmd/
│   ├── server/          # Main application entry point
│   └── init-admin/      # Admin initialization script
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Database connection
│   ├── dto/             # Data Transfer Objects
│   ├── entity/          # Business entities
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── model/           # Database models
│   ├── repository/      # Data access layer
│   └── usecase/         # Business logic
├── docker compose.yml
├── Dockerfile
└── go.mod
```
