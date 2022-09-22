APP_NAME=wr-latest-daily-redirect

.PHONY: all
all: clean build

.PHONY: clean
clean:
	@rm -f ../bin/${APP_NAME}

build:
	@env GOOS=linux GOARCH=amd64 go build -o ../bin/${APP_NAME} ./cmd/main.go
