.PHONY: build
build:
	go build -v ./cmd/todo-list

.DEFAULT_GOAL := build

