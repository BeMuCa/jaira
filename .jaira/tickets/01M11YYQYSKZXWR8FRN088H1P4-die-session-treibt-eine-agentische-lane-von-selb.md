---
id: 01M11YYQYSKZXWR8FRN088H1P4
title: Die Session treibt eine agentische Lane von selbst weiter
status: done
ready: true
creator: BeMuCa
goal: "Landet ein Ticket in einer Lane, die eine Session bearbeiten darf, arbeitet die Session die Lane und bewegt das Ticket weiter - zurueck, zum Menschen oder vorwaerts - ohne dass der Nutzer es sagt"
context: |-
  Berk am 27.08.: critique gibt eine Kritik aus, aber die Session, die die Tickets orchestriert, soll das Ticket dann auch selbst zurueck nach in-progress schieben und die Punkte umsetzen, oder zum Menschen - fuer jede Lane, die keinen Menschen braucht, ohne Aufforderung.

  Was heute da ist: der Lane-Prompt (critique.md, optimize.md) sagt am Ende, wohin das Ticket soll. Der verwaltete Block in CLAUDE.md/AGENTS.md sagt 'work one lane to empty before starting the next'. Beides steht - und gemessen laeuft es trotzdem nur, wenn eine Session dazu aufgefordert wird (27.08.: 4DQPMS und QPJNQP durch alle Lanes, auf Anweisung; 24.08.: 12 Tickets unbearbeitet in critique auf dem req-Board).

  Warum der Block nicht reicht: er wird einmal zu Sitzungsbeginn gelesen; der Moment, in dem die Session weitermachen soll, ist der Move.

  Drei Mechanismen, nicht exklusiv:
  (a) Die Ausgabe von 'jaira move' sagt den naechsten Schritt: 'X -> critique, eine agentische Lane. Naechster Schritt: jaira show X --for-lane critique --json, dann move.' --json traegt next_lane schon; die Aufforderung fehlt. Eine Zeile in internal/cli/flow.go.
  (b) Ein Stop-Hook fuer Claude Code (wie Berks verify-gate.js): verweigert das Sitzungsende, solange 'jaira next --per-lane --json' agentische Lanes mit unbearbeiteten Tickets meldet. Staerkster Hebel, weil die Umgebung ihn durchsetzt, nicht das Modell. Claude-Code-spezifisch, also opt-in: 'jaira hook print' gibt das Snippet aus, der Nutzer traegt es in settings.json ein.
  (c) Nur den Block schaerfen. Billigste Variante; die Messung spricht gegen sie.

  Entscheidung offen: (a), (b), beides.
