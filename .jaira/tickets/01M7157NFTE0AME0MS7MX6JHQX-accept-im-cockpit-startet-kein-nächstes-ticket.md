---
id: 01M7157NFTE0AME0MS7MX6JHQX
title: Accept im Cockpit startet kein nächstes Ticket mehr
status: testing
ready: true
creator: Alexander Sacharov
goal: "Accept im Cockpit öffnet weiter den PR, startet aber nie von selbst das nächste Ticket; was als Nächstes läuft, entscheidet der Teamlead"
context: |-
  Falsch heute: Accept im Cockpit startet nach done selbst das nächste Ticket in einem neuen Herdr-Tab.
  Ausgelöst am 2026-11-09 in einem anderen Projekt: Toast '6XTM56 accepted; 992RJX started in a Herdr tab'. Alex wusste nicht, was da startete und warum.
  Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx, afterAccept(), Zweig step === 'next' (jaira next --all, nextToStart, spawn.sh --no-parent).
  Eingeführt mit GPFA7B (Commit 0a2b3eb9, 2026-10-12).
  Entscheidung von Alex: der Teil 'nächstes Ticket starten' fliegt raus. Der Teil 'PR öffnen' (step === 'pr') bleibt genau wie er ist.
  Wer das nächste Ticket startet: der Teamlead (Skill jaira-teamlead) — er fragt die Person oder startet selbst.
  Wenn kein PR fällig ist, sagt Accept nur das: z.B. '<id> accepted' plus den bestehenden Hinweis (prWaits / prUnchecked), ohne Tab zu öffnen.
  Texte in hooks/texts.ts (nextStarted, nothingNext, nextNeedsHerdr) und Tests in tests/board.test.ts (Toasts mit 'started in a Herdr tab') ziehen mit.
  Client-sichtbar: Release-Note core/release/notes/<ID>.md nötig.
definition-of-done: Accept auf ein Ticket ohne fälligen PR öffnet keinen Herdr-Tab und ruft spawn.sh nicht auf
tags:
  - acceptance
blocked-by: []
related: []
commits: []
created-at: 2026-11-09T22:00:36Z
updated-at: 2026-11-09T22:08:13Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-64644
claimed-at: 2026-11-09T22:01:08Z
outcome-what: "Accept im Cockpit öffnet nur noch den PR, wenn einer fällig ist; der Zweig, der das nächste Ticket per jaira next + spawn.sh in einem Herdr-Tab startete, ist entfernt, samt nextToStart, dispatchLane und drei Texten"
outcome-why: "Das automatisch gestartete Ticket kam ohne Ankündigung und ohne Grund; welches Ticket als Nächstes läuft, entscheidet der Teamlead"
outcome-resolves: Accept ohne fälligen PR startet nichts mehr; Toast nennt nur Ticket und PR-Hinweis
review-summary: "none"
round: 2
review-gaps: "register.tsx afterAccept: hint aus stuck() entfernt — nach dem Wegfall des Next-Zweigs war er dort immer leer (PR-Schritt hat nie einen Hinweis, Fehlerstring auch nicht); tests/board.test.ts: settings.json-Fixture {dispatch:worker} im PR-schon-da-Test entfernt, liest keiner mehr. Gelassen: w.next/nextFails-Fixture (Overview liest jaira next), tabsOf/settingsJson (andere Aufrufer), stepAfterAccept als 'pr'|'none' statt boolean (Benennung trägt), herdr()-Prüfung im PR-Zweig (bestehendes Muster)."
---

# Accept im Cockpit startet kein nächstes Ticket mehr

## Definition of Done

- [x] Accept auf ein Ticket ohne fälligen PR öffnet keinen Herdr-Tab und ruft spawn.sh nicht auf
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:596; tests 'nothing ahead of the target: no pull request is asked for, no tab opens…' and 'a pull request in any state already there: nothing starts…' (no spawn.sh, no jaira next)
- [x] Accept mit fälligem PR öffnet den PR-Tab wie bisher
  proof: test 'a branch ahead of the target without a pull request: the press opens one, by a person, in the ticket's worktree' unchanged and green
