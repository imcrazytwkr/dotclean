#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

case_setup "05_preserve"
seed_pair_stub "photo.jpg"
run_both -p --keep=native
assert_exists "${APPLE_DIR}" "._photo.jpg"
assert_exists "${OURS_DIR}" "._photo.jpg"
assert_sidecar_trees_match
