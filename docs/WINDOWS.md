# jaira on Windows: Defender, signing and winget

Windows Defender has reported jaira as a trojan — `Trojan:Win32/Bearfoos.A!ml`
on `jaira.exe` and `Trojan:Script/Wacatac.C!ml` on `install.ps1`. Both names end
in `!ml`: they come from Defender's cloud machine-learning heuristics, not from a
signature of known malware. The flagged files were the official release
artifacts, byte for byte.

What the heuristics react to: an unsigned executable with no version
information that starts other programs (`git`, `gh`, `bash`, `wsl.exe`), and an
install script that downloads, unpacks and writes the user `PATH`.

## What a user can do

- Check the download is the release artifact: compare the zip's SHA256
  (`Get-FileHash jaira_<version>_windows_amd64.zip`) with its line in
  `checksums.txt` on the release page — that file lists the archives, not
  `jaira.exe` itself. `install.ps1` does this check on every install.
- Restore it from quarantine under *Windows Security > Virus & threat
  protection > Protection history*, or report it yourself at the form below.

## What every release already does

- `jaira.exe` carries a version resource (Properties > Details: product,
  company, version) — `cmd/jaira/winres/winres.json`, linked by the
  `go-winres` hook in `.goreleaser.yaml`.
- `install.ps1` refreshes `PATH` without compiling C# or calling into native
  Windows libraries, and is attached to every release.
- The release run collects what a false-positive report needs: the run page's
  summary lists every Windows file with its SHA256, and the artifact
  `wdsi-submission` holds the files themselves.

## Reporting a release to Microsoft (manual)

Microsoft takes false-positive reports only through a web form; its submission
API is limited to Defender for Endpoint enterprise tenants. After a release:

1. Open the release run in the Actions tab and download the artifact
   `wdsi-submission`.
2. Go to <https://www.microsoft.com/en-us/wdsi/filesubmission>, sign in, and
   choose *Software developer*.
3. Submit each `.exe` and `install.ps1` as *Incorrectly detected as malware*,
   with the detection name Defender showed and the SHA256 from the summary.
4. Answers usually take one to three days. Other vendors that flag it on
   VirusTotal (Ikarus, Trellix, Fortinet) have their own forms.

## Signing (off)

`release.yaml` has a `sign-windows` job that has SignPath sign `jaira.exe` and
`install.ps1`. It stays off until the repository variable
`SIGNPATH_ORGANIZATION_ID` is set; the release footer names SignPath only for a
release it actually signed. SignPath Foundation turned the free open-source
application down in November 2026 for lack of public visibility — a new
application, or a paid SignPath plan, turns it on. Setup is listed in the
header of `.github/workflows/release.yaml`.

## winget (off until a token is set)

goreleaser writes winget manifests (`BeMuCa.jaira`, both Windows zips) for every
release and, once a token is set, opens a pull request on
`microsoft/winget-pkgs`. A repository admin sets it up once:

1. Fork `microsoft/winget-pkgs` under the account or organisation that will
   open the pull requests.
2. Set the repository variable `WINGET_PKGS_OWNER` to that owner (default
   `BeMuCa`).
3. Create a token that may push to that fork and open pull requests on
   `microsoft/winget-pkgs` (classic PAT with `public_repo`), and set it as the
   secret `WINGET_GITHUB_TOKEN`.

Without the secret the manifests are only written to `dist/` and the release
runs as before. winget-pkgs validates every package with Defender: while
`jaira.exe` is still flagged, the pull request fails that check and needs the
false-positive report above first. Once a release is listed, users install
with `winget install BeMuCa.jaira`.
