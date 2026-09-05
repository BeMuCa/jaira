---
id: 01M1Q59NW9Q45XAW820Q7WY0YT
title: jaira claim --release gibt einen Claim frei
status: done
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "jaira claim <id> --release loescht den Claim auf dem Ticket - den eigenen sofort, einen fremden nur wenn er abgelaufen ist - und die Boards hoeren auf, tote Claims bei jedem Befehl zu melden"
context: "Berk am 04.09. (Feedback L123): verwaiste Claims toter Sessions noergeln in jeder Ausgabe (SGPDYK trug 162h einen), und es gibt keinen Befehl, sie loszuwerden - man wartet oder editiert. Lease-Dauer (30min) bleibt unangetastet (Design: abgestuerzte Sessions geben schnell frei); nur der Aufraeumweg fehlt. Fremde AKTIVE Claims bleiben tabu (Ownership-Rail)."
definition-of-done: claim --release entfernt den eigenen Claim; einen fremden nur abgelaufenen ebenso; ein fremder aktiver wird verweigert; Tests decken alle drei Faelle; --help erklaert es
tags: []
blocked-by: []
commits:
  - 5f7556f39d124b217c58d3eb66fe180f089150ae
created-at: 2026-09-04T21:30:49Z
updated-at: 2026-09-05T02:29:18Z
claimed-by: EE-3NX6GL3-397647
claimed-at: 2026-09-05T02:25:14Z
updated-by: BeMuCa
outcome-what: "claim --release abgesichert und dokumentiert: eigener Claim frei, fremder nur abgelaufen (Guard mit Halter+Alter in der Verweigerung), --help-Langtext erklaert den Befehl erstmals"
outcome-why: "Feedback L123 - Befund beim Bauen korrigiert: --release existierte schon, war aber undokumentiert UND loeste fremde aktive Claims bedingungslos; der fehlende Teil war Schutz+Sichtbarkeit, nicht der Befehl"
outcome-resolves: "Drei DoD-Faelle je mit Test (runCLI, echter Befehl); Lease-Dauer unangetastet (Design); der abgelaufene SGPDYK-162h-Fall von gestern waere damit ein Einzeiler gewesen"
executed-by: fable
review-summary: "Kritik: der Guard nutzt die existierende ClaimActive-Wahrheit statt einer zweiten Alters-Logik; Verweigerung nennt Halter+Alter (handlungsfaehig); Lease-Design (30min, selbstheilend) bewusst unangetastet - nur der Aufraeumweg wurde sicher und sichtbar."
review-gaps: "Nichts entfernt. Gelassen: --steal unveraendert (bewusst getrennte Geste fuer LEBENDE Claims); kein Auto-Purge alter Claims beim Board-Laden (das Werkzeug bewegt/loescht nichts unaufgefordert)."
review-check: |-
  1. jaira claim <id>; jaira claim <id> --release -> Released, Feld leer.
  2. Fremden Claim simulieren (--session x), claimed-at per set altern: --release raeumt ihn ab - die Noergel-Zeile verschwindet.
  3. Fremden AKTIVEN Claim (--session x, frisch): --release verweigert mit Halter und Alter (still live).
  4. jaira claim --help erklaert --release jetzt.
test-verdict: "pass: voller Lauf -race nach Cache-Loeschung RC=0 (test17, 15 Pakete); DoD am Baum (Guard in claim.go, Long-Text); Verhalten per runCLI in allen drei Faellen wirklich ausgefuehrt"
review-verdict: "accept (koordinator-verifiziert, Direktdurchlauf per Berks Auftrag: kleiner Guard+Doku-Schnitt, drei echte CLI-Faelle getestet, voller Lauf -race RC=0; Feedback-Irrtum offen korrigiert - der Befehl existierte, ihm fehlten Schutz und Sichtbarkeit)"
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