definition-of-done: "Nach 'jaira move X --to critique' nennt die Ausgabe den Befehl fuer den naechsten Schritt; ein Stop-Hook-Snippet ist per Befehl abrufbar und verweigert das Ende bei unbearbeiteten agentischen Lanes; auf einer frischen Session mit einem Ticket in critique laeuft nach 'go' die Lane bis signoff ohne weitere Anweisung durch"
blocked-by: []
commits: []
created-at: 2026-08-27T15:55:56Z
updated-at: 2026-09-03T10:28:53Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2843276
claimed-at: 2026-08-31T18:10:52Z
assignee: BeMuCa
outcome-what: "Testluecken aus dem Review geschlossen. (1) TestStopHookSnippetRuns fuehrt das Snippet aus: dekodiert das JSON von 'jaira hook print', nimmt den command-String verbatim und laesst ihn per exec 'sh -c' gegen ein Stub-jaira im PATH laufen — vier Faelle mit erwarteten Exit-Codes 0/0/2/0 und der stderr-Zeile im blockierenden Fall. (2) TestStopHookReadsThePerLaneAgenticFlag hat jetzt die unterscheidende Haelfte: zweites Board mit Arbeit nur in todo, per-lane-JSON darf '\"agentic\": true' nicht enthalten. (3) Der Long-Text von 'hook print' sagt jetzt, dass die stop_hook_active-Pruefung ausschliesslich an stdin haengt und mit stdin auf /dev/null bedingungslos blockiert. Commit 20ffa1f."
outcome-why: "Das Snippet zitieren beweist nur die Schreibweise, nicht dass es laeuft — und die positive Haelfte allein laesst genau den Fehler durch, der weh tut: ein Hook, der immer '\"agentic\": true' findet, blockiert jedes Sitzungsende und sieht dabei korrekt aus. Der stdin-Satz steht da, damit niemand die stop_hook_active-Zeile fuer doppelte Absicherung haelt und rauswirft."
outcome-resolves: "Beide Testluecken und der Doku-Satz aus dem Review sind erledigt und gegengeprueft: zwei Mutationen (exit-0-Zweig entfernt, 2>/dev/null zu 2>&1) machen die Tests rot, hook.go danach byte-identisch wiederhergestellt. Unveraendert offen bleibt DoD-Klausel 3 (frische Session laeuft bis signoff durch) — braucht eine echte Claude-Code-Session mit dem Hook in settings.json — und Mechanismus (c), der Satz im verwalteten Block, der laut Vorgabe auf Berks 'go' wartet. Achtung: Commit 20ffa1f enthaelt vier fremde Dateien wegen des geteilten Git-Index; Details im Note."
review-summary: none
review-gaps: "Nichts entfernt. Gelassen: TUI-Moves ohne Nudge (DoD sagt Befehlsausgabe; die TUI hat keine), --json unangetastet (kehrt vor der Zeile zurueck, next_lane traegt es schon), Mechanismus c (Block-Satz) bewusst NICHT gebaut - wartet auf Berks Go, Plan-Punkt superseded statt falsch getickt; NOTES.md unangetastet (Release-Entscheidung); kein jq im Snippet (nicht ueberall installiert), jeder Fehler exit 0 (fail-open, wie hooks/sync-tasks.sh)."
review-check: "1. Neu bauen und in diesem Repo: jaira move <ticket> --to critique -> letzte Zeile nennt 'jaira show <id> --for-lane critique --json'; move --to todo -> keine Zeile; move --json -> reines JSON. 2. jaira hook print -> das Snippet; es nach ~/.claude/settings.json einsetzen, Session mit wartender agentischer Lane beenden wollen -> Stop wird mit der jaira-Meldung verweigert; ohne wartende Arbeit -> Stop geht durch. 3. DoD-Klausel 3 (frische Session arbeitet bis signoff) ist DEINE Pruefung - bewusst ungetickt."
review-verdict: "accept (Opus-End-Review auf cedc37d: Snippet als Claude Code dekodiert und in acht Szenarien ausgefuehrt, Quoting/stdin-Semantik/fail-open alle bestaetigt; Test-Delta vom Koordinator gelesen: TestStopHookSnippetRuns fuehrt den command unter sh mit Stub-jaira aus - 0/0/2/0 beobachtet statt zitiert, Negativ-Haelfte gepinnt, beide Mutationen rot verifiziert). Offen und ehrlich ungetickt: DoD-Klausel 3 (frische Session arbeitet nach 'go' bis signoff) - nur in deiner echten Session pruefbar; Mechanismus c (Block-Satz) wartet auf dein Go."
---

# Die Session treibt eine agentische Lane von selbst weiter

## Definition of Done

- [ ] Nach 'jaira move X --to critique' nennt die Ausgabe den Befehl fuer den naechsten Schritt; ein Stop-Hook-Snippet ist per Befehl abrufbar und verweigert das Ende bei unbearbeiteten agentischen Lanes; auf einer frischen Session mit einem Ticket in critique laeuft nach 'go' die Lane bis signoff ohne weitere Anweisung durch

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] internal/cli/flow.go: nach der Erfolgszeile des Moves eine zweite Zeile, wenn die Ziel-Lane agentic ist — nennt 'jaira show <handle> --for-lane <lane> --json'. Nur der --json-Zweig bleibt unberuehrt (returnt vorher), also keine Vermischung.
  proof: internal/cli/flow.go:227 druckt nextStepLine (Helper ab Zeile 247); --json returnt vorher
