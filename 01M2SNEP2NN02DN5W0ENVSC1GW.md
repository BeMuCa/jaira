---
id: 01M2SNEP2NN02DN5W0ENVSC1GW
title: "Die Pruefschleife gehoert ins Binary, nicht in den Katalog"
status: brainstorm
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
updated-at: 2026-09-18T10:30:32Z
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

- [x] brainstorm
- [x] planning

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
- **2026-09-18 08:12 · Alexander Sacharov** — brainstorm, am Code nachgesehen - was WIRKLICH im Weg steht, und warum Form 2 nicht das ist, was sie im Kontext zu sein scheint.

WAS DER CODE ZEIGT

1. 'eingebettet' und 'auf jedem Board voreingestellt' sind heute EIN Bit, nicht zwei. Builtins() (core/lane/lane.go:385) liest go:embed builtin/*.md, und setUp() (core/lane/lane.go:614) schreibt beim ersten Oeffnen eines Boards genau die Lanes als Dateien, fuer die l.Builtin true ist - sofern kein Default-Board existiert. Wer critique.md einfach nach core/lane/builtin/ legt, hat damit NICHT Form 2 gebaut, sondern Form 4: dreizehn Lanes auf jedem frischen Board. Form 2 verlangt, das Feld Lane.Builtin (lane.go:136) in zwei Tatsachen zu trennen - 'reist im Binary mit' und 'steht in der Vorauswahl'. Das ist ein Frontmatter-Feld mehr (etwa 'default-board: false') oder ein zweites embed-Verzeichnis, plus die eine Zeile in setUp, die statt l.Builtin die Vorauswahl fragt.

2. Ist das getan, ist Offline geschenkt. Installable() (core/lane/order.go:175) bietet schon heute built-ins PLUS ~/.jaira/lanes an, und lane.Add() kopiert von dort auf das Board. Eine eingebettete critique-Lane ist damit ohne jedes Netz per 'jaira lanes add critique' installierbar - 'market adopt' entfaellt fuer sie ganz. Das ist DoD 3 erledigt, ohne eine Fehlermeldung schreiben zu muessen.

3. Form 1 allein erreicht Offline NICHT. Die Namen der Katalog-Lanes kennt das Binary heute nirgends; die Liste kommt ausschliesslich aus market.List() ueber die Contents-API (core/market/market.go:45). 'jaira lanes' koennte also ohne Netz gar nicht nennen, was fehlt - es sei denn, man verdrahtet die drei Namen fest, was ein halbes Einbetten ohne dessen Nutzen ist.

4. Form 3 hat die Mechanik schon, braucht aber trotzdem neuen Stoff. DefaultBoard (core/lane/defaultboard.go) ist per Benutzer und kann jede Lane-Auswahl halten - ein zweites 'agentic'-Profil waere kein neues Konzept im Kern. Neu waeren die Frage in 'jaira init' und ein zweiter mitgelieferter Auswahl-Satz; und sie hilft nur dem, der beim Anlegen schon weiss, was er waehlen soll.

5. DoD 6, Versionierung: version steht in cmd/jaira/main.go:11 (ldflags, 'dev' im Source-Build) und ist dem Paket market heute gar nicht bekannt. Weg A - '?ref=v<version>' an apiBase() - sind wenige Zeilen plus Durchreichen der Version, Fallback auf den Default-Branch mit lauter Ansage, wenn die Version 'dev' ist. Weg B (Release-Assets) liefert dasselbe und kostet den Release-Workflow dazu. Unter Form 2 braucht der Markt fuer die drei Lanes gar nicht mehr versioniert zu sein - er bleibt noetig fuer den Rest des Katalogs (changelog-writer, secrets-scan).

EMPFEHLUNG: Form 2, und Form 1 obendrauf, weil sie durch 2 erst moeglich und dann fast gratis ist. Dazu Weg A.

Form 2 einbetten, aber nicht voreinstellen: kein Board bekommt eine Lane mehr als heute (DoD 4 bleibt gewahrt), kein neues Kommando, keine Frage in 'jaira init' - es wird ein Bool in zwei Bools zerlegt und eine Netzabhaengigkeit entfernt statt eine hinzugefuegt. An der Scope-Regel gemessen ist das die einzige der vier Formen, die das Binary nicht um ein Konzept wachsen laesst.
Form 1 obendrauf: 'jaira lanes' nennt am Fuss die mitgelieferten, nicht installierten Lanes mit der Zeile, die sie holt. Ohne das ist DoD 2 nicht erfuellt - eingebettet und ungenannt ist so unsichtbar wie im Katalog. Mit 2 im Ruecken ist es eine Liste aus dem Binary statt aus dem Netz.
Form 4 faellt an der Scope-Regel: dreizehn Lanes als Voreinstellung zwingt jedem Tracker-Benutzer einen Konveyer auf, den niemand faehrt.
Form 3 faellt am Aufwand gegen den Nutzen: init-Dialog plus zweites Profil, und der Nutzen trifft nur den, der die Wahl schon versteht. Sie bleibt spaeter nachruestbar - ein 'agentic'-Default-Board ist nach Form 2 nur noch eine Datei mit dreizehn Namen darin.
Weg B faellt: eigener Auslieferungsweg fuer das Ergebnis, das A mit einem Query-Parameter liefert.

OFFEN UND MENSCHENSACHE: die Wahl selbst (DoD 1). Alex hat sie ausdruecklich nicht vorweggenommen.
- **2026-09-18 10:30 · Alexander Sacharov** — ENTSCHIEDEN von Alex am 2026-09-18, in der Sitzung mitlesend (conversational): FORM 2, mit Form 1 obendrauf und Weg A fuer den Rest des Katalogs.

Gewaehlt: die drei Lanes reisen im Binary mit, stehen aber NICHT in der Vorauswahl. Dazu nennt 'jaira lanes' die mitgelieferten, nicht installierten Lanes, damit ein frisches Board von ihnen erfaehrt. 'jaira lanes market' bekommt '?ref=v<version des laufenden Binaries>' fuer den Katalog, der weiterhin aus dem Netz kommt.

Verworfen, an der Scope-Regel aus CLAUDE.md gemessen:
- Form 1 allein: erreicht DoD 3 nicht. Das Binary kennt die Namen der Katalog-Lanes nirgends, sie kommen nur aus market.List() (core/market/market.go:45). Ohne Netz kann es nicht nennen, was fehlt - ausser man verdrahtet die drei Namen fest, was ein halbes Einbetten ohne dessen Nutzen ist.
- Form 3: die Mechanik traegt zwar (DefaultBoard, core/lane/defaultboard.go), neu waeren aber ein Dialog in 'jaira init' und ein zweiter mitgelieferter Auswahl-Satz - ein Konzept mehr im Binary, und der Nutzen trifft nur den, der die Wahl beim Anlegen schon versteht. Nach Form 2 bleibt sie jederzeit nachruestbar: ein 'agentic'-Default-Board ist dann nur noch eine Datei mit dreizehn Namen darin.
- Form 4: dreizehn Lanes als Voreinstellung zwingt jedem, der jaira als reinen Tracker aufsetzt, einen Konveyer auf, den niemand faehrt. Das ist genau das Wachstum, gegen das die Scope-Regel steht.
- Weg B (Release-Assets): ein eigener Auslieferungsweg fuer das Ergebnis, das Weg A mit einem Query-Parameter liefert.

Warum Form 2 als einzige die Scope-Regel besteht: kein Board bekommt eine Lane mehr als heute, kein neues Kommando, keine Frage in 'jaira init'. Es wird ein bestehendes Bool (Lane.Builtin, core/lane/lane.go:136) in zwei Tatsachen zerlegt - 'reist im Binary mit' und 'steht in der Vorauswahl' - und eine Netzabhaengigkeit entfernt statt eine hinzugefuegt.
