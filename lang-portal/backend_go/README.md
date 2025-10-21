# Language Learning Portal - Backend API

A Go-based REST API backend for a Japanese language learning portal. This system manages vocabulary inventory, tracks learning progress, and provides a unified interface for various learning activities.

## Features

- **Vocabulary Management**: Store and retrieve Japanese words with romaji and English translations
- **Thematic Groups**: Organize words into thematic collections (greetings, numbers, colors, etc.)
- **Study Sessions**: Track learning sessions with different study activities
- **Progress Tracking**: Monitor learning progress with statistics and success rates
- **Learning Record Store (LRS)**: Record correct/incorrect answers for each word practice

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: SQLite3
- **Task Runner**: Mage

## Project Structure

```
backend_go/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/                  # Data models and DTOs
│   ├── handlers/                # HTTP request handlers
│   ├── database/                # Database connection & queries
│   ├── middleware/              # CORS, logging middleware
│   └── utils/                   # Helper functions
├── db/
│   ├── migrations/              # SQL migration files
│   └── seeds/                   # Seed data JSON files
├── magefile.go                  # Mage task definitions
├── go.mod                       # Go module dependencies
└── words.db                     # SQLite database (generated)
```

## Prerequisites

- Go 1.21 or higher
- Mage (Go task runner)

### Install Mage

```bash
go install github.com/magefile/mage@latest
```

## Getting Started

### 1. Install Dependencies

```bash
cd backend_go
go mod download
```

### 2. Setup Database

Initialize the database, run migrations, and seed data:

```bash
mage setup
```

Or run steps individually:

```bash
mage init      # Create database file
mage migrate   # Run migrations
mage seed      # Import seed data
```

### 3. Run the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Available Mage Tasks

- `mage init` - Initialize the database file
- `mage migrate` - Run database migrations
- `mage seed` - Import seed data from JSON files
- `mage setup` - Run init, migrate, and seed in sequence
- `mage clean` - Remove the database file
- `mage reset` - Clean and setup from scratch

## API Endpoints

### Dashboard

- `GET /api/dashboard/last_study_session` - Get most recent study session
- `GET /api/dashboard/study_progress` - Get overall progress statistics
- `GET /api/dashboard/quick_stats` - Get quick summary stats

### Words

- `GET /api/words` - List all words (paginated)
- `GET /api/words/:id` - Get word details with groups

### Groups

- `GET /api/groups` - List all groups (paginated)
- `GET /api/groups/:id` - Get group details
- `GET /api/groups/:id/words` - Get words in a group
- `GET /api/groups/:id/study_sessions` - Get study sessions for a group

### Study Sessions

- `GET /api/study_sessions` - List all study sessions (paginated)
- `GET /api/study_sessions/:id` - Get study session details
- `GET /api/study_sessions/:id/words` - Get words in a study session
- `POST /api/study_sessions/:id/words/:word_id/review` - Record word review

### Study Activities

- `GET /api/study_activities/:id` - Get study activity details
- `GET /api/study_activities/:id/study_sessions` - Get sessions for an activity
- `POST /api/study_activities` - Create new study session

### Utility

- `POST /api/reset_history` - Reset all study history
- `POST /api/full_reset` - Full database reset
- `GET /health` - Health check endpoint

## API Response Format

All endpoints return JSON responses with consistent structure:

### Paginated Responses

```json
{
  "items": [...],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

### Error Responses

```json
{
  "error": "Error message description"
}
```

## Database Schema

### Tables

- **words** - Vocabulary words (id, japanese, romaji, english, parts)
- **groups** - Thematic groups (id, name)
- **words_groups** - Many-to-many relationship (id, word_id, group_id)
- **study_activities** - Learning activity types (id, name, thumbnail, description)
- **study_sessions** - Study session records (id, group_id, created_at, study_activity_id)
- **word_review_items** - Word practice records (word_id, study_session_id, correct, created_at)

## Adding Seed Data

Create a JSON file in `db/seeds/` with the following format:

```json
[
  {
    "kanji": "こんにちは",
    "romaji": "konnichiwa",
    "english": "hello"
  }
]
```

Update `magefile.go` to include your seed file:

```go
seedConfigs := []SeedConfig{
    {File: "db/seeds/your_file.json", GroupName: "Your Group Name"},
}
```

Then run:

```bash
mage seed
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o server cmd/server/main.go
./server
```

## CORS Configuration

CORS is enabled for all origins by default. To restrict origins, modify `internal/middleware/cors.go`.

## Notes

- No authentication/authorization (single-user system)
- All timestamps are in UTC
- Pagination default: 100 items per page
- Foreign key constraints are enabled

## Example Usage

### Create a Study Session

```bash
curl -X POST http://localhost:8080/api/study_activities \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "study_activity_id": 1
  }'
```

### Record a Word Review

```bash
curl -X POST http://localhost:8080/api/study_sessions/1/words/5/review \
  -H "Content-Type: application/json" \
  -d '{
    "correct": true
  }'
```

### Get Dashboard Stats

```bash
curl http://localhost:8080/api/dashboard/quick_stats
```

## License

MIT
