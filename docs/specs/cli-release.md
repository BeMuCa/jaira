# CLI, Updates and Releases
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## Exit codes and output
CLI-01  The CLI shall exit 0 on success, 1 on an unexpected error, 2 on a usage error, 3 when a gate refuses, 4 on unresolved dependencies and 5 on an unknown or ambiguous id. [4539MZ, X8KZQ7]
CLI-03  Where --json is given, the CLI shall write the payload to stdout and a failure as {"error":{reason,code,message}} on stderr.
CLI-04  The CLI shall write hints and warnings only to stderr, so stdout stays the payload an agent parses. [T8A8KM, NH788Q]
CLI-05  If a ticket lives only on its ref, then a command needing its file shall refuse with exit 3 (on_ref_only). [4539MZ]
CLI-06  If the ticket file on disk is older than its ref, then a write shall refuse with exit 3 (stale_file). [F3VNVC]

## Reading the board
CLI-10  When 'jaira next' runs, the CLI shall report the one furthest-along actionable ticket no other live session has claimed.
CLI-11  Where --per-lane is given, 'jaira next' shall report one ticket per lane with waiting work, in pipeline order.
CLI-12  When a ticket waits only on an unmerged pull request, 'jaira next' shall list it apart (waiting_on_pr in JSON). [2EF48S]
CLI-13  Where --actionable is given, 'jaira list' shall show only tickets that could be started right now.
CLI-14  Every ticket in a --json payload shall carry next_lane, empty while parked, awaiting an answer or at the end.
CLI-15  When 'jaira show <id> --for-lane <lane> --json' runs, the CLI shall return prompt, input, model_tier, produces, mode and context.
CLI-16  Where --notes-last n is given, 'jaira show' shall print the newest n notes and say how many it hid.
CLI-17  Where --brief is given, 'jaira show' shall print one JSON line with title, goal, question, review and the newest four notes. [TQ5G6R]
CLI-18  When 'jaira json get <path>' reads JSON from stdin, the CLI shall print the value at the dotted path, and exit 5 for a missing path unless --default is given. [TQ5G6R]

## Writing tickets
CLI-20  When 'jaira create' runs on a board with a usable remote, the CLI shall put the ticket on its ref only and name 'jaira pull <id>'.
CLI-21  If the board has no usable remote, then 'jaira create' shall write a file and say why it is not on a ref (file-only-reason).
CLI-22  While on a branch whose name carries a ticket handle, 'jaira create' shall set follows to it and say so, unless --follows none.
CLI-23  When 'jaira create --lane' skips lanes, the CLI shall name the skipped lanes and the inputs left empty. [8WFNQY]
CLI-24  When 'jaira create --tag' names a tag new to the board, the CLI shall say so and point at 'jaira tags'. [79GEPW]
CLI-25  When 'jaira set' writes a field, the CLI shall rewrite only that field's bytes and leave every other line as it was.
CLI-26  Where --append is given, 'jaira set' shall add below a scalar or at the end of a list instead of overwriting. [K97RAV]
CLI-27  When 'jaira claim' runs, the CLI shall take a 30-minute lease and report taking over an expired one (took_over). [K1W2JR]
CLI-28  If another session's claim is still live, then 'jaira claim --release' shall refuse with exit 3. [7WY0YT]
CLI-29  If 'jaira claim' targets a rejected ticket, then the CLI shall refuse with exit 3 (rejected). [X8KZQ7]
CLI-30  If 'jaira reject' runs without --force, then the CLI shall refuse with exit 3: rejecting is a person's call. [A0C081]
CLI-31  When 'jaira reject' targets a ref-only ticket, the CLI shall fetch, write the rejection on the ref and clear the assignee. [4539MZ]
CLI-32  When 'jaira release <id>' runs, the CLI shall clear the assignee on the ref first and remove the file second; another person's ticket needs --force.

## Keeping a board in step with the binary
CLI-40  When the first command runs under a new binary version, jaira shall replace unedited shipped lane files, renew the jaira block where a file carries it, and say so in one stderr line. [NH788Q]
CLI-41  If a lane, role or skill file matches no version listed in core/shipped/shipped.sha256, then jaira shall not replace it. [PWAF22, 0K8Y6S]
CLI-42  While the board's version stamp differs from the binary, every command shall print a stderr line pointing at 'jaira update'. [RNQRJB]
CLI-43  While the binary is a dev build, jaira shall neither refresh lanes nor print the stale-board line. [DNAEPN]
CLI-44  When 'jaira update' runs, it shall regenerate the jaira section, stamp the version and print the notes since the last stamp. [PWAF22]
CLI-45  The 'jaira update' command shall never change whether a board is shared or private. [ETR0PX, 8W1R94, NNV39Y]
CLI-46  Where --roles is given, 'jaira update' shall refresh roles, cockpit mod and jaira skill wherever the roles are installed. [PWAF22]
CLI-47  Where --lanes is given, 'jaira update' shall replace this board's lane files that are an older shipped version. [PWAF22]
CLI-48  Where --dry-run is given, 'jaira update' shall say what it would write and write nothing, not even the stamp. [PWAF22]
CLI-49  If 'jaira update' leaves an edited file alone, then it shall name the file and exit 3. [PWAF22]
CLI-50  Where --agent-file agents|claude|both is given, 'jaira update' shall move the jaira section there, and the choice shall hold. [HE3WWW]
CLI-51  If the board's stamp is empty, dev or unknown, then 'jaira update' shall show every release's notes. [AR97K9]

