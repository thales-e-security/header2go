# Test Instructions

These C tests are very manual at present. They are designed to be used during development to double-check the
generated Go code actually does what we think it does.

For this test:

1. Replace the contents of `goFunctionA` with:

    ```go
    a1v := a1.Slice(1)
    a1v[0].a = 42
    return 0
    ```
2. Build the shared object file:

    ```
    go build -o testbuild.so -buildmode=c-shared .
    ```
   
3. Change directory to `ctest` and build the C code:

     ```
     gcc -o main main.c ../testbuild.so
     ```
    
4. Run the test, it should print this:

    ```
    $ ./main 
    Should be 42: 42
    ```