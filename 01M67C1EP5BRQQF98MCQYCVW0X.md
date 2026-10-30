---
id: 01M67C1EP5BRQQF98MCQYCVW0X
title: "Ein Skill schaut, wie die Person arbeitet, und schlaegt vor, ob und wie jaira in ihren Ablauf passt"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Eine Person, die jaira neu einsetzt, bekommt aus ihrer eigenen Arbeitsweise einen Vorschlag: welche Boards, welche Lanes, oder ob jaira fuer sie gar nicht passt."
context: |-
  Idee von berkc am 30.10.2026: 'should we add a skill, where when started, the agent checks how the user is working and then checks if and how jaira should be used in his workflow? ... let an agent brainstorm a little about that skill.'
  Ueberschneidet sich mit 7BWHCN (jaira lernt Arbeitsmuster aus Sessions/Commits) und PP24H1 (Interview-Skill baut Board und Milestones).
  Quellen fuer 'wie die Person arbeitet': Claude-Code-Transkripte (~/.claude/projects), git log, vorhandene Skills/Hooks.
definition-of-done: "Brainstorm: Abgrenzung zu 7BWHCN/PP24H1, Quellen, Ergebnisform; Empfehlung ob ein eigener Skill"
tags:
  - cli
blocked-by: []
related:
  - 01M5GFAWXX2K6E8MEBGZ7BWHCN
  - 01M55T0M58DDQ7BVE2MPPP24H1
commits: []
created-at: 2026-10-30T21:39:17Z
updated-at: 2026-10-30T21:42:44Z
updated-by: berkc
---

# Ein Skill schaut, wie die Person arbeitet, und schlaegt vor, ob und wie jaira in ihren Ablauf passt

## Definition of Done

- [ ] Brainstorm: Abgrenzung zu 7BWHCN/PP24H1, Quellen, Ergebnisform; Empfehlung ob ein eigener Skill

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 21:42 · berkc** — Brainstorm YCVW0X (Agent, 30.10.2026) - Abgrenzung, Quellen, Ergebnisform, Empfehlung.

Kurz: Ein Skill /jaira-role-fit. Er liest lokal Spuren der Arbeitsweise (git log, jaira-Metriken, Struktur der Claude-Code-Sitzungen), stellt dann wenige Interviewfragen und endet mit einem von drei Urteilen: passt so / passt mit diesen Boards und Lanes / passt nicht. Geschrieben wird nichts ohne Ja.

