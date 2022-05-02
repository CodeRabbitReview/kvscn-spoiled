//nolint
package handlers

import (
	"bytes"
	"github.com/mishaprokop4ik/storage/internal/models"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		body           io.Reader
		method         string
		url            string
		expectedBody   string
		expectedStatus int
	}{
		{
			name:           "more than one slash in URL",
			body:           nil,
			method:         http.MethodDelete,
			url:            "/api/blabla/",
			expectedBody:   "",
			expectedStatus: http.StatusNotAcceptable,
		},
		{
			name:           "get all data from empty storage",
			body:           nil,
			method:         http.MethodGet,
			url:            "/api/",
			expectedBody:   `{"response":"no data in storage"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "incorrect URL",
			body:           nil,
			method:         http.MethodGet,
			url:            "/api",
			expectedBody:   "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "without /api prefix in url",
			body:           nil,
			method:         http.MethodGet,
			url:            "/apinn",
			expectedBody:   "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "add new data into storage",
			body: bytes.NewBuffer([]byte(`{
					"key": "key value",
					"entity": {
						"name": "misha",
						"age": 20
					}
			}`)),
			method:         http.MethodPut,
			url:            "/api/",
			expectedBody:   "",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "delete value with key that is not in storage",
			body:           bytes.NewBuffer([]byte(`{"key": "misha_prokopchyk"}`)),
			method:         http.MethodDelete,
			url:            "/api/",
			expectedBody:   `{"response":"no such value by this key"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "get value by id from storage",
			body:           bytes.NewBuffer([]byte(`{"key":"key_value"}`)),
			method:         http.MethodGet,
			url:            "/api/id",
			expectedBody:   `{"response":{"key":"key value","entity":{"name":"misha","age":20}}}`,
			expectedStatus: http.StatusOK,
		},
		{
			name: "add new data into storage with a object key",
			body: bytes.NewBuffer([]byte(`{
					"key": {
						"country": "Ukraine",
						"city": "Kharkov"
					},
					"entity": {
						"name": "misha",
						"age": 20
					}
			}`)),
			method:         http.MethodPut,
			url:            "/api/",
			expectedBody:   "",
			expectedStatus: http.StatusCreated,
		},
		{
			name: "add new data into storage with a incorrect key",
			body: bytes.NewBuffer([]byte(`{
					"key": "{}",
					"entity": {
						"name": "misha",
						"age": 20
					}
			}`)),
			method:         http.MethodPut,
			url:            "/api/",
			expectedBody:   `{"response":"empty key value"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "get all data from non empty storage",
			body:           nil,
			method:         http.MethodGet,
			url:            "/api/",
			expectedBody:   `[{"key":"key value","entity":{"name":"misha","age":20}},{"key":{"country":"Ukraine","city":"Kharkov"},"entity":{"name":"misha","age":20}}]`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get data by key that does not exist",
			body:           bytes.NewBuffer([]byte(`{"key":"123123"}`)),
			method:         http.MethodGet,
			url:            "/api/id",
			expectedBody:   `{"response":"no such value by this key"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "delete data",
			body:           bytes.NewBuffer([]byte(`{"key": "key_value"}`)),
			method:         http.MethodDelete,
			url:            "/api/",
			expectedBody:   "",
			expectedStatus: http.StatusNoContent,
		},
	}
	storage := NewStorage(nil, storage.NewStorage())

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.url, tt.body)
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(storage.ServeHTTP)
		handler.ServeHTTP(w, req)
		if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
			strings.ReplaceAll(tt.expectedBody, " ", "")) {
			t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
		}

		if !reflect.DeepEqual(w.Code, tt.expectedStatus) {
			t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatus, w.Code)
		}
	}
}

func TestGetAllInEmptyStorage(t *testing.T) {
	storage := NewStorage(nil, storage.NewStorage())
	tests := []struct {
		name               string
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:               "empty storage",
			expectedBody:       `{"response":"no data in storage"}`,
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/", nil)
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(storage.ServeHTTP)
		handler.ServeHTTP(w, req)
		if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
			strings.ReplaceAll(tt.expectedBody, " ", "")) {
			t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
		}

		if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
			t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
		}
	}
}

func TestGetAllWithOneObjectInStorage(t *testing.T) {
	s := storage.NewStorage()

	if err := s.Put(storage.Pair{
		Key: models.NewKey("misha"),
		Entity: models.NewEntity(20, []byte(`{
    "key":"misha",
    "entity": {
		"misha": 20
    }
}
`)),
	}); err != nil {
		t.Error(err)
	}

	storageServer := NewStorage(nil, s)
	tests := []struct {
		name               string
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:               "with one object in storageServer",
			expectedBody:       `[{"key":"misha","entity":{"misha":20}}]`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/", nil)
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(storageServer.ServeHTTP)
		handler.ServeHTTP(w, req)
		if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
			strings.ReplaceAll(tt.expectedBody, " ", "")) {
			t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
		}

		if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
			t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
		}
	}
}