## Upgrading the binary
CLI-55  When 'jaira self upgrade' runs, it shall download the latest release, verify its SHA-256 against checksums.txt, replace the binary and run the new binary's 'jaira update --roles'. [PWAF22]
CLI-56  If Homebrew or 'go install' owns the binary or its directory is unwritable, then self upgrade shall exit 3 naming the right way.
CLI-57  If the running binary is a dev build, then 'jaira self upgrade' without --check shall refuse with exit 3.
CLI-58  Where --check is given, self upgrade shall report latest and current and install nothing; with --version it shall exit 2.
CLI-59  Where --version is given, self upgrade shall install that exact release and, on a downgrade, leave skills and boards alone.
CLI-60  When the latest release is already installed, self upgrade shall exit 0 and download nothing.
CLI-61  If JAIRA_RELEASE_API or JAIRA_RELEASE_DOWNLOADS is set, then jaira shall accept only https or loopback http and warn on stderr.
CLI-62  Where jaira runs on Windows, self upgrade shall rename the running exe to .old-<pid> and sweep those leftovers on the next run. [5Q20DC]

## Newer-release notice
CLI-65  The CLI shall check for a newer release at most once an hour per machine, in a detached process, cached in ~/.jaira. [T8A8KM]
CLI-66  When the cache knows a newer release, 'jaira next', 'resume', 'claim' and 'move' shall print one stderr line telling the agent to preview, ask the user, and upgrade only on yes. [T8A8KM, EDSQZP]
CLI-67  If JAIRA_NO_UPDATE_CHECK=1 is set or the build is dev, then jaira shall neither check nor print the notice. [T8A8KM]

## Release notes
CLI-70  The binary shall embed its release notes and read them only from itself, never from a copy on disk. [AR97K9]
CLI-71  The binary shall read a note file directly in core/release/notes/ as Unreleased, and one in notes/<X.Y.Z>/ as that release. [AR97K9]
CLI-72  The binary shall read each line starting '- ' as one change and '- ! ' as one the reader must act on, ignoring all else. [AR97K9]

## Release workflow and CI
CLI-75  When the release-notes workflow runs on master with X.Y.Z, it shall move every note into notes/<version>/, run go test and the shipped-manifest check, and open a PR from release/<version>. [AR97K9, 776466, 0K8Y6S]
CLI-76  If the version is not X.Y.Z, already has a folder, has no notes or its tag exists, then the cut shall fail. [AR97K9]
CLI-77  When a release/<X.Y.Z> PR from this repository is merged, release-tag shall tag the merge commit and start release.yaml. [AR97K9]
CLI-78  The tag build shall run no tests and build linux, darwin and windows for amd64 and arm64, versioned from the tag. [776466]
CLI-79  The CI workflow shall skip a push or PR that changes only core/release/notes/**. [776466]
CLI-80  When a PR changes code, CI shall test on ubuntu with -race, macOS and Windows in three shards, and build Windows on ubuntu. [CQ1DCC, QFZDY7, APABM4]
CLI-81  If a PR changes only tickets or prose, then CI shall skip every test job. [CQ1DCC]
CLI-82  When code lands on master, CI shall run only the ubuntu -race job. [CQ1DCC]
CLI-83  If core/shipped/shipped.sha256 misses a file version that shipped in a v* tag, then CI shall fail. [0K8Y6S]
CLI-84  The repository shall pin LF endings for every embedded text file, release notes included. [5HDAJZ, APABM4]

## Install and Windows
CLI-88  When install.sh runs, it shall verify the archive's SHA-256, put jaira in ~/.local/bin (JAIRA_INSTALL_DIR) and run 'jaira roles install --global' unless JAIRA_NO_ROLES=1. [V38STG, Q9F9W4]
CLI-89  When install.ps1 runs, it shall verify SHA-256, put jaira.exe in %LOCALAPPDATA%\Programs\jaira without admin rights and add it to the user PATH once. [8Z6ZN0]
CLI-90  The install.ps1 script shall refresh PATH without Add-Type or DllImport. [32TTSS]
CLI-91  The release build shall give jaira.exe a version resource and publish a wdsi-submission artifact per run. [32TTSS]
CLI-92  Where WINGET_GITHUB_TOKEN is set, the release build shall open a winget-pkgs pull request. [32TTSS]
CLI-93  If a lock file is being deleted on Windows, then jaira shall wait for it instead of failing with Access is denied. [8X9SHE, 32TTSS]
CLI-94  If the tree holds a known Windows trap, then a Go test shall fail on Linux and name the remedy. [APABM4]
