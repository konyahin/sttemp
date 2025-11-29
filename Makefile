.POSIX:
.SUFFIXES:
.PHONY: all test install clean

all: check sttemp test install

sttemp: *.go
	gofmt -w . 
	go build

test: sttemp
	go test

install: sttemp
	go install .

check:
	go vet
	staticcheck

clean:
	rm sttemp
