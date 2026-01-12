package models

import (
	"time"
)

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
}
type Session struct {
	SessionId  string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	LineCount  int
	SizeBytes  int
	IsUploaded bool
	Mode       string
	Tag        string
	Url		   string
}

type SessionCreateResponse struct {
    SessionId string `json:"session_id"`
}
type UserCreateResponse struct{
	UserId string `json:"user_id"`
}

type Configurations struct{
	UserId string `json:"user_id"`
	ApiKey string `json:"api_key"`
}
