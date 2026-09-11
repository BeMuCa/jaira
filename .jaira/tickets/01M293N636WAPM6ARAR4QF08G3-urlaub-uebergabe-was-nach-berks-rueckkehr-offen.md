---
id: 01M293N636WAPM6ARAR4QF08G3
title: "Urlaub-Uebergabe: was nach Berks Rueckkehr offen ist"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Alles, was am 11.09. offen blieb und auf Berks Rueckkehr wartet, steht an einer Stelle - mit den Ticket-Handles, die dazugehoeren, und dem was jeweils zu tun ist"
context: "Am 11.09. hat Alex das Kartenlayout der TUI umgebaut. Vier Commits auf worktree-logbook-is-a-decision, noch kein PR. Drei Punkte blieben offen, und alle drei brauchen einen Menschen.\n\nPUNKT 1 - Berk weiss noch nichts, und zwei seiner Entscheidungen sind zurueckgedreht.\nKarten tragen keinen Rahmen mehr; sie sind gefuellte Baender mit der Tag-Farbe in einer Zelle links (QQ3EX4, Commit 784ca78). Das kippt:\n  81XRXX - jede Karte hat einen eigenen Rahmen. Berk am 08.09. woertlich: 'das mit den 2 zeilen ist ok, ich will den rahmen'.\n  VS5DFW - die Tag-Farbe sitzt in der linken Rahmenkante, die Auswahl kriegt den ganzen Rahmen plus Balken.\nBeide sind archiviert, mit einer Notiz auf dem Ticket, warum. 'jaira restore' holt sie zurueck, falls Berk den Rahmen wiederhaben will - dann faellt QQ3EX4.\nDer Tausch, um den es geht: Rahmen weg, dafuer Titelbreite w-6 auf w-4 und Kartenhoehe 5 auf 3 Zeilen. Vorher war auf dem Board kein Titel zu Ende lesbar.\nZu tun: Berk den Stand zeigen. Sagt er ja, ist der Punkt erledigt. Sagt er nein, beide Tickets restoren.\n\nPUNKT 2 - dieses Board legt fertige Tickets noch automatisch ab.\n.jaira/lanes/done.md traegt 'logbook-on-entry: true'. Das ist die alte Tuer, die 74VM40 gerade abgeschafft hat: seitdem ist das Ablegen eine Entscheidung ('jaira logbook --all'), nicht was Fertigwerden nebenbei tut. Auf diesem Board steht sie noch offen, darum sind ANX7Y5, XK4124 und QQ3EX4 beim Uebergang nach done sofort in den Logbuch-Ordner gewandert.\nZu tun: entscheiden ob das Absicht ist. Wenn nicht, die Zeile aus .jaira/lanes/done.md nehmen.\n\nPUNKT 3 - sieben Tickets sind kaputt abgelegt, und jaira meckert bei JEDEM Befehl darueber.\nSechs liegen doppelt, einmal auf dem Board und einmal in .jaira/logbook/as-20260911/. jaira weigert sich zu waehlen: 'one of the two has to go, and only you can say which'. Ausserdem blockiert es das Trimmen der done-Lane.\nGemessen, welche Kopie welchen Stand hat:\n  P1AE82   Board human     Logbuch done\n  76WCCW   Board signoff   Logbuch done\n  NJPQWE   Board signoff   Logbuch done\n  KA9CFA   Board signoff   Logbuch done\n  TYTVBZ   Board signoff   Logbuch done\n  ZEFXXM   Board signoff   Logbuch done\nDie Logbuch-Kopie ist ueberall die weiter fortgeschrittene. Das SIEHT danach aus, dass die Board-Kopie ein Rest vom Ablegen ist und weg kann - aber genau diese Entscheidung will jaira nicht selbst treffen, und ungeprueft zu loeschen ist der eine Weg, auf dem hier wirklich Arbeit verlorengeht. Vor dem Loeschen einmal diffen.\nDazu: 51RNC1 liegt auf seinem ref und nicht auf der Platte - 'jaira pull 51RNC1' holt es.\nZu tun: pro Ticket entscheiden welche Kopie bleibt, die andere loeschen, dann 51RNC1 pullen.\n\nVerwandte Tickets stehen hier im Text und nicht in einem Link-Feld: 'blocked-by' waere eine echte Abhaengigkeit und wuerde dieses Ticket faelschlich blockieren, und 'follows' nimmt nur eins - und nimmt ausserdem kein Ticket an, das schon im Logbuch liegt, was auf QQ3EX4 zutrifft. Das Link-Fenster, das sie alle zeigen wuerde, ist 51RNC1 und liegt selbst noch in testing.\nFertig und im Logbuch: ANX7Y5 (Fuellung), XK4124 (Glow, Taste c), QQ3EX4 (Baender statt Rahmen).\nArchiviert: 81XRXX, VS5DFW (von QQ3EX4 ueberholt), VHQ0F4 (ging in Punkt 1 auf)."
definition-of-done: "Berk hat den Layout-Wechsel gesehen und gesagt was gilt; 'logbook-on-entry' auf diesem Board ist entschieden; die sieben doppelt oder gar nicht abgelegten Tickets sind aufgeloest und jaira meldet beim Start keine unlesbaren Tickets mehr"
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-11T20:48:28Z
updated-at: 2026-09-11T21:01:12Z
updated-by: Alexander Sacharov
---

