package utils

import (
	"database/sql"
	"time"
	"strings"
)


func NullTimeToString(n sql.NullTime) interface{} {
	if n.Valid {
		return n.Time.Format(time.RFC3339)
	}
	return nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
func DetectLevel(line string) string {
	up := strings.ToUpper(line)

	switch {
	case strings.Contains(up, "ERROR"):
		return "ERROR"
	case strings.Contains(up, "WARN"):
		return "WARN"
	case strings.Contains(up, "DEBUG"):
		return "DEBUG"
	case strings.Contains(up, "INFO"):
		return "INFO"
	default:
		return "INFO"
	}
}
