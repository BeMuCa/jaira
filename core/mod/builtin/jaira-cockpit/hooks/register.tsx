import { atom, read, update } from 'claude-code'
import type { EngineInterface, Register } from 'claude-code'
import type { BoardView, CacheState, ThreadTally } from '../types'
import { explanationNotes, explanationOf, lanesOf, type Action, type Option, noteOf, whoOf, mainRepoOf, otherRefusal, snapshotOf, withoutWarnings, upwards, worktreeFor, type Agent, type Waiting } from './board.ts'
import { cacheView, dollars, guard, isPersonsPrompt, ttlMs, windowsUsed } from './cache.ts'
import { digest } from './digest.ts'
import {
  agentName, AGENT_ORIGINS, clipReport, DIGESTED, isAgentReport, isUnscoped, keyOf, listReports, NOTIFICATION,
  projectKeys, REPORTS_PER_TICKET, rewriteRow, showReport, stamp, TICKETS, type Report,
} from './filter.ts'
import { body, source } from './frame.ts'
import { add, appended, emptyTotals, row, rowsOf, type Totals } from './metrics.ts'
import {
  backMessage, carrierOf, carrierTab, chooseWorker, DEFAULT_CONTEXT_SHARE, isExpired, RECORD_PREFIX, selfLine, startOf, tabsOf,
  workerLines, type Role, type SessionRecord,
} from './sessions.ts'
import { t } from './texts.ts'
import { classifyPrompt, file, listThreads, parseFiling, tally, type Thread } from './threads.ts'

// The engine follows `$` only into functions declared in this same file, and
// a module hooks each event once: so every call on `$` is here, one hook per
// event, and the other files are rules and words that take plain values.
//
// The parts: the noise filter (agent reports as one line, /cockpit report);
// the band above the prompt (tickets waiting for a person with Accept / Give
// back / Explain / Discuss, context and cache, conversation threads); the
// cache guard; the status line; the pane per ticket; the message pane for an
// error or an explanation; metrics per session. Whatever needs a board — a
// .jaira/ and jaira on PATH — stays out of sight without one; the rest works
// in any project.
//
// A session that carries a ticket — a worker, a lane, the dispatcher, started
// by the line spawn.sh types — is an agent's, not a person's: there the mod
// writes the session record Give back reads and its metrics, and nothing
// else. An instruction to it arrives as it was written, and it has no band,
// no status line, no threads and no cache guard.

const boardAtom = atom({ plugin: 'jaira-cockpit', key: 'board' } as const, null)
const cacheAtom = atom({ plugin: 'jaira-cockpit', key: 'cache' } as const, { tokens: 0, lastAnswerAt: null, ttl: 0 })
const threadsAtom = atom({ plugin: 'jaira-cockpit', key: 'threads' } as const, { questions: 0, stale: 0 })
const givingAtom = atom({ plugin: 'jaira-cockpit', key: 'giving' } as const, null)
const paneAtom = atom({ plugin: 'jaira-cockpit', key: 'paneTicket' } as const, null)
const minuteAtom = atom({ plugin: 'jaira-cockpit', key: 'minute' } as const, 0)
const noticeAtom = atom({ plugin: 'jaira-cockpit', key: 'notice' } as const, null)
const messageAtom = atom({ plugin: 'jaira-cockpit', key: 'message' } as const, null)
const explainedAtom = atom({ plugin: 'jaira-cockpit', key: 'explained' } as const, {})
const decidedAtom = atom({ plugin: 'jaira-cockpit', key: 'decided' } as const, {})
const explainingAtom = atom({ plugin: 'jaira-cockpit', key: 'explaining' } as const, null)
const answeringAtom = atom({ plugin: 'jaira-cockpit', key: 'answering' } as const, null)
const usedAtom = atom({ plugin: 'jaira-cockpit', key: 'used' } as const, '')
const expandedAtom = atom({ plugin: 'jaira-cockpit', key: 'expanded' } as const, false)
const theirsOpenAtom = atom({ plugin: 'jaira-cockpit', key: 'theirsOpen' } as const, false)

const PANE = 'cockpit-ticket'
const MESSAGE = 'cockpit-message'
const GIVE_BACK = 'cockpit-give-back'
const EXPLANATION = 'cockpit-explanation'
const ANSWER = 'cockpit-answer'
const BAND_TICKETS = 3
// Between the parts of a line: the band's one separator.
const SEP = ' | '

// ── what this module keeps between events (a reload starts it over) ──────

