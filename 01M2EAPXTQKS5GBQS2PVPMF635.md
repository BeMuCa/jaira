---
id: 01M2EAPXTQKS5GBQS2PVPMF635
title: "Die Rollen-Prompts werden im Repository gepflegt, nicht im Heimverzeichnis"
status: human
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "core/role/builtin ist die einzige Quelle der Rollen-Prompts: wer eine Rolle aendert, aendert die Datei im Repository, und ~/.claude/skills entsteht danach nur noch aus 'jaira roles install --global'"
context: |-
  Die ausgelieferten Prompts sind schon veraltet, am Tag der Auslieferung.
  C9V7ZV hat die sieben Rollen aus ~/.claude/skills nach core/role/builtin eingefroren. Danach wurden zwei davon im Heimverzeichnis weitergeschrieben, das Binary kennt diese Aenderungen nicht.
  Zahlen von heute: ~/.claude/skills/dispatcher/SKILL.md 170 Zeilen gegen core/role/builtin/jaira-dispatcher/SKILL.md 122. teamlead 92 gegen 82. Die fuenf role-* stimmen byteweise ueberein.
  Was im Binary fehlt: die verschaerfte Drei-Runden-Regel, der Abschnitt 'A human typing in a workers tab is not a fault', die Pflicht pro Lane eine Zeile zu melden, das Schliessen eines Worker-Tabs sobald seine Lane fertig ist, und in jaira-teamlead 'Close what you started'.
  Solange das so steht, darf niemand die unpraefixierten Ordner loeschen: 'roles install --global' wuerde die Fassung vor der Verschaerfung zurueckschreiben.
  Das Ziel ist nicht nur ein einmaliger Abgleich, sondern die Richtung umzudrehen - ab hier wird im Repository editiert und von dort installiert, nicht umgekehrt.
  Beim Uebernehmen gilt dieselbe Regel wie beim Einfrieren: der Ordnername ist der Kommandoname, frontmatter name: muss denselben String tragen, und jeder Querverweis auf eine Schwesterrolle braucht das jaira-Praefix.
  Nicht Teil dieses Tickets: eine Automatik, die beide Seiten dauerhaft synchron haelt.
definition-of-done: "core/role/builtin/jaira-dispatcher/SKILL.md und jaira-teamlead/SKILL.md tragen den Text aus ~/.claude/skills, mit jaira-Praefix in name: und in jeder Querverweis-Zeile; core/role/role_test.go TestCrossReferencesCarryThePrefix gruen; go test ./... -race gruen; 'jaira roles install --global --force' legt die sieben jaira-<id> Ordner an und ein diff gegen core/role/builtin zeigt keinen Unterschied; die sieben unpraefixierten Ordner teamlead, dispatcher, role-brainstorm, role-lane, role-pr, role-research, role-tester sind aus ~/.claude/skills entfernt; eine Zeile in core/release/NOTES.md unter ## Unreleased"
tags:
  - cli
