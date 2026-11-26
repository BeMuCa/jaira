---
id: 01M85APT0AZR3XRMFNVJHR3ZK3
title: "Marktplatz fuer Lanes und Boards: fuenf Entscheidungen, bevor gebaut wird"
status: human
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Berk hat die fuenf offenen Punkte zum Marktplatz entschieden, und die Antworten stehen als Notiz auf diesem Ticket."
context: "Der Milestone extension-marketplace kann nicht starten, bevor feststeht, wo der Katalog liegt und wer mergt.\nAusloeser: Berks Idee vom 24.11.2026 (Erweiterungs-Repo, Pull Requests, Website, Installation per ID). Der Entwurf steht im Ticket zu docs/MARKET.md (YTT25Y).\n1. Wo liegt der Katalog?\n a) Neues Repo BeMuCa/jaira-extensions. Kosten: ein zweites Katalog-Repo neben dem fuer Plugins, zwei Websites oder eine, die zwei Quellen liest.\n b) Im bestehenden Plugin-Katalog github.com/panda-bobr/jaira-releases, neben plugins.json. Ein Ort fuer alles, was man zu jaira dazuholt, eine Website, eine CI. Kosten: das Repo gehoert Alex' Konto; Berk braucht dort Merge-Rechte.\n c) Bleibt in BeMuCa/jaira/lanes/. Kosten: jeder Beitrag laeuft durch die ganze Go-Testsuite und ist erst nach einem Release sichtbar.\n Empfehlung: b, wenn Berk und Alex dort beide mergen duerfen; sonst a. Nicht c.\n2. Wer reviewt und mergt Beitraege (CODEOWNERS)? Empfehlung: Berk und Alex, einer genuegt.\n3. Darf ein Agent allein eine Lane aus dem Marktplatz uebernehmen? Empfehlung: nein - ohne Terminal nur mit --yes, und die Agenten-Notiz sagt, dass --yes Sache einer Person ist (Ticket R0YRGP). Kosten: ein Mensch muss einmal je Lane ja sagen.\n4. Einreichen ueber die Website: ein vorbelegter GitHub-Link (kein Server, GitHub-Konto noetig) oder ein Portal mit eigener Anmeldung (Server, Token, Spam-Abwehr, Dauerbetrieb)? Empfehlung: der Link; das Portal wird nicht gebaut (Ticket 3WYNV4).\n5. Lanes mit script: im oeffentlichen Katalog - nie, oder nach besonderem Review? Empfehlung: nie; Firmen-Kataloge duerfen es mit --allow-script (Ticket Z9JNAV). Und: Skills verteilt jaira nicht selbst, dafuer gibt es Claude-Code-Plugin-Marktplaetze - einverstanden?"
definition-of-done: Zu jedem der fuenf Punkte steht Berks Antwort als Notiz auf dem Ticket
tags:
  - docs
blocked-by:
  - 01M8BQFBMS778FE9WK6YDAW7WP
related:
  - 01M6NHNRNKP678V30QDJAN1XR5
  - 01M5GFAV1C3DV9NWZFXC4G07X9
commits: []
created-at: 2026-11-23T23:08:55Z
updated-at: 2026-11-26T10:50:33Z
updated-by: Berk
question: "Marktplatz fuer Lanes und Boards - fuenf Entscheidungen (Optionen und Kosten im Kontext): 1) Wo liegt der Katalog: a) neues Repo BeMuCa/jaira-extensions, b) im Plugin-Katalog panda-bobr/jaira-releases neben plugins.json, c) bleibt in BeMuCa/jaira/lanes? Empfehlung b, sonst a. 2) Wer mergt Beitraege? Empfehlung Berk und Alex. 3) Darf ein Agent allein uebernehmen? Empfehlung nein, --yes ist Sache einer Person. 4) Einreichen ueber die Website: vorbelegter GitHub-Link oder Portal mit Server? Empfehlung Link. 5) script:-Lanes im oeffentlichen Katalog nie, und Skills verteilt jaira nicht selbst - einverstanden?"
---

# Marktplatz fuer Lanes und Boards: fuenf Entscheidungen, bevor gebaut wird

## Definition of Done

- [ ] Zu jedem der fuenf Punkte steht Berks Antwort als Notiz auf dem Ticket
- [ ] docs/MARKET.md und die betroffenen Tickets im Milestone sind an die Antworten angepasst

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

