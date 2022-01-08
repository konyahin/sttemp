.POSIX:
.SUFFIXES:
.PHONY: clean uninstall test 

BIN    = sttemp
CC     = cc
CFLAGS = -Wall -Werror -O

all: $(BIN)

strings.o: src/strings.c src/strings.h
	$(CC) $(CFLAGS) -c src/strings.c

files.o : src/files.c src/files.h
	$(CC) $(CFLAGS) -c src/files.c

token.o: src/token.c src/token.h
	$(CC) $(CFLAGS) -c src/token.c

main.o : src/main.c
	$(CC) $(CFLAGS) -c src/main.c

$(BIN): main.o files.o strings.o token.o
	$(CC) main.o files.o strings.o token.o -o $(BIN)

clean:
	rm $(BIN)
	rm *.o

install: $(BIN)
	cp sttemp /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/sttemp

test: $(BIN)
	./sttemp test && cat test && rm -f test

