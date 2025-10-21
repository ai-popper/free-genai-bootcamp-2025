# Backend Server Technical Specs

## Business Goal:

A language learning school wants to build a prototype of learning portal which will act as three things:
-Inventory of possible vocabulary that can be learned
-Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
-A unified launchpad to launch different learning apps

## Technical Requirements

- The backend will be built using Go language
- The database will be SQLite3
- The API will be built using Gin
- Mage is a task runner for Go
- The API will always return JSON response
- there will be no authentication or authorization
- everything be treated as a single user

## Project Structure

```
backend_go/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── models/                  # Database models
│   ├── handlers/                # HTTP handlers (controllers)
│   ├── database/                # Database connection & queries
│   ├── middleware/              # CORS, logging, etc.
│   └── utils/                   # Helper functions
├── db/                          # Unified database folder
│   ├── migrations/              # SQL migration files
│   └── seeds/                   # Seed data JSON files
├── magefile.go                  # Mage task definitions
├── go.mod                       # Go module file
├── go.sum                       # Dependencies
├── words.db                     # SQLite database (created by tasks)
└── README.md
```

## Database Schema

Our database will be a single SQLite database called `words.db`
that will be in a root of the project folder of `backend_go`

we have following tables

- words - stores vocabulary words
  - id integer
  - japenese string
  - romaji string
  - english string
  - parts json
- words_groups - joins table for words and groups many-to-many
  - id integer
  - word_id integer
  - group_id integer
- groups - thematic groups of words
  - id integer
  - name string
- study_sessions- records of study_sessions grouping word_review_items
  - id integer
  - group_id integer
  - created_at datetime
  - study_activity_id integer
- study_activities- a specific study_activities, linking a study_sessions to groups
  - id integer
  - study_sessions_id integer
  - group_id integer
  - created_at datetime
- word_review_items- a record of word practice, detrermining if the word was correct or not
  - word_id integer
  - study_sessions_id integer
  - correct boolean
  - created_at datetime

## API Endpoints

### Dashboard Endpoints

#### GET /api/dashboard/last_study_session
Returns the most recent study session with its details.

**Response:**
```json
{
  "id": 42,
  "group_id": 5,
  "created_at": "2025-10-19T15:30:00Z",
  "study_activity_id": 12,
  "group_id": 456,
  "group_name": "Basic Greetings"
}
```

#### GET /api/dashboard/study_progress
Returns overall study progress statistics across all groups.
Please note that the frontend will determine progress bar based on total word studied and total available words.

**Response:**
```json
{
  "total_words_studied": 124,
  "total_words_available": 850,
}
```

#### GET /api/dashboard/quick_stats
Returns quick summary statistics for dashboard display.

**Response:**
```json
{
  "success_rate": 80.0,
  "total_study_sessions": 127,
  "total_active_groups": 8,
  "study_streak_days": 7
}
```

### Study Activities Endpoints

#### GET /api/study_activities/:id
Returns details of a specific study activity including name, thumbnail, and description.

**Response:**
```json
{
  "id": 12,
  "name": "Flashcard Practice",
  "thumbnail": "https://example.com/thumbnails/flashcard.png",
  "description": "Practice vocabulary using interactive flashcards"
}  
```

#### GET /api/study_activities/:id/study_sessions
Returns paginated list of study sessions for a specific study activity.

**Response:**
```json
{
  "items": [
    {
      "id": 38,
      "activity_name": "Flashcard Practice",
      "group_name": "Basic Greetings",
      "start_time": "2025-10-18T14:20:00Z",
      "end_time": "2025-10-18T14:35:00Z",
      "review_items_count": 15
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 1,
    "total_items": 2
  }
}
```

#### POST /api/study_activities
Creates a new study session for a study activity with a selected group.

*Request Params*
  - group_id: integers,
  - study_activity_id: integers


**Response:**
```json
{
  "id": 43,
  "group_id": 5
}
```

### Words Endpoints

#### GET /api/words
Returns paginated list of vocabulary words (100 per page).

**Query Parameters:**
- `page` (optional, default: 1)

**Response:**
```json
{
  "items": [
    {
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 38,
      "wrong_count": 7
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 9,
    "total_items": 850
  }
}
```

#### GET /api/words/:id
Returns details of a specific word.

**Response:**
```json
{
  "id": 1,
  "japanese": "こんにちは",
  "romaji": "konnichiwa",
  "english": "hello",
  "correct_count": 38,
  "wrong_count": 7,
  "groups": [
    {
      "id": 5,
      "name": "Basic Greetings"
    },
    {
      "id": 12,
      "name": "Daily Conversation"
    }
  ]
}
```

### Groups Endpoints

#### GET /api/groups
Returns paginated list of word groups (100 per page).

**Query Parameters:**
- `page` (optional, default: 1)

