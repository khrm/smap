.Phony: all

all: build


build:
	go build -a -ldflags "-s -w" -o smap cmd/smap/*.go

image:
	docker build .  -t smap:latest

