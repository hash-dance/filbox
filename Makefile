SHELL=/usr/bin/env bash

all: build
.PHONY: all

BINS:=

filbox:
	go build -mod=vendor -o filbox main.go
BINS+=filbox

build: $(BINS)

clean:
	rm -f $(BINS)
