---
id: 01M4V25R2BDN0HVWF879NA413Y
title: Das Board laedt nur geaenderte Tickets neu und parst jedes nur einmal
status: in-progress
ready: true
creator: BeMuCa
goal: "Ein offenes Board braucht im Leerlauf kaum CPU und friert nicht mehr ein, weil es nur liest und dekodiert, was sich geaendert hat."
context: |-
  Heute braucht ein offenes Board auf diesem Repo 50-80 % eines CPU-Kerns im Leerlauf und friert bis 3,9 s ein (RAM-Analyse vom 2026-10-13, https://claude.ai/artifact/FrBA4U7LDwQQZhnDh9hf1E).
  Warum:
  - Alle 2 s laedt es alles neu (internal/tui/model.go:973 tick, :375 reload): Lanes, Tags, Milestones, jedes Ticket samt Refs, Logbuch, Sitzungen.
  - Decode (core/ticket/schema.go:545) liest ~34 Felder einzeln; jeder Zugriff parst den ganzen YAML-Kopf neu (core/ticket/frontmatter.go:172). Also ~34 Parses pro Ticket und Takt.
  - Der Datei-Waechter (model.go:854) erfaehrt, welche Datei sich geaendert hat, wirft den Namen weg (case _, ok := <-w.Events) und laedt alles.
  Gemessen in der Analyse: 'einmal parsen' allein CPU 64 -> 9 %; 'nur Geaendertes dekodieren' CPU 66 -> 6 %. RAM sinkt nur wenige MB.
  Berk am 2026-10-13: Refs pruefen, das Ganze einmal parsen, gezielt nach 'zuletzt angefasst' schauen; 'CPU-Gewinn passt mir'.
  Grenzen, die bekannt sind: Ein Ref-Hash sieht keine Handaenderung und kein git checkout; ein Fingerabdruck ueber ALLE Eingaben war in der Analyse fehleranfaellig (Prototyp vergass HEAD, Branch-Ref, ~/.jaira/lanes, Datum). Darum pro Ticket pruefen, nicht global.
definition-of-done: "Decode parst den YAML-Kopf eines Tickets genau einmal; gleiche Ausgabe wie vorher, belegt durch einen Test ueber vorhandene Ticket-Dateien"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T16:40:17Z
updated-at: 2026-10-13T21:14:30Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3239170
claimed-at: 2026-10-13T21:01:47Z
---

# Das Board laedt nur geaenderte Tickets neu und parst jedes nur einmal

## Definition of Done

- [x] Decode parst den YAML-Kopf eines Tickets genau einmal; gleiche Ausgabe wie vorher, belegt durch einen Test ueber vorhandene Ticket-Dateien
  proof: core/ticket/schema.go Decode/decodeFields, frontmatter.go Fields; TestDecodeParsesOnceAndMatchesBefore, TestDecodeMatchesBeforeOnEdgeCases
- [ ] Ein Ticket, dessen Datei (Zeit und Groesse) und Ref-Hash sich seit dem letzten Laden nicht geaendert haben, wird nicht neu gelesen und dekodiert; eine Aenderung gleicher Groesse wird trotzdem gefunden; Tests decken beide Faelle
- [ ] Der Datei-Waechter gibt die geaenderten Dateien weiter, und nur diese werden neu gelesen; das volle Neuladen laeuft nur noch als seltenes Sicherheitsnetz
- [ ] Gemessen auf diesem Repo: CPU des offenen Boards im Leerlauf vorher und nachher, Zahlen als Notiz am Ticket; go test ./... -race gruen
- [ ] CPU-Test: ein Test belegt, dass ein Takt ohne Aenderung kein Ticket dekodiert und keinen YAML-Kopf parst (gezaehlt, nicht in CPU-Prozent, damit er im CI stabil ist), und dass eine Aenderung trotzdem im naechsten Takt erscheint

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] A (DoD 1): Test zuerst mit decodeBefore als Kopie des heutigen Decode ueber .jaira/tickets + logbook + Randfaelle; dann Fields/ticketFields/decodeFields + Counts(); Decode parst einmal; ein gruener Commit
- [x] B: lane.parse liest ueber ein einziges d.Fields() (str, list, Has, rejects-to); go test ./core/lane gruen; ein Commit
- [x] C: Ticket.Clone() mit eigenen Slices und Doc-Kopie + Test, dass Aendern der Kopie das Original nicht beruehrt; ein Commit
- [~] D (DoD 2, Dateien): Kopf-Cache im Store (EnableCache, cached, racyWindow 2 s, Forget/ForgetAll, loadHeader + idOf, Clone pro Treffer); Tests: unveraendert = 0 Lesen/0 Decode, gleich grosse Aenderung mit gleicher mtime wird gefunden, ForgetAll findet sie auch ausserhalb des Fensters, Kopie aendern laesst Cache heil, geloeschte Datei verschwindet, Cache aus = alter Weg; ein Commit
- [ ] E (DoD 2, Refs): gitref.ListHashes hinter ReadMany + Hash-Cache in refsync.Extra; Test mit fileOnRef: zweites Extra ohne Decode, geaenderter Ref erscheint; ein Commit
- [ ] F1 (DoD 5): Board schaltet den Cache ein (New, switchBoard); Test in internal/tui: Takt ohne Aenderung laesst Counts unveraendert, Titel per store.Mutate geaendert -> naechster Takt zeigt ihn mit genau 1 Decode; ein Commit
- [ ] F2 (DoD 3, Waechter): chan []string, changeMsg{paths}, Forget vor reload; Tests: changeMsg mit einem Pfad liest nur diesen neu (1 Decode, auch bei gleicher Groesse+mtime), echtes fsnotify-Ereignis kommt mit Pfad an (5 s Timeout, Skip ohne Watcher); ein Commit
- [ ] F3 (DoD 3, Sicherheitsnetz): fullReparseEvery 15 min + m.fullAt + ForgetAll im Takt, Kommentare an reload/tick umschreiben, core/release/notes/NA413Y.md im selben Commit; Test: unsichtbare Aenderung erscheint erst, wenn fullAt aelter als 15 min ist
- [ ] G (DoD 4): CPU im Leerlauf vorher (82b2252a) und nachher auf einer Kopie dieses Repos ohne Remote messen, Rezept in der Bauplan-Notiz; Zahlen + go test ./... -race als Notiz; kein Commit

