---
id: 01M2E86MS66ZZ58RRV3EY0A7DT
title: "Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira bringt ein lauffaehiges Beispiel-Hookskript im Binary mit und druckt es auf Zuruf, damit der erste Lane-Wechsel ohne eigene Erfindung bei einem Menschen ankommt - und nur dann, wenn sich ohne ihn nichts bewegt."
context: |-
  core/hook/hook.go ruft bei jedem move und jedem claim ein Skript auf und uebergibt JAIRA_EVENT, JAIRA_TICKET, JAIRA_TITLE, JAIRA_STATUS, JAIRA_ASSIGNEE, JAIRA_ACTOR und JAIRA_ROOT als Umgebungsvariablen. Der Mechanismus ist fertig und gut: jaira bringt bewusst keine eigene Abhaengigkeit mit, die Anbindung gehoert dem Nutzer.
  Was fehlt, ist der erste Schritt. Wer 'hook' in den Einstellungen sieht, hat ein leeres Feld und keinen Anhaltspunkt, was hineingehoert. Praktisch schreibt es deshalb niemand.
  Am 2026-09-13 ist hier eines von Hand entstanden, ~/.jaira/notify.sh, rund 25 Zeilen: es prueft HERDR_ENV, verwandelt den Lane-Wechsel in 'herdr notification show', und waehlt den Ton nach Ziel-Lane - ein Lane, der einem Menschen gehoert, bekommt request, done bekommt done, alles andere bleibt stumm. Die Regel dahinter, die das Beispiel transportieren soll: einen Ton verdient nur der Zustand, in dem sich ohne den Menschen nichts bewegt.
  Zu klaeren, absichtlich nicht entschieden:
  - ob das Beispiel eingebettet wird (wie die Lanes, go:embed) und per Befehl geschrieben, oder ob es nur im Repository unter einem Beispielordner liegt und in der Dokumentation genannt wird.
  - ob mehr als ein Beispiel beiliegt. Herdr ist das, was hier laeuft; ntfy.sh und eine Terminal-Glocke sind die naechstliegenden, weil sie keine Konten brauchen. Slack braucht eine Webhook-URL und ist damit kein Beispiel, sondern eine Einrichtung.
  - ob der Einsetz-Befehl auch settings.json schreibt oder nur die Datei hinlegt und sagt, welche Zeile fehlt. Fremde Einstellungen ungefragt zu aendern ist die unangenehmere Variante.
  Nicht Teil dieses Tickets: an core/hook selbst etwas aendern. Der Vertrag steht.
definition-of-done: "Ein Befehl legt ein lauffaehiges Hook-Beispiel ab und sagt, wie es scharfgeschaltet wird; das Beispiel laeuft ohne Herdr fehlerfrei durch und tut dann nichts; ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar; ein Hinweis auf das Beispiel steht dort, wo 'hook' dokumentiert ist; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
parent: 01M2E248SM9X1JRZBNTHC9V7ZV
related: []
commits:
  - e03f8f08e71752da054ea672fefa3d96a0d8a4ec
  - a5a924365c667141d1a2f765d639422de88625c3
  - 1115a7ac45e894dbeb465f99a62cb9a8f4a15161
  - 7b8a64bf66b10304d9c90a411d035ced2dacdd7a
  - a19a5f3929b5faa69232959127cbd5a0495a9c47
  - fd86a96d7f945cf8eee08c1a1dc2f45d416fd23c
  - 503c8995721b8575898f3310b27d5e4c2aa1da84
  - 11fa2a4332a4a39eed01e914c35797583c1ab909
  - cb5ed6af20adf1a8a63f02b389da2990834c4814
  - 1c6be9d38c71f7574c25899ac16cd599f5577117
  - 0c5da60a7dfe4e552d2ce8c64712d8e4b0d257ac
