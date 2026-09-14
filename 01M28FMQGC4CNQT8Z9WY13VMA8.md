---
id: 01M28FMQGC4CNQT8Z9WY13VMA8
title: Jede Aenderung faehrt auf einem Branch und kommt durch einen PR
status: testing
ready: true
creator: Alexander Sacharov
goal: "Es steht als Regel des Projekts geschrieben, dass Arbeit auf einem Branch mit ihrem Ticket faehrt und ueber einen PR ankommt - und dass das Pruefen dieses PRs dem Maintainer gehoert, nicht dem, der ihn aufmacht"
context: |-
  Heute gilt die Regel schon, aber nur als Gewohnheit: die Arbeit an den Ref-Tickets lief auf einem Branch, das Ticket ritt in denselben Commits mit, und der PR (#9) ging ans Original. Nirgends steht das aufgeschrieben.

  Was fehlt, konkret: der generierte jaira-Block sagt 'das Ticket faehrt im selben Commit wie der Code', aber nicht, dass dieser Commit auf einem Branch liegt und master nur durch einen PR erreicht. Ein Agent, der die Anweisungen liest, darf daraus schliessen, direkt auf master zu committen - was heute niemand tut, aber nichts verbietet.

  Und die zweite Haelfte, die genauso fehlt: wer den PR aufmacht, prueft ihn nicht selbst ab. Das Abnehmen gehoert dem Maintainer. Das Projekt sagt dasselbe schon einmal an anderer Stelle - 'a review agent cannot certify its own work' - nur eben nicht ueber PRs.

  Der Platz dafuer ist der Bereich hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md: alles dahinter ueberlebt die naechste Regeneration des Blocks. Dort steht schon die NOTES.md-Regel, also gehoert diese daneben.
definition-of-done: "hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und das Abnehmen des PRs gehoert dem Maintainer - ein Agent macht ihn auf und merged ihn nie; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt; dieser Branch und sein PR sind selbst das erste Beispiel dafuer"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-11T14:58:42Z
updated-at: 2026-09-14T20:00:52Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Die Regel steht jetzt in CLAUDE.md, AGENTS.md und README: ein Agent pusht seinen Branch und macht den PR nicht auf. Die Rollen-Prompts in core/role/builtin sagen aber weiter das Gegenteil - jaira-role-pr/SKILL.md ist als ganze Rolle 'mach den PR auf' gebaut, jaira-teamlead/SKILL.md:79-80 traegt woertlich den alten Satz, und :89 laesst den Dispatcher schliessen, 'sobald der PR offen ist'. Deine Entscheidung, weil beides vertretbar ist: (A) Die Regel gilt nur fuer dieses Projekt - hinter jaira:local, wo sie steht. Dann bleibt core/role/builtin unangetastet, kostet aber, dass jeder Agent hier ein ausgeliefertes Prompt liest, das ihm das Gegenteil sagt. (B) Die Regel gilt jaira-weit. Dann muessen die drei Prompt-Stellen gedreht werden, jaira-role-pr wird auf reines Pushen zurueckgebaut oder abgeschafft, und es braucht eine NOTES.md-Zeile plus 'jaira roles install --global --force' fuer alle - das ist ein eigenes Ticket, nicht mehr dieses."
outcome-what: "Die doppelte Begruendung aus der neuen PR-Sektion gefaltet - CLAUDE.md und AGENTS.md sagen sie jetzt in der README-Formulierung"
outcome-why: "Absatz 2 war woertlich der Satz, den der generierte jaira-Block 40 Zeilen darueber schon traegt (core/board/announce.go:88-91)"
outcome-resolves: "Die Regel steht unveraendert an allen drei Stellen, nur ohne die Wiederholung; go test ./... gruen"
claimed-by: DESKTOP-RFTCH11-690668
claimed-at: 2026-09-14T19:50:20Z
review-summary: "core/role/builtin/jaira-role-pr/SKILL.md:3,6,10-12 - die Rolle heisst 'Open a pull request for finished ticket work' und sagt woertlich 'Opening a pull request is a contributor's job. Accepting one is the maintainer's. You are the contributor: you may open'. Das ist genau der Satz, den CLAUDE.md:166 jetzt umdreht. Ein ausgeliefertes Prompt weist den Agenten an zu tun, was die Projektregel ihm verbietet. || core/role/builtin/jaira-teamlead/SKILL.md:79-80 - 'Opening one is a contributor's job, accepting it is the maintainer's' - derselbe alte Satz ein zweites Mal, ebenfalls ausgeliefert. || core/role/builtin/jaira-teamlead/SKILL.md:89 - 'Close it once the pull request is open, not once it is merged' nennt als Schlusspunkt einen Moment, den ein Agent unter der neuen Regel nie erlebt, weil er den PR nicht aufmacht. || Der Diff schreibt die Regel an drei Dokumentationsstellen auf, aber die Prompts, die Agenten tatsaechlich ausfuehren, liegen im selben Repository und tragen weiter die alte. Wer die Regel in CLAUDE.md liest und die Rolle /jaira-role-pr benutzt, bekommt zwei Anweisungen, die sich widersprechen."
review-gaps: "Docs-only change, one duplication removed: the new section's second paragraph repeated, word for word, a sentence the generated jaira block already carries 40 lines above it in the same file (core/board/announce.go:88-91, shipped into every board) - folded into the first paragraph using the phrasing the README already uses at the same place, so all three copies of the rule now read alike. Left alone and why: the rule standing in three files is the DoD, not duplication (README serves readers who never run an agent); the missing blank line before <!-- jaira:end --> is harmless, the marker parser matches the line and Markdown swallows the comment; the NOTES.md rule sitting behind jaira:end instead of jaira:local in CLAUDE.md is pre-existing and outside the regenerated block either way; core/role/builtin still carries the old rule - that is the critique's open A/B question and DoD item 2, out of this lane. Nothing dead, nothing to hoist, go test ./... green."
---

# Jede Aenderung faehrt auf einem Branch und kommt durch einen PR

## Definition of Done

- [x] hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und der PR gehoert von Anfang an dem Maintainer - ein Agent pusht seinen Branch und hoert dort auf, er macht den PR nicht auf, merged ihn nicht und gibt ihn nicht frei; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt
  proof: CLAUDE.md:156-169 und AGENTS.md:166-179 hinter dem jaira:local-Marker, README.md:842-851 unter Development
- [ ] Die ausgelieferten Rollen-Prompts sagen dasselbe wie die Dokumentation: core/role/builtin/jaira-role-pr/SKILL.md:11 ('You are the contributor: you may open, push, and answer') und die entsprechenden Stellen in jaira-teamlead/SKILL.md tragen die neue Regel - der Agent schiebt den Zweig und haelt an, das Aufmachen gibt der Mensch frei. Nachgestellt: kein 'may open' mehr in core/role/builtin.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] CLAUDE.md: turn the PR rule around behind the jaira:local marker - the agent pushes the branch and stops
- [x] AGENTS.md: add the same section behind the jaira:local marker, it is missing there entirely
- [x] README.md under Development: turn the same two sentences around so a reader who never uses an agent finds the rule

