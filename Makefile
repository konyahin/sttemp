.POSIX:
.SUFFIXES:

BIN    = sttemp
CC     = cc
CFLAGS = -W -O

all: src/main.c src/config.h
	$(CC) $(CFLAGS) src/main.c -o $(BIN)

clean:
	rm $(BIN)
