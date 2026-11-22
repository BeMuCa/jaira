import { atom, read, update } from 'claude-code'
import type { ElementTable, EngineInterface, Register } from 'claude-code'
import type { BoardView, CacheState, ThreadTally } from '../types'
import {
  branchesFor, carrierBranch, commitNaming, defaultBranchOf, remoteFor, type Forge, newWorktreeFor, worktreeOfBranch, pendingOf, ticketsNamed, forgeOf, hasPullRequest, repoViewCommand, targetRepoOf,
  explanationNotes, explanationOf, huntOf, isOlder, lightOf, lanesOf, MIN_JAIRA, otherBoardsOf, pullRequestsCommand, stepAfterAccept, type StepAfterAccept, noteOf, whoOf, mainRepoOf, otherRefusal, ownTreeFor, diskOf, snapshotOf, versionOf, withoutWarnings, upwards, worktreeFor,
  offeredRelease, upgradeCheckOf,
  commitsOf, createdHandle, decisionLines, executionOf, oneLine, decisionsOf, followTitle, gapItems, languageOf, NOTE_MARK, reasonOf, reasonPrompt, releaseLines, headPrompt, rephrasedOf, rephrasePrompt, standOf, storyOf,
  type Action, type Agent, type GapItem, type Story, type Hunt, type HuntState, type NoteClass, type Option, type Waiting,
} from './board.ts'
import { cacheView, dollars, guard, isPersonsPrompt, ttlMs, windowsUsed } from './cache.ts'
import { digest, reportLine } from './digest.ts'
import {
  agentName, AGENT_ORIGINS, clipReport, DIGESTED, isAgentReport, isUnscoped, keyOf, listReports, NOTIFICATION,
  projectKeys, reportArgs, REPORTS_PER_TICKET, rewriteRow, showReport, stamp, TICKETS, type Report,
} from './filter.ts'
import { body, source } from './frame.ts'
import { add, appended, emptyTotals, row, rowsOf, type Totals } from './metrics.ts'
import { COLORS, HATS, overviewOf, pickOf, stoneOf, workerKey, workersOf } from './overview.ts'
import { drawOverview, headRow, type OverviewTab } from './overview-pane.tsx'
import {
  backMessage, carrierOf, carrierTab, chooseWorker, DEFAULT_CONTEXT_SHARE, isExpired, RECORD_PREFIX, selfLine, startOf, tabsOf,
  workerLines, type Role, type SessionRecord, type Tab,
} from './sessions.ts'
import { overviewWords, t } from './texts.ts'
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
const choosingAtom = atom({ plugin: 'jaira-cockpit', key: 'choosing' } as const, null)
const paneAtom = atom({ plugin: 'jaira-cockpit', key: 'paneTicket' } as const, null)
const minuteAtom = atom({ plugin: 'jaira-cockpit', key: 'minute' } as const, 0)
const noticeAtom = atom({ plugin: 'jaira-cockpit', key: 'notice' } as const, null)
const messageAtom = atom({ plugin: 'jaira-cockpit', key: 'message' } as const, null)
const explainedAtom = atom({ plugin: 'jaira-cockpit', key: 'explained' } as const, {})
const decidedAtom = atom({ plugin: 'jaira-cockpit', key: 'decided' } as const, {})
const storyAtom = atom({ plugin: 'jaira-cockpit', key: 'story' } as const, null)
const typedAtom = atom({ plugin: 'jaira-cockpit', key: 'typed' } as const, '')
const refusalAtom = atom({ plugin: 'jaira-cockpit', key: 'refusal' } as const, null)
const reasonsAtom = atom({ plugin: 'jaira-cockpit', key: 'reasons' } as const, {})
const rephrasedAtom = atom({ plugin: 'jaira-cockpit', key: 'rephrased' } as const, {})
const huntsAtom = atom({ plugin: 'jaira-cockpit', key: 'hunts' } as const, {})
const explainingAtom = atom({ plugin: 'jaira-cockpit', key: 'explaining' } as const, null)
const answeringAtom = atom({ plugin: 'jaira-cockpit', key: 'answering' } as const, null)
const usedAtom = atom({ plugin: 'jaira-cockpit', key: 'used' } as const, '')
const expandedAtom = atom({ plugin: 'jaira-cockpit', key: 'expanded' } as const, false)
const theirsOpenAtom = atom({ plugin: 'jaira-cockpit', key: 'theirsOpen' } as const, false)
const busyAtom = atom({ plugin: 'jaira-cockpit', key: 'busy' } as const, {})
const overviewAtom = atom({ plugin: 'jaira-cockpit', key: 'overview' } as const, null)
const overviewTabAtom = atom({ plugin: 'jaira-cockpit', key: 'overviewTab' } as const, 'queue')
const overviewFoldedAtom = atom({ plugin: 'jaira-cockpit', key: 'overviewFolded' } as const, {})
const overviewPeopleAtom = atom({ plugin: 'jaira-cockpit', key: 'overviewPeople' } as const, [])
const overviewLanesAtom = atom({ plugin: 'jaira-cockpit', key: 'overviewLanes' } as const, [])
const releaseAtom = atom({ plugin: 'jaira-cockpit', key: 'release' } as const, null)
const upgradeAtom = atom({ plugin: 'jaira-cockpit', key: 'upgrade' } as const, null)
const workersAtom = atom({ plugin: 'jaira-cockpit', key: 'workers' } as const, null)
const workersSinceAtom = atom({ plugin: 'jaira-cockpit', key: 'workersSince' } as const, {})
const workersReportsAtom = atom({ plugin: 'jaira-cockpit', key: 'workersReports' } as const, {})
const workersFellAtom = atom({ plugin: 'jaira-cockpit', key: 'workersFell' } as const, -1)
const workersColorsAtom = atom({ plugin: 'jaira-cockpit', key: 'workersColors' } as const, {})
const workersHatsAtom = atom({ plugin: 'jaira-cockpit', key: 'workersHats' } as const, {})
const workersFrameAtom = atom({ plugin: 'jaira-cockpit', key: 'workersFrame' } as const, 0)

const PANE = 'cockpit-ticket'
const MESSAGE = 'cockpit-message'
const ACCEPTANCE = 'cockpit-acceptance'
// Explain all's lock, beside the row's: a chat holds its press until it is
// done, and Execute in the same pane must not wait on it.
const CHAT = ':chat'
const EXPLANATION = 'cockpit-explanation'
// The Acceptance pane's light in the theme's colours.
const LIGHT_COLOUR = { green: 'success', yellow: 'warning', red: 'error' } as const
const ANSWER = 'cockpit-answer'
const OVERVIEW = 'cockpit-overview'
const UPGRADE = 'cockpit-upgrade'
// The store key of the release the person dismissed: kept across sessions,
// so that release is never offered again, and a newer one is.
const DISMISSED_RELEASE = 'dismissed-release'
// The release row's key among the pressed rows: Update and the yes in its
// pane lock it as a ticket's buttons lock its row.
const RELEASE = 'release'
const BAND_TICKETS = 3
// How long a pressed row stays locked when its action has not ended by then.
const BUSY_MS = 20_000
// The Workers tab's frame: the figures hammer each frame, and the board is
// read again every READ_FRAMES, so an agent that came or went shows without
// a reopening.
const FRAME_MS = 250
const READ_FRAMES = 20
// The Overview pane's own rows, beside the Workers tab's picture: its head,
// the tabs and Close on one line.
const OVERVIEW_ROWS = 1
// Between the parts of a line: the band's one separator.
const SEP = ' | '
// A filter's choice with one item pressed: out if it was in, in if it was not.
const toggled = (chosen: string[], x: string) => (chosen.includes(x) ? chosen.filter(y => y !== x) : [...chosen, x])
// A pane left with nothing to show (its ticket gone from the board, its state
// cleared) keeps its head: [ Close ] closes it as in every other state.
const nothingPane = ({ Box, Button, Text }: Pick<ElementTable, 'Box' | 'Button' | 'Text'>, head: string, close: { key: string; onPress: () => void }) => (
  <Box flexDirection="column">{headRow({ Box, Button }, head, { ...close, label: t.close })}<Text dimColor>{t.paneNothing}</Text></Box>
)

// ── what this module keeps between events (a reload starts it over) ──────

// The board this session's directory belongs to: undefined until looked
// for, null when there is none.
let board: { dir: string; repo: string } | null | undefined
// Who the person in this session is, by `jaira whoami`: their name and
// aliases; undefined until asked, and when jaira could not say.
let me: string[] | undefined
let meAsked = false
// The jaira on PATH, by `jaira --version`: read once a session and again
// after an update from the band; undefined for a dev build.
let installed: string | undefined
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
// copy is the ticket file pulled into the main checkout for the chat
// (ticketHere), which the chat removes itself as its last step; a minute
// after the chat is done (doneAt) the band names it if it is still there.
const chats = new Map<string, { since: number; copy?: string; doneAt?: number }>()
const COPY_CHECK_MS = 60_000
// The Workers tab's frames, running while it is open and shown, and the
// board the Overview was opened on.
let frames: { cancel: () => void } | undefined
let overviewDir: string | undefined
// Whether the Overview pane is on screen: placed by its open, until it closes.
let overviewOpen = false
// Whether the Workers tab is reading the board right now: a read that takes
// longer than READ_FRAMES is not started a second time beside it.
let reading = false
// When reap.sh last ran from here: it asks the forge about every ticket that
// has a tab, so it runs at most every REAP_EVERY_MS, not with each refresh.
let reapedAt: number | undefined
const REAP_EVERY_MS = 10 * 60_000

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

// Native Windows often has no HOME, only USERPROFILE, and spells it with
// backslashes: forward ones are read alike by Windows and by Git Bash.
async function home($: EngineInterface): Promise<string> {
  return ((await $.env.get('HOME')) || (await $.env.get('USERPROFILE')) || '').replaceAll('\\', '/')
}

// Where jaira keeps the person's own files, as jaira finds it: JAIRA_HOME,
// else ~/.jaira.
async function jairaHome($: EngineInterface): Promise<string> {
  return (await $.env.get('JAIRA_HOME')) || `${await home($)}/.jaira`
}

// ── the noise filter ─────────────────────────────────────────────────────

// Every read-modify-write of the kept reports runs one after the other:
// keep() deletes a key before it sets it again, and a reader between the two
// sees nothing — a writer there writes that nothing back, and the reports are
// gone. Other processes sharing the store are beyond this lock.
let reportsLock: Promise<unknown> = Promise.resolve()
function locked<T>(work: () => Promise<T>): Promise<T> {
  const run = reportsLock.then(work, work)
  reportsLock = run.catch(() => undefined)
  return run
}