## Progress
- **2026-09-14 18:22 · Alexander Sacharov** — Alex hat die offene Frage am 14.09. beantwortet: ein Agent macht den Pull Request NICHT auf. Er pusht den Branch und hoert dort auf; das Aufmachen gibt der Mensch in Auftrag. Mergen und Freigeben waren schon verboten - das verschiebt die Linie nach vorn, aufs Aufmachen. Damit steht die Regel heute an drei Stellen falsch herum und alle drei muessen gedreht werden: CLAUDE.md hinter dem jaira:local-Marker, AGENTS.md hinter dem jaira:local-Marker, README.md unter Development. Die Formulierung des DoD-Kriteriums ('ein Agent macht ihn auf und merged ihn nie') ist damit selbst ueberholt - erfuellt ist es jetzt durch die gedrehte Regel, nicht durch die woertliche.
- **2026-09-14 18:25 · Alexander Sacharov** — Der Abschnitt stand schon in CLAUDE.md und README - nur falsch herum ('an agent may open a pull request'). In AGENTS.md stand er ebenfalls schon, aber weiter unten hinter der NOTES.md-Regel und ohne Leerzeile vor der Ueberschrift, weshalb eine Suche nach dem Marker plus 30 Zeilen ihn nicht findet. Wer hier nachsieht, ob eine Regel fehlt, muss nach dem Ueberschriftentext greppen, nicht nach dem Marker. Die Leerzeile ist mitgefixt. Kein NOTES.md-Eintrag: nur Dokumentation, von aussen am Binary nichts zu beobachten.
- **2026-09-14 18:28 · Alexander Sacharov** — critique: Die gedrehte Regel steht in CLAUDE.md, AGENTS.md und README - aber die in diesem Repository gepflegten Rollen-Prompts sagen weiter das Gegenteil. core/role/builtin/jaira-role-pr/SKILL.md ist als ganze Rolle 'mach den PR auf' gebaut (Zeile 3, 6, 10-12), core/role/builtin/jaira-teamlead/SKILL.md:79-80 traegt woertlich den alten Satz, und Zeile 89 setzt den Schlusspunkt des Dispatchers auf 'sobald der PR offen ist' - ein Moment, den ein Agent unter der neuen Regel nie erreicht. Nicht selbst gefixt, weil es eine Entscheidung braucht: CLAUDE.md/AGENTS.md hinter jaira:local sind die Regel DIESES Projekts, core/role/builtin geht an jedes jaira-Board. Ob die Regel jaira-weit gilt oder nur hier, aendert den Umfang von 'zwei Saetze umschreiben' auf 'die Rolle jaira-role-pr abschaffen oder auf reines Pushen zurueckbauen, plus NOTES.md-Zeile'. Sonst nichts gefunden: die drei Textstellen sagen dasselbe, die Platzierung hinter dem Marker stimmt, die NOTES.md-Regel in CLAUDE.md liegt hinter jaira:end und war vorher wie nachher ausserhalb des regenerierten Blocks - kein Verlust durch den neuen Marker.
- **2026-09-14 19:50 · Alexander Sacharov** — Dispatcher 14.09.: Ticket lag in human, obwohl nur critique gelaufen war - optimize, testing und review haben nie ein Feld hinterlassen. Auftrag von Alex: die Route nachholen, den Inhalt nicht wieder aufmachen (die gedrehte PR-Regel steht und gilt). Die critique-Frage A/B bleibt offen und geht mit dem Bericht nach oben; fuer diesen Durchlauf gilt A - der Umbau von core/role/builtin ist nicht dieses Ticket. Das Ticket KSGSKK fasst jaira-dispatcher/SKILL.md und jaira-teamlead/SKILL.md ohnehin an; die PR-Regel in jaira-role-pr braucht davon unabhaengig eine eigene Entscheidung.
- **2026-09-14 19:51 · Alexander Sacharov** — Scope-Frage aus critique am 2026-09-14 entschieden, Variante B: dieses Ticket aendert auch die ausgelieferten Prompts, nicht nur die Dokumentation.

