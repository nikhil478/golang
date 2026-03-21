package internals

import "fmt"

// A Map is hash table with buckets, where each bucket holds multiple key-value pairs

// m[key] = value

// nternally:
// Hash function computes hash of key
// Hash → determines bucket index
// Bucket stores key-value pairs
// If full → overflow bucket is used

// Internally, Go uses a structure similar to:

// type hmap struct {
//     count     int           // number of elements
//     B         uint8         // log2(number of buckets)
//     buckets   unsafe.Pointer
//     oldbuckets unsafe.Pointer // during resizing
//     hash0     uint32
//     nevacuate   int // how much has been moved used during resizing
// }

// Key concepts:
// 	buckets → main storage
// 	oldbuckets → used during resizing
// 	B → controls number of buckets (2^B)

// Each bucket stores:
// 	Up to 8 key-value pairs
// 	Plus overflow pointer if needed

// If two keys go to same bucket:
// 	Stored in same bucket
// 	If full → overflow bucket is created

// the unsafe.Pointer for buckets in hmap points to a contiguous array of bucket structures, not a linked list.
// if bucket is full key hash collides then overflow bucket is created not in regular file regular will be indexed by key thats it

// Load factor controls when resizing happens:
// entries / buckets
// High → triggers resize
// Low → efficient

// Memory Behavior
// Maps:
// Do NOT shrink automatically (Ends up in memory leak instead create new maps)
// Retain memory even after deletions
// because they reuse it for future allocations instead or releading it to os

//. using new or var m map[string]int
// if u try to use it this will panic as unsafe.pointer value is nil

// Key rules in map: 
	// 	Allowed: int, string, bool, arrays

// Value rules in Map: Value Can Be Anything (map, slice, interface, function...)

// Additional notes , struct and interface allowed as keys , but value inside struct should be comparable
// or in case of interface allow those contrect types allowed where value is comparable

func MapCreation() {

	// using make
	m := make(map[string]int)

	// using literals
	n := map[string]int{
		"a": 1,
	}

	fmt.Printf("m %v & n : %v", m, n)
}

func MapOperations() {
	m := make(map[string]int)
	m["a"] = 1

	// Read operation
	val := m["a"]
	fmt.Printf("Read operation %v \n", val)

	// Safe check operation
	val, ok := m["key"]
	fmt.Printf("Safe check operation , val : %v, ok : %v \n", val, ok)

	// delete : safe operation not panic even ele not exist
	fmt.Println("Before delete where key not exist")
	delete(m, "key")
	fmt.Printf("After delete where key not exist %v \n", m)

	// delete : safe operation not panic even ele not exist
	fmt.Println("Before delete where key exist")
	delete(m, "a")
	fmt.Printf("After delete where key exist %v \n", m)

	// Thats how u can iterate over map
	// for k, v := range m {

	// }
}

// Map is NOT Thread-Safe
// throws error fatal error: concurrent map writes

// Safe alternatives
	// sync.Mutex
	// sync.RWMutex
	// sync.Map

func MapNotThreadSafe() {
	m := make(map[string]int)
	for range 50 {
		go func() { m["a"] = 1 }()
		go func() { m["b"] = 2 }()
	}
}