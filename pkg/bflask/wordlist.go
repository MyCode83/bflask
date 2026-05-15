package bflask

import (
	"bufio"
	"context"
	"os"
	"strings"
)

func CountCandidates(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count int64
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			count++
		}
	}
	return count, scanner.Err()
}

func StreamCandidates(ctx context.Context, path string, out chan<- string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		candidate := strings.TrimRight(scanner.Text(), "\r\n")
		if strings.TrimSpace(candidate) == "" {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- candidate:
		}
	}
	return scanner.Err()
}
