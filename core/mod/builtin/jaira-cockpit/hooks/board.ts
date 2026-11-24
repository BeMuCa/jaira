// What the cockpit knows about a jaira board it reads off the board itself,
// through the jaira binary: `jaira lanes --json` says whose lane is whose,
// `jaira list --json` where the tickets stand. Not from the agents' reports:
// so it is the same in both dispatch modes and with Herdr tabs. register.tsx
// runs the commands; this file only reads what they printed.

import type { Tab } from './sessions.ts'

export type Waiting = {
  id: string
  title: string
  // What the ticket is for, so a pane about it starts with the why.
  goal: string
  lane: string
  verdict: string
  // Which way the verdict goes by the words the board's review lane declares
  // (`verdicts` in lanes --json): accept, sendBack, or other for free text.
  // undefined where no lane declares words — a jaira before that key.
  verdictWay?: VerdictWay
  question: string
  // Whether every definition-of-done item is ticked; undefined where jaira
  // did not say (a build before dod_complete).
  dodComplete?: boolean
  // The definition-of-done items still open, apart: those only a person can
  // tick ('seen by a person: ...') and every other one; undefined where jaira
  // listed no items.
  dodOpen?: { person: number; other: number }
  // Where Accept moves it: the ticket's next lane, terminal (signoff ->
  // done) or not (script-approval -> build, human -> review: M5JSE9) - one
  // step, so a press never skips a lane; null where jaira named none.
  acceptTo: string | null
  // Whether acceptTo is the board's terminal lane: only there does Accept
  // end the work - a pull request, the signoff hold - and does a question
  // stop being what the ticket waits for.
  acceptEnds: boolean
  // Whose it is; '' for nobody's.
  assignee: string
  // The review's findings, one per line: its blocker lines, and every other
  // one. The band counts them; the Verdict pane lists them whole.
  blockers: string[]
  notes: string[]
  // False where the review's findings are unknown: a verdict of an earlier
  // round whose gaps jaira did not keep. The band then counts nothing.
  findingsKnown: boolean
  // The review's own gaps as written, blockers included: what the Give back
  // pane lists one by one for the person to decide. '' before a review.
  gaps: string
  // The round the verdict was written in, when the ticket has moved on a
  // round since; undefined while the verdict is this round's.
  reviewedIn?: number
}

export type Counts = { worker: number; review: number; person: number }

export type Lane = {
  id: string; agentic: boolean; terminal?: boolean; context: string; human_exit?: boolean; requires_question?: boolean
  verdicts?: { accept?: string[]; send_back?: string[] }
}
export type VerdictWay = 'accept' | 'sendBack' | 'other'
type Ticket = {
  handle: string
  title: string
  goal?: string
  status: string
  dod_complete?: boolean
  dod_items?: Item[]
  next_lane?: string
  question?: string
  assignee?: string
  // verdict_gaps: review-gaps as the review left them beside its verdict,
  // kept by jaira when a later round's loop lane rewrites review-gaps.
  review?: { verdict?: string; gaps?: string; verdict_gaps?: string }
  // The ticket's round, and the round its verdict was written in (0 or
  // absent: unknown, a verdict from before jaira counted rounds).
  round?: number
  review_round?: number
}

export function parse<T>(json: string): T | undefined {
  try {
    return JSON.parse(json) as T
  } catch {
    return undefined
  }
}

export const oneLine = (s: string | undefined) => (s ?? '').replace(/\s+/g, ' ').trim()

// A jaira built before `lanes --json` carried human_exit says nothing about
// whose lane is whose: then the shipped board's person lanes are taken.
const PERSON_LANES = new Set(['human', 'signoff'])

// A person's lane: one only a person moves work out of (signoff), or one
// where work waits for a person's answer (human), which jaira marks
// requires_question and not human_exit.
export function isPersons(lane: Lane): boolean {
  if (lane.human_exit === undefined && lane.requires_question === undefined) return PERSON_LANES.has(lane.id)
  return lane.human_exit === true || lane.requires_question === true
}

// The boards `jaira boards --json` names besides today's, whose tickets the
// band reads with --board. Today's board is read without it, so a jaira
// that knows no boards still shows that one.
export function otherBoardsOf(boardsJson: string): string[] {
  const boards = parse<{ boards?: { name?: string }[] }>(withoutWarnings(boardsJson))?.boards ?? []
  return boards.map(b => (b.name ?? '').trim()).filter(n => n !== '' && n !== 'default')
}

export type Agent = { description: string; name?: string; status: string }

// The chats a person opens from the band are theirs, not an agent at work.
export const PERSONS_TABS = new Set(['explain', 'discuss'])
// A subagent at work: running a turn, about to, or held on its own work.
export const WORKING = new Set(['pending', 'running', 'waiting'])
// A Herdr tab whose agent is in its turn: at work, or held at an approval
// dialog (`blocked`) that ends only when a person answers it.
export const IN_TURN = new Set(['working', 'blocked'])

// The tickets waiting in a lane a person owns: those of the person this
// session is (`me`, their name and aliases from `jaira whoami`) or of nobody,
// counted as waiting for them, and the teammates' apart. Without `me` all
// count as the person's: rather a button too many than a ticket hidden. And
// the agents at work on this board's tickets, by kind: a Herdr tab spawn.sh
// labelled `<id>/<lane>`, and a subagent of this session whose description
// or name carries a ticket's id. One whose lane is the fresh one is a review,
// every other a worker. A tab or an agent naming no ticket of this board is
// not counted.
export function snapshotOf(
  lanesJson: string,
  listJson: string,
  tabs: readonly Tab[],
  agents: readonly Agent[],
  me?: readonly string[],
): { waiting: Waiting[]; theirs: Waiting[]; counts: Counts } {
  const mine = new Set((me ?? []).map(n => n.toLowerCase()))
  const lanes = parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? []
  const tickets = parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? []
  const byId = new Map(lanes.map(l => [l.id, l]))
  const words = lanes.flatMap(l => (l.verdicts === undefined ? [] : [l.verdicts]))
  const waiting: Waiting[] = []
  const theirs: Waiting[] = []
  const counts: Counts = { worker: 0, review: 0, person: 0 }
  for (const t of tickets) {
    const lane = byId.get(t.status)
    if (lane === undefined || !isPersons(lane)) continue
    const assignee = (t.assignee ?? '').trim()
    const isMine = me === undefined || assignee === '' || mine.has(assignee.toLowerCase())
    if (isMine) counts.person++
    const next = t.next_lane === undefined ? undefined : byId.get(t.next_lane)
    const into = isMine ? waiting : theirs
    const verdict = oneLine(t.review?.verdict)
    into.push({
      id: t.handle,
      title: t.title,
      goal: (t.goal ?? '').trim(),
      lane: t.status,
      verdict,
      ...(words.length > 0 && verdict !== '' ? { verdictWay: wayOf(t.review?.verdict ?? '', words) } : {}),
      question: oneLine(t.question),
      ...(typeof t.dod_complete === 'boolean' ? { dodComplete: t.dod_complete } : {}),
      ...(t.dod_items !== undefined && t.dod_items.length > 0 ? { dodOpen: dodOpenOf(t.dod_items) } : {}),
      acceptTo: next?.id ?? null,
      acceptEnds: next?.terminal === true,
      assignee,
      // Before a review there are no review notes; after a later round's
      // optimize lane rewrote review-gaps they are in the copy jaira kept.
      blockers: verdict === '' ? [] : blockersOf(reviewGaps(t)),
      notes: verdict === '' ? [] : notesOf(reviewGaps(t)),
      findingsKnown: verdict === '' || sameRound(t) || reviewGaps(t) !== undefined,
      gaps: verdict === '' ? '' : reviewGaps(t) ?? '',
      ...(sameRound(t) ? {} : { reviewedIn: t.review_round }),
    })
  }
  const ids = tickets.map(t => t.handle)
  const fresh = lanes.filter(l => l.agentic && l.context === 'fresh').map(l => l.id)
  // A worker's tab stays open after its turn, for Give back: only one whose
  // agent is in its turn counts — working, or held at an approval dialog.
  for (const { label, agent_status } of tabs) {
    const [id = '', lane = ''] = label.split('/')
    if (!IN_TURN.has(agent_status ?? '') || !ids.includes(id) || lane === '' || PERSONS_TABS.has(lane)) continue
    if (fresh.includes(lane)) counts.review++
    else counts.worker++
  }
  for (const a of agents) {
    const said = `${a.description} ${a.name ?? ''}`
    if (!WORKING.has(a.status) || !ids.some(id => said.includes(id))) continue
    if (fresh.some(lane => new RegExp(`\\b${lane}\\b`, 'i').test(said))) counts.review++
    else counts.worker++
  }
  return { waiting, theirs, counts }
}

// Each ticket's lane, by handle, from `jaira list --json`.
export function lanesOf(listJson: string): Map<string, string> {
  return new Map((parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? []).map(t => [t.handle, t.status]))
}