func TestGetAllWithManyObjectsInStorage(t *testing.T) {
	s := storage.NewStorage()

	if err := s.Put(storage.Pair{
		Key: models.NewKey("misha"),
		Entity: models.NewEntity(20, []byte(`{
    "key":"misha",
    "entity": {
		"misha": 20
    }}`)),
	}); err != nil {
		t.Error(err)
	}
	if err := s.Put(storage.Pair{
		Key: models.NewKey("dasha"),
		Entity: models.NewEntity(20, []byte(`{
    "key":"dasha",
    "entity": {
		"dasha": 20
    }}`)),
	}); err != nil {
		t.Error(err)
	}

	storageServer := NewStorage(nil, s)
	tests := []struct {
		name               string
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:               "with more than one object in storageServer",
			expectedBody:       `[{"key":"dasha","entity":{"dasha":20}},{"key":"misha","entity":{"misha":20}}]`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/", nil)
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(storageServer.ServeHTTP)
		handler.ServeHTTP(w, req)
		if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
			strings.ReplaceAll(tt.expectedBody, " ", "")) {
			t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
		}

		if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
			t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
		}
	}
}

func TestPut(t *testing.T) {
	s := storage.NewStorage()
	storageServer := NewStorage(nil, s)
	tests := []struct {
		name               string
		input              *bytes.Buffer
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "simple input",
			input: bytes.NewBuffer([]byte(`{
    "key": "developer",
    "entity": {
        "name": "misha",
        "age": 20
    }
}`)),
			expectedBody:       "",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "input with empty entity",
			input: bytes.NewBuffer([]byte(`{
    "key": "developer"
}`)),
			expectedBody:       `{"response":"nil in input data"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "input with empty key",
			input: bytes.NewBuffer([]byte(`{
    "entity": {
        "name": "misha",
        "age": 20
    }
}`)),
			expectedBody:       `{"response":"empty key value"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/", tt.input)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedBody, " ", "")) {
				t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
			}

			if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
				t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
			}
		})
	}
}

func TestGetFromEmptyStorage(t *testing.T) {
	s := storage.NewStorage()
	l := &log.Logger{}
	storageServer := NewStorage(l, s)
	tests := []struct {
		name               string
		expectedStatusCode int
		key                string
		expectedBody       string
	}{
		{
			name:               "get by id from empty storage",
			expectedStatusCode: http.StatusNotFound,
			key:                `{"key": "value"}`,
			expectedBody:       `{"response":"no data in storage"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/id",
				bytes.NewBuffer([]byte(tt.key)))
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedBody, " ", "")) {
				t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
			}

			if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
				t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
			}
		})
	}
}

func TestGet(t *testing.T) {
	s := storage.NewStorage()
	s.Put(storage.Pair{
		Key: models.NewKey("developer"),
		Entity: models.NewEntity(struct {
			name string
			lvl  string
		}{
			name: "misha",
			lvl:  "trainee",
		}, []byte(`{
    "key": "developer",
    "entity": {
        "name": "misha",
        "lvl": "trainee"
    }
}`)),
	})
	storageServer := NewStorage(nil, s)
	tests := []struct {
		name               string
		expectedStatusCode int
		key                string
		expectedBody       string
	}{
		{
			name:               "get by id",
			expectedStatusCode: http.StatusOK,
			key:                `{"key": "developer"}`,
			expectedBody:       `{"response":{"key":"developer","entity":{"name":"misha","lvl":"trainee"}}}`,
		},
		{
			name:               "get by id incorrect id",
			expectedStatusCode: http.StatusNotFound,
			key:                `{"key": "key_value"}`,
			expectedBody:       `{"response":"no such value by this key"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/id",
				bytes.NewBuffer([]byte(tt.key)))
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedBody, " ", "")) {
				t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
			}

			if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
				t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	s := storage.NewStorage()
	s.Put(storage.Pair{
		Key: models.NewKey("developer"),
		Entity: models.NewEntity(struct {
			name string
			lvl  string
		}{
			name: "misha",
			lvl:  "trainee",
		}, []byte(`{
    "key": "developer",
    "entity": {
        "name": "misha",
        "lvl": "trainee"
    }
}`)),
	})
	storageServer := NewStorage(nil, s)
	tests := []struct {
		name               string
		expectedStatusCode int
		key                string
		expectedBody       string
	}{
		{
			name:               "delete by that does not exist",
			expectedStatusCode: http.StatusNotFound,
			key:                `{"key": "misha"}`,
			expectedBody:       `{"response":"no such value by this key"}`,
		},
		{
			name:               "simple delete",
			expectedStatusCode: http.StatusNoContent,
			key:                `{"key": "developer"}`,
			expectedBody:       ``,
		},
		{
			name:               "from empty storage",
			expectedStatusCode: http.StatusNotFound,
			key:                `{"key": "user"}`,
			expectedBody:       `{"response":"no data in storage"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/",
				bytes.NewBuffer([]byte(tt.key)))
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedBody, " ", "")) {
				t.Errorf("Test failed %s expected body: %v, got: %v", tt.name, tt.expectedBody, w.Body.String())
			}

			if !reflect.DeepEqual(w.Code, tt.expectedStatusCode) {
				t.Errorf("Test failed %s expected code: %v, got: %v", tt.name, tt.expectedStatusCode, w.Code)
			}
		})
	}
}
