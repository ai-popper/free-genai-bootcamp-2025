package database

import (
	"database/sql"
	"time"

	"github.com/lang-portal/backend/internal/models"
)

// Word queries

func GetWords(page, limit int) ([]models.WordResponse, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get words with correct/wrong counts
	query := `
		SELECT 
			w.id,
			w.japanese,
			w.romaji,
			w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []models.WordResponse
	for rows.Next() {
		var w models.WordResponse
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CorrectCount, &w.WrongCount)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, w)
	}

	return words, total, nil
}

func GetWordByID(id int) (*models.WordDetailResponse, error) {
	query := `
		SELECT 
			w.id,
			w.japanese,
			w.romaji,
			w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`

	var word models.WordDetailResponse
	err := DB.QueryRow(query, id).Scan(
		&word.ID, &word.Japanese, &word.Romaji, &word.English,
		&word.CorrectCount, &word.WrongCount,
	)
	if err != nil {
		return nil, err
	}

	// Get groups
	groupQuery := `
		SELECT g.id, g.name
		FROM groups g
		JOIN words_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?
		ORDER BY g.name
	`

	rows, err := DB.Query(groupQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	word.Groups = []models.GroupSummary{}
	for rows.Next() {
		var g models.GroupSummary
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		word.Groups = append(word.Groups, g)
	}

	return &word, nil
}

// Group queries

func GetGroups(page, limit int) ([]models.GroupResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			g.id,
			g.name,
			COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		GROUP BY g.id
		ORDER BY g.id
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []models.GroupResponse
	for rows.Next() {
		var g models.GroupResponse
		err := rows.Scan(&g.ID, &g.Name, &g.WordCount)
		if err != nil {
			return nil, 0, err
		}
		groups = append(groups, g)
	}

	return groups, total, nil
}

func GetGroupByID(id int) (*models.GroupResponse, error) {
	query := `
		SELECT 
			g.id,
			g.name,
			COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		WHERE g.id = ?
		GROUP BY g.id
	`

	var group models.GroupResponse
	err := DB.QueryRow(query, id).Scan(&group.ID, &group.Name, &group.WordCount)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func GetGroupWords(groupID, page, limit int) ([]models.WordResponse, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`
	var total int
	err := DB.QueryRow(countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			w.id,
			w.japanese,
			w.romaji,
			w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []models.WordResponse
	for rows.Next() {
		var w models.WordResponse
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CorrectCount, &w.WrongCount)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, w)
	}

	return words, total, nil
}

// Study Session queries

