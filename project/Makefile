VERSION=$(shell git describe --tags --always)

.PHONY: graph
graph:
	go run github.com/99designs/gqlgen generate  

.PHONY: trans 
trans:
	go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
	go mod tidy
	goi8n extract; goi18n merge active.en.toml translate.zh.toml

.PHONY: wire
wire:
	go install github.com/google/wire/cmd/wire@latest
	go mod tidy
	cd ./app ; wire

.PHONY: build
build:
	mkdir -p build/ && \
	go mod tidy && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./build/app ./cmd