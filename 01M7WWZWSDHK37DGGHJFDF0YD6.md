---
id: 01M7WWZWSDHK37DGGHJFDF0YD6
title: "Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket traegt ein Blocker-Kriterium und die Option stop-after-fix; testing, bughunt und review lesen sie, und das Cockpit bietet Explain nur bei einem blocker an."
context: |-
  Heute: ein Teamlead schreibt mehrmals am Tag von Hand "fix -> test -> kein Bughunt -> review, Blocker nur bei falschem Ergebnis oder Absturz". Das hat jedes Mal sofort konvergiert, steht aber nur in Notizen und gilt fuer das naechste Ticket nicht.
  Ohne Schwelle ist jeder Randfall ein Blocker; jeder Fix ist neuer Code, den der naechste Bughunt jagt (das Ziel wandert mit).
  review-gaps hat immer 2-6 note:-Zeilen; das Cockpit zeigt Explain schon bei note: allein (core/mod/builtin/jaira-cockpit/hooks/board.ts:196 explainable/namesGap), deshalb hat jedes Ticket in signoff ein Explain.
  Die Schwelle gehoert in DoD/Plan, gesetzt von pre-process (siehe Ticket zu Randfaellen vor dem Code); der Mensch kann sie aendern.
  Verwandt: H6CD33 (max-rounds einer Schleifen-Lane).
definition-of-done: "Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)"
tags:
  - gates
blocked-by:
  - 01M7WWZRD7WNC5NT1PS5TW5FE2
related:
  - 01M6M7EA5EG5TWJSTVNQH6CD33
  - 01M5ZJ4KA5N2PH7Q4WY4F1V59P
  - 01M7WWZRD7WNC5NT1PS5TW5FE2
commits: []
created-at: 2026-11-20T16:35:17Z
updated-at: 2026-11-20T22:55:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91898
claimed-at: 2026-11-20T22:44:37Z
---

# Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt

## Definition of Done

