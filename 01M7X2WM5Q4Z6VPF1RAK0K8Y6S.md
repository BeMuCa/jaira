---
id: 01M7X2WM5Q4Z6VPF1RAK0K8Y6S
title: "shipped.sha256 kennt jede Version, die in einem Release-Tag ausgeliefert wurde, und CI haelt das fest"
status: backlog
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
updated-at: 2026-11-20T20:00:59Z
updated-by: Alexander Sacharov
---

# shipped.sha256 kennt jede Version, die in einem Release-Tag ausgeliefert wurde, und CI haelt das fest

## Definition of Done

- [ ] Alle Datei-Versionen aus allen v*-Tags stehen in core/shipped/shipped.sha256 (Test oder Skript-Pruefung)
- [ ] Ein CI-Schritt prueft gegen die Tags und schlaegt fehl, wenn eine ausgelieferte Version fehlt; laeuft auch im release-notes-Workflow vor dem Release-PR
- [ ] jaira update --roles auf einer unveraenderten v0.5.2-Installation von board.ts ueberschreibt sie statt Skipped (Test)
- [ ] Eingabeklassen: Datei in einem Tag geloescht und spaeter wieder angelegt, CRLF-Version derselben Datei, Tag ohne die Datei, Merge-Konflikt im Manifest
- [ ] .gitattributes setzt 'core/shipped/shipped.sha256 merge=union', sodass zwei PRs, die beide Zeilen anhaengen, nicht mehr in Konflikt geraten (heute bei jedem PR: #147, feat/tw5fe2)

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
