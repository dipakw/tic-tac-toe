package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetRequestPayloadAs[T any](r *http.Request) (*T, error) {
	var dest T

	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		return nil, fmt.Errorf("failed to get the payload: %s", err.Error())
	}

	return &dest, nil
}

func HttpRequest[T any](method, url string, payload any) (*T, error) {
	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	var result T

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