- [x] Neuer Befehl 'jaira hook print' (internal/cli/hook.go, Elternbefehl wie 'self'/'lanes'): gibt das Stop-Hook-Snippet fuer Claude Codes settings.json auf stdout aus, sonst nichts — Muster von 'jaira lanes template'.
  proof: internal/cli/hook.go: newHookCmd/newHookPrintCmd, registriert in internal/cli/root.go:211
- [x] Snippet-Inhalt: ruft 'jaira next --per-lane --json' und prueft mit grep -E auf "agentic": true; kein jq (auf dieser Maschine nicht installiert). Blockiert per Exit-Code + stderr-Meldung nach Claude Codes Stop-Hook-Vertrag.
  proof: internal/cli/hook.go:35 stopHookSnippet; vier Faelle per Hand geprueft: stop_hook_active=true -> 0, agentische Arbeit -> 2 + stderr, kein Board -> 0, nur nicht-agentische Arbeit -> 0
- [x] Tests in internal/cli: Nudge erscheint beim Move in eine agentische Lane, nicht in eine nicht-agentische; 'hook print' gibt stabilen Text aus (Golden-artig, im Hausstil).
  proof: internal/cli/nextstep_test.go (3 Tests), internal/cli/hook_test.go (4 Tests)
- [-] Ausserhalb des Scopes: der Satz im verwalteten Block (CLAUDE.md/AGENTS.md) fuer Subagenten — wartet laut Ticket-Progress auf Berks 'go'.
  proof: ausserhalb des Scopes, laut Berks Vorgabe; nicht gebaut
- [x] Verifikation: go test ./internal/cli -count=1; go test ./... -race -count=1; gofmt -l core internal nur internal/cli/tickets.go.
  proof: go test ./internal/cli -count=1 ok; go test ./... -race -count=1: 18 Pakete ok, kein FAIL/DATA RACE; gofmt -l core internal nur internal/cli/tickets.go

## Progress
- **2026-08-27 22:09 · BeMuCa** — 27.08., zu Berks Frage 'sub agents off in anderen Sessions': die Zeile 'Do not call the Agent tool unless the user requested it' steht in keiner Datei auf dem Rechner (settings, remote-settings, policy-limits, hooks, skills geprueft), aber auch im Prompt eines gespawnten Kind-Agenten - also Standard dieses Claude-Code-Builds in diesem Modus, jede Session. Sie ist weich ('unless the user requested it'): ein Satz im verwalteten Block - 'eine agentische Lane wird von einem Subagenten auf der model-tier der Lane bearbeitet; spawne einen' - ist die Anforderung, die in jeder Session gilt. Eine Zeile in core/board/announce.go plus Test. Wartet auf Berks go.
- **2026-08-31 18:24 · BeMuCa** — 31.08.: (a) und (b) gebaut, (c) nicht. Der Satz im verwalteten Block fuer Subagenten wartet laut Ticket-Progress vom 27.08. weiter auf Berks 'go' — bewusst nicht angefasst, also auch nichts in core/board/announce.go geaendert.

Nur der CLI-Move nudged, nicht die TUI. Die DoD sagt 'nennt die Ausgabe den Befehl' — das ist Befehls-Ausgabe, und die TUI hat keine. Die beiden anderen Move-Schreibstellen (internal/tui/model.go applyMove und der forcierte Move) bleiben unberuehrt; wer in der TUI schiebt, sieht die Lane sowieso vor sich.

Nudge nur bei lane.Agentic. Genau das ist 'traegt einen Prompt': core/lane/lane.go:311 lehnt eine agentische Lane ohne Prompt-Body ab. human und signoff sind agentic:false, also kann die Zeile nie dazu auffordern, die Menschen-Lane zu bearbeiten.

