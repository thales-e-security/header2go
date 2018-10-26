#include "../input.h"
#include <stdio.h>

int main() {
    TypeA a = {0};
    functionA(&a, 0);
    printf("Should be 42: %d\n", a.a);
}