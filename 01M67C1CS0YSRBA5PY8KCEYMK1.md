---
id: 01M67C1CS0YSRBA5PY8KCEYMK1
title: "Ein Board-Typ 'Projektidee': Idee rein, aehnliche Ideen recherchieren, Nutzen pruefen (abwaehlbar), weitere Schritte"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Eine Projektidee laeuft durch ein eigenes Board: Recherche nach Aehnlichem, optional eine Nutzen-Pruefung, und was sonst noetig ist, bevor gebaut wird."
context: |-
  Wunsch von berkc am 30.10.2026: 'We should create a Project Idea board ... take in a Projekt idea, then we do steps like researching similar ideas, validating the usefullness of the idea (this is a lane that we can skip, its asked while creating the ticket) and other steps necessary (start an agent to brainstorm this board idea a little)'.
  Werkzeug: Board-Typen aus Y2C9QN (core/lane/boardtypes/<typ>/: description, order, Lane-Dateien) - Ordner statt Go-Code.
  Abwaehlbare Lane = requires-option wie brainstorm/pre-process (core/lane/lane.go:120).
  Verwandt: P5ABVR (Ideen-Board fuer den Hackathon).
definition-of-done: "Brainstorm: Lanes mit Zweck, Prompt-Skizze und welche abwaehlbar sind"
tags:
  - cli
blocked-by: []
related:
  - 01M5659DPY9FK72XVMZFY2C9QN
  - 01M5GFB5KWABVZVWWDVNP5ABVR
commits: []
created-at: 2026-10-30T21:39:15Z
updated-at: 2026-10-30T21:44:00Z
updated-by: berkc
---

# Ein Board-Typ 'Projektidee': Idee rein, aehnliche Ideen recherchieren, Nutzen pruefen (abwaehlbar), weitere Schritte

## Definition of Done

- [ ] Brainstorm: Lanes mit Zweck, Prompt-Skizze und welche abwaehlbar sind

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
- Nur lokal, nichts verlassen den Rechner ausser der Zusammenfassung im eigenen Chat.
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
- **2026-10-30 21:44 · berkc** — KORREKTUR: Die Notiz von 21:42 direkt davor ('Brainstorm YCVW0X ... /jaira-role-fit') gehoert zu YCVW0X, nicht hierher - zwei parallele Agenten haben dieselbe Scratch-Datei benutzt. Sie steht auf YCVW0X auch richtig. Fuer CEYMK1 ignorieren; es gilt die folgende.

Brainstorm CEYMK1 (Agent, 30.10.2026): Vorschlag Board-Typ 'projektidee' - eine Idee rein, am Ende Go/No-Go mit einseitigem Brief.

Dateien: core/lane/boardtypes/projektidee/ mit description, order und je einer eigenen Lane-Datei (keine eingebaute Lane passt, siehe unten).
order: idee, schaerfen, recherche, nutzen, zuschnitt, entscheidung, abgelegt - 7 Lanes, davon 4 agentisch.

1. idee - Mensch, kein Agentenschritt.
Zweck: Rohidee festhalten, eine Zeile reicht. Hier wird auch die Option 'nutzen' an- oder abgewaehlt.
Liefert: title, context (kommen von create).

2. schaerfen - agentisch, Tier strong, output-produces: goal.
Zweck: aus dem Satz eine pruefbare Idee machen.
Prompt-Skizze: Schreib die Idee in drei Zeilen: welches Problem, fuer wen, was danach anders ist. Liste die drei Annahmen, an denen sie scheitert, die riskanteste zuerst. Ist die Notiz zu duenn, frag per human-Frage statt zu raten. Das goal ist ein Satz.

3. recherche - agentisch, Tier strong, output-produces: prior-art. Nicht abwaehlbar.
Zweck: gibt es das schon (Produkte, Open-Source, aufgegebene Versuche)?
Prompt-Skizze: Arbeite nach /jaira-role-research mit der Frage 'Gibt es das schon, und was fehlt den vorhandenen?'. Nenne 3-5 echte Projekte, je: was es kann, was fehlt, Quelle, Datum, Konfidenz. Schluss in einer Zeile: existiert schon / Luecke ist X / unklar. Kein Projekt gefunden ist ein Befund, keine Erlaubnis.

4. nutzen - agentisch, Tier strong, output-produces: nutzen-urteil, requires-option: nutzen. ABWAEHLBAR.
Zweck: Nutzen pruefen, bevor Zeit hineingeht.
Prompt-Skizze: Pruefe die riskanteste Annahme aus schaerfen gegen Belege: wer hat das Problem, wie oft, was tut er heute stattdessen (ein Workaround ist der beste Beleg fuer echten Schmerz). Belege aus Issues, Foren, Bewertungen der Konkurrenz aus recherche - keine erfundenen Interviews. Urteil stark/schwach/unbelegt, dazu der Beleg, der es am ehesten kippen wuerde.

