---
id: 01M5ZJ555685ZXARWNEF5QR36R
title: Die Abnahme-Rolle fuehrt das Interview im Chat; die Seite ist nur noch auf Wunsch
status: pre-process
ready: true
creator: BeMuCa
goal: "jaira-role-acceptance fragt standardmaessig im Chat ab - je Ticket ID, Einzeiler, wie pruefen, Empfehlung, dann 'accept X / explain X / fix X' - und baut die statische Seite nur mit --page."
context: |-
  Was schief ist: die Abnahme-Seite veraltete mit jedem Rebuild, der Read-back lief nie; das Interview im Chat war fuer 17 Tickets schneller als jede Seite (Feedback L324/L335/O24, 27.10.2026).
  Heute: core/role/builtin/jaira-role-acceptance/SKILL.md baut eine Seite und liest Markierungen zurueck; der Chat-Weg steht nirgends.
  Was gut lief: pro Ticket ID, ein Satz was, wie pruefen, Empfehlung, Frage oder Entscheidung; der Mensch antwortet 'accept 1TFFG2' oder stellt die Frage.
  Zusammenhang: Q63EEV (Cockpit-Pane je Ticket, Frage offen in human) - dieses Ticket aendert nur die Rolle, nicht das Cockpit.
definition-of-done: SKILL.md beschreibt das Chat-Interview als Default mit dem Format je Ticket und den drei Antworten
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:15Z
updated-at: 2026-10-31T02:34:12Z
assignee: berkc
updated-by: berkc
claimed-by: berk-ca-3129429
claimed-at: 2026-10-31T01:27:42Z
---

# Die Abnahme-Rolle fuehrt das Interview im Chat; die Seite ist nur noch auf Wunsch

## Definition of Done

- [ ] SKILL.md beschreibt das Chat-Interview als Default mit dem Format je Ticket und den drei Antworten
- [ ] --page baut die Seite wie bisher; ohne --page wird keine Datei geschrieben
- [ ] Die statische Aufzeichnung der Abnahme entsteht auch im Chat-Weg (gleicher Pfad wie heute)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 23:45 · berkc** — Frage berkc 31.10.2026: Das Interview in einer neuen Session bzw. als eigener Interviewer-Subagent fuehren? Vorteil: der Kontext der Hauptsession bleibt schlank, und das Interview laeuft parallel zur normalen Arbeit. Tokens: das Interview kostet gleich viel, aber nicht im teuren, langen Hauptkontext - jede Antwort dort wuerde den ganzen Verlauf neu mitschicken. Empfehlung: ja, eigener Agent/Tab, der nur die Tickets in human/signoff liest; die Hauptsession bekommt am Ende eine Zeile pro Ticket (angenommen/zurueck/Folgeticket).
- **2026-10-31 01:25 · berkc** — Maintainer 31.10.2026: umsetzen. Der Skill heisst interReview (Interview-Abnahme im Chat). Danach geht der Maintainer mit der Session die Tickets in review/signoff damit durch.
- **2026-10-31 02:34 · berkc** — pre-process decisions (planner):
- Name: role dir jaira-role-interreview -> invoked as /jaira-role-interreview. Not 'interReview': the dir name IS the command (core/role/role.go:12), every role carries jaira-role- (TestCrossReferencesCarryThePrefix, wantRoles), frontmatter name must equal dir (role_test.go TestBuiltinsAreTheShippedRoles), lowercase avoids case trouble on Windows/macOS filesystems.
- Rename, not a second role: git mv jaira-role-acceptance -> jaira-role-interreview (SKILL.md, build.py, template.html, report.html). One role with two modes is what the DoD describes (--page builds the page as before). A second role would leave /jaira-role-acceptance building a page by default, against the goal.
- Old installs: install() (core/role/install.go:62) only writes the files of current Builtins and never deletes, so a renamed dir would leave every user's old jaira-role-acceptance working on with page-by-default. Fix: keep core/role/builtin/jaira-role-acceptance/ with a short forwarding stub SKILL.md only (forwards to /jaira-role-interreview <args> --page). On 'jaira update --roles' / 'roles install' an untouched old SKILL.md matches a shipped hash -> Updated to the stub without --force; an edited one is Skipped as always. Old build.py/template.html/report.html stay on disk in the old dir, unused, harmless. The stub also keeps every v0.4.0-fixture test (core/role/shipped_test.go, internal/cli/roles_test.go, skillsnudge_test.go, updatesetup_test.go) valid unchanged.
- Shipped hashes: append-only. Append the lines TestEveryInstalledFileIsAShippedVersion prints (new interreview keys, stub, every edited SKILL.md). Never remove the old acceptance lines. Do not run scripts/shipped-history.sh (git rename detection with a new file at the old path is ambiguous).
- DoD 2 vs 3: 'no file without --page' means no page files (data.json/page.html/db.json/decisions.json/url.txt). The one file the chat path writes is the record report.md at the same path as today, ~/.claude/acceptance/<repo>/<name>/report.md, appended per interview so earlier records stay.
- review-lane tickets have no review-check yet (the review lane produces it): the chat interview takes 'how to check' from review-check when present, otherwise from the DoD proofs + test-verdict, and says which.
- accept: human -> next_lane (route continues, model review still runs); review and signoff -> the terminal lane, because the person is the sign-off and the maintainer wants to clear review/signoff with it. Uses --force only for the person-lane refusal; dry-run first, any other refusal shown verbatim, no move.
- Called by an agent (dispatcher/worker/teamlead): no interview in the agent's long context (berkc note 30.10): answer with the line to paste into a new session. Interview ends with one line per ticket for the main session.
