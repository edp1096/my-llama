#include <stdlib.h>
#include <string.h>

#include "llama.h"
#include "binding_api.h"

char* llama_get_system_info() {
    const char* result = llama_print_system_info();

    char* c_result = (char*)malloc(strlen(result) + 1);
    strcpy(c_result, result);

    return c_result;
}