.PHONY: all clean test

LDFLAGS=-g -Wall -Werror -fPIC
CFLAGS=-g -Wall -Werror -fPIC -I. -Wno-unused-function -ansi -pedantic

all: capn.so capnpc-c test

clean:
	rm -f *.o *.so capnpc-c compiler/*.o

%.o: %.c *.h *.inc compiler/*.h
	$(CC) $(CFLAGS) -c $< -o $@

capn.so: capn-malloc.o capn-stream.o capn.o
	$(CC) -shared $(LDFLAGS) $^ -o $@

capnpc-c: compiler/capnpc-c.o compiler/schema.capnp.o compiler/str.o capn.so
	$(CC) $(LDFLAGS) $^ -o $@

test: capn-test
	./capn-test

%-test.o: %-test.cpp *.h *.c *.inc
	$(CXX) -g -Wall -Werror -I. `gtest-config --cppflags --cxxflags` -o $@ -c $<

capn-test: capn-test.o capn-stream-test.o compiler/schema-test.o compiler/schema.capnp.o
	$(CXX) -g -Wall -Werror -I. `gtest-config --ldflags --libs` -o $@ $^
