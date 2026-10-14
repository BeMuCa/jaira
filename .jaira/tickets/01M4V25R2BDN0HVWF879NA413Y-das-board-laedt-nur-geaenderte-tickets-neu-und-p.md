---
id: 01M4V25R2BDN0HVWF879NA413Y
title: Das Board laedt nur geaenderte Tickets neu und parst jedes nur einmal
status: signoff
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
updated-at: 2026-10-13T22:35:10Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3239170
claimed-at: 2026-10-13T21:01:47Z
outcome-what: "Decode parses a ticket head once (Fields); lane files likewise; Ticket.Clone; Store.EnableCache header cache keyed by file time+size with a 2 s racy-window byte compare and Forget/ForgetAll; ref-only tickets cached by ref commit hash (gitref.ListHashes); the board enables the cache, the watcher passes changed paths, and a 15 min full re-decode is the safety net; release note core/release/notes/NA413Y.md. Round 2: Extra builds its result in one pass over the sorted ids so the order no longer depends on which refs moved; Doc.find reuses Fields.find. Round 3: Extra drops a cached ref-only ticket whose ref no longer has its file"
outcome-why: "An idle board used 50-80 % of a CPU core because every 2 s tick re-read and re-decoded every ticket, parsing the YAML head about 34 times each"
outcome-resolves: "DoD1: TestDecodeParsesOnceAndMatchesBefore compares Decode with a verbatim copy of the old one over all repo tickets and counts 1 decode + 1 parse each. DoD2: cache.go skips unchanged time+size (TestCacheUnchangedReadsAndDecodesNothing), same-size change found by racy-window compare/Forget/ForgetAll, ref hash cache TestExtraSkipsUnchangedRefs. DoD3: watcher sends paths, only they are Forgotten (TestChangeMsgRereadsOnlyTheNamedFile); full reparse every 15 min (TestFullReparseFindsInvisibleChange). DoD4: idle CPU 61.9/64.0 % before to 1.7/1.4 % after on a 34-ticket copy, go test ./... -race green. DoD5: TestIdleTickDecodesNothingAndChangeAppearsNextTick counts 0 decodes/parses on idle ticks and shows a change next tick. Round 2: critique findings fixed, TestExtraKeepsOrderWhenOneRefMoved pins the Extra order. Round 3: TestExtraDropsCachedTicketWhoseRefLostItsFile pins that a ref without its file is dropped"
review-summary: |-
  Ein offenes Board dekodiert im Leerlauf kein Ticket mehr.
  1. Einmal parsen: Decode liest den YAML-Kopf eines Tickets mit einem Parse (neuer Typ Fields in core/ticket/frontmatter.go) statt mit ~34 Parses, einen pro Feld. Lane-Dateien (core/lane/lane.go) ebenso.
  2. Kopf-Cache fuer Ticket-Dateien (core/ticket/cache.go): pro Pfad gemerkt werden Aenderungszeit, Groesse, Bytes und das dekodierte Ticket. Gleiche Zeit und Groesse heisst: nicht lesen, nicht dekodieren. Ist die Datei beim Lesen juenger als 2 s, werden beim naechsten Mal die Bytes verglichen (die "racy git"-Regel). Wer liest, bekommt immer eine Kopie (Ticket.Clone).
  3. Nur das Board schaltet den Cache ein (internal/tui/model.go New und switchBoard). Die CLI ruft EnableCache nie und liest wie vorher.
  4. Ref-only-Tickets (nur auf refs/jaira/tickets, ohne Datei) merkt sich der Syncer mit dem Commit-Hash ihres Refs (gitref.ListHashes). Gleicher Hash: kein Lesen, kein Decode. Ein Ref ohne Ticket-Datei faellt raus.
  5. Der Datei-Waechter gibt die Pfade weiter, die er gesehen hat; nur die werden vergessen und neu gelesen.
  6. Sicherheitsnetz: alle 15 min leert der 2-s-Takt den Cache ganz. Das faengt Aenderungen, die Groesse und Zeitstempel behalten (cp -p, rsync -t), wo kein Waechter laeuft (WSL-Laufwerke unter /mnt).
  Ergebnis: Leerlauf-CPU auf einer 34-Ticket-Kopie dieses Repos von ~56-64 % auf ~1-2 % eines Kerns. Release-Note: core/release/notes/NA413Y.md.
