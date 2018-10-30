typedef struct {
    int val3;
} StructB;

struct StructA {
    int val1;
    StructB val2;
};

void functionA(struct StructA a1, long a2);