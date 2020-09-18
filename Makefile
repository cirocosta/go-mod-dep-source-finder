VERSION = $(shell cat ./VERSION)
NAME = go-mod-dep-source-finder
PIPELINE_NAME = go-mod-dep-source-finder


all: install

install:
	go install -v

test:
	go test -v ./...

fmt:
	go fmt ./...

image:
	docker build -t cirocosta/$(NAME):$(VERSION) .
	docker tag cirocosta/$(NAME):$(VERSION) cirocosta/$(NAME):latest

push-image:
	docker push cirocosta/$(NAME):$(VERSION)
	docker push cirocosta/$(NAME):latest

pipeline:
	jsonnet \
		--ext-code 'repositories=$(shell cat ./ci/repositories.json)' \
		./ci/pipeline.jsonnet > ./ci/pipeline.json
