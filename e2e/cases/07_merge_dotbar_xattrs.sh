#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

case_setup "07_merge_dotbar_xattrs"

if ! is_merge_fs "${CASE_DIR}"; then
	println "e2e: skip 07_merge_dotbar_xattrs (TEST_DIR/TMPDIR is not APFS/HFS+)"
	exit 0
fi

if [[ ! -f "${GOLDEN_PDF}" ]]; then
	println "e2e: missing golden ${GOLDEN_PDF}" >&2
	exit 1
fi

seed_pair_golden "doc.pdf" "${GOLDEN_PDF}"
run_both --keep=dotbar
assert_missing "${APPLE_DIR}" "._doc.pdf"
assert_missing "${OURS_DIR}" "._doc.pdf"
assert_exists "${APPLE_DIR}" "doc.pdf"
assert_exists "${OURS_DIR}" "doc.pdf"
assert_sidecar_trees_match
assert_xattrs_match "doc.pdf"