blocked-by: []
parent: 01M2E85S75MEF7YJJRJ6C9QS8F
related: []
commits: []
created-at: 2026-09-13T21:27:58Z
updated-at: 2026-09-13T22:01:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-44645
claimed-at: 2026-09-13T21:37:52Z
outcome-what: "core/role/builtin/jaira-dispatcher/SKILL.md und jaira-teamlead/SKILL.md tragen jetzt den Text aus ~/.claude/skills. Geaendert wurde daran nur 'name:' im frontmatter und vier Querverweis-Zeilen, die /role-lane und /role-tester auf /jaira-role-lane und /jaira-role-tester heben. Eine Zeile in core/release/NOTES.md unter ## Unreleased. Danach 'roles install --global --force' gelaufen, diff -r fuer alle sieben jaira-<id> Ordner sauber, und die sieben unpraefixierten Ordner aus ~/.claude/skills entfernt."
outcome-why: "Die ausgelieferten Prompts waren am Tag der Auslieferung schon veraltet: dem Binary fehlten die Drei-Runden-Regel, 'A human typing in a workers tab is not a fault', die Meldepflicht pro Lane, das Schliessen eines Worker-Tabs und 'Close what you started'. Solange die unpraefixierten Ordner danebenlagen, war unklar welche Fassung gilt; jetzt ist core/role/builtin die einzige Quelle."
outcome-resolves: "Definition of Done Punkt 1 vollstaendig: Text uebernommen, TestCrossReferencesCarryThePrefix und go test ./... -race gruen, Installation deckungsgleich, alte Ordner weg, NOTES.md-Zeile geschrieben."
review-summary: none
test-verdict: "pass: go test ./... -race gruen (RC=0), TestCrossReferencesCarryThePrefix gruen, 'roles install --global --force' RC=0 und diff -r gegen core/role/builtin fuer alle sieben jaira-<id> Ordner ohne Unterschied, die sieben unpraefixierten Ordner sind weg, NOTES.md-Zeile unter ## Unreleased vorhanden"
question: "Die Uebernahme ist geprueft: installierte jaira-* Dateien und core/role/builtin sind deckungsgleich, die sieben unpraefixierten Ordner sind weg, go test ./... -race gruen. Zwei Dinge brauchen dein Ja, weil kein Test sie abdecken kann. Erstens der Inhalt der Prompts: geprueft wurde nur, dass jedes name: und jeder Querverweis das jaira-Praefix traegt - ob der verschaerfte Text vollstaendig und richtig uebernommen wurde, sieht kein Test, das ist dein Text und deine Abnahme. Zweitens: nach der Testing-Runde kam auf denselben Zweig noch f3bc433 - der Abschnitt 'Where a worktree goes' in jaira-dispatcher, die .worktrees-Pfadformel in spawn.sh und eine NOTES-Zeile. Der war nicht Teil der geprueften Definition of Done und will separat abgenommen werden. Offen dabei geblieben und bewusst nicht entschieden: die letzte Zeile deines Textentwurfs ('Remove a worktree once its pull request is merged') ist nicht eingefuegt, weil sie einen frueheren Loeschzeitpunkt nennt als Zeile 129 des Dispatchers ('only once the ticket is off the board') - sag, welche der beiden gilt."
---

# Die Rollen-Prompts werden im Repository gepflegt, nicht im Heimverzeichnis

## Definition of Done

