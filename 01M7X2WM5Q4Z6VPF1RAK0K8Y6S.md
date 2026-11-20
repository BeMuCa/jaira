---
id: 01M7X2WM5Q4Z6VPF1RAK0K8Y6S
title: "shipped.sha256 kennt jede Version, die in einem Release-Tag ausgeliefert wurde, und CI haelt das fest"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "core/shipped/shipped.sha256 enthaelt den Hash jeder Datei-Version, die je in einem v*-Tag stand; ein CI-Schritt schlaegt fehl, sobald eine ausgelieferte Version fehlt, sodass jaira update --roles eine unveraenderte alte Installation nie mehr als Handaenderung ueberspringt."
context: |-
  Gefunden vom bughunt auf TW5FE2 (20.11.2026, Note 16:50): 8 Datei-Versionen aus Release-Tags fehlen im Manifest, z.B. core/mod/builtin/jaira-cockpit/hooks/board.ts aus v0.5.2/v0.5.3, dazu Versionen von register.tsx, texts.ts, types/index.d.ts und jaira-role-acceptance/template.html.
  Folge: jaira update --roles (core/role/install.go:147 writeFile) haelt eine unberuehrte Installation dieser Version fuer von Hand geaendert und ueberspringt sie (Skipped, nur mit --force). Der Nutzer bekommt still keine neue Rolle/kein neues Cockpit.
  Nachweis: scripts/shipped-history.sh erzeugt 59 Zeilen, die im committeten Manifest fehlen; 8 davon standen in v*-Tags (git ls-tree je Tag); git log -S<hash> auf dem Manifest ist leer, sie waren nie drin.
  Ursache vermutlich: das Manifest wird bei Merges von Hand zusammengefuehrt (Konflikte append-only aufgeloest), Zwischenstaende gehen verloren.
  Alex 20.11.2026: eigenes Ticket, ja.
definition-of-done: "Alle Datei-Versionen aus allen v*-Tags stehen in core/shipped/shipped.sha256 (Test oder Skript-Pruefung)"
tags:
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-11-20T18:18:21Z
updated-at: 2026-11-20T22:54:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2906
claimed-at: 2026-11-20T22:45:03Z
outcome-what: "Manifest neu erzeugt (8 Tag-Versionen fehlten, u.a. cockpit board.ts v0.5.2/0.5.3); scripts/shipped-history.sh --check haelt es gegen jeden v*-Tag und laeuft in ci.yaml (Job shipped-manifest) und im release-notes-Workflow; .gitattributes merge=union fuer das Manifest; LC_ALL=C im Skript"
outcome-why: Hand-gemergte Manifest-Konflikte verloren Zeilen; jaira update --roles hielt dann eine unberuehrte Installation fuer handgeaendert und uebersprang sie still
outcome-resolves: "Jede ausgelieferte Version steht im Manifest, CI bricht ab, wenn eine fehlt, und Manifest-Merges konfliktieren nicht mehr"
review-summary: |-
  scripts/shipped-history.sh:61-67 prueft das Zeilenformat (Konfliktmarker) ein zweites Mal - core/shipped/shipped_test.go TestManifestLinesAreWellFormed tut das schon in CI und im release-notes-Lauf; Block und den Konfliktmarker-Fall in history_test.go streichen
  core/role/shipped_test.go TestInstallReplacesTheV052CockpitBoard wiederholt TestInstallReplacesAnOlderShippedVersion fuer ein zweites Fixture; beide als Tabelle (key, Fixture, Pin-Hash) in einem Test
round: 2
---

# shipped.sha256 kennt jede Version, die in einem Release-Tag ausgeliefert wurde, und CI haelt das fest

## Definition of Done

- [x] Alle Datei-Versionen aus allen v*-Tags stehen in core/shipped/shipped.sha256 (Test oder Skript-Pruefung)
  proof: scripts/shipped-history.sh --check (exit 0 on 34 tags; was: 8 missing); core/shipped/shipped.sha256 regenerated, 0 lines lost
