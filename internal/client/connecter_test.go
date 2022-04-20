package client

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func BenchmarkPutConcurrently(b *testing.B) {
	var err error
	param := `{"key":"user1","entity": {"misha": 20}}`
	expectedResult := []byte(`[{"key":"user1","entity":{"misha":20}}]`)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err = AddOrUpdate(param)
			if err != nil {
				b.Fatal(err)
			}
		}()
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()

	resp, err := GetAll()
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
	for i := 0; i < b.N; i++ {
		_, err = AddOrUpdate(param)
		if err != nil {
			b.Fatal(err)
		}
	}

	resp, err := GetAll()
	if err != nil {
		b.Error(err)
	}

	if !reflect.DeepEqual(resp.Body, expectedResult) {
		b.Fatalf("expected: %s; got: %s", expectedResult, resp.Body)
	}
}

func TestGetAll(t *testing.T) {

}

func TestGetByID(t *testing.T) {

}

func TestAddOrUpdate(t *testing.T) {

}

func TestDelete(t *testing.T) {

}
