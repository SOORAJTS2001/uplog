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

func BatchUpload(batch []models.LogEntry,userId string, sessionId string, tag string, retries int, backendDisabled *bool) error {
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
		req.Header.Set("User-Id",userId)
		if key := os.Getenv("UPLOG_API_KEY"); key != "" {
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

	resp, err := client.Do(req,)
	if err != nil {
		return "", fmt.Errorf("create session request failed: %w", err)
	}
	defer resp.Body.Close()

	var sessionCreateResponse models.SessionCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionCreateResponse); err != nil {
		return "", fmt.Errorf("decode user create response: %w", err)
	}

	return sessionCreateResponse.SessionId, nil
}

// SendConfig returns (apiKey, userID, error)
func SetupSession() (string,string,string, error) {
	var err error
	client := &http.Client{Timeout: 3 * time.Second}
	fileBytes, err := os.ReadFile(constants.CredentialsFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var config models.Configurations
	config_err:=json.Unmarshal(fileBytes,&config)
	if config_err!=nil{
		// error even after the file has content
		if len(fileBytes) != 0{
			fmt.Println("Error",config_err)
		}
	}
	apiKey := config.ApiKey
	userId := config.UserId

	if userId == "" {
		var err error
		userId, _ = getUserId(client)
		config.UserId = userId
		configBytes,err :=json.Marshal(config)
		if err != nil {
			return "","","", err
		}
		os.WriteFile(constants.CredentialsFile,configBytes,0o700)
	}
	sessionId, err := createSession(client, apiKey, userId)
	sessionUrl:="subject."+userId+"."+sessionId
	fmt.Printf("Uploading Logs at %s%s\n",constants.Domain,sessionUrl)
	time.Sleep(2*time.Second)
	if err != nil {
		return "","","", err
	}

	return apiKey,userId,sessionId, nil
}
