package bflask

import "testing"

const timedCookie = "eyJzb21lX251bWJlciI6IDMsICJzb21lX3dvcmRzIjogIlNvbWUgVElNRUQgc2hvcnQgcGF5bG9hZCJ9.CsxT0w.ErzMtlBiK4ro_UDkvLJDbT5AlNc"

func TestVerifier(t *testing.T) {
	verifier, err := NewVerifier(timedCookie, "cookie-session", "sha1")
	if err != nil {
		t.Fatal(err)
	}

	payload, ok, err := verifier.Verify("super secret 1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected key to verify")
	}
	if string(payload) != `{"some_number": 3, "some_words": "Some TIMED short payload"}` {
		t.Fatalf("unexpected payload: %s", payload)
	}

	if _, ok, err := verifier.Verify("wrong"); err != nil || ok {
		t.Fatalf("expected wrong key to fail without error, ok=%v err=%v", ok, err)
	}
}
