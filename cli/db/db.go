package db

import (
		"database/sql"
		"fmt"
		"path/filepath"
		"os"
		"cli/constants"
		_ "github.com/mattn/go-sqlite3"


)

func dbPath() string {
	return filepath.Join(constants.ConfigDir, constants.SqliteFileName)
}

func InitDB() {
	var err error
	db, err := sql.Open("sqlite3", dbPath())
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
		expired_at   DATETIME,
		line_count   INTEGER DEFAULT 0,
		size_bytes   INTEGER DEFAULT 0,
		is_uploaded  INTEGER DEFAULT 0,
		mode         TEXT NOT NULL,
		tag			 TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at);
	CREATE INDEX IF NOT EXISTS idx_sessions_expired_at ON sessions(expired_at);
	`
	_, err = db.Exec(create)
	if err != nil {
		fmt.Printf("failed to create tables: %v\n", err)
		os.Exit(1)
	}
}
