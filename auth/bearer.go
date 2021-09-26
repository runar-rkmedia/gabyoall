package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

type BearerTokenCreatorOptions struct {
	ImpersionationCredentials
	ClientID     string
	RedirectUri  string
	Endpoint     string
	Token        string
	EndpointType string
	Expires      *time.Time
	Payload      map[string]interface{}
	ClientSecret string
}
type BearerTokenCreator struct {
	BearerTokenCreatorOptions
	l           logger.AppLogger
	clientToken string
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

func (bc *BearerTokenCreator) Retrieve() error {
	bc.l.Info().Msg("Retrieving bearer-token with clientID/secret")
	if bc.EndpointType != "keycloak" {
		return fmt.Errorf("Not implemented for auth.endpointType: '%s'. Supported types are: keycloak", bc.EndpointType)
	}
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", bc.ClientID)
	data.Set("client_secret", bc.ClientSecret)
	data.Set("username", bc.Username)
	data.Set("password", bc.Password)
	url := bc.Endpoint + "protocol/openid-connect/token"
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("faild to create request: %w", err)
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Unexpected statuscode returned for endpoint %s: %d %s %s %s", url, res.StatusCode, res.Status, contentType, string(body))
	}
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return err
	}
	token := jsonResponse["access_token"]
	if token == nil {
		return fmt.Errorf("failed to find access_token in payload %#v", jsonResponse)
	}
	tokenStr, ok := token.(string)
	if !ok {
		return fmt.Errorf("token does not seem to be valid: %#v", tokenStr)
	}

	bc.clientToken = tokenStr
	bc.l.Info().Msg("successfully retrieved impersonator-token")
	if bc.ImpersionationCredentials.UserIDToImpersonate != "" {
		return bc.Impersonate(bc.ImpersionationCredentials.UserIDToImpersonate)
	}
	sid, err := bc.GetUserIDByUsername(bc.ImpersionationCredentials.UserNameToImpersonate)
	if err != nil {
		panic(err)
	}
	fmt.Println(sid)
	return bc.Impersonate(sid)
}
func (bc *BearerTokenCreator) GetUserIDByUsername(userName string) (string, error) {
	bc.l.Info().Str("userName", userName).Msg("Attempting to find user by userName")
	if bc.clientToken == "" {
		return "", fmt.Errorf("missing client-token, please run Retrieve first")
	}

	uri := strings.Replace(bc.Endpoint, "/realms/", "/admin/realms/", 1) + "users?briefRepresentation=true&username=" + userName
	r, err := http.NewRequest(http.MethodGet, uri, nil)
	r.Header.Set("Authorization", "Bearer "+bc.clientToken)
	if err != nil {
		return "", fmt.Errorf("failed to create request to %s: %w", uri, err)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {

		return "", fmt.Errorf("failed to perform request to %s: %w", uri, err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of request to %s: %w", uri, err)
	}
	switch {
	case res.StatusCode == 403:
		return "", fmt.Errorf("Receeived forbidden as status-code for %s %d %s (%s). This might mean you are missing the 'query-users/manage-users'-role for the user '%s'. This is needed to look up users. You may also use userID directly. ", uri, res.StatusCode, res.Status, string(body), bc.ImpersionationCredentials.Username)
	case res.StatusCode >= 300:
		return "", fmt.Errorf("Unexpected status-code for %s %d %s %s", uri, res.StatusCode, res.Status, string(body))
	}
	var users []KeycloakBriedUserInfo
	err = json.Unmarshal(body, &users)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response from users-endpoint %s: %w", uri, err)
	}
	// The search is case-sensitive
	userLowered := strings.ToLower(userName)
	foundLowered := -1
	for i := 0; i < len(users); i++ {
		// Prefer case-sensitive match
		if users[i].Username == userName {
			return users[i].ID, nil
		}
		// ... but allow insensitive match
		if strings.ToLower(users[i].Username) == userLowered {
			foundLowered = i
		}

	}
	if foundLowered >= 0 {
		return users[foundLowered].Username, nil
	}
	return "", nil

}
func (bc *BearerTokenCreator) Impersonate(userID string) error {
	bc.l.Info().Str("userID", userID).Msg("attempting to imppersonate user with impersonator-token")
	if bc.clientToken == "" {
		return fmt.Errorf("clientToken not set. Please call Retrieve first")
	}
	if bc.RedirectUri == "" {
		return fmt.Errorf("redirectUri is required for impersonation")
	}
	uri := strings.Replace(bc.Endpoint, "/realms/", "/admin/realms/", 1) + "users/" + userID + "/impersonation"

	r, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create http-request: %w", err)
	}
	r.Header.Set("Authorization", "Bearer "+bc.clientToken)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("failed while performing http-request: %w", err)
	}
	contentType := res.Header.Get("Content-Type")
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Unexpected statuscode returned for endpoint %s: %d %s %s %s", uri, res.StatusCode, res.Status, contentType, string(body))
	}
	cookies := res.Cookies()
	var validCookies []*http.Cookie
	now := time.Now()
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Expires.After(now) || cookies[i].RawExpires == "" {
			validCookies = append(validCookies, cookies[i])
		}
	}
	var jsonResponse struct {
		Redirect  string
		sameRealm bool
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	if jsonResponse.Redirect == "" {
		return fmt.Errorf("Missing redirect-value in json-response: %s", string(body))
	}

	// redirectUrl := "    .get(uri"$authUri?response_mode=fragment&response_type=token&client_id=${config.clientId}&redirect_uri=${config.redirectUri}")"
	redirectUrl := jsonResponse.Redirect
	redirectUrl = bc.Endpoint + "protocol/openid-connect/auth?response_mode=fragment&response_type=token&client_id=" + bc.ClientID + "&redirect_uri=" + bc.RedirectUri
	reqRedirect, err := http.NewRequest(http.MethodGet, redirectUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create redirect-request: %s %w", redirectUrl, err)
	}

	for i := 0; i < len(validCookies); i++ {
		reqRedirect.AddCookie(validCookies[i])
	}
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resRedirect, err := c.Do(reqRedirect)
	if err != nil {
		return fmt.Errorf("failed to perform redirect-request: %s %w", redirectUrl, err)
	}
	defer resRedirect.Body.Close()
	if resRedirect.StatusCode >= 400 {
		redirectBody, err := ioutil.ReadAll(resRedirect.Body)
		return fmt.Errorf("Unexpected status-code for redirect-call: %d %s %s %s %w", resRedirect.StatusCode, resRedirect.Status, redirectUrl, string(redirectBody), err)
	}
	location := resRedirect.Header.Get("Location")
	if location == "" {
		return fmt.Errorf("no location found in headers")
	}
	u, err := url.ParseQuery(location)
	if err != nil {
		return fmt.Errorf("failed to parse location-uri: %w", err)
	}
	tokenStr := u.Get("access_token")
	if tokenStr == "" {
		var keys []string
		for k := range u {
			keys = append(keys, k)
		}
		return fmt.Errorf("could not find 'access_token' in query-params: %v. This may be because the client does not have Implicit Flow enabled within keycloak", keys)
	}
	bc.Token = tokenStr
	_, err = bc.ParseToken()
	if err != nil {
		return err
	}
	user := getPayloadUser(bc.Payload)
	bc.l.Info().Str("user", user).Msg("Successfully performed impersonation")
	return nil

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
func (bc *BearerTokenCreator) ParseToken() (map[string]interface{}, error) {
	if bc.Payload != nil {
		return bc.Payload, nil
	}

	payload, err := ParseToken(bc.Token)
	if err != nil {
		return payload, err
	}
	bc.Payload = payload
	if exp, ok := payload["exp"]; ok {
		expN, ok := exp.(float64)

		if !ok || expN <= 0 {
			return payload, ErrExpiresParseError
		}
		pex := time.Unix(int64(expN), 0)
		bc.Expires = &pex

	} else {
		return payload, ErrMissingExpires
	}

	return payload, err
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
