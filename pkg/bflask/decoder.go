package bflask

import "fmt"

func DecodeCookie(raw string) ([]byte, error) {
	cookie, err := ParseCookie(raw)
	if err != nil {
		return nil, err
	}

	payload, err := DecodePayload(cookie)
	if err != nil {
		return nil, fmt.Errorf("decode cookie: %w", err)
	}
	return payload, nil
}
