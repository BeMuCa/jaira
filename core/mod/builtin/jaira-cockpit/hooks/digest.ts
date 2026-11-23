// digest.ts reads one agent report and boils it down to the line the
// teamlead chat shows in its place. Everything here is pure: no `$`, so the
// rules can be read and tested without an engine.

import { filterWords as words } from './texts.ts'

// A handle is the last six characters of a ticket's ULID, Crockford base32
// (no I, L, O, U). Six capital letters alone are too often an ordinary word
// (STATUS, PASSED), so a handle without a digit counts only where jaira's
// own wording puts one: after "ticket", in a commit subject, after a jaira
// verb.
const HANDLE = '[0-9A-HJKMNP-TV-Z]{6}'
const ULID = /\b[0-9A-HJKMNP-TV-Z]{20}([0-9A-HJKMNP-TV-Z]{6})\b/g
// A number alone is never one: a task-notification's <usage> block carries
// six-digit counts (total_tokens, duration_ms) that would outvote the handle.
const WITH_DIGIT = new RegExp(`\\b(?=[A-Z]*[0-9])(?=[0-9]*[A-Z])${HANDLE}\\b`, 'g')
const IN_CONTEXT = new RegExp(
  `(?:\\b[Tt]icket:?\\s+|\\w\\(|\\bjaira\\s+[a-z-]+\\s+(?:--\\S+\\s+)*)(${HANDLE})\\b`,
  'g',
)

// The ticket a report is about: the handle named most often, the first one
// named on a tie. Undefined when none is recognisable — the report is then
// passed on unchanged rather than guessed at.
export function ticketOf(text: string): string | undefined {
  const seen = new Map<string, number>()
  const count = (id: string) => seen.set(id, (seen.get(id) ?? 0) + 1)
  for (const m of text.matchAll(ULID)) count(m[1]!)
  for (const m of text.matchAll(WITH_DIGIT)) count(m[0])
  for (const m of text.matchAll(IN_CONTEXT)) count(m[1]!)
  let best: string | undefined
  for (const [id, n] of seen) if (best === undefined || n > seen.get(best)!) best = id
  return best
}

// A report that has to stay whole: a task-notification whose status is not
// "completed", or a text that says something broke, was aborted, or waits at
// an approval dialog — a "completed" status included, since an agent that
// ends its turn normally can still report a failed run. "0 failed" in a test
// summary is not a failure; "3 failed" is. FAIL in capitals is how go test
// reports one, and only in capitals, so it is matched without /i.
const FAILURE =
  /\b(?:fehlgeschlagen|abgebrochen|aborted|crashed|timed out|killed|panic:|approval dialog)|(?<!\b0 )\b(?:failed|failures?)\b|\bexit (?:code |status )?[1-9]\d*\b|^\s*error\b/im
const GO_FAIL = /\bFAIL\b/
const STATUS = /<status>\s*([a-z_-]+)\s*<\/status>/i

export function isFailure(text: string): boolean {
  const status = STATUS.exec(text)
  if (status && status[1]!.toLowerCase() !== 'completed') return true
  return FAILURE.test(text) || GO_FAIL.test(text)
}

