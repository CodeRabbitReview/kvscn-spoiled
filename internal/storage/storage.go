package storage

import (
	"encoding/json"
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/models"
	"reflect"
)

type Value interface {
	Type() reflect.Type
	Entity() interface{}
}

type Keyer interface {
	Value
}

type Entitier interface {
	Value
	JSON() json.RawMessage
}

// Pair combines a key and input value
type Pair struct {
	Key    Keyer
	Entity Entitier
}

func (p Pair) emptyKey() bool {
	if v, ok := p.Key.Entity().(string); ok {
		return v == ""
	}
	return p.Key.Entity() == nil
}

func (p Pair) nilEntity() bool {
	if v, ok := p.Entity.Entity().(string); ok {
		return v == ""
	}
	return p.Entity.Entity() == nil
}

type Storage struct {
	pairs map[Keyer]Entitier
}

func NewStorage() *Storage {
	return &Storage{
		pairs: make(map[Keyer]Entitier),
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
func (s *Storage) Get(key Keyer) (Entitier, error) {
	if len(s.pairs) == 0 {
		return nil, fmt.Errorf("no data in storage")
	}
	if v, ok := s.pairs[key]; ok {
		return v, nil
	}

	return nil, models.ErrNoSuchKey
}

// Delete remove data from storage
// A key is ignored if it does not exist
func (s *Storage) Delete(key Keyer) error {
	fmt.Println(key.Entity(), len(s.pairs) == 0)
	if len(s.pairs) == 0 {
		return fmt.Errorf("no data in storage")
	}
	if _, ok := s.pairs[key]; !ok {
		return models.ErrNoSuchKey
	}
	delete(s.pairs, key)
	return nil
}

func (s *Storage) GetAll() (map[Keyer]Entitier, error) {
	if len(s.pairs) == 0 {
		return nil, fmt.Errorf("no data in storage")
	}
	return s.pairs, nil
}
