---
id: 01M6NHNWE5JR91RRD6Q8D44NVR
title: "panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING"
status: signoff
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Autor bringt sein Plugin per PR in den Katalog, und die CI sagt ihm vor jedem Menschen, was fehlt."
context: |-
  Das Repo https://github.com/panda-bobr/jaira-releases ist seit 05.11.2026 leer.
  Inhalt: plugins.json (leeres Array), README, CONTRIBUTING.md (Fork, Eintrag anhaengen, PR), GitHub-Action, die einen PR prueft.
  CI-Pruefung: id eindeutig und ohne 'jaira'; repo existiert und ist oeffentlich; jaira-plugin.json am Default-Branch passt zum Schema (R6X6FH); Release v<version> hat Assets fuer linux/darwin/windows amd64+arm64 und checksums.txt; PR-Autor ist Besitzer des repo (wie Obsidians GitHub-Verknuepfung).
  Arbeit liegt in einem anderen Repo als jaira: Branch und PR dort, das Ticket faehrt mit einem Notiz-Commit hier.
definition-of-done: "plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases"
tags:
  - ci
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:12Z
updated-at: 2026-11-05T12:23:10Z
updated-by: Alexander Sacharov
question: ""
outcome-what: "panda-bobr/jaira-releases: README auf main, per PR #1 plugins.json ([]), schema/ (Kopien aus R6X6FH), CONTRIBUTING.md, .github/scripts/check_catalog.py mit 23 Unit-Tests und Workflow check-catalog (pull_request_target, nur plugins.json-PRs, Kommentar mit allen Problemen). Fixture-Repos panda-bobr/jaira-plugin-ci-fixture (Release ohne checksums.txt) und -norelease. Test-PR #2 zeigt die Ablehnung."
outcome-why: "Ein Plugin-Autor soll seinen Eintrag per PR einreichen koennen und vor jedem Menschen erfahren, was fehlt; der Katalog aus R6X6FH existierte nur als leeres Repo."
outcome-resolves: "DoD 1 (main 44d293d) und DoD 2 (PR #2, Run 42024115800, Bot-Kommentar mit doppelter id, fehlendem Release, fehlender Pruefsumme)"
review-summary: "Im Repo panda-bobr/jaira-releases liegt jetzt der Plugin-Katalog: plugins.json (leer), schema/ (Kopien der R6X6FH-Schemas), README, CONTRIBUTING und .github/scripts/check_catalog.py mit Workflow check-catalog. Ein PR auf plugins.json startet pull_request_target: der Workflow checkt main aus, liest plugins.json des PR nur als Daten und fuehrt keinen PR-Code aus; das Skript prueft Schema, eindeutige ids ohne 'jaira' und ohne jaira-Befehlsnamen, oeffentliches Repo, Besitz durch den PR-Autor (oder oeffentliche Org-Mitgliedschaft), jaira-plugin.json auf dem Default-Branch, Release v<version> mit sechs Binaries und checksums.txt, und schreibt alle Probleme als einen Bot-Kommentar, der bei jedem Push aktualisiert wird. Aendert ein PR repo oder author einer gelisteten id, muss der Autor auch das bisherige repo besitzen (check_takeover). GitHub-Fehler (Rate-Limit, 5xx) erscheinen als Problem statt als leerer Kommentar. Ein eigener Job unit laeuft auf pull_request mit paths .github/** und schema/** nur mit Lese-Token und fuehrt die 29 Unit-Tests aus. Im jaira-Repo selbst aendert der Branch nur die Ticketdatei."
round: 3
review-gaps: |-
  note: 56da600 (optimize-Refactor: changed_entries inline) liegt auf origin/fix/d44nvr-errors, ist aber nicht in main - PR #3 wurde bei 52c090e gemergt; Verhalten gleich, also entweder nachziehen oder verwerfen
  note: author ist Freitext (Schema: string 1-80) - ein neuer Eintrag kann jeden Autorennamen tragen; geprueft wird nur der Besitz von repo
  note: wird ein PR zwischen Event und Lauf erneut gepusht, findet git show den HEAD_SHA nicht und die Meldung lautet 'plugins.json must be a JSON array' statt 'neu ausfuehren'; der folgende synchronize-Lauf ersetzt den Kommentar
  note: DoD 3 ist nur per Unit-Test belegt (4 Tests mit aufgezeichneten API-Antworten), nicht per Live-PR - die DoD verlangt genau einen Test
