.PHONY: all

all:
	cd capngoc && \
	go build && \
	./capngoc -pkg msgs ../msgs/*.capnp && \
	./capngoc -pkg test ../test/*.capnp && \
	go fmt ../msgs && \
	go fmt ../test && \
	go build ../test
