package types

import "github.com/runar-rkmedia/gabyoall/cmd"

type DynamicAuth struct {
	Requests []DynamicRequest `json:"requests,omitempty"`
}

type Secret string
type Secrets map[string]Secret

var (
	Redacted = "**REDACTED**"
	// When marchsalling a Secret, the value will be replaced with this constant value (**REDACTED**)
	RedactedSecret = []byte(`"` + Redacted + `"`)
)

func IsRedacted(s string) bool {
	return s == Redacted
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return RedactedSecret, nil
}

type DynamicRequest struct {
	Method         string            `json:"method,omitempty"`
	Uri            string            `json:"uri,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	JsonRequest    bool              `json:"json_request,omitempty"`
	JsonResponse   bool              `json:"json_response,omitempty"`
	ResultJmesPath string            `json:"result_jmes_path,omitempty"`
	Body           interface{}       `json:"body,omitempty"`
}
type AuthConfig struct {
	// With dynamic, a series of requests can be made before the stress-test is performed.
	// The output can be piped into eachother, and then into the requests made during the stress-test.
	Dynamic DynamicAuth `json:"dynamic,omitempty"`
	// Bearer or Dynamic
	Kind string `json:"kind,omitempty"`
	// Used to impersonate (currently only works with keycloak)
	ImpersionationCredentials ImpersionationCredentials `json:"impersionation_credentials,omitempty"`
	// ClientID to use for
	ClientID     string `json:"client_id,omitempty"`
	RedirectUri  string `json:"redirect_uri,omitempty"`
	ClientSecret Secret `json:"client_secret,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`
	// Used with kind=Bearer and impersonation. Currenly, only keycloak is supported
	EndpointType string `json:"endpoint_type,omitempty"`
	// The header-key to use. Defaults to Authorization
	HeaderKey string `json:"header_key,omitempty"`
	Token     Secret `json:"token,omitEmpty"`
}
type ImpersionationCredentials struct {
	// Username to impersonate with. Needs to have the impersonation-role
	Username string `json:"username,omitempty"`
	Password Secret `json:"password,omitempty"`
	// UserID to impersonate as. This is preferred over UserNameToImpersonate
	UserIDToImpersonate string `json:"user_id_to_impersonate,omitempty"`
	// Will perform a lookup to get the ID of the username.
	UserNameToImpersonate string `json:"user_name_to_impersonate,omitempty"`
}
type Config struct {
	// Auth can be used to configure the authentication performed on a request.
	Auth *AuthConfig `json:"auth,omitempty"`
	// A list of http-status-codes to consider OK. Defaults to 200 and 204.
	OkStatusCodes *[]int `json:"ok_status_codes,omitempty"`
	// Whether or not Response-data should be stored.
	ResponseData *bool `json:"response_data,omitempty"`
	// Concurrency for the requests to be made
	Concurrency *int `json:"concurrency,omitempty"`
	// Number of requests to be performaed
	RequestCount *int     `json:"request_count,omitempty"`
	Secrets      *Secrets `json:"secrets,omitempty"`
}

// MergeWith will overwrite values with values in argument c.
func (c Config) MergeInto(config cmd.Config) cmd.Config {
	if c.Concurrency != nil && *c.Concurrency > 0 {
		config.Concurrency = *c.Concurrency
	}
	if c.RequestCount != nil && *c.RequestCount > 0 {
		config.RequestCount = *c.RequestCount
	}
	if c.OkStatusCodes != nil {
		config.OkStatusCodes = *c.OkStatusCodes
	}
	if c.ResponseData != nil {
		config.ResponseData = *c.ResponseData
	}
	if c.Auth != nil {
		if c.Auth.Endpoint != "" {
			config.Auth.Endpoint = c.Auth.Endpoint
		}
		if c.Auth.EndpointType != "" {
			config.Auth.EndpointType = c.Auth.EndpointType
		}
		if c.Auth.HeaderKey != "" {
			config.Auth.HeaderKey = c.Auth.HeaderKey
		}
		if c.Auth.Kind != "" {
			config.Auth.Kind = c.Auth.Kind
		}
		if c.Auth.RedirectUri != "" {
			config.Auth.RedirectUri = c.Auth.RedirectUri
		}
		if c.Auth.Token != "" {
			config.Auth.Token = string(c.Auth.Token)
		}
		if len(c.Auth.Dynamic.Requests) != 0 {
			for _, v := range c.Auth.Dynamic.Requests {
				config.Auth.Dynamic.Requests = make([]cmd.DynamicRequest, len(c.Auth.Dynamic.Requests))
				if vv := cmd.DynamicRequest(v); true {
					config.Auth.Dynamic.Requests = append(config.Auth.Dynamic.Requests, vv)
				}
			}
		}
		if c.Auth.ClientID != "" {
			config.Auth.ClientID = c.Auth.ClientID
		}
		if c.Auth.ClientSecret != "" && !IsRedacted(string(c.Auth.ClientSecret)) {
			config.Auth.ClientSecret = string(c.Auth.ClientSecret)
		}
		if c.Auth.ImpersionationCredentials.Password != "" && IsRedacted(string(c.Auth.ImpersionationCredentials.Password)) {
			config.Auth.ImpersionationCredentials.Password = string(c.Auth.ImpersionationCredentials.Password)
		}
		if c.Auth.ImpersionationCredentials.UserIDToImpersonate != "" {
			config.Auth.ImpersionationCredentials.UserIDToImpersonate = c.Auth.ImpersionationCredentials.UserIDToImpersonate
		}
		if c.Auth.ImpersionationCredentials.UserNameToImpersonate != "" {
			config.Auth.ImpersionationCredentials.UserNameToImpersonate = c.Auth.ImpersionationCredentials.UserNameToImpersonate
		}
		if c.Auth.ImpersionationCredentials.Username != "" {
			config.Auth.ImpersionationCredentials.Username = c.Auth.ImpersionationCredentials.Username
		}
	}
	return config
}
