package storage

import "errors"

var ErrNoSuchKey = errors.New("no such value by this key")
var ErrNilInput = errors.New("nil error in key data")
var ErrEmptyKeyString = errors.New("empty key value")

type Key string

// Pair combines a key and input value
type Pair struct {
	Key   Key
	Value interface{}
}

func (p Pair) emptyKey() bool {
	return p.Key == ""
}

func (p Pair) nilInput() bool {
	return p.Value == nil
}

type Storage struct {
	pairs map[Key]interface{}
}

func NewStorage() *Storage {
	return &Storage{
		pairs: make(map[Key]interface{}),
	}
}

//Put add new value to storage
func (s *Storage) Put(p Pair) error {
	if p.emptyKey() {
		return ErrEmptyKeyString
	}
	if p.nilInput() {
		return ErrNilInput
	}

	s.pairs[p.Key] = p.Value

	return nil
}

// Get returns a copy of value from storage
// If no such data by key returns NoSuchValue error
func (s *Storage) Get(k Key) (interface{}, error) {
	if v, ok := s.pairs[k]; ok {
		return v, nil
	}

	return nil, ErrNoSuchKey
}

// Delete remove data from storage
// A key is ignored if it does not exist
func (s *Storage) Delete(k Key) {
	delete(s.pairs, k)
}
