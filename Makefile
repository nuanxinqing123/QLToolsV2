.PHONY: gen docs

run:
	go run cmd/main.go

gen:
	go run cmd/generate/generate.go

docs:
	swag init -g cmd/main.go -o docs

build:
	goreleaser build --clean --single-target --snapshot