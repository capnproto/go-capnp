#!/bin/bash
set -ev

# Install go-capnproto
export GOPATH="$HOME/gopath"
mkdir -p "$GOPATH/src/zombiezen.com/go"
mv "$TRAVIS_BUILD_DIR" "$GOPATH/src/zombiezen.com/go/capnproto2"
go get -v -t -d zombiezen.com/go/capnproto2/...
