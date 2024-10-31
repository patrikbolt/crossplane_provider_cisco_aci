.PHONY: generate build test

generate:
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./apis/..."

build:
	go build ./...

test:
	go test ./...

