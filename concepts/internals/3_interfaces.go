package internals


// An interface in Go defines a set of method signatures
// A type satisfies an interface implicitly((Duck Typing))

// Empty interface
	// var x interface{}
	// Holds any type
	// Equivalent to any

// (interface)
//    ├── type  → dynamic type
//    └── data  → pointer to actual value

// 5. Nil Interface vs Nil Value (Common Interview Trap)
	// var x *int = nil
	// var i interface{} = x
// i != nil 
	// Why?
	// i has:
	// type = *int
	// data = nil
// So:
	// interface is NOT nil if type ≠ nil


// Go uses a method table (itab):
// 	interface → itab → method pointers
// It enables dynamic dispatch

// Type Assertion
	// val, ok := i.(int)
// val → extracted value
// ok → whether assertion succeeded
// Without ok → panic if wrong type:
	// val := i.(int) // panics if not int

// 8. Type Switch
// switch v := i.(type) {
// case int:
// case string:
// default:
// }
//  Used to handle multiple possible types

// Interface Composition
// type Writer interface {
//     Write([]byte) (int, error)
// }

// type ReadWriter interface {
//     Reader
//     Writer
// }
// Interfaces can embed other interfaces

// . Interface as Function Parameter
// func Process(r Reader) {}
// ✔ Enables loose coupling
// ✔ Improves testability

