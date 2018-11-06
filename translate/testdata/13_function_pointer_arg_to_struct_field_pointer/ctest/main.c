#include "../input.h"
#include <stdio.h>

int main() {
    int x = 0;

    Buffer b = {.buffer = &x, .len = 1};
    functionA(&b);
    printf("Should be 42: %d\n", x);
}