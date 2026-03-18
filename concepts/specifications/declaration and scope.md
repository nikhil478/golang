## Declarations and scope

source : https://go.dev/ref/spec#Declarations_and_scope

- A declaration binds a non-blank identifier to a constant, type, type parameter, variable, function, label or package
- Every identifier in a program must be declared.
- No identifier may be declared twice in same block and no identifier may be declared both the file and package block

The blank identifier may be used like any other identifier in a declaration, but it does not introduce a binding and thus is not declared.

In the package block, the identifier init may only be used for init function declarations, and like the blank identifiers it does not introduce a new binding

Declaration = ConstDecl | TypeDecl | VarDecl
TopLevelDecl = Declaration | FunctionDecl | MethodDecl

Go is lexically scoped using blocks:

    1. The scope of a predeclared identifier (https://go.dev/ref/spec#Predeclared_identifiers) is the universe block
    2. The scope of an identifier denoting a constant, type, variable, or function (but not method) declared at top level (outside any function) is the package block
    3. The scope of the package name of an imported package is the file block of the file containing the import declaration
    4. The scope of an identifier denoting a method receiver (The variable before a method — represents the instance), function parameter, or result variable is the function body
    5. The scope of an identifier denoting a type parameter of a function or declared by a method receiver begins after the name of the function and ends at the end of the function body

    valid : 
    func printValue[T any](v T) T {
    var x T
    return x
}
    wrong : 
    var x T // ❌ ERROR: T not in scope yet

func printValue[T any](v T) {}


    6. The scope of a constant or a variable identifier denoting a type parameter of a type begins after the name of the type and ends at the end of the TypeSpec

    valid: type Box[T any] struct {
    value T
}

wrong:

type Box[T any] struct {
    value T
}

7. The scope of a constant or variable identifier declared inside a function begins at the end of the ConstSpec or VarSpec (ShortVarDecl for short variable declarations) and ends at the end of the innermost containing block.

[Variables and constants in Go exist only after they are declared and only inside the nearest block they belong to.]

8. The scope of a type identifier declared inside a function begins at the identifier in the TypeSpec and ends at the end of the innermost containing block.

Variable:
“I need to exist first before I can be used”
Type:
“I don’t need to ‘exist’ — I’m just a label”

They are discussed separately because their starting point of scope is different, even though their ending point is the same.


An identifier declared in a block may be redeclared in an inner block. While the identifier of the inner declaration is in scope, it denotes the entity declared by the inner declaration.

Even trickier (real-world bug)
err := doSomething()

if err != nil {
    err := handleError() // ❗ shadows outer err
    _ = err
}
👉 Outer err is untouched → can cause bugs


The package clause is not a declaration; the package name does not appear in any scope. Its purpose is to identify the files belonging to the same package and to specify the default package name for import declarations.

Labels are declared by labeled statements and are used in the "break", "continue", and "goto" statements. It is illegal to define a label that is never used. In contrast to other identifiers, labels are not block scoped and do not conflict with identifiers that are not labels. The scope of a label is the body of the function in which it is declared and excludes the body of any nested function.


| Feature          | Variables | Labels                    |
| ---------------- | --------- | ------------------------- |
| Block scoped?    | ✅ Yes     | ❌ No                      |
| Function scoped? | ❌ No      | ✅ Yes                     |
| Must be used?    | ❌ No      | ✅ Yes                     |
| Name conflict?   | ❌ Yes     | ❌ No (separate namespace) |

Ahh — you caught a **mistake in wording earlier**, and you’re right to question it. Let’s correct it cleanly.

---

# ❌ What I said earlier (incorrect wording)

> “package main = label”

👉 That is **wrong / misleading** ❗
Thanks for pointing it out.

---

# ✅ Correct understanding

## 1️⃣ `package main`

```go id="9od5ti"
package main
```

✔ NOT a label
✔ NOT a declaration
✔ NOT an identifier in scope

👉 It is just:

> **a package clause (metadata)**

---

## 2️⃣ Label (real label)

```go id="x92m0r"
start:
    fmt.Println("hi")
```

✔ This **is a label**
✔ Used with `goto`, `break`, `continue`
✔ Has **function scope**

---

## 3️⃣
