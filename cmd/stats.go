package cmd

import (
	"time"

	"github.com/google/uuid"
)

type RequestStat struct {
	ErrorType    `json:",omitempty"`
	Response     map[string]interface{} `json:",omitempty"`
	RawResponse  string                 `json:",omitempty"`
	StatusCode   int                    `json:",omitempty"`
	Error        error                  `json:",omitempty"`
	Duration     time.Duration
	DurationText string
	RequestID    string
	Start        time.Time
	EndTime      time.Time
}

func (r *RequestStat) End(errorType ErrorType, err error) RequestStat {
	r.ErrorType = errorType
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.Start)
	r.DurationText = r.Duration.String()
	if err != nil {
		r.Error = err
	}
	return *r
}

type Stats struct {
	Total       time.Duration
	Min         time.Duration
	Max         time.Duration
	Average     time.Duration
	TotalText   string
	MinText     string
	MaxText     string
	AverageText string
}

func NewStat() RequestStat {
	return RequestStat{
		RequestID: "srv-test-" + uuid.NewString(),
		Start:     time.Now(),
	}
}

type GqlResponse struct {
	Errors []Error                `json:"errors,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}
