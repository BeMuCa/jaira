---
id: 01M7157NFTE0AME0MS7MX6JHQX
title: Accept im Cockpit startet kein nächstes Ticket mehr
status: done
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
commits:
  - 2ec06416a413df92811505ab6079af71dd8b4aa0
  - 55a92796fd92cad7bfa82fb41fe6bc1d42549bbb
  - 85f42ad0b704ba75d654e59253c9ea5e5313be5c
  - 35415ce22a9a932bf30d8495d2f629f725273707
created-at: 2026-11-09T22:00:36Z
updated-at: 2026-11-09T22:31:03Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-64644
claimed-at: 2026-11-09T22:01:08Z
outcome-what: "Accept im Cockpit öffnet nur noch den PR, wenn einer fällig ist; der Zweig, der das nächste Ticket per jaira next + spawn.sh in einem Herdr-Tab startete, ist entfernt, samt nextToStart, dispatchLane und drei Texten"
outcome-why: "Das automatisch gestartete Ticket kam ohne Ankündigung und ohne Grund; welches Ticket als Nächstes läuft, entscheidet der Teamlead"
outcome-resolves: Accept ohne fälligen PR startet nichts mehr; Toast nennt nur Ticket und PR-Hinweis
review-summary: "Accept im Cockpit startet kein nächstes Ticket mehr. afterAccept() (register.tsx:592) kennt nur noch einen Schritt: ist ein PR fällig (Branch vor dem Ziel, alle Tickets darauf in done, noch kein PR), öffnet es ihn wie bisher über spawn.sh pr; sonst bleibt der Toast '<id> accepted: done', und nur wenn es einen Grund gibt, warum kein PR geöffnet wurde (andere Tickets nicht in done, Forge nicht erreichbar), kommt '<id> accepted; <Grund>' dazu. Weg sind der Aufruf 'jaira next --all', 'herdr tab list', spawn.sh --no-parent, nextToStart/dispatchLane (board.ts) und die Texte nextStarted/nothingNext/nextNeedsHerdr. stepAfterAccept liefert 'pr' | 'none'. Ohne fälligen PR braucht Accept auch kein Herdr mehr. Tests in board.test.ts umgestellt (kein SPAWN, kein jaira next), Release-Note X6JHQX.md, drei Hashes in shipped.sha256."
round: 3
review-gaps: |-
  note: Kommentar register.tsx:583 hat eine überlange Zeile ('… where there is a reason. A step that cannot run') — nur Umbruch, kein Inhalt
  note: worktree enthält das unversionierte .jaira/milestones/abnahme-2026-11-A4140M.md — gehört nicht zu diesem Ticket, nicht mit committen
test-verdict: "pass: go build ./... ok, go test ./... RC=0, claude plugin test . im Cockpit 175/0 RC=0; DoD 1-5 im Baum geprüft (Next-Zweig, Helfer und Texte weg, PR-Zweig-Tests unverändert grün, Toast-Tests, Release-Note eine '- ! '-Zeile). Live im Cockpit nicht gedrückt — nur über die Press-Simulation der Tests"
review-verdict: "Diff erfüllt alle fünf DoD-Punkte; keine Defekte gefunden. claude plugin test . 175/0, go test ./core/role ok. Der PR-Zweig ist unverändert (nur aus dem if herausgezogen), der Hinweis ist im Fehlerpfad zu Recht weg, weil decideStep ihn nur bei step 'none' setzt."
review-round: 3
review-check: |-
  1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  → 175 pass, 0 fail
  2. In einem Repo mit Cockpit ein Ticket in signoff nehmen, dessen Branch schon einen PR hat (oder nichts vor master ist), im Cockpit Accept drücken
  3. Erwartet: Toast '<ID> accepted: done', kein neuer Herdr-Tab, kein anderes Ticket gestartet
  4. Ein Ticket nehmen, dessen Branch vor master liegt und noch keinen PR hat, Accept drücken
  5. Erwartet: wie früher ein Herdr-Tab '<id>/pr' und Toast '<ID>: opening the pull request for <branch> in a Herdr tab'
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
- **2026-11-09 22:09 · Alexander Sacharov** — testing (Runde 1): go build ./... ok; claude plugin test . 175/0; go test ./... exit 1: core/role shipped_test.go:74 TestEveryInstalledFileIsAShippedVersion — core/shipped/shipped.sha256 fehlt die Version von hooks/board.ts, register.tsx, texts.ts. Fix: scripts/shipped-history.sh > core/shipped/shipped.sha256 bzw. die drei Zeilen anhängen, im selben Commit wie der Code.
- **2026-11-09 22:11 · Alexander Sacharov** — in-progress Runde 3: drei Hashes an core/shipped/shipped.sha256 angehängt (append-only laut Dateikopf; scripts/shipped-history.sh hätte ~40 fremde Altversionen mitgezogen). go test ./core/role ok.
- **2026-11-09 22:11 · Alexander Sacharov** — critique (Runde 3): nur die Antwort auf testing gelesen (3 angehängte Hashes in shipped.sha256, wie 736c80ee) — nichts.
- **2026-11-09 22:11 · Alexander Sacharov** — optimize (Runde 2): nur shipped.sha256 neu, nichts zu entfernen; review-gaps aus Runde 1 gilt weiter.
- **2026-11-09 22:14 · Alexander Sacharov** — review: keine Blocker; zwei Notes (Kommentar-Umbruch register.tsx:583, fremde Datei .jaira/milestones/ unversioniert im Worktree). Weiter nach signoff.
- **2026-11-09 22:30 · Alexander Sacharov** — explain: Der Review fand keine Fehler; Accept öffnet nur noch den fälligen PR und startet nie ein nächstes Ticket. Zwei Kleinigkeiten: eine überlange Kommentarzeile in register.tsx:583 und eine fremde, unversionierte Datei .jaira/milestones/abnahme-2026-11-A4140M.md im Worktree.
Warum: Beides ändert kein Verhalten; die Kommentarzeile ist nur Optik, die Milestone-Datei gehört nicht zu diesem Ticket und darf nur nicht mit committet werden.
1) Zurückgeben: Kommentar in register.tsx:583 umbrechen, danach erneut prüfen [give back]
2) Ticket so annehmen, die Milestone-Datei bleibt uncommittet [accept]
Empfehlung: 2 - alle fünf DoD-Punkte sind erfüllt, Tests laufen 175/0, und der Umbruch ist keine eigene Runde wert.
- **2026-11-09 22:30 · Alexander Sacharov** — decision: 2) Ticket so annehmen, die Milestone-Datei bleibt uncommittet
