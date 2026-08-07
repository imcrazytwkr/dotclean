# E2E regression tests

Darwin-only regression tests: run Apple `dot_clean` (resolved via `command -v`) and `./dotclean` on identical trees and compare the results.

```bash
make test-e2e
# or
./e2e/run.sh
```

Skipped (exit 0) when not on Darwin or when a suitable `dot_clean` is not on `PATH`.

> :warning: **These tests must run on an APFS volume under real macOS**
>
> Set `TEST_DIR` env variable to an APFS path (e.g. under `$HOME`) so that merge/xattr case is not skipped.

**Compared:**

- sidecar presence
- exit status
- native xattr names/values after `--keep=dotbar` (case 07)

**Ours-only:**

- case 06 checks `-N` / `--dry-run` (Apple has no dry-run):
  - exits with 0
  - lists the sidecar
  - leaves the tree unchanged

**Not compared:** path whitespace handling (intentional difference); unimplemented Apple modes.
