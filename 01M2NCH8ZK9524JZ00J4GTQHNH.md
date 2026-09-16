---
id: 01M2NCH8ZK9524JZ00J4GTQHNH
title: "Der Dispatcher bekommt einen Gespraechsmodus, statt dass eine zweite Rolle daneben entsteht"
status: brainstorm
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket, in dem Gestaltungsentscheidungen noch offen sind, wird im Gespraech mit dem Menschen gefuehrt statt geraten: der Dispatcher haelt vor der Plan-Lane an, fragt, schreibt die Antwort auf das Ticket, zeigt den Code bevor darauf aufgebaut wird, und committet nicht von allein - und das ist ein Modus derselben Rolle, keine zweite Rolle daneben"
context: |-
  Alex am 2026-09-16, aus zwei Laeufen dieser Woche.

  Was schiefgeht: laesst man ein Ticket autonom durch die Lanes laufen, dessen Form noch nicht feststeht - UI-Aenderungen und Datenbank-/Dateiformat-Aenderungen sind seine Beispiele -, dann raet der Worker. Am Ende steht viel Arbeit, die neu gemacht werden muss.

  Der Beleg liegt auf diesem Board. 9ZZSFT lief autonom in einem Durchgang durch. 0YGWXQ brauchte sieben critique-Runden und drei Entscheidungen von Alex. 0YGWXQ ist nicht groesser - es war weniger entschieden.

  Damit ist Alex' erster Vorschlag ('wenn ein Ticket gross aussieht') das falsche Mass. Das Mass ist, WIE VIELE ENTSCHEIDUNGEN IM TICKET NOCH OFFEN SIND: ein DoD-Punkt, den man auf zwei verschiedene Arten erfuellen kann und beide kommen durch den Gate, ist eine offene Entscheidung. Null davon heisst, der Lauf braucht niemanden.

  Warum es KEINE zweite Kommandozeile wird, obwohl Alex zuerst danach gefragt hat: was er beschreibt, existiert schon - es ist der Teamlead, auf ein Ticket verengt. Genau so lief diese Woche: mit Alex ueber die Form von 0YGWXQ reden, seine Antwort mit 'jaira note' auf das Ticket schreiben, dann Worker starten. Zwei Prompts, die fast dasselbe tun, driften auseinander. Dieses Repository leidet schon daran: 3YRPXJ und 7MG5GB beschreiben beide, dass spawn.sh nicht erweiterbar ist und deshalb geforkt wird. Alex hat dem Modus-statt-Rolle am 16.09. zugestimmt.

  Zwei Eigenschaften hat Alex am 16.09. nachgereicht, und die zweite hat eine Falle:

  Erstens, der Mensch sieht den Code, bevor darauf aufgebaut wird. Nicht am Ende, wenn das Zurueckdrehen am teuersten ist.

  Zweitens, in diesem Modus wird nicht automatisch committet. Die Falle: 'kein Commit' ist die Abwesenheit eines Mechanismus, nicht einer. Committet der Mensch von Hand, faellt mit hoher Wahrscheinlichkeit der Handle aus der Betreffzeile - und CLAUDE.md sagt seit 9ZZSFT ausdruecklich, dass genau dieser Handle die Quelle ist, aus der jaira die Commit-Liste ableitet, weil die Historie der Ticket-Datei bewusst duenn geworden ist. Ohne Handle bleibt die Liste leer und der Zug in die Endlane wird verweigert. Der Agent muss die fertige Commit-Zeile samt Handle zurueckgeben, so wie jaira-role-pr es mit der PR-Beschreibung schon macht - das Muster gibt es also.

  Noch nicht entschieden, Sache der Brainstorm-Lane: in welchen Haeppchen der Mensch den Code sieht. Der DoD-Punkt ist die naheliegende Einheit, weil er ohnehin das Inkrement ist. Zu fein und der Modus ist langsamer, als es selbst zu tippen.

  Nicht vor dem naechsten Release. 0YGWXQ und 7KX89C sind ein erklaerbares Paar; diese Rolle ist der Release danach. Und sie entsteht im Repository, nicht in ~/.claude - 7KX89C raeumt gerade genau das auf.
definition-of-done: "Der Eintritt in den Modus haengt an einer nachpruefbaren Bedingung, nicht am Bauchgefuehl: vor der Plan-Lane zaehlt der Dispatcher die noch offenen Entscheidungen des Tickets auf. Keine offene - er laeuft weiter wie heute. Mindestens eine - er haelt an und fragt den Menschen."
tags:
  - docs
