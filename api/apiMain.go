package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	_ "embed"

	"github.com/go-playground/validator/v10"
	"github.com/runar-rkmedia/gabyoall/api/bboltStorage"
	_ "github.com/runar-rkmedia/gabyoall/api/docs"
	"github.com/runar-rkmedia/gabyoall/api/requestContext"
	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/logger"
)

var (
	//go:embed generated-swagger.yml
	swaggerYml string
)

type ApiConfig struct {
	Address string
	Port    int
	logger.LogConfig
}

//go:generate swagger generate spec -o generated-swagger.yml
func main() {
	cfg := ApiConfig{
		Address: "0.0.0.0",
		Port:    80,
		LogConfig: logger.LogConfig{
			Level:      "debug",
			Format:     "human",
			WithCaller: true,
		},
	}
	logger.InitLogger(cfg.LogConfig)
	l := logger.GetLogger("main")
	db, err := bboltStorage.NewBbolt(l, "db.bbolt")
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to initialize storage")
	}
	ctx := requestContext.Context{
		l,
		&db,
		validator.New(),
	}
	address := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	// handler := http.NewServeMux()
	handler := http.StripPrefix("/api/", EndpointsHandler(ctx))
	l.Info().Str("address", cfg.Address).Int("port", cfg.Port).Msg("Creating listeners")
	err = http.ListenAndServe(address, handler)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create listener")
	}
}

func EndpointsHandler(ctx requestContext.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rc := requestContext.NewReqContext(&ctx, r, rw)
		var body []byte
		var err error
		isGet := r.Method == http.MethodGet
		isPost := r.Method == http.MethodPost
		path := r.URL.Path
		paths := strings.Split(path, "/")

		if rc.ContentKind > 0 && isPost {
			body, err = ioutil.ReadAll(r.Body)
			if err != nil {
				rc.WriteErr(err, requestContext.CodeErrReadBody)
			}
		}

		switch paths[0] {
		case "swagger", "swagger.yaml", "swagger.yml":
			rw.Header().Set("Content-Type", "text/vnd.yaml")
			rw.Header().Set("Content-Disposition", `attachment; filename="swagger-gobyoall.yaml"`)
			rw.Write([]byte(swaggerYml))
			return
		case "endpoint":
			// Create endpoint
			if isPost && paths[1] == "" {
				var input types.EndpointPayload
				err := rc.Unmarshal(body, &input)
				if err != nil {
					rc.WriteErr(err, requestContext.CodeErrMarhal)
					return
				}
				err = rc.ValidateStruct(input)
				if err != nil {
					rc.WriteErr(err, requestContext.CodeErrInputValidation)
					return
				}
				e, err := ctx.DB.CreateEndpoint(input)
				rc.WriteAuto(e, err, requestContext.CodeErrDBCreateEndpoint)
				return

			}
			// List endpoints
			if isGet && len(paths) == 1 {
				es, err := ctx.DB.Endpoints()
				rc.WriteAuto(es, err, requestContext.CodeErrEndpoint)
				return
			}
			// Get endpoint
			if isGet && len(paths) == 2 {
				es, err := ctx.DB.Endpoint(paths[1])
				rc.WriteAuto(es, err, requestContext.CodeErrEndpoint)
				return
			}
		}

		rc.WriteError("No route here", requestContext.CodeErrNoRoute)
	}
}

type OkResponse struct {
	Ok bool `json:"ok"`
}
type CreateResponse struct {
	Ok bool   `json:"ok"`
	ID string `json:"id"`
}

var (
	okResponse = OkResponse{true}
)

var idRegex = regexp.MustCompile(`^[a-zA-Z0-9-_]{6,40}$`)

func validateIDInput(id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return id, requestContext.ErrIDEmpty
	}
	if len(id) > 40 {
		return id, requestContext.ErrIDTooLong
	}
	if !idRegex.MatchString(id) {
		return id, requestContext.ErrIDNonValid
	}
	return id, nil
}

type GeneralResponse struct {
	Ok bool
}

func NewGeneralResponse(ok bool) GeneralResponse {
	return NewGeneralResponse(ok)
}
