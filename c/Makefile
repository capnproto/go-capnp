.PHONY: all clean test

CFLAGS=-g -Wall -Werror -fPIC -I. -Wno-unused-function

all: capn.so test

clean:
	rm -f *.o *.so capnpc-c compiler/*.o

%.o: %.c *.h *.inc
	$(CC) $(CFLAGS) -c $< -o $@

capn.so: capn-malloc.o capn-stream.o capn.o
	$(CC) -shared $(CFLAGS) $^ -o $@

test: capn-test
	./capn-test

%-test.o: %-test.cpp *.h *.c *.inc
	$(CXX) `gtest-config --cppflags --cxxflags` -o $@ -c $<

capn-test: capn-test.o
	$(CXX) `gtest-config --ldflags --libs` -o $@ $^
