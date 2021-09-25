package cmd

import "time"

// TODO: use general map[string]interface{} instead to support all kinds of jwt's
type JwtPayload struct {
	Exp       int64      `json:"exp"`
	ExpiresAt *time.Time `json:"-"`
	Iat       int64      `json:"iat"`
	AuthTime  int64      `json:"auth_time"`
	Jti       string     `json:"jti"`
	Iss       string     `json:"iss"`
	// Aud                                    string         `json:"aud"`
	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	Nonce             string         `json:"nonce"`
	SessionState      string         `json:"session_state"`
	ACR               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	EmailVerified     bool           `json:"email_verified"`
	GraphqlEndpoint   string         `json:"graphqlEndpoint"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account RealmAccess `json:"account"`
}
