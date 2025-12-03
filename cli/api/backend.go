package api

import (
	"bytes"
	"cli/constants"
	"cli/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func BatchUpload(batch []models.LogEntry, sessionId string, tag string, retries int, backendDisabled *bool) error {
	fmt.Println("Trying to batch upload")
	if len(batch) == 0 {
		return nil
	}
	body, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}
	client := &http.Client{Timeout: 12 * time.Second}

	var lastErr error

	for attempt := 1; attempt <= retries; attempt++ {

		req, err := http.NewRequest(
			"POST",
			constants.BackendUploadEndpoint+"?session_id="+sessionId+"&tag="+tag,
			bytes.NewReader(body),
		)
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		if key := os.Getenv("UPLOG_API_KEY"); key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
		if key := os.Getenv("USER_ID"); key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error",err)
			lastErr = err
			time.Sleep(300 * time.Millisecond) // small backoff
			continue
		}

		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		lastErr = fmt.Errorf(
			"status=%d body=%s",
			resp.StatusCode,
			string(data),
		)

		time.Sleep(300 * time.Millisecond)
	}

	log.Printf("Uplog Backend unreachable after %d retries. Permanently disabling backend.\n", retries)
	*backendDisabled = true

	return lastErr
}

func getUserId(client *http.Client) (string, error) {
	req, _ := http.NewRequest("POST", constants.BackendUserCreateEndpoint, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("create user request failed: %w", err)
	}
	defer resp.Body.Close()

	var userCreateResponse models.UserCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&userCreateResponse); err != nil {
		return "", fmt.Errorf("decode user create response: %w", err)
	}

	return userCreateResponse.UserId, nil
}

func createSession(client *http.Client, apiKey string, userId string) (string, error) {
	req, _ := http.NewRequest("POST", constants.BackendSessionCreateEndpoint, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Id", userId)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("create sessoin request failed: %w", err)
	}
	defer resp.Body.Close()

	var sessionCreateResponse models.SessionCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionCreateResponse); err != nil {
		return "", fmt.Errorf("decode user create response: %w", err)
	}

	return sessionCreateResponse.SessionId, nil
}

// SendConfig returns (apiKey, userID, error)
func SetupSession() (string, error) {
	var err error
	client := &http.Client{Timeout: 3 * time.Second}

	apiKey := os.Getenv("UPLOG_API_KEY")
	userId := os.Getenv("UPLOG_USER_ID")

	if userId == "" {
		var err error
		userId, err = getUserId(client)
		os.Setenv("UPLOG_USER_ID",userId)
		if err != nil {
			return "", err
		}
	}
	sessionId, err := createSession(client, apiKey, userId)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}
