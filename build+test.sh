#!/bin/bash -e
cd "$(dirname $0)"
PATH=$HOME/go/bin:$PATH
unset GOPATH
export GOARCH=${1}

function v
{
  echo
  echo $@
  $@
}

if ! type -p goveralls; then
  v go get     github.com/mattn/goveralls
  v go install github.com/mattn/goveralls
  go mod tidy
fi

echo period...
v go test -v -covermode=count -coverprofile=period.out .
v go tool cover -func=period.out
#[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=period.out -service=travis-ci -repotoken $COVERALLS_TOKEN

v gofmt -l -w *.go

v go vet ./...
