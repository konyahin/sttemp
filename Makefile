.POSIX:
.SUFFIXES:
.PHONY: clean install uninstall test 

BIN       = sttemp
CFLAGS    = -Wall -Werror -Os
LDFLAGS   = -s

PREFIX    = /usr/local
MANPREFIX = ${PREFIX}/share/man

all: $(BIN) README.md

strings.o: src/strings.c src/strings.h
	$(CC) $(CFLAGS) -c src/strings.c

files.o : src/files.c src/files.h
	$(CC) $(CFLAGS) -c src/files.c

token.o: src/token.c src/token.h
	$(CC) $(CFLAGS) -c src/token.c

main.o : src/main.c src/config.h
	$(CC) $(CFLAGS) -c src/main.c

$(BIN): main.o files.o strings.o token.o
	$(CC) main.o files.o strings.o token.o -o $(BIN)

README.md: $(BIN).1
	pandoc $(BIN).1 -o README.md

clean:
	rm -f $(BIN)
	rm -f *.o

install: $(BIN)
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	mkdir -p $(DESTDIR)$(MANPREFIX)/man1
	install -m 775 $(BIN) $(DESTDIR)$(PREFIX)/bin/
	install -m 644 $(BIN).1 $(DESTDIR)$(MANPREFIX)/man1/

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BIN)
	rm -f $(DESTDIR)$(MANPREFIX)/man1/$(BIN).1

test: $(BIN)
	./$(BIN) test && cat test && rm -f test

