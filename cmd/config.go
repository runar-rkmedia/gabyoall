package cmd

import (
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Url               string `cfg:"url" description:"The url to make requests to"`
	NoTokenValidation bool   `cfg:"no-token-validation" description:"If set, will skip validation of token"`
	PrintTable        bool   `cfg:"print-table" description:"If set, will print table while running"`
	AuthToken         string `cfg:"auth-token" description:"Set to use a token"`
	OperationName     string `cfg:"operation-name" description:"For Graphql, you may set an operation-name"`
	// TODO: implement
	Body string `cfg:"data" short:"d" description:"Data to include in requests."`
	// TODO: implement
	Header map[string]string `cfg:"header" short:"H" description:"Additional headers to include"`
	// TODO: implement
	Method string `cfg:"method" short:"X" description:"Http-method"`

	Query         string                 `cfg:"query" description:"For Graphql, you may set a query"`
	Variables     map[string]interface{} `cfg:"variables" description:"For Graphql, you may add variables"`
	LogLevel      string                 `cfg:"log-level" default:"info" description:"Log-level to use. Can be trace,debug,info,warn(ing),error or panic"`
	LogFormat     string                 `cfg:"log-format" default:"human" description:"Format of the logs. Can be human or json"`
	Output        string                 `cfg:"output" description:"File to output results to"`
	OkStatusCodes []int                  `cfg:"ok-status-codes" description:"list of status-codes to consider ok. If none is provided, any status-code within 200-299 is considered ok."`
	ResponseData  bool                   `cfg:"response-data" description:"Set to include response-data in output"`
	Mock          bool                   `cfg:"mock" description:"Enable to mock the requests."`
	Concurrency   int                    `cfg:"concurrency" description:"Amount of concurrent requests." default:"100" short:"c"`
	RequestCount  int                    `cfg:"request-count" default:"200" description:"Number of request to make total" short:"n"`
}

func GetConfig(l logger.AppLogger) *Config {
	var cfg Config
	viper.Unmarshal(&cfg)
	return &cfg
}