// Who `jaira whoami --json` says the person is: the identity and its aliases.
export function whoOf(whoamiJson: string): string[] | undefined {
  const who = parse<{ identity?: string; aliases?: string[] }>(withoutWarnings(whoamiJson))
  const names = [who?.identity ?? '', ...(who?.aliases ?? [])].map(n => n.trim()).filter(n => n !== '')
  return names.length === 0 ? undefined : names
}

// A review verdict a person should have explained: one that names a
// reservation, or does not accept the work at all. A verdict is a sentence
// ("Approve. The diff meets the DoD ...", "the diff matches the criteria; no
// defects found"); a plain yes, however long, gets no Explain.
const RESERVATION = /vorbehalt|reservation|caveat|changes|änderungen|\bbut\b|\baber\b|however|jedoch|except|ausser|außer|allerdings|offen(e|er)? (punkt|frage)|open (point|question)/i
const REFUSAL = /^\s*(request|reject|refuse|fail|zur(ü|ue)ck|ablehn|nicht|not)\b/i

// A reservation word under a negation is none: "approve, no changes needed",
// "no open questions", "ohne Vorbehalt". The negation reaches two words at
// most and never past a comma or a full stop, and a connective is never
// negated: "no blockers, but the docs are missing" keeps its reservation.
const CONNECTIVE = /\b(but|aber|however|jedoch)\b/i
const NEGATED = new RegExp(`\\b(no|without|kein(e|en|er)?|ohne)\\s+([^\\s,;:.]+\\s+){0,2}?(${RESERVATION.source})`, 'gi')

function hasReservation(verdict: string): boolean {
  const unnegated = verdict.replace(NEGATED, m => (CONNECTIVE.test(m) ? m : ''))
  return RESERVATION.test(unnegated) || REFUSAL.test(verdict)
}

// Which way a review-verdict goes, read as jaira reads it (core/lane
// verdict.go SendsBackBy): its first line before the first ':', folded —
// case, a space like a dash, a closing full stop — must begin with a declared
// word, followed by its end or by anything that is not a letter.
const fold = (s: string) => s.trim().replace(/[.!;,]+$/, '').toLowerCase().replace(/-/g, ' ').split(/\s+/).filter(Boolean).join('-')
function wayOf(verdict: string, words: NonNullable<Lane['verdicts']>[]): VerdictWay {
  const head = fold((verdict.trim().split('\n')[0] ?? '').split(':')[0] ?? '')
  const begins = (w: string) => head.startsWith(w) && !/^\p{L}/u.test(head.slice(w.length))
  if (words.some(v => (v.send_back ?? []).some(begins))) return 'sendBack'
  if (words.some(v => (v.accept ?? []).some(begins))) return 'accept'
  return 'other'
}

export type Item = { n?: number; done?: boolean; state?: string; text?: string }

// The items still open: neither done nor superseded.
const isOpen = (i: Item) => i.state !== 'superseded' && i.done !== true
export function openOf(items: readonly Item[] | undefined): number {
  return (items ?? []).filter(isOpen).length
}

// The open definition-of-done items: those a person ticks at acceptance
// ('seen by a person: ...') and the rest.
const BY_A_PERSON = /^\s*seen by a person\s*:/i
// The open 'seen by a person' items themselves, by number and text, off
// `jaira show --json`: the Acceptance pane lists them for the person to tick
// (DKMG91). An item without a number cannot be ticked and is left out.
export type SeenItem = { n: number; text: string; checked: boolean }
export function seenItemsOf(showJson: string): SeenItem[] {
  const items = parse<{ dod_items?: Item[] }>(withoutWarnings(showJson))?.dod_items ?? []
  return items
    .filter(i => typeof i.n === 'number' && isOpen(i) && BY_A_PERSON.test(i.text ?? ''))
    .map(i => ({ n: i.n!, text: oneLine(i.text), checked: false }))
}

// Which of the items the person ticked in the pane Accept still ticks, read
// against the ticket as it is now: the pane's list was read when it opened,
// and jaira dod --done overwrites whatever an item holds. One ticked
// meanwhile keeps its own proof; one superseded, reworded or gone is changed,
// and Accept ticks nothing then - the person saw another text.
export function seenToTick(ticked: readonly SeenItem[], showJson: string): { tick: number[]; changed: string[] } {
  const now = new Map((parse<{ dod_items?: Item[] }>(withoutWarnings(showJson))?.dod_items ?? []).map(i => [i.n, i]))
  const tick: number[] = []
  const changed: string[] = []
  for (const x of ticked) {
    const i = now.get(x.n)
    if (i === undefined || i.state === 'superseded' || oneLine(i.text) !== x.text) changed.push(x.text)
    else if (i.done !== true) tick.push(x.n)
  }
  return { tick, changed }
}

// What Accept writes as the proof of a 'seen by a person' item it ticks.
export const seenProof = (who: string | undefined) => `seen by ${who?.trim() || 'a person'} in the cockpit`

function dodOpenOf(items: readonly Item[]): { person: number; other: number } {
  const person = openOf(items.filter(i => BY_A_PERSON.test(i.text ?? '')))
  return { person, other: openOf(items) - person }
}

// The light on top of the Verdict pane: green where the person can accept,
// yellow where something wants a look first, red where the ticket should go
// back. Worked out from what the ticket already says — the review's verdict
// and findings, the newest bug hunt, the definition of done, the question —
// never from a field of its own. The reasons in the order they count, every
// red one before any yellow one; the first is the one the light names.
export type LightReason =
  | { kind: 'refused' } | { kind: 'blockers'; n: number } | { kind: 'staleHunt' } | { kind: 'noHunt' } | { kind: 'dodOpen'; n: number }
  | { kind: 'notReviewed' } | { kind: 'freeVerdict' } | { kind: 'reservation' } | { kind: 'oldReview'; round: number } | { kind: 'notesUnknown' }
  | { kind: 'notes'; n: number } | { kind: 'byAPerson'; n: number } | { kind: 'dodIncomplete' } | { kind: 'question' } | { kind: 'huntUnknown' }
export type Light = { color: 'green' | 'yellow' | 'red'; reasons: LightReason[] }

// classes: the notes' marks as the Acceptance pane shows them, where it
// knows them. Then the light follows the marks (Alex 22.11.2026): a 🔴 makes
// it red, a 🟡 yellow; a known limit ⚪ or a note for another ticket 🔵 colours
// nothing. Without them the review's prefixes count, as on the band's row.
export function lightOf(w: Waiting, hunt: HuntState | undefined, classes?: readonly NoteClass[]): Light {
  const blockers = classes === undefined ? w.blockers.length : classes.filter(c => c === 'blocker').length
  const notes = classes === undefined ? w.notes.filter(n => (classOf(n) ?? 'edge') === 'edge').length : classes.filter(c => c === 'edge').length
  const red: LightReason[] = []
  const yellow: LightReason[] = []
  // The lane's own words where it declares them; a jaira that does not say
  // them leaves only the verdict's wording to go by.
  const refused = w.verdictWay !== undefined ? w.verdictWay === 'sendBack' : w.verdict !== '' && REFUSAL.test(w.verdict)
  if (refused) red.push({ kind: 'refused' })
  if (blockers > 0) red.push({ kind: 'blockers', n: blockers })
  // A hunt that ran and no longer covers HEAD is red; none finished at all
  // (at '') is a look: a ticket that changed no code has none, and one that
  // did reaches a person's lane without one only past the gate, by --force.
  if (hunt?.state === 'fail' && hunt.at !== '') red.push({ kind: 'staleHunt' })
  // The terminal lane's gate reads the definition of done: an open item other
  // than a person's own refuses the move. A jaira that listed no items says
  // only whether all are ticked, and then an open one may be the person's.
  if (w.dodOpen !== undefined && w.dodOpen.other > 0) red.push({ kind: 'dodOpen', n: w.dodOpen.other })
  if (w.verdict === '') yellow.push({ kind: 'notReviewed' })
  else if (w.verdictWay === 'other') yellow.push({ kind: 'freeVerdict' })
  else if (!refused && hasReservation(w.verdict)) yellow.push({ kind: 'reservation' })
  if (w.reviewedIn !== undefined) yellow.push({ kind: 'oldReview', round: w.reviewedIn })
  if (!w.findingsKnown) yellow.push({ kind: 'notesUnknown' })
  if (notes > 0) yellow.push({ kind: 'notes', n: notes })
  if (w.dodOpen !== undefined && w.dodOpen.person > 0) yellow.push({ kind: 'byAPerson', n: w.dodOpen.person })
  if (w.dodOpen === undefined && w.dodComplete === false) yellow.push({ kind: 'dodIncomplete' })
  if (w.question !== '') yellow.push({ kind: 'question' })
  if (hunt?.state === 'fail' && hunt.at === '') yellow.push({ kind: 'noHunt' })
  // unknown: a hunted file is changed and not committed; undefined: a jaira
  // that does not say, or a fail in a checkout without the hunt's commit.
  if (hunt === undefined || hunt.state === 'unknown') yellow.push({ kind: 'huntUnknown' })
  // A ticket waiting on a question in a lane that does not end the work
  // (human: Accept would not move it to done) is there for the answer, not
  // for acceptance: the question decides, yellow, whatever else is open.
  if (w.question !== '' && !w.acceptEnds) {
    return { color: 'yellow', reasons: [{ kind: 'question' }, ...red, ...yellow.filter(r => r.kind !== 'question')] }
  }
  const color = red.length > 0 ? 'red' : yellow.length > 0 ? 'yellow' : 'green'
  return { color, reasons: [...red, ...yellow] }
}

