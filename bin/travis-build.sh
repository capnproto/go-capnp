#!/bin/bash
set -ev

# Install gcc 4.7
sudo apt-get install -qq gcc-4.7 g++-4.7
export CXX=g++-4.7

# Install capnp
cd "$HOME"
wget -O capnproto.tar.gz https://capnproto.org/capnproto-c++-0.5.1.2.tar.gz
tar zxf capnproto.tar.gz
cd capnproto-c++-0.5.1.2
./configure && make -j6 check
sudo make install

# Install go-capnproto
export GOPATH="$HOME/gopath"
mkdir -p "$GOPATH/src/zombiezen.com/go"
mv "$TRAVIS_BUILD_DIR" "$GOPATH/src/zombiezen.com/go/capnproto"
go get -v -t -d zombiezen.com/go/capnproto
