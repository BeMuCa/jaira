import type { EngineInterface, Register } from 'claude-code'
import { digest } from './digest.ts'
import {
  agentName, clipReport, DIGESTED, isAgentReport, isUnscoped, keyOf, listReports, NOTIFICATION, projectKeys,
  REPORT_USAGE, REPORTS_PER_TICKET, rewriteRow, showReport, TICKETS, AGENT_ORIGINS, type Report,
} from './filter.ts'
import { source } from './frame.ts'
import { langOf, type Lang } from './lang.ts'

// The engine follows `$` only into functions declared in this same file, and
// a module hooks each event once: so every call on `$` is here, one hook per
// event, and the other files are rules and words that take plain values.

// Mod commands share one namespace with every other mod's, so the cockpit
// takes one name and hangs its subcommands under it.
const USAGE = `Usage: /cockpit <subcommand>\n\n${REPORT_USAGE}`

async function language($: EngineInterface): Promise<Lang> {
  try {
    return langOf((await $.settings.read()).language)
  } catch {
    return 'en'
  }
}

// The project a report belongs to: the repository's main working tree, so
// every worktree of one repository shares its reports; else the directory.
async function project($: EngineInterface): Promise<string> {
  try {
    const repo = await $.session.repo()
    if (repo !== null) return repo.root
  } catch {}
  try {
    return await $.session.cwd()
  } catch {
    return ''
  }
}

async function sessionId($: EngineInterface): Promise<string> {
  try {
    return (await $.session.id()).slice(0, 8)
  } catch {
    return ''
  }
}

// ── the noise filter ─────────────────────────────────────────────────────

async function keep($: EngineInterface, ticket: string, report: Report): Promise<void> {
  const root = await project($)
  const key = keyOf(root, ticket)
  const before = ((await $.store.get(key)) as Report[] | undefined) ?? []
  // Delete first so the key moves to the end of the store's insertion order:
  // the first keys are then the tickets heard from longest ago.
  await $.store.delete(key)
  await $.store.set(key, [{ ...report, text: clipReport(report.text) }, ...before].slice(0, REPORTS_PER_TICKET))
  const tickets = projectKeys(await $.store.keys(), root)
  for (const old of tickets.slice(0, Math.max(0, tickets.length - TICKETS))) await $.store.delete(old)
}

// The text the conversation gets instead, with the report kept; undefined
// when it passes on unchanged.
async function shorten($: EngineInterface, text: string, origin: string, agent?: string): Promise<string | undefined> {
  if (DIGESTED.test(text)) return undefined
  const report = NOTIFICATION.exec(text)?.[0] ?? text
  const short = digest(report, await language($))
  if (short === undefined) return undefined
  await keep($, short.ticket, {
    at: new Date().toISOString(),
    origin: source(report, origin, agent),
    session: await sessionId($),
    text: report,
  })
  return text.replace(report, () => short.line)
}

async function reports($: EngineInterface, args: string, lang: Lang): Promise<string> {
  const root = await project($)
  const ticket = args.trim().toUpperCase()
  if (ticket === '') {
    const keys = projectKeys(await $.store.keys(), root).reverse()
    const rows = await Promise.all(keys.map(async k => [k, ((await $.store.get(k)) as Report[] | undefined) ?? []] as const))
    return listReports(rows, lang)
  }
  return showReport(ticket, ((await $.store.get(keyOf(root, ticket))) as Report[] | undefined) ?? [], lang)
}

// What session.receive left alone on purpose — a comment on an artifact, a
// pull request's event — reaches prompt.submit next, under the same origin a
// report has; it must not be shortened there either.
const untouched = new Set<string>()

// What prompt.submit shortened, by the text it was handed, until the row
// carrying it is appended. A subagent's SendMessage to the main conversation
// reaches prompt.submit alone (origin peer, no session.receive), and the
// engine stores that row with the text as it arrived, wrapped in its own
// lead-in, whatever the hook answered: the row the model reads and the chat
// shows is rewritten at session.append too.
const shortened = new Map<string, string>()

export const register: Register = on => {
  on('session.start', async ($, e, next) => {
    const started = await next(e)
    // Kept before reports were scoped to a project; no project may show them.
    for (const k of await $.store.keys()) if (isUnscoped(k)) await $.store.delete(k)
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

  on('session.receive', async ($, e, next) => {
    if (!isAgentReport(e)) {
      untouched.add(e.text)
      return next(e)
    }
    const text = await shorten($, e.text, e.origin.kind)
    return next(text === undefined ? e : { ...e, text })
  }).catch(($, e, next) => next(e))

  on('prompt.submit', async ($, e, next) => {
    if (untouched.delete(e.text) || !AGENT_ORIGINS.has(e.origin.kind)) return next(e)
    // A background task of this session always reports in its own block.
    if (e.origin.kind === 'task-notification' && !NOTIFICATION.test(e.text)) return next(e)
    const text = await shorten($, e.text, e.origin.kind, agentName(e))
    if (text === undefined) return next(e)
    shortened.set(e.text, text)
    // A row that never came keeps nothing alive for long.
    for (const old of [...shortened.keys()].slice(0, Math.max(0, shortened.size - 20))) shortened.delete(old)
    return next({ ...e, text })
  }).catch(($, e, next) => next(e))

  on('session.append', ($, e, next) => next(rewriteRow(e, shortened))).catch(($, e, next) => next(e))

  on('command.run', { command: 'cockpit' }, async ($, e) => {
    const [sub = '', ...rest] = (e.args ?? '').trim().split(/\s+/)
    if (sub === 'report') return { text: await reports($, rest.join(' '), await language($)) }
    return { text: USAGE }
  })

  on('tool.call', { tool: 'mcp__jaira-cockpit__report' }, async ($, e) => {
    // A tool's arguments ride on `e` beside `tool`, not under a key of their own.
    const { ticket } = e as { ticket?: unknown }
    return { result: await reports($, typeof ticket === 'string' ? ticket : '', 'en') }
  })
}