blocked-by: []
related:
  - 01M2KBPPVH5PAKZ98B0X7KX89C
commits: []
created-at: 2026-09-16T15:14:31Z
updated-at: 2026-09-16T15:28:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-53684
claimed-at: 2026-09-16T15:25:11Z
---

# Der Dispatcher bekommt einen Gespraechsmodus, statt dass eine zweite Rolle daneben entsteht

## Definition of Done

- [ ] Der Eintritt in den Modus haengt an einer nachpruefbaren Bedingung, nicht am Bauchgefuehl: vor der Plan-Lane zaehlt der Dispatcher die noch offenen Entscheidungen des Tickets auf. Keine offene - er laeuft weiter wie heute. Mindestens eine - er haelt an und fragt den Menschen.
- [ ] Es gibt weiterhin genau eine Dispatcher-Rolle. Kein zweiter Skill und keine zweite Kommandozeile daneben; der Modus steht im selben Prompt.
- [ ] Was im Gespraech entschieden wird, steht mit 'jaira note' auf dem Ticket, BEVOR die Arbeit daran beginnt - nicht hinterher. Nachgestellt an einem Ticket, dessen Sitzung mittendrin abgebrochen wird: die Entscheidung ist danach noch da.
- [ ] In diesem Modus committet der Agent nicht selbst. Er legt die Aenderungen bereit und gibt eine fertige Commit-Zeile zurueck, die den Ticket-Handle im Betreff traegt und die Ticket-Datei mitnimmt. Nachgestellt: nach dem Commit des Menschen leitet jaira die Commit-Liste vollstaendig ab und der Zug in die Endlane wird nicht verweigert.
- [ ] Der Mensch sieht den Code, bevor darauf aufgebaut wird: der Dispatcher legt ihn nach jedem Inkrement vor und wartet, statt am Ende alles auf einmal zu zeigen.
- [ ] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer den Modus und fuer das, was ein Benutzer beim Committen anders tut.

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-16 15:28 · Alexander Sacharov** — Brainstorm, Befund aus dem Code — nicht aus der Notiz.

Die Rolle ist seit 7KX89C eine Datei im Repository: core/role/builtin/jaira-dispatcher/SKILL.md, ausgeliefert per 'jaira roles install'. Der Modus ist damit eine Prompt-Aenderung an EINER Datei, kein Go-Code. Die Kopie in ~/.claude ist heute byte-identisch (diff leer), es gibt also noch keine Drift zu heilen.

Was der Prompt heute kann und was nicht: der Abschnitt 'When to stop and hand back' listet drei Haltepunkte — Worker am Freigabedialog, dieselbe Lane dreimal zurueck, Menschen-Lane erreicht. Alle drei sind REAKTIV: der Dispatcher haelt an, nachdem etwas passiert ist. Ein Halt VOR der Plan-Lane, weil im Ticket noch etwas offen ist, existiert nirgends. Das ist die eigentliche Luecke, und sie ist genau die, die DoD-Punkt 1 verlangt.

Der Haken, den die Ticket-Notiz nicht sieht: 'nicht automatisch committen' ist keine Eigenschaft des Dispatchers. Das Committen steht im Worker-Prompt (jaira-role-lane/SKILL.md: 'move the ticket, then git add -A and commit ... whose message names the ticket id'). Ein Modus, der nur im Dispatcher-Prompt steht, erreicht den Worker nie: spawn.sh tippt fest '/jaira-role-lane $ticket $lane' in den Tab und hat keine Stelle, an der ein weiteres Argument mitfaehrt (dieselbe Zeile, die 3YRPXJ aufmacht). Der Modus muss also einen Weg zum Worker haben — das ist die Entscheidung, die dieses Ticket traegt, und sie steht in der naechsten Notiz.

Was NICHT erfunden werden muss: der Kanal zum Menschen gibt es schon. Das Frontmatter-Feld 'question' traegt heute eine Frage an den Menschen (auf 7KX89C im Einsatz), die human-Lane ist der Ort, an dem eine Person antwortet, und 4XHZ6N baut daneben den gesammelten human-note-Kanal. Der Gespraechsmodus soll keinen vierten daneben stellen.
