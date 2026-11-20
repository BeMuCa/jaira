---
id: 01M2HZKCQ48BPQYN1YPTXB97NJ
title: "Ein Board sagt, ob Arbeit durch einen Pull Request kommt oder direkt auf den Hauptzweig"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Wer allein an einem Repository arbeitet, bekommt keine Rolle vorgesetzt, die ihn zu einem Pull Request schickt, den er nicht will - und ein Team, das Pull Requests benutzt, bekommt die Regel weiterhin ueberall, wo ein Agent sie liest."
definition-of-done: "Ein Board haelt fest, wie Arbeit bei ihm ankommt - durch einen Pull Request oder direkt auf dem Hauptzweig - und die Rollen richten sich danach. Nachgestellt an zwei Boards mit verschiedener Einstellung: auf dem einen schiebt die Rolle den Zweig und haelt an, auf dem anderen fuehrt sie die Arbeit ohne Pull Request zu Ende."
tags:
  - cli
blocked-by:
  - 01M7WWYB5XZMD01DZY3MJG3REW
related:
  - 01M7WWYB5XZMD01DZY3MJG3REW
commits: []
created-at: 2026-09-15T07:30:45Z
updated-at: 2026-11-20T23:02:51Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
context: |-
  13VMA8 hat die Regel 'ein Agent schiebt den Zweig und haelt an, der Mensch macht den Pull Request auf' am 2026-09-15 verbindlich gemacht - in CLAUDE.md, AGENTS.md, dem README und in den ausgelieferten Rollen-Prompts. Fuer ein Team ist das richtig. Fuer einen Menschen, der allein an einem Repository arbeitet und auf den Hauptzweig committet, ist es schlicht falsch: die Rolle haelt ihn an und reicht ihm eine 'gh pr create'-Zeile fuer etwas, das er nie wollte.

  Alex am 2026-09-15: ein Projekt kann git benutzen und trotzdem ohne Pull Requests arbeiten; das Board soll sagen, was bei ihm gilt.

  Wo die Einstellung NICHT hingehoert, und das ist an zwei Faellen aus denselben zwei Tagen gelernt:
  - ~/.jaira/settings.json gilt pro Rechner. Genau daran ist 9ET6NC gescheitert: ein 'remote' fuer alle Boards hat jedes andere Board lahmgelegt.
  - 'git config' gilt pro Klon. Dort liegt 'jaira.forge' richtig, WEIL es persoenlich ist - in welchen Remote DU schiebst, ist deine Sache. 'Dieses Projekt arbeitet mit Pull Requests' ist nicht persoenlich: gilt es fuer das Team, gilt es fuer alle. Es muss mit dem Repository reisen.

  Der Mechanismus dafuer ist schon da: der jaira-Block in CLAUDE.md und AGENTS.md wird ERZEUGT. 'jaira update' kann den Absatz schreiben, der zur Einstellung des Boards passt.

  Damit steht die Frage, die dieses Ticket beantworten muss: die Rollen-Prompts werden fuer alle gleich ausgeliefert und koennen nicht pro Projekt verschieden sein. Wie erfaehrt ein ausgelieferter Prompt, dass dieses Board es anders haelt? Entweder er liest die Einstellung zur Laufzeit, oder der erzeugte Block gilt vor ihm und der Prompt sagt das ausdruecklich. Das ist die Entscheidung, nicht die Umsetzung.

  Verwandt, aber nicht dasselbe: PWAF22 will, dass 'jaira update' das ganze Setup nachzieht. Wenn dieses Ticket dem Befehl etwas zu schreiben gibt, bekommt PWAF22 einen Grund mehr zu existieren.
question: ""
---

# Ein Board sagt, ob Arbeit durch einen Pull Request kommt oder direkt auf den Hauptzweig

## Definition of Done

- [x] Ein Board haelt fest, wie Arbeit bei ihm ankommt - durch einen Pull Request oder direkt auf dem Hauptzweig - und die Rollen richten sich danach. Nachgestellt an zwei Boards mit verschiedener Einstellung: auf dem einen schiebt die Rolle den Zweig und haelt an, auf dem anderen fuehrt sie die Arbeit ohne Pull Request zu Ende.
  proof: Nachgestellt 21.11. an zwei Boards (solo: .jaira/pr=none, team: ohne): roles list sagt none/person; pr-gate laesst none bei done durch (internal/cli/prgate_test.go TestPRGateLetsANoneBoardLandOnlyDoneWork); jaira-role-pr/SKILL.md:280 'Land it' nur bei none, sonst push + Zeile zurueck
