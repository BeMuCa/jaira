---
id: 01M7668DV84RT7ZT0A8Y32TTSS
title: "Windows: Defender meldet jaira.exe und install.ps1 als Trojaner (ML-Heuristik)"
status: critique
ready: true
creator: Alexander Sacharov
goal: "jaira.exe und der Installer werden auf einem Windows-Rechner mit Microsoft Defender nicht mehr als Malware gemeldet, auch bei jeder neuen Version."
context: |-
  Was kaputt ist: Defender blockiert jaira auf Windows. Gemeldet am 2026-10-30 mit v0.5.7, installiert per 'irm .../scripts/install.ps1 | iex' und per Zip-Download in Edge.
  Zwei Detektionen, beide '!ml' (Cloud-ML-Heuristik, keine Signatur): Trojan:Win32/Bearfoos.A!ml auf jaira.exe, Trojan:Script/Wacatac.C!ml auf dem Installer-Skript.
  VirusTotal: v0.5.8 jaira.exe (SHA256 1e85ecb8...dd65) 3/70 - Ikarus Trojan.WinGo.Spy, Trellix Artemis, Fortinet Malicious_Behavior.SB; Microsoft dort sauber. v0.5.12 (3ddf00e8...4678) gleiches Bild.
  Hashes geprüft: die gemeldeten Dateien sind die offiziellen Release-Artefakte, nichts untergeschoben.
  Warum es auslöst: jaira.exe ist unsigniert, ohne Versionsinfo (.goreleaser.yaml baut mit -s -w, kein Resource), startet git/gh/bash/wsl.exe und liest Dateien im Home. scripts/install.ps1 sieht aus wie ein Dropper: irm|iex, Download nach %TEMP%, Expand-Archive, exe kopieren, HKCU\Environment Path schreiben, Add-Type mit DllImport user32 SendMessageTimeout.
  Schon getan: False-Positive-Meldung an Microsoft (wdsi/filesubmission) und Ikarus/Fortinet/Trellix macht Alex von Hand.
  Offen: eine Firmen-Identität für Azure Trusted Signing (Azure-Konto und Prüfung der juristischen Person) - das kann kein Agent.
definition-of-done: "Windows-Binaries im Release sind Authenticode-signiert (z. B. Azure Trusted Signing im release.yaml), signtool verify /pa ist grün"
tags:
  - release
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-11-11T20:54:42Z
updated-at: 2026-11-17T21:36:31Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-46806
claimed-at: 2026-11-17T21:31:20Z
question: ""
outcome-what: "Without a signature: winget manifests and a winget-pkgs PR per release via goreleaser (off unless WINGET_GITHUB_TOKEN), a wdsi-submission artifact plus SHA256 summary per release run (scripts/wdsi-bundle.sh), docs/WINDOWS.md for the manual steps; SignPath line only on a release that was really signed"
outcome-why: "SignPath Foundation turned the application down on 2026-11-17; Alex chose to lower false positives without a signature and add winget as a distribution path"
outcome-resolves: "DoD 5, 6, 7; DoD 1 and the signing half of DoD 3 stay open for Alex at signoff"
review-summary: "none - round 2 reads only the testing finding: install.ps1:66-70 comment reworded, no API names left, meaning unchanged"
review-gaps: "none - round 2: one comment reworded, nothing to remove"
test-verdict: "pass so far: go vet + go test ./... green (rc=0, round 1; round 2 changed only a .ps1 comment); linux + GOOS=windows build ok; DoD 2 read from /tmp/j.exe (CompanyName 'jaira project', ProductName 'jaira', FileVersion 0.5.13.0, asInvoker manifest); install.ps1: 0 hits for Add-Type|DllImport; DoD 4/5 present; signed-windows.sh exercised on a scratch tree (arm64 swap, amd64 no-op without signed file, linux no-op, stage plain/signed, bad arg exit 2). Not run: goreleaser and release.yaml (first tag), install.ps1 on Windows (no pwsh), SignPath signing (DoD 1 + signing part of DoD 3, waiting for approval)."
round: 2
blocked-reason: "wartet auf SignPath-Genehmigung, Antrag 2026-11-11"
---

# Windows: Defender meldet jaira.exe und install.ps1 als Trojaner (ML-Heuristik)

## Definition of Done

- [ ] Windows-Binaries im Release sind Authenticode-signiert (z. B. Azure Trusted Signing im release.yaml), signtool verify /pa ist grün
- [x] jaira.exe trägt eine Versionsinfo (CompanyName, ProductName, FileVersion), sichtbar unter Eigenschaften > Details
  proof: /tmp/j.exe (go-winres v0.3.3 make + GOOS=windows go build, run by Alex 2026-11-11): CompanyName 'jaira project', ProductName 'jaira', FileVersion 0.5.13.0, manifest asInvoker; source cmd/jaira/winres/winres.json
