# Task API

A Golang Learning Exercise

## Tech Stack

- **Language:** Go 1.25
- **Web Framework:** Gin
- **Database:** MySQL 8.0
- **Driver:** go-sql-driver/mysql
- **Configuration:** godotenv
- **Containerization:** Docker & Docker Compose

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (for containerized setup)
- MySQL 9.5 (if running without Docker)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/task-api.git
cd task-api
```

2. **Set up environment variables**
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```bash
# Server Configuration
SERVER_PORT=8080
ENVIRONMENT=development
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=taskuser
DB_PASSW=taskpass
DB_NAME=taskdb
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

3. **Install dependencies**
```bash
go mod download
```

## Running the Application

### Option 1: Using Docker (Recommended)

Start all services with Docker Compose:
```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

To run in background:
```bash
docker-compose up -d --build
```

View logs:
```bash
docker-compose logs -f api
```

Stop services:
```bash
docker-compose down
```

### Option 2: Running Locally

1. **Start MySQL**
```bash
docker-compose up -d mysql
```

2. **Run the application**
```bash
go run cmd/api/main.go
```

## API Endpoints

### Health Check
```bash
GET /health
```

**Response:**
```json
{
  "status": "healthy"
}
```

### Create Task
```bash
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Learn Go",
  "description": "Build a REST API",
  "status": "pending",
  "priority": 5
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "title": "Learn Go",
  "description": "Build a REST API",
  "status": "pending",
  "priority": 5,
  "created_at": "2026-01-18T20:00:00Z",
  "updated_at": "2026-01-18T20:00:00Z"
}
```

### Get All Tasks
```bash
GET /api/v1/tasks?status=pending&limit=10&offset=0
```

**Query Parameters:**
- `status` (optional): Filter by task status (`pending`, `in_progress`, `completed`)
- `limit` (optional): Number of results (default: 10, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response:** `200 OK`
```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Learn Go",
      "description": "Build a REST API",
      "status": "pending",
      "priority": 5,
      "created_at": "2026-01-18T20:00:00Z",
      "updated_at": "2026-01-18T20:00:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

### Get Task by ID
```bash
GET /api/v1/tasks/:id
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Learn Go",
  "description": "Build a REST API",
  "status": "pending",
  "priority": 5,
  "created_at": "2026-01-18T20:00:00Z",
  "updated_at": "2026-01-18T20:00:00Z"
}
```

### Update Task
```bash
PUT /api/v1/tasks/:id
Content-Type: application/json

{
  "status": "completed",
  "priority": 10
}
```

**Note:** All fields are optional. Only provided fields will be updated.

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Learn Go",
  "description": "Build a REST API",
  "status": "completed",
  "priority": 10,
  "created_at": "2026-01-18T20:00:00Z",
  "updated_at": "2026-01-18T20:30:00Z"
}
```

### Delete Task
```bash
DELETE /api/v1/tasks/:id
```

**Response:** `200 OK`
```json
{
  "message": "Task deleted successfully"
}
```

## Task Status Values

- `pending` - Task is pending (default)
- `in_progress` - Task is in progress
- `completed` - Task is completed

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Building the Application
```bash
# Build for current OS
go build -o bin/task-api ./cmd/api

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/task-api-linux ./cmd/api

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/task-api.exe ./cmd/api
```

## Docker Commands
```bash
# Build and start services
docker-compose up --build

# Start services in background
docker-compose up -d

# View logs
docker-compose logs -f api
docker-compose logs -f mysql

# Stop services
docker-compose down

# Stop and remove volumes (fresh start)
docker-compose down -v

# Rebuild just the API service
docker-compose up -d --build api

# Execute commands in running containers
docker-compose exec api sh
docker-compose exec mysql mysql -u root -p
```