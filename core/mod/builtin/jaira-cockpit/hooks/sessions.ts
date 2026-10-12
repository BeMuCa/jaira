// Every session the cockpit runs in writes down where it stands — which
// ticket it carries, its Herdr pane, when the model last answered, how big
// its context is — so the session a person works from can tell whether the
// worker that carried a ticket is still worth sending work back to. The
// records live in `$.store`, which every session on the machine shares, under
// `session:<id>`. Pure; register.tsx writes and reads them.

export type Role = 'worker' | 'dispatch' | 'lane'

export const RECORD_PREFIX = 'session:'

// A record a day old describes a worker whose cache is long cold: it decides
// nothing any more and goes.
export const RECORD_MAX_AGE_MS = 24 * 60 * 60_000

export type SessionRecord = {
  session: string
  pane: string
  cwd: string
  // The repository's main working tree, shared by every worktree of it: the
  // board the record belongs to. Missing in a record of an older build.
  repo?: string
  ticket: string
  role: Role
  lane: string
  lastAnswerAt?: number
  tokens: number
  // The model's context window; missing in a record of an older build.
  window?: number
  ttl: number
  updatedAt: number
}

// The line spawn.sh types to start an agent on a ticket names the ticket:
// `/jaira-worker <id>`, `/jaira-dispatcher <id>`, `/jaira-role-lane <id> <lane>`.
const START = /^\s*\/jaira-(worker|dispatcher|role-lane)\s+([0-9A-Za-z]{6})\b(?:\s+([a-z][a-z-]*))?/

export function startOf(prompt: string): { ticket: string; role: Role; lane: string } | undefined {
  const m = START.exec(prompt)
  if (m === null) return undefined
  const role: Role = m[1] === 'worker' ? 'worker' : m[1] === 'dispatcher' ? 'dispatch' : 'lane'
  return { ticket: m[2]!.toUpperCase(), role, lane: role === 'lane' ? (m[3] ?? '') : '' }
}

// The newest record of the agent that carries `ticket` as a whole: a worker
// or a dispatcher, never a single lane's agent.
export function carrierOf(records: readonly SessionRecord[], ticket: string): SessionRecord | undefined {
  return records
    .filter(r => r.ticket === ticket && r.role !== 'lane')
    .sort((a, b) => b.updatedAt - a.updatedAt)[0]
}

// Past this share of its model's context window, in percent, an old worker's
// context costs more than a fresh one reading the ticket again. A setting:
// userConfig workerContextShare.
export const DEFAULT_CONTEXT_SHARE = 50

// A record that does not say its window counts as the smallest one in use:
// rather a new worker too many than one carried on too full.
export const UNKNOWN_WINDOW = 200_000

export function contextLimit(record: SessionRecord, sharePercent: number): number {
  return Math.round(((record.window ?? UNKNOWN_WINDOW) * sharePercent) / 100)
}

export type Choice = { reuse: boolean; why: string }

// Whom a ticket coming back for a new round goes to: the worker that carried
// it while its tab lives, its cache is warm and its context under the limit;
// otherwise a new one. Give back asks it, and the teamlead reads its answer
// off `workers` — one rule for both. The reason is written onto the ticket,
// so it says what was measured.
export function chooseWorker(input: { record?: SessionRecord; tabAlive: boolean; now: number; share: number }): Choice {
  const { record, tabAlive, now, share } = input
  if (!tabAlive) return { reuse: false, why: 'the old worker has no Herdr tab open any more' }
  if (record === undefined) return { reuse: false, why: 'nothing recorded about the old worker session (cache and context unknown)' }
  const k = Math.round(record.tokens / 1000)
  const limit = contextLimit(record, share)
  if (record.tokens >= limit) return { reuse: false, why: `old worker's context ${k}k >= ${Math.round(limit / 1000)}k (${share}% of its window)` }
  const left = record.lastAnswerAt === undefined ? -1 : record.lastAnswerAt + record.ttl - now
  if (left <= 0) {
    const cold = record.lastAnswerAt === undefined ? '' : ` for ${Math.max(1, Math.round(-left / 60_000))} min`
    return { reuse: false, why: `old worker's prompt cache cold${cold}` }
  }
  return { reuse: true, why: `tab alive, cache warm for ${Math.ceil(left / 60_000)} more min, context ${k}k` }
}

