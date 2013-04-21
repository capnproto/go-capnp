.PHONY: all

all:
	go fmt && \
	cd capngoc && \
	go fmt && \
	go build && \
	./capngoc -pkg msgs ../msgs/*.capnp && \
	./capngoc -pkg test ../test/*.capnp && \
	go fmt ../msgs && \
	go fmt ../test && \
	go build ../test