function keep($: EngineInterface, ticket: string, report: Report): Promise<void> {
  return locked(async () => {
    const root = await project($)
    const key = keyOf(root, ticket)
    const before = await stored($, key)
    // Delete first so the key moves to the end of the store's insertion order:
    // the first keys are then the tickets heard from longest ago.
    await $.store.delete(key)
    await $.store.set(key, [{ ...report, text: clipReport(report.text) }, ...before].slice(0, REPORTS_PER_TICKET))
    const tickets = projectKeys(await $.store.keys(), root)
    for (const old of tickets.slice(0, Math.max(0, tickets.length - TICKETS))) await $.store.delete(old)
  })
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

async function stored($: EngineInterface, key: string): Promise<Report[]> {
  return ((await $.store.get(key)) as Report[] | undefined) ?? []
}

// Read under the lock too, or a read between keep()'s delete and set finds
// the ticket empty. The project is taken inside it: read and write one key.
async function reports($: EngineInterface, args: string): Promise<string> {
  const { ticket, all } = reportArgs(args)
  return locked(async () => {
    const root = await project($)
    if (ticket === '') {
      const keys = projectKeys(await $.store.keys(), root).reverse()
      return listReports(await Promise.all(keys.map(async k => [k, await stored($, k)] as const)))
    }
    return showReport(ticket, await readReports($, keyOf(root, ticket)), all)
  })
}

// A ticket's kept reports, marked read: shown once, the next report call
// brings only what arrived after it. Call it under the lock. The mark is
// best-effort: a store that refuses the write (full, I/O) still shows the
// reports, which then come again on the next read.
async function readReports($: EngineInterface, key: string): Promise<Report[]> {
  const kept = await stored($, key)
  if (kept.some(r => r.seen !== true)) {
    try {
      await $.store.set(key, kept.map(r => ({ ...r, seen: true })))
    } catch {}
  }
  return kept
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

// Workers and reviews are the agents that run: Herdr tabs (none where there
// is no Herdr or its list failed), and this session's own subagents.
async function agentsOf($: EngineInterface, dir: string): Promise<{ bin: string | undefined; tabs: Tab[]; listed: boolean; agents: Agent[] }> {
  const bin = await herdr($)
  const listed = bin === undefined ? undefined : await run($, [bin, 'tab', 'list'], dir)
  let agents: Agent[] = []
  try {
    agents = await $.agent.list()
  } catch {}
  return { bin, tabs: listed?.exitCode === 0 ? tabsOf(listed.stdout) : [], listed: listed?.exitCode === 0, agents }
}

async function refreshBoard($: EngineInterface): Promise<void> {
  if (await carries($)) return
  if (board === undefined) {
    board = await findBoard($)
    if (board !== null) await checkVersion($, board.dir)
  }
  if (board === null) {
    await update($, boardAtom, () => null)
    return
  }
  await refreshRelease($)
  const [lanes, list, boards] = await Promise.all([
    run($, ['jaira', 'lanes', '--json'], board.dir),
    run($, ['jaira', 'list', '--json'], board.dir),
    run($, ['jaira', 'boards', '--json'], board.dir),
  ])
  if (lanes?.exitCode !== 0 || list?.exitCode !== 0) return
  const { bin, tabs, agents } = await agentsOf($, board.dir)
  if (!meAsked) {
    meAsked = true
    const who = await run($, ['jaira', 'whoami', '--json'], board.dir)
    me = who?.exitCode === 0 ? whoOf(who.stdout) : undefined
  }
  const view: BoardView = snapshotOf(lanes.stdout, list.stdout, tabs, agents, me)
  // Every other board of the repository, each by its own lanes: a ticket
  // there waits in a lane today's board may not have. A board whose read
  // fails — a jaira whose lanes take no --board — is left out, not the band.
  const b = board
  for (const name of boards?.exitCode === 0 ? otherBoardsOf(boards.stdout) : []) {
    const [l, li] = await Promise.all([
      run($, ['jaira', 'lanes', '--board', name, '--json'], b.dir),
      run($, ['jaira', 'list', '--board', name, '--json'], b.dir),
    ])
    if (l?.exitCode !== 0 || li?.exitCode !== 0) continue
    const other = snapshotOf(l.stdout, li.stdout, tabs, agents, me)
    view.waiting.push(...other.waiting)
    view.theirs.push(...other.theirs)
    view.counts.worker += other.counts.worker
    view.counts.review += other.counts.review
    view.counts.person += other.counts.person
  }
  await update($, boardAtom, () => view)
  // What an Explain chat wrote back onto the tickets the band shows, Explain
  // or not: an explanation stays readable after a later round, only asking
  // for a new one waits for that round's verdict.
  const shownIds = (await read($, expandedAtom)) ? view.waiting : view.waiting.slice(0, BAND_TICKETS)
  const explained: Record<string, { text: string; decided?: string }> = {}
  // The newest bug hunt of each, as jaira works it out against the checkout
  // the ticket's code is in: its own worktree, where there is one.
  const hunts: Record<string, HuntState> = {}
  const trees = shownIds.length > 0 ? await run($, ['git', 'worktree', 'list', '--porcelain'], board.dir) : undefined
  for (const w of shownIds) {
    const shown = await run($, ['jaira', 'show', w.id, '--json'], board.dir)
    const notes = shown?.exitCode === 0 ? explanationNotes(shown.stdout) : undefined
    if (notes !== undefined) explained[w.id] = notes
    const dir = trees?.exitCode === 0 ? worktreeFor(trees.stdout, w.id, board.repo) : board.repo
    const brief = await run($, ['jaira', 'show', w.id, '--brief', '--json'], dir)
    const hunt = brief?.exitCode === 0 ? huntOf(brief.stdout) : undefined
    if (hunt === undefined) continue
    // jaira says fail, too, for a checkout without the hunt's commit — the
    // main one, where the ticket's worktree is gone or named otherwise. That
    // fail says nothing about the code: '?', not a red cross.
    if (hunt.state === 'fail' && hunt.at !== '') {
      const held = await run($, ['git', 'merge-base', '--is-ancestor', hunt.at, 'HEAD'], dir)
      if (held?.exitCode !== 0) continue
    }
    hunts[w.id] = hunt
  }
  await update($, explainedAtom, () => explained)
  await update($, huntsAtom, () => hunts)
  // A chat has done its part once its note is on the ticket — stamped to the
  // minute, so from the minute it was opened: an Explain chat's explanation,
  // a Discuss chat's decision, which the band then offers to read. A Discuss
  // chat whose question went away without one is done too. Its tab goes, but
  // only once its agent stopped working: a Discuss chat notes the decision
  // first and then still moves the ticket and tells this session.
  const open = tabs
  const now = await $.clock.now()
  for (const [label, { since, copy, doneAt }] of chats) {
    const [id = '', kind] = label.split('/')
    if (doneAt !== undefined) {
      if (now - doneAt < COPY_CHECK_MS) continue
      chats.delete(label)
      if (copy !== undefined && await $.fs.exists(copy).catch(() => false)) await say($, t.copyLeft(id, copy))
      continue
    }
    const tab = open.find(x => x.label === label)
    if (tab?.agent_status === 'working') continue
    const from = since - (since % 60_000)
    const note = await noteOn($, board.dir, id, kind === 'explain' ? 'explain' : 'decision')
    const fresh = note !== undefined && note.at >= from ? note : undefined
    // A teammate's ticket can be talked through too, and sits in theirs.
    const w = [...view.waiting, ...view.theirs].find(x => x.id === id)
    const done = fresh !== undefined || (kind === 'discuss' && (w === undefined || w.question === ''))
    if (!done) continue
    // The chat removes its copy itself; whether it did is looked at a
    // minute later, when even a chat without a tab has written its last line.
    if (copy === undefined) chats.delete(label)
    else chats.set(label, { since, copy, doneAt: now })
    if (kind === 'discuss' && fresh !== undefined) await update($, decidedAtom, d => ({ ...d, [id]: fresh.text }))
    if (tab !== undefined) await run($, [bin!, 'tab', 'close', tab.tab_id], board.dir)
  }
}

// What spawn.sh opened for a ticket whose pull request is merged — its tabs,
// worktree and branch — goes without a teamlead session: a teamlead runs
// reap.sh every round, a person's session runs it here, from its first minute
// on and then every REAP_EVERY_MS — not at the start itself: a
// tab spawn.sh opens comes up before its role line is typed, and only that
// line says the session carries a ticket. Only with Herdr: without it reap.sh
// keeps everything.
// Its lines say nothing a person has to act on, so none is shown.
async function reapMerged($: EngineInterface): Promise<void> {
  if (await carries($)) return
  if (!board || (await herdr($)) === undefined) return
  const now = await $.clock.now()
  if (reapedAt !== undefined && now - reapedAt < REAP_EVERY_MS) return
  reapedAt = now
  const script = await spawnScript($)
  if (script === undefined) return
  const [sh = 'bash', spawn = ''] = script
  await run($, [sh, spawn.replace(/spawn\.sh$/, 'reap.sh'), board.repo], board.repo, 120_000)
}

// Once a session, when the board is found: a jaira older than the cockpit
// needs is said out loud, with both versions, before a press fails on it.
async function checkVersion($: EngineInterface, dir: string): Promise<void> {
  const r = await run($, ['jaira', '--version'], dir)
  installed = r?.exitCode === 0 ? versionOf(r.stdout) : undefined
  if (installed !== undefined && isOlder(installed, MIN_JAIRA)) await say($, t.tooOld(installed, MIN_JAIRA))
}

// Each minute with the board, as the TUI reads it: the release jaira's own
// cache knows, offered while it is newer than the jaira installed and not the
// one dismissed. The hourly check of the TUI and of next/resume/claim/move
// writes that cache, so the cockpit starts no check of its own.
async function refreshRelease($: EngineInterface): Promise<void> {
  // Off for the cockpit as for the TUI: a cache left from before still names a release.
  if ((await $.env.get('JAIRA_NO_UPDATE_CHECK')) === '1') return update($, releaseAtom, () => null)
  let cache: string | undefined
  try {
    cache = await $.fs.read(`${await jairaHome($)}/update-check.json`)
  } catch {}
  const dismissed = (await $.store.get(DISMISSED_RELEASE)) as string | undefined
  const offer = offeredRelease(installed, cache, dismissed) ?? null
  await update($, releaseAtom, () => offer)
}

// Dismiss: this release is not offered again, in this session or a later
// one; the next one is.
async function dismissRelease($: EngineInterface, version: string): Promise<void> {
  await $.store.set(DISMISSED_RELEASE, version)
  await update($, releaseAtom, () => null)
  await say($, t.dismissed(version))
}

// Update, before anything is written: what is released, and what bringing
// this repository's setup up to date would change. The check only refreshes
// the release cache and the dry run writes nothing; the pane asks before the
// binary is replaced.
async function previewUpgrade($: EngineInterface, dir: string): Promise<void> {
  const r = await run($, ['jaira', 'self', 'upgrade', '--check', '--json'], dir, 30_000)
  const check = r?.exitCode === 0 ? upgradeCheckOf(r.stdout) : undefined
  if (check === undefined) return fail($, t.updateJaira, withoutWarnings(output(r)))
  if (check.upToDate) {
    installed = check.current
    await refreshRelease($)
    return say($, t.upToDate(check.current))
  }
  const dry = await run($, ['jaira', 'update', '--roles', '--lanes', '--dry-run'], dir, 30_000)
  const preview = dry === undefined ? 'jaira did not run' : withoutWarnings(`${dry.stdout}\n${dry.stderr}`)
  // Exit 3 is a preview too: it names the edited files an update leaves
  // alone. Anything else is no preview, and without one nothing is offered.
  if (dry?.exitCode !== 0 && dry?.exitCode !== 3) return fail($, t.updateJaira, `${t.dryRunFailed}\n${preview}`)
  await update($, upgradeAtom, () => ({ from: check.current, to: check.latest, preview }))
  const opened = await $.ui.open({ id: UPGRADE, title: t.upgradeTitle(check.current, check.latest), focus: true, closeOnEscape: true })
  if (!opened.isPlaced) await fail($, t.updateJaira, t.noRoomUpdate)
}

// Update, on the person's yes in that pane: the binary, then this
// repository's skills and lane files — the two commands the line after
// next/resume/claim/move names.
async function upgradeJaira($: EngineInterface, dir: string, to: string): Promise<void> {
  await $.ui.close({ id: UPGRADE })
  await update($, upgradeAtom, () => null)
  await say($, t.upgrading(to))
  const up = await run($, ['jaira', 'self', 'upgrade'], dir, 300_000)
  if (up?.exitCode !== 0) return fail($, t.updateJaira, withoutWarnings(output(up)))
  const setup = await run($, ['jaira', 'update', '--roles', '--lanes'], dir, 60_000)
  await checkVersion($, dir)
  await refreshRelease($)
  await update($, noticeAtom, () => null)
  await show($, t.upgraded(to), withoutWarnings(`${up.stdout}\n${up.stderr}\n${setup === undefined ? 'jaira update did not run' : `${setup.stdout}\n${setup.stderr}`}`))
}

async function noteOn($: EngineInterface, dir: string, id: string, kind: 'explain' | 'decision') {
  const shown = await run($, ['jaira', 'show', id, '--json'], dir)
  return shown?.exitCode === 0 ? noteOf(shown.stdout, kind) : undefined
}

// Where a ticket's work lives: running jaira anywhere else answers "on its
// ref and not on your disk".
async function ticketDir($: EngineInterface, b: { dir: string; repo: string }, id: string): Promise<string> {
  return (await ownTree($, b, id)) ?? b.repo
}

async function ownTree($: EngineInterface, b: { dir: string; repo: string }, id: string): Promise<string | undefined> {
  const r = await run($, ['git', 'worktree', 'list', '--porcelain'], b.dir)
  return r?.exitCode === 0 ? ownTreeFor(r.stdout, id) : undefined
}

// Where a press writes, with the ticket on that disk as its ref has it: a
// file on the ref alone or older than the ref refuses every write ("on its
// ref and not on your disk", a stale file), so it is pulled first. Pulled
// into the ticket's own worktree it stays, as the branch's; pulled into the
// main checkout only for this press, leave() takes it away again once the
// press is over, refused or not: a copy left there outranks the ref and shows
// the ticket in a lane it left long ago. Every write already sits on the ref,
// or whole in the outbox, so nothing goes with it. A refused pull (somebody
// else holds the ticket, no remote) stops the press: it never steals.
// copy names the file pulled into the main checkout, the one leave() removes.
type Here = { ok: true; dir: string; copy?: string; leave: () => Promise<void> } | { ok: false; why: string }

async function ticketHere($: EngineInterface, b: { dir: string; repo: string }, id: string): Promise<Here> {
  const own = await ownTree($, b, id)
  const dir = own ?? b.repo
  const stays: Here = { ok: true, dir, leave: async () => {} }
  const disk = async () => {
    const shown = await run($, ['jaira', 'show', id, '--json'], dir)
    return diskOf(shown?.exitCode === 0 ? shown.stdout : '')
  }
  const before = await disk()
  if (before.state === 'current') return stays
  const pulled = await jaira($, dir, ['pull', id])
  if (!pulled.ok) return { ok: false, why: pulled.why }
  if (own !== undefined || before.state !== 'ref-only') return stays
  const { path } = await disk()
  if (path === '') return stays
  return {
    ok: true,
    dir,
    copy: path,
    leave: async () => {
      const r = await run($, ['git', 'clean', '-fqx', '--', path], dir)
      // A line beside the band, never the message pane: that already holds
      // what the press itself had to say (a refusal, a note that was lost).
      if (r?.exitCode !== 0) await say($, t.copyNotRemoved(id, path, output(r)))
    },
  }
}

type Done = { ok: true; out: string } | { ok: false; why: string }

async function jaira($: EngineInterface, cwd: string, args: string[]): Promise<Done> {
  const r = await run($, ['jaira', ...args], cwd, 30_000)
  if (r === undefined) return { ok: false, why: 'jaira did not run' }
  if (r.exitCode !== 0) return { ok: false, why: withoutWarnings(r.stderr || r.stdout) || `jaira exited ${r.exitCode}` }
  return { ok: true, out: withoutWarnings(r.stdout) }
}

// Whether a press's move found the ticket moved by another session first:
// jaira answers a move to the lane the ticket already sits in with exit 0 and
// '<handle> is already in <lane>' (internal/cli/flow.go), and a second Give
// back or Accept that went on from there would write its note, create its
// follow-ups and hand the work out a second time.
function movedBefore(d: Done, id: string, to: string): boolean {
  return d.ok && d.out.split(/\r?\n/).some(l => l.trim() === `${id} is already in ${to}`)
}

// Someone else's move came first: nothing of this press is written.
async function movedFirst($: EngineInterface, what: string, id: string, to: string): Promise<void> {
  await fail($, `${what} ${id}`, t.alreadyIn(id, to))
  await refreshBoard($)
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

// A ticket waiting for a person, this one's or a teammate's.
async function waitingTicket($: EngineInterface, id: string | null): Promise<Waiting | undefined> {
  const view = await read($, boardAtom)
  return [...(view?.waiting ?? []), ...(view?.theirs ?? [])].find(w => w.id === id)
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
// refused: where a gate's refusal goes instead of the message pane - the
// Acceptance pane shows it beside the choices it keeps. True once the ticket
// moved.
type Refused = (why: string) => Promise<void>

async function accept($: EngineInterface, w: Waiting, after?: (dir: string) => Promise<void>, refused?: Refused): Promise<boolean> {
  if (!board || w.acceptTo === null) return false
  const here = await ticketHere($, board, w.id)
  if (!here.ok) {
    await fail($, `${t.accept} ${w.id}`, here.why)
    return false
  }
  try {
    return await accepting($, w, w.acceptTo, here.dir, after, refused)
  } finally {
    await here.leave()
  }
}

async function accepting($: EngineInterface, w: Waiting, to: string, dir: string, after?: (dir: string) => Promise<void>, onRefused?: Refused): Promise<boolean> {
  const dry = await run($, ['jaira', 'move', w.id, '--to', to, '--dry-run', '--json'], dir, 30_000)
  const refused = dry === undefined ? 'jaira did not run' : dry.exitCode === 0 ? undefined : otherRefusal(dry.stderr || dry.stdout)
  if (refused !== undefined) {
    await (onRefused !== undefined ? onRefused(refused) : fail($, `${t.accept} ${w.id}`, refused))
    return false
  }
  const done = await jaira($, dir, ['move', w.id, '--to', to, '--force'])
  if (movedBefore(done, w.id, to)) {
    await movedFirst($, t.accept, w.id, to)
    return false
  }
  if (done.ok) {
    // What the person decided on the review's notes goes onto the ticket once
    // it is accepted, as Give back writes it once the ticket went back: a
    // refused or failed Accept writes nothing and creates no follow-up.
    if (after !== undefined) await after(dir)
    await say($, t.accepted(w.id, to))
    const next = await afterAccept($, w.id, to)
    // One line into this session, the teamlead's, once the press has done
    // all it does: it goes on from there — logbook, the next worker — without
    // the person typing it, and never starts the pull request Accept started.
    void $.prompt.submit({ text: t.toSessionAccepted(w.id, w.lane, to, next) })
  } else await fail($, `${t.accept} ${w.id}`, done.why)
  await refreshBoard($)
  return done.ok
}

// What a command said when it failed, for the message pane: its error, else
// its output, else that it never ran.
function output(r: Awaited<ReturnType<typeof run>>): string {
  return (r?.stderr || r?.stdout || 'did not run').trim()
}

// The repository a pull request of this board goes to and the branch it
// goes into, fetched: `upstream/main` in a fork, `origin/master` otherwise —
// what "ahead" is counted against. undefined where origin is missing or no
// forge can be told (the caller counts against the local base then, and a
// pull request cannot be opened anyway); a string says why it failed.
type Target = { on: Forge; repo: string; base: string }

// The forge of origin and the repository its pull requests go to - the
// parent in a fork; undefined where no forge can be told, a string where the
// forge did not answer.
async function forgeRepo($: EngineInterface, repoDir: string): Promise<{ on: Forge; repo: string } | string | undefined> {
  // Unset, `git config` exits 1 and the host of origin decides.
  const forge = await run($, ['git', 'config', 'jaira.forge'], repoDir)
  const origin = await run($, ['git', 'remote', 'get-url', 'origin'], repoDir)
  const url = origin?.exitCode === 0 ? origin.stdout.trim() : ''
  const on = url ? forgeOf(forge?.exitCode === 0 ? forge.stdout : '', url) : undefined
  if (on === undefined) return undefined
  const view = await run($, repoViewCommand(on, url), repoDir, 30_000)
  const repo = view?.exitCode === 0 ? targetRepoOf(on, view.stdout) : undefined
  return repo === undefined ? `${on === 'github' ? 'gh' : 'glab'} repo view: ${output(view)}` : { on, repo }
}

async function pullTarget($: EngineInterface, repoDir: string): Promise<Target | string | undefined> {
  const found = await forgeRepo($, repoDir)
  if (found === undefined || typeof found === 'string') return found
  const { on, repo } = found
  const tool = on === 'github' ? 'gh' : 'glab'
  const targetView = await run($, repoViewCommand(on, repo), repoDir, 30_000)
  const branch = targetView?.exitCode === 0 ? defaultBranchOf(on, targetView.stdout) : undefined
  if (branch === undefined) return `${tool} repo view ${repo}: ${(targetView?.stderr || 'no default branch').trim()}`
  const remotes = await run($, ['git', 'remote', '-v'], repoDir)
  const remote = remotes?.exitCode === 0 ? remoteFor(remotes.stdout, repo) : undefined
  if (remote === undefined) return t.noRemoteFor(repo)
  const fetched = await run($, ['git', 'fetch', remote, branch], repoDir, 60_000)
  if (fetched?.exitCode !== 0) return `git fetch ${remote} ${branch}: ${output(fetched)}`
  return { on, repo, base: `${remote}/${branch}` }
}

// Which step the press starts, and on which branch: the ticket's branch
// found, counted against the pull request's target, its tickets read, its
// pull requests asked for — stepAfterAccept (board.ts) decides from that.
// A forge that cannot be asked leaves the pull request unchecked: the rule
// gets undefined for it, and the hint says why no pull request was looked
// at. A string is a failure: no pull request is opened, and the pane says why.
//
// A board whose .jaira/pr says none (XB97NJ) takes no pull requests: the pr
// lane lands the branch instead, behind 'jaira roles pr-gate'. Landing needs
// a remote, not a forge: a host that names no forge lands on origin, as
// jaira-role-pr does there. What the role cannot land past stops here too,
// with its reason — no origin at all, a forge that does not answer (whether
// origin is a fork cannot be told), a fork whose parent no remote points at.
// The step is the same 'pr' and `land` says which of the two it is.
type Decided = { step: StepAfterAccept; branch?: string; hint: string; land?: boolean }

async function decideStep($: EngineInterface, b: { dir: string; repo: string }, id: string, terminal: string, porcelain: string): Promise<Decided | string> {
  const found = branchesFor(porcelain, id)
  const asked = await pullTarget($, b.repo)
  // Under none a reason is held until there is something to land: a branch
  // with nothing ahead, or a ticket still open, is no failure.
  const target = typeof asked === 'string' ? undefined : asked
  // Without a forge to name it, the base is origin's default branch — what
  // 'pr-gate --land' and jaira-role-pr read — and the local one only where
  // origin has none to show.
  const base = target?.base ?? await originDefault($, b.repo) ?? found.base
  // The base decides, as the gate does: this checkout's .jaira/pr may be the
  // branch's own word for itself.
  const land = base !== undefined && await landsWithoutPullRequest($, b.repo, base)
  if (typeof asked === 'string' && !land) return { step: stepAfterAccept(0, [], undefined), hint: t.prUnchecked(asked) }
  let branch = found.ticket
  // No worktree: the branch spawn.sh named <prefix>/<slug>, left behind
  // when the worktree went; else the branch whose commits name the ticket,
  // where it was worked on another ticket's branch.
  if (branch === undefined) {
    const listed = await run($, ['git', 'branch', '--list', '--format=%(refname:short)', `*/${id.toLowerCase()}`], b.repo)
    if (listed?.exitCode === 0) branch = carrierBranch(listed.stdout, found.base ?? '')
  }
  if (branch === undefined && base !== undefined) {
    const log = await run($, ['git', 'log', '--branches', '--not', base, '--format=%H %s'], b.repo)
    const sha = log?.exitCode === 0 ? commitNaming(log.stdout, id) : undefined
    const holders = sha ? await run($, ['git', 'branch', '--contains', sha, '--format=%(refname:short)'], b.repo) : undefined
    if (holders?.exitCode === 0) branch = carrierBranch(holders.stdout, found.base ?? base)
  }
  if (branch === undefined || base === undefined || branch === found.base) return { step: 'none', hint: '' }
  // One subject per commit ahead: the count and the tickets in one walk.
  const subjects = await run($, ['git', 'log', '--format=%s', `${base}..${branch}`], b.repo)
  if (subjects?.exitCode !== 0) return `git log ${base}..${branch}: ${output(subjects)}`
  const ahead = subjects.stdout.split('\n').filter(Boolean).length
  if (ahead === 0) return { step: 'none', branch, hint: '' }
  const list = await run($, ['jaira', 'list', '--json'], b.dir, 30_000)
  if (list?.exitCode !== 0) return `jaira list: ${withoutWarnings(list?.stderr ?? '') || 'did not run'}`
  const pending = pendingOf(list.stdout, ticketsNamed(subjects.stdout), terminal)
  if (pending.length > 0) return { step: stepAfterAccept(ahead, pending, false), branch, hint: land ? t.landWaits(branch, pending) : t.prWaits(branch, pending), land }
  if (land && typeof asked === 'string') return t.cannotLand(asked)
  if (land && target === undefined) {
    const origin = await run($, ['git', 'remote', 'get-url', 'origin'], b.repo)
    if (origin?.exitCode !== 0) return t.cannotLand(t.noOrigin)
  }
  if (land) return { step: stepAfterAccept(ahead, pending, false), branch, hint: '', land }
  if (target === undefined) return t.noForge(branch)
  const cmd = pullRequestsCommand(target.on, target.repo, branch)
  const listed = await run($, cmd, b.repo, 30_000)
  const has = listed?.exitCode === 0 ? hasPullRequest(listed.stdout) : undefined
  return { step: stepAfterAccept(ahead, pending, has), branch, hint: has === undefined ? t.prUnchecked(`${cmd[0]}: ${output(listed)}`) : '' }
}

// What the press starts once the ticket is in done, without a further word
// from the person: a pull request, where stepAfterAccept (board.ts) says the
// ticket's branch needs one, and nothing else. A pull request is opened only
// from here — a press of Accept is a person's — by the pr lane of spawn.sh.
// The next ticket is never started from here: which one runs next is the
// teamlead's call (X6JHQX), so without a pull request to open the press adds
// nothing to its "accepted" but why none was opened, where there is a reason. A step that cannot run
// starts nothing and says why in the message pane; the ticket stays in done,
// where the press put it: accepted is accepted, and what follows is the
// person's to start by hand then. It answers what it said, for the line
// Accept hands this session; '' where it said nothing.
async function afterAccept($: EngineInterface, id: string, terminal: string): Promise<string> {
  if (!board) return ''
  const b = board
  const trees = await run($, ['git', 'worktree', 'list', '--porcelain'], b.dir)
  const decided = await decideStep($, b, id, terminal, trees?.exitCode === 0 ? trees.stdout : '')
  const stuck = async (why: string) => {
    await fail($, t.afterAccept(id), `${why}\n${t.staysDone(id)}`)
    return t.afterAcceptFailed(id)
  }
  const said = async (text: string) => {
    await say($, text)
    return text
  }
  if (typeof decided === 'string') return stuck(decided)
  const { step, branch, hint, land } = decided
  // No pull request to open: only why not, where there is a reason.
  if (step !== 'pr' || branch === undefined) return hint ? said(t.acceptedNoPr(id, hint)) : ''
  const bin = await herdr($)
  const script = await spawnScript($)
  if (bin === undefined || script === undefined) return stuck(land ? t.landNeedsHerdr(branch) : t.prNeedsHerdr(branch))
  // In the worktree that carries the branch — the role pushes from there —
  // checked out again beside the repository where it went after done.
  let wt = trees?.exitCode === 0 ? worktreeOfBranch(trees.stdout, branch) : undefined
  if (wt === undefined) {
    wt = newWorktreeFor(b.repo, branch)
    const added = await run($, ['git', 'worktree', 'add', wt, branch], b.repo, 60_000)
    if (added?.exitCode !== 0) return stuck(`git worktree add ${wt} ${branch}: ${output(added)}`)
  }
  const r = await run($, [...script, '--no-worktree', id.toLowerCase(), id, 'pr', wt], wt, 120_000)
  if (r?.exitCode !== 0) return stuck(`spawn.sh: ${output(r).slice(-160)}`)
  return said(land ? t.landStarted(id, branch) : t.prStarted(id, branch))
}

// origin's default branch as a remote-tracking name: origin/HEAD where it
// resolves, else origin/main or origin/master; undefined when origin has
// none of them here.
async function originDefault($: EngineInterface, repo: string): Promise<string | undefined> {
  const head = await run($, ['git', 'rev-parse', '--abbrev-ref', '--verify', '-q', 'refs/remotes/origin/HEAD'], repo)
  if (head?.exitCode === 0 && head.stdout.trim()) return head.stdout.trim()
  const refs = await run($, ['git', 'for-each-ref', '--format=%(refname:short)', 'refs/remotes/origin/main', 'refs/remotes/origin/master'], repo)
  if (refs?.exitCode !== 0) return undefined
  return refs.stdout.split('\n').map(l => l.trim()).find(l => l === 'origin/main') ?? refs.stdout.split('\n').map(l => l.trim()).find(l => l === 'origin/master')
}

// Whether .jaira/pr says none as base holds it, read as jaira reads it: the
// first line that is neither blank nor a comment, a byte order mark dropped,
// exactly 'none'. No such file in base, or anything else in it, is the pull
// request path, as before. Read from the base, as 'jaira roles pr-gate
// --land' reads it, never off this checkout's disk.
async function landsWithoutPullRequest($: EngineInterface, repo: string, base: string): Promise<boolean> {
  const shown = await run($, ['git', 'show', `${base}:.jaira/pr`], repo)
  if (shown?.exitCode !== 0) return false
  const word = shown.stdout.replace(/^\uFEFF/, '').split('\n').map(l => l.trim()).find(l => l !== '' && !l.startsWith('#'))
  return word === 'none'
}

// The person's jaira settings, where jaira reads them: JAIRA_HOME, else
// ~/.jaira. undefined where there is no file.
async function settingsJson($: EngineInterface): Promise<string | undefined> {
  try {
    return await $.fs.read(`${await jairaHome($)}/settings.json`)
  } catch {
    return undefined
  }
}

// A press on Accept, Give back, Explain or Discuss: the row's buttons lock, drawn dim,
// until the action has done its part or BUSY_MS has passed; a press on a
// locked row does nothing, so a slow answer is not asked for twice. An action
// that answers 'hold' is done at once but keeps the row locked to BUSY_MS: a
// request sent to this session has no end to wait for, and a double click
// would send it twice.
async function once($: EngineInterface, id: string, action: () => Promise<void | 'hold'>): Promise<void> {
  const at = await $.clock.now()
  let took = false
  await update($, busyAtom, busy => {
    took = busy[id] === undefined
    return took ? { ...busy, [id]: at } : busy
  })
  if (!took) return
  const free = () => update($, busyAtom, busy => {
    if (busy[id] !== at) return busy
    const { [id]: _, ...rest } = busy
    return rest
  })
  $.clock.after(BUSY_MS, () => void free())
  let held = false
  try {
    held = (await action()) === 'hold'
  } finally {
    if (!held) await free()
  }
}

// The spawn.sh the roles were installed with: the project's own first, as
// the dispatcher's run-lane.sh finds it beside itself. Started through bash,
// not by its name: Windows runs no .sh by itself.
async function spawnScript($: EngineInterface): Promise<string[] | undefined> {
  const tail = '.claude/skills/jaira-dispatcher/scripts/spawn.sh'
  for (const dir of [board?.repo, await home($)]) {
    if (!dir) continue
    if (await $.fs.exists(`${dir}/${tail}`).catch(() => false)) return [await bash($), `${dir}/${tail}`]
  }
  return undefined
}

// On native Windows the Git Bash spawn.sh is written for: a bare `bash` there
// is System32's, which starts WSL. Claude Code names Git Bash in
// CLAUDE_CODE_GIT_BASH_PATH; else it sits three levels above git's exec path
// (<Git>/mingw64/libexec/git-core -> <Git>/bin/bash.exe).
async function bash($: EngineInterface): Promise<string> {
  if ((await $.env.get('OS')) !== 'Windows_NT') return 'bash'
  const named = await $.env.get('CLAUDE_CODE_GIT_BASH_PATH')
  if (named) return named
  const r = await run($, ['git', '--exec-path'], board?.repo ?? '.')
  const git = r?.exitCode === 0 ? r.stdout.trim().replaceAll('\\', '/').replace(/\/[^/]+\/[^/]+\/[^/]+$/, '') : ''
  return git ? `${git}/bin/bash.exe` : 'bash'
}

// Herdr App (herdr-core) sets no HERDR_ENV and has no herdr binary: inside it
// (HC_TERMINAL_ID set) the dispatcher's herdr-hc.sh carries out the same calls
// through hc, beside the spawn.sh this mod starts.
async function herdr($: EngineInterface): Promise<string | undefined> {
  if ((await $.env.get('HERDR_ENV')) === '1') return (await $.env.get('HERDR_BIN_PATH')) || 'herdr'
  if (!(await $.env.get('HC_TERMINAL_ID'))) return undefined
  const script = await spawnScript($)
  return script?.[1]?.replace(/spawn\.sh$/, 'herdr-hc.sh')
}

// Explain a review finding, or talk an open question through, in a Herdr tab
// of its own: the teamlead's chat is not where that conversation belongs.
// Without Herdr this session is asked instead: Explain by a subagent, which
// writes the same note a tab would; Discuss here, since a subagent cannot
// talk with a person.
async function openChat($: EngineInterface, id: string, kind: 'explain' | 'discuss'): Promise<void | 'hold'> {
  const what = `${kind === 'explain' ? t.explain : t.discuss} ${id}`
  if (!board) return fail($, what, t.noBoard)
  // One chat at a time on a ticket a chat holds a copy for: that chat removes
  // the copy when it ends, from under any other chat still writing there —
  // held until the copy was looked at, a minute after the chat was done.
  const holder = [...chats].find(([l, c]) => l.startsWith(`${id}/`) && c.copy !== undefined)
  if (holder !== undefined) return fail($, what, t.chatHoldsCopy(id, holder[0].split('/')[1] === 'explain' ? t.explain : t.discuss))
  // The chat writes its note onto the ticket, which a file on the ref alone
  // or older than the ref refuses: pulled first, as for Accept (76ERHH). A
  // refused pull stops the press; it never steals.
  const here = await ticketHere($, board, id)
  if (!here.ok) return fail($, what, here.why)
  const { dir } = here
  const label = `${id}/${kind}`
  // A copy pulled into the main checkout is the chat's to remove, as its last
  // step: only the chat knows when it wrote its last line. A press that opens
  // no chat removes it at once.
  const { copy } = here
  const opened = async () => chats.set(label, { since: await $.clock.now(), copy })
  if ((await herdr($)) === undefined) {
    await opened()
    await say($, kind === 'explain' ? t.explainAsked(id) : t.discussHere(id))
    void $.prompt.submit({ text: kind === 'explain' ? t.explainBySubagent(id, dir, copy) : t.discussInThisChat(id, dir, copy) })
    return 'hold'
  }
  const script = await spawnScript($)
  if (script === undefined) {
    await here.leave()
    return fail($, what, t.noSpawn)
  }
  // A Discuss chat sends its decision to this session: its messaging socket
  // is the address SendMessage takes.
  const socket = kind === 'discuss' ? await $.env.get('CLAUDE_CODE_MESSAGING_SOCKET') : undefined
  const parent = socket ? ['--parent', `uds:${socket}`] : []
  const removes = copy === undefined ? [] : ['--copy', copy]
  // Registered before spawn.sh runs, which waits up to two minutes for
  // claude while the row unlocks after BUSY_MS: a second press in that time
  // already finds this chat holding the copy.
  await opened()
  const r = await run($, [...script, '--no-worktree', ...parent, ...removes, id.toLowerCase(), id, kind, dir], dir, 120_000)
  if (r?.exitCode !== 0) {
    // spawn.sh may have failed after its tab came up (claude slow to start):
    // that chat runs and removes the copy itself, so it is not taken from
    // under it. No such tab: no chat started, and the copy goes at once.
    const tabs = await agentsOf($, dir)
    if (!(copy !== undefined && tabs.listed && tabs.tabs.some(x => x.label === label))) {
      chats.delete(label)
      await here.leave()
    }
    return fail($, what, (r?.stderr || r?.stdout || 'spawn.sh did not run').trim())
  }
  await say($, t.opened(id))
}

// '+N more' opens every waiting ticket in the band, with its buttons; the
// band scrolls where they do not fit. Opened, the explanations of all of
// them are read too.
async function toggleMore($: EngineInterface, open: boolean): Promise<void> {
  await update($, expandedAtom, () => open)
  if (open) await refreshBoard($)
}

// The Overview pane, on the board as it is now: read afresh at each press.
// A command that failed leaves the pane saying so, not half a board.
async function openOverview($: EngineInterface, dir: string): Promise<void> {
  overviewDir = dir
  const [lanes, list, next] = await Promise.all([
    run($, ['jaira', 'lanes', '--json'], dir),
    run($, ['jaira', 'list', '--json'], dir),
    run($, ['jaira', 'next', '--all', '--json'], dir),
  ])
  const ok = lanes?.exitCode === 0 && list?.exitCode === 0 && next?.exitCode === 0
  const view = ok ? overviewOf(lanes.stdout, list.stdout, next.stdout, await $.clock.now()) : null
  await update($, overviewAtom, () => view)
  const opened = await $.ui.open({ id: OVERVIEW, title: overviewWords.title, focus: true, closeOnEscape: true })
  // Unplaced, the pane is not drawn and hears no close: no frames for it.
  overviewOpen = opened.isPlaced
  if (!opened.isPlaced) {
    $.ui.toast(overviewWords.noRoom)
    return stopFrames()
  }
  await showTab($, await read($, overviewTabAtom))
}

// A tab of the Overview pane; Workers reads who stands where and starts its
// frames, every other tab stops them.
async function showTab($: EngineInterface, tab: OverviewTab): Promise<void> {
  await update($, overviewTabAtom, () => tab)
  if (tab !== 'workers') return stopFrames()
  await readWorkers($)
  // The person may have shown another tab or closed the pane while it read.
  if (!overviewOpen || (await read($, overviewTabAtom)) !== 'workers') return
  frames ??= $.clock.every(FRAME_MS, () => void frame($))
}

function stopFrames(): void {
  frames?.cancel()
  frames = undefined
}

async function frame($: EngineInterface): Promise<void> {
  const n = await update($, workersFrameAtom, f => f + 1)
  if (n % READ_FRAMES === 0) await readWorkers($)
}

// How many frames ago the stone lost its last chunk; -1 when none ever fell.
async function fallingOf($: EngineInterface): Promise<number> {
  const fellAt = await read($, workersFellAtom)
  return fellAt < 0 ? -1 : (await read($, workersFrameAtom)) - fellAt
}

// The first line of the last report kept for each ticket, under the
// reports' lock like every other read of them; marks none read.
// A store that cannot be read leaves the lines without reports, not the tab
// without agents.
async function lastReports($: EngineInterface, ids: readonly string[]): Promise<Record<string, string>> {
  return locked(async () => {
    const lines: Record<string, string> = {}
    try {
      const root = await project($)
      for (const id of ids) {
        const last = (await stored($, keyOf(root, id)))[0]
        const line = typeof last?.text === 'string' ? reportLine(last.text) : ''
        if (line !== '') lines[id] = line
      }
    } catch {}
    return lines
  })
}

// Who is on the board: its tickets, Herdr's tabs and this session's
// subagents. A ticket seen before keeps its colour and its hat, an agent
// the time it was first seen; one gone drops out. The stone loses a chunk
// when the open work at it is less than at the read before.
async function readWorkers($: EngineInterface): Promise<void> {
  const dir = overviewDir
  if (dir === undefined || reading) return
  reading = true
  try {
    const [lanes, list] = await Promise.all([run($, ['jaira', 'lanes', '--json'], dir), run($, ['jaira', 'list', '--json'], dir)])
    if (lanes?.exitCode !== 0 || list?.exitCode !== 0) {
      await update($, workersAtom, () => null)
      return
    }
    const { tabs, agents } = await agentsOf($, dir)
    const workers = workersOf(lanes.stdout, list.stdout, tabs, agents)
    const ids = [...new Set(workers.map(f => f.id))]
    const keys = [...new Set(workers.map(workerKey))]
    const now = await $.clock.now()
    const frameNow = await read($, workersFrameAtom)
    const before = await read($, workersAtom)
    const reports = await lastReports($, ids)
    await update($, workersReportsAtom, () => reports)
    await update($, workersAtom, () => workers)
    await update($, workersColorsAtom, kept => pickOf(ids, kept, COLORS))
    await update($, workersHatsAtom, kept => pickOf(ids, kept, Object.keys(HATS)))
    await update($, workersSinceAtom, kept => Object.fromEntries(keys.map(k => [k, kept[k] ?? now])))
    if (before !== null && stoneOf(workers) < stoneOf(before)) await update($, workersFellAtom, () => frameNow)
  } finally {
    reading = false
  }
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
  // One pull for the note and the press it carries out: the copy pulled for
  // the note is the one Accept and Give back then find on the disk.
  const here = await ticketHere($, board, id)
  if (!here.ok) return fail($, `${t.decide} ${id}`, here.why)
  try {
    await deciding($, id, text, action, share, here.dir)
  } finally {
    await here.leave()
  }
}

async function deciding($: EngineInterface, id: string, text: string, action: Action | undefined, share: number, dir: string): Promise<void> {
  const done = await jaira($, dir, ['note', id, `decision: ${text}`])
  if (!done.ok) return fail($, `${t.decide} ${id}`, done.why)
  await update($, decidedAtom, d => ({ ...d, [id]: text }))
  // Answered from here on, before the next look at the notes says so too.
  await update($, explainedAtom, e => (e[id] === undefined ? e : { ...e, [id]: { ...e[id], decided: text } }))
  const w = (await read($, boardAtom))?.waiting.find(x => x.id === id)
  if (action === 'giveBack' && w !== undefined) return giveBack($, id, w.lane, text, share)
  if (action === 'accept' && w !== undefined && w.acceptTo !== null) {
    await accept($, w)
    return
  }
  await say($, t.decidedOn(id))
  void $.prompt.submit({ text: t.decisionFor(id, text) })
}

// The Acceptance pane (2V50Y5): one window instead of Accept, Give back and
// Verdict. It opens with focus, so its first note's buttons take the
// keyboard; the band holds it already when the row's button is pressed. The
// head and history are read once, on opening: the ticket, its commits, its
// release note and its pull request do not change while a person reads.
async function openAcceptance($: EngineInterface, w: Waiting): Promise<void> {
  if (!board) return
  await update($, givingAtom, () => w.id)
  await update($, typedAtom, () => '')
  await update($, refusalAtom, () => null)
  // The sentences Explain gave on its notes: asked again on a press, so a
  // changed language or a failed call is not shown on.
  await update($, reasonsAtom, ({ [w.id]: _, ...rest }) => rest)
  const dir = await ticketDir($, board, w.id)
  const shown = await run($, ['jaira', 'show', w.id, '--json'], dir)
  const json = shown?.exitCode === 0 ? shown.stdout : undefined
  const items = itemsOf(w, json)
  await update($, choosingAtom, () => (items.length === 0 ? null : { id: w.id, items }))
  const story = json === undefined ? undefined : storyOf(json)
  await update($, storyAtom, () => ({ id: w.id, ...(story !== undefined ? { story } : {}), release: [], pr: false }))
  const opened = await $.ui.open({ id: ACCEPTANCE, title: t.acceptanceOf(w.id), focus: true, closeOnEscape: true })
  if (!opened.isPlaced) return fail($, t.acceptanceOf(w.id), t.noRoom)
  // The rest fills in while the pane is up, outside the row's lock: Execute
  // pressed meanwhile is not lost to a lock the reading still holds.
  void readStory($, w, story, items)
}

// The Acceptance pane's history from git and the forge, and its texts in the
// person's language: added to what the pane already shows.
async function readStory($: EngineInterface, w: Waiting, story: Story | undefined, items: readonly GapItem[]): Promise<void> {
  const code = await codeDir($, w.id)
  const log = await run($, ['git', 'log', '--format=%h', '--name-only', '-F', `--grep=${w.id}`, 'HEAD'], code)
  let note = ''
  try {
    note = await $.fs.read(`${code}/core/release/notes/${w.id}.md`)
  } catch {}
  const pr = await openPullRequest($, w.id)
  await update($, storyAtom, s => (s?.id !== w.id ? s : { ...s, ...(log?.exitCode === 0 ? { commits: commitsOf(log.stdout) } : {}), release: releaseLines(note), pr }))
  await rephrase($, w.id, items.map(x => x.text))
  const reason = lightOf(w, (await read($, huntsAtom))[w.id], items.map(x => x.cls)).reasons[0]
  await rephrase($, `${w.id}${HEAD}`, [...(story === undefined ? [w.title] : [story.title, story.goal, story.what, story.why]), ...(reason === undefined ? [] : [t.lightReason(reason)])].filter(x => x !== ''), headPrompt)
}

// The checkout the ticket's code is in: its own worktree, where there is one.
async function codeDir($: EngineInterface, id: string): Promise<string> {
  if (!board) return ''
  const trees = await run($, ['git', 'worktree', 'list', '--porcelain'], board.dir)
  return trees?.exitCode === 0 ? worktreeFor(trees.stdout, id, board.repo) : board.repo
}

// Whether the ticket's branch has a pull request, in any state; false where
// nobody can tell - no forge, no branch, a failed call: the history leaves
// it out. Two calls to the forge and no fetch: the pane only looks.
async function openPullRequest($: EngineInterface, id: string): Promise<boolean> {
  if (!board) return false
  const trees = await run($, ['git', 'worktree', 'list', '--porcelain'], board.repo)
  const branch = trees?.exitCode === 0 ? branchesFor(trees.stdout, id).ticket : undefined
  if (branch === undefined) return false
  const found = await forgeRepo($, board.repo)
  if (found === undefined || typeof found === 'string') return false
  const listed = await run($, pullRequestsCommand(found.on, found.repo, branch), board.repo, 30_000)
  return listed?.exitCode === 0 && hasPullRequest(listed.stdout) === true
}

// Explain on one note: one everyday sentence from the cheap model on why it
// matters for accepting the ticket, kept per ticket, note and language so a
// second press asks nothing. A press is the person asking: without a
// language set the sentence comes in the note's own.
const REASON_PREFIX = 'reason:'
async function explainNote($: EngineInterface, id: string, title: string, text: string): Promise<void> {
  // Asked already and waiting: nothing twice. A sentence shown, or a failure,
  // is asked again - from the store where it is kept, in today's language.
  if ((await read($, reasonsAtom))[id]?.[text] === '') return
  const put = (said: string | undefined) => update($, reasonsAtom, r => {
    const { [text]: _, ...rest } = r[id] ?? {}
    return { ...r, [id]: said === undefined ? rest : { ...rest, [text]: said } }
  })
  await put('')
  const language = languageOf(await settingsJson($))
  const key = `${REASON_PREFIX}${id}`
  const kept = ((await $.store.get(key)) ?? {}) as Record<string, string>
  const keyOfNote = `${language ?? ''}\n${text}`
  if (kept[keyOfNote] !== undefined) return void await put(kept[keyOfNote])
  let said: string | undefined
  try {
    const result = await $.model.complete({ model: 'haiku', prompt: reasonPrompt(language, title, text), maxTokens: 300 })
    said = result.isAnswered ? reasonOf(result.text) : undefined
  } catch {}
  await put(said ?? t.noReason)
  if (said !== undefined) await storeNewest($, REASON_PREFIX, key, { ...kept, [keyOfNote]: said })
}

// Execute: what the person chose, carried out as executionOf says. The pane
// stays open, the choices in it, until the ticket really moved: a gate that
// refuses done is shown in it, and nothing is written.
// canType: whether the surface has a field to type a reason into (mobile has none).
async function execute($: EngineInterface, w: Waiting, typed: string, share: number, canType = true): Promise<void> {
  const c = await read($, choosingAtom)
  const items = c?.id === w.id ? c.items : []
  const e = executionOf(items, typed, w.acceptTo)
  const refuse = async (text: string) => void await update($, refusalAtom, () => ({ id: w.id, text }))
  await update($, refusalAtom, () => null)
  if (e.step === 'notTerminal') return refuse(t.notTerminal(w.id, w.lane))
  if (e.step === 'needWhy') return refuse(canType ? t.blockerNeedsWhy(w.id) : t.blockerNeedsTyping(w.id))
  if (e.step === 'giveBack') {
    closeAcceptance($)
    if (items.length === 0) return giveBack($, w.id, w.lane, typed.trim(), share)
    return giveBackWith($, w, items, typed.trim(), share)
  }
  const moved = await acceptWith($, w, items, e.asksWhy ? typed : '', why => refuse(`${t.refused(w.id)}\n${why}`))
  if (moved) closeAcceptance($)
}

function closeAcceptance($: EngineInterface): void {
  void update($, givingAtom, () => null)
  void update($, choosingAtom, () => null)
  void update($, typedAtom, () => '')
  void update($, refusalAtom, () => null)
  void $.ui.close({ id: ACCEPTANCE })
}

// ── review notes in the person's language ────────────────────────────────

// The Verdict pane's order: what stops the ticket first.
const CLASS_ORDER: readonly NoteClass[] = ['blocker', 'edge', 'limit', 'elsewhere']
const REPHRASED_PREFIX = 'rephrased:'
const REPHRASED_TICKETS = 50
// The newest rephrase of each ticket: an older one still waiting for the
// model shows nothing once it is answered. No call is held back for one
// that waits - the newest press wins, a double press may ask twice
// (decision 22.11.2026: a rule that skips calls kept growing edges).
const newest = new Map<string, number>()
let rephrases = 0

// The review notes of a ticket in the language the person set in their own
// settings.json, short and plain, for the Verdict and Give back panes to show
// instead of the ticket's words. The ticket is never written: what Give back
// and Accept write onto it stays the original text. One call of the cheap
// model per ticket and notes: the answer is kept in the store for as long as
// the language and the notes stay the same. No language set, a failed call
// or an answer that does not fit: nothing is shown but the original, a
// rewrite shown before included - the language may have been taken away.
// The head's rewrite is kept under the ticket's id with this after it.
const HEAD = ':head'

async function rephrase($: EngineInterface, id: string, texts: readonly string[], prompt = rephrasePrompt): Promise<void> {
  const mine = ++rephrases
  newest.set(id, mine)
  const language = texts.length === 0 ? undefined : languageOf(await settingsJson($))
  const show = (map?: Record<string, string>) =>
    newest.get(id) !== mine ? undefined : update($, rephrasedAtom, ({ [id]: _, ...rest }) => (map === undefined ? rest : { ...rest, [id]: map }))
  if (language === undefined) return show()
  const shown = await rephrased($, REPHRASED_PREFIX + id, language, texts, prompt)
  await show(shown === undefined ? undefined : Object.fromEntries(texts.map((x, i) => [x, shown[i]!])))
}

// The notes in that language: kept in the store, or asked of the cheap model
// and kept; undefined where the call fails or its answer does not fit.
async function rephrased($: EngineInterface, key: string, language: string, texts: readonly string[], prompt = rephrasePrompt): Promise<string[] | undefined> {
  const kept = (await $.store.get(key)) as { language: string; source: string[]; shown: string[] } | undefined
  if (kept?.language === language && kept.source.length === texts.length && kept.source.every((x, i) => x === texts[i])) return kept.shown
  let shown: string[] | undefined
  try {
    const result = await $.model.complete({ model: 'haiku', prompt: prompt(language, texts), maxTokens: 4000 })
    shown = result.isAnswered ? rephrasedOf(result.text, texts.length) : undefined
  } catch {
    return undefined
  }
  if (shown === undefined) return undefined
  await storeNewest($, REPHRASED_PREFIX, key, { language, source: [...texts], shown })
  return shown
}

// A model's answer kept in the store under key: deleted first so the key
// moves to the end of the store's insertion order, as keep() does, and the
// tickets read longest ago under the same prefix drop out.
async function storeNewest($: EngineInterface, prefix: string, key: string, value: unknown): Promise<void> {
  await $.store.delete(key)
  await $.store.set(key, value)
  const keys = (await $.store.keys()).filter(k => k.startsWith(prefix))
  for (const old of keys.slice(0, Math.max(0, keys.length - REPHRASED_TICKETS))) await $.store.delete(old)
}

// The review's notes, each with what it starts as: its class, or how it was
// decided at an earlier acceptance, read off the ticket's notes. The classes
// of an Explain note count only while it explains this round's review.
function itemsOf(w: Waiting, showJson: string | undefined): GapItem[] {
  if (w.gaps.trim() === '') return []
  const items = gapItems(w.gaps, showJson === undefined ? undefined : explanationNotes(showJson)?.text, showJson === undefined ? undefined : decisionsOf(showJson))
  // What stops the ticket first, as the Verdict pane ordered them.
  return items.sort((a, b) => CLASS_ORDER.indexOf(a.cls) - CLASS_ORDER.indexOf(b.cls))
}

// The person's answers on the review's notes, carried out. Both read the
// ticket's lane again first: a second session or a teammate may have decided
// it since the pane opened, and then nothing is written twice.
async function stillIn($: EngineInterface, w: Waiting, what: string): Promise<{ board: string } | undefined> {
  if (!board) return undefined
  const shown = await run($, ['jaira', 'show', w.id, '--json'], await ticketDir($, board, w.id))
  const stand = shown?.exitCode === 0 ? standOf(shown.stdout) : undefined
  if (stand?.status === w.lane) return { board: stand.board }
  await fail($, `${what} ${w.id}`, stand?.status === undefined ? output(shown) : t.movedMeanwhile(w.id, w.lane, stand.status))
  await refreshBoard($)
  return undefined
}

// A follow-up ticket per note marked so, following this one, the note word
// for word as its context; the handle of each, or why it was not created.
async function createFollowUps($: EngineInterface, dir: string, id: string, onBoard: string, items: readonly GapItem[]): Promise<Map<string, string>> {
  const out = new Map<string, string>()
  const lost: string[] = []
  // A note an earlier acceptance already made a follow-up keeps that ticket.
  for (const i of items.filter(x => x.choice === 'later' && x.follow === undefined)) {
    // Every value with '=' and the title, which jaira create takes only as
    // its argument, after '--': a finding that starts with '-' ('-race
    // fails') is never read as a flag.
    const r = await run($, ['jaira', 'create', `--context=${i.text}`, `--follows=${id}`, ...(onBoard ? [`--board=${onBoard}`] : []), '--json', '--', followTitle(i.text)], dir, 30_000)
    const handle = r?.exitCode === 0 ? createdHandle(r.stdout) : undefined
    if (handle === undefined) lost.push(i.text)
    out.set(i.text, handle ?? `(not created: ${oneLine(withoutWarnings(output(r))).slice(0, 120).replace(/[()]/g, '')})`)
  }
  // Not only in the note: the person hears it, as for a note that was lost.
  if (lost.length > 0) await fail($, t.followUpsOf(id), t.followUpsLost(id, lost))
  return out
}

async function giveBackWith($: EngineInterface, w: Waiting, items: readonly GapItem[], extra: string, share: number): Promise<void> {
  const stand = await stillIn($, w, t.giveBack)
  if (stand === undefined) return
  await giveBack($, w.id, w.lane, async dir => decisionLines(items, await createFollowUps($, dir, w.id, stand.board, items), extra), share)
}

// why: what the person typed to accept a blocker not fixed, which goes onto
// the ticket as a decision of its own.
async function acceptWith($: EngineInterface, w: Waiting, items: readonly GapItem[], why: string, refused: Refused): Promise<boolean> {
  const stand = await stillIn($, w, t.accept)
  if (stand === undefined) return false
  return accept($, w, async dir => {
    const lines = decisionLines(items, await createFollowUps($, dir, w.id, stand.board, items))
    const notes = [...(lines === '' ? [] : [`Accepted from ${w.lane}: ${lines}`]), ...(why.trim() === '' ? [] : [`decision: ${t.acceptedDespite}: ${oneLine(why)}`])]
    for (const note of notes) {
      const noted = await jaira($, dir, ['note', w.id, note])
      if (!noted.ok) await fail($, t.noteOn(w.id), t.noteLost(w.id, noted.why, note))
    }
  }, refused)
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

// Give back: move, then note, then the work goes to the worker that
// carried it while that is cheaper — its tab alive, its cache warm, its
// context small — else its tab is closed FIRST and a new worker started, so
// a ticket never has two worker tabs. The choice and why go onto the ticket.
// The review lane is not touched here: the worker always hands it to a new
// agent, whose tab run-lane.sh closes after the lane.
async function giveBack($: EngineInterface, id: string, lane: string, given: string | ((dir: string) => Promise<string>), share: number): Promise<void> {
  if (!board) return
  const here = await ticketHere($, board, id)
  if (!here.ok) return fail($, `${t.giveBack} ${id}`, here.why)
  try {
    await givingBack($, id, lane, given, share, here.dir)
  } finally {
    await here.leave()
  }
}

async function givingBack($: EngineInterface, id: string, lane: string, given: string | ((dir: string) => Promise<string>), share: number, dir: string): Promise<void> {
  if (!board) return
  // Back to an earlier lane needs no --force since jaira 0.3.7 (MIN_JAIRA),
  // so any other refusal of the move still stops it. The move goes first: a
  // note written before a refused move would say the ticket went back.
  const moved = await jaira($, dir, ['move', id, '--to', 'in-progress'])
  if (movedBefore(moved, id, 'in-progress')) return movedFirst($, t.giveBack, id, 'in-progress')
  if (!moved.ok) return fail($, `${t.giveBack} ${id}`, moved.why)
  // Decisions on the review's notes are written out only now: a follow-up
  // ticket is created for a ticket that really went back.
  const reason = typeof given === 'string' ? given : await given(dir)
  // The reason is the one thing the person typed: a note that did not land is
  // shown with its text, and the handover still goes on — the ticket already
  // sits in in-progress, and nothing here could hand it out again.
  const note = `Given back from ${lane}: ${reason}`
  const noted = await jaira($, dir, ['note', id, note])
  if (!noted.ok) await fail($, t.noteOn(id), t.noteLost(id, noted.why, note))
  // Who got it, as it happened: the note and the toast say only that.
  const handedTo = async (who: 'old' | 'new' | 'none', why: string) => {
    const to = who === 'old' ? 'the old worker' : who === 'new' ? 'a new worker' : 'nobody (no worker started)'
    await jaira($, dir, ['note', id, `cockpit: given back to ${to} - ${why}`])
    await say($, t.returnedTo(id, who))
    // One line into this session, the teamlead's: the work is back and who
    // carries it now, so it starts no second worker beside this one.
    void $.prompt.submit({ text: t.toSessionGivenBack(id, lane, who, reason) })
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
  // In its own tab, like the one after Accept: not this chat's child.
  const started = await run($, [...script, '--no-parent', id.toLowerCase(), id, 'worker', board.repo], board.repo, 120_000)
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
    const settings = JSON.parse((await settingsJson($)) ?? '') as { dispatch?: unknown }
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
    // A reload drops the timers that free a pressed row; the lock is the host's and stays.
    await update($, busyAtom, () => ({}))
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
        'Returns the full text of the agent reports and task-notifications about one jaira ticket that reached this conversation shortened to one line, newest first — only those not read yet, unless all is true. Without a ticket, lists the tickets that have kept reports.',
      inputSchema: {
        type: 'object',
        properties: {
          ticket: { type: 'string', description: 'The ticket handle, e.g. MG677R.' },
          all: { type: 'boolean', description: 'Also the reports already read before. Default false: only the new ones.' },
        },
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
      void refreshBoard($).then(() => reapMerged($))
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

  // Only the doors a delivery arrives by. The engine's own security default
  // (cc-plugin-sec-default, outermost in Team and Enterprise organisations)
  // passes every person's plugin over on attachment, hook-context and
  // compaction rows, and a hook on those names the cockpit as skipped. A
  // literal, so the engine reads the doors off the source.
  on('session.append', { door: ['prompt', 'delivery'] }, ($, e, next) => next(rewriteRow(e, shortened))).catch(($, e, next) => next(e))

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
      const { ticket, all } = e as { ticket?: unknown; all?: unknown }
      return { result: await reports($, `${typeof ticket === 'string' ? ticket : ''}${all === true ? ' all' : ''}`) }
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
    const busy = await read($, busyAtom)
    const release = await read($, releaseAtom)
    const hunts = await read($, huntsAtom)
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
    if (head.length === 0 && waiting.length === 0 && theirs.length === 0 && decisions.length === 0 && !hasCache && !hasThreads && notice === null && release === null) return next(e)
    const table = $.ui.resolve(e)
    const { Box, Button, Text } = table
    const plain = e.surface === 'terminal' ? true : undefined
    const bracketed = (label: string) => (plain ? `[ ${label} ]` : label)
    // The session line: context, cache, the advice, threads; what asks for
    // action in the warning colour.
    const session: { text: string; color?: 'warning' | 'error' }[] = [
      { text: t.context(c.k) },
      ...(c.minutesLeft !== undefined ? [{ text: t.cacheLeft(c.minutesLeft) }] : cache.lastAnswerAt !== null ? [{ text: t.cacheCold, color: 'warning' as const }] : []),
      ...(c.hint === 'new' ? [{ text: t.hintNew, color: 'error' as const }] : c.hint === 'compact' ? [{ text: t.hintCompact, color: 'warning' as const }] : []),
      ...(hasThreads ? [{ text: t.threadsOpen(counts.questions, counts.stale), color: 'warning' as const }] : []),
    ]
    const shown = expanded ? waiting : waiting.slice(0, BAND_TICKETS)
    // A ticket's marks, in the theme's colours: a green tick or a red cross
    // per check, a dim '?' where nobody can tell, the review's blockers red.
    const tick = (key: string, ok: boolean | undefined) =>
      ok === undefined ? <Text key={key} dimColor>?</Text> : <Text key={key} color={ok ? 'success' : 'error'}>{ok ? '✓' : '✗'}</Text>
    const marks = (w: Waiting, hunt: Hunt | undefined) => [
      `${t.markDod} `, tick('dod', w.dodComplete),
      `${SEP}${t.markHunt} `, hunt === 'skipped' ? <Text key="hunt" dimColor>–</Text> : tick('hunt', hunt === 'ok' ? true : hunt === 'fail' ? false : undefined),
      ...(w.verdict !== ''
        ? [
          SEP, <Text key="blockers" {...(w.blockers.length > 0 ? { color: 'error' as const, bold: true } : { dimColor: true })}>{w.findingsKnown ? t.markBlockers(w.blockers.length) : t.markBlockersUnknown}</Text>,
          SEP, <Text key="notes" dimColor>{w.findingsKnown ? t.markNotes(w.notes.length) : t.markNotesUnknown}</Text>,
        ]
        : []),
      // A verdict an earlier round wrote: said, so a ticket moved on past a
      // refusing review does not reach the person silent.
      ...(w.reviewedIn !== undefined ? [SEP, <Text key="old" color="warning">{t.markOldReview(w.reviewedIn)}</Text>] : []),
      ...(w.question !== '' ? [SEP, <Text key="question" color="warning">?</Text>] : []),
    ]
    // A frame of its own, so the band reads as the cockpit's and not as chat.
    return (
      <Box flexDirection="column" borderStyle="round" borderDimColor paddingX={1}>
        <Box flexDirection="row">
          <Text bold color="claude">{t.bandTitle}</Text>
          {head.length > 0 ? <Text dimColor wrap="truncate-end">{'   '}{head.join(SEP)}</Text> : null}
          {board ? <Button key="overview" label={overviewWords.button} onPress={() => void openOverview($, board!.dir)} /> : null}
        </Box>
        {/* What just happened on top, this person's own work below it, the
            teammates' tickets last of all lists. */}
        {notice !== null ? <Text dimColor wrap="truncate-end">{notice}</Text> : null}
        {release !== null && board
          ? (
            <Box flexDirection="row">
              <Text color="warning" wrap="truncate-end">{t.released(release.latest, release.installed)}{' '}</Text>
              <Button key="release-update" label={t.update} plain={plain} onPress={() => void once($, RELEASE, () => previewUpgrade($, board!.dir))}><Text color="success" dimColor={busy[RELEASE] !== undefined}>{bracketed(t.update)}</Text></Button>
              <Button key="release-dismiss" label={t.dismiss} onPress={() => void dismissRelease($, release.latest)} />
            </Box>
          )
          : null}
        {waiting.length > 0 ? <Box marginTop={1}><Text bold>{t.waiting}</Text></Box> : null}
        {/* One row per ticket: its marks say whether it can be accepted, the
            verdict, the notes and the question are a press on Verdict away. */}
        {shown.map(w => (
          <Box key={`row-${w.id}`} flexDirection="column">
            <Box flexDirection="row">
              <Text wrap="truncate-end">
                <Text bold>{w.id}</Text>{SEP}{w.lane}{SEP}{marks(w, hunts[w.id]?.state)}{' '}
              </Text>
              {/* Told apart by colour before they are read: the theme's own, so a
                  light terminal reads them too; the focus inverts them as any.
                  The engine draws a Button's brackets uncoloured, so the
                  terminal gets them plain and in the label's colour. */}
              {/* One button for Accept, Give back and the verdict: the
                  Acceptance pane shows the ticket and does what is chosen in it. */}
              <Button key={`acceptance-${w.id}`} label={t.acceptance} plain={plain} onPress={() => void once($, w.id, () => openAcceptance($, w))}><Text color="claude" dimColor={busy[w.id] !== undefined}>{bracketed(t.acceptance)}</Text></Button>
              {explained[w.id] !== undefined
                ? <Button key={`read-${w.id}`} label={t.readExplanation} onPress={() => void readExplanation($, w.id, explained[w.id]!)} />
                : null}
              {w.question !== ''
                ? <Button key={`discuss-${w.id}`} label={t.discuss} onPress={() => void once($, w.id, () => openChat($, w.id, 'discuss'))} />
                : null}
            </Box>
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
                  ? <Button key={`discuss-${w.id}`} label={t.discuss} onPress={() => void once($, w.id, () => openChat($, w.id, 'discuss'))} />
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

  // The Acceptance pane, as laid out with Alex on 22.11.2026: the light and
  // the ticket on top, how it got here, each review note with ok / fix /
  // ticket and Explain, then what Execute will do and the button.
  on('ui.render', { component: 'Pane', requestId: ACCEPTANCE }, async ($, e) => {
    const id = await read($, givingAtom)
    const w = (await read($, boardAtom))?.waiting.find(x => x.id === id)
    const table = $.ui.resolve(e)
    const { Box, Button, Text } = table
    const Input = 'Input' in table ? table.Input : undefined
    if (w === undefined) return nothingPane(table, 'acceptance-head', { key: 'acceptance-close', onPress: () => closeAcceptance($) })
    const choosing = await read($, choosingAtom)
    const items = choosing?.id === w.id ? choosing.items : []
    const story = await read($, storyAtom)
    const s = story?.id === w.id ? story : undefined
    const typed = await read($, typedAtom)
    const refusal = await read($, refusalAtom)
    const reasons = (await read($, reasonsAtom))[w.id] ?? {}
    const rephrased = await read($, rephrasedAtom)
    const said = { ...rephrased[`${w.id}${HEAD}`], ...rephrased[w.id] }
    const tr = (x: string) => said[x] ?? x
    const hunt = (await read($, huntsAtom))[w.id]
    const light = lightOf(w, hunt, items.map(x => x.cls))
    const first = light.reasons[0]
    const exec = executionOf(items, typed, w.acceptTo)
    const now = await $.clock.now()
    const choose = (n: number, choice: GapItem['choice']) =>
      void update($, choosingAtom, c => (c === null || c.id !== w.id ? c : { ...c, items: c.items.map((x, i) => (i === n ? { ...x, choice } : x)) }))
    const go = () => void once($, w.id, async () => execute($, w, await read($, typedAtom), share, Input !== undefined))
    // A choice's button: the chosen one bold in its colour, the others dim.
    const CHOICES: [GapItem['choice'], string, 'success' | 'error' | 'suggestion'][] = [['accept', t.choiceOk, 'success'], ['now', t.choiceFix, 'error'], ['later', t.choiceTicket, 'suggestion']]
    const story_ = s?.story
    const history = [
      ...(story_ !== undefined ? [[t.round(story_.round), ...(story_.back !== undefined ? [t.backFrom(story_.back.count, story_.back.from)] : []), ...(story_.created !== undefined ? [t.age(now - story_.created)] : [])].join(SEP)] : []),
      [...(story_ !== undefined && story_.dod.total > 0 ? [t.dodCount(story_.dod.done, story_.dod.total, story_.dod.superseded)] : []), t.huntOnHead(hunt?.state ?? ''), t.review(w.verdict)].join(SEP),
      ...((s?.commits?.commits ?? 0) > 0 || s?.pr ? [[...((s?.commits?.commits ?? 0) > 0 ? [t.commitsFiles(s!.commits!.commits, s!.commits!.files)] : []), ...(s?.pr ? [t.prOpen] : [])].join(SEP)] : []),
    ]
    return (
      <Box flexDirection="column">
        {headRow({ Box, Button }, 'acceptance-head', { key: 'acceptance-close', label: t.close, onPress: () => closeAcceptance($) })}
        <Text key="light" wrap="wrap">
          <Text bold color={LIGHT_COLOUR[light.color]}>● {t.light[light.color]}</Text>
          : {first === undefined ? t.lightAccept : tr(t.lightReason(first))}{light.reasons.length > 1 ? t.lightMore(light.reasons.length - 1) : ''}
        </Text>
        <Text key="title" bold wrap="wrap">{tr(story_?.title || w.title)}</Text>
        {(story_?.goal || w.goal) ? <Text key="goal" wrap="wrap"><Text bold>{t.goal}</Text> {tr(story_?.goal || w.goal)}</Text> : null}
        {story_?.what ? <Text key="what" wrap="wrap"><Text bold>{t.done}</Text> {tr(story_.what)}</Text> : null}
        {story_?.why ? <Text key="why" wrap="wrap"><Text bold>{t.why}</Text> {tr(story_.why)}</Text> : null}
        {w.question !== '' ? <Text key="question" color="warning" wrap="wrap">? {w.question}</Text> : null}
        <Box marginTop={1}><Text bold>{t.historyHeading}</Text><Text dimColor wrap="wrap">{`${SEP}${w.id}${SEP}${w.lane}${w.reviewedIn !== undefined ? `${SEP}${t.reviewedIn(w.reviewedIn)}` : ''}`}</Text></Box>
        {history.map((line, i) => <Text key={`history-${i}`} wrap="wrap">{line}</Text>)}
        {(s?.release.length ?? 0) > 0 ? <Text key="release" wrap="wrap"><Text bold>{t.releaseNote}</Text> {s!.release.join(' / ')}</Text> : null}
        {items.length > 0 ? <Box marginTop={1}><Text bold>{t.notesChoose}</Text></Box> : null}
        {!w.findingsKnown ? <Text dimColor wrap="wrap">{t.notesNotKept}</Text> : null}
        {items.map((x, i) => (
          <Box key={`note-${w.id}-${i}`} flexDirection="column" marginTop={i > 0 ? 1 : 0}>
            <Text key={`text-${w.id}-${i}`} wrap="wrap">
              {`${NOTE_MARK[x.cls]} `}{tr(x.text)}
              {x.earlier !== undefined ? <Text dimColor>{` ${t.decidedBefore}`}</Text> : null}
            </Text>
            <Box flexDirection="row">
              {CHOICES.map(([c, label, colour], j) => (
                <Button key={`${c}-${w.id}-${i}`} label={label} plain {...(i === 0 && j === 0 ? { autoFocus: true as const } : {})} onPress={() => choose(i, c)}>
                  <Text {...(x.choice === c ? { bold: true, color: colour } : { dimColor: true })}>{x.choice === c ? `[•${label}]` : `[ ${label} ]`}</Text>
                </Button>
              ))}
              <Button key={`why-${w.id}-${i}`} label={t.explainNote} onPress={() => void explainNote($, w.id, story_?.title || w.title, x.text)} />
            </Box>
            {reasons[x.text] !== undefined ? <Text dimColor wrap="wrap">{`  ${reasons[x.text] === '' ? t.explaining : reasons[x.text]}`}</Text> : null}
          </Box>
        ))}
        <Box marginTop={1}><Text bold wrap="wrap">{t.summary(exec.step, exec.fix, exec.followUps, w.acceptTo ?? '')}</Text></Box>
        {refusal?.id === w.id ? <Text key="refusal" color="error" wrap="wrap">{refusal.text}</Text> : null}
        {Input !== undefined
          ? (
            <Input
              key={`typed-${w.id}`}
              label={exec.asksWhy ? t.whyLabel : items.length > 0 ? t.moreLabel : t.reasonLabel}
              placeholder={exec.asksWhy ? t.whyPlaceholder : items.length > 0 ? t.morePlaceholder : t.reasonPlaceholder}
              submitLabel={t.execute}
              {...(items.length === 0 ? { autoFocus: true as const } : {})}
              onInput={(text: string) => void update($, typedAtom, () => text)}
              onSubmit={(text: string) => {
                void update($, typedAtom, () => text)
                // Enter in the field never accepts by itself - one Enter too
                // many would take the ticket to done: that is Execute's.
                // Words typed, or a note to fix, it carries out.
                const plan = executionOf(items, text, w.acceptTo)
                if (plan.step === 'accept' && !plan.asksWhy) return void say($, t.nothingToFix(w.id))
                void once($, w.id, () => execute($, w, text, share))
              }}
            />
          )
          : null}
        <Box flexDirection="row">
          <Button key={`execute-${w.id}`} label={t.execute} onPress={go} />
          <Button key={`explain-all-${w.id}`} label={t.explainAll} onPress={() => void once($, `${w.id}${CHAT}`, () => openChat($, w.id, 'explain'))} />
        </Box>
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: EXPLANATION }, async ($, e) => {
    const explaining = await read($, explainingAtom)
    const { Box, Button, Text } = $.ui.resolve(e)
    if (explaining === null) return nothingPane({ Box, Button, Text }, 'explanation-head', { key: 'explanation-close', onPress: () => void $.ui.close({ id: EXPLANATION }) })
    const { id, decided } = explaining
    const x = explanationOf(explaining.text)
    const ticket = await waitingTicket($, id)
    const open = decided === undefined
    const pick = (o: Option) => {
      void $.ui.close({ id: EXPLANATION })
      void decide($, id, `${o.key}) ${o.text}`, o.action, share)
    }
    // An option whole: what it says, what picking it does in the band's
    // colours for that step, whether the review recommends it, highlighted;
    // the rest of the row neutral. Drawn as many rows as it needs.
    const mark = (a: Action) => (a === 'accept' ? t.actionAccept : t.actionGiveBack)
    // A review note's class in the theme's colours, so a light terminal reads
    // them too: a blocker red, a small edge yellow, a known limit dim, one
    // that belongs to another ticket blue.
    const noteMark = (cls: NoteClass, label: string) =>
      cls === 'limit'
        ? <Text dimColor>{`[${label}]`}</Text>
        : <Text bold color={cls === 'blocker' ? 'error' : cls === 'edge' ? 'warning' : 'suggestion'}>{`[${label}]`}</Text>
    const said = (o: Option) =>
      `${o.text}${o.action !== undefined ? ` ${mark(o.action)}` : ''}${o.key === x.recommended ? ` ${t.recommended}` : ''}`
    const line = (o: Option) => (
      <Text wrap="wrap">
        {o.text}
        {o.action !== undefined ? <Text color={o.action === 'accept' ? 'success' : 'error'}>{` ${mark(o.action)}`}</Text> : null}
        {o.key === x.recommended ? <Text color="suggestion">{` ${t.recommended}`}</Text> : null}
      </Text>
    )
    // Each row's digit is its place, whatever the note numbered it: the
    // engine draws it as '1:' before the row — before the label too, where a
    // surface draws no children, so the label leaves it out — and takes it
    // as the key that picks it, which a letter could not be in capitals. Own
    // answer is the digit after the last option.
    const digit = (key: string) => String(x.options.findIndex(o => o.key === key) + 1)
    const ownDigit = String(x.options.length + 1)
    return (
      <Box flexDirection="column">
        {headRow({ Box, Button }, 'explanation-head', { key: 'explanation-close', label: t.close, onPress: () => void $.ui.close({ id: EXPLANATION }) })}
        {about(Text, ticket)}
        <Box marginTop={ticket !== undefined ? 1 : 0}><Text bold>{t.explanationOf(id)}</Text></Box>
        {x.parts.map((p, i) => p.kind === 'option'
          // Open, each option is a row of its own to pick: the engine's
          // Select draws one line per option and cuts a long one, so the
          // choice would hide what is chosen. Its digit picks it too.
          ? open
            ? <Button key={`pick-${id}-${p.key}`} plain hotkey={digit(p.key)} label={said(p)} {...(p.key === (x.recommended ?? x.options[0]?.key) ? { autoFocus: true as const } : {})} onPress={() => pick(p)}>{line(p)}</Button>
            : <Box key={`part-${i}`} flexDirection="row"><Text bold>{`  ${p.key}) `}</Text>{line(p)}</Box>
          // The status line: a green tick or a red cross per check.
          : p.kind === 'status'
            ? <Box key={`part-${i}`} flexDirection="row" marginBottom={1}><Text wrap="wrap">{p.checks.map((c, j) => (
                <Text key={`check-${j}`}>
                  {j > 0 ? '   ' : ''}
                  {c.unknown ? <Text dimColor>{'? '}</Text> : c.ok === undefined ? null : <Text color={c.ok ? 'success' : 'error'}>{c.ok ? '✓ ' : '✗ '}</Text>}
                  {c.label}
                </Text>
              ))}</Text></Box>
          // A review note: only its class mark coloured, the text plain.
          : p.kind === 'note'
            ? <Box key={`part-${i}`} flexDirection="row"><Text wrap="wrap">{noteMark(p.cls, p.label)} {p.text}</Text></Box>
          // The recommendation, bold in a colour of its own.
          : p.kind === 'recommendation'
            ? <Box key={`part-${i}`} marginTop={1}><Text wrap="wrap" bold color="claude">{p.text}</Text></Box>
          // A paragraph a line apart; the options one under the other.
          : <Box key={`part-${i}`} marginTop={i > 0 && x.parts[i - 1]?.kind !== 'status' ? 1 : 0} marginBottom={x.parts[i + 1]?.kind === 'option' ? 1 : 0}><Text wrap="wrap">{p.text}</Text></Box>)}
        {/* Answered already, on the ticket: the choice is read, not offered again. */}
        {open
          ? <Box marginTop={1}><Button key={`pick-${id}-own`} plain hotkey={ownDigit} {...(x.options.length === 0 ? { autoFocus: true as const } : {})} label={t.ownAnswer} onPress={() => void askOwnAnswer($, id)} /></Box>
          : <Box marginTop={1}><Text wrap="wrap"><Text bold>{t.decidedAs}</Text> {decided}</Text></Box>}
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

  on('ui.render', { component: 'Pane', requestId: OVERVIEW }, async ($, e) =>
    drawOverview($.ui.resolve(e), await read($, overviewAtom), { tab: await read($, overviewTabAtom), folded: await read($, overviewFoldedAtom), people: await read($, overviewPeopleAtom), lanes: await read($, overviewLanesAtom) }, {
      close: () => {
        overviewOpen = false
        stopFrames()
        void $.ui.close({ id: OVERVIEW })
      },
      tab: tab => void showTab($, tab),
      fold: key => void update($, overviewFoldedAtom, f => ({ ...f, [key]: !f[key] })),
      person: name => void update($, overviewPeopleAtom, p => toggled(p, name)),
      everybody: () => void update($, overviewPeopleAtom, () => []),
      lane: id => void update($, overviewLanesAtom, l => toggled(l, id)),
      allLanes: () => void update($, overviewLanesAtom, () => []),
    }, {
      workers: await read($, workersAtom),
      colors: await read($, workersColorsAtom),
      hats: await read($, workersHatsAtom),
      since: await read($, workersSinceAtom),
      reports: await read($, workersReportsAtom),
      now: await $.clock.now(),
      falling: await fallingOf($),
      frame: await read($, workersFrameAtom),
      columns: e.props.bodyColumns,
      rows: e.props.scroll.bodyRows - OVERVIEW_ROWS,
    }))

  // The Workers tab's frames stop with the pane when the person closes it;
  // its own Close button stops them itself, as this hook does not hear it.
  on('ui.close', async ($, e, next) => {
    if (e.id === OVERVIEW) {
      overviewOpen = false
      stopFrames()
    }
    return next(e)
  }).catch(($, e, next) => next(e))

  // What the ticket is for, before anything about it: a pane opened from the
  // band is read by someone who may not remember the ticket.
  const about = (Text: ElementTable['Text'], w: Waiting | undefined) =>
    w === undefined
      ? []
      : [
        <Text key="about-title" bold wrap="wrap">{w.title}</Text>,
        ...(w.goal !== '' ? [<Text key="about-goal" wrap="wrap"><Text bold>{t.goal}</Text> {w.goal}</Text>] : []),
      ]

  on('ui.render', { component: 'Pane', requestId: MESSAGE }, async ($, e) => {
    const message = await read($, messageAtom)
    const { Box, Button, Text } = $.ui.resolve(e)
    return (
      <Box flexDirection="column">
        {headRow({ Box, Button }, 'message-head', { key: 'message-close', label: t.close, onPress: () => void $.ui.close({ id: MESSAGE }) })}
        <Text bold>{message?.title ?? ''}</Text>
        <Text wrap="wrap">{message?.text ?? ''}</Text>
      </Box>
    )
  })

  // Update's preview, and the one press that writes: Cancel holds the focus,
  // so an Enter pressed once too often writes nothing.
  on('ui.render', { component: 'Pane', requestId: UPGRADE }, async ($, e) => {
    const u = await read($, upgradeAtom)
    const { Box, Button, Text } = $.ui.resolve(e)
    if (u === null || !board) return <Box><Text dimColor>{t.paneNothing}</Text></Box>
    const dir = board.dir
    return (
      <Box flexDirection="column">
        <Text wrap="wrap">{t.upgradeIntro(u.from, u.to)}</Text>
        <Text wrap="wrap">{u.preview}</Text>
        <Box flexDirection="row" marginTop={1}>
          <Button key="upgrade-confirm" label={t.confirmUpgrade(u.to)} onPress={() => void once($, RELEASE, () => upgradeJaira($, dir, u.to))} />
          <Button key="upgrade-cancel" label={t.cancel} autoFocus onPress={() => void $.ui.close({ id: UPGRADE })} />
        </Box>
      </Box>
    )
  })

  on('ui.render', { component: 'Pane', requestId: PANE }, async ($, e) => {
    const id = await read($, paneAtom)
    const { Box, Button, Text } = $.ui.resolve(e)
    if (id === null) return nothingPane({ Box, Button, Text }, 'ticket-head', { key: 'ticket-close', onPress: () => void $.ui.close({ id: PANE }) })
    // Opening a ticket reads its reports: they count as read, as after report.
    const kept = await locked(async () => readReports($, keyOf(await project($), id)))
    let csv = ''
    try {
      csv = await $.fs.read(`${await home($)}/.jaira/metrics.csv`)
    } catch {}
    const usage = rowsOf(csv, id)
    const waiting = (await read($, boardAtom))?.waiting.find(w => w.id === id)
    return (
      <Box flexDirection="column">
        {headRow({ Box, Button }, 'ticket-head', { key: 'ticket-close', label: t.close, onPress: () => void $.ui.close({ id: PANE }) })}
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
