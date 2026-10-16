---
id: 01M522QPD7TPY1GQDXR1EDSQZP
title: "Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'"
status: review
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
updated-at: 2026-10-16T11:37:07Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-482020
claimed-at: 2026-10-16T11:16:45Z
outcome-what: "jaira-cockpit-Band: neue Zeile 'jaira X is released (this is Y)' mit [ Update ] und Dismiss, sobald der Release-Cache ($JAIRA_HOME|~/.jaira/update-check.json) eine neuere Version als 'jaira --version' kennt; jede Minute neu gelesen. Update: 'jaira self upgrade --check --json' und 'jaira update --roles --lanes --dry-run', Ergebnis im Pane cockpit-upgrade (Cancel hat den Fokus); erst 'Update to X' dort fuehrt 'jaira self upgrade' und dann 'jaira update --roles --lanes' aus und zeigt die Ausgabe. Dismiss legt die Version im Store (dismissed-release) ab: diese Version nie wieder, die naechste schon. Helfer offeredRelease/upgradeCheckOf in board.ts, Woerter in texts.ts, State release/upgrade im Vertrag, 12 Tests, README-Absatz, Release-Notiz, shipped.sha256 neu."
outcome-why: "Die stderr-Zeile nach next/resume/claim/move erreicht nur Sitzungen, die gerade jaira-Befehle fahren; das Band erreicht jede Sitzung mit offenem Cockpit. Berk wollte den Knopf, mit Vorschau vor jedem Schreiben und einem Verwerfen, das sich die Version merkt."
outcome-resolves: "DoD 1: Hinweis mit Update und Verwerfen (Dismiss) bei neuerer Version, Update = Vorschau dann Ja-Pfad, Verwerfen nur fuer genau diese Version; getestet mit claude plugin test (155 pass) und einem node-Lauf der Helfer auf echter jaira-Ausgabe."
review-summary: none
review-gaps: "nothing removed: the one duplicate (JAIRA_HOME-or-~/.jaira in settingsJson) is already folded into jairaHome(); every new text key and helper has a caller; considered one helper for 'stdout+stderr without warnings' in previewUpgrade/upgradeJaira (register.tsx:366,385) - two callers with different fallbacks, left; pre-existing and left alone: metrics.csv is read and written under HOME/.jaira, not JAIRA_HOME (register.tsx:913,1429)"
test-verdict: "pass: go build ./... ok, go test -race -count=1 ./... RC=0 (32 ok, 0 FAIL); claude plugin validate passed, claude plugin test 155 pass 0 fail (12 in 'a newer jaira release'); DoD 1 proof lines checked in the tree; yes-path run with real binaries built as 0.5.1/0.6.0 against a local stub release in a scratch HOME/JAIRA_HOME: check + dry-run wrote nothing (snapshot equal), self upgrade replaced only the scratch binary, update --roles --lanes ran; ~/.local/bin/jaira unchanged; not seen: the band live in a Claude Code terminal"
---

# Cockpit-Band: Update-Hinweis mit Knopf 'Update' und Knopf 'Verwerfen'

## Definition of Done

- [x] Band zeigt bei neuerer Version Hinweis mit Knopf Update und Knopf Verwerfen; Update fuehrt Vorschau und danach Ja-Pfad aus; Verwerfen unterdrueckt den Hinweis fuer genau diese Version; Test mit Node-Lauf im Cockpit
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1177 band row Update/Dismiss, :356 previewUpgrade (self upgrade --check --json + update --roles --lanes --dry-run, pane), :375 upgradeJaira (self upgrade + update --roles --lanes, only from the pane's yes), :346 dismissRelease (store dismissed-release), board.ts:520 offeredRelease; tests/board.test.ts:1559 describe 'a newer jaira release' (12 tests) - claude plugin test: 155 pass 0 fail; node 24 run of board.ts on real jaira output in scratch

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] board.ts: releaseOf(cache json) und offeredRelease(installed, latest, dismissed); texts.ts: Woerter fuer Hinweis, Vorschau, Ergebnis
- [x] register.tsx: Cache jede Minute lesen (jairaHome geteilt mit settingsJson), Hinweiszeile mit Update und Dismiss im Band, Dismiss merkt Version im Store
- [x] register.tsx: Update = self upgrade --check --json + update --roles --lanes --dry-run, Pane mit Vorschau; Bestaetigen = self upgrade + update --roles --lanes; types/index.d.ts
- [x] tests/board.test.ts: Hinweis, Update mit Vorschau, nichts vor Bestaetigung, Ja-Pfad, Dismiss fuer genau diese Version; claude plugin test + validate
- [x] release note core/release/notes/EDSQZP.md, shipped.sha256 neu, Ja-Pfad mit echtem Binary gegen Stub-Release in Scratch-JAIRA_HOME pruefen

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
- **2026-10-16 11:34 · BeMuCa** — critique (erster Pass, ganzer Diff 9d2cfd51): kein Befund. Geprueft und stehen gelassen: (1) eigener Pane cockpit-upgrade statt $.ui.ask wie beim Cache-Guard - ask ist der AskUserQuestion-Dialog fuer eine kurze Frage mit 'Other'-Freitext; die Vorschau ist mehrzeilig und folgt dem Muster des Explanation-Panes (Wahl in einem Pane mit Fokus). (2) run+output statt des jaira()-Helfers, weil die Vorschau stdout braucht und jaira() nur ok/why liefert. (3) 'update --roles --lanes' nach 'self upgrade' wiederholt --roles: noetig, weil ein Upgrade von 0.5.1 vom alten Binary gemacht wird, das update --roles noch nicht startet (Release-Notiz PWAF22). (4) Dismiss im Store statt in ~/.jaira: Entscheidung 3 oben.
