---
id: 01M4ATQ46NZ08TFVQARZ9WSS49
title: "Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo"
status: signoff
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer jaira neu benutzt, findet in der README in einer Zeile pro Skill, was er tut, und im Repo eine Seite, die zeigt, welcher Weg fuer welche Aufgabe passt"
context: |-
  Heute beschreibt README.md (Abschnitt 'Roles: a teamlead, a dispatcher, and the workers', ca. Zeile 615) nur vier Rollen ausfuehrlich; die anderen fuenf stehen nur hinter 'jaira roles list names the rest'.
  Berk am 06.10.: hat nicht gewusst, welche Skills es gibt und ob man mit Claude oder mit den Skills redet. Daraus entstand im Chat eine Liste mit einem Einzeiler pro Skill und eine HTML-Seite (Workflow-Baum, alle Rollen, Abnahmeseite, Einrichten).
  Berk am 07.10.: die Einzeiler sollen so in die README; die Seite soll ins Repo unter docs/.
  Die Seite ist auf Deutsch (so gewuenscht), README und docs/*.md sind Englisch - die Einzeiler kommen deshalb auf Englisch in die README.
  Die Seite wurde ausserhalb des Repos gebaut und liegt bisher unversioniert in /home/berk/git/jAIra/docs/how-to-work-with-jaira.html (auf Branch feat/PAP369-link-lines, falscher Ort).
  Kein Release-Notes-Eintrag: README und docs/ sind nicht im Binary.
definition-of-done: "README.md listet im Rollen-Abschnitt jeden Skill (jaira plus die neun aus 'jaira roles list') mit genau einer Zeile, was er tut; die bestehenden ausfuehrlichen Absaetze bleiben"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-07T09:22:07Z
updated-at: 2026-10-07T09:45:20Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3654451
claimed-at: 2026-10-07T09:38:39Z
outcome-what: "review: docs-only diff read in full; README table, link and reword plus the German tutorial page checked against SKILL.md prompts, NOTES.md, update.go, announce.go, signoff.go and --help; artifact v5 read back and equals the file"
outcome-why: "both DoD items hold and every spot-checked factual claim is backed by source; one imprecise sentence (page line 242, agent-file default) noted as a gap, not a defect"
outcome-resolves: "n/a"
executed-by: opus
review-summary: "Three docs-only commits (6ab980c, a02ca51, c492f35), no Go code touched. README.md:623-634 gains a Skill|What it does table with ten rows: the base skill jaira (noted as shipping in .claude/skills/jaira/, not via roles install) plus the nine ids that jaira roles list prints on 0.3.5; names and one-liners match the installed SKILL.md descriptions. README.md:636-638 links docs/how-to-work-with-jaira.html, says it is German and must be opened in a browser. README.md:662 bullet \"jaira roles list names the rest\" reworded to \"shows the roles the binary you run carries\", since the table now names them all. The four detailed role paragraphs below are untouched. docs/how-to-work-with-jaira.html (624 lines) is the German tutorial: setup, a decision tree A-F plus common tail Z, who-to-talk-to, all-roles table, role hierarchy, lane strip, acceptance page, mode, what is new 0.3.0-0.3.5, rules, CLI list, copy buttons. It carries no doctype on purpose (the artifact host adds one); table{font:inherit} and lang=de on .wrap cover the quirks-mode consequences in a clone. The published artifact FAnMJiZz6dG2hb3it2GM2c, read back in this lane, equals the file apart from the host wrapper (the two CSS rules c492f35 removed are absent in both), so DoD 2 holds."
review-gaps: "Spot-checked the page against the sources it names and found every checked claim backed: roles install exit 3 and \"edited here\" (internal/cli/roles.go:141,151); jaira:local handling and block-from-lanes (core/board/announce.go:306-318, noteFor); \"what is new\" state per person and working tree under home (core/ticket/store.go:638-646); --agent-file (update.go:65-85); hook print blocks once per stop (hook print --help); HERDR_ENV=1, three-returns stop, approval dialog never confirmed, worktree ../.worktrees/<repo>-<slug> and branch feat/<slug> (dispatcher SKILL.md:221,385; scripts/spawn.sh:59,65); python3/gh/glab needs; brainstorm questions 1-3 and user-visible slicing; research writes via jaira note and never moves; tester separates pre-existing failures; pr pushes only when agent-invoked; logbook-summary default last calendar week, no paths/hashes; acceptance: no ids = all human lanes, --name/milestone/branch naming, machine vs person split, ok/no/skipped with no accept button, --force on read-back, wrong-step correction; board key a goes to the terminal lane from any RequiresHumanExit lane (internal/tui/signoff.go:166-196,267-272); jaira with no args opens the home screen (root.go:185); ten lanes on init, lanes add, lanes footer, 0.3.0-0.3.5 rows all in core/release/NOTES.md. Two things to know, neither a reason to send it back: (1) how-to-work-with-jaira.html:242 \"Ohne die Option bleibt er in beiden\" is imprecise: without --agent-file the block goes only into the files that already carry it (announce.go chosenAgentFiles), both only on a board that never chose; after --agent-file agents a plain jaira update keeps AGENTS.md alone. Fixing costs a republish. (2) Visual render not checked in a browser here (no headless browser in this sandbox); the review-check below has the person do it. Unrelated: worktree_diff carries three untracked .jaira/milestones/demo-*.md files that belong to no commit of this ticket; leave them out of any git add."
test-verdict: "pass: docs only, no suite run. README table 10 skills match 'roles list' (0.3.5 build) incl. jaira; link target docs/how-to-work-with-jaira.html exists; HTML 624 lines, tags balanced, 12 #links all have ids; every jaira command+flag on the page exists in --help. Not checked: visual render (no browser, libasound.so.2 missing), DoD 2 artifact v5 identity"
review-verdict: "The diff satisfies both definition-of-done items; no defect found. One imprecise sentence on the page (agent-file default, line 242) is worth a follow-up line, not a return. Accept unless the browser check in review-check shows a layout problem."
review-check: "1. cd /home/berk/git/.worktrees/jAIra-howto-docs && git log origin/master..HEAD --oneline  -> exactly three lines, each starting docs(9WSS49).  2. sed -n 623,638p README.md  -> a table \"Skill | What it does\" with 10 rows (jaira, /jaira-teamlead, /jaira-dispatcher, /jaira-role-lane, /jaira-role-brainstorm, /jaira-role-research, /jaira-role-tester, /jaira-role-acceptance, /jaira-role-pr, /jaira-role-logbook-summary), then a line linking docs/how-to-work-with-jaira.html.  3. jaira roles list  -> the nine ids printed are the nine /jaira-* rows of that table, no id missing or extra.  4. sed -n 662p README.md  -> \"- jaira roles list shows the roles the binary you run carries.\"  5. Open file:///home/berk/git/.worktrees/jAIra-howto-docs/docs/how-to-work-with-jaira.html in a browser  -> heading \"Mit jaira arbeiten\", a row of 11 pill links under it, the \"Alle Rollen\" table in the same sans font as the body (not a serif default), every command block with a \"Kopieren\" button.  6. Click the pill \"Die Lanes\"  -> the page jumps to a strip of 13 lanes in order backlog..blocked, with human and signoff in orange.  7. Click \"Kopieren\" on the jaira self upgrade block  -> the button reads \"Kopiert\" (or \"Markiert\" if the browser refused clipboard) and the clipboard holds the command.  8. Open https://claude.ai/artifact/FAnMJiZz6dG2hb3it2GM2c  -> the same page; its footer names internal/cli/update.go and core/board/announce.go, same as the local file.  9. After the branch is pushed, open the README on GitHub and scroll to \"Roles: a teamlead, a dispatcher, and the workers\"  -> the table renders as a table and the docs/how-to-work-with-jaira.html link opens the HTML source (expected; the line says so)."
---

# Die README nennt jede Rolle in einer Zeile, und das Tutorial liegt im Repo

## Definition of Done

- [x] README.md listet im Rollen-Abschnitt jeden Skill (jaira plus die neun aus 'jaira roles list') mit genau einer Zeile, was er tut; die bestehenden ausfuehrlichen Absaetze bleiben
  proof: README.md:623-634 Tabelle Skill|What it does, 10 Zeilen (jaira + die neun aus 'jaira roles list' von 0.3.5); Absaetze darunter unveraendert
- [x] docs/how-to-work-with-jaira.html ist committet und deckt sich mit der veroeffentlichten Seite
  proof: docs/how-to-work-with-jaira.html; dieselbe Datei als Version 6 nach https://claude.ai/artifact/FAnMJiZz6dG2hb3it2GM2c veroeffentlicht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-07 09:31 · BeMuCa** — critique (pass 1, whole diff of 6ab980c read): three findings, fix is clear for each.
1) README.md:658 bullet "`jaira roles list` names the rest: ..." is obsolete - the new table names every role. It is a pointer, not a detailed paragraph the DoD protects. Reword to "`jaira roles list` shows the roles the binary you run carries."
2) Nothing links to docs/how-to-work-with-jaira.html (git grep: no reference outside .jaira). Goal says a newcomer finds it. Add one English line after the table (README.md:635), in the style of README.md:680, saying the page is German and must be opened in a browser (GitHub shows .html as source).
3) docs/how-to-work-with-jaira.html:1 starts at <meta charset>: no <!doctype html>, no <html lang="de">. claude.ai adds both when publishing; a file opened from a clone gets neither -> quirks mode + wrong language. Add both above line 1, then republish the same file to the artifact (FAnMJiZz6dG2hb3it2GM2c) so DoD 2 still holds. HTML parser ignores a doctype inside an already-open document and merges lang onto the existing <html>, so the published copy is unaffected.
Let stand on purpose: README table repeating the four detailed bullets (DoD keeps the bullets); German page duplicating the README table (person asked for both); the data-theme CSS that only the artifact host uses (needed for the published copy, which must match); hand-maintained page drifting by version (footer names the refresh steps). Content correctness (e.g. jaira pull, cp of SKILL.md) not judged here - that is review/testing.
- **2026-10-07 09:32 · BeMuCa** — Antwort auf critique Runde 1: L1 (README 'names the rest' umformuliert) und L2 (Link auf die Seite hinter der Tabelle) wie vorgeschlagen. L3 anders geloest: kein <!doctype>/<html lang> in die Datei - die Artifact-Vorgabe des Publishers verbietet eigene doctype/html/head/body-Tags ausdruecklich, er wrappt selbst. Stattdessen die sichtbare Quirks-Folge behoben: table { font: inherit } (Quirks setzt Tabellen-Schrift auf initial) und lang=de auf dem Inhalts-Wrapper. Datei als Version 4 veroeffentlicht, DoD 2 haelt.
- **2026-10-07 09:36 · BeMuCa** — critique pass 2 (read only round 1's three findings and a02ca51): all three closed, nothing new. L1 README 'jaira roles list shows the roles the binary you run carries' and L2 link line after the table (README.md:636-638) as proposed. L3 closed by the different fix in the 09:32 note: both harms the finding named are gone - table { font: inherit } undoes the quirks-mode table reset of font/line-height (font-size set right after, text-align set on th/td), and lang=de on .wrap covers every visible element incl. the JS copy button. Keeping the file in the publisher's input format (the Artifact host wraps the page in its own skeleton) is a defensible reason not to add a doctype; not reopened. DoD 2 republish taken from the 09:32 note, not re-checked here.
- **2026-10-07 09:38 · BeMuCa** — Nach optimize: die zwei wirkungslosen CSS-Zeilen, die optimize nur wegen des Neu-Veroeffentlichens stehen liess (gap: 0 auf ol.chain, reduced-motion-Regel ohne Transitions auf der Seite), entfernt und als Version 5 veroeffentlicht. Ein Republish kostet einen Aufruf; kein Grund, toten Code zu behalten.
- **2026-10-07 09:39 · BeMuCa** — testing: pass. Compared README.md:623-634 names with jaira-master roles list (9 + jaira): equal. Link target exists. Python HTMLParser: no unclosed/mismatched tags, 12 in-page anchors all resolve. Commands checked against --help: hook print, lanes add, logbook, milestone create, note, pull, resume, roles install/list (--project --global), self upgrade, show --for-lane --json, update --agent-file, create --goal/--context/--dod, next --per-lane --json, set, board, init: all exist. NOT verified: visual render (no headless browser, libasound.so.2 missing); DoD 2 (file == artifact v5) not checkable here.
- **2026-10-07 09:45 · BeMuCa** — Nach review, Ticket bleibt in signoff: den Gap aus review-gaps (Zeile 242) behoben. Ohne --agent-file schreibt jaira in die Dateien, die den Block schon tragen, und nur wenn keine ihn hat in beide (core/board/announce.go:354 chosenAgentFiles) - die Seite sagte 'bleibt in beiden'. Zeile 236 im selben Sinn auf 'und/oder' gezogen. Sonst nichts geaendert; als Version 6 veroeffentlicht. Dieser eine Satz ist nach dem Review geaendert und von keinem Worker mehr gelesen worden.
