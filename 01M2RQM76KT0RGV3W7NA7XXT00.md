---
id: 01M2RQM76KT0RGV3W7NA7XXT00
title: Der Lane-Payload liefert einen Ausschnitt des Diffs und meldet ihn als vollstaendig
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Kritik- oder Review-Lane beurteilt den ganzen Branch, oder sie erfaehrt im Payload, dass sie es nicht tut - statt einen Bruchteil zu unterschreiben, ohne es zu merken"
context: |-
  Gefunden von der review-Lane auf GTQHNH am 2026-09-17, nachdem sie selbst darauf hereingefallen war.

  Was passiert: internal/cli/flow.go:589 nimmt 'shas := t.Commits' und leitet die Commit-Liste NUR dann aus git ab, wenn das Feld leer ist. Steht im Ticket-Frontmatter ein 'commits:' mit drei SHAs, waehrend auf dem Branch einundzwanzig liegen, dann zeigt 'jaira show <id> --for-lane review --json' den Diff dieser drei - und meldet complete:true. Nichts im Payload sagt, dass achtzehn Commits fehlen.

  Nachgemessen an drei Tickets desselben Tages: GTQHNH fuehrte 3 SHAs bei 21 auf dem Branch, D8CSAA 1 bei 4, 8XHZT5 hatte das Feld leer und bekam deshalb den ganzen Branch. Die review-Lane von GTQHNH hat den Ausschnitt bemerkt und stattdessen 'git diff origin/HEAD...HEAD' geurteilt; die von D8CSAA hat es nicht bemerkt.

  Warum das gefaehrlich und nicht bloss unpraktisch ist: der Ausfall ist still und sieht aus wie Erfolg. complete:true heisst fuer jeden Leser 'du hast alles'. Eine Review, die einen Bruchteil sieht, findet nichts und unterschreibt - und das ist genau der Fall, gegen den die Review-Lane ueberhaupt existiert.

  Wie 'commits:' teilweise gefuellt wird, ist noch nicht untersucht - vermutlich stempelt ein 'jaira move' den Stand von damals hinein, und spaetere Commits kommen nicht mehr dazu. Das gehoert zur Untersuchung.

  Im Prompt ist das seit f6e8dba (GTQHNH) nur BESCHRIEBEN: jaira-role-lane/SKILL.md sagt der Lane, die den Diff beurteilt, sie solle 'jaira show <id> --json | jq .commits | length' gegen 'git log origin/HEAD..HEAD --oneline | wc -l' zaehlen. Eine Handanweisung ist kein Ersatz dafuer, dass das Werkzeug die Wahrheit sagt.
definition-of-done: "Der Payload einer Lane, die einen Diff beurteilt, zeigt den ganzen Branch - oder er sagt, dass er es nicht tut: entweder leitet showForLane die Liste immer aus git ab, oder complete ist false und 'missing' nennt die Zahl der nicht enthaltenen Commits. Ein Leser kann nicht mehr einen Ausschnitt fuer das Ganze halten."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-17T22:26:05Z
updated-at: 2026-09-18T06:33:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-99942
claimed-at: 2026-09-18T06:30:48Z
mode: conversational
---

# Der Lane-Payload liefert einen Ausschnitt des Diffs und meldet ihn als vollstaendig

## Definition of Done

- [ ] Der Payload einer Lane, die einen Diff beurteilt, zeigt den ganzen Branch - oder er sagt, dass er es nicht tut: entweder leitet showForLane die Liste immer aus git ab, oder complete ist false und 'missing' nennt die Zahl der nicht enthaltenen Commits. Ein Leser kann nicht mehr einen Ausschnitt fuer das Ganze halten.
- [ ] Nachgestellt an dem Fall, der es gezeigt hat: ein Ticket mit drei SHAs in 'commits:' und einundzwanzig Commits auf dem Branch. Mit Test.
- [ ] Untersucht und auf dem Ticket festgehalten, WIE 'commits:' teilweise gefuellt wird - welcher Schreibpfad den Stand einfriert und warum er spaetere Commits nicht nachtraegt. Ohne diese Antwort ist jede Reparatur geraten.
- [ ] Eine Zeile in core/release/NOTES.md, wenn sich aendert, was ein Benutzer im Payload sieht.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-18 06:30 · Alexander Sacharov** — Alex am 2026-09-18: dieses Ticket wird 0.3.1. Der Zweig fix/7XXT00 haengt an release/0.3.0 (88b6816), nicht an master - 0.3.0 ist noch nicht gemerged und noch nicht getaggt. Die Zeile fuer NOTES.md gehoert deshalb unter das leere '## Unreleased' ganz oben, NICHT unter '## 0.3.0': diese Sektion wird gleich getaggt und ist damit geschlossene Historie. Das Umbenennen von '## Unreleased' nach '## 0.3.1' und eine frische leere darueber macht ein Mensch in dem Commit, den er taggt - kein Agent.
- **2026-09-18 06:33 · Alexander Sacharov** — Entscheidung zu DoD 1 (Alex, 2026-09-18, vor der Plan-Lane): Weg A - showForLane leitet die Commit-Liste im case "diff" IMMER aus git ab, statt t.Commits als fertige Antwort zu lesen. t.Commits wird mit dem aus git abgeleiteten Stand vereinigt, damit ein eingetragener SHA, den die Ableitung nicht findet, nicht verloren geht: shas := union(t.Commits, env.DeriveCommits(t)). complete:true bedeutet damit wieder 'du siehst den ganzen Branch'. Kein 'missing'-Feld mit einer Zahl fehlender Commits, und kein zusaetzliches Warnfeld - der stille Ausfall wird beseitigt, nicht gemeldet. Begruendung: ein Flag, das niemand liest, ersetzt einen stillen Fehler nur durch einen lauten Umweg. Preis, der bewusst akzeptiert wird: der beurteilte Diff ist nicht mehr allein aus dem Frontmatter reproduzierbar.
