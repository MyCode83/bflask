package bflask

import (
	"strings"
	"testing"
)

func TestParseCookie(t *testing.T) {
	cookie := "eyJ1c2VyIjoiYWRtaW4ifQ.CsxT0w.signature"
	parsed, err := ParseCookie(cookie)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Payload != "eyJ1c2VyIjoiYWRtaW4ifQ" || parsed.Timestamp != "CsxT0w" || parsed.Signature != "signature" {
		t.Fatalf("unexpected parsed cookie: %#v", parsed)
	}
}

func TestParseCompressedCookie(t *testing.T) {
	cookie := ".eJyrVkrLz1eyUkpKLFKqBQAdegQ0.CsxT0w.signature"
	parsed, err := ParseCookie(cookie)
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.Compressed || !strings.HasPrefix(parsed.Payload, ".") {
		t.Fatalf("expected compressed cookie payload, got %#v", parsed)
	}
}

func TestParseInvalidCookie(t *testing.T) {
	if _, err := ParseCookie("not-a-cookie"); err == nil {
		t.Fatal("expected invalid cookie error")
	}
}
