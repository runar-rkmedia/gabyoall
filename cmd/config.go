package cmd

import (
	"log"
	"os"

	"github.com/runar-rkmedia/gabyoall/logger"
)

type Config struct {
	Url               string
	NoTokenValidation *bool
	Token             string
	OperationName     string
	Query             string
	LogLevel          string
	Output            string
	OkStatusCodes     []int
	ResponseData      *bool
	Concurrency       int
	Count             int
	Variables         *map[string]interface{}
}

var c = map[string]*Config{}

func GetConfig(l logger.AppLogger, path string) *Config {

	if c[path] != nil {
		return c[path]
	}
	var conf Config
	err := ReadYamlFile(path, &conf)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Fatalf("failed to get config %v", err)
	}
	c[path] = &conf

	return c[path]
}
