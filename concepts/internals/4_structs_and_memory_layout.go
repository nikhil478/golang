package internals

import (
	"fmt"
	"unsafe"
)

// A struct is a collection of fields grouped together
// Struct is: Just a contigous block of memory with offsets

type User struct {
	ID   int
	Name string
}

func StructValueType() {
	a := User{ID: 1}
	b := a
	a.ID = 2
	fmt.Printf("a val : %v , b val : %v \n", a, b)
}

// Memory Layout
// Struct fields are stored contiguously in memory

type A struct {
	x int32
	y int64
}

// | x (4 bytes) | padding (4 bytes) | y (8 bytes) |

// Padding & Alignment (CRITICAL)
// Go aligns fields for performance

// Rules:
// * Each field is aligned to its natural boundary
// * Struct size is multiple of largest field alignment

type Bad struct {
	a bool  // 1 byte
	b int64 // 8 bytes
}

// Layout: | a (1) | padding (7) | b (8) | → total = 16 bytes

type Good struct {
	b int64
	a bool
}

// Layout: | b (8) | a (1) | padding (7) | → total = 16 bytes

// Field Ordering Optimization
// 	 Always order fields: largest → smallest

// Why?
// 	* Reduces padding
// 	* Improves cache efficiency

// Memory Alignment Example (Important)
type Example struct {
    a int8
    b int64
    c int8
}

// Layout: | a (1) | padding (7) | b (8) | c (1) | padding (7) |
// 👉 Total = 24 bytes

// 🔥 Optimized:
type Better struct {
    b int64
    a int8
    c int8
}
// 👉 Total = 16 bytes

// 🔹 17. Cache Efficiency
// 👉 Struct layout affects:
// * CPU cache usag

// Nested struct also follows alignment rules

func StructSize() {
	println(unsafe.Sizeof(User{}))
}

// Struct vs Pointer

// Value:
func FS(u User) {}

// 👉 Copies entire struct

// Pointer:
func FP(u *User) {}

// 👉 Passes reference

// Struct is comparable if and only if all the value inside is comparable
// ❌ Not allowed: u1 == u2
// type A struct {
//     s []int
// }

// Struct Tags
// Used by refelection and json/encoding
type ST struct {
	Name string `json:"name"`
}

// Anonymous Struct
func AnonymousStruct() {
	a := struct {
		Name string
	}{
		Name: "Go",
	}
	fmt.Println(a)
}