test-verdict: "pass: unit RC=0 (29 Tests, 56da600); DoD 3 live: Verschieben von hello (panda-bobr/jaira-hello-demo -> sashasoft90/jaira) durch octocat abgelehnt mit 'needs the owner of the listed repository', durch sashasoft90 (jetzt oeffentliches panda-bobr-Mitglied, anon 204) korrekt erlaubt; DoD 4: PR #3 Run 42026044679 job unit success, check skipped"
claimed-by: DESKTOP-RFTCH11-12777
claimed-at: 2026-11-05T12:10:19Z
review-verdict: "Alle vier DoD-Punkte erfuellt und auf GitHub nachgeprueft: main 85953d3 enthaelt plugins.json, README, CONTRIBUTING, Skript und Workflow (PR #1 und #3 gemergt); PR #2 Lauf 42024115800 lehnt mit Bot-Kommentar ab (doppelte id, fehlendes Release, fehlende checksums.txt); check_takeover mit Tests; Lauf 42026044679 (pull_request) Job unit success, 29 Tests lokal auf origin/main gruen. Workflow-Sicherheit geprueft: kein PR-Code unter pull_request_target, PR-Werte nur ueber env. Kein Blocker; vier Notes."
review-round: 3
review-check: "1. Oeffne https://github.com/panda-bobr/jaira-releases - auf main liegen plugins.json, README.md, CONTRIBUTING.md, schema/ und .github/.  2. Oeffne https://github.com/panda-bobr/jaira-releases/pull/2 - unten steht ein Bot-Kommentar 'Catalogue check failed' mit den Zeilen 'the id is taken', 'no release v0.0.1' und 'has no checksums.txt'.  3. Oeffne https://github.com/panda-bobr/jaira-releases/actions/runs/42026044679 - Job unit gruen, Job check uebersprungen (DoD 4).  4. DoD 3 lokal: gh repo clone panda-bobr/jaira-releases /tmp/jr && cd /tmp/jr && uv run --no-project --with jsonschema==4.23.0 python -m unittest discover -s .github/scripts -v 2>&1 | grep takeover - vier Tests mit takeover/owner_of_both/description_change, alle 'ok'; die letzte Zeile der vollen Ausgabe ist 'OK' bei 29 Tests."
---

# panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING

## Definition of Done

