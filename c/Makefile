.PHONY: all clean

all: capn.so

clean:
	rm -f *.o *.so

%.o: %.c *.h *.inc
	$(CC) -Wall -Werror -g -O2 -c $< -o $@

capn.so: capn-malloc.o capn-stream.o capn.o
	$(CC) -shared -Wall -Werror -fPIC -g -O2 $^ -o $@
