//nolint
package storage

import (
	"reflect"
	"testing"
)

func TestStorage_Put(t *testing.T) {
	tests := []struct {
		name            string
		input           Pair
		expectedStorage Storage
		gotStorage      Storage
		expectedError   error
	}{
		{
			"nil key",
			Pair{
				Key:   "nil",
				Value: nil,
			},
			*NewStorage(),
			*NewStorage(),
			ErrNilInput,
		},
		{
			"empty string key",
			Pair{
				Key:   "empty string value",
				Value: "",
			},
			Storage{
				pairs: map[Key]interface{}{
					"empty string value": "",
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"simple string key",
			Pair{
				Key:   "simple string",
				Value: "simple",
			},
			Storage{
				pairs: map[Key]interface{}{
					"simple string": "simple",
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"number key",
			Pair{
				Key:   "number",
				Value: 5,
			},
			Storage{
				pairs: map[Key]interface{}{
					"number": 5,
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"slice",
			Pair{
				Key:   "slice",
				Value: []int{1, 2, 3},
			},
			Storage{
				pairs: map[Key]interface{}{
					"slice": []int{1, 2, 3},
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"struct key",
			Pair{
				Key: "struct",
				Value: struct {
					name string
					age  int
				}{
					"misha",
					20,
				},
			},
			Storage{
				pairs: map[Key]interface{}{
					"struct": struct {
						name string
						age  int
					}{
						"misha",
						20,
					},
				},
			},
			*NewStorage(),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.gotStorage.Put(tt.input)
			if !reflect.DeepEqual(gotErr, tt.expectedError) {
				t.Fatal(gotErr)
			}

			if !reflect.DeepEqual(tt.gotStorage, tt.expectedStorage) {
				t.Fatalf("expected storage: %v, got %v", tt.expectedStorage, tt.gotStorage)
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	tests := []struct {
		name          string
		key           Key
		expectedValue interface{}
		expectedError error
		storage       Storage
	}{
		{
			"string return",
			"string data",
			"data",
			nil,
			Storage{
				pairs: map[Key]interface{}{
					"string data": "data",
				},
			},
		},
		{
			"int return",
			"number",
			5,
			nil,
			Storage{
				pairs: map[Key]interface{}{
					"number": 5,
				},
			},
		},
		{
			"slice return",
			"slice",
			[]int{1, 2, 3},
			nil,
			Storage{
				pairs: map[Key]interface{}{
					"slice": []int{1, 2, 3},
				},
			},
		},
		{
			"struct return ",
			"struct",
			struct {
				name string
				age  int
			}{
				"misha",
				20,
			},
			nil,
			Storage{
				pairs: map[Key]interface{}{
					"struct": struct {
						name string
						age  int
					}{
						"misha",
						20,
					},
				},
			},
		},
		{
			"no such key return",
			"abab",
			nil,
			ErrNoSuchKey,
			Storage{
				pairs: make(map[Key]interface{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.storage.Get(tt.key)

			if !reflect.DeepEqual(err, tt.expectedError) {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tt.expectedValue) {
				t.Fatalf("expected value: %v, got %v", tt.expectedValue, got)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	tests := []struct {
		name          string
		key           Key
		storage       Storage
		expectedError error
	}{
		{
			"simple delete",
			"key",
			Storage{
				pairs: map[Key]interface{}{
					"key": "value",
				},
			},
			ErrNoSuchKey,
		},
		{
			"delete with no such key",
			"abab",
			Storage{
				pairs: map[Key]interface{}{
					"key": "value",
				},
			},
			ErrNoSuchKey,
		},
		{
			"empty key",
			"",
			Storage{
				pairs: make(map[Key]interface{}),
			},
			ErrNoSuchKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.Delete(tt.key)

			_, err := tt.storage.Get(tt.key)
			if !reflect.DeepEqual(err, tt.expectedError) {
				t.Fatal(err)
			}
		})
	}
}
