package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/MyCode83/bflask/internal/config"
	"github.com/MyCode83/bflask/internal/logging"
	"github.com/MyCode83/bflask/internal/output"
	"github.com/MyCode83/bflask/pkg/bflask"
)

var crackCmd = &cobra.Command{
	Use:   "crack",
	Short: "Crack a Flask signed session cookie with a SECRET_KEY wordlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load(v)
		log := logging.New(cfg.Verbose, cfg.Quiet, cfg.JSON)
		printer := output.New(cfg.JSON, cfg.Quiet)
		printer.Banner()

		if cfg.Cookie == "" {
			return errors.New("--cookie is required")
		}
		if cfg.Wordlist == "" {
			return errors.New("--wordlist is required")
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if cfg.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
		}

		engine, err := bflask.NewEngine(bflask.Options{
			Cookie:   cfg.Cookie,
			Wordlist: cfg.Wordlist,
			Threads:  cfg.Threads,
			Salt:     cfg.Salt,
			Digest:   cfg.Digest,
		})
		if err != nil {
			return err
		}

		loaded, err := engine.CountCandidates()
		if err != nil {
			return fmt.Errorf("load wordlist: %w", err)
		}
		log.Info().Int64("candidates", loaded).Msg("Loaded candidate keys")
		log.Info().Int("threads", cfg.Threads).Str("digest", cfg.Digest).Str("salt", cfg.Salt).Msg("Starting bruteforce")

		result, err := engine.Crack(ctx, loaded)
		if errors.Is(err, bflask.ErrNotFound) {
			return printer.NotFound(result.Stats)
		}
		if err != nil {
			return err
		}

		if cfg.Output != "" {
			if writeErr := writeOutput(cfg.Output, result); writeErr != nil {
				return writeErr
			}
		}

		return printer.Found(result)
	},
}

func init() {
	flags := crackCmd.Flags()
	flags.StringP("cookie", "c", "", "Flask signed session cookie")
	flags.StringP("wordlist", "w", "", "path to SECRET_KEY wordlist")
	flags.IntP("threads", "t", 50, "number of concurrent workers")
	flags.StringP("salt", "s", "cookie-session", "itsdangerous signer salt")
	flags.StringP("digest", "d", "sha1", "digest algorithm: sha1, sha224, sha256, sha384, sha512, md5")
	flags.BoolP("verbose", "v", false, "enable verbose logging")
	flags.Duration("timeout", 0, "overall timeout, for example 30s or 5m")
	flags.StringP("output", "o", "", "write successful result to a file")
	flags.BoolP("json", "j", false, "emit JSON result")

	mustBind("cookie")
	mustBind("wordlist")
	mustBind("threads")
	mustBind("salt")
	mustBind("digest")
	mustBind("verbose")
	mustBind("timeout")
	mustBind("output")
	mustBind("json")
}

func mustBind(key string) {
	if err := v.BindPFlag(key, crackCmd.Flags().Lookup(key)); err != nil {
		panic(err)
	}
}

func writeOutput(path string, result bflask.Result) error {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
