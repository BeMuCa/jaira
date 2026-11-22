---
id: 01M2N12WVABVRYM7C923WATJDQ
title: Der GitLab-Arm der pr-Rolle benutzt ein abgekuendigtes Flag und laesst --head weg
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Auf einem GitLab-Board oeffnet /jaira-role-pr den Merge Request im richtigen Projekt und aus dem richtigen Fork-Branch, mit Flags, die glab noch kennt."
context: |-
  In core/role/builtin/jaira-role-pr/SKILL.md steht der GitLab-Zweig der Zielrepository-Leiter. Die review-Lane von A3R6YC hat ihn am echten glab geprueft und zwei Fehler gefunden:

  1. 'glab mr create --target-project <pfad>' ist in glab 1.114 abgekuendigt. Nachstellen: 'glab mr create --target-project foo/bar --title x --description y' - erste Ausgabezeile ist 'Flag --target-project has been deprecated, Use --repo instead.' Der Befehl legt nichts an, er bricht danach ab, weil hier kein GitLab-Remote konfiguriert ist.
  2. Im Fork-Fall fehlt '--head <fork>'. 'glab mr create --help' zeigt im Beispielblock 'glab mr create --repo upstream/project --head your-namespace/project ...' - zwei Flags, wo SKILL.md nur eines nennt. Ohne --head sucht glab den Quell-Branch im Zielprojekt und findet ihn nicht.

  Der GitHub-Arm ist nachgewiesen richtig und wurde auf dem Fork sashasoft90/jaira durchgespielt; nur GitLab ist ungeprueft. Dieses Board liegt auf GitHub, deshalb wurde A3R6YC mit diesem Rest angenommen statt zurueckgeschickt - der Fehler faellt hier niemandem auf und trifft erst das erste Board auf GitLab.

  Dazu offen und NICHT nachgestellt: die review-Lane vermutet, dass 'gh pr create' ohne '--head' aus einem Fork-Clone haengen bleibt. Wer dieses Ticket nimmt, kann das gleich mitpruefen - es ist dieselbe Stelle.

  Die Aenderung ist reiner Text in einer Datei.
definition-of-done: "core/role/builtin/jaira-role-pr/SKILL.md nennt im GitLab-Zweig '--repo' statt '--target-project' und ergaenzt '--head <namespace>/<projekt>' fuer den Fork-Fall; 'glab mr create --help' belegt beide Flags; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen"
tags:
  - docs
  - cli
blocked-by: []
related:
  - 01M7WWYB5XZMD01DZY3MJG3REW
commits: []
created-at: 2026-09-16T11:54:25Z
updated-at: 2026-11-22T18:51:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-426824
claimed-at: 2026-11-22T18:44:01Z
outcome-what: "jaira-role-pr SKILL.md: create lines use glab --repo/--head and gh --head; open-request lookup lists by bare branch and filters the head owner (gh headRepositoryOwner, glab source_project_id); release notes; shipped.sha256"
outcome-why: "glab 1.114 deprecates --target-project and needs --head for a fork MR; gh pr list --head owner:branch listed nothing, so a fork never saw its open PR"
outcome-resolves: WATJDQ
review-summary: "none - the fix sits at the one place the create command is spelled; both forges get one rule (always name the head) instead of a fork branch, and the person-path gh line now matches what pr-round.sh already passes on the agent path. path_with_namespace checked against a real glab repo view -F json."
review-gaps: "none - two changed command lines, one paragraph saying why, one note line; nothing duplicated or unused"
test-verdict: "pass: go vet ./... RC=0, go test -race ./core/role/... ./core/shipped/... ./core/release/... RC=0; DoD 1 verified in the tree (SKILL.md:440,446,473, notes/WATJDQ.md); glab 1.114: 'mr create --repo <host>/up/proj --head <host>/me/proj' parses with no deprecation and reaches the API (fails only on the unreachable host, recovery file removed), '--target-project' prints 'Flag --target-project has been deprecated, Use --repo instead.'; no live MR (glab not authenticated, no GitLab project here); the change writes nothing, so no two-writer case"
hunt-verdict: "fail: bughunt pre-existing SKILL.md:390 gh pr list --head <owner>:<branch> returns [] (folded into DoD 2, same subject: fork head in the pr role)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Der GitLab-Arm der pr-Rolle benutzt ein abgekuendigtes Flag und laesst --head weg

