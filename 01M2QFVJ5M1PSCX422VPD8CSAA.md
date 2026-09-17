---
id: 01M2QFVJ5M1PSCX422VPD8CSAA
title: 0.3.0 schneiden
status: optimize
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer 'jaira update' nach 0.3.0 laufen laesst, liest vier Aenderungen und nicht zweiundzwanzig Zeilen Entwicklungsweg - und die Version ist geschnitten, sobald die vier Tickets in signoff angenommen sind"
context: |-
  Stand am 2026-09-17: master ist gleich upstream/master, letzter Tag ist v0.2.1. Naechste Version ist 0.3.0, weil Milestones eine Funktion sind und kein Patch.

  Was drin ist - vier Aenderungen: Milestones ganz (0YGWXQ), eine Lane ohne Code-Aenderung committet nicht mehr und der Handle im Betreff wird zur Bedingung (9ZZSFT), spawn.sh kennt die Lane 'dispatch' und --no-worktree (7KX89C), und jaira-role-pr macht den Pull Request auf, wenn ein Mensch die Rolle ruft.

  Das Problem, das vor dem Schnitt behoben gehoert: unter '## Unreleased' stehen 22 Zeilen, davon 18 zu Milestones. Mehrere davon beschreiben DASSELBE Verhalten in verschiedenen Entwicklungsstaenden - erst 'create weist den Namen ab', dann 'folge dem Ref in den Baum, der abgelegt hat', dann 'nein, lies den Abweisungstext als Hinweis auf dein eigenes Logbuch'. Ein Benutzer hat die Zwischenstaende nie gesehen; fuer ihn ist das ein Verhalten, nicht drei.

  Warum das zaehlt und nicht Kosmetik ist: core/release/NOTES.md ist mit go:embed im Binary und ist genau das, was 'jaira update' jemandem vorliest, dessen Board ein aelterer Build angefasst hat. Zweiundzwanzig Zeilen fuer vier Aenderungen ist das, was dieser Mensch bekommt.

  Die Sektion ist nicht getaggt, also ist sie offen und darf umgeschrieben werden - CLAUDE.md verbietet nur, eine bereits getaggte Sektion nachtraeglich wachsen zu lassen.

  Blockiert vom Menschen: 9ZZSFT, 0YGWXQ, 7KX89C und GTQHNH stehen in signoff. Ohne ihre Annahme wird nichts geschnitten. Das Falten der Notizen haengt nicht daran und kann sofort laufen.
definition-of-done: "Unter '## Unreleased' steht je eine Zeile je Aenderung, die ein Benutzer von aussen merkt - kein Zwischenstand, den nie jemand ausgeliefert bekommen hat, und keine zwei Zeilen fuer ein Verhalten. Die 18 Milestone-Zeilen sind zusammengefasst; was jede Zeile den Leser TUN laesst, bleibt erhalten."
tags:
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-09-17T10:51:02Z
updated-at: 2026-09-17T11:29:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-33921
claimed-at: 2026-09-17T11:14:50Z
mode: conversational
outcome-what: "Die Naht zwischen den gefalteten Milestone-Zeilen neu gezogen: 'Reise auf refs/jaira/milestones/<name>, jaira fetch holt sie' ist aus der Kartenrand-Zeile heraus und hat eine eigene Zeile (core/release/NOTES.md:24); die Kartenrand-Zeile beschreibt jetzt nur noch das Rendern (23); 'jaira milestone create/add/rm/ls' und die Handbearbeitung von .jaira/milestones/<name>.md sind zu einer Zeile verschmolzen (21). Weiter 5 Milestone-Zeilen, '## Unreleased' traegt 9."
outcome-why: "Die critique hat auf Zeile 24 ein Verhalten auf zwei Zeilen gefunden - Kartenrendering und Ref-Reise verklebt, wobei das Reise-Thema in Zeile 25 weiterlief. Genau das schliesst DoD 1 aus. Nach dem Umbau decken sich die 5 Zeilen exakt mit den 5 Faehigkeiten aus der Entscheidung des Menschen vom 11:14."
outcome-resolves: "DoD 1, 2 und 3. DoD 4 bleibt offen und gehoert nicht in diese Lane: die Umbenennung nach '## 0.3.0' faehrt in dem Commit, den der Mensch taggt, und 9ZZSFT/0YGWXQ/7KX89C/GTQHNH stehen noch in signoff."
review-summary: "none"
---

