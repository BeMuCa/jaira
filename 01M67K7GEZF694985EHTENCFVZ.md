---
id: 01M67K7GEZF694985EHTENCFVZ
title: "Neue Lane-Arten: eine Lane, die nur einen Skill startet, und Automations-Lanes (Artefakte, E-Mail, API-Aufrufe, Telegram senden/empfangen, Video-Upload)"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "A lane can name a skill and a repository script for the agent working it, and can refuse to let a ticket leave until the ticket carries a receipt for this round, so a side effect such as a message or an upload happens at most once per round."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 25): 'einen neuen Board/Lane-Typ, wo man nur einen Skill reinhaengt und der dann gestartet wird, wenn das Ticket da landet. z.B. Ponytail-Skill soll etwas zum Code optimieren sein. Eine Lane, die nur Ponytail + vllt einen Mini-Prompt hat.' Weitere Ideen: Artefakt-Lanes (Bilder, HTML als Output), Lanes, die E-Mails versenden, bestimmte APIs aufrufen, eine Telegram-Nachricht an ihn senden oder eine von ihm erwarten, Videos per API hochladen. 'Ich denk in Richtung N8N-Automation.'
  Heute: eine Lane ist Prompt + Tier + Ein-/Ausgabevertrag; Artefakte fehlen (eigenes altes Ticket zu .jaira/artifacts/). Eine Lane, die auf eine Nachricht wartet, waere eine menschliche Lane mit externem Kanal.
  Brainstorm durch ein starkes Modell gewuenscht.
definition-of-done: "Brainstorm: Lane-Arten mit Format, Sicherheit (Secrets, Seiteneffekte), was im Binary bleibt und was ein Skill ist"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T23:44:56Z
updated-at: 2026-11-04T21:31:32Z
updated-by: berkc
claimed-by: berk-ca-1279809
claimed-at: 2026-11-04T21:24:25Z
---

# Neue Lane-Arten: eine Lane, die nur einen Skill startet, und Automations-Lanes (Artefakte, E-Mail, API-Aufrufe, Telegram senden/empfangen, Video-Upload)

## Definition of Done

- [x] Brainstorm: Lane-Arten mit Format, Sicherheit (Secrets, Seiteneffekte), was im Binary bleibt und was ein Skill ist
  proof: Progress note 2026-10-30 23:49 (brainstorm), decisions in notes 2026-11-03 17:18 and 2026-11-04 00:36
- [x] A lane file may declare 'skill: <name>'; an agentic lane with a skill needs no prompt body (the body is then the extra instruction beside the skill). Parser test covers it.
  proof: core/lane/lane.go parse() Skill; core/lane/lane_test.go TestSkillLaneNeedsNoPromptBody, TestAgenticLaneNeedsPromptOrSkill, TestSkillRefusesAPath
