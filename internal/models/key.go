package models

import (
	"fmt"
	"reflect"
	"strings"
)

type Key struct {
	entity string
}

func NewKey(e interface{}) Key {
	if key, ok := e.(string); ok {
		if strings.HasSuffix(key, " ") {
			key = key[1:]
		}
		if strings.HasPrefix(key, " ") {
			key = key[:len([]rune(key))]
		}
		key = strings.ReplaceAll(key, " ", "_")
		return Key{entity: key}
	}

	k := fmt.Sprintf("%v", e)

	return Key{entity: k}
}

func (k Key) Type() reflect.Type {
	return reflect.TypeOf(k.entity)
}

func (k Key) Entity() interface{} {
	return k.entity
}
