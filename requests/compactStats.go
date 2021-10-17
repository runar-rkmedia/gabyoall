package requests

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/pelletier/go-toml"
)

type CompactRequestStatistics struct {
	Stats
	// required: true
	StartTime       time.Time
	ResponseHashMap ByteHashMap `json:"response_hash_map,omitempty"`
	Requests        map[string]CompactStat
}

type CompactStat struct {
	ErrorType    `json:"errorType,omitempty"`
	ResponseHash *Hash         `json:"response_hash,omitempty"`
	StatusCode   int16         `json:"status_code,omitempty"`
	Error        string        `json:"error,omitempty"`
	RequestID    string        `json:"request_id,omitempty"`
	Offset       time.Duration `json:"offset,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
}

// should match uuids, and the request-id that this project creates for each request.
var requestIdRegex = regexp.MustCompile(`([a-fA-F0-9-]{25,32}|srv-test-[-_\w]*)`)

func hash256(content []byte) *[32]byte {
	if len(content) == 0 {
		return nil
	}
	s := sha256.New()
	// TODO: this should probably be a setting on endpoint/request.
	// The point is to somehow ignore any references to unique request ids so that the same response
	// even though they contain unique request-ids still hashes to the
	normalized := requestIdRegex.ReplaceAll(content, []byte("__UID__"))
	s.Write([]byte(normalized))
	sum := s.Sum(nil)
	fmt.Println("hash256", len(sum), string(sum), string(content))
	c := (*[32]byte)(sum)
	return c
}

// This is only for use will the original RequestStat
func (rs *CompactRequestStatistics) AddStat(stat RequestStat) {
	s := CompactStat{
		ErrorType:  stat.ErrorType,
		StatusCode: int16(stat.StatusCode),
		RequestID:  stat.RequestID,
		Offset:     stat.Start.Sub(rs.StartTime),
		Duration:   stat.Duration,
	}
	if len(stat.RawResponse) > 0 {
		body := stat.RawResponse
		bodyHash := hash256(body)
		if bodyHash != nil {
			rs.ResponseHashMap[*bodyHash] = ByteContent{
				Content:     body,
				ContentType: stat.ContentType,
			}
			h := Hash(*bodyHash)
			s.ResponseHash = &h
		}
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
}

type Hash [32]byte
type ByteHashMap map[[32]byte]ByteContent

type ByteContent struct {
	Content     []byte `json:"content,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

func (bm ByteHashMap) Add(contentType string, body []byte) *Hash {
	if len(body) == 0 {
		return nil
	}

	bodyHash := hash256(body)
	b := *bodyHash
	if bodyHash == nil {
		return nil
	}
	fmt.Println("\n\nhashy Add", string(body), string(b[:]))
	bm[*bodyHash] = ByteContent{
		Content:     body,
		ContentType: contentType,
	}
	h := Hash(*bodyHash)
	return &h
}

func NewCompactRequestStatistics() CompactRequestStatistics {
	return CompactRequestStatistics{
		StartTime: time.Now(),
		Stats: Stats{
			Min: time.Hour * 100,
		},
		ResponseHashMap: ByteHashMap{},
		Requests:        map[string]CompactStat{},
	}
}

func (c ByteHashMap) MarshalJSON() ([]byte, error) {
	u := map[string]ByteContent{}
	for kb, vb := range c {
		u[base64.URLEncoding.EncodeToString(kb[:])] = vb
	}
	return json.Marshal(u)
}
func (c ByteContent) MarshalJSON() ([]byte, error) {
	t := struct {
		Content     interface{} `json:"content,omitempty"`
		ContentType string      `json:"contentType,omitempty"`
	}{
		ContentType: string(c.ContentType),
	}
	switch {
	case c.ContentType == "application/json", strings.Contains(c.ContentType, "json"):
		var j map[string]interface{}
		err := json.Unmarshal(c.Content, &j)
		if err == nil {
			t.Content = j
		}
		// TODO: test that these work
	case c.ContentType == "text/vnd.yaml", strings.Contains(c.ContentType, "yaml"):
		var j map[string]interface{}
		err := yaml.Unmarshal(c.Content, &j)
		if err == nil {
			t.Content = j
		}
		// TODO: test that these work
	case c.ContentType == "application/toml", strings.Contains(c.ContentType, "toml"):
		var j map[string]interface{}
		err := toml.Unmarshal(c.Content, &j)
		if err == nil {
			t.Content = j
		}
	}
	if t.Content == nil || t.Content == "" {
		t.Content = string(c.Content)
	}
	return json.Marshal(t)
}

func (c Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + base64.URLEncoding.EncodeToString(c[:]) + `"`), nil
}
