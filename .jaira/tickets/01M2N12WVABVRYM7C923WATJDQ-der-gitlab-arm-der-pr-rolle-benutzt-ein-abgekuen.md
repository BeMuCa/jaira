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
updated-at: 2026-11-22T18:47:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-426824
claimed-at: 2026-11-22T18:44:01Z
outcome-what: "jaira-role-pr SKILL.md: GitLab create lines use --repo and --head <origin path>; the person-path gh pr create line names --head <owner>:<branch>; release note; shipped.sha256 regenerated"
outcome-why: "glab 1.114 deprecates --target-project and without --head looks for the source branch in the target project, so a fork MR fails"
outcome-resolves: WATJDQ
---

# Der GitLab-Arm der pr-Rolle benutzt ein abgekuendigtes Flag und laesst --head weg

## Definition of Done

- [x] core/role/builtin/jaira-role-pr/SKILL.md nennt im GitLab-Zweig '--repo' statt '--target-project' und ergaenzt '--head <namespace>/<projekt>' fuer den Fork-Fall; 'glab mr create --help' belegt beide Flags; eine Zeile unter ## Unreleased in core/release/NOTES.md; go test ./core/role/... gruen
  proof: core/role/builtin/jaira-role-pr/SKILL.md:446,473 (--repo + --head; glab mr create --help lines 41,54,63 on glab 1.114); core/release/notes/WATJDQ.md (NOTES.md is closed); go test ./core/role/... ./core/shipped/... ok

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
