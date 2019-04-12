VERSION = $(shell cat ./VERSION)
NAME = go-mod-license-finder


all: install

install:
	go install -v

test:
	go test -v ./...

fmt:
	go fmt ./...

image:
	docker build -t cirocosta/$(NAME):$(VERSION) .
	docker tag cirocosta/$(NAME):$(VERSION) ciroocsta/$(NAME):latest

