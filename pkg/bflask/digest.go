package bflask

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
)

type DigestFactory func() hash.Hash

func Digest(name string) (hash.Hash, error) {
	factory, err := DigestFactoryFor(name)
	if err != nil {
		return nil, err
	}
	return factory(), nil
}

func DigestFactoryFor(name string) (DigestFactory, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "sha1":
		return sha1.New, nil
	case "sha224":
		return sha256.New224, nil
	case "sha256":
		return sha256.New, nil
	case "sha384":
		return sha512.New384, nil
	case "sha512":
		return sha512.New, nil
	case "md5":
		return md5.New, nil
	default:
		return nil, fmt.Errorf("unsupported digest %q", name)
	}
}
