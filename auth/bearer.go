package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/utils"
)

var (
	ErrExpired           = errors.New("Token expired")
	ErrMissingToken      = errors.New("Token missing")
	ErrMissingExpires    = errors.New("Token payload does not have the exp-field")
	ErrExpiresParseError = errors.New("Failed to parse the exp-field in token payload ")
	ErrTokenParseError   = errors.New("Token looks invalid. Expected to find JWT with 3 parts seperated by '.'")
)

type ImpersionationCredentials struct {
	// Username to impersonate with. Needs to have the impersonation-role
	Username,
	Password,
	// UserID to impersonate as. This is preferred over UserNameToImpersonate
	UserIDToImpersonate string
	// Will perform a lookup to get the ID of the username.
	UserNameToImpersonate string
}

type TokenPayload = struct {
	Raw       map[string]interface{}
	Expires   time.Time
	ExpiresIn string
	Token     string
}

type BearerTokenCreatorOptions struct {
	ImpersionationCredentials
	ClientID     string
	RedirectUri  string
	Endpoint     string
	EndpointType string
	Expires      *time.Time
	ClientSecret string
}
type BearerTokenCreator struct {
	BearerTokenCreatorOptions
	l                  logger.AppLogger
	clientToken        string
	ClientTokenExpires *time.Time
	sync.Mutex
}

func NewBearerTokenCreator(l logger.AppLogger, options BearerTokenCreatorOptions) BearerTokenCreator {
	if !strings.HasSuffix(options.Endpoint, "/") {
		options.Endpoint = options.Endpoint + "/"
	}
	return BearerTokenCreator{
		BearerTokenCreatorOptions: options,
		l:                         l,
	}
}

func IsRequired(key string) error {
	return fmt.Errorf("%s is required", key)
}

func getPayloadUser(payload map[string]interface{}) string {
	if payload == nil {
		return ""
	}
	for _, k := range []string{"name", "given_name", "preferred_username", "email", "sub"} {
		value, ok := payload[k]
		if !ok {
			continue
		}
		s, ok := value.(string)
		if !ok {
			continue
		}
		return s

	}
	return ""
}

func (bc *BearerTokenCreator) Renew() error {
	panic("Not implemented")
}
func (bc *BearerTokenCreator) Validate() error {
	_, err := bc.PrintValidity()
	return err
}

type ImpersionateOptions struct {
	UserName, UserID string
}

func (bc *BearerTokenCreator) Impersionate(options ImpersionateOptions) (t TokenPayload, err error) {
	switch bc.EndpointType {
	case "keycloak":
		if options.UserID == "" && options.UserName == "" {
			return t, fmt.Errorf("To impersonate, a userid/username is required")
		}
		kc := Keycloak{
			l:            bc.l,
			ClientID:     bc.ClientID,
			ClientSecret: bc.ClientSecret,
			Username:     bc.Username,
			Password:     bc.Password,
			Endpoint:     bc.Endpoint,
			RedirectUri:  bc.RedirectUri,
		}
		ct, err := kc.Retrieve()
		if err != nil {
			return t, err
		}
		if options.UserID == "" {
			id, err := kc.GetUserIDByUsername(ct, options.UserName)
			if err != nil {
				return t, err
			}
			options.UserID = id
		}
		// The client-token from above is spent after impersonation.
		return kc.Impersonate(ct, options.UserID)
	default:
		return t, fmt.Errorf("EndpointType '%s' does not support impersonation", bc.EndpointType)
	}
}

func (bc *BearerTokenCreator) PrintValidity() (string, error) {
	if bc.Expires == nil {
		return "", nil
	}
	now := time.Now()
	diff := bc.Expires.Sub(now)
	if diff <= 0 {
		return fmt.Sprintf("Token expired %s ago", utils.PrettyDuration(diff)), ErrExpired
	}
	return fmt.Sprintf("Token expires in %s", utils.PrettyDuration(diff)), nil
}
func (bc *BearerTokenCreator) ParseToken(s string) (TokenPayload, error) {
	return TokenStringToTokenPayload(s)
}

func TokenStringToTokenPayload(s string) (TokenPayload, error) {
	var tp TokenPayload
	if s == "" {
		return tp, nil
	}

	payload, err := ParseToken(s)
	if err != nil {
		return tp, err
	}
	tp.Raw = payload
	tp.Token = s
	if exp, ok := payload["exp"]; ok {
		expN, ok := exp.(float64)

		if !ok || expN <= 0 {
			return tp, ErrExpiresParseError
		}
		pex := time.Unix(int64(expN), 0)
		tp.Expires = pex

	} else {
		return tp, ErrMissingExpires
	}

	return tp, err
}

func ParseToken(token string) (map[string]interface{}, error) {
	if token == "" {
		return nil, ErrMissingToken
	}

	// To copy straight from firefox
	token = strings.TrimPrefix(token, "Authorization: ")
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}
	rawToken := strings.TrimPrefix(token, "Bearer ")
	tokenSplitted := strings.Split(rawToken, ".")
	if len(tokenSplitted) != 3 {
		return nil, ErrTokenParseError
	}
	data := tokenSplitted[1]
	payloadRaw, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode token: %w", err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal([]byte(payloadRaw), &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload in token: %w", err)
	}
	return payload, nil
}
