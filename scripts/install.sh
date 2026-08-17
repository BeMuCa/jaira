#!/bin/sh
# Installs the latest jaira release binary.
#
#   curl -fsSL https://raw.githubusercontent.com/BeMuCa/jaira/master/scripts/install.sh | sh
#
# Downloads the release archive for this OS/arch, verifies its checksum, and
# puts the binary in ~/.local/bin (override with JAIRA_INSTALL_DIR). No root,
# no dependencies beyond curl or wget plus tar. On Windows, download the zip
# from the releases page instead, or use: go install github.com/BeMuCa/jaira/cmd/jaira@latest
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

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
linux | darwin) ;;
*)
	echo "error: unsupported OS '$os' — on Windows, grab the zip from https://github.com/$REPO/releases" >&2
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
fetch "$base/$archive" >"$tmp/$archive"
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
