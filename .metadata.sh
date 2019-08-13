# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
BINARY="secspy"
# github username
GHUSER="davidnewhall"
# Github repo containing homebrew formula repo.
HBREPO="golift/homebrew-mugs"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift"
DESC="Command Line Interface for SecuritySpy (IP Camera NVR)"
GOLANGCI_LINT_ARGS="--enable-all -D gochecknoglobals"
CONFIG_FILE="secspy.conf"
LICENSE="MIT"
# FORMULA is either 'service' or 'tool'. Services run as a daemon, tools do not.
# This affects the homebrew formula (launchd) and linux packages (systemd).
FORMULA="tool"

export BINARY GHUSER HBREPO MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE LICENSE FORMULA

# The rest is mostly automatic.
# Fix the repo if it doesn't match the binary name.
# Provide a better URL if one exists.

# Used as go import path in docker and homebrew builds.
IMPORT_PATH="github.com/${GHUSER}/${BINARY}"
# Used for source links and wiki links.
SOURCE_URL="https://${IMPORT_PATH}"
# Used for documentation links.
URL="https://github.com/${GHUSER}/${BINARY}"

VERSION_PATH="${IMPORT_PATH}/cli.Version"

# Dynamic. Recommend not changing.
VVERSION=$(git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1))
VERSION="$(echo $VVERSION | tr -d v | grep -E '^\S+$' || echo development)"
# This produces a 0 in some envirnoments (like Homebrew), but it's only used for packages.
ITERATION=$(git rev-list --count --all || echo 0)
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
COMMIT="$(git rev-parse --short HEAD || echo 0)"

# Used by homebrew downloads.
SOURCE_PATH=https://codeload.${IMPORT_PATH}/tar.gz/v${VERSION}

export IMPORT_PATH SOURCE_URL URL VERSION_PATH VVERSION VERSION ITERATION DATE COMMIT SOURCE_PATH
