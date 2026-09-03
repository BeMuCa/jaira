---
id: 01M0YT4RXZHYN3XT90RA2MM32Y
title: "Eine Regel fuer fertige Tickets: archive oder sync"
status: done
ready: true
creator: BeMuCa
goal: "Vier Doku-Stellen sagen dasselbe darueber, wo ein fertiges Ticket hingeht"
context: "Widerspruch, entstanden beim Mergen von PR #5. Die done-Lane, SKILL.md und docs/AGENTS.md sagen 'jaira archive, sobald gepusht'. Der generierte Block in CLAUDE.md sagt 'jaira sync'. Beide stempeln inzwischen die Commit-Liste; Unterschied ist nur der Zielordner und dass sync die Endlane verlangt. Begruendung fuer sync stand in Issue #4: 'archive reads as abandoned' und ein Tagesordner pro Person fuer die Stundenbuchung - das ist die Praxis des Beitragenden, nicht zwingend Berks. Berk entscheidet die Regel, dann werden vier Stellen gleichgezogen."
definition-of-done: "done-Lane, SKILL.md, docs/AGENTS.md und der generierte Block nennen denselben Weg; der jeweils andere Befehl wird als Sonderfall benannt, nicht verschwiegen"
blocked-by: []
commits:
  - 7173a454706067eadd531b1d1d378dec76affe0d
  - 5170c917295b35a83823d2c353bdf13b4539d5ab
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-08-26T10:34:07Z
updated-at: 2026-09-03T15:48:25Z
claimed-by: EE-3NX6GL3-3492202
claimed-at: 2026-08-27T15:56:35Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Die vier Stellen (Skill, README-Absatz, README-Befehlsliste, docs/AGENTS.md) sagen jetzt dieselbe Regel: fertige Arbeit ins Logbuch, nicht bearbeitete Tickets ins Archiv, aus jeder Lane. Kein Code."
outcome-why: "Die einzige erzwungene Differenz zwischen beiden Befehlen ist, dass logbook ein Ticket vor der Terminal-Lane ablehnt (syncout.go:88-105, jetzt logbook.go). Und der datierte Ordner bekommt mit der Tagesstatistik (SPDWGH) einen Leser. Beides sagt: das Logbuch ist der Weg fuer Fertiges."
outcome-resolves: "DoD: 'eine Regel, an allen Stellen gleich' - grep nach 'archived predecessor|Archive it before' findet nichts mehr; jede Stelle nennt beide Befehle und wann welcher."
review-summary: "Reiner Doku-Commit 7173a45: Skill-Abschnitt 'Taking a ticket off the board' neu geschrieben mit der Regel als erster Zeile und beiden Befehlen im Codeblock; README-Absatz nach dem signoff-Bild und Befehlsliste ergaenzt; docs/AGENTS.md:80 'archived' -> 'logged or archived'. archive --help war schon in b93a85c dran."
review-gaps: "Nur Text geaendert, nichts zu falten oder zu loeschen. Gesucht: weitere Stellen, die 'archive' fuer fertige Tickets empfehlen - grep 'archived predecessor|Archive it before|jaira archive <id>' ueber README, docs/, SKILL.md, announce.go: keine mehr. Stehen gelassen: die globale Skill-Kopie unter ~/.claude/skills/jaira/SKILL.md ist Berks eigene Datei und sagt noch 'archive'; sie liegt nicht im Repo."
review-check: |-
  1. jaira archive --help lesen: der zweite Absatz sagt, archive stempelt Commits und ist fuer Tickets, die nicht bearbeitet werden; fuer Fertiges verweist er auf jaira logbook.
  2. grep -n 'logbook' README.md - der Absatz nach dem signoff-Bild nennt 'jaira logbook <id>' fuer das abgenommene Ticket und 'jaira archive <id>' fuer nicht bearbeitete; die Befehlsliste hat beide Zeilen.
  3. grep -n 'Taking a ticket off the board' -A 12 .claude/skills/jaira/SKILL.md - erste Zeile des Abschnitts ist die Regel in Fettdruck.
  4. grep -n 'archived predecessor' README.md docs/*.md .claude/skills/jaira/SKILL.md - kein Treffer.
review-verdict: "Angenommen. Gleiche Einschraenkung wie zuvor: dieselbe Sitzung. Geprueft: grep nach den alten Formulierungen leer; go test ./core/board/ gruen (announce_test prueft den Block-Text, der hier nicht geaendert wurde). Die Entscheidung selbst - Logbuch fuer Fertiges - ist Berks Frage 'we answered that right?' vom 27.08.; sie folgt aus der einzigen erzwungenen Differenz (Terminal-Lane-Gate) und aus dem Leser des datierten Ordners (Tagesstatistik SPDWGH). Wenn er anders entscheidet, sind es dieselben vier Stellen noch einmal."
---

# Eine Regel fuer fertige Tickets: archive oder sync

## Definition of Done

- [x] done-Lane, SKILL.md, docs/AGENTS.md und der generierte Block nennen denselben Weg; der jeweils andere Befehl wird als Sonderfall benannt, nicht verschwiegen
  proof: SKILL.md 'Taking a ticket off the board', README 'Accepting it moves the ticket to done' + command list, docs/AGENTS.md:80, archive --help (b93a85c) - alle vier sagen: fertig -> logbook, nicht bearbeitet -> archive

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