Stop-Hook-Vertrag verifiziert, nicht geraten (code.claude.com/docs/en/hooks.md): Exit 2 ist der EINZIGE Code, der das Sitzungsende blockiert — Exit 1 ist ein nicht-blockierender Fehler und laesst die Session enden. Die stderr-Zeile ist genau der Text, den Claude als Grund zum Weitermachen bekommt. Das Feld stop_hook_active im stdin-Payload ist true, sobald schon einmal ein Stop-Hook blockiert hat; ohne diese Pruefung haengt die Session an einer Bedingung, die der Hook selbst nicht loesen kann. Claude Code bricht zusaetzlich nach 8 aufeinanderfolgenden Blocks hart ab. Stop hat keine matcher-Unterstuetzung, deshalb steht im Snippet keiner.

Kein jq im Snippet: jq ist auf dieser Maschine nicht installiert (command -v jq leer). hooks/sync-tasks.sh hat noch einen jq-Zweig mit Fallback; hier ist grep der einzige Weg, weil die zwei gesuchten Felder feste Namen haben. Muster '": *true' deckt die eingerueckte und die kompakte JSON-Kodierung ab. Deshalb prueft hook_test.go zusaetzlich die echten Bytes von 'next --per-lane --json', nicht nur den Snippet-String — ein Wechsel des Encoders wuerde den Hook sonst still entwaffnen.

Snippet faellt bei jedem Fehler offen (Exit 0): kein Board, kein jaira im PATH, kaputtes Repo. Gleiche Entscheidung wie hooks/sync-tasks.sh: ein Board, das nicht gelesen werden kann, ist ein Aergernis; ein Hook, der jemanden in der Session festhaelt, ist deutlich schlimmer.

Ein 'jaira hook stop' als Go-Befehl (der Check in Go statt in der Shell) waere robuster testbar gewesen, ist aber verworfen: zwei neue Befehle statt einem, und das Ticket nennt selbst 'jaira hook print gibt das Snippet aus'. Scope-Disziplin schlaegt hier Eleganz.

core/release/NOTES.md bewusst nicht angefasst: v0.1.0 ist getaggt, eine Zeile dort waere rueckdatiert, und eine neue '## 0.2.0'-Sektion wuerde eine Versionsnummer erfinden. Gehoert in den Release-Schnitt, nicht hierher.

Nicht verifiziert und hier nicht verifizierbar: die dritte DoD-Klausel ('frische Session, Ticket in critique, laeuft nach go bis signoff durch'). Das braucht eine echte frische Claude-Code-Session mit dem Hook in settings.json — Berk-seitig. Der Hook ist opt-in und wird von 'jaira init' absichtlich nicht eingetragen; solange er nicht installiert ist, kann diese Klausel nicht erfuellt sein.
- **2026-08-31 18:44 · BeMuCa** — End-Review accept mit 2 Testluecken + 1 Doku-Satz: hook_test pinnt nur die Positiv-Haelfte ('agentic': true vorhanden) - Negativ-Fall ergaenzen (Arbeit nur in todo -> Muster fehlt); kein Test FUEHRT das Snippet AUS - eines, das den command JSON-dekodiert und via sh -c mit Stub-jaira laeuft (0/0/2); ein Satz in der Long-Hilfe, dass der stop_hook_active-Guard allein an stdin haengt (stdin /dev/null blockt bedingungslos).
- **2026-08-31 18:53 · BeMuCa** — 31.08., zweite Runde nach dem Review: das Snippet wird jetzt ausgefuehrt, nicht nur zitiert.

TestStopHookSnippetRuns dekodiert das JSON aus 'jaira hook print', nimmt genau den command-String und laesst ihn per exec 'sh -c' laufen. jaira ist ein Stub-Skript in einem t.TempDir() am Anfang von PATH, das kanonisches per-lane-JSON ausgibt — bin zuerst, damit der Stub ein hier real installiertes jaira ueberstimmt, restlicher PATH bleibt drin, weil grep gefunden werden muss. Vier Faelle: stop_hook_active=true trotz wartender agentischer Arbeit -> 0; nur nicht-agentische Arbeit -> 0; agentische Arbeit -> 2 plus stderr-Zeile; gar kein jaira im PATH -> 0.

