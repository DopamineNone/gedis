.PHONY: build

build:
	@wire ./cmd/wire.go && go build -o ./bin/main ./cmd