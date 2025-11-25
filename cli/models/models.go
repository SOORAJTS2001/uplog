package models


import (
	"time"
	"database/sql"

)
type Session struct {
	SessionID  string
	CreatedAt  time.Time
	ExpiredAt  sql.NullTime
	LineCount  int64
	SizeBytes  int64
	IsUploaded bool
	Mode       string
	Tag        string
}

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
}

type SessionCreateResponse struct {
    SessionID string `json:"session_id"`
}