5. zuschnitt - agentisch, Tier strong, output-produces: plan (eingebautes Feld, Checkliste).
Zweck: der kleinste baubare Schnitt, der die riskanteste Annahme testet.
Prompt-Skizze: 3-7 Plan-Schritte; jeder wird spaeter ein Backlog-Ticket, also je Titel, goal und ein dod-Satz in einer Notiz. Was ausdruecklich NICHT gebaut wird. Aufwand S/M/L und ein Abbruchkriterium. Stack nur nennen, wenn recherche ihn nahelegt. Die Gegenrede (warum es scheitert) gehoert hier hinein, nicht in eine eigene Lane.

6. entscheidung - Mensch, requires-human-exit: true, output-produces: entscheidung.
Zweck: Go / No-Go / zurueck (zu recherche oder schaerfen mit neuer Frage).
Der Mensch setzt 'jaira set <id> entscheidung="go: <ein Satz>"'. Ein Agent kann hier nicht weiter, nur hineinlegen.

7. abgelegt - terminal, kein Agentenschritt. Eigene Datei, weil die eingebaute done-Lane requires-commits und requires-outcome hat (.jaira/lanes/done.md) - eine Idee hat keine Commits. Go und No-Go landen beide hier; ein No-Go mit Begruendung ist genauso wertvoll (man prueft dieselbe Idee nicht zweimal).

Endprodukt: das Ticket selbst ist der einseitige Brief - goal + prior-art + nutzen-urteil + plan + entscheidung, lesbar mit 'jaira show <id>'. Keine eigene Datei, kein neues Format. Bei Go legt der Teamlead oder die Person je Plan-Schritt ein Ticket auf dem Bau-Board an: 'jaira create --board <code-board> --follows <idee-id> --goal ... --dod ...'. Das ist Handarbeit, weil es keine boarduebergreifende Weitergabe gibt (8NHD8Q, P5ABVR); zuschnitt kann die fertigen create-Zeilen als Notiz liefern.

Technisch geprueft (Code gelesen, nicht ausgefuehrt):
- Eigene output-Felder (prior-art, nutzen-urteil, entscheidung) gehen ohne Go-Code: core/gate/gate.go fieldFilled faellt bei unbekannten Feldern auf das Frontmatter-Skalar zurueck; geschrieben mit 'jaira set <id> prior-art=...'.
- Keine Lane des Typs darf requires-specified tragen (wie todo), sonst verlangt das Board dod und assignee, die eine Idee nicht hat.
- Id ohne eigene Datei = eingebaute Lane (core/lane/boardtypes.go auf origin/feat/Y2C9QN-boards-create); hier passt keine, also 7 eigene Dateien.

Verworfen:
A) Das Code-Board mit Option brainstorm plus /jaira-role-research benutzen. Spart den Typ, aber: keine Nutzen-Pruefung, done verlangt Commits, und Ideen verstopfen das Board, auf dem gebaut wird.
B) Volle Startup-Pipeline (Markt, Wettbewerb, Monetarisierung, Tech-Spike, Pitch, Kritik als je eigene Lanes, 10+). Jede Lane ist eine Uebergabe mehr; fuer eine Person zu schwer und groesser als paca. Wettbewerb steckt in recherche, Kritik in zuschnitt.

Offene Fragen an berkc:
1. 'nutzen' standardmaessig an oder aus? Und wie wird beim Anlegen gefragt: heute hat 'jaira create' kein --option; es gibt nur 'jaira dod <id> --option nutzen' danach und 'jaira lanes --options' als Board-Vorgabe. Braucht create ein --option-Flag bzw. eine Frage im TUI-Dialog (Code-Aenderung, eigenes Ticket)?
2. Wohin geht ein Go: Tickets auf das Code-Board (per Teamlead, von Hand) oder baut man auf demselben Board weiter?
3. No-Go eigener Terminal-Lane 'verworfen' oder nur das Feld entscheidung in 'abgelegt'? (Ob zwei terminale Lanes erlaubt sind, nicht geprueft.)
4. Mit P5ABVR zusammenlegen? Das Hackathon-Ideen-Board waere dann nur eine Stufe davor: Thema -> N Tickets auf einem projektidee-Board.
5. Soll recherche auch die eigenen Repos und alten Ideen/Logbuecher durchsuchen, nicht nur das Web?
6. Lane-Prompts auf Englisch wie die eingebauten, oder Deutsch?
