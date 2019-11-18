SHELL=/bin/bash -o pipefail
PKG_LIST := $(shell go list ./...)

fmt:
		@go fmt ${PKG_LIST}

lint:
		@go get golang.org/x/lint/golint
		@golint -set_exit_status ${PKG_LIST}

test:
		@go test -coverpkg=./... -coverprofile=c.out ${PKG_LIST}  -race
		@go tool cover -func=c.out | grep 'total:'
