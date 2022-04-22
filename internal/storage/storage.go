package storage

import (
	"encoding/json"
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/models"
	"github.com/mishaprokop4ik/storage/internal/recoverer"
	"reflect"
	"sync"
)

// Value responds to simple value in storage
type Value interface {
	Type() reflect.Type
	Entity() interface{}
}

// Keyer responds to key value in storage
// Keyer has to have Value methods
type Keyer interface {
	Value
}

// Entitier responds to entity value in storage
// Entitier has to have Value methods and JSON that gives
// entity value in json format
type Entitier interface {
	Value
	JSON() json.RawMessage
}

// Pair combines a Keyer and input Entitier interfaces
type Pair struct {
	Key    Keyer
	Entity Entitier
}

func (p Pair) emptyKey() bool {
	if v, ok := p.Key.Entity().(string); ok {
		return v == "" || v == "{}" || v == "<nil>"
	}
	return p.Key.Entity() == nil
}

func (p Pair) nilEntity() bool {
	if v, ok := p.Entity.Entity().(string); ok {
		return v == "" || v == "{}" || v == "<nil>"
	}
	return p.Entity.Entity() == nil
}

type Storage struct {
	pairs   map[Keyer]Entitier
	mu      *sync.RWMutex
	resumer resumer
}

type resumer interface {
	RecoverData(action, data string, actions recoverer.Actions) error
}

func NewStorage(r resumer) *Storage {
	return &Storage{
		pairs:   make(map[Keyer]Entitier),
		mu:      &sync.RWMutex{},
		resumer: r,
	}
}

//Put adds new value to storage or update old value by key
//If key or value is empty - return errors models.ErrNilInput and models.ErrEmptyKey
func (s *Storage) Put(p Pair) error {
	if p.nilEntity() {
		return models.ErrNilInput
	}

	if p.emptyKey() {
		return models.ErrEmptyKey
	}

	if v, ok := s.pairs[p.Key]; !ok ||
		v != nil && string(v.JSON()) != string(p.Entity.JSON()) {
		if s.resumer != nil {
			err := s.resumer.RecoverData("p", string(p.Entity.JSON()), recoverer.DefaultActions)
			if err != nil {
				return err
			}
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		s.pairs[p.Key] = p.Entity
	}
	return nil
}

// Get returns a copy of value from storage
// If no such data by key returns models.ErrNilInput error
// If there is not any data in storage returns no data in storage error
func (s *Storage) Get(key Keyer) (Entitier, error) {
	if len(s.pairs) == 0 {
		return nil, fmt.Errorf("no data in storage")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.pairs[key]; ok {
		return v, nil
	}
	return nil, models.ErrNoSuchKey
}

// Delete removes data from storage
// If there is not any data by key
// If there is not any data in storage returns no data in storage error
func (s *Storage) Delete(key Keyer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pairs) == 0 {
		return fmt.Errorf("no data in storage")
	}
	if _, ok := s.pairs[key]; ok {
		b, err := json.Marshal(key.Entity())
		if err != nil {
			return err
		}
		if s.resumer != nil {
			err = s.resumer.RecoverData("d",
				fmt.Sprintf(`{"key": %s}`, string(b)), recoverer.DefaultActions)
			if err != nil {
				return err
			}
		}
	} else {
		return models.ErrNoSuchKey
	}

	delete(s.pairs, key)
	return nil
}

// GetAll returns all data from storage
// If there is not any data in storage returns no data in storage error
func (s *Storage) GetAll() (map[Keyer]Entitier, error) {
	if len(s.pairs) == 0 {
		return nil, fmt.Errorf("no data in storage")
	}
	return s.pairs, nil
}
