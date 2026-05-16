package bflask

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNotFound = errors.New("secret key not found")

type Options struct {
	Cookie   string
	Wordlist string
	Threads  int
	Salt     string
	Digest   string
}

type Stats struct {
	Loaded  int64         `json:"loaded"`
	Checked int64         `json:"checked"`
	Elapsed time.Duration `json:"elapsed"`
}

type Result struct {
	Found      bool   `json:"found"`
	SecretKey  string `json:"secret_key,omitempty"`
	Payload    string `json:"payload,omitempty"`
	RawPayload string `json:"raw_payload,omitempty"`
	Stats      Stats  `json:"stats"`
}

type Engine struct {
	opts     Options
	verifier *Verifier
}

func NewEngine(opts Options) (*Engine, error) {
	if opts.Threads <= 0 {
		opts.Threads = 1
	}
	verifier, err := NewVerifier(opts.Cookie, opts.Salt, opts.Digest)
	if err != nil {
		return nil, err
	}
	return &Engine{opts: opts, verifier: verifier}, nil
}

func (e *Engine) CountCandidates() (int64, error) {
	return CountCandidates(e.opts.Wordlist)
}

func (e *Engine) Crack(ctx context.Context, loaded int64) (Result, error) {
	start := time.Now()
	parentCtx := ctx
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan string, e.opts.Threads*2)
	result := make(chan Result, 1)
	errs := make(chan error, 1)

	var checked atomic.Int64
	stats := func() Stats {
		return Stats{Loaded: loaded, Checked: checked.Load(), Elapsed: time.Since(start)}
	}
	sendErr := func(err error) {
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}
		select {
		case errs <- err:
			cancel()
		default:
		}
	}

	go func() {
		defer close(jobs)
		sendErr(StreamCandidates(ctx, e.opts.Wordlist, jobs))
	}()

	var wg sync.WaitGroup
	for range e.opts.Threads {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case candidate, ok := <-jobs:
					if !ok {
						return
					}

					payload, ok, err := e.verifier.Verify(candidate)
					checked.Add(1)
					if err != nil {
						sendErr(err)
						return
					}
					if !ok {
						continue
					}

					r := Result{
						Found:      true,
						SecretKey:  candidate,
						Payload:    PrettyPayload(payload),
						RawPayload: string(payload),
						Stats:      stats(),
					}
					select {
					case result <- r:
						cancel()
					default:
					}
					return
				}
			}
		}()
	}

	workersDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(workersDone)
	}()

	select {
	case r := <-result:
		cancel()
		<-workersDone
		r.Stats = stats()
		return r, nil
	case err := <-errs:
		cancel()
		<-workersDone
		return Result{Stats: stats()}, err
	case <-workersDone:
		select {
		case r := <-result:
			r.Stats = stats()
			return r, nil
		case err := <-errs:
			return Result{Stats: stats()}, err
		default:
		}
		if err := parentCtx.Err(); err != nil {
			return Result{Stats: stats()}, err
		}
		return Result{Stats: stats()}, ErrNotFound
	case <-ctx.Done():
		<-workersDone
		select {
		case r := <-result:
			r.Stats = stats()
			return r, nil
		case err := <-errs:
			return Result{Stats: stats()}, err
		default:
		}
		if err := parentCtx.Err(); err != nil {
			return Result{Stats: stats()}, err
		}
		return Result{Stats: stats()}, ErrNotFound
	}
}
