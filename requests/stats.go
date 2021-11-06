package requests

import (
	"encoding/json"
	"time"

	"github.com/runar-rkmedia/gabyoall/utils"
)

type RequestStat struct {
	ts          TimeSeriePusher
	ErrorType   `json:"errorType,omitempty"`
	RawResponse []byte    `json:"rawResponse,omitempty"`
	ContentType string    `json:"-"`
	Start       time.Time `json:"-"`
	RequestID   string
	Duration    time.Duration `json:"duration,omitempty"`
	CompactStat
}

func (r *RequestStat) End(body []byte, errorType ErrorType, err error) RequestStat {
	r.ErrorType = errorType
	endTime := time.Now()
	r.Duration = endTime.Sub(r.Start)
	r.ts.Push(string(errorType), endTime, float64(r.Duration))
	r.RawResponse = body
	if err != nil {
		r.Error = err.Error()
	}
	return *r
}

type Durationable time.Duration

type Stats struct {
	Total   time.Duration
	Min     time.Duration
	Max     time.Duration
	Average time.Duration
}

func NewStat(offset time.Duration, ts TimeSeriePusher) RequestStat {
	id, _ := utils.ForceCreateUniqueId()
	return RequestStat{
		CompactStat: CompactStat{},
		ts:          ts,
		Start:       time.Now(),
		RequestID:   id,
	}
}

func (c *Stats) MarshalJSON() ([]byte, error) {
	t := struct {
		Total   int64 `json:"total,omitempty"`
		Min     int64 `json:"min,omitempty"`
		Max     int64 `json:"max,omitempty"`
		Average int64 `json:"average,omitempty"`
	}{
		Total:   c.Total.Milliseconds(),
		Min:     c.Min.Milliseconds(),
		Max:     c.Max.Milliseconds(),
		Average: c.Average.Milliseconds(),
	}
	return json.Marshal(t)

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
	return s

}
