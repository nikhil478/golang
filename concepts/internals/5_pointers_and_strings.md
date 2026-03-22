## Pointers
A pointer stores the memory address of a value

To modify value
*p = 20

# When to Use Pointers?

1. Avoid copying large structs
2. Modify original data
3. Share data between functions
4. For optional values (nil check)

# When NOT to Use

1. Small structs 
2. Simple values (int, bool) 
3. Over-optimization

# Pointer Internals

1. 8 bits on 64-bit system


## String
A string is a read-only slice of bytes

1. Strings are not mutable

# Internal Representation

type string struct {
    data *byte
    len  int
}

Feature	    string	        []byte
Mutable	    ❌	            ✅
Copy	    cheap (header)	depends
Use case	read-only	    modification

#  String Slicing
sub := s[0:5]
sub → points to SAME underlying data

String slicing does NOT copy data

# Memory Retention Issue
big := "very large string..."
small := big[:5]
👉 small keeps reference to entire big string

✅ Use:
var b strings.Builder

# Strings are UTF-8
len("你好") // 6
👉 bytes, not characters

# Rune handling:
for _, r := range s {}

# String Conversion
[]byte(s)
string(b)
👉 usually creates copy

# String Comparison
s1 == s2
👉 compares bytes