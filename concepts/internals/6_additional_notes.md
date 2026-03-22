## 9. Zero Values

### What are Zero Values?

Every variable in Go is automatically initialized to a default value if not explicitly assigned.

| Type                                 | Zero Value |
| ------------------------------------ | ---------- |
| int, float                           | 0          |
| string                               | ""         |
| bool                                 | false      |
| pointer, slice, map, func, interface | nil        |

###  Why Go Uses Zero Values

* Eliminates uninitialized memory bugs
* Simplifies code (no need for constructors everywhere)
* Enables predictable behavior

### Design Philosophy

* "Make the zero value useful"
* Structs should be usable without explicit initialization

Example:

```go
type Mutex struct {
    locked bool
}
```

Zero value is already usable (unlocked state)

### 🔹 Common Interview Scenarios

* Difference between nil slice vs empty slice
* Using zero value structs safely
* Checking nil vs empty

---

## 10. Type System

### 🔹 Static Typing

* Types are checked at compile time
* No implicit conversions

```go
var a int = 10
var b float64 = float64(a)
```

### 🔹 Dynamic Behavior via Interfaces

* Interfaces enable polymorphism

```go
type Shape interface {
    Area() float64
}
```

### 🔹 Duck Typing (Implicit Implementation)

* No explicit "implements"

### 🔹 Type Conversion vs Type Assertion

| Concept    | When Used                            |
| ---------- | ------------------------------------ |
| Conversion | Between known types                  |
| Assertion  | Extract concrete type from interface |

### 🔹 Edge Cases

* Losing precision in conversions
* Converting between custom types

---

## 11. Type Assertion

### 🔹 Syntax

```go
value, ok := i.(ConcreteType)
```

### 🔹 Safe Assertion

* Returns (value, ok)
* Prevents panic

### 🔹 Unsafe Assertion

```go
value := i.(ConcreteType) // panic if wrong
```

### 🔹 Type Switch

```go
switch v := i.(type) {
case int:
case string:
default:
}
```

### 🔹 Runtime Behavior

* Happens at runtime
* Uses type metadata stored in interface

### 🔹 Common Pitfalls

* Panics due to wrong type
* Confusing nil interface vs typed nil

---

## 🚀 12. Runtime Understanding

### 🔹 Go Runtime Responsibilities

* Goroutine scheduling
* Memory allocation
* Garbage collection

### 🔹 Stack vs Heap

| Stack        | Heap       |
| ------------ | ---------- |
| Fast         | Slower     |
| Auto managed | GC managed |

### 🔹 Escape Analysis

* Decides stack vs heap allocation

```bash
go build -gcflags="-m"
```

### 🔹 Garbage Collection (GC)

#### Basic Model

* Mark-and-sweep
* Concurrent GC

#### Steps

1. Mark reachable objects
2. Sweep unreachable memory

#### Key Features

* Low pause times
* Runs concurrently with program

### 🔹 Common Interview Angles

* What triggers GC?
* How to reduce GC pressure?

---

## 🚀 13. Performance Awareness

### 🔹 Avoid Unnecessary Allocations

* Reuse objects
* Use sync.Pool when appropriate

### 🔹 Value vs Reference Types

| Type            | Behavior   |
| --------------- | ---------- |
| Arrays, structs | Copied     |
| Slices, maps    | Referenced |

### 🔹 Copy vs Reference Implications

#### Copy Example

```go
b := a // full copy for structs
```

#### Reference Example

```go
b := a // slice points to same backing array
```

### 🔹 Slice Gotchas

* Underlying array sharing
* Append may reallocate

### 🔹 Common Performance Traps

* Large struct copies
* Excessive allocations in loops
* Interface boxing/unboxing

### 🔹 Optimization Mindset

* Measure before optimizing
* Use pprof

```bash
go tool pprof
```