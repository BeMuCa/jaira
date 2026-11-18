#!/bin/sh
# Gathers what a Defender false-positive report needs, after a release
# build (32TTSS). Microsoft takes such reports only through the web form at
# https://www.microsoft.com/en-us/wdsi/filesubmission - there is no API
# outside an enterprise tenant - so this puts the files and their SHA256 in
# one place and prints a markdown summary for the run page.
#
# Run from the repository root after goreleaser: takes every
# dist/jaira_windows_<arch>*/jaira.exe - the signed one when
# signed-windows.sh swapped it in - and dist-extra/install.ps1, and writes
# them with SHA256SUMS to wdsi/.
set -eu

out=wdsi

rm -rf "$out"
mkdir -p "$out"
found=0
for exe in dist/jaira_windows_*/jaira.exe; do
	[ -f "$exe" ] || continue
	arch=$(basename "$(dirname "$exe")" | sed 's/^jaira_windows_\([a-z0-9]*\).*/\1/')
	cp "$exe" "$out/jaira_windows_$arch.exe"
	found=$((found + 1))
done
if [ "$found" -eq 0 ]; then
	echo "wdsi-bundle: no jaira.exe under dist/" >&2
	exit 1
fi
if [ -f dist-extra/install.ps1 ]; then
	cp dist-extra/install.ps1 "$out/install.ps1"
fi
(cd "$out" && sha256sum -- * > SHA256SUMS)

echo "## Defender false-positive report"
echo
echo "Submit each file at https://www.microsoft.com/en-us/wdsi/filesubmission"
echo "(\"Software developer\", \"Incorrectly detected as malware\") - the files"
echo "are in the artifact \`wdsi-submission\`. Steps: docs/WINDOWS.md."
echo
echo '| file | SHA256 |'
echo '| --- | --- |'
while read -r sum name; do
	echo "| $name | \`$sum\` |"
done < "$out/SHA256SUMS"