// The lines of a review's gaps, each on its own, the leading list marker
// ('- ', '1. ') dropped. A line that says there is nothing — 'none',
// 'keine Lücken', 'No gaps found.' — is no finding.
const NO_GAP = /^(none|keine|nothing|nichts|n\/a|-+|no (gaps?|defects?|issues?|findings?)( found)?|keine (lücken?|mängel|befunde))\s*([—–\-:(.,;!]|$)/i
const MARKER = /^\s*(?:[-*•]|\d+[.)])\s+/
// A blocker line as a review writes it, bare or in markdown: 'blocker:',
// 'Blocker:', '**blocker:**', '**Blocker**:', '`blocker`:' - the list marker
// is gone already (MARKER). Emphasis around the label is formatting, not a
// different finding.
const BLOCKER = /^[*_`]*\s*blocker\s*[*_`]*\s*:/i

function findings(gaps: string | undefined): string[] {
  return (gaps ?? '').split(/\r?\n/).map(l => oneLine(l.replace(MARKER, ''))).filter(l => l !== '' && !NO_GAP.test(l))
}

// The findings that are a blocker, and those that are not.
function blockersOf(gaps: string | undefined): string[] {
  return findings(gaps).filter(l => BLOCKER.test(l))
}

function notesOf(gaps: string | undefined): string[] {
  return findings(gaps).filter(l => !BLOCKER.test(l))
}

// The review's own gaps: review-gaps in the verdict's round, the copy jaira
// kept beside the verdict after that. Undefined for an earlier-round verdict
// without a copy: review-gaps then may be optimize's, and neither the notes
// under the row nor Explain read a blocker out of them — only a verdict that
// reserves or refuses still offers Explain.
function reviewGaps(t: Ticket): string | undefined {
  if (sameRound(t)) return t.review?.gaps
  return t.review?.verdict_gaps || undefined
}

// Whether the verdict was written in the ticket's current round. A later
// round rewrites review-gaps in its loop lanes while the old verdict stands,
// so the two belong together only then. A verdict whose round jaira did not
// record is taken as current, as it was before rounds were counted.
function sameRound(t: Ticket): boolean {
  return !t.review_round || t.review_round === (t.round ?? 1)
}

// A jaira older than the board it reads says so on a line of its own before
// anything else it prints; that line is no part of an answer or an error.
const OLDER_BUILD = /^jaira: this board was set up by an older version.*$\n?/gm

export function withoutWarnings(text: string): string {
  return text.replace(OLDER_BUILD, '').trim()
}

