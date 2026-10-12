import type { EngineInterface, Register, SessionReceiveInput } from 'claude-code'
import { digest } from './digest.ts'

// The noise filter: an agent's report or a task-notification that names a
// ticket reaches the conversation as one line, and its full text is kept in
// the store, where `/cockpit report <id>` (and the model's `report` tool)
// reads it back. Reports with no recognisable ticket, failures and anything a
// person wrote pass through unchanged.

// What a report is kept as, newest first under `report:<ticket>`.
type Report = { at: string; origin: string; text: string }

const PREFIX = 'report:'
// The store holds 4 MiB of JSON in all: 20 tickets of 10 reports of 12,000
// characters stay well below it, escaping included.
const REPORTS_PER_TICKET = 10
const TICKETS = 20
const REPORT_CHARS = 12_000

// Deliveries an agent sent: a background task's notification, another
// session's or agent's message. They reach the session two ways: a relayed
// one through `session.receive`, a local one — this session's own background
// task finishing — through `prompt.submit` alone, so both are hooked. A
// person's own message (typed, or through Remote Control as `bridge`), a
// Slack ping, a scheduled prompt and anything unclassified are not reports
// and are never touched.
const AGENT_ORIGINS = new Set(['task-notification', 'peer', 'peer-send-message', 'coordinator'])

function isAgentReport(e: SessionReceiveInput): boolean {
  // An agentId means the delivery is for one of this session's agents, not
  // for the conversation the person reads. An event is a relay's wake — a
  // pull request, a comment on an artifact — written by people.
  return e.agentId === undefined && e.event === undefined && AGENT_ORIGINS.has(e.origin.kind)
}

async function keep($: EngineInterface, ticket: string, report: Report): Promise<void> {
  const key = PREFIX + ticket
  const before = ((await $.store.get(key)) as Report[] | undefined) ?? []
  const text = report.text.length <= REPORT_CHARS
    ? report.text
    : `${report.text.slice(0, REPORT_CHARS)}\n[… cut at ${REPORT_CHARS} characters]`
  // Delete first so the key moves to the end of the store's insertion order:
  // the first keys are then the tickets heard from longest ago.
  await $.store.delete(key)
  await $.store.set(key, [{ ...report, text }, ...before].slice(0, REPORTS_PER_TICKET))
  const tickets = (await $.store.keys()).filter(k => k.startsWith(PREFIX))
  for (const old of tickets.slice(0, Math.max(0, tickets.length - TICKETS))) await $.store.delete(old)
}

// A delivery already shortened at `session.receive` reaches `prompt.submit`
// again; it is not a report a second time.
const DIGESTED = /— full report: \/cockpit report [0-9A-Z]{6}$/
// The notification itself, without the warning the engine wraps it in.
const NOTIFICATION = /<task-notification>[\s\S]*?<\/task-notification>/

// The text the conversation gets instead, with the report kept; undefined
// when it passes on unchanged.
async function shorten($: EngineInterface, text: string, origin: string): Promise<string | undefined> {
  if (DIGESTED.test(text)) return undefined
  const report = NOTIFICATION.exec(text)?.[0] ?? text
  const short = digest(report)
  if (short === undefined) return undefined
  await keep($, short.ticket, { at: new Date().toISOString(), origin, text: report })
  return text.replace(report, () => short.line)
}

async function showReports($: EngineInterface, args: string): Promise<string> {
  const ticket = args.trim().toUpperCase()
  if (ticket === '') {
    const tickets = (await $.store.keys()).filter(k => k.startsWith(PREFIX)).reverse()
    if (tickets.length === 0) return 'No agent reports kept yet.'
    const rows = await Promise.all(tickets.map(async k => {
      const reports = ((await $.store.get(k)) as Report[] | undefined) ?? []
      return `${k.slice(PREFIX.length)}  ${reports.length} report(s), last ${reports[0]?.at ?? '?'}`
    }))
    return `Tickets with kept reports, most recent first:\n${rows.join('\n')}\n\n/cockpit report <id> shows one.`
  }
  const reports = ((await $.store.get(PREFIX + ticket)) as Report[] | undefined) ?? []
  if (reports.length === 0) return `No reports kept for ${ticket}.`
  return reports
    .map((r, i) => `── ${ticket} · report ${i + 1} of ${reports.length} · ${r.at} · ${r.origin} ──\n${r.text}`)
    .join('\n\n')
}

