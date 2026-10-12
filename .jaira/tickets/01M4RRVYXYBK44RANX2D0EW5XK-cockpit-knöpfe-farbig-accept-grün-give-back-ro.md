---
id: 01M4RRVYXYBK44RANX2D0EW5XK
title: "Cockpit-Knöpfe farbig: Accept grün, Give back rot, Explain gelb"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Knöpfe im Cockpit-Band sind auf einen Blick unterscheidbar: Accept grün, Give back rot, Explain gelb"
context: |-
  Im Cockpit-Band sehen Accept, Give back und Explain gleich aus; man muss lesen, um den richtigen zu treffen.
  - Wunsch von Alex 2026-10-12: Accept grün, Give back rot, Explain gelb.
  - Code: Knöpfe in core/mod/builtin/jaira-cockpit/hooks/register.tsx (~Zeile 856-868, Branch feat/13pqb3). Accept hat heute variant="primary" (blau hinterlegt, siehe Screenshots), die anderen keine variant.
  - Unbekannt: ob <Button> der Mod-Engine eine Farbe annimmt (color-Prop oder weitere variants wie success/danger/warning). Erst die Engine-Typen prüfen; gibt es keine Farbe, die Beschriftung in einem farbigen <Text> rendern oder melden, was geht.
  - Die Farben sind Theme-Farben (success/error/warning wie bei 'threads:' in gelb), keine festen Hex-Werte, damit helles und dunkles Terminal-Theme funktionieren.
  - Read explanation und Verdict bleiben neutral.
definition-of-done: "Accept ist grün, Give back rot, Explain gelb im Band, über Theme-Farben der Engine; Test"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T19:19:11Z
updated-at: 2026-10-12T20:04:49Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-94221
claimed-at: 2026-10-12T20:03:11Z
outcome-what: "Accept, Give back und Explain im Cockpit-Band zeichnen ihre Beschriftung als <Text> in den Theme-Farben success, error, warning; Read explanation und Verdict bleiben neutral"
outcome-why: "Die Knöpfe sahen gleich aus; man musste lesen, um den richtigen zu treffen"
outcome-resolves: "Accept grün, Give back rot, Explain gelb über Theme-Farben; Fokus bleibt durch die Inversion der Engine erkennbar; Release-Note 0EW5XK.md"
review-summary: "Im Cockpit-Band (core/mod/builtin/jaira-cockpit/hooks/register.tsx:858-864) tragen Accept, Give back und Explain kein label= mehr, sondern ein <Text> als Kind in den Theme-Farben success (grün), error (rot), warning (gelb). Accept verliert variant=\"primary\", damit es nicht blau hinterlegt mit grüner Schrift erscheint. Read explanation und Verdict bleiben unverändert neutral. Ein Test (tests/board.test.ts:239) liest die Farbe des Text-Kinds pro Knopf-key und prüft, dass Accept keine variant hat. Release-Note core/release/notes/0EW5XK.md, eine Zeile."
review-gaps: "1) DoD 2 (Fokus erkennbar) ist nur aus der Engine-Doku belegt, nie im Terminal gesehen; ob die Inversion der Engine auch ein farbiges Text-Kind invertiert, zeigt kein Test. 2) cbb8a775 (optimize) hat label= an den drei Knöpfen entfernt, obwohl die Notiz 19:23 label als Adresse/Name des Knopfs behalten wollte; die Begründung 'Text-Kind benennt den Knopf, wenn label fehlt' konnte ich nicht an claude-code.d.ts prüfen (Datei hier nicht auffindbar). Tests laufen grün, aber kein Test prüft den Namen des Knopfs. 3) Der Test prüft Props, nicht die gezeichnete Farbe. Nicht Teil des Tickets: worktree enthält unversionierte .jaira/milestones/demo-*.md — nicht mitcommitten."
test-verdict: "pass: go build RC=0, go test -race -count=1 ./... RC=0, claude plugin test 89/89 RC=0, validate passed; DoD 1 im Baum und per Mutation geprüft (success->text lässt den Test rot werden); DoD 3 Datei vorhanden; DoD 2 nur aus der Engine-Doku (Fokus = inverse), live nicht gesehen — Handprobe in signoff"
review-verdict: "Diff erfüllt DoD 1 und 3; 89/89 Tests in claude plugin test grün. DoD 2 (Fokus) ist plausibel, aber unbewiesen — bitte einmal live prüfen (Schritt 4-5 im Check). Keine Defekte gefunden."
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test  → 89 pass, 0 fail  2. Claude Code in einem Repo mit jaira-Board starten, in dem ein Ticket in signoff/human wartet  3. Cockpit-Band ansehen: Accept grün, Give back rot, Explain gelb (Explain nur bei erklärbarem Ticket); Read explanation und Verdict in normaler Farbe; Accept NICHT blau hinterlegt  4. ctrl+x tab ins Band, mit Tab/Pfeil auf Accept → der Knopf erscheint invertiert, klar als fokussiert erkennbar  5. dasselbe für Give back und Explain  6. optional: helles Terminal-Theme → alle drei Farben lesbar"
---

