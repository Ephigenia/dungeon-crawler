.PHONY: install run build

GO			:= $(shell which go)

install:
	@$(GO) version
	@$(GO) mod download
	@$(GO) mod tidy

run:
	go run .

build:
	go build -o dungeon-crawler .
