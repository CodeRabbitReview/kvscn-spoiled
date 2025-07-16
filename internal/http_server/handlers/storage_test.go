//nolint
package handlers

import (
	"bytes"
	"github.com/mishaprokop4ik/storage/internal/models"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// changeIndexPath only for test usage
func changeIndexPath(p string) {
	indexPath = p
}

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
	s := storage.NewStorage(nil)
	storage := NewStorage(s)

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
	s := storage.NewStorage(nil)
	storage := NewStorage(s)
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
	s := storage.NewStorage(nil)
	e, err := models.NewClearEntity(20, []byte(`{
    "key":"misha",
    "entity": {
		"misha": 20
    }
}
`))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Put(storage.Pair{
		Key:    models.NewKey("misha"),
		Entity: e,
	}); err != nil {
		t.Error(err)
	}

	storageServer := NewStorage(s)
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
	s := storage.NewStorage(nil)
	e, err := models.NewClearEntity(20, []byte(`{
    "key":"misha",
    "entity": {
		"misha": 20
    }}`))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Put(storage.Pair{
		Key:    models.NewKey("misha"),
		Entity: e,
	}); err != nil {
		t.Error(err)
	}

	e, err = models.NewClearEntity(20, []byte(`{
    "key":"dasha",
    "entity": {
		"dasha": 20
    }}`))
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Put(storage.Pair{
		Key:    models.NewKey("dasha"),
		Entity: e,
	}); err != nil {
		t.Error(err)
	}

	storageServer := NewStorage(s)
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
	s := storage.NewStorage(nil)

	storageServer := NewStorage(s)
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
	s := storage.NewStorage(nil)
	storageServer := NewStorage(s)
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
	s := storage.NewStorage(nil)

	e, err := models.NewClearEntity(struct {
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
}`))
	if err != nil {
		t.Fatal(err)
	}

	s.Put(storage.Pair{
		Key:    models.NewKey("developer"),
		Entity: e,
	})
	storageServer := NewStorage(s)
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
	s := storage.NewStorage(nil)

	e, err := models.NewClearEntity(struct {
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
}`))
	if err != nil {
		t.Fatal(err)
	}

	s.Put(storage.Pair{
		Key:    models.NewKey("developer"),
		Entity: e,
	})
	storageServer := NewStorage(s)
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

func TestOutHTMLWithoutData(t *testing.T) {
	s := storage.NewStorage(nil)
	storageServer := NewStorage(s)
	changeIndexPath("static/index.gohtml")
	tests := []struct {
		name        string
		expectedOut string
	}{
		{
			name: "empty input",
			expectedOut: `<!DOCTYPE HTML>
        <html>
        <head>
            <title>Storage</title>
            <style>
                .storage_out {
                    font-family: arial, sans-serif;
                    border-collapse: collapse;
                    width: 100%;
                }
        
                .storage_out-data, .storage_out-header {
                    border: 1px solid #dddddd;
                    text-align: left;
                    padding: 8px;
                }
        
                .storage_out-row:nth-child(even) {
                    background-color: #dddddd;
                }
            </style>
        </head>
        <body>
        <header>
            <h1>
                Key value storage
            </h1>
        </header>
        <div class="main">
            <table class="storage_out">
                
                
                    <div class="empty_storage">
                        The key value storage is empty
                    </div>
                
        
                
            </table>
        </div>
        </body>
        </html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/out", nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedOut, " ", "")) {
				t.Errorf("Test failed %s with rendering expected: %s, got: %s", tt.name,
					tt.expectedOut, w.Body.String())
			}
		})
	}
}

func TestOutHTMLWithOneData(t *testing.T) {
	s := storage.NewStorage(nil)

	e, err := models.NewClearEntity(map[string]interface{}{
		"age":  20,
		"name": "misha",
	}, []byte(`{
    "key": "person",
    "entity": {
        "name": "misha",
        "age": 20
    }
}`))
	if err != nil {
		t.Fatal(err)
	}

	s.Put(storage.Pair{
		Key:    models.NewKey("person"),
		Entity: e,
	})
	storageServer := NewStorage(s)
	changeIndexPath("static/index.gohtml")
	tests := []struct {
		name        string
		expectedOut string
	}{
		{
			name: "one data",
			expectedOut: `<!DOCTYPE HTML>
        <html>
        <head>
            <title>Storage</title>
            <style>
                .storage_out {
                    font-family: arial, sans-serif;
                    border-collapse: collapse;
                    width: 100%;
                }
        
                .storage_out-data, .storage_out-header {
                    border: 1px solid #dddddd;
                    text-align: left;
                    padding: 8px;
                }
        
                .storage_out-row:nth-child(even) {
                    background-color: #dddddd;
                }
            </style>
        </head>
        <body>
        <header>
            <h1>
                Key value storage
            </h1>
        </header>
        <div class="main">
            <table class="storage_out">
                
                
        
                
                    <tr class="storage_out-row">
                        <th class="storage_out-header">
                            Key
                        </th>
                        <th class="storage_out-header">
                            Entity
                        </th>
                    </tr>
                        <tr class="storage_out-row">
                            <td class="storage_out-data">
                                person
                            </td>
                            <td class="storage_out-data">
                                map[age:20 name:misha]
                            </td>
                        </tr>
                    
                
            </table>
        </div>
        </body>
        </html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/out", nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedOut, " ", "")) {
				t.Errorf("Test failed %s with rendering expected: %s, got: %s", tt.name,
					tt.expectedOut, w.Body.String())
			}
		})
	}
}

