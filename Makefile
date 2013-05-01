.PHONY: all

all:
	go fmt && \
	cd capngoc && \
	go fmt && \
	go build && \
	./capngoc -pkg test ../test/*.capnp && \
	./capngoc -lang=c -pkg test ../test/*.capnp && \
	go fmt ../test && \
	go build ../test
