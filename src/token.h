/* See LICENSE file for copyright and license details. */

#include <stdlib.h>
#include <string.h>

void free_tokens();
void add_new_token(char* name, char* value);
char* find_in_tokens(const char* name, size_t length);