**Response:**
```json
{
  "items": [
    {
      "id": 5,
      "name": "Basic Greetings",
      "word_count": 25
    },
    {
      "id": 7,
      "name": "Numbers",
      "word_count": 30
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 1,
    "total_items": 15
  }
}
```

#### GET /api/groups/:id
Returns details of a specific group.

**Response:**
```json
{
  "id": 5,
  "name": "Basic Greetings",
  "word_count": 25
}
```

#### GET /api/groups/:id/study_sessions
Returns paginated list of study sessions for a specific group.

**Response:**
```json
{
  "items": [
    {
      "id": 42,
      "activity_name": "Flashcard Practice",
      "group_name": "Basic Greetings",
      "start_time": "2025-10-19T15:30:00Z",
      "end_time": "2025-10-19T15:42:00Z",
      "review_items_count": 15
    },
    {
      "id": 38,
      "activity_name": "Flashcard Practice",
      "group_name": "Basic Greetings",
      "start_time": "2025-10-18T14:20:00Z",
      "end_time": "2025-10-18T14:35:00Z",
      "review_items_count": 15
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 1,
    "total_items": 23
  }
}
```

#### GET /api/groups/:id/words
Returns paginated list of words in a specific group.

**Response:**
```json
{
  "items": [
    {
      "id": 1,
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 38,
      "wrong_count": 7
    },
    {
      "id": 2,
      "japanese": "ありがとう",
      "romaji": "arigatou",
      "english": "thank you",
      "correct_count": 42,
      "wrong_count": 3
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 1,
    "total_items": 25
  }
}
```

### Study Sessions Endpoints

#### GET /api/study_sessions
Returns paginated list of study sessions (100 per page).

**Query Parameters:**
- `page` (optional, default: 1)

**Response:**
```json
{
  "items": [
    {
      "id": 42,
      "activity_name": "Flashcard Practice",
      "group_name": "Basic Greetings",
      "start_time": "2025-10-19T15:30:00Z",
      "end_time": "2025-10-19T15:42:00Z",
      "review_items_count": 15
    },
    {
      "id": 41,
      "activity_name": "Quiz Mode",
      "group_name": "Numbers",
      "start_time": "2025-10-19T14:15:00Z",
      "end_time": "2025-10-19T14:35:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 2,
    "total_items": 127
  }
}
```

#### GET /api/study_sessions/:id
Returns details of a specific study session.

**Response:**
```json
{
  "id": 42,
  "activity_name": "Flashcard Practice",
  "group_name": "Basic Greetings",
  "start_time": "2025-10-19T15:30:00Z",
  "end_time": "2025-10-19T15:42:00Z",
  "review_items_count": 15
}
```

#### GET /api/study_sessions/:id/words
Returns paginated list of word review items in a specific study session.

**Response:**
```json
{
  "items": [
    {
      "id": 1,
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 38,
      "wrong_count": 7
    },
    {
      "id": 2,
      "japanese": "ありがとう",
      "romaji": "arigatou",
      "english": "thank you",
      "correct_count": 42,
      "wrong_count": 3
    },
    {
      "id": 3,
      "japanese": "さようなら",
      "romaji": "sayounara",
      "english": "goodbye",
      "correct_count": 25,
      "wrong_count": 8
    }
  ],
  "pagination": {
    "current_page": 1,
    "items_per_page": 100,
    "total_pages": 1,
    "total_items": 15
  }
}
```

### Utility Endpoints

#### POST /api/reset_history
Resets all study history (sessions and review items) but keeps vocabulary and groups.

**Response:**
```json
{
  "success": true,
  "message": "Study history reset successfully"
}
```

#### POST /api/full_reset
Performs a complete database reset, removing all data and recreating with seed data.

**Response:**
```json
{
  "success": true,
  "message": "Full database reset completed successfully"
}
```

#### POST /api/study_sessions/:id/words/:word_id/review
Records a review result for a word in a study session.

*Request Parameters:*
  -id(study_session_id): integers,
  -word_id: integers,
  -correct: boolean
  

**Request Payload:**
```json
{
  "correct": true
}
```

**Response:**
```json
{
  "success": true,
  "word_id": 1,
  "study_session_id":123,
  "correct": true,
  "created_at": "2025-10-19T15:42:00Z"
}
```

## Task Runner Tasks
Mage is a task runner for Go.
Lets list out possible tasks we need for our lang Portal.

### Inititalize Database
This task will Inititalize the SQLite database called `words.db`

### Migrate Database
This task will run a series of migrations sql files on the database.

migrations live in the `migrations` folder.
The migration file will run in the order of their file name.
The file name should look like this:

```sql
0001_init_db.sql
0002_create_words_table.sql
```

### Seed Data
this task will import seed data and transform them into target data for our database.

All seed files live in the `seeds` folder.

In our task we should  have DSL to specify each seed file and its expected group name.

```json
[
  {
    "kanji": "払う",
    "romaji": "harau",
    "english": "to pay",
  },
  ...
]
```



