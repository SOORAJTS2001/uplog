package main

import (
	"cli/tail/utils"
	"database/sql"
	"time"
	 _ "github.com/mattn/go-sqlite3"
)

func InsertSession(db *sql.DB, sessionID string, createdAt time.Time, expiredAt sql.NullTime, lines int64, bytes int64, isUploaded bool, mode string, tag string) error {
	tx, err := db.Begin()
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

func IncrementSessionStats(db *sql.DB, sessionID string, addLines int64, addBytes int64) error {
	_, err := db.Exec(`
	UPDATE sessions SET
		line_count = line_count + ?,
		size_bytes = size_bytes + ?
	WHERE session_id = ?
	`, addLines, addBytes, sessionID)
	return err
}

func MarkSessionUploaded(db *sql.DB, sessionID string) error {
	_, err := db.Exec(`UPDATE sessions SET is_uploaded = 1 WHERE session_id = ?`, sessionID)
	return err
}