Was heute da ist (geprueft mit jaira 0.5.7):
- Rollen sind Markdown unter core/role/builtin/<name>/SKILL.md, teils mit Python-Helfer (logbook-summary/collect.py, acceptance/build.py); 'jaira roles install' schreibt sie in den Skills-Ordner. Ein Helfer-Skript braucht also keine neue Runtime.
- 'jaira boards create' gibt es in 0.5.7 NICHT. Ein Board ist heute ein Verzeichnis .jaira/boards/<name>/ mit Lane-Dateien, order und description. 'jaira boards create <name> --type' bringt erst Y2C9QN (steht in signoff, nicht auf master).
- Vorhandene Bausteine fuer den Vorschlag: 'jaira lanes add|remove|move|use', 'jaira lanes market', 'jaira lanes default --lanes', 'jaira lanes template --board'. 'jaira lanes add --help' zeigt kein --board (nicht geprueft, ob es vom Elternbefehl geerbt wird).
- Sitzungsprotokolle: ~/.claude/projects/<cwd mit - statt />/<uuid>.jsonl, eine JSON-Zeile pro Eintrag (type, timestamp, cwd, gitBranch, message.content[] mit tool_use name/input). Subagenten liegen in .../<uuid>/subagents/*.jsonl.
- Auf diesem Rechner: 107 Dateien in 14 Projekten, davon 97 Subagenten-Dateien. Claude Code raeumt alte Sitzungen weg -> die Datenlage ist duenn, das Interview ist kein Extra, sondern Pflicht.
- Nur die Namen ausgezaehlt (kein Text gelesen): Skill-Aufrufe disciplined-implementation 9, exploration-router 7, jaira-worker 5, rigorous-reasoning 4; Agent-Aufrufe meist ohne subagent_type (80). Der Lane-Zweck steckt also eher in Skill-Namen und Agent-'description' als im Agententyp.
- ~/.jaira/metrics.csv hat pro Lauf lane, loops, verdict, tool_uses - zeigt, welche Lanes wirklich gelaufen sind und welche Schleifen sich drehen.
- ~/.claude/settings.json: Hooks (PreToolUse, Stop, SubagentStop ...) und ~97 Skills - Namen sagen viel (gsd-*, quality-run = critique/optimize/test).
- Vorbild ausserhalb: gsd-profile-user / gsd-user-profiler bauen ein Profil aus Sitzungsnachrichten - genau das, was wir NICHT wollen (liest Text).

Ablauf des Skills:
1. Einwilligung pro Quelle einholen: "Darf ich git log dieses Repos lesen? Die Struktur deiner Claude-Sitzungen fuer dieses Repo? Fuer alle Repos?" Ohne Ja wird die Quelle uebersprungen.
2. Repo-Spuren: git log (Praefixe feat/fix/refactor/docs, Anteil Commits mit Ticket-Id, Branch-Namen, PR-Merges ja/nein, wie viele Autoren), CLAUDE.md/AGENTS.md, ob .jaira/ schon da ist, 'jaira boards', 'jaira lanes'.
3. Sitzungs-Spuren ueber ein mitgeliefertes collect.py: liest NUR type, timestamp, gitBranch, tool_use.name, input.subagent_type, input.skill. Agent-'description' nur per Stichwortabgleich gegen Lane-Ids (critique, review, test, optimize, plan ...) zu einer Zahl verdichtet, nie ausgegeben. Ergebnis: Folgen wie "implementieren -> Kritik -> Optimieren -> Test" mit Haeufigkeit.
4. jaira-Spuren: metrics.csv (welche Lanes laufen, Schleifenzahl, Abbrueche), Tickets die lange in einer Lane stehen.
5. Befund zeigen: 5-8 Zeilen, nur Muster und Zahlen, jede mit Quelle ("in 9 von 12 Sitzungen nach Implementierung ein Review-Skill").
6. Interview, eine Frage pro Runde (Technik aus PP24H1): nur die Luecken, die die Spuren offen lassen - arbeitest du allein oder im Team, gibt es Arbeit ausser Code (Video, Recherche, Doku) die ein eigenes Board braucht, wo soll ein Mensch zwingend drauf schauen.
7. Urteil, eins von drei: (a) jaira passt nicht - z.B. ein Mensch, ein Kurzprojekt, keine mehrstufige Arbeit, oder schon ein anderes System (GSD .planning/) das dasselbe tut; Begruendung in einer Zeile, Ende. (b) passt mit Standard-Board. (c) passt mit diesen Aenderungen.
8. Vorschlag als fertige Befehle, die die Person bestaetigt oder streicht, z.B.:
   jaira init
   jaira lanes add critique; jaira lanes add optimize; jaira lanes add testing
   jaira lanes remove pre-process
   jaira lanes default --lanes backlog,todo,in-progress,critique,testing,signoff,done   (wenn es fuer alle Repos der Person gelten soll)
   jaira boards create research --type research   (erst nach Y2C9QN; bis dahin: "lege .jaira/boards/research/ an" als Schritt mit 'jaira lanes template --board')
   jaira roles install
9. Nach Ja: Befehle ausfuehren, 'jaira lanes' und 'jaira boards' zur Kontrolle zeigen. Nichts committen (Lane-Dateien reisen mit dem naechsten Code-Commit bzw. per 'jaira share').

Datenschutz-Grenzen (gehoeren woertlich in die SKILL.md):
- Nur lokal, nichts verlaesst den Rechner ausser der Zusammenfassung im eigenen Chat.
- Nie Gespraechstext, Prompts, Tool-Ein-/Ausgaben, Dateiinhalte aus Sitzungen zitieren oder speichern - nur Namen, Zeitstempel, Zaehlungen, Reihenfolgen.
- Standard nur das Projektverzeichnis des aktuellen Repos; andere Projekte nur nach eigenem Ja, und dann nur aggregiert ohne Projektnamen.
- Der Befund landet nicht auf dem Board und nicht in einer Datei, ausser die Person verlangt es.
- Andere Agenten (Codex, opencode ...): kein Transkript-Format bekannt -> Schritt 3 entfaellt, git + Interview reichen. Damit bleibt die CLI-Regel "jeder bash-faehige Agent" gewahrt.

Verhaeltnis zu den anderen Tickets - Empfehlung: zusammenlegen, ein Skill.
- 7BWHCN und YCVW0X sind dieselbe Idee; YCVW0X fuegt nur das Urteil "passt nicht" hinzu. 7BWHCN in YCVW0X aufgehen lassen; sein DoD-Punkt "an diesem Repo findet er critique/optimize/testing nach in-progress" wird Abnahmetest von YCVW0X.
- PP24H1 aufteilen: der Board/Lane-Teil ist Schritt 6-8 hier (gleiche Interviewtechnik). Was bleibt - Milestones, Epics, Roadmap fuer ein neues Vorhaben - ist ein anderer Moment (Planung eines Vorhabens, nicht Einrichtung eines Arbeitsablaufs) und bleibt als kleinerer PP24H1, der den Fit-Skill am Ende aufrufen kann oder umgekehrt.
- Abhaengigkeit: Y2C9QN (boards create --type) zuerst mergen, sonst muss der Skill Board-Verzeichnisse von Hand beschreiben.

Alternative 1: Alles (YCVW0X + 7BWHCN + PP24H1) in eine Rolle mit Modi board/milestone.
- Kosten: eine SKILL.md von ~250+ Zeilen, zwei Ausloeser in einem Prompt; Milestone-Fragen stoeren jemanden, der nur seinen Ablauf einrichten will. Spart eine Rolle, verliert Schaerfe.
Alternative 2: Kein Skill, sondern 'jaira suggest' im Go-Binary, das Transkripte selbst parst.
- Kosten: hoch - bindet das Binary an das interne JSONL-Format von Claude Code (aendert sich ohne Ankuendigung), widerspricht "jeder bash-faehige Agent", waechst gegen "kleiner als paca". Vorteil nur: deterministisch und testbar. Ablehnen.
(Billigste Variante, falls gar kein neuer Skill: 3 Fragen im jaira-teamlead, wenn das Board frisch ist. Kosten fast null, aber ohne Spuren - nur Interview.)

Empfehlung: Eigener Skill /jaira-role-fit (Name offen), SKILL.md ~120 Zeilen im Ton der anderen Rollen + collect.py (~80 Zeilen, nur Strukturfelder), ausgeliefert unter core/role/builtin. 7BWHCN darin aufgehen lassen, PP24H1 auf Milestones/Epics verkleinern, nach Y2C9QN bauen.

Offene Fragen an berkc:
1. Darf der Skill Agent-'description'-Felder per Stichwort abgleichen, oder strikt nur Tool-/Skill-Namen? (Ohne description sind 80 von ~120 Agent-Aufrufen stumm.)
2. Nur dieses Repo oder alle Projekte unter ~/.claude/projects (mit Ja)? Alle zeigt den persoenlichen Stil, ein Repo den Projektstil.
3. Soll ein Ergebnis auch 'jaira lanes default' (fuer alle kuenftigen Repos) setzen duerfen, oder nur dieses Board?
4. 7BWHCN als Duplikat schliessen und PP24H1 verkleinern - einverstanden?
5. Ausliefern (core/role/builtin) oder erst als persoenlicher Skill ausprobieren?
6. Name: jaira-role-fit, jaira-role-onboard oder jaira-role-setup?
7. Soll "passt nicht" auch fremde Systeme erkennen (GSD .planning/, GitHub Issues) und dann zur Koexistenz raten statt Ersatz?
