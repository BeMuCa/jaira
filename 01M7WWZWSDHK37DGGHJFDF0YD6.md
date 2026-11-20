---
id: 01M7WWZWSDHK37DGGHJFDF0YD6
title: "Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt"
status: backlog
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
updated-at: 2026-11-20T17:45:38Z
updated-by: Alexander Sacharov
---

# Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt

## Definition of Done

- [ ] Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)
- [ ] Mit stop-after-fix laeuft nach einem Fix-Durchgang nur der Testlauf, kein neuer Bughunt
- [ ] Das Cockpit zeigt Explain nur, wenn review-gaps eine blocker:-Zeile hat oder das Urteil einen Vorbehalt nennt; reine note:-Zeilen stehen als Liste ohne Knopf (Test in board.test.ts)
- [ ] jaira zeichnet bei jedem Rueckweg die Herkunfts-Lane auf und ob eine Person oder eine Schleifen-Lane zurueckgeschickt hat; jaira show --json gibt beides aus, und die Schleifenzaehlung zaehlt nur Rueckgaben derselben Schleifen-Lane (Test: signoff->in-progress zaehlt nicht)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:48 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 17:05 · Alexander Sacharov** — Aus TW5FE2 (bughunt e12f36d5): der Prompt-Zaehler in in-progress kann eine Rueckgabe der Person aus signoff nicht von einer review-Rueckgabe unterscheiden - jaira speichert die Herkunfts-Lane nicht (show --json unterscheidet sich nur in status/next_lane/round; sendsBack in core/move/move.go:133 erhoeht round bei jedem Rueckweg, updated_by ist derselbe git-Nutzer). Folge: nach zwei echten review-Rueckgaben eskaliert in-progress, statt zu tun, was die Person wollte. Die echte Zaehlung gehoert ins Binary - neuer DoD-Punkt hier.
