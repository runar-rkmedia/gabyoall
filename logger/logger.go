// package logger is a simple wrapper for a log-interface.
package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Currently, this is using zerolog as an experiment.
// Previously, logrus has been used, which was very convenient, but why not try something else?

var (
	al AppLogger
)

type AppLogger struct {
	zerolog.Logger
}

func With(l zerolog.Logger) AppLogger {
	return AppLogger{l}
}

type LogConfig struct {
	Level      string
	Format     string
	WithCaller bool
}

func convertLevelStr(s string) (zerolog.Level, bool) {
	switch strings.ToLower(s) {
	case "panic", "5":
		return zerolog.PanicLevel, true
	case "fatal", "4":
		return zerolog.FatalLevel, true
	case "error", "3":
		return zerolog.ErrorLevel, true
	case "warn", "warning", "2":
		return zerolog.WarnLevel, true
	case "info", "1":
		return zerolog.InfoLevel, true
	case "debug", "0":
		return zerolog.DebugLevel, true
	case "trace", "-1":
		return zerolog.TraceLevel, true
	}
	return zerolog.InfoLevel, false
}

func InitLogger(cfg LogConfig) AppLogger {
	var l zerolog.Logger
	switch cfg.Format {
	case "human":
		out := zerolog.ConsoleWriter{Out: os.Stderr}
		out.FormatTimestamp = func(i interface{}) string {
			return i.(string)
		}
		l = log.Output(out)
	default:
		l = log.Logger
	}

	if cfg.WithCaller {
		l = l.With().Caller().Logger()
	}
	if level, ok := convertLevelStr(cfg.Level); ok {
		l = l.Level(level)
	}
	al = AppLogger{l}

	return al
}

func GetLoggerWithLevel(label, level string) AppLogger {
	l := GetLogger(label)
	lvl, ok := convertLevelStr(level)
	if !ok {
		l.Warn().Str("level", level).Msg("The level was not correct")
	}
	l = AppLogger{l.Level(lvl)}
	return l
}
func GetLogger(label string) AppLogger {
	l := al.With().Str("label", label).Logger()
	return AppLogger{l}
}

func (al *AppLogger) HasDebug() bool {

	return al.GetLevel() <= zerolog.DebugLevel
}
func (al *AppLogger) HasTrace() bool {
	return al.GetLevel() <= zerolog.TraceLevel
}
func (al *AppLogger) ErrErr(err error) *zerolog.Event {
	return al.Error().Err(err)
}
func (al *AppLogger) ErrWarn(err error) *zerolog.Event {
	return al.Warn().Err(err)
}
func (al *AppLogger) WithStringPairs(pairs ...string) AppLogger {
	l := al.With()
	for i := 0; i < len(pairs)-1; i += 2 {
		l = l.Str(pairs[i], pairs[i+1])
	}
	return AppLogger{l.Logger()}
}

func Debug(s string, js ...interface{}) {
	fmt.Println("\n\n\t\t", s)
	for _, j := range js {
		v, _ := yaml.Marshal(j)
		fmt.Println(string(v))
	}
}
