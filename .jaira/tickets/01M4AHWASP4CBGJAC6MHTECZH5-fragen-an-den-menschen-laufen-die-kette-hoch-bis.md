---
id: 01M4AHWASP4CBGJAC6MHTECZH5
title: Fragen an den Menschen laufen die Kette hoch bis zum Teamlead
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Der Mensch beantwortet jede Frage in der einen Sitzung, die er selbst gestartet hat: ein Worker gibt seine Frage an den Dispatcher, ein Dispatcher mit Teamlead an den Teamlead; nur die oberste Sitzung fragt selbst."
context: |-
  Was falsch ist: Worker und Dispatcher laufen in eigenen Herdr-Tabs und haben dort AskUserQuestion. Jeder fragt in seinem Tab. Alex muss suchen, welcher Tab ihn gerade etwas gefragt hat.
  Heute so geregelt: core/role/builtin/jaira-dispatcher/SKILL.md, Abschnitt 'How you ask' (~Z.164): der Dispatcher fragt selbst mit AskUserQuestion; nur ohne AskUserQuestion (Subagent, headless) gibt er nummerierte Optionen in seinem Bericht nach oben. In einem Herdr-Tab hat er AskUserQuestion — also fragt er immer selbst. jaira-role-lane sagt gar nichts dazu, ein Worker fragt also auch selbst.
  Der Teamlead (core/role/builtin/jaira-teamlead/SKILL.md, Punkt 4) kann schon Optionen eines Dispatchers als Wahl stellen und die Antwort mit 'jaira note' aufs Ticket schreiben — er bekommt sie nur nie, weil der Dispatcher im Tab selbst fragt.
  Regel (mit Alex abgestimmt, 2026-10-07): wer von einer anderen Sitzung gestartet wurde, fragt den Menschen nie selbst. Er schreibt die Frage aufs Ticket (jaira note, Optionen nummeriert, Empfehlung zuerst), gibt sie an seinen Starter und wartet. Nur wer keinen Starter hat — der Teamlead, oder ein Dispatcher, den der Mensch direkt gestartet hat — fragt mit AskUserQuestion.
  Woher eine Rolle ihren Starter kennt: spawn.sh (core/role/builtin/jaira-dispatcher/scripts/spawn.sh:145-149) schickt heute nur '/jaira-dispatcher <id>' bzw. '/jaira-role-lane <id> <lane>'. Es muss den Namen des Starters mitgeben, z.B. als Argument.
  Weg nach oben/unten, zu prüfen: SendMessage/ListAgents zwischen lokalen Sitzungen (Herdr-Tabs sind solche); die Frage liegt zusätzlich auf dem Ticket, damit sie einen gestorbenen Tab überlebt. Die Antwort geht als jaira note aufs Ticket und als Nachricht zurück nach unten.
  Grenze: Freigabe-Dialoge des Harness (Permission-Prompts) in einem Worker-Tab lassen sich nicht weiterreichen. Der Teamlead kann nur sagen, in welchem Tab einer wartet — das bleibt so.
  Verwandt, nicht dasselbe: 4XHZ6N (Zeilen für den Menschen pro Lane), ZWK0PF (offene Frage auf der Karte sichtbar).
definition-of-done: "jaira-role-lane: ein Worker fragt nie selbst; er schreibt die Frage als nummerierte Optionen aufs Ticket, gibt sie in seinem Bericht an den Dispatcher und hält an"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-07T06:47:41Z
updated-at: 2026-10-07T06:52:22Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-61773
claimed-at: 2026-10-07T06:50:35Z
---

# Fragen an den Menschen laufen die Kette hoch bis zum Teamlead

## Definition of Done

- [ ] jaira-role-lane: ein Worker fragt nie selbst; er schreibt die Frage als nummerierte Optionen aufs Ticket, gibt sie in seinem Bericht an den Dispatcher und hält an
- [ ] jaira-dispatcher: mit Starter reicht er jede Frage (eigene und die seiner Worker) an den Starter weiter und wartet; ohne Starter fragt er selbst wie heute
- [ ] jaira-teamlead: nimmt weitergereichte Fragen an, stellt sie als Wahl, schreibt die Antwort mit jaira note aufs Ticket und gibt sie an den Dispatcher zurück
- [ ] spawn.sh gibt dem gestarteten Dispatcher bzw. Worker den Namen seines Starters mit
- [ ] Ein Probelauf Teamlead -> Dispatcher -> Worker in Herdr: eine Frage des Workers erscheint nur im Teamlead-Tab, die Antwort kommt beim Worker an
- [ ] core/release/NOTES.md hat eine Zeile unter Unreleased
- [ ] Die Skills sagen die Richtung: der Teamlead ist der Ort, an dem der Mensch antwortet; ein direkt gestarteter Dispatcher ist der Ausnahmefall
- [ ] Schreibt der Mensch direkt in den Dispatcher-Tab, antwortet der Dispatcher dort, hält das Entschiedene mit jaira note fest und meldet es seinem Teamlead; der Teamlead fragt es nicht noch einmal

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-07 06:49 · Alexander Sacharov** — Entscheidung Alex 2026-10-07: mit einem Worker will er nie sprechen — ein Worker fragt ausnahmslos seinen Dispatcher, nie den Menschen. Mit dem Dispatcher geht beides: ohne Teamlead fragt er selbst, mit Teamlead reicht er weiter. Richtung auf Dauer: alles läuft über den Teamlead.
- **2026-10-07 06:49 · Alexander Sacharov** — Ergänzung Alex 2026-10-07: es wird trotzdem Momente geben, in denen er direkt in den Dispatcher-Tab schreibt. Das ist erlaubt, kein Fehler. Der Dispatcher antwortet dann dort, schreibt das Entschiedene mit jaira note aufs Ticket und meldet es an seinen Teamlead weiter — sonst arbeitet der Teamlead mit einem veralteten Stand und fragt dasselbe noch einmal.
- **2026-10-07 06:51 · Alexander Sacharov** — Präzisierung (Vorschlag an Alex, 2026-10-07): spricht der Mensch selbst im Dispatcher-Tab, antwortet der Dispatcher dort; eine Rückfrage, die aus diesem Gespräch entsteht, stellt er auch dort. Fragen, die aus seiner eigenen Arbeit entstehen (Worker hängt, dritte Runde einer Lane), gehen immer an den Teamlead. Was im Tab entschieden wurde: jaira note + eine Zeile per SendMessage an den Teamlead. Mechanik geprüft: ListAgents nennt jeder Sitzung ihren eigenen Namen ('This session is jaira-b4'), SendMessage erreicht lokale Sitzungen über diesen Namen — spawn.sh kann ihn als Starter durchreichen.
- **2026-10-07 06:52 · Alexander Sacharov** — Stand bei Übergabe an den Dispatcher (2026-10-07): run-lane.sh erkennt jetzt einen Worker, der im Lane mit gefülltem question-Feld stehen bleibt (Exit 5); spawn.sh nimmt --parent <session> und hängt es an '/jaira-dispatcher <id>'. Geprüft: 'jaira set <id> question=...' geht in jeder Lane, 'question=' leert es. Noch offen: die Texte in jaira-role-lane, jaira-dispatcher, jaira-teamlead (DoD 1-3, 7, 8), NOTES.md. Alex will es einfach halten und flexibel lassen: keine neuen Mechanismen über diese zwei hinaus; wer mit dem Menschen spricht, antwortet ihm dort.
