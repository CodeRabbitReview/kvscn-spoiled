package storage

import "errors"

var ErrNoSuchValue = errors.New("no such value by this key")

type Key string

type Pair struct {
	Key   Key
	Value interface{}
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
func (s *Storage) Put(p *Pair) {
	s.pairs[p.Key] = p.Value
}

// Get returns a copy of value from storage
// If no such data by key returns NoSuchValue error
func (s *Storage) Get(k Key) (interface{}, error) {
	if v, ok := s.pairs[k]; ok {
		return v, nil
	}

	return nil, ErrNoSuchValue
}

// Delete remove data from storage
// A key is ignored if it does not exist
func (s *Storage) Delete(k Key) {
	delete(s.pairs, k)
}