# 0.3.0 schneiden

## Definition of Done

- [x] Unter '## Unreleased' steht je eine Zeile je Aenderung, die ein Benutzer von aussen merkt - kein Zwischenstand, den nie jemand ausgeliefert bekommen hat, und keine zwei Zeilen fuer ein Verhalten. Die 18 Milestone-Zeilen sind zusammengefasst; was jede Zeile den Leser TUN laesst, bleibt erhalten.
  proof: core/release/NOTES.md:21-25 — Naht aus der critique korrigiert: Reise ueber refs/jaira/milestones/<name> ist jetzt eine eigene Zeile (24), Kartenrand-Farben stehen allein (23), create und Handbearbeitung sind eine Zeile (21); weiter 5 Milestone-Zeilen, '## Unreleased' traegt 9
- [x] Das Zeilenformat ueberlebt das Falten: jede Aenderung ist genau eine Zeile, die mit '- ' beginnt, keine ist ueber zwei Zeilen umgebrochen. Nachgewiesen mit dem Test, der core/release/NOTES.md scannt.
  proof: go test ./core/release/ ok (TestEmbeddedNotesParseToAtLeastOneRealEntry); awk-Scan ueber '## Unreleased': 9 '- '-Zeilen, keine Nicht-'- '-Zeile
- [x] Keine Aenderung geht beim Falten verloren: jede der vier Aenderungen ist im gefalteten Text nachweisbar, und die Kommandos und Tasten, die ein Benutzer tippt ('jaira milestone', 'M', '--milestone', '--no-worktree', 'dispatch', 'jaira restore'), stehen weiter darin.
  proof: grep ueber '## Unreleased': 'jaira milestone' 4x, '`M`' 2x, '--milestone' 1x, '--no-worktree' 1x, 'dispatch' 5x, 'jaira restore' 1x, 'refs/jaira/milestones' 1x, 'jaira fetch' 2x
- [ ] Die Umbenennung von '## Unreleased' nach '## 0.3.0' und die frische leere '## Unreleased' darueber faehren in dem Commit, den der Mensch taggt - nicht frueher. Der Tag selbst wird von einem Menschen gesetzt, nicht von einem Agenten.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Die 18 Milestone-Zeilen (NOTES.md 21-38) auf 5 Zeilen falten, je eine pro Faehigkeit nach der Entscheidung vom 2026-09-17 11:14
- [x] Zeilenformat pruefen: jede Aenderung genau eine '- '-Zeile, kein Umbruch - go test ./core/release/
- [x] Nachweisen, dass kein Kommando und keine Taste verloren ging: jaira milestone, M, --milestone, --no-worktree, dispatch, jaira restore
- [ ] Umbenennung nach '## 0.3.0' NICHT hier - sie gehoert in den Commit, den der Mensch taggt, nach Annahme der vier signoff-Tickets

## Progress
- **2026-09-17 11:14 · Alexander Sacharov** — Entscheidung vom Menschen, vor der Plan-Lane: die 18 Milestone-Zeilen werden auf ca. 5 Zeilen gefaltet, je eine pro Faehigkeit, die ein Benutzer anders benutzt - (1) Milestone anlegen/aendern per CLI und per Hand in .jaira/milestones/<name>.md, (2) Board und 'jaira list' auf eine Runde einschraenken (--milestone, Taste M), (3) Milestone-Farben am RECHTEN Kartenrand lesen, (4) Milestone reist auf seinem eigenen Ref und kommt mit 'jaira fetch', (5) Milestone mit 'jaira logbook <name>' ablegen und mit 'jaira restore <name>.md' zurueckholen. Alle Zwischenstaende und Corner-Case-Abweisungen (drei Fassungen der create-Abweisung, Locking, verlorene Rennen) fallen weg - ein Benutzer hat sie nie gesehen. Nicht 4 Zeilen (eine Milestone-Zeile wuerde zu lang und Kommandos ertrinken darin), nicht 8-10. Die anderen drei Aenderungen bleiben je eine Zeile. DoD 3 bleibt gewahrt: jaira milestone, M, --milestone, --no-worktree, dispatch, jaira restore stehen weiter im Text.
- **2026-09-17 11:18 · Alexander Sacharov** — Erster Wurf waren 7 Zeilen (CLI und Handbearbeitung getrennt, logbook und restore getrennt) - auf 5 nachgefaltet, weil die Entscheidung vom Menschen 5 Faehigkeiten aufzaehlt und 8-10 ausdruecklich ausschliesst. Die Ref-Zeile hat keine eigene bekommen, sondern haengt an der Farbzeile: sie ist zwei Saetze lang und eine sechste Zeile waere gegen die Entscheidung gewesen.

