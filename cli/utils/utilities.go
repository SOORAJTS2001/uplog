package utils

import (
	"flag"
	"time"
	"fmt"
	"os"
	"strings"
	"net/http"
	"encoding/json"
	"database/sql"
	"uplog/config"
	"uplog/models"
)

    // command := fs.Arg(0)
    // commandArgs := fs.Args()[1:]

    // commands.RunCmd(sql_object,home, command, commandArgs,respJSON,pollInterval,tag)
func RunCmdWithFlags(home string, args []string,pollInterval time.Duration,tag*string) (string, []string, time.Duration, *string) {
    fs := flag.NewFlagSet("uplog run", flag.ExitOnError)

    // poll is an int (in milliseconds)
    poll := fs.Int("poll", int(config.PollIntervalLimit/time.Millisecond), "Polling interval in milliseconds")
	tag = fs.String("tag","","Optional tag name, to tag the session")
    fs.Parse(args)

    if fs.NArg() < 1 {
        fmt.Println("usage: uplog run [--poll N] <command> [args...]")
        os.Exit(1)
    }

    // convert ms → time.Duration
    pollInterval = time.Duration(*poll) * time.Millisecond

    // enforce minimum poll interval (to avoid hammering backend)
    if pollInterval < config.PollIntervalLimit {
        fmt.Printf("Cannot poll below %v ms. Try --poll >= %v.\n",
            config.PollIntervalLimit/time.Millisecond,
            config.PollIntervalLimit/time.Millisecond)
        os.Exit(1)
    }
	command := fs.Arg(0)
    commandArgs := fs.Args()[1:]
	return command, commandArgs, pollInterval, tag
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

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func NullTimeToString(n sql.NullTime) interface{} {
	if n.Valid {
		return n.Time.Format(time.RFC3339)
	}
	return nil
}

func ConstructShareURL(sessionID string) string {
	return config.LogsDomain + sessionID
}

// placeholder: ask backend for session id
// Replace with your real API call to create a session. For now return UUID.
func RequestSessionIDFromBackend(respJSON models.SessionCreateResponse) (string, error) {
	// Example: do POST to backend create session, pass auth header if present, parse returned id
	// For now we do a best-effort call to /create-session; if it fails, return error and caller will fallback to local UUID.
	// Try a quick request to backend (non-fatal)
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("POST", config.BackendSessionCreateEndpoint, nil)
	if err != nil {
		return "", err
	}
	if key := os.Getenv("UPLOG_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// read body as id (assume plain text or JSON id)
		json.NewDecoder(resp.Body).Decode(&respJSON)
		sessionID := respJSON.SessionID
		if sessionID != "" {
			return sessionID, nil
		}
	}
	// fallback: return error so caller can use uuid
	return "", fmt.Errorf("backend returned status %d", resp.StatusCode)
}