## Definition of Done

- [x] core/role/builtin/jaira-role-pr/SKILL.md nennt im GitLab-Zweig '--repo' statt '--target-project' und ergaenzt '--head <namespace>/<projekt>' fuer den Fork-Fall; 'glab mr create --help' belegt beide Flags; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen
  proof: core/role/builtin/jaira-role-pr/SKILL.md:446,473 (--repo + --head; glab mr create --help lines 41,54,63 on glab 1.114); core/release/notes/WATJDQ.md (NOTES.md is closed); go test ./core/role/... ./core/shipped/... ok
- [x] SKILL.md 'Does it already have one open': gh pr list filters the head owner via --json headRepositoryOwner (gh pr list --head takes no <owner>: prefix and returns [] for it), glab mr list filters source_project_id against origin's project - a same-named branch in another fork does not count as this branch's request
  proof: core/role/builtin/jaira-role-pr/SKILL.md:390-409; gh pr list --repo BeMuCa/jaira --head feat/4539mz --jq select owner==sashasoft90 -> #189, owner==nobody -> empty; glab mr list --repo gitlab-org/cli -F json --jq select(source_project_id==45049979) -> MR URLs

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-22 18:44 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:46 · Alexander Sacharov** — decision: release note goes to core/release/notes/WATJDQ.md, not core/release/NOTES.md as the DoD says - NOTES.md is closed history since AR97K9 (CLAUDE.md). decision: gh pr create (person path) gets --head <owner-of-origin>:<branch> too - same fork-clone class the context names; pr-round.sh (agent path) already passed --head, so only the SKILL.md line was missing it. Not reproduced live: no fork clone with a TTY-less gh here; gh pr create --help itself says --head skips the fork/push prompt. decision: --head is named always, not only in a fork - in a non-fork it equals the target and glab/gh accept it, one line instead of two branches. Teamlead (jaira-d0) said 22.11.: do not wait for blocked-by CRASZS - it touches core/pr and the cockpit, this ticket only SKILL.md.
- **2026-11-22 18:46 · Alexander Sacharov** — blocked-by CRASZS removed on the teamlead's word (22.11.): different files, do not wait.
- **2026-11-22 18:49 · Alexander Sacharov** — bughunt: pre-existing: core/role/builtin/jaira-role-pr/SKILL.md:390 - class: the open-request check passes gh pr list --head in <owner>:<branch> form, which gh does not support for list (only for create), so it never finds an open pull request from a fork; example: PR #189 (head feat/4539mz from the fork) is found with --head feat/4539mz and missed with --head sashasoft90:feat/4539mz, so the role takes 'Nothing listed', writes a description and runs gh pr create instead of going to Answering review comments; blame: 8f537cd329 (2026-09-16, before the ticket); ran: gh pr list --repo BeMuCa/jaira --head sashasoft90:feat/4539mz --state all (rc 0, []) vs --head feat/4539mz (rc 0, #189); gh 2.83.1 help: '("<owner>:<branch>" syntax not supported)'
- **2026-11-22 18:49 · Alexander Sacharov** — bughunt: done at d115846a8 - 1 finding (0 change, 1 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/WATJDQ.md core/role/builtin/jaira-role-pr/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 18:50 · Alexander Sacharov** — decision: bughunt's pre-existing finding (gh pr list --head owner:branch finds nothing) folded into this ticket as DoD 2 - same module, same subject (naming the fork head), a few lines. glab list gets the matching source_project_id filter.
