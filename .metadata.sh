# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
TITLE="SecuritySpy CLI"
# github username
GHUSER="davidnewhall"
# docker hub username
DHUSER="golift"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift"
DESC="Command Line Interface for SecuritySpy (IP Camera NVR)"
GOLANGCI_LINT_ARGS="--enable-all -D gochecknoglobals"
CONFIG_FILE="secspy.conf"
LICENSE="MIT"

# The rest if mostly automatic.

BINARY="$(basename -s .git $(git config --get remote.origin.url))"
[ "$BINARY" != "" ] || BINARY="$(basename $(pwd))"
GHREPO="${GHUSER}/${BINARY}"
URL="https://github.com/${GHREPO}"

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
VERSION_PATH="github.com/${GHREPO}/cli.Version"

# Dynamic. Recommend not changing.
VERSION="$(git tag -l --merged | tail -n1 | tr -d v || echo development)"
# This produces a 0 in some envirnoments (like Homebrew), but it's only used for packages.
ITERATION=$(git rev-list --count --all || echo 0)
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
COMMIT="$(git rev-parse --short HEAD || echo 0)"

export TITLE GHUSER DHUSER MAINT VENDOR DESC CONFIG_FILE LICENSE GOLANGCI_LINT_ARGS
export BINARY GHREPO URL VERSION_PATH VERSION ITERATION DATE COMMIT
