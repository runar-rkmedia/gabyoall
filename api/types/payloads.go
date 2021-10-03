package types

import (
	"time"

	"github.com/runar-rkmedia/gabyoall/requests"
)

type Storage interface {
	Endpoints() (es map[string]EndpointEntity, err error)
	Endpoint(id string) (EndpointEntity, error)
	CreateEndpoint(e EndpointPayload) (EndpointEntity, error)
}

type Entity struct {
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	ID        string     `json:"id,omitempty"`
}

// TODO: implement
type Schedule struct{}

type EndpointPayload struct {
	// required: true
	// example: https://example.com
	Url     string              `json:"url,omitempty" validate:"required,uri"`
	Headers map[string][]string `json:"headers,omitempty" validate:"dive,max=1000"`
}

type EndpointEntity struct {
	requests.Endpoint
	Entity
}
