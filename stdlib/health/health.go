package health

import (
	"context"
	"fmt"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"sync"
	"time"
)

const (
	infoHealth string = `Health:`
	errHealth  string = `%s Health Error`

	_OK     string = `[OK]`
	_FAILED string = `[FAILED]`
)

type probe string

const (
	ready probe = `[READINESS CHECK]`
	live  probe = `[LIVENESS CHECK]`
)

type Health interface {
	// HealthEndpoint returns url endpoint of liveness check
	HealthEndpoint() string
	// ReadyEndpoint returns url endpoint of readiness check
	ReadyEndpoint() string
	// IsReady returns app readiness status as error interface
	IsReady() error
	// IsHealthy returns app liveness status as error interface
	IsHealthy() error
	// IsReadyAndHealthy returns readiness and liveness status as error interface
	IsReadyAndHealthy() error
	// Stop readiness and liveness checker goroutines.
	// This should be called during app termination to free resource.
	Stop()
}

type Options struct {
	// WaitBeforeContinue Blocking all process before app is ready and healthy
	WaitBeforeContinue bool
	// MaxWaitingTime defines max waiting time during blocking time ( WaitBeforeContinue )
	MaxWaitingTime time.Duration
	// Liveness probe configuration
	Liveness ProbeOptions
	// Readiness probe configuration
	Readiness ProbeOptions
}

type ProbeOptions struct {
	Enabled bool
	// Success Threshold at respective probe
	SuccessThreshold int
	// Failure Threshold at respective probe
	FailureThreshold int
	// InitialDelaySec
	InitialDelaySec time.Duration
	// PeriodSec determines how frequent, we need to do probe checking
	PeriodSec time.Duration
	// CheckTimeout timeout in single probe check
	CheckTimeout time.Duration
	// CheckF function called to define probe status
	CheckF func(ctx context.Context, cancel context.CancelFunc) error
	// Endpoint endpoint that will be registered in httpmux
	Endpoint string
}

type health struct {
	logger log.Logger
	status *status
	opt Options
	termReady chan struct{}
	termLive chan struct{}
}

type status struct {
	mu *sync.RWMutex
	isHealthy bool
	isReady bool
}

func Init(logger log.Logger, opt Options) Health{
	health := &health{
		logger: logger,
		status: &status{
			mu:        &sync.RWMutex{},
			isHealthy: false,
			isReady:   false,
		},
		opt:       opt,
		termReady: make(chan struct{}, 1),
		termLive:  make(chan struct{}, 1),
	}
	health.runCheckers()
	if opt.WaitBeforeContinue {
		if err := health.waitUntilReadyAndHealthy(); err != nil {
			logger.Panic(ErrAppInitFailed)
			return nil
		}
	}
	return health
}

func (h *health) waitUntilReadyAndHealthy() error {
	ticker := time.NewTicker(1 * time.Second)
	maxAttempt := int64(h.opt.MaxWaitingTime)/int64(time.Second) + 1
	attempt := int64(0)
	for range ticker.C {
		if err := h.IsReady(); err == nil {
			if err := h.IsHealthy(); err == nil {
				ticker.Stop()
				break
			}
		}
		attempt++
		if attempt > maxAttempt {
			return ErrAppInitFailed
		}
	}
	return nil
}

func (h *health) runCheckers() {
	if h.opt.Readiness.Enabled {
		if h.opt.Readiness.CheckF == nil {
			h.logger.Panic(ErrReadyCheckFuncIsNil)
			return
		}
		go h.AsyncChecker(ready)
	}
	if h.opt.Liveness.Enabled {
		if h.opt.Liveness.CheckF == nil {
			h.logger.Panic(ErrHealthCheckFuncIsNil)
			return
		}
		go h.AsyncChecker(live)
	}
}

func (h *health) AsyncChecker(mode probe) {
	var (
		ticker           *time.Ticker
		successThreshold int
		failureThreshold int
		initialDelay     time.Duration
		timeout          time.Duration
		f                func(ctx context.Context, canc context.CancelFunc) error
		failure          int
		success          int
		termSig          chan struct{}
	)

	switch mode {
	case ready:
		ticker = time.NewTicker(h.opt.Readiness.PeriodSec)
		successThreshold = h.opt.Readiness.SuccessThreshold
		failureThreshold = h.opt.Readiness.FailureThreshold
		initialDelay := h.opt.Readiness.InitialDelaySec - h.opt.Readiness.PeriodSec
		if initialDelay < 0 {
			initialDelay = 0
		}
		timeout = h.opt.Readiness.CheckTimeout
		f = h.opt.Readiness.CheckF
		termSig = h.termReady

	case live:
		ticker = time.NewTicker(h.opt.Liveness.PeriodSec)
		successThreshold = h.opt.Liveness.SuccessThreshold
		failureThreshold = h.opt.Liveness.FailureThreshold
		initialDelay = h.opt.Liveness.InitialDelaySec - h.opt.Liveness.PeriodSec
		if initialDelay < 0 {
			initialDelay = 0
		}
		timeout = h.opt.Liveness.CheckTimeout
		f = h.opt.Liveness.CheckF
		termSig = h.termLive

	}
	//initial delay sec
	time.Sleep(initialDelay)

	h.logger.Info(_OK, infoHealth, fmt.Sprintf("%s starts after %s delay", string(mode), initialDelay))

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := f(ctx, cancel)
			if err != nil {
				failure++
				if failure > failureThreshold {
					h.setStatus(mode, false)
					failure = 0
					success = 0
				}
			} else {
				success++
				if success > successThreshold {
					h.setStatus(mode, true)
					success = 0
					failure = 0
				}
			}
		case <-termSig:
			ticker.Stop()
			return
		}
	}
}

func (h *health) setStatus(mode probe, _OK bool) {
	s := h.status
	s.mu.Lock()
	defer s.mu.Unlock()
	if mode == live {
		s.isHealthy = _OK
		return
	}
	s.isReady = _OK
	return
}

func (h *health) IsReady() error {
	h.status.mu.RLock()
	defer h.status.mu.RUnlock()
	if h.status.isReady {
		return nil
	}
	return ErrNotReady
}

func (h *health) HealthEndpoint() string {
	return h.opt.Liveness.Endpoint
}

func (h *health) ReadyEndpoint() string {
	return h.opt.Readiness.Endpoint
}

func (h *health) IsHealthy() error {
	h.status.mu.RLock()
	defer h.status.mu.RUnlock()
	if h.status.isHealthy {
		return nil
	}
	return ErrNotHealthy
}

func (h *health) IsReadyAndHealthy() error {
	h.status.mu.RLock()
	defer h.status.mu.RUnlock()
	if h.status.isReady && h.status.isHealthy {
		return nil
	}
	return ErrNotReadyAndHealthy
}

func (h *health) Stop() {
	if h.opt.Readiness.Enabled {
		h.setStatus(ready, false)
		close(h.termReady)
	}
	if h.opt.Liveness.Enabled {
		h.setStatus(live, false)
		close(h.termLive)
	}
}