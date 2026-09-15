---
id: 01M2HZKCQ48BPQYN1YPTXB97NJ
title: "Ein Board sagt, ob Arbeit durch einen Pull Request kommt oder direkt auf den Hauptzweig"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Wer allein an einem Repository arbeitet, bekommt keine Rolle vorgesetzt, die ihn zu einem Pull Request schickt, den er nicht will - und ein Team, das Pull Requests benutzt, bekommt die Regel weiterhin ueberall, wo ein Agent sie liest."
definition-of-done: "Ein Board haelt fest, wie Arbeit bei ihm ankommt - durch einen Pull Request oder direkt auf dem Hauptzweig - und die Rollen richten sich danach. Nachgestellt an zwei Boards mit verschiedener Einstellung: auf dem einen schiebt die Rolle den Zweig und haelt an, auf dem anderen fuehrt sie die Arbeit ohne Pull Request zu Ende."
tags:
  - cli
blocked-by: []
related:
  - 01M2HX6SYCFB27V26R5APWAF22
commits: []
created-at: 2026-09-15T07:30:45Z
updated-at: 2026-09-15T07:31:22Z
assignee: ""
updated-by: Alexander Sacharov
context: |-
  13VMA8 hat die Regel 'ein Agent schiebt den Zweig und haelt an, der Mensch macht den Pull Request auf' am 2026-09-15 verbindlich gemacht - in CLAUDE.md, AGENTS.md, dem README und in den ausgelieferten Rollen-Prompts. Fuer ein Team ist das richtig. Fuer einen Menschen, der allein an einem Repository arbeitet und auf den Hauptzweig committet, ist es schlicht falsch: die Rolle haelt ihn an und reicht ihm eine 'gh pr create'-Zeile fuer etwas, das er nie wollte.

  Alex am 2026-09-15: ein Projekt kann git benutzen und trotzdem ohne Pull Requests arbeiten; das Board soll sagen, was bei ihm gilt.

  Wo die Einstellung NICHT hingehoert, und das ist an zwei Faellen aus denselben zwei Tagen gelernt:
  - ~/.jaira/settings.json gilt pro Rechner. Genau daran ist 9ET6NC gescheitert: ein 'remote' fuer alle Boards hat jedes andere Board lahmgelegt.
  - 'git config' gilt pro Klon. Dort liegt 'jaira.forge' richtig, WEIL es persoenlich ist - in welchen Remote DU schiebst, ist deine Sache. 'Dieses Projekt arbeitet mit Pull Requests' ist nicht persoenlich: gilt es fuer das Team, gilt es fuer alle. Es muss mit dem Repository reisen.

  Der Mechanismus dafuer ist schon da: der jaira-Block in CLAUDE.md und AGENTS.md wird ERZEUGT. 'jaira update' kann den Absatz schreiben, der zur Einstellung des Boards passt.

  Damit steht die Frage, die dieses Ticket beantworten muss: die Rollen-Prompts werden fuer alle gleich ausgeliefert und koennen nicht pro Projekt verschieden sein. Wie erfaehrt ein ausgelieferter Prompt, dass dieses Board es anders haelt? Entweder er liest die Einstellung zur Laufzeit, oder der erzeugte Block gilt vor ihm und der Prompt sagt das ausdruecklich. Das ist die Entscheidung, nicht die Umsetzung.

  Verwandt, aber nicht dasselbe: PWAF22 will, dass 'jaira update' das ganze Setup nachzieht. Wenn dieses Ticket dem Befehl etwas zu schreiben gibt, bekommt PWAF22 einen Grund mehr zu existieren.
---

# Ein Board sagt, ob Arbeit durch einen Pull Request kommt oder direkt auf den Hauptzweig

## Definition of Done

- [ ] Ein Board haelt fest, wie Arbeit bei ihm ankommt - durch einen Pull Request oder direkt auf dem Hauptzweig - und die Rollen richten sich danach. Nachgestellt an zwei Boards mit verschiedener Einstellung: auf dem einen schiebt die Rolle den Zweig und haelt an, auf dem anderen fuehrt sie die Arbeit ohne Pull Request zu Ende.
- [ ] Die Einstellung reist mit dem Repository, nicht mit dem Rechner und nicht mit dem Klon: wer klont, erbt sie. Nachgestellt an einem frischen Klon.
- [ ] Ein ausgelieferter Rollen-Prompt richtet sich danach, obwohl er fuer alle gleich ausgeliefert wird - und wie er davon erfaehrt, steht begruendet im Ticket, nicht nur im Code.
- [ ] Ein Board ohne Einstellung verhaelt sich wie heute: durch einen Pull Request. Wer nichts tut, merkt nichts.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

