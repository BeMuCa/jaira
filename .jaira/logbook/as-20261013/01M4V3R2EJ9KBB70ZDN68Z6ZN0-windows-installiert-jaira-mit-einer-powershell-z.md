---
id: 01M4V3R2EJ9KBB70ZDN68Z6ZN0
title: Windows installiert jaira mit einer PowerShell-Zeile
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Auf Windows installiert eine Zeile 'irm https://raw.githubusercontent.com/BeMuCa/jaira/master/scripts/install.ps1 | iex' das neueste jaira-Release, ohne Adminrechte, so wie install.sh es auf Linux und macOS tut"
context: |-
  Was heute falsch ist: auf Windows gibt es keinen Installer. scripts/install.sh bricht bei uname != linux/darwin ab und schickt Leute zum Zip auf der Releases-Seite oder zu 'go install'.
  Folge: ein Teamkollege auf Windows muss das Zip von Hand holen, entpacken und den PATH selbst setzen - das widerspricht 'jeder kann es trivial installieren' (CLAUDE.md, Constraints).
  Ausloeser: Alex am 13.10. nach Release 0.4.1: 'Installer fuer Windows' als naechste Aufgabe.
  Entschieden (Alex, 13.10.): ein install.ps1 als Spiegel von install.sh - kein MSI/EXE (braeuchte Code-Signatur wegen SmartScreen), kein winget/scoop (Pflege bei jedem Release).
  Was es schon gibt: .goreleaser.yaml baut windows amd64/arm64 als Zip (goos windows, Zeile ~25-31) und schreibt eine Checksum-Datei; install.sh zeigt die Reihenfolge: Release finden, Archiv laden, Checksum pruefen, entpacken, ablegen.
  Ziel-Ort: %LOCALAPPDATA%\Programs\jaira (ueberschreibbar mit JAIRA_INSTALL_DIR wie bei install.sh), danach in den User-PATH, nicht den System-PATH.
  README: der Installations-Abschnitt nennt heute nur install.sh und den Zip-Hinweis.
  Achtung Windows: PowerShell 5.1 (vorinstalliert) muss reichen, nicht nur pwsh 7; TLS 1.2 dort explizit setzen.
definition-of-done: "scripts/install.ps1 laedt das neueste Release-Zip fuer die Architektur (amd64/arm64), prueft die SHA-256 gegen die Checksum-Datei des Releases und bricht bei Abweichung ab"
tags:
  - release
blocked-by: []
related: []
commits:
  - 7f996ba6cf7fc22128654f61ae72411fc5db5e5f
  - 97c019dc8f8cc9886cd8fe23bfc1839a4a75ffc1
  - 6f62ffc0e1c21a9cbdf9ff75d387575d2b49f77b
