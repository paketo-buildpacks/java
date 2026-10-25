#!/usr/bin/env bash

set -eu
set -o pipefail

readonly ROOT_DIR="$(cd "$(dirname "${0}")/.." && pwd)"
readonly BIN_DIR="${ROOT_DIR}/.bin"

# shellcheck source=SCRIPTDIR/.util/tools.sh
source "${ROOT_DIR}/scripts/.util/tools.sh"

# shellcheck source=SCRIPTDIR/.util/print.sh
source "${ROOT_DIR}/scripts/.util/print.sh"

function main {
  local archive_path image_ref token
  token=""

  while [[ "${#}" != 0 ]]; do
    case "${1}" in
    --archive-path | -a)
      archive_path="${2}"
      shift 2
      ;;

    --image-ref | -i)
      image_ref="${2}"
      shift 2
      ;;

    --token | -t)
      token="${2}"
      shift 2
      ;;

    --help | -h)
      shift 1
      usage
      exit 0
      ;;

    "")
      # skip if the argument is empty
      shift 1
      ;;

    *)
      util::print::error "unknown argument \"${1}\""
      ;;
    esac
  done

  if [[ -z "${image_ref:-}" ]]; then
    usage
    util::print::error "--image-ref is required"
  fi

  if [[ -z "${archive_path:-}" ]]; then
    util::print::info "Using default archive path: ${ROOT_DIR}/build/buildpack-release-artifact.tgz"
    archive_path="${ROOT_DIR}/build/buildpack-release-artifact.tgz"
  else
    archive_path="${archive_path}"
  fi

  repo::prepare

  tools::install "${token}"

  buildpack::publish "${image_ref}" "${archive_path}"
}

function usage() {
  cat <<-USAGE
Publishes a composite buildpack to a registry.

OPTIONS
  -a, --archive-path <filepath>       Path to the buildpack release artifact (default: ${ROOT_DIR}/build/buildpack-release-artifact.tgz) (optional)
  -h, --help                          Prints the command usage
  -i, --image-ref <ref>               List of image reference to publish to (required)
  -t, --token <token>                 Token used to download assets from GitHub (e.g. jam, pack, etc) (optional)

USAGE
}

function repo::prepare() {
  util::print::title "Preparing repo..."

  mkdir -p "${BIN_DIR}"

  export PATH="${BIN_DIR}:${PATH}"
}

function tools::install() {
  local token
  token="${1}"

  util::tools::pack::install \
    --directory "${BIN_DIR}" \
    --token "${token}"

  util::tools::yj::install \
    --directory "${BIN_DIR}" \
    --token "${token}"
}

function buildpack::publish() {
  local image_ref archive_path tmp_dir
  image_ref="${1}"
  archive_path="${2}"

  util::print::title "Publishing composite buildpack..."

  util::print::info "Extracting archive..."
  tmp_dir=$(mktemp -d -p $ROOT_DIR)
  tar -xvf $archive_path -C $tmp_dir

  util::print::info "Publishing buildpack to ${image_ref}"

  current_dir=$(pwd)
  cd $tmp_dir

  # If package.toml has no targets we must specify one on the command line, otherwise pack will complain.
  # This is here for backward compatibility but eventually all package.toml files should have targets defined.
  targets=""
  if cat package.toml | yj -tj | jq -r .targets | grep -q null; then
    # Use the local architecture to support running locally and in CI, which will be linux/amd64 by default.
    arch=$(util::tools::arch)
    targets="--target linux/${arch}"
    echo "package.toml has no targets so ${targets} will be used"
  fi

  pack \
    buildpack package "${image_ref}" \
    --config package.toml \
    --flatten \
    --format image \
    --publish \
    ${targets}

  cd $current_dir
  rm -rf $tmp_dir
}

main "${@:-}"
