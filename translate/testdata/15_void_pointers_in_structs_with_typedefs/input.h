typedef unsigned char TypeB;

// This is what PKCS #11 does and it tripped up the original logic
typedef void *VOID_PTR;

typedef struct {
  unsigned long len;
  VOID_PTR ptr;
} TypeA;

int functionA(TypeA a1, long a2);
