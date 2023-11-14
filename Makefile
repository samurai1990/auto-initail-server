.PHONY: build 


all:build

build:
	-rm auto-initail-server
	@go build -o auto-initail-server
