---
id: 01M2EAPXTQKS5GBQS2PVPMF635
title: "Die Rollen-Prompts werden im Repository gepflegt, nicht im Heimverzeichnis"
status: in-progress
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
updated-at: 2026-09-13T21:31:23Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30768
claimed-at: 2026-09-13T21:29:20Z
---

# Die Rollen-Prompts werden im Repository gepflegt, nicht im Heimverzeichnis

## Definition of Done

- [ ] core/role/builtin/jaira-dispatcher/SKILL.md und jaira-teamlead/SKILL.md tragen den Text aus ~/.claude/skills, mit jaira-Praefix in name: und in jeder Querverweis-Zeile; core/role/role_test.go TestCrossReferencesCarryThePrefix gruen; go test ./... -race gruen; 'jaira roles install --global --force' legt die sieben jaira-<id> Ordner an und ein diff gegen core/role/builtin zeigt keinen Unterschied; die sieben unpraefixierten Ordner teamlead, dispatcher, role-brainstorm, role-lane, role-pr, role-research, role-tester sind aus ~/.claude/skills entfernt; eine Zeile in core/release/NOTES.md unter ## Unreleased

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-13 21:31 · Alexander Sacharov** — Uebernahme war reine Textkopie: nur 'name:' im frontmatter und die vier Querverweis-Zeilen (/role-lane, /role-tester) bekamen das jaira-Praefix, sonst kein Satz angefasst. Diff gegen ~/.claude/skills belegt genau diese vier Stellen.
Nicht mitgenommen, bewusst: ~/.claude/skills/teamlead/scripts/spawn.sh. Das Skript liegt im Repository unter jaira-dispatcher/scripts/spawn.sh - da gehoert es hin, weil der Dispatcher-Prompt es als 'scripts/spawn.sh' aufruft - und die Repo-Fassung ist die neuere: sie ueberspringt den Port-Offset in einem Repo ohne .env und sendet '/jaira-role-lane'. Die Heimfassung haette beides zurueckgedreht.
Die sieben unpraefixierten Ordner wurden erst geloescht, nachdem 'roles install --global --force' lief und 'diff -r' fuer alle sieben jaira-<id> Ordner keinen Unterschied zeigte.
