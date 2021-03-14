.PHONY deps

deps:
	go get ./...

run: deps
  go run main.go

