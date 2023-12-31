package main

import (
	"bufio"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

const (
	cwd           = "mimic"
	oldFileSuffix = ".old"
)

var embedfs embed.FS

type compoundReadCloser struct {
	closer     io.Closer
	readcloser io.ReadCloser
}

func (c *compoundReadCloser) Read(p []byte) (n int, err error) {
	return c.readcloser.Read(p)
}

func (c *compoundReadCloser) Close() error {
	if err := c.readcloser.Close(); err != nil {
		return err
	}
	if err := c.closer.Close(); err != nil {
		return err
	}
	return nil
}

func Extract(p string) (*compoundReadCloser, error) {
	f, err := embedfs.Open(path.Join(cwd, p))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", p)
	}
	r, err := gzip.NewReader(bufio.NewReader(f))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build reader")
	}
	return &compoundReadCloser{closer: f, readcloser: r}, nil
}

func main() {
	var src = "repro.exe"
	var dest = "D:\\dropgz\\repro-dropgz\\repro.exe"
	fmt.Println("Starting mimic")

	fmt.Println("Extracting file to replace with")
	rc, err := Extract(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rc.Close()
	fmt.Println("Extracted file")

	fmt.Println("Stat'ing file")
	if _, err := os.Stat(dest); err == nil {
		fmt.Println("File exists")
		fmt.Println("Renaming file")
		var newName = dest + oldFileSuffix
		if err = os.Rename(dest, newName); err != nil {
			fmt.Println(errors.Wrapf(err, "failed to rename the %s to %s", dest, newName))
			return
		}
	}

	fmt.Println("Open file")
	target, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755) //nolint:gomnd // executable file bitmask
	if err != nil {
		fmt.Println(errors.Wrapf(err, "failed to create file %s", dest))
		return
	}
	defer target.Close()
	fmt.Println("io.Copy (write) to file")
	_, err = io.Copy(bufio.NewWriter(target), rc)
	fmt.Println(errors.Wrapf(err, "failed to copy %s to %s", src, dest))
}
