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
  v go install github.com/mattn/goveralls
fi

if ! type -p shadow; then
  v go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
fi

echo period...
v go test -v -covermode=count -coverprofile=period.out .
v go tool cover -func=period.out
#[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=period.out -service=travis-ci -repotoken $COVERALLS_TOKEN

v gofmt -l -w *.go

v go vet ./...

v shadow ./...
