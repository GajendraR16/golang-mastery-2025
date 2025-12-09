# Go Learning Notes

## Variables, Functions, Constants, and Type Overflows

Variables in Go are declared using `var` keyword or short declaration `:=`. Functions are declared with `func` keyword. Constants use `const` keyword and are immutable. Type overflows occur when a value exceeds the range of its data type.

## Looping Constructs

For loop is the only looping construct in Go.

**Syntax:** `for init; condition; increment/decrement { //code }`

- `init` runs just once at the beginning of the loop
- `condition` is evaluated every single time from the beginning. The loop won't begin if condition is false
- At the end of each iteration, `increment/decrement` is executed

## Conditions

**if-else statements:** Go supports standard if-else conditional statements.

Go also has a short statement which can be written in the if clause to compute anything and check the condition. The variable initialized in this short statement has a lifetime limited to just the if-else block.

Example: `if x := getValue(); x > 10 { // use x }`

## Types

### Pointers

- `p = &x` → returns the address of variable x
- `*p` → dereferences the pointer and returns the value itself

### Structs

Structs can be described as a grouping of other generic or aggregate types to form a new user-defined type.

Struct fields can be accessed using dot notation: `structName.fieldName`

### Arrays

The length and capacity of arrays are defined in the declaration. Using `len()` and `cap()` functions have no computational complexity overhead while accessing array properties.

### Slices

Slices are a dynamic way of declaring arrays. They provide a more flexible alternative to arrays with variable length and built-in methods for manipulation.

**Important Memory Detail:**
When we check the slice address versus the underlying array address, they are different:
- The **array address** points to the first element of the underlying array
- The **slice address** points to the slice struct, which contains:
  - Pointer to the first element of the underlying array
  - `len` field (current length)
  - `cap` field (capacity)

This means a slice is actually a struct that references an underlying array, not the array itself.

**Slices are Non-Comparable:**
Slices cannot be compared using the `==` operator (except for comparing to `nil`). This is because:
- Slices contain a pointer to the underlying array
- Comparing slices would require deep comparison of elements
- Go doesn't provide built-in deep comparison for slices

To compare slices, you must iterate through elements manually or use helper functions.

**Append Functionality:**
The `append()` function is used to add elements to a slice. Behind the scenes:
- If the slice has sufficient capacity, `append()` adds elements to the existing underlying array
- If capacity is exceeded, `append()` allocates a new larger array and **copies** all existing elements to it
- The new slice is returned with updated length and capacity
- This copy operation can be expensive for large slices

Example:
```go
s := []int{1, 2, 3}
s = append(s, 4, 5)  // May trigger copy if capacity exceeded
```

### Maps

Maps are Go's built-in associative data type (hash tables or dictionaries in other languages). They map keys to values and provide fast lookups.

**Syntax:** `map[KeyType]ValueType`

**Creating Maps:**
```go
// Using make
m := make(map[string]int)

// Using map literal
m := map[string]int{"key1": 1, "key2": 2}
```

**Checking if a Key Exists:**
Maps return two values when accessing a key: the value and a boolean indicating if the key exists.

```go
value, ok := m["key"]
if !ok {
    // Key doesn't exist
}

// Or compare in one line
if val, ok := m["key"]; !ok || val != expectedValue {
    // Key doesn't exist OR value doesn't match expected
}
```

**Important Map Properties:**
- Maps are **reference types** - passing a map to a function passes a reference, not a copy
- Maps are **not safe for concurrent use** - require synchronization (mutex) for concurrent access
- The **zero value** of a map is `nil` - a nil map behaves like an empty map for reads but causes panic on writes
- Maps are **not comparable** using `==` (except for comparing to `nil`)
- Iteration order is **not guaranteed** - maps iterate in random order

**Deleting from Maps:**
```go
delete(m, "key")  // Removes key from map, safe even if key doesn't exist
```

### Closures

Closures are anonymous functions that can access variables from their outer scope. They "close over" variables from the enclosing function, maintaining access to them even after the outer function returns.

**Basic Closure Example:**
```go
func counter() func() int {
    count := 0
    return func() int {
        count++  // Accesses variable from outer scope
        return count
    }
}

// Usage
c := counter()
fmt.Println(c()) // 1
fmt.Println(c()) // 2
```