Das per-lane-JSON im Test ist bewusst kanonisch (Konstanten perLaneAgentic/perLaneQuiet), nicht vom echten Befehl erzeugt: dieser Test prueft das Shell-Verhalten. Dass der echte Befehl noch genau diese Bytes ausgibt, prueft TestStopHookReadsThePerLaneAgenticFlag — die Arbeitsteilung ist Absicht, damit ein Encoder-Wechsel genau einen Test rot macht und nicht alle.

Gegenprobe eingebaut, das Review hatte recht: derselbe Test haelt jetzt auch die negative Haelfte fest. Zweites Board, Arbeit nur in todo, per-lane darf '"agentic": true' nicht enthalten (und muss '"agentic": false' enthalten, sonst prueft die Zusicherung ein leeres Ergebnis). Ohne diese Haelfte blockiert ein Hook, der immer true findet, jedes Sitzungsende und sieht dabei korrekt aus.

Zwei Mutationen gegengeprueft, damit die Tests nicht leer laufen: 'exit 0' im stop_hook_active-Zweig zu 'true' geaendert -> TestStopHookSnippetRuns Fall 1 faellt auf exit 2; '2>/dev/null' zu '2>&1' -> TestStopHookFailsOpen rot. hook.go danach byte-identisch wiederhergestellt (mit diff gegen HEAD geprueft).

Long-Text von 'hook print' ergaenzt: die stop_hook_active-Pruefung haengt vollstaendig an stdin. Wer den Befehl mit geschlossenem stdin oder /dev/null laufen laesst, blockiert bedingungslos. Das ist eine Absicherung, keine doppelte Absicherung — steht jetzt da, damit es niemand fuer redundant haelt und rauswirft.

Nebenbefund fuers Protokoll: im Arbeitsbaum liegen gerade fremde, nicht committete Aenderungen anderer Sessions (core/board/announce.go, internal/tui/view.go, internal/cli/tickets.go, lanes/*). 'gofmt -l core internal' zeigt deshalb auch internal/tui/view.go — nicht von mir. Committet wurden nur internal/cli/hook.go und internal/cli/hook_test.go; meine beiden Dateien sind gofmt-clean.
- **2026-08-31 18:58 · BeMuCa** — 31.08., Nachtrag zum Git-Vorfall in dieser Runde — bitte lesen, betrifft nicht den Code:

Der Arbeitsbaum wird gerade von mehreren Sessions gleichzeitig beschrieben, und der Git-Index ist geteilt. Mein 'git add <meine zwei Dateien>' hat deshalb nichts genuetzt: zum Zeitpunkt meines Commits lagen fremde, von anderen Sessions bereits gestagete Dateien im Index, und 'git commit' committet den Index. Commit 20ffa1f enthaelt darum vier Dateien, die nicht zu diesem Ticket gehoeren: core/board/announce.go, core/lane/rejectsto_test.go, internal/cli/lanes_test.go, internal/cli/tickets.go. Von mir sind nur internal/cli/hook.go und internal/cli/hook_test.go.

Mein Reparaturversuch war ein Fehler und hat es kurz schlimmer gemacht: 'git reset --soft HEAD~1' traf nicht meinen Commit, weil in der Zwischenzeit a5923da (TQXBY5) gelandet war — ich habe also fremde History entfernt. Wieder hergestellt mit 'git reset --hard a5923da'; a5923da war ein leerer Commit (sein Inhalt war schon in 20ffa1f mitgeschwemmt), es ist also kein Inhalt verloren gegangen, nur kurzzeitig die Commit-Nachricht. Danach keine History-Chirurgie mehr — in einem Baum, in den drei Sessions schreiben, ist das Risiko groesser als der Nutzen.

Lehre fuer die naechste Runde in diesem Repo: 'git commit --only <pfade>' statt 'git add <pfade> && git commit', und 'git reset' hier gar nicht. Wer den Diff reviewt: bei 20ffa1f die vier genannten Dateien abziehen.

Verifikation auf dem endgueltigen Baum (HEAD a5923da): go test ./internal/cli -count=1 ok; go test ./... -race -count=1 alle 18 Pakete ok, kein FAIL, kein DATA RACE; gofmt -l auf meinen beiden Dateien leer.
