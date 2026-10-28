---
id: 01M5B0PCJW8KCF3XFCTSWVYG9H
title: Tickets und Logbook nennen keine anderen Projekte mehr
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Kein Text unter .jaira/ (Tickets, Logbook, Archiv, Ticket-Refs) nennt mehr den Namen eines anderen Projekts des Maintainers."
context: |-
  Was falsch ist: Ticket-Texte und Notes nennen andere, private Projekte des Maintainers beim Namen. jaira ist oeffentlich auf GitHub, das ist ein Datenleck.
  Die Namen stehen NICHT hier, damit dieses Ticket das Leck nicht wiederholt. Liste: ~/.jaira/foreign-names (eine Zeile je Name, nur auf Alex' Rechner).
  Wo sie stehen (Stand 19.10.2026): upstream/master: gut 20 Dateien unter .jaira/tickets/, .jaira/logbook/ und .jaira/archive/ (git grep -ciE -f ~/.jaira/foreign-names upstream/master). Refs: 9 Tickets unter refs/jaira/tickets/* (aktueller Inhalt jedes Refs mit git show <ref>:<datei> pruefen).
  Anlass: Alex 19.10.2026 - 'nie Projektnamen verwenden'. Auf CRDH6B wurde eine Note schon per sed anonymisiert (Commit 54583386 auf feat/crdh6b).
  Entschieden von Alex: nur der heutige Text wird bereinigt, mit normalen Commits und einem PR. Die git-Historie (master und Refs) wird NICHT umgeschrieben, kein force-push.
  Ausnahme: die CLI kann Notes nicht aendern. Deshalb ist hier ausdruecklich erlaubt, Dateien unter .jaira/ von Hand zu aendern - nur das Ersetzen der Namen, sonst nichts.
  Ersatz: neutral beschreiben, was gemeint war ('ein anderes Projekt', 'ein Epic in einem anderen Projekt', 'das Video-Board' usw.), damit der Satz verstaendlich bleibt.
  Offen und zuerst klaeren: wie eine von Hand geaenderte Ticket-Datei auf ihren Ref kommt (refsync, z.B. core/refsync/refsync.go Record), ohne dass der naechste jaira-Aufruf den alten Text vom Ref zurueckschreibt. Pruefen mit git show auf den Ref nach dem Schreiben.
  Im Commit und PR selbst die Namen auch nicht nennen.
definition-of-done: git grep -iE -f ~/.jaira/foreign-names auf dem Branch findet unter .jaira/ nichts mehr
tags:
  - docs
blocked-by: []
related: []
commits:
  - c81e83d46174653d667e27d5737c90456985c1e7
created-at: 2026-10-19T21:22:16Z
updated-at: 2026-10-26T22:40:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30382
claimed-at: 2026-10-26T21:58:04Z
question: ""
mode: conversational
outcome-what: "Namen anderer Projekte des Maintainers, davon abgeleitete Kuerzel und der Firmen-Host sind ersetzt: in 25 Dateien unter .jaira/, in 9 Dateien ausserhalb (.planning, zwei Tests, ein Skript, ein Rollen-Skill) und auf 9 Ticket-Refs auf upstream, jeweils mit einem fast-forward-Commit"
outcome-why: "jaira ist oeffentlich, die Namen waren ein Datenleck. Die Historie bleibt unangetastet (so entschieden), bereinigt wird nur der heutige Text"
outcome-resolves: "Kein heutiger Text unter .jaira/, auf den Ticket-Refs oder in den genannten Dateien nennt mehr eines der Projekte aus ~/.jaira/foreign-names"
review-summary: "Reiner Textersatz, kein Verhalten geaendert: in 39 Dateien (22 unter .jaira/logbook und .jaira/archive, 8 Ticket-Kopien unter .jaira/tickets, 9 ausserhalb: .planning, zwei Testkommentare, ein Docstring in scripts/migrate-tickets.py, ein Beispiel-Host in der PR-Rolle) sind Namen anderer Projekte des Maintainers, das Board-Kuerzel und der Firmen-Host durch neutrale Umschreibungen ersetzt ('ein anderes Projekt', 'test project board', git.example.de). Auf upstream sind die Ticket-Refs per fast-forward-Commit bereinigt; die Historie bleibt unveraendert, wie entschieden."
review-gaps: |-
  note: outcome-what nennt 25 Dateien unter .jaira/, der Diff aendert 30 (22 Logbook/Archiv + 8 Ticket-Kopien); die Arbeit ist vollstaendig, nur die Zahl stimmt nicht
  note: noch kein Commit - commits ist leer und die Ticket-Datei dieses Tickets ist untracked; ohne einen Commit mit WVYG9H im Betreff wird der Zug nach done verweigert
  note: die 8 Ticket-Kopien unter .jaira/tickets behalten ihr updated_at; ein Klon mit alter Kopie kann den alten Text wieder auf den Ref schreiben (in der Note vom 19.10. schon als Restrisiko festgehalten)
  note: zwischen 21.10. und 26.10. sind auf 4 Refs neue Nennungen entstanden - ohne eine Pruefung beim Schreiben kommt das Leck wieder
test-verdict: "pass: go vet ./... RC=0, go build ./cmd/jaira RC=0, go test ./... -race RC=0 (auf upstream/master 254b9085); DoD 1 grep .jaira exit 1, DoD 2 0 Treffer auf 139 upstream-Refs nach Push, DoD 3 30 Dateien +86/-86 und kein Frontmatter-Feld geaendert, DoD 5 grep ausserhalb exit 1; gebautes Binary: jaira list exit 0, jaira show auf jedes geaenderte Ticket ok, Frontmatter aller 30 Dateien intakt"
review-verdict: "Der Diff erfuellt alle fuenf DoD-Punkte, unabhaengig nachgeprueft: grep mit ~/.jaira/foreign-names unter .jaira/ und im ganzen Baum 0 Treffer, Board-Kuerzel und Firmen-Host 0 Treffer, 136 Ticket-Refs auf upstream 0 Treffer, eigene Ticket-Datei 0 Treffer, go test ./internal/cli/ ./internal/tui/ ./core/role/ gruen. Jede Aenderung ist eine Namensersetzung. Kein Blocker; offen ist nur der Commit."
review-round: 1
review-check: "1. cd ~/projects/.worktrees/jaira-wvyg9h  2. grep -rliE -f ~/.jaira/foreign-names .jaira ; git grep -liE -f ~/.jaira/foreign-names -- .  -> beide Befehle geben nichts aus  3. git diff --word-diff=plain -- . | grep -E '\\[-' | less  -> jede Zeile tauscht nur einen Namen gegen eine neutrale Umschreibung, der Satz bleibt lesbar  4. git fetch upstream '+refs/jaira/*:refs/chk/*' ; for r in $(git for-each-ref --format='%(refname)' refs/chk/); do git show $r:$(git ls-tree --name-only $r | head -1); done | grep -ciE -f ~/.jaira/foreign-names  -> 0 ; danach git for-each-ref --format='%(refname)' refs/chk/ | xargs -n1 git update-ref -d  5. go test ./internal/cli/ ./internal/tui/ ./core/role/  -> dreimal ok"
---

# Tickets und Logbook nennen keine anderen Projekte mehr

## Definition of Done

- [x] git grep -iE -f ~/.jaira/foreign-names auf dem Branch findet unter .jaira/ nichts mehr
  proof: grep -rniE -f ~/.jaira/foreign-names .jaira -> exit 1 (keine Treffer, inkl. der drei neuen Namen); ebenso Praefix, Workspace-Label und Firmen-Host
- [x] Der aktuelle Inhalt jedes Refs unter refs/jaira/ enthaelt keinen Namen aus der Liste, auf upstream geprueft nach dem Push
  proof: nach push-refs.sh + push-refs2.sh (Alex, 26.10.): git fetch upstream refs/jaira/*, git show auf den aktuellen Inhalt jedes Refs | grep -iE -f ~/.jaira/foreign-names -> 0 Treffer auf 139 Refs; 9 Refs per fast-forward-Commit bereinigt, keine Historie umgeschrieben
- [x] Jeder ersetzte Satz bleibt verstaendlich; kein Feld ausser dem Namen ist veraendert (git diff --word-diff zeigt nur Namensersetzungen)
  proof: git diff --word-diff=plain -- .jaira: 25 Dateien, 78/78 Zeilen, jede Aenderung ist eine Namens-/Host-/Praefix-Ersetzung; Paare in fester Liste, kein Frontmatter-Feld ausser Text beruehrt
- [x] Commit-Nachrichten, Ticket-Notes und PR-Text dieses Tickets nennen keinen der Namen
  proof: Ticket-Datei und Ref dieses Tickets: grep -ciE -f ~/.jaira/foreign-names -> 0; vorgeschlagene Commit-Zeile 'docs(WVYG9H): replace names of other projects...' -> 0; PR-Text schreibt der Maintainer, dort gilt dieselbe Liste
- [x] Die Dateien ausserhalb .jaira/ (.planning/HANDOFF.md, internal/cli/staleclaim_test.go, internal/tui/home_refs_test.go, scripts/migrate-tickets.py) nennen keinen Namen aus der Liste mehr; go test ./internal/cli/ ./internal/tui/ bleibt gruen
  proof: git grep -niE -f ~/.jaira/foreign-names -- . ':!.jaira' -> nur noch Binaer-Rauschen in docs/img/demo-chain.gif (kein Listen-Name); go test ./internal/cli/ ./internal/tui/ ./core/role/ ok

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] ~/.jaira/foreign-names um die drei weiteren Projektnamen ergaenzen (Rechner-lokal, nicht im Repo)
- [x] Fundstellen-Liste: Listen-Namen, abgeleitete Kuerzel (Compose-/Worktree-Praefix, zwei Seitennamen, Workspace-Label) und Firmen-Host - auf dem Branch und auf jedem Ref
- [x] Ersetzungen unter .jaira/ auf dem Branch, Satz fuer Satz von Hand (python-Skript mit exakten Alt/Neu-Paaren im Scratchpad); word-diff pruefen -> DoD 1, 3
- [x] Die vier Dateien ausserhalb .jaira/ (drei aus der Entscheidung plus ein Testkommentar in internal/tui/home_refs_test.go); go test der beiden Pakete -> DoD 5
- [x] Ref-Inhalte: dieselben Paare auf den aktuellen Inhalt jedes betroffenen Refs anwenden, je Ref einen neuen Commit oben drauf bauen (hash-object/mktree/commit-tree), lokal noch nicht pushen; diff zeigen
- [x] Nach Freigabe: git push upstream <neu>:<ref> --force-with-lease=<ref>:<alt> je Ref, update-ref lokal, danach git show auf upstream pruefen -> DoD 2
- [x] Commit-Nachricht, Notes und PR-Text gegen die Liste greppen -> DoD 4

## Progress
- **2026-10-19 21:22 · Alexander Sacharov** — Teamlead 19.10.2026: der Worktree startet auf dem lokalen master 7f008dec, upstream/master ist weiter (36bde2df, release 0.5.4). Vor dem ersten Commit auf upstream/master aufsetzen und dort greppen - sonst fehlen die neuesten Logbook-Dateien.
- **2026-10-19 21:22 · Alexander Sacharov** — parent session: jaira-20
- **2026-10-19 21:25 · Alexander Sacharov** — decision (Ref-Weg, aus dem Code geklaert): ein Ticket-Ref ist eine Commit-Kette, jeder Commit hat genau eine Datei <ULID>.md als ganzen Baum (core/gitref/gitref.go refWrite/commitFile). jaira schreibt einen Ref nur aus dem Outbox, und nur fuer Tickets, die hier geschrieben wurden - ein Ref, fuer den nichts in der Queue liegt, wird von keinem jaira-Aufruf ueberschrieben. Weg also ohne CLI: hash-object + mktree + commit-tree -p <alter Ref> + git push upstream <neu>:<ref> --force-with-lease=<ref>:<alt>, danach update-ref lokal. Das ist ein normaler Commit oben drauf (fast-forward), keine Historie umgeschrieben. jaira pull scheidet aus: es claimt und aendert damit assignee auf dem Ref.
- **2026-10-19 21:25 · Alexander Sacharov** — decision: updated_at bleibt unveraendert (DoD 3: nur Namen ersetzt). Restrisiko: ein Klon, der eine alte Kopie derselben Datei auf der Platte hat und nach einem fetch darauf schreibt, schickt den alten Text mit gueltigem Lease wieder hoch (CheckFresh vergleicht nur updated_at). Betrifft nur Tickets, die auch als Datei auf master liegen; mit dem Merge dieses PRs ist auch die master-Kopie sauber. Pruefung nach dem Push: git show auf jeden Ref.
- **2026-10-19 21:25 · Alexander Sacharov** — question (vor dem Plan, an Alex): Umfang. Gezaehlt auf master+Refs: 1) drei weitere Namen eigener Projekte, die NICHT in ~/.jaira/foreign-names stehen (je 1 Fundstelle, eine davon auch auf einem Ref), dazu ein Compose-Praefix und ein Herdr-Workspace-Label, die vom Hauptnamen abgeleitet sind. 2) der GitLab-Host der Firma, 11 Stellen. 3) drei Dateien ausserhalb .jaira/ nennen einen Listen-Namen (.planning/HANDOFF.md, internal/cli/staleclaim_test.go, scripts/migrate-tickets.py). Optionen je Punkt: F1 a) mitnehmen und die drei Namen in ~/.jaira/foreign-names ergaenzen [empfohlen] b) nur die Liste. F2 a) Host neutral ersetzen ('ein firmeneigenes GitLab') [empfohlen] b) Host bleibt. F3 a) mitnehmen, DoD-Punkt dazu [empfohlen] b) eigenes Ticket c) bleibt.
- **2026-10-21 07:47 · Alexander Sacharov** — decision (Alex, 2026-10-21, ueber Teamlead) zum Umfang: F1 a) die drei weiteren Projektnamen samt abgeleitetem Compose-Praefix und Workspace-Label mitnehmen und die drei Namen in ~/.jaira/foreign-names ergaenzen. F2 a) den Firmen-GitLab-Host neutral ersetzen ('ein firmeneigenes GitLab'). F3 a) die drei Dateien ausserhalb .jaira/ (.planning/HANDOFF.md, internal/cli/staleclaim_test.go, scripts/migrate-tickets.py) hier mitnehmen, DoD-Punkt dazu. Frage damit geschlossen.
- **2026-10-21 07:48 · Alexander Sacharov** — Teamlead jaira-46, 2026-10-21: alter Worker-Tab (w3:tJ3) geschlossen - Session nicht mehr erreichbar (ListAgents kennt sie nicht, workers-Tool: Kontext unbekannt). Neuer Worker in derselben Worktree, Parent jaira-46. Der Ref dieses Tickets war lokal kaputt (leeres git-Objekt nach WSL-Absturz), von upstream neu geholt.
- **2026-10-21 07:48 · Alexander Sacharov** — parent session: jaira-46
- **2026-10-21 07:52 · Alexander Sacharov** — plan (pre-process): Ersatz nach festen Regeln, damit die Pruefung nach DoD 3 mechanisch geht. Hauptname -> 'ein anderes Projekt' in der passenden Beugung; Video-Repo -> 'das Video-Projekt'; die drei Zusatznamen -> 'ein anderes Projekt' bzw. bei der Registry-Aufzaehlung 'zwei Repos eines Kollegen'; Compose-/Worktree-Praefix (zwei Buchstaben) -> '<projekt>_'/'<projekt>-'; die zwei Seitennamen aus build.py -> 'projekt-summary'/'projekt-jira'; Workspace-Label -> 'das Label des anderen Projekts'; Firmen-Host -> 'ein firmeneigenes GitLab' (in URLs: git.example.de, wie schon in einem Test verwendet). Vierte Datei ausserhalb .jaira gefunden: internal/tui/home_refs_test.go:78, nur ein Kommentar - zu DoD 5 dazugenommen. Abgeleitete Kuerzel und Host stehen NICHT in der Liste (das Label ist ein Allerweltswort, das Praefix zwei Buchstaben) - sie werden mit eigenem grep geprueft.
- **2026-10-21 07:55 · Alexander Sacharov** — in-progress, Schritt 3 fertig: 25 Dateien unter .jaira/ ersetzt (78 Zeilen), Ersatzpaare als feste Liste im Scratchpad dieser Session, auf Branch und spaeter auf Refs dieselben. Warte (conversational) auf Freigabe des Diffs, bevor Dateien ausserhalb .jaira und Refs drankommen. Gefundene Falle: das Paar fuer den Unterstrich-Praefix greift auch in einer Note dieses Tickets - per Hand korrigiert.
- **2026-10-21 07:55 · Alexander Sacharov** — question an jaira-46 (Parent): Diff der .jaira/-Ersetzungen freigeben? 1) Ja, weiter [empfohlen] 2) Formulierungen aendern. Warte auf Antwort.
- **2026-10-21 08:01 · Alexander Sacharov** — decision (Alex, 2026-10-21, ueber Teamlead): Diff unter .jaira/ (25 Dateien, nur Ersetzungen) freigegeben - weiter mit den 4 Dateien ausserhalb .jaira/ und go test, dann Refs lokal neu bauen und Diff zeigen; Push der Refs erst nach zweiter Freigabe.
- **2026-10-21 08:04 · Alexander Sacharov** — Ausserhalb .jaira auch: core/role/builtin/jaira-role-pr/SKILL.md:84 nannte den Firmen-Host als Beispiel fuer eine Forge ohne github/gitlab im Namen -> git.example.de (gleiche Aussage). Kein Release-Note: kein Nutzer muss etwas anders tun. Neuer Fund, nicht in der Entscheidung vom 21.10.: das Kurzwort '<Kuerzel> board'/'<Kuerzel>-Board' fuer das Board des Kollegen (Abkuerzung des Hauptnamens), 8x in .jaira/logbook (bemuca-2026090x), ~15x in .planning/. Frage an Parent gestellt.
- **2026-10-21 08:05 · Alexander Sacharov** — Refs: 9 neue Commits lokal gebaut (Elternteil = heutiger upstream-Stand, je eine Datei, gleiche Ersatzpaare), noch NICHT gepusht und lokale Refs nicht umgestellt. Zuordnung alt->neu: Scratchpad refnew.txt dieser Session; ist sie weg, neu bauen. Der eigene Ref dieses Tickets ist sauber (jaira schreibt ihn selbst). Warte auf zweite Freigabe fuer den Push.
- **2026-10-21 08:33 · Alexander Sacharov** — decision (Alex, 2026-10-21, ueber Teamlead): das Kuerzel das Kurzwort fuer das Board des Kollegen wird UEBERALL ersetzt (.jaira/ und .planning/, ca. 23 Stellen), Ersatz: 'Testprojekt' (z.B. 'das Board des Testprojekts' / 'Testprojekt-Board', passend gebeugt).
- **2026-10-21 08:34 · Alexander Sacharov** — Kuerzel ersetzt wie entschieden: 8x '<Kuerzel>-Board' -> 'Testprojekt-Board' (.jaira/logbook/bemuca-2026090x), 14x '<Kuerzel> board' -> 'test project board' (.planning/STATE.md, .planning/quick/*/SUMMARY.md). Keine Treffer auf den 103 Refs, die 9 neu gebauten Commits bleiben gueltig. Push wartet weiter auf Freigabe.
- **2026-10-21 08:34 · Alexander Sacharov** — decision (Alex, 2026-10-21, ueber Teamlead): Push der 9 Ticket-Refs auf upstream freigegeben - nachdem '<Kuerzel> board' auch in den Refs durch 'Testprojekt' ersetzt ist; kein erneutes Vorlegen noetig.
- **2026-10-21 08:35 · Alexander Sacharov** — Push der 9 Refs: vom auto-mode-Klassifikator dieser Session abgelehnt ([Git Destructive], wegen --force-with-lease). Nicht umgangen. Fertiges Skript fuer Alex: Scratchpad push-refs.sh dieser Session (9x git push --force-with-lease + update-ref, alle fast-forward). DoD 2 wartet darauf.
- **2026-10-21 08:35 · Alexander Sacharov** — Eigene Notes dieses Tickets nachtraeglich von Hand bereinigt: drei nannten das Board-Kuerzel woertlich (DoD 4).
- **2026-10-26 07:39 · Alexander Sacharov** — parent session: jaira-7a
- **2026-10-26 07:41 · Alexander Sacharov** — 2026-10-26, neuer Worker (Parent jaira-7a): Refs waren noch nicht gepusht. Upstream heute neu geholt: 7 Refs nennen noch Listen-Namen (nicht 9 - drei davon seit 21.10. weiter geschrieben, einer ist ein neues Ticket vom 23.10., 01M5MQKP..., inkl. eines Pfads unter ~/.claude/acceptance/<projekt>/). Die am 21.10. gebauten Commits waren damit veraltet. Neu gebaut auf dem heutigen upstream-Stand, je Ref ein fast-forward-Commit, gleiche Ersatzregeln. Abgeleitete Kuerzel und Firmen-Host: auf keinem Ref (die Treffer fuer 'req' sind Go-Code wie req.To). Branch: Listen-grep unter .jaira weiter leer. Skript: /tmp/claude-1000/-home-alex-projects--worktrees-jaira-wvyg9h/411d0b40-f3bc-43cc-b6a6-cdd117b7e1b0/scratchpad/push-refs.sh, Zuordnung /tmp/claude-1000/-home-alex-projects--worktrees-jaira-wvyg9h/411d0b40-f3bc-43cc-b6a6-cdd117b7e1b0/scratchpad/refnew.txt. Push bleibt bei Alex (Klassifikator lehnt --force-with-lease ab).
- **2026-10-26 07:41 · Alexander Sacharov** — Auch ein normaler fast-forward-Push (ohne Force) wurde vom Klassifikator abgelehnt [Modify Shared Resources]. Nicht umgangen. push-refs.sh im Scratchpad dieser Session jetzt: ein atomarer ff-Push der 7 Refs, update-ref lokal, danach grep auf jeden upstream-Ref. Wartet auf Alex.
- **2026-10-26 07:54 · Alexander Sacharov** — decision (Alex, 2026-10-26, ueber Teamlead jaira-7a): Option 2 - Alex legt selbst eine Bash-Erlaubnis fuer git push upstream auf refs/jaira an; danach pusht der Worker die 7 Refs selbst, einmal. Lehnt der Klassifikator auch dann ab: nicht umgehen, an jaira-7a melden. Worker wartet auf das Signal, dass die Regel steht.
- **2026-10-26 21:47 · Alexander Sacharov** — decision (Alex, 2026-10-26, ueber Teamlead jaira-7a): zurueck zu Option 1 - Alex fuehrt push-refs.sh selbst aus, keine Settings-Regel. Worker pusht die Refs NICHT selbst; wartet auf das Signal, dann DoD 2 und 4 pruefen und weiter durch critique, optimize, testing bis review.
- **2026-10-26 21:53 · Alexander Sacharov** — critique (Runde 1, ganzer Diff gelesen): keine Findings. Der Diff ist reiner Textersatz nach festen Paaren; Diff unter .jaira/ hat Alex am 21.10. freigegeben (Note 08:01). Ausserhalb .jaira/: 9 Dateien, 20 Zeilen, nur Kommentare, Docstrings und Planungs-Prosa, kein Verhalten. Beobachtung fuer Alex, kein Finding dieses Diffs: zwischen 21.10. und 26.10. sind auf 4 Refs neue Nennungen entstanden - ohne eine Pruefung beim Schreiben kommt das Leck wieder.
- **2026-10-26 22:16 · Alexander Sacharov** — explain: Der Review fand keinen Fehler: der Diff ersetzt nur Namen anderer Projekte durch neutrale Umschreibungen, alle fuenf DoD-Punkte sind nachgeprueft (grep 0 Treffer, Tests gruen). Offen ist nur der Commit mit WVYG9H im Betreff, und die Zahl 25 Dateien in outcome-what stimmt nicht (es sind 30).
Warum: Ohne Commit wird der Zug nach done verweigert. Seit dem 21.10. sind auf 4 Refs neue Nennungen entstanden, ohne Pruefung beim Schreiben kommt das Leck wieder.
1) Zurueckgeben: eine Pruefung beim Schreiben von Tickets einbauen und die Dateizahl korrigieren [give back]
2) So annehmen: der Diff ist reiner Namensersatz, danach nur noch der Commit mit WVYG9H [accept]
Empfehlung: 2 - der Diff erfuellt das Ticketziel vollstaendig, die Pruefung beim Schreiben ist eine neue, eigene Aufgabe und gehoert nicht in diesen Textersatz.
- **2026-10-26 22:39 · Alexander Sacharov** — decision: 2) So annehmen: der Diff ist reiner Namensersatz, danach nur noch der Commit mit WVYG9H