func TestOutHTMLWithManyData(t *testing.T) {
	s := storage.NewStorage(nil)

	e, err := models.NewClearEntity(map[string][]map[string]string{
		"students": {
			map[string]string{
				"name": "misha",
			},
			map[string]string{
				"name": "dasha",
			},
		},
	}, []byte(`{
    "key": "students",
    "entity": {
        "students": [
            {
                "name": "misha"
            }, 
            {
                "name": "dasha"
            }
        ]
    }
}
}`))
	if err != nil {
		t.Fatal(err)
	}

	s.Put(storage.Pair{
		Key:    models.NewKey("students"),
		Entity: e,
	})

	e, err = models.NewClearEntity(map[string]string{
		"name": "sergei",
	}, []byte(`{
    "key": "teacher",
    "entity": {
        "name": "sergei",
    }
}`))
	if err != nil {
		t.Fatal(err)
	}
	s.Put(storage.Pair{
		Key:    models.NewKey("teacher"),
		Entity: e,
	})
	storageServer := NewStorage(s)
	changeIndexPath("static/index.gohtml")
	tests := []struct {
		name        string
		expectedOut string
	}{
		{
			name: "multiply data from storage",
			expectedOut: `<!DOCTYPE HTML>
        <html>
        <head>
            <title>Storage</title>
            <style>
                .storage_out {
                    font-family: arial, sans-serif;
                    border-collapse: collapse;
                    width: 100%;
                }
        
                .storage_out-data, .storage_out-header {
                    border: 1px solid #dddddd;
                    text-align: left;
                    padding: 8px;
                }
        
                .storage_out-row:nth-child(even) {
                    background-color: #dddddd;
                }
            </style>
        </head>
        <body>
        <header>
            <h1>
                Key value storage
            </h1>
        </header>
        <div class="main">
            <table class="storage_out">
                
                
        
                
                    <tr class="storage_out-row">
                        <th class="storage_out-header">
                            Key
                        </th>
                        <th class="storage_out-header">
                            Entity
                        </th>
                    </tr>
                        <tr class="storage_out-row">
                            <td class="storage_out-data">
                                students
                            </td>
                            <td class="storage_out-data">
                                map[students:[map[name:misha] map[name:dasha]]]
                            </td>
                        </tr>
                    
                        <tr class="storage_out-row">
                            <td class="storage_out-data">
                                teacher
                            </td>
                            <td class="storage_out-data">
                                map[name:sergei]
                            </td>
                        </tr>
                    
                
            </table>
        </div>
        </body>
        </html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/out", nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedOut, " ", "")) {
				t.Errorf("Test failed %s with rendering expected: %s, got: %s", tt.name,
					tt.expectedOut, w.Body.String())
			}
		})
	}
}

func TestOutHTMLEmptyHTMLPath(t *testing.T) {
	s := storage.NewStorage(nil)
	storageServer := NewStorage(s)
	changeIndexPath("")
	tests := []struct {
		name        string
		expectedOut string
	}{
		{
			name:        "empty html path",
			expectedOut: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/out", nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(storageServer.ServeHTTP)
			handler.ServeHTTP(w, req)
			if !reflect.DeepEqual(strings.ReplaceAll(w.Body.String(), " ", ""),
				strings.ReplaceAll(tt.expectedOut, " ", "")) {
				t.Errorf("Test failed %s with rendering expected: %s, got: %s", tt.name,
					tt.expectedOut, w.Body.String())
			}
			if !reflect.DeepEqual(w.Code, http.StatusInternalServerError) {
				t.Errorf("Test failed %s with rendering expected status code: %d, got: %d", tt.name, http.StatusInternalServerError, w.Code)
			}
		})
	}
}
