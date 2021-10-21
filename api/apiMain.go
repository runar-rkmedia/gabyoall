package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "embed"

	"github.com/NYTimes/gziphandler"
	"github.com/go-playground/validator/v10"
	"github.com/runar-rkmedia/gabyoall/api/bboltStorage"
	_ "github.com/runar-rkmedia/gabyoall/api/docs"
	"github.com/runar-rkmedia/gabyoall/api/requestContext"
	"github.com/runar-rkmedia/gabyoall/api/scheduler"
	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/frontend"
	"github.com/runar-rkmedia/gabyoall/logger"
)

var (
	//go:embed generated-swagger.yml
	swaggerYml string
	// These are added at build...
	Version      string
	BuildDateStr string
	BuildDate    time.Time
	GitHash      string
	isDev        = true
	IsDevStr     = "1"

	serverStartTime = time.Now()
)

func init() {
	if BuildDateStr != "" {
		t, err := time.Parse("2006-01-02T15:04:05", BuildDateStr)
		if err != nil {
			panic(fmt.Errorf("Failed to parse build-date: %w", err))
		}
		BuildDate = t
	}
	if IsDevStr != "1" {
		isDev = false
	}
}

type ApiConfig struct {
	Address      string
	RedirectPort int
	Port         int
	CertFile     string
	CertKey      string
	logger.LogConfig
}

//go:generate swagger generate spec -o generated-swagger.yml
//go:generate sh -c "cd ../frontend && yarn gen"
func main() {
	err := cmd.InitConfig()
	if err != nil {
		panic(err)
	}
	cfg := ApiConfig{
		Address:      "0.0.0.0",
		Port:         443,
		RedirectPort: 80,
		CertFile:     "server.crt",
		CertKey:      "server.key",
		LogConfig: logger.LogConfig{
			Level:      "info",
			Format:     "human",
			WithCaller: true,
		},
	}
	logger.InitLogger(cfg.LogConfig)
	l := logger.GetLogger("main")
	l.Info().Str("version", Version).Time("buildDate", BuildDate).Time("buildDateLocal", BuildDate.Local()).Str("gitHash", GitHash).Msg("Starting")
	db, err := bboltStorage.NewBbolt(l, "db.bbolt")
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to initialize storage")
	}
	ctx := requestContext.Context{
		L:               l,
		DB:              &db,
		StructValidater: validator.New(),
	}

	s := scheduler.NewScheduler(l, &db, cmd.GetConfig(l))

	s.Run()
	address := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	handler := http.NewServeMux()
	// http.Handle("/api/", http.StripPrefix("/api/", EndpointsHandler(ctx)))
	handler.Handle("/api/",
		gziphandler.GzipHandler(
			http.StripPrefix("/api/", EndpointsHandler(ctx))))

	useCert := false
	if cfg.CertFile != "" {
		_, err := os.Stat(cfg.CertFile)
		if err == nil {
			useCert = true
		}

	}

	if isDev {
		// In development, we serve the file directly.
		handler.Handle("/", http.FileServer(http.Dir("./frontend/dist/")))
	} else {
		handler.Handle("/", frontend.DistServer)
	}
	l.Info().Str("address", cfg.Address).Int("port", cfg.Port).Bool("redirectHttpToHttps", useCert && cfg.RedirectPort != 0).Bool("tls", useCert).Msg("Creating listener")
	srv := http.Server{Addr: address, Handler: handler}
	if useCert {
		// TODO: re-read the certificate before it expires.
		if cfg.RedirectPort != 0 {
			redirectTLS := func(w http.ResponseWriter, r *http.Request) {
				newAddress := "https://" + r.Host
				if cfg.Port != 443 {
					newAddress += fmt.Sprintf(":%d", cfg.Port)
				}
				http.Redirect(w, r, newAddress+r.RequestURI, http.StatusMovedPermanently)
			}
			go func() {
				redirectAddress := fmt.Sprintf("%s:%d", cfg.Address, cfg.RedirectPort)
				if err := http.ListenAndServe(redirectAddress, http.HandlerFunc(redirectTLS)); err != nil {
					l.Fatal().Err(err).Str("redirectAddress", redirectAddress).Msg("Failed to create redirect-listener")

				}
			}()

		}
		err = srv.ListenAndServeTLS("server.crt", "server.key")
	} else {
		srv.ListenAndServe()
	}
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create listener")
	}

}

