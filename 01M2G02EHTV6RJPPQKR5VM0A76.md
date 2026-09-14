---
id: 01M2G02EHTV6RJPPQKR5VM0A76
title: "Ein Board im Datei-Modus sagt es nicht, kommt nicht zurueck und laesst sich nicht pruefen"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Wer ein Ticket anlegt, sieht in derselben Zeile, ob es auf einem Ref liegt oder als Datei; ein Datei-Ticket kommt mit einem Befehl auf seinen Ref; und ein Befehl sagt, in welchem Modus dieses Board laeuft und warum."
context: |-
  Alex hat am 2026-09-14 eine Stunde verloren, weil siebzehn Tickets auf dem requirementsgenie-Board still als Datei angelegt wurden statt auf ihren Refs. Kein Befehl hat das gesagt, und zurueck ging es auch nicht.

  Die Ursache selbst - der Remote-Name gilt pro Rechner statt pro Board - ist 9ET6NC und wird dort behoben. Dieses Ticket ist alles andere: dass derselbe Fehler wieder eine Stunde kosten wuerde, auch wenn er aus einem anderen Grund auftritt.

  Fuenf Stellen, an denen das Werkzeug es haette sagen koennen:

  1. 'jaira create' nennt den gewaehlten Modus nicht. internal/cli/refs.go:69 fileOnRefOnly() gibt bei fehlendem Remote einfach false zurueck, und create schweigt darueber. Im Ref-Modus steht 'On its ref, not on your disk' da; im Datei-Modus steht nichts.

  2. Der Diagnosetext existiert, steht aber an der falschen Stelle. Nur 'jaira release' nennt die wirkliche Ursache: 'this board does not carry tickets on refs: gitref: no repository or no such remote: no remote "upstream"'. release laeuft am Ende, der Fehler entsteht bei create am Anfang.

  3. Es gibt keinen Weg zurueck. Ein Datei-Ticket laesst sich nicht auf seinen Ref legen. 'jaira release' darauf antwortet 'gitref: no ref for this ticket' - eine Feststellung, kein Hinweis. Der einzige Weg war, siebzehn Tickets mit neuen Ids neu anzulegen.

  Das Stueck dafuer ist schon da und ist universell: Record() (core/refsync/refsync.go:108) stellt jede Schreibung in die Outbox, und ein Ticket ohne Ref least den leeren String, also 'ich erwarte, dass dieser Ref nicht existiert'. fileOnRefOnly() (internal/cli/refs.go:69) schickt und entfernt danach die Datei. An create gebunden ist davon nur der Funktionsname.

  4. Kein Befehl zeigt den Zustand des Boards. 'jaira whoami' zeigt die Identitaet, aber nichts von der git-Seite: welcher Remote erwartet wird, ob er da ist, in welchem Modus das Board laeuft, wie viele Tickets nur als Datei liegen.

  5. 'jaira create --dod' nimmt genau einen Eintrag. Mehrzeilig uebergeben wird alles ab dem zweiten Absatz zu loser Prosa im Rumpf, ohne Kaestchen. Weitere Punkte gehen nur ueber 'jaira dod --add', und das verlangt die Datei - also 'jaira pull', und das schreibt den Aufrufer als assignee ein. Ein Ticket mit sieben Kriterien laesst sich im Ref-Modus also nicht anlegen, ohne es sich selbst zuzuweisen.

  Nachgestellt am 2026-09-14: APABM4, S1VM40 und 9ET6NC wurden mit fuenf bis sechs Kriterien angelegt und trugen danach je genau eines. Bei 9ET6NC lief der Gate der Endlane dadurch gegen ein Fuenftel der Bedingungen, waehrend ein Worker den Rest als Kontext las.
definition-of-done: "'jaira create' nennt den Modus in beiden Faellen. Im Datei-Modus sagt es, dass das Ticket als Datei und nicht auf einem Ref liegt, nennt den Remote-Namen, nach dem gesucht wurde, und den Grund - nicht nur das Schweigen von heute."
tags:
  - cli
blocked-by: []
related:
  - 01M2FYEHZYFCW7DSNRHY9ET6NC
commits: []
created-at: 2026-09-14T13:00:30Z
updated-at: 2026-09-14T13:23:44Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-152036
claimed-at: 2026-09-14T13:23:44Z
---

# Ein Board im Datei-Modus sagt es nicht, kommt nicht zurueck und laesst sich nicht pruefen

## Definition of Done

- [ ] 'jaira create' nennt den Modus in beiden Faellen. Im Datei-Modus sagt es, dass das Ticket als Datei und nicht auf einem Ref liegt, nennt den Remote-Namen, nach dem gesucht wurde, und den Grund - nicht nur das Schweigen von heute.
- [ ] Ein Datei-Ticket kommt mit einem Befehl auf seinen Ref und die lokale Datei verschwindet dabei. Die Logik dafuer ist die vorhandene fileOnRefOnly (internal/cli/refs.go:69), nicht eine zweite Kopie davon. Nachgestellt auf einem Fixture-Board, dessen Ticket im Datei-Modus entstanden ist.
- [ ] Ein Befehl zeigt den git-Zustand des Boards in einem Aufruf: den eingestellten Remote-Namen, die Remotes die dieses Repository hat, ob der Ref-Modus laeuft, und wie viele Tickets nur als Datei liegen. Nachgestellt auf einem Board mit passendem und auf einem mit fehlendem Remote.
- [ ] 'jaira create --dod' nimmt den Schalter mehrfach, wie --tag es tut. Nachgestellt: ein Ticket mit drei Kriterien wird im Ref-Modus mit einem Aufruf angelegt und traegt danach drei Kaestchen, ohne dass es dafuer gepullt wurde.
- [ ] Der Diagnosetext, der heute nur aus 'jaira release' kommt, erscheint dort wo der Zustand entsteht. Nachgestellt: auf einem Board ohne passenden Remote nennt schon der erste 'jaira create' den Grund, nicht erst ein Befehl am Ende der Kette.
- [ ] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer jede von aussen sichtbare Aenderung: die Modus-Zeile in create, der neue Befehl fuer den Ref-Nachtrag, der neue Zustandsbefehl, und der wiederholbare --dod.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

