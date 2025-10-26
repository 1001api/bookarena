package pkg

import (
	"github.com/golang-jwt/jwt/v5"
)

type JSONB map[string]interface{}
type RoleType string

const (
	RoleUser       RoleType = "user"
	RoleAdmin      RoleType = "admin"
	RoleSuperAdmin RoleType = "superadmin"

	// JWT
	AccessTokenName  = "__hk_asid"
	RefreshTokenName = "__hk_rsid"

	// Cookie
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	TokenIssuer      = "bookarena"
)

type StateClaims struct {
	ClientID    string `json:"client_id"`
	RedirectURL string `json:"redirect_url"`
	jwt.RegisteredClaims
}
