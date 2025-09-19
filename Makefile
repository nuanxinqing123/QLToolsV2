.PHONY: gen docs

gen:
	go run cmd/generate/generate.go

docs:
	swag init -g cmd/main.go -o cmd/docs
