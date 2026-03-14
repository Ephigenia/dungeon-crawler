.PHONY: install run build

install:
	go mod download

run:
	go run .

build:
	go build -o dungeon-crawler .
