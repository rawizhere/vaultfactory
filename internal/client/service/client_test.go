package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

func TestClientService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v1/auth/register", r.URL.Path)

			var req map[string]string
			_ = json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "test@example.com", req["email"])
			assert.Equal(t, "password123", req["password"])

			response := map[string]interface{}{
				"user": map[string]interface{}{
					"id":    uuid.New().String(),
					"email": "test@example.com",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &ClientService{
			baseURL:    server.URL + "/api/v1",
			httpClient: &http.Client{},
		}

		user, err := client.Register(context.Background(), "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("User already exists"))
		}))
		defer server.Close()

		client := &ClientService{
			baseURL:    server.URL + "/api/v1",
			httpClient: &http.Client{},
		}

		user, err := client.Register(context.Background(), "test@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "request failed with status 400")
	})
}

func TestClientService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v1/auth/login", r.URL.Path)

			var req map[string]string
			_ = json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "test@example.com", req["email"])
			assert.Equal(t, "password123", req["password"])

			response := map[string]interface{}{
				"user": map[string]interface{}{
					"id":    uuid.New().String(),
					"email": "test@example.com",
				},
				"access_token":  "access-token-123",
				"refresh_token": "refresh-token-123",
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &ClientService{
			baseURL:    server.URL + "/api/v1",
			httpClient: &http.Client{},
		}

		user, accessToken, refreshToken, err := client.Login(context.Background(), "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "access-token-123", accessToken)
		assert.Equal(t, "refresh-token-123", refreshToken)
		assert.Equal(t, "access-token-123", client.accessToken)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid credentials"))
		}))
		defer server.Close()

		client := &ClientService{
			baseURL:    server.URL + "/api/v1",
			httpClient: &http.Client{},
		}

		user, accessToken, refreshToken, err := client.Login(context.Background(), "test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Contains(t, err.Error(), "request failed with status 401")
	})
}

func TestClientService_AddData(t *testing.T) {
	t.Run("successful data addition", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v1/data", r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

			var req map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, string(models.LoginPassword), req["type"])
			assert.Equal(t, "test-password", req["name"])
			assert.Equal(t, "test-metadata", req["metadata"])

			response := models.DataItem{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Type:     models.LoginPassword,
				Name:     "test-password",
				Metadata: "test-metadata",
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &ClientService{
			baseURL:     server.URL + "/api/v1",
			accessToken: "test-token",
			httpClient:  &http.Client{},
		}

		validJSON := `{"username": "test", "password": "secret"}`
		item, err := client.AddData(context.Background(), models.LoginPassword, "test-password", "test-metadata", validJSON)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, models.LoginPassword, item.Type)
		assert.Equal(t, "test-password", item.Name)
	})

	t.Run("invalid JSON data", func(t *testing.T) {
		client := &ClientService{
			baseURL:     "http://localhost:8080/api/v1",
			accessToken: "test-token",
			httpClient:  &http.Client{},
		}

		invalidJSON := `{"username": "test", "password": "secret"`
		item, err := client.AddData(context.Background(), models.LoginPassword, "test-password", "test-metadata", invalidJSON)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "invalid JSON data")
	})

	t.Run("not authenticated", func(t *testing.T) {
		client := &ClientService{
			baseURL:     "http://localhost:8080/api/v1",
			accessToken: "",
			httpClient:  &http.Client{},
		}

		validJSON := `{"username": "test", "password": "secret"}`
		item, err := client.AddData(context.Background(), models.LoginPassword, "test-password", "test-metadata", validJSON)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "not authenticated")
	})
}

func TestClientService_ListData(t *testing.T) {
	t.Run("successful data listing", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/data", r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

			response := []models.DataItem{
				{
					ID:       uuid.New(),
					UserID:   uuid.New(),
					Type:     models.LoginPassword,
					Name:     "test-password",
					Metadata: "test-metadata",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &ClientService{
			baseURL:     server.URL + "/api/v1",
			accessToken: "test-token",
			httpClient:  &http.Client{},
		}

		items, err := client.ListData(context.Background())

		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, models.LoginPassword, items[0].Type)
		assert.Equal(t, "test-password", items[0].Name)
	})

	t.Run("not authenticated", func(t *testing.T) {
		client := &ClientService{
			baseURL:     "http://localhost:8080/api/v1",
			accessToken: "",
			httpClient:  &http.Client{},
		}

		items, err := client.ListData(context.Background())

		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "not authenticated")
	})
}

func TestClientService_TokenManagement(t *testing.T) {
	t.Run("save and load token", func(t *testing.T) {
		// Создаем временную директорию для тестов
		tempDir := t.TempDir()

		client := &ClientService{
			baseURL:     "http://localhost:8080/api/v1",
			accessToken: "test-token",
			configDir:   tempDir,
		}

		// Сохраняем токен
		err := client.saveToken()
		assert.NoError(t, err)

		// Проверяем, что файл создан
		tokenFile := filepath.Join(tempDir, "token")
		_, err = os.Stat(tokenFile)
		assert.NoError(t, err)

		// Создаем новый клиент и загружаем токен
		newClient := &ClientService{
			baseURL:   "http://localhost:8080/api/v1",
			configDir: tempDir,
		}

		err = newClient.loadToken()
		assert.NoError(t, err)
		assert.Equal(t, "test-token", newClient.accessToken)
	})

	t.Run("logout clears token", func(t *testing.T) {
		tempDir := t.TempDir()

		client := &ClientService{
			baseURL:     "http://localhost:8080/api/v1",
			accessToken: "test-token",
			configDir:   tempDir,
		}

		// Сохраняем токен
		_ = client.saveToken()

		// Выходим
		err := client.Logout()
		assert.NoError(t, err)
		assert.Empty(t, client.accessToken)

		// Проверяем, что файл удален
		tokenFile := filepath.Join(tempDir, "token")
		_, err = os.Stat(tokenFile)
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})
}