## Progress
- **2026-10-13 16:56 · BeMuCa** — Berk 2026-10-13: will einen CPU-Test (DoD 5). Seine Idee zu den Intervallen: Refs alle 2 s pruefen, Zeit+Groesse nur alle paar Minuten (15?), weil Refs fuer alles KI-Getriebene reichen. Einschaetzung Claude: Refs fangen jeden Schreibvorgang ueber jaira (Agenten und Menschen ueber CLI/TUI), aber keine Handaenderung an Ticket-Dateien, kein git checkout/pull/merge und keine Dateien ohne Ref (lanes, tags, settings). Gemessen 2026-10-13 auf diesem Repo: stat aller 34 Ticket-Dateien 0,04 ms, git for-each-ref refs/jaira (77 Refs) 1,5 ms. Beides ist alle 2 s praktisch gratis - teuer ist nur das Parsen. Vorschlag: Refs + Zeit/Groesse alle 2 s, geparst wird nur, was sich geaendert hat; das volle Neu-Parsen als Sicherheitsnetz alle ~15 min (dort passen Berks 15). Der Datei-Waechter meldet Dateiaenderungen ohnehin sofort. Endgueltig beim Start nach PR #74.
- **2026-10-13 21:03 · BeMuCa** — Intervall entschieden (pre-process 2026-10-13):
- Alle 2 s: stat jeder Ticket-Datei und des Logbuch-Fensters, dazu git for-each-ref mit Hashes. Geparst wird nur, was sich geaendert hat.
- Alle 15 min: volles Neu-Parsen (Kopf-Cache leeren). Dort passen Berks 15.
- Dazwischen meldet der Datei-Waechter sofort und nennt die Dateien.
Warum stat nicht nur alle 15 min (Berks Idee):
- Refs sehen keine Handaenderung und kein git checkout/pull/merge. Nur stat sieht das.
- Auf WSL2-Laufwerken unter /mnt meldet fsnotify nichts. Genau dafuer gibt es den 2-s-Takt (Kommentar internal/tui/model.go ~:970). Mit stat alle 15 min waere ein git pull dort bis zu 15 min unsichtbar.
- stat aller 34 Dateien kostet 0,04 ms (Messung in der Notiz oben). Da gibt es nichts zu sparen; teuer war nur das Parsen.
Warum das volle Neu-Parsen alle 15 min reicht:
- Es faengt nur, was stat nicht sieht: gleiche Groesse UND gleiche mtime ausserhalb des 2-s-Fensters (cp -p, tar, rsync -t, Uhrversatz auf Netzlaufwerken). Selten; 15 min Verzug ist vertretbar.
- Nach "einmal parsen" kostet es einen Parse pro Ticket, einmal in 15 min.
- Es raeumt nebenbei Cache-Eintraege fuer Pfade weg, die es nicht mehr gibt (umbenannt, ins Logbuch verschoben).
Der Ref-Hash braucht kein Sicherheitsnetz: gleicher Commit-Hash heisst gleicher Inhalt.
- **2026-10-13 21:03 · BeMuCa** — Bauplan fuer die Schritte im ## Plan (pre-process 2026-10-13). in-progress muss hier nichts mehr entscheiden. Buchstaben = Plan-Schritte.

