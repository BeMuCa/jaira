---
id: 01M4RRVYXYBK44RANX2D0EW5XK
title: "Cockpit-Knöpfe farbig: Accept grün, Give back rot, Explain gelb"
status: done
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
commits:
  - 9e348d4145e37efb8dda6ad4fc2f3e3af3a70b79
  - eff028e0745caa86d73294c8d5f7245c23541e32
  - cbb8a77552b6b87f77afbedab0fa1640cf1fd1c9
  - d285852e9550ebbbb8248f892d967298a617f72e
  - 62288fb36300c4202d767fc2df5ca42f3a34f5a1
  - 530dbe19ef8f4cb9f194e47619db8bda2fd112d3
created-at: 2026-10-12T19:19:11Z
updated-at: 2026-10-12T20:46:15Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-81961
claimed-at: 2026-10-12T20:14:51Z
outcome-what: "Accept, Give back und Explain im Cockpit-Band zeichnen ihre Beschriftung als <Text> in den Theme-Farben success, error, warning; Read explanation und Verdict bleiben neutral"
outcome-why: "Die Knöpfe sahen gleich aus; man musste lesen, um den richtigen zu treffen"
outcome-resolves: "Accept grün, Give back rot, Explain gelb über Theme-Farben; Fokus bleibt durch die Inversion der Engine erkennbar; Release-Note 0EW5XK.md"
review-summary: "Im Cockpit-Band (core/mod/builtin/jaira-cockpit/hooks/register.tsx) tragen Accept, Give back und Explain ihre Beschriftung als <Text> in den Theme-Farben success/error/warning; Accept hat kein variant=primary mehr. Im Terminal sind die Knöpfe plain und zeichnen '[ Wort ]' samt Klammern im farbigen Text; Desktop/Mobile behalten den nativen Knopf. Neu in dieser Runde (62288fb3, 530dbe19): ein Druck auf einen der drei Knöpfe läuft durch once(): ein Token pro Ticket im State (busyAtom) sperrt die drei Knöpfe dieser Zeile, ihr Text wird dimColor, weitere Drücke tun nichts. Frei wird die Zeile, wenn accept/askWhatIsMissing/openChat zurückkehrt, spätestens nach 20 s; ein spätes Ende hebt nur seine eigene Sperre auf. session.start leert busyAtom, damit ein Reload keine Zeile dauerhaft sperrt. Release-Note hat eine zweite Zeile."
review-gaps: "1) Sichtbar inaktiv nur über dimColor am Kind-<Text>; auf Desktop/Mobile zeichnet die Engine den nativen Knopf, ob der das Text-dimColor zeigt, ist ungeprüft - Band ist aber ein Terminal-Band, gering. 2) Give back gibt die Sperre sofort nach ui.open frei (Pane offen = Rückmeldung) - passt zur Entscheidung 20:15, nur wird man das Dimmen dort kaum sehen. 3) Gezeichnetes Bild (Fokus-Inversion auf plain-Knopf, Abstände, wie stark dim auf farbigem Text wirkt) deckt kein Test. Kein Defekt gefunden: update() wiederholt den Updater bei Versionskonflikt, took nimmt den Wert des geschriebenen Laufs - korrekt (index.d.ts:14583). Untracked .jaira/milestones/ im worktree_diff gehört nicht zu diesem Ticket."
test-verdict: "pass: go build RC=0, go test -race -count=1 ./... RC=0, claude plugin test 93/93 RC=0 (3x stabil), validate passed; DoD 5 per Mutation geprüft (ohne took-Sperre, ohne Token-Vergleich, ohne 20-s-Timer, ohne Reset in session.start wird je ein Test rot); live nicht gesehen: wie gedimmt die Zeile im Terminal aussieht, Fokus-Inversion, Abstände"
review-verdict: "Diff 9e348d41..530dbe19 erfüllt DoD 1-5; claude plugin test 93/93 grün. Kein Defekt gefunden. Unsicher bleibt nur das Live-Bild: Fokus und Abstände der farbigen plain-Knöpfe und ob das Dimmen nach einem Klick deutlich genug zu sehen ist."
review-check: "1. Im Worktree /home/alex/projects/.worktrees/jaira-0ew5xk Claude Code im Terminal starten, so dass das Cockpit-Band ein wartendes Ticket mit Accept zeigt. 2. Ansehen: '[ Accept ]' ganz grün (Klammern inklusive), '[ Give back ]' ganz rot, '[ Explain ]' ganz gelb; Read explanation und Verdict normal. 3. ctrl+x tab ins Band, mit Tab über die drei Knöpfe: der fokussierte erscheint invertiert. 4. Auf Explain klicken: alle drei Knöpfe dieser Zeile werden sofort blass; ein zweiter Klick startet keinen zweiten Tab; die Zeile eines anderen Tickets bleibt kräftig. 5. Nach dem Start des Tabs (spätestens 20 s) sind die Knöpfe wieder kräftig. 6. Automatisch: cd core/mod/builtin/jaira-cockpit && claude plugin test . -> 93 pass."
---

