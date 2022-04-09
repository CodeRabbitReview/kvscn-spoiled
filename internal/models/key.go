package models

import (
	"reflect"
)

type Key struct {
	entity interface{}
}

func NewKey(e interface{}) Key {
	return Key{entity: e}
}

func (k Key) Type() reflect.Type {
	return reflect.TypeOf(k.entity)
}

func (k Key) Entity() interface{} {
	return k.entity
}
