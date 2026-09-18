---
id: 01M2SNEP2NN02DN5W0ENVSC1GW
title: "Die Pruefschleife gehoert ins Binary, nicht in den Katalog"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer jaira installiert, um Agenten arbeiten zu lassen, bekommt die Pruefschleife sofort - statt die Haelfte des Produkts im Katalog vermuten zu muessen"
context: |-
  Alex am 2026-09-18, nachdem 'jaira lanes' auf diesem Board dreizehn Lanes zeigte und zehn davon built-in waren.

  Was ist: 'critique', 'optimize' und 'testing' sind KEINE built-ins. Built-in sind genau zehn: backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked (core/lane/lane.go:34, go:embed builtin/*.md). Die drei liegen im Katalog unter lanes/ und erreichen ein Board nur ueber 'jaira lanes market adopt <id>' plus 'jaira lanes add <id>'. lanes/README.md nennt das ausdruecklich Absicht: built-ins stehen auf jedem Board, Katalog-Lanes nur, wenn jemand sie absichtlich holt.

  Warum das heute stoert: ohne diese drei ist ein frisches Board ein Tracker, kein Agenten-Konveyer. Die Schleife critique -> in-progress ist das, was auf diesem Board die echten Defekte findet - an einem einzigen Tag unter anderem einen Unterscheider, der einen schreibenden Worker still zum stummen Leser gemacht haette, und einen fest verdrahteten Branch-Namen 'master' in einem ausgelieferten Prompt. Wer jaira wegen der Agenten installiert und die drei nicht kennt, bekommt nichts davon und erfaehrt auch nicht, dass es sie gibt: 'jaira lanes' zeigt nur, was installiert ist, und der Katalog meldet sich von selbst nirgends.

  Das Gegenargument steht in CLAUDE.md und ist ernst zu nehmen: 'Every feature is measured against is this smaller than paca - the project fails by growing'. Dreizehn Lanes als Voreinstellung zwingt jedem einen Konveyer auf, auch dem, der jaira als reinen Tracker aufsetzt.

  Die Entscheidung ist deshalb NICHT 'alles einbauen', sondern: was ist die Voreinstellung fuer wen. Moegliche Formen, in der Reihenfolge wachsenden Eingriffs - (1) nur sichtbar machen: ein frisches Board oder 'jaira lanes' nennt die Katalog-Lanes, die es nicht hat; (2) die drei ins Binary einbetten, aber nicht auf die Voreinstellung setzen, also ohne Netz adoptierbar; (3) ein zweites Default-Board 'agentic' neben dem schlanken, bei 'jaira init' waehlbar; (4) sie als Voreinstellung fuer jedes neue Board. Welche - das ist die eigentliche Frage dieses Tickets und gehoert in die brainstorm-Lane.
definition-of-done: "Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-18T07:07:21Z
updated-at: 2026-09-18T08:10:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-28090
claimed-at: 2026-09-18T08:09:46Z
mode: conversational
---

# Die Pruefschleife gehoert ins Binary, nicht in den Katalog

## Definition of Done

- [ ] Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl.
- [ ] Ein frisch mit 'jaira init' angelegtes Board erfaehrt von der Pruefschleife, ohne dass jemand den Katalog kennt. Nachgestellt an einer leeren Testdoska: der Weg von 'jaira init' zu einer arbeitenden critique-Lane ist ohne Vorwissen gehbar.
- [ ] Wer jaira ohne Netz benutzt, kommt an die drei Lanes heran, oder die Fehlermeldung sagt, dass sie aus dem Netz kommen und wie man sie sonst bekommt. 'jaira lanes market' holt heute von GitHub.
- [ ] Ein Board, das jaira als reinen Tracker benutzt, wird nicht mit einem Konveyer beladen, den niemand faehrt - die Wahl bleibt eine Wahl.
- [ ] Je eine Zeile in core/release/NOTES.md fuer das, was ein Benutzer dadurch anders tut.
- [ ] Der Katalog ist an eine Version gebunden oder sagt, dass er es nicht ist: 'jaira lanes market' holt heute von 'https://api.github.com/repos/BeMuCa/jaira/contents/lanes' ohne '?ref=' und bekommt damit den HEAD des Default-Branches, egal wie alt das laufende Binary ist. Entweder fragt der Aufruf den Tag des laufenden Binaries ab, oder er sagt dem Benutzer, dass er die Entwicklungsfassung bekommt. Nachgestellt mit einem Binary, das eine aeltere Version meldet.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-18 07:10 · Alexander Sacharov** — Alex am 2026-09-18, beim Durchdenken der Katalog-Idee: 'kann man dann nicht alle Lanes in den Markt legen, und werden sie aus dem Release geladen?'

Zweite Frage nachgesehen, und die Antwort ist nein: core/market/market.go:45 baut die Adresse als 'https://api.github.com/repos/BeMuCa/jaira/contents/lanes' - ohne '?ref='. Die GitHub-Contents-API liefert ohne ref den HEAD des Default-Branches. Der Katalog ist also gar nicht versioniert: wer 0.1.4 installiert hat, bekommt heute die Lanes von heutigem master. Nennt eine Lane spaeter ein Feld oder eine Option, die sein Binary nicht kennt, kommt sie kaputt bei ihm an - und umgekehrt laesst sich der Katalog fuer ein altes Release nicht reparieren, es gibt nur einen fuer alle.

Erste Frage, warum nicht ALLES in den Katalog kann: die zehn built-ins sind der Offline-Boden. 'jaira init' funktioniert heute im Flugzeug. Mit allem im Katalog haengt das Anlegen eines Boards an GitHub - nicht erreichbar, Rate-Limit erschoepft, Repository umgezogen, und es entsteht kein Board. Dazu waere der Weg zu einem arbeitenden Board eine Kette aus zehnmal 'market adopt' plus 'lanes add' statt einem 'jaira init'.

Was das fuer dieses Ticket bedeutet: die Grenze ist nicht 'wichtig oder unwichtig', sondern 'ohne das gibt es kein Board' gegen 'das aendert, wie man arbeitet'. Und Form 2 aus dem Kontext - einbetten, aber nicht voreinstellen - ist genau deshalb ein ernster Kandidat: sie loest Offline und Versionierung in einem, weil eine eingebettete Lane mit dem Binary mitreist und damit automatisch zu ihm passt.
- **2026-09-18 07:12 · Alexander Sacharov** — Alex fragt nach, wie die Versionierung technisch aussieht: 'aendern wir dann market UND das Erstellen des Releases, damit alle Dateien als Download im Release liegen?' Drei Wege durchgesehen, damit brainstorm nicht bei null anfaengt.

A - EIN PARAMETER. Die GitHub-Contents-API nimmt '?ref='. Aus core/market/market.go:45 wird 'https://api.github.com/repos/BeMuCa/jaira/contents/lanes?ref=v<version des laufenden Binaries>'. Am Release aendert sich NICHTS - GitHub liefert den Baum jedes Tags von selbst. Ein Source-Build hat keinen Tag: dann auf den Default-Branch zurueckfallen und es laut sagen, nicht stillschweigend.

B - DATEIEN ALS RELEASE-ASSETS. Das war die Frage. Der meiste Aufwand: der Release-Workflow muss sie hochladen, und market muss statt der Verzeichnisliste die Asset-Liste lesen. Heraus kommt genau das, was A schon liefert. Der einzige Mehrwert waere Unabhaengigkeit von der Repository-Struktur - dafuer lohnt kein eigener Auslieferungsweg.

C - EINBETTEN. Fuer critique, optimize und testing der eigentliche Punkt: sie reisen mit dem Binary, die Version stimmt per Konstruktion, es braucht gar kein Netz. market bleibt dann das, was es sein sollte - der Ort fuer FREMDE Lanes, nicht fuer die eigenen.

Empfehlung des Dispatchers, nicht Entscheidung: C fuer die drei, A fuer den Rest des Katalogs, B gar nicht. Die Entscheidung gehoert der brainstorm-Lane und ist DoD 1.
- **2026-09-18 08:09 · Alexander Sacharov** — Dispatcher, beim Start am 2026-09-18: Alex ist weg und hat gesagt, dieses Ticket soll anlaufen, sobald 7XXT00 eine Menschen-Lane erreicht - das ist jetzt passiert. Er hat die Entscheidung dieses Tickets NICHT vorweggenommen, also ist sie offen und dieses Ticket besteht praktisch nur aus ihr. Nicht raten: die vier moeglichen Formen stehen im Kontext, die drei technischen Wege zur Versionierung (A ?ref=, B Release-Assets, C einbetten) in der Notiz darueber, samt der Empfehlung des vorigen Dispatchers (C fuer die drei Lanes, A fuer den Rest, B gar nicht) - eine Empfehlung, keine Entscheidung. Arbeite brainstorm aus, lege die Frage sauber vor und halte dann an. Der Zweig feat/VSC1GW haengt an master (3d12b7a, 0.3.0 gemerged).