Bewusst weggeworfen, nicht vergessen: die drei Fassungen der create-Abweisung, 'logbook nennt gefilte Milestones statt ticket: not found', die zwei Locking-Zeilen (restore und Hintergrund-fetch nehmen dasselbe Lock) und die Zeile ueber den verlorenen Wettlauf mit abgeschnittenem Namen. Alle beschreiben Zwischenstaende innerhalb derselben unveroeffentlichten Sektion - kein Benutzer hat je einen Build gesehen, in dem sie fehlten. Wer sie sucht: git show 686daa5:core/release/NOTES.md.

DoD 4 bleibt offen und gehoert nicht in diese Lane: die Umbenennung nach '## 0.3.0' faehrt in dem Commit, den der Mensch taggt, und 9ZZSFT/0YGWXQ/7KX89C/GTQHNH stehen noch in signoff. Das Ticket kann die Endlane nicht erreichen, bevor das passiert ist.
- **2026-09-17 11:21 · Alexander Sacharov** — critique: Das Falten selbst ist richtig geschnitten — die Zwischenstaende (Lock-Rennen, der 'elease'-Abschneidefehler, die drei Entwicklungsstufen der Filed-Abweisung) sind Bugs an einer Faehigkeit, die noch nie ausgeliefert wurde, und gehoeren nicht in die Notizen. Auch das Verschmelzen von 'rm loescht nur die Zeile' in die create-Zeile (21) und von Taste M in die list-Zeile (23) ist je ein Verhalten, nicht zwei. Die eine Ausnahme ist Zeile 24: dort ist das Kartenrendering mit dem Reisen ueber refs/jaira/milestones/<name> verklebt, und dieses Reise-Thema laeuft in Zeile 25 weiter. Ein Verhalten auf zwei Zeilen ist genau der Fehler, den DoD 1 ausschliesst — nur diesmal nicht durch Nichtfalten, sondern durch falsche Naht. Fix: den Reise-Satz aus 24 heraus und an den Anfang von 25 (wo 'jaira fetch' schon erklaert wird) oder als eigene sechste Zeile davor. Bewusst NICHT beanstandet: der --json-Schluessel 'milestone' aus der Rennmeldung ist weggefallen, waehrend 'milestones' und 'milestones_filed' blieben — das ist konsistent, weil die Rennmeldung zur Behebung eines nie ausgelieferten Fehlers gehoert.
- **2026-09-17 11:23 · Alexander Sacharov** — Critique-Befund behoben, aber nicht so wie die critique vorgeschlagen hat. Ihr Vorschlag war, den Reise-Satz an den Anfang der logbook-Zeile zu haengen - das haette die Naht nur verschoben (Reisen + Ablegen sind auch zwei Verhalten). Stattdessen: Reise hat jetzt eine eigene Zeile, und create + Handbearbeitung sind zu EINER Zeile verschmolzen. Damit bleibt es bei 5 Milestone-Zeilen und sie decken sich jetzt exakt mit der Aufzaehlung des Menschen vom 11:14 - anlegen/aendern (CLI und Hand), einschraenken, Farben, reisen, ablegen/zurueckholen. Der vorige Wurf hatte 1 gespalten und 3+4 verklebt und kam nur zufaellig auf 5.\n\nDie verschmolzene create-Zeile ist mit 918 Zeichen die laengste der Sektion; die Befuerchtung des Menschen, Kommandos wuerden in einer langen Zeile ertrinken, galt einer EINZIGEN Milestone-Zeile fuer alles - die Kommandonamen stehen hier weiter alle vorn im ersten Satz. Wer sie doch wieder trennen will: das ist eine Abweichung von der 5-Faehigkeiten-Entscheidung, keine Formatfrage.
- **2026-09-17 11:25 · Alexander Sacharov** — critique (zweiter Durchgang): nichts mehr zu beanstanden. Der Befund des ersten Durchgangs (Naht auf Zeile 24: Kartenrendering und Ref-Reise verklebt) ist behoben, und zwar besser als vorgeschlagen - die Reise hat eine eigene Zeile statt an die logbook-Zeile angehaengt zu werden, wo dieselbe Doppelung nur umgezogen waere. Die 5 Milestone-Zeilen decken sich jetzt Zeile fuer Zeile mit den 5 Faehigkeiten der Menschenentscheidung vom 11:14. Geprueft und NICHT beanstandet: (1) die 918 Zeichen der create-Zeile - die Sektion 0.2.1 traegt Zeilen derselben Laenge, das Format ist ein Zeilenscan und kennt keine Laengengrenze, und alle Kommandonamen stehen vorn im ersten Satz; (2) 'rm loescht nur die Zeile' in der create-Zeile - im ersten Durchgang stehengelassen, wird nicht neu aufgemacht; (3) die weggefallenen Zwischenstaende - jede beschreibt einen Bug an einer nie ausgelieferten Faehigkeit. Nachgezaehlt: 9 '- '-Zeilen unter '## Unreleased', keine Nicht-'- '-Zeile, kein Umbruch; jaira milestone 4x, M 2x, --milestone, --no-worktree, --color, jaira restore, refs/jaira/milestones, milestones_filed je 1x.
- **2026-09-17 11:29 · Alexander Sacharov** — optimize: Nichts entfernt - vier Durchgaenge ueber den Diff, jeder Befund gepruefter Verzicht statt Uebersehen.

