package printer

import (
	"fmt"
	"sync"
	"time"

	tm "github.com/buger/goterm"
	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/utils"
	"github.com/tj/go-spin"
)

type ValidityStringer interface {
	PrintValidity() (string, error)
}

type Printer struct {
	config           cmd.Config
	validityStringer ValidityStringer
	startTime        time.Time
	spinner          *spin.Spinner
	out              cmd.Output
	operationName    string
	lastOut          time.Time
	quit             chan struct{}
	sync.Mutex
}

func NewPrinter(config cmd.Config, validityStringer ValidityStringer, operationName string, out cmd.Output, startTime time.Time) *Printer {
	printer := Printer{
		config:           config,
		validityStringer: validityStringer,
		startTime:        startTime,
		spinner:          spin.New(),
		operationName:    operationName,
		quit:             make(chan struct{}),
		out:              out,
	}
	return &printer
}

func (p *Printer) Animate() chan struct{} {
	go func(quit chan struct{}) {
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(800 * time.Millisecond)
				p.spinner.Next()
			}
		}

	}(p.quit)
	return p.quit
}
func (p *Printer) Update(i int, successes int) {

	now := time.Now()
	p.Lock()
	if now.Sub(p.lastOut) < 500*time.Millisecond {
		p.Unlock()
		return
	}
	p.lastOut = now
	p.Unlock()
	p.update(i, successes)
}
func (p *Printer) update(i int, successes int) {
	if p.config.PrintTable {
		tm.Clear()
		p.out.PrintTable()
		tm.Flush()
	}
	validStr := ""
	if p.validityStringer != nil {
		_validStr, _ := p.validityStringer.PrintValidity()
		validStr = _validStr
	}
	fraction := float64(i) / float64(p.config.RequestCount)
	dur := time.Now().Sub(p.startTime)
	estimatedCompletion := time.Duration(float64(dur)/fraction) - dur
	fails := ""
	failures := i - successes
	if failures > 0 {
		fails = fmt.Sprintf("\033[31m[%d (%.2f%%)\033[0m", failures, float64(failures)/float64(i)*100)
	}
	fmt.Printf("\r\033[36m[%d/%d (%.2f%%) %s -c=%d] %s Waiting for result from: %s (%s) \033[m %s (%s) %s", i, p.config.RequestCount, fraction*100, fails, p.config.Concurrency, p.spinner.Current(), p.config.Url, p.operationName, utils.PrettyDuration(dur), utils.PrettyDuration(estimatedCompletion), validStr)

}

func (p *Printer) Complete(i int, successes int) {
	p.update(i, successes)
	close(p.quit)
	fmt.Println("")
}
