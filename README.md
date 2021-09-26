# Gobyoall

Opinionated stress-tester for servers.

Features:
 - Concurrent requests
 - Reads config-files as well as cli-arguments (viper/cobra). Env-variables are also supported.
   - Makes it a lot easier to create multiple config-files for different requests and running them at a schedule
   - Store templates for requests in seperate files for reuse.
 - Outputs statistics as well as detailed result of each requests.
   - Output-path can be configured with go-templating for various needs.
 - Categorizes request-errors into buckets.
 - Integrates with GraphQL.

```
Flags:
      --auth-token string       Set to use a token
  -c, --concurrency int         Amount of concurrent requests. (default 100)
      --config string           config file (default is $HOME/.config/gobyoall-conf.yaml)
  -d, --data string             Data to include in requests.
  -H, --header stringToString   Additional headers to include (default [])
  -h, --help                    help for gobyoall
      --log-format string       Format of the logs. Can be human or json (default "human")
      --log-level string        Log-level to use. Can be trace,debug,info,warn(ing),error or panic (default "info")
  -X, --method string           Http-method
      --mock                    Enable to mock the requests.
      --no-token-validation     If set, will skip validation of token
      --ok-status-codes ints    list of status-codes to consider ok. If none is provided, any status-code within 200-299 is considered ok.
      --operation-name string   For Graphql, you may set an operation-name
      --output string           File to output results to
      --print-table             If set, will print table while running
      --query string            For Graphql, you may set a query
  -n, --request-count int       Number of request to make total (default 200)
      --response-data           Set to include response-data in output
      --url string              The url to make requests to
```

Example:

```
gobyoall \
  --query 'query{me}' \
  --auth-token "eyJ..." \
  --url example.com
```

## Install

```
go install github.com/runar-rkmedia/gabyoall
```

## Environment-variables

Environment-variables are all prefixed with `GOBYOALL_`. For instance, to set `auth-token` use `GOBYOALL_AUTHTOKEN`.

## JWT-tokens

Tokens can be either set manually, but it is vary practical to have the gobyoall generate them for you via impersonation. This lets you automate more stress-tests without the hassle of handling the tokens yourself, or opening up your server-application with special logic.

To set manually, see the help-section.

### Impersonation

> Impersonation currenly only works with Keycloak.

#### Keycloak

To impersonate a user, you need:

- A keycloak client with:
  - Direct Access Grants enabled
  - Implicit Flow Enabled
  - A Valid Redirect URI. The exact URI does not matter, as it will not be followed, e.g. the call wont be made.
- A user with:
  - Role Mapping: `impersonation` (under realm-managent)
  - (Optional) Role Mapping: manage-users. (under realm-management)

You then need to add these items to the config. Env-variables can also be used.

```yaml
auth:
  type: Bearer
  endpoint: https://keycloak-server.com/auth/realms/example/
  endpointType: keycloak
  # The exact value does not matter, as the redirect will never be followed.
  redirectUri: http://localhost:3000 
  clientId: test-client
  clientSecret: 4ac...8
  impersionationCredentials: 
    username: test
    password: test
    # prefer to set the userID over UserNameToImpersonate as it does not require a lookup
    userIDToImpersonate: 638492ff-282e-4ccd-8e4c-f65db4093d12
#     userNameToImpersonate: johndoe
```