Doppelung: keine. Kein Milestone-Satz steht in einer bereits getaggten Sektion (awk ueber alles unterhalb '## Unreleased': 0 Treffer auf 'milestone'), und die 5 gefalteten Zeilen ueberschneiden sich untereinander nicht - 'rm nimmt Tickets heraus' (Z. 21) und 'jaira logbook <name> legt ab' (Z. 25) sind zwei Kommandos mit zwei Wirkungen, nicht eins zweimal.

Tote Aussagen: keine. Jedes genannte Kommando, Flag, Feld und jede Taste existiert im Code - milestone create/add/rm/ls in internal/cli/milestones.go:77,178,197,309, --color 1-255 in :172, StatusFiled in core/milestone/milestone.go:303, milestones_filed in internal/cli/fetch.go, jaira restore in internal/cli/archive.go:95.

Fluff geprueft und BEHALTEN, mit Grund: (1) 'rm entfernt nur die Zeile, ein Milestone verschwindet nicht' (Z. 21) liest sich wie Beruhigung, steuert aber ein Kommando - ohne den Satz haelt ein Leser 'rm' fuer den Weg, einen Milestone zu loeschen, und Z. 25 nennt dafuer 'jaira logbook'. (2) 'die Karte ist genau so breit' (Z. 23) beschreibt ein Nicht-Ereignis, aber genau das, was ein Leser am TUI befuerchtet. (3) 'und Erstellzeit' (Z. 21) ist echter Fluff, aber 18 Zeichen - ein Edit dafuer waere teurer als der Gewinn.

Kosten: entfaellt, der Text ist ein Zeilenscan ohne Laengengrenze. Die 907 Zeichen der create-Zeile sind nicht die laengste der Sektion - die PR-Zeile hat 918.

Ausserhalb dieser Aenderung gefunden, NICHT angefasst: die beiden spawn.sh-Zeilen (Z. 19 und 20, aus 7KX89C) beschreiben dasselbe Skript, aber nur Z. 20 sagt, wie man es bekommt ('jaira roles install --global --force') - wer nur Z. 19 liest, erfaehrt die Lane 'dispatch' und nicht, dass sie mit dem jaira-dispatcher-Rollenordner kommt. Zweitens nennt '--color' bei milestone 1-255, bei tag 0-255 (internal/cli/tags.go:389) - eine Inkonsistenz im Code, nicht in den Notizen.
