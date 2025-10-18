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
- The API will always return JSON response
- there will be no authentication or authorization
- everything be treated as a single user

## Database Schema

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

### API Endpoints

- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions
- POST /api/study_activities
  - required params: group_id, study_activity_id
- GET /api/words
  - Pagination with 100 words per page  
- GET /api/words/:id
- GET /api/groups
  - Pagination with 100 words per page 
- GET /api/groups/:id/study_sessions
- GET /api/groups/:id
- GET /api/groups/:id/words
- GET /api/study_sessions
  - Pagination with 100 study sessions per page
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words
- POST /api/reset_history
- POST /api/full_reset
- POST /api/study_sessions/:id/words/:word_id/review
  required params: correct 