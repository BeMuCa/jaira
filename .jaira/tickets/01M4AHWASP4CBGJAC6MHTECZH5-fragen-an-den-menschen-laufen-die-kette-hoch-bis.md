---
id: 01M4AHWASP4CBGJAC6MHTECZH5
title: Fragen an den Menschen laufen die Kette hoch bis zum Teamlead
status: critique
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
updated-at: 2026-10-07T06:58:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65032
claimed-at: 2026-10-07T06:52:58Z
outcome-what: "run-lane.sh zählt nur eine neue Frage (q0-Basis) als Exit 5; jaira-role-lane beschreibt, was ein Worker mit einer Rollenfrage tatsächlich tut"
outcome-why: "jaira move leert question nie, ein Ticket aus human hätte sofort falsch Exit 5 gemeldet; der Worker hat kein --parent und damit keinen direkten Weg zum Dispatcher"
outcome-resolves: "critique Runde 1, beide Punkte"
review-summary: |-
  core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:205 nimmt jedes nichtleere question-Feld als Frage dieses Workers, aber jaira move leert das Feld nie (erledigte Tickets im logbook tragen ihre human-Frage noch) — ein Ticket, das aus human zurückkommt, meldet beim ersten Idle eines Workers falsch Exit 5 mit einer schon beantworteten Frage; vor der Schleife q0="$(question)" merken und nur eine Frage zählen, die != q0 ist
  core/role/builtin/jaira-role-lane/SKILL.md, Absatz "A question about your role ... goes to the dispatcher directly": der Worker bekommt laut Note 06:55 bewusst kein --parent und hat damit keinen direkten Weg zum Dispatcher; schreiben, was er tatsächlich tut — die Frage als letzte Zeile in seinen Pane, ohne question-Feld, und anhalten (der Dispatcher sieht sie, wenn run-lane.sh mit Exit 3 zurückkommt)
---

# Fragen an den Menschen laufen die Kette hoch bis zum Teamlead

## Definition of Done

- [x] jaira-role-lane: ein Worker fragt nie selbst; er schreibt die Frage als nummerierte Optionen aufs Ticket, gibt sie in seinem Bericht an den Dispatcher und hält an
  proof: core/role/builtin/jaira-role-lane/SKILL.md: section 'You never ask the person yourself'
- [x] jaira-dispatcher: mit Starter reicht er jede Frage (eigene und die seiner Worker) an den Starter weiter und wartet; ohne Starter fragt er selbst wie heute
  proof: core/role/builtin/jaira-dispatcher/SKILL.md: 'How you ask' — --parent case, workers' questions (run-lane exit 5)
- [x] jaira-teamlead: nimmt weitergereichte Fragen an, stellt sie als Wahl, schreibt die Antwort mit jaira note aufs Ticket und gibt sie an den Dispatcher zurück
  proof: core/role/builtin/jaira-teamlead/SKILL.md: 'What you actually decide' point 4 — SendMessage up, note, SendMessage back down
- [x] spawn.sh gibt dem gestarteten Dispatcher bzw. Worker den Namen seines Starters mit
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh: --parent -> '/jaira-dispatcher <id> --parent <s>'; teamlead SKILL.md passes --parent <own name>
- [ ] Ein Probelauf Teamlead -> Dispatcher -> Worker in Herdr: eine Frage des Workers erscheint nur im Teamlead-Tab, die Antwort kommt beim Worker an
- [x] core/release/NOTES.md hat eine Zeile unter Unreleased
  proof: core/release/NOTES.md: line under ## Unreleased
- [x] Die Skills sagen die Richtung: der Teamlead ist der Ort, an dem der Mensch antwortet; ein direkt gestarteter Dispatcher ist der Ausnahmefall
  proof: core/role/builtin/jaira-teamlead/SKILL.md 'You are where the person answers'; core/role/builtin/jaira-dispatcher/SKILL.md 'Where the question goes depends on who started you'
- [x] Schreibt der Mensch direkt in den Dispatcher-Tab, antwortet der Dispatcher dort, hält das Entschiedene mit jaira note fest und meldet es seinem Teamlead; der Teamlead fragt es nicht noch einmal
  proof: core/role/builtin/jaira-dispatcher/SKILL.md 'The person may still type into your tab'; core/role/builtin/jaira-teamlead/SKILL.md same paragraph, no re-asking

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
- **2026-10-07 06:55 · Alexander Sacharov** — in-progress: Worker bekommt bewusst KEIN --parent. Er fragt nie und braucht keinen Adressaten: sein Kanal ist das question-Feld + jaira note, das run-lane.sh (Exit 5) liest; DoD 4 ist damit für den Worker durch den Dispatcher erfüllt. Dispatcher schreibt den Parent-Namen selbst als note aufs Ticket, damit ein nach jaira resume neu gestarteter Dispatcher ihn kennt — kein neues Frontmatter-Feld (Alex: keine neuen Mechanismen). DoD 5 (Probelauf in Herdr) offen: braucht installierte Rollen aus diesem Branch (jaira roles install --force) und einen Menschen, der sieht, in welchem Tab die Frage erscheint — Sache der testing-Lane. go test ./core/role/... grün.
- **2026-10-07 06:56 · Alexander Sacharov** — critique (Runde 1): core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:205 nimmt jedes nichtleere question-Feld als Frage dieses Workers, aber jaira move leert das Feld nie (erledigte Tickets im logbook tragen ihre human-Frage noch) — ein Ticket, das aus human zurückkommt, meldet beim ersten Idle eines Workers falsch Exit 5 mit einer schon beantworteten Frage; vor der Schleife q0="$(question)" merken und nur eine Frage zählen, die != q0 ist
core/role/builtin/jaira-role-lane/SKILL.md, Absatz "A question about your role ... goes to the dispatcher directly": der Worker bekommt laut Note 06:55 bewusst kein --parent und hat damit keinen direkten Weg zum Dispatcher; schreiben, was er tatsächlich tut — die Frage als letzte Zeile in seinen Pane, ohne question-Feld, und anhalten (der Dispatcher sieht sie, wenn run-lane.sh mit Exit 3 zurückkommt). Stehen gelassen: Parent-Name als note statt Frontmatter-Feld (Note 06:55, Alex: keine neuen Mechanismen); Worker ohne --parent (Note 06:55); question-Feld + note als Doppelkanal — das Feld ist das maschinenlesbare Signal für run-lane.sh, die note das Überlebende.
- **2026-10-07 06:58 · Alexander Sacharov** — in-progress Runde 2: beide critique-Punkte umgesetzt. run-lane.sh merkt q0=$(question) nach dem Startstatus und zählt nur eine Frage != q0 als Exit 5 — ein aus human zurückgekehrtes Ticket mit alter Frage wartet also normal weiter. Bewusst in Kauf genommen: stellt ein Worker wörtlich dieselbe Frage noch einmal, sieht run-lane.sh sie nicht und läuft bis Timeout (Exit 3, Pane bleibt offen) — harmlos. jaira-role-lane: Rollenfrage (welche critique) als letzte Zeile im Pane, ohne question-Feld/note; Dispatcher liest den Pane bei Exit 3 (passt zu dispatcher SKILL.md 'Exit 3 … look at its tab'). DoD 5 (Probelauf Herdr) bleibt offen für testing/Mensch.
