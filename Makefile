.POSIX:
.SUFFIXES:

BIN    = sttemp
CC     = cc
CFLAGS = -Wall -Werror -O

all: $(BIN)

files.o : src/files.c src/files.h
	$(CC) $(CFLAGS) -c src/files.c

main.o : src/main.c src/config.h
	$(CC) $(CFLAGS) -c src/main.c

$(BIN): main.o files.o
	$(CC) $(CFLAGS) main.o files.o -o $(BIN)

clean:
	rm $(BIN)
	rm *.o
