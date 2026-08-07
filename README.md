[![license](https://img.shields.io/github/license/imcrazytwkr/dotclean)](LICENSE)

A drop-in cross-platform reimplementation of macOS [`dot_clean`](https://keith.github.io/xcode-man-pages/dot_clean.1.html): merge or remove AppleDouble `._*` sidecars.

## Features

- Usage is backwards-compatible to the macOS implementation.
- **Merges xattrs on macOS HFS+/APFS only**; elsewhere (including Linux), discards and removes sidecars on FAT32/exFAT and other filesystems. Merge behavior is not documented in the original but has been derived empirically.
- `-N` / `--dry-run` lists deletion targets without changing the filesystem.
- `-Q` / `--set-quarantine` applies `com.apple.quarantine` from AppleDouble when merging (off by default, matching Apple).
- `--keep=mostrecent|dotbar|native` controls which side wins when merging xattrs (default `mostrecent`).
- `-D` / `--deep` and `-S` / `--spotlight` optionally remove Finder / Spotlight junk (disabled by `-p` / `--preserve`).
- Opaque whitespace-safe path arguments (I've had problems with how `dot_clean` handles whitespace in the past)
- GNU-style flags via [pflag](https://github.com/spf13/pflag)

## Build

```bash
make build
```

## Testing

```bash
make test          # go test ./...
make test-e2e      # Darwin only: compare against Apple's dot_clean (skips elsewhere)
```

Darwin unit tests (`//go:build darwin`) use `TEST_DIR`, defaulting to system temp directory. Set `TEST_DIR` to an APFS or HFS+ path so that xattr merge tests are not skipped. The same applies to end-to-end tests; see [`e2e/README.md`](e2e/README.md).

## Usage

```bash
dotclean /path/to/volume
dotclean -N '/path/with spaces' "$HOME/another path with spaces"
# macOS mounts
dotclean -m /Volumes/USB
# Linux GVFS mounts
dotclean -m "/run/media/$USER/USB"
dotclean -h
```

Full manual: [docs/dotclean.1.md](docs/dotclean.1.md) (generate man page with `make man` if [pandoc](https://pandoc.org/) is installed).

## Dependencies

- Go toolchain
- `github.com/spf13/pflag`
- `golang.org/x/sys` (Darwin xattr / `statfs`)