- [~] scripts/install.ps1 kommt ohne Add-Type/DllImport aus und ist ebenfalls signiert
  proof: scripts/install.ps1: grep Add-Type|DllImport = 0; Send-JairaEnvChange uses [Environment]::SetEnvironmentVariable; signing pending SignPath
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/32TTSS.md
- [x] README nennt SignPath nicht, solange nicht signiert wird; der Release-Text (.goreleaser.yaml release.footer) nennt 'Free code signing provided by SignPath.io, certificate by SignPath Foundation' nur bei einem Release, das sign-windows tatsaechlich signiert hat
  proof: .goreleaser.yaml release.footer (envOrDefault WINDOWS_SIGNED), .github/workflows/release.yaml goreleaser step env WINDOWS_SIGNED; README.md install section links docs/WINDOWS.md instead
- [x] winget: jedes Release legt per release.yaml automatisch einen PR mit dem Manifest auf microsoft/winget-pkgs an (z. B. winget-releaser oder wingetcreate), aktiv nur wenn das Token-Secret gesetzt ist; ohne Secret laeuft das Release unveraendert
  proof: .goreleaser.yaml winget (BeMuCa.jaira, skip_upload unless WINGET_GITHUB_TOKEN non-empty); release.yaml env WINGET_GITHUB_TOKEN/WINGET_PKGS_OWNER; goreleaser v2.12.7 check ok, snapshot wrote dist/winget/manifests/b/BeMuCa/jaira/*/ (zip, portable, x64+arm64)
- [x] WDSI-Meldung pro Release so weit automatisiert, wie Microsoft es zulaesst: der Release-Job legt Dateien und SHA256 fuer die Fehlalarm-Meldung bereit und nennt den Link; was manuell bleibt, steht in docs
  proof: scripts/wdsi-bundle.sh (run on a snapshot dist: 2 exe + install.ps1 + SHA256SUMS, summary table, exit 1 without exe); release.yaml steps 'Collect what a false-positive report needs' + artifact wdsi-submission; manual steps docs/WINDOWS.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Versionsinfo: winres.json + go-winres im goreleaser before-hook, .syso gitignored
  proof: .goreleaser.yaml before.hooks go-winres; cmd/jaira/winres/winres.json
- [x] Versionsinfo im gebauten jaira.exe prüfen (PE-Ressource auslesen)
  proof: /tmp/j.exe read with a VS_VERSIONINFO parser: all strings + fixed 0.5.13.0 present
- [x] install.ps1: Broadcast ohne Add-Type/DllImport
  proof: scripts/install.ps1 Send-JairaEnvChange
- [x] release.yaml: SignPath-Signierjob für jaira.exe + install.ps1, aktiv nur mit Secret
  proof: .github/workflows/release.yaml sign-windows; scripts/signed-windows.sh; .signpath/artifact-configuration.xml
- [x] Release-Note 32TTSS.md
  proof: core/release/notes/32TTSS.md

## Progress
- **2026-11-11 21:08 · Alexander Sacharov** — False-Positive-Meldungen an Microsoft am 2026-11-11 abgeschickt: (1) jaira.exe v0.5.7, Trojan:Win32/Bearfoos.A!ml, SHA256 bc57ef9c...1b29; (2) scripts/install.ps1 (master, SHA256 d865b65d...c470), Trojan:Script/Wacatac.C!ml. Antwort von Microsoft steht aus, meist 1-3 Tage.
- **2026-11-11 21:11 · Alexander Sacharov** — decision: jaira ist Open Source (MIT, öffentlich), keine Firmensoftware - signiert wird über SignPath Foundation (kostenlos für OSS, kein Azure, keine juristische Person). Antrag bei signpath.org stellt der Repo-Owner von Hand; bis dahin zuerst Versionsinfo in jaira.exe und install.ps1 ohne Add-Type/DllImport, der Signierschritt in release.yaml wird für SignPath vorbereitet.
- **2026-11-11 21:12 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/460629.sock
- **2026-11-11 21:12 · Alexander Sacharov** — decision: PATH-Broadcast ohne Add-Type - [Environment]::SetEnvironmentVariable mit Target User sendet WM_SETTINGCHANGE selbst; install.ps1 setzt und löscht damit eine Hilfsvariable statt SendMessageTimeout per DllImport. Versionsinfo: go-winres erzeugt beim Release eine .syso in cmd/jaira (goreleaser before-hook, Version aus dem Tag), nicht eingecheckt; CompanyName 'jaira project', ProductName 'jaira'.
- **2026-11-11 21:13 · Alexander Sacharov** — question: DoD 1 (Authenticode-Signatur, signtool verify grün) und der Signatur-Teil von DoD 3 kann kein Agent erfüllen, solange SignPath den Antrag nicht genehmigt hat - der Gate verweigert done. Optionen: (1) Empfohlen: DoD 1 und 3 umformulieren auf 'Signierschritt für SignPath in release.yaml vorbereitet, abschaltbar per Secret' + Folgeticket fürs echte Signieren nach Genehmigung; (2) Ticket bleibt offen: Versionsinfo + install.ps1 jetzt, dann blocked bis SignPath genehmigt; (3) Signier-Teil ganz raus, eigenes Ticket, dieses nur Versionsinfo + install.ps1.
- **2026-11-11 21:15 · Alexander Sacharov** — decision (Alex, 2026-11-11): Option 2 - DoD bleibt wie sie ist. Jetzt Versionsinfo + install.ps1 ohne Add-Type, Signierschritt in release.yaml für SignPath vorbereiten; danach Ticket nach blocked (wartet auf SignPath-Genehmigung), kein Folgeticket.
- **2026-11-11 21:15 · Alexander Sacharov** — decision: Path bleibt wie heute roh gelesen und per Set-ItemProperty -Type ExpandString geschrieben - %USERPROFILE%-Einträge bleiben unexpandiert. Nur der Broadcast ändert sich: [Environment]::SetEnvironmentVariable setzt und löscht eine Hilfsvariable JAIRA_INSTALL_ENV_REFRESH (Target User); .NET schickt dabei selbst WM_SETTINGCHANGE 'Environment'. Path wird nie über SetEnvironmentVariable geschrieben.
- **2026-11-11 21:18 · Alexander Sacharov** — SignPath Foundation verlangt den Hinweis auf README und Download-Seite (Teamlead 2026-11-11, Antrag bei signpath.org heute gestellt) - DoD 5.
- **2026-11-11 21:20 · Alexander Sacharov** — Signierweg: sign-windows-Job baut mit 'goreleaser build' dieselben Binaries (reproduzierbar: CGO off, -trimpath, gleiche Version/Resource), lädt windows_amd64/arm64 jaira.exe + install.ps1 als Artifact hoch, SignPath-Action v2 signiert (Foundation verlangt Herkunft aus GitHub Actions, darum Artifact und nicht REST-Upload). Der Release-Job lädt 'signed' herunter; goreleaser build post-hook (scripts/signed-windows.sh exe) ersetzt jaira.exe vor dem Zippen, before-hook 'stage' legt install.ps1 (signiert falls vorhanden) als Release-Asset bei. Aus solange vars.SIGNPATH_ORGANIZATION_ID leer. Offen nach Genehmigung: Variablen/Secret setzen, artifact-configuration.xml in SignPath einspielen, README-Einzeiler auf releases/latest/download/install.ps1 umstellen (raw master bleibt unsigniert).
- **2026-11-11 21:20 · Alexander Sacharov** — question: DoD 2 braucht einen Beleg, dass jaira.exe die Versionsinfo trägt. Lokal go-winres auszuführen (go run github.com/tc-hib/go-winres@v0.3.3, externer Code) hat der Auto-Mode verweigert. Optionen: (1) Empfohlen: Alex erlaubt es einmal oder führt es selbst aus: '! cd <worktree> && go run github.com/tc-hib/go-winres@v0.3.3 make --in cmd/jaira/winres/winres.json --out cmd/jaira/rsrc --arch amd64 --file-version 0.5.13.0 --product-version 0.5.13.0 && GOOS=windows go build -o /tmp/j.exe ./cmd/jaira' - danach lese ich die Ressource aus; (2) Beleg erst beim ersten Release nach dem Merge (Windows: Eigenschaften > Details) - DoD 2 bleibt bis dahin offen, das Ticket ist ohnehin blocked.
- **2026-11-11 21:27 · Alexander Sacharov** — SignPath-Foundation-Antrag am 2026-11-11 von Alex abgeschickt (Homepage = GitHub-Repo, Download-URL = Releases-Seite, Privacy leer - keine Nutzerdaten). blocked_reason beim Parken: wartet auf SignPath-Genehmigung, Antrag 2026-11-11.
- **2026-11-11 21:28 · Alexander Sacharov** — decision (Alex): DoD-2-Beleg über Option 1 - Alex führt go-winres + GOOS=windows go build selbst aus, ich lese danach /tmp/j.exe aus.
- **2026-11-11 21:33 · Alexander Sacharov** — testing: parked in blocked out of testing. On return: after SignPath approval set vars SIGNPATH_ORGANIZATION_ID/PROJECT_SLUG/SIGNING_POLICY_SLUG + secret SIGNPATH_API_TOKEN, paste .signpath/artifact-configuration.xml into the SignPath project, cut a release, check signtool verify /pa on jaira.exe and Get-AuthenticodeSignature on the release's install.ps1, switch the README one-liner to releases/latest/download/install.ps1, tick DoD 1 and 3, then human -> review.
- **2026-11-11 21:33 · Alexander Sacharov** — testing: install.ps1:69 the new comment names Add-Type and DllImport; Defender's script ML reads text, comments included. Fix: say 'no compiled C# and no native user32 call' instead.
- **2026-11-11 21:48 · Alexander Sacharov** — SignPath hat den Antrag 2026-11-11 bestätigt: wird bearbeitet. Kein PR vor der Genehmigung - Alex wartet.
- **2026-11-11 21:59 · Alexander Sacharov** — PR #135 auf Alex' Anweisung vor der SignPath-Genehmigung geöffnet: Signierschritt bleibt aus, bis vars.SIGNPATH_ORGANIZATION_ID gesetzt ist. Variablen/Secret kann nur ein Repo-Admin setzen (Alex hat push, kein admin).
- **2026-11-16 22:44 · Alexander Sacharov** — CI fix on PR #135: (1) ticket-commits failed on 8a2d540c, a commit carrying only this ticket file - folded into the next code commit, branch force-pushed with lease. (2) test (windows-latest) TestRegisterTagsLosesNoTagUnderConcurrentWriters: 'Access is denied' on tags.lock. Cause in core/ticket/store.go lock(): on Windows a just-removed lock file is delete-pending and O_EXCL create returns access denied, not exists; it is now waited for like a busy lock, bounded by lockTimeout. Not caused by 32TTSS, but it blocks this PR and every other PR's Windows run.
- **2026-11-17 21:21 · Alexander Sacharov** — SignPath Foundation hat am 17.11. abgelehnt: das Projekt hat noch zu wenig oeffentliche Sichtbarkeit (Stars, Forks, Contributors, externe Erwaehnungen). Erneut bewerben, sobald es bekannter ist; sonst SignPath als bezahltes Abo. Signieren ueber die Foundation faellt damit vorerst weg.
- **2026-11-17 21:30 · Alexander Sacharov** — decision (Alex, 17.11.): nach der SignPath-Absage Weg 3 jetzt umsetzen - ohne Signatur Fehlalarme senken, so automatisch wie moeglich, und winget als Verteilweg. Versionsinfo und install.ps1 ohne Add-Type sind schon da. Signatur (DoD 1) bleibt offen; ob sie fuer done noetig ist, entscheidet Alex in signoff. Achtung: winget-pkgs prueft jedes Paket mit Defender - solange jaira.exe gemeldet wird, faellt der winget-PR durch; das Token-Secret kann nur ein Repo-Admin (Berk) setzen. Zuerst auf upstream/master rebasen.
- **2026-11-17 21:31 · Alexander Sacharov** — Achtung DoD 5: README und release.footer nennen 'Free code signing provided by SignPath.io' - nach der Absage stimmt das nicht mehr. Den Satz entfernen oder erst bei aktiver Signatur ausgeben; DoD 5 entsprechend anpassen.
- **2026-11-17 21:31 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/51013.sock
- **2026-11-17 21:32 · Alexander Sacharov** — decision: rebased on upstream/master (0.5.13). Weg 3 settled without a question: (a) DoD 5 - SignPath line out of README; release.footer prints it only when sign-windows signed this release (env WINDOWS_SIGNED), so a later SignPath approval needs no edit. (b) winget via goreleaser's own winget publisher (OSS, zip archives, use ''), not winget-releaser/wingetcreate: it runs in the existing release job and writes manifests for both arches; package identifier BeMuCa.jaira (publisher = GitHub owner, permanent once in winget-pkgs - Alex can change before the secret is set). PR from a winget-pkgs fork (owner = vars.WINGET_PKGS_OWNER, default BeMuCa) into microsoft/winget-pkgs; skipped unless secret WINGET_GITHUB_TOKEN is non-empty. (c) WDSI has no public submission API (the MDE submission API needs an enterprise tenant): release job bundles jaira.exe x2 + install.ps1 + SHA256SUMS as artifact wdsi-submission and writes hashes and the wdsi/filesubmission link into the job summary; manual steps in docs/WINDOWS.md.
