// digest.ts reads one agent report and boils it down to the line the
// teamlead chat shows in its place. Everything here is pure: no `$`, so the
// rules can be read and tested without an engine.

// The lanes jaira ships. A board may rename or add lanes; an arrow between
// two lowercase words (`review -> signoff`) is read whatever the names are.
const LANES = [
  'backlog', 'brainstorm', 'todo', 'pre-process', 'in-progress', 'critique',
  'optimize', 'testing', 'human', 'review', 'signoff', 'done', 'blocked',
]

// A handle is the last six characters of a ticket's ULID, Crockford base32
// (no I, L, O, U). Six capital letters alone are too often an ordinary word
// (STATUS, PASSED), so a handle without a digit counts only where jaira's
// own wording puts one: after "ticket", in a commit subject, after a jaira
// verb.
const HANDLE = '[0-9A-HJKMNP-TV-Z]{6}'
const ULID = /\b[0-9A-HJKMNP-TV-Z]{20}([0-9A-HJKMNP-TV-Z]{6})\b/g
const WITH_DIGIT = new RegExp(`\\b(?=[A-Z]*[0-9])${HANDLE}\\b`, 'g')
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

const ARROW = /\b([a-z][a-z-]*)\s*(?:->|→|=>)\s*([a-z][a-z-]*)\b/
const INTO = new RegExp(
  `\\b(?:to|into|in|nach|in die lane|sits in|steht in|liegt in)\\s+\`?(${LANES.join('|')})\\b`,
  'i',
)

// The lane change, "from → to", or the lane the ticket ended in alone.
export function laneOf(text: string): string | undefined {
  for (const line of text.split('\n')) {
    const m = ARROW.exec(line)
    if (m && (LANES.includes(m[1]!) || LANES.includes(m[2]!))) return `${m[1]} → ${m[2]}`
  }
  const m = INTO.exec(text)
  return m ? `→ ${m[1]!.toLowerCase()}` : undefined
}

const VERDICT_LINE = /(?:review-verdict\s*=\s*|\b(?:verdict|urteil)\s*[:=]\s*)["“]?([^"\n”]+)/i
const VERDICT_WORD =
  /\b(approve[sd]?(?: with (?:reservations?|caveats?)| mit vorbehalt)?|request(?:s|ed)? changes|changes requested|reject(?:ed)?|sent back|zur(?:ü|ue)ckgegeben|mit vorbehalt)\b/i

export function verdictOf(text: string): string | undefined {
  const line = VERDICT_LINE.exec(text)
  if (line) return clip(line[1]!.trim(), 60)
  const word = VERDICT_WORD.exec(text)
  return word ? word[1]!.toLowerCase() : undefined
}

const QUESTION = /(?:--question\s+["“]|\b(?:open question|offene frage|question|frage)\s*[:=]\s*["“]?)([^"\n”]+)/i
const NO_QUESTION = /^(?:none|no|keine|nein|nichts|-|—)\b/i

export function questionOf(text: string): string | undefined {
  const m = QUESTION.exec(text)
  if (!m) return undefined
  const q = m[1]!.trim()
  return q === '' || NO_QUESTION.test(q) ? undefined : clip(q, 100)
}

// The line the chat shows instead of the report, or undefined when the
// report must be passed on whole: no ticket recognisable, or a failure.
export function digest(text: string): { ticket: string; line: string } | undefined {
  if (isFailure(text)) return undefined
  const ticket = ticketOf(text)
  if (ticket === undefined) return undefined
  const parts = [ticket, laneOf(text), verdictOf(text)]
  const question = questionOf(text)
  if (question !== undefined) parts.push(`question: ${question}`)
  const line = parts.filter(p => p !== undefined).join(' · ')
  return { ticket, line: `${line} — full report: /flight ${ticket}` }
}

function clip(s: string, n: number): string {
  return s.length <= n ? s : `${s.slice(0, n - 1)}…`
}
