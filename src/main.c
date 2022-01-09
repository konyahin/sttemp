/* See LICENSE file for copyright and license details. */

#include "config.h"
#include "files.h"
#include "strings.h"
#include "token.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

const int pat_start_len = sizeof(pattern_start) / sizeof(pattern_start[0]) - 1;
const int pat_end_len = sizeof(pattern_end) / sizeof(pattern_end[0]) - 1;

enum search {WO_ENV, W_ENV};
typedef enum search Search;

void show_usage() {
    printf("sttemp - simple template manager\n");
    printf("Usage:\n\tsttemp template_name [file_name]\n");
}

FILE* open_template(const char* template_name) {
    const char* dir;
    if (template_dir[0] == '~') {
        char* homedir = getenv("HOME");
        dir = strconcat(homedir, template_dir + 1);
    } else {
        dir = template_dir;
    }
    char *template_path = strconcat(dir, template_name);
    FILE *template = fopen(template_path, "rb");
    free(template_path);
    return template;
}

char* get_placeholder_value(Search search_type, const char* name, size_t length) {
    char* value = find_in_tokens(name, length);
    if (value != NULL) {
        return value;
    }
    char* new_name = strndup(name, length);
    if (search_type == W_ENV) {
        value = getenv(new_name);
        if (value != NULL) {
            free(new_name);
            return value;
        }
    }

    printf("Enter value for {%.*s}: ", (int) length, name);
    value =  freadline(stdin);
    add_new_token(new_name, value);
    return value;
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        show_usage();
        return 1;
    }

    if (strcmp("-h", argv[1]) == 0) {
        show_usage();
        return 0;
    }

    size_t first_file_arg = 1;
    Search search_type = WO_ENV;
    if (strcmp("-e", argv[1]) == 0) {
        search_type = W_ENV;
        first_file_arg++;
    }
    
    char* template_name = argv[first_file_arg++];
    char* target_name;
    if (first_file_arg < argc) {
        target_name  = argv[first_file_arg];
    } else {
        target_name = template_name;
    }

    FILE* template = open_template(template_name);
    if (template == NULL) {
        fprintf(stderr, "Template doesn't exist: %s\n", template_name);
        return 1;
    }

    size_t buf_len = 0;
    char *buf = freadall(template, &buf_len);
    fclose(template);

    FILE* output = fopen(target_name, "w");

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
        char *value = get_placeholder_value(search_type, start, token_length);
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