- [x] core/role/builtin/jaira-dispatcher/SKILL.md und jaira-teamlead/SKILL.md tragen den Text aus ~/.claude/skills, mit jaira-Praefix in name: und in jeder Querverweis-Zeile; core/role/role_test.go TestCrossReferencesCarryThePrefix gruen; go test ./... -race gruen; 'jaira roles install --global --force' legt die sieben jaira-<id> Ordner an und ein diff gegen core/role/builtin zeigt keinen Unterschied; die sieben unpraefixierten Ordner teamlead, dispatcher, role-brainstorm, role-lane, role-pr, role-research, role-tester sind aus ~/.claude/skills entfernt; eine Zeile in core/release/NOTES.md unter ## Unreleased
  proof: core/role/role_test.go TestCrossReferencesCarryThePrefix; go test ./... -race green; diff -r core/role/builtin/jaira-* ~/.claude/skills clean after 'roles install --global --force'; core/release/NOTES.md:18

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-13 21:31 · Alexander Sacharov** — Uebernahme war reine Textkopie: nur 'name:' im frontmatter und die vier Querverweis-Zeilen (/role-lane, /role-tester) bekamen das jaira-Praefix, sonst kein Satz angefasst. Diff gegen ~/.claude/skills belegt genau diese vier Stellen.
Nicht mitgenommen, bewusst: ~/.claude/skills/teamlead/scripts/spawn.sh. Das Skript liegt im Repository unter jaira-dispatcher/scripts/spawn.sh - da gehoert es hin, weil der Dispatcher-Prompt es als 'scripts/spawn.sh' aufruft - und die Repo-Fassung ist die neuere: sie ueberspringt den Port-Offset in einem Repo ohne .env und sendet '/jaira-role-lane'. Die Heimfassung haette beides zurueckgedreht.
Die sieben unpraefixierten Ordner wurden erst geloescht, nachdem 'roles install --global --force' lief und 'diff -r' fuer alle sieben jaira-<id> Ordner keinen Unterschied zeigte.
- **2026-09-13 21:38 · Alexander Sacharov** — critique: Uebernahme gegen die Vorher-Differenz geprueft (aus dem Transkript der Uebernahme-Session rekonstruiert, da die unpraefixierten Ordner geloescht sind). Jede Stelle, die ~/.claude/skills gegenueber core/role/builtin voraus hatte, steht jetzt im Repository: Dispatcher 'These instructions can change while you run', der Tab-Schliess-Absatz, 'A human typing in a worker's tab is not a fault', 'Your own tab is not yours to close', die Drei-Runden-Regel, 'Do not swallow what the human should hear', 'Three lines at the end'; Teamlead 'Close what you started' und 'Bringing one back'. Kein Rest der alten Fassung im Text, keine unpraefixierte /role-* Referenz mehr in core/role/builtin. Die fuenf role-* wurden vom Commit nicht angefasst (letzte Aenderung b9eae90). Kein Befund.
Nebenbefund, nicht Teil des Diffs: die Zeilenzahlen im Ticket-Kontext (dispatcher 170, teamlead 92) waren schon beim Schreiben veraltet - die Heimfassungen hatten 180 und 120 Zeilen, und genau diese Zahlen haben die Repo-Dateien jetzt. Das ist ein Beleg fuer die Vollstaendigkeit, kein Fehler an der Uebernahme.
- **2026-09-13 21:41 · Alexander Sacharov** — testing: bestanden, am Arbeitsbaum auf feat/C9V7ZV-roles-install geprueft, nichts geaendert.
Gates: 'go test ./... -race' gruen, kein einziges FAIL, alle Pakete ok (core/lane 5.4s, internal/cli 25.4s, core/release 1.0s). 'go test ./core/role -race -run TestCrossReferencesCarryThePrefix -count=1' PASS.
Demand, Punkt fuer Punkt am Baum: (1) name: traegt in allen sieben SKILL.md denselben String wie der Ordner; (2) grep ueber core/role/builtin findet keine unpraefixierte Querverweis-Zeile mehr (/role-*, /dispatcher, /teamlead); (3) dispatcher 180 Zeilen, teamlead 120 - genau die Zahlen der geloeschten Heimfassungen, also die vollstaendige Uebernahme; (4) 'go run ./cmd/jaira roles install --global --force' RC=0, Ausgabe '0 written, 8 unchanged, 0 skipped, 0 overwritten', danach 'diff -r' fuer alle sieben jaira-<id> Ordner ohne Unterschied - auch schon vor dem Lauf deckungsgleich, die Installation ist also idempotent; (5) keiner der sieben unpraefixierten Ordner liegt noch in ~/.claude/skills; (6) core/release/NOTES.md:18 traegt die Zeile unter ## Unreleased, eine Zeile, nicht umgebrochen.
Funktion: die Installation selbst war der Funktionstest - das Binary aus diesem Baum schreibt genau den Inhalt von core/role/builtin nach ~/.claude/skills, damit ist die im Goal geforderte Richtung belegt.
Deckungsluecke, als Befund und nicht als Fehler: die Aenderung ist Prompt-Text, und der einzige Test darauf ist TestCrossReferencesCarryThePrefix - der prueft Praefixe, nicht Inhalt. Dass der verschaerfte Text vollstaendig und richtig uebernommen wurde, deckt kein Test ab; das bleibt die Sache der human-Lane.
- **2026-09-13 22:01 · Alexander Sacharov** — Antwort aus der human-Lane am 2026-09-14: der Loeschzeitpunkt eines Worktrees ist "wenn der Ticket von der Tafel ist". Damit gilt jaira-dispatcher/SKILL.md:129 unveraendert weiter, und die letzte Zeile des Textentwurfs ("Remove a worktree once its pull request is merged") bleibt bewusst draussen - es gibt keine zweite, konkurrierende Regel im ausgelieferten Text. Keine Codeaenderung noetig, der Zustand entspricht der Entscheidung bereits. Uebernahme und f3bc433 sind damit abgenommen, das Ticket geht weiter nach review.
