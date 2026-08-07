#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

case_setup "03_orphan_cleanup"
seed_orphan_stub "lonely.jpg"
run_both -n --keep=native
assert_missing "${APPLE_DIR}" "._lonely.jpg"
assert_missing "${OURS_DIR}" "._lonely.jpg"
assert_sidecar_trees_match
