package jwt

import (
	"time"
	"github.com/kataras/iris"
	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/go.uuid"
	"github.com/kamva/go-libs/redis"
	"github.com/kamva/go-libs/contracts"
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
)

type ResponseCodes struct {
	InvalidToken     string
	BlacklistedToken string
	InactiveToken    string
	ExpiredToken     string
	InternalIssue    string
}

type Config struct {
	TTL    int
	TTR    int
	Secret string
	Codes  ResponseCodes
}

type JWT struct {
	config         Config
	authErrMessage string
	cache          *redis.Redis
}

func (t *JWT) Sign(claims *ServiceClaims) string {
	id := uuid.Must(uuid.NewV4())

	claims.Issuer = "kite"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Unix() + int64(t.config.TTL*60)
	claims.Id = id.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(t.config.Secret))

	if err != nil {
		t.internalErrorException(err.Error())
	}

	return tokenString
}

func (t *JWT) SignAuthenticatable(authenticatable contracts.Authenticatable) string {
	userId := authenticatable.GetUserId()
	instanceId := authenticatable.GetInstanceId()

	var claims = new(ServiceClaims)
	claims.UserId = userId

	if instanceId != "" {
		claims.InstanceId = instanceId
	}

	return t.Sign(claims)
}

func (t *JWT) Logout(tokenString string) bool {
	token := t.Tokenize(tokenString)

	return t.Block(token)
}

func (t *JWT) Verify(tokenString string, silent bool) *jwt.Token {
	token := t.Tokenize(tokenString)

	if !t.isBlacklisted(token) {
		if claims, ok := token.Claims.(*ServiceClaims); ok {
			if claims.Validate() {
				return token
			}

			if claims.IsExpired(t.config.TTR) {
				t.expiredException()
			} else if silent {
				return token
			}

			t.inactivatedException()
		} else {
			t.invalidTokenException("Invalid Claims")
		}
	} else {
		t.blackListedException()
	}

	return token
}

func (t *JWT) Tokenize(tokenString string) *jwt.Token {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.config.Secret), nil
	})

	if err != nil {
		if cErr, ok := err.(ClaimError); ok {
			switch cErr.ErrorType {
			case CEInactiveToken:
				t.inactivatedException()
			case CEInvalidToken:
				t.invalidTokenException(cErr.Message)
			}
		}
		t.invalidTokenException(err.Error())
	}

	return token
}

func (t *JWT) RefreshToken(tokenString string) string {
	token := t.Verify(tokenString, true)
	if t.Block(token) {
		if claims, ok := token.Claims.(*ServiceClaims); ok {
			return t.Sign(claims)
		}
	} else {
		t.internalErrorException("Cannot blacklist token.")
	}

	return ""
}

func (t *JWT) Block(token *jwt.Token) bool {
	key := "jwt." + t.getTokenKey(token)
	return t.cache.Set(key, "blacklisted")
}

func (t *JWT) getTokenKey(token *jwt.Token) string {
	if claims, ok := token.Claims.(*ServiceClaims); ok {
		return claims.GetId()
	} else {
		t.invalidTokenException("Invalid Claims")
		return ""
	}
}

func (t *JWT) isBlacklisted(token *jwt.Token) bool {
	key := "jwt." + t.getTokenKey(token)
	return t.cache.Exists(key) && t.cache.Get(key) == "blacklisted"
}

func (t *JWT) GetUserClaims(tokenString string) (string, string) {
	token := t.Verify(tokenString, false)

	claim, ok := token.Claims.(*ServiceClaims)

	if !ok {
		t.invalidTokenException("Invalid token claims.")
	}

	return claim.UserId, claim.InstanceId
}

func (t *JWT) invalidTokenException(errString string) {
	panic(exceptions.Exception{
		Message:         errString,
		ResponseMessage: t.authErrMessage,
		Code:            t.config.Codes.InvalidToken,
		StatusCode:      iris.StatusUnauthorized,
	})
}

func (t *JWT) blackListedException() {
	panic(exceptions.Exception{
		Message:         "Token Blacklisted!",
		ResponseMessage: t.authErrMessage,
		Code:            t.config.Codes.BlacklistedToken,
		StatusCode:      iris.StatusUnauthorized,
	})
}

func (t *JWT) expiredException() {
	panic(exceptions.Exception{
		Message:         "Token Expired!",
		ResponseMessage: t.authErrMessage,
		Code:            t.config.Codes.ExpiredToken,
		StatusCode:      iris.StatusUnauthorized,
	})
}

func (t *JWT) inactivatedException() {
	panic(exceptions.Exception{
		Message:         "Token is inactive!",
		ResponseMessage: t.authErrMessage,
		Code:            t.config.Codes.InactiveToken,
		StatusCode:      iris.StatusUnauthorized,
	})
}

func (t *JWT) internalErrorException(errString string) {
	panic(exceptions.Exception{
		Message:         errString,
		ResponseMessage: translation.Translate("internal_error"),
		Code:            t.config.Codes.InternalIssue,
		StatusCode:      iris.StatusInternalServerError,
	})
}

func (t *JWT) validateConfig() {
	if t.config.TTL <= 0 {
		t.internalErrorException("JWT TTL must be greater than 0")
	}

	if t.config.TTR <= 0 {
		t.internalErrorException("JWT TTR must be greater than 0")
	}

	if t.config.TTR < t.config.TTL {
		t.internalErrorException("JWT TTR must be greater than TTL")
	}

	if len(t.config.Secret) != 32 {
		t.internalErrorException("Length of JWT secret should be 32")
	}
}

func NewJWT(config Config, authErrMessage string, cache *redis.Redis) *JWT {
	jwtObject := JWT{
		config:         config,
		cache:          cache,
		authErrMessage: authErrMessage,
	}

	jwtObject.validateConfig()

	return &jwtObject
}
