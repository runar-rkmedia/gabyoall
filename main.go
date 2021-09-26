package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/printer"
	"github.com/runar-rkmedia/gabyoall/queries"
	"github.com/runar-rkmedia/gabyoall/utils"
	"github.com/runar-rkmedia/gabyoall/worker"
)

// TODO: Reread token every minute? In case of short-lived tokens.
// For instance, if they are read from env-variables, the user could set them.

type TemplateVars struct {
	cmd.Config
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
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

	if query.Query == "" {
		l.Fatal().Interface("query", query).Msg("Missing query")
	}
	vars := TemplateVars{*config}
	// TODO: Clean the path
	outputPath := config.Output
	if outputPath != "" {
		outputPath = utils.RunTemplating(l, config.Output, "output", vars)
		err := os.MkdirAll(path.Dir(outputPath), 0755)
		if err != nil {
			l.Fatal().Err(err).Str("dir", path.Dir(outputPath)).Msg("Failed to create directories for output")
		}
	}

	token := utils.RunTemplating(l, config.AuthToken, "token", vars)

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
	if payload.GraphqlEndpoint != "" && !utils.UrlsAreEqual(payload.GraphqlEndpoint, config.Url) {

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
	endpoint := cmd.NewEndpoint(logger.GetLogger("gql"), config.Url)
	endpoint.Headers.Add("Authorization", token)

	l.Info().Str("url", config.Url).Str("operationName", query.OperationName).Int("count", config.RequestCount).Int("paralism", config.Concurrency).Msg("Running requests with paralism")
	SetupCloseHandler(func(signal os.Signal) {
		out.Write()
	})

	wt := worker.WorkThing{}
	ch := wt.Run(endpoint, *config, query)

	successes := 0
	startTime := time.Now()
	print := printer.NewPrinter(
		*config,
		payload,
		query.OperationName,
		out,
		startTime,
	)
	print.Animate()

	for i := 0; i < config.RequestCount; i++ {
		stat := <-ch
		if stat.ErrorType == "" {
			successes++
		}
		out.AddStat(stat)
		print.Update(i, successes)

	}
	out.CalculateStats()

	print.Complete(config.RequestCount, successes)
	err = out.Write()
	if err != nil {
		l.Fatal().Err(errors.Unwrap(err)).Msg("Failed to write output")
	}
	l.Info().Msg("All done")
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

type GraphqlRequest struct {
	cmd.Endpoint
	queries.GraphQLQuery
}
