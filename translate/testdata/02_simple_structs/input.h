typedef struct {
    int val3;
} StructB;

typedef struct _StructA {
    int val1;
    StructB val2;
} StructA;

void functionA(StructA a1, long a2);