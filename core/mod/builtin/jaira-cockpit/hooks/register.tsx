import { atom, read, update } from 'claude-code'
import type { EngineInterface, Register } from 'claude-code'
import type { BoardView, CacheState, ThreadTally } from '../types'
import { hasReservation, mainRepoOf, otherRefusal, snapshotOf, upwards, worktreeFor, type Waiting } from './board.ts'
import { cacheView, dollars, guard, isPersonsPrompt, ttlMs } from './cache.ts'
import { digest } from './digest.ts'
import {
  agentName, AGENT_ORIGINS, clipReport, DIGESTED, isAgentReport, isUnscoped, keyOf, listReports, NOTIFICATION,
  projectKeys, REPORTS_PER_TICKET, rewriteRow, showReport, stamp, TICKETS, type Report,
} from './filter.ts'
import { body, source } from './frame.ts'
import { langOf, type Lang } from './lang.ts'
import { add, appended, emptyTotals, row, rowsOf, type Totals } from './metrics.ts'
import {
  backMessage, carrierOf, carrierTab, chooseWorker, isExpired, RECORD_PREFIX, startOf, type Role, type SessionRecord,
} from './sessions.ts'
import { texts } from './texts.ts'
import { classifyPrompt, file, listThreads, parseFiling, tally, type Thread } from './threads.ts'

// The engine follows `$` only into functions declared in this same file, and
// a module hooks each event once: so every call on `$` is here, one hook per
// event, and the other files are rules and words that take plain values.
//
// The parts: the noise filter (agent reports as one line, /cockpit report);
// the band above the prompt (tickets waiting for a person with Accept / Give
// back / Explain / Discuss, context and cache, conversation threads); the
// cache guard; the status line; the pane per ticket; metrics per session.
// Whatever needs a board — a .jaira/ and jaira on PATH — stays out of sight
// without one; the rest works in any project.

const boardAtom = atom({ plugin: 'jaira-cockpit', key: 'board' } as const, null)
const cacheAtom = atom({ plugin: 'jaira-cockpit', key: 'cache' } as const, { tokens: 0, lastAnswerAt: null, ttl: 0 })
const threadsAtom = atom({ plugin: 'jaira-cockpit', key: 'threads' } as const, { questions: 0, stale: 0 })
const givingAtom = atom({ plugin: 'jaira-cockpit', key: 'giving' } as const, null)
const paneAtom = atom({ plugin: 'jaira-cockpit', key: 'paneTicket' } as const, null)
const langAtom = atom({ plugin: 'jaira-cockpit', key: 'lang' } as const, 'en')
const minuteAtom = atom({ plugin: 'jaira-cockpit', key: 'minute' } as const, 0)
const noticeAtom = atom({ plugin: 'jaira-cockpit', key: 'notice' } as const, null)

const PANE = 'cockpit-ticket'
const BAND_TICKETS = 3

// ── what this module keeps between events (a reload starts it over) ──────

// The board this session's directory belongs to: undefined until looked
// for, null when there is none.
let board: { dir: string; repo: string } | null | undefined
// The ticket this session carries, read off the line that started it.
let carrying: { ticket: string; role: Role; lane: string } | undefined
let lastAnswerAt: number | undefined
let startedAt = 0
let totals: Totals = emptyTotals()
let turn = 0
// The person's message of the turn running, for the threads.
let personSaid: string | undefined
let threads: Thread[] = []

// ── shared helpers ───────────────────────────────────────────────────────

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

async function run($: EngineInterface, argv: string[], cwd: string, timeoutMs = 15_000) {
  try {
    return await $.process.run(argv, { cwd, timeoutMs })
  } catch {
    return undefined
  }
}

