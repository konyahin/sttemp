.POSIX:
.SUFFIXES:

BIN    = sttemp
CC     = cc
CFLAGS = -Wall -Werror -O

all: $(BIN)

strings.o: src/strings.c src/strings.h
	$(CC) $(CFLAGS) -c src/strings.c

files.o : src/files.c src/files.h
	$(CC) $(CFLAGS) -c src/files.c

main.o : src/main.c src/config.h
	$(CC) $(CFLAGS) -c src/main.c

$(BIN): main.o files.o strings.o
	$(CC) main.o files.o strings.o -o $(BIN)

clean:
	rm $(BIN)
	rm *.o

install: $(BIN)
	cp sttemp /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/sttemp