round: 3
review-gaps: |-
  Nichts, was den Merge blockiert. Schwaecher als es klingt oder offen:
  1. DoD 1, der Beleg: decodeBefore (core/ticket/decode_once_test.go) ruft Doc.Validate/Scalar/List auf, und genau die baut dieser Diff mit um. Der Test vergleicht also nicht gegen den alten Parser-Code; ein Fehler in diesen geteilten Helfern bliebe unentdeckt. Ich habe es darum selbst gegengeprueft: Decode plus Doc.Scalar/List/Has/Keys/Validate von 82b2252a und von HEAD ueber 146 echte Ticket-Dateien (Tickets + Logbuch) und 12 Randfaelle: Ausgabe byte-gleich. lane.Load auf den 13 Lanes dieses Boards: ebenfalls gleich.
  2. Windows: TestExtraDropsCachedTicketWhoseRefLostItsFile (core/refsync/refsync_test.go:973) ruft `git hash-object -t tree /dev/null`. Auf Windows konnte ich es nicht laufen lassen. Git-Quelltext spricht dafuer, dass es geht: compat/mingw.c mingw_open macht aus "/dev/null" "nul", hash-object oeffnet per xopen und liest eine Nicht-Datei ueber index_pipe. Sicher ist es erst nach dem ersten windows-latest-Lauf der CI. Falls rot: der run-Helfer hat kein stdin, also entweder einen Helfer mit stdin fuer `git hash-object -t tree --stdin`, oder den festen Leerbaum-Hash 4b825dc642cb6eb9a060e54bf8d69288fbee4904 (gilt nur fuer SHA-1-Repos).
  3. Release-Note: "while nothing changes it rereads no ticket" stimmt nicht ganz; alle 15 min liest das Board jedes Ticket neu. Kleinigkeit, aber Text, den Nutzer lesen.
  4. CLI: Der Kopf-Cache bleibt fuer die CLI aus (EnableCache nur in internal/tui/model.go:326 und :359; tui.New nur aus internal/cli/board.go, home.go und scripts/shotgen). Der Ref-Cache im Syncer ist dagegen immer an, auch in der CLI. Harmlos, weil inhaltsadressiert: gleicher Hash = gleicher Inhalt.
  5. Vorbestehend, NICHT von diesem Diff: mein Stresstest (List, Load, Forget, ForgetAll parallel, -race) meldet einen Data Race auf Store.dupIDs in idIndex (core/ticket/store.go:1217 gegen :835). Auf 82b2252a derselbe Race ohne Cache. Der Hintergrund-Fetch des Boards ruft Store.Load (core/refsync/refsync.go:883, :1215, :1280) und kann ihn treffen. Eigenes Ticket wert, nicht dieses.