- [x] Der Toast nach Accept ohne PR nennt das akzeptierte Ticket und den PR-Hinweis, kein gestartetes Ticket
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:492 ('<id> accepted: done') + :596 / texts.ts:22 ('<id> accepted; <hint>'); tests 'two tickets on one branch…' and 'no gh…'
- [x] Tote Texte und Hilfsfunktionen des Next-Zweigs sind entfernt; claude plugin test . im Cockpit ist grün
  proof: nextToStart/dispatchLane removed from core/mod/builtin/jaira-cockpit/hooks/board.ts, nextStarted/nothingNext/nextNeedsHerdr from texts.ts; claude plugin test . in core/mod/builtin/jaira-cockpit: 175 pass, 0 fail, exit 0
- [x] Release-Note in core/release/notes/<ID>.md sagt: nach Accept startet der Teamlead das nächste Ticket
  proof: core/release/notes/X6JHQX.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-09 22:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35688.sock
- **2026-11-09 22:01 · Alexander Sacharov** — Dispatcher: keine offenen Entscheidungen vor dem Plan — Alex hat entschieden (Next-Zweig raus, PR-Zweig bleibt), Toast-Form steht im Kontext. Autonom.
- **2026-11-09 22:02 · Alexander Sacharov** — Teamlead: zuerst versehentlich als dispatcher gestartet (dispatch-Einstellung ist worker). Auf Alex' Wunsch beide Tabs geschlossen und als jaira-worker neu gestartet; in-progress wurde mitten in der Arbeit abgebrochen — git diff im Worktree prüfen.
- **2026-11-09 22:04 · Alexander Sacharov** — Worker übernimmt nach Neustart: claim --steal wurde vom auto-mode-Klassifikator verweigert, Claim liegt noch beim geschlossenen Dispatcher (DESKTOP-RFTCH11-64644). Weiter ohne neuen Claim. Diff des abgebrochenen in-progress übernommen und fertig gemacht. Entscheidung: Toast ohne PR heißt immer '<id> accepted', mit '; <Hinweis>' nur wenn es einen gibt (vorher kein Toast bei fehlendem Hinweis) — DoD 3 verlangt, dass der Toast das Ticket nennt. Step heißt jetzt 'none' statt 'next'. nextToStart/dispatchLane aus board.ts entfernt; settingsJson und tabsOf bleiben (anderswo genutzt). Overview-Test ruft weiter jaira next — gehört nicht zu Accept.
- **2026-11-09 22:05 · Alexander Sacharov** — Korrektur meiner obigen Entscheidung: ein Toast '<id> accepted' gibt es schon — accept() sagt t.accepted(id, lane) = '<id> accepted: done' vor afterAccept (register.tsx:492). Mein zweiter Schlüssel 'accepted' in texts.ts hat den doppelt belegt und überschrieben (bun meldet doppelte Keys nicht). Zurück zum Ansatz des Vorgängers: ohne PR nur ein zweiter Toast '<id> accepted; <Hinweis>', wenn es einen Hinweis gibt; sonst bleibt der erste Toast stehen.
- **2026-11-09 22:06 · Alexander Sacharov** — critique (Runde 1): core/mod/builtin/jaira-cockpit/hooks/register.tsx:535-538 beschreibt noch zwei Schritte ('never the board', 'stops both steps'); es gibt nur noch den PR-Schritt — Satz auf 'A string is a failure: no pull request is opened, and the pane says why' kürzen, 'never the board' streichen. Sonst nichts: spawn.sh --no-parent bleibt (Give back nutzt es, register.tsx:931), Doku in spawn.sh stimmt weiter; tabsOf/settingsJson bleiben (andere Aufrufer); 'pr'|'none' als Typ gelassen.
- **2026-11-09 22:06 · Alexander Sacharov** — in-progress Runde 2: Kommentar register.tsx:535-537 auf einen Schritt gekürzt, nur Kommentar, Tests unberührt.
- **2026-11-09 22:07 · Alexander Sacharov** — critique (Runde 2): einziger Befund aus Runde 1 (Kommentar register.tsx:535) behoben in 55a92796; nichts weiter.