# Urlaub-Uebergabe: was nach Berks Rueckkehr offen ist

## Definition of Done

- [~] Berk hat den Layout-Wechsel gesehen und gesagt was gilt; 'logbook-on-entry' auf diesem Board ist entschieden; die sieben doppelt oder gar nicht abgelegten Tickets sind aufgeloest und jaira meldet beim Start keine unlesbaren Tickets mehr

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 21:00 · Alexander Sacharov** — PUNKT 3 ist erledigt, und die Ursache ist gefunden - sie war nicht das, wonach es aussah. Die Board-Kopien sind nicht 'nochmal aufgetaucht': sie lagen seit dem 03.09. ununterbrochen dort, git kennt fuer sie kein einziges Loeschen ('git log --diff-filter=D' auf die Pfade ist leer). Was am 11.09. passierte, war das Gegenteil - Commit 733df8f legte fuenf Logbuch-Kopien DANEBEN, als reine Additionen: '5 files changed, 290 insertions(+)', keine einzige Loeschung. P1AE82 kam genauso in 4feaf14. Das Ablegen hat also kopiert statt verschoben, und ausgerechnet in dem Commit, der 'a filed ticket stays filed' im Titel traegt.

Geloescht wurden die sechs Board-Kopien. Vorher geprueft, dass dabei nichts verlorengeht: der Diff jeder Board-Kopie gegen ihre Logbuch-Kopie enthielt ausschliesslich aeltere Werte von vier Feldern - status (signoff bzw. human gegen done), commits (leer gegen gestempelt), updated-at und updated-by. Keine Zeile Text, kein Feld, das nur auf dem Board stand. Die Logbuch-Kopie war ueberall echt die vollstaendigere.

Danach meldet jaira nichts mehr, weder auf stderr noch als Zaehler in der Fusszeile des Boards.

Was offen bleibt und nicht zu diesem Ticket gehoert: ob das Kopieren-statt-Verschieben noch im Code steckt. Geprueft ist nur, dass die Dateien weg sind, nicht dass ein naechstes 'jaira logbook' es nicht wieder tut. Dafuer steht jetzt ein eigenes Ticket.
- **2026-09-11 21:01 · Alexander Sacharov** — Korrektur zur vorigen Notiz: es gibt KEIN eigenes Ticket fuer 'kopiert statt verschoben', weil es nichts zu tun gibt. Nachgesehen, wie die drei heute abgelegten Tickets in git liegen - Commit a860e46 zeigt sie als 'rename .jaira/{tickets => logbook/as-20260911}/...', nicht als Add plus stehengebliebene Datei. Der Code verschiebt also korrekt.

Damit bleibt fuer 733df8f nur eine Erklaerung, die zu allen Befunden passt: dort wurden die Additionen gestaged und die Loeschungen nicht. Auf der Platte war die Datei weg, im Commit nicht - und beim naechsten Checkout war sie wieder da. Das ist ein Bedienfehler beim Committen, kein Fehler im Ablegen, und es erklaert auch, warum git fuer die Board-Kopien nie eine Loeschung kennt.

Wer das kuenftig vermeiden will: 'git add -A .jaira' statt einzelner Pfade, oder vor dem Commit einmal 'git status' auf Loeschungen ansehen.