test-verdict: "pass: suite green (RC=0, 29 ok, no FAIL), build/gofmt/vet clean, DoD 1-5 named tests pass -race, pass-1 stale-ref throwaway now returns 0; TUI on 34-ticket copy shows tickets, CLI title set and hand edit appear <4 s; idle CPU 61.0% before (82b2252a) -> 1.0% after"
review-verdict: "Der Diff erfuellt alle fuenf DoD-Punkte. Keinen Fehler gefunden, der Berk nicht erreichen darf. Selbst gelaufen: go test ./... -race -count=1 gruen (29 ok, kein FAIL); gofmt und go vet sauber; go vet mit GOOS=windows und GOOS=darwin sauber; die neuen Tests 5x mit -race stabil; Stresstest auf dem Cache ohne Race im Cache selbst; Leerlauf-CPU ueber 30 s auf einer Kopie ohne Remote: 56,2 % (82b2252a) -> 2,1 % (HEAD); im laufenden Board waren ein CLI-Titel, ein sed-Edit und ein Edit mit gleicher Groesse und gleichem Zeitstempel nach ~0,3 s sichtbar. Unsicher nur: der /dev/null-Test auf Windows-CI (hier nicht ausfuehrbar, Quelltext spricht dafuer). Empfehlung: annehmen."
review-round: 3
review-check: |-
  1. cd /home/berk/git/jAIra-NA413Y && go build -race -o /tmp/jaira-na413y ./cmd/jaira  -> keine Ausgabe.
  2. cd /home/berk/git/jAIra-NA413Y && go build -o /tmp/jaira-na413y-cpu ./cmd/jaira  -> keine Ausgabe (ohne -race, nur fuer die CPU-Messung).
  3. rm -rf /tmp/na413y-board && cp -a /home/berk/git/jAIra /tmp/na413y-board && git -C /tmp/na413y-board remote remove origin  -> keine Ausgabe.
  4. git -C /tmp/na413y-board remote  -> gibt NICHTS aus. Steht da etwas: stoppen, das Board wuerde Refs pushen.
  5. mkdir -p /tmp/na413y-home /tmp/na413y-jhome  -> keine Ausgabe.
  6. Terminal A: cd /tmp/na413y-board && HOME=/tmp/na413y-home JAIRA_HOME=/tmp/na413y-jhome /tmp/jaira-na413y  -> Projekt-Screen; Enter druecken -> Board mit Lanes, Backlog mit ~32 Karten.
  7. Terminal B: cd /tmp/na413y-board && HOME=/tmp/na413y-home JAIRA_HOME=/tmp/na413y-jhome /tmp/jaira-na413y set FCMP17 title="NA413Y Test eins"  -> in Terminal A heisst die Backlog-Karte FCMP17 binnen 2 s "NA413Y Test eins".
  8. Terminal B: sed -i 's/^title: .*/title: NA413Y Sed zwei/' /tmp/na413y-board/.jaira/tickets/*FCMP17*.md  -> die Karte heisst binnen 2 s "NA413Y Sed zwei".
  9. Terminal B (gleiche Groesse UND gleicher Zeitstempel, das sieht nur der Waechter): F=$(ls /tmp/na413y-board/.jaira/tickets/*FCMP17*.md); touch -r $F /tmp/na413y-ref; sed -i 's/Sed zwei/Sed drei/' $F; touch -r /tmp/na413y-ref $F  -> die Karte heisst binnen 2 s "NA413Y Sed drei".
  10. Terminal A: q druecken (oder ctrl+c)  -> Board zu.
  11. Terminal A: cd /tmp/na413y-board && HOME=/tmp/na413y-home JAIRA_HOME=/tmp/na413y-jhome /tmp/jaira-na413y-cpu, dann Enter  -> Board offen. Ab jetzt nichts mehr druecken.
  12. Terminal B, 10 s spaeter: top -p "$(pgrep -n -f /tmp/jaira-na413y-cpu)"  -> %CPU von jaira-na413y-cpu meist 0-3 %. Vor dem Change waren es auf diesem Repo 55-65 %. Mit q aus top raus.
  13. Genauer, Terminal B: P=$(pgrep -n -f /tmp/jaira-na413y-cpu); a=$(awk '{print $14+$15}' /proc/$P/stat); sleep 30; b=$(awk '{print $14+$15}' /proc/$P/stat); echo "$(( (b-a)*100/(30*$(getconf CLK_TCK)) )) %"  -> nach 30 s etwa "1 %" oder "2 %" (meine Messung: 2,1 %; vorher 56 %).
  14. Terminal A: q druecken  -> Board zu.
  15. Was keinen UI-Weg hat (einmal parsen, Ref-Cache, 15-min-Netz): cd /home/berk/git/jAIra-NA413Y && go test ./core/ticket ./core/refsync ./internal/tui -race -count=1 -run 'TestDecodeParsesOnceAndMatchesBefore|TestDecodeMatchesBeforeOnEdgeCases|TestCache|TestCloneIsIndependent|TestExtraSkipsUnchangedRefs|TestExtraKeepsOrderWhenOneRefMoved|TestExtraDropsCachedTicketWhoseRefLostItsFile|TestIdleTickDecodesNothingAndChangeAppearsNextTick|TestChangeMsgRereadsOnlyTheNamedFile|TestWatcherReportsChangedPath|TestFullReparseFindsInvisibleChange'  -> drei Zeilen "ok" (core/ticket braucht ~10-50 s).
  16. Aufraeumen: rm -rf /tmp/na413y-board /tmp/na413y-home /tmp/na413y-jhome /tmp/na413y-ref /tmp/jaira-na413y /tmp/jaira-na413y-cpu  -> keine Ausgabe.
---

# Das Board laedt nur geaenderte Tickets neu und parst jedes nur einmal

## Definition of Done

- [x] Decode parst den YAML-Kopf eines Tickets genau einmal; gleiche Ausgabe wie vorher, belegt durch einen Test ueber vorhandene Ticket-Dateien
  proof: core/ticket/schema.go Decode/decodeFields, frontmatter.go Fields; TestDecodeParsesOnceAndMatchesBefore, TestDecodeMatchesBeforeOnEdgeCases
- [x] Ein Ticket, dessen Datei (Zeit und Groesse) und Ref-Hash sich seit dem letzten Laden nicht geaendert haben, wird nicht neu gelesen und dekodiert; eine Aenderung gleicher Groesse wird trotzdem gefunden; Tests decken beide Faelle
  proof: core/ticket/cache.go:72 cached; TestCacheUnchangedReadsAndDecodesNothing (0 reads/0 decodes), TestCacheRacyWindowComparesBytes + TestCacheSameSizeSameTimeFoundAfterForget (same-size change found); core/refsync/refsync.go:1092 Extra, TestExtraSkipsUnchangedRefs (ref hash)
- [x] Der Datei-Waechter gibt die geaenderten Dateien weiter, und nur diese werden neu gelesen; das volle Neuladen laeuft nur noch als seltenes Sicherheitsnetz
  proof: internal/tui/model.go:867 startWatching (chan []string), :1052 changeMsg Forget(paths); :1002 fullReparseEvery; TestChangeMsgRereadsOnlyTheNamedFile, TestWatcherReportsChangedPath, TestFullReparseFindsInvisibleChange
- [x] Gemessen auf diesem Repo: CPU des offenen Boards im Leerlauf vorher und nachher, Zahlen als Notiz am Ticket; go test ./... -race gruen
  proof: measured, see note: idle CPU 61.9/64.0 % before (82b2252a) -> 1.7/1.4 % after on a 34-ticket copy; go test ./... -race all 29 packages ok
- [x] CPU-Test: ein Test belegt, dass ein Takt ohne Aenderung kein Ticket dekodiert und keinen YAML-Kopf parst (gezaehlt, nicht in CPU-Prozent, damit er im CI stabil ist), und dass eine Aenderung trotzdem im naechsten Takt erscheint
  proof: internal/tui/idle_test.go:27 TestIdleTickDecodesNothingAndChangeAppearsNextTick (Counts decodes+headerParses 0 over idle ticks; changed title shown next tick with 1 decode)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] A (DoD 1): Test zuerst mit decodeBefore als Kopie des heutigen Decode ueber .jaira/tickets + logbook + Randfaelle; dann Fields/ticketFields/decodeFields + Counts(); Decode parst einmal; ein gruener Commit
- [x] B: lane.parse liest ueber ein einziges d.Fields() (str, list, Has, rejects-to); go test ./core/lane gruen; ein Commit
- [x] C: Ticket.Clone() mit eigenen Slices und Doc-Kopie + Test, dass Aendern der Kopie das Original nicht beruehrt; ein Commit
- [x] D (DoD 2, Dateien): Kopf-Cache im Store (EnableCache, cached, racyWindow 2 s, Forget/ForgetAll, loadHeader + idOf, Clone pro Treffer); Tests: unveraendert = 0 Lesen/0 Decode, gleich grosse Aenderung mit gleicher mtime wird gefunden, ForgetAll findet sie auch ausserhalb des Fensters, Kopie aendern laesst Cache heil, geloeschte Datei verschwindet, Cache aus = alter Weg; ein Commit
- [x] E (DoD 2, Refs): gitref.ListHashes hinter ReadMany + Hash-Cache in refsync.Extra; Test mit fileOnRef: zweites Extra ohne Decode, geaenderter Ref erscheint; ein Commit
- [x] F1 (DoD 5): Board schaltet den Cache ein (New, switchBoard); Test in internal/tui: Takt ohne Aenderung laesst Counts unveraendert, Titel per store.Mutate geaendert -> naechster Takt zeigt ihn mit genau 1 Decode; ein Commit
- [x] F2 (DoD 3, Waechter): chan []string, changeMsg{paths}, Forget vor reload; Tests: changeMsg mit einem Pfad liest nur diesen neu (1 Decode, auch bei gleicher Groesse+mtime), echtes fsnotify-Ereignis kommt mit Pfad an (5 s Timeout, Skip ohne Watcher); ein Commit
- [x] F3 (DoD 3, Sicherheitsnetz): fullReparseEvery 15 min + m.fullAt + ForgetAll im Takt, Kommentare an reload/tick umschreiben, core/release/notes/NA413Y.md im selben Commit; Test: unsichtbare Aenderung erscheint erst, wenn fullAt aelter als 15 min ist
- [x] G (DoD 4): CPU im Leerlauf vorher (82b2252a) und nachher auf einer Kopie dieses Repos ohne Remote messen, Rezept in der Bauplan-Notiz; Zahlen + go test ./... -race als Notiz; kein Commit

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
- **2026-10-13 21:41 · BeMuCa** — DoD 4 measured 2026-10-13 on a copy of this repo (34 tickets, origin removed, HOME and JAIRA_HOME in scratchpad, binaries built without -race): idle board CPU over 60 s after 8 s warm-up, two runs each. BEFORE (82b2252a): 61.9 % and 64.0 % of one core. AFTER (3d1c6bca): 1.7 % and 1.4 %. Command: tmux new-session -d -x 200 -y 50 'cd <copy> && HOME=.. JAIRA_HOME=.. exec <bin>'; Enter on the project screen; utime+stime from /proc/<pid>/stat now and +60 s, / (CLK_TCK*60). Dead end: with an empty projects.json the binary stops on the project screen, which idles at ~0.9 % for both binaries; Enter opens the board, only then the numbers mean anything. go test ./... -race -count=1 -timeout 40m on the final tree: all 29 packages ok, exit 0 (internal/tui 83.6s).
- **2026-10-13 21:46 · BeMuCa** — critique (pass 1, whole diff read): three findings, all small, fix is clear. (1) core/refsync/refsync.go:1131 comment claims the sort keeps List's order, but cache hits are appended (:1145) before fresh refs (:1168) — decode wanted into extraCache first, then append in one pass over sorted ids. No caller breaks today (board sorts at internal/tui/model.go:525, CLI starts with an empty cache), the comment is what is wrong. (2) core/ticket/frontmatter.go:246-257 Doc.find duplicates Fields.find (:191) — Doc.find = Fields() + mappingErr() + f.find(key). (3) internal/tui/idle_test.go:57 comment says agedTicket, func is agedBoard. Let stand, per the 21:03 Bauplan/Intervall notes: 15 min ForgetAll in the TUI tick rather than per-entry expiry in cache.go (ForgetAll also drops entries of renamed/filed paths), 2 s racy window with byte compare, cache off for the CLI and Load uncached, exported Counts + unexported fileReads for tests, ListHashes placed after ReadMany instead of on listNames/namesFrom (milestone branch works there), defensive Ticket.Clone. Also let stand: lane.go parse-once (same Fields, not counted), watcher retry on a full channel instead of dropping names, raw bytes kept per entry for the racy compare. Next pass reads only these three.
- **2026-10-13 21:49 · BeMuCa** — in-progress round 2: fixed critique findings. Extra now decodes wanted into extraCache first, then builds out in one pass over the sorted kept ids; TestExtraKeepsOrderWhenOneRefMoved fails on the old code (moved lowest id came out last), passes now. Doc.find reuses Fields.find. Comment typo agedBoard.
- **2026-10-13 21:50 · BeMuCa** — critique (pass 2, read only pass-1 findings + 759cf9b6): all three addressed. (1) core/refsync/refsync.go:1142 collects kept ids in sorted order, :1164-1169 builds out in one pass over kept, so the comment at :1131 now holds; TestExtraKeepsOrderWhenOneRefMoved (refsync_test.go:928) pins it. (2) core/ticket/frontmatter.go:246-255 Doc.find = Fields() + mappingErr() + f.find(key), same first-match semantics as before. (3) internal/tui/idle_test.go:57 comment says agedBoard. Repair breaks no DoD item. Nothing left; trade-offs pass 1 let stand stay closed.
- **2026-10-13 21:58 · BeMuCa** — optimize: flagged, not fixed (behaviour, not a cleanup): in Syncer.Extra, when a cached ref-only ticket's ref moves and ReadMany returns no content for it (the new commit's tree has no <id>.md, or the ref vanished between ListHashes and ReadMany), the continue at core/refsync/refsync.go:1151 leaves the old extraCache entry and the loop at :1163-1166 hands back the previously decoded ticket. Before this change that ticket was dropped. jaira never writes a ref without <id>.md, and the vanished-ref race self-heals on the next call (the prune drops ids no longer listed), so only a foreign writer reaches it. If review wants it: delete(y.extraCache, id) in the !ok branch. Also: Doc.Has and Doc.Validate lost their last non-test callers on this branch (lane.go boolOr, Decode); kept as Doc's read API with their own tests in frontmatter_test.go.
- **2026-10-13 22:04 · BeMuCa** — testing: gates green (build, gofmt, vet, go test ./... -race RC=0, 29 ok, no FAIL), DoD1-5 tests pass, hand check on scratch board copy ok (CLI set + hand sed edit show within 5 s, idle CPU 1.33 %). ONE FINDING: core/refsync/refsync.go Extra, loop over wanted (~:1149-1152) 'if !ok { continue }' leaves the old extraCache entry. Throwaway test (not committed): fileOnRef, Extra (1 ticket), then git commit-tree of the empty tree + update-ref refs/jaira/tickets/<id> to it, Extra again. Branch HEAD: 'STALE: Extra still returns 1'. Base 82b2252a same test: 'Extra returned 0 tickets' (PASS). Fix: delete(y.extraCache, id) in that !ok branch (then kept loop yields nothing) plus a regression test in core/refsync/refsync_test.go using the same steps.
- **2026-10-13 22:09 · BeMuCa** — critique (pass 3, read only the testing finding + 516b6dfa, and c0a514b4 for DoD breaks only): testing finding addressed. core/refsync/refsync.go:1152 delete(y.extraCache, id) in the !ok branch, so the kept loop at :1164-1168 finds no entry and drops the ticket; TestExtraDropsCachedTicketWhoseRefLostItsFile (refsync_test.go:966) uses the testing note's steps and passes. Side effect let stand: a ref without its file is not cached, so every Extra call asks ReadMany for it again - no decode, no parse (DoD5 holds), and only a foreign writer makes such a ref. c0a514b4 breaks no DoD item: Fields.validate lost its nil guard, but Fields() never returns nil without an error and both callers (schema.go:586-590 Decode, cache.go:101-109) check that error first; titleOf -> byID (milestonebar_test.go:307) only changes how a missing ticket fails the test. Trade-offs passes 1 and 2 let stand stay closed.
- **2026-10-13 22:11 · BeMuCa** — optimize pass 2: TestExtraDropsCachedTicketWhoseRefLostItsFile now builds its no-file ref through the package's run() instead of its own git closure. Checked it still guards the fix: with the delete at refsync.go:1152 removed it fails (Extra after the ref lost its file = 1 tickets; want 0).
- **2026-10-13 22:19 · BeMuCa** — testing (pass 2): idle CPU over 60 s on a 34-ticket remote-less copy of /home/berk/git/jAIra, HOME/JAIRA_HOME in scratchpad, tmux, non-race builds: before 3718 ticks = 61.0 %, after 97 ticks = 1.0 %. Title via CLI (set T8A8KM) and sed hand edit both visible within 4 s. Gate: go test ./... -race -count=1 RC=0, 29 ok. Throwaway stale-ref test (pass-1 finding) passes: Extra returned 0 after ref moved to empty tree.
- **2026-10-13 22:35 · BeMuCa** — review: independent checks beyond the tests. (1) Decode + Doc.Scalar/List/Has/Keys/Validate built from 82b2252a and from HEAD, run over 146 real ticket files and 12 edge cases: output byte-identical; lane.Load on this board's 13 lanes identical. Done because decodeBefore calls the refactored Doc helpers and cannot catch a regression in them. (2) Stress test (List/Load/Forget/ForgetAll in parallel, -race): no race in cache.go or Clone; the only race is on Store.dupIDs in idIndex (store.go:1217 vs :835), present on 82b2252a too - pre-existing, worth its own ticket. (3) /dev/null in TestExtraDropsCachedTicketWhoseRefLostItsFile not run on Windows; git compat/mingw.c maps /dev/null to nul in mingw_open, so expected fine - first windows-latest CI run decides. (4) Idle CPU 30 s, copy without remote: 56.2 % before, 2.1 % after.
