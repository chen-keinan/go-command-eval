#!/usr/bin/env bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
## need to generate mock before lint
go install github.com/golang/mock/mockgen@latest
go install -v github.com/golang/mock/mockgen
export PATH=$GOPATH/bin:$PATH
go generate ./...
$GOPATH/bin/golangci-lint run -v