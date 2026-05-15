package bflask

import (
	"fmt"

	itsdangerous "github.com/octopart/go-itsdangerous"
)

func SignCookie(payload []byte, secret, salt, digest string) (string, error) {
	if len(payload) == 0 {
		return "", fmt.Errorf("payload is empty")
	}
	if secret == "" {
		return "", fmt.Errorf("secret is empty")
	}

	factory, err := DigestFactoryFor(digest)
	if err != nil {
		return "", err
	}

	derivedKey := deriveHMACKey([]byte(secret), []byte(salt), factory)
	signer := itsdangerous.NewTimestampSignature(string(derivedKey), salt, ".", "none", factory(), hmacAlgorithm{digest: factory})
	signed, err := signer.SignB64(payload)
	if err != nil {
		return "", err
	}
	return string(signed), nil
}
