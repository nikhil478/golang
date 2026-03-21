package internals

import (
	"fmt"
	"runtime"
)

// A slice is a dynamic → its size can change
// reference-like → it points to data instead of owning it
// data type → defined structure in Go
// represents a view → shows part of data
// over an underlying array → actual data is stored in an array

// Slice does not store data + It points to an array

// Internal Structure
// type slice struct {
//     ptr *T   // pointer to array
//     len int  // number of elements
//     cap int  // total capacity from ptr
// }

// slice creation
func SliceCreation() {

	// Using make : len = 3, cap = 4
	a := make([]int, 3, 4)

	// Using literals : len = 3, cap = 3
	s := []int{1, 2, 3}

	arr := [5]int{1, 2, 3, 4, 5}

	// Using Array: len =3, cap = 4
	slice := arr[1:4]

	fmt.Printf("a %v , s: %v, slice : %v", a, s, slice)
}

// slice append behavior
func AppendBehavior() {

	a := make([]int, 3, 4)
	b := append(a, 5)
	a = append(a, 6)
	fmt.Printf("a val : %v and b val %v", a, b) // a and b sharing same underlying array so overrides

	c := append(a, 5)
	a = append(a, 6)
	fmt.Printf("a val : %v and c val %v", a, c) // but here not as cap limit reached so for c new array allocation happens

	d := append(c, 7)
	c = append(c, 8)
	fmt.Printf("d val : %v and c val %v \n", d, c) // here it overrides again

	d = append(c, 7)
	c = append(c, 8)
	fmt.Printf("d val : %v and c val %v \n", d, c) // here it overrides again

	d = append(c, 7)
	c = append(c, 8)
	fmt.Printf("d val : %v and c val %v \n", d, c) // here it overrides again

	d = append(c, 7)
	c = append(c, 8)
	fmt.Printf("d val : %v and c val %v \n", d, c) // here it not overrides because of growth strategy

	// Growth strategy : An internal runtime growth algorithm that balances:
	// memory usage
	// allocation cost
	// performance
	// Small slices → grow faster (often ~2x)
	// Large slices -> slow growth there is no fixed formulae determine by go runtime internally
}

func DeepCopy() {
	a := []int{1, 2, 3}
	b := make([]int, len(a))
	copy(b, a)
	a[2] = 1
	b[2] = 4
	fmt.Printf("a val : %v and b val %v \n", a, b)
}

// Nil slice vs empty slice

// Type			Value			len	cap
// nil slice	nil	/ new		0	0
// empty slice	[]int{} / make	0	0

// Behavior:
// var s []int // nil
// Can still append
// Cannot index

// Slicing Operations
// max = control how much of the array your slice is allowed to “see” and hold in memory
func SlicingOperation() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	low := 0
	high := 6
	// max := 3 // array out of bound
	max := 6

	a := s[low:high]
	b := s[low:high:max]
	fmt.Printf("a val : %v and b val %v \n", a, b)
}

// Memory leak safeguard using max 
var store [][]int // keeps slices alive

func LeakyCase() {
	big := make([]int, 1000000)
	small := big[:10]
	store = append(store, small)
}

func SafeCase() {
	big := make([]int, 1000000)
	small := big[:10:10] 
	store = append(store, small)
}

func printMem(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s → Alloc = %v KB\n", msg, m.Alloc/1024)
}

func main() {
}
