// Package service содержит бизнес-логику клиентского приложения.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// ClientService предоставляет методы для взаимодействия с сервером.
type ClientService struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
	configDir   string
}

// NewClientService создает новый экземпляр ClientService.
func NewClientService() *ClientService {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".vaultfactory")

	// Создаём директорию конфигурации если её нет
	_ = os.MkdirAll(configDir, 0755)

	client := &ClientService{
		baseURL: "http://localhost:8080/api/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		configDir: configDir,
	}

	// Загружаем сохранённый токен
	_ = client.loadToken()

	return client
}

// Register регистрирует нового пользователя на сервере.
func (c *ClientService) Register(ctx context.Context, email, password string) (*models.User, error) {
	req := map[string]string{
		"email":    email,
		"password": password,
	}

	resp, err := c.makeRequest(ctx, "POST", "/auth/register", req)
	if err != nil {
		return nil, err
	}

	var authResp struct {
		User *models.User `json:"user"`
	}

	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return authResp.User, nil
}

// Login выполняет аутентификацию пользователя на сервере.
func (c *ClientService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	req := map[string]string{
		"email":    email,
		"password": password,
	}

	resp, err := c.makeRequest(ctx, "POST", "/auth/login", req)
	if err != nil {
		return nil, "", "", err
	}

	var authResp struct {
		User         *models.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
	}

	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	c.accessToken = authResp.AccessToken
	_ = c.saveToken()
	return authResp.User, authResp.AccessToken, authResp.RefreshToken, nil
}

func (c *ClientService) AddData(ctx context.Context, dataType models.DataType, name, metadata, data string) (*models.DataItem, error) {
	var jsonData json.RawMessage
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %w", err)
	}

	req := map[string]interface{}{
		"type":     dataType,
		"name":     name,
		"metadata": metadata,
		"data":     jsonData,
	}

	resp, err := c.makeAuthenticatedRequest(ctx, "POST", "/data", req)
	if err != nil {
		return nil, err
	}

	var item models.DataItem
	if err := json.Unmarshal(resp, &item); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &item, nil
}

func (c *ClientService) ListData(ctx context.Context) ([]*models.DataItem, error) {
	resp, err := c.makeAuthenticatedRequest(ctx, "GET", "/data", nil)
	if err != nil {
		return nil, err
	}

	var items []*models.DataItem
	if err := json.Unmarshal(resp, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (c *ClientService) GetData(ctx context.Context, id string) (*models.DataItem, error) {
	resp, err := c.makeAuthenticatedRequest(ctx, "GET", fmt.Sprintf("/data/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var item models.DataItem
	if err := json.Unmarshal(resp, &item); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &item, nil
}

func (c *ClientService) DeleteData(ctx context.Context, id string) error {
	_, err := c.makeAuthenticatedRequest(ctx, "DELETE", fmt.Sprintf("/data/%s", id), nil)
	return err
}

func (c *ClientService) Sync(ctx context.Context) error {
	_, err := c.makeAuthenticatedRequest(ctx, "GET", "/data/sync?last_sync=1970-01-01T00:00:00Z", nil)
	return err
}

func (c *ClientService) makeRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *ClientService) makeAuthenticatedRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	if c.accessToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *ClientService) saveToken() error {
	if c.accessToken == "" {
		return nil
	}

	tokenFile := filepath.Join(c.configDir, "token")
	return os.WriteFile(tokenFile, []byte(c.accessToken), 0600)
}

func (c *ClientService) loadToken() error {
	tokenFile := filepath.Join(c.configDir, "token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return err
	}

	c.accessToken = string(data)
	return nil
}

func (c *ClientService) Logout() error {
	c.accessToken = ""
	tokenFile := filepath.Join(c.configDir, "token")
	os.Remove(tokenFile)
	return nil
}
