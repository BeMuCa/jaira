# Handoff — 2026-08-31, Abend

Stand nach der Backlog-Flotte; dieses File ist das Gedächtnis. Der vorige
Handoff (Morgen desselben Tags) ist in der git-History dieses Files.

## Kontrakt: Berks erstes "go" nach dem Clear (03.09.)

= die saubere Auflistung aller offenen Punkte, Learnings und Fragen zeigen
(dieses File + Memory jaira-open-decisions), NICHT bauen. Der Schema-Bau
startet weiterhin erst nach seinem Spec-Review.

Neu eingefangen (todo, 8W1R94): Update-Prozess durchspielen und vereinfachen;
Share-Zustand bleibt auf JEDEM Update-Pfad erhalten (ETR0PX deckte nur
'jaira update'+.gitignore; self upgrade, Binary-Wechsel, Migration offen).

## Wo alles steht

- **master gepusht** bis einschließlich des 0.1.1-NOTES-Blocks. ~17 Commits
  heute: Schema-Spec, HREQJR (Wrap überall), dann die Flotte
  (DNAEPN, MFD7P3, MR3CN8, A1TZ4N, YBC0MT, TQXBY5, 88H1P4, GEC3TK).
  `go test ./... -race` (Cache geleert) grün auf dem Push-Stand;
  `gofmt -l` nur das bekannte `internal/cli/tickets.go`.
- **Binary** `~/.local/bin/jaira` neu gebaut; nach dem NOTES-Commit einmal
  nachziehen, falls der Stempel nicht HEAD ist.
- Zwei Spitzen-Commits wurden vor dem Push **re-splittet** (geteilter Index
  zweier paralleler Agents hatte sie verheddert); Baum-Identität bewiesen.

## Was Berk jetzt entscheidet (sein Inbox-Stand)

1. **signoff: 13 Tickets** warten auf Abnahme — die 9 der Flotte (DNAEPN,
   MFD7P3, MR3CN8, A1TZ4N, YBC0MT, TQXBY5, GEC3TK, 88H1P4, QA3GN1) plus
   D9SC5A (Herkunfts-Label), 9H265S (Commits im Signoff), 79GEPW (Tags
   Kern) und AXTFG3 (Tag-Boxen + t-Legende). Abnahme im Board, oder auf
   Zuruf per `--force` (Invariante 11 gilt; Erfolg an der Lane pruefen,
   nicht an der letzten Ausgabezeile).
2. **human: SGPDYK** — done-cap. Empfehlung liegt als Notiz auf dem Ticket
   (Logbuch / holds:10 im Lane-File / updated-at / Trigger beim done-Move,
   gemeldet); ein "ja" startet den Bau. Sein 1-Monat-nach-archive-Vorschlag
   wurde begruendet abgeraten (Antwort vom 02.09. im Chat, Kern in der
   Ticket-Notiz).
3. **human: CD9TCB** — drei Optionen (durch 743737f erledigt erklären /
   Doku-Feld `overrides:` / Warnung bewusst reaktivieren).
4. **human: 7ZQ0ZN** — NOTES-0.1.1-Block ist auf master; Release = zwei
   Befehle: `git tag v0.1.1 && git push origin v0.1.1`.
5. **Schema-Spec-Review**: `docs/superpowers/specs/2026-08-31-schema-design.md`
   (`1e8fc57`) — Cut 1 startet erst auf sein Go. Offene Mini-Frage dazu:
   Entscheidung 10 nennt `verdict` nicht in der Reviewer-Reihenfolge; der
   Code zeigt es zwischen gaps und check.
6. **AXTFG3-Optik ist sein Call im Signoff**: eine Box kostet 2 Zeilen je
   Karte (nur bei Registry-Farbe); falls zu teuer, war der benannte
   Ausweich ein linker Farbbalken — dann neues Ticket.
7. **backlog, neu**: „lanes add setzt neben den Anker statt ans Ende" —
   Design-Frage (order.go ist Invarianten-Territorium), bewusst nicht gefixt;
   README sagt inzwischen ehrlich „Belongs, not lands".
8. **DoD-Klausel 3 von 88H1P4** ungetickt: frische Session + `jaira hook
   print`-Snippet in settings.json → arbeitet nach „go" bis signoff. Nur er
   kann das prüfen.

## Tags (01.09., Berks Wunsch) — komplett, in signoff

- **79GEPW** (Feld, `.jaira/tags`-Farb-Registry, `jaira tags`/`jaira tag`,
  Filter `tag:`, Doku): drei Commits, zwei adversariale Opus-Reviews mit
  eigenen Proben, 17 Funde gefixt (Concurrency-Verlust, /tags-Anker,
  Zaehl-Dedup, Suggest-Rueckzug …). Implementierer starb am Opus-Limit nach
  fertigen Edits; Koordinator hat verifiziert und committet.
- **AXTFG3** (Karten-Boxen in Tag-Farbe, `t`-Legende, Karten-Budget nach
  echter Hoehe statt /3): Sonnet-Implementierung, Koordinator-Review
  (Opus gedeckelt bis 19:00); ein Fund (Colour-Lookup unnormalisiert)
  koordinator-gefixt.
- **Zwei Koordinator-Pannen, beide fix-forward und dokumentiert:**
  84df2a4 ging mit nicht kompilierendem TESTfile auf master (set -e-Kette
  stoppte nicht; seitdem explizite || exit-Guards); am Vortag ein per
  tail -1 fehlgelesener Erfolg (BDV0HM-Fehlalarm, zurueckgezogen).
- Neu im Backlog: WriteAtomic fuer die uebrigen geteilten Board-Dateien
  (L19 des Tag-Reviews, v.a. core/lane/order.go).
- Der Agent-Block ist auf JEDEM Board stale bis `jaira update` (steht im
  0.1.1-NOTES-Block).

## Zurückgestellt, mit Grund

- **NFJCTK, B4MGTP, FCMP17**: stecken wörtlich im Schema-Spec (Cuts 1/4/6
  bzw. eigener Cut) — einzeln bauen hieße den abgestimmten Plan zerlegen.
- **88H1P4 Mechanismus c** (Block-Satz für Subagents): wartet auf sein Go.

## Arbeitsweise der Flotte (falls wiederholt)

Sequenziell je Ticket: Subagent implementiert und fährt das Board selbst
(claim → lanes → critique), Koordinator macht critique/optimize am selbst
gelesenen Diff, unabhängiges Opus-Review in Batches, Funde als Schleife
zurück. Lektionen (auch im Memory `jaira-code-learnings`): geteilter
git-Index verschluckt fremde staged Files → file-genau adden und notfalls
re-splitten; transiente FAILs ohne Testnamen = paralleler Schreiber
mid-compile; lipgloss polstert mehrzeilige Blöcke (styleLines).

## Ungeprüft / bekannt offen

- `human`-Lane weiterhin nicht exit-gated (nur signoff trägt
  requires-human-exit) — unverändert aus dem Morgen-Handoff, nicht ticketiert.
- Die grüne Lücke aus Berks Screenshot: Code kann sie nicht erzeugen
  (fitWindow-Repro leftover=0 auf jeder Breite); Diagnose: laufender Prozess
  mit veralteter Breite (Resize-Event verloren). Sein Wiggle-Test steht aus.