created-at: 2026-10-13T17:07:46Z
updated-at: 2026-10-13T17:44:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-201480
claimed-at: 2026-10-13T17:08:28Z
outcome-what: "Neues scripts/install.ps1: findet das neueste Release, laedt das Zip fuer amd64/arm64, prueft SHA-256 gegen checksums.txt, legt jaira.exe in %LOCALAPPDATA%\\Programs\\jaira (oder JAIRA_INSTALL_DIR) und haengt das Verzeichnis einmal an den User-PATH. README und install.sh nennen die Windows-Zeile, Release-Note 8Z6ZN0.md."
outcome-why: Auf Windows gab es keinen Installer; install.sh brach ab und schickte zum Zip - ein Windows-Kollege musste entpacken und PATH von Hand setzen.
outcome-resolves: "Windows installiert jaira mit 'irm .../install.ps1 | iex' ohne Adminrechte; ein zweiter Lauf aktualisiert. Getestet unter PowerShell 5.1 und 7.4."
review-summary: "Neues scripts/install.ps1 (137 Zeilen, reines ASCII): ermittelt die Architektur (PROCESSOR_ARCHITEW6432/ARCHITECTURE -> amd64/arm64, sonst Abbruch), holt den neuesten Tag ueber api.github.com, laedt jaira_<v>_windows_<arch>.zip und checksums.txt, prueft SHA-256 per .NET und wirft bei Abweichung, entpackt, benennt ein vorhandenes (evtl. laufendes) jaira.exe nach jaira.exe.old-<pid> um wie self upgrade, kopiert das neue hinein und haengt das Verzeichnis roh ueber HKCU:\\Environment (REG_EXPAND_SZ, nur wenn noch nicht drin) an den User-PATH, plus WM_SETTINGCHANGE-Broadcast und PATH der laufenden Sitzung. Alles in Funktionen mit throw statt exit, damit 'irm | iex' das Fenster nicht schliesst. README nennt die PowerShell-Zeile, install.sh verweist in Kopfkommentar und Fehlermeldung darauf, Release-Note core/release/notes/8Z6ZN0.md."
round: 2
review-gaps: "1) Kein Test im Repo: unit.ps1/e2e.ps1 lagen nur im Scratchpad, eine kuenftige Aenderung an install.ps1 ist durch nichts abgesichert (DoD verlangt keinen, daher kein Rueckweisungsgrund). 2) Kein Rollback: scheitert Copy-Item nach dem Umbenennen, liegt nur noch jaira.exe.old-<pid> da und jaira ist weg, bis man das Skript erneut startet. 3) Randfall: dritter Lauf in derselben PowerShell-Sitzung, waehrend das Ergebnis des zweiten laeuft - .old-<PID> ist dann gesperrt und Move-Item -Force darauf scheitert. Alles klein; nichts davon verletzt die DoD."
test-verdict: "pass: go build + go test ./... RC=0; Arch- und Checksum-Pruefungen (Treffer, Abweichung bricht ab, fehlende Zeile bricht ab) gruen unter PowerShell 5.1.22621 und 7.4.6; echtes 'irm <url> | iex' (Skript per HTTP aus WSL) installiert v0.4.1 zweimal je Shell, 1 PATH-Eintrag, User-PATH danach wiederhergestellt. Nicht live: der Standardordner %LOCALAPPDATA%\\Programs\\jaira (Tests liefen mit JAIRA_INSTALL_DIR, um den echten Ordner nicht anzufassen) und arm64-Hardware."
review-verdict: "Diff erfuellt alle fuenf DoD-Punkte; keine Defekte im Normalpfad. Selbst geprueft unter Windows PowerShell 5.1 (aus WSL): Datei parst, Test-JairaChecksum akzeptiert die richtige Summe und wirft 'checksum mismatch' bei falscher, Get-JairaArch liefert amd64. Archivname passt zum goreleaser-Default (jaira_<version>_windows_<arch>.zip, checksums.txt). Den echten Download-und-PATH-Lauf habe ich nicht wiederholt - dafuer steht nur das e2e des Implementierers (5.1.22621 und 7.4.6)."
review-round: 2
review-check: "1. Windows: PowerShell oeffnen (nicht als Admin). 2. irm https://raw.githubusercontent.com/BeMuCa/jaira/feat/8z6zn0/scripts/install.ps1 | iex  (vor dem Merge der Branch, danach master). 3. Erwartet: 'downloading jaira v0.4.1 (windows/amd64)...', dann 'installed jaira ... to C:\\Users\\<du>\\AppData\\Local\\Programs\\jaira\\jaira.exe' und 'added ... to your user PATH'. 4. Neues Terminal: jaira --version gibt die Version aus. 5. Zeile aus Schritt 2 nochmal: wieder 'installed ...', aber KEINE 'added ... PATH'-Zeile. 6. [Environment]::GetEnvironmentVariable('Path','User') -split ';' | Select-String jaira  -> genau ein Treffer. 7. Fenster ist nach beiden Laeufen noch offen."
---

# Windows installiert jaira mit einer PowerShell-Zeile

## Definition of Done

- [x] scripts/install.ps1 laedt das neueste Release-Zip fuer die Architektur (amd64/arm64), prueft die SHA-256 gegen die Checksum-Datei des Releases und bricht bei Abweichung ab
  proof: scripts/install.ps1 Get-JairaArch + Test-JairaChecksum; Install-Jaira laedt jaira_<v>_windows_<arch>.zip + checksums.txt; Test 'checksum mismatch throws' unter 5.1 und 7
