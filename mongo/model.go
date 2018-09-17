package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"kamva.ir/libraries/contracts"
	"kamva.ir/libraries/exceptions"
	"kamva.ir/libraries/translation"
	"github.com/kataras/iris"
)

type Model struct {
	Session    *Session
	collection *mgo.Collection
	scheme     contracts.Scheme
}

func (m *Model) Collection() *mgo.Collection {
	return m.collection
}

func (m *Model) Initiate(collectionName string, indices ...mgo.Index) {
	if m.Session == nil {
		panic(exceptions.Exception{
			ResponseMessage: translation.Translate("internal_error"),
			Message:         "Mongodb session is not set.",
			Code:            "LIB_ERR",
			StatusCode:      iris.StatusInternalServerError,
		})

	}

	m.collection = m.Session.GetCollection(collectionName)

	for _, index := range indices {
		m.collection.EnsureIndex(index)
	}
}

func (m *Model) Create(data contracts.Scheme) error {
	m.scheme = data
	m.scheme.SetId(bson.NewObjectId())
	data = m.scheme
	return m.collection.Insert(&m.scheme)
}

func (m *Model) Find(id string) (contracts.Scheme, error) {
	err := m.collection.FindId(bson.ObjectIdHex(id)).One(&m.scheme)
	return m.scheme, err
}

func (m *Model) Update(data contracts.Scheme) error {
	err := m.collection.UpdateId(data.GetId(), &data)
	if err == nil {
		m.scheme = data
	}

	return err
}

func (m *Model) Delete(data contracts.Scheme) error {
	err := m.collection.Remove(&data)
	if err == nil {
		m.scheme = nil
	}

	return err
}

func (m *Model) Query(query interface{}) *mgo.Query {
	return m.collection.Find(query)
}
