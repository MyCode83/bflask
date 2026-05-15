package bflask

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/alitto/pond/v2"
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan string, e.opts.Threads*2)
	result := make(chan Result, 1)
	errs := make(chan error, 1)

	var checked atomic.Int64

	pool := pond.NewPool(e.opts.Threads)
	defer pool.StopAndWait()

	go func() {
		defer close(jobs)
		if err := StreamCandidates(ctx, e.opts.Wordlist, jobs); err != nil && !errors.Is(err, context.Canceled) {
			select {
			case errs <- err:
			default:
			}
		}
	}()

	submitDone := make(chan struct{})
	go func() {
		defer close(submitDone)
		for candidate := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
			}
			candidate := candidate
			task := pool.SubmitErr(func() error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				payload, ok, err := e.verifier.Verify(candidate)
				checked.Add(1)
				if err != nil || !ok {
					return err
				}

				cancel()
				r := Result{
					Found:      true,
					SecretKey:  candidate,
					Payload:    PrettyPayload(payload),
					RawPayload: string(payload),
					Stats: Stats{
						Loaded:  loaded,
						Checked: checked.Load(),
						Elapsed: time.Since(start),
					},
				}
				select {
				case result <- r:
				default:
				}
				return nil
			})
			_ = task
		}
	}()

	select {
	case r := <-result:
		<-submitDone
		r.Stats.Checked = checked.Load()
		r.Stats.Elapsed = time.Since(start)
		return r, nil
	case err := <-errs:
		cancel()
		<-submitDone
		return Result{Stats: Stats{Loaded: loaded, Checked: checked.Load(), Elapsed: time.Since(start)}}, err
	case <-submitDone:
		pool.StopAndWait()
		select {
		case r := <-result:
			r.Stats.Checked = checked.Load()
			r.Stats.Elapsed = time.Since(start)
			return r, nil
		default:
		}
		return Result{Stats: Stats{Loaded: loaded, Checked: checked.Load(), Elapsed: time.Since(start)}}, ErrNotFound
	case <-ctx.Done():
		<-submitDone
		return Result{Stats: Stats{Loaded: loaded, Checked: checked.Load(), Elapsed: time.Since(start)}}, ctx.Err()
	}
}
