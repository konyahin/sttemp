package main

import (
	"io"
)

type OutputFile interface {
	io.Writer
	io.Closer
}

type Stdout struct {
	ioh *IOHandler
}

func (s *Stdout) Write(p []byte) (n int, err error) {
	return s.ioh.Write(p)
}

func (*Stdout) Close() error {
	return nil
}

func StdoutInstance(ioh *IOHandler) OutputFile {
	return &Stdout{ioh}
}
