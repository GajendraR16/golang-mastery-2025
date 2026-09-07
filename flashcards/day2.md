# Day 2 Flashcards

## Bugs Fixed from Day 1

- **Stack Pop**: Captured the value before deleting it
- **Valid Parentheses**: Added base case check — `len(stack) == 0`

---

## Morning Session

### Task Manager CLI

- **What was easy**: Almost everything
- **What you struggled with**: `Delete` logic — used wrong `copy` syntax
  - Used: `s[:idx], s[:idx+1]`
  - Should be: `s[idx:], s[idx+1:]`
  - Also forgot to trim the last duplicate element: `s = s[:len(s)-1]`
- **Key insight**: `copy(dst, src)` shifts elements left — slice the tail correctly then shrink the slice

---

## Evening Session

### Interfaces

- **What was easy**: Writing the logic for all methods
- **What you struggled with**:
  - Figuring out the right data type to use — slept on it and came up with an approach
  - Return types for `Save` and `List` — they don't necessarily produce a storage error
  - Key insight: don't create errors just to satisfy an interface — only return errors that are meaningful

---

### LeetCode

| Problem | Time | Mistakes |
|---------|------|----------|
| #217 Contains Duplicate | 2 min | None |
| #88 Merge Sorted Array | 7 min | Missing condition in `else` block — caught after tracing |
| #26 Remove Duplicates | 1 min 20 sec | None |
| #27 Remove Element | — | Couldn't solve |
| #349 Intersection of Two Arrays | 2 min | None |
