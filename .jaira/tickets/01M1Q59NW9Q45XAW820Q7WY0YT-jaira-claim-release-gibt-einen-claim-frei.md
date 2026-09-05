---
id: 01M1Q59NW9Q45XAW820Q7WY0YT
title: jaira claim --release gibt einen Claim frei
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "jaira claim <id> --release loescht den Claim auf dem Ticket - den eigenen sofort, einen fremden nur wenn er abgelaufen ist - und die Boards hoeren auf, tote Claims bei jedem Befehl zu melden"
context: "Berk am 04.09. (Feedback L123): verwaiste Claims toter Sessions noergeln in jeder Ausgabe (SGPDYK trug 162h einen), und es gibt keinen Befehl, sie loszuwerden - man wartet oder editiert. Lease-Dauer (30min) bleibt unangetastet (Design: abgestuerzte Sessions geben schnell frei); nur der Aufraeumweg fehlt. Fremde AKTIVE Claims bleiben tabu (Ownership-Rail)."
definition-of-done: claim --release entfernt den eigenen Claim; einen fremden nur abgelaufenen ebenso; ein fremder aktiver wird verweigert; Tests decken alle drei Faelle; --help erklaert es
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:49Z
updated-at: 2026-09-05T02:26:32Z
claimed-by: EE-3NX6GL3-397647
claimed-at: 2026-09-05T02:25:14Z
updated-by: BeMuCa
outcome-what: "claim --release abgesichert und dokumentiert: eigener Claim frei, fremder nur abgelaufen (Guard mit Halter+Alter in der Verweigerung), --help-Langtext erklaert den Befehl erstmals"
outcome-why: "Feedback L123 - Befund beim Bauen korrigiert: --release existierte schon, war aber undokumentiert UND loeste fremde aktive Claims bedingungslos; der fehlende Teil war Schutz+Sichtbarkeit, nicht der Befehl"
outcome-resolves: "Drei DoD-Faelle je mit Test (runCLI, echter Befehl); Lease-Dauer unangetastet (Design); der abgelaufene SGPDYK-162h-Fall von gestern waere damit ein Einzeiler gewesen"
executed-by: fable
---

# jaira claim --release gibt einen Claim frei

## Definition of Done

- [x] claim --release entfernt den eigenen Claim; einen fremden nur abgelaufenen ebenso; ein fremder aktiver wird verweigert; Tests decken alle drei Faelle; --help erklaert es
  proof: claim.go: Release-Guard (fremd+aktiv verweigert mit Halter+Alter, 'still live'); Long-Text erklaert --release; drei Tests: eigener/fremd-abgelaufen/fremd-aktiv, alle -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-05 02:25 · BeMuCa** — Befund beim Lesen: 'claim --release' EXISTIERT bereits (claim.go, Flag registriert) - mein Feedback L123 war insofern falsch informiert; es fehlte aber der Schutz: der Release-Zweig loest auch fremde AKTIVE Claims bedingungslos (keine ClaimActive-Pruefung), und weder --help-Langtext noch der verwaltete Block erwaehnen ihn - daher unauffindbar. Der Schnitt wird: Guard (fremd+aktiv -> Verweigerung mit Halter+Alter, eigener/abgelaufener -> frei), Long-Text-Absatz, drei Tests (eigener/fremd-abgelaufen/fremd-aktiv).