- [x] Die Einstellung reist mit dem Repository, nicht mit dem Rechner und nicht mit dem Klon: wer klont, erbt sie. Nachgestellt an einem frischen Klon.
  proof: frischer Klon eines solo-Boards mit committeter .jaira/pr=none: nach jaira init sagt roles list 'none', jaira update schreibt den Absatz; die Datei liegt in .jaira/ (core/board/arrival.go:52 PRModeFor), nicht in settings.json oder git config
- [ ] Ein ausgelieferter Rollen-Prompt richtet sich danach, obwohl er fuer alle gleich ausgeliefert wird - und wie er davon erfaehrt, steht begruendet im Ticket, nicht nur im Code.
- [ ] Ein Board ohne Einstellung verhaelt sich wie heute: durch einen Pull Request. Wer nichts tut, merkt nichts.
- [ ] Eine Zeile in core/release/notes/XB97NJ.md (die Notiz-Dateien ersetzen seit AR97K9 core/release/NOTES.md).
- [ ] Eingabeklassen der Datei .jaira/pr per Tabellentest: fehlt, leer, nur Kommentar, CRLF, Gross-/Kleinschreibung ('None'), Leerzeichen um den Wert, Tippfehler ('nnone'), zwei Zeilen - nur exakt 'none' heisst kein PR, alles Unbekannte bleibt person.
- [ ] Der Block ist byte-gleich zu heute, wo .jaira/pr fehlt oder person sagt; ein Board mit none/agent bekommt den passenden Absatz durch jaira update, und jaira validate meldet den Block als veraltet, sobald die Datei sich aendert.
- [ ] Bei none landet die Arbeit nur durch das Tor: jaira roles pr-gate laesst erst durch, wenn jedes Ticket des Zweigs in done ist, liest .jaira/pr vom Basis-Zweig, und das Landen ist ein fast-forward push - ein abgelehnter (nicht fast-forward) Push wird nie mit --force wiederholt.

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] core/board/arrival.go: PR-Modus person|agent|none parsen (PRModeIn/PRModeFor) + arrivalSection(mode); settings.PRFor/PRIn delegieren, PRNone dazu
- [ ] announce.go: noteFor/announceInAgentFile/NoteIsCurrent haengen arrivalSection(PRModeFor(root)) an; person -> leer, Block unveraendert
- [ ] roles.go: roles list zeigt none mit Satz; pr-gate: none laesst bei allen Tickets in done durch (pr_when ignoriert), json traegt pr
- [ ] Rollentexte: jaira-role-pr Abschnitt 'Land it' fuer none (gate, merge base, push HEAD:<branch>, nie --force); teamlead/dispatcher/worker je ein Satz; shipped.sha256 neu
- [ ] Tests: Tabellentest Eingabeklassen, Block byte-gleich ohne Datei, Absatz bei none/agent, NoteIsCurrent-Drift, pr-gate none; Release-Note; Nachstellen an zwei Boards + frischem Klon

## Progress
- **2026-09-15 07:34 · Alexander Sacharov** — Alex am 2026-09-15 zur Richtung und zum Zeitpunkt: der Weg ueber den erzeugten Block in CLAUDE.md ist der, den er ausprobiert haben will - nicht weil er sicher der richtige ist, sondern weil er sich ohne Umbau probieren laesst. Das Ticket geht in Release 0.2.2, und was sich im Gebrauch als falsch herausstellt, wird danach schrittweise nachgezogen.

Fuer die Plan-Lane heisst das: die Entscheidung 'erzeugter Block statt Laufzeit-Einstellung im Prompt' ist vorgegeben und nicht neu aufzurollen. Offen bleibt, WIE der ausgelieferte Prompt dem Block weicht - ob er ausdruecklich sagt, dass der Block vor ihm gilt, oder ob der Block den betreffenden Absatz des Prompts ersetzt.

