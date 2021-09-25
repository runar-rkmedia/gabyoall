package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

func WriteAuto(outpath string, content interface{}) error {
	ext := filepath.Ext(outpath)
	var marshal Marshal
	switch ext {
	case ".yaml", ".yml":
		marshal = yaml.Marshal
	default:
		marshal = func(j interface{}) ([]byte, error) {
			return json.MarshalIndent(j, "", "  ")
		}
	}
	return Write(marshal, outpath, content)
}

func Write(marshal Marshal, outpath string, content interface{}) error {
	b, err := marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml for output: %w", err)
	}
	err = os.WriteFile(outpath, b, 0666)
	if err != nil {
		return fmt.Errorf("failed to write file-contents to path %s: %w", outpath, err)
	}
	return nil
}
