package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func isSafePath(name string) bool {
	for _, part := range strings.Split(name, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func extractFile(e entry, src *os.File, dest string) error {
	if _, err := src.Seek(int64(e.offset), io.SeekStart); err != nil {
		return err
	}
	buf := make([]byte, len(e.path)+2)
	if _, err := src.Read(buf); err != nil {
		return nil
	}
	expect := make([]byte, len(e.path)+2)
	copy(expect, buf)
	if !bytes.Equal(buf, expect) {
		return errors.New("out of sync with file")
	}
	if dest == "" {
		return nil
	}
	fp, err := os.Create(dest)
	if err != nil {
		return err
	}
	exe := exec.Command("uncompress")
	exe.Stdin = &io.LimitedReader{R: src, N: int64(e.cmpsize)}
	exe.Stdout = fp
	exe.Stderr = os.Stderr
	if err := exe.Run(); err != nil {
		return err
	}
	return nil
}

func extractDirectory(e entry, dest string) error {
	if dest == "" {
		return nil
	}
	return os.Mkdir(dest, 0777)
}

func extractLink(e entry, dest string) error {
	if dest == "" {
		return nil
	}
	return os.Symlink(e.symval, dest)
}

func extractEntry(e entry, src *os.File, dest string) error {
	name := path.Clean(e.path)
	if !isSafePath(name) {
		return errors.New("invalid path")
	}
	if dest != "" {
		dest = path.Join(dest, name)
		if err := os.MkdirAll(path.Dir(dest), 0777); err != nil {
			return err
		}
	}
	switch e.ty {
	case 'f':
		return extractFile(e, src, dest)
	case 'd':
		return extractDirectory(e, dest)
	case 'l':
		return extractLink(e, dest)
	default:
		return fmt.Errorf("unknown type: %q", e.ty)
	}
}

func extract(entries []entry, src, dest string) error {
	fp, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fp.Close()
	for _, e := range entries {
		if err := extractEntry(e, fp, dest); err != nil {
			return fmt.Errorf("%s: %v", e.path, err)
		}
	}
	return nil
}
