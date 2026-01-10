package main

import (
	"io"
	"os"
)

type OutputFile interface {
	io.Writer
	io.Closer
}

type Stdout struct{}

var stdoutInstance = &Stdout{}

func (*Stdout) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (*Stdout) Close() error {
	return nil
}

func StdoutInstance() OutputFile {
	return stdoutInstance
}
