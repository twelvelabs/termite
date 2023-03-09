package ui

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

func NewProgressIndicator(ios *IOStreams) *ProgressIndicator {
	return &ProgressIndicator{
		ios: ios,
	}
}

type ProgressIndicator struct {
	ios *IOStreams

	spin *spinner.Spinner
	mu   sync.Mutex
}

func (pi *ProgressIndicator) IsEnabled() bool {
	return pi.ios.IsStdoutTTY() && pi.ios.IsStderrTTY()
}

func (pi *ProgressIndicator) Start() {
	pi.StartWithLabel("")
}

func (pi *ProgressIndicator) StartWithLabel(label string) {
	if !pi.IsEnabled() {
		return
	}

	pi.mu.Lock()
	defer pi.mu.Unlock()

	if pi.spin != nil {
		if label == "" {
			pi.spin.Prefix = ""
		} else {
			pi.spin.Prefix = label + " "
		}
		return
	}

	// https://github.com/briandowns/spinner#available-character-sets
	dotStyle := spinner.CharSets[11]
	sp := spinner.New(
		dotStyle,
		120*time.Millisecond,
		spinner.WithWriter(pi.ios.Err),
		spinner.WithColor("fgCyan"),
	)
	if label != "" {
		sp.Prefix = label + " "
	}

	sp.Start()
	pi.spin = sp
}

func (pi *ProgressIndicator) Stop() {
	pi.mu.Lock()
	defer pi.mu.Unlock()
	if pi.spin == nil {
		return
	}
	pi.spin.Stop()
	pi.spin = nil
}
