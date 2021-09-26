package auth

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
