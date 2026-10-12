---
id: 01M4RPX9C61168AXRHPZGPFA7B
title: "Accept im Cockpit stösst den nächsten Schritt an: PR öffnen oder nächstes Ticket starten"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Druck auf Accept bringt die Arbeit weiter, statt das Ticket nur nach done zu schieben"
context: "Accept im Cockpit-Band schiebt ein Ticket heute nur nach done. Danach passiert nichts.\n- Code: accept() in core/mod/builtin/jaira-cockpit/hooks/register.tsx (~Zeile 323): dry-run move, dann 'jaira move <id> --to <acceptTo> --force', dann Toast '<id> accepted: done'.\n- Wunsch von Alex (2026-10-12): je nachdem, was entschieden wurde, soll Accept den nächsten Schritt starten:\n  a) den Pull Request öffnen (Branch des Tickets ist gepusht, PR fehlt noch), oder\n  b) das nächste Ticket starten (Dispatcher/Worker per spawn.sh, wie der Teamlead es tut).\n- Regel im Repo: ein Agent öffnet nie selbst einen PR, das ist die Entscheidung einer Person. Ein Accept-Druck IST die Entscheidung einer Person; /jaira-role-pr öffnet den PR, wenn eine Person die Rolle aufruft. Das ist der Weg für a).\n- Ungeklärt, für brainstorm: woran erkennt das Cockpit, ob a) oder b) gilt? Kandidaten: gibt es schon einen PR zum Branch (gh/glab), eine 'decision:'-Note auf dem Ticket, oder fragt Accept mit zwei Knöpfen.\n- Ungeklärt: soll das Ticket erst nach dem Merge nach done (dann wäre 'done bei Accept' falsch), oder bleibt done bei Accept."
definition-of-done: Nach Accept startet das Cockpit den entschiedenen nächsten Schritt (PR öffnen oder nächstes Ticket) ohne weiteren Befehl der Person; Test
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T18:44:57Z
updated-at: 2026-10-12T19:07:59Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-48016
claimed-at: 2026-10-12T18:49:39Z
question: ""
mode: conversational
---

# Accept im Cockpit stösst den nächsten Schritt an: PR öffnen oder nächstes Ticket starten

## Definition of Done

- [x] Nach Accept startet das Cockpit den entschiedenen nächsten Schritt (PR öffnen oder nächstes Ticket) ohne weiteren Befehl der Person; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:344 afterAccept() called from accept() :332 after the move; tests 'after Accept' > 'a branch ahead of master without a pull request…' and 'a pull request in any state already there: the next ticket starts…' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:550,563; claude plugin test: 92 pass, exit 0)
- [x] Die Regel, welcher Schritt gilt, steht an einer Stelle im Code und ist im Kommentar begründet
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:242-254 stepAfterAccept() with the comment giving the reason; the only place the step is decided (register.tsx afterAccept calls it)
- [x] Ein PR wird nur über den Accept-Druck einer Person geöffnet, nie von einem Agenten allein; Test
  proof: spawn.sh lane pr typed only by register.tsx afterAccept (after Accept's move); jaira-role-pr SKILL.md --by-person paragraph; tests 'an agent's session has no Accept, so it can never open a pull request', 'a refused Accept starts no next step at all' (board.test.ts:651,642), core/role/role_test.go TestSpawnTypesTheRoleForItsLaneName pr
- [x] Scheitert der Folgeschritt (kein gh, kein Remote, kein nächstes Ticket), sagt das Cockpit es sichtbar und das Ticket bleibt im richtigen Zustand; Test
  proof: register.tsx afterAccept stuck(): message pane '<why>\n<id> stays in done'; tests 'no gh…', 'no remote…', 'a host neither GitHub nor GitLab…', 'no next ticket…', 'no Herdr…' (board.test.ts:595-641)
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/GPFA7B.md
- [ ] Mehrere Tickets auf einem Branch: ein PR wird erst geöffnet, wenn jedes Ticket, das die Commits des Branches nennen (master..branch), in done ist; vorher -> 'next' mit Hinweis, welches Ticket noch fehlt. Ein Accept eines Tickets ohne eigenen Branch, dessen Id in einem anderen Branch steht, findet diesen Branch. Der PR nennt alle Tickets; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] board.ts: reine Regel stepAfterAccept(branch, ahead, pr) + dispatchLane(settings) mit begründendem Kommentar
- [x] spawn.sh: Lane 'pr' tippt /jaira-role-pr <id> --by-person; jaira-role-pr SKILL: --by-person = Accept-Druck einer Person, ein Agent setzt es nie
- [x] register.tsx accept(): nach dem Move Branch/ahead/PR ermitteln (git rev-list, gh/glab), dann PR-Tab oder nächstes Ticket (jaira next --all, erstes ohne Worker-Tab) per spawn.sh; jedes Scheitern sichtbar, Ticket bleibt done
- [x] Tests: board.test.ts (PR, nächstes, kein gh, kein Remote, kein nächstes, kein Herdr) + role_test.go spawn pr
- [x] core/release/notes/GPFA7B.md