// The board this session's directory belongs to: undefined until looked
// for, null when there is none.
let board: { dir: string; repo: string } | null | undefined
// Who the person in this session is, by `jaira whoami`: their name and
// aliases; undefined until asked, and when jaira could not say.
let me: string[] | undefined
let meAsked = false
// The ticket this session carries, read off the line that started it.
let carrying: { ticket: string; role: Role; lane: string } | undefined
let lastAnswerAt: number | undefined
let startedAt = 0
let totals: Totals = emptyTotals()
let turn = 0
// The person's message of the turn running, for the threads.
let personSaid: string | undefined
let threads: Thread[] = []
// The Explain and Discuss chats this session opened, by their tab's label
// (`<id>/explain`), with when: the tab goes once the chat has done its part.
const chats = new Map<string, number>()

// ── shared helpers ───────────────────────────────────────────────────────

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

// Whether this session carries a ticket. The line that started it comes once;
// a reload of the mod — a new install, a hot reload — starts the module over,
// so what that line said is kept in the store under the session's id and
// read back once.
const CARRYING_PREFIX = 'carrying:'
let carryingRead = false

async function carries($: EngineInterface): Promise<boolean> {
  if (!carryingRead) {
    carryingRead = true
    const kept = (await $.store.get(CARRYING_PREFIX + (await sessionId($)))) as typeof carrying
    if (carrying === undefined && kept !== undefined) carrying = kept
  }
  return carrying !== undefined
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
  const short = digest(report)
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

async function reports($: EngineInterface, args: string): Promise<string> {
  const root = await project($)
  const ticket = args.trim().toUpperCase()
  if (ticket === '') {
    const keys = projectKeys(await $.store.keys(), root).reverse()
    const rows = await Promise.all(keys.map(async k => [k, ((await $.store.get(k)) as Report[] | undefined) ?? []] as const))
    return listReports(rows)
  }
  return showReport(ticket, await keptReports($, ticket))
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
  if (await carries($)) return
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
  // Workers and reviews are the agents that run: Herdr tabs, and this
  // session's own subagents where there is no Herdr.
  const bin = await herdr($)
  const tabs = bin === undefined ? undefined : await run($, [bin, 'tab', 'list'], board.dir)
  let agents: Agent[] = []
  try {
    agents = await $.agent.list()
  } catch {}
  if (!meAsked) {
    meAsked = true
    const who = await run($, ['jaira', 'whoami', '--json'], board.dir)
    me = who?.exitCode === 0 ? whoOf(who.stdout) : undefined
  }
  const view: BoardView = snapshotOf(lanes.stdout, list.stdout, tabs?.exitCode === 0 ? tabsOf(tabs.stdout) : [], agents, me)
  await update($, boardAtom, () => view)
  // What an Explain chat wrote back onto the tickets the band shows, Explain
  // or not: an explanation stays readable after a later round, only asking
  // for a new one waits for that round's verdict.
  const shownIds = (await read($, expandedAtom)) ? view.waiting : view.waiting.slice(0, BAND_TICKETS)
  const explained: Record<string, { text: string; decided?: string }> = {}
  for (const w of shownIds) {
    const shown = await run($, ['jaira', 'show', w.id, '--json'], board.dir)
    const notes = shown?.exitCode === 0 ? explanationNotes(shown.stdout) : undefined
    if (notes !== undefined) explained[w.id] = notes
  }
  await update($, explainedAtom, () => explained)
  // A chat has done its part once its note is on the ticket — stamped to the
  // minute, so from the minute it was opened: an Explain chat's explanation,
  // a Discuss chat's decision, which the band then offers to read. A Discuss
  // chat whose question went away without one is done too. Its tab goes, but
  // only once its agent stopped working: a Discuss chat notes the decision
  // first and then still moves the ticket and tells this session.
  const open = bin !== undefined && tabs?.exitCode === 0 ? tabsOf(tabs.stdout) : []
  for (const [label, since] of chats) {
    const [id = '', kind] = label.split('/')
    const tab = open.find(x => x.label === label)
    if (tab?.agent_status === 'working') continue
    const from = since - (since % 60_000)
    const note = await noteOn($, board.dir, id, kind === 'explain' ? 'explain' : 'decision')
    const fresh = note !== undefined && note.at >= from ? note : undefined
    // A teammate's ticket can be talked through too, and sits in theirs.
    const w = [...view.waiting, ...view.theirs].find(x => x.id === id)
    const done = fresh !== undefined || (kind === 'discuss' && (w === undefined || w.question === ''))
    if (!done) continue
    chats.delete(label)
    if (kind === 'discuss' && fresh !== undefined) await update($, decidedAtom, d => ({ ...d, [id]: fresh.text }))
    if (tab !== undefined) await run($, [bin!, 'tab', 'close', tab.tab_id], board.dir)
  }
}

async function noteOn($: EngineInterface, dir: string, id: string, kind: 'explain' | 'decision') {
  const shown = await run($, ['jaira', 'show', id, '--json'], dir)
  return shown?.exitCode === 0 ? noteOf(shown.stdout, kind) : undefined
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
  if (r.exitCode !== 0) return { ok: false, why: withoutWarnings(r.stderr || r.stdout) || `jaira exited ${r.exitCode}` }
  return { ok: true }
}

async function say($: EngineInterface, text: string): Promise<void> {
  await update($, noticeAtom, () => text)
  $.ui.toast(text)
}

// Text a person reads whole, an error or an explanation: a pane of its own
// with a heading, the text wrapped, and Close. Where no pane fits, a toast
// carries it.
async function show($: EngineInterface, title: string, text: string): Promise<void> {
  await update($, messageAtom, () => ({ title, text }))
  const opened = await $.ui.open({ id: MESSAGE, title, focus: true, closeOnEscape: true })
  if (!opened.isPlaced) $.ui.toast(`${title}: ${text}`)
}

async function fail($: EngineInterface, what: string, why: string): Promise<void> {
  await update($, noticeAtom, () => null)
  await show($, t.failed(what), why)
}

// The person pressed Accept: that press is the human exit the lane asks for,
// and nothing more. Accept moves the ticket to its next lane, and only where
// that lane is the terminal one (snapshotOf sets acceptTo): from signoff to
// done, never from human past review and signoff. The gates run first as a
// dry run; --force follows only when the human exit is all they refuse, so
// an open definition-of-done item or plan step still stops the press.
async function accept($: EngineInterface, w: Waiting): Promise<void> {
  if (!board || w.acceptTo === null) return
  const dir = await ticketDir($, board, w.id)
  const dry = await run($, ['jaira', 'move', w.id, '--to', w.acceptTo, '--dry-run', '--json'], dir, 30_000)
  const refused = dry === undefined ? 'jaira did not run' : dry.exitCode === 0 ? undefined : otherRefusal(dry.stderr || dry.stdout)
  if (refused !== undefined) return fail($, `${t.accept} ${w.id}`, refused)
  const done = await jaira($, dir, ['move', w.id, '--to', w.acceptTo, '--force'])
  if (done.ok) await say($, t.accepted(w.id, w.acceptTo))
  else await fail($, `${t.accept} ${w.id}`, done.why)
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
  const bin = await herdr($)
  const script = await spawnScript($)
  const what = `${kind === 'explain' ? t.explain : t.discuss} ${id}`
  if (bin === undefined || script === undefined || !board) return fail($, what, t.needsHerdr)
  const dir = await ticketDir($, board, id)
  // A Discuss chat sends its decision to this session: its messaging socket
  // is the address SendMessage takes.
  const socket = kind === 'discuss' ? await $.env.get('CLAUDE_CODE_MESSAGING_SOCKET') : undefined
  const parent = socket ? ['--parent', `uds:${socket}`] : []
  const r = await run($, [script, '--no-worktree', ...parent, id.toLowerCase(), id, kind, dir], dir, 120_000)
  if (r?.exitCode !== 0) return fail($, what, (r?.stderr || r?.stdout || 'spawn.sh did not run').trim())
  chats.set(`${id}/${kind}`, await $.clock.now())
  await say($, t.opened(id))
}

// '+N more' opens every waiting ticket in the band, with its buttons; the
// band scrolls where they do not fit. Opened, the explanations of all of
// them are read too.
async function toggleMore($: EngineInterface, open: boolean): Promise<void> {
  await update($, expandedAtom, () => open)
  if (open) await refreshBoard($)
}

// A decision is read once: the pane shows it, the band lets it go.
async function readDecision($: EngineInterface, id: string, text: string): Promise<void> {
  await update($, decidedAtom, d => Object.fromEntries(Object.entries(d).filter(([k]) => k !== id)))
  await show($, t.decisionOf(id), text)
}

// An explanation is read in a pane that takes the keyboard: its options are
// a choice there, with an own answer beside them, and the choice goes onto
// the ticket as a decision.
async function readExplanation($: EngineInterface, id: string, notes: { text: string; decided?: string }): Promise<void> {
  await update($, explainingAtom, () => ({ id, ...notes }))
  const opened = await $.ui.open({ id: EXPLANATION, title: t.explanationOf(id), focus: true, closeOnEscape: true })
  if (!opened.isPlaced) $.ui.toast(`${t.explanationOf(id)}: ${notes.text}`)
}

// An own answer is typed in a pane of its own, opened with focus, as Give
// back asks: a field drawn into the open explanation pane after the pick is
// not sure to get the keyboard.
async function askOwnAnswer($: EngineInterface, id: string): Promise<void> {
  await $.ui.close({ id: EXPLANATION })
  await update($, answeringAtom, () => id)
  const opened = await $.ui.open({ id: ANSWER, title: t.answerOn(id), focus: true, closeOnEscape: true })
  if (!opened.isPlaced) await fail($, `${t.decide} ${id}`, t.noRoom)
}

// The choice is the decision itself: noted on the ticket first, as a Discuss
// chat notes it, then carried out — Give back with the choice as its reason,
// Accept where the band's own Accept would move it. A choice that does
// neither goes to this session's conversation, the teamlead's, so it is
// heard and not only left on the ticket.
async function decide($: EngineInterface, id: string, text: string, action: Action | undefined, share: number): Promise<void> {
  if (!board) return
  const done = await jaira($, await ticketDir($, board, id), ['note', id, `decision: ${text}`])
  if (!done.ok) return fail($, `${t.decide} ${id}`, done.why)
  await update($, decidedAtom, d => ({ ...d, [id]: text }))
  // Answered from here on, before the next look at the notes says so too.
  await update($, explainedAtom, e => (e[id] === undefined ? e : { ...e, [id]: { ...e[id], decided: text } }))
  const w = (await read($, boardAtom))?.waiting.find(x => x.id === id)
  if (action === 'giveBack' && w !== undefined) return giveBack($, id, w.lane, text, share)
  if (action === 'accept' && w !== undefined && w.acceptTo !== null) return accept($, w)
  await say($, t.decidedOn(id))
  void $.prompt.submit({ text: t.decisionFor(id, text) })
}

// Give back asks what is missing in a pane of its own, opened with focus: a
// pane that takes the keyboard puts it on its autoFocus field. The same field
// in the band never got it — the band holds the keyboard already when Give
// back is pressed, and the ring stays on the button.
async function askWhatIsMissing($: EngineInterface, w: Waiting): Promise<void> {
  await update($, givingAtom, () => w.id)
  const opened = await $.ui.open({ id: GIVE_BACK, title: `${t.giveBack} ${w.id}`, focus: true, closeOnEscape: true })
  if (!opened.isPlaced) await fail($, `${t.giveBack} ${w.id}`, t.noRoom)
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
    repo: await project($),
    ticket: carrying.ticket,
    role: carrying.role,
    lane: carrying.lane,
    ...(lastAnswerAt !== undefined ? { lastAnswerAt } : {}),
    tokens: usage.context.tokens ?? 0,
    ...(usage.context.window ? { window: usage.context.window } : {}),
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
async function giveBack($: EngineInterface, id: string, lane: string, reason: string, share: number): Promise<void> {
  if (!board) return
  const dir = await ticketDir($, board, id)
  const noted = await jaira($, dir, ['note', id, `Given back from ${lane}: ${reason}`])
  // Back to an earlier lane needs no --force since jaira 0.3.7, which this
  // mod ships inside: so any other refusal of the move still stops it.
  const moved = noted.ok ? await jaira($, dir, ['move', id, '--to', 'in-progress']) : noted
  if (!moved.ok) return fail($, `${t.giveBack} ${id}`, moved.why)
  // Who got it, as it happened: the note and the toast say only that.
  const handedTo = async (who: 'old' | 'new' | 'none', why: string) => {
    const to = who === 'old' ? 'the old worker' : who === 'new' ? 'a new worker' : 'nobody (no worker started)'
    await jaira($, dir, ['note', id, `cockpit: given back to ${to} - ${why}`])
    await say($, t.returnedTo(id, who))
    await refreshBoard($)
  }
  const bin = await herdr($)
  // No Herdr: nobody to hand it to from here; the board says where it is.
  if (bin === undefined) return handedTo('none', 'no Herdr in this session')
  const tabs = await run($, [bin, 'tab', 'list'], dir)
  const tab = tabs?.exitCode === 0 ? carrierTab(tabsOf(tabs.stdout), id) : undefined
  const record = carrierOf(await readRecords($), id)
  const choice = chooseWorker({ ...(record ? { record } : {}), tabAlive: tab !== undefined, now: await $.clock.now(), share })
  let why = choice.why
  if (choice.reuse && record) {
    // The pane may have gone since the tab list: only text that arrived
    // counts as handed to the old worker.
    const sent = await run($, [bin, 'pane', 'send-text', record.pane, backMessage(id, lane, reason)], dir)
    if (sent?.exitCode === 0) {
      await $.clock.sleep(1000)
      await run($, [bin, 'pane', 'send-keys', record.pane, 'enter'], dir)
      return handedTo('old', why)
    }
    why = `${why}, but its pane took no text`
  }
  const script = await spawnScript($)
  // Nothing to start a new one with: the old tab stays, so the ticket keeps
  // whatever worker it has.
  if (script === undefined) return handedTo('none', `${why}; no spawn.sh found (jaira roles install), old tab left open`)
  if (tab !== undefined) await run($, [bin, 'tab', 'close', tab.tab_id], dir)
  const started = await run($, [script, id.toLowerCase(), id, 'worker', board.repo], board.repo, 120_000)
  if (started?.exitCode === 0) return handedTo('new', why)
  await handedTo('none', `${why}; spawn.sh failed: ${(started?.stderr ?? '').trim().slice(-160)}`)
}

// The workers of this board, for the teamlead (the tool) and a person
// (/cockpit workers) alike: the asking session first, from its live figures,
// then every agent carrying a ticket, with whom a new round would go to.
async function workers($: EngineInterface, share: number): Promise<string> {
  let self = ''
  try {
    const usage = await $.session.usage()
    self = selfLine(usage.context.tokens ?? 0, usage.context.window, ttlMs(usage.rateLimits))
  } catch {}
  // A worker asking never looked for the board: it draws nothing.
  if (board === undefined) board = await findBoard($)
  if (board === null) return [self, t.workersNoBoard].filter(Boolean).join('\n')
  const list = await run($, ['jaira', 'list', '--json'], board.dir)
  const bin = await herdr($)
  const tabs = bin === undefined ? undefined : await run($, [bin, 'tab', 'list'], board.dir)
  const lines = workerLines({
    records: await readRecords($),
    repo: await project($),
    ...(tabs?.exitCode === 0 ? { tabs: tabsOf(tabs.stdout) } : {}),
    lanes: lanesOf(list?.exitCode === 0 ? list.stdout : ''),
    now: await $.clock.now(),
    share,
  })
  return [self, lines.length === 0 ? t.workersNone : t.workersTitle, ...lines].filter(Boolean).join('\n')
}

// ── context, cache, what was used ──────────────────────────────────────────

async function refreshCache($: EngineInterface): Promise<void> {
  if (await carries($)) return
  let usage
  try {
    usage = await $.session.usage()
  } catch {
    return
  }
  const state: CacheState = { tokens: usage.context.tokens ?? 0, lastAnswerAt: lastAnswerAt ?? null, ttl: ttlMs(usage.rateLimits) }
  await update($, cacheAtom, () => state)
  // What the person used: the five-hour and seven-day windows on a
  // subscription, else what this session cost. The band's head line shows
  // it beside the board's counts; the mod keeps no status line of its own.
  const used = windowsUsed(usage.rateLimits) || (usage.cost !== undefined ? `${t.usage} ${dollars(usage.cost.usd)}` : '')
  await update($, usedAtom, () => used)
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
  // The share of a model's window, in percent, past which a ticket's next
  // round goes to a new worker.
  const share = typeof options.workerContextShare === 'number' ? options.workerContextShare : DEFAULT_CONTEXT_SHARE

  on('session.start', async ($, e, next) => {
    const started = await next(e)
    startedAt = await $.clock.now()
    // Kept before reports were scoped to a project; no project may show them.
    for (const k of await $.store.keys()) if (isUnscoped(k)) await $.store.delete(k)
    // A session that ended without session.end leaves its mark; the 50 heard
    // from last are plenty.
    const carryingKeys = (await $.store.keys()).filter(k => k.startsWith(CARRYING_PREFIX))
    for (const old of carryingKeys.slice(0, Math.max(0, carryingKeys.length - 50))) await $.store.delete(old)
    // Every session keeps its threads; the 20 heard from last are enough.
    const threadKeys = (await $.store.keys()).filter(k => k.startsWith(THREADS_PREFIX))
    for (const old of threadKeys.slice(0, Math.max(0, threadKeys.length - THREAD_SESSIONS))) await $.store.delete(old)
    await readRecords($)
    const kept = (await $.store.get(await threadsKey($))) as { turn: number; threads: Thread[] } | undefined
    if (kept !== undefined) {
      turn = kept.turn
      threads = kept.threads
      const counts: ThreadTally = tally(threads, turn)
      await update($, threadsAtom, () => counts)
    }
    await $.command.register({
      name: USAGE_COMMAND,
      description: 'The jaira cockpit: report [id], threads, ticket <id>, workers.',
      argumentHint: 'report [id] | threads | ticket <id> | workers',
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
    await $.tool.register({
      name: 'workers',
      description:
        "Lists this session's own context and prompt cache, then every agent carrying a ticket of this jaira board (worker or dispatcher, with an open Herdr tab): ticket, lane, context tokens, minutes its prompt cache stays warm, and whether the ticket's next round would go to the same worker or a new one, and why.",
      inputSchema: { type: 'object', properties: {} },
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
    if (await carries($)) return next(e)
    if (!isAgentReport(e)) {
      untouched.add(e.text)
      return next(e)
    }
    const text = await shorten($, e.text, e.origin.kind)
    return next(text === undefined ? e : { ...e, text })
  }).catch(($, e, next) => next(e))

  on('prompt.submit', async ($, e, next) => {
    if (untouched.delete(e.text)) return next(e)
    if (await carries($)) return next(e)
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
    if (start !== undefined) {
      // From here on the session is an agent's: what session.start drew goes.
      carrying = start
      await $.store.set(CARRYING_PREFIX + (await sessionId($)), start)
      await update($, boardAtom, () => null)
      return next(e)
    }
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
    await writeRecord($)
    if (await carries($)) return result
    await refreshCache($)
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
      return { result: await reports($, typeof ticket === 'string' ? ticket : '') }
    }
    if (e.tool === 'mcp__jaira-cockpit__workers') return { result: await workers($, share) }
    totals = { ...totals, tools: totals.tools + 1 }
    return next(e)
   }).catch(($, e, next) => next(e))

  on('session.end', async ($, e, next) => {
    await writeMetrics($)
    await $.store.delete(CARRYING_PREFIX + (await sessionId($)))
    return next(e)
  }).catch(($, e, next) => next(e))

  on('command.run', { command: USAGE_COMMAND }, async ($, e) => {
    const [sub = '', ...rest] = (e.args ?? '').trim().split(/\s+/)
    if (sub === 'report') return { text: await reports($, rest.join(' ')) }
    if (sub === 'workers') return { text: await workers($, share) }
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
    if (e.props.hasSurvey || (await carries($))) return next(e)
    const view = await read($, boardAtom)
    const cache = await read($, cacheAtom)
    const counts = await read($, threadsAtom)
    const notice = await read($, noticeAtom)
    const explained = await read($, explainedAtom)
    const decided = await read($, decidedAtom)
    const used = await read($, usedAtom)
    const expanded = await read($, expandedAtom)
    const theirsOpen = await read($, theirsOpenAtom)
    await read($, minuteAtom)
    const now = await $.clock.now()
    const waiting = view?.waiting ?? []
    const theirs = view?.theirs ?? []
    const decisions = Object.entries(decided)
    const c = cacheView(cache.tokens, now, cache.lastAnswerAt ?? undefined, cache.ttl)
    const hasCache = cache.tokens > 0
    const hasThreads = counts.questions + counts.stale > 0
    // The head line: the board's counts and what was used, beside the title.
    const head = [
      ...(view !== null ? [t.workers(view.counts.worker), t.reviews(view.counts.review), t.waitingFor(view.counts.person)] : []),
      ...(used !== '' ? [used] : []),
    ]
    if (head.length === 0 && waiting.length === 0 && theirs.length === 0 && decisions.length === 0 && !hasCache && !hasThreads && notice === null) return next(e)
    const table = $.ui.resolve(e)
    const { Box, Button, Text } = table
    // Mobile has no text field: there the band offers no Give back.
    const hasInput = 'Input' in table
    // The session line: context, cache, the advice, threads; what asks for
    // action in the warning colour.
    const session: { text: string; color?: 'warning' | 'error' }[] = [
      { text: t.context(c.k) },
      ...(c.minutesLeft !== undefined ? [{ text: t.cacheLeft(c.minutesLeft) }] : cache.lastAnswerAt !== null ? [{ text: t.cacheCold, color: 'warning' as const }] : []),
      ...(c.hint === 'new' ? [{ text: t.hintNew, color: 'error' as const }] : c.hint === 'compact' ? [{ text: t.hintCompact, color: 'warning' as const }] : []),
      ...(hasThreads ? [{ text: t.threadsOpen(counts.questions, counts.stale), color: 'warning' as const }] : []),
    ]
    const shown = expanded ? waiting : waiting.slice(0, BAND_TICKETS)
    // A frame of its own, so the band reads as the cockpit's and not as chat.
    return (
      <Box flexDirection="column" borderStyle="round" borderDimColor paddingX={1}>
        <Box flexDirection="row">
          <Text bold color="claude">{t.bandTitle}</Text>
          {head.length > 0 ? <Text dimColor wrap="truncate-end">{'   '}{head.join(SEP)}</Text> : null}
        </Box>
        {/* What just happened on top, this person's own work below it, the
            teammates' tickets last of all lists. */}
        {notice !== null ? <Text dimColor wrap="truncate-end">{notice}</Text> : null}
        {waiting.length > 0 ? <Box marginTop={1}><Text bold>{t.waiting}</Text></Box> : null}
        {shown.map(w => (
          <Box key={`row-${w.id}`} flexDirection="column">
            <Box flexDirection="row">
              <Text wrap="truncate-end">
                <Text bold>{w.id}</Text>{SEP}{w.lane}{w.verdict !== '' ? `${SEP}${w.verdict}` : ''}{' '}
              </Text>
              {w.acceptTo !== null
                ? <Button key={`accept-${w.id}`} label={t.accept} variant="primary" onPress={() => void accept($, w)} />
                : null}
              {hasInput
                ? <Button key={`back-${w.id}`} label={t.giveBack} onPress={() => void askWhatIsMissing($, w)} />
                : null}
              {w.explainable
                ? <Button key={`explain-${w.id}`} label={t.explain} onPress={() => void openChat($, w.id, 'explain')} />
                : null}
              {explained[w.id] !== undefined
                ? <Button key={`read-${w.id}`} label={t.readExplanation} onPress={() => void readExplanation($, w.id, explained[w.id]!)} />
                : null}
              {w.verdict !== ''
                ? <Button key={`verdict-${w.id}`} label={t.verdict} onPress={() => void show($, t.verdictOf(w.id), w.verdict)} />
                : null}
              {w.question !== ''
                ? <Button key={`discuss-${w.id}`} label={t.discuss} onPress={() => void openChat($, w.id, 'discuss')} />
                : null}
            </Box>
            {w.question !== '' ? <Text dimColor wrap="truncate-end">  ? {w.question}</Text> : null}
          </Box>
        ))}
        {waiting.length > BAND_TICKETS
          ? (
            <Box flexDirection="row">
              <Button
                key={expanded ? 'fewer' : 'more'}
                label={expanded ? t.fewer : t.more(waiting.length - BAND_TICKETS)}
                onPress={() => void toggleMore($, !expanded)}
              />
            </Box>
          )
          : null}
        {decisions.map(([id, text]) => (
          <Box key={`decided-${id}`} flexDirection="row">
            <Text><Text bold>{id}</Text>{SEP}{t.decided}{' '}</Text>
            <Button key={`decision-${id}`} label={t.readDecision} onPress={() => void readDecision($, id, text)} />
          </Box>
        ))}
        {theirs.length > 0
          ? (
            <Box flexDirection="row" marginTop={1}>
              <Button
                key="theirs"
                label={t.theirs(theirs.length, [...new Set(theirs.map(w => w.assignee))])}
                onPress={() => void update($, theirsOpenAtom, open => !open)}
              />
            </Box>
          )
          : null}
        {theirsOpen
          ? theirs.map(w => (
            <Box key={`their-${w.id}`} flexDirection="column">
              <Box flexDirection="row">
                <Text dimColor wrap="truncate-end">
                  {w.id}{SEP}{w.lane}{SEP}{w.assignee}{w.verdict !== '' ? `${SEP}${w.verdict}` : ''}{' '}
                </Text>
                {w.question !== ''
                  ? <Button key={`discuss-${w.id}`} label={t.discuss} onPress={() => void openChat($, w.id, 'discuss')} />
                  : null}
              </Box>
              {w.question !== '' ? <Text dimColor wrap="truncate-end">  ? {w.question}</Text> : null}
            </Box>
          ))
          : null}
        {hasCache || hasThreads
          ? (
            <Box marginTop={1} flexDirection="row">
              {session.flatMap((part, i) => [
                ...(i > 0 ? [<Text key={`sep-${i}`} dimColor>{SEP}</Text>] : []),
                part.color !== undefined
                  ? <Text key={`part-${i}`} color={part.color}>{part.text}</Text>
                  : <Text key={`part-${i}`} dimColor>{part.text}</Text>,
              ])}
            </Box>
          )
          : null}
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: GIVE_BACK }, async ($, e) => {
    const id = await read($, givingAtom)
    const w = (await read($, boardAtom))?.waiting.find(x => x.id === id)
    const table = $.ui.resolve(e)
    const { Box, Text } = table
    const Input = 'Input' in table ? table.Input : undefined
    if (w === undefined || Input === undefined) return <Box><Text dimColor>{t.paneNothing}</Text></Box>
    return (
      <Box flexDirection="column">
        <Text>{w.id} · {w.lane}{w.verdict !== '' ? ` · ${w.verdict}` : ''}</Text>
        <Input
          key={`reason-${w.id}`}
          label={t.reasonLabel}
          placeholder={t.reasonPlaceholder}
          submitLabel={t.giveBack}
          autoFocus
          onSubmit={(reason: string) => {
            // Nothing typed: the pane stays open for the reason.
            if (reason.trim() === '') return
            void update($, givingAtom, () => null)
            void $.ui.close({ id: GIVE_BACK })
            void giveBack($, w.id, w.lane, reason.trim(), share)
          }}
        />
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: EXPLANATION }, async ($, e) => {
    const explaining = await read($, explainingAtom)
    const { Box, Button, Text } = $.ui.resolve(e)
    if (explaining === null) return <Box><Text dimColor>{t.paneNothing}</Text></Box>
    const { id, decided } = explaining
    const x = explanationOf(explaining.text)
    const open = decided === undefined
    const pick = (o: Option) => {
      void $.ui.close({ id: EXPLANATION })
      void decide($, id, `${o.key}) ${o.text}`, o.action, share)
    }
    // An option whole: what it says, what picking it does, whether the
    // review recommends it. Drawn as many rows as it needs.
    const said = (o: Option) =>
      `${o.text}${o.action !== undefined ? ` ${o.action === 'accept' ? t.actionAccept : t.actionGiveBack}` : ''}${o.key === x.recommended ? ` ${t.recommended}` : ''}`
    const line = (o: Option) => <Text wrap="wrap">{said(o)}</Text>
    return (
      <Box flexDirection="column">
        <Text bold>{t.explanationOf(id)}</Text>
        {x.parts.map((p, i) => p.kind === 'option'
          // Open, each option is a row of its own to pick: the engine's
          // Select draws one line per option and cuts a long one, so the
          // choice would hide what is chosen. Its letter picks it too.
          ? open
            ? <Button key={`pick-${id}-${p.key}`} plain hotkey={p.key} label={`${p.key}) ${said(p)}`} {...(p.key === (x.recommended ?? x.options[0]?.key) ? { autoFocus: true as const } : {})} onPress={() => pick(p)}>{line(p)}</Button>
            : <Box key={`part-${i}`} flexDirection="row"><Text bold>{`  ${p.key}) `}</Text>{line(p)}</Box>
          // A paragraph a line apart; the options one under the other.
          : <Box key={`part-${i}`} marginTop={i > 0 ? 1 : 0} marginBottom={x.parts[i + 1]?.kind === 'option' ? 1 : 0}><Text wrap="wrap">{p.text}</Text></Box>)}
        {/* Answered already, on the ticket: the choice is read, not offered again. */}
        {open
          ? <Box marginTop={1}><Button key={`pick-${id}-own`} plain hotkey="o" {...(x.options.length === 0 ? { autoFocus: true as const } : {})} label={t.ownAnswer} onPress={() => void askOwnAnswer($, id)} /></Box>
          : <Box marginTop={1}><Text wrap="wrap"><Text bold>{t.decidedAs}</Text> {decided}</Text></Box>}
        <Button key="explanation-close" label={t.close} onPress={() => void $.ui.close({ id: EXPLANATION })} />
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: ANSWER }, async ($, e) => {
    const id = await read($, answeringAtom)
    const table = $.ui.resolve(e)
    const { Box, Text } = table
    const Input = 'Input' in table ? table.Input : undefined
    if (id === null || Input === undefined) return <Box><Text dimColor>{t.paneNothing}</Text></Box>
    return (
      <Box flexDirection="column">
        <Input
          key={`answer-${id}`}
          label={t.ownLabel}
          placeholder={t.ownPlaceholder}
          submitLabel={t.decide}
          autoFocus
          onSubmit={(answer: string) => {
            // Nothing typed: the pane stays open for the answer.
            if (answer.trim() === '') return
            void update($, answeringAtom, () => null)
            void $.ui.close({ id: ANSWER })
            // An own answer asks for something the options did not offer:
            // it goes back, with the answer as the reason.
            void decide($, id, answer.trim(), 'giveBack', share)
          }}
        />
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: MESSAGE }, async ($, e) => {
    const message = await read($, messageAtom)
    const { Box, Button, Text } = $.ui.resolve(e)
    return (
      <Box flexDirection="column">
        <Text bold>{message?.title ?? ''}</Text>
        <Text wrap="wrap">{message?.text ?? ''}</Text>
        <Button key="message-close" label={t.close} onPress={() => void $.ui.close({ id: MESSAGE })} />
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: PANE }, async ($, e) => {
    const id = await read($, paneAtom)
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
