package cmd

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/spf13/viper"
)

type AuthConfig struct {
	Dynamic DynamicAuth
	Kind    string
	ImpersionationCredentials
	ClientID     string
	RedirectUri  string
	ClientSecret string
	Endpoint     string
	EndpointType string
	Token        string
	Expires      *time.Time
	HeaderKey    string
	Payload      map[string]interface{}
}
type ImpersionationCredentials struct {
	// Username to impersonate with. Needs to have the impersonation-role
	Username,
	Password,
	// UserID to impersonate as. This is preferred over UserNameToImpersonate
	UserIDToImpersonate string
	// Will perform a lookup to get the ID of the username.
	UserNameToImpersonate string
}
type Config struct {
	Auth              AuthConfig             `cfg:"-"`
	Url               string                 `cfg:"url" description:"The url to make requests to"`
	NoTokenValidation bool                   `cfg:"no-token-validation" description:"If set, will skip validation of token"`
	PrintTable        bool                   `cfg:"print-table" description:"If set, will print table while running"`
	AuthToken         string                 `cfg:"auth-token" description:"Set to use a token"`
	OperationName     string                 `cfg:"operation-name" description:"For Graphql, you may set an operation-name"`
	Body              interface{}            `cfg:"data" short:"d" description:"Data to include in requests."`
	Header            map[string]string      `cfg:"header" short:"H" description:"Additional headers to include"`
	Method            string                 `cfg:"method" short:"X" description:"Http-method"`
	Query             string                 `cfg:"query" description:"For Graphql, you may set a query"`
	Variables         map[string]interface{} `cfg:"variables" description:"For Graphql, you may add variables"`
	LogLevel          string                 `cfg:"log-level" default:"info" description:"Log-level to use. Can be trace,debug,info,warn(ing),error or panic"`
	LogFormat         string                 `cfg:"log-format" default:"human" description:"Format of the logs. Can be human or json"`
	Output            string                 `cfg:"output" description:"File to output results to"`
	OkStatusCodes     []int                  `cfg:"ok-status-codes" description:"list of status-codes to consider ok. If none is provided, any status-code within 200-299 is considered ok."`
	ResponseData      bool                   `cfg:"response-data" description:"Set to include response-data in output"`
	Mock              bool                   `cfg:"mock" description:"Enable to mock the requests."`
	Concurrency       int                    `cfg:"concurrency" description:"Amount of concurrent requests." default:"100" short:"c"`
	RequestCount      int                    `cfg:"request-count" default:"200" description:"Number of request to make total" short:"n"`
	Api               ApiConfig              `cfg:"api" description:"Used with the api-server"`
}

type ApiConfig struct {
	Address      string `cfg:"address" default:"0.0.0.0" description:"Address (interface) to listen to)"`
	RedirectPort int    `cfg:"redirect-port" default:"80" description:"Used normally to redirect from http to https. Will be ignored if zero or same as listening-port"`
	Port         int    `cfg:"port" default:"80" description:"Port to listen to"`
	CertFile     string `cfg:"cert-file" default:"" description:"Number of request to make total"`
	CertKey      string `cfg:"cert-key" default:"" description:"Number of request to make total"`
	DBLocation   string `cfg:"db-path" default:"./storage/db.bbolt" description:"Filepath to where to store the database"`
}

type DynamicAuth struct {
	Requests []DynamicRequest
}

type DynamicRequest struct {
	Method         string
	Uri            string
	Headers        map[string]string
	JsonRequest    bool
	JsonResponse   bool
	ResultJmesPath string
	Body           interface{}
}

func GetConfig(l logger.AppLogger) *Config {
	var cfg Config
	viper.Unmarshal(&cfg)
	if cfg.Concurrency > cfg.RequestCount {
		cfg.Concurrency = cfg.RequestCount
	}
	l.Info().Str("config-file", viper.ConfigFileUsed()).Msg("Using config-file")
	return &cfg
}

func InitConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		viper.SetConfigName("gobyoall-conf")
		viper.AddConfigPath(path.Join(home, "gobyall"))
		viper.AddConfigPath(path.Join(home, ".config", "gobyall"))
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix("gobyoall")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			return fmt.Errorf("Fatal error config file: %w \n", err)
		}
	}
	return nil
}
