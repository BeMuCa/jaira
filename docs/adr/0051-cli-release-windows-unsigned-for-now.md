# Windows releases ship unsigned for now, with winget and a false-positive bundle
Status: Accepted
Date: 2026-11-18

## Context
Defender flagged jaira.exe and install.ps1 as trojans (ML heuristics on unsigned binaries and installers).
jaira is open source (MIT); paid signing infrastructure and a legal entity are not available.

## Decision
No signature yet: jaira.exe gets a version resource (go-winres), install.ps1 refreshes PATH without Add-Type/DllImport,
each release publishes a wdsi-submission artifact and, with WINGET_GITHUB_TOKEN, a winget-pkgs PR.
Signing via SignPath Foundation is a later step; release.yaml signs only when SIGNPATH_ORGANIZATION_ID is set.

## Rejected
- Azure Trusted Signing (needs Azure and a legal entity).
- Native Windows calls in install.ps1 (Add-Type/SendMessageTimeout, a heuristic trigger).

## Consequences
Users may still see a Defender warning; docs/WINDOWS.md holds the manual report steps.

## Sources
32TTSS (decision: 2026-11-11 21:11, 21:12, 21:15; 2026-11-17 21:32; 2026-11-18 08:47), 8Z6ZN0
