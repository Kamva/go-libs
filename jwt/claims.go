package jwt

import (
	"github.com/kamva/go-libs/utils"
	"time"
)

const (
	CEInactiveToken = iota
	CEInvalidToken  = iota
)

type ClaimError struct {
	Message   string
	ErrorType int
}

func (e ClaimError) Error() string {
	return e.Message
}

type ServiceClaims struct {
	Audience    string              `json:"aud,omitempty"`
	ExpiresAt   int64               `json:"exp,omitempty"`
	Id          string              `json:"jti,omitempty"`
	IssuedAt    int64               `json:"iat,omitempty"`
	Issuer      string              `json:"iss,omitempty"`
	NotBefore   int64               `json:"nbf,omitempty"`
	Subject     string              `json:"sub,omitempty"`
	InstanceId  string              `json:"iid,omitempty"`
	UserId      string              `json:"uid,omitempty"`
	Permissions map[string][]string `json:"prm,omitempty"`
	Services    []string            `json:"srv,omitempty"`
	Admin       bool                `json:"adm,omitempty"`
	God         bool                `json:"god,omitempty"`
}

func (c ServiceClaims) Valid() error {

	if c.IsInactive() {
		return ClaimError{
			Message:   "Authentication token has expired.",
			ErrorType: CEInactiveToken,
		}
	}

	if !c.ValidateClaims() {
		return ClaimError{
			Message:   "Authentication token is invalid.",
			ErrorType: CEInvalidToken,
		}
	}

	return nil
}

func (c ServiceClaims) Validate() bool {
	return !c.IsInactive() && c.ValidateClaims()
}

func (c *ServiceClaims) ValidateClaims() bool {
	return c.VerifyIssuedAt() && c.VerifyNotBefore() && c.VerifyJWTId() && c.VerifyIssuer()
}

func (c *ServiceClaims) IsInactive() bool {
	now := time.Now().Unix()
	return now >= c.ExpiresAt
}

func (c *ServiceClaims) IsExpired(ttr int) bool {
	now := time.Now().Unix()
	return (now - c.IssuedAt) > int64(ttr*60)
}

func (c *ServiceClaims) VerifyIssuedAt() bool {
	now := time.Now().Unix()
	return now >= c.IssuedAt
}

func (c *ServiceClaims) VerifyNotBefore() bool {
	now := time.Now().Unix()
	return now >= c.NotBefore || c.NotBefore == 0
}

func (c *ServiceClaims) VerifyJWTId() bool {
	return c.Id != ""
}

func (c *ServiceClaims) VerifyIssuer() bool {
	return utils.StringInSlice(c.Issuer, utils.GetTrustedAuthIssuers())
}

func (c ServiceClaims) GetId() string {
	return c.Id
}
