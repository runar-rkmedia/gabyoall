package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
)

type KeycloakBriedUserInfo struct {
	ID                         string        `json:"id"`
	CreatedTimestamp           int64         `json:"createdTimestamp"`
	Username                   string        `json:"username"`
	Enabled                    bool          `json:"enabled"`
	Totp                       bool          `json:"totp"`
	EmailVerified              bool          `json:"emailVerified"`
	FirstName                  string        `json:"firstName"`
	LastName                   string        `json:"lastName"`
	Email                      string        `json:"email"`
	DisableableCredentialTypes []string      `json:"disableableCredentialTypes"`
	RequiredActions            []interface{} `json:"requiredActions"`
	NotBefore                  int64         `json:"notBefore"`
	Access                     Access        `json:"access"`
	Attributes                 *Attributes   `json:"attributes,omitempty"`
}

type Access struct {
	ManageGroupMembership bool `json:"manageGroupMembership"`
	View                  bool `json:"view"`
	MapRoles              bool `json:"mapRoles"`
	Impersonate           bool `json:"impersonate"`
	Manage                bool `json:"manage"`
}

type Attributes struct {
	Locale []string `json:"locale"`
}

type Keycloak struct {
	l                                                                 logger.AppLogger
	ClientID, ClientSecret, Username, Password, Endpoint, RedirectUri string
	sync.Mutex
}

func (bc *Keycloak) Retrieve() (t TokenPayload, err error) {
	if bc.ClientID == "" {
		return t, IsRequired("ClientID")
	}
	if bc.Password == "" {
		return t, IsRequired("Password")
	}
	if bc.Username == "" {
		return t, IsRequired("Username")
	}
	if bc.Password == "" {
		return t, IsRequired("Password")
	}
	bc.Lock()
	defer bc.Unlock()
	bc.l.Info().Msg("Retrieving bearer-token with clientID/secret")
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", bc.ClientID)
	data.Set("client_secret", bc.ClientSecret)
	data.Set("username", bc.Username)
	data.Set("password", bc.Password)
	url := bc.Endpoint + "protocol/openid-connect/token"
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return t, fmt.Errorf("faild to create request: %w", err)
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return t, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return t, err
	}
	if res.StatusCode > 299 {
		return t, fmt.Errorf("Unexpected statuscode returned for endpoint %s: %d %s %s %s", url, res.StatusCode, res.Status, contentType, string(body))
	}
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return t, err
	}
	token := jsonResponse["access_token"]
	if token == nil {
		return t, fmt.Errorf("failed to find access_token in payload %#v", jsonResponse)
	}
	tokenStr, ok := token.(string)
	if !ok {
		return t, fmt.Errorf("token does not seem to be valid: %#v", tokenStr)
	}

	t, err = TokenStringToTokenPayload(tokenStr)
	if err != nil {
		return t, err
	}

	bc.l.Info().Msg("successfully retrieved impersonator-token")
	return t, nil
}
func (bc *Keycloak) GetUserIDByUsername(t TokenPayload, userName string) (string, error) {
	bc.l.Info().Str("userName", userName).Msg("Attempting to find user by userName")
	if t.Token == "" {
		return "", fmt.Errorf("missing client-token, please run Retrieve first")
	}

	uri := strings.Replace(bc.Endpoint, "/realms/", "/admin/realms/", 1) + "users?briefRepresentation=true&username=" + userName
	r, err := http.NewRequest(http.MethodGet, uri, nil)
	r.Header.Set("Authorization", "Bearer "+t.Token)
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
		return "", fmt.Errorf("Receeived forbidden as status-code for %s %d %s (%s). This might mean you are missing the 'query-users/manage-users'-role for the user '%s'. This is needed to look up users. You may also use userID directly. ", uri, res.StatusCode, res.Status, string(body), userName)
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
func (bc *Keycloak) Impersonate(clientToken TokenPayload, userID string) (TokenPayload, error) {
	var tp TokenPayload
	bc.Lock()
	defer bc.Unlock()

	bc.l.Info().Str("userID", userID).Msg("attempting to imppersonate user with impersonator-token")
	if clientToken.Token == "" {
		return tp, fmt.Errorf("clientToken not set. Please call Retrieve first")
	}
	if clientToken.Token == "" {
		return tp, fmt.Errorf("redirectUri is required for impersonation")
	}
	uri := strings.Replace(bc.Endpoint, "/realms/", "/admin/realms/", 1) + "users/" + userID + "/impersonation"

	r, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return tp, fmt.Errorf("failed to create http-request: %w", err)
	}
	r.Header.Set("Authorization", "Bearer "+clientToken.Token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return tp, fmt.Errorf("failed while performing http-request: %w", err)
	}
	contentType := res.Header.Get("Content-Type")
	defer res.Body.Close()
	if res.StatusCode > 299 {
		body, err := io.ReadAll(res.Body)
		additional := ""
		if err != nil {
			additional = fmt.Sprintf(". Additionally, reading the body failed: %s", err.Error())
		} else {
			additional = fmt.Sprintf("body: %s", string(body))
		}
		return tp, fmt.Errorf("Unexpected statuscode returned for endpoint %s: %d %s %s %s", uri, res.StatusCode, res.Status, contentType, additional)
	}
	cookies := res.Cookies()
	var validCookies []*http.Cookie
	now := time.Now()
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Expires.After(now) || cookies[i].RawExpires == "" {
			validCookies = append(validCookies, cookies[i])
		}
	}
	redirectUrl := bc.Endpoint + "protocol/openid-connect/auth?response_mode=fragment&response_type=token&client_id=" + bc.ClientID + "&redirect_uri=" + bc.RedirectUri
	reqRedirect, err := http.NewRequest(http.MethodGet, redirectUrl, nil)
	if err != nil {
		return tp, fmt.Errorf("failed to create redirect-request: %s %w", redirectUrl, err)
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
		return tp, fmt.Errorf("failed to perform redirect-request: %s %w", redirectUrl, err)
	}
	defer resRedirect.Body.Close()
	if resRedirect.StatusCode >= 400 {
		redirectBody, err := ioutil.ReadAll(resRedirect.Body)
		return tp, fmt.Errorf("Unexpected status-code for redirect-call: %d %s %s %s %w", resRedirect.StatusCode, resRedirect.Status, redirectUrl, string(redirectBody), err)
	}
	location := resRedirect.Header.Get("Location")
	if location == "" {
		return tp, fmt.Errorf("no location found in headers")
	}
	u, err := url.ParseQuery(location)
	if err != nil {
		return tp, fmt.Errorf("failed to parse location-uri: %w", err)
	}
	tokenStr := u.Get("access_token")
	if tokenStr == "" {
		var keys []string
		for k := range u {
			keys = append(keys, k)
		}
		return tp, fmt.Errorf("could not find 'access_token' in query-params: %v. This may be because the client does not have Implicit Flow enabled within keycloak", keys)
	}
	_tp, err := TokenStringToTokenPayload(tokenStr)
	if err != nil {
		return _tp, err
	}
	user := getPayloadUser(_tp.Raw)
	bc.l.Info().Str("user", user).Msg("Successfully performed impersonation")
	return _tp, nil

}
