#!/usr/bin/env bash

# Shared helpers for e2e cases. Sourced by each case script.
set -euo pipefail

# shellcheck source=common.sh
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

: "${E2E_ROOT:?}"
: "${E2E_WORK:?}"
: "${DOTCLEAN:?}"
: "${APPLE_DOT_CLEAN:?}"
: "${REPO_ROOT:?}"

GOLDEN_PDF="${REPO_ROOT}/internal/appledouble/testdata/rich.pdf.appledouble"
export GOLDEN_PDF

# Xattrs that may differ or fail under SIP / policy; ignored in comparisons.
XATTR_SKIP_REGEX='^(com\.apple\.macl|com\.apple\.provenance)$'


case_setup() {
	local name="$1"
	CASE_DIR="${E2E_WORK}/${name}"
	rm -rf "${CASE_DIR}"
	APPLE_DIR="${CASE_DIR}/apple"
	OURS_DIR="${CASE_DIR}/ours"
	mkdir -p "${APPLE_DIR}" "${OURS_DIR}"
	println "==> ${name}"
}

seed_pair_stub() {
	# native + stub sidecar (use with --keep=native so neither tool merges)
	local base="$1"
	local sub="${2:-}"
	_seed_pair "${base}" "${sub}" stub
}

seed_pair_golden() {
	local base="$1"
	local golden="$2"
	local sub="${3:-}"
	_seed_pair "${base}" "${sub}" "${golden}"
}

_seed_pair() {
	local base="$1"
	local sub="$2"
	local payload="$3"

	local apple_d ours_d
	if [[ -n "${sub}" ]]; then
		apple_d="${APPLE_DIR}/${sub}"
		ours_d="${OURS_DIR}/${sub}"
		mkdir -p "${apple_d}" "${ours_d}"
	else
		apple_d="${APPLE_DIR}"
		ours_d="${OURS_DIR}"
	fi

	printf 'payload' >"${apple_d}/${base}"
	printf 'payload' >"${ours_d}/${base}"

	local side="._${base}"
	if [[ "${payload}" == "stub" ]]; then
		printf 'AD' >"${apple_d}/${side}"
		printf 'AD' >"${ours_d}/${side}"
	else
		cp "${payload}" "${apple_d}/${side}"
		cp "${payload}" "${ours_d}/${side}"
	fi
}

seed_orphan_stub() {
	local base="$1"
	printf 'AD' >"${APPLE_DIR}/._${base}"
	printf 'AD' >"${OURS_DIR}/._${base}"
}

run_both() {
	# Usage: run_both [flags...]
	# Runs Apple then ours against their trees. Extra args are tool flags.
	local -a flags=("$@")
	local ac oc
	set +e
	"${APPLE_DOT_CLEAN}" "${flags[@]}" "${APPLE_DIR}"
	ac=$?
	"${DOTCLEAN}" "${flags[@]}" "${OURS_DIR}"
	oc=$?
	set -e
	if [[ "${ac}" -ne 0 || "${oc}" -ne 0 ]]; then
		println "e2e: non-zero exit apple=${ac} ours=${oc}" >&2
		return 1
	fi
}

list_sidecars() {
	# Relative paths of ._ files under $1, sorted.
	local root="$1"
	# -print0 / sort -z would be safer; names here have no newlines.
	(cd "${root}" && find . -name '._*' -print | sed 's|^\./||' | sort)
}

assert_sidecar_trees_match() {
	local a b
	a="$(list_sidecars "${APPLE_DIR}")"
	b="$(list_sidecars "${OURS_DIR}")"
	if [[ "${a}" != "${b}" ]]; then
		println "e2e: sidecar tree mismatch" >&2
		println "--- apple ---" >&2
		println "${a}" >&2
		println "--- ours ---" >&2
		println "${b}" >&2
		return 1
	fi
}

assert_exists() {
	local root="$1"
	local rel="$2"
	if [[ ! -e "${root}/${rel}" ]]; then
		println "e2e: missing ${root}/${rel}" >&2
		return 1
	fi
}

assert_missing() {
	local root="$1"
	local rel="$2"
	if [[ -e "${root}/${rel}" ]]; then
		println "e2e: still present ${root}/${rel}" >&2
		return 1
	fi
}

is_merge_fs() {
	# True if path's volume looks like APFS or HFS(+).
	local path="$1"
	local mnt personality
	mnt="$(df "${path}" | awk 'NR==2 {print $NF}')"
	personality="$(diskutil info "${mnt}" 2>/dev/null | awk -F': *' '/File System Personality/ {print $2}')"
	case "${personality}" in
	*APFS* | *HFS*) return 0 ;;
	*) return 1 ;;
	esac
}

xattr_names() {
	# macOS `xattr FILE` lists names one per line.
	local file="$1"
	xattr "${file}" 2>/dev/null | grep -vE "${XATTR_SKIP_REGEX}" | sort -u || true
}

xattr_dump() {
	local file="$1"
	local name
	while IFS= read -r name; do
		[[ -z "${name}" ]] && continue
		println "${name}"
		xattr -px "${name}" "${file}" 2>/dev/null || println '<missing>'
	done < <(xattr_names "${file}")
}

assert_xattrs_match() {
	local rel="$1"
	local a b
	a="$(xattr_dump "${APPLE_DIR}/${rel}")"
	b="$(xattr_dump "${OURS_DIR}/${rel}")"
	if [[ "${a}" != "${b}" ]]; then
		println "e2e: xattr mismatch for ${rel}" >&2
		println "--- apple ---" >&2
		println "${a}" >&2
		println "--- ours ---" >&2
		println "${b}" >&2
		return 1
	fi
	if [[ -z "${a}" ]]; then
		println "e2e: expected some xattrs on ${rel} after merge" >&2
		return 1
	fi
}
