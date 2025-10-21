package models

import "time"

// Word represents a vocabulary word
type Word struct {
	ID       int    `json:"id"`
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts"` // JSON string
}

// Group represents a thematic group of words
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WordGroup represents the many-to-many relationship
type WordGroup struct {
	ID      int `json:"id"`
	WordID  int `json:"word_id"`
	GroupID int `json:"group_id"`
}

// StudyActivity represents a type of learning activity
type StudyActivity struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
}

// StudySession represents a study session
type StudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
}

// WordReviewItem represents a word practice record
type WordReviewItem struct {
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// Response DTOs

type WordResponse struct {
	ID           int    `json:"id,omitempty"`
	Japanese     string `json:"japanese"`
	Romaji       string `json:"romaji"`
	English      string `json:"english"`
	CorrectCount int    `json:"correct_count"`
	WrongCount   int    `json:"wrong_count"`
}

type WordDetailResponse struct {
	ID           int            `json:"id"`
	Japanese     string         `json:"japanese"`
	Romaji       string         `json:"romaji"`
	English      string         `json:"english"`
	CorrectCount int            `json:"correct_count"`
	WrongCount   int            `json:"wrong_count"`
	Groups       []GroupSummary `json:"groups"`
}

type GroupSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GroupResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count"`
}

type StudySessionResponse struct {
	ID               int       `json:"id"`
	ActivityName     string    `json:"activity_name"`
	GroupName        string    `json:"group_name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	ReviewItemsCount int       `json:"review_items_count"`
}

type DashboardLastSession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
	GroupName       string    `json:"group_name"`
}

type DashboardProgress struct {
	TotalWordsStudied   int `json:"total_words_studied"`
	TotalWordsAvailable int `json:"total_words_available"`
}

type DashboardQuickStats struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ReviewResponse struct {
	Success        bool      `json:"success"`
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateSessionResponse struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
}
