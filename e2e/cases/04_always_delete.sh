#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

case_setup "04_always_delete"
seed_orphan_stub "lonely.jpg"
seed_pair_stub "photo.jpg"
run_both -m --keep=native
assert_missing "${APPLE_DIR}" "._lonely.jpg"
assert_missing "${OURS_DIR}" "._lonely.jpg"
assert_missing "${APPLE_DIR}" "._photo.jpg"
assert_missing "${OURS_DIR}" "._photo.jpg"
assert_exists "${APPLE_DIR}" "photo.jpg"
assert_exists "${OURS_DIR}" "photo.jpg"
assert_sidecar_trees_match