- [ ] A lane file may declare 'script: <path>' relative to the repository root; the loader refuses an absolute path, a path leaving the repository, and a script on a lane that is not agentic. Parser test covers it.
- [ ] 'receipt: required' on a lane: a forward move out of it is refused (code needs_receipt, field receipt-<lane>) unless the ticket field receipt-<lane> is '<key> <result>' with key = <ticket-id>-<lane>-r<round>; a receipt from an earlier round does not count; any other receipt: value fails the load. Gate test covers it.
- [ ] 'jaira show <id> --for-lane <lane> --json' carries skill, skill_path (found in the skills directories jaira roles installs into, else listed under missing), script, idempotency_key, receipt_required, receipt and a run line passing JAIRA_IDEMPOTENCY_KEY to the script; when the current key already has a receipt the run line is empty and script_done is true, so the script is not run again. The text output says the same. CLI test covers it.
- [ ] 'jaira receipt <id> <result>' writes receipt-<lane> = '<key> <result>' for the ticket's current lane; it refuses with exit 3 when the lane does not require a receipt and when the current key already has one. CLI test covers it.
- [ ] 'jaira lanes show <id>' and the generated lanes block in CLAUDE.md/AGENTS.md name a lane's skill, script and receipt requirement.
- [ ] README 'Lanes are pipeline stages' documents skill, script and receipt: required in the lane file format, including that jaira never runs the script itself.
- [ ] Release note core/release/notes/ENCFVZ.md, one line per client-facing change.
- [x] Follow-up tickets created with 'jaira create' (follows ENCFVZ, full context): max-rounds per lane, wait-for: reply lanes, lane artifacts (output-artifacts/input-artifacts).
  proof: follow-ups H6CD33 (max-rounds), M5ABXE (wait-for: reply); artifacts already ticket FCMP17, noted there instead of a duplicate

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/lane/lane.go: add Lane fields Skill, Script string and ReceiptRequired bool; parse 'skill', 'script', 'receipt' in parse(); agentic needs prompt body OR skill; script must be relative, filepath.Clean must not start with '..', lane must be agentic; receipt: only 'required' (case-insensitive) or absent, anything else an error; add the three to lanesEquivalent. Tests in core/lane/lane_test.go.
- [ ] core/gate/gate.go: add CodeNeedsReceipt = "needs_receipt", exported ReceiptField(laneID) = "receipt-"+laneID and ReceiptKey(t, laneID) = t.ID+"-"+laneID+"-r"+CurrentRound(t), and HasReceipt(t, laneID) (raw field's first space-separated word == key). In CheckAdvance, beside the OutputOwed block (to > from, !parking, lane not skipped): cur.ReceiptRequired && !HasReceipt -> violation naming 'jaira receipt <handle> <result>'. Tests in core/gate (new receipt_test.go): no receipt refused, current-key receipt passes, round-1 receipt refused in round 2, send-back not refused.
- [ ] core/role: add FindSkill(root, name string) string returning the first <dir>/<name>/SKILL.md that exists, over every agent's ProjectTarget(root) + also-dirs, then GlobalTarget + also-dirs (target.go). Test with t.TempDir and HOME set to a temp dir.
- [ ] internal/cli/flow.go showForLane: JSON keys only when the lane declares them: skill, skill_path (missing entry 'skill <name> (no SKILL.md found in ...)' when FindSkill returns empty); script, idempotency_key, receipt_required, receipt (current field value), script_done (HasReceipt), run = 'JAIRA_IDEMPOTENCY_KEY=<key> JAIRA_TICKET=<id> JAIRA_LANE=<lane> ./<script>' or empty when script_done; missing entry when the script file does not exist under s.Root. Text output: a '## Skill / script' section with the same facts and, when script_done, 'already ran for this key: do not run it again'. Tests in internal/cli.
- [ ] internal/cli: new 'jaira receipt <id> <result...>' command (copy the shape of newNoteCmd in resume.go, register in root.go AddCommand): loads env, refuses ExitValidation when the ticket's lane lacks receipt: required or HasReceipt already true; else SetScalar(ReceiptField, key+" "+result) via s.Mutate. Prints 'Receipt for <handle> in <lane>: <key> <result>'. --json emits ticketJSON. Tests: write, double write refused, lane without receipt refused. Add it to docs/COMMANDS.md if that file lists every command.
- [ ] jaira lanes show (internal/cli/tickets.go ~1440): JSON keys skill, script, receipt_required; text lines Skill:/Script:/Receipt:. core/board/announce.go LaneFact gets Skill, Script, Receipt; per-lane marks 'runs skill <name>', 'runs script <path> (JAIRA_IDEMPOTENCY_KEY from --for-lane)', 'leaving needs a receipt: jaira receipt <id> <result>'; core/lane/facts.go fills them. Test that the marks appear only for lanes that declare them.
- [ ] README.md 'Lanes are pipeline stages': a short subsection with an example lane file using skill, script and receipt: required, the key format, that jaira never runs the script (the agent does, from --for-lane's run line), and that the script must hold no secrets.
- [ ] core/release/notes/ENCFVZ.md: one line each for skill:, script:, receipt: required + 'jaira receipt'.
- [ ] Run go vet and go test for ./core/lane ./core/gate ./core/role ./core/board ./internal/cli (with the -skip of the three hanging tests) and ./core/shipped; fix what fails.

## Progress
- **2026-10-30 23:49 · berkc** — Brainstorm (starkes Modell, 31.10.2026). Gelesen: .jaira/lanes/*.md, core/lane/lane.go (parse), core/lane/template.go, core/lane/share.go, core/hook/hook.go, core/settings/settings.go, jaira-worker SKILL.md, --for-lane JSON, FCMP17 (Artefakte), ATDG74 (Schwarm).

Ausgangslage heute:
- Eine Lane ist Frontmatter + Prompt-Body. lane.go verlangt fuer agentic: true einen Prompt-Body und ein model-tier.
- input-requires / output-produces meinen Ticket-Felder; das Gate verweigert den Move, solange ein Output fehlt.
- jaira fuehrt nichts aus: ein Agent liest --for-lane und arbeitet. Einzige Ausnahme: der Hook (core/hook) - ein Skript aus ~/.jaira/settings.json, pro Person, 5 s Timeout, bei move/claim. Das ist das Vorbild fuer alles unten.
- Geteilte Lanes (.jaira/shared/) sind absichtlich inert, weil committeter Prompt fremder Text ist (T-5-02). Dieselbe Regel muss fuer alles gelten, was etwas ausfuehrt.

(1) Skill-Lane
- Neues Feld: skill: ponytail. Der Prompt-Body wird optional und ist dann der Mini-Prompt, der zusaetzlich gilt.
- lane.go: agentic verlangt Prompt-Body ODER skill. model-tier bleibt Pflicht.
- jaira loest den Namen auf (.claude/skills/<name>/SKILL.md, dann ~/.claude/skills/<name>/SKILL.md) und gibt im --for-lane JSON skill und skill_path aus. So kann auch ein Agent ohne Skill-Tool die Datei lesen - die CLI bleibt fuer jeden bash-faehigen Agenten nutzbar.
- Nicht gefunden: jaira lanes / validate warnt; --for-lane setzt missing: [skill ponytail]; die Lane wird nicht still uebersprungen.
- Worker/Dispatcher/role-lane: ein Schritt mehr in "Read the lane's brief": steht skill im Payload, den Skill mit Ticket-ID + Mini-Prompt + bounded input starten, dann wie immer output-produces setzen und moven.
- Ausgabevertrag: unveraendert output-produces. Ein Skill kennt jaira-Felder nicht - der Mini-Prompt sagt, wohin das Ergebnis gehoert (z.B. "schreib die Befunde nach review-gaps"). Ohne output-produces faellt der Default auf note: mindestens eine jaira note der Lane, sonst Move verweigert.
- Kombinierbar mit rejects-to: ponytail findet etwas -> zurueck nach in-progress, wie optimize heute.

(2) Automations-Lanes
Was in jaira gehoert und was nicht:
- Ins Binary: nur das Benennen, das Ausloesen, die Quittung. Kein HTTP-Client, kein SMTP, kein Telegram-SDK, kein Secret-Speicher, kein Scheduler.
- Alles, was nach aussen spricht, ist ein Skript oder ein Skill der Person.
Format: die Lane nennt eine Absicht, nicht ein Programm: action: telegram-ask. Welches Skript dahinter steht, steht pro Person in ~/.jaira/settings.json unter actions: {"telegram-ask": "~/bin/tg-ask.sh"} - genau wie der Hook heute. Grund: ein committetes Skript, das jaira beim Move ausfuehrt, ist Code-Ausfuehrung per Pull Request auf dem Rechner jedes Teammitglieds. Die Repo-Datei sagt "hier wird eine Telegram-Frage gestellt", der Rechner der Person entscheidet, ob und womit.
Fehlt die Zuordnung bei einer Person: die Lane ist fuer sie nicht ausfuehrbar; next --per-lane zeigt sie als "nicht eingerichtet", das Ticket wartet (nicht still durchreichen - eine E-Mail, die nie rausging, darf nicht wie gesendet aussehen).
Ausfuehren: kein Daemon, also loest ein Agent oder eine Person aus: jaira run <id>. jaira ruft das Skript mit Umgebungsvariablen wie beim Hook (JAIRA_TICKET, JAIRA_LANE, ...) plus nur den Feldern aus input-requires (JAIRA_FIELD_GOAL usw.) und den Artefakt-Pfaden. Nur was die Lane deklariert, verlaesst den Rechner.
Secrets:
- Nie im Repo, nie in der Lane-Datei, nie im Ticket, nie in einer Note. jaira kennt keinen Token und reicht keinen weiter.
- Das Skript liest seine eigenen Secrets (Umgebung, ~/.config, pass, OS-Keychain). jaira kann das nicht pruefen, nur dokumentieren.
- jaira run gibt stdout des Skripts nicht ungefiltert ins Ticket: nur eine Zeile JAIRA_RESULT=... bzw. eine Datei als Artefakt. Sonst landet ein versehentlich geloggter Token im Ticket und damit in git.
Seiteneffekte und Idempotenz:
- Quittung auf dem Ticket: action-<lane>: pending <key> vor dem Aufruf, action-<lane>: done <key> <ergebnis> danach. key = ticket-id + lane + Durchlauf-Nummer (wie oft das Ticket die Lane betreten hat).
- jaira run verweigert (Exit 3), wenn fuer denselben key schon done steht; --again nur mit Absicht.
- pending ohne done (Absturz zwischen Senden und Quittung): jaira weiss nicht, ob gesendet wurde -> keine automatische Wiederholung, sondern Frage an die Person. Ehrlich: genau-einmal ist ohne Server unmoeglich; jaira garantiert hoechstens-einmal pro key, das Skript bekommt JAIRA_IDEMPOTENCY_KEY und kann ihn beim Anbieter mitschicken (Message-ID bei E-Mail, Idempotency-Key bei APIs).
- Ein Ruecksprung (critique -> in-progress -> wieder hier) ist ein neuer Durchlauf und darf wieder senden - das ist gewollt, z.B. eine neue Demo.
- Die Quittung ist ein Ticket-Feld, wird also mit git geteilt: ein Teammitglied sieht, dass schon gesendet wurde, und sendet nicht doppelt.
Warte-Lane (Telegram-Antwort) gegen heutige human-Lane:
- Heute: human verlangt --question und requires-human-exit; die Person liest im Board und moved selbst.
- Vorschlag: keine neue Lane-Art, sondern eine human-Lane mit action: telegram-ask. Beim Betreten (oder per jaira run) schickt das Skript die offene Frage.
- Rueckweg: jaira pollt nicht. Ein Listener der Person (Telegram-Bot, cron, systemd) ruft auf ihrem Rechner jaira note <id> "<antwort>" und jaira move <id> --to <rejects-to> auf. Damit handelt die Person selbst (ihr Rechner, ihre Identitaet) und requires-human-exit bleibt wahr.
- Pflicht im Beispielskript: nur Nachrichten aus der eigenen chat-id annehmen, sonst kann jeder, der den Bot findet, Tickets bewegen.
- Unterschied zur heutigen human-Lane also nur: wie die Frage hinausgeht und von wo die Antwort kommt. Der Gate-Vertrag bleibt gleich.

(3) Artefakte (baut auf FCMP17 auf)
- Ablage: .jaira/artifacts/<handle>/<ticket-kurz>/<lane>-<name>.<ext> fuer geteilte, .jaira/artifacts/local/... gitignoriert fuer grosse oder private.
- Die Lane deklariert output-artifacts: [demo.html]; das Gate prueft beim Move, dass die Datei existiert.
- Die Folge-Lane deklariert input-artifacts: [demo.html]; --for-lane gibt Pfade, nie Inhalte (bounded input bleibt bounded).
- Groessengrenze fuer geteilte Artefakte (Vorschlag 1 MB, Text/HTML/PNG/SVG). Videos und alles Grosse nur local/ oder extern; dann steht auf dem Ticket nur ein Link-Feld (artifact-url-<name>).
- Geteiltes Artefakt = fremder Text, der als Eingabe in ein Modell laeuft; dieselbe Warnung wie bei shared/ Lanes.
- Eine action-Lane bekommt Artefakt-Pfade als JAIRA_ARTIFACT_<NAME>; so laedt ein Upload-Skript das Video hoch, ohne dass jaira je davon weiss.

(4) Weitere Lane-Ideen, nach Nutzen sortiert
1. Skill-Lanes fuer vorhandene Skills: ponytail/simplify nach in-progress, security-review mit requires-option. Hoechster Nutzen, fast keine Kosten.
2. Research-Lane vor pre-process mit skill: jaira-role-research (requires-option: research) - die Rolle gibt es schon.
3. Telegram/ntfy-Frage aus der human-Lane (action + Listener) - die Person kann unterwegs entscheiden; heute schlaeft ein Ticket, bis jemand ins Board schaut.
4. Demo-Artefakt-Lane vor signoff: Screenshot oder HTML-Seite, die zeigt, was gebaut wurde; signoff liest sie statt des Diffs. Passt zu jaira-role-acceptance.
5. Release-Notes-Lane (projektspezifisch): schreibt core/release/notes/<TICKET>.md, wenn die Aenderung client-facing ist - heute vergisst man das am leichtesten.
6. PR-Lane mit skill: jaira-role-pr: pusht den Branch und gibt die gh-Zeile zurueck (oeffnet nie selbst, wie CLAUDE.md verlangt).
7. CI-Warte-Lane: wartet auf gruene GitHub Actions; der "Listener" ist ein gh run watch, der jaira move ausfuehrt. Selbes Muster wie die Telegram-Warte-Lane.
8. Docs-Lane: aktualisiert docs/ und README fuer client-facing Aenderungen, requires-option: docs.
9. Benachrichtigung bei Ankunft in signoff (action: notify). Ueberschneidet sich mit dem Hook - nur bauen, wenn der Hook pro Lane filtern koennen soll.
10. Schwarm-Lane (ATDG74) als Variante der Skill-Lane: swarm: 5 an critique/review. Nur lesende Lanes. Eigenes Ticket, nicht hier.
11. Upload-Lane (Video/Bild per API) - geringster Nutzen fuer ein Code-Board; geht mit action + Artefakt ohne Extra-Code, also kein eigenes Feature.

(5) Was draussen bleibt, damit es kleiner als paca bleibt
- Kein Daemon, kein Scheduler, kein Polling in jaira. Ausloeser ist immer eine Session, ein jaira run oder das Skript der Person.
- Keine Netzwerk-Clients und keine Anbieter-Integrationen im Binary (kein SMTP, Telegram, YouTube, OpenAI).
- Kein Secret-Speicher, keine verschluesselten Felder.
- Keine Ablauf-Sprache: keine Bedingungen, Verzweigungen, Schleifen, Variablen in Lane-Dateien ausser dem, was rejects-to schon kann. Das waere N8N nachgebaut.
- Kein grafischer Ablauf-Editor, kein Marketplace fuer Actions.
- Keine Skripte aus dem Repo, die jaira ausfuehrt. Repo nennt, Person fuehrt aus.
- Keine Binaerdateien ueber der Grenze in git.

Format-Skizze einer Lane-Datei
Skill-Lane:
---
id: ponytail
after: in-progress
agentic: true
model-tier: strong
skill: ponytail
rejects-to: in-progress
input-requires: [goal, diff, notes]
output-produces: [review-gaps]
---
# Prompt
Nur den Diff dieses Tickets optimieren. Befunde nach review-gaps, je Zeile einer. Nichts gefunden: review-gaps="keine".

Automations-Lane (Frage per Telegram):
---
id: ask-phone
after: in-progress
agentic: false
requires-question: true
requires-human-exit: true
action: telegram-ask
input-requires: [title, question]
---
Dazu pro Person in ~/.jaira/settings.json: "actions": {"telegram-ask": "~/bin/jaira-tg-ask.sh"}

Artefakt-Lane:
---
id: demo
after: testing
agentic: true
model-tier: cheap
output-artifacts: [demo.html]
---
# Prompt
Baue eine HTML-Seite, die zeigt, was dieses Ticket geaendert hat ...

Folge-Tickets in Bau-Reihenfolge
1. skill-Feld: lane.go (Prompt ODER skill), Aufloesung skill_path, --for-lane-Feld, lanes template, validate-Warnung, jaira-worker/dispatcher/role-lane lesen es, Release-Note. Klein, sofort nuetzlich, keine Sicherheitsfrage (Skill liegt ohnehin lokal).
2. Artefakte = FCMP17 mit den Ergaenzungen oben (Pfadschema pro Ticket, Groessengrenze, local/ vs. Handle-Ordner, Pfade statt Inhalte).
3. action-Feld + actions-Zuordnung in ~/.jaira/settings.json + jaira run <id> + Quittung (pending/done, key, Exit 3 bei Wiederholung) + "nicht eingerichtet" in next --per-lane.
4. Warte-Muster: human-Lane mit action; Doku + Beispielskript (Telegram, chat-id-Pruefung) unter docs/, nicht im Binary. Kein neuer jaira-Code ausser vielleicht action beim Betreten automatisch ausloesen (offene Frage).
5. Beispiel-Lanes (ponytail, research, demo, release-notes) als shared/-Vorlagen zum Adoptieren.

Zwei Alternativen mit Kosten
A) "Alles ist ein Skill": nur Schritt 1 + 2. Automation (Mail, Telegram, Upload) macht ein Skill, der eigene Skripte aufruft; jaira weiss nichts davon. Kosten: klein (ein Feld, ein JSON-Feld, Skill-Text, Artefakte aus FCMP17). Schwaeche: keine Quittung - ob die Mail doppelt rausgeht, haengt an der Disziplin des Agenten; ein abgebrochener Lauf wiederholt sie. Kein Lauf ohne Agent (ein Agent muss die Lane oeffnen, Tokens fuer eine Mail).
B) "Action-Lanes im Binary" (Schritte 1-4): zusaetzlich action-Feld, Zuordnung pro Person, jaira run mit Quittung. Kosten: mittel - ein neuer Befehl, Settings-Erweiterung, Gate-Regel, Tests fuer Wiederholung/Absturz, Doku; geschaetzt ein Vielfaches von A. Gewinn: hoechstens-einmal pro Durchlauf, laeuft ohne Modell (cron, Listener), Quittung im git sichtbar fuer das Team.
(Verworfen: eine N8N-artige Ablauf-Engine in jaira - widerspricht git-only, kein Daemon, kleiner als paca.)

Empfehlung
A jetzt bauen (Schritt 1, dann FCMP17). Ponytail als erste echte Skill-Lane benutzen. B erst, wenn eine konkrete Automation (vermutlich die Telegram-Frage) laeuft und eine doppelte oder verlorene Nachricht wirklich wehgetan hat - dann Schritt 3 in kleinster Form (Quittung + jaira run, sonst nichts). Die Warte-Lane braucht in beiden Faellen keinen neuen Lane-Typ, nur ein Skript der Person, das jaira note + jaira move aufruft.

Offene Fragen an berkc
1. Was ist "Ponytail" genau - ein vorhandener Skill (bei dir nicht in .claude/skills oder ~/.claude/skills gefunden) oder einer, der erst geschrieben werden soll? Soll er nur Befunde liefern (wie optimize, zurueck nach in-progress) oder selbst Code aendern?
2. Darf eine Lane eine action beim Betreten automatisch ausloesen (jaira move ruft das Skript, wie heute der Hook), oder immer nur explizit per jaira run? Automatisch ist bequemer, aber dann sendet jeder Move eines Teammitglieds mit derselben Zuordnung.
3. Fehlt einer Person die actions-Zuordnung: Ticket wartet (Vorschlag) oder Lane wird fuer sie uebersprungen?
4. Sollen Artefakte ueberhaupt in git (Handle-Ordner) oder grundsaetzlich nur local/ + Link? Groessengrenze 1 MB ok?
5. Telegram-Antwort: reicht "Antwort wird Note, Ticket geht nach rejects-to", oder soll eine Antwort wie "ok" auch weiter nach vorne moven duerfen (dann entscheidet ein Bot, was heute nur eine Person darf)?
6. A oder B - oder A jetzt und B als Ticket im Backlog?
- **2026-10-31 01:25 · berkc** — Antworten Maintainer 31.10.2026: Ponytail = github.com/DietrichGebert/ponytail, ein Skill 'lazy senior dev' (minimaler Code, /ponytail-review, -audit, -debt) - passt als Skill-Lane neben optimize. Skript pro Person: der Maintainer will ein mitgeliefertes Skript (z.B. Telegram-API-Aufruf), das der Agent nutzt - keine Zuordnung pro Person noetig, solange das Skript keine Secrets enthaelt (Token aus lokalem Speicher, CMN1YT). Quittung gegen Doppelsenden: entfaellt vorerst, der Text auf dem Ticket soll ohnehin weg (eigenes Ticket). Entscheidung: Variante A (alles ist ein Skill, der Agent ruft die API).
- **2026-10-31 01:25 · berkc** — Korrektur zur Notiz davor: A oder B ist NICHT entschieden - der Maintainer hat nur gefragt, was A/B bedeuten. Unklar auch, welcher 'Text auf dem Ticket' weg soll (sein Hinweis zu 'No double sends'); nachfragen.
- **2026-11-03 17:18 · berkc** — Maintainer 03.11.2026: Entscheidung Variante A (alles ist ein Skill; der Agent ruft die API per mitgeliefertem Skript). Keine Quittung, kein jaira run. Offen bleibt nur, welcher 'Text auf dem Ticket' weg soll.
- **2026-11-04 00:36 · berkc** — Maintainer 04.11.2026: Variante A (Skill-Lane, der Agent ruft ein mitgeliefertes Skript) PLUS Quittung als Gate: wer eine solche Lane verlaesst, braucht eine Quittung auf dem Ticket, sonst verweigert der Move - wie die anderen Lane-Gates. Vorschlag Lane-Attribute (Brainstorm): skill: <name>; script: <Pfad im Repo>; receipt: required (Gate prueft Feld receipt-<lane> = '<key> <ergebnis>'); key = Ticket+Lane+Durchlauf, ans Skript als JAIRA_IDEMPOTENCY_KEY, gleicher key mit Quittung -> Skript laeuft nicht nochmal; max-rounds: <n> (heute fest 3 Ruecksendungen); wait-for: reply (Lane wartet auf Antwort, z.B. Telegram, Antwort als Note); output-artifacts / input-artifacts (aus FCMP17); vorhandene: rejects-to, requires-option, model-tier, input-requires, output-produces.
- **2026-11-04 21:28 · berkc** — Brainstorm 04.11.2026 (worker). Direction is the maintainer's: Variant A + receipt gate. Scope here: skill, script, receipt: required, --for-lane hand-over, 'jaira receipt', docs. Out: max-rounds, wait-for, artifacts (follow-ups). Sub-decisions I settled, disagree here: (1) jaira never runs the script - the agent working the lane does, from the 'run' line in --for-lane; that keeps 'no daemon, nothing runs by itself' true. (2) key = <full ticket id>-<lane>-r<round>, round = the existing 'round' field (counts send-backs, core/move/move.go:126), so a lane re-entered after a send-back is a new key and may run again; full id, not the 6-char handle, because providers dedupe on it across boards. (3) The gate only accepts a receipt whose key is the current one: a stale round-1 receipt does not let round 2 leave. (4) 'jaira receipt <id> <result>' computes the key itself so nobody types it; refuses exit 3 if the key already has a receipt or the lane requires none. (5) script path relative to repo root, must stay inside it, lane must be agentic (nobody else would run it). (6) skill resolved to a SKILL.md in the same directories 'jaira roles install' uses (core/role/target.go), so a non-Claude agent can read the file; not found -> listed in missing, the lane is not silently skipped. (7) Agents learn of it from the per-lane line in the generated CLAUDE.md/AGENTS.md block, not from edits to the hash-pinned role skills.
- **2026-11-04 21:29 · berkc** — Pre-process: plan written. Artifacts follow-up is the existing FCMP17 (note added there), not a new duplicate. Follow-ups H6CD33, M5ABXE have no assignee yet - the maintainer assigns.
