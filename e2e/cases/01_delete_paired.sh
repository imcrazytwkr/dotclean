#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

case_setup "01_delete_paired"
seed_pair_stub "photo.jpg"
run_both --keep=native
assert_missing "${APPLE_DIR}" "._photo.jpg"
assert_missing "${OURS_DIR}" "._photo.jpg"
assert_exists "${APPLE_DIR}" "photo.jpg"
assert_exists "${OURS_DIR}" "photo.jpg"
assert_sidecar_trees_match