- [x] jaira.exe landet in %LOCALAPPDATA%\Programs\jaira oder in JAIRA_INSTALL_DIR; das Verzeichnis kommt in den User-PATH, nur wenn es noch nicht drin ist; keine Adminrechte noetig
  proof: scripts/install.ps1 Install-Jaira (LOCALAPPDATA\Programs\jaira / JAIRA_INSTALL_DIR) + Add-JairaUserPath (HKCU:\Environment, nur wenn fehlt); e2e: 1 PATH-Eintrag nach 2 Laeufen, ohne Admin
- [x] Laeuft unter Windows PowerShell 5.1 und unter pwsh 7; ein zweiter Lauf aktualisiert statt zu doppeln
  proof: e2e-Lauf zweimal hintereinander unter powershell.exe 5.1.22621 und pwsh 7.4.6: beide installieren v0.4.1, zweiter Lauf ersetzt jaira.exe (Umbenennen nach .old-<pid>), PATH-Eintrag einmal
- [x] README nennt die Windows-Zeile neben install.sh; install.sh verweist bei Windows auf install.ps1 statt aufs Zip
  proof: README.md Install-Abschnitt (powershell-Block nach install.sh); scripts/install.sh Kopfkommentar + Fehlermeldung im case $os nennen install.ps1
- [x] Release-Note als Datei core/release/notes/<Handle>.md
  proof: core/release/notes/8Z6ZN0.md; go test ./core/release/... gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 17:08 · Alexander Sacharov** — parent session: jaira-6a
- **2026-10-13 17:08 · Alexander Sacharov** — Offene Entscheidungen: keine. Selbst entschieden: Architektur ueber PROCESSOR_ARCHITEW6432/PROCESSOR_ARCHITECTURE (AMD64->amd64, ARM64->arm64, sonst Abbruch - 32-bit x86 gibt es nicht als Release). Version ueber api.github.com/releases/latest wie install.sh. PATH-Vergleich ohne Gross/Klein und ohne abschliessenden Backslash. User-PATH ueber [Environment]::SetEnvironmentVariable(...,'User') und zusaetzlich in der laufenden Sitzung. Getestet wird mit powershell.exe 5.1 aus WSL (/mnt/c/Windows/System32/WindowsPowerShell).
- **2026-10-13 17:14 · Alexander Sacharov** — Gefunden: (1) Windows PowerShell 5.1, aus einer Umgebung mit pwsh-7-PSModulePath gestartet (so auf Alex' Rechner), findet Get-FileHash nicht - daher SHA256 ueber .NET. (2) Der zweite Lauf scheiterte zuerst an Copy-Item: das gerade installierte jaira.exe lief noch (vermutlich der Update-Check im Hintergrund). Loesung wie self upgrade: altes exe nach jaira.exe.old-<pid> umbenennen, dann kopieren, Reste wegfegen. (3) Datei ist bewusst reines ASCII: ohne BOM liest 5.1 eine lokal gespeicherte Datei als ANSI, und ein UTF-8-Gedankenstrich enthaelt Byte 0x94, das PowerShell als Anfuehrungszeichen liest. Tests: scratchpad unit.ps1 (Funktionen, Wegwerf-Key HKCU:\Software\JairaInstallTest) und e2e.ps1 (echter Release v0.4.1, zwei Laeufe, User-PATH vorher gesichert und danach exakt wiederhergestellt), beide gruen unter 5.1.22621 und 7.4.6.
- **2026-10-13 17:16 · Alexander Sacharov** — critique runde 1: scripts/install.ps1 Add-JairaUserPath: der Parameter $Key existiert nur fuer einen Test, der nicht im Repo liegt (scratchpad unit.ps1) - ein einziger Aufrufer, immer HKCU:\Environment; Parameter streichen und den Schluessel im Funktionskoerper festschreiben. Bewusst stehen gelassen: Add-Type/SendMessageTimeout statt [Environment]::SetEnvironmentVariable - letzteres schreibt REG_SZ und liest expandiert, wuerde also %USERPROFILE%-Eintraege im User-PATH fest verdrahten (rustup und scoop gehen denselben Weg). Expand-Archive bleibt: lief im e2e unter 5.1 mit dem pwsh-7-Modulpfad.
- **2026-10-13 17:16 · Alexander Sacharov** — critique runde 2: der eine Befund (Parameter $Key) ist erledigt, HKCU:\Environment steht jetzt im Funktionskoerper; e2e unter 5.1 und 7 erneut gruen. Nichts weiter.