A. Einmal parsen (core/ticket/frontmatter.go, schema.go)
- Neuer Typ Fields = ein Parse des Kopfes: root ast.Node (f.Docs[0].Body, nil bei leerem Kopf) und vals []*ast.MappingValueNode (nil, wenn root keine Mapping ist; MappingValueNode-Wurzel = ein Element wie in topLevel).
- func (d *Doc) Fields() (*Fields, error): der einzige Lese-Aufruf von parser.ParseBytes([]byte(d.fm), parser.ParseComments). Fehlertext wortgleich: "ticket: frontmatter is not valid YAML: %w".
- *Fields bekommt Has, Scalar, List. Die Rumpfe von Doc.Scalar/Doc.List hinter dem find wandern in Helfer scalarOf(key, node) / listOf(key, node), die Doc und Fields beide nutzen. Erster Treffer gewinnt, wie Doc.find. Ein nil *Fields verhaelt sich wie ein leerer Kopf.
- f.validate(): Wurzeltyp-Pruefung und Anker/Alias/Merge-Walk aus Doc.Validate, gleiche Meldungen. Doc.Validate = Fields() + validate(). Doc.topLevel = Fields() + die bisherige Meldung "root is %T". Die Setter bleiben unveraendert.
- Decode(d, path) = d.ticketFields() + decodeFields(d, f, path). decodeFields prueft zuerst (gleiche Fehlerreihenfolge wie heute), liest dann alle str/list ueber f statt ueber d.
- Zaehler (sync/atomic): decodes in decodeFields, headerParses in d.ticketFields() (= Fields-Aufruf fuer einen Ticket-Kopf: Decode, Cache, idOf). Lesbar ueber func Counts() (decodes, headerParses int64). Lane-Dateien zaehlen nicht: sie sind kein Ticket-Kopf (so lese ich "keinen YAML-Kopf" in DoD 5).
- Test zuerst: decodeBefore = woertliche Kopie des heutigen Decode (Validate + d.Scalar/d.List) im Testfile, angelegt VOR dem Umbau. Ueber ../../.jaira/tickets/*.md und ../../.jaira/logbook/*/*.md (weniger als 1 Datei = Fehler): reflect.DeepEqual(Decode, decodeBefore), Fehlertexte gleich, je Datei Counts-Differenz decodes 1 und headerParses 1. Dazu Randfaelle: kaputtes YAML, Anker, Wurzel keine Mapping, ein einzelner Schluessel, Block-Skalar |, null-Wert, Liste wo Skalar erwartet und umgekehrt, fehlende id, fehlender status.

B. Lanes (core/lane/lane.go parse ~:278): einmal f, _ := d.Fields(); str, list, das Has bei ~:301 und die rejects-to-Lesung bei ~:406 lesen ueber f. Heute parst jeder dieser Aufrufe die Lane-Datei neu, und lane.Load laeuft jeden Takt.

C. Ticket.Clone() (schema.go): Struct-Kopie mit eigenen Kopien von BlockedBy, Commits, Tags, Related, DoDItems, PlanItems, Options; doc als Kopie des Doc-Structs (drei Strings). Grund: der Cache gibt nie denselben *Ticket zweimal aus. Aufrufer setzen Felder: store.extra setzt ReadOnly, applyMove (model.go ~1854) setzt Assignee.

D. Kopf-Cache im Store (core/ticket/store.go)
- Aus, bis jemand func (s *Store) EnableCache() ruft. Die CLI ruft es nie und bleibt wie heute; nur das Board schaltet ein.
- Feld cache *headerCache: eigener sync.Mutex + map[path]*cacheEntry. Mutex, weil fetchRefs in einer Goroutine Store.Load -> resolve -> idIndex -> idOf ruft, waehrend Update List laufen laesst. Zeiger-Feld, damit ein kopierter Store keinen Mutex kopiert.
- cacheEntry: mod time.Time, size int64, seen time.Time, raw []byte, t *Ticket, err error, id string, idErr error.
- s.cached(path): os.Stat. Eintrag mit gleicher mod und size: wenn seen.Sub(mod) >= racyWindow -> Treffer ohne Lesen. Sonst Datei lesen und mit raw vergleichen: gleich -> seen = Zeit vor dem Lesen, Treffer; ungleich -> neu bauen. Neu bauen: Stat, seen := time.Now(), ReadFile, ParseDoc, f := d.ticketFields(), id/idErr wie heutiges idOf (NormalizeIDPrefix des Scalar(id)), t/err = decodeFields. Stat-Fehler -> Eintrag loeschen, Fehler zurueck. Unexportierter Zaehler fileReads beim Lesen, nur fuer Tests.
- racyWindow = 2 * time.Second. Grund: FAT/exFAT speichern mtime in 2-s-Schritten, Linux stempelt mit grober Uhr; in diesem Fenster kann ein gleich grosses Neuschreiben dieselbe mtime behalten. Das ist die "racy git"-Regel.
- loadHeader (List und LoggedSince): Cache an -> cached(path), dann e.err oder e.t.Clone(). idOf: Cache an -> e.id, e.idErr. Cache aus -> heutiger Code unveraendert.
- Forget(paths ...string) loescht diese Eintraege, ForgetAll() leert alles. Bei Cache aus ohne Wirkung.
- Store.Load bleibt ungecacht: das ist der Mutationspfad.

