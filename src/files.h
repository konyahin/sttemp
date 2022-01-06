/* See LICENSE file for copyright and license details. */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#define BUF_SIZE 2 * 1024 * 1024

/**
 * Read line from file `input` and return it content
 * without new line symbol.
 */
char* freadline(FILE *input);

/**
 * Read all content of file `input` and return it.
 */
char* freadall(FILE* input, size_t* length);
