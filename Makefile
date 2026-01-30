.POSIX:
.SUFFIXES:
.PHONY: all test install coverage clean

all: check sttemp test install

sttemp: *.go
	gofmt -w . 
	go build

test: sttemp
	go test ./...

install: sttemp
	go install .

check:
	go vet
	staticcheck

coverage:
	go test -coverprofile=coverage.out && go tool cover -func=coverage.out
	rm coverage.out

coverage-html:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	rm coverage.out

clean:
	rm sttemp