created-at: 2026-09-13T20:44:07Z
updated-at: 2026-09-14T15:15:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-60947
claimed-at: 2026-09-13T21:46:35Z
outcome-what: "Zweiter Review-Durchgang: die drei rueckweisenden Maengel des ersten Durchgangs sind einzeln nachgeprueft und behoben - Kopfkommentar und case-Zweig des Beispielskripts sagen dasselbe, NOTES.md:17 beschreibt das Verhalten des Binaries, und die Behauptung ueber zwei hoerbare Toene ist zurueckgenommen statt vorgetaeuscht. review-summary, review-gaps, review-verdict und review-check sind gesetzt."
outcome-why: "Der Diff seit 0945a8a besteht ausschliesslich aus diesen drei Korrekturen, enthaelt keine Codezeile und 'go test ./... -race' ist gruen. Das Ticket erfuellt damit seinen Zweck: das Beispiel widerspricht der Regel nicht mehr, die es weitergeben soll. Ein nicht rueckweisender Rest (README.md:353-357 traegt dieselbe Aussage in der ueberholten Fassung) ist in review-gaps und in einer Notiz festgehalten, damit ein Mensch entscheidet, ob er noch hier nachgezogen wird."
outcome-resolves: "Alle DoD-Zeilen abgenommen: Befehl mit Anleitung, Lauf ohne Werkzeug (env -i, leerer PATH, exit 0), hoerbare Unterscheidung Menschen-Lane gegen gewoehnlichen Wechsel (Ton gegen Stille, feinere Unterscheidung ehrlich ueber die gedruckte Zeile), README-Hinweis beim 'hook'-Absatz, genau eine Zeile unter ## Unreleased (nachgezaehlt), go test ./... -race gruen."
review-summary: "Zweiter Durchgang. Der Diff seit 0945a8a aendert genau zwei Dateien und keine Zeile Code: core/hook/example/notify.sh bekommt zwei umgeschriebene Kommentarbloecke, core/release/NOTES.md eine ersetzte Zeile. Der Kopfkommentar des Beispielskripts (Zeilen 18-24) nennt jetzt zwei bewusste Ausnahmen von seiner eigenen Regel - die Lane, die auf einen Menschen wartet (human, signoff), und die Endlane (done) - und deckt sich damit mit dem case-Zweig in Zeile 55-59, der genau diese drei Lanes durchlaesst. Der Kommentar ueber deliver() (35-42) behauptet nicht mehr, zwei Glocken gegen eine seien der hoerbare Unterschied, sondern sagt, dass die gedruckte Zeile ('jaira human:' gegen 'jaira done:') ihn traegt und die Glockenzahl nur ein Hinweis ist, weil Terminals dicht aufeinander folgende BEL verschmelzen. NOTES.md:17 beschreibt jetzt dasselbe Verhalten wie das Binary: stumm fuer jede Agenten-Lane, Ton bei human, signoff und done, unterschieden an der gedruckten Zeile. Befehl, go:embed, /dev/tty-Probe und die drei Tests sind unangetastet."
review-gaps: |-
  Die drei rueckweisenden Maengel des ersten Durchgangs sind erledigt, jeder einzeln nachgeprueft:

  1. Widerspruch Kopfkommentar gegen Code - behoben. notify.sh:18-24 nennt done ausdruecklich als zweite, bewusste Ausnahme und begruendet sie. Der case in 55-59 laesst human, signoff, done durch, alles andere exit 0. Kopf und Code sagen jetzt dasselbe. Gegengeprueft am Skript, das das Binary wirklich druckt ('jaira hook example'), nicht nur an der Quelldatei - go:embed zieht die geaenderte Fassung.

  2. NOTES.md:17 - behoben. Die Zeile sagt jetzt 'stays silent for every agent lane and rings ... only when a ticket reaches a person (human, signoff) or finishes (done), telling the two apart by the line it prints'. Das ist das Verhalten des Skripts. Weiterhin genau eine Zeile fuer dieses Ticket unter ## Unreleased (nachgezaehlt), der Rest der Datei unberuehrt.

  3. Unbelegte Behauptung ueber zwei hoerbare Toene - behoben, und zwar durch Zuruecknehmen statt durch einen Umbau, der hier nicht zu haben ist. notify.sh:35-42 sagt jetzt, dass die Textzeile den Unterschied traegt und die Glockenzahl ein Hinweis ist; outcome-resolves des Implementierers sagt dasselbe. Die Begruendung, warum eine Pause zwischen den Glocken nicht in Frage kommt (sleep = Prozess plus Wartezeit in einem Skript, das jaira nach 5 s abschiesst), steht in der Notiz und ist nachvollziehbar.

  'go test ./... -race': alle Pakete ok, kein FAIL.

  Nicht rueckweisend, aber offen und beim naechsten Anfassen mitzunehmen:

  - README.md:353-357 ist die vierte Stelle mit derselben Aussage und ist nicht mitkorrigiert worden. Dort steht weiterhin 'twice for a lane that belongs to a person, once for a finished ticket ... because only a state in which nothing moves without you is worth a sound'. Das verkauft erstens die Glockenzahl als die Unterscheidung, die Punkt 3 gerade zurueckgenommen hat, und zweitens begruendet der Nebensatz nur die human/signoff-Haelfte, waehrend der Satz davor done mit aufzaehlt - derselbe Bruch wie in Punkt 1, nur milder. Ich weise das nicht zurueck, weil der erste Durchgang diesen Absatz gesehen und fuer richtig befunden hat und dieser Durchgang ausdruecklich nur die drei benannten Stellen pruefen sollte. Es ist aber Nutzertext, und wer die drei Stellen korrigiert hat, sollte diese vierte hinterherziehen.
  - Unveraendert offen aus dem ersten Durchgang: 'jaira hook example' druckt auch unter Windows ein /bin/sh-Skript ohne Hinweis, und der Test traegt //go:build unix; kein Test ruft das Skript ueber Shebang und x-Bit auf, worauf hook.Run per exec.Command angewiesen ist. Beides war schon damals nicht rueckweisend und gehoert in ein eigenes Ticket.
