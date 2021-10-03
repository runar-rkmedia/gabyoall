package bboltStorage

import (
	"github.com/teris-io/shortid"
)

func CreateUniqueId() (string, error) {
	return shortid.Generate()
}
