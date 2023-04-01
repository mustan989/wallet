.PHONY: all test clean

run:
	go run app/main.go

test:
	go test -v ./...

coverage:
	go test -coverprofile=test/coverage.out ./...
	go tool cover -html=test/coverage.out
