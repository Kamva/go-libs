package crypt

import (
	"fmt"
	"io"
	"errors"
	"crypto/aes"
	"crypto/rand"
	"crypto/cipher"
	"encoding/base64"
	"github.com/kataras/iris"
	"kamva.ir/libraries/exceptions"
	"kamva.ir/libraries/translation"
)

type Crypt struct {
	key       []byte
	issueCode string
}

func (c *Crypt) Encrypt(text string) string {
	plaintext := []byte(text)

	block := c.getCipherBlock()

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		c.throwException(err.Error())
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(cipherText)
}

func (c *Crypt) Decrypt(text string) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(text)

	block := c.getCipherBlock()

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}

func (c *Crypt) getCipherBlock() cipher.Block {
	block, err := aes.NewCipher(c.key)

	if err != nil {
		c.throwException(err.Error())
	}

	return block
}

func (c *Crypt) throwException(message string) {
	panic(exceptions.Exception{
		ResponseMessage: translation.Translate("internal_error"),
		Message:         message,
		Code:            c.issueCode,
		StatusCode:      iris.StatusInternalServerError,
	})
}

func NewCrypt(key string, exceptionCode string) *Crypt {
	return &Crypt{
		key:       []byte(key),
		issueCode: exceptionCode,
	}
}
