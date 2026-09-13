---
id: 01M2E86MS66ZZ58RRV3EY0A7DT
title: "Ein Beispiel-Hook liegt bei, damit Lane-Wechsel jemanden erreichen"
status: review
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
commits: []
created-at: 2026-09-13T20:44:07Z
updated-at: 2026-09-13T22:03:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-60947
claimed-at: 2026-09-13T21:46:35Z
outcome-what: "Neuer Befehl 'jaira hook example' druckt ein lauffaehiges Benachrichtigungs-Skript fuer die Einstellung \"hook\". Das Skript liegt als echte Datei unter core/hook/example/notify.sh und wird per go:embed ueber core/hook/example.go ins Binary gezogen; core/hook/hook.go bleibt unberuehrt, der Vertrag aendert sich nicht. Es klingelt zweimal fuer human und signoff, einmal fuer done und bleibt fuer alle Agenten-Lanes und fuer claim stumm; die Glocke geht nach /dev/tty, weil core/hook stdout des Skripts verwirft. Tests in internal/cli/hook_example_test.go fahren das gedruckte Skript wirklich aus. README.md nennt das Beispiel im hook-Absatz und grenzt es gegen 'hook print' ab, eine Zeile steht unter ## Unreleased in core/release/NOTES.md."
outcome-why: "Wer 'hook' in den Einstellungen sah, hatte ein leeres Feld und keinen Anhaltspunkt, was hineingehoert - praktisch schrieb es deshalb niemand. Jetzt gibt es einen ersten Schritt, der ohne jede Installation laeuft, und mit ihm die Regel, die das Beispiel transportieren soll: einen Ton verdient nur der Zustand, in dem sich ohne den Menschen nichts bewegt."
outcome-resolves: "Die DoD-Zeile ist abgehakt: der Befehl legt das Beispiel ab und nennt die Scharfschalt-Zeile (internal/cli/hook.go, newHookExampleCmd), das Skript laeuft mit leerem PATH fehlerfrei durch und tut nichts (TestHookExampleRunsOnAMachineWithNothingInstalled), human/signoff sind von done hoerbar unterscheidbar (TestHookExampleSoundsOnlyForThePersonsLanes), der Hinweis steht im hook-Absatz der README, eine Zeile steht in core/release/NOTES.md, und go test ./... -race ist gruen (exit 0)."
review-summary: "Neuer Unterbefehl 'jaira hook example' (internal/cli/hook.go:newHookExampleCmd) druckt ein 65-zeiliges POSIX-Shellskript auf stdout und sonst nichts - keine Datei, keine Einstellung. Das Skript liegt echt unter core/hook/example/notify.sh und wird von der neuen Datei core/hook/example.go per go:embed ins Binary gezogen; core/hook/hook.go ist unveraendert, der Vertrag steht. Das Skript beendet sich sofort bei JAIRA_EVENT=claim, waehlt dann nach JAIRA_STATUS: human und signoff geben zwei BEL-Zeichen plus eine Textzeile, done eines, jede andere Lane exit 0 ohne Ausgabe. Ausgegeben wird nach /dev/tty, sofern eine Subshell die tty oeffnen kann, sonst nach stdout - weil core/hook/hook.go Stdout und Stderr des Skripts auf nil setzt und alles nach stdout im echten Betrieb verschwindet. Drei Tests in internal/cli/hook_example_test.go (//go:build unix) fahren das gedruckte Skript wirklich aus, einer davon mit leerem PATH. README.md ergaenzt den hook-Absatz und grenzt 'hook example' gegen das unverwandte 'hook print' ab; eine Zeile steht unter ## Unreleased."
review-gaps: |-
  Drei Mangel, alle an derselben Stelle: die Regel, die das Beispiel transportieren soll, stimmt an drei Stellen nicht mit dem ueberein, was das Skript tut.

  1. core/hook/example/notify.sh:18-20 sagt im Kopf: 'only a state in which nothing moves without a person is worth a sound. Every other lane stays silent.' Zeile 47 laesst 'done' klingeln. done ist die Endlane - dort bewegt sich weder mit noch ohne Menschen etwas, es ist keine Bitte. Die Kopfzeile wird also vom eigenen Code 27 Zeilen spaeter widerlegt. Das ist nicht kosmetisch: das Ticket existiert, um genau diese Regel weiterzugeben, und wer das Skript liest, liest zuerst die Regel und dann den Gegenbeweis. Die Brainstorm-Notiz hat die Gefahr fuer 'blocked' selbst benannt ('sonst ist die Regel, die es lehren soll, im Beispiel selbst schon gebrochen') und dieselbe Ausnahme fuer done dann nicht in den Satz eingearbeitet.

  2. core/release/NOTES.md:17 sagt, das Skript klingele 'only for the lanes that wait on a person'. Es klingelt auch fuer done, und done wartet auf niemanden. Das ist die Zeile, die 'jaira update' einem Nutzer vorliest - sie beschreibt ein Verhalten, das das Binary nicht hat.

  3. Die hoerbare Unterscheidung human/signoff gegen done haengt allein an der Anzahl der BEL-Zeichen, und beide gehen in einem einzigen printf ohne Pause hintereinander raus (notify.sh:33-35). Viele Terminals (VTE/GNOME Terminal, iTerm2) fassen Glocken in kurzem Abstand zusammen oder drosseln sie, sodass zwei BEL als ein Ton ankommen. Der Test zaehlt Bytes im String, nicht Toene - er kann diesen Unterschied nicht widerlegen. Die DoD-Zeile im engeren Sinn (Menschen-Lane gegen gewoehnlichen Wechsel = Ton gegen Stille) ist erfuellt; die feinere Unterscheidung, die outcome-resolves ausdruecklich behauptet ('human/signoff sind von done hoerbar unterscheidbar'), ist vom Diff nicht gedeckt. Ohne externes Werkzeug laesst sie sich auch nicht herstellen - die ehrliche Fassung ist, dass die Textzeile ('jaira human:' gegen 'jaira done:') den Unterschied traegt und der Glockenzaehler nur ein Hinweis ist. Das gehoert dann so ins Skript und ins outcome, statt behauptet zu werden.

  Kleiner, nicht rueckweisend: 'jaira hook example' druckt auch unter Windows ein /bin/sh-Skript, ohne das irgendwo zu sagen (Windows wird gebaut, siehe .goreleaser.yaml und core/selfupdate/replace_windows.go); der Test traegt //go:build unix, dort ist der Befehl also ganz ungetestet. Und der Test fuehrt 'sh script' aus, nie das Skript direkt - die Shebang-Zeile und das x-Bit, auf die hook.Run per exec.Command(script) angewiesen ist, sind von keinem Test gedeckt (ich habe den direkten Aufruf von Hand geprueft, er laeuft).

  Kein Defekt gefunden: das Skript laeuft mit leerem PATH und leerem Environment fehlerfrei durch (von Hand nachgestellt, exit 0 fuer alle Lanes), benutzt nur Shell-Builtins, ist gegen Prozentzeichen im Titel sicher, und 'go test ./... -race' ist gruen. Es steht genau eine Zeile fuer dieses Ticket unter ## Unreleased. Nichts im Diff ist ueberfluessig.

  Anmerkung zum Weg, nicht zum Diff: critique, optimize und testing sind auf diesem Board nicht installiert, das Ticket ist von in-progress direkt nach review gegangen. Niemand hat vor mir draufgeschaut.
