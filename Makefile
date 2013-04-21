.PHONY: all

all:
	cd capngoc && \
	go build && \
	./capngoc -pkg msgs ../msgs/*.capnp && \
	./capngoc -pkg test ../test/*.capnp
