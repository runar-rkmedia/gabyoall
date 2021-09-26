package main

import (
	"errors"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/runar-rkmedia/gabyoall/auth"
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
	// TODO: Refactor so this is a bit more general. (but still support graphql)
	var query = queries.GraphQLQuery{
		Body:          config.Body,
		Query:         config.Query,
		Variables:     config.Variables,
		OperationName: config.OperationName,
		Headers:       config.Header,
		Method:        config.Method,
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

	if query.Query == "" && query.Body == "" {
		l.Fatal().Interface("query", query).Msg("Missing query/body")
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
	bearerC := auth.NewBearerTokenCreator(
		logger.GetLogger("token-handler"),
		auth.BearerTokenCreatorOptions{
			ImpersionationCredentials: auth.ImpersionationCredentials{
				Username:              config.Auth.ImpersionationCredentials.Username,
				Password:              config.Auth.ImpersionationCredentials.Password,
				UserIDToImpersonate:   config.Auth.ImpersionationCredentials.UserIDToImpersonate,
				UserNameToImpersonate: config.Auth.ImpersionationCredentials.UserNameToImpersonate,
			},
			ClientID:     config.Auth.ClientID,
			RedirectUri:  config.Auth.RedirectUri,
			Endpoint:     config.Auth.Endpoint,
			Token:        config.AuthToken,
			EndpointType: config.Auth.EndpointType,
			ClientSecret: config.Auth.ClientSecret,
		})
	if config.Auth.ImpersionationCredentials.UserIDToImpersonate != "" || config.Auth.ImpersionationCredentials.UserNameToImpersonate != "" {
		err := bearerC.Retrieve()
		if err != nil {
			l.Fatal().Err(err).Msg("failed to retrieve token")
		}
	} else {

	}

	if token != "" && bearerC.Token == "" {
		bearerC.Token = token
	} else {
		token = bearerC.Token
	}
	// bearerC.Retrieve()
	payloadJson, err := bearerC.ParseToken()
	if err != nil {
		l.Fatal().Interface("payload", payloadJson).Err(err).Msg("Could not parse token")
	}
	err = bearerC.Validate()
	if err != nil {
		l.Fatal().Err(err).Msg("Token is not valid")
	}
	// os.Exit(1)

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
		&bearerC,
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
