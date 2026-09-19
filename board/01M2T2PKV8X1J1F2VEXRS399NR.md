---
id: 01M2T2PKV8X1J1F2VEXRS399NR
title: "Follow-up: Der Lane-Payload liefert einen Ausschnitt des Diffs und meldet ihn als vollstaendig"
status: backlog
ready: false
creator: Alexander Sacharov
assignee: Alexander Sacharov
context: "Raised from the review of 7XXT00.\n\nGefunden von der review-Lane auf GTQHNH am 2026-09-17, nachdem sie selbst darauf hereingefallen war.\n\nWas passiert: internal/cli/flow.go:589 nimmt 'shas := t.Commits' und leitet die Commit-Liste NUR dann aus git ab, wenn das Feld leer ist. Steht im Ticket-Frontmatter ein 'commits:' mit drei SHAs, waehrend auf dem Branch einundzwanzig liegen, dann zeigt 'jaira show <id> --for-lane review --json' den Diff dieser drei - und meldet complete:true. Nichts im Payload sagt, dass achtzehn Commits fehlen.\n\nNachgemessen an drei Tickets desselben Tages: GTQHNH fuehrte 3 SHAs bei 21 auf dem Branch, D8CSAA 1 bei 4, 8XHZT5 hatte das Feld leer und bekam deshalb den ganzen Branch. Die review-Lane von GTQHNH hat den Ausschnitt bemerkt und stattdessen 'git diff origin/HEAD...HEAD' geurteilt; die von D8CSAA hat es nicht bemerkt.\n\nWarum das gefaehrlich und nicht bloss unpraktisch ist: der Ausfall ist still und sieht aus wie Erfolg. complete:true heisst fuer jeden Leser 'du hast alles'. Eine Review, die einen Bruchteil sieht, findet nichts und unterschreibt - und das ist genau der Fall, gegen den die Review-Lane ueberhaupt existiert.\n\nWie 'commits:' teilweise gefuellt wird, ist noch nicht untersucht - vermutlich stempelt ein 'jaira move' den Stand von damals hinein, und spaetere Commits kommen nicht mehr dazu. Das gehoert zur Untersuchung.\n\nIm Prompt ist das seit f6e8dba (GTQHNH) nur BESCHRIEBEN: jaira-role-lane/SKILL.md sagt der Lane, die den Diff beurteilt, sie solle 'jaira show <id> --json | jq .commits | length' gegen 'git log origin/HEAD..HEAD --oneline | wc -l' zaehlen. Eine Handanweisung ist kein Ersatz dafuer, dass das Werkzeug die Wahrheit sagt.\n\nThe reviewer said:\n\n> Der Diff erfuellt die acht DoD-Punkte, und er tut es an der Wurzel statt an der Meldung: showForLane und der Signoff-Schirm leiten beide immer aus git ab, der Worktree haengt mit dran, und die Herkunft steht als Token daneben. Tests decken jeden Zweig, den man ohne kaputtes Repo erreichen kann, und die Notizen halten zu jeder Entscheidung die verworfene Alternative fest. Keine Defekte gefunden.\n> \n> Ein Befund gehoert vor die Annahme, nicht danach: Befund 1 (SKILL.md:122-135 beschreibt weiter das alte Verhalten und weist die nebenher laufende Kritik an, genau den Diff zu ignorieren, den dieser Change ihr gerade gibt). Die Datei wird per go:embed ausgeliefert und die NOTES.md-Zeile fordert 'jaira roles install --global --force' - es wuerde also eine Anweisung ausgerollt, die der Binary widerspricht, mit der sie kommt. Das sind zehn Zeilen Prosa in der Datei, die dieser Change ohnehin anfasst, kein neuer Mechanismus. Befund 2 und 3 sind Nachtraege, kein Grund zurueckzuschicken.\n> \n> Unsicher bin ich bei genau einer Sache und sage es lieber, als sie zu runden: ob der Worktree-Anteil auf einem geteilten Worktree stoert, kann ich nicht pruefen - hier stimmt 'ein Worktree je Ticket', und wo das nicht gilt, beurteilt eine Lane fremde Aenderungen mit. Die Notiz vom 07:35 nennt den Preis, ein Test kann ihn nicht abbilden."
blocked-by: []
follows: 01M2RQM76KT0RGV3W7NA7XXT00
commits: []
created-at: 2026-09-18T10:58:52Z
updated-at: 2026-09-18T10:58:52Z
---

# Follow-up: Der Lane-Payload liefert einen Ausschnitt des Diffs und meldet ihn als vollstaendig

## Definition of Done

- [ ] <What must be true that is not true yet>

## Progress

