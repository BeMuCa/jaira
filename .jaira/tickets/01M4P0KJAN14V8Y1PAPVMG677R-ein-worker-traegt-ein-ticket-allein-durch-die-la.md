---
id: 01M4P0KJAN14V8Y1PAPVMG677R
title: "Ein Worker traegt ein Ticket allein durch die Lanes, nur der Review laeuft frisch"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein einziger Agent fuehrt ein Ticket durch alle agentischen Lanes bis zur Review; eine Lane mit 'context: fresh' bekommt einen neuen Agenten ohne Vorwissen; Boards mit eigenen Lanes laufen ohne Migration weiter"
context: |-
  Heute startet jaira-dispatcher fuer JEDE Lane einen neuen Agenten (core/role/builtin/jaira-dispatcher). Das kostet pro Lane ca. 55-60k Tokens Kaltstart, auch fuer critique, die nur den Diff liest.
  A/B-Versuch am 11.10.2026 auf drei Tickets (GDB5YJ, 8WFNQY, K1W2JR), Report: https://claude.ai/artifact/7hu5TjpnChgLF6PkcuugfS
  - Arm A (ein Agent pro Lane): 17 Agenten, 1,11 Mio Tokens, 69 min.
  - Arm B (ein Agent pro Ticket, Review frisch): 6 Agenten, 450k Tokens, 55 min.
  - Blinder Richter: 2:1 fuer A, jedes Mal knapp; Code oft zeilengleich.
  - Fehler in A sassen an den Uebergaben: veralteter --proof nach optimize, --amend nach move, veraltete Kommentare, die keine Lane fixte.
  Arm B lief auf den UNVERAENDERTEN 13 Lanes: der Worker folgt next_lane und holt sich jede Lane mit 'jaira show --for-lane'. Lanes muessen also nicht migriert werden.
  Offen ist nur: woher weiss der Worker, welche Lane er NICHT selbst machen darf? Heute heisst sie review, auf fremden Boards anders.
  Vorschlag: optionales Frontmatter-Feld 'context: fresh' in der Lane-Datei (.jaira/lanes/*.md). Fehlt es ueberall, gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh. 'jaira update' bietet das Feld fuer die eingebaute review-Lane an, angepasste Lanes bleiben unberuehrt.
  Was erhalten bleiben muss: offene Entscheidungen vor dem Plan zaehlen und mode=conversational setzen, Neustart ueber 'jaira resume', notes an jeder Pause, Menschen-Lanes.
  Risiko aus dem Versuch: ein Langlauf-Agent hat einmal Tests mit relativem HOME laufen lassen, go brach ab, der FAIL-Grep sah es nicht. Deshalb bleibt der frische Reviewer Pflicht.
definition-of-done: "Lane-Frontmatter kennt 'context: fresh' (Schema, Parser, 'jaira lanes' zeigt es, 'show --for-lane --json' gibt es aus); ohne das Feld auf dem ganzen Board gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh; Tests decken beides"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-11T17:36:41Z
updated-at: 2026-10-11T18:31:48Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-521795
claimed-at: 2026-10-11T17:38:06Z
outcome-what: "Ungueltiger context-Wert laesst die Lane geladen und fresh, zaehlt als deklariert, Warnung + 'fresh!' in jaira lanes + validate-Fehler (Exit 3); NOTES ehrlich ueber den Agenten-Block"
outcome-why: "Review-Befund: Tippfehler verschob die fresh-Grenze auf in-progress, NOTES behauptete 'nichts aendert sich'"
outcome-resolves: "DoD 7 und 8"
review-summary: "Neues Lane-Feld context: fresh|shared (core/lane/lane.go parse+impliedFresh): jeder andere Wert ist ein Parse-Fehler; setzt keine Lane das Feld, wird die letzte agentische Lane vor der ersten requires-human-exit-Lane implizit fresh (Standard-Board: review; HITL zaehlt nicht). Sichtbar in jaira lanes (Spalte CONTEXT, fresh* + Fussnote), lanes show, show --for-lane (json context) und im generierten Agenten-Block (review: 'fresh context - hand it to a new agent'); laneFacts aus CLI und TUI zu lane.Set.Facts vereint. Opt-in-Schalter dispatch=lanes|worker in core/settings (Default lanes, Tippfehler -> lanes), umschaltbar im TUI-Settings-Screen ('Agents per ticket'), lesbar via jaira roles list (+json dispatch/dispatch_role). Neue Rolle jaira-worker (traegt das Ticket selbst durch alle Lanes, zaehlt offene Entscheidungen vor dem Plan, mode=conversational inkl. lesender critique, Neustart ueber resume/notes, fresh-Lane per run-lane.sh/spawn.sh im eigenen Herdr-Tab, sonst Agent()-Fallback). spawn.sh kennt Lane-Name worker -> /jaira-worker. jaira-dispatcher/SKILL.md unveraendert; teamlead liest den Schalter und startet das Passende. Zwei NOTES.md-Zeilen unter Unreleased."
review-gaps: "1. Tippfehler im context-Feld verschiebt die fresh-Grenze: review.md mit 'context: Fresh' -> Lane faellt mit stderr-Warnung weg, der Default markiert dann in-progress als fresh*, die Fussnote von jaira lanes behauptet faelschlich 'no lane on this board declares context:', jaira validate nennt es nicht und endet mit 0. Im Worker-Modus ginge damit die Implementierung an einen frischen Agenten und keine unabhaengige Pruefung fande statt. Das Wegfallen der Lane ist bestehendes Parse-Fehler-Verhalten (precedence genauso), die verschobene Markierung ist neu. Billige Abhilfe: ein gesetzter, aber ungueltiger Wert zaehlt als 'deklariert' (Default greift nicht) oder wird mit Warnung ignoriert statt die Lane zu verwerfen. Kleiner Defekt, kein Blocker. 2. Opt-in ist nicht vollstaendig: der generierte Agenten-Block (CLAUDE.md/AGENTS.md) markiert review auf JEDEM Board als 'fresh context - hand it to a new agent' statt 'yours to work', unabhaengig von dispatch; nach dem Upgrade meldet jaira validate den Block als veraltet bis jaira update. Fuer jaira-dispatcher aendert das nichts (neuer Agent pro Lane ohnehin), eine Session ohne Rollen, die das Board laut Block selbst faehrt, soll review nun aber abgeben. NOTES-Zeile 1 'Without the setting nothing changes' ist deshalb leicht ueberzogen. 3. NOTES.md bekommt zwei Zeilen statt einer; beide einzeilig, Instruktionsform, die erste nennt die Einstellung. 4. spawn.sh worker nur ueber Herdr-Stub getestet; ohne Herdr bricht es sauber ab ('not inside a Herdr pane')."
test-verdict: "pass: go build ./... und go test ./... -race -count=1 mit isoliertem absolutem HOME RC=0; DoD 1-6 im Baum geprueft (Proofs geoeffnet, Tests laufen); auf einem Scratch-Board mit gebautem Binary gesehen: jaira lanes zeigt review als fresh* mit Fussnote, lanes show review 'Context: fresh — ...', show --for-lane review --json context=fresh und Kopfzeile 'context: fresh', roles list 'dispatch: lanes ... /jaira-dispatcher', nach {\"dispatch\":\"worker\"} --json worker/jaira-worker, context: shared auf review -> keine fresh-Lane, context: fresch -> Ladewarnung"
review-verdict: "DoD durch den Diff gedeckt, go test ./... -race mit isoliertem HOME exit 0, Scratch-Board verhaelt sich wie beschrieben. Kein blockierender Defekt; Luecke 1 (Tippfehler im context verschiebt fresh auf in-progress, irrefuehrende Fussnote) ist ein kleiner echter Defekt, den eine Person entweder als Nacharbeit zurueckschickt oder akzeptiert. Freigabe mit Vorbehalt."
review-check: "1. cd /home/alex/projects/jaira/.claude/worktrees/mg677r && go build -o /tmp/jn ./cmd/jaira  2. H=$(mktemp -d); B=$(mktemp -d); cd $B && git init -q && git commit -q --allow-empty -m i; export HOME=$H JAIRA_HOME=$H  3. /tmp/jn init  4. /tmp/jn lanes -> Spalte CONTEXT, review zeigt fresh*, darunter die Fussnote  5. /tmp/jn roles list -> letzte Zeile 'dispatch: lanes ... /jaira-dispatcher'  6. echo '{\"dispatch\":\"worker\"}' > $H/settings.json; /tmp/jn roles list -> 'dispatch: worker ... /jaira-worker'  7. /tmp/jn (Board) -> S druecken -> Eintrag 'Agents per ticket' zeigt worker, enter schaltet auf lanes  8. Luecke 1 sehen: in .jaira/lanes/review.md unter agentic: die Zeile 'context: Fresh' einfuegen, /tmp/jn lanes -> review fehlt, in-progress zeigt fresh*  9. /tmp/jn roles install --project; ls .claude/skills -> jaira-worker ist dabei, jaira-dispatcher/SKILL.md unveraendert (git diff master -- core/role/builtin/jaira-dispatcher/SKILL.md ist leer)"
---

# Ein Worker traegt ein Ticket allein durch die Lanes, nur der Review laeuft frisch

## Definition of Done

- [x] Lane-Frontmatter kennt 'context: fresh' (Schema, Parser, 'jaira lanes' zeigt es, 'show --for-lane --json' gibt es aus); ohne das Feld auf dem ganzen Board gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh; Tests decken beides
  proof: core/lane/lane.go parse+impliedFresh; internal/cli/tickets.go lanes CONTEXT column, flow.go for-lane context; tests core/lane/context_test.go, internal/cli/context_test.go TestLanesAndForLaneReportContext, core/board/lanesection_test.go TestLaneSectionMarksAFreshLane
- [x] Rolle jaira-worker in core/role/builtin als Alternative zu jaira-dispatcher (der unveraendert bleibt): traegt das Ticket selbst durch alle Lanes bis vor die erste fresh-Lane, startet dafuer einen frischen Agenten mit jaira-role-lane, und behaelt Zaehlen offener Entscheidungen, mode=conversational, Neustart ueber resume und notes
  proof: core/role/builtin/jaira-worker/SKILL.md (Before the plan, conversational mode, The fresh lane, restartable); TestBuiltinsAreTheShippedRoles, TestWorkerNamesScriptsThatShip
- [x] jaira-teamlead liest die Einstellung und startet jaira-worker oder jaira-dispatcher; ohne Einstellung bleibt alles wie heute (dispatcher, ein Agent pro Lane)
  proof: core/role/builtin/jaira-teamlead/SKILL.md 'Which one is the person's choice' reads jaira --json roles list; internal/cli/roles.go; TestRolesListReportsTheDispatchSetting
- [x] Zeile in core/release/NOTES.md unter Unreleased; go test ./... -race gruen (mit isoliertem HOME, siehe bekannte Basis-Fehler)
  proof: core/release/NOTES.md Unreleased (2 Zeilen); go test ./... -race -count=1 mit HOME/USERPROFILE/JAIRA_HOME = absolutes leeres Temp-Verzeichnis: exit 0
- [x] Opt-in-Schalter in core/settings: ~/.jaira/settings.json "dispatch": "lanes" | "worker", Default "lanes"; im TUI-Einstellungsbildschirm (internal/tui/settings.go) umschaltbar; ueber ein bestehendes CLI-Kommando lesbar (jaira roles list); NOTES.md nennt die Einstellung und wie man sie einschaltet
  proof: core/settings/settings.go Dispatch/DispatchMode/DispatchRole (TestDispatchDefaultsToLanes); internal/tui/settings.go 'Agents per ticket' (TestSettingsEnterSwitchesDispatch); jaira roles list dispatch/dispatch_role; NOTES.md Unreleased line 2
- [x] Herdr bleibt: im Worker-Modus bekommt der Worker einen eigenen Herdr-Tab (Teamlead via spawn.sh, Lane-Name 'worker'), der frische Agent der fresh-Lane ebenfalls einen eigenen Tab via spawn.sh/run-lane.sh, nach seiner Lane geschlossen; ohne HERDR_ENV=1 gilt der bisherige Fallback
  proof: jaira-dispatcher/scripts/spawn.sh lane 'worker' -> /jaira-worker (TestSpawnTypesTheRoleForItsLaneName, Herdr-Stub); jaira-worker/SKILL.md 'Starting another agent' uses run-lane.sh/spawn.sh, Fallback Agent()
- [x] Ein gesetzter, aber ungueltiger context-Wert (Tippfehler wie 'Fresh') verschiebt die fresh-Grenze nicht: die Lane bleibt geladen bzw. der Wert zaehlt als deklariert, mit sichtbarer Warnung; der Default springt nicht auf in-progress; 'jaira lanes' behauptet nicht, keine Lane deklariere context; 'jaira validate' meldet ihn mit Exit != 0; Tests decken es
  proof: core/lane/lane.go parse: ContextInvalid, liest als fresh; Load-Warning ContextInvalidMessage; internal/cli/tickets.go laneContext 'fresh!' + invalidContextNote; internal/cli/validate.go lane_context error (Exit 3); Tests TestContextUnknownValueRunsFresh, TestContextTypoKeepsTheFreshBoundary, TestAnInvalidContextIsReportedNotDropped
- [x] Agenten-Block in CLAUDE.md/AGENTS.md und NOTES.md widersprechen sich nicht: entweder stimmt die 'fresh context'-Markierung fuer beide dispatch-Modi oder die NOTES-Zeile sagt ehrlich, was sich ohne Einstellung aendert; Begruendung der Wahl als Note
  proof: core/release/NOTES.md Unreleased Zeile 1 ('the roles work as before') und Zeile 2 ('every session that drives the board itself hands that lane to a new agent, whatever dispatch is set to'); Begruendung in Note Runde 2

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] lane.Lane.Context: parse (fresh|shared), Fehler fuer anderes, lanesEquivalent, Template
- [x] Set: Default-Regel in Load nach der Reihenfolge (ContextImplied)
- [x] Tests in core/lane fuer Feld, Fehler und Default (inkl. shared schaltet ab, Board ohne Menschen-Lane)
- [x] CLI: jaira lanes (Tabelle+JSON), lanes show, show --for-lane --json context; Tests
- [x] announce: fresh-Lane in der Lane-Liste markieren
- [x] Rolle jaira-worker schreiben (nutzt die Skripte von jaira-dispatcher), spawn.sh-Lane 'worker', role_test anpassen
- [x] teamlead liest Schalter; role-lane-Wortlaut; README/docs/roles.go
- [x] NOTES.md Zeile; go test ./... -race mit isoliertem HOME
- [x] Schalter dispatch in core/settings + TUI-Eintrag + 'jaira roles list' (Text+JSON), Tests

## Progress
- **2026-10-11 17:40 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan — alle selbst entschieden, keine braucht den Menschen:
1. 'erste Menschen-Lane' = erste Lane mit requires-human-exit, Lanes die nur eine Frage nehmen (requires-question, HITL) zaehlen nicht — wie lane.Next sie ueberspringt. Woertlich genommen waere HITL die erste und testing fresh; das widerspricht dem Ziel (review frisch). Auf dem Standard-Board ergibt die Regel review.
2. Werte: fresh oder shared. shared ist noetig, damit ein Board ausdruecklich KEINE frische Lane haben kann (sonst greift immer der Default). Leer = nicht gesetzt. Jeder andere Wert ist ein Parse-Fehler wie ein falscher precedence.
3. Der Default greift nur, wenn KEINE Lane des Boards context: setzt.
4. Eingebaute review.md bleibt unveraendert, 'jaira update' bietet nichts an: Default ergibt dasselbe, und ein Feld in der eingebauten Lane wuerde die review.md jedes bestehenden Boards als 'weicht ab' (Overrides) markieren. Nicht in der DoD.
5. jaira-dispatcher bleibt als Stub-SKILL.md, der auf /jaira-worker verweist; spawn.sh/run-lane.sh ziehen nach jaira-worker/scripts. spawn.sh: Lane-Name 'worker' startet /jaira-worker, 'dispatch' bleibt als Alias.
6. Conversational-Modus: die lesende critique neben der Arbeit bleibt, der Worker startet sie.
7. 'jaira lanes': neue Spalte CONTEXT; abgeleitetes fresh mit * und Fussnote. JSON: context + context_implied. show --for-lane --json: context.
8. Die generierte CLAUDE.md-Lane-Liste markiert eine fresh-Lane, damit jeder Agent es liest.
- **2026-10-11 17:43 · Alexander Sacharov** — Teil 1 (Lane-Feld) committed in 8f1707b. Naechstes: Rolle jaira-worker.
- **2026-10-11 17:45 · Alexander Sacharov** — Scope-Erweiterung von Alex am 2026-10-11 (ueber den Koordinator): (1) Worker-Modus ist opt-in, Default bleibt der Dispatcher mit einem Agenten pro Lane; Schalter 'dispatch' in ~/.jaira/settings.json, im TUI-Settings-Screen, per CLI lesbar; Teamlead startet das Passende. (2) Herdr bleibt: Worker und frischer Review-Agent bekommen je einen eigenen Tab via spawn.sh. DoD 2 und 3 entsprechend umformuliert (vorher: worker ERSETZT dispatcher, alter Name nur noch Hinweis) — das widersprach der Erweiterung. Entscheidung: eigene Rolle jaira-worker statt Zweig im Dispatcher — der Dispatcher-Prompt bleibt byte-gleich, also aendert sich fuer niemanden ohne Schalter etwas, und zwei Arbeitsweisen in einer Datei wuerden jeden Absatz an eine Bedingung haengen. Der Worker benutzt die Skripte von jaira-dispatcher/scripts (gleiches Skills-Verzeichnis, roles install schreibt beide). CLI-Leseweg: 'jaira roles list' (--json: Schluessel dispatch und dispatch_role) — die Stelle, an der der Teamlead ohnehin die Rollen sieht.
- **2026-10-11 17:52 · Alexander Sacharov** — critique Runde 1: internal/cli/tickets.go:1427 impliedContextNote/laneContext stehen zwischen dem Doc-Kommentar von laneFacts und laneFacts — der Kommentar haengt jetzt an der Konstante; Block ueber den Kommentar ziehen
internal/tui/settings.go switchDispatch gibt settingsActionSaved zurueck und loest m.reload() aus, dessen Kommentar (Logbuch-Fenster) nicht passt und das fuer dispatch nichts tut; settingsActionNone zurueckgeben
core/lane/lane.go impliedFresh-Kommentar: eine Zeile ueber 150 Zeichen; umbrechen. internal/cli/roles.go Long: eine ueberlange Zeile; umbrechen. core/board/announce.go: verirrte Leerzeile nach dem tier-Block entfernen
core/role/builtin/jaira-role-lane/SKILL.md: 'A worker starts you for a lane marked fresh' klingt wie der einzige Fall — der Worker startet auch die lesende critique; 'When a worker starts you for ...' schreiben
- **2026-10-11 17:54 · Alexander Sacharov** — critique Runde 2: alle fuenf Funde aus Runde 1 in 5b01536 behoben (laneFacts-Kommentar, settingsActionNone, zwei Umbrueche, Leerzeile, role-lane-Wortlaut); die Reparatur fuehrt nichts Neues ein. Nichts mehr zu sagen.
- **2026-10-11 17:57 · Alexander Sacharov** — optimize: laneFacts gab es zweimal (internal/cli/tickets.go und internal/tui/browse.go); die TUI-Kopie kannte Fresh nicht, sodass ein im TUI neu geschriebener Agenten-Block die fresh-Markierung verlor — beide durch lane.Set.Facts in core/lane/facts.go ersetzt. Lane.Fresh() entfernt (Einzeiler, ein Aufrufer). Gelassen: laneContext und impliedContextNote (zwei Aufrufer bzw. echte Logik), Settings.DispatchRole (zwei Aufrufer in roles.go, getestet). Suite mit isoliertem HOME gruen.
- **2026-10-11 17:58 · Alexander Sacharov** — testing: bestanden. Beobachtet, nicht als Fehler gewertet: ein ungueltiger context-Wert laesst die Lane wie jeder andere Parse-Fehler (z. B. precedence) mit Warnung fallen; fehlt dadurch review, greift der Default und macht in-progress zur fresh-Lane. Das Board ist dann ohnehin kaputt und warnt laut; ein Fehler-tolerantes Parsen waere eine eigene Designentscheidung. spawn.sh nur ueber den Herdr-Stub in TestSpawnTypesTheRoleForItsLaneName geprueft, kein echter Tab geoeffnet (Auftrag).
- **2026-10-11 17:58 · Alexander Sacharov** — Worker-Lauf bis vor review abgeschlossen: in-progress x2, critique x2 (1 Fund-Runde), optimize x1, testing x1. review ist die fresh-Lane und gehoert einem neuen Agenten.
- **2026-10-11 18:03 · Alexander Sacharov** — review (frischer Agent): Freigabe mit Vorbehalt. Kleiner Defekt: ungueltiger context-Wert (z. B. 'Fresh') verwirft die Lane, der Board-Default springt dann auf in-progress, und die lanes-Fussnote sagt faelschlich, keine Lane deklariere context. Ausserdem markiert der Agenten-Block review auf jedem Board als fresh, auch ohne dispatch=worker. Details in review-gaps.
- **2026-10-11 18:03 · Alexander Sacharov** — Review-Befund, Rueckgabe empfohlen: ein gesetzter, aber ungueltiger context-Wert (Tippfehler wie 'Fresh') wirft die review-Lane raus, der Board-Default macht dann in-progress zur fresh-Lane, 'jaira lanes' und 'jaira validate' melden nichts. Fix: ungueltiger Wert zaehlt als deklariert bzw. die Lane bleibt mit Warnung geladen, der Default verschiebt sich nicht; validate meldet ihn mit Exit != 0. Zweitens: der generierte Agent-Block markiert review auch ohne dispatch=worker als 'fresh context' - NOTES-Zeile 'ohne Einstellung aendert sich nichts' entsprechend ehrlich machen.
- **2026-10-11 18:30 · Alexander Sacharov** — Runde 2, Entscheidungen: (1) Ungueltiger context-Wert (z. B. 'Fresh') wirft die Lane nicht mehr raus: Parse behaelt sie, liest den Wert als fresh (ContextInvalid haelt den Rohwert), das zaehlt als deklariert, also greift der Board-Default nicht. Warum fresh und nicht shared: fresh ist die sichere Lesart - ein Tippfehler kostet dann einen neuen Agenten, nie die unabhaengige Pruefung. Gemeldet in Load-Warnings (stderr), 'jaira lanes' Spalte 'fresh!' + Fussnote, lanes --json context_invalid, 'jaira validate' als error lane_context (Exit 3). (2) Agenten-Block: NOTES ehrlich gemacht statt Block an dispatch zu koppeln. Grund: der Block ist Board-Datei im Repo, dispatch ist pro Person (~/.jaira/settings.json) - ein vom Schalter abhaengiger Block wuerde zwischen Teammitgliedern hin und her kippen und validate staendig 'veraltet' melden. Die Markierung ist zudem in beiden Modi wahr: der Dispatcher startet fuer review ohnehin einen neuen Agenten; nur eine Session ohne Rollen gibt review jetzt ab - genau das soll sie (Ziel: frischer Reviewer Pflicht). NOTES Zeile 1/2 sagen das jetzt.