// Mod commands share one namespace with every other mod's, so the cockpit
// takes one name and hangs its subcommands under it.
const USAGE = 'Usage: /cockpit <subcommand>\n\n  report [id]   full text of the agent reports shortened for a ticket, newest first; without an id, the tickets that have some'

export const register: Register = on => {
  on('session.start', async ($, e, next) => {
    const started = await next(e)
    await $.command.register({
      name: 'cockpit',
      description: 'The jaira cockpit. report [id]: full text of the agent reports it shortened, newest first.',
      argumentHint: 'report [ticket id]',
      immediate: true,
    })
    await $.tool.register({
      name: 'report',
      description:
        'Returns the full text of the agent reports and task-notifications about one jaira ticket that reached this conversation shortened to one line, newest first. Without a ticket, lists the tickets that have kept reports.',
      inputSchema: {
        type: 'object',
        properties: { ticket: { type: 'string', description: 'The ticket handle, e.g. MG677R.' } },
      },
    })
    return started
  })

  // What session.receive left alone on purpose — a comment on an artifact,
  // a pull request's event — reaches prompt.submit next, under the same
  // origin a report has; it must not be shortened there either.
  const untouched = new Set<string>()

  on('session.receive', async ($, e, next) => {
    if (!isAgentReport(e)) {
      untouched.add(e.text)
      return next(e)
    }
    const text = await shorten($, e.text, e.origin.kind)
    return next(text === undefined ? e : { ...e, text })
  }).catch(($, e, next) => next(e))

  // What prompt.submit shortened, by the text it was handed, until the row
  // carrying it is appended. A subagent's SendMessage to the main
  // conversation reaches prompt.submit alone (origin peer, no
  // session.receive), and the engine stores that row with the text as it
  // arrived, wrapped in its own lead-in, whatever the hook answered: the row
  // the model reads and the chat shows is rewritten at session.append too.
  const shortened = new Map<string, string>()

  on('prompt.submit', async ($, e, next) => {
    if (untouched.delete(e.text) || !AGENT_ORIGINS.has(e.origin.kind)) return next(e)
    // A background task of this session always reports in its own block.
    if (e.origin.kind === 'task-notification' && !NOTIFICATION.test(e.text)) return next(e)
    const text = await shorten($, e.text, e.origin.kind)
    if (text === undefined) return next(e)
    shortened.set(e.text, text)
    // A row that never came keeps nothing alive for long.
    for (const old of [...shortened.keys()].slice(0, Math.max(0, shortened.size - 20))) shortened.delete(old)
    return next({ ...e, text })
  }).catch(($, e, next) => next(e))

  on('session.append', ($, e, next) => {
    const isDelivery = e.door === 'prompt' || e.door === 'delivery'
    if (!isDelivery || e.agentId !== undefined || e.message.type !== 'user' || !AGENT_ORIGINS.has(e.origin.kind)) return next(e)
    let isChanged = false
    const content = e.message.content.map(block => {
      if (block.type !== 'text') return block
      for (const [original, short] of shortened) {
        if (!block.text.includes(original)) continue
        shortened.delete(original)
        isChanged = true
        return { ...block, text: block.text.replace(original, () => short) }
      }
      return block
    })
    return next(isChanged ? { ...e, message: { ...e.message, content } } : e)
  }).catch(($, e, next) => next(e))

  on('command.run', { command: 'cockpit' }, async ($, e) => {
    const [sub = '', ...rest] = (e.args ?? '').trim().split(/\s+/)
    if (sub === 'report') return { text: await showReports($, rest.join(' ')) }
    return { text: USAGE }
  })

  on('tool.call', { tool: 'mcp__jaira-cockpit__report' }, async ($, e) => {
    // A tool's arguments ride on `e` beside `tool`, not under a key of their own.
    const { ticket } = e as { ticket?: unknown }
    return { result: await showReports($, typeof ticket === 'string' ? ticket : '') }
  })
}
