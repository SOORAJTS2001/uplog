package db

import (
		"database/sql"
		"time"
		"uplog/models"
		"uplog/utils"


)
func InsertSession(sql_object *sql.DB,sessionID string, createdAt time.Time, expiredAt sql.NullTime, lines int64, bytes int64, isUploaded bool, mode string,tag string) error {
	tx, err := sql_object.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
	INSERT OR REPLACE INTO sessions(session_id, created_at, expired_at, line_count, size_bytes, is_uploaded, mode, tag)
	VALUES(?,?,?,?,?,?,?,?)
	`, sessionID, createdAt.Format(time.RFC3339), utils.NullTimeToString(expiredAt), lines, bytes, utils.BoolToInt(isUploaded), mode, tag)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func IncrementSessionStats(sql_object *sql.DB,sessionID string, addLines int64, addBytes int64) error {
	_, err := sql_object.Exec(`
	UPDATE sessions SET
		line_count = line_count + ?,
		size_bytes = size_bytes + ?
	WHERE session_id = ?
	`, addLines, addBytes, sessionID)
	return err
}

func MarkSessionUploaded(sql_object *sql.DB,sessionID string) error {
	_, err := sql_object.Exec(`UPDATE sessions SET is_uploaded = 1 WHERE session_id = ?`, sessionID)
	return err
}

func ListSessions(sql_object *sql.DB) ([]models.Session, error) {
	rows, err := sql_object.Query(`SELECT session_id, created_at, expired_at, line_count, size_bytes, is_uploaded, mode, tag FROM sessions ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Session
	for rows.Next() {
		var s models.Session
		var createdAtStr string
		var expiredAtStr sql.NullString
		var isUploadedInt int
		if err := rows.Scan(&s.SessionID, &createdAtStr, &expiredAtStr, &s.LineCount, &s.SizeBytes, &isUploadedInt, &s.Mode, &s.Tag); err != nil {
			return nil, err
		}
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		if expiredAtStr.Valid {
			t, _ := time.Parse(time.RFC3339, expiredAtStr.String)
			s.ExpiredAt = sql.NullTime{Time: t, Valid: true}
		} else {
			s.ExpiredAt = sql.NullTime{Valid: false}
		}
		s.IsUploaded = isUploadedInt != 0
		res = append(res, s)
	}
	return res, nil
}

// -------------------- Utilities --------------------
