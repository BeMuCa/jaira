---
id: 01M28FMQGC4CNQT8Z9WY13VMA8
title: Jede Aenderung faehrt auf einem Branch und kommt durch einen PR
status: human
ready: true
creator: Alexander Sacharov
goal: "Es steht als Regel des Projekts geschrieben, dass Arbeit auf einem Branch mit ihrem Ticket faehrt und ueber einen PR ankommt - und dass das Pruefen dieses PRs dem Maintainer gehoert, nicht dem, der ihn aufmacht"
context: |-
  Heute gilt die Regel schon, aber nur als Gewohnheit: die Arbeit an den Ref-Tickets lief auf einem Branch, das Ticket ritt in denselben Commits mit, und der PR (#9) ging ans Original. Nirgends steht das aufgeschrieben.

  Was fehlt, konkret: der generierte jaira-Block sagt 'das Ticket faehrt im selben Commit wie der Code', aber nicht, dass dieser Commit auf einem Branch liegt und master nur durch einen PR erreicht. Ein Agent, der die Anweisungen liest, darf daraus schliessen, direkt auf master zu committen - was heute niemand tut, aber nichts verbietet.

  Und die zweite Haelfte, die genauso fehlt: wer den PR aufmacht, prueft ihn nicht selbst ab. Das Abnehmen gehoert dem Maintainer. Das Projekt sagt dasselbe schon einmal an anderer Stelle - 'a review agent cannot certify its own work' - nur eben nicht ueber PRs.

  Der Platz dafuer ist der Bereich hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md: alles dahinter ueberlebt die naechste Regeneration des Blocks. Dort steht schon die NOTES.md-Regel, also gehoert diese daneben.
definition-of-done: "hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und der PR gehoert von Anfang an dem Maintainer - ein Agent pusht seinen Branch und hoert dort auf, er macht den PR nicht auf, merged ihn nicht und gibt ihn nicht frei; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-11T14:58:42Z
updated-at: 2026-09-14T20:14:15Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Nichts blockiert: testing hat mit pass bestaetigt, dass kein Rollen-Prompt mehr erlaubt, einen PR aufzumachen, dass die drei Dokumentationsstellen und die Prompts dasselbe sagen, und dass jaira-role-pr seine uebrige Arbeit behalten hat. Beim Pruefen kam heraus, dass das Frontmatter-Feld dieses Tickets noch die alte Regel trug, waehrend das Kaestchen im Rumpf schon die neue hatte - ich habe es angeglichen und daraus Ticket NYW4M7 gemacht, weil es der dritte Fall an einem Tag war. Du musst hier nur sagen, ob du die Arbeit annimmst."
outcome-what: "Die doppelte Begruendung aus der neuen PR-Sektion gefaltet - CLAUDE.md und AGENTS.md sagen sie jetzt in der README-Formulierung"
outcome-why: "Absatz 2 war woertlich der Satz, den der generierte jaira-Block 40 Zeilen darueber schon traegt (core/board/announce.go:88-91)"
outcome-resolves: "Die Regel steht unveraendert an allen drei Stellen, nur ohne die Wiederholung; go test ./... gruen"
claimed-by: DESKTOP-RFTCH11-690668
claimed-at: 2026-09-14T19:50:20Z
review-summary: "core/role/builtin/jaira-role-pr/SKILL.md:3,6,10-12 - die Rolle heisst 'Open a pull request for finished ticket work' und sagt woertlich 'Opening a pull request is a contributor's job. Accepting one is the maintainer's. You are the contributor: you may open'. Das ist genau der Satz, den CLAUDE.md:166 jetzt umdreht. Ein ausgeliefertes Prompt weist den Agenten an zu tun, was die Projektregel ihm verbietet. || core/role/builtin/jaira-teamlead/SKILL.md:79-80 - 'Opening one is a contributor's job, accepting it is the maintainer's' - derselbe alte Satz ein zweites Mal, ebenfalls ausgeliefert. || core/role/builtin/jaira-teamlead/SKILL.md:89 - 'Close it once the pull request is open, not once it is merged' nennt als Schlusspunkt einen Moment, den ein Agent unter der neuen Regel nie erlebt, weil er den PR nicht aufmacht. || Der Diff schreibt die Regel an drei Dokumentationsstellen auf, aber die Prompts, die Agenten tatsaechlich ausfuehren, liegen im selben Repository und tragen weiter die alte. Wer die Regel in CLAUDE.md liest und die Rolle /jaira-role-pr benutzt, bekommt zwei Anweisungen, die sich widersprechen."
review-gaps: "Docs-only change, one duplication removed: the new section's second paragraph repeated, word for word, a sentence the generated jaira block already carries 40 lines above it in the same file (core/board/announce.go:88-91, shipped into every board) - folded into the first paragraph using the phrasing the README already uses at the same place, so all three copies of the rule now read alike. Left alone and why: the rule standing in three files is the DoD, not duplication (README serves readers who never run an agent); the missing blank line before <!-- jaira:end --> is harmless, the marker parser matches the line and Markdown swallows the comment; the NOTES.md rule sitting behind jaira:end instead of jaira:local in CLAUDE.md is pre-existing and outside the regenerated block either way; core/role/builtin still carries the old rule - that is the critique's open A/B question and DoD item 2, out of this lane. Nothing dead, nothing to hoist, go test ./... green."
test-verdict: "pass: alle sieben ausgelieferten Prompts sagen jetzt dasselbe wie die Dokumentation - kein 'may open' und keine andere Formulierung von 'mach den PR auf' mehr in core/role/builtin; jaira-role-pr behaelt Ordner- und frontmatter-Namen (core/role/role_test.go:18) und kann weiter zu einem offenen PR pushen und Review-Kommentare beantworten; NOTES.md traegt eine einzeilige Unreleased-Zeile; go test ./... -race gruen, RC=0"
---

# Jede Aenderung faehrt auf einem Branch und kommt durch einen PR

## Definition of Done

- [x] hinter dem jaira:local-Marker in CLAUDE.md und AGENTS.md steht die Regel: Arbeit laeuft auf einem Branch, das Ticket faehrt in denselben Commits mit, master wird nur durch einen PR erreicht, und der PR gehoert von Anfang an dem Maintainer - ein Agent pusht seinen Branch und hoert dort auf, er macht den PR nicht auf, merged ihn nicht und gibt ihn nicht frei; dieselbe Regel steht im README unter Development, damit sie auch findet, wer nie einen Agenten benutzt
  proof: CLAUDE.md:156-169 und AGENTS.md:166-179 hinter dem jaira:local-Marker, README.md:842-851 unter Development
- [x] Die ausgelieferten Rollen-Prompts sagen dasselbe wie die Dokumentation: core/role/builtin/jaira-role-pr/SKILL.md:11 ('You are the contributor: you may open, push, and answer') und die entsprechenden Stellen in jaira-teamlead/SKILL.md tragen die neue Regel - der Agent schiebt den Zweig und haelt an, das Aufmachen gibt der Mensch frei. Nachgestellt: kein 'may open' mehr in core/role/builtin.
  proof: core/role/builtin/jaira-role-pr/SKILL.md:10-13 and core/role/builtin/jaira-teamlead/SKILL.md:79-80,90-91; grep -rni 'may open' core/role/builtin returns nothing; go test ./... -race green

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
- **2026-09-14 20:04 · Alexander Sacharov** — testing 14.09.: Gates gruen - go test ./... -race, RC=0, kein FAIL, internal/tui 131s, internal/wintrap ok. Die go:embed-Prompts bauen also.

DoD 1 ERFUELLT, im Baum nachgelesen: CLAUDE.md:154 <!-- jaira:local -->, Abschnitt 'Work rides on a branch and arrives through a PR' 156-169, Endmarke :170 - der Abschnitt liegt also wirklich zwischen local und end und ueberlebt die Regeneration (core/board/announce.go:24, jairaMarkerLocal). AGENTS.md:127 local, Abschnitt 166-179, :180 end - ebenfalls drin. README.md:842-851 unter Development. Alle drei sagen dasselbe in derselben Reihenfolge: Branch, Ticket faehrt in denselben Commits, master nur ueber PR, und 'the pull request belongs to the maintainer from the moment it exists - an agent pushes its branch and stops there: it does not open the pull request, does not merge one, and does not approve one'. Kein Widerspruch im Detail zwischen den drei Kopien; die einzige Abweichung ist stilistisch ('the person who wrote' vs 'whoever wrote').

DoD 2 NICHT ERFUELLT - das ist der Befund. Die Prompts tragen weiter die alte Regel, obwohl der Fortschrittseintrag vom 14.09. 19:51 Variante B entschieden hat (dieses Ticket aendert auch die ausgelieferten Prompts). grep 'may open' core/role/builtin:
- core/role/builtin/jaira-role-pr/SKILL.md:11 'You are the contributor: you may open, push, and answer.'
- core/role/builtin/jaira-role-pr/SKILL.md:3 description 'Open a pull request for finished ticket work'
- core/role/builtin/jaira-role-pr/SKILL.md:6 Ueberschrift 'Open it, answer it, never accept it'
- core/role/builtin/jaira-role-pr/SKILL.md:10 'Opening a pull request is a contributor's job.'
- core/role/builtin/jaira-teamlead/SKILL.md:79-80 'You never merge a pull request and never approve your own. Opening one is a contributor's job, accepting it is the maintainer's.'
- core/role/builtin/jaira-teamlead/SKILL.md:89 'Close it once the pull request is open, not once it is merged.' - der Dispatcher wartet damit auf einen Moment, den ein Agent unter der neuen Regel nie erlebt.
- core/role/builtin/jaira-dispatcher/SKILL.md:157 'not when the pull request merges' ist unauffaellig, betrifft nur Worktree-Lebensdauer.

Damit ist genau der Zustand da, der schlimmer ist als vorher: Dokumentation sagt 'macht den PR nicht auf', der Prompt, den der Agent ausfuehrt, sagt 'you may open'. Ein Leser weiss nicht, was gilt.

Vorschlag fuer in-progress: jaira-role-pr auf reines Pushen zurueckbauen (Name, description, Ueberschrift und Zeile 10-12 drehen: Branch pushen, anhalten, das Aufmachen gibt der Mensch in Auftrag; 'Never gh pr merge' bleibt), jaira-teamlead:79-80 auf dieselbe Formulierung drehen und :89 auf 'sobald der Branch gepusht und gemeldet ist'. Dazu eine NOTES.md-Zeile unter ## Unreleased - die Rollen-Prompts sind client-facing (jaira roles install). Kein Lane-Wechsel von mir: der Auftrag war nur testing.
- **2026-09-14 20:05 · Alexander Sacharov** — testing hat mit fail zurueckgeschickt, und der Befund ist genau der, den ich beim Entscheiden der A/B-Frage benannt habe: die Dokumentation sagt, ein Agent macht keinen PR auf, und der Prompt, den derselbe Agent ausfuehrt, sagt 'you may open, push, and answer'. Das ist schlechter als der alte Zustand, weil ein Leser nicht mehr weiss, was gilt.

Betroffen sind mehr Stellen als die eine, die ich genannt hatte: jaira-role-pr/SKILL.md Zeilen 3, 6, 10 und 11, jaira-teamlead/SKILL.md:79-80 ('Opening one is a contributor's job') und :89 ('Close it once the pull request is open') - der letzte Satz beschreibt einen Moment, den ein Agent unter der neuen Regel nie erreicht, der Abschnitt muss also umgeschrieben und nicht nur korrigiert werden.

Reihenfolge, die eingehalten werden muss: dieses Ticket zuerst, KSGSKK danach. KSGSKK traegt die installierten Fassungen aus ~/.claude/skills nach core/role/builtin - und die installierten Fassungen tragen die alte PR-Regel immer noch, weil ich dort nur den Transport-Abschnitt gepatcht habe. Wer KSGSKK vor diesem Ticket macht, holt die alte Regel zurueck.
- **2026-09-14 20:09 · Alexander Sacharov** — in-progress nach testing-fail: die ausgelieferten Prompts gedreht. jaira-role-pr behaelt Namen und Ordner - core/role/role_test.go:18 listet 'jaira-role-pr' namentlich und prueft, dass der frontmatter-name gleich dem Ordnernamen ist; die Rolle abschaffen haette also den Test mitgeaendert. Stattdessen zurueckgebaut: description, Ueberschrift und Zeile 10-13 sagen jetzt 'push und halt an', der Abschnitt 'Before you open anything' heisst 'Before you push anything' und endet mit 'git push -u origin HEAD und stop', und der PR-Beschreibungs-Abschnitt bleibt erhalten, aber als Uebergabe an den Menschen ('Hand back the description, do not open it yourself') - der Inhalt war gut und wird weiter gebraucht, nur macht ihn jetzt nicht der Agent auf. 'Never gh pr create' ist zu den Boundaries dazugekommen. jaira-teamlead:89 war kein Wortfix: der Schlusspunkt des Dispatchers heisst jetzt 'sobald der Branch gepusht und die Arbeit gemeldet ist'. Bewusst nicht angefasst: jaira-dispatcher/SKILL.md:157 ('not when the pull request merges') - das begrenzt die Lebensdauer des Worktrees, nicht wer den PR aufmacht, und widerspricht der neuen Regel nicht; ausserdem liegt es im Transport-Abschnitt, der KSGSKK gehoert. Die vier uebrigen Rollen (brainstorm, lane, research, tester) nennen PRs gar nicht, nur 'git merge-base' im tester.
- **2026-09-14 20:12 · Alexander Sacharov** — testing 14.09. (2. Durchlauf, nach dem fail um 20:04): Befund selbst nachgeprueft, nicht uebernommen.

1. grep -rni 'may open' core/role/builtin: leer. Breiter gesucht (open/opening/create a pull request, 'gh pr create', 'open a PR') ueber ALLE sieben Prompts, nicht nur die zwei geaenderten: jeder Treffer sagt jetzt das Richtige. jaira-role-pr/SKILL.md:3,6,10-14,36-37,75 ('Never gh pr create'), jaira-teamlead/SKILL.md:79-80,90-91. brainstorm, lane, research nennen PRs gar nicht; tester nur 'git merge-base' (kein PR).

2. Die drei Dokumentationskopien und die Prompts sagen dasselbe in der Sache: CLAUDE.md:156-169 (zwischen :154 local und :170 end), AGENTS.md:166-179 (zwischen :127 und :180), README.md:842-851 unter ## Development (Ueberschrift :813). Ueberall: Branch pushen und anhalten, der Mensch macht den PR auf, mergen und freigeben bleiben verboten. Kein 'macht auf, merged aber nicht' mehr uebrig.

3. jaira-role-pr verliert seine restliche Arbeit nicht: :59-72 'Answering review comments' steht vollstaendig ('once a person has opened the pull request, the branch is yours to push to and the thread is yours to answer'), :39-57 uebergibt die fertige PR-Beschreibung, :79-80 meldet die PR-URL, wenn schon eine offen ist.

4. Ordner core/role/builtin/jaira-role-pr und frontmatter name: jaira-role-pr unveraendert - core/role/role_test.go:18 haelt.

5. core/release/NOTES.md: erste Zeile unter ## Unreleased, 420 Zeichen auf EINER Zeile, als Anweisung formuliert (inkl. 'jaira roles install --force'), kein Umbruch.

6. go test ./... -race: RC=0, kein FAIL, core/role ok, core/lane 5.2s, internal/wintrap 1.5s. Die go:embed-Prompts bauen.

Ausdruecklich NICHT als Fehler gewertet (gehoert KSGSKK): jaira-dispatcher/SKILL.md:157 'not when the pull request merges' begrenzt die Lebensdauer des Worktrees, nicht wer den PR aufmacht.

Offen, aber kein Grund fuer fail: die frontmatter-Zeile definition-of-done traegt weiter den ueberholten Wortlaut 'ein Agent macht ihn auf und merged ihn nie' - die Checkbox im Body ist gedreht, die frontmatter-Kopie nicht. Wer spaeter nur die Frontmatter liest, liest die alte Regel.
