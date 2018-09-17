package heap

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Heap will check every checkEvery for memory obtained from the system, by the process
// - using the metric Sys at https://golang.org/pkg/runtime/#MemStats -
// whether it reached or exceeds the threshold and take a memory profile to the path directory,
// but no more often than every minTimeDiff seconds
// any errors will be sent to the errors channel
type Heap struct {
	cfg      Config
	lastUnix int64
	Errors   chan error
}

// Config is the config for triggering profile
type Config struct {
	Path        string
	Threshold   int
	MinTimeDiff int
	CheckEvery  time.Duration
}

// New creates a new Heap trigger. use a nil channel if you don't care about any errors
func New(cfg Config, errors chan error) (*Heap, error) {
	heap := Heap{
		cfg:      cfg,
		lastUnix: int64(0),
		Errors:   errors,
	}
	return &heap, nil
}

func (heap Heap) logError(err error) {
	if heap.Errors != nil {
		heap.Errors <- err
	}
}

// Run runs the trigger. encountered errors go to the configured channel (if any).
// you probably want to run this in a new goroutine.
func (heap Heap) Run() {
	cfg := heap.cfg
	tick := time.NewTicker(cfg.CheckEvery)
	m := &runtime.MemStats{}
	for ts := range tick.C {
		runtime.ReadMemStats(m)
		unix := ts.Unix()
		if m.Sys >= uint64(cfg.Threshold) && unix >= heap.lastUnix+int64(cfg.MinTimeDiff) {
			f, err := os.Create(fmt.Sprintf("%s/%d.profile-heap", cfg.Path, unix))
			if err != nil {
				heap.logError(err)
				continue
			}
			err = pprof.WriteHeapProfile(f)
			if err != nil {
				heap.logError(err)
			}
			heap.lastUnix = unix
			err = f.Close()
			if err != nil {
				heap.logError(err)
			}
		}
	}
}
