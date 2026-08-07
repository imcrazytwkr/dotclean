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
- native xattr names/values after `--keep=dotbar` (case 07), excluding attrs that may still differ under SIP (`com.apple.macl`, `com.apple.provenance` — product merge also skips these; regex is a safety net). Default merge omits quarantine on both tools (Apple never applies it; ours requires `-Q`).

**Ours-only:**

- case 06 checks `-N` / `--dry-run` (Apple has no dry-run):
  - exits with 0
  - lists the sidecar
  - leaves the tree unchanged
- `-Q` / `--set-quarantine` (opt-in quarantine from AppleDouble; not used in parity cases)

**Not compared:** path whitespace handling (intentional difference); unimplemented Apple modes.
