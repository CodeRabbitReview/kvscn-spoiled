package models

import (
	"reflect"
)

type Entity struct {
	entity interface{}
}

func NewEntity(e interface{}) Entity {
	return Entity{entity: e}
}

func (e Entity) Type() reflect.Type {
	return reflect.TypeOf(e.entity)
}

func (e Entity) Entity() interface{} {
	return e.entity
}
