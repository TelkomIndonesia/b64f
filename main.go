package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/goreleaser/fileglob"
)

const extension = ".b64"

func main() {
	base := os.DirFS(".")
	isDecode := true
	if len(os.Args) > 1 && os.Args[1] == "encode" {
		isDecode = false
	}

	if err := b64f(base, isDecode); err != nil {
		log.Fatal(err)
	}
}

func b64f(base fs.FS, isDecode bool) (err error) {
	files, err := listFiles(base, isDecode)
	if err != nil {
		return err
	}

	for _, f := range files {
		var err error
		if isDecode {
			err = decodeFile(base, f)
		}

		if !isDecode && !strings.HasSuffix(f, extension) {
			err = encodeFile(base, f)
		}

		if err != nil {
			return err
		}
	}
	return
}

func listFiles(base fs.FS, isDecode bool) (files []string, err error) {
	patterns, err := listPatternsFromFile(base)
	if err != nil {
		log.Printf("[WARN] failed to list from .b64f file with err: %v.\nWill try to read from stdin instead.", err)
		patterns, err = listPatternsFromStdin()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load patterns: %w", err)
	}

	for _, pattern := range patterns {
		if isDecode {
			pattern = pattern + extension
		}

		f, err := fileglob.Glob(pattern, fileglob.WithFs(base))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("failed to find files based on pattern '%s': %w", pattern, err)
		}
		files = append(files, f...)
	}

	return
}

func listPatternsFromFile(base fs.FS) (patterns []string, err error) {
	b, err := base.Open(".b64f")
	if err != nil {
		return nil, err
	}

	return listPatternsFromReader(b)
}

func listPatternsFromStdin() (patterns []string, err error) {
	return listPatternsFromReader(os.Stdin)
}

func listPatternsFromReader(r io.Reader) (patterns []string, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		patterns = append(patterns, scanner.Text())
	}
	return patterns, scanner.Err()
}

func decodeFile(base fs.FS, f string) (err error) {
	r, err := base.Open(f)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", f, err)
	}
	defer r.Close()

	dec := strings.TrimSuffix(f, extension)
	w, err := os.OpenFile(dec, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for write: %w", f, err)
	}
	defer w.Close()

	rB64 := base64.NewDecoder(base64.StdEncoding.WithPadding(base64.StdPadding), r)
	if _, err := io.Copy(w, rB64); err != nil {
		return fmt.Errorf("failed to decode file '%s' to '%s': %w", f, dec, err)
	}

	return
}

func encodeFile(base fs.FS, f string) (err error) {
	r, err := base.Open(f)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", f, err)
	}
	defer r.Close()

	enc := f + extension
	w, err := os.OpenFile(enc, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file '%s' for write: %w", f, err)
	}
	defer w.Close()

	wB64 := base64.NewEncoder(base64.StdEncoding, w)
	defer wB64.Close()
	if _, err := io.Copy(wB64, r); err != nil {
		return fmt.Errorf("failed to decode file '%s' to '%s': %w", f, enc, err)
	}

	return
}
