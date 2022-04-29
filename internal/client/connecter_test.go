//nolint
package client

import (
	"fmt"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func init() {
	zlog.Init("stderr")
}

func BenchmarkPutConcurrently(b *testing.B) {
	var err error
	param := `{"key":"user1","entity": {"misha": 20}}`
	expectedResult := []byte(`[{"key":"user1","entity":{"misha":20}}]`)
	c := NewAPI("https://localhost:8080", "./../../localhost.pem")
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err = c.AddOrUpdate(param)
			if err != nil {
				b.Error(err)
			}
		}()
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()
	b.StopTimer()
	resp, err := c.GetAll()
	if err != nil {
		b.Error(err)
	}

	if !reflect.DeepEqual(resp.Body, expectedResult) {
		b.Fatalf("expected: %s; got: %s", expectedResult, resp.Body)
	}
}

func BenchmarkPutSequentially(b *testing.B) {
	var err error
	param := `{"key":"user1","entity": {"misha": 20}}`
	expectedResult := []byte(`[{"key":"user1","entity":{"misha":20}}]`)
	c := NewAPI("https://localhost:8080", "./../../localhost.pem")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = c.AddOrUpdate(param)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	resp, err := c.GetAll()
	if err != nil {
		b.Error(err)
	}

	if !reflect.DeepEqual(resp.Body, expectedResult) {
		b.Fatalf("expected: %s; got: %s", expectedResult, resp.Body)
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name               string
		expectedOut        []byte
		expectedStatusCode int
		certPath           string
		server             *httptest.Server
	}{
		{
			name:               "with correct cert path",
			expectedOut:        []byte(`{"key":"person","entity": {"misha": 20}}`),
			expectedStatusCode: http.StatusOK,
			certPath:           "./../../localhost.pem",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodGet) {
					t.Error("incorrect method")
				}
				_, _ = rw.Write([]byte(`{"key":"person","entity": {"misha": 20}}`))
				rw.WriteHeader(http.StatusOK)
			})),
		},
		{
			name:               "with incorrect cert path",
			expectedOut:        []byte(`{"key":"person","entity": {"misha": 20}}`),
			expectedStatusCode: http.StatusOK,
			certPath:           "source",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodGet) {
					t.Error("incorrect method")
				}
				_, _ = rw.Write([]byte(`{"key":"person","entity": {"misha": 20}}`))
				rw.WriteHeader(http.StatusOK)
			})),
		},
		{
			name:               "not http.StatusOK server response",
			expectedOut:        []byte{},
			expectedStatusCode: http.StatusInternalServerError,
			certPath:           "./../localhost.pem",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodGet) {
					t.Error("incorrect method")
				}
				rw.WriteHeader(http.StatusInternalServerError)
			})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewAPI(tt.server.URL, tt.certPath)

			out, err := c.GetAll()
			if err != nil && err.Error() != "incorrect status code want: 200; get: 500" {
				t.Error(err)
			}

			if !reflect.DeepEqual(strings.ReplaceAll(string(out.Body), " ", ""),
				strings.ReplaceAll(string(tt.expectedOut), " ", "")) {
				t.Errorf("Test failed get all client expected body: %v, got: %v", tt.expectedOut, string(out.Body))
			}

			if !reflect.DeepEqual(tt.expectedStatusCode, out.StatusCode) {
				t.Errorf("Test failed get all client expected code: %v, got: %v", tt.expectedStatusCode, out.StatusCode)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name               string
		server             *httptest.Server
		expectedOut        []byte
		expectedStatusCode int
		expectedErr        error
		key                string
	}{
		{
			name: "simple get",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/id") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodGet) {
					t.Error("incorrect method")
				}
				_, _ = rw.Write([]byte(`{"key":"person","entity": {"misha": 20}}`))
				rw.WriteHeader(http.StatusOK)
			})),
			expectedOut:        []byte(`{"key":"person","entity": {"misha": 20}}`),
			expectedStatusCode: http.StatusOK,
			key:                `{"key":"person"}`,
		},
		{
			name: "not correct http status",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/id") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodGet) {
					t.Error("incorrect method")
				}
				rw.WriteHeader(http.StatusInternalServerError)
			})),
			expectedOut:        []byte{},
			expectedStatusCode: http.StatusInternalServerError,
			key:                `{"key":"person2"}`,
			expectedErr:        fmt.Errorf("incorrect status code want: 200; get: 500"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewAPI(tt.server.URL, "./../../localhost.pem")

			out, err := c.GetByID(tt.key)
			if err != nil &&
				err.Error() != tt.expectedErr.Error() {
				t.Error(err)
			}

			if !reflect.DeepEqual(strings.ReplaceAll(string(out.Body), " ", ""),
				strings.ReplaceAll(string(tt.expectedOut), " ", "")) {
				t.Errorf("Test failed get by id client expected body: %v, got: %v", tt.expectedOut, string(out.Body))
			}

			if !reflect.DeepEqual(tt.expectedStatusCode, out.StatusCode) {
				t.Errorf("Test failed get by id client expected code: %v, got: %v", tt.expectedStatusCode, out.StatusCode)
			}
		})
	}
}

func TestAddOrUpdate(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		certPath           string
		param              string
		server             *httptest.Server
	}{
		{
			name:               "simple insertion",
			expectedStatusCode: http.StatusCreated,
			certPath:           "./../../localhost.pem",
			param:              `{"key":"person","entity": {"misha": 20}}`,
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodPost) {
					t.Error("incorrect method")
				}
				rw.WriteHeader(http.StatusCreated)
			})),
		},
		{
			name:               "simple insertion with incorrect response status code",
			expectedStatusCode: http.StatusBadGateway,
			certPath:           "./../../localhost.pem",
			param:              `{"key":"person","entity": {"misha": 20}}`,
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadGateway)
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodPost) {
					t.Error("incorrect method")
				}
			})),
		},
	}
	for _, tt := range tests {
		c := NewAPI(tt.server.URL, tt.certPath)

		out, err := c.AddOrUpdate(`{"key":"person","entity": {"misha": 20}}`)
		if err != nil && err.Error() != "incorrect status code want: 201; get: 502" {
			t.Error(err)
		}

		if !reflect.DeepEqual(tt.expectedStatusCode, out.StatusCode) {
			t.Errorf("Test failed put client expected code: %v, got: %v", tt.expectedStatusCode, out.StatusCode)
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		expectedErr        error
		key                string
		server             *httptest.Server
	}{
		{
			name:               "simple delete",
			expectedStatusCode: http.StatusNoContent,
			expectedErr:        nil,
			key:                `{"key":"person"}`,
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodDelete) {
					t.Error("incorrect method")
				}
				rw.WriteHeader(http.StatusNoContent)
			})),
		},
		{
			name:               "server not http.StatusNoContent",
			expectedStatusCode: http.StatusCreated,
			expectedErr: fmt.Errorf("incorrect status code want: %d; get: %d",
				http.StatusNoContent, http.StatusCreated),
			key: `{"key":"person"}`,
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if !reflect.DeepEqual(req.URL.String(), "/api/") {
					t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
				}
				if !reflect.DeepEqual(req.Method, http.MethodDelete) {
					t.Error("incorrect method")
				}
				rw.WriteHeader(http.StatusCreated)
			})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewAPI(tt.server.URL, "./../../localhost.pem")
			out, err := c.Delete(tt.key)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.expectedStatusCode, out.StatusCode) {
				t.Errorf("Test failed delete client expected code: %v, got: %v", tt.expectedStatusCode, out.StatusCode)
			}
		})
	}
}
