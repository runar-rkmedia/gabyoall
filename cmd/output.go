package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	tm "github.com/buger/goterm"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/requests"
	queries "github.com/runar-rkmedia/gabyoall/requests"
)

type Output struct {
	l          logger.AppLogger
	Url        string
	Query      queries.Request
	Details    map[requests.ErrorType]requests.RequestStats
	JwtPayload map[string]interface{}
	Count      map[requests.ErrorType]int
	Stats      map[requests.ErrorType]requests.Stats
	path       string
}

type Marshal func(j interface{}) ([]byte, error)

func (o *Output) AddStat(stat requests.RequestStat) *Output {
	o.Details[stat.ErrorType] = append(o.Details[stat.ErrorType], stat)
	o.Count[stat.ErrorType]++
	return o
}
func (o *Output) Write() error {
	if o.path == "" {
		return nil
	}
	err := WriteAuto(o.path, o)
	if err != nil {
		return err
	}
	o.l.Info().Str("path", o.path).Msg("Wrote to file")
	return nil
}

func (o *Output) CalculateStats() {
	for errorType, r := range o.Details {
		s := r.Calculate()
		o.Stats[errorType] = s
	}
}
func (o *Output) GetPath() string {
	return o.path
}

func (out *Output) PrintTable() {

	totals := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(totals, "\nCount\tErrorType\tMin\tAverage\tMax\tTotal\n")
	countSort := []struct {
		Count     int
		ErrorType requests.ErrorType
	}{}

	for k, v := range out.Count {
		countSort = append(countSort, struct {
			Count     int
			ErrorType requests.ErrorType
		}{v, k})
	}
	sort.SliceStable(countSort, func(i, j int) bool {
		return countSort[i].Count > countSort[j].Count
	})
	out.CalculateStats()
	for _, c := range countSort {
		s := out.Stats[c.ErrorType]
		fmt.Fprintf(totals, "%d\t%s\t%s\t%s\t%s\t%s\n", c.Count, c.ErrorType, s.MinText, s.AverageText, s.MaxText, s.TotalText)
	}
	tm.Println(totals)
}

func NewOutput(l logger.AppLogger, path, url string, query queries.Request, JwtPayload map[string]interface{}) (Output, error) {
	abs := ""
	if path != "" {

		_abs, err := filepath.Abs(path)
		if err != nil {
			return Output{}, fmt.Errorf("filepath %s returned error: %w", path, err)
		}
		abs = _abs
	}
	return Output{
		l:          l,
		path:       abs,
		Url:        url,
		Query:      query,
		JwtPayload: JwtPayload,
		Details:    map[requests.ErrorType]requests.RequestStats{},
		Count:      map[requests.ErrorType]int{},
		Stats:      map[requests.ErrorType]requests.Stats{},
	}, nil
}
