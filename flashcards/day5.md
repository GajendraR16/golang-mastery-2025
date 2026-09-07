# Day 5 Flashcards — Concurrency

## Morning Session — Goroutines & Channels

| Exercise | Done |
|---|---|
| WaitGroup | ✓ |
| Pipeline | ✓ |
| Worker Pool | ✓ |
| Select with Timeout | ✓ |

- **Biggest struggle**:
  - Forgot to `range` over the channel in the worker pool
  - Used `time.Sleep` instead of `WaitGroup` for synchronization

---

## Evening Session — Context & Graceful Shutdown

| Checkpoint | Result |
|---|---|
| Context added to DB methods | ✓ |
| Graceful shutdown working | ✓ |
| Ctrl+C waits for in-flight requests | ✓ |

---

## Key Insights

- **Channels**: Didn't focus enough while writing — caused subtle bugs. Need to be more deliberate with channel direction and ranging over them
- **Context**: Mostly straightforward — just threading `ctx context.Context` through the call chain. Mistake: used bare `ctx` in the handler instead of `r.Context()`
- **Graceful shutdown**: Referenced a Medium tutorial to implement — worth revisiting to internalize the pattern
