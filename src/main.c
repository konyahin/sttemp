/* See LICENSE file for copyright and license details. */

#include "config.h"
#include "files.h"
#include "strings.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

void show_usage() {
    printf("sttemp - simple template manager\n");
    printf("Usage:\n\tsttemp template_name file_name\n");
}

FILE* open_template(const char* template_name) {
    char *template_path = strconcat(template_dir, template_name);
    FILE *template = fopen(template_path, "rb");
    free(template_path);
    return template;
}

struct token {
    char* name;
    char* value;
};
typedef struct token Token;

void free_token(Token* token) {
    free(token->name);
    free(token->value);
    free(token);
}

Token** tokens = NULL;
size_t tokens_len = 0;

void free_tokens() {
    for (size_t i = 0; i < tokens_len; i++) {
        free_token(tokens[i]);
    }
    free(tokens);
    tokens = NULL;
}

char* get_placeholder_value(const char* name, size_t length) {
    for (size_t i = 0; i < tokens_len; i++) {
        if (strncmp(tokens[i]->name, name, length) == 0) {
            return tokens[i]->value;
        }
    }
    printf("Enter value for {%.*s}: ", (int) length, name);
    char* value =  freadline(stdin);

    Token* token = malloc(sizeof(Token));
    token->name = malloc(length);
    strncpy(token->name, name, length);
    token->value = value;

    tokens = realloc(tokens, sizeof(Token) * ++tokens_len);
    tokens[tokens_len - 1] = token;
    return value;
}

int main(int argc, char *argv[]) {
    if (argc < 3) {
        show_usage();
        return 1;
    }

    char* template_name = argv[1];
    FILE* template = open_template(template_name);
    if (template == NULL) {
        fprintf(stderr, "Template doesn't exist: %s\n", template_name);
        return 1;
    }

    size_t buf_len = 0;
    char *buf = freadall(template, &buf_len);
    fclose(template);

    const int pat_start_len = strlen(pattern_start);
    const int pat_end_len = strlen(pattern_end);

    FILE* output = fopen(argv[2], "w");

    char *start = buf;
    char *last = start;
    while ((start = strstr(start, pattern_start)) != NULL) {
        fwrite(last, sizeof(char), start - last, output);
        start = start + pat_start_len;

        char* end = strstr(start, pattern_end);
        if (end == NULL) {
            fprintf(stderr, "Unfinished pattern: %10s", start);
            fclose(output);
            free(buf);
            return 1;
        }

        size_t token_length = end - start;
        char *value = get_placeholder_value(start, token_length);
        fwrite(value, sizeof(char), strlen(value), output);

        start = end + pat_end_len;
        last = start;
    }

    fwrite(last, sizeof(char), buf_len - (last - buf), output);
    fclose(output);
    free(buf);
    free_tokens();

    return 0;
}