## Progress
- **2026-10-12 18:49 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 18:49 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/gpfa7b (Worktree /home/alex/projects/.worktrees/jaira-gpfa7b) ist von feat/qrhb35 (82aa5654) abgezweigt, nicht von master: der Cockpit-Code mit workers()/giveBack ohne --force existiert nur dort. feat/qrhb35 ist noch nicht gepusht/gemergt; der PR dieses Tickets stapelt sich darauf. Parallel laufen WP1CHK (giveBack/Version), GPFA7B (accept), TWB0R6 (Band-Render) - alle in hooks/register.tsx; nur den eigenen Bereich anfassen.
- **2026-10-12 18:49 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/gpfa7b ist von feat/qrhb35 (82aa5654) abgezweigt, nicht von master: Cockpit-Code existiert nur dort; PR stapelt sich darauf. Parallel: WP1CHK (giveBack), TWB0R6 (Band-Render) in hooks/register.tsx - nur accept() und was daran hängt anfassen.
- **2026-10-12 18:50 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan (nur Alex kann sie treffen): (1) woran Accept erkennt, ob PR oder nächstes Ticket gilt; (2) ob done bei Accept bleibt oder erst nach dem Merge; (3) wie der PR geöffnet wird - /jaira-role-pr behandelt eine von spawn.sh getippte Zeile heute als Agenten-Aufruf und öffnet dann keinen PR. Selbst entschieden: 'nächstes Ticket' = jaira next --json, gestartet per spawn.sh mit worker bzw. dispatch je nach settings.json, wie der Teamlead es tut.
- **2026-10-12 18:50 · Alexander Sacharov** — Frage 1/3 an Alex: Woran erkennt Accept, welcher Folgeschritt gilt? 1) Regel im Code (Empfohlen): Ticket-Branch hat Commits vor master und noch keinen PR -> PR öffnen; sonst -> nächstes Ticket starten. Ein Knopf, DoD 2 verlangt genau eine solche Regel. 2) Zwei Knöpfe 'Accept + PR' und 'Accept + nächstes': die Person wählt bei jedem Druck, mehr Platz im Band. 3) Beides nacheinander: PR öffnen, wenn er fehlt, und danach immer das nächste Ticket starten. 4) Eine 'decision:'-Note auf dem Ticket bestimmt es: früher entschieden, aber jemand muss die Note schreiben.
- **2026-10-12 18:56 · Alexander Sacharov** — Alex 2026-10-12, Frage 1 entschieden: Regel im Code - hat der Ticket-Branch Commits vor master und noch keinen PR -> PR öffnen, sonst -> nächstes Ticket. Ein Knopf.
- **2026-10-12 18:56 · Alexander Sacharov** — Frage 2/3 an Alex: Wann kommt das Ticket nach done? 1) Bei Accept, wie heute (Empfohlen): kein neuer Zustand; der Merge ist Sache des Maintainers, und logbook folgt wie bisher. 2) Erst nach dem Merge: Accept öffnet den PR und lässt das Ticket in signoff; jemand muss es nach dem Merge nach done schieben, und das Cockpit muss den Merge erkennen.
- **2026-10-12 18:56 · Alexander Sacharov** — Frage 3/3 an Alex: Wie öffnet Accept den PR? 1) Neue spawn.sh-Lane 'pr' startet /jaira-role-pr <id> --by-person in einem Tab (Empfohlen): die Rolle schreibt die Beschreibung und öffnet den PR; die Rolle lernt, dass diese Marke einen Accept-Druck bedeutet. Nachteil: jeder, der die Zeile tippt, kann die Marke setzen. 2) Das Cockpit selbst ruft 'gh pr create --fill' bzw. 'glab mr create --fill': kein Agent, sofort, aber nur Commit-Titel als Beschreibung, keine Ticket-Zusammenfassung. 3) Das Cockpit pusht nur und zeigt die fertige create-Zeile zum Kopieren: sicherste Variante, aber dann braucht es doch einen zweiten Befehl der Person (widerspricht DoD 1).
- **2026-10-12 18:57 · Alexander Sacharov** — Alex 2026-10-12, Fragen 2+3 entschieden: (2) done bei Accept wie bisher; (3) neue Lane 'pr' in spawn.sh öffnet einen Tab mit /jaira-role-pr <id> --by-person, die Rolle liest --by-person als Accept-Druck einer Person und öffnet den PR selbst.
- **2026-10-12 18:59 · Alexander Sacharov** — Selbst entschieden: master = Branch des Haupt-Worktrees aus 'git worktree list' (der Mock und die Praxis haben dort master), nicht origin/HEAD - ein Remote ist für die Regel nicht nötig. PR-Prüfung: gh pr list --head <branch> --state all (auch ein gemergter/geschlossener PR zählt als vorhanden -> nächstes Ticket); GitLab: glab mr list --source-branch <b> --all -F json; Forge aus remote.origin.url. Scheitert die Prüfung, startet nichts und das Cockpit zeigt warum; das Ticket bleibt done (Alex: done bei Accept). Nächstes Ticket: erstes aus 'jaira next --all --json' ohne Herdr-Tab <id>/worker|dispatch, damit nie ein zweiter Worker auf ein laufendes Ticket kommt; Rolle worker oder dispatch nach ~/.jaira/settings.json (JAIRA_HOME) wie settings.DispatchMode.
- **2026-10-12 19:04 · Alexander Sacharov** — Gelernt: 'claude plugin test .' im Mod-Ordner ist der Testlauf (vorher 81, jetzt 92 grün), 'claude plugin validate .' prüft die $-Aufrufe. Ohne Herdr zeigt jetzt jedes Accept den Pane 'Next step after <id> failed' - so verlangt DoD 4 es; der alte Test 'Accept reads the dry run past the warning…' prüft deshalb nur noch, dass kein 'Accept … failed' kommt. writeMetrics liest settings.json jetzt über dieselbe settingsJson() (JAIRA_HOME beachtet) statt eines zweiten Pfads. Tab-Liste nicht lesbar -> kein nächstes Ticket (sonst könnte ein zweiter Worker auf ein laufendes Ticket kommen).
- **2026-10-12 19:05 · Alexander Sacharov** — Pause (conversational): Diff ungecommittet im Worktree an jaira-f5 geschickt mit der Commit-Zeile; warte auf Alex' OK, bevor ich nach critique ziehe. Critique-neben-der-Arbeit läuft in Herdr-Tab w3:pEC (GPFA7B/critique), Antworten kommen per SendMessage.
- **2026-10-12 19:07 · Alexander Sacharov** — Alex 2026-10-12: Was, wenn zwei Tickets in einem Branch gemacht werden? Mit der Regel (Commits vor master + kein PR -> PR) öffnet das erste Accept den PR, während das zweite Ticket noch läuft, und das zweite Accept findet keinen eigenen Branch (sucht */<slug>) und geht still auf 'next'. Daher DoD 6: Branch über die Ticket-Ids in seinen Commits finden, PR erst wenn alle diese Tickets done sind. Gestapelte Branches (feat/twb0r6 auf feat/qrhb35) sind ein anderer Fall: dort hat jedes Ticket seinen Branch.
- **2026-10-12 19:07 · Alexander Sacharov** — Critique neben der Arbeit, 2 Befunde, beide behoben: (1) PR-Prüfung jetzt wie jaira-role-pr gegen den Ziel-Repo: gh repo view <origin> --json isFork,parent,nameWithOwner bzw. glab repo view -F json (forked_from_project), dann gh pr list --repo <ziel> --head <branch> --state all. '<owner>:<branch>' wie in SKILL.md nicht übernommen: gh pr list --help sagt dazu 'syntax not supported', headRefName passt auch für Fork-PRs. (2) Ohne Worktree: git branch --list --format=%(refname:short) '*/<slug>'.
