---
id: 01M2MFEEE9VZAABPKKS27MG5GB
title: "spawn.sh schreibt Werte, die nur einem Projekt gehoeren"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "spawn.sh legt den worktree an, oeffnet die Kachel und startet den Arbeiter - mehr nicht. Alles, was ein Projekt zum Loslaufen braucht, macht ein Skript im Repository, das spawn.sh aufruft und mit worktree-Pfad, Slug, Port-Versatz und Wurzel versorgt. Ohne dieses Skript legt spawn.sh nur den worktree an."
context: |-
  spawn.sh legt fuer jeden Arbeiter einen worktree an und schreibt ihm eine eigene .env. Darin stehen heute HTTP_PORT und DB_PORT_HOST - Namen, die sich ein einzelnes Projekt ausgedacht hat. Ein anderes Projekt nennt seine Ports anders oder hat keine, und bekommt die Zeilen trotzdem.

  Die Folge sieht man an einem Board, das jaira benutzt: dessen spawn.sh ist von 100 auf 159 Zeilen gewachsen, weil dort IMAGE_NS, COMPOSE_PROFILES, VITE_PORT_HOST und BACKEND_PORT_HOST nachgetragen wurden. Alles vier gehoert diesem einen Projekt.

  Das Schlimme daran ist nicht das Ueberschreiben - 'roles install' laesst eine bearbeitete Datei stehen und meldet sie. Es ist das Gegenteil: die Datei ist eingefroren. Jede Verbesserung an spawn.sh wird bei diesem Anwender still uebersprungen, und er merkt es erst, wenn er die Ausgabe von 'roles install' liest. Der Fork ist auch nicht sortenrein: --no-worktree und der Lane-Name 'dispatch' stehen darin und sind allgemein, sie gehoeren hierher.

  Am 16.09.2026 aufgekommen, als in einem worktree die COMPOSE_PROFILES des Hauptverzeichnisses mitkopiert wurden und der Zweig acht Container zu viel hochfuhr.

  Noch offen und Teil der Loesung: welche Angaben der Hook braucht. Vorschlag: Slug und Port-Versatz als Argumente, damit die Datei ihre Ports genauso verteilen kann wie spawn.sh heute.
definition-of-done: "spawn.sh ruft .jaira/worktree-env auf, wenn die Datei existiert und ausfuehrbar ist, und haengt ihre Ausgabe an die .env des worktree"
tags:
  - cli
blocked-by:
  - 01M7WWYKK8CMJSM2S45WJDDCEZ
related:
  - 01M7WWYKK8CMJSM2S45WJDDCEZ
commits: []
created-at: 2026-09-16T06:46:09Z
updated-at: 2026-11-22T11:15:04Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-33230
claimed-at: 2026-11-22T11:08:10Z
outcome-what: "spawn.sh no longer copies .env or writes port names; without .jaira/worktree-setup it makes the worktree only, hands out no offset, and tells a repo with a .env how to get the old behaviour from the shipped worktree-setup.example; roles install --help says an edited file is frozen"
outcome-why: "Port names like HTTP_PORT belonged to one project; repos forked spawn.sh to add theirs, and an edited spawn.sh never gets jaira's fixes again"
outcome-resolves: "spawn.sh schreibt Werte, die nur einem Projekt gehoeren - DoD 1-7"
---

# spawn.sh schreibt Werte, die nur einem Projekt gehoeren

## Definition of Done

- [x] spawn.sh ruft ein Einrichtungsskript des Repositorys auf und uebergibt ihm worktree-Pfad, Slug, Port-Versatz und Repository-Wurzel
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:284 (since 167HY8); TestSpawnLetsTheRepositorySetUpItsWorktree/setup checks the four args
- [x] Ohne dieses Skript legt spawn.sh den worktree an und richtet sonst nichts ein
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:206,304; TestSpawnLetsTheRepositorySetUpItsWorktree (no-setup loop: no .env, no jaira-offset, no lock)
- [x] Weder .env-Kopie noch HTTP_PORT oder DB_PORT_HOST stehen noch in spawn.sh; die Dokumentation zeigt ein Beispielskript, das beides macht
  proof: grep HTTP_PORT|DB_PORT_HOST|cp .env in core/role/builtin/jaira-dispatcher/scripts/spawn.sh = 0; example core/role/builtin/jaira-dispatcher/scripts/worktree-setup.example, named in spawn.sh --help
