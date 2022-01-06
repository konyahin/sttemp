#include "strings.h"


char* strconcat(const char* first, const char* second) {
    size_t first_len = strlen(first);
    size_t second_len = strlen(second);
    char *buf = malloc(first_len + second_len + 1);
    memcpy(buf, first, first_len);
    memcpy(buf + first_len, second, second_len + 1);
    return buf;
}
