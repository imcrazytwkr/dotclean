#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

case_setup "02_flat"
seed_pair_stub "root.jpg"
seed_pair_stub "nested.jpg" "sub"
run_both -f --keep=native
assert_missing "${APPLE_DIR}" "._root.jpg"
assert_missing "${OURS_DIR}" "._root.jpg"
assert_exists "${APPLE_DIR}" "sub/._nested.jpg"
assert_exists "${OURS_DIR}" "sub/._nested.jpg"
assert_sidecar_trees_match
