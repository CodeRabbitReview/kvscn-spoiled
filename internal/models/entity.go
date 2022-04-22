package models

import (
	"encoding/json"
	"reflect"
	"regexp"
)

type Entity struct {
	entity interface{}
	json   json.RawMessage
}

func NewClearEntity(e interface{}, body json.RawMessage) (Entity, error) {
	removeAllSpaces, err := regexp.Compile(`\r|\t|\n| `)
	if err != nil {
		return Entity{}, err
	}

	entity := removeAllSpaces.ReplaceAllString(string(body), "")
	return Entity{entity: e, json: json.RawMessage(entity)}, nil
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
