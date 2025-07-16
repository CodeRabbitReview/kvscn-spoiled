//nolint
package recoverer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRecover_RecoverData(t *testing.T) {
	file, err := ioutil.TempFile("", "storage_test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	removeAllSpaces, err := regexp.Compile(`\r|\t|\n| `)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name          string
		action        string
		data          string
		logger        *TransactionLogger
		expectedError error
		expectedOut   string
	}{
		{
			name:          "incorrect action",
			action:        "g",
			data:          `{"key": {"1": 20},"entity": {"misha": 20}}`,
			logger:        NewTransactionLogger(file.Name()),
			expectedError: fmt.Errorf(`incorrect action type: g; want one of this: [p d]`),
			expectedOut:   ``,
		},
		{
			name:          "correct insert",
			action:        "p",
			data:          `{"key":{"1": 20},"entity": {"misha": 20}}`,
			logger:        NewTransactionLogger(file.Name()),
			expectedError: nil,
			expectedOut:   `p{"key":{"1":20},"entity":{"misha":20}}`,
		},
		{
			name:          "correct delete",
			action:        "d",
			data:          `{"key": {"1": 20}`,
			logger:        NewTransactionLogger(file.Name()),
			expectedError: nil,
			expectedOut:   `p{"key":{"1":20},"entity":{"misha":20}}d{"key":{"1":20}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.logger.RecoverData(tt.action, tt.data, DefaultActions)
			if !reflect.DeepEqual(tt.expectedError, err) {
				t.Errorf("expected error: %v; got: %v", tt.expectedError, err)
			}

			readFile, err := ioutil.ReadFile(file.Name())
			if err != nil {
				t.Fatal(err)
			}
			out := removeAllSpaces.ReplaceAllString(string(readFile), "")
			if !reflect.DeepEqual(tt.expectedOut, out) {
				t.Errorf("expected out: %s; got: %s", tt.expectedOut, out)
			}
		})
	}
}

func TestRecover_SendRecovered(t *testing.T) {
	f, err := os.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(`p	{"key":{"54":20},"entity":{"misha":20}}`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("\n")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(`p	{"key":{"47":20},"entity":{"misha":20}}`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("\n")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(`d	{"key": "map[54:20]"}`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("\n")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !reflect.DeepEqual(req.URL.String(), "/api/") {
			t.Errorf("incorrect url %s; want: %s", req.URL.String(), "/api/")
		}
		if !(req.Method != http.MethodDelete) || !(req.Method != http.MethodPut) {
			rw.WriteHeader(http.StatusBadRequest)
		}

		rw.WriteHeader(http.StatusOK)
	}))

	transactionLogger := NewTransactionLogger(f.Name())
	port := strings.Split(server.URL, ":")[2]
	transactionLogger.SendRecovered(":" + port)
}
