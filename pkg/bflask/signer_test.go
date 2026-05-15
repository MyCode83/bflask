package bflask

import "testing"

func TestSignCookieCanBeVerified(t *testing.T) {
	const payload = `{"user":"admin"}`
	const secret = "supersecret"
	const salt = "cookie-session"
	const digest = "sha256"

	cookie, err := SignCookie([]byte(payload), secret, salt, digest)
	if err != nil {
		t.Fatal(err)
	}

	verifier, err := NewVerifier(cookie, salt, digest)
	if err != nil {
		t.Fatal(err)
	}

	decoded, ok, err := verifier.Verify(secret)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected signed cookie to verify")
	}
	if string(decoded) != payload {
		t.Fatalf("decoded payload = %s, want %s", decoded, payload)
	}
}
