package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
)

type Redis struct {
	connection redis.Conn
	issueCode  string
}

func (r Redis) Set(key string, value string) bool {
	result, err := r.connection.Do("SET", key, value)
	r.throwErrorIfAny(err)

	return result.(string) == "OK"
}

func (r Redis) Get(key string) string {
	value, err := redis.String(r.connection.Do("GET", key))
	r.throwErrorIfAny(err)

	return value
}

func (r Redis) LPush(key string, value interface{}) int {
	result, err := redis.Int(r.connection.Do("LPUSH", key, value))
	r.throwErrorIfAny(err)

	return result
}

func (r Redis) LRange(key string, start int, end int) []interface{} {
	values, err := redis.Values(r.connection.Do("LRANGE", key, start, end))
	r.throwErrorIfAny(err)

	return values
}

func (r Redis) Exists(key string) bool {
	result, err := redis.Int(r.connection.Do("EXISTS", key))
	r.throwErrorIfAny(err)

	return result > 0
}

func (r Redis) FlushAll() {
	_, err := r.connection.Do("FLUSHALL")
	r.throwErrorIfAny(err)
}

func (r Redis) throwErrorIfAny(err error) {
	if err != nil {
		panic(exceptions.Exception{
			Message         : err.Error(),
			ResponseMessage : translation.Translate("internal_error"),
			Code            : r.issueCode,
			StatusCode      : iris.StatusInternalServerError,
		})
	}
}

func NewRedis(url string, exceptionCode string) *Redis {
	connection, err := redis.DialURL(url)

	if err != nil {
		panic(exceptions.Exception{
			Message         : err.Error(),
			ResponseMessage : translation.Translate("internal_error"),
			Code            : exceptionCode,
			StatusCode      : iris.StatusInternalServerError,
		})
	}

	return &Redis{
		connection: connection,
		issueCode: exceptionCode,
	}
}
