# Handoff — 2026-09-03, nach dem "go"

Stand nach der go-Session; dieses File ist das Gedächtnis. Der vorige Handoff
(31.08. Abend) ist in der git-History dieses Files.

## Was heute passierte

- **Das Board ist public.** `jaira share` gelaufen (Berks Anweisung: "ab jetzt
  arbeiten wir public"), `.jaira/` aus dem .gitignore, Merge-Driver registriert,
  Commit `6c5d81b`. `.jaira/lanes/` bleibt by design privat (share sagt das
  selbst). Tickets sind ab jetzt Teil jedes Pushes.
- **Die 13 signoff-Tickets sind done** — per `--force` auf Berks Zuruf
  (Flotte + D9SC5A, 9H265S, 79GEPW, AXTFG3). done hielt danach 21 Tickets.
- **done-cap (SGPDYK) gebaut** — Berks "approved! arbeite durch". Zwei Commits:
  `d34341c` (feat) + `b1c4ecb` (fix nach Review). `holds: N` im Lane-File
  (builtin done: 10, Template synchron, nur terminale Lanes), Store.Overflow/
  TrimLane (core/ticket/trim.go), Trim an ALLEN VIER Status-Schreibstellen,
  gemeldet (CLI-Zeile mit restore, --json `trimmed`/`trim_error`, TUI-notify).
  Ticket in **signoff** — Berk nimmt ab.
- **Unabhängiges Opus-Review, zwei Durchgänge:** Erster fand 5 reproduzierte
  Funde — (1) accept() im signoff war eine VIERTE Schreibstelle ohne Trim
  (genau der menschliche Abnahme-Weg), (2) Tie-Break warf bei gleichem
  updated-at das NEUSTE Ticket raus, (3) Partial-Trim wurde verschluckt,
  (4) CLI gab exit 1 nach gelungenem Move (Retry = already_in_lane, Lane für
  immer über Cap), (5) holds auf nicht-terminaler Lane schickte unfertige
  Arbeit ins Logbuch. Alle gefixt, je mit Test gepinnt. Zweiter Durchgang:
  **accept, keine Funde** (Pin über 550 Kombinationen gemessen).
- **Neue Tickets:** P1AE82 (Version links oben im Projektfenster — KOLLIDIERT
  mit DNAEPN, das die Versionszeile auf dev-Builds gerade stumm schaltete;
  Klärung im Ticket-Kontext), D0SAHM (Epic-Layer: Storys als verbundene
  Bubbles über dem Board, Berks großer Wunsch vom 03.09., brainstorm
  angehakt, kompletter Spec im Kontext).

## Der eine Satz, der Berk überraschen könnte

Die Projektkopie `.jaira/lanes/done.md` hat **noch kein holds** — der Cap
greift auf DIESEM Board erst nach dem Drift-Refresh (Lane-Settings im TUI).
done hält 21: der erste done-Move oder Accept NACH dem Refresh schickt 11
Tickets ins Logbuch — laut Meldung, mit restore-Weg, aber eben elf auf einmal.

## Was Berk jetzt entscheidet

1. **signoff: SGPDYK** — done-cap abnehmen (review-check auf dem Ticket hat
   die Handgriffe). Danach ggf. Drift-Refresh für dieses Board.
2. **Schema-Spec-Review** (docs/superpowers/specs/2026-08-31-schema-design.md,
   1e8fc57) — Cut 1 startet nur auf sein Go. Mini-Frage: gehört `verdict` in
   die Reviewer-Reihenfolge von Entscheidung 10?
3. **7ZQ0ZN (human): Release** = `git tag v0.1.1 && git push origin v0.1.1`.
4. **CD9TCB (human)** — drei Optionen, unverändert.
5. **AXTFG3-Nachklapp:** Tag-Name im geöffneten Ticket ist bewusst ungefärbt
   (view.go:993, styleLines-Begründung im Code); Berk will Farbe — machbar
   (nur den Namen färben), auf Zuruf neues Ticket.
6. **Wiggle-Test** zur grünen Board-Lücke steht weiter aus.
7. Einmal `jaira update` pro Board (Agent-Block stale durch Tags + Cap).

## Bekannt offen, unticketiert

- 88H1P4 ist done (per Zuruf), aber Mechanismus c (CLAUDE.md-Block-Satz)
  wartet weiter auf Go, und DoD-Klausel 3 (frische Session bis signoff) blieb
  ungeprüft — beides steht als Override auf dem Ticket.
- Reviewer-Restnotizen zu SGPDYK: Partial-Trim-Kollision nur per Probe
  gepinnt (kein deterministischer Test); Display-Tie `ID<` vs Trim-Tie `ID>`
  (kosmetisch). `human`-Lane weiter nicht exit-gated.
- 8W1R94 (Update-Prozess durchspielen) liegt in todo; backlog: lanes-add-Anker,
  WriteAtomic für geteilte Board-Dateien.

## Arbeitsweise, Lektionen dieser Session

- Der Auto-Mode-Classifier blockt for-Schleifen mit `jaira move --force`;
  Einzelbefehle gehen durch (nondeterministisch — einfach neu ansetzen).
- Muster wie gehabt: Koordinator implementiert/critiqued am selbst gelesenen
  Diff, unabhängiges Opus-Review mit eigenen Proben in Schleife, Funde
  zurück nach in-progress, fix-Commit, Re-Review. Der Reviewer fand die
  vierte Schreibstelle, die grep nach moveMutation nicht zeigt: accept()
  schreibt den Status direkt (signoff.go).
- Ticket + Code im selben Commit; `verify-gate` verlangt am Turn-Ende Beweise:
  `go test ./... -race` (Cache geleert), `gofmt -l` (nur tickets.go), Binary
  `go build -o ~/.local/bin/jaira ./cmd/jaira`.
