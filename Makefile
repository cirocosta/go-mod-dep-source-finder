all: install

install:
	go install -v

test:
	go test ./...

fmt:
	go fmt ./...
