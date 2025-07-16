package storage

import (
	"github.com/mishaprokop4ik/storage/internal/storage/models"
	"reflect"
)

type Value interface {
	Type() reflect.Type
	Entity() interface{}
}

// Pair combines a key and input value
type Pair struct {
	Key    Value
	Entity Value
}

func (p Pair) emptyKey() bool {
	return p.Key.Entity() == "" || p.Key.Entity() == nil
}

func (p Pair) nilEntity() bool {
	return p.Entity.Entity() == "" || p.Entity.Entity() == nil
}

type Storage struct {
	pairs map[Value]Value
}

func NewStorage() *Storage {
	return &Storage{
		pairs: make(map[Value]Value),
	}
}

//Put add new value to storage
func (s *Storage) Put(p Pair) error {
	if p.nilEntity() {
		return models.ErrNilInput
	}

	if p.emptyKey() {
		return models.ErrEmptyKeyString
	}
	s.pairs[p.Key] = p.Entity

	return nil
}

// Get returns a copy of value from storage
// If no such data by key returns NoSuchValue error
func (s *Storage) Get(key Value) (Value, error) {
	if v, ok := s.pairs[key]; ok {
		return v, nil
	}

	return nil, models.ErrNoSuchKey
}

// Delete remove data from storage
// A key is ignored if it does not exist
func (s *Storage) Delete(key Value) {
	delete(s.pairs, key)
}
