package setup

import (
	"os"
	"fmt"
	"database/sql"
	"path/filepath"
	"uplog/config"
	_ "github.com/mattn/go-sqlite3"

)

func UserHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	// fallback
	dir, _ := os.UserHomeDir()
	return dir
}

func SetupDirectories(home string) {
	base := filepath.Join(home, config.BaseDir)
	_ = os.MkdirAll(filepath.Join(base, config.ConfigDir), 0o700)
	_ = os.MkdirAll(filepath.Join(base, config.TmpDir), 0o700)
}

func DbPath(home string) string {
	return filepath.Join(home, config.BaseDir, config.ConfigDir, config.SqliteFileName)
}

func InitDB(sql_object *sql.DB,home string) (*sql.DB){
	var err error
	sql_object, err = sql.Open("sqlite3", DbPath(home))
	if err != nil {
		fmt.Printf("failed to open sqlite: %v\n", err)
		os.Exit(1)
	}
	// busy timeout
	_, _ = sql_object.Exec("PRAGMA busy_timeout = 5000;")

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
	_, err = sql_object.Exec(create)
	if err != nil {
		fmt.Printf("failed to create tables: %v\n", err)
		os.Exit(1)
	}
	return sql_object
}
