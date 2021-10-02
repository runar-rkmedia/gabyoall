package requests

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

type RequestStats []RequestStat

func (r RequestStats) Calculate() Stats {

	s := Stats{
		Min: time.Hour * 24,
	}
	for i := 0; i < len(r); i++ {
		if r[i].Duration > s.Max {
			s.Max = r[i].Duration
		}
		if r[i].Duration < s.Min {
			s.Min = r[i].Duration
		}
		s.Total += r[i].Duration
	}
	s.Average = s.Total / time.Duration(len(r))
	s.TotalText = s.Total.String()
	s.MinText = s.Min.String()
	s.MaxText = s.Max.String()
	s.AverageText = s.Average.String()
	return s

}
