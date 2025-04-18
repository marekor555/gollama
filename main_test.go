package gollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "data"}`))
	}))
	defer server.Close()

	// Test successful request
	data, err := getRequest(server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(data) != `{"test": "data"}` {
		t.Errorf("Expected {\"test\": \"data\"}, got %s", string(data))
	}

	// Test invalid URL
	_, err = getRequest("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL, got none")
	}
}

func TestListModels(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("Expected path /api/tags, got %s", r.URL.Path)
		}

		mockResponse := ModelList{
			Models: []ModelInfo{
				{Name: "model1"},
				{Name: "model2"},
			},
		}
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	model := &Model{
		Name: "test-model",
		Addr: server.URL,
	}

	models, err := model.ListModels()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedModels := []string{"model1", "model2"}
	if len(models) != len(expectedModels) {
		t.Errorf("Expected %d models, got %d", len(expectedModels), len(models))
	}

	for i, m := range models {
		if m != expectedModels[i] {
			t.Errorf("Expected model %s, got %s", expectedModels[i], m)
		}
	}
}

func TestCreateModel(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockResponse := ModelList{
			Models: []ModelInfo{
				{Name: "model1"},
				{Name: "test-model"},
			},
		}
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	tests := []struct {
		name        string
		modelName   string
		addr        string
		expectError bool
	}{
		{
			name:        "Valid model",
			modelName:   "test-model",
			addr:        server.URL,
			expectError: false,
		},
		{
			name:        "Invalid model name",
			modelName:   "non-existent-model",
			addr:        server.URL,
			expectError: true,
		},
		{
			name:        "Empty address",
			modelName:   "test-model",
			addr:        "",
			expectError: true, // Will fail because localhost won't be running
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := CreateModel(tt.modelName, tt.addr)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if model == nil {
					t.Error("Expected model instance, got nil")
				}
				if model != nil && model.Name != tt.modelName {
					t.Errorf("Expected model name %s, got %s", tt.modelName, model.Name)
				}
			}
		})
	}
}
