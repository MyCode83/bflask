package bflask

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEngineFindsKeyAndCancels(t *testing.T) {
	dir := t.TempDir()
	wordlist := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(wordlist, []byte("wrong\nsuper secret 1\nunused\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	engine, err := NewEngine(Options{
		Cookie:   timedCookie,
		Wordlist: wordlist,
		Threads:  2,
		Salt:     "cookie-session",
		Digest:   "sha1",
	})
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := engine.CountCandidates()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result, err := engine.Crack(ctx, loaded)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Found || result.SecretKey != "super secret 1" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestEngineHonorsCanceledContext(t *testing.T) {
	dir := t.TempDir()
	wordlist := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(wordlist, []byte("wrong\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	engine, err := NewEngine(Options{
		Cookie:   timedCookie,
		Wordlist: wordlist,
		Threads:  1,
		Salt:     "cookie-session",
		Digest:   "sha1",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = engine.Crack(ctx, 1)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
