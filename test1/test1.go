package test1

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu      sync.Mutex
	limiter *rate.Limiter
}

func NewRateLimiter(interval time.Duration, maxEvents int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Every(interval), maxEvents),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.limiter.Allow()
}

func WatchAndExecute(path string, interval time.Duration, maxCalls int, fn func(), triggers ...fsnotify.Op) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return fmt.Errorf("failed to watch directory: %v", err)
	}

	limiter := NewRateLimiter(interval, maxCalls)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("watcher events channel closed unexpectedly")
			}
			for _, trigger := range triggers {
				if event.Op&trigger == trigger {
					if limiter.Allow() {
						fn()
					} else {
						slog.Debug("Function call suppressed due to rate limiting")
					}
					break
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher errors channel closed unexpectedly")
			}
			slog.Error("error watching directory", "err", err)
		}
	}
}

func RunTest1() {
	dir := "/tmp"
	interval := 10 * time.Second
	maxCallsPerInterval := 1

	triggers := []fsnotify.Op{fsnotify.Write, fsnotify.Create, fsnotify.Remove, fsnotify.Chmod}
	err := WatchAndExecute(dir, interval, maxCallsPerInterval, func() {
		slog.Debug("function called within time limit after file system change")
	}, triggers...)
	if err != nil {
		panic(err)
	}
}
