% dotclean(1) | User Commands
%
% August 2026

# NAME

`dotclean` - merge or remove AppleDouble **`._*`** sidecar files

# SYNOPSIS

**dotclean** [*options*] *directory* ...

# DESCRIPTION

**`dotclean`** is a cross-platform reimplementation of macOS **`dot_clean`**(1). For each *directory*, it finds AppleDouble sidecar files (names beginning with **`._`**) and either merges their metadata into the corresponding native file or discards that metadata, optionally deleting the sidecar file.

**Merge** (applying extended attributes from the AppleDouble file) is performed **only** under macOS and only on **HFS+** and **APFS** volumes. On Linux, HFS+ and APFS support is not consistent enough among various filesystem drivers, so **ext\*** and **Btrfs** extended attributes are not supported. On all other filesystems, including **FAT32** and **exFAT**, metadata cannot be stored natively without recreating a `._*` file, so **`dotclean`** removes the sidecar file unless the **`--preserve`** flag is supplied. Missing extended-attribute support is not an error: a successful clean-up exits 0.

## CAVEATS

**`dotclean`** never applies **`com.apple.macl`** or **`com.apple.provenance`** attributes as they are considered protected by macOS. The list of protected attributes within the code will get updated as more are discovered. As a precaution, permission errors arising when setting an attribute are skipped.

**`com.apple.quarantine`** attribute is applied only when the **`-Q`** flag is supplied. By default it is ignored to conform to **`dot_clean`**(1).

# OPTIONS

**`-f`**, **`--flat`**
: Do not recurse into subdirectories.

**`-h`**, **`--help`**
: Print help and exit.

**`-m`**, **`--always-delete`**
: Always delete AppleDouble files (including orphans). Overrides **`--preserve`**.

**`-n`**, **`--cleanup`**
: Delete an AppleDouble file if there is no matching native file.

**`-p`**, **`--preserve`**
: Preserve AppleDouble files after handling (do not delete paired sidecars). Disables **`--deep`** and **`--spotlight`**.

**`-s`**, **`--follow-symlinks`**
: Follow symbolic links when they point at AppleDouble files.

**`-v`**, **`--verbose`**
: Print verbose progress to standard error.

**`-N`**, **`--dry-run`**
: List paths that would be deleted; do not merge or delete anything.

**`-Q`**, **`--set-quarantine`**
: When merging on HFS+/APFS, apply **com.apple.quarantine** from the AppleDouble sidecar. Off by default (matches **`dot_clean`**(1), which does not re-apply quarantine).

**`-D`**, **`--deep`**
: Remove **`.DS_Store`**, **`.AppleDouble`**, **`.Trashes`**, and **`.TemporaryItems`** (files or directories). Off by default. Ignored when **`-p`** / **`--preserve`** is set.

**`-S`**, **`--spotlight`**
: Remove **`.fseventsd`** and directories whose names start with **`.Spotlight`**. Off by default. Ignored when **`-p`** / **`--preserve`** is set.

**`--keep`**=*`mode`*
: When merging on HFS+/APFS: **mostrecent** (default, prefer existing native attributes), **dotbar** (prefer AppleDouble), or **native** (skip applying AppleDouble attributes).

# DELETE RULES

| **Native file** | Default | **-m** | **-p** | **-n** |
|-----------------|---------|--------|--------|--------|
| Exists          | Delete  | Delete | Keep   | N/A    |
| Does not exist  | Keep    | Delete | Keep   | Delete |

# EXAMPLES

Remove AppleDouble sidecars on a USB volume under macOS (typical FAT/exFAT clean-up):

```
dotclean /Volumes/USB
```

Remove AppleDouble sidecars on a USB volume under Linux (GVFS):

```
dotclean /run/media/user/USB
```

Preview deletions:

```
dotclean -N '/run/media/user/My Drive/Pictures'
```

Merge preferring AppleDouble data on an APFS path, then delete sidecars:

```
dotclean --keep=dotbar ~/Documents
```

# EXIT STATUS

**0**
: Success, including a successful dry-run or help.

**1**
: No directories were given, or an error occurred while cleaning.

**2**
: Invalid options (bad flags or an unrecognized **`--keep`** value).

# SEE ALSO

**dot_clean**(1) on macOS, AppleSingle/AppleDouble Formats (Apple Developer Note).

# NOTES

This implementation is not affiliated with Apple. Behavior aims to match **`dot_clean`**(1) as much as possible while also providing Linux and BSD support and additional options for deeper cleaning of macOS resource files.
