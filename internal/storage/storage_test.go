//nolint
package storage

import (
	"github.com/mishaprokop4ik/storage/internal/storage/models"
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
				Key:    models.NewKey(nil),
				Entity: models.NewEntity(nil),
			},
			*NewStorage(),
			*NewStorage(),
			models.ErrNilInput,
		},
		{
			"empty key value",
			Pair{
				Key:    models.NewKey(""),
				Entity: models.NewEntity("empty string value"),
			},
			Storage{
				pairs: map[Value]Value{},
			},
			*NewStorage(),
			models.ErrEmptyKeyString,
		},
		{
			"simple string key",
			Pair{
				Key:    models.NewKey("simple string"),
				Entity: models.NewEntity("simple"),
			},
			Storage{
				pairs: map[Value]Value{
					models.NewKey("simple string"): models.NewEntity("simple"),
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"number key",
			Pair{
				Key:    models.NewKey("number"),
				Entity: models.NewEntity(5),
			},
			Storage{
				pairs: map[Value]Value{
					models.NewKey("number"): models.NewEntity(5),
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"slice",
			Pair{
				Key:    models.NewKey("slice"),
				Entity: models.NewEntity([]int{1, 2, 3}),
			},
			Storage{
				pairs: map[Value]Value{
					models.NewKey("slice"): models.NewEntity([]int{1, 2, 3}),
				},
			},
			*NewStorage(),
			nil,
		},
		{
			"struct key",
			Pair{
				Key: models.NewKey("struct"),
				Entity: models.NewEntity(struct {
					name string
					age  int
				}{
					"misha",
					20,
				}),
			},
			Storage{
				pairs: map[Value]Value{
					models.NewKey("struct"): models.NewEntity(struct {
						name string
						age  int
					}{
						"misha",
						20,
					}),
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
				t.Fatalf("expected error: %v, got %v", tt.expectedError, gotErr)
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
		key           Value
		expectedValue interface{}
		expectedError error
		storage       Storage
	}{
		{
			"string return",
			models.NewKey("string data"),
			"data",
			nil,
			Storage{
				pairs: map[Value]Value{
					models.NewKey("string data"): models.NewEntity("data"),
				},
			},
		},
		{
			"int return",
			models.NewKey("number"),
			5,
			nil,
			Storage{
				pairs: map[Value]Value{
					models.NewKey("number"): models.NewEntity(5),
				},
			},
		},
		{
			"slice return",
			models.NewKey("slice"),
			[]int{1, 2, 3},
			nil,
			Storage{
				pairs: map[Value]Value{
					models.NewKey("slice"): models.NewEntity([]int{1, 2, 3}),
				},
			},
		},
		{
			"struct return ",
			models.NewKey("struct"),
			struct {
				name string
				age  int
			}{
				"misha",
				20,
			},
			nil,
			Storage{
				pairs: map[Value]Value{
					models.NewKey("struct"): models.NewEntity(struct {
						name string
						age  int
					}{
						"misha",
						20,
					}),
				},
			},
		},
		{
			"no such key return",
			models.NewKey("abab"),
			models.NewEntity(nil),
			models.ErrNoSuchKey,
			Storage{
				pairs: make(map[Value]Value),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.storage.Get(tt.key)

			if !reflect.DeepEqual(err, tt.expectedError) {
				t.Fatal(err)
			}
			if got != nil && !reflect.DeepEqual(got.Entity(), tt.expectedValue) {
				t.Fatalf("expected value: %v, got %v", tt.expectedValue, got)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	tests := []struct {
		name          string
		key           Value
		storage       Storage
		expectedError error
	}{
		{
			"simple delete",
			models.NewKey("key"),
			Storage{
				pairs: map[Value]Value{
					models.NewKey("key"): models.NewEntity("value"),
				},
			},
			models.ErrNoSuchKey,
		},
		{
			"delete with no such key",
			models.NewKey("abab"),
			Storage{
				pairs: map[Value]Value{
					models.NewKey("key"): models.NewEntity("value"),
				},
			},
			models.ErrNoSuchKey,
		},
		{
			"empty key",
			models.NewKey(""),
			Storage{
				pairs: make(map[Value]Value),
			},
			models.ErrNoSuchKey,
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