Begruendung: das Ziel des Tickets ist, dass die Regel GILT, nicht dass sie irgendwo steht. CLAUDE.md, AGENTS.md und das README lesen Menschen. core/role/builtin liest der Agent, und dort steht heute das Gegenteil - jaira-role-pr/SKILL.md:11 sagt woertlich 'you may open, push, and answer'. Eine Regel, der der massgebliche Text widerspricht, ist keine Regel: der naechste Agent liest seinen eigenen Prompt und macht den PR auf, voellig korrekt nach dem, was er vor sich hat.

Damit waechst das Ticket um zwei Dateien, nicht um ein Thema. KSGSKK fasst dieselben zwei SKILL.md-Dateien gleich danach im selben Worktree an; das nacheinander zu tun ist billiger, als es zu trennen.
- **2026-09-14 19:53 · Alexander Sacharov** — optimize: Nur Dokumentation, nichts zu loeschen ausser einer Wiederholung. Absatz 2 der neuen Sektion ('The ticket travelling with the code is what makes the pull request readable: a reviewer sees the change and what it was for in one place, instead of a diff whose ticket is still in whatever state the last commit left it.') stand woertlich schon 40 Zeilen weiter oben in derselben Datei - der generierte jaira-Block sagt denselben Satz, er kommt aus core/board/announce.go:88-91 und wird in jedes Board geschrieben. In CLAUDE.md und AGENTS.md zu einem Nebensatz gefaltet, in der Formulierung, die das README an derselben Stelle schon benutzt; damit sagen alle drei Kopien der Regel dasselbe in derselben Laenge. Bewusst stehen gelassen: (1) die Regel steht dreimal in drei Dateien - das fordert das DoD, README ist fuer Leser ohne Agenten; (2) die fehlende Leerzeile vor <!-- jaira:end --> in beiden Dateien - der Parser in core/board/announce.go matcht die Markerzeile, Markdown schluckt den Kommentar, kein Defekt; (3) die NOTES.md-Regel liegt in CLAUDE.md hinter jaira:end statt hinter jaira:local - ausserhalb des Blocks, ueberlebt die Regeneration genauso, vorbestehend; (4) core/role/builtin - das ist die offene A/B-Frage der critique, nicht diese Lane. go test ./... gruen.
- **2026-09-14 20:00 · Alexander Sacharov** — Wie bei 74VM40: zum zweiten Mal aus human zurueck, wieder ohne test-verdict, und wieder nicht aus Nachlaessigkeit des Dispatchers. jaira erzwingt keine Reihenfolge - 'after' ordnet nur die Anzeige (core/lane/lane.go:43), und human traegt precedence 40, steht also vor den drei Schleifen-Lanes. Wer die naechste Lane sucht, landet regelkonform dort. Siehe D28H7V.
