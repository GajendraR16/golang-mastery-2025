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

### Methods

Methods are functions with a special receiver argument. They allow you to define functions on types, enabling object-oriented programming patterns in Go.

**Syntax:** `func (receiver Type) methodName() returnType { }`

### Interfaces

Interfaces define a set of method signatures. A type implements an interface by implementing all the methods in the interface. Go uses implicit interface satisfaction - no explicit declaration needed.

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

