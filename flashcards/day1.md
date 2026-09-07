# Day 1 Flashcards

## Morning Session

### 1. Linked List — Insert, Delete, Reverse

- **Delete**: Used `curr.Next` as loop condition instead of `curr`
- **Delete**: Pointed `prev.Next = curr.Next` correctly but struggled with the logic initially
- **Reverse**: Needed to track `prev`, `curr`, and `next` pointers carefully

---

### 2. Stack — Push, Pop, Peek

- **Pop issue**: Removed the value before capturing it
  ```
  [1, 2, 3, 4] → pop() removed 4, then returned 3 instead of 4
  ```
- Fix: Capture the value first, then remove it

---

### 3. Two Sum (Map)

- Solved correctly but flipped the return indices
- Remember: return `[i, map[complement]]` not the other way around

---

### 4. FizzBuzz

- Missing `else` clause — wasn't returning the number when neither Fizz nor Buzz

---

### 5. Reverse a String

- Initial approach used two maps — simplified to one pass with two pointers
- Two-pointer pattern: swap `s[left]` and `s[right]`, move inward

---

## Evening Session

### 1. Two Sum *(repeat)*

- Reinforces: map lookups, one-pass thinking
- Target: ≤ 10 min
- No new mistakes — repetition helping

---

### 2. Valid Parentheses

- Reinforces: stack + conditionals
- Target: ≤ 15 min
- Missed the final stack length check — need to verify stack is empty at the end

---

### 3. Reverse Linked List *(repeat)*

- Reinforces: pointers/references, linked-list manipulation
- Target: ≤ 15 min
- Already done in morning — felt easy the second time

---

### 4. Valid Anagram

- Reinforces: strings, hash maps, frequency counting
- Target: ≤ 15 min
- No mistakes ✓