// A question ends at a quote, the line's end, or a "·" — the separator of a
// status line, whose next field is not part of the question.
const QUESTION = /(?:--question\s+["“]|\b(?:open question|offene frage|question|frage)\s*[:=]\s*["“]?)([^"\n”·]+)/i
const NO_QUESTION = /^(?:none|no|keine|nein|nichts|-|—)\b/i

export function questionOf(text: string): string | undefined {
  const m = QUESTION.exec(text)
  if (!m) return undefined
  const q = m[1]!.trim()
  return q === '' || NO_QUESTION.test(q) ? undefined : clip(q, 100)
}

// The status line a worker or dispatcher ends its report with (see the
// Report section of jaira-worker and jaira-dispatcher):
//
//   jaira: 55CS2M in-progress -> signoff · branch: feat/55cs2m · commits: a1b2c3d..e4f5a6b (3) · gates: ok · notes: 1 · verdict: approve
//
// Fixed keys in a fixed order, "-" for an empty value; verdict comes last and
// keeps everything after its key, so its text may hold a "·" of its own.
export type StatusLine = {
  ticket: string
  from: string
  to: string
  branch?: string
  commits?: string
  gates?: string
  notes?: string
  verdict?: string
}

const FIELDS = ['branch', 'commits', 'gates', 'notes'] as const
// A lane id is what core/lane accepts: lowercase letters, digits and hyphens
// (qa2, 2-review). A narrower pattern drops the whole line on such a board.
const LANE = '[a-z0-9][a-z0-9-]*'
const STATUS_HEAD = new RegExp(
  `^\\s*jaira:\\s+(${HANDLE})\\s+(${LANE})\\s*(?:->|→)\\s*(${LANE})\\s*(?:·\\s*(.*))?$`,
  'i',
)
// The engine's frame around a message — <result>, <agent-message from="…">
// and their closing tags — can share the physical line with the status line,
// on either side. They are not part of it.
const FRAME_OPEN = /^(?:\s*<[a-z][a-z-]*(?:\s[^>]*)?>)+/i
const FRAME_CLOSE = /(?:<\/[a-z][a-z-]*>\s*)+$/i

// The one line above the status line that says what waits for a person
// (see the Report section of jaira-worker): "needs you: <what>". It is all a
// report carries besides the status line, so a digest that dropped it would
// drop the only thing the teamlead has to act on. Read like the status line:
// the last such line, frame tags around it stripped — in a two-line report
// it is the first line, the one the engine's frame opens on. Only an empty
// value is no ask: the worker writes no line when nothing waits, and a word
// list would also drop a real ask that starts "no answer from CI…".
const NEEDS_YOU = /^[ \t]*needs you:[ \t]*(.*?)[ \t]*$/i

export function needsYouOf(text: string): string | undefined {
  const lines = text.split(/\r?\n/)
  for (let i = lines.length - 1; i >= 0; i--) {
    const m = NEEDS_YOU.exec(lines[i]!.replace(FRAME_OPEN, '').replace(FRAME_CLOSE, ''))
    if (!m) continue
    return m[1] === '' ? undefined : clip(m[1]!, 140)
  }
  return undefined
}

// The last such line in the report, or undefined when there is none: a
// report quoting an earlier one still ends with its own. A key that is
// missing or holds "-" is left out; a key not in FIELDS (or verdict) is ignored.
export function statusLineOf(text: string): StatusLine | undefined {
  const lines = text.split(/\r?\n/)
  for (let i = lines.length - 1; i >= 0; i--) {
    const m = STATUS_HEAD.exec(lines[i]!.replace(FRAME_OPEN, '').replace(FRAME_CLOSE, ''))
    if (!m) continue
    const line: StatusLine = { ticket: m[1]!.toUpperCase(), from: m[2]!.toLowerCase(), to: m[3]!.toLowerCase() }
    let rest = (m[4] ?? '').trim()
    const verdict = /(?:^|·)\s*verdict:\s*(.*)$/i.exec(rest)
    if (verdict) {
      line.verdict = verdict[1]!.trim()
      rest = rest.slice(0, verdict.index)
    }
    for (const field of rest.split('·')) {
      const kv = /^\s*([a-z]+):\s*(.*?)\s*$/i.exec(field)
      if (!kv) continue
      const key = kv[1]!.toLowerCase()
      if ((FIELDS as readonly string[]).includes(key)) line[key as (typeof FIELDS)[number]] = kv[2]!
    }
    for (const key of [...FIELDS, 'verdict'] as const) if (line[key] === '' || line[key] === '-') delete line[key]
    return line
  }
  return undefined
}

// The line the chat shows instead of the report, or undefined when the
// report must be passed on whole: no ticket recognisable, or a failure.
// A report that ends with a status line is read from that line alone. One
// without it shows the ticket and nothing else: a lane or a verdict picked
// out of free text ("reject? reason: …", "nothing rejected", "all the way to
// signoff") was wrong often enough that an id alone is the better line.
export function digest(text: string): { ticket: string; line: string } | undefined {
  if (isFailure(text)) return undefined
  const status = statusLineOf(text)
  const ticket = status?.ticket ?? ticketOf(text)
  if (ticket === undefined) return undefined
  const parts = status
    ? [
        ticket,
        `${status.from} → ${status.to}`,
        status.verdict,
        status.branch,
        status.commits === undefined ? undefined : `${words.commits} ${status.commits}`,
        status.gates === undefined ? undefined : `${words.gates} ${status.gates}`,
        status.notes === undefined ? undefined : `${words.notes} ${status.notes}`,
      ].map(p => (p === undefined ? undefined : clip(p, 60)))
    : [ticket]
  // A question inside the status line's verdict is shown twice, once there
  // and once here: a known limit, chosen over a dedupe that lost questions.
  const question = questionOf(text)
  if (question !== undefined) parts.push(`${words.question}: ${question}`)
  const needs = needsYouOf(text)
  if (needs !== undefined) parts.push(`${words.needsYou}: ${needs}`)
  const line = parts.filter(p => p !== undefined).join(' · ')
  return { ticket, line: `${words.marker} ${line} — ${words.fullReport}: /cockpit report ${ticket}` }
}

function clip(s: string, n: number): string {
  return s.length <= n ? s : `${s.slice(0, n - 1)}…`
}

// The first line of a kept report worth reading: a task notification's
// result, else the text, past blank lines and lines that are nothing but
// tags and whole elements (`<task-id>a1</task-id> <status>done</status>`).
// A line of content keeps every angle bracket it has: `Map<K, V>` or
// `n<3 then x>2` is what the agent wrote, not markup. The line is cut into
// its tags first and each piece matched alone, with patterns that cannot
// match one text two ways: an agent writes the report, and a pattern that
// backtracks on it would hang the session at every read of the tab.
const TAG = /^<[^<>]*>$/
const ELEMENT = /^<([a-z][\w-]*)(?:\s[^<>]*)?>[^<>]*<\/\1>$/i
const tagsOnly = (line: string) => line.replace(/>\s*</g, '>\n<').split('\n').every(p => TAG.test(p) || ELEMENT.test(p))
export function reportLine(text: string): string {
  const body = /<result>([\s\S]*?)(?:<\/result>|$)/.exec(text)?.[1] ?? text
  return body.split('\n').map(l => l.replace(/\s+/g, ' ').trim()).find(l => l !== '' && !tagsOnly(l)) ?? ''
}