async function home($: EngineInterface): Promise<string> {
  return (await $.env.get('HOME')) ?? ''
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

async function keptReports($: EngineInterface, ticket: string): Promise<Report[]> {
  return ((await $.store.get(keyOf(await project($), ticket))) as Report[] | undefined) ?? []
}

async function reports($: EngineInterface, args: string, lang: Lang): Promise<string> {
  const root = await project($)
  const ticket = args.trim().toUpperCase()
  if (ticket === '') {
    const keys = projectKeys(await $.store.keys(), root).reverse()
    const rows = await Promise.all(keys.map(async k => [k, ((await $.store.get(k)) as Report[] | undefined) ?? []] as const))
    return listReports(rows, lang)
  }
  return showReport(ticket, await keptReports($, ticket), lang)
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

// ── the board ────────────────────────────────────────────────────────────

async function findBoard($: EngineInterface): Promise<{ dir: string; repo: string } | null> {
  let cwd = ''
  try {
    cwd = await $.session.cwd()
  } catch {
    return null
  }
  for (const dir of upwards(cwd)) {
    let has = false
    try {
      has = await $.fs.exists(`${dir}/.jaira`)
    } catch {
      return null
    }
    if (!has) continue
    // jaira itself must answer too: a board with no binary is no board here.
    const lanes = await run($, ['jaira', 'lanes', '--json'], dir)
    if (lanes === undefined || lanes.exitCode !== 0) return null
    const common = await run($, ['git', 'rev-parse', '--path-format=absolute', '--git-common-dir'], dir)
    return { dir, repo: common?.exitCode === 0 ? mainRepoOf(common.stdout, dir) : dir }
  }
  return null
}

async function refreshBoard($: EngineInterface): Promise<void> {
  if (board === undefined) board = await findBoard($)
  if (board === null) {
    await update($, boardAtom, () => null)
    return
  }
  const [lanes, list] = await Promise.all([
    run($, ['jaira', 'lanes', '--json'], board.dir),
    run($, ['jaira', 'list', '--json'], board.dir),
  ])
  if (lanes?.exitCode !== 0 || list?.exitCode !== 0) return
  const view: BoardView = snapshotOf(lanes.stdout, list.stdout)
  await update($, boardAtom, () => view)
}

// Where a ticket's work lives: running jaira anywhere else answers "on its
// ref and not on your disk".
async function ticketDir($: EngineInterface, b: { dir: string; repo: string }, id: string): Promise<string> {
  const r = await run($, ['git', 'worktree', 'list', '--porcelain'], b.dir)
  return r?.exitCode === 0 ? worktreeFor(r.stdout, id, b.repo) : b.repo
}

type Done = { ok: true } | { ok: false; why: string }

async function jaira($: EngineInterface, cwd: string, args: string[]): Promise<Done> {
  const r = await run($, ['jaira', ...args], cwd, 30_000)
  if (r === undefined) return { ok: false, why: 'jaira did not run' }
  if (r.exitCode !== 0) return { ok: false, why: (r.stderr || r.stdout).replace(/\s+/g, ' ').trim().slice(0, 200) }
  return { ok: true }
}

async function say($: EngineInterface, text: string): Promise<void> {
  await update($, noticeAtom, () => text)
  $.ui.toast(text)
}

// The person pressed Accept: that press is the human exit the lane asks for,
// and nothing more. Accept moves the ticket to its next lane, and only where
// that lane is the terminal one (snapshotOf sets acceptTo): from signoff to
// done, never from human past review and signoff. The gates run first as a
// dry run; --force follows only when the human exit is all they refuse, so
// an open definition-of-done item or plan step still stops the press.
async function accept($: EngineInterface, w: Waiting): Promise<void> {
  const t = texts(await language($))
  if (!board || w.acceptTo === null) return
  const dir = await ticketDir($, board, w.id)
  const dry = await run($, ['jaira', 'move', w.id, '--to', w.acceptTo, '--dry-run', '--json'], dir, 30_000)
  const refused = dry === undefined ? 'jaira did not run' : dry.exitCode === 0 ? undefined : otherRefusal(dry.stderr || dry.stdout)
  if (refused !== undefined) {
    await say($, t.failed(`${t.accept} ${w.id}`, refused))
    return
  }
  const done = await jaira($, dir, ['move', w.id, '--to', w.acceptTo, '--force'])
  await say($, done.ok ? t.accepted(w.id, w.acceptTo) : t.failed(`${t.accept} ${w.id}`, done.why))
  await refreshBoard($)
}

// The spawn.sh the roles were installed with: the project's own first, as
// the dispatcher's run-lane.sh finds it beside itself.
async function spawnScript($: EngineInterface): Promise<string | undefined> {
  const tail = '.claude/skills/jaira-dispatcher/scripts/spawn.sh'
  for (const dir of [board?.repo, await home($)]) {
    if (!dir) continue
    try {
      if (await $.fs.exists(`${dir}/${tail}`)) return `${dir}/${tail}`
    } catch {}
  }
  return undefined
}

async function herdr($: EngineInterface): Promise<string | undefined> {
  if ((await $.env.get('HERDR_ENV')) !== '1') return undefined
  return (await $.env.get('HERDR_BIN_PATH')) || 'herdr'
}

// Explain a review finding, or talk an open question through, in a Herdr tab
// of its own: the teamlead's chat is not where that conversation belongs.
async function openChat($: EngineInterface, id: string, kind: 'explain' | 'discuss'): Promise<void> {
  const t = texts(await language($))
  const bin = await herdr($)
  const script = await spawnScript($)
  if (bin === undefined || script === undefined || !board) {
    await say($, t.needsHerdr)
    return
  }
  const dir = await ticketDir($, board, id)
  const r = await run($, [script, '--no-worktree', id.toLowerCase(), id, kind, dir], dir, 120_000)
  await say($, r?.exitCode === 0 ? t.opened(id) : t.failed(kind, (r?.stderr ?? '').trim().slice(-200)))
}

// Every session with a Herdr pane writes where it stands, so the session a
// person works from can tell whether the worker that carried a ticket is
// worth sending work back to.
async function writeRecord($: EngineInterface): Promise<void> {
  if (carrying === undefined) return
  const pane = await $.env.get('HERDR_PANE_ID')
  if (!pane) return
  const usage = await $.session.usage()
  const id = await sessionId($)
  const record: SessionRecord = {
    session: id,
    pane,
    cwd: await $.session.cwd(),
    ticket: carrying.ticket,
    role: carrying.role,
    lane: carrying.lane,
    ...(lastAnswerAt !== undefined ? { lastAnswerAt } : {}),
    tokens: usage.context.tokens ?? 0,
    ttl: ttlMs(usage.rateLimits),
    updatedAt: await $.clock.now(),
  }
  await $.store.set(RECORD_PREFIX + id, record)
}

// Every session's record, those a day old dropped on the way.
async function readRecords($: EngineInterface): Promise<SessionRecord[]> {
  const now = await $.clock.now()
  const out: SessionRecord[] = []
  for (const key of await $.store.keys()) {
    if (!key.startsWith(RECORD_PREFIX)) continue
    const record = (await $.store.get(key)) as SessionRecord | undefined
    if (record === undefined || isExpired(record, now)) await $.store.delete(key)
    else out.push(record)
  }
  return out
}

// Give back: note and move first, then the work goes to the worker that
// carried it while that is cheaper — its tab alive, its cache warm, its
// context small — else its tab is closed FIRST and a new worker started, so
// a ticket never has two worker tabs. The choice and why go onto the ticket.
// The review lane is not touched here: the worker always hands it to a new
// agent, whose tab run-lane.sh closes after the lane.
async function giveBack($: EngineInterface, id: string, lane: string, reason: string): Promise<void> {
  const t = texts(await language($))
  if (!board) return
  const dir = await ticketDir($, board, id)
  const noted = await jaira($, dir, ['note', id, `Zurueckgegeben aus ${lane}: ${reason}`])
  const moved = noted.ok ? await jaira($, dir, ['move', id, '--to', 'in-progress', '--force']) : noted
  if (!moved.ok) {
    await say($, t.failed(`${t.giveBack} ${id}`, moved.why))
    return
  }
  // Who got it, as it happened: the note and the toast say only that.
  const handedTo = async (who: 'old' | 'new' | 'none', why: string) => {
    const to = who === 'old' ? 'den alten Worker' : who === 'new' ? 'einen neuen Worker' : 'niemanden (kein Worker gestartet)'
    await jaira($, dir, ['note', id, `cockpit: Rueckgabe an ${to} - ${why}`])
    await say($, t.returnedTo(id, who))
    await refreshBoard($)
  }
  const bin = await herdr($)
  // No Herdr: nobody to hand it to from here; the board says where it is.
  if (bin === undefined) return handedTo('none', 'kein Herdr in dieser Sitzung')
  const tabs = await run($, [bin, 'tab', 'list'], dir)
  const tab = tabs?.exitCode === 0 ? carrierTab(tabs.stdout, id) : undefined
  const record = carrierOf(await readRecords($), id)
  const choice = chooseWorker({ ...(record ? { record } : {}), tabAlive: tab !== undefined, now: await $.clock.now() })
  let who: 'old' | 'new' | 'none' = 'old'
  let why = choice.why
  if (choice.reuse && record) {
    await run($, [bin, 'pane', 'send-text', record.pane, backMessage(id, lane, reason)], dir)
    await $.clock.sleep(1000)
    await run($, [bin, 'pane', 'send-keys', record.pane, 'enter'], dir)
  } else {
    const script = await spawnScript($)
    if (script === undefined) {
      // Nothing to start a new one with: the old tab stays, so the ticket
      // keeps whatever worker it has.
      who = 'none'
      why = `${choice.why}; kein spawn.sh gefunden (jaira roles install), alter Tab bleibt offen`
    } else {
      if (tab !== undefined) await run($, [bin, 'tab', 'close', tab.tab_id], dir)
      const started = await run($, [script, id.toLowerCase(), id, 'worker', board.repo], board.repo, 120_000)
      who = started?.exitCode === 0 ? 'new' : 'none'
      if (who === 'none') why = `${choice.why}; spawn.sh schlug fehl: ${(started?.stderr ?? '').trim().slice(-160)}`
    }
  }
  await handedTo(who, why)
}

// ── context, cache, status line ──────────────────────────────────────────

async function refreshCache($: EngineInterface): Promise<void> {
  let usage
  try {
    usage = await $.session.usage()
  } catch {
    return
  }
  const state: CacheState = { tokens: usage.context.tokens ?? 0, lastAnswerAt: lastAnswerAt ?? null, ttl: ttlMs(usage.rateLimits) }
  await update($, cacheAtom, () => state)
  const view = await read($, boardAtom)
  const t = texts(await read($, langAtom))
  // What the person used: the five-hour window on a subscription, else what
  // this session cost.
  const window = usage.rateLimits.find(r => /five|5h/i.test(r.kind))
  const used = window !== undefined
    ? `5h ${Math.round(window.percentUsed)}%`
    : usage.cost !== undefined ? `${t.usage} ${dollars(usage.cost.usd)}` : ''
  const lanes = view === null ? '' : t.status(view.counts.worker, view.counts.review, view.counts.person)
  $.ui.status([lanes, used].filter(s => s !== '').join(' · ') || undefined)
}

// ── conversation threads ─────────────────────────────────────────────────

const THREADS_PREFIX = 'threads:'
const THREAD_SESSIONS = 20

async function threadsKey($: EngineInterface): Promise<string> {
  return THREADS_PREFIX + (await sessionId($))
}

// One call of the cheap model per turn of the person's, never per agent
// report: it files the turn under a thread and says where the thread stands.
async function fileTurn($: EngineInterface, human: string, answer: string): Promise<void> {
  const result = await $.model.complete({ model: 'haiku', prompt: classifyPrompt(threads, human, answer, Boolean(board)), maxTokens: 400 })
  if (!result.isAnswered) return
  const filing = parseFiling(result.text)
  if (filing === undefined) return
  threads = file(threads, board ? filing : { ...filing, ticket: null }, turn)
  const key = await threadsKey($)
  // Delete first so the key moves to the end of the store's insertion order,
  // as keep() does: session.start then drops the sessions heard from longest ago.
  await $.store.delete(key)
  await $.store.set(key, { turn, threads })
  const counts: ThreadTally = tally(threads, turn)
  await update($, threadsAtom, () => counts)
}

// ── metrics ──────────────────────────────────────────────────────────────

async function writeMetrics($: EngineInterface): Promise<void> {
  if (!board || totals.input + totals.cacheRead + totals.cacheWrite === 0) return
  const path = `${await home($)}/.jaira/metrics.csv`
  let existing: string | undefined
  try {
    existing = await $.fs.read(path)
  } catch {}
  let arm = 'live'
  try {
    const settings = JSON.parse(await $.fs.read(`${await home($)}/.jaira/settings.json`)) as { dispatch?: unknown }
    if (typeof settings.dispatch === 'string') arm = settings.dispatch
  } catch {}
  const line = row({
    arm,
    ticket: carrying?.ticket ?? '',
    session: await sessionId($),
    lane: carrying?.role === 'lane' ? carrying.lane : carrying?.role ?? 'teamlead',
    totals,
    durationMs: (await $.clock.now()) - startedAt,
  })
  await $.fs.write(path, appended(existing, line))
}

// ── the cache guard ──────────────────────────────────────────────────────

// Hold a person's prompt onto a cold cache with a large context: send it,
// compact first, start over, or cancel. Answers the event's result when it
// holds the prompt back; undefined lets it go on.
async function guarded($: EngineInterface, text: string, origin: string, threshold: number, usdPerMTok: number) {
  const cache = await read($, cacheAtom)
  const held = guard({
    origin,
    text,
    tokens: cache.tokens,
    now: await $.clock.now(),
    lastAnswerAt: lastAnswerAt,
    ttl: cache.ttl,
    threshold,
    usdPerMTok,
  })
  if (held === undefined) return undefined
  const t = texts(await language($))
  const options = [t.guardSend, t.guardCompact, t.guardClear, t.guardCancel]
  let choice: string
  try {
    choice = await $.ui.ask(t.guardQuestion(held.k, held.coldMinutes, dollars(held.usd)), { options, header: t.guardHeader })
  } catch {
    // Nobody to ask (a -p run) or the question dismissed: send, as without
    // the guard.
    return undefined
  }
  if (choice === t.guardSend) return undefined
  // Whatever else they chose, their words wait in the prompt for them.
  await $.prompt.fill({ text })
  if (choice === t.guardCompact) $.clock.after(0, () => void $.session.compact({}))
  else if (choice === t.guardClear) $.clock.after(0, () => void $.command.run({ command: 'clear' }))
  return { drop: choice }
}

// ── the hooks ────────────────────────────────────────────────────────────

const USAGE_COMMAND = 'cockpit'

export const register: Register = (on, options) => {
  const threshold = typeof options.cacheGuardTokens === 'number' ? options.cacheGuardTokens : 30_000
  const usdPerMTok = typeof options.usdPerMTok === 'number' ? options.usdPerMTok : 5

  on('session.start', async ($, e, next) => {
    const started = await next(e)
    startedAt = await $.clock.now()
    // Kept before reports were scoped to a project; no project may show them.
    for (const k of await $.store.keys()) if (isUnscoped(k)) await $.store.delete(k)
    // Every session keeps its threads; the 20 heard from last are enough.
    const threadKeys = (await $.store.keys()).filter(k => k.startsWith(THREADS_PREFIX))
    for (const old of threadKeys.slice(0, Math.max(0, threadKeys.length - THREAD_SESSIONS))) await $.store.delete(old)
    await readRecords($)
    const lang = await language($)
    await update($, langAtom, () => lang)
    const kept = (await $.store.get(await threadsKey($))) as { turn: number; threads: Thread[] } | undefined
    if (kept !== undefined) {
      turn = kept.turn
      threads = kept.threads
      const counts: ThreadTally = tally(threads, turn)
      await update($, threadsAtom, () => counts)
    }
    await $.command.register({
      name: USAGE_COMMAND,
      description: 'The jaira cockpit: report [id], threads, ticket <id>.',
      argumentHint: 'report [id] | threads | ticket <id>',
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
    await refreshBoard($)
    await refreshCache($)
    // The board moves while nobody types; the countdown moves with the clock.
    $.clock.every(60_000, () => {
      void refreshBoard($)
      void refreshCache($)
      void update($, minuteAtom, n => n + 1)
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
    if (untouched.delete(e.text)) return next(e)
    if (AGENT_ORIGINS.has(e.origin.kind)) {
      // A background task of this session always reports in its own block.
      if (e.origin.kind === 'task-notification' && !NOTIFICATION.test(e.text)) return next(e)
      const text = await shorten($, e.text, e.origin.kind, agentName(e))
      if (text === undefined) return next(e)
      shortened.set(e.text, text)
      // A row that never came keeps nothing alive for long.
      for (const old of [...shortened.keys()].slice(0, Math.max(0, shortened.size - 20))) shortened.delete(old)
      return next({ ...e, text })
    }
    const start = startOf(e.text)
    if (start !== undefined && carrying === undefined) carrying = start
    const held = await guarded($, e.text, e.origin.kind, threshold, usdPerMTok)
    if (held !== undefined) return held
    personSaid = isPersonsPrompt(e.origin.kind, e.text) ? e.text : undefined
    return next(e)
  }).catch(($, e, next) => next(e))

  on('session.append', ($, e, next) => next(rewriteRow(e, shortened))).catch(($, e, next) => next(e))

  on('turn.complete', async ($, e, next) => {
    const result = await next(e)
    if (e.agentId !== undefined) return result
    if (result.usage !== undefined) totals = add(totals, result.usage)
    lastAnswerAt = await $.clock.now()
    turn++
    await refreshCache($)
    await writeRecord($)
    await refreshBoard($)
    const human = personSaid
    personSaid = undefined
    if (human !== undefined && e.answer !== '') $.clock.after(0, () => void fileTurn($, human, e.answer))
    return result
  }).catch(($, e, next) => next(e))

  on('tool.call', async ($, e, next) => {
    if (e.tool === 'mcp__jaira-cockpit__report') {
      // A tool's arguments ride on `e` beside `tool`, not under a key of their own.
      const { ticket } = e as { ticket?: unknown }
      return { result: await reports($, typeof ticket === 'string' ? ticket : '', 'en') }
    }
    totals = { ...totals, tools: totals.tools + 1 }
    return next(e)
   }).catch(($, e, next) => next(e))

  on('session.end', async ($, e, next) => {
    await writeMetrics($)
    return next(e)
  }).catch(($, e, next) => next(e))

  on('command.run', { command: USAGE_COMMAND }, async ($, e) => {
    const lang = await language($)
    const t = texts(lang)
    const [sub = '', ...rest] = (e.args ?? '').trim().split(/\s+/)
    if (sub === 'report') return { text: await reports($, rest.join(' '), lang) }
    if (sub === 'threads') return { text: listThreads(threads, turn, t.threadsTitle, t.threadsNone) }
    if (sub === 'ticket' && rest[0]) {
      const id = rest[0].toUpperCase()
      await update($, paneAtom, () => id)
      await $.ui.open({ id: PANE, title: t.paneTitle(id) })
      return { text: t.paneTitle(id) }
    }
    return { text: t.usageText }
  })

  on('ui.render', { component: 'AbovePrompt' }, async ($, e, next) => {
    if (e.props.hasSurvey) return next(e)
    const view = await read($, boardAtom)
    const cache = await read($, cacheAtom)
    const counts = await read($, threadsAtom)
    const giving = await read($, givingAtom)
    const notice = await read($, noticeAtom)
    await read($, minuteAtom)
    const t = texts(await read($, langAtom))
    const now = await $.clock.now()
    const waiting = view?.waiting ?? []
    const c = cacheView(cache.tokens, now, cache.lastAnswerAt ?? undefined, cache.ttl)
    const hasCache = cache.tokens > 0
    const hasThreads = counts.questions + counts.stale > 0
    if (waiting.length === 0 && !hasCache && !hasThreads && notice === null) return next(e)
    const table = $.ui.resolve(e)
    const { Box, Button, Text } = table
    // Mobile has no text field: there the band offers no Give back.
    const Input = 'Input' in table ? table.Input : undefined
    const cacheLine = [
      t.context(c.k),
      c.minutesLeft !== undefined ? t.cacheLeft(c.minutesLeft) : cache.lastAnswerAt !== null ? t.cacheCold : '',
      c.hint === 'new' ? t.hintNew : c.hint === 'compact' ? t.hintCompact : '',
      hasThreads ? t.threadsOpen(counts.questions, counts.stale) : '',
    ].filter(s => s !== '').join(' · ')
    return (
      <Box flexDirection="column">
        {waiting.length > 0 ? <Text bold>{t.waiting}</Text> : null}
        {waiting.slice(0, BAND_TICKETS).map(w => (
          <Box key={`row-${w.id}`} flexDirection="column">
            <Box flexDirection="row">
              <Text wrap="truncate-end">
                {w.id} · {w.lane}{w.verdict !== '' ? ` · ${w.verdict}` : ''}{' '}
              </Text>
              {w.acceptTo !== null
                ? <Button key={`accept-${w.id}`} label={t.accept} variant="primary" onPress={() => void accept($, w)} />
                : null}
              {Input !== undefined
                ? <Button key={`back-${w.id}`} label={t.giveBack} onPress={() => void update($, givingAtom, () => w.id)} />
                : null}
              {hasReservation(w.verdict)
                ? <Button key={`explain-${w.id}`} label={t.explain} onPress={() => void openChat($, w.id, 'explain')} />
                : null}
              {w.question !== ''
                ? <Button key={`discuss-${w.id}`} label={t.discuss} onPress={() => void openChat($, w.id, 'discuss')} />
                : null}
            </Box>
            {w.question !== '' ? <Text dimColor wrap="truncate-end">  ? {w.question}</Text> : null}
            {giving === w.id && Input !== undefined
              ? (
                <Input
                  key={`reason-${w.id}`}
                  label={t.reasonLabel}
                  placeholder={t.reasonPlaceholder}
                  autoFocus
                  onSubmit={(reason: string) => {
                    void update($, givingAtom, () => null)
                    if (reason.trim() !== '') void giveBack($, w.id, w.lane, reason.trim())
                  }}
                />
              )
              : null}
          </Box>
        ))}
        {waiting.length > BAND_TICKETS ? <Text dimColor>{t.more(waiting.length - BAND_TICKETS)}</Text> : null}
        {notice !== null ? <Text dimColor wrap="truncate-end">{notice}</Text> : null}
        {hasCache || hasThreads ? <Text dimColor wrap="truncate-end">{cacheLine}</Text> : null}
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: PANE }, async ($, e) => {
    const id = await read($, paneAtom)
    const t = texts(await read($, langAtom))
    const { Box, Text } = $.ui.resolve(e)
    if (id === null) return <Box><Text dimColor>{t.paneNothing}</Text></Box>
    const kept = await keptReports($, id)
    let csv = ''
    try {
      csv = await $.fs.read(`${await home($)}/.jaira/metrics.csv`)
    } catch {}
    const usage = rowsOf(csv, id)
    const waiting = (await read($, boardAtom))?.waiting.find(w => w.id === id)
    return (
      <Box flexDirection="column">
        {waiting !== undefined ? <Text bold>{waiting.lane}{waiting.verdict !== '' ? ` · ${waiting.verdict}` : ''}</Text> : null}
        <Text bold>{t.paneReports}</Text>
        {kept.length === 0 ? <Text dimColor>{t.paneNothing}</Text> : null}
        {kept.map((r, i) => (
          <Box key={`r-${i}`} flexDirection="column">
            <Text dimColor>── {stamp(r.at)} · {r.origin} ──</Text>
            <Text>{body(r.text)}</Text>
          </Box>
        ))}
        <Text bold>{t.paneUsage}</Text>
        {usage.length === 0 ? <Text dimColor>{t.paneNothing}</Text> : usage.map((u, i) => <Text key={`u-${i}`}>{u}</Text>)}
      </Box>
    )
  })
}