review-verdict: "Zurueck nach in-progress. Der Mechanismus ist richtig gebaut und laeuft - der Befehl, die Einbettung, die Tests und die leere Maschine stimmen alle. Was nicht stimmt, ist genau das, wofuer das Ticket da ist: das Beispiel soll eine Regel weitergeben, und es widerspricht ihr in seinem eigenen Kopfkommentar (notify.sh:18-20 gegen :47) und in der Release-Notiz, die der Nutzer zu lesen bekommt (NOTES.md:17 behauptet 'only for the lanes that wait on a person', done klingelt trotzdem). Dazu behauptet outcome-resolves eine hoerbare Unterscheidung human/done, die nur aus zwei pausenlosen BEL-Zeichen besteht und die kein Test pruefen kann. Drei Textstellen, keine Architektur - aber es sind die drei Stellen, an denen das Ticket seinen Zweck erfuellt oder nicht. Zu reparieren in einer Sitzung: den Abschluss-Ton als bewusste zweite Ausnahme in die Regel schreiben, NOTES.md:17 auf das korrigieren, was das Skript tut, und im Skript sagen, dass die Textzeile den Unterschied traegt."
review-check: |-
  1. In /home/alex/projects/jaira-Y0A7DT: 'PATH=$PATH:/usr/local/go/bin go run ./cmd/jaira hook example > /tmp/notify.sh && chmod +x /tmp/notify.sh'. Es entsteht eine Datei, es erscheint keine Ausgabe.
  2. Zeilen 18 bis 20 von /tmp/notify.sh lesen. Dort steht: nur ein Zustand, in dem sich ohne einen Menschen nichts bewegt, ist einen Ton wert - jede andere Lane bleibt stumm.
  3. Jetzt Zeile 47 derselben Datei lesen. Dort steht 'done) tone=finished'. done ist die Endlane und wartet auf niemanden. Das ist der Widerspruch: der Kopf sagt stumm, der Code klingelt.
  4. Zeile 17 von core/release/NOTES.md lesen: 'rings the terminal bell only for the lanes that wait on a person'. Das ist dieselbe Aussage, und sie ist genauso falsch - der Nutzer bekommt sie von 'jaira update' vorgelesen.
  5. In einem echten Terminalfenster (nicht in einem Editor-Panel) nacheinander ausfuehren: 'env -i JAIRA_EVENT=move JAIRA_STATUS=human JAIRA_TITLE=Test /tmp/notify.sh', dann dasselbe mit JAIRA_STATUS=done. Hinhoeren: klingt der erste Aufruf hoerbar nach zwei Toenen und der zweite nach einem? Auf vielen Terminals verschmelzen die zwei zu einem - dann ist die Unterscheidung, die outcome-resolves behauptet, auf diesem Rechner nicht da, und nur die gedruckte Textzeile ('jaira human:' gegen 'jaira done:') trennt die Faelle.
  6. Gegenprobe, dass sonst nichts kaputt ist: 'env -i PATH= JAIRA_EVENT=move JAIRA_STATUS=in-progress /tmp/notify.sh; echo $?'. Es erscheint keine Ausgabe und es steht 0 da - das Skript laeuft auf einer Maschine ohne jedes Werkzeug durch und tut nichts.
  7. 'PATH=$PATH:/usr/local/go/bin go test ./... -race'. Alle Pakete melden ok, kein FAIL.
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
