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

type Printer struct {
	config        cmd.Config
	jwtPayload    cmd.JwtPayload
	startTime     time.Time
	spinner       *spin.Spinner
	out           cmd.Output
	operationName string
	lastOut       time.Time
	sync.Mutex
}

func NewPrinter(config cmd.Config, jwtPayload cmd.JwtPayload, operationName string, out cmd.Output, startTime time.Time) *Printer {
	printer := Printer{
		config:        config,
		jwtPayload:    jwtPayload,
		startTime:     startTime,
		spinner:       spin.New(),
		operationName: operationName,
		out:           out,
	}
	return &printer
}

func (p *Printer) Animate() chan struct{} {
	spinCh := make(chan struct{})
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

	}(spinCh)
	return spinCh
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
	payloadExp := ""
	if p.config.NoTokenValidation != true && p.jwtPayload.ExpiresAt != nil {
		payloadExp = fmt.Sprintf("Token expires: %s", p.jwtPayload.ExpiresAt.Sub(time.Now()).String())
	}
	fraction := float64(i) / float64(p.config.RequestCount)
	dur := time.Now().Sub(p.startTime)
	estimatedCompletion := time.Duration(float64(dur)/fraction) - dur
	fails := ""
	failures := i - successes
	if failures > 0 {
		fails = fmt.Sprintf("\033[31m[%d (%.2f%%)\033[0m", failures, float64(failures)/float64(i)*100)
	}
	fmt.Printf("\r\033[36m[%d/%d (%.2f%%) %s -c=%d] %s Waiting for result from: %s (%s) \033[m %s (%s) %s", i, p.config.RequestCount, fraction*100, fails, p.config.Concurrency, p.spinner.Current(), p.config.Url, p.operationName, utils.PrettyDuration(dur), utils.PrettyDuration(estimatedCompletion), payloadExp)

}

func (p *Printer) Complete(i int, successes int) {
	p.update(i, successes)
	fmt.Println("")
}
