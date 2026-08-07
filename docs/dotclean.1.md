% dotclean(1) | User Commands
%
% August 2026

# NAME

dotclean - merge or remove AppleDouble `._*` sidecar files

# SYNOPSIS

**dotclean** [*options*] *directory* ...

# DESCRIPTION

**dotclean** is a cross-platform reimplementation of macOS **dot_clean**(1). For each *directory*, it finds AppleDouble sidecar files (names beginning with `._`) and either merges their metadata into the corresponding native file or discards that metadata, then optionally deletes the sidecar.

**Merge** (apply extended attributes from the AppleDouble file) is performed **only** on **HFS+** and **APFS** volumes. On all other filesystems—including **FAT32** and **exFAT**—metadata cannot be stored natively without recreating a `._*` file, so **dotclean** discards the sidecar contents and deletes the sidecar (unless **--preserve** is set). Missing extended-attribute support is not an error; a successful cleanup exits 0.

Unlike Apple’s tool, path operands are treated as opaque: whitespace inside a single argument is preserved.

`.DS_Store` and other Finder litter are **not** removed in the default (normal) mode.

# OPTIONS

**-f**, **--flat**
: Do not recurse into subdirectories.

**-h**, **--help**
: Print help and exit.

**-m**, **--always-delete**
: Always delete AppleDouble files (including orphans). Overrides **--preserve**.

**-n**, **--cleanup**
: Delete an AppleDouble file if there is no matching native file.

**-p**, **--preserve**
: Preserve AppleDouble files after handling (do not delete paired sidecars).

**-s**, **--follow-symlinks**
: Follow symbolic links when they point at AppleDouble files.

**-v**, **--verbose**
: Print verbose progress to standard error.

**-N**, **--dry-run**
: List paths that would be deleted; do not merge or delete anything.

**-Q**, **--set-quarantine**
: When merging on HFS+/APFS, apply **com.apple.quarantine** from the AppleDouble sidecar. Off by default (matches Apple **dot_clean**, which does not re-apply quarantine).

**--keep**=*mode*
: When merging on HFS+/APFS: **mostrecent** (default; prefer existing native attributes), **dotbar** (prefer AppleDouble), or **native** (skip applying AppleDouble attributes).

# DELETE RULES

| Situation | Default | **-m** | **-p** | **-n** |
|-----------|---------|--------|--------|--------|
| Paired sidecar (native exists), after handle | Delete | Delete | Keep | — |
| Orphan sidecar (no native) | Keep | Delete | Keep | Delete |

# EXAMPLES

Remove AppleDouble sidecars under a USB volume (typical FAT/exFAT cleanup):

```
dotclean /Volumes/USB
```

Preview deletions:

```
dotclean -N "/media/user/My Drive/Pictures"
```

Merge preferring AppleDouble data on an APFS path, then delete sidecars:

```
dotclean --keep=dotbar ~/Documents
```

# DIAGNOSTICS

**dotclean** exits 0 on success (including a successful dry-run), and \>0 if an error occurs.

# SEE ALSO

**dot_clean**(1) on macOS, AppleSingle/AppleDouble Formats (Apple Developer Note).

# NOTES

This implementation is not affiliated with Apple. Behavior aims to match useful **dot_clean** semantics for removable media while fixing path whitespace handling and documenting **--preserve** (**-p**), which appears in Apple’s **-h** output but not always in the published man page.
