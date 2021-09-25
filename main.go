package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"syscall"
	"text/template"
	"time"

	spin "github.com/tj/go-spin"

	"github.com/Masterminds/sprig"
	tm "github.com/buger/goterm"
	"github.com/rs/zerolog/log"

	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/queries"
)

var envRegex = regexp.MustCompile(`(\$({([^}]*)}))`)

// TODO: Reread token every minute? In case of short-lived tokens.
// For instance, if they are read from env-variables, the user could set them.

func expandEnv(s string) string {
	envRegex.FindAllStringSubmatch(s, -1)
	return envRegex.ReplaceAllStringFunc(s, func(str string) string {
		_envKey := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(str, "$"), "{"), "}")
		split := strings.Split(_envKey, ":")
		def := ""
		envKey := split[0]
		if len(split) > 1 {
			def = split[1]
		}

		val := os.Getenv(envKey)
		if val == "" {
			return def
		}
		return val
	})
}
func executeTemplate(l logger.AppLogger, templateString, name string, vars interface{}) string {
	buf := new(bytes.Buffer)
	t := template.New(name)
	t.Funcs(sprig.TxtFuncMap())
	tmpl, err := t.Parse(templateString)
	if err != nil {
		l.Error().Err(err).Str("templateString", templateString).Str("name", name).Msg("Failed to parse templateString to template")
	}

	err = tmpl.Execute(buf, vars)
	if err != nil {
		l.Error().Err(err).Str("templateString", templateString).Str("name", name).Msg("Failed to execute templateString to template")
	}
	return buf.String()
}

func runTemplating(l logger.AppLogger, templateString, name string, vars interface{}) string {
	templateString = expandEnv(templateString)
	return executeTemplate(l, templateString, name, vars)
}

func prepareUrl(s string) string {
	if !strings.HasPrefix(s, "http") {
		s = "https://" + s
	}
	s = strings.TrimSuffix(s, "/")
	return s
}

func urlsAreEqual(a, b string) bool {
	A := prepareUrl(a)
	B := prepareUrl(b)
	return A == B
}

func prettyDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%03dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%02.1fs", d.Seconds())
	}
	if d < time.Hour {
		_s := d.Seconds()
		m := math.Floor(_s / 60)
		s := _s - (m * 60)
		return fmt.Sprintf("%02.0f:%02.0f", m, s)
	}
	return d.String()

}

type TemplateVars struct {
	cmd.Config
}

