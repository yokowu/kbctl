.PHONY: build

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kbctl cmd/main.go
	@echo "------ build go success ------"