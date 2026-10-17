#!/bin/sh
# Installs the latest jaira release binary.
#
#   curl -fsSL https://raw.githubusercontent.com/BeMuCa/jaira/master/scripts/install.sh | sh
#
# Downloads the release archive for this OS/arch, verifies its checksum, and
# puts the binary in ~/.local/bin (override with JAIRA_INSTALL_DIR). No root,
# no dependencies beyond curl or wget plus tar. Then it runs
# `jaira roles install --global` so the role skills and the cockpit are set up
# too; JAIRA_NO_ROLES=1 installs only the binary. On Windows, use install.ps1 in
# PowerShell instead:
#   irm https://raw.githubusercontent.com/BeMuCa/jaira/master/scripts/install.ps1 | iex
set -eu

REPO="BeMuCa/jaira"
INSTALL_DIR="${JAIRA_INSTALL_DIR:-$HOME/.local/bin}"

fetch() {
	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$1"
	elif command -v wget >/dev/null 2>&1; then
		wget -qO- "$1"
	else
		echo "error: neither curl nor wget is available" >&2
		exit 1
	fi
}

# download is fetch with a progress bar, shown only when stderr is a terminal
# (curl only; wget stays quiet), so `curl | sh` in CI logs stays clean.
download() {
	if [ -t 2 ] && command -v curl >/dev/null 2>&1; then
		curl -fL --progress-bar "$1"
	else
		fetch "$1"
	fi
}

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
linux | darwin) ;;
*)
	echo "error: unsupported OS '$os' — on Windows, run this in PowerShell instead:" >&2
	echo "  irm https://raw.githubusercontent.com/$REPO/master/scripts/install.ps1 | iex" >&2
	exit 1
	;;
esac

arch=$(uname -m)
case "$arch" in
x86_64 | amd64) arch=amd64 ;;
aarch64 | arm64) arch=arm64 ;;
*)
	echo "error: unsupported architecture '$arch'" >&2
	exit 1
	;;
esac

tag=$(fetch "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null |
	sed -n 's/^ *"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)
if [ -z "$tag" ]; then
	echo "error: no release found — install from source instead:" >&2
	echo "  go install github.com/$REPO/cmd/jaira@latest" >&2
	exit 1
fi
version=${tag#v}

archive="jaira_${version}_${os}_${arch}.tar.gz"
base="https://github.com/$REPO/releases/download/$tag"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "downloading jaira $tag ($os/$arch)…"
download "$base/$archive" >"$tmp/$archive"
fetch "$base/checksums.txt" >"$tmp/checksums.txt"

# Verify against the published checksum so a truncated or tampered download
# never lands on PATH. sha256sum on Linux, shasum on macOS.
(
	cd "$tmp"
	want=$(grep " $archive\$" checksums.txt | awk '{print $1}')
	if [ -z "$want" ]; then
		echo "error: $archive is not in checksums.txt" >&2
		exit 1
	fi
	if command -v sha256sum >/dev/null 2>&1; then
		got=$(sha256sum "$archive" | awk '{print $1}')
	else
		got=$(shasum -a 256 "$archive" | awk '{print $1}')
	fi
	if [ "$want" != "$got" ]; then
		echo "error: checksum mismatch for $archive" >&2
		exit 1
	fi
)

tar -xzf "$tmp/$archive" -C "$tmp" jaira
mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmp/jaira" "$INSTALL_DIR/jaira"

echo "installed $("$INSTALL_DIR/jaira" --version) to $INSTALL_DIR/jaira"
case ":$PATH:" in
*":$INSTALL_DIR:"*) ;;
*) echo "note: $INSTALL_DIR is not on your PATH" ;;
esac

# Roles and cockpit: one command writes both, so it is one step, not two. A
# failure here is printed as it came and does not undo the install above.
if [ -n "${JAIRA_NO_ROLES:-}" ]; then
	echo "skipped the skills and the Cockpit (JAIRA_NO_ROLES is set) — run: jaira roles install --global"
	exit 0
fi
echo "Setting up the skills and the Cockpit…"
if out=$("$INSTALL_DIR/jaira" roles install --global 2>&1); then
	echo "✓ $(printf '%s\n' "$out" | grep ' written, ' | tail -n1)"
	echo "active now: the jaira role skills and the Cockpit in ~/.claude/skills"
	echo "the Cockpit loads in your next Claude Code session"
else
	echo "✗ could not set up the skills and the Cockpit:" >&2
	printf '%s\n' "$out" >&2
	echo "the binary is installed; retry with: jaira roles install --global" >&2
fi
