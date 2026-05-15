package bflask

import (
	"crypto/hmac"
	"fmt"
	"hash"

	itsdangerous "github.com/octopart/go-itsdangerous"
)

type hmacAlgorithm struct {
	digest DigestFactory
}

func (a hmacAlgorithm) GetSignature(key, value []byte) []byte {
	mac := hmac.New(a.digest, key)
	mac.Write(value)
	return mac.Sum(nil)
}

func (a hmacAlgorithm) VerifySignature(key, value, sig []byte) bool {
	return hmac.Equal(sig, a.GetSignature(key, value))
}

type Verifier struct {
	cookie  Cookie
	salt    string
	digest  string
	factory DigestFactory
}

func NewVerifier(rawCookie, salt, digest string) (*Verifier, error) {
	cookie, err := ParseCookie(rawCookie)
	if err != nil {
		return nil, err
	}
	factory, err := DigestFactoryFor(digest)
	if err != nil {
		return nil, err
	}
	return &Verifier{cookie: cookie, salt: salt, digest: digest, factory: factory}, nil
}

func (v *Verifier) Verify(secret string) ([]byte, bool, error) {
	derivedKey := deriveHMACKey([]byte(secret), []byte(v.salt), v.factory)
	signer := itsdangerous.NewTimestampSignature(string(derivedKey), v.salt, ".", "none", v.factory(), hmacAlgorithm{digest: v.factory})
	payload, err := signer.UnsignB64([]byte(v.cookie.Raw), 0)
	if err != nil {
		return nil, false, nil
	}
	return payload, true, nil
}

func deriveHMACKey(secret, salt []byte, digest func() hash.Hash) []byte {
	mac := hmac.New(digest, secret)
	mac.Write(salt)
	return mac.Sum(nil)
}

func (v *Verifier) DecodeUnsignedPayload() ([]byte, error) {
	payload, err := DecodePayload(v.cookie)
	if err != nil {
		return nil, fmt.Errorf("decode unsigned payload: %w", err)
	}
	return payload, nil
}
