/* See LICENSE file for copyright and license details. */

#include "files.h"


char* freadall(FILE* input) {
    char *buf = malloc(BUF_SIZE);
    size_t used = 0;
    size_t len  = 0;

    do {
        buf = realloc(buf, BUF_SIZE + used);
        len = fread(buf + used, 1, BUF_SIZE, input);
        used = used + len;
    } while (len != 0);

    buf = realloc(buf, used + 1);
    buf[used] = '\0';

    return buf;
}

char* freadline(FILE *input) {
    char *buf = malloc(BUF_SIZE);
    size_t used = 0;

    while(1) {
        buf = realloc(buf, BUF_SIZE + used);
        if (fgets(buf + used, BUF_SIZE, input) == NULL) {
            break;
        }
        char *new_line = strchr(buf + used, '\n');
        if (new_line != NULL) {
            used = new_line - buf;
            break;
        } else {
            used += BUF_SIZE - 1;
        }
    }

    buf = realloc(buf, used + 1);
    buf[used] = '\0';

    return buf;
}

