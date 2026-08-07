#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=../lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/../lib.sh"

# Apple's dot_clean has no -N; this case only exercises our --dry-run.
case_setup "06_dry_run"
seed_pair_stub "photo.jpg"

set +e
out="$("${DOTCLEAN}" -N --keep=native "${OURS_DIR}" 2>&1)"
oc=$?
set -e
if [[ "${oc}" -ne 0 ]]; then
	println "e2e: ours dry-run exited ${oc}" >&2
	println "${out}" >&2
	exit 1
fi

assert_exists "${OURS_DIR}" "._photo.jpg"
assert_exists "${OURS_DIR}" "photo.jpg"

case "${out}" in
*"${OURS_DIR}/._photo.jpg"*) ;;
*)
	println "e2e: dry-run output missing sidecar path" >&2
	println "${out}" >&2
	exit 1
	;;
esac
