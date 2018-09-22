package contracts

import (
	"reflect"
)

type BaseTaggable struct{}

func (r BaseTaggable) GetTag(caller interface{}, field string, tag string) string {
	taggedField, _ := reflect.TypeOf(caller).Elem().FieldByName(field)
	return taggedField.Tag.Get(tag)
}

type BaseRequest struct {
	BaseTaggable
}
