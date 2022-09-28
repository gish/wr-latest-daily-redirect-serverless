APP_NAME=wr-latest-daily-redirect
GIT_COMMIT=`git rev-parse --short HEAD`

.PHONY: all
all: clean build

.PHONY: clean
clean:
	@rm -f ../bin/${APP_NAME}

build:
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.commit=${GIT_COMMIT}" -o ../bin/${APP_NAME} ./cmd/main.go