Nicht als Freibrief lesen: 'spaeter nachziehen' gilt fuer die Form, nicht fuer Kriterium 4. Ein Board ohne Einstellung muss sich vom ersten Tag an wie heute verhalten - wer nichts tut, darf nichts merken.
- **2026-09-15 07:34 · Alexander Sacharov** — Am 2026-09-15 auf Alex' Wunsch als zu besprechen markiert: die erste Loesung wird gemeinsam festgelegt, bevor das Ticket in eine Lane geht. Die Frage steht im question-Feld, damit sie auf dem Board sichtbar ist und nicht nur in einer Notiz, die niemand oeffnet.
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 22:50 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 22:51 · Alexander Sacharov** — Worker 2026-11-20, vor dem Plan - offene Entscheidungen. Lage: JG3REW (done, Branch feat/jg3rew, NOCH NICHT auf master) hat .jaira/pr (person|agent) eingefuehrt; die Rollen lesen es zur Laufzeit ueber 'jaira roles list'. Das ist genau die Board-Datei, die dieses Ticket braucht, und ein Laufzeit-Weg statt des im September gewaehlten erzeugten Blocks.
F1 Wo steht die Einstellung?
 1) (empfohlen) dritter Wert in .jaira/pr: person|agent|none - 'none' = kein PR, Arbeit direkt auf den Hauptzweig. Eine Datei, ein Begriff; Branch baut auf feat/jg3rew auf oder wartet auf dessen Merge.
 2) eigene Datei .jaira/flow (pr|direct) - unabhaengig von JG3REW, aber zwei Dateien fuer ein Thema.
F2 Wie erfaehrt der ausgelieferte Prompt davon?
 1) (empfohlen) beides: der Prompt liest es zur Laufzeit (roles list, wie JG3REW) UND jaira update schreibt den Absatz im CLAUDE.md/AGENTS.md-Block passend (PR-Absatz bzw. 'direkt auf master, kein PR'). Agenten ohne Rollen-Skill lesen den Block.
 2) nur erzeugter Block, Prompt sagt 'der Block gilt vor mir' - Entscheidung vom 15.09., aber dann gibt es zwei Mechanismen fuer .jaira/pr.
 3) nur Laufzeit, Block bleibt fest.
Den woertlichen Absatz entwerfe ich selbst, er wird im Review sichtbar.
- **2026-11-20 22:53 · Alexander Sacharov** — decision: Teamlead 20.11. 23:00: 1/1 - F1: dritter Wert 'none' in .jaira/pr. F2: beide Wege - die Rolle liest zur Laufzeit (roles list), und jaira update schreibt den passenden Absatz in den CLAUDE.md/AGENTS.md-Block. Baut auf feat/jg3rew auf (PR #149, noch nicht gemergt). Modus bleibt autonom: Alex schlaeft, Teamlead will bis signoff durchgezogen; weitere Fragen an den Teamlead.
- **2026-11-20 22:55 · Alexander Sacharov** — decision (Plan, Begruendung fuer DoD 3): Der ausgelieferte Prompt erfaehrt es auf zwei Wegen, beide aus derselben Datei .jaira/pr. (a) Laufzeit: jaira-role-pr und teamlead lesen 'jaira roles list --json | jaira json get pr' wie seit JG3REW - darum gilt es sofort nach dem Merge der Datei, ohne dass jemand update laufen laesst, und eine Rolle kann nicht mit einem veralteten Block auseinanderlaufen. (b) Block: jaira update schreibt fuer agent/none einen Absatz 'Wie Arbeit ankommt' in den CLAUDE.md/AGENTS.md-Block, fuer Agenten ohne jaira-Skill. Bei person/fehlender Datei schreibt der Block NICHTS Neues - byte-gleich zu heute (DoD 4): die Regel 'Agent pusht und haelt an' steht heute in den Rollen und im local-Teil dieses Repos, nicht im erzeugten Block, und ein neuer Absatz bei allen Boards wuerde jedes Board als 'veraltet' melden. Parsing liegt in core/board (settings importiert board, nicht umgekehrt), settings delegiert. none heisst: ein Zweig pro Ticket bleibt (Worktrees brauchen ihn), aber er landet ohne PR: jaira-role-pr laesst pr-gate (alle Tickets in done; pr-when ist bei none bedeutungslos) laufen, merged <remote>/<branch> in den Zweig und schiebt 'git push <remote> HEAD:<branch>' - nur fast-forward, nie --force. Die Annahme in done ist die Entscheidung, wie bei agent.
