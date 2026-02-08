package main

import (
	"io"
)

type OutputFile interface {
	io.Writer
	io.Closer
}

type Stdout struct {
	writer io.Writer
}

func (s *Stdout) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

func (*Stdout) Close() error {
	return nil
}

func StdoutInstance(writer io.Writer) OutputFile {
	return &Stdout{writer}
}
