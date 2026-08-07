#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=common.sh
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
E2E_ROOT="${REPO_ROOT}/e2e"
cd "${REPO_ROOT}"

resolve_apple_dot_clean() {
	local path
	path="$(command -v dot_clean 2>/dev/null)" || return 1

	case "${path}" in
	/usr/bin/* | /usr/sbin/* | /usr/local/bin/* | /usr/local/sbin/*) ;;
	*)
		println "e2e: skip (dot_clean at ${path} is not under /usr/bin or /usr/local/bin)" >&2
		return 1
		;;
	esac

	if [[ ! -x "${path}" ]]; then
		println "e2e: skip (requires dot_clean on PATH under /usr/bin or /usr/local/bin)" >&2
		return 1
	fi

	println "${path}"
}

if [[ "$(uname -s)" != "Darwin" ]]; then
	println "e2e: skip (requires Darwin)"
	exit 0
fi

APPLE_DOT_CLEAN="$(resolve_apple_dot_clean)" || exit 0

BASE="${TEST_DIR:-${TMPDIR:-/tmp}}"
E2E_WORK="$(mktemp -d "${BASE}/dotclean-e2e.XXXXXX")"
trap 'rm -rf "${E2E_WORK}"' EXIT

BIN_DIR="${E2E_WORK}/bin"
mkdir -p "${BIN_DIR}"

DOTCLEAN="${BIN_DIR}/dotclean"
go build -o "${BIN_DIR}/dotclean" .

export REPO_ROOT E2E_ROOT E2E_WORK APPLE_DOT_CLEAN DOTCLEAN

failed=0
shopt -s nullglob
for case in "${E2E_ROOT}/cases/"*.sh; do
	if ! bash "${case}"; then
		println "FAIL: ${case}" >&2
		failed=1
	fi
done

if [[ "${failed}" -ne 0 ]]; then
	println "e2e: one or more cases failed" >&2
	exit 1
fi
println "e2e: all cases passed (apple=${APPLE_DOT_CLEAN})"
