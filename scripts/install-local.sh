#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Install the current parley-deck-cli checkout locally.

Usage:
  scripts/install-local.sh [--prefix DIR] [--bin-dir DIR] [--dry-run]

Defaults:
  prefix:  $PARLEY_DECK_HOME or $HOME/.parley-deck
  bin-dir: <prefix>/bin

The installer builds ./cmd/parley from the current checkout and installs:
  <bin-dir>/parley

It also keeps Go build/module caches under <prefix>/cache by default so sandboxed
or restricted home-cache setups do not block local installs.
USAGE
}

repo_root() {
  local script_dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  cd "${script_dir}/.." && pwd
}

abs_path() {
  case "$1" in
    /*) printf '%s\n' "$1" ;;
    *) printf '%s/%s\n' "$(pwd)" "$1" ;;
  esac
}

prefix="${PARLEY_DECK_HOME:-${HOME}/.parley-deck}"
bin_dir=""
dry_run=0

while [ "$#" -gt 0 ]; do
  case "$1" in
    --prefix)
      if [ "$#" -lt 2 ]; then
        echo "missing value for --prefix" >&2
        exit 2
      fi
      prefix="$2"
      shift 2
      ;;
    --bin-dir)
      if [ "$#" -lt 2 ]; then
        echo "missing value for --bin-dir" >&2
        exit 2
      fi
      bin_dir="$2"
      shift 2
      ;;
    --dry-run)
      dry_run=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

root="$(repo_root)"
prefix="$(abs_path "${prefix}")"
if [ -z "${bin_dir}" ]; then
  bin_dir="${prefix}/bin"
fi
bin_dir="$(abs_path "${bin_dir}")"
target="${bin_dir}/parley"
build_cache="${GOCACHE:-${prefix}/cache/go-build}"
mod_cache="${GOMODCACHE:-${prefix}/cache/go-mod}"
metadata="${prefix}/INSTALLATION"

if ! command -v go >/dev/null 2>&1; then
  echo "go is required but was not found on PATH" >&2
  exit 1
fi

git_commit="unknown"
if command -v git >/dev/null 2>&1; then
  git_commit="$(cd "${root}" && git rev-parse --short HEAD 2>/dev/null || true)"
  if [ -z "${git_commit}" ]; then
    git_commit="unknown"
  fi
fi

echo "parley-deck-cli local install"
echo "  repo:       ${root}"
echo "  commit:     ${git_commit}"
echo "  prefix:     ${prefix}"
echo "  binary:     ${target}"
echo "  go:         $(go env GOVERSION)"
echo "  GOCACHE:    ${build_cache}"
echo "  GOMODCACHE: ${mod_cache}"

if [ "${dry_run}" -eq 1 ]; then
  echo "dry-run: no files written"
  exit 0
fi

mkdir -p "${bin_dir}" "${build_cache}" "${mod_cache}"

tmp="$(mktemp "${target}.tmp.XXXXXX")"
cleanup() {
  rm -f "${tmp}"
}
trap cleanup EXIT

(
  cd "${root}"
  GOCACHE="${build_cache}" GOMODCACHE="${mod_cache}" go build -trimpath -o "${tmp}" ./cmd/parley
)

chmod 0755 "${tmp}"
mv "${tmp}" "${target}"
trap - EXIT

mkdir -p "${prefix}"
installed_version="$("${target}" version 2>/dev/null || true)"
{
  echo "repo=${root}"
  echo "commit=${git_commit}"
  echo "version=${installed_version:-unknown}"
  echo "installed_at=$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  echo "binary=${target}"
  echo "go=$(go env GOVERSION)"
} > "${metadata}"

echo
echo "Installed ${target}"
echo "Version: ${installed_version:-unknown}"
echo
if ! command -v parley >/dev/null 2>&1 || [ "$(command -v parley)" != "${target}" ]; then
  echo "Add this to your shell profile if needed:"
  echo "  export PATH=\"${bin_dir}:\$PATH\""
  echo
fi
echo "Verify with:"
echo "  ${target} version"
