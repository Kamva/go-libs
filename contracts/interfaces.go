package contracts

import "github.com/globalsign/mgo/bson"

type Taggable interface {
	GetTag(interface{}, string, string) string
}

type RequestData interface {
	Taggable
}

type Serializable interface {
	Serialize() string
}

type Pushable interface {
	Serializable
}

type Scheme interface {
	Taggable
	GetId() bson.ObjectId
	SetId(bson.ObjectId)
}

type Authenticatable interface {
	Scheme
	GetUserId() string
	GetInstanceId() string
	GetPassword() string
	SetPassword(string)
}
