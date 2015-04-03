#!/bin/bash
set -ev

# Install gcc
sudo apt-get install -qq g++-4.8 libstdc++-4.8-dev
sudo update-alternatives --quiet --install /usr/bin/gcc  gcc  /usr/bin/gcc-4.8  60 --slave   /usr/bin/g++  g++  /usr/bin/g++-4.8 --slave   /usr/bin/gcov gcov /usr/bin/gcov-4.8
sudo update-alternatives --quiet --set gcc /usr/bin/gcc-4.8

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
