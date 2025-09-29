build:
	go build -o bin/app .

run:
	go run .

test:
	go test ./... -v
