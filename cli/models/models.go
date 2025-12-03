package models

import (
	"time"
	"database/sql"
)

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
}
type Session struct {
	SessionId  string
	CreatedAt  time.Time
	ExpiredAt  sql.NullTime
	LineCount  int64
	SizeBytes  int64
	IsUploaded bool
	Mode       string
	Tag        string
}

type SessionCreateResponse struct {
    SessionId string `json:"session_id"`
}
type UserCreateResponse struct{
	UserId string `json:"user_id"`
}
