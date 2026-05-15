package bflask

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Cookie struct {
	Raw        string
	Payload    string
	Timestamp  string
	Signature  string
	Compressed bool
}

func ParseCookie(raw string) (Cookie, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Cookie{}, errors.New("cookie is empty")
	}

	parts := strings.Split(raw, ".")
	if strings.HasPrefix(raw, ".") {
		if len(parts) != 4 || parts[1] == "" || parts[2] == "" || parts[3] == "" {
			return Cookie{}, errors.New("invalid compressed Flask cookie format")
		}
		return Cookie{
			Raw:        raw,
			Payload:    "." + parts[1],
			Timestamp:  parts[2],
			Signature:  parts[3],
			Compressed: true,
		}, nil
	}

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return Cookie{}, errors.New("invalid Flask cookie format, expected payload.timestamp.signature")
	}

	return Cookie{Raw: raw, Payload: parts[0], Timestamp: parts[1], Signature: parts[2]}, nil
}

func DecodePayload(c Cookie) ([]byte, error) {
	payload := c.Payload
	compressed := c.Compressed || strings.HasPrefix(payload, ".")
	if strings.HasPrefix(payload, ".") {
		payload = payload[1:]
	}

	decoded, err := rawURLBase64Decode(payload)
	if err != nil {
		return nil, err
	}

	if !compressed {
		return decoded, nil
	}

	reader, err := zlib.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func PrettyPayload(data []byte) string {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return string(data)
	}
	pretty, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(data)
	}
	return string(pretty)
}

func rawURLBase64Decode(s string) ([]byte, error) {
	for len(s)%4 != 0 {
		s += "="
	}
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	return b, nil
}
