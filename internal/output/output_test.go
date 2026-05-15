package output

import (
	"bytes"
	"testing"

	"github.com/MyCode83/bflask/pkg/bflask"
)

func TestFoundQuietPrintsOnlySecretKey(t *testing.T) {
	var out bytes.Buffer
	printer := New(false, true)
	printer.out = &out

	err := printer.Found(bflask.Result{
		SecretKey: "super secret",
		Payload:   `{"user":"admin"}`,
	})
	if err != nil {
		t.Fatalf("Found returned error: %v", err)
	}

	want := "super secret\n"
	if got := out.String(); got != want {
		t.Fatalf("quiet output = %q, want %q", got, want)
	}
}

func TestNotFoundQuietPrintsNothing(t *testing.T) {
	var out bytes.Buffer
	printer := New(false, true)
	printer.out = &out

	err := printer.NotFound(bflask.Stats{Checked: 10})
	if err != nil {
		t.Fatalf("NotFound returned error: %v", err)
	}

	if got := out.String(); got != "" {
		t.Fatalf("quiet not-found output = %q, want empty", got)
	}
}
