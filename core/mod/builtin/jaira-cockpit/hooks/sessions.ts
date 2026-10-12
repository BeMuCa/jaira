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
  ticket: string
  role: Role
  lane: string
  lastAnswerAt?: number
  tokens: number
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

// Past this the old worker's context costs more than a fresh one reading
// the ticket again.
export const CONTEXT_LIMIT = 150_000

export type Choice = { reuse: boolean; why: string }

// Whom a ticket given back goes to: the worker that carried it while its
// tab lives, its cache is warm and its context is small; otherwise a new one.
// The reason is written onto the ticket, so it says what was measured.
export function chooseWorker(input: { record?: SessionRecord; tabAlive: boolean; now: number }): Choice {
  const { record, tabAlive, now } = input
  if (!tabAlive) return { reuse: false, why: 'the old worker has no Herdr tab open any more' }
  if (record === undefined) return { reuse: false, why: 'nothing recorded about the old worker session (cache and context unknown)' }
  const k = Math.round(record.tokens / 1000)
  if (record.tokens >= CONTEXT_LIMIT) return { reuse: false, why: `old worker's context ${k}k >= ${CONTEXT_LIMIT / 1000}k` }
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
  return tabsOf(tabListJson).find(t => t.label === `${ticket}/worker` || t.label === `${ticket}/dispatch`)
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
