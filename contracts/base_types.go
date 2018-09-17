package contracts

import "reflect"

type BaseRequest struct {}

func (r BaseRequest) GetTag(field string, tag string) string {
	taggedField, _ := reflect.TypeOf(r).FieldByName(field)
	return taggedField.Tag.Get(tag)
}