- [x] Die Hilfe von 'roles install' sagt, dass eine bearbeitete Datei nicht nur unveraendert bleibt, sondern auch keine Verbesserungen mehr bekommt
  proof: internal/cli/roles.go:410-414
- [x] Ein Board, das bisher ohne Einrichtungsskript lief, bekommt beim naechsten Arbeiter eine verstaendliche Meldung statt einer stillen Verhaltensaenderung
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:304-307 stderr line with install -m 755 line; TestSpawnLetsTheRepositorySetUpItsWorktree checks it is said with .env and not without
- [x] worktree-setup.example reproduces the old .env byte for byte (copy, '# worker stack' line, COMPOSE_PROJECT_NAME, HTTP_PORT, DB_PORT_HOST, DB_PORT_TEST_HOST); the offset tests use it as their setup, in a repository whose path holds a space
  proof: TestSpawnGivesEveryWorktreeItsOwnOffset byte compare of repo-8r8vts/.env; newOffsetRepo installs the example, base path 'with space'
- [x] no setup: no port offset, no offset lock, no .env in the worktree; with a .env in the root one stderr line says so and names the cp line - a test for both repos (with and without .env)
  proof: TestSpawnLetsTheRepositorySetUpItsWorktree no-setup loop over dotenv true/false

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] move the .env/port block out of spawn.sh into scripts/worktree-setup.example (executable when copied)
- [x] offset only with a setup; the DoD-5 stderr line where root has .env
- [x] spawn.sh --help, roles install --help, release note
- [x] tests: example as setup in role_test.go offset tests + no-setup cases
- [x] regenerate core/shipped/shipped.sha256, go test ./core/role ./internal/cli

## Progress
- **2026-09-16 06:46 · Alexander Sacharov** — Praezisiert am 16.09.2026 auf Ansage: jaira soll den worktree anlegen und sonst nichts. Die erste Fassung dieses Tickets hat noch von einer Datei gesprochen, die Zeilen fuer die .env ausgibt - das ist zu eng gedacht.

Warum zu eng: das Kopieren der .env ist selbst schon projektspezifisch. Ein Projekt ohne .env braucht es nicht, ein anderes will daneben noch 'npm install' oder eine Datenbank hochfahren. Wenn jaira die .env kopiert und ein Hook nur Zeilen anhaengen darf, bleibt die eine Entscheidung, die am haeufigsten falsch ist, weiter in jaira.

Also: ein Einrichtungsskript, kein Zeilenlieferant. jaira ruft es auf, gibt ihm worktree-Pfad, Slug, Port-Versatz und Repository-Wurzel, und was danach in dem Verzeichnis steht, ist Sache des Projekts. Rueckwaerts vertraeglich bleibt es dadurch, dass ein Board ohne dieses Skript einen worktree ohne .env bekommt - wer die alte Bequemlichkeit will, schreibt sich das Skript einmal hin, und jaira zeigt eins als Vorlage.
- **2026-09-16 06:48 · Alexander Sacharov** — Reihenfolge, entschieden von der Teamlead-Sitzung am 2026-09-16, damit niemand hier anfaengt, bevor 7KX89C durch ist.

Drei Tickets fassen dieselbe Datei an: 7KX89C (spawn.sh faehrt im Repository mit), dieses hier (Projektwerte raus aus spawn.sh), 3YRPXJ (spawn.sh kann nur eine Rolle starten). Das ist eine Arbeit an einem File, keine drei nebeneinander.

7KX89C geht zuerst, und der Grund steht im Kontext DIESES Tickets: spawn.sh liegt heute nur in ~/.claude/skills/jaira-dispatcher/ auf einer Maschine. Wer hier zuerst baut, aendert eine Datei, die in keinem Diff steht, in keinem PR auftaucht und bei Berk nicht ankommt - und muss dieselbe Aenderung ein zweites Mal machen, sobald 7KX89C die Datei ins Repository holt. Genau das Einfrieren, das dieses Ticket beklagt, wuerde die Behebung dieses Tickets selbst treffen.

Danach dieses Ticket, dann 3YRPXJ.

