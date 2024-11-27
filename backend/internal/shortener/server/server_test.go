package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gusevgrishaem1/url-shortener/internal/shortener/config"
	"github.com/gusevgrishaem1/url-shortener/internal/shortener/model"
)

// MockStorage is a mock implementation of the Storage interface.
type MockStorage struct {
	data map[string]string
}

func (m *MockStorage) Save(url model.URL) error {
	m.data[url.Short] = url.Original
	return nil
}

func (m *MockStorage) Get(short string) (original string, exists bool) {
	original, exists = m.data[short]
	return original, exists
}

func TestShortenURL(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	cfg := &config.Config{
		Port:      "8080",
		ServerURL: "http://localhost:8080",
	}
	handler := &Handler{storage: mockStorage, config: cfg}

	// Create a sample URL to shorten
	url := model.URL{Original: "https://www.example.com"}

	// Create a new POST request with the URL to be shortened
	body, _ := json.Marshal(url)
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler.shortenURL(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respURL model.URL
	err = json.NewDecoder(rr.Body).Decode(&respURL)
	assert.NoError(t, err)

	// Check that the short URL was generated and saved correctly
	assert.NotEmpty(t, respURL.Short)
	assert.Equal(t, url.Original, mockStorage.data[respURL.Short])
}

func TestGetOriginalURL(t *testing.T) {
	mockStorage := &MockStorage{data: make(map[string]string)}
	cfg := &config.Config{
		Port:      "8080",
		ServerURL: "http://localhost:8080",
	}
	handler := &Handler{storage: mockStorage, config: cfg}

	// Add a sample short URL to the mock storage
	short := "http://localhost:8080/short123"
	original := "https://www.example.com"
	mockStorage.data[short] = original

	// Create a new GET request to retrieve the original URL
	req, err := http.NewRequest("GET", "/short123", nil)
	assert.NoError(t, err)

	// Record the response
	rr := httptest.NewRecorder()
	handler.getOriginalURL(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusFound, rr.Code)

	// Check the location header
	assert.Equal(t, original, rr.Header().Get("Location"))
}