review-verdict: "Freigabe zur Abnahme. Die drei Maengel, wegen derer das Ticket zurueckging, sind an genau den benannten Stellen behoben, und der Diff seit 0945a8a besteht nur aus diesen Korrekturen - keine Codezeile, kein Umbau, kein Nebeneffekt. Beachtenswert ist, wie Punkt 3 geloest wurde: nicht durch Erfinden einer hoerbaren Unterscheidung, die POSIX-Shell ohne externes Werkzeug nicht hergibt, sondern durch ehrliches Zuruecknehmen der Behauptung. Das ist die richtige Entscheidung fuer ein Beispiel, dessen Zweck das Weitergeben einer Regel ist. Ein Rest bleibt: README.md:353-357 traegt dieselbe Aussage in ihrer alten, jetzt ueberholten Form. Ich weise deswegen nicht zurueck - der Absatz war im ersten Durchgang abgenommen und lag ausserhalb dessen, was dieser Durchgang pruefen sollte - aber es ist Text, den ein Nutzer liest, und er sollte nachgezogen werden, bevor oder kurz nachdem das hier landet. Wer abnimmt, entscheidet, ob das noch in dieses Ticket gehoert oder in ein eigenes."
review-check: |-
  1. Im Verzeichnis /home/alex/projects/jaira-Y0A7DT ausfuehren: 'PATH=$PATH:/usr/local/go/bin go run ./cmd/jaira hook example > /tmp/notify.sh && chmod +x /tmp/notify.sh'. Es erscheint keine Ausgabe, die Datei entsteht.
  2. 'sed -n "18,24p" /tmp/notify.sh' lesen. Dort muss stehen, dass die Regel genau zwei bewusste Ausnahmen hat: die Lane, die auf einen Menschen wartet (human, signoff), und die Lane, die das Ticket beendet (done). Das war der erste Mangel - frueher stand dort, jede andere Lane bleibe stumm, waehrend done klingelte.
  3. 'sed -n "55,59p" /tmp/notify.sh' lesen. Dort muessen genau diese drei Lanes stehen: 'human | signoff' und 'done', alles andere 'exit 0'. Schritt 2 und Schritt 3 muessen sich decken - das ist der Kern der Pruefung.
  4. 'sed -n "35,42p" /tmp/notify.sh' lesen. Dort muss stehen, dass die gedruckte Zeile den Unterschied traegt und die Glockenzahl nur ein Hinweis ist. Das war der dritte Mangel - frueher wurde die Glockenzahl als hoerbare Unterscheidung behauptet.
  5. 'sed -n "17p" core/release/NOTES.md' lesen. Die Zeile muss human, signoff und done nennen und sagen, dass die Faelle an der gedruckten Zeile auseinanderzuhalten sind. Das war der zweite Mangel.
  6. Das Skript wirklich laufen lassen, nacheinander: 'env -i JAIRA_EVENT=move JAIRA_STATUS=human JAIRA_TITLE=Test /tmp/notify.sh', dann dasselbe mit JAIRA_STATUS=done, dann mit JAIRA_STATUS=in-progress. Erwartet: 'jaira human: Test', 'jaira done: Test', und beim dritten gar nichts.
  7. 'env -i PATH= JAIRA_EVENT=move JAIRA_STATUS=in-progress /tmp/notify.sh; echo $?' - keine Ausgabe, dann eine 0. Das Skript laeuft auf einer Maschine ohne jedes Werkzeug durch.
  8. 'PATH=$PATH:/usr/local/go/bin go test ./... -race' - alle Pakete melden ok, kein FAIL.
  9. Wenn Sie entscheiden wollen, ob der offene Punkt noch hier hineingehoert: 'sed -n "351,357p" README.md' lesen. Dort steht die alte Fassung der Aussage ('twice for a lane that belongs to a person, once for a finished ticket'), die Schritt 4 gerade zurueckgenommen hat.