- [x] Ein CI-Schritt prueft gegen die Tags und schlaegt fehl, wenn eine ausgelieferte Version fehlt; laeuft auch im release-notes-Workflow vor dem Release-PR
  proof: .github/workflows/ci.yaml job shipped-manifest; .github/workflows/release-notes.yaml step 'Every version a release shipped is in the manifest' (fetch-depth: 0)
- [x] jaira update --roles auf einer unveraenderten v0.5.2-Installation von board.ts ueberschreibt sie statt Skipped (Test)
  proof: core/role/shipped_test.go TestInstallReplacesAnOlderShippedVersion/cockpit_board.ts_v0.5.2 (fails on the old manifest)
- [x] Eingabeklassen: Datei in einem Tag geloescht und spaeter wieder angelegt, CRLF-Version derselben Datei, Tag ohne die Datei, Merge-Konflikt im Manifest
  proof: core/shipped/history_test.go TestCheckHoldsTheManifestAgainstEveryTag (deleted+recreated, CRLF, conflict markers, CRLF manifest), TestCheckFailsWithoutTags, TestCheckFailsOnAnUnreadableTag; tag without the file = v0.2.0 in the fixture
- [x] .gitattributes setzt 'core/shipped/shipped.sha256 merge=union', sodass zwei PRs, die beide Zeilen anhaengen, nicht mehr in Konflikt geraten (heute bei jedem PR: #147, feat/tw5fe2)
  proof: .gitattributes core/shipped/shipped.sha256 merge=union; core/shipped/history_test.go TestManifestMergesWithoutConflict (fails without merge=union)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 18:18 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 18:18 · Alexander Sacharov** — Teamlead: erst starten, wenn PR #145 gemergt ist - der CI-Schritt aendert .github/workflows/ci.yaml, das #145 gerade umbaut.
- **2026-11-20 22:45 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock; PR #145 ist seit 20.11. 17:40 gemergt, Startbedingung erfuellt
- **2026-11-20 22:46 · Alexander Sacharov** — Offene Entscheidungen: keine, selbst entschieden: (1) Pruefung als 'scripts/shipped-history.sh --check' im selben Skript, damit Pfad-Mapping (lanes/, .claude/skills/jaira) nur einmal existiert. (2) Geprueft wird gegen die Tags, nicht gegen jede History-Version: Zwischenstaende auf Branches wurden nie ausgeliefert, und die Version im Arbeitsbaum haelt schon TestEveryInstalledFileIsAShippedVersion fest. (3) Keine v*-Tags im Klon = Fehler, nicht gruen: ein Check, der nichts prueft, darf nicht bestehen. (4) Blob-Bytes roh hashen (cat-file), keine CRLF-Umwandlung - das ist, was der Build einbettete. (5) Skript-Test als Go-Test mit Temp-Repo, auf Windows uebersprungen (CI-Schritt laeuft auf ubuntu). Stand vor dem Fix: Generator liefert 75 Zeilen, die im Manifest fehlen; kein Manifest-Eintrag fehlt im Generator (nichts geht durch Regenerieren verloren).
- **2026-11-20 22:53 · Alexander Sacharov** — Gefunden: die alten 8 fehlenden Zeilen standen unsortiert am Ende des Manifests (von Hand angehaengt) - Regenerieren sortiert sie ein, daher der grosse Diff ohne verlorene Zeile. Das Skript setzt jetzt LC_ALL=C: unter anderer Locale sortierte Manifeste wuerden bei merge=union Duplikate erzeugen. Bekannte Grenze: Pfade mit Leerzeichen kann das Manifest-Format nicht (Regex [^ ]+ hier wie in TestManifestLinesAreWellFormed); keine solche Datei existiert. Ein Duplikat nach einem union-Merge laesst TestManifestLinesAreWellFormed scheitern -> regenerieren.
- **2026-11-20 22:53 · Alexander Sacharov** — critique: (1) Formatpruefung im Skript doppelt TestManifestLinesAreWellFormed -> streichen. (2) zwei fast gleiche Install-Tests -> Tabellentest. Bewusst stehen gelassen: Pruefung gegen Tags statt gegen die ganze History (Note 'Offene Entscheidungen' Punkt 2).