Was aus der Sitzung noch dazugehoert: der Teamlead hat dreimal eine gepatchte Kopie von spawn.sh im Scratchpad gefahren, weil die Prompt-Zeile fest verdrahtet ist. Das ist 3YRPXJ und dort notiert - aber es ist dasselbe Muster wie hier: wer spawn.sh nicht erweitern kann, forkt es. Eine Loesung fuer dieses Ticket, die nur die .env herausnimmt, laesst den Fork-Grund von 3YRPXJ stehen. Die Plan-Lane sollte beide Haken - Einrichtungs-Hook und Prompt-Argument - als eine Erweiterbarkeitsfrage ansehen, auch wenn sie in zwei Tickets gebaut werden.
- **2026-09-16 07:50 · Alexander Sacharov** — Die Schnittstelle ist unten schon in Betrieb, im Board eines anderen Projekts - sie muss nicht erfunden, nur uebernommen werden. Aufruf, wie er dort laeuft: setup="$root/.jaira/worktree-setup"; if [ -x "$setup" ]; then off=$(( ( $(printf '%s' "$slug" | cksum | cut -d' ' -f1) % 40 ) + 1 )); "$setup" "$wt" "$slug" "$off" "$root"; fi

Vier Dinge, die der Umzug nicht von selbst mitbringt: (1) Der Port-Versatz bleibt bei jaira - die einzige Zusicherung, die ein Projekt nicht geben kann: zwei Arbeiter duerfen nie einen Stapel teilen. jaira rechnet ihn aus dem Slug und uebergibt ihn; das Skript des Projekts verteilt ihn nur auf Namen, die es selbst kennt. (2) Bei --no-worktree darf das Skript nicht laufen: dort ist wt == root, und es wuerde in die lebende .env des Repositorys schreiben. Heute schuetzt das nur zufaellig das 'if [ ! -d "$wt" ]' darum herum. (3) Ein fehlschlagendes Skript darf den Arbeiter nicht toeten: spawn.sh hat 'set -euo pipefail', ein Exit ungleich 0 aus dem Hook bricht ab, bevor die Kachel ueberhaupt aufgeht. Abfangen und melden, nicht abbrechen. (4) 'fehlt' und 'nicht ausfuehrbar' getrennt melden, sonst sieht ein vergessenes chmod +x aus wie ein Projekt ohne Einrichtung.

Was der Umzug loest, gemessen: das Exemplar von spawn.sh in einem anderen Projekt war auf 159 Zeilen gewachsen - IMAGE_NS, COMPOSE_PROFILES, AZIMUTT_PORT, VITE_PORT_HOST, BACKEND_PORT_HOST, dazu ein pg_dump der Hauptdatenbank. Alles davon traegt Namen, die nur dieses eine Projekt kennt. 'roles install' ueberschreibt so eine Datei nicht, sie friert nur ein: jede Verbesserung an spawn.sh wird bei dem Anwender still uebersprungen.

Nicht projektspezifisch und daher in spawn.sh gehoerend, nicht in den Hook: --no-worktree und der Lane-Name 'dispatch'.
- **2026-11-20 16:50 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-22 11:06 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2784.sock
- **2026-11-22 11:08 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 11:09 · Alexander Sacharov** — decisions (worker, no open question for the person): (1) DoD 1 already holds since 167HY8 (spawn.sh:182-237 calls .jaira/worktree-setup with the four args); this ticket removes the .env fallback only. (2) Without a setup no port offset is handed out at all - no setup means no stack. (3) DoD 5: a new worktree without setup in a repository that has a .env gets one stderr line naming the change and the cp line for the example; a repo without .env gets silence (nothing changed for it). (4) The example ships beside spawn.sh as scripts/worktree-setup.example (installed by roles install, so the message can name a path that exists on the user's disk) and reproduces the old behaviour byte for byte; the offset tests run with it as their setup, which proves it. (5) The offset scan still reads '# worker stack: offset=' from the .env of worktrees made before offsets were written down - compatibility for standing worktrees, not a project value.
- **2026-11-22 11:15 · Alexander Sacharov** — in-progress: shipped.sha256 regenerated after 'git add' of the new example (shipped-history.sh reads tracked files only - untracked it is silently left out). .gitattributes needed '*.example text eol=lf' or internal/wintrap fails. Kept the 'skipped ... (edited here; --force to replace)' output line unchanged: the cockpit and updateauto_test parse it; DoD 4 asks for the help only.
