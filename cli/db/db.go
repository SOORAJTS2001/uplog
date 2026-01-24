package db

import (
	"cli/constants"
	"cli/models"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

func dbPath() string {
	return filepath.Join(constants.ConfigDir, constants.SqliteFileName)
}

func InitDB() *sql.DB {
	var err error
	db, err := sql.Open("sqlite", dbPath())
	if err != nil {
		fmt.Printf("failed to open sqlite: %v\n", err)
		os.Exit(1)
	}
	// busy timeout
	_, _ = db.Exec("PRAGMA busy_timeout = 5000;")

	create := `
	CREATE TABLE IF NOT EXISTS sessions (
		session_id   TEXT PRIMARY KEY,
		created_at   DATETIME NOT NULL,
		expires_at   DATETIME NOT NULL,
		updated_at   DATETIME,
		line_count   INTEGER DEFAULT 0,
		size_bytes   INTEGER DEFAULT 0,
		is_uploaded  INTEGER DEFAULT 0,
		mode         TEXT NOT NULL,
		tag			 TEXT,
		url			 TEXT

	);
	CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	`
	_, err = db.Exec(create)
	if err != nil {
		fmt.Printf("failed to create tables: %v\n", err)
		os.Exit(1)
	}
	return db
}
func ListSessions(db *sql.DB) ([]models.Session, error) {
	rows, err := db.Query(`
		SELECT
			session_id,
			created_at,
			expires_at,
			line_count,
			size_bytes,
			is_uploaded,
			mode,
			tag,
			url
		FROM sessions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		err := rows.Scan(
			&s.SessionId,
			&s.CreatedAt,
			&s.ExpiresAt,
			&s.LineCount,
			&s.SizeBytes,
			&s.IsUploaded,
			&s.Mode,
			&s.Tag,
			&s.Url,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, rows.Err()
}
func GetSessionById(db *sql.DB, sessionID string) (*models.Session, error) {
	query := `
	SELECT
		session_id,
		created_at,
		expires_at,
		line_count,
		size_bytes,
		is_uploaded,
		mode,
		tag,
		url
	FROM sessions
	WHERE session_id = ?
	`

	row := db.QueryRow(query, sessionID)

	var s models.Session
	err := row.Scan(
		&s.SessionId,
		&s.CreatedAt,
		&s.ExpiresAt,
		&s.LineCount,
		&s.SizeBytes,
		&s.IsUploaded,
		&s.Mode,
		&s.Tag,
		&s.Url,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // session not found
		}
		return nil, err
	}

	return &s, nil
}
func UpsertSessionById(db *sql.DB, s *models.Session) error {
	query := `
	INSERT INTO sessions (
		session_id,
		created_at,
		expires_at,
		line_count,
		size_bytes,
		is_uploaded,
		mode,
		tag,
		url
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(session_id) DO UPDATE SET
		created_at  = excluded.created_at,
		expires_at  = excluded.expires_at,
		line_count  = excluded.line_count,
		size_bytes  = excluded.size_bytes,
		is_uploaded = excluded.is_uploaded,
		mode        = excluded.mode,
		tag         = excluded.tag,
		url			= excluded.url
	`

	_, err := db.Exec(query,
		s.SessionId,
		s.CreatedAt,
		s.ExpiresAt,
		s.LineCount,
		s.SizeBytes,
		s.IsUploaded,
		s.Mode,
		s.Tag,
		s.Url,
	)

	return err
}
func UpdateSessionById(db *sql.DB, sessionID string, batch []models.LogEntry) {
	session, err := GetSessionById(db, sessionID)
	if err != nil {
		fmt.Println("Session update failed",err)
	}
	line_count, size_bytes := session.LineCount, session.SizeBytes
	for _, line := range batch {
		line_count += 1
		size_bytes += len(line.Message)
	}
	session.LineCount, session.SizeBytes = line_count, size_bytes
	UpsertSessionById(db, session)
}
func DeleteSessionById(sessionId string,db *sql.DB)(error){
	res, err := db.Exec(
		`DELETE FROM sessions WHERE session_id = ?`,
		sessionId,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		fmt.Println("No matching session with id",sessionId,"found!")
	}

	return nil
}
func DeleteAllSessions(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS sessions`)
	if err != nil {
		return err
	}
	return err
}
