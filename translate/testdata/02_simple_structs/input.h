typedef struct {
 char val1;
 signed char val2;
 unsigned char val3;
 short val4;
 unsigned short val5;
 int val6;
 unsigned int val7;
 long val8;
 unsigned long val9;
 long long val10;
 unsigned long long val11;
 float val12;
 double val13;
} StructB;

typedef struct _StructA {
    int val1;
    StructB val2;
} StructA;

void functionA(StructA a1, long a2);