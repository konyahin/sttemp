.POSIX:
.SUFFIXES:

BIN    = sttemp
CC     = cc
CFLAGS = -Wall -Werror -O

all: $(BIN)

$(BIN): src/main.c src/config.h
	$(CC) $(CFLAGS) src/main.c -o $(BIN)

clean:
	rm $(BIN)