**Key Characteristics:**
- **Variable Capture:** Closures capture variables by reference, not by value
- **Lifetime Extension:** Variables captured by closures remain alive even after the outer function returns
- **State Preservation:** Each closure maintains its own copy of captured variables

**Common Use Cases:**
- **Event Handlers:** Capturing context for callback functions
- **Factory Functions:** Creating specialized functions with pre-configured behavior
- **Iterators:** Maintaining state between function calls
- **Decorators:** Wrapping functions with additional behavior

**Variable Capture Gotcha:**
```go
// Common mistake - all closures capture the same variable
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() {
        fmt.Println(i) // All print 3 (final value of i)
    }
}

// Correct approach - capture by value
for i := 0; i < 3; i++ {
    i := i  // Create new variable in loop scope
    funcs[i] = func() {
        fmt.Println(i) // Each prints its own value
    }
}
```

**Memory Considerations:**
- Closures keep references to captured variables, preventing garbage collection
- Be mindful of memory leaks when closures capture large objects or long-lived references

### Methods

Methods are functions with a special receiver argument. They allow you to define functions on types, enabling object-oriented programming patterns in Go.

**Syntax:** `func (receiver Type) methodName() returnType { }`

**Value vs Pointer Receivers:**

**Value Receiver:**
```go
func (p Person) getName() string {
    return p.name  // Receives a copy of the struct
}
```
- Method receives a **copy** of the value
- Cannot modify the original struct
- Use when you don't need to modify the receiver
- More memory efficient for small structs

**Pointer Receiver:**
```go
func (p *Person) setName(name string) {
    p.name = name  // Modifies the original struct
}
```
- Method receives a **pointer** to the original value
- Can modify the original struct
- Use when you need to modify the receiver or for large structs (avoids copying)
- Required for methods that modify the receiver

**Method Values and Method Expressions:**

**Method Value:**
A method bound to a specific receiver instance.
```go
p := Person{name: "John"}
methodValue := p.getName  // Method bound to instance p
result := methodValue()   // Calls p.getName()
```

**Method Expression:**
A function that takes the receiver as its first argument.
```go
methodExpr := Person.getName     // Method expression
result := methodExpr(p)          // Pass receiver as first argument
// Equivalent to: result := p.getName()
```

Method expressions are useful for:
- Passing methods as function arguments
- Creating generic functions that work with different receiver types

### Interfaces

Interfaces define a set of method signatures. A type implements an interface by implementing all the methods in the interface. Go uses implicit interface satisfaction - no explicit declaration needed.

**Basic Interface Example:**
```go
type Writer interface {
    Write([]byte) (int, error)
}

type File struct { /* ... */ }

func (f *File) Write(data []byte) (int, error) {
    // File now implements Writer interface
}
```

**Smaller Interfaces are Better:**
- **Interface Segregation Principle:** Prefer many small, focused interfaces over large ones
- Small interfaces are easier to implement and compose
- Better abstraction - clients depend only on methods they need
- Example: `io.Reader` has just one method `Read()`, making it highly reusable

**Interface Internal Structure:**
An interface value has two components:
1. **Dynamic Type (Concrete Type):** The actual type of the value stored
2. **Dynamic Value:** The actual value of the concrete type

```go
var w io.Writer
w = os.Stdout  // Dynamic type: *os.File, Dynamic value: os.Stdout

// Interface is nil only when both type and value are nil
var w io.Writer  // w is nil (type=nil, value=nil)
```

**Pointer Receivers and Interface Satisfaction:**

**Why Pointer Receiver Methods Need Pointers:**
```go
type Counter struct { count int }

func (c *Counter) Increment() {  // Pointer receiver
    c.count++
}

var c Counter
var inc interface{ Increment() }

inc = &c  // ✓ Works - &c has type *Counter
inc = c   // ✗ Fails - c has type Counter, doesn't have Increment method
```
- Only `*Counter` satisfies the interface, not `Counter`
- The method is defined on the pointer type
- You cannot call a pointer receiver method on a non-addressable value

**Why Value Receiver Methods Work with Both:**
```go
type Counter struct { count int }

func (c Counter) Get() int {  // Value receiver
    return c.count
}

var c Counter
var getter interface{ Get() int }

getter = c   // ✓ Works - c has Get method
getter = &c  // ✓ Also works - Go automatically dereferences
```
- Both `Counter` and `*Counter` satisfy the interface
- Go can automatically dereference `&c` to call the value receiver method
- Value receiver methods are in the method set of both the type and its pointer