# Cockpit-Knöpfe farbig: Accept grün, Give back rot, Explain gelb

## Definition of Done

- [x] Accept ist grün, Give back rot, Explain gelb im Band, über Theme-Farben der Engine; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:862-868; test 'waiting for you > Accept, Give back and Explain are told apart by the theme's colours; Read explanation and Verdict stay plain' (tests/board.test.ts:239)
- [x] Der fokussierte Knopf bleibt erkennbar (Fokus-Ring/Hervorhebung geht durch die Farbe nicht verloren)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:829,862-868: plain Buttons, no inverse set; ButtonProps.plain (2.1.296): 'the focus and the pointer still invert it'
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/0EW5XK.md:1
- [x] Die eckigen Klammern '[ ' und ' ]' der Knöpfe tragen dieselbe Farbe wie die Beschriftung (heute nur das Wort farbig, die Klammern neutral - wirkt abgeschnitten); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:829-830,862-868 (plain + '[ label ]' in the coloured Text on the terminal); test 'waiting for you > on the terminal the coloured buttons draw their brackets in the label's colour; elsewhere the native button stays' (tests/board.test.ts:259)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 19:19 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/0ew5xk ist von feat/13pqb3 abgezweigt (Kette QRHB35 -> TWB0R6 -> 13PQB3): 13PQB3 hat die Knopfreihe (Verdict-Knopf, Explain-Regel) umgebaut. Nur die Knöpfe einfärben; GPFA7B arbeitet parallel an accept()/afterAccept in derselben Datei.
- **2026-10-12 19:19 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 19:23 · Alexander Sacharov** — Keine offenen Entscheidungen. Engine 2.1.296: ButtonProps hat keine color-Prop, variant nur 'primary'|'secondary'. Weg: Button behält label (Adresse/Name) und bekommt ein <Text color=success|error|warning> als Kind — erlaubt laut ButtonProps.label-Doku. Fokus: die Engine markiert Fokus und Pointer per inverse, unabhängig von der Textfarbe; wir setzen kein plain/inverse. Accept behält variant=primary (Desktop zeichnet so den Hauptknopf; im Terminal bleiben die Klammern in Akzentfarbe, die Beschriftung grün).
- **2026-10-12 19:25 · Alexander Sacharov** — DoD 2 ist aus der Engine-Doku belegt, nicht live gesehen: die Test-API zeigt keinen Fokus-Zustand. Im Terminal prüfen: ctrl+x tab ins Band, Fokus auf Accept — der Knopf muss invertiert erscheinen.
- **2026-10-12 19:25 · Alexander Sacharov** — critique: register.tsx:858 variant="primary" von Accept streichen. Primary zeichnet den Knopf in der Akzentfarbe (laut Kontext blau hinterlegt); grüne Schrift darin ist weder 'grün' noch gut lesbar. Die Farbe übernimmt jetzt, was primary markierte. Meine frühere Notiz (primary behalten wegen Desktop) ist damit überholt: das Band ist ein Terminal-Band, der Desktop zeichnet ohnehin seinen nativen Knopf.
- **2026-10-12 19:26 · Alexander Sacharov** — critique Runde 2: der einzige Befund (variant=primary an Accept) ist mit eff028e0 erledigt, Test sichert es. Nichts weiter.
- **2026-10-12 19:34 · Alexander Sacharov** — review (fresh): DoD 1+3 durch Diff und Test gedeckt, claude plugin test 89/89 grün. Offen für den Menschen: Fokus-Inversion bei farbigem Text-Kind live prüfen; label= wurde in cbb8a775 entfernt, Begründung aus ButtonProps-Doku hier nicht nachprüfbar.
- **2026-10-12 19:34 · Alexander Sacharov** — Zu review-gap 2: geprüft in claude-code.d.ts (Engine 2.1.296, ButtonProps.label): 'Children of strings and Text ... are drawn in its place ... The label then names the control; absent, it is the children's text.' Die Adresse kommt aus key=, das bleibt gesetzt. Die Notiz 19:23 ist damit durch optimize überholt, nicht verletzt.
- **2026-10-12 19:40 · Alexander Sacharov** — explain: Befund: Keine Defekte; die drei Knöpfe tragen ihre Farbe, 89/89 Tests grün. Offen ist nur, ob der fokussierte farbige Knopf im Terminal wirklich invertiert erscheint (DoD 2) – das stützt sich auf die Engine-Doku, gesehen hat es niemand.
Warum: Ist der Fokus nicht erkennbar, trifft man per Tastatur leicht den falschen Knopf – ausgerechnet zwischen Accept und Give back. Die Tests prüfen nur Props, nicht was gezeichnet wird, und die Test-API zeigt keinen Fokus.
a) Einmal live prüfen (Check-Schritte 4-5: ctrl+x tab, Fokus auf Accept/Give back/Explain), dann annehmen
b) Ohne Live-Prüfung annehmen und der Engine-Doku vertrauen
c) Zurück nach in-progress, um den Fokus zusätzlich selbst zu markieren (z. B. Farbe bei Fokus abschalten)
Empfehlung: a - zwei Minuten im Terminal klären die einzige offene Frage; review-gap 2 (label) ist laut Notiz 19:34 schon an claude-code.d.ts belegt.
- **2026-10-12 19:45 · Alexander Sacharov** — decision: a) Einmal live prüfen (Check-Schritte 4-5: ctrl+x tab, Fokus auf Accept/Give back/Explain), dann annehmen
- **2026-10-12 20:02 · Alexander Sacharov** — Alex 2026-10-12: zurückgegeben - die [ ] müssen auch farbig sein, sonst sieht das Wort falsch/abgeschnitten aus. DoD 4. Status-Tippfehler 'in-progres' (19:46:49, Herkunft unklar) per Ref-Commit e3a3caff auf in-progress korrigiert.
- **2026-10-12 20:03 · Alexander Sacharov** — DoD 4, Weg: die Klammern zeichnet die Engine selbst um ein Button-Kind, ohne Farbprop (ButtonProps 2.1.296). Darum auf dem Terminal plain (keine Chrome) und die Klammern '[ ' ' ]' im selben farbigen <Text>; laut Doku invertieren Fokus und Pointer auch einen plain Button, DoD 2 bleibt. Nur bei e.surface === 'terminal': Desktop/Mobile zeichnen ihren nativen Knopf und bekämen sonst die Klammern als Text hinein.
- **2026-10-12 20:04 · Alexander Sacharov** — Runde 2 offen fürs Live-Bild: die Test-API zeigt keine Malerei. Ob zwischen einem plain '[ Give back ]' und dem nächsten Engine-Knopf derselbe Abstand bleibt wie zwischen zwei Engine-Knöpfen, ist nur live zu sehen. label= ist wieder gesetzt: mit Klammern im Text wäre der Name des Knopfs sonst '[ Accept ]'.
