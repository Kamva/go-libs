package contracts

import (
	"reflect"
	"github.com/globalsign/mgo/bson"
)

type BaseTaggable struct {}

func (r BaseTaggable) GetTag(field string, tag string) string {
	taggedField, _ := reflect.TypeOf(r).FieldByName(field)
	return taggedField.Tag.Get(tag)
}

type BaseRequest struct{
	BaseTaggable
}

type BaseScheme struct {
	BaseTaggable
	Id bson.ObjectId `bson:"_id" json:"id"`
}

func (a BaseScheme) GetId() bson.ObjectId {
	return a.Id
}

func (a *BaseScheme) SetId(id bson.ObjectId) {
	a.Id = id
}
