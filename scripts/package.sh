#!/bin/bash

set -e
set -u
set -o pipefail

readonly ROOT_DIR="$(cd "$(dirname "${0}")/.." && pwd)"
readonly BIN_DIR="${ROOT_DIR}/.bin"
readonly BUILD_DIR="${ROOT_DIR}/build"

# shellcheck source=SCRIPTDIR/.util/tools.sh
source "${ROOT_DIR}/scripts/.util/tools.sh"

# shellcheck source=SCRIPTDIR/.util/print.sh
source "${ROOT_DIR}/scripts/.util/print.sh"

function main {
  local version output token flags
  token=""

  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --version|-v)
        version="${2}"
        shift 2
        ;;

      --output|-o)
        output="${2}"
        shift 2
        ;;

      --token|-t)
        token="${2}"
        shift 2
        ;;

      --label)
        flags+=("--label" "${2}")
        shift 2
        ;;

      --help|-h)
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
    esac
  done

  if [[ -z "${version:-}" ]]; then
    usage
    echo
    util::print::error "--version is required"
  fi

  if [[ -z "${output:-}" ]]; then
    output="${BUILD_DIR}/buildpackage.cnb"
  fi

  repo::prepare

  tools::install "${token}"

  buildpack::archive "${version}"
  buildpack::release::archive
  buildpackage::create "${output}" "${flags[@]}"
}

function usage() {
  cat <<-USAGE
package.sh --version <version> [OPTIONS]

Packages the buildpack into a buildpackage .cnb file.

OPTIONS
  --help               -h            prints the command usage
  --version <version>  -v <version>  specifies the version number to use when packaging the buildpack
  --output <output>    -o <output>   location to output the packaged buildpackage artifact (default: ${ROOT_DIR}/build/buildpackage.cnb)
  --token <token>                    Token used to download assets from GitHub (e.g. jam, pack, etc) (optional)
USAGE
}

function repo::prepare() {
  util::print::title "Preparing repo..."

  rm -rf "${BUILD_DIR}"

  mkdir -p "${BIN_DIR}"
  mkdir -p "${BUILD_DIR}"

  export PATH="${BIN_DIR}:${PATH}"
}

function tools::install() {
  local token
  token="${1}"

  util::tools::jam::install \
    --directory "${BIN_DIR}" \
    --token "${token}"

  util::tools::pack::install \
    --directory "${BIN_DIR}" \
    --token "${token}"

  util::tools::yj::install \
    --directory "${BIN_DIR}" \
    --token "${token}"
}

function buildpack::archive() {
  local version
  version="${1}"

  util::print::title "Packaging buildpack into ${BUILD_DIR}/buildpack.tgz..."

  jam pack \
    --buildpack "${ROOT_DIR}/buildpack.toml" \
    --version "${version}" \
    --offline \
    --output "${BUILD_DIR}/buildpack.tgz"
}

function buildpack::release::archive() {
  local tmp_dir

  util::print::title "Packaging buildpack into ${BUILD_DIR}/buildpack-release-artifact.tgz..."

  tmp_dir=$(mktemp -d -p $BUILD_DIR)

  cat <<'README_EOF' > $tmp_dir/README.md
# Composite buildpack release artifact

This is a buildpack release artifact that contains everything needed to package and publish a composite buildpack. Composite buildpacks are a logic grouping of other buildpacks.

It contains the following files:

* `buildpack.toml` - this is needed because it contains the buildpacks and ordering information for the composite buildpack
* `package.toml` - this is needed because it contains the dependencies (and URIs) that let pack know where to find the buildpacks referenced in `buildpack.toml`.
  * `package.toml` can contain targets (platforms) for multi-arch support
* `build/buildpack.tgz` - this is added because it is referenced in `package.toml` by some buildpacks

## package locally

To package this buildpack to local .cnb file(s) run the following.

```
pack buildpack package mybuildpack.cnb --format file --config package.toml
```

## package and publish to a registry

To package this buildpack and publish it to a registry run the following.

* Note that as of pack v0.38.2 at least one target is required in package.toml or on the command line when publishing to a registry with `--publish`.

* replace SOME-REGISTRY with your registry (e.g. index.docker.io/yourdockerhubusername)
* replace SOME-VERSION with the version you want to publish (e.g. 0.0.1)

```
pack buildpack package SOME-REGISTRY/mybuildpack:SOME-VERSION --format image --config package.toml --publish
```
README_EOF

  mkdir -p $tmp_dir/build
  cp ${BUILD_DIR}/buildpack.tgz $tmp_dir/build
  cp ${ROOT_DIR}/package.toml $tmp_dir/
  # add the buildpack.toml from the tgz file because it has the version populated
  tar -xzf ${BUILD_DIR}/buildpack.tgz -C $tmp_dir/ buildpack.toml

  tar -czf ${BUILD_DIR}/buildpack-release-artifact.tgz -C $tmp_dir $(ls $tmp_dir)
  rm -rf $tmp_dir
}

function buildpackage::create() {
  local output flags release_archive_path tmp_dir
  output="${1}"
  flags=("${@:2}")
  release_archive_path="${BUILD_DIR}/buildpack-release-artifact.tgz"

  util::print::title "Packaging buildpack..."

  util::print::info "Extracting release archive..."
  tmp_dir=$(mktemp -d -p $BUILD_DIR)
  tar -xvf $release_archive_path -C $tmp_dir

  current_dir=$(pwd)
  cd $tmp_dir

  args=(
      --config package.toml
      --format file
    )

  args+=("${flags[@]}")

  # Use the local architecture to support running locally and in CI, which will be linux/amd64 by default.
  arch=$(util::tools::arch)

  # If package.toml has no targets we must specify one on the command line, otherwise pack will complain.
  # This is here for backward compatibility but eventually all package.toml files should have targets defined.
  if cat package.toml | yj -tj | jq -r .targets | grep -q null; then
    echo "package.toml has no targets so --target linux/${arch} will be passed to pack"
    args+=("--target linux/${arch}")
  fi

  pack \
    buildpack package "${output}" \
    ${args[@]}

  if [[ -e "${BUILD_DIR}/buildpackage-linux-${arch}.cnb" ]]; then
    echo "Copying linux-${arch} buildpackage to buildpackage.cnb"
    cp "${BUILD_DIR}/buildpackage-linux-${arch}.cnb" "${BUILD_DIR}/buildpackage.cnb"
  fi

  cd $current_dir
  rm -rf $tmp_dir
}

main "${@:-}"
