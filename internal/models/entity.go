package models

import (
	"encoding/json"
	"reflect"
)

type Entity struct {
	entity interface{}
	json   json.RawMessage
}

func NewEntity(e interface{}, body json.RawMessage) Entity {
	return Entity{entity: e, json: body}
}

func (e Entity) Type() reflect.Type {
	return reflect.TypeOf(e.entity)
}

func (e Entity) Entity() interface{} {
	return e.entity
}

func (e Entity) JSON() json.RawMessage {
	return e.json
}