# Cockpit-Knöpfe farbig: Accept grün, Give back rot, Explain gelb

## Definition of Done

- [x] Accept ist grün, Give back rot, Explain gelb im Band, über Theme-Farben der Engine; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:890-896; test 'waiting for you > Accept, Give back and Explain are told apart by the theme's colours; Read explanation and Verdict stay plain' (tests/board.test.ts:239)
- [x] Der fokussierte Knopf bleibt erkennbar (Fokus-Ring/Hervorhebung geht durch die Farbe nicht verloren)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:857,890-896: plain Buttons, no inverse set; ButtonProps.plain (2.1.296): 'the focus and the pointer still invert it'
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/0EW5XK.md:1-2
- [x] Die eckigen Klammern '[ ' und ' ]' der Knöpfe tragen dieselbe Farbe wie die Beschriftung (heute nur das Wort farbig, die Klammern neutral - wirkt abgeschnitten); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:857-858,890-896 (plain + '[ label ]' in the coloured Text on the terminal); test 'waiting for you > on the terminal the coloured buttons draw their brackets in the label's colour; elsewhere the native button stays' (tests/board.test.ts:259)
- [x] Nach einem Klick auf Accept, Give back oder Explain sind die Knoepfe im Band gesperrt (nicht klickbar, sichtbar inaktiv), bis die Aktion Wirkung zeigt oder eine Zeitgrenze ablaeuft - kein Doppelklick, und man sieht, dass etwas passiert; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:344-366 once() + 890-896 (dimColor from busy) + 693 (a reload frees the lock); tests 'waiting for you > a pressed row locks, drawn dim, until its action is done; a second press does nothing' (tests/board.test.ts:290), '... free again after 20 seconds, and the late end frees only its own lock' (:310), 'a reload frees a row whose press it cut off' (:331)

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
- **2026-10-12 20:05 · Alexander Sacharov** — critique Runde 3 (nach DoD 4): nur d285852e gelesen; plain/bracketed nur am Terminal passt zum bestehenden Muster, eine Oberfläche per Fähigkeit/Surface zu unterscheiden (hasInput); label= zurück ist begründet (Notiz 22:0x). Kein Befund. optimize: nichts zu entfernen.
- **2026-10-12 20:12 · Alexander Sacharov** — Alex 2026-10-12: versehentlich nach done geraten (vermutlich Accept-Klick beim Testen). Vom Teamlead zurueck nach review: frisches Review von d285852e (DoD 4, Klammern) bis zum Ende, Verdict fuer diese Runde neu schreiben - der vorhandene ist aus Runde 1.
- **2026-10-12 20:12 · Alexander Sacharov** — review (fresh, Runde 2): DoD 1-4 durch Diff 08d26f68..d285852e und Tests gedeckt, 90/90 grün. plain-Doku (index.d.ts:1105-1111) stützt die Fokus-Inversion auch ohne Chrome. Offen nur das Live-Bild: Fokus auf farbigem plain Knopf, Abstände. Das untracked .jaira/milestones/ im worktree_diff gehört nicht zu diesem Ticket.
- **2026-10-12 20:14 · Alexander Sacharov** — Alex 2026-10-12 (signoff): zurueckgegeben fuer DoD 5 - nach dem Klick kommt die Rueckmeldung nicht sofort, man sieht nicht, dass etwas passiert ist, und kann nochmal klicken. Der Knopf soll fuer eine Weile inaktiv werden. Vom Teamlead in dieses Ticket gefaltet (dieselben Knoepfe, dieselbe Datei register.tsx). Achtung: GPFA7B arbeitet parallel an accept()/afterAccept in derselben Datei.
- **2026-10-12 20:15 · Alexander Sacharov** — DoD 5, entschieden ohne Rückfrage: (1) gesperrt werden die drei Aktionsknöpfe (Accept, Give back, Explain) der Zeile des geklickten Tickets, nicht das ganze Band - andere Tickets bleiben bedienbar. (2) 'Wirkung zeigt' = die Aktion (accept/askWhatIsMissing/openChat) ist fertig: Accept hat das Board neu gelesen, Give back den Pane geöffnet, Explain den Tab gestartet oder einen Fehler gezeigt. (3) Zeitgrenze 20 s, danach wieder frei, auch wenn die Aktion noch läuft (spawn.sh darf bis 120 s). (4) Sichtbar inaktiv: dimColor auf dem farbigen Text; Button hat kein disabled, darum schluckt ein Wrapper den Klick. Die Sperre ist ein Token pro Ticket im State (busy), damit das späte Ende einer alten Aktion eine neue Sperre nicht aufhebt. accept()/afterAccept bleiben unberührt (GPFA7B).
- **2026-10-12 20:20 · Alexander Sacharov** — critique (DoD 5): Sperre im $.state überlebt einen Hot-Reload, der Timer von $.clock.after und das finally der laufenden Aktion nicht -> Zeile bleibt für die Session gesperrt. Fix: in session.start busyAtom auf {} setzen (reload feuert session.start erneut, laut plugin-authoring). Test dazu.
- **2026-10-12 20:22 · Alexander Sacharov** — critique Runde 2 (DoD 5): der Reload-Befund ist mit 530dbe19 erledigt, Test 'a reload frees a row whose press it cut off' (Mutation geprüft). Nichts weiter.
- **2026-10-12 20:30 · Alexander Sacharov** — review (fresh, Runde 3, DoD 5): 62288fb3+530dbe19 gelesen, once()-Sperre korrekt auch bei update()-Retry, Reload-Reset in session.start. 93/93 grün. Offen nur Live-Bild (Dimmen sichtbar genug?).
- **2026-10-12 20:43 · Alexander Sacharov** — Abnahme Runde 1 (Seite HFjLu3hRk4JWen6w9shguz): Farben, Abstaende, Sperre nach Explain, Ruecksetzen nach 20 s ok. Fokus (s02) uebersprungen: ctrl+x ist in Herdr die Praefix-Taste. Accept (s06) scheiterte am Aufbau der Probe (F16YHM ohne Commit), nicht an 0EW5XK - Alex: 'Accept F16YHM failed / lane done requires the commits that carry this change'. Beide Schritte neu: b1-s07 (Fokus-Hervorhebung per Maus-Hover), b1-s08 (Accept nach Probe-Commit 719cddc).
- **2026-10-12 20:45 · Alexander Sacharov** — Abnahme Runde 1 abgeschlossen: b1-s07 (Fokus-Hervorhebung per Maus-Hover) ok - Alex im Chat: 'мышь работет'; b1-s08 (Accept nach Probe-Commit) ok - Alex: 'принял', F16YHM in der Probe in done. Alle sechs Schritte ok, beide Maschinenchecks gruen -> angenommen.