func GetStudySessions(page, limit int) ([]models.StudySessionResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			NULL as end_time,
			COUNT(wri.word_id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []models.StudySessionResponse
	for rows.Next() {
		var s models.StudySessionResponse
		var endTime sql.NullTime
		err := rows.Scan(&s.ID, &s.ActivityName, &s.GroupName, &s.StartTime, &endTime, &s.ReviewItemsCount)
		if err != nil {
			return nil, 0, err
		}
		if endTime.Valid {
			s.EndTime = &endTime.Time
		}
		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

func GetStudySessionByID(id int) (*models.StudySessionResponse, error) {
	query := `
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			NULL as end_time,
			COUNT(wri.word_id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.id = ?
		GROUP BY ss.id
	`

	var s models.StudySessionResponse
	var endTime sql.NullTime
	err := DB.QueryRow(query, id).Scan(&s.ID, &s.ActivityName, &s.GroupName, &s.StartTime, &endTime, &s.ReviewItemsCount)
	if err != nil {
		return nil, err
	}
	if endTime.Valid {
		s.EndTime = &endTime.Time
	}

	return &s, nil
}

func GetStudySessionWords(sessionID, page, limit int) ([]models.WordResponse, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(DISTINCT wri.word_id)
		FROM word_review_items wri
		WHERE wri.study_session_id = ?
	`
	var total int
	err := DB.QueryRow(countQuery, sessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT DISTINCT
			w.id,
			w.japanese,
			w.romaji,
			w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []models.WordResponse
	for rows.Next() {
		var w models.WordResponse
		err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CorrectCount, &w.WrongCount)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, w)
	}

	return words, total, nil
}

func GetGroupStudySessions(groupID, page, limit int) ([]models.StudySessionResponse, int, error) {
	offset := (page - 1) * limit

	countQuery := "SELECT COUNT(*) FROM study_sessions WHERE group_id = ?"
	var total int
	err := DB.QueryRow(countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			NULL as end_time,
			COUNT(wri.word_id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.group_id = ?
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []models.StudySessionResponse
	for rows.Next() {
		var s models.StudySessionResponse
		var endTime sql.NullTime
		err := rows.Scan(&s.ID, &s.ActivityName, &s.GroupName, &s.StartTime, &endTime, &s.ReviewItemsCount)
		if err != nil {
			return nil, 0, err
		}
		if endTime.Valid {
			s.EndTime = &endTime.Time
		}
		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

func CreateStudySession(groupID, activityID int) (int, error) {
	result, err := DB.Exec(
		"INSERT INTO study_sessions (group_id, study_activity_id, created_at) VALUES (?, ?, ?)",
		groupID, activityID, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func CreateWordReview(sessionID, wordID int, correct bool) (*models.WordReviewItem, error) {
	now := time.Now()
	_, err := DB.Exec(
		"INSERT INTO word_review_items (study_session_id, word_id, correct, created_at) VALUES (?, ?, ?, ?)",
		sessionID, wordID, correct, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.WordReviewItem{
		StudySessionID: sessionID,
		WordID:         wordID,
		Correct:        correct,
		CreatedAt:      now,
	}, nil
}

// Study Activity queries

func GetStudyActivityByID(id int) (*models.StudyActivity, error) {
	var activity models.StudyActivity
	err := DB.QueryRow(
		"SELECT id, name, thumbnail, description FROM study_activities WHERE id = ?",
		id,
	).Scan(&activity.ID, &activity.Name, &activity.Thumbnail, &activity.Description)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func GetActivityStudySessions(activityID, page, limit int) ([]models.StudySessionResponse, int, error) {
	offset := (page - 1) * limit

	countQuery := "SELECT COUNT(*) FROM study_sessions WHERE study_activity_id = ?"
	var total int
	err := DB.QueryRow(countQuery, activityID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			NULL as end_time,
			COUNT(wri.word_id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.study_activity_id = ?
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, activityID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []models.StudySessionResponse
	for rows.Next() {
		var s models.StudySessionResponse
		var endTime sql.NullTime
		err := rows.Scan(&s.ID, &s.ActivityName, &s.GroupName, &s.StartTime, &endTime, &s.ReviewItemsCount)
		if err != nil {
			return nil, 0, err
		}
		if endTime.Valid {
			s.EndTime = &endTime.Time
		}
		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

// Dashboard queries

func GetLastStudySession() (*models.DashboardLastSession, error) {
	query := `
		SELECT 
			ss.id,
			ss.group_id,
			ss.created_at,
			ss.study_activity_id,
			g.name as group_name
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		ORDER BY ss.created_at DESC
		LIMIT 1
	`

	var session models.DashboardLastSession
	err := DB.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetStudyProgress() (*models.DashboardProgress, error) {
	var progress models.DashboardProgress

	// Total words available
	err := DB.QueryRow("SELECT COUNT(*) FROM words").Scan(&progress.TotalWordsAvailable)
	if err != nil {
		return nil, err
	}

	// Total unique words studied
	err = DB.QueryRow("SELECT COUNT(DISTINCT word_id) FROM word_review_items").Scan(&progress.TotalWordsStudied)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func GetQuickStats() (*models.DashboardQuickStats, error) {
	var stats models.DashboardQuickStats

	// Total study sessions
	err := DB.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&stats.TotalStudySessions)
	if err != nil {
		return nil, err
	}

	// Total active groups (groups with at least one session)
	err = DB.QueryRow(`
		SELECT COUNT(DISTINCT group_id) FROM study_sessions
	`).Scan(&stats.TotalActiveGroups)
	if err != nil {
		return nil, err
	}

	// Success rate
	var totalCorrect, totalWrong int
	err = DB.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END), 0)
		FROM word_review_items
	`).Scan(&totalCorrect, &totalWrong)
	if err != nil {
		return nil, err
	}

	total := totalCorrect + totalWrong
	if total > 0 {
		stats.SuccessRate = float64(totalCorrect) / float64(total) * 100
	}

	// Study streak (consecutive days with sessions)
	stats.StudyStreakDays = calculateStreak()

	return &stats, nil
}

func calculateStreak() int {
	query := `
		SELECT DATE(created_at) as study_date
		FROM study_sessions
		GROUP BY DATE(created_at)
		ORDER BY study_date DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return 0
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0
	}

	streak := 1
	for i := 0; i < len(dates)-1; i++ {
		diff := dates[i].Sub(dates[i+1]).Hours() / 24
		if diff == 1 {
			streak++
		} else {
			break
		}
	}

	return streak
}

// Seed functions

func InsertWord(japanese, romaji, english, parts string) (int64, error) {
	result, err := DB.Exec(
		"INSERT INTO words (japanese, romaji, english, parts) VALUES (?, ?, ?, ?)",
		japanese, romaji, english, parts,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func InsertGroup(name string) (int64, error) {
	result, err := DB.Exec("INSERT OR IGNORE INTO groups (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil || id == 0 {
		// If no insert happened, get the existing ID
		err = DB.QueryRow("SELECT id FROM groups WHERE name = ?", name).Scan(&id)
	}
	return id, err
}

func LinkWordToGroup(wordID, groupID int64) error {
	_, err := DB.Exec(
		"INSERT OR IGNORE INTO words_groups (word_id, group_id) VALUES (?, ?)",
		wordID, groupID,
	)
	return err
}
