package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/config"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(serverURL, token string) *Client {
	if serverURL == "" {
		if envURL := os.Getenv("ORBIT_SERVER"); envURL != "" {
			serverURL = envURL
		} else {
			serverURL = "http://127.0.0.1:8787"
		}
	}
	serverURL = strings.TrimRight(serverURL, "/")

	if token == "" {
		if envToken := os.Getenv("ORBIT_TOKEN"); envToken != "" {
			token = envToken
		} else {
			// Auto-discover the admin token on this node. The master key is never used as a token.
			root := config.DefaultRootPath()
			tokenPath := filepath.Join(root, "secrets", "admin.token")
			if b, err := os.ReadFile(tokenPath); err == nil {
				token = strings.TrimSpace(string(b))
			}
		}
	}

	return &Client{
		BaseURL: serverURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Do(method, path string, body io.Reader, contentType string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", url, err)
	}

	return resp, nil
}

func (c *Client) PostJSON(path string, payload any, target any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Do("POST", path, bytes.NewReader(data), "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}
	return nil
}

func (c *Client) GetJSON(path string, target any) error {
	resp, err := c.Do("GET", path, nil, "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

var _ = errors.New
