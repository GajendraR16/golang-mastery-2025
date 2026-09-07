# Day 4 Flashcards — PostgreSQL

## Fixes from Day 3

- **Search handler**: Removed fake error return
- **Server logging**: Migrated to `slog`

---

## Morning Session — PostgresStore

| Checkpoint | Result |
|---|---|
| Remembered `sql.NullTime` | ✓ |
| Remembered `$1` placeholders | ✓ |
| Remembered `rows.Close()` | ✓ |

- **Time to implement**: 3–4 hours

---

## Evening Session — Swap & Test

- **One-line swap**: ✗ — didn't work
- **curl commands**: All passing ✓
- **Bug hit**: Task wasn't syncing after update — a duplicate was being created because `Save` was still being called inside `Complete`

---

## LeetCode

No problems assigned today.
