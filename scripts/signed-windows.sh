#!/bin/sh
# Swaps the SignPath-signed Windows files into a release build (32TTSS).
#
# release.yaml's sign-windows job leaves them in ./signed with the layout of
# .signpath/artifact-configuration.xml:
#   signed/windows_amd64/jaira.exe  signed/windows_arm64/jaira.exe
#   signed/install.ps1
# Without ./signed - signing not set up, a local snapshot, the sign-windows
# job itself - every command here leaves the build as it is.
#
#   signed-windows.sh exe <binary path> <goreleaser target>
#       goreleaser build post hook: overwrite a windows binary with its
#       signed copy
#   signed-windows.sh stage
#       goreleaser before hook: put install.ps1, signed when there is a
#       signed one, in dist-extra/ for the release to attach
set -eu

case "${1:-}" in
exe)
	path=$2 target=$3
	case $target in
	windows_amd64*) arch=amd64 ;;
	windows_arm64*) arch=arm64 ;;
	*) exit 0 ;;
	esac
	signed="signed/windows_$arch/jaira.exe"
	[ -f "$signed" ] || exit 0
	cp "$signed" "$path"
	echo "signed-windows: $path is the signed build"
	;;
stage)
	rm -rf dist-extra
	mkdir -p dist-extra
	if [ -f signed/install.ps1 ]; then
		cp signed/install.ps1 dist-extra/install.ps1
	else
		cp scripts/install.ps1 dist-extra/install.ps1
	fi
	;;
*)
	echo "usage: signed-windows.sh exe <path> <target> | stage" >&2
	exit 2
	;;
esac
