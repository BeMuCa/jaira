// What the cockpit knows about a jaira board it reads off the board itself,
// through the jaira binary: `jaira lanes --json` says whose lane is whose,
// `jaira list --json` where the tickets stand. Not from the agents' reports:
// so it is the same in both dispatch modes and with Herdr tabs. register.tsx
// runs the commands; this file only reads what they printed.

import type { Tab } from './sessions.ts'

export type Waiting = {
  id: string
  title: string
  lane: string
  verdict: string
  question: string
  // Where Accept moves it: the ticket's next lane, when that lane is the
  // board's terminal one (signoff -> done); null where the next lane is
  // still work to do (human -> review), so a press never skips a review.
  acceptTo: string | null
  // Whose it is; '' for nobody's.
  assignee: string
  // Whether the band offers Explain: the verdict names a reservation, or
  // the review names a gap.
  explainable: boolean
}

export type Counts = { worker: number; review: number; person: number }

type Lane = { id: string; agentic: boolean; terminal?: boolean; context: string; human_exit?: boolean }
type Ticket = {
  handle: string
  title: string
  status: string
  next_lane?: string
  question?: string
  assignee?: string
  review?: { verdict?: string; gaps?: string }
  // The ticket's round, and the round its verdict was written in (0 or
  // absent: unknown, a verdict from before jaira counted rounds).
  round?: number
  review_round?: number
}

function parse<T>(json: string): T | undefined {
  try {
    return JSON.parse(json) as T
  } catch {
    return undefined
  }
}

const oneLine = (s: string | undefined) => (s ?? '').replace(/\s+/g, ' ').trim()

// A jaira built before `lanes --json` carried human_exit says nothing about
// whose lane is whose: then the shipped board's person lanes are taken.
const PERSON_LANES = new Set(['human', 'signoff'])

function isPersons(lane: Lane): boolean {
  return lane.human_exit ?? PERSON_LANES.has(lane.id)
}

export type Agent = { description: string; name?: string; status: string }

// The chats a person opens from the band are theirs, not an agent at work.
const PERSONS_TABS = new Set(['explain', 'discuss'])
const WORKING = new Set(['pending', 'running', 'waiting'])

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
    into.push({
      id: t.handle,
      title: t.title,
      lane: t.status,
      verdict: oneLine(t.review?.verdict),
      question: oneLine(t.question),
      acceptTo: next?.terminal === true ? next.id : null,
      assignee,
      explainable: explainable(oneLine(t.review?.verdict), t.review?.gaps, sameRound(t)),
    })
  }
  const ids = tickets.map(t => t.handle)
  const fresh = lanes.filter(l => l.agentic && l.context === 'fresh').map(l => l.id)
  // A worker's tab stays open after its turn, for Give back: only one whose
  // agent is working counts.
  for (const { label, agent_status } of tabs) {
    const [id = '', lane = ''] = label.split('/')
    if (agent_status !== 'working' || !ids.includes(id) || lane === '' || PERSONS_TABS.has(lane)) continue
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

// A review's gaps name one unless they are empty or say there is none, alone
// or before a remark: 'none', 'none — the docs could say more', 'keine
// (blockierend) ...', 'No gaps found.', 'keine Lücken'. A verdict that
// reads as a yes ("the one weak spot ... (gap 1)") still has something a
// person should have explained.
const NO_GAP = /^(none|keine|nothing|nichts|n\/a|-+|no (gaps?|defects?|issues?|findings?)( found)?|keine (lücken?|mängel|befunde))\s*([—–\-:(.,;!]|$)/i

function namesGap(gaps: string | undefined): boolean {
  const g = oneLine(gaps)
  return g !== '' && !NO_GAP.test(g)
}

// Whether the verdict was written in the ticket's current round. A later
// round rewrites review-gaps in its loop lanes while the old verdict stands,
// so the two belong together only then. A verdict whose round jaira did not
// record is taken as current, as it was before rounds were counted.
function sameRound(t: Ticket): boolean {
  return !t.review_round || t.review_round === (t.round ?? 1)
}

// Whether the band offers Explain. The gaps count only beside a verdict of
// the same round: before the review, the field holds what the optimize lane
// removed, not a gap — whatever the lane is called.
function explainable(verdict: string, gaps: string | undefined, current: boolean): boolean {
  return hasReservation(verdict) || (verdict !== '' && current && namesGap(gaps))
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
export function noteOf(showJson: string, kind: 'explain' | 'decision'): { text: string; at: number } | undefined {
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
  return { text, at: Date.parse(`${stamp.replace(' ', 'T')}:00Z`) }
}

// An Explain chat's note, in the shape spawn.sh asks for: one line per part —
// the finding, why it counts, one line per option as `a) ...`, ending in
// `[give back]` or `[accept]` when picking it is that step, and
// `Empfehlung: <letter> - ...`. Its parts in the order written, the options
// apart for the choice with what each does, and the option recommended. A
// note with no option line — one written before the shape was asked for —
// is one part, whole, as it was always shown.
export type Action = 'giveBack' | 'accept'
export type Option = { key: string; text: string; action?: Action }
export type Part = { kind: 'text'; text: string } | ({ kind: 'option' } & Option)
export type Explanation = { parts: Part[]; options: Option[]; recommended?: string }

const OPTION = /^\s*\(?([a-f])\)\s+(.+)$/
const MARK = /\s*\[(give back|accept)\]\s*$/i
const RECOMMENDED = /^\s*(?:empfehlung|recommendation)\s*:\s*\(?([a-f])\b/i

export function explanationOf(text: string): Explanation {
  const lines = text.split('\n').map(l => l.trim()).filter(l => l !== '')
  const parts: Part[] = lines.map(l => {
    const m = OPTION.exec(l)
    if (!m) return { kind: 'text', text: l }
    const mark = MARK.exec(m[2]!)?.[1]?.toLowerCase()
    const text = m[2]!.replace(MARK, '').trim()
    return { kind: 'option', key: m[1]!, text, ...(mark ? { action: mark === 'accept' ? 'accept' : 'giveBack' } : {}) }
  })
  const options = parts.flatMap(({ kind, ...o }) => (kind === 'option' ? [o as Option] : []))
  if (options.length === 0) return { parts: [{ kind: 'text', text: text.trim() }], options }
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
// your disk".
export function worktreeFor(porcelain: string, id: string, fallback: string): string {
  const slug = id.toLowerCase()
  for (const block of porcelain.split(/\n\n+/)) {
    const path = /^worktree (.+)$/m.exec(block)?.[1]
    const branch = /^branch (.+)$/m.exec(block)?.[1] ?? ''
    if (path !== undefined && (branch.endsWith(`/${slug}`) || path.endsWith(`-${slug}`))) return path
  }
  return fallback
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
