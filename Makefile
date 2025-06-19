.PHONY: run build test test-cover clean

GO_PACKAGES=./...
APP_NAME=task-runner
MAIN_PACKAGE=./cmd/task-runner

run:
	go run $(MAIN_PACKAGE)

build: clean
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(MAIN_PACKAGE)

test:
	go test -v $(GO_PACKAGES)

test-cover:
	go test -cover -coverprofile=coverage.out $(GO_PACKAGES)
	go tool cover -func=coverage.out

clean:
	rm -rf bin coverage.out
