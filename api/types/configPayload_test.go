package types

import (
	"github.com/ghodss/yaml"
	"github.com/r3labs/diff/v2"
	"github.com/runar-rkmedia/gabyoall/cmd"
	"testing"
)

func TestConfig_MergeInto(t *testing.T) {
	type fields struct {
		Auth          *AuthConfig
		OkStatusCodes *[]int
		ResponseData  *bool
		Concurrency   *int
		RequestCount  *int
		Secrets       *Secrets
	}
	tests := []struct {
		name   string
		fields fields
		target cmd.Config
		want   cmd.Config
	}{
		{
			"Should merge values",
			fields{Concurrency: pint(2)},
			cmd.Config{Concurrency: 1, RequestCount: 1},
			cmd.Config{Concurrency: 2, RequestCount: 1},
		},
		{
			"Should not merge null-values",
			fields{},
			cmd.Config{Concurrency: 1, RequestCount: 3},
			cmd.Config{Concurrency: 1, RequestCount: 3},
		},
		{
			"Should merge int-arrays",
			fields{OkStatusCodes: pints(2, 3)},
			cmd.Config{OkStatusCodes: []int{1, 2}, RequestCount: 3},
			cmd.Config{OkStatusCodes: []int{2, 3}, RequestCount: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Auth:          tt.fields.Auth,
				OkStatusCodes: tt.fields.OkStatusCodes,
				ResponseData:  tt.fields.ResponseData,
				Concurrency:   tt.fields.Concurrency,
				RequestCount:  tt.fields.RequestCount,
				Secrets:       tt.fields.Secrets,
			}
			w1 := Pretty(c)
			DiffTest(t, c.MergeInto(tt.target), tt.want)
			w2 := Pretty(c)
			if w1 != w2 {
				t.Errorf("Original config-struct was modified")
			}
		})
	}
}

func DiffTest(t *testing.T, got, want interface{}) {
	changelog, err := diff.Diff(got, want)
	if err != nil {
		t.Errorf("failed during changelog-generation: %s", err)
	}
	if len(changelog) == 0 {
		return
	}
	t.Errorf("Diff: \n%s\nGot:\n%s\nWant:\n%s", Pretty(changelog), Pretty(got), Pretty(want))
}

func Pretty(v interface{}) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func pint(n int) *int {
	return &n
}
func pints(n ...int) *[]int {
	return &n
}
