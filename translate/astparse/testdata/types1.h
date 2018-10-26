struct Struct1 {
    int a;
    long b;
};

typedef struct {
  struct Struct1 a;
  int *b;
} Type2;

typedef struct Struct1 Type1, *PtrType1;

typedef struct Struct3 {
  PtrType1 a;
} Type3;

typedef Type3 *PtrType3;

typedef PtrType3 AnotherPtrType3;


void functionA(PtrType3 a1, long a2);

char *functionB(struct Struct3 a1, Type2 *a2);