---

# Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen

## Definition of Done

- [x] Ein Befehl legt ein lauffaehiges Hook-Beispiel ab und sagt, wie es scharfgeschaltet wird; das Beispiel laeuft ohne Herdr fehlerfrei durch und tut dann nichts; ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar; ein Hinweis auf das Beispiel steht dort, wo 'hook' dokumentiert ist; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen
  proof: internal/cli/hook.go:newHookExampleCmd (der Befehl 'jaira hook example' und sein Long-Text mit der Scharfschalt-Zeile); TestHookExampleRunsOnAMachineWithNothingInstalled (laeuft mit leerem PATH fehlerfrei durch); TestHookExampleSoundsOnlyForThePersonsLanes (human/signoff zwei Glocken, done eine, Agenten-Lanes stumm); README.md:344-352 im hook-Absatz; core/release/NOTES.md:17; go test ./... -race exit 0

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Beispielskript core/hook/example/notify.sh schreiben: Zustellzeile ist die Terminal-Glocke, ntfy- und notify-send-Einzeiler auskommentiert direkt darunter; Lane-Tabelle human|signoff -> Bitte-Ton, done -> Abschluss-Ton, blocked und alle Agenten-Lanes stumm; Kopfkommentar englisch mit den JAIRA_*-Variablen, der 5s-Kappung und dem Satz, welche Zeile zu ersetzen ist
- [x] Skript einbetten: neue Datei core/hook/example.go mit go:embed example/notify.sh und einer Funktion ExampleScript() string; core/hook/hook.go bleibt unberuehrt, der Vertrag aendert sich nicht
- [x] Befehl 'jaira hook example' in internal/cli/hook.go ergaenzen: druckt das Skript auf stdout, schreibt keine Datei und keine Einstellung; Long-Text nennt 'jaira hook example > ~/.jaira/notify.sh && chmod +x ~/.jaira/notify.sh' und die fehlende Zeile "hook": "/home/<du>/.jaira/notify.sh"
- [x] Test internal/cli/hook_example_test.go nach dem Muster von runStopHook: Ausgabe von 'hook example' in eine Datei schreiben und mit sh und gesetzten JAIRA_*-Variablen ausfuehren; human, signoff und done geben hoerbare Ausgabe, in-progress, todo und blocked geben nichts, jeder Fall endet mit exit 0
- [x] Test fuer den leeren Rechner: ohne jedes Zustellwerkzeug auf PATH laeuft das Skript fehlerfrei durch und tut nichts (das ist die DoD-Zeile 'laeuft ohne Herdr fehlerfrei durch')
- [x] README.md: im Hook-Absatz bei Zeile 344-350 zwei Saetze auf 'jaira hook example'; bei Zeile 561 einen Halbsatz, dass 'hook print' und 'hook example' zwei unverwandte Dinge sind
- [x] core/release/NOTES.md: eine englische Zeile unter ## Unreleased, als Anweisung formuliert
- [x] go test ./... -race gruen, danach die DoD-Zeile mit --proof abhaken
- [x] Folgeticket erfassen: 'jaira lanes --json' gibt requires-question und requires-human-exit nicht aus, deshalb stehen die Lane-Namen im Beispiel hart und es ist auf einem Board mit eigenen Lanes stumm

## Progress
- **2026-09-13 21:42 · Alexander Sacharov** — Befund im Code (Lane brainstorm; gelesen: core/hook/hook.go, core/settings/settings.go, internal/cli/hook.go, core/lane/lane.go, README.md 338-350 und 561-567, hooks/sync-tasks.sh).

Der Vertrag steht wirklich: hook.Run nimmt einen Skriptpfad, setzt sieben JAIRA_*-Variablen, kappt bei 5s, verschluckt Ausgabe und Exit-Code. Daran fehlt nichts.

Fehlend ist nicht Dokumentation, sondern ein erster Lauf. Ueber den Hook existiert heute: das Feld "hook" in core/settings/settings.go:43 und vier Zeilen README. Beides beschreibt die Schnittstelle und liefert keine einzige Zeile, die man starten koennte. Wer ein leeres Feld und eine Variablenliste sieht, muss das Skript selbst erfinden - und schreibt es deshalb nicht.

Zwei Dinge aus dem Code, die die Ticket-Notiz nicht kennt:

1. Der Name "hook" ist schon belegt. internal/cli/hook.go hat "jaira hook print", und das meint etwas ganz anderes: ein Claude-Code-Stop-Hook-Snippet fuer ~/.claude/settings.json. Ein Beispiel fuer den move/claim-Hook muss unter diesem Elternbefehl unterscheidbar heissen, sonst stehen zwei unverwandte Dinge unter einem Wort.

2. Ein Skript kann heute nicht herausfinden, welche Lane einem Menschen gehoert. "jaira lanes --json" liefert agentic, terminal, precedence - aber nicht die Gates requires-question / requires-human-exit, die human und signoff erst zu Menschen-Lanes machen (core/lane/builtin/30-human.md, 45-signoff.md). agentic:false allein reicht nicht: backlog, todo, done und blocked sind ebenfalls agentic:false. Jedes Beispiel muss die Lane-Namen also hart hinschreiben und ist auf einem Board mit eigenen Lanes falsch.
- **2026-09-13 21:43 · Alexander Sacharov** — Drei Wege, wie das Beispiel ausgeliefert werden kann.

A) Nur im Repository, unter examples/ oder hooks/, in der README genannt.
Kostet fast nichts und aendert kein Byte am Binary. Gibt auf, wen es erreicht: jaira wird als einzelne Datei verteilt, "jaira self update" holt ein Binary und kein Arbeitsverzeichnis. Wer nicht geklont hat, sieht die Datei nie - und das ist genau die Person, die das leere Feld anstarrt. Praezedenzfall hooks/sync-tasks.sh zeigt den Defekt: das Skript existiert seit Langem und muss in der README erklaert werden, weil es nirgends greifbar ist.

