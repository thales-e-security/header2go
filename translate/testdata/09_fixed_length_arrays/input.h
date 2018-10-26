typedef struct {
    int val3[10];
} StructB;

typedef struct _StructA {
    int val1;
    StructB val2[2];
} StructA;

void functionA(StructA a1, long a2);