// What stops a move besides the person's own exit, read off the error
// `jaira move --dry-run --json` prints when a gate refuses: every other
// refusal, or undefined when the human exit (needs_human) is the only one.
// That one Accept stands for; every other gate holds, as for a move made by
// hand in the board. An error that lists no refusals is a refusal of its own.
export function otherRefusal(errorJson: string): string | undefined {
  type Refused = { error?: { message?: string; violations?: { code?: string; message?: string }[] } }
  const text = withoutWarnings(errorJson)
  let error: Refused['error']
  try {
    error = (JSON.parse(text.slice(Math.max(0, text.search(/^\{/m)))) as Refused).error
  } catch {
    return text || 'jaira refused the move'
  }
  if (error?.violations === undefined) return error?.message || 'jaira refused the move'
  const others = error.violations.filter(v => v.code !== 'needs_human')
  return others.length === 0 ? undefined : others.map(v => v.message ?? v.code).join('\n')
}

// What a chat the band opened wrote back onto the ticket, the newest of its
// kind, read off `jaira show --json`: an Explain chat's `explain:` note, a
// Discuss chat's `decision:` note, up to the next note or section, and when
// it was written (jaira stamps notes to the minute, in UTC). undefined when
// there is none yet.
export function noteOf(showJson: string, kind: 'explain' | 'decision' | 'parent session'): { text: string; at: number; line: number } | undefined {
  const lead = new RegExp(`^- \\*\\*(\\d{4}-\\d\\d-\\d\\d \\d\\d:\\d\\d) · [^*]*\\*\\* — ${kind}:\\s*`)
  const body = parse<{ body?: string }>(withoutWarnings(showJson))?.body ?? ''
  const lines = body.split('\n')
  let at = -1
  for (let i = 0; i < lines.length; i++) if (lead.test(lines[i]!)) at = i
  if (at < 0) return undefined
  const stamp = lead.exec(lines[at]!)![1]!
  const rest = lines.slice(at + 1)
  const end = rest.findIndex(l => l.startsWith('- **') || l.startsWith('#'))
  const text = [lines[at]!.replace(lead, ''), ...(end < 0 ? rest : rest.slice(0, end))].join('\n').trim()
  return { text, at: Date.parse(`${stamp.replace(' ', 'T')}:00Z`), line: at }
}

// A ticket's explanation and, once it was answered, the answer: read off the
// notes, never off a session, so a restart cannot offer an answered choice
// again. Answered means a decision: note written after the explain: note —
// by its place in the notes, since both are stamped to the minute. A newer
// explanation, a new round's finding, is open again. An explanation a give-
// back came after explains a review of a round that is over: it is no
// explanation any more, until a new one is written (Alex 22.11.2026). jaira
// stamps no round change, so the give-back is read off the note every Give
// back of the cockpit writes; a ticket moved back by hand without one keeps
// its explanation.
const GIVEN_BACK = /^- \*\*\d{4}-\d\d-\d\d \d\d:\d\d · [^*]*\*\* — Given back from /
export function explanationNotes(showJson: string): { text: string; decided?: string } | undefined {
  const explain = noteOf(showJson, 'explain')
  if (explain === undefined) return undefined
  const lines = (parse<{ body?: string }>(withoutWarnings(showJson))?.body ?? '').split('\n')
  if (lines.some((l, i) => i > explain.line && GIVEN_BACK.test(l))) return undefined
  const decision = noteOf(showJson, 'decision')
  return decision !== undefined && decision.line > explain.line ? { text: explain.text, decided: decision.text } : { text: explain.text }
}

// An Explain chat's note, in the shape spawn.sh asks for: one line per part —
// the finding, why it counts, one line per option as `1) ...` (or, in an
// older note, `a) ...`), ending in
// `[give back]` or `[accept]` when picking it is that step, and
// `Empfehlung: <number> - ...`. Its parts in the order written, the options
// apart for the choice with what each does, and the option recommended. A
// note with no option line — one written before the shape was asked for —
// is one part, whole, as it was always shown.
//
// Since DF0YD6 the note also has a status line, `Status: DoD 7/7 ok | tests
// ok | blocker fail | bughunt on HEAD ok`, one line per review note opening
// with its class — `[blocker]`, `[small edge]`, `[known limit]`, `[not this
// ticket]` — and the recommendation as a part of its own, so the pane can
// colour each. A line in neither shape stays a text part: an older note reads
// as it always did.
export type Action = 'giveBack' | 'accept'
export type Option = { key: string; text: string; action?: Action }
export type NoteClass = 'blocker' | 'edge' | 'limit' | 'elsewhere'
// unknown: the check could not be decided - drawn as a dim '?', never red.
export type Check = { label: string; ok?: boolean; unknown?: true }
export type Part =
  | { kind: 'text'; text: string }
  | ({ kind: 'option' } & Option)
  | { kind: 'status'; checks: Check[] }
  | { kind: 'note'; cls: NoteClass; label: string; text: string }
  | { kind: 'recommendation'; text: string }
export type Explanation = { parts: Part[]; options: Option[]; recommended?: string }

const OPTION = /^\s*\(?([1-6]|[a-f])\)\s+(.+)$/
const MARK = /\s*\[(give back|accept)\]\s*$/i
const RECOMMENDED = /^\s*(?:empfehlung|recommendation)\s*:\s*\(?([1-6]|[a-f])\b/i
const RECOMMENDATION = /^\s*(?:empfehlung|recommendation)\s*:/i
const STATUS = /^\s*(?:status|stand)\s*:\s*(.+)$/i
const CLASS = /^\s*[*_`]*\[(blocker|small edge|kleiner rand|known limit|bekannte grenze|not this ticket|nicht aus diesem ticket)\][*_`]*\s*(.*)$/i
const CLASSES: Record<string, NoteClass> = {
  'blocker': 'blocker', 'small edge': 'edge', 'kleiner rand': 'edge', 'known limit': 'limit',
  'bekannte grenze': 'limit', 'not this ticket': 'elsewhere', 'nicht aus diesem ticket': 'elsewhere',
}
const PASSED = /^(ok|pass|yes|ja|✓)$/i
const FAILED = /^(fail|no|nein|✗)$/i
const UNKNOWN = /^(unknown|unbekannt|\?)$/i

// One check of the status line: its last word says whether it passed, and
// is drawn as the mark instead; a check with neither word is shown as is.
function checkOf(segment: string): Check {
  const words = segment.trim().split(/\s+/)
  const last = words[words.length - 1] ?? ''
  if (words.length > 1 && PASSED.test(last)) return { label: words.slice(0, -1).join(' '), ok: true }
  if (words.length > 1 && FAILED.test(last)) return { label: words.slice(0, -1).join(' '), ok: false }
  if (words.length > 1 && UNKNOWN.test(last)) return { label: words.slice(0, -1).join(' '), unknown: true }
  return { label: segment.trim() }
}

function partOf(l: string): Part {
  const m = OPTION.exec(l)
  if (m) {
    const mark = MARK.exec(m[2]!)?.[1]?.toLowerCase()
    const text = m[2]!.replace(MARK, '').trim()
    return { kind: 'option', key: m[1]!, text, ...(mark ? { action: mark === 'accept' ? 'accept' : 'giveBack' } : {}) }
  }
  const status = STATUS.exec(l)
  if (status && status[1]!.includes('|')) return { kind: 'status', checks: status[1]!.split('|').map(checkOf).filter(c => c.label !== '') }
  const note = CLASS.exec(l)
  if (note) return { kind: 'note', cls: CLASSES[note[1]!.toLowerCase()]!, label: note[1]!, text: note[2]!.trim() }
  if (RECOMMENDATION.test(l)) return { kind: 'recommendation', text: l }
  return { kind: 'text', text: l }
}

export function explanationOf(text: string): Explanation {
  const lines = text.split('\n').map(l => l.trim()).filter(l => l !== '')
  const parts: Part[] = lines.map(partOf)
  const options = parts.flatMap(p => (p.kind === 'option' ? [{ key: p.key, text: p.text, ...(p.action ? { action: p.action } : {}) }] : []))
  const shaped = parts.some(p => p.kind === 'status' || p.kind === 'note')
  if (options.length === 0 && !shaped) return { parts: [{ kind: 'text', text: text.trim() }], options }
  const recommended = lines.map(l => RECOMMENDED.exec(l)?.[1]).find(k => k !== undefined && options.some(o => o.key === k))
  return { parts, options, ...(recommended !== undefined ? { recommended } : {}) }
}

// The main working tree, from `git rev-parse --git-common-dir`: spawn.sh puts
// a ticket's worktree beside it, never beside another worktree.
export function mainRepoOf(gitCommonDir: string, fallback: string): string {
  const dir = gitCommonDir.trim().replace(/\/+$/, '')
  if (!dir.endsWith('/.git')) return fallback
  return dir.slice(0, -'/.git'.length)
}

// Where the ticket's work lives, from `git worktree list --porcelain`: the
// worktree on a branch ending in its lower-case id, else the main working
// tree. Running `jaira move` anywhere else answers "on its ref and not on
// your disk". Herdr's review worktree of feat/<id> sits on review/feat/<id>
// and ends in the id just the same, without the ticket on its disk: of all
// that match, the branch with the fewest parts is the ticket's own. A
// detached tree — the ticket's own in the middle of a rebase — ranks after
// <prefix>/<id> and before any branch nested deeper than that.
export function worktreeFor(porcelain: string, id: string, fallback: string): string {
  return ownTreeFor(porcelain, id) ?? fallback
}

// The ticket's own worktree, undefined where it has none: a copy pulled into
// any other tree is there only for the press, a copy here belongs on the branch.
export function ownTreeFor(porcelain: string, id: string): string | undefined {
  return ticketsTree(worktreesOf(porcelain), id)?.path
}

// What the disk a press writes from holds of the ticket, off `jaira show
// --json` there: no file (on its ref alone), a file older than the ref, or a
// file a write can go to. Either of the first two refuses every write — "on
// its ref and not on your disk", a stale file — until `jaira pull` brings it.
// path is the file, once there is one. Unreadable output is taken as current:
// the write itself then says what is wrong.
export function diskOf(showJson: string): { state: 'ref-only' | 'behind' | 'current'; path: string } {
  const t = parse<{ 'on-ref-only'?: boolean; 'behind-ref'?: boolean; path?: string }>(withoutWarnings(showJson))
  const state = t?.['on-ref-only'] ? 'ref-only' : t?.['behind-ref'] ? 'behind' : 'current'
  return { state, path: t?.path ?? '' }
}

// A ticket's title out of `jaira show --json`; undefined where it says none.
export function titleOf(showJson: string): string | undefined {
  const title = parse<{ title?: unknown }>(withoutWarnings(showJson))?.title
  return typeof title === 'string' && title.trim() !== '' ? title.trim() : undefined
}

function ticketsTree(trees: Worktree[], id: string): Worktree | undefined {
  const parts = (w: Worktree) => (w.branch === '' ? 2.5 : w.branch.split('/').length)
  return trees.filter(w => isTicketsTree(w, id)).sort((a, b) => parts(a) - parts(b))[0]
}

type Worktree = { path: string; branch: string }

function worktreesOf(porcelain: string): Worktree[] {
  const out: Worktree[] = []
  for (const block of porcelain.split(/\n\n+/)) {
    const path = /^worktree (.+)$/m.exec(block)?.[1]
    const branch = (/^branch (.+)$/m.exec(block)?.[1] ?? '').replace(/^refs\/heads\//, '')
    if (path !== undefined) out.push({ path, branch })
  }
  return out
}

function isTicketsTree(w: Worktree, id: string): boolean {
  const slug = id.toLowerCase()
  return w.branch.endsWith(`/${slug}`) || w.path.endsWith(`-${slug}`)
}

// The branch a ticket's work rides on and the main working tree's, from the
// same porcelain; undefined where there is none. The main tree's branch is
// the base only where no forge can name one: in a fork the local master
// lags the upstream by hundreds of commits, every one of them "ahead".
export function branchesFor(porcelain: string, id: string): { ticket?: string; base?: string } {
  const trees = worktreesOf(porcelain)
  const ticket = ticketsTree(trees, id)?.branch
  const base = trees[0]?.branch
  return { ...(ticket ? { ticket } : {}), ...(base ? { base } : {}) }
}

// Where a branch is checked out, off the same porcelain; undefined where it
// is in no worktree. Without one, newWorktreeFor is where spawn.sh would have
// put it: beside the repository, never inside it.
export function worktreeOfBranch(porcelain: string, branch: string): string | undefined {
  return worktreesOf(porcelain).find(w => w.branch === branch)?.path
}

// Which branch a ticket's work rides on, in this order and nothing more:
// the branch its own worktree stands on, else a branch named <prefix>/<slug>
// (its worktree removed), else the one branch that holds the newest commit
// that is the ticket's own. Several holding that commit, none of them the
// ticket's by worktree or name, is never guessed at: ambiguous names them and
// the person picks (Alex, QTQW15). Three rounds of rules that guessed —
// git's order, the distance to a tip, which tip lies under which — each
// started a worker in another ticket's worktree in a case the next round
// found. undefined: no branch of the ticket's at all, and spawn.sh makes
// <prefix>/<slug>. base is the default branch, never a carrier; holders is
// `git branch --contains <that commit>`, slugBranches `git branch --list
// '*/<slug>'`, both as --format=%(refname:short) lists them.
export type Carrier = { branch: string; by: 'commit' | 'worktree' | 'name' } | { ambiguous: string[] }

export function carrierFor(input: { own?: string; base: string; holders?: string; slugBranches?: string }): Carrier | undefined {
  const { own, base, holders, slugBranches } = input
  if (own !== undefined && own !== base) return { branch: own, by: 'worktree' }
  // Two named after it — feat/<slug> and fix/<slug> — are as much a guess as
  // two holders: neither is picked.
  const named = slugBranches === undefined ? [] : holdersOf(slugBranches, base)
  if (named.length === 1) return { branch: named[0]!, by: 'name' }
  if (named.length > 1) return { ambiguous: named }
  const branches = holders === undefined ? [] : holdersOf(holders, base)
  if (branches.length === 1) return { branch: branches[0]!, by: 'commit' }
  return branches.length > 1 ? { ambiguous: branches } : undefined
}

// Where Accept moves a ticket (VHGWQF): its own worktree, else the worktree
// of the branch its commits ride on, as carrierFor settles it. A follow-up
// worked on the branch of the ticket it followed has its commits there and
// nowhere the main checkout reaches, and a move into done from the main
// checkout is refused for want of them. Several branches that could each be
// the carrier: no tree, they are named, and nothing moves (QTQW15). A carrier
// with no worktree, or no carrier: no tree either, and the move runs in the
// main checkout as before, whose refusal then names the branch.
export function acceptTreeFor(porcelain: string, id: string, carrier: Carrier | undefined): { dir?: string } | { ambiguous: string[] } {
  const own = ownTreeFor(porcelain, id)
  if (own !== undefined) return { dir: own }
  if (carrier === undefined) return {}
  if ('ambiguous' in carrier) return { ambiguous: carrier.ambiguous }
  const dir = worktreeOfBranch(porcelain, carrier.branch)
  return dir === undefined ? {} : { dir }
}

// The branches a `git branch --format=%(refname:short)` listing names, in
// its order, without the default branch, a detached HEAD — listed as "(HEAD
// detached at <sha>)", which is no branch — or a branch that only wraps
// another one listed, Herdr's review/<branch>.
function holdersOf(holders: string, base: string): string[] {
  const listed = holders.split('\n').map(l => l.trim()).filter(b => b && b !== base && !b.startsWith('('))
  return listed.filter(b => !listed.some(o => b.endsWith(`/${o}`)))
}

export function newWorktreeFor(repo: string, branch: string): string {
  const cut = repo.lastIndexOf('/')
  return `${repo.slice(0, cut)}/.worktrees/${repo.slice(cut + 1)}-${branch.slice(branch.lastIndexOf('/') + 1)}`
}

// What an Accept starts once the ticket is in done — the one rule, decided
// by Alex on GPFA7B. A branch with commits master does not have and no pull
// request yet is finished work nobody can see: the press that accepted it is
// a person's, which is the only thing allowed to open a pull request, so it
// opens one. But a branch can carry several tickets, and a pull request
// opened while one of them is still being worked proposes half the work: so
// only once every ticket its commits name is in done (`pending` empty). Any
// other case — a ticket still pending, the branch already proposed (a pull
// request in any state, merged and closed included), or nothing ahead of
// master — has no further step of its own, and the press only says the
// ticket is accepted: which ticket runs next is the teamlead's call, never
// the button's (X6JHQX). One button, so the rule decides, never the
// person per press. The caller asks the forge only where the answer can
// still change the step, and passes false otherwise. A forge that could not
// be asked — offline, logged out, its repository no remote here — is
// undefined: whether a pull request exists is unknown, so none is opened
// rather than wait for the network (Alex, GPFA7B DoD 8); the press says why.
export type StepAfterAccept = 'pr' | 'none'

export function stepAfterAccept(ahead: number, pending: readonly string[], hasPullRequest: boolean | undefined): StepAfterAccept {
  return ahead > 0 && pending.length === 0 && hasPullRequest === false ? 'pr' : 'none'
}

// The tickets a branch carries, off `git log --format=%s <base>..<branch>`:
// every handle a subject names the way this board's commits do,
// `type(A3K9QP): ...`, or several at once, `chore(A3K9QP,B7XQ2M): ...`,
// once each.
export function ticketsNamed(subjects: string): string[] {
  const out: string[] = []
  for (const line of subjects.split('\n')) {
    const scope = /^[a-z]+\(([^)]*)\)!?:/.exec(line.trim())?.[1] ?? ''
    for (const id of scope.split(/[\s,]+/)) if (/^[0-9A-Z]{6}$/.test(id) && !out.includes(id)) out.push(id)
  }
  return out
}

// The newest commit that is the ticket's own, off `git log --format=%H %s`:
// its subject's scope is the handle and nothing else, `fix(A3K9QP): ...`.
// A scope that lists it beside others — `chore(0GMYQB,A3K9QP)`, a follow-up's
// commit on the follow-up's branch — is not its own and never counts, nor
// does the handle in the text (Alex, QTQW15). The body is never read — a
// merge commit quotes the subjects it brought in there.
export function commitNaming(log: string, id: string): string | undefined {
  for (const line of log.split('\n')) {
    const cut = line.indexOf(' ')
    if (cut > 0 && /^[a-z]+\(([^)]*)\)!?:/.exec(line.slice(cut + 1).trim())?.[1]?.trim() === id) return line.slice(0, cut)
  }
  return undefined
}

// Which of those are not in the terminal lane yet — or, given several lanes,
// in none of them (signoff up to done, where .jaira/pr-when says signoff) —
// off `jaira list --json`. A handle the board does not list is filed away in the logbook — past done
// — or no ticket of this board, and holds nothing up.
export function pendingOf(listJson: string, named: readonly string[], arrived: string | readonly string[]): string[] {
  const tickets = parse<{ tickets?: { handle?: string; status?: string }[] }>(withoutWarnings(listJson))?.tickets ?? []
  const lanes = typeof arrived === 'string' ? [arrived] : arrived
  return named.filter(id => {
    const status = tickets.find(t => t.handle === id)?.status
    return status !== undefined && !lanes.includes(status)
  })
}

// The forge jaira-role-pr would open a pull request on: `git config
// jaira.forge` where it is set, else the host of origin. undefined where
// neither settles it.
export type Forge = 'github' | 'gitlab'

export function forgeOf(forge: string, originUrl: string): Forge | undefined {
  const on = forge.trim() || (/github\.com/i.test(originUrl) ? 'github' : /gitlab/i.test(originUrl) ? 'gitlab' : '')
  return on === 'github' || on === 'gitlab' ? on : undefined
}

// A git remote URL taken apart, as core/pr's parseRemote does: the host,
// without user or port, and the repository's path there without .git — for
// https://, ssh://, git+ssh://, git:// and the scp form [user@]host:path,
// whatever the user before the @ is called. A GitHub Enterprise remote often
// reads alice@ghe.example.com:team/app.git, which gh cannot resolve as it
// stands (XSGWPX). undefined where it names no host or no owner/repo.
export function parseRemote(raw: string): { host: string; repo: string } | undefined {
  const s = raw.trim()
  let host: string
  let path: string
  const scheme = s.indexOf('://')
  if (scheme > 0) {
    const rest = s.slice(scheme + 3)
    const slash = rest.indexOf('/')
    if (slash < 0) return undefined
    host = rest.slice(0, slash).replace(/^.*@/, '')
    if (!host.endsWith(']')) host = host.replace(/:[^:]*$/, '')
    path = rest.slice(slash + 1)
  } else {
    const colon = s.indexOf(':')
    if (colon <= 0 || /[/\\]/.test(s.slice(0, colon))) return undefined
    host = s.slice(0, colon).replace(/^.*@/, '')
    // A drive letter, C:/repos/app, is no host.
    if (host.length === 1) return undefined
    path = s.slice(colon + 1)
  }
  path = path.replace(/^\/+|\/+$/g, '').replace(/\.git$/, '')
  return host && path.includes('/') ? { host: host.toLowerCase(), repo: path } : undefined
}

// How gh is told the repository of origin: on a host other than github.com
// as host/owner/repo — gh's own [HOST/]OWNER/REPO, which `gh repo view` and
// --repo both take — with that host, which the answers' owner/repo then
// carry again. github.com — every *.github.com, ssh.github.com among them —
// and a URL that cannot be taken apart, go as they are: gh reads those itself.
export function ghRepoOf(originUrl: string): { repo: string; host?: string } {
  const r = parseRemote(originUrl)
  return r === undefined || r.host === 'github.com' || r.host.endsWith('.github.com') ? { repo: originUrl } : { repo: `${r.host}/${r.repo}`, host: r.host }
}

// owner/repo as gh is to be told it on host.
export function onHost(repo: string, host: string | undefined): string {
  return host ? `${host}/${repo}` : repo
}

// Where a pull request of this branch would be, as jaira-role-pr settles it:
// in a fork, the parent — asking origin there answers for the fork, finds
// nothing, and a second pull request gets opened onto a merged branch. The
// command describes a repository (origin's URL, or owner/repo);
// targetRepoOf reads the target off origin's answer, defaultBranchOf the
// branch pull requests go into off the target's.
export function repoViewCommand(on: Forge, repo: string): string[] {
  return on === 'github'
    ? ['gh', 'repo', 'view', repo, '--json', 'isFork,parent,nameWithOwner,defaultBranchRef']
    : ['glab', 'repo', 'view', repo, '-F', 'json']
}

export function targetRepoOf(on: Forge, viewJson: string): string | undefined {
  if (on === 'github') {
    const v = parse<{ nameWithOwner?: string; parent?: { name?: string; owner?: { login?: string } } | null }>(withoutWarnings(viewJson))
    const parent = v?.parent?.owner?.login && v.parent.name ? `${v.parent.owner.login}/${v.parent.name}` : undefined
    return parent ?? v?.nameWithOwner
  }
  const v = parse<{ path_with_namespace?: string; forked_from_project?: { path_with_namespace?: string } | null }>(withoutWarnings(viewJson))
  return v?.forked_from_project?.path_with_namespace ?? v?.path_with_namespace
}

export function defaultBranchOf(on: Forge, viewJson: string): string | undefined {
  const v = parse<{ defaultBranchRef?: { name?: string }; default_branch?: string }>(withoutWarnings(viewJson))
  return (on === 'github' ? v?.defaultBranchRef?.name : v?.default_branch) || undefined
}

// The local remote that points at a repository, off `git remote -v`: the
// one whose URL ends in owner/repo, SSH or HTTPS. undefined where none does —
// a fork cloned without its upstream added.
export function remoteFor(remotes: string, repo: string): string | undefined {
  const tail = new RegExp(`[:/]${repo.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}(\\.git)?$`, 'i')
  for (const line of remotes.split('\n')) {
    const [name, url] = line.trim().split(/\s+/)
    if (name && url && tail.test(url)) return name
  }
  return undefined
}

// The pull requests of a branch in that repository, in every state. Both
// tools print a JSON array, empty where there is none.
export function pullRequestsCommand(on: Forge, repo: string, branch: string): string[] {
  return on === 'github'
    ? ['gh', 'pr', 'list', '--repo', repo, '--head', branch, '--state', 'all', '--json', 'number']
    : ['glab', 'mr', 'list', '--repo', repo, '--source-branch', branch, '--all', '-F', 'json']
}

export function hasPullRequest(listJson: string): boolean | undefined {
  const list = parse<unknown>(withoutWarnings(listJson))
  return Array.isArray(list) ? list.length > 0 : undefined
}

// The directory holding .jaira/, from `cwd` upwards: the candidates in the
// order they are tried.
export function upwards(cwd: string): string[] {
  const dirs: string[] = []
  let dir = cwd.replace(/\/+$/, '') || '/'
  for (let depth = 0; depth < 12; depth++) {
    dirs.push(dir)
    if (dir === '/') break
    const cut = dir.lastIndexOf('/')
    dir = cut <= 0 ? '/' : dir.slice(0, cut)
  }
  return dirs
}

// The oldest jaira the cockpit relies on: Give back moves a ticket back out
// of a person's lane without --force, which jaira refuses before 0.3.7. The
// mod also installs on its own from the marketplace, beside any jaira, so
// this is the one place that says which one it needs.
export const MIN_JAIRA = '0.3.7'

// The version `jaira --version` printed (`jaira version 0.3.7`), or undefined
// where it names none — a build from source says `dev`, which is no release
// to compare with.
export function versionOf(stdout: string): string | undefined {
  return /\bversion\s+(\d+\.\d+\.\d+)\b/.exec(stdout)?.[1]
}

// Whether `version` comes before `min`, part by part.
export function isOlder(version: string, min: string): boolean {
  const a = version.split('.').map(Number)
  const b = min.split('.').map(Number)
  for (let i = 0; i < 3; i++) if (a[i] !== b[i]) return a[i]! < b[i]!
  return false
}

// The release the band offers an update to: the one the release cache names
// (`update-check.json`, written by `jaira self upgrade --check` — the hourly
// check the TUI and next/resume/claim/move start), when it is newer than the
// jaira installed and not the one the person dismissed. undefined in every
// other case: a dev build, a cache never written, the same version.
export function offeredRelease(
  installed: string | undefined, cacheJson: string | undefined, dismissed: string | undefined,
): { installed: string; latest: string } | undefined {
  const latest = parse<{ latest?: unknown }>(cacheJson ?? '')?.latest
  if (installed === undefined || typeof latest !== 'string' || !/^\d+\.\d+\.\d+$/.test(latest)) return undefined
  if (!isOlder(installed, latest) || latest === dismissed) return undefined
  return { installed, latest }
}

// What `jaira self upgrade --check --json` answered: the version installed,
// the latest release, and whether they are the same.
export function upgradeCheckOf(json: string): { current: string; latest: string; upToDate: boolean } | undefined {
  const c = parse<{ current?: unknown; latest?: unknown; up_to_date?: unknown }>(withoutWarnings(json))
  if (typeof c?.current !== 'string' || typeof c.latest !== 'string') return undefined
  return { current: c.current, latest: c.latest, upToDate: c.up_to_date === true }
}

// ── what the person decides per review note, in the Give back pane ───────
//
// Every finding of the review's gaps is a row of its own, and the person
// gives each one of three answers: fix it now (it goes back to the worker
// with Give back), later in a follow-up ticket, or accept it as it is. The
// answers are written onto the ticket in the note Give back or Accept writes
// anyway, so the next acceptance reads them back instead of asking again.

export type Choice = 'now' | 'later' | 'accept'
// follow: the ticket an earlier follow-up decision already created, so a
// second acceptance that keeps it a follow-up creates no second ticket.
export type GapItem = { text: string; cls: NoteClass; choice: Choice; earlier?: Choice; follow?: string }
export type Decided = { choice: Choice; follow?: string }

// A finding as the review lane writes it: 'blocker:' or 'note:', bare or in
// markdown emphasis; 'note: pre-existing:' is a fault this change did not make.
const PREFIXED = /^[*_`]*\s*(blocker|note)\s*[*_`]*\s*:/i
// The class a note states itself, right after 'note:' (GY8T52): edge, limit
// or elsewhere - the four words, with 'blocker:', that the review lane, the
// signoff role and the teamlead share. It outranks the Explain note.
const NOTE_CLASS = /^[*_`]*\s*note\s*[*_`]*\s*:[\s*_`]*(edge|limit|elsewhere)[\s*_`]*(?::|$)/i
// The two older marks: 'pre-existing' is elsewhere, 'fix-now' a small edge.
// An Explain note that classes every finding still overrides them, as it did
// before the four words.
const OLD_MARK = /^[*_`]*\s*note\s*[*_`]*\s*:[\s*_`]*(pre-?existing|fix-?now)\b/i
const WORD_CLASS: Record<string, NoteClass> = { edge: 'edge', limit: 'limit', elsewhere: 'elsewhere', preexisting: 'elsewhere', fixnow: 'edge' }

// The class a finding carries in its own words: 'blocker:', a class word
// after 'note:', or an older mark. Undefined for a note without any - an
// older review's - which the caller classes as before.
export function classOf(text: string): NoteClass | undefined {
  if (BLOCKER.test(text)) return 'blocker'
  const m = NOTE_CLASS.exec(text) ?? OLD_MARK.exec(text)
  return m ? WORD_CLASS[m[1]!.toLowerCase().replace('-', '')] : undefined
}

// What each class starts as: a blocker fixed now, every other note accepted
// (Alex 21.11.2026). A follow-up ticket is created only for a row the person
// turned to it: a press on Accept alone files nothing on the shared board.
const PRESELECTED: Record<NoteClass, Choice> = { blocker: 'now', edge: 'accept', limit: 'accept', elsewhere: 'accept' }

// The findings of a review's gaps, each with its class and what it starts
// as. A line without 'blocker:' or 'note:' continues the finding above it —
// a note wrapped over two lines — but only where the gaps use the prefixes
// at all; an older review's gaps, without them, are a finding per line.
// The class is the one the finding states (classOf): 'blocker:', or 'note:'
// followed by edge, limit or elsewhere, which the review lane writes since
// GY8T52 - so the pane shows the class the teamlead and the signoff role
// read. A finding without one, or with only the older 'pre-existing:' or
// 'fix-now:' mark, takes the Explain note's class when that note classes
// exactly as many notes as there are findings, in order; otherwise its mark's
// class, else edge. A 'blocker:' line stays a blocker whatever the Explain
// note says at its position: matching by position may be off, and a blocker
// read as a known limit would be accepted unseen. A
// note decided before, word for word, starts as it was decided. The acceptance rule: a finding classed
// wrong is better than one lost or two merged — no line of the gaps but a
// "none" ever drops out.
export function gapItems(gaps: string | undefined, explainText?: string, decided?: ReadonlyMap<string, Decided>): GapItem[] {
  const lines = findings(gaps)
  const prefixed = lines.some(l => PREFIXED.test(l))
  const texts: string[] = []
  for (const l of lines) {
    if (prefixed && !PREFIXED.test(l) && texts.length > 0) texts[texts.length - 1] += ` ${l}`
    else texts.push(l)
  }
  const explained = explainText === undefined ? [] : explanationOf(explainText).parts.flatMap(p => (p.kind === 'note' ? [p.cls] : []))
  return texts.map((text, i) => {
    const stated = classOf(text)
    const firm = BLOCKER.test(text) || NOTE_CLASS.test(text)
    const cls: NoteClass = firm ? stated! : explained.length === texts.length ? explained[i]! : stated ?? 'edge'
    const earlier = decided?.get(text)
    if (earlier === undefined) return { text, cls, choice: PRESELECTED[cls] }
    return { text, cls, choice: earlier.choice, earlier: earlier.choice, ...(earlier.follow !== undefined ? { follow: earlier.follow } : {}) }
  })
}

// The headings of the decision lines, in the words the ticket asked for.
const HEADINGS: Record<Choice, string> = { now: 'beheben:', later: 'Folge-Ticket:', accept: 'angenommen:' }
const EXTRA = 'dazu: '

// The decision lines the Give back or Accept note carries: a heading per
// answer given, each finding word for word under it, a follow-up with the
// ticket it became (or why none was created). What the person typed is part
// of the work to do: it stands under 'beheben:', after the notes to fix, so
// a give-back with words only still opens with that list.
export function decisionLines(items: readonly GapItem[], follows: ReadonlyMap<string, string> = new Map(), extra = ''): string {
  const out: string[] = []
  const typed = extra.trim() === '' ? [] : [`${EXTRA}${oneLine(extra)}`]
  for (const c of ['now', 'later', 'accept'] as const) {
    const of = items.filter(i => i.choice === c)
    const more = c === 'now' ? typed : []
    if (of.length === 0 && more.length === 0) continue
    out.push(HEADINGS[c], ...of.map(i => `- ${c === 'later' ? `${follows.get(i.text) ?? i.follow ?? '(not created)'}: ` : ''}${i.text}`), ...more)
  }
  return out.join('\n')
}

// A note the cockpit wrote when a person gave a ticket back or accepted it
// with decisions, in the shape decisionLines writes, after its first words.
const DECIDED_LEAD = /^(?:Given back|Accepted) from [^:\s]+:\s*/
const NOTE_START = /^- \*\*\d{4}-\d\d-\d\d \d\d:\d\d [^*]*\*\* — /
const FOLLOW = /^(?:([0-9A-Z]{6})|\(not created[^)]*\)):\s*/

// What was decided before, read off the notes of `jaira show --json`: each
// finding by its words, the newest decision winning. A note that is not the
// cockpit's, or a line after the decisions, is read past.
export function decisionsOf(showJson: string): Map<string, Decided> {
  const body = parse<{ body?: string }>(withoutWarnings(showJson))?.body ?? ''
  const out = new Map<string, Decided>()
  let heading: Choice | undefined
  let inNote = false
  for (const raw of body.split('\n')) {
    let line = raw
    if (NOTE_START.test(line)) {
      line = line.replace(NOTE_START, '')
      inNote = DECIDED_LEAD.test(line)
      heading = undefined
      if (!inNote) continue
      line = line.replace(DECIDED_LEAD, '')
      if (line === '') continue
    } else if (!inNote) continue
    const h = (Object.keys(HEADINGS) as Choice[]).find(c => HEADINGS[c] === line.trim())
    if (h !== undefined) heading = h
    else if (heading !== undefined && line.startsWith('- ')) {
      const text = line.slice(2).trim()
      const follow = heading === 'later' ? FOLLOW.exec(text) : null
      out.set(follow ? text.slice(follow[0].length) : text, { choice: heading, ...(follow?.[1] ? { follow: follow[1] } : {}) })
    } else if (!line.startsWith(EXTRA)) inNote = false
  }
  return out
}

// The handle `jaira create --json` printed, for the follow-up line.
export function createdHandle(createJson: string): string | undefined {
  return parse<{ handle?: string }>(withoutWarnings(createJson))?.handle
}

// The lane and board `jaira show --json` gives, read again right before
// decisions are carried out: someone else may have moved the ticket since the
// pane opened, and a follow-up goes on the ticket's own board.
export function standOf(showJson: string): { status?: string; board: string } {
  const s = parse<{ status?: string; board?: string }>(withoutWarnings(showJson))
  return { ...(s?.status !== undefined ? { status: s.status } : {}), board: s?.board ?? '' }
}

// A follow-up's title: the finding without its 'note:', its class word and
// an older mark after it - 'pre-existing:' or the review's 'fix-now:' - once
// each, so a finding whose own text opens with such a word keeps it; cut to
// fit a list row; the
// whole finding is its context.
export function followTitle(text: string): string {
  const bare = text.replace(PREFIXED, '').trim().replace(/^(?:[\s*_`]*(?:edge|limit|elsewhere)[\s*_`]*:\s*)?(?:[\s*_`]*(?:pre-?existing|fix-?now)[\s*_`]*:\s*)?/i, '').trim() || text
  return bare.length <= 72 ? bare : `${bare.slice(0, 71).trimEnd()}…`
}

// The newest bug hunt's state, read off `jaira show <id> --brief --json`:
// ok, fail, unknown or skipped as jaira worked it out against the checkout
// it ran in; fail where no hunt finished (bughunt null); undefined where the
// brief has no bughunt key — a jaira older than the hunt state — or did not
// parse, so the band says it does not know rather than guess. `at` is the
// commit the hunt ran at, '' where none finished: jaira also answers fail
// for a checkout that does not hold that commit, which the caller checks.
export type Hunt = 'ok' | 'fail' | 'unknown' | 'skipped'
export type HuntState = { state: Hunt; at: string }
const HUNTS = new Set<string>(['ok', 'fail', 'unknown', 'skipped'])
export function huntOf(briefJson: string): HuntState | undefined {
  const brief = parse<{ bughunt?: { state?: string; at?: string } | null }>(withoutWarnings(briefJson))
  if (brief === undefined || !('bughunt' in brief)) return undefined
  if (brief.bughunt === null) return { state: 'fail', at: '' }
  const state = brief.bughunt?.state ?? ''
  return HUNTS.has(state) ? { state: state as Hunt, at: brief.bughunt?.at ?? '' } : undefined
}

// ── following a pull request (CRASZS) ────────────────────────────────────

// The pull requests a report names by address — GitHub's /pull/<n>,
// GitLab's /-/merge_requests/<n> — once each, in order. A bare "#12" is not
// taken: an issue, a list item and a pull request all write it. Whose they
// are, jaira pr says (followNamed).
export function pullRequestsIn(text: string): { number: number; url: string }[] {
  const out: { number: number; url: string }[] = []
  for (const m of text.matchAll(/https?:\/\/[^\s)>\]]+?\/(?:pull|-\/merge_requests)\/(\d+)\b/g)) {
    if (!out.some(x => sameAddress(x.url, m[0]))) out.push({ number: Number(m[1]), url: m[0] })
  }
  return out
}

// Whether two addresses of a pull request are the same one: scheme and case
// aside.
export function sameAddress(a: string, b: string): boolean {
  const norm = (u: string) => u.toLowerCase().replace(/^https?:\/\//, '').replace(/\/+$/, '')
  return norm(a) === norm(b)
}

// One pull request as `jaira pr --json` prints it, as far as the cockpit
// reads it.
export type PullRequest = { number: number; url: string; tickets: string[]; settled: boolean; unsettled: string[] }

// The open pull requests `jaira pr --json` lists; undefined where the
// answer is not that command's.
export function pullRequestsOf(prJson: string): PullRequest[] | undefined {
  const st = parse<{ prs?: (Partial<PullRequest> & { state?: string })[] }>(withoutWarnings(prJson))
  if (!Array.isArray(st?.prs)) return undefined
  return st.prs
    .filter(p => typeof p.number === 'number' && (p.state ?? 'open') === 'open')
    .map(p => ({ number: p.number!, url: p.url ?? '', tickets: p.tickets ?? [], settled: p.settled === true, unsettled: p.unsettled ?? [] }))
}

// When a board's pull request opens, from `jaira roles list --json`:
// 'signoff' where .jaira/pr-when says so, 'done' otherwise and for a jaira
// too old to say.
export function prWhenOf(rolesJson: string): 'signoff' | 'done' {
  return parse<{ pr_when?: string }>(withoutWarnings(rolesJson))?.pr_when === 'signoff' ? 'signoff' : 'done'
}

// The lanes from `from` up to `to`, both included, in the board's order off
// `jaira lanes --json` — where a ticket has arrived once the pull request may
// open in signoff. Empty where either lane is not on the board, or `to`
// comes first.
export function lanesBetween(lanesJson: string, from: string, to: string): string[] {
  const ids = (parse<{ lanes?: { id?: string }[] }>(withoutWarnings(lanesJson))?.lanes ?? []).map(l => l.id ?? '')
  const a = ids.indexOf(from)
  const b = ids.indexOf(to)
  return a < 0 || b < a ? [] : ids.slice(a, b + 1)
}

// ── review notes as a person reads them ─────────────────────────────────

// The mark each note opens with, the one the teamlead's signoff reports use:
// a blocker red, a small edge yellow, a known limit white, one that belongs
// to another ticket blue. It says the class before a word is read.
export const NOTE_MARK: Record<NoteClass, string> = { blocker: '🔴', edge: '🟡', limit: '⚪', elsewhere: '🔵' }

// The language the person reads the notes in, from their own settings.json
// ("language"); undefined where it names none, or the file does not parse:
// then the notes stay as the ticket has them and nothing is sent anywhere.
export function languageOf(settingsJson: string | undefined): string | undefined {
  const s = parse<{ language?: unknown }>(settingsJson ?? '')
  const lang = typeof s?.language === 'string' ? oneLine(s.language) : ''
  return lang === '' ? undefined : lang
}

// The cheap model's prompt for the Acceptance pane's Translate: the head -
// the ticket's title, goal, what was done and why, the reason the light
// names - and every review note, in one call, back as one JSON array of as
// many strings in the same order. A note's class prefix goes - the mark says
// it.
export function translatePrompt(language: string, head: readonly string[], notes: readonly string[]): string {
  return [
    `Translate each text below into ${language}, plainly, for a person who reads slowly and was not there.`,
    `The first ${head.length} are the title of a ticket, its goal, what was done and why, and the reason a traffic light shows; a title stays short.`,
    ...(notes.length === 0 ? [] : [
      `The last ${notes.length} are review notes: each one or two short, plain sentences saying what is wrong or missing and what it means.`,
      'Drop a leading "blocker:", "note:", "edge:", "limit:", "elsewhere:", "pre-existing:" or "fix-now:" from a review note - a mark in front of it already says it.',
    ]),
    'Keep code identifiers, file paths, commands and ticket handles exactly as written. Add nothing a text does not say.',
    `Answer with a JSON array of exactly ${head.length + notes.length} strings, in the same order, and nothing else.`,
    '',
    JSON.stringify([...head, ...notes]),
  ].join('\n')
}

// The model's answer as one shown text per note, or undefined when it is no
// array of exactly that many non-empty strings: a note shown wrong, or two
// notes swapped, is worse than the original text.
export function rephrasedOf(answer: string, n: number): string[] | undefined {
  const start = answer.indexOf('[')
  const end = answer.lastIndexOf(']')
  if (start < 0 || end <= start) return undefined
  const out = parse<unknown>(answer.slice(start, end + 1))
  if (!Array.isArray(out) || out.length !== n) return undefined
  if (!out.every(x => typeof x === 'string' && oneLine(x) !== '')) return undefined
  return out.map(x => oneLine(x as string))
}

// ── the Acceptance pane ─────────────────────────────────────────────────

// What Execute does with the person's choices: one rule over every
// combination, so the pane's summary line and the press agree. A note marked
// fix now gives the ticket back, follow-ups and all; words typed beside no
// such note are what is missing, and give it back too. A blocker not marked
// fix now is accepted only with a reason, which the field then asks for: its
// words are that reason, not work to do. Nothing to fix accepts — through the
// same gates as Accept, into the next lane whether or not it ends the work
// (M5JSE9); where jaira named no next lane there is nothing to accept: words
// typed give it back, nothing typed does nothing. An accept with 'seen by a person' items
// still unticked in the pane (unseen) is no accept yet: it waits for them.
export type Step = 'giveBack' | 'accept' | 'needWhy' | 'noNext' | 'needSeen'
export type Execution = { step: Step; fix: number; followUps: number; asksWhy: boolean }

export function executionOf(items: readonly GapItem[], typed: string, acceptTo: string | null, unseen = 0): Execution {
  const fix = items.filter(i => i.choice === 'now').length
  const followUps = items.filter(i => i.choice === 'later' && i.follow === undefined).length
  // Only where the ticket can be accepted at all: elsewhere words typed are
  // what is missing, never a reason to accept.
  const asksWhy = acceptTo !== null && fix === 0 && items.some(i => i.cls === 'blocker' && i.choice !== 'now')
  const words = typed.trim() !== ''
  const step: Step = fix > 0 || (words && !asksWhy) ? 'giveBack' : acceptTo === null ? 'noNext' : asksWhy && !words ? 'needWhy' : unseen > 0 ? 'needSeen' : 'accept'
  return { step, fix, followUps, asksWhy }
}

// The pane's head and history, read off one `jaira show --json`: what the
// ticket is for, what was done and why, and how it got here. Commits, files,
// the bug hunt, the pull request and the release note come from elsewhere
// and are added by the caller.
export type Story = {
  title: string
  goal: string
  what: string
  why: string
  round: number
  // How often in a row it came back, and from where; undefined: never.
  back?: { count: number; from: string }
  dod: { done: number; superseded: number; total: number }
  verdict: string
  // created_at in ms; undefined where jaira did not say.
  created?: number
}

export function storyOf(showJson: string): Story | undefined {
  type Shown = {
    title?: string; goal?: string; outcome?: { what?: string; why?: string }; round?: number
    back?: { count?: number; from?: string }; dod_items?: Item[]; review?: { verdict?: string }; created_at?: string
  }
  const s = parse<Shown>(withoutWarnings(showJson))
  if (s === undefined || s === null || typeof s !== 'object') return undefined
  const items = s.dod_items ?? []
  const superseded = items.filter(i => i.state === 'superseded').length
  const created = s.created_at === undefined ? NaN : Date.parse(s.created_at)
  const count = s.back?.count ?? 0
  return {
    title: oneLine(s.title),
    goal: oneLine(s.goal),
    what: oneLine(s.outcome?.what),
    why: oneLine(s.outcome?.why),
    round: s.round ?? 1,
    ...(count > 0 ? { back: { count, from: s.back?.from ?? '' } } : {}),
    dod: { done: items.filter(i => i.state !== 'superseded' && i.done === true).length, superseded, total: items.length },
    verdict: oneLine(s.review?.verdict),
    ...(Number.isNaN(created) ? {} : { created }),
  }
}

// The lines of a release note file that are changes ('- ...'), the mark of
// one the reader must act on ('- ! ') kept; prose and comments read past.
export function releaseLines(note: string): string[] {
  return note.split(/\r?\n/).filter(l => l.startsWith('- ')).map(l => oneLine(l.slice(2))).filter(l => l !== '')
}

// The commits and the files they changed, from `git log --format=%h
// --name-only` over the commits naming the ticket: a hash line opens each.
export function commitsOf(log: string): { commits: number; files: number } {
  const lines = log.split(/\r?\n/).map(l => l.trim()).filter(l => l !== '')
  const hashes = lines.filter(l => /^[0-9a-f]{7,40}$/.test(l))
  return { commits: hashes.length, files: new Set(lines.filter(l => !hashes.includes(l))).size }
}

// Explain on one note: the cheap model says in one everyday sentence why it
// matters for accepting the ticket, or why it does not - told what the
// ticket is for, what was done, and the change to the files the note names,
// so it can say what happens in the code rather than repeat the note.
export type NoteContext = { title: string; goal: string; what: string; diff: string }
export const DIFF_LINES = 200

export function reasonPrompt(language: string | undefined, ctx: NoteContext, note: string): string {
  return [
    `A person decides whether to accept the ticket "${ctx.title}". A reviewer left the note below.`,
    ...(ctx.goal !== '' ? [`The ticket's goal: ${ctx.goal}`] : []),
    ...(ctx.what !== '' ? [`What was done: ${ctx.what}`] : []),
    `In ${language ?? 'the language the note is written in'}, say in one short everyday sentence why this note matters for that decision, or why it does not - name concretely what the code does, as the change below shows.`,
    'Keep code identifiers, file paths and ticket handles exactly as written. Answer with the sentence alone.',
    '',
    note,
    ...(ctx.diff !== '' ? ['', `The change to the files the note names (at most ${DIFF_LINES} lines):`, ctx.diff] : []),
  ].join('\n')
}

// The files a note names: words with a file extension, a ':line' after them
// dropped - 'register.tsx:1214', 'core/release/notes/<ID>.md'. A placeholder
// with '<' or '>' names no file. An absolute path names the file by its
// name: it is matched at any depth, and a checkout elsewhere has another root.
export function notePaths(note: string): string[] {
  const found = note.match(/[\w./-]+\.(?:tsx?|jsx?|go|md|sh|py|json|ya?ml|toml|css|html)(?::\d+)?\b/g) ?? []
  return [...new Set(found.map(p => p.replace(/:\d+$/, '')).map(p => (p.startsWith('/') ? p.slice(p.lastIndexOf('/') + 1) : p)))]
}

// A diff cut to its first n lines, with a line saying how many were left out.
export function clipLines(text: string, n: number): string {
  const lines = text.replace(/\n$/, '').split('\n')
  return lines.length <= n ? lines.join('\n') : [...lines.slice(0, n), `... ${lines.length - n} more lines`].join('\n')
}

// The model's sentence, on one line; undefined when it said nothing.
export function reasonOf(answer: string): string | undefined {
  const line = oneLine(answer)
  return line === '' ? undefined : line.length <= 400 ? line : `${line.slice(0, 399).trimEnd()}…`
}
