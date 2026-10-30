---
id: 01M67K6Z3NSBTBXAK4EVCMN1YT
title: "jaira web wird Beobachtungs- und Steuerzentrale fuer PMs ohne Claude: Onboarding mit Anbieterwahl, lokale Tokens, Projekte aus Repos, Tabs fuer Tickets, Aktive, Milestones/Epics, Schreiben"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Ein Projektmanager ohne Claude und ohne CLI-Kenntnisse verfolgt und steuert die Arbeit seiner Entwickler komplett im Browser, und Code und Ticket-Tracking (jaira, Jira, GitHub/GitLab Issues) passen zusammen."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkte 2-5, wortgetreu gespeichert in der Claude-Memory prompt-2026-10-31-web-pm-and-lanes).
  Web-Nutzer heissen ab jetzt PM.
  - Nur git noetig; Infos an bestimmte Entwickler evtl. ueber Git-Refs weiterreichen.
  - curl-Onboarding: Wahl git / Jira / andere grosse Anbieter (YouTrack, GitLab/GitHub Issues, Epics, Projects), dann verbinden (API-Token bzw. git). Fehlt git: Link zu git ausgeben, oder Pfad zu git einfuegen bzw. in eine Umgebungsvariable schreiben. Danach oeffnet das Web.
  - Web-Settings: der PM legt weitere Tokens an; sie bleiben nur bei ihm, werden nirgends geteilt (Sicherheit zuerst).
  - Neues Projekt: Anbieter waehlen; bei git alle gefundenen Repos zur Auswahl oder Clone-Link einfuegen (Default, wenn keins gefunden). Das Projekt wird angelegt, das Repo geholt, und es zeigt Tickets, Milestone-Roadmap mit Strassennamen und die uebrigen Milestones.
  - Begriffe: Milestone = Epic. jaira-Milestones muessen sich als Epics nach Jira uebertragen lassen; im Web muss klar sein, dass Milestones Epics sind.
  - Tabs: 'Alle Tickets'; 'Aktiv' - kleine Avatare je Nutzer am Board, darunter seine Tickets, Farbe nach Stand: backlog grau, in Arbeit orange, in review orange-gruen, done gruen, blockiert rot; 'Milestones/Epics' mit Jira-(o.ae.)-Verbindung; 'Schreiben' - Tickets und Milestones anlegen, auch direkt aus einem Jira-Ticket uebernehmen.
  - Der PM weist Tickets und Milestones zu (ueber den festen jaira-Namen, siehe Ticket zum festen Namen).
  - Start von jaira web gibt eine schoen gestaltete Ausgabe, die den Server vorstellt.
  - Punkt 5 der Nachricht bricht ab: 'nach einer Zuweisung durch Kollegen oder PM siehst du es auf dem Board als ...' - Ende fehlt.
  Heute: jaira web gibt es (Filter, Assignee auf Karten, Lane-Farben, Log - ZDBXN1, TPCKZB, HEY6T7).
definition-of-done: "Brainstorm (starkes Modell): Funktionsliste, Aufteilung in Tickets, Jira-Mapping (Epics/Tasks/Subtasks vs. Strassenname/Milestone/Ticket), Token-Ablage, was fehlt"
tags:
  - cli
blocked-by: []
related:
  - 01M5MNWJF240306ZVDA41K3RA8
commits: []
created-at: 2026-10-30T23:44:38Z
updated-at: 2026-10-30T23:44:38Z
---

# jaira web wird Beobachtungs- und Steuerzentrale fuer PMs ohne Claude: Onboarding mit Anbieterwahl, lokale Tokens, Projekte aus Repos, Tabs fuer Tickets, Aktive, Milestones/Epics, Schreiben

## Definition of Done

- [ ] Brainstorm (starkes Modell): Funktionsliste, Aufteilung in Tickets, Jira-Mapping (Epics/Tasks/Subtasks vs. Strassenname/Milestone/Ticket), Token-Ablage, was fehlt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