B) Eingebettet (go:embed) und von einem Befehl auf stdout gedruckt.
Kostet rund eine Seite Text im Binary und einen Befehl. Gibt: das Beispiel ist da, wo das Programm ist, und die Zeile "jaira <befehl> > ~/.jaira/notify.sh && chmod +x" ist der erste Lauf. Genau das Muster, das "jaira hook print" schon benutzt, und genau das Muster von core/lane/builtin/*.md - jaira bettet seine Vorlagen ein, das ist hier Hausstil und kein Sonderfall.

C) Eingebettet und der Befehl schreibt Datei plus ~/.jaira/settings.json.
Gibt: null Handgriffe. Kostet: der Befehl fasst eine Datei an, die schon einen funktionierenden Hook eintragen kann, und ueberschreibt ihn im schlimmsten Fall. internal/cli/hook.go begruendet fuer den Claude-Code-Fall ausdruecklich das Gegenteil ("settings.json is the user's file"), und ~/.jaira/settings.json ist zwar jairas Datei, aber der eingetragene Hook ist es nicht.

Empfehlung: B. Erreichbarkeit ist der ganze Zweck des Tickets, also faellt A. Der Unterschied zwischen B und C ist ein einziger cp-Handgriff, C bezahlt ihn mit einem Schreibzugriff auf fremde Konfiguration - das ist kein Tausch, den ein Beispiel wert ist.
- **2026-09-13 21:43 · Alexander Sacharov** — Die drei offenen Fragen, entschieden.

Frage 1 - eingebettet und per Befehl, oder nur im Repository: eingebettet, per Befehl auf stdout. Begruendung oben (Weg B). Der Befehl heisst nicht "jaira hook print" - das ist vergeben - sondern muss den move/claim-Hook benennen, Vorschlag "jaira hook example". Umbenennen von "hook print" ist kein Teil dieses Tickets.

Frage 2 - wie viele Beispiele: genau eines, und zwar keines, das auf diesem Rechner haengt. ~/.jaira/notify.sh braucht HERDR_ENV und HERDR_BIN_PATH und beendet sich bei allen anderen mit exit 0 - als mitgeliefertes Beispiel waere es ein Skript, das fuer fast jeden nichts tut, und das ist schlimmer als keins, weil man den Fehler nicht sieht. Was ohne Konto, ohne URL und ohne Installation laeuft, ist die Terminal-Glocke; ntfy.sh braucht ein Topic und Slack eine Webhook-URL, beides ist Einrichtung und kein Beispiel. Also: ein Skript, dessen Zustellzeile die eine offensichtlich zu ersetzende Zeile ist, mit ntfy und einem Desktop-Toast als auskommentierte Einzeiler direkt darunter. Eine Datei, drei Zustellwege, kein Menue aus Befehlen - und damit auch die Antwort auf "ist das kleiner als paca".

Frage 3 - schreibt der Befehl settings.json: nein. Er druckt und nennt die fehlende Zeile ("hook": "/home/<du>/.jaira/notify.sh"), so wie "jaira hook print" es vormacht. Kein --write in der ersten Fassung.

Die Regel, die das Beispiel tragen muss, und was sie konkret heisst: einen Ton verdient nur der Zustand, in dem sich ohne den Menschen nichts bewegt. Das sind die Lanes mit requires-question oder requires-human-exit (human, signoff) - dort steht die Arbeit bis jemand entscheidet. done bekommt einen anderen, abschliessenden Ton, weil es das Ende ist und nicht die Bitte. blocked ist der Grenzfall: die Lane wartet laut Definition auf eine externe Abhaengigkeit, nicht auf den Menschen - das handgeschriebene notify.sh gibt ihr trotzdem request. Fuer das mitgelieferte Beispiel bleibt blocked stumm, weil sonst die Regel, die es lehren soll, im Beispiel selbst schon gebrochen ist. Alles andere - jeder Wechsel zwischen Agenten-Lanes - ist stumm.

Offen und bewusst nicht hier geloest: dass die Lane-Namen im Beispiel hart stehen muessen, weil "jaira lanes --json" die Gates nicht ausgibt. Auf einem Board mit eigenen Lanes ist das Beispiel dann stumm statt falsch, was die ertraeglichere Richtung ist. Gehoert als eigenes Ticket erfasst, nicht hier.
- **2026-09-13 21:47 · Alexander Sacharov** — Warum der Plan so aussieht (Lane pre-process; gelesen: internal/cli/hook.go, internal/cli/hook_test.go, core/hook/hook.go, core/release/NOTES.md, README.md 338-350 und 550-567, ~/.jaira/notify.sh).

Drei Entscheidungen, die die Brainstorm-Notizen offen gelassen haben und die der Plan jetzt festlegt:

1. Wo das Skript liegt: core/hook/example/notify.sh, eingebettet von einer neuen Datei core/hook/example.go. Nicht in internal/cli, weil jaira seine Vorlagen in core/ einbettet (core/lane/builtin/*.md, core/release/NOTES.md) und internal/cli heute kein go:embed hat. core/hook/hook.go selbst wird nicht angefasst - das Ticket verbietet Aenderungen am Vertrag, eine zusaetzliche Datei im selben Paket ist keine.

2. Warum die Terminal-Glocke die Zustellzeile ist und nicht nur ein Kommentar: sie ist das Einzige, was ohne Konto, URL und Installation wirklich etwas tut - und damit das Einzige, was ein Test beobachten kann. Ein Beispiel, dessen aktive Zeile auskommentiert ist, laeuft zwar fehlerfrei durch, aber die DoD-Zeile "ein Lane-Wechsel in eine Lane, die einem Menschen gehoert, ist hoerbar von einem gewoehnlichen unterscheidbar" waere dann nicht pruefbar. Mit der Glocke ist sie es: das Byte ist da oder nicht.

3. Wie getestet wird: exakt nach dem Muster, das internal/cli/hook_test.go fuer 'hook print' schon aufgebaut hat - die Ausgabe des Befehls wird als das ausgefuehrt, was der Nutzer einfuegt, nicht eine Go-Umschreibung davon. Der Unterschied: hier gibt es keinen Stub auf PATH, sondern gesetzte JAIRA_*-Variablen und die Frage, ob Ausgabe entsteht.

Was der Plan bewusst nicht tut: 'jaira hook print' umbenennen. Die zwei Befehle unter 'hook' meinen weiterhin Unterschiedliches (Claude-Code-Stop-Hook vs. move/claim-Hook); der Plan loest das mit einem Halbsatz in der README statt mit einer Umbenennung, die kein Teil dieses Tickets ist.

Reihenfolge: Skript zuerst, dann Einbettung, dann Befehl, dann Tests. Umgekehrt haette man einen Befehl, der eine Datei druckt, deren Inhalt noch nicht entschieden ist - und die Zustellzeile ist die einzige echte Entscheidung im Skript.
- **2026-09-13 21:52 · Alexander Sacharov** — Umsetzung (Lane in-progress). Zwei Dinge, die der Code nicht sagt:
1. Die Glocke geht nach /dev/tty, nicht nach stdout, weil core/hook/hook.go Stdout und Stderr des Skripts auf nil setzt - alles, was ein Hook nach stdout schreibt, ist im echten Betrieb weg. Wer das Beispiel spaeter auf 'echo' umstellt, macht es damit stumm.
2. Der Test muss dem Skript das Controlling Terminal wegnehmen (SysProcAttr{Setsid:true}), sonst schreibt es genau deswegen nach /dev/tty und der Test sieht leeres stdout. Erster Versuch war '[ -w /dev/tty ]' als Probe: untauglich, weil access() nur den Geraeteknoten prueft und auch ohne Terminal wahr ist. Jetzt oeffnet eine Subshell die tty wirklich - '( : >/dev/tty ) 2>/dev/null' - und ein Fehlschlag kostet nur die Subshell, nicht das Skript.
Deshalb traegt internal/cli/hook_example_test.go '//go:build unix'; Setsid gibt es unter Windows nicht.
- **2026-09-13 21:56 · Alexander Sacharov** — Lane-Abweichung, damit es niemand zweimal sucht: CLAUDE.md beschreibt Lanes critique, optimize und testing - auf diesem Board sind sie nicht installiert. 'jaira lanes' kennt nur backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked. Der Zug in-progress -> critique wurde abgelehnt. Ziel wurde deshalb review: human liegt zwar mit Rang 40 dazwischen, hat aber requires-question: true und ist die Lane fuer eine offene Entscheidung - es gab keine.
Folgeticket 3MJNYS erfasst (Lane-Rollen statt harter Lane-Namen im Beispielskript). Es liegt nur auf seiner Ref, nicht auf der Platte: 'jaira pull 3MJNYS'.
- **2026-09-13 22:03 · Alexander Sacharov** — Review schickt zurueck nach in-progress (gelesen: der ganze Diff 103f551^..0945a8a, dazu core/hook/hook.go und internal/cli/refs.go fuer den Aufrufweg; von Hand ausgefuehrt: das gedruckte Skript mit env -i und leerem PATH fuer human, signoff, done, in-progress, todo, einmal direkt ueber die Shebang-Zeile statt ueber sh; 'go test ./... -race' gruen).

Was am Diff stimmt und nicht noch einmal geprueft werden muss: der Befehl druckt und schreibt nichts, die Einbettung ueber core/hook/example.go laesst core/hook/hook.go unberuehrt, das Skript benutzt ausschliesslich Shell-Builtins und laeuft mit leerem Environment fuer jede Lane mit exit 0 durch, die /dev/tty-Probe als Subshell ist richtig begruendet (core/hook verwirft stdout des Skripts), Prozentzeichen im Titel sind harmlos, es steht genau eine Zeile unter ## Unreleased, und der Diff enthaelt nichts Ueberfluessiges.

Warum es trotzdem zurueckgeht - drei Textstellen, an denen das Skript der Regel widerspricht, die es beibringen soll:
1. notify.sh:18-20 sagt 'only a state in which nothing moves without a person is worth a sound. Every other lane stays silent', notify.sh:47 laesst done klingeln. done ist die Endlane und wartet auf niemanden. Die Brainstorm-Notiz hat genau diese Falle fuer 'blocked' selbst benannt und die Ausnahme fuer done dann nicht in den Satz eingearbeitet. Der Abschluss-Ton ist eine bewusste Entscheidung und darf bleiben - er muss nur in der Regel stehen, statt sie zu widerlegen.
2. core/release/NOTES.md:17 traegt denselben Fehler nach aussen: 'rings the terminal bell only for the lanes that wait on a person'. Das liest 'jaira update' einem Nutzer vor, und es beschreibt ein Verhalten, das das Binary nicht hat.
3. outcome-resolves behauptet, human/signoff seien von done hoerbar unterscheidbar. Der Unterschied sind zwei BEL-Zeichen in einem einzigen printf ohne Pause (notify.sh:33-35). VTE/GNOME Terminal und iTerm2 drosseln oder verschmelzen Glocken in kurzem Abstand, das kommt dann als ein Ton an. Der Test zaehlt Bytes, nicht Toene, kann das also nicht widerlegen. Ohne externes Werkzeug ist die Pause in POSIX sh auch nicht herstellbar - die tragfaehige Fassung ist, dass die gedruckte Textzeile den Unterschied traegt und die Glockenzahl nur ein Hinweis ist. Das gehoert so ins Skript und ins outcome, statt behauptet zu werden.

Nicht rueckweisend, aber beim naechsten Durchgang billig mitzunehmen: der Test fuehrt immer 'sh script' aus, nie das Skript direkt - Shebang und x-Bit, auf die hook.Run per exec.Command(script) angewiesen ist, deckt kein Test ab (ich habe den direkten Aufruf von Hand geprueft, er laeuft). Und 'jaira hook example' druckt unter Windows ein /bin/sh-Skript ohne jeden Hinweis darauf, waehrend der Test //go:build unix traegt und dort gar nicht laeuft; Windows wird laut .goreleaser.yaml gebaut.

Nicht am Diff, aber fuer den naechsten Leser: critique, optimize und testing sind auf diesem Board nicht installiert, das Ticket kam von in-progress direkt hierher. Vor diesem Review hat niemand auf die Aenderung geschaut.
- **2026-09-13 22:06 · Alexander Sacharov** — Zweiter Durchgang in-progress, nur die drei Textstellen aus review-gaps. Was dabei herauskam und nicht im Code steht:
- Die Regel des Beispiels hat zwei Ausnahmen, nicht eine. 'Nur ein Zustand, in dem sich ohne einen Menschen nichts bewegt' beschreibt human/signoff, aber nicht done. done wartet auf niemanden und klingelt trotzdem, weil ein Abschluss das zweite ist, das man lieber gesagt bekommt als nachsieht. Der Kopfkommentar nennt das jetzt als bewusste zweite Ausnahme, statt sie 27 Zeilen spaeter vom eigenen case-Zweig widerlegen zu lassen. Wer hier spaeter kuerzt: der Satz muss beide Ausnahmen nennen, sonst ist der Widerspruch zurueck.
- Die hoerbare Unterscheidung human/signoff gegen done kann die Glockenzahl nicht tragen. Zwei BEL gehen im selben printf ohne Pause raus, und VTE/iTerm2 fassen Glocken in kurzem Abstand zu einem Ton zusammen. Nicht reparierbar ohne externes Werkzeug: eine Pause dazwischen wuerde 'sleep' brauchen, also einen Prozess und eine Wartezeit in einem Skript, das jaira nach 5 s abschiesst - dafuer ist der Effekt zu klein. Die ehrliche Fassung steht jetzt in deliver(): den Unterschied traegt die gedruckte Zeile ('jaira human:' gegen 'jaira done:'), der Glockenzaehler ist ein Hinweis. Ein Kanal mit eigenem Titel (ntfy, notify-send) bekommt die Unterscheidung sauber zurueck.
- NOTES.md:17 stand auf 'only for the lanes that wait on a person' und beschrieb damit ein Binary, das es nicht gibt. Jetzt nennt die Zeile beide Faelle und sagt, woran man sie auseinanderhaelt. Auf dieser Datei arbeitet parallel ein anderer Branch, deshalb nur diese eine Zeile ersetzt, nichts drumherum angefasst.
Nicht angefasst, weil review es ausdruecklich fuer richtig befunden hat: Befehl, go:embed, /dev/tty-Probe, die drei Tests. Die zwei kleinen Anmerkungen von review (Windows druckt ein /bin/sh-Skript; kein Test ruft das Skript ueber die Shebang-Zeile auf) sind bewusst offen - sie waren nicht rueckweisend und gehoeren in ein eigenes Ticket.
- **2026-09-13 22:11 · Alexander Sacharov** — Zweiter Review-Durchgang: freigegeben. Die drei Maengel sind an den benannten Stellen behoben und einzeln nachgeprueft - Kopfkommentar (notify.sh:18-24) deckt sich jetzt mit dem case-Zweig (55-59), NOTES.md:17 beschreibt das Verhalten des Binaries, und die Behauptung ueber zwei hoerbare Toene ist zurueckgenommen statt vorgetaeuscht. Geprueft wurde am Skript, das 'jaira hook example' wirklich druckt, nicht an der Quelldatei.

Was fuer den naechsten Leser nicht im Diff steht: es gibt eine vierte Stelle mit derselben Aussage, README.md:353-357, und sie ist nicht mitgezogen worden. Dort steht weiterhin die Glockenzahl als die Unterscheidung und eine Begruendung, die nur human/signoff deckt, waehrend der Satz done mit aufzaehlt - genau der Bruch, wegen dessen das Ticket beim ersten Mal zurueckging, nur milder. Ich habe deswegen nicht zurueckgewiesen: der erste Durchgang hat diesen Absatz gesehen und abgenommen, und dieser Durchgang sollte ausdruecklich nur die drei benannten Stellen pruefen. Wer abnimmt, entscheidet, ob die Zeile noch hier nachgezogen wird oder ein eigenes Ticket bekommt - sie ist Nutzertext, also nicht folgenlos.

Weiterhin offen und bewusst nicht angefasst: Windows druckt ein /bin/sh-Skript ohne Hinweis, und kein Test ruft das Skript ueber Shebang und x-Bit auf.