// The Herdr tab of the agent carrying a ticket, from `herdr tab list`:
// spawn.sh labels it `<id>/worker` or `<id>/dispatch`. Herdr says whether an
// agent works in it: `working`, `idle` or `done`.
export type Tab = { tab_id: string; label: string; agent_status?: string }

export function tabsOf(tabListJson: string): Tab[] {
  try {
    return (JSON.parse(tabListJson) as { result?: { tabs?: Tab[] } }).result?.tabs ?? []
  } catch {
    return []
  }
}

export function carrierTab(tabListJson: string, ticket: string): Tab | undefined {
  return carrierTabIn(tabsOf(tabListJson), ticket)
}

function carrierTabIn(tabs: readonly Tab[], ticket: string): Tab | undefined {
  return tabs.find(t => t.label === `${ticket}/worker` || t.label === `${ticket}/dispatch`)
}

const size = (tokens: number, window: number | undefined) =>
  `context ${Math.round(tokens / 1000)}k${window ? ` of ${Math.round(window / 1000)}k` : ''}`

// The session asking, from its live figures: in the middle of its own turn
// its cache is warm, and stays so for one TTL after the turn's answer.
export function selfLine(tokens: number, window: number | undefined, ttl: number): string {
  return `this session: ${size(tokens, window)}, cache warm while it works and ${Math.round(ttl / 60_000)} min after its last answer`
}

// Every agent carrying a ticket of this board — a worker or a dispatcher —
// with its ticket's lane, its context, its cache, and whom a new round of
// that ticket would go to by chooseWorker. With Herdr (`tabs` given) only
// those whose tab is open: a closed tab carries nothing any more. One at work
// is never judged: nothing is decided about a worker while it works.
export function workerLines(input: {
  records: readonly SessionRecord[]
  repo: string
  tabs?: readonly Tab[]
  lanes: ReadonlyMap<string, string>
  now: number
  share: number
}): string[] {
  const { records, repo, tabs, lanes, now, share } = input
  const ours = records.filter(r => r.repo === repo && r.role !== 'lane')
  const tickets = new Set(ours.map(r => r.ticket))
  for (const tab of tabs ?? []) {
    const [id = '', kind] = tab.label.split('/')
    if ((kind === 'worker' || kind === 'dispatch') && lanes.has(id)) tickets.add(id)
  }
  const lines: string[] = []
  for (const ticket of [...tickets].sort()) {
    const tab = tabs === undefined ? undefined : carrierTabIn(tabs, ticket)
    if (tabs !== undefined && tab === undefined) continue
    const record = carrierOf(ours, ticket)
    const head = `${ticket} · ${lanes.get(ticket) ?? '?'} · ${record?.role ?? tab?.label.split('/')[1] ?? '?'}`
    if (record === undefined) {
      lines.push(`${head} · nothing recorded (context and cache unknown)`)
      continue
    }
    const left = record.lastAnswerAt === undefined ? 0 : record.lastAnswerAt + record.ttl - now
    const cache = tab?.agent_status === 'working' ? 'cache warm (working)' : left > 0 ? `cache ${Math.ceil(left / 60_000)} min left` : 'cache cold'
    const next = tab?.agent_status === 'working'
      ? 'working: nothing to decide'
      : (c => `next round: ${c.reuse ? 'same worker' : 'new worker'} (${c.why})`)(chooseWorker({ record, tabAlive: tab !== undefined, now, share }))
    lines.push(`${head} · ${size(record.tokens, record.window)} · ${cache} · ${next}`)
  }
  return lines
}

// What the cockpit types into a pane leads with this, so the cockpit in that
// session knows the words are its own and no person's: no cache guard, no
// thread for them.
export const TYPED_BY_COCKPIT = '[jaira-cockpit]'

// What the old worker is told when work comes back to it.
export function backMessage(ticket: string, lane: string, reason: string): string {
  return `${TYPED_BY_COCKPIT} Ticket ${ticket} came back from ${lane} to in-progress. Reason: ${reason} - also a note on the ticket. Carry it on as before (jaira show ${ticket} --json).`
}

export function isExpired(record: SessionRecord, now: number): boolean {
  return now - record.updatedAt > RECORD_MAX_AGE_MS
}
