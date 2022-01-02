/* See LICENSE file for copyright and license details. */

#include "config.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

#define BUF_SIZE 2097152

void show_usage() {
    printf("sttemp - simple template manager\n");
    printf("Usage:\n\tsttemp template_name\n");
}

char* strconcat(const char* first, const char* second) {
    size_t first_len = strlen(first);
    size_t second_len = strlen(second);
    char *buf = malloc(first_len + second_len + 1);
    memcpy(buf, first, first_len);
    memcpy(buf + first_len, second, second_len + 1);
    return buf;
}

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

int main(int argc, char *argv[]) {
    if (argc < 2) {
        show_usage();
        return 1;
    }

    char* template_name = argv[1];
    if (strcmp("-h", template_name) == 0) {
        show_usage();
        return 0;
    }

    char *temp_path = strconcat(template_dir, template_name);
    FILE *template = fopen(temp_path, "rb");
    if (template == NULL) {
        fprintf(stderr, "Template doesn't exist: %s", temp_path);
        return 1;
    }
    free(temp_path);

    char *buf = freadall(template);
    fclose(template);

    printf("%s", buf);
    printf("==================================\n");

    const int pat_start_len = strlen(pattern_start);
    const int pat_end_len = strlen(pattern_end);

    char *start = buf;
    while ((start = strstr(start, pattern_start)) != NULL) {
        start = start + pat_start_len;
        char* end = strstr(start, pattern_end);
        if (end == NULL) {
            fprintf(stderr, "Unfinished pattern: %10s", start);
            return 1;
        }

        size_t token_length = end - start;
        char* token_name = malloc(token_length + 1);
        memcpy(token_name, start, token_length);
        token_name[token_length] = '\0';

        printf("%s\n", token_name);
        free(token_name);

        start = end + pat_end_len;
    }

    free(buf);

    return 0;
}
