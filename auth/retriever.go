package auth

import (
	"fmt"
	"strings"

	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/logger"
)

type ValidityStringer interface {
	PrintValidity() (string, error)
}

func Retrieve(l logger.AppLogger, a cmd.AuthConfig) (err error, token string, payload *TokenPayload, validityStringer ValidityStringer) {
	switch strings.ToLower(a.Kind) {
	case "dynamic":
		if len(a.Dynamic.Requests) == 0 {
			err = fmt.Errorf("With auth.kind set to dynamic, at least one request must be added")
			l.Error().Err(err).Msg("failed during Retrieve.dynamic")
			return
		}
		da := DynamicAuth{
			Requests:  make([]DynamicRequest, len(a.Dynamic.Requests)),
			HeaderKey: a.HeaderKey,
		}
		// This is not ideal, but it keeps the config from importing DynamicRequests, and vice verca.
		for i := 0; i < len(a.Dynamic.Requests); i++ {
			da.Requests[i] = DynamicRequest(a.Dynamic.Requests[i])
		}
		res, err := da.Retrieve()
		if err != nil {
			l.Error().Err(err).Msg("Failed during dynamic-auth-request")
			return err, token, payload, validityStringer
		}
		token = res.Token
	case "bearer":
		if a.ClientID != "" {
			bearerC := NewBearerTokenCreator(
				logger.GetLogger("token-handler"),
				BearerTokenCreatorOptions{
					ImpersionationCredentials: ImpersionationCredentials{
						Username:              a.ImpersionationCredentials.Username,
						Password:              a.ImpersionationCredentials.Password,
						UserIDToImpersonate:   a.ImpersionationCredentials.UserIDToImpersonate,
						UserNameToImpersonate: a.ImpersionationCredentials.UserNameToImpersonate,
					},
					ClientID:     a.ClientID,
					RedirectUri:  a.RedirectUri,
					Endpoint:     a.Endpoint,
					EndpointType: a.EndpointType,
					ClientSecret: a.ClientSecret,
				})
			validityStringer = &bearerC
			if a.ImpersionationCredentials.UserIDToImpersonate != "" || a.ImpersionationCredentials.UserNameToImpersonate != "" {

				tokenPayload, err := bearerC.Impersionate(ImpersionateOptions{
					UserName: a.UserNameToImpersonate,
					UserID:   a.UserIDToImpersonate,
				})
				if err != nil {
					l.Error().Err(err).Msg("failed to retrieve token")
					return err, token, payload, validityStringer
				}
				token = tokenPayload.Token
				payload = &tokenPayload
			}
		}
	case "":
		l.Warn().Msg("No auth.kind set, not using any authentication")
	default:
		l.Debug().Str("auth.kind", a.Kind).Msg("unrecognized option.")
	}
	return
}
