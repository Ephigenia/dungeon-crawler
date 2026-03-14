.PHONY: init install run build

init:
	go mod init github.com/ephigenia/ebit-engine-game-1

install:
	go mod download

run:
	go run .

build:
	go build -o dungeon-crawler .