- [ ] Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)
- [ ] Mit stop-after-fix laeuft nach einem Fix-Durchgang nur der Testlauf, kein neuer Bughunt
- [ ] Das Cockpit zeigt Explain nur, wenn review-gaps eine blocker:-Zeile hat oder das Urteil einen Vorbehalt nennt; reine note:-Zeilen stehen als Liste ohne Knopf (Test in board.test.ts)
- [ ] jaira zeichnet bei jedem Rueckweg die Herkunfts-Lane auf und ob eine Person oder eine Schleifen-Lane zurueckgeschickt hat; jaira show --json gibt beides aus, und die Schleifenzaehlung zaehlt nur Rueckgaben derselben Schleifen-Lane (Test: signoff->in-progress zaehlt nicht)
- [ ] Die Explain-Vorlage (core/mod/builtin/jaira-cockpit/hooks/texts.ts:76-80) ordnet jede Note zuerst ein (Defekt / Rand in diesem Ticket / vorbestehend / nur Hinweis), sagt ausdruecklich, wenn nichts davon die Abnahme verlangt, und gibt keine feste Reihenfolge 1) [give back] 2) [accept] mehr vor - heute zwingt die Vorlage 'Warum: why it matters' und give back als erste Option, das liest sich immer wie 'muss geaendert werden'
- [ ] Steht ein Ticket in signoff/human und ist sein review-verdict aus einer frueheren Runde (review_round < round), schweigt das Cockpit nicht: das Band sagt 'Review von Runde N, seitdem geaendert' und zeigt die alten Funde mit Explain; heute blendet sameRound() in board.ts:189 Explain dann aus (Test)
- [ ] Die Explain-Ansicht im Cockpit ist farbig und auf einen Blick lesbar: Statuszeile mit gruenem Haken/rotem Kreuz je Pruefung (DoD, Tests, Blocker, bughunt auf HEAD); jede Note mit farbiger Marke nach Einordnung - rot 'Blocker', gelb 'kleiner Rand', grau/gedimmt 'bekannte Grenze', blau 'nicht aus diesem Ticket'; die Empfehlung fett mit eigener Farbe; die Notiz-Texte normal, nur die Marken farbig; funktioniert im hellen und dunklen Terminal-Thema (Theme-Farben statt fester Hex-Werte, wie das Band heute)
- [ ] class: the loop brake reads recorded facts, never prose — every return is recorded by jaira move, whatever path moved it (CLI, cockpit, --force); parking in blocked and back records nothing; a ticket without back-* fields reads as no return (table test)
- [ ] input classes: blocker-rule with quotes, umlauts, multiple lines and CRLF survives set/show round-trip; stop-after-fix rejects anything but true/false/empty (exit 2); two parallel moves on one ticket leave consistent back-* fields (one refused or both counted)
- [ ] acceptance rule for the explain parser: an explain note line without a known class mark or status shape is shown as plain text, never dropped - an old note must still read as before

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] schema: fields blocker-rule, stop-after-fix, back-from, back-by, back-count in core/ticket/schema.go; jaira set checks stop-after-fix is true/false/empty like mode
- [x] move: on every send-back write back-from, back-by (person when the from-lane is not agentic or needs a human exit/question, else lane) and back-count (consecutive returns from the same lane since a person's or another lane's return); table test in core/move
- [x] show --json prints the five fields; --for-lane payload carries blocker_rule (ticket's, else the default four classes) and stop_after_fix beside mode, also in the plain-text header; test in internal/cli
- [ ] lane prompts + twins: testing skips the hunt in a fix round when stop_after_fix, reads blocker_rule; review reads blocker_rule; bughunt skill marks a change finding that meets the rule; in-progress third-return reads back-from/back-by/back-count instead of counting notes; pre-process may set the two fields; shipped.sha256
- [ ] cockpit board.ts: Explain only for a blocker: line or a reservation; note: lines listed under the row without a button; a review of an older round shows 'Review von Runde N, seitdem geaendert' with Explain; tests in board.test.ts
- [ ] explain template in texts.ts and spawn.sh: status line, one line per note with [class], one Empfehlung, no Warum:, no fixed give-back-first order, neutral case says it can be accepted
- [ ] explain pane: parse status line and class marks, colour them with theme colours (error/warning/dim/suggestion), recommendation bold; old notes without marks render as before; tests
- [ ] release note core/release/notes/DF0YD6.md; go test ./... and claude plugin test in the cockpit dir, exit codes checked

## Progress
- **2026-11-20 16:48 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 17:05 · Alexander Sacharov** — Aus TW5FE2 (bughunt e12f36d5): der Prompt-Zaehler in in-progress kann eine Rueckgabe der Person aus signoff nicht von einer review-Rueckgabe unterscheiden - jaira speichert die Herkunfts-Lane nicht (show --json unterscheidet sich nur in status/next_lane/round; sendsBack in core/move/move.go:133 erhoeht round bei jedem Rueckweg, updated_by ist derselbe git-Nutzer). Folge: nach zwei echten review-Rueckgaben eskaliert in-progress, statt zu tun, was die Person wollte. Die echte Zaehlung gehoert ins Binary - neuer DoD-Punkt hier.
- **2026-11-20 17:45 · Alexander Sacharov** — Alex 20.11.2026: 'aus dem Text der Beschreibung im Cockpit habe ich immer das Gefuehl, dass man es korrigieren muss.' Ursache in texts.ts:78-80: Explain gibt es nur bei einem Fund, die Vorlage verlangt 'Warum: why it matters' und setzt [give back] als Option 1. Alex schickt deshalb (und weil er keine offenen Raender will) jedes Ticket zurueck.
- **2026-11-20 17:57 · Alexander Sacharov** — Beobachtet 20.11.2026 in einem anderen Projekt: review Runde 7 schickt mit Blocker zurueck, Stopp-Regel nach drei Rueckgaben, Teamlead entscheidet 'Befund richtig, Verhalten so gewollt' und schiebt das Ticket ohne neues Review nach signoff. round=8, review_round=7 -> sameRound() false -> kein Explain, obwohl review-verdict 'Zurueck an in-progress: ein Blocker' und eine blocker:-Zeile auf dem Ticket stehen. Die Person sieht nichts. Verwandt: F1V59P (Verdikt 'zurueck' sperrt signoff).
- **2026-11-20 18:00 · Alexander Sacharov** — Alex 20.11.2026 zur Explain-Form (DoD 5): 'so wie du es hier erklaerst, ist es sofort verstaendlich; die Beschreibung im Cockpit klingt wie eine Katastrophe.' Vorbild ist die Form, die der Teamlead im Chat benutzt:
- eine Kopfzeile mit dem Stand: Lane, DoD x/y, Tests, 'kein Blocker' / 'Blocker: ...', bughunt auf HEAD ja/nein
- eine Tabelle, eine Zeile je Note: | # | was (ein Satz, Alltagssprache) | Einordnung |; Einordnung ist eine von: bekannte Grenze (nicht loesbar, dokumentiert), kleiner Rand (formal ok, eine Runde ~10 min), nicht aus diesem Ticket (wo es erledigt wird), Blocker
- dann genau eine Empfehlung mit Grund in einem Satz, und welche Note sie betrifft
- kein 'Warum: why it matters' pro Note, keine feste Reihenfolge give back zuerst; der neutrale Fall (nur bekannte Grenzen / fremde Notes) sagt ausdruecklich 'kann so angenommen werden'
- **2026-11-20 22:44 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 22:45 · Alexander Sacharov** — decision: feat/df0yd6 builds on origin/feat/tw5fe2 (merged in), because TW5FE2 (blocked_by) is done but not on master yet and changes the same lane prompts (testing, review, bughunt). Once TW5FE2's PR lands, its commits drop out of this branch's diff.
- **2026-11-20 22:47 · Alexander Sacharov** — decision (worker, settled from code): 1) two frontmatter fields: 'blocker-rule' (one sentence, free text) and 'stop-after-fix' (true/false). Empty blocker-rule = the four classes TW5FE2 already put into testing/review prompts (wrong result, data loss, crash, security hole); --for-lane hands both to testing, review and the bughunt role. 2) return origin: move writes 'back-from: <lane>', 'back-by: person|lane' (person = the from-lane is owned by a person, e.g. signoff/human) and 'back-count: <n>' = consecutive returns from that same lane, reset by a return from another lane or a person; show --json prints all three; the in-progress prompt's note-based count is replaced by reading them. 3) explain note format: one line per note with its class in brackets ([blocker] / [kleiner Rand] / [bekannte Grenze] / [nicht aus diesem Ticket]), a status line, one Empfehlung; no 'Warum:' and no fixed [give back]-first order; the cockpit colours the class marks with theme colours. Same template in texts.ts and spawn.sh.
- **2026-11-20 22:47 · Alexander Sacharov** — question: stop-after-fix default for a ticket where nobody set it? 1) on (Recommended) - matches the rule the teamlead writes by hand daily: after a fix round only tests run, no hunt; the field turns it off per ticket. 2) off - TW5FE2's diff-only hunt stays the default; pre-process or a person turns stop-after-fix on per ticket.
- **2026-11-20 22:48 · Alexander Sacharov** — decision: stop-after-fix is off by default (teamlead jaira-64, 2026-11-20). The diff-only hunt after a fix stays the default: it found real defects in fixes (TW5FE2 rebase fallback, JG3REW foreign comment delete, NNV39Y not-ignored). stop-after-fix=true is a deliberate per-ticket opt-out set by pre-process or a person, e.g. for text-only tickets. The brake for endless loops is the count (3 returns from the same lane -> teamlead), not skipping the hunt. Mode stays autonomous.
- **2026-11-20 22:55 · Alexander Sacharov** — in-progress: CRLF in any 'jaira set' value made the write fail ('value position is at end of line') - pre-existing for every field; fixed at the cause in Doc.SetScalar (core/ticket/frontmatter.go), CR/CRLF stored as LF. Concurrency test needs one Store per goroutine: a shared Store races in idIndex, but jaira never shares one (two processes = two stores).
