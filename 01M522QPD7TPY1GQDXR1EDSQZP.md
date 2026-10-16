---
id: 01M522QPD7TPY1GQDXR1EDSQZP
title: "Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'"
status: in-progress
ready: true
creator: BeMuCa
goal: "Wenn eine neue jaira-Version da ist, zeigt das Cockpit-Band einen Hinweis mit zwei Knoepfen: Update (fuehrt den Ja-Pfad aus) und Verwerfen (blendet den Hinweis fuer diese Version aus)"
context: "Berk wollte: im Cockpit einen Update-Hinweis mit zwei Knoepfen. Das Band zeigt heute keinen Update-Hinweis. T8A8KM hat ihn verworfen, weil das Cockpit nur mit installiertem Cockpit laeuft und das TS nicht getestet war. Das ist kein Grund gegen den Knopf, nur einer fuer den Weg. Die stderr-Zeile in next/resume/claim/move erreicht nur Sitzungen, die gerade einen jaira-Befehl fahren. Der Knopf erreicht jede Sitzung, die das Cockpit offen hat. Ja-Pfad: erst 'jaira self upgrade --check', dann eine Vorschau 'jaira update --roles --lanes --dry-run', erst nach Ja das echte Update. Verwerfen merkt sich die Version, damit die Zeile nicht wiederkommt."
definition-of-done: Band zeigt bei neuerer Version Hinweis mit Knopf Update und Knopf Verwerfen; Update fuehrt Vorschau und danach Ja-Pfad aus; Verwerfen unterdrueckt den Hinweis fuer genau diese Version; Test mit Node-Lauf im Cockpit
tags:
  - release
  - tui
blocked-by: []
related: []
follows: 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-10-16T10:04:46Z
updated-at: 2026-10-16T11:32:37Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-482020
claimed-at: 2026-10-16T11:16:45Z
---

# Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'

## Definition of Done

- [ ] Band zeigt bei neuerer Version Hinweis mit Knopf Update und Knopf Verwerfen; Update fuehrt Vorschau und danach Ja-Pfad aus; Verwerfen unterdrueckt den Hinweis fuer genau diese Version; Test mit Node-Lauf im Cockpit

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] board.ts: releaseOf(cache json) und offeredRelease(installed, latest, dismissed); texts.ts: Woerter fuer Hinweis, Vorschau, Ergebnis
- [x] register.tsx: Cache jede Minute lesen (jairaHome geteilt mit settingsJson), Hinweiszeile mit Update und Dismiss im Band, Dismiss merkt Version im Store
- [x] register.tsx: Update = self upgrade --check --json + update --roles --lanes --dry-run, Pane mit Vorschau; Bestaetigen = self upgrade + update --roles --lanes; types/index.d.ts
- [x] tests/board.test.ts: Hinweis, Update mit Vorschau, nichts vor Bestaetigung, Ja-Pfad, Dismiss fuer genau diese Version; claude plugin test + validate
- [ ] release note core/release/notes/EDSQZP.md, shipped.sha256 neu, Ja-Pfad mit echtem Binary gegen Stub-Release in Scratch-JAIRA_HOME pruefen

## Progress
- **2026-10-16 11:20 · BeMuCa** — Entscheidungen vor der Arbeit (keine offen fuer Berk, alle aus Code/Ticket/Default gesetzt):
1. Versionsdaten: Cockpit liest dieselbe Cache-Datei wie TUI und stderr-Zeile ($JAIRA_HOME bzw. ~/.jaira/update-check.json, Feld latest) und vergleicht mit 'jaira --version'; kein zweiter Check. Aufgefrischt wird der Cache weiter von TUI und next/resume/claim/move (PollCache, stuendlich); das Cockpit liest ihn jede Minute, wie das TUI.
2. Knopf-Texte englisch: 'Update' und 'Dismiss' statt 'Verwerfen' - texts.ts: jeder Cockpit-Text ist englisch, egal welche Sprache die Person spricht.
3. Verwerfen merkt sich die Version im Store des Cockpits ($.store, ueber Sitzungen). Gilt nur fuer das Band; die stderr-Zeile der CLI bleibt. Eine neuere Version zeigt den Hinweis wieder.
4. Update: 'jaira self upgrade --check --json', dann 'jaira update --roles --lanes --dry-run', beides im Board-Verzeichnis; Vorschau in einem Pane mit Knopf zum Bestaetigen und Abbrechen. Erst die Bestaetigung startet 'jaira self upgrade' und danach 'jaira update --roles --lanes'. --check schreibt nur den Cache, sonst nichts.
5. Ohne Board und in dev-Builds kein Hinweis (wie die TUI-Pille und die Regel des Cockpits: was ein Board braucht, bleibt ohne Board unsichtbar).
6. Mod-Dateien geaendert -> core/shipped/shipped.sha256 mit scripts/shipped-history.sh neu erzeugen (Test erzwingt es).
- **2026-10-16 11:32 · BeMuCa** — Gepruefte Wege (fuer critique/testing/review):
- claude plugin test im Mod-Ordner: 155 pass, 0 fail; 7 der 12 neuen Tests fallen mit dem alten register.tsx (die 5 'kein Hinweis'-Faelle sind Gegenproben und bestehen ohne Feature).
- Ja-Pfad mit echten Binaries in Scratch (HOME, JAIRA_HOME, PATH auf scratchpad/yes): Build dieses Branches als 0.5.1 und 0.6.0 (-X main.version, nicht release.Current - das setzt cli.Execute), Stub-Release per python http.server auf 127.0.0.1 ueber JAIRA_RELEASE_API/JAIRA_RELEASE_DOWNLOADS. jaira next schrieb den Cache (PollCache), --check JSON wie vom Cockpit gelesen, dry-run aenderte nichts (Snapshot gleich), self upgrade ersetzte nur das Scratch-Binary und zog die alte register.tsx in beiden Skills-Ordnern nach. ~/.local/bin/jaira vorher/nachher gleicher sha256.
- node 24 laedt hooks/board.ts direkt (Type-Stripping): offeredRelease/upgradeCheckOf auf der echten Ausgabe geprueft.
- tsc (npx typescript@5.6.3, tsconfig ausserhalb des Mod-Ordners): hooks/ und types/ fehlerfrei; tests/board.test.ts hat 5 alte Fehler (doppeltes toasts im World-Typ Z.42/69/142, Mounted.input Z.780/782) - nicht von hier, nicht angefasst.
Bekannte Grenze: die Vorschau ist der dry-run des HEUTE installierten jaira; was erst das neue Binary mitbringt, steht nicht darin (self upgrade laesst danach ohnehin das neue 'update --roles' laufen). Ein Upgrade im Terminal waehrend der Sitzung: das Cockpit merkt sich die Version je Sitzung, der Hinweis bleibt, bis Update gedrueckt wird - dann sagt --check 'already the latest' und der Hinweis geht.
