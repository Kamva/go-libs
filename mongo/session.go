package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
	"github.com/kataras/iris"
)

type Session struct {
	session *mgo.Session
	db      string
	issueCode string
}

func (s *Session) Copy() *Session {
	return &Session{session: s.session.Copy()}
}

func (s *Session) GetCollection(collection string) *mgo.Collection {
	return s.session.DB(s.db).C(collection)
}

func (s *Session) GetDB() *mgo.Database {
	return s.session.DB(s.db)
}

func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

func NewSession(url string, db string, exceptionCode string) *Session {
	session, err := mgo.Dial(url)

	if err != nil {
		panic(exceptions.Exception{
			ResponseMessage: translation.Translate("internal_error"),
			Message:         err.Error(),
			Code:            exceptionCode,
			StatusCode:      iris.StatusInternalServerError,
		})

	}

	return &Session{
		session: session,
		db: db,
		issueCode:exceptionCode,
	}
}
