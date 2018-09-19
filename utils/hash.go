package utils

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
	"github.com/kataras/iris"
)

type Hash struct {
	salt      string
	issueCode string
}

func (h *Hash) Generate(password string) string {
	saltedBytes := []byte(password + h.salt)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.MinCost)

	if err != nil {
		panic(exceptions.Exception{
			Message         : err.Error(),
			ResponseMessage : translation.Translate("internal_error"),
			Code            : h.issueCode,
			StatusCode      : iris.StatusInternalServerError,
		})
	}

	hash := string(hashedBytes[:])
	return hash
}

func (h *Hash) Compare(hash string, password string) error {
	incoming := []byte(password + h.salt)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func NewHash(salt string, exceptionCode string) *Hash {
	return &Hash{
		salt: salt,
		issueCode: exceptionCode,
	}
}
