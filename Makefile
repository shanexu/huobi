.DEFAULT_GOAL := build

build:
	go build -o bin/huobi

.PHONY: clean
clean:
	@rm -rf bin
