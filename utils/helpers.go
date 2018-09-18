package utils

import (
	"os"
	"github.com/kamva/go-libs/exceptions"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"regexp"
)

func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ThrowException(responseMessage string, message string, code string, status int) {
	panic(exceptions.Exception{
		Message:         message,
		ResponseMessage: responseMessage,
		Code:            code,
		StatusCode:      status,
	})
}

func ThrowValidationException(responseMessage string, data interface{}, code string) {
	panic(exceptions.ValidationException{
		ResponseMessage: responseMessage,
		Data:            data,
		Code:            code,
	})
}

func StringInSlice(search string, list []string) bool {
	for _, item := range list {
		if search == item {
			return true
		}
	}
	return false
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ToCamelCase(snake string) string {
	parts := strings.Split(snake, "_")

	for i := 0; i < len(parts); i++ {
		parts[i+1] = strings.Title(parts[i+1])
	}

	return strings.Join(parts, "")
}

func ValidateObjectId(instanceId string) bool {
	result, _ := regexp.MatchString("^[0-9a-fA-F]{24}$", instanceId)

	return result
}
