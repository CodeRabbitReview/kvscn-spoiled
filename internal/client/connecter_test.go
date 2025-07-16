package client

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func BenchmarkPutConcurrently(b *testing.B) {
	var err error
	param := `{"key":"user1","entity": {"misha": 20}}`
	expectedResult := []byte(`[{"key":"user1","entity":{"misha":20}}]`)
	c := NewAPI("http://localhost:8080")
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
	c := NewAPI("http://localhost:8080")
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
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !reflect.DeepEqual(req.URL.String(), "/api/") {
			t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
		}
		if !reflect.DeepEqual(req.Method, http.MethodGet) {
			t.Error("incorrect method")
		}
		_, _ = rw.Write([]byte(`{"key":"person","entity": {"misha": 20}}`))
		rw.WriteHeader(http.StatusOK)
	}))
	expectedOut := []byte(`{"key":"person","entity": {"misha": 20}}`)
	expectedStatusCode := 200
	c := NewAPI(server.URL)

	out, err := c.GetAll()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(strings.ReplaceAll(string(out.Body), " ", ""),
		strings.ReplaceAll(string(expectedOut), " ", "")) {
		t.Errorf("Test failed get all client expected body: %v, got: %v", expectedOut, string(out.Body))
	}

	if !reflect.DeepEqual(expectedStatusCode, out.StatusCode) {
		t.Errorf("Test failed get all client expected code: %v, got: %v", expectedStatusCode, out.StatusCode)
	}
}

func TestGetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !reflect.DeepEqual(req.URL.String(), "/api/id") {
			t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
		}
		if !reflect.DeepEqual(req.Method, http.MethodGet) {
			t.Error("incorrect method")
		}
		_, _ = rw.Write([]byte(`{"key":"person","entity": {"misha": 20}}`))
		rw.WriteHeader(http.StatusOK)
	}))
	expectedOut := []byte(`{"key":"person","entity": {"misha": 20}}`)
	expectedStatusCode := 200
	c := NewAPI(server.URL)

	out, err := c.GetByID(`{"key":"person"}`)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(strings.ReplaceAll(string(out.Body), " ", ""),
		strings.ReplaceAll(string(expectedOut), " ", "")) {
		t.Errorf("Test failed get all client expected body: %v, got: %v", expectedOut, string(out.Body))
	}

	if !reflect.DeepEqual(expectedStatusCode, out.StatusCode) {
		t.Errorf("Test failed get all client expected code: %v, got: %v", expectedStatusCode, out.StatusCode)
	}
}

func TestAddOrUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !reflect.DeepEqual(req.URL.String(), "/api/") {
			t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
		}
		if !reflect.DeepEqual(req.Method, http.MethodPost) {
			t.Error("incorrect method")
		}
		rw.WriteHeader(http.StatusCreated)
	}))
	expectedStatusCode := http.StatusCreated
	c := NewAPI(server.URL)

	out, err := c.AddOrUpdate(`{"key":"person","entity": {"misha": 20}}`)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedStatusCode, out.StatusCode) {
		t.Errorf("Test failed get all client expected code: %v, got: %v", expectedStatusCode, out.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !reflect.DeepEqual(req.URL.String(), "/api/") {
			t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
		}
		if !reflect.DeepEqual(req.Method, http.MethodDelete) {
			t.Error("incorrect method")
		}
		rw.WriteHeader(http.StatusNoContent)
	}))
	expectedStatusCode := http.StatusNoContent
	c := NewAPI(server.URL)

	out, err := c.Delete(`{"key":"person"}`)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedStatusCode, out.StatusCode) {
		t.Errorf("Test failed get all client expected code: %v, got: %v", expectedStatusCode, out.StatusCode)
	}
}