E. Ref-only-Tickets (core/gitref/gitref.go, core/refsync/refsync.go Extra ~:1081)
- gitref: func (r *Repo) ListHashes() (map[string]string, error) = for-each-ref --format="%(objectname) %(refname)" Prefix; Schluessel = Name ohne Prefix, Wert = Commit-Hash. Direkt hinter ReadMany einfuegen, nicht neben List/ListMilestones: dort arbeitet der Milestone-Zweig parallel.
- Extra nimmt ListHashes statt List. Cache im Syncer: Mutex + map[id]{hash, t}. Gewollte id mit gleichem Hash -> t.Clone() (t nil = kaputter Ref, ueberspringen); nur die uebrigen gehen an ReadMany und Decode. Nicht mehr gelistete ids fliegen raus. Kein Sicherheitsnetz noetig: Hash ist inhaltsadressiert.
- Test mit dem Helfer fileOnRef (refsync_test.go ~:551).

F. Board (internal/tui)
- New und switchBoard: s.EnableCache() vor dem ersten reload; m.fullAt = time.Now().
- Waechter (model.go startWatching ~:853): m.watch wird chan []string. Die Goroutine sammelt ev.Name in einer Menge. Beim Debounce nicht blockierend senden; ist der Kanal voll, Menge behalten und Debounce neu stellen. waitForChange liefert changeMsg{paths}. Handler: m.store.Forget(msg.paths...), dann reload. Session-Pfade stehen nicht im Cache; Forget uebergeht sie.
- Takt (tickMsg ~:992): wenn time.Since(m.fullAt) >= fullReparseEvery -> m.store.ForgetAll(), m.fullAt = time.Now(); dann reload. const fullReparseEvery = 15 * time.Minute, Kommentar mit dem Grund aus der Intervall-Notiz.
- Kommentare an reload (~:369-374, "A full rescan ... second source of truth") und tick (~:970) auf das neue Verhalten umschreiben. Berk hat diese Abwaegung am 2026-10-13 umgedreht ("CPU-Gewinn passt mir").
- Tests ohne t.Parallel: Counts ist paketweit.
- Release-Note core/release/notes/NA413Y.md, eine Zeile, im Commit von F3:
  - Leave the board open as long as you like: while nothing changes it rereads no ticket, a changed ticket still shows within 2 seconds, and an edit that keeps both a file's size and its timestamp shows within 15 minutes.

G. Messung DoD 4 (kein Commit). Nie jaira gegen /home/berk/git/* laufen lassen, nur gegen die Kopie.
- S=<scratchpad>; cp -a /home/berk/git/jAIra $S/cpu-board; jedes Remote der Kopie entfernen (git -C $S/cpu-board remote, dann remote remove). Grund: das Board schiebt Refs.
- cp -a ~/.jaira $S/jhome und JAIRA_HOME=$S/jhome setzen, damit die Projektliste in ~/.jaira unberuehrt bleibt.
- vorher: git -C /home/berk/git/jAIra-NA413Y archive 82b2252a | tar -x -C $S/base; dort go build -o $S/jaira-before ./cmd/jaira. nachher: go build -o $S/jaira-after ./cmd/jaira im Worktree. Ohne -race: der Race-Detektor verfaelscht die CPU.
- je Binary: tmux new-session -d -s cpu -x 200 -y 50 "cd $S/cpu-board && JAIRA_HOME=$S/jhome <binary>"; 10 s warten; pid = tmux list-panes -t cpu -F '#{pane_pid}'; Felder 14+15 aus /proc/<pid>/stat jetzt und nach 60 s; CPU% = Differenz / (getconf CLK_TCK * 60) * 100; tmux kill-session -t cpu.
- Notiz: beide Zahlen, Ticket-Zahl der Kopie, Befehl; dazu go test ./... -race (Ausgabe-Ende woertlich).
