package utils

import (
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
)

// Returns a unique id.
// The id is guaranteed to be created and can be used even if there is an error.
// The error should be logged.
func ForceCreateUniqueId() (string, error) {
	id, err := shortid.Generate()
	var didErr error
	// If for some reason this crashes, (I dont know why it would), lets try again once.
	if err != nil {
		didErr = err
		id, err = shortid.Generate()
		// Not really sure how the above could fail, but if it does, I dont really care too much.
		// Lets just generate a uuid then
		if err != nil {
			id = uuid.NewString()
		}
	}
	if err != nil {
		err = didErr
	}
	return id, err
}