**Type Assertions:**
Type assertions extract the concrete value from an interface.

```go
var w io.Writer = os.Stdout

// Type assertion
f := w.(*os.File)  // Extracts *os.File from interface

// Safe type assertion with ok check
f, ok := w.(*os.File)
if !ok {
    // w doesn't hold *os.File
}
```

**Why Type Assertions are Necessary:**
```go
var w io.Writer = os.Stdout

// Cannot call File-specific methods directly
w.Sync()  // ✗ Error: io.Writer has no Sync method

// Must use type assertion first
if f, ok := w.(*os.File); ok {
    f.Sync()  // ✓ Works - f is *os.File
}
```
- Interface only exposes methods in its definition
- Concrete type may have additional methods
- Type assertion is required to access methods not in the interface
- Without type assertion, you're limited to the interface's method set

**Type Switches:**
```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case *os.File:
        fmt.Printf("File: %s\n", v.Name())
    default:
        fmt.Printf("Unknown type\n")
    }
}
```

**Empty Interface:**
```go
var any interface{}  // or any (Go 1.18+)
```
- Can hold values of any type
- Provides no methods
- Useful for generic containers or functions that accept any type
- Requires type assertion to use the underlying value

**Interface Comparison:**
- Interfaces are comparable if their dynamic types are comparable
- Two interface values are equal if they have identical dynamic types and equal dynamic values
- Comparing interfaces with non-comparable dynamic types (like slices) causes panic

**Best Practice: "Accept Interfaces, Return Concrete Types"**

**Why Return Concrete Types (Not Interfaces):**
```go
// ✗ Bad - Returns interface
func NewWriter() io.Writer {
    return &bytes.Buffer{}
}

// ✓ Good - Returns concrete type
func NewWriter() *bytes.Buffer {
    return &bytes.Buffer{}
}
```

**Reasons to return concrete types:**
1. **Flexibility for Callers:** Callers can access all methods of the concrete type, not just interface methods
2. **No Premature Abstraction:** Don't create abstractions until you need them
3. **Easier to Extend:** Adding methods to concrete types doesn't break existing code
4. **Better Documentation:** Concrete types show exactly what's being returned
5. **Avoid Interface Pollution:** Prevents creating unnecessary interfaces

**Example of the Problem:**
```go
// Returns interface - limits caller
func GetConfig() io.Reader {
    return &bytes.Buffer{}
}

config := GetConfig()
config.Reset()  // ✗ Error: io.Reader has no Reset method

// Returns concrete type - caller has full access
func GetConfig() *bytes.Buffer {
    return &bytes.Buffer{}
}

config := GetConfig()
config.Reset()  // ✓ Works: *bytes.Buffer has Reset method
```

**Why Accept Interfaces as Parameters:**
```go
// ✓ Good - Accepts interface
func WriteData(w io.Writer, data []byte) error {
    _, err := w.Write(data)
    return err
}

// Can be called with any type that implements io.Writer
WriteData(os.Stdout, data)
WriteData(&bytes.Buffer{}, data)
WriteData(httpResponseWriter, data)
```

**Reasons to accept interfaces:**
1. **Maximum Flexibility:** Function works with any type implementing the interface
2. **Testability:** Easy to pass mock implementations for testing
3. **Decoupling:** Function doesn't depend on concrete implementations
4. **Composability:** Promotes code reuse across different types
5. **Dependency Inversion:** Depend on abstractions, not concretions

**Summary:**
- **Parameters:** Use interfaces for flexibility and testability
- **Return Values:** Use concrete types to avoid limiting callers
- **Proverb:** "Be conservative in what you send, be liberal in what you accept"

## Bit Operations

Bit operations are fundamental operations that work directly on binary representations of numbers:

- `x & 1` → checks the last bit (determines if number is odd/even)
- `x >> 1` → right shifts by 1 bit, removes the last bit. Equivalent to `x / 2^1`
- `x << 1` → left shifts by 1 bit. Equivalent to `x * 2^1`
- `x & (x - 1)` → clears the rightmost 1-bit

**Example:**
```
x = 5 (binary: 101)
x - 1 = 4 (binary: 100)
x & (x - 1) = 101 & 100 = 100 (decimal: 4)
```

This operation is commonly used in algorithms to count set bits or check if a number is a power of 2.

