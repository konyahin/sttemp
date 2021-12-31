#include "config.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

void show_usage() {
    printf("sttemp - simple template manager\n");
    printf("Usage:\n\tsttemp template_name\n");
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

    printf("Template directory: %s\n", template_dir);
    printf("Template: %s var %s\n", pattern_start, pattern_end);

    size_t temp_len = strlen(template_name) + strlen(template_dir) + 1;
    char *temp_path = (char*) malloc(temp_len);
    strcpy(temp_path, template_dir);
    strcat(temp_path, template_name);

    FILE *template = fopen(temp_path, "r");
    if (template == NULL) {
        fprintf(stderr, "Template doesn't exist: %s", temp_path);
        return 1;
    }

    fclose(template);
    free(temp_path);

    return 0;
}
