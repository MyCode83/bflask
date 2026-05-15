package bflask

import "testing"

func TestDecodeCookie(t *testing.T) {
	payload, err := DecodeCookie(timedCookie)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"some_number": 3, "some_words": "Some TIMED short payload"}`
	if string(payload) != want {
		t.Fatalf("decoded payload = %s, want %s", payload, want)
	}
}
