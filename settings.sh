# Application Builder Configuration File. Customized for: hello-world
# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
# Powered by Application Builder: https://github.com/golift/application-builder

# Bring in dynamic repo/pull/source info.
source $(dirname "${BASH_SOURCE[0]}")/init/buildinfo.sh

# Must match the repo name to make things easy. Otherwise, fix some other paths.
BINARY="secspy"
# github username
REPO="davidnewhall/secspy"
# Github repo containing homebrew formula repo.
HBREPO="golift/homebrew-mugs"
MAINT="David Newhall II <david at sleepers dot pro>"
DESC="Command Line Interface for SecuritySpy (IP Camera NVR)"
GOLANGCI_LINT_ARGS="--enable-all -D forbidigo,nlreturn,exhaustivestruct"
CONFIG_FILE="secspy.conf"
LICENSE="MIT"
# FORMULA is either 'service' or 'tool'. Services run as a daemon, tools do not.
# This affects the homebrew formula (launchd) and linux packages (systemd).
FORMULA="tool"
# Used for source links and wiki links.
SOURCE_URL="https://github.com/${REPO}"

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
# This makes a path like golift.io/version.Version=1.3.3
# Name the Version-containing library the same as the github repo, without dashes.
VERSION_PATH="golift.io/version"

# Used by homebrew downloads.
SOURCE_PATH=https://codeload.github.com/${GHUSER}/${REPO}/tar.gz/v${VERSION}

export BINARY HBREPO MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE
export LICENSE FORMULA SOURCE_URL VERSION_PATH SOURCE_PATH

### Optional ###

# Import this signing key only if it's in the keyring.
gpg --list-keys 2>/dev/null | grep -q B93DD66EF98E54E2EAE025BA0166AD34ABC5A57C
[ "$?" != "0" ] || export SIGNING_KEY=B93DD66EF98E54E2EAE025BA0166AD34ABC5A57C

export WINDOWS_LDFLAGS=""
export MACAPP=""
export EXTRA_FPM_FLAGS=""