- [x] plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases
  proof: panda-bobr/jaira-releases main 44d293d (Merge PR #1): plugins.json, README.md, CONTRIBUTING.md; check-catalog push run success on 44d293d
- [x] Die Action lehnt einen PR mit doppelter id, fehlendem Release oder fehlender Pruefsumme ab; in einem Test-PR gezeigt
  proof: panda-bobr/jaira-releases PR #2 (draft test PR): check-catalog run 42024115800 failed; bot comment issuecomment-6657527419 lists duplicate id ci-fixture, no release v0.0.1 (ci-fixture-norelease), no checksums.txt (ci-fixture); unit tests .github/scripts/test_check_catalog.py (23)
- [x] Ein PR, der repo oder author eines schon gelisteten id aendert, wird abgelehnt, wenn der Autor nicht beide repos besitzt; Test
  proof: jaira-releases fix/d44nvr-errors 52c090e: .github/scripts/check_catalog.py check_takeover; tests test_takeover_of_a_listed_id_is_refused, test_takeover_with_foreign_old_repo_is_refused, test_owner_of_both_repos_may_move_the_entry, test_description_change_needs_no_old_owner_check (29 OK)
- [x] Die Unit-Tests der Katalogpruefung laufen in jedem PR, der .github/** aendert; im PR fix/d44nvr-errors gezeigt
  proof: panda-bobr/jaira-releases PR #3 (fix/d44nvr-errors 52c090e): run 42026044679, event pull_request, job unit success 'Ran 29 tests ... OK', job check skipped; .github/workflows/check-catalog.yml on.pull_request paths .github/**, schema/**

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Repo jaira-releases lokal anlegen, README-Start-Commit auf main pushen
- [x] plugins.json, schema/, CONTRIBUTING.md, check_catalog.py + Unit-Tests, Workflow auf feat/d44nvr
- [x] Fixture-Repo sashasoft90/jaira-plugin-ci-fixture mit Release v0.0.1 ohne checksums.txt
- [x] Test-Branches (doppelte id, fehlendes Release, fehlende Pruefsumme) pushen, gh-Zeilen fuer Alex
- [x] DoD 3: Besitzerpruefung des alten repo bei geaendertem repo/author eines gelisteten id, mit Tests
- [x] DoD 4: Unit-Job auf pull_request (paths .github/**, schema/**), concurrency je Event
- [x] Push-Zeile und gh pr create fuer fix/d44nvr-errors an Teamlead; nach PR-Lauf DoD 4 mit Proof

## Progress
- **2026-11-05 10:23 · Alexander Sacharov** — Teamlead (05.11.2026): der Standard (R6X6FH, abgenommen) liegt noch nicht auf master, sondern auf origin/feat/r6x6fh: docs/PLUGINS.md und docs/schema/*.json - lesen mit 'git show origin/feat/r6x6fh:docs/PLUGINS.md'. Die Arbeit selbst gehoert in das Repo panda-bobr/jaira-releases; Alex ist dort admin. Ein Agent pusht seinen Branch und oeffnet keinen PR.
- **2026-11-05 10:24 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 10:24 · Alexander Sacharov** — decision (worker, settled): Pruefung als Python-Skript (stdlib + jsonschema per pip) in .github/scripts/check_catalog.py; Workflow auf pull_request_target, checkt NUR main aus und liest plugins.json des PR-Kopfs als Daten (git show), fuehrt nie PR-Code aus - sonst koennte ein PR das Pruefskript selbst aendern und sich durchwinken. PR, der andere Dateien als plugins.json aendert, wird abgelehnt.
- **2026-11-05 10:24 · Alexander Sacharov** — decision (worker, settled): Schemas aus R6X6FH (docs/schema/*.json, Branch feat/r6x6fh) werden nach jaira-releases/schema/ kopiert, nicht zur CI-Zeit aus jaira geholt - Repo bleibt eigenstaendig; bei Schema-Aenderung in jaira muss die Kopie nachgezogen werden (steht im README). Besitzer-Pruefung: PR-Autor == Besitzer von repo, oder Besitzer ist eine Organisation und der Autor ist dort oeffentliches Mitglied (/orgs/<org>/public_members/<user>, ohne Token abfragbar). jaira-Befehlsnamen als feste Liste im Skript.
- **2026-11-05 10:25 · Alexander Sacharov** — question (open, an Alex): 1) Wie kommt der Inhalt auf main des leeren Repos? (a, empfohlen) Worker pusht nur README als Start-Commit auf main, alles andere auf Branch feat/d44nvr, Alex oeffnet+merged den PR. (b) Worker pusht alles direkt auf main, ohne Review. 2) Test-PR fuer DoD 2: (a, empfohlen) Worker legt oeffentliches Fixture-Repo sashasoft90/jaira-plugin-ci-fixture an (Release v0.0.1 mit jaira-plugin.json, ohne checksums.txt), pusht Test-Branches, Alex oeffnet die Test-PRs mit gelieferten gh-Zeilen. (b) Wie a, aber Worker oeffnet die Test-PRs selbst (als Test markiert, danach geschlossen) - Ausnahme von 'Agents oeffnen keine PRs'. (c) Kein Fixture-Repo: Test-PR zeigt nur doppelte id und fehlendes Release; fehlende Pruefsumme nur im Unit-Test mit aufgezeichneten API-Antworten.
- **2026-11-05 10:27 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): 1a - nur README als Start-Commit auf main von panda-bobr/jaira-releases, alles andere auf feat/d44nvr, Alex oeffnet und merged den PR. 2a - oeffentliches Fixture-Repo sashasoft90/jaira-plugin-ci-fixture (Release v0.0.1 ohne checksums.txt), Test-Branches pushen, Alex oeffnet die Test-PRs mit fertigen gh-Zeilen.
- **2026-11-05 10:29 · Alexander Sacharov** — dead end: git push origin main nach panda-bobr/jaira-releases wurde vom Auto-Mode-Klassifikator dieser Sitzung verweigert. Nicht umgangen. Alles liegt lokal im Scratchpad-Clone (main: README-Commit eb57b38, Branch feat/d44nvr); ein neu gestarteter Worker findet den Clone nicht - Inhalt muss dann neu erzeugt oder Alex pusht. Push ist Alex' Entscheidung.
- **2026-11-05 10:33 · Alexander Sacharov** — dead end: auch 'gh repo create sashasoft90/jaira-plugin-ci-fixture --public' vom Klassifikator verweigert ([Create Public Surface]). Nicht umgangen, nicht an den Teamlead weitergereicht (waere Permission-Laundering). Alles liegt lokal fertig in /tmp/claude-1000/-home-alex-projects--worktrees-jaira-d44nvr/e9c86c25-536c-4790-a6d8-3e163fb430b1/scratchpad: jaira-releases (main eb57b38 README; feat/d44nvr e4f71cd; test/d44nvr-catalog-check fe5c04b), fixtures/ (zwei Repos), publish-1-before-merge.sh und publish-2-after-merge.sh fuer Alex.
- **2026-11-05 10:33 · Alexander Sacharov** — decision (worker): ZWEI Fixture-Repos statt einem - jaira-plugin-ci-fixture (Release v0.0.1 ohne checksums.txt) und jaira-plugin-ci-fixture-norelease (Manifest, kein Release). Grund: die Pruefung liest das Manifest vom Default-Branch, ein Repo kann nicht zugleich 'Release fehlt' und 'Pruefsumme fehlt' zeigen. So zeigt EIN Test-PR alle drei Faelle (DoD 2: doppelte id, fehlendes Release, fehlende Pruefsumme). Pruefskript gegen echte GitHub-API geprueft: nicht existierendes Repo, Org-Repo ohne oeffentliche Mitgliedschaft, fehlendes Manifest werden korrekt gemeldet; 22 Unit-Tests gruen (python -m unittest discover -s .github/scripts).
- **2026-11-05 10:33 · Alexander Sacharov** — gotcha: ein Branch, der als erster in ein leeres Repo gepusht wird, wird Default-Branch. Darum feat/d44nvr erst nach main pushen (publish-1 tut das in dieser Reihenfolge). check-catalog laeuft auf pull_request_target nur mit paths: plugins.json - der PR feat/d44nvr selbst loest ihn nicht aus (Workflow liegt noch nicht auf main); Unit-Job laeuft bei push auf main.
- **2026-11-05 10:37 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): Veroeffentlichung freigegeben; Fixture-Repos nach panda-bobr statt sashasoft90 (panda-bobr/jaira-plugin-ci-fixture, panda-bobr/jaira-plugin-ci-fixture-norelease). publish-1-before-merge.sh und test/d44nvr-catalog-check (neuer Commit 2d830c6, repo-URLs) angepasst. Teamlead fuehrt publish-1 auf Alex' Anweisung aus; der Worker pusht nichts. Erwartung Test-PR: zusaetzlich zu doppelter id, fehlendem Release, fehlender Pruefsumme meldet die Pruefung 'sashasoft90 is not a public member of panda-bobr', solange Alex' Mitgliedschaft nicht oeffentlich ist. Trockenlauf gegen panda-bobr/jaira-hello-demo (id hello): Manifest, Release v0.1.0, Binaries und checksums.txt bestehen; einzig die Mitgliedschaftspruefung schlaegt fehl - Kandidat fuer einen gruenen PR, sobald die Mitgliedschaft oeffentlich ist.
- **2026-11-05 10:39 · Alexander Sacharov** — Trockenlauf mit den echten Fixture-Repos (nach publish-1): Pruefung meldet doppelte id, fehlendes Release (norelease), fehlendes checksums.txt (ci-fixture) und die Mitgliedschaft. Gefunden: eine doppelte id wurde zweimal gegen GitHub geprueft und doppelt gemeldet -> Commit 6840b4f auf feat/d44nvr prueft je id nur einmal (+Test, 23 gruen). 6840b4f ist LOKAL und muss vor dem Merge noch nach feat/d44nvr gepusht werden; test-Branch darauf rebased.
- **2026-11-05 10:39 · Alexander Sacharov** — PR https://github.com/panda-bobr/jaira-releases/pull/1 offen (feat/d44nvr bei 6840b4f), auf Alex' Anweisung vom Teamlead geoeffnet; Merge macht Alex, danach publish-2 und Test-PR.
- **2026-11-05 11:44 · Alexander Sacharov** — publish-2 umgebaut: kein Force-Push mehr. Worker hat test/d44nvr-catalog-check lokal auf origin/main 44d293d rebased (Diff = nur plugins.json); der Branch war nie gepusht, also reicht 'git push origin test/d44nvr-catalog-check'. Die Ablehnung im Teamlead lag am -f.
- **2026-11-05 12:06 · Alexander Sacharov** — critique (Runde 1): 2 Funde - (1) Nicht-404-HTTP-Fehler ergeben leeren Bot-Kommentar: fangen und als Problem melden, Kommentar-Schritt an report.md binden; (2) event/manifest aus JAIRA_COMMANDS streichen. Bewusst stehen gelassen: Schema-Kopien statt Abruf aus jaira (Notiz 'decision (worker, settled)' zu schema/), Python statt Go (separates Repo, kein Code aus jaira wiederverwendbar). Fix geht als neuer Branch fix/d44nvr-errors in jaira-releases und braucht wieder einen PR von Alex.
- **2026-11-05 12:06 · Alexander Sacharov** — Runde-1-Fix: Branch fix/d44nvr-errors in jaira-releases, Commit 5bace82 (lokal, ungepusht, auf origin/main 44d293d). Dead end: 'if: hashFiles(runner.temp/report.md)' geht nicht - hashFiles sieht nur Dateien unter GITHUB_WORKSPACE; daher Shell-Guard '[ -f report.md ] || exit 0' im Kommentar-Schritt. 25 Unit-Tests gruen.
- **2026-11-05 12:07 · Alexander Sacharov** — critique (Runde 2): beide Funde aus Runde 1 in 5bace82 erledigt (HTTPError/URLError in run() gefangen und als Problem gemeldet, Kommentar-Schritt hinter report.md; event/manifest gestrichen). Keine neuen Funde.
- **2026-11-05 12:08 · Alexander Sacharov** — offen nach testing: fix/d44nvr-errors (5bace82, 310de33) liegt lokal im Scratchpad-Clone, ungepusht; braucht Push + PR + Merge durch Alex. DoD ist schon durch main 44d293d erfuellt, der Branch verbessert nur Fehlermeldungen. Push-Zeile: git -C /tmp/claude-1000/-home-alex-projects--worktrees-jaira-d44nvr/e9c86c25-536c-4790-a6d8-3e163fb430b1/scratchpad/jaira-releases push origin fix/d44nvr-errors
- **2026-11-05 12:10 · Alexander Sacharov** — review (frische Sitzung): kein Blocker, DoD 1+2 auf GitHub nachgeprueft (main 44d293d, PR #2 Run 42024115800, Schemas byte-gleich mit origin/feat/r6x6fh, 23 Tests gruen). Vor done: fix/d44nvr-errors ist noch nicht gepusht - auf main postet die Pruefung bei Rate-Limit/5xx einen leeren Kommentar.
- **2026-11-05 12:11 · Alexander Sacharov** — question geleert (Push-Freigabe ist erteilt und erledigt); Hinweis aus review.
- **2026-11-05 12:13 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): beide Review-Hinweise werden in diesem Ticket behoben, auf fix/d44nvr-errors (noch nicht gepusht). (1) Ein PR, der bei einem schon gelisteten id das Feld repo (oder author) aendert, wird abgelehnt, ausser der PR-Autor ist Besitzer des bisherigen UND des neuen repo. (2) Die Unit-Tests von .github/scripts laufen auch in jedem PR, der .github/** aendert - nicht nur bei Push auf main. Achtung: pull_request_target checkt main aus, fuer die Tests braucht es einen eigenen Job auf pull_request, der keine Geheimnisse hat.
- **2026-11-05 12:15 · Alexander Sacharov** — DoD 3/4 umgesetzt in 52c090e auf fix/d44nvr-errors (lokal). Entscheidungen: verschwundenes altes repo -> Ablehnung mit 'ask a maintainer' (Besitz nicht beweisbar); nur repo/author-Aenderung loest die Altbesitzer-Pruefung aus, description/name nicht. Unit-Job: pull_request mit paths .github/** und schema/**, permissions nur contents: read; pull-requests: write jetzt nur am check-Job. Concurrency-Gruppe enthaelt event_name, sonst bricht pull_request den pull_request_target-Lauf desselben PRs ab. DoD 4 wird im PR fix/d44nvr-errors gezeigt: der PR aendert .github/** -> unit laeuft; check laeuft nicht (kein plugins.json).
- **2026-11-05 12:19 · Alexander Sacharov** — critique (Runde 3, nur 52c090e gelesen): keine Funde. check_takeover steht neben check_entry und nutzt check_owner wieder; pull_request-Job laeuft PR-Code nur mit Lese-Token - Standardmuster. Hinweis: outcome-resolves nennt nur DoD 1+2; DoD 3+4 tragen ihre Proofs selbst.
- **2026-11-05 12:20 · Alexander Sacharov** — testing: sashasoft90 ist seit ~12:20 oeffentliches Mitglied von panda-bobr (public_members anon 204) - Mitgliedschafts-Zeile faellt fuer panda-bobr-Repos jetzt weg. 56da600 (optimize) ist lokal und noch nicht in PR #3; Push: git -C /tmp/claude-1000/-home-alex-projects--worktrees-jaira-d44nvr/e9c86c25-536c-4790-a6d8-3e163fb430b1/scratchpad/jaira-releases push origin fix/d44nvr-errors
- **2026-11-05 12:23 · Alexander Sacharov** — review (frische Sitzung, Runde 2): kein Blocker. Nachgeprueft auf GitHub: PR #1/#3 gemergt (main 85953d3), PR #2 geschlossen mit Bot-Kommentar, Lauf 42026044679 unit success; 29 Tests auf origin/main lokal gruen. Offen nur: 56da600 liegt auf origin/fix/d44nvr-errors, nicht in main (reiner Refactor).
