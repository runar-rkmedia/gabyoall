package cmd

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

func ReadYamlFile(file string, j interface{}, replacor ...func(b []byte) []byte) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	for _, f := range replacor {
		f(b)
	}

	err = yaml.Unmarshal(b, j)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal file: %s %w", file, err)
	}
	return nil
}
