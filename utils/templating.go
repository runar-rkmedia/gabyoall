package utils

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/runar-rkmedia/gabyoall/logger"
)

var envRegex = regexp.MustCompile(`(\$({([^}]*)}))`)

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

func RunTemplating(l logger.AppLogger, templateString, name string, vars interface{}) string {
	templateString = expandEnv(templateString)
	return executeTemplate(l, templateString, name, vars)
}