type AccessControl struct {
	AllowOrigin string
	MaxAge      time.Duration
}

var (
	accessControl = AccessControl{
		AllowOrigin: "_any_",
		MaxAge:      24 * time.Hour,
	}
	pingByte = []byte{}
)

func EndpointsHandler(ctx requestContext.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "ping" {
			rw.Write(pingByte)
			return
		}
		h := rw.Header()
		switch accessControl.AllowOrigin {
		case "_any_":
			h.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		default:
			h.Set("Access-Control-Allow-Origin", accessControl.AllowOrigin)
		}
		h.Set("Access-Control-Allow-Headers", "x-request-id, content-type, jmes-path")
		h.Set("Access-Control-Max-Age", fmt.Sprintf("%0.f", accessControl.MaxAge.Seconds()))
		if r.Method == "OPTIONS" {
			h.Set("Cache-Control", fmt.Sprintf("public, max-age=%0.f", accessControl.MaxAge.Seconds()))
			h.Set("Vary", "origin")

			return
		}
		rc := requestContext.NewReqContext(&ctx, r, rw)
		var body []byte
		var err error
		isGet := r.Method == http.MethodGet
		isPost := r.Method == http.MethodPost
		isPut := r.Method == http.MethodPut
		path := r.URL.Path
		paths := strings.Split(strings.TrimSuffix(path, "/"), "/")

		if rc.ContentKind > 0 && (isPost || isPut) {
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
		case "serverInfo":
			if isGet && len(paths) == 1 {
				info := types.ServerInfo{
					ServerStartedAt: serverStartTime,
					GitHash:         GitHash,
					Version:         Version,
					BuildDate:       BuildDate,
				}
				rc.WriteAuto(info, err, "serverInfo")
				return
			}
		case "endpoint":
			// Create endpoint
			if isPost && len(paths) == 1 {
				var input types.EndpointPayload
				if err := rc.ValidateBytes(body, &input); err != nil {
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
		case "request":
			// Create request
			if isPost && len(paths) == 1 {
				var input types.RequestPayload
				if err := rc.ValidateBytes(body, &input); err != nil {
					return
				}
				e, err := ctx.DB.CreateRequest(input)
				rc.WriteAuto(e, err, requestContext.CodeErrDBCreateRequest)
				return
			}
			// List requests
			if isGet && len(paths) == 1 {
				es, err := ctx.DB.Requests()
				rc.WriteAuto(es, err, requestContext.CodeErrRequest)
				return
			}
			// Get request
			if isGet && len(paths) == 2 {
				es, err := ctx.DB.Request(paths[1])
				rc.WriteAuto(es, err, requestContext.CodeErrRequest)
				return
			}
		case "stat":
			// List stats
			if isGet && len(paths) == 1 {
				es, err := ctx.DB.CompactStats()
				rc.WriteAuto(es, err, requestContext.CodeErrRequest)
				return
			}
		case "schedule":
			// Create schedule
			if isPost && len(paths) == 1 {
				var input types.SchedulePayload
				if err := rc.ValidateBytes(body, &input); err != nil {
					return
				}
				e, err := ctx.DB.CreateSchedule(input)
				rc.WriteAuto(e, err, requestContext.CodeErrDBCreateSchedule)
				return
			}
			// Update schedule
			if isPut && len(paths) == 2 {
				var input types.SchedulePayload
				if err := rc.ValidateBytes(body, &input); err != nil {
					return
				}
				e, err := ctx.DB.UpdateSchedule(paths[1], types.Schedule{SchedulePayload: input})
				rc.WriteAuto(e, err, requestContext.CodeErrDBCreateSchedule)
				return
			}
			// List schedules
			if isGet && len(paths) == 1 {
				es, err := ctx.DB.Schedules()
				rc.WriteAuto(es, err, requestContext.CodeErrSchedule)
				return
			}
			// Get schedule
			if isGet && len(paths) == 2 {
				es, err := ctx.DB.Schedule(paths[1])
				rc.WriteAuto(es, err, requestContext.CodeErrSchedule)
				return
			}
		}
		// http.FileServer(frontend.StaticFiles).ServeHTTP(rc.Rw, rc.rw)

		rc.WriteError(fmt.Sprintf("No route registerd for: %s %s", r.Method, r.URL.Path), requestContext.CodeErrNoRoute)
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
