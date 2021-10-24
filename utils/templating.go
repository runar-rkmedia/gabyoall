package utils

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/jmespath/go-jmespath"
	"github.com/runar-rkmedia/gabyoall/logger"
)

var envRegex = regexp.MustCompile(`(\$({([^}]*)}))`)
var funcs = sprig.TxtFuncMap()

func init() {
	funcs["jmes"] = func(jmesPath string, obj interface{}) string {
		result, err := jmespath.Search(jmesPath, obj)
		if err != nil {
			return fmt.Sprintf("failed in jmes-path '%s' %s for %#v", jmesPath, err.Error(), obj)
		}
		switch result.(type) {
		case string:
			return result.(string)
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("%d", result.(int))
		case float64, float32:
			return fmt.Sprintf("%f", result.(float64))
		}
		return fmt.Sprintf("%v", result)
	}
}

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
	t.Funcs(funcs)
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

func RunTemplating(l logger.AppLogger, templateString, name string, vars interface{}) string {
	if templateString == "" {
		return ""
	}
	templateString = expandEnv(templateString)
	return executeTemplate(l, templateString, name, vars)
}