func main() {
	cmd.Execute()
	config := cmd.GetConfig(logger.GetLogger("initial"))
	var query = queries.GraphQLQuery{
		Query:         config.Query,
		Variables:     config.Variables,
		OperationName: config.OperationName,
	}
	logger.InitLogger(logger.LogConfig{
		Level:      config.LogLevel,
		Format:     config.LogFormat,
		WithCaller: true,
	})
	l := logger.GetLogger("main")
	if config.RequestCount == 0 {
		l.Fatal().Msg("Request-count cannot be 0")

	}
	var okStatusCodes []int

	if query.Query == "" {
		l.Fatal().Interface("query", query).Msg("Missing query")
	}
	vars := TemplateVars{*config}
	// TODO: Clean the path
	outputPath := runTemplating(l, config.Output, "output", vars)
	err := os.MkdirAll(path.Dir(outputPath), 0755)
	if err != nil {
		l.Fatal().Err(err).Str("dir", path.Dir(outputPath)).Msg("Failed to create directories for output")
	}

	token := runTemplating(l, config.AuthToken, "token", vars)

	if token == "" {
		l.Fatal().Msg("Token is requried.")
	}
	// To copy straight from firefox
	token = strings.TrimPrefix(token, "Authorization: ")
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}
	rawToken := strings.TrimPrefix(token, "Bearer ")
	tokenSplitted := strings.Split(rawToken, ".")
	if len(tokenSplitted) != 3 {
		l.Fatal().Str("rawToken", rawToken).Msg("Token looks invalid. Expected to find JWT with 3 parts seperated by '.'")
	}
	data := tokenSplitted[1]
	payloadRaw, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		l.Fatal().Err(err).Str("raw", tokenSplitted[1]).Msg("failed to base64-decode payload")
	}
	var payload cmd.JwtPayload
	err = json.Unmarshal([]byte(payloadRaw), &payload)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to decode payload in token")
	}
	var payloadJson map[string]interface{}
	err = json.Unmarshal([]byte(payloadRaw), &payloadJson)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to decode payload in token")
	}
	if payload.GraphqlEndpoint != "" && !urlsAreEqual(payload.GraphqlEndpoint, config.Url) {

		// l.Fatal().Str("url", *url).Interface("payload", payloadJson).Str("graphqlEndpoint", payload.GraphqlEndpoint).Msg("Graphql-endpoint from token does not match match url")
	}
	pex := time.Unix(payload.Exp, 0)
	payload.ExpiresAt = &pex
	l.Info().Str("exp", payload.ExpiresAt.String()).Str("in", payload.ExpiresAt.Sub(time.Now()).String()).Msg("Token expires at {{exp}} in {{in}} ")
	if config.NoTokenValidation != true && payload.ExpiresAt.Before(time.Now()) {
		l.Fatal().Str("url", config.Url).Interface("payload", payloadJson).Msg("Token expired. Please retreive a new one.")
	}
	out, err := cmd.NewOutput(l, outputPath, config.Url, query, payloadJson)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to set up output")
	}
	l.Info().Str("path", out.GetPath()).Msg("Will write output to path:")
	endpoint := cmd.NewGraphQLEndpoint(logger.GetLogger("gql"), config.Url)
	endpoint.Headers.Add("Authorization", token)
	ch := make(chan cmd.RequestStat)

	hasWorkChan := make(chan struct{}, config.Concurrency)
	workChan := make(chan Work, config.Concurrency)
	l.Info().Str("url", config.Url).Str("operationName", query.OperationName).Int("count", config.RequestCount).Int("paralism", config.Concurrency).Msg("Running requests with paralism")
	SetupCloseHandler(func(signal os.Signal) {
		out.Write()
	})

	go func() {
		for {
			select {
			case work := <-workChan:
				go work()
				hasWorkChan <- struct{}{}

			}

		}
	}()
	go func() {
		for j := 0; j < config.RequestCount; j++ {

			workChan <- func(i int) Work {

				return func() {
					if config.Mock {
						// TODO: replace with a mocked http-client-interface
						stat := cmd.NewStat()
						time.Sleep(time.Millisecond * time.Duration(rand.Int63n(80)+1))
						errorType := cmd.Unknwon
						n := rand.Intn(7)
						switch n {
						case 1:
							errorType = cmd.NonOK
						case 2:
							errorType = cmd.ServerTestError
						case 3:
							errorType = "RandomErr"
						case 4:
							errorType = "OtherErr"
						case 6:
							errorType = "MadeUpError"

						}
						ch <- stat.End(errorType, nil)
						return
					}
					_, stat, _ := endpoint.RunQuery(query, okStatusCodes)
					if stat.Response != nil {

						if config.ResponseData != false && stat.Response != nil && stat.Response["error"] == nil && stat.Response["data"] != nil {
							delete(stat.Response, "data") //stat.Response[data]
						}
					}
					ch <- stat

				}
			}(j)
		}
	}()
	i := 0
	successes := 0
	startTime := time.Now()
	quitSpinner := make(chan struct{})
	shouldSpin := true
	if shouldSpin {
		spinner := spin.New()
		go func(quit chan struct{}) {
			for {
				select {
				case <-quit:
					return
				default:
					payloadExp := ""
					if config.NoTokenValidation != true && payload.ExpiresAt != nil {
						payloadExp = fmt.Sprintf("Token expires: %s", payload.ExpiresAt.Sub(time.Now()).String())
					}
					fraction := float64(i) / float64(config.RequestCount)
					dur := time.Now().Sub(startTime)
					estimatedCompletion := time.Duration(float64(dur)/fraction) - dur
					// TODO: sync these values. This will likely crash.
					fails := ""
					failures := i - successes
					if failures > 0 {
						fails = fmt.Sprintf("\033[31m[%d (%.2f%%)\033[0m", failures, float64(failures)/float64(i)*100)
					}
					fmt.Printf("\r\033[36m[%d/%d (%.2f%%) %s -c=%d] %s Waiting for result from: %s (%s) \033[m %s (%s) %s", i, config.RequestCount, fraction*100, fails, config.Concurrency, spinner.Next(), config.Url, query.OperationName, prettyDuration(dur), prettyDuration(estimatedCompletion), payloadExp)
					time.Sleep(300 * time.Millisecond)
				}
			}
		}(quitSpinner)
	}
	var lastOut time.Time
outer:
	for {
		select {
		case stat := <-ch:
			// bar.SetCurrent(int64(i))
			i++
			if stat.ErrorType == "" {
				successes++
			}
			if config.PrintTable && time.Now().Sub(lastOut) > 500*time.Millisecond {
				if shouldSpin {
					shouldSpin = false
					go func() { quitSpinner <- struct{}{} }()
				}

				tm.Clear()
				out.PrintTable()
				now := time.Now()
				tm.Printf("\nFinished %d of %d (%.2f%%) %s (%s) \n", i, config.RequestCount, float64(i)/float64(config.RequestCount)*100, now.Format("15:04:05.0"), query.OperationName)

				tm.Flush()
				lastOut = time.Now()
			}
			out.AddStat(stat)
			select {
			case <-hasWorkChan:
			default:
			}

			if i >= config.RequestCount {
				out.CalculateStats()
				break outer
			}
		default:

		}
	}

	if config.PrintTable {
		tm.Clear()
		out.PrintTable()
		now := time.Now()
		tm.Printf("\nFinished %d of %d (%.2f%%) %s \n", i, config.RequestCount, float64(i)/float64(config.RequestCount)*100, now.Format("15:04:05.0"))

		tm.Flush()
		lastOut = time.Now()
	}
	fmt.Println("")
	err = out.Write()
	if err != nil {
		l.Fatal().Err(errors.Unwrap(err)).Msg("Failed to write output")
	}
	log.Info().Interface("stat", out.Stats).Msg("All done")
}

func SetupCloseHandler(f func(signal os.Signal)) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		signal := <-c
		log.Warn().Msg("- Ctrl+C pressed in Terminal")
		f(signal)
		os.Exit(0)
	}()
}

type Work func()

type GraphqlRequest struct {
	cmd.GraphQlEndpoint
	queries.GraphQLQuery
}
