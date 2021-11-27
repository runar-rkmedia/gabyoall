package internal

import (
	"fmt"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/go-test/deep"
	"github.com/gookit/color"
)

var (
	DefaultCompareOptions = CompareOptions{true, true, true}
	// red                   = color.FgRed.Render
	cName = color.Green.Render
	cGot  = color.Yellow.Render
	cDiff = color.Red.Render
	cWant = color.Cyan.Render
)

func Compare(name string, got, want interface{}, options ...CompareOptions) error {

	opts := DefaultCompareOptions
	if options != nil && len(options) > 0 {
		opts = options[0]
	}

	// At least on of these must be applied
	if !opts.Reflect && !opts.Diff {
		panic("Invlid options")
	}

	if opts.Diff {
		if diff := deep.Equal(got, want); diff != nil {
			var g interface{} = got
			var w interface{} = want
			var d interface{} = diff
			if opts.Yaml {
				g = MustYaml(got)
				w = MustYaml(want)
				d = MustYaml(diff)
			}
			return fmt.Errorf("YAML: %s: \n%v\ndiff:\n%s\nwant:\n%v", cName(name), cGot(g), cDiff(d), cWant(w))
		}
	}
	if opts.Reflect {
		if !reflect.DeepEqual(got, want) {
			var g interface{} = got
			var w interface{} = want
			if opts.Yaml {
				g = MustYaml(got)
				w = MustYaml(want)
			}
			return fmt.Errorf("YAML: %s: \n%v\nwant:\n%v", cName(name), cGot(g), cWant(w))
		}
	}
	return nil
}

type CompareOptions struct {
	// Produces a diff of the result, but may in some edgecases not detect all errors (like differences in timezones)
	Diff,
	// Uses a traditional reflect.DeepEqual to perform tests.
	Reflect,
	// Produces output in yaml for readability
	Yaml bool
}

func MustYaml(j interface{}) string {
	b, err := yaml.Marshal(j)
	if err != nil {
		panic(fmt.Errorf("Failed to marshal: %w\n\n%v", err, j))
	}
	return string(b)
}
