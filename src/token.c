/* See LICENSE file for copyright and license details. */

#include "token.h"

struct token {
    char* name;
    char* value;
};
typedef struct token Token;

Token* new_token(char* name, char* value) {
    Token* token = malloc(sizeof(Token));
    token->name = name;
    token->value = value;
    return token;
}

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

void add_token(Token* token) {
    tokens = realloc(tokens, sizeof(Token) * ++tokens_len);
    tokens[tokens_len - 1] = token;
}

void add_new_token(char* name, char* value) {
    Token* token = new_token(name, value);
    add_token(token);
}

char* find_in_tokens(const char* name, size_t length) {
    // O(n) = n, but I don't worry about it right now
    for (size_t i = 0; i < tokens_len; i++) {
        if (strncmp(tokens[i]->name, name, length) == 0) {
            return tokens[i]->value;
        }
    }
    return NULL;
}

