package types

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/runar-rkmedia/gabyoall/requests"
)

type CompactRequestStatisticsEntity struct {
	Entity
	CompactRequestStatistics
}
type CompactRequestStatistics struct {
	requests.Stats
	// required: true
	StartTime       time.Time
	ResponseHashMap ByteHashMap
	Requests        map[string]CompactStat
}

type CompactStat struct {
	requests.ErrorType
	ResponseHash Hash          `json:"response_hash,omitempty"`
	StatusCode   int16         `json:"status_code,omitempty"`
	Error        string        `json:"error,omitempty"`
	RequestID    string        `json:"request_id,omitempty"`
	Offset       time.Duration `json:"offset,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
}

func hash256(content []byte) [32]byte {
	sum := sha256.New().Sum(content)
	c := (*[32]byte)(sum)
	return *c
}

// This is only for use will the original RequestStat
func (rs *CompactRequestStatistics) AddStat(stat requests.RequestStat) {
	s := CompactStat{
		ErrorType:  stat.ErrorType,
		StatusCode: int16(stat.StatusCode),
		RequestID:  stat.RequestID,
		Offset:     stat.Start.Sub(rs.StartTime),
		Duration:   stat.EndTime.Sub(stat.Start),
	}
	if stat.Error != nil {
		s.Error = stat.Error.Error()
	}
	if len(stat.RawResponse) > 0 {
		body := []byte(stat.RawResponse)
		bodyHash := hash256(body)
		rs.ResponseHashMap[bodyHash] = body
		s.ResponseHash = bodyHash
	}
	if stat.Duration > rs.Max {
		rs.Max = stat.Duration
	}
	if stat.Duration < rs.Min {
		rs.Min = stat.Duration
	}
	rs.Total += stat.Duration
	rs.Requests[stat.RequestID] = s
}
func (rs *CompactRequestStatistics) RecalculateAll() {
	// to offset the min-calculation
	rs.Min = 100 * time.Hour
	for _, req := range rs.Requests {
		if req.Duration > rs.Max {
			rs.Max = req.Duration
		}
		if req.Duration < rs.Min {
			rs.Min = req.Duration
		}
		rs.Total += req.Duration
	}
	rs.Calculate()
}
func (rs *CompactRequestStatistics) Calculate() {
	length := len(rs.Requests)
	if length == 0 {
		return
	}
	rs.Average = rs.Total / time.Duration(length)
	rs.TotalText = rs.Total.String()
	rs.MinText = rs.Min.String()
	rs.MaxText = rs.Max.String()
	rs.AverageText = rs.Average.String()
}

type Hash [32]byte
type ByteHashMap map[[32]byte][]byte

func NewCompactRequestStatistics() CompactRequestStatistics {
	return CompactRequestStatistics{
		StartTime: time.Now(),
		Stats: requests.Stats{
			Min: time.Hour * 100,
		},
		ResponseHashMap: map[[32]byte][]byte{},
		Requests:        map[string]CompactStat{},
	}
}

type bob string

func (c ByteHashMap) MarshalJSON() ([]byte, error) {
	u := map[string]string{}
	for kb, vb := range c {
		u[base64.URLEncoding.EncodeToString(kb[:])] = string(vb)
	}
	return json.Marshal(u)
}

func (c Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + base64.URLEncoding.EncodeToString(c[:]) + `"`), nil
}
