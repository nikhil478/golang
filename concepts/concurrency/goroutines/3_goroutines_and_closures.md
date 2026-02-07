## Goroutines and Closure

1. Goroutines execute within the same address space they are created in

2. they can directly modify variables in the enclosing lexical block 
    go compiler and the runtime takes care of pinning the variable , moving the variable from stack to heap, to facilitate goroutines, to have access to the variables even after the enclosing function has return