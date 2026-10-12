import type { On, SessionRateLimit } from 'claude-code'
import { describe, expect, mock, test } from 'claude-code/testing'

// A world beneath the plugin: a repository at /repo with a board, jaira, git
// and Herdr answering from memory, and every command the plugin runs
// recorded. `board: false` is a project with no .jaira/.

type World = {
  ran: string[][]
  files: Map<string, string>
  usage: { tokens: number; rateLimits: SessionRateLimit[] }
  tickets: unknown[]
  tabs: { tab_id: string; label: string; agent_status?: string }[]
  // The old worker's pane went away after the tab list was read.
  paneGone: boolean
  // This session's subagents, as $.agent.list() gives them.
  agents: { id: string; description: string; type: string; status: string; name?: string }[]
  // `jaira lanes --json` of a build before human_exit.
  oldJaira: boolean
  asked: string[]
  answer: string
  model: string
  modelCalls: number
  filled: string[]
  store: Map<string, unknown>
  status: (string | undefined)[]
  // What a dry run of a move refuses, by code: a person's lane always refuses
  // needs_human, the exit Accept stands for.
  refusals: { code: string; message: string }[]
  // A jaira older than the board: it warns on a line before what it prints.
  warns: boolean
  // What `jaira whoami --json` says; undefined: it says nothing.
  me?: { identity: string; aliases: string[] }
  // The ticket bodies `jaira show --json` gives, by id.
  bodies: Record<string, string>
  // The panes the plugin opened, and closed, by id and title.
  opened: { id: string; title?: string; focus?: boolean }[]
  closed: string[]
  clock: ReturnType<typeof mock.clock>
  // After Accept: commits the ticket's branch has that the target lacks, the
  // forge's answer for its pull requests (null: the tool failed), origin
  // (null: none), `git config jaira.forge`, and `jaira next --all`.
  ahead: number
  prs: unknown[] | null
  origin: string | null
  forge: string
  next: string[]
  // `git branch --list '*/<slug>'`, read where the ticket has no worktree.
  branches: string[]
  // `git log --format=%s master..<branch>`; the commits on branches the
  // base lacks, which `git log -E --grep` searches with the pattern it is
  // given, and the branches `git branch --contains` gives the one it finds.
  subjects: string[]
  branchCommits: { sha: string; subject: string }[]
  holders: string[]
  toasts: string[]
  // `git remote -v`.
  remotes: string
  // The forge answers `repo view` but refuses `pr list`.
  prListFails: boolean
}

const LANES = {
  lanes: [
    { id: 'in-progress', agentic: true, context: '' },
    { id: 'critique', agentic: true, context: '' },
    { id: 'review', agentic: true, context: 'fresh' },
    { id: 'signoff', agentic: false, context: '', human_exit: true },
    { id: 'human', agentic: false, context: '', human_exit: true },
    { id: 'done', agentic: false, context: '', terminal: true },
  ],
}

const SIGNOFF = { handle: 'MG677R', title: 'cache key', status: 'signoff', next_lane: 'done', question: '', review: { verdict: 'approve mit Vorbehalt: der Schluessel ignoriert den Branch' } }
const HUMAN = { handle: 'K7P2QX', title: 'tags popup', status: 'human', next_lane: 'review', question: 'Popup oder eigene Seite?', review: { verdict: '' } }
const WORKING = { handle: 'T8QS1Y', title: 'filter', status: 'in-progress', question: '', review: {} }

const OLDER = "jaira: this board was set up by an older version of jaira — run 'jaira update' to bring it up to date"
const GIVE = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-give-back', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as const
const MESSAGE = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-message', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never

const WORKTREES = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n\nworktree /.worktrees/repo-mg677r\nHEAD 2\nbranch refs/heads/feat/mg677r\n'
const SPAWN = '/repo/.claude/skills/jaira-dispatcher/scripts/spawn.sh'

function world(on: On, opts: { board?: boolean; herdr?: boolean; store?: Record<string, unknown>; spawn?: boolean } = {}): World {
  const w: World = {
    ran: [],
    files: new Map(),
    usage: { tokens: 0, rateLimits: [{ kind: 'five_hour', percentUsed: 34 }] },
    tickets: [SIGNOFF, HUMAN, WORKING],
    tabs: [],
    paneGone: false,
    agents: [],
    oldJaira: false,
    asked: [],
    answer: '',
    model: '',
    modelCalls: 0,
    filled: [],
    store: new Map(Object.entries(opts.store ?? {})),
    status: [],
    refusals: [{ code: 'needs_human', message: 'lane "signoff" is a human checkpoint' }],
    warns: false,
    bodies: {},
    opened: [],
    closed: [],
    clock: mock.clock(on, { now: 1_000_000 }),
    ahead: 0,
    prs: [],
    origin: 'git@github.com:acme/repo.git',
    forge: '',
    next: [],
    branches: [],
    subjects: ['feat(MG677R): cache key'],
    branchCommits: [],
    holders: [],
    toasts: [],
    prListFails: false,
    remotes: 'origin\tgit@github.com:acme/repo.git (fetch)\norigin\tgit@github.com:acme/repo.git (push)\nupstream\thttps://github.com/upstream/repo.git (fetch)\n',
  }
  const hasBoard = opts.board ?? true
  // The store, readable by the test; insertion order kept as the engine's is.
  on('store.get', ($, e) => ({ value: w.store.get(e.key) }))
  on('store.set', ($, e) => {
    w.store.set(e.key, e.value)
    return { value: undefined } as never
  })
  on('store.delete', ($, e) => {
    w.store.delete(e.key)
    return { value: undefined } as never
  })
  on('store.keys', () => ({ value: [...w.store.keys()] }))
  mock.env(on, { HOME: '/home/a', ...(opts.herdr ? { HERDR_ENV: '1', HERDR_BIN_PATH: 'herdr', HERDR_PANE_ID: 'w1:p1', CLAUDE_CODE_MESSAGING_SOCKET: '/tmp/cc-socks-1000/77.sock' } : {}) })
  on('session.start', ($, e) => ({ cwd: e.cwd }))
  on('command.register', () => ({ value: undefined }) as never)
  on('tool.register', () => ({ value: undefined }) as never)
  on('ui.render', () => ({ type: 'Box', props: {}, children: [] }) as never)
  on('ui.open', ($, e) => {
    const { id, title, focus } = e as unknown as { id: string; title?: string; focus?: boolean }
    w.opened.push({ id, ...(title !== undefined ? { title } : {}), ...(focus ? { focus } : {}) })
    return { value: { isPlaced: true } }
  })
  on('ui.close', ($, e) => {
    w.closed.push((e as unknown as { id: string }).id)
    return { value: undefined } as never
  })
  on('ui.status', ($, e) => {
    w.status.push((e as { text?: string }).text)
    return { value: undefined } as never
  })
  on('prompt.fill', ($, e) => {
    w.filled.push(e.text)
    return { isFilled: true } as never
  })
  on('ui.toast', ($, e) => {
    w.toasts.push((e as unknown as { text: string }).text)
    return { value: undefined }
  })
  on('session.end', ($, e) => ({ sessionId: 'abcdef123456' }))
  on('session.receive', ($, e) => ({ text: e.text }))
  on('prompt.submit', ($, e) => ({ text: e.text, origin: e.origin }))
  on('turn.complete', ($, e) => ({ text: e.answer, ...(e.usage ? { usage: e.usage } : {}) }))
  on('session.cwd', () => ({ value: '/repo' }))
  on('session.repo', () => ({ value: { root: '/repo', remote: null, internal: false, name: null } }))
  on('session.id', () => ({ value: 'abcdef123456' }))
  on('session.usage', () => ({ value: { startedAt: 0, context: { tokens: w.usage.tokens, window: 1_000_000 }, rateLimits: w.usage.rateLimits } }))
  on('fs.exists', ($, e) => ({ value: (hasBoard && e.path === '/repo/.jaira') || (e.path === SPAWN && opts.spawn !== false) || w.files.has(e.path) }))
  on('fs.read', ($, e) => (w.files.has(e.path) ? { value: w.files.get(e.path)! } : { deny: 'no such file' }))
  on('fs.write', ($, e) => {
    w.files.set(e.path, e.text)
    return { value: undefined }
  })
  on('process.run', ($, e) => {
    w.ran.push([...e.argv, `@${e.init?.cwd ?? ''}`])
    const ok = (stdout = '') => ({ value: { exitCode: 0, stdout, stderr: '', isStdoutTruncated: false, isStderrTruncated: false } })
    const [cmd, ...args] = e.argv
    if (cmd === 'jaira' && args[0] === 'lanes') return ok(JSON.stringify(w.oldJaira ? { lanes: LANES.lanes.map(({ human_exit, ...l }) => l) } : LANES))
    if (cmd === 'jaira' && args[0] === 'list') return ok(JSON.stringify({ count: w.tickets.length, tickets: w.tickets }))
    if (cmd === 'jaira' && args[0] === 'move' && args.includes('--dry-run')) {
      const error = (w.warns ? `${OLDER}\n` : '') + JSON.stringify({ error: { code: 3, reason: 'gate_refused', violations: w.refusals } })
      return { value: { exitCode: 3, stdout: '', stderr: error, isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'whoami') return ok(w.me === undefined ? '' : JSON.stringify(w.me))
    if (cmd === 'jaira' && args[0] === 'show') return ok((w.warns ? `${OLDER}\n` : '') + JSON.stringify({ handle: args[1], body: w.bodies[args[1]!] ?? '' }))
    if (cmd === 'git' && args[0] === 'rev-parse') return ok('/repo/.git\n')
    if (cmd === 'git' && args[0] === 'worktree') return ok(WORKTREES)
    const failed = (stderr: string) => ({ value: { exitCode: 1, stdout: '', stderr, isStdoutTruncated: false, isStderrTruncated: false } })
    if (cmd === 'git' && args[0] === 'config') return w.forge ? ok(`${w.forge}\n`) : failed('')
    if (cmd === 'git' && args[0] === 'remote' && args[1] === 'get-url') return w.origin === null ? failed("error: No such remote 'origin'") : ok(`${w.origin}\n`)
    if ((cmd === 'gh' || cmd === 'glab') && w.prs === null) return failed(`${cmd}: command not found`)
    if ((cmd === 'gh' || cmd === 'glab') && args[1] === 'list' && w.prListFails) return failed("HTTP 401: Bad credentials")
    // origin is a fork of upstream/repo: pull requests are asked of the parent.
    // Its default branch is main; the fork's own lags behind as master.
    if (cmd === 'gh' && args[0] === 'repo') {
      return ok(JSON.stringify(args[2] === 'upstream/repo'
        ? { isFork: false, nameWithOwner: 'upstream/repo', parent: null, defaultBranchRef: { name: 'main' } }
        : { isFork: true, nameWithOwner: 'acme/repo', parent: { name: 'repo', owner: { login: 'upstream' } }, defaultBranchRef: { name: 'master' } }))
    }
    if (cmd === 'glab' && args[0] === 'repo') {
      return ok(JSON.stringify(args[2] === 'upstream/repo'
        ? { path_with_namespace: 'upstream/repo', default_branch: 'main' }
        : { path_with_namespace: 'acme/repo', forked_from_project: { path_with_namespace: 'upstream/repo' }, default_branch: 'master' }))
    }
    if (cmd === 'git' && args[0] === 'remote' && args[1] === '-v') return ok(w.remotes)
    if (cmd === 'git' && args[0] === 'branch' && args[1] === '--contains') return ok(w.holders.join('\n'))
    if (cmd === 'git' && args[0] === 'branch') return ok(w.branches.join('\n'))
    if (cmd === 'git' && args[0] === 'log' && args.includes('--branches')) {
      const pattern = new RegExp(args.find(a => a.startsWith('--grep='))!.slice('--grep='.length), 'm')
      return ok(w.branchCommits.find(c => pattern.test(c.subject))?.sha ?? '')
    }
    // As many subjects as commits ahead: the first `ahead` of them.
    if (cmd === 'git' && args[0] === 'log') return ok(w.subjects.concat(Array(Math.max(0, w.ahead - w.subjects.length)).fill('chore: more')).slice(0, w.ahead).join('\n'))
    // The move a press forces lands: jaira list reads the ticket in done after it.
    if (cmd === 'jaira' && args[0] === 'move' && args.includes('--force')) {
      w.tickets = w.tickets.map(x => ((x as { handle: string }).handle === args[1] ? { ...(x as object), status: args[3] } : x))
      return ok()
    }
    if (cmd === 'gh' || cmd === 'glab') return ok(JSON.stringify(w.prs))
    if (cmd === 'jaira' && args[0] === 'next') return ok(JSON.stringify({ count: w.next.length, tickets: w.next.map(handle => ({ handle })) }))
    if (cmd === 'herdr' && args[0] === 'tab' && args[1] === 'list') return ok(JSON.stringify({ result: { tabs: w.tabs } }))
    if (cmd === 'herdr' && args[0] === 'pane' && w.paneGone) return { value: { exitCode: 1, stdout: '', stderr: 'no such pane', isStdoutTruncated: false, isStderrTruncated: false } }
    return ok()
  })
  on('tool.call', { tool: 'AskUserQuestion' }, ($, e) => {
    w.asked.push(JSON.stringify(e))
    const { questions } = e as { questions: { question: string }[] }
    return { result: { questions, answers: { [questions[0]!.question]: w.answer } }, text: w.answer }
  })
  on('agent.list', () => ({ value: w.agents }) as never)
  on('model.complete', () => {
    w.modelCalls++
    return { value: { isAnswered: true, text: w.model, usage: { input_tokens: 1, output_tokens: 1, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 } } }
  })
  return w
}

const BAND = { plugin: 'jaira-cockpit', surface: 'terminal', component: 'AbovePrompt', props: { hasSurvey: false, isWorking: false, maxRows: 20, bodyColumns: 120, scroll: { offset: 0, bodyRows: 20 }, view: {} } } as const

// The engine's `$` a test body is handed.
type Test$ = Parameters<Extract<Parameters<typeof test>[1], (...args: never[]) => unknown>>[0]

const typed = (args = '') =>
  ({ command: 'cockpit', args, origin: { kind: 'composer' }, presentation: { isFullscreen: false, columns: 120 } }) as const

async function startSession($: Test$) {
  await $.session.start({ source: 'startup', cwd: '/repo', surface: null, isInteractive: true } as never)
}

async function answer($: Test$, text = 'ok', tokens = 0) {
  await $.turn.complete({
    answer: text, durationMs: 10, isAborted: false, turnId: `t${Math.random()}`, reason: 'answer',
    usage: { input_tokens: 100, output_tokens: 20, cache_creation_input_tokens: tokens, cache_read_input_tokens: 500, model: 'claude-opus-5-5' },
  } as never)
}

// The band's head line beside its title: the board's counts and what was used.
async function headLine($: Test$) {
  const ui = await $.ui.mount(BAND)
  return (await ui.find({ type: 'Text', text: /^ {3}\S/ }))?.text.trim()
}

// The band's session line, its parts as drawn with the colour of each: the
// band's last row, from the context on.
async function sessionLine(ui: { findAll: (q: { type: string }) => Promise<{ text: string; props: Record<string, unknown> }[]> }) {
  const texts = await ui.findAll({ type: 'Text' })
  const parts = texts.slice(texts.findLastIndex(x => x.text.startsWith('context ')))
  return { text: parts.map(x => x.text).join(''), colors: parts.filter(x => x.text !== ' | ').map(x => (x.props.color as string | undefined) ?? '') }
}

const ran = (w: World, ...prefix: string[]) => w.ran.filter(r => prefix.every((p, i) => r[i] === p))

describe('waiting for you', () => {
  test('the band lists the tickets in a person\'s lanes, with their verdict', async ($, on) => {
    world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    // A frame of its own, titled, so it reads as the cockpit and not as chat.
    expect((await ui.find({ type: 'Box' }))?.props.borderStyle).toBe('round')
    expect(await ui.find({ type: 'Text', text: 'jaira-cockpit' })).toBeDefined()
    expect((await ui.find({ text: 'Waiting for you' }))).toBeDefined()
    expect((await ui.find({ type: 'Text', text: /MG677R \| signoff \| approve mit Vorbehalt/ }))).toBeDefined()
    expect((await ui.find({ type: 'Text', text: /K7P2QX \| human/ }))).toBeDefined()
    expect(await ui.find({ text: /T8QS1Y/ })).toBeUndefined()
  })

  test('Accept moves the ticket to done in its own worktree', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R')).toEqual([
      ['jaira', 'move', 'MG677R', '--to', 'done', '--dry-run', '--json', '@/.worktrees/repo-mg677r'],
      ['jaira', 'move', 'MG677R', '--to', 'done', '--force', '@/.worktrees/repo-mg677r'],
    ])
  })

  test('Accept is offered only where the next lane ends the board: never in human, past review and signoff', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'accept-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'accept-K7P2QX' })).toBeUndefined()
    expect(await ui.find({ key: 'discuss-K7P2QX' })).toBeDefined()
    void w
  })

  test('Accept forces only the human exit: any other gate refusing moves nothing and says why', async ($, on) => {
    const w = world(on)
    w.refusals = [...w.refusals, { code: 'needs_nonmodel_signal', message: 'the definition of done is not met: criterion 1 is still open' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-MG677R' })
    expect(ran(w, 'jaira', 'move').filter(r => r.includes('--force'))).toEqual([])
    expect(w.opened).toEqual([{ id: 'cockpit-message', title: 'Accept MG677R failed', focus: true }])
    const pane = await $.ui.mount(MESSAGE('Accept MG677R failed'))
    expect(await pane.find({ type: 'Text', text: 'Accept MG677R failed' })).toBeDefined()
    expect(await pane.find({ type: 'Text', text: 'the definition of done is not met: criterion 1 is still open' })).toBeDefined()
    await pane.press({ key: 'message-close' })
    expect(w.closed).toEqual(['cockpit-message'])
  })

  test('Accept reads the dry run past the warning an older jaira prints before it', async ($, on) => {
    const w = world(on)
    w.warns = true
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R').map(r => r[5])).toEqual(['--dry-run', '--force'])
    expect(w.opened.filter(o => o.title === 'Accept MG677R failed')).toEqual([])
  })

  test('an error is shown whole in a pane of its own, never cut to a line', async ($, on) => {
    const w = world(on)
    const long = `the definition of done is not met: ${'criterion still open, '.repeat(20)}and one more`
    w.refusals = [...w.refusals, { code: 'needs_nonmodel_signal', message: long }, { code: 'needs_commits', message: 'no commit names MG677R' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-MG677R' })
    const pane = await $.ui.mount(MESSAGE('Accept MG677R failed'))
    expect(await pane.find({ type: 'Text', text: `${long}\nno commit names MG677R` })).toBeDefined()
    expect(await ui.find({ type: 'Text', text: /criterion still open/ })).toBeUndefined()
  })

  test('a plain yes gets no Explain, however long the sentence', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...SIGNOFF, review: { verdict: 'Approve. Der Diff erfuellt die DoD: Befehl existiert, Test deckt ihn ab.' } },
      { ...SIGNOFF, handle: 'AB12CD', review: { verdict: 'the diff matches the criteria; no defects found' } },
      { ...SIGNOFF, handle: 'EF34GH', review: { verdict: 'request changes: the button skips review' } },
    ]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'explain-MG677R' })).toBeUndefined()
    expect(await ui.find({ key: 'explain-AB12CD' })).toBeUndefined()
    expect(await ui.find({ key: 'explain-EF34GH' })).toBeDefined()
  })

  // The band shows three tickets, so each case below fits on it.
  test('a negated reservation gets no Explain', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...SIGNOFF, handle: 'GH56JK', review: { verdict: 'approve, no changes needed' } },
      { ...SIGNOFF, handle: 'JK78LM', review: { verdict: 'Approve: no open questions, ohne Vorbehalt' } },
    ]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    for (const id of ['GH56JK', 'JK78LM']) {
      expect(await ui.find({ key: `accept-${id}` })).toBeDefined()
      expect(await ui.find({ key: `explain-${id}` })).toBeUndefined()
    }
  })

  test('a negation ends at a comma and never takes the connective with it', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...SIGNOFF, handle: 'LM90NP', review: { verdict: 'Approve - no blockers, but the docs are missing' } },
      { ...SIGNOFF, handle: 'NP12QR', review: { verdict: 'Approve: no defects, however the flaky test remains' } },
      { ...SIGNOFF, handle: 'QR34ST', review: { verdict: 'Annehmen, keine Fehler, aber DoD 2 nie live gesehen' } },
    ]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    for (const id of ['LM90NP', 'NP12QR', 'QR34ST']) expect(await ui.find({ key: `explain-${id}` })).toBeDefined()
  })

  test('a verdict with a reservation gets Explain, which opens a Herdr tab of its own', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'explain-K7P2QX' })).toBeUndefined()
    await ui.press({ key: 'explain-MG677R' })
    expect(ran(w, SPAWN)).toEqual([[SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'explain', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('an open question shows under its ticket with Discuss', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ type: 'Text', text: /\? Popup oder eigene Seite\?/ })).toBeDefined()
    await ui.press({ key: 'discuss-K7P2QX' })
    // The chat sends its decision to this session, by its messaging socket.
    expect(ran(w, SPAWN)[0]?.slice(0, 7)).toEqual([SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'k7p2qx', 'K7P2QX', 'discuss'])
  })

  test('a Discuss chat\'s decision comes back: its tab closes and the band offers it to read, once', async ($, on) => {
    const w = world(on, { herdr: true })
    // 1_000_000 ms is 1970-01-01 00:16:40 UTC: jaira stamps notes to the minute.
    w.bodies.K7P2QX = '## Progress\n\n- **1970-01-01 00:05 · Alex** — decision: an older one\n'
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'working' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'discuss-K7P2QX' })
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    // The chat noted the decision and moved the ticket back itself.
    // Noted first: the chat still moves the ticket and tells this session.
    w.bodies.K7P2QX += '- **1970-01-01 00:17 · Alex** — decision: a page of its own, back to in-progress\n'
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(await (await $.ui.mount(BAND)).find({ key: 'decision-K7P2QX' })).toBeUndefined()
    w.tickets = [SIGNOFF, WORKING, { ...HUMAN, status: 'in-progress', question: '' }]
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'idle' }]
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([['herdr', 'tab', 'close', 'w1:t5', '@/repo']])
    const band = await $.ui.mount(BAND)
    expect(await band.find({ type: 'Text', text: /K7P2QX \| decided in Discuss/ })).toBeDefined()
    await band.press({ key: 'decision-K7P2QX' })
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: 'Decision on K7P2QX', focus: true })
    const pane = await $.ui.mount(MESSAGE('Decision on K7P2QX'))
    expect(await pane.find({ type: 'Text', text: 'a page of its own, back to in-progress' })).toBeDefined()
    expect(await (await $.ui.mount(BAND)).find({ key: 'decision-K7P2QX' })).toBeUndefined()
  })

  test('without Herdr, Explain says so instead of failing', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'explain-MG677R' })
    expect(ran(w, SPAWN)).toEqual([])
    expect(w.opened).toEqual([{ id: 'cockpit-message', title: 'Explain MG677R failed', focus: true }])
    const pane = await $.ui.mount(MESSAGE('Explain MG677R failed'))
    expect(await pane.find({ type: 'Text', text: /needs Herdr/ })).toBeDefined()
  })

  test('what the Explain chat wrote onto the ticket is read in the cockpit', async ($, on) => {
    const w = world(on)
    w.warns = true
    w.bodies.MG677R = [
      '## Progress',
      '',
      '- **2026-10-12 14:00 · Alex** — explain: an old one',
      '- **2026-10-12 15:00 · Alex** — explain: The key leaves out the branch.',
      'Two branches share one cache entry.',
      '- **2026-10-12 15:05 · Alex** — something else',
    ].join('\n')
    await startSession($)
    const ui = await $.ui.mount(BAND)
    // No explanation on K7P2QX: no button there.
    expect(await ui.find({ key: 'read-K7P2QX' })).toBeUndefined()
    await ui.press({ key: 'read-MG677R' })
    expect(w.opened).toEqual([{ id: 'cockpit-message', title: 'Explanation of MG677R', focus: true }])
    const pane = await $.ui.mount(MESSAGE('Explanation of MG677R'))
    expect(await pane.find({ type: 'Text', text: 'The key leaves out the branch.\nTwo branches share one cache entry.' })).toBeDefined()
  })

  test('the Explain tab closes once its explanation is on the ticket, not on an older one', async ($, on) => {
    const w = world(on, { herdr: true })
    // 1_000_000 ms is 1970-01-01 00:16:40 UTC: jaira stamps notes to the minute.
    w.bodies.MG677R = '## Progress\n\n- **1970-01-01 00:10 · Alex** — explain: an older one\n'
    w.tabs = [{ tab_id: 'w1:t4', label: 'MG677R/explain', agent_status: 'working' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'explain-MG677R' })
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    w.bodies.MG677R += '- **1970-01-01 00:16 · Alex** — explain: the key leaves out the branch\n'
    // Still answering in its tab: it stays until its turn is over.
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    w.tabs = [{ tab_id: 'w1:t4', label: 'MG677R/explain', agent_status: 'idle' }]
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([['herdr', 'tab', 'close', 'w1:t4', '@/repo']])
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toHaveLength(1)
  })

  test('the Discuss tab closes once the ticket\'s question is answered', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'working' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'discuss-K7P2QX' })
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    w.tickets = [SIGNOFF, { ...HUMAN, question: '' }, WORKING]
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'done' }]
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([['herdr', 'tab', 'close', 'w1:t5', '@/repo']])
  })

  test('only the person\'s own and nobody\'s tickets get buttons and count; the teammates\' sit in one line that opens', async ($, on) => {
    const w = world(on)
    w.me = { identity: 'Alexander Sacharov', aliases: ['Alexander Sacharov', 'alex@example.org'] }
    w.tickets = [
      { ...SIGNOFF, assignee: 'alex@EXAMPLE.org' },
      { ...SIGNOFF, handle: 'AB12CD', assignee: '' },
      { ...SIGNOFF, handle: 'EF34GH', assignee: 'Berk' },
      { ...HUMAN, assignee: 'Carla' },
      WORKING,
    ]
    await startSession($)
    expect(await headLine($)).toBe('0 worker | 0 review | 2 waiting for you | 5h 34%')
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'accept-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'accept-AB12CD' })).toBeDefined()
    expect(await ui.find({ key: 'accept-EF34GH' })).toBeUndefined()
    expect(await ui.find({ key: 'discuss-K7P2QX' })).toBeUndefined()
    expect((await ui.find({ key: 'theirs' }))?.props.label).toBe('2 waiting for Berk, Carla')
    await ui.press({ key: 'theirs' })
    const open = await $.ui.mount(BAND)
    expect(await open.find({ type: 'Text', text: /EF34GH \| signoff \| Berk/ })).toBeDefined()
    // A teammate's ticket: talked through, never accepted or given back from here.
    expect(await open.find({ key: 'discuss-K7P2QX' })).toBeDefined()
    for (const id of ['EF34GH', 'K7P2QX']) {
      expect(await open.find({ key: `accept-${id}` })).toBeUndefined()
      expect(await open.find({ key: `back-${id}` })).toBeUndefined()
    }
  })

  test('a Discuss chat on a teammate\'s ticket stays open while it waits for the person', async ($, on) => {
    const w = world(on, { herdr: true })
    w.me = { identity: 'Alexander Sacharov', aliases: ['Alexander Sacharov'] }
    w.tickets = [SIGNOFF, { ...HUMAN, assignee: 'Carla' }, WORKING]
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'idle' }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'theirs' })
    await (await $.ui.mount(BAND)).press({ key: 'discuss-K7P2QX' })
    await w.clock.advance(60_000)
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    w.tickets = [SIGNOFF, { ...HUMAN, assignee: 'Carla', question: '' }, WORKING]
    await w.clock.advance(60_000)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([['herdr', 'tab', 'close', 'w1:t5', '@/repo']])
  })

  test('without a word from jaira whoami every waiting ticket counts as the person\'s', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, assignee: 'Berk' }, HUMAN, WORKING]
    await startSession($)
    expect(await headLine($)).toBe('0 worker | 0 review | 2 waiting for you | 5h 34%')
    expect(await (await $.ui.mount(BAND)).find({ key: 'theirs' })).toBeUndefined()
  })

  test('+N more opens every waiting ticket with its buttons, and show fewer closes them again', async ($, on) => {
    const w = world(on)
    w.tickets = [
      SIGNOFF,
      HUMAN,
      { ...SIGNOFF, handle: 'AB12CD', review: { verdict: 'approve' } },
      { ...SIGNOFF, handle: 'EF34GH', review: { verdict: 'approve' } },
      { ...SIGNOFF, handle: 'GH56JK', review: { verdict: 'approve' } },
    ]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'accept-GH56JK' })).toBeUndefined()
    expect((await ui.find({ key: 'more' }))?.props.label).toBe('+2 more')
    await ui.press({ key: 'more' })
    const open = await $.ui.mount(BAND)
    expect(await open.find({ key: 'accept-EF34GH' })).toBeDefined()
    expect(await open.find({ key: 'back-GH56JK' })).toBeDefined()
    await open.press({ key: 'fewer' })
    const closed = await $.ui.mount(BAND)
    expect(await closed.find({ key: 'accept-GH56JK' })).toBeUndefined()
    expect(await closed.find({ key: 'more' })).toBeDefined()
  })

  test('a ticket nobody explained yet offers no explanation to read', async ($, on) => {
    world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'explain-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'read-MG677R' })).toBeUndefined()
  })
})

describe('after Accept', () => {
  const press = async ($: Test$) => {
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-MG677R' })
  }
  const prLane = (w: World) => ran(w, SPAWN).filter(r => r.includes('pr'))
  const failure = async ($: Test$) => {
    const pane = await $.ui.mount(MESSAGE('Next step after MG677R failed'))
    return (await pane.findAll({ type: 'Text' })).map(x => x.text).join('\n')
  }

  test('a branch ahead of the target without a pull request: the press opens one, by a person, in the ticket\'s worktree', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await press($)
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'upstream/main..feat/mg677r', '@/repo']])
    expect(ran(w, 'gh')).toEqual([
      ['gh', 'repo', 'view', 'git@github.com:acme/repo.git', '--json', 'isFork,parent,nameWithOwner,defaultBranchRef', '@/repo'],
      ['gh', 'repo', 'view', 'upstream/repo', '--json', 'isFork,parent,nameWithOwner,defaultBranchRef', '@/repo'],
      ['gh', 'pr', 'list', '--repo', 'upstream/repo', '--head', 'feat/mg677r', '--state', 'all', '--json', 'number', '@/repo'],
    ])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
    // After the move, never before it: the press accepted first.
    const order = w.ran.map(r => r[0] === SPAWN ? 'spawn' : r.includes('--force') ? 'move' : '').filter(Boolean)
    expect(order).toEqual(['move', 'spawn'])
    expect(w.opened).toEqual([])
  })

  test('a pull request in any state already there: the next ticket starts, under the dispatch setting', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prs = [{ number: 12 }]
    w.next = ['T8QS1Y', 'AB12CD']
    w.tabs = [{ tab_id: 't9', label: 'T8QS1Y/worker' }]
    w.files.set('/home/a/.jaira/settings.json', '{"dispatch":"worker"}')
    await press($)
    expect(prLane(w)).toEqual([])
    expect(ran(w, 'jaira', 'next')).toEqual([['jaira', 'next', '--all', '--json', '@/repo']])
    // T8QS1Y has a worker tab already: never a second agent on it.
    expect(ran(w, SPAWN)).toEqual([[SPAWN, 'ab12cd', 'AB12CD', 'worker', '/repo', '@/repo']])
    expect(w.opened).toEqual([])
  })

  test('nothing ahead of the target: no pull request is asked for, and without a worker setting the dispatcher starts the next one', async ($, on) => {
    const w = world(on, { herdr: true })
    w.next = ['AB12CD']
    await press($)
    expect(ran(w, 'gh', 'pr')).toEqual([])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, 'ab12cd', 'AB12CD', 'dispatch', '/repo', '@/repo']])
  })

  test('GitLab, by jaira.forge or by host, asks glab', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 1
    w.forge = 'gitlab'
    await press($)
    expect(ran(w, 'glab', 'mr')).toEqual([['glab', 'mr', 'list', '--repo', 'upstream/repo', '--source-branch', 'feat/mg677r', '--all', '-F', 'json', '@/repo']])
    expect(prLane(w)).toHaveLength(1)
  })

  test('a worktree removed after done leaves its branch: the branch is found, checked out again beside the repo, and the pull request opens there', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }]
    w.branches = ['feat/zz99zz']
    w.subjects = ['fix(ZZ99ZZ): one']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-ZZ99ZZ' })
    expect(ran(w, 'git', 'branch')).toEqual([['git', 'branch', '--list', '--format=%(refname:short)', '*/zz99zz', '@/repo']])
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'upstream/main..feat/zz99zz', '@/repo']])
    expect(ran(w, 'git', 'worktree', 'add')).toEqual([['git', 'worktree', 'add', '/.worktrees/repo-zz99zz', 'feat/zz99zz', '@/repo']])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, '--no-worktree', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-zz99zz', '@/.worktrees/repo-zz99zz']])
  })

  test('two tickets on one branch: no pull request while the other is not in done, the next ticket starts, and the press says which one is missing', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 4
    w.subjects = ['feat(MG677R): cache key', 'fix(T8QS1Y): filter', 'feat(T8QS1Y): filter', 'docs: no ticket']
    w.next = ['AB12CD']
    await press($)
    expect(ran(w, 'gh', 'pr')).toEqual([])
    expect(prLane(w)).toEqual([])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, 'ab12cd', 'AB12CD', 'dispatch', '/repo', '@/repo']])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; AB12CD started in a Herdr tab; no pull request for feat/mg677r yet: T8QS1Y not in done')
  })

  test('two tickets on one branch, both in done: the pull request opens; a handle the board no longer lists is filed away and holds nothing up', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 4
    w.tickets = [SIGNOFF, { ...WORKING, status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'fix(T8QS1Y): filter', 'chore(OLD123): filed away long ago']
    await press($)
    expect(prLane(w)).toHaveLength(1)
  })

  test('a ticket worked on another ticket\'s branch finds that branch by the commits naming it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }, { ...WORKING, handle: 'MG677R', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'fix(ZZ99ZZ): on the same branch']
    w.branchCommits = [{ sha: 'abc123', subject: 'fix(ZZ99ZZ): on the same branch' }]
    w.holders = ['master', 'feat/mg677r']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-ZZ99ZZ' })
    expect(ran(w, 'git', 'log', '--branches')).toEqual([['git', 'log', '--branches', '--not', 'upstream/main', '--format=%H', '-n', '1', '-E', '--grep=^[a-z]+\\(([^)]*[ ,])?ZZ99ZZ([ ,][^)]*)?\\)!?:', '@/repo']])
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([['git', 'branch', '--contains', 'abc123', '--format=%(refname:short)', '@/repo']])
    // Its worktree is the other ticket's: the role pushes from there.
    expect(ran(w, SPAWN)).toEqual([[SPAWN, '--no-worktree', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('the branch search finds a ticket named beside others in one subject, and never a handle that only starts the same', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }, { ...WORKING, handle: 'MG677R', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'chore(MG677R,ZZ99ZZ): both']
    w.branchCommits = [
      { sha: 'wrong1', subject: 'fix(ZZ99ZZX): a longer handle' },
      { sha: 'wrong2', subject: 'docs: mentions (ZZ99ZZ) in passing' },
      { sha: 'abc123', subject: 'chore(MG677R,ZZ99ZZ): both' },
    ]
    w.holders = ['feat/mg677r']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-ZZ99ZZ' })
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([['git', 'branch', '--contains', 'abc123', '--format=%(refname:short)', '@/repo']])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, '--no-worktree', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('in a fork, ahead counts against the parent\'s default branch, fetched first — never the lagging local master', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await press($)
    expect(ran(w, 'git', 'fetch')).toEqual([['git', 'fetch', 'upstream', 'main', '@/repo']])
    const order = w.ran.map(r => r[1] === 'fetch' ? 'fetch' : r[1] === 'log' ? 'log' : '').filter(Boolean)
    expect(order).toEqual(['fetch', 'log'])
    expect(ran(w, 'git', 'log')).toEqual([['git', 'log', '--format=%s', 'upstream/main..feat/mg677r', '@/repo']])
  })

  test('not a fork: the base is origin\'s own default branch', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 1
    w.origin = 'git@github.com:upstream/repo.git'
    w.remotes = 'origin\tgit@github.com:upstream/repo.git (fetch)\n'
    await press($)
    expect(ran(w, 'git', 'fetch')).toEqual([['git', 'fetch', 'origin', 'main', '@/repo']])
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'origin/main..feat/mg677r', '@/repo']])
  })

  test('a fork without its parent as a remote: nothing is counted, the next ticket starts, and the press says which remote to add', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.remotes = 'origin\tgit@github.com:acme/repo.git (fetch)\n'
    w.next = ['AB12CD']
    await press($)
    expect(ran(w, 'git', 'log')).toEqual([])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, 'ab12cd', 'AB12CD', 'dispatch', '/repo', '@/repo']])
    expect(w.toasts.at(-1)).toContain('Pull requests go to upstream/repo, and no remote here points at it')
  })

  test('no gh: the pull request goes unchecked, the next ticket starts anyway, and the press says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prs = null
    w.next = ['AB12CD']
    await press($)
    expect(prLane(w)).toEqual([])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, 'ab12cd', 'AB12CD', 'dispatch', '/repo', '@/repo']])
    expect(ran(w, 'jaira', 'move', 'MG677R').map(r => r[5])).toEqual(['--dry-run', '--force'])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; AB12CD started in a Herdr tab; pull requests not checked (gh repo view: gh: command not found)')
  })

  test('a forge that answers the repository but not the pull requests: the next ticket starts all the same', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prListFails = true
    w.next = ['AB12CD']
    await press($)
    expect(prLane(w)).toEqual([])
    expect(ran(w, SPAWN)).toEqual([[SPAWN, 'ab12cd', 'AB12CD', 'dispatch', '/repo', '@/repo']])
    expect(w.toasts.at(-1)).toContain('pull requests not checked (gh: HTTP 401')
  })

  test('a subject naming several tickets holds the pull request for each of them', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [SIGNOFF, WORKING, { ...WORKING, handle: 'AB12CD', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'chore(AB12CD,T8QS1Y): both at once']
    w.next = ['QR34ST']
    await press($)
    expect(prLane(w)).toEqual([])
    expect(w.toasts.at(-1)).toContain('no pull request for feat/mg677r yet: T8QS1Y not in done')
  })

  test('no remote: nothing starts, and the pane says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    await press($)
    expect(ran(w, SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('a host neither GitHub nor GitLab: nothing starts, and the pane names the setting that settles it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = 'git@git.example.de:acme/repo.git'
    await press($)
    expect(ran(w, SPAWN)).toEqual([])
    expect(await failure($)).toContain('git config jaira.forge')
  })

  test('no next ticket: nothing starts, and the pane says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.next = ['T8QS1Y']
    w.tabs = [{ tab_id: 't9', label: 'T8QS1Y/dispatch' }]
    await press($)
    expect(ran(w, SPAWN)).toEqual([])
    expect(await failure($)).toContain('No ticket can be started next')
  })

  test('no Herdr: nothing starts, and the pane says the step needs it', async ($, on) => {
    const w = world(on)
    w.ahead = 3
    await press($)
    expect(ran(w, SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs a pull request')
  })

  test('a refused Accept starts no next step at all', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.refusals = [...w.refusals, { code: 'needs_commits', message: 'no commit names MG677R' }]
    await press($)
    expect(ran(w, 'git', 'log')).toEqual([])
    expect(ran(w, SPAWN)).toEqual([])
  })

  test('an agent\'s session has no Accept, so it can never open a pull request', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await startSession($)
    await $.prompt.submit({ text: '/jaira-worker MG677R --parent jaira-4c', origin: { kind: 'composer' }, wait: false })
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'accept-MG677R' })).toBeUndefined()
    expect(prLane(w)).toEqual([])
  })
})

describe('give back', () => {
  async function giveBack($: Test$, w: World) {
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'back-MG677R' })
    const pane = await $.ui.mount(GIVE('Give back MG677R'))
    const typed = pane.input({ key: 'reason-MG677R', text: 'der Branch fehlt im Schluessel' })
    // Herdr is given a second between the text and its Enter.
    await w.clock.advance(1000)
    await typed
    return w
  }

  test('asks what is missing in a pane that takes the keyboard, never in a band field the press leaves unfocused', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'back-MG677R' })
    // Live, the band's own field never got the ring: the press keeps it on
    // the button. A pane opened with focus hands it to its autoFocus field.
    expect(await ui.find({ key: 'reason-MG677R' })).toBeUndefined()
    expect(w.opened).toEqual([{ id: 'cockpit-give-back', title: 'Give back MG677R', focus: true }])
    const pane = await $.ui.mount(GIVE('Give back MG677R'))
    expect((await pane.find({ key: 'reason-MG677R' }))?.props.autoFocus).toBe(true)
    // Enter on an empty field gives nothing back and keeps the pane open.
    await pane.input({ key: 'reason-MG677R', text: '  ' })
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(w.closed).toEqual([])
    await pane.input({ key: 'reason-MG677R', text: 'fehlt' })
    expect(w.closed).toEqual(['cockpit-give-back'])
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.slice(0, 5)).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress'])
  })

  test('notes why and moves the ticket back to in-progress, without --force, so any other refusal still stops it', async ($, on) => {
    const w = world(on)
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'note', 'MG677R')[0]?.[3]).toBe('Given back from signoff: der Branch fehlt im Schluessel')
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress', '@/.worktrees/repo-mg677r'])
  })

  // The limit is a share of the old worker's window, never a fixed size.
  const big = (window: number) => ({ 'session:old': {
    session: 'old', pane: 'w1:p9', cwd: '/x', repo: '/repo', ticket: 'MG677R', role: 'worker', lane: '',
    lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 300_000, window, ttl: 3_600_000, updatedAt: 1_000_000 - 20 * 60_000,
  } })

  test('a context under half of its window goes back to the same worker, however many tokens that is', async ($, on) => {
    const w = world(on, { herdr: true, store: big(1_000_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')).toHaveLength(1)
    expect(ran(w, SPAWN)).toEqual([])
  })

  test('past half of its window a new worker takes the round', async ($, on) => {
    const w = world(on, { herdr: true, store: big(500_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe("cockpit: given back to a new worker - old worker's context 300k >= 250k (50% of its window)")
  })

  test('the share is a setting', { options: { workerContextShare: 70 } }, async ($, on) => {
    const w = world(on, { herdr: true, store: big(500_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')).toHaveLength(1)
  })

  test('goes to the worker that carried it while its tab lives, its cache is warm and its context small', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 20 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')).toHaveLength(1)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(ran(w, SPAWN)).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to the old worker - tab alive, cache warm for 40 more min, context 80k')
  })

  test('otherwise closes the old worker\'s tab first, then starts a new one', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 70 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 70 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await giveBack($, w)
    const close = w.ran.findIndex(r => r[0] === 'herdr' && r[1] === 'tab' && r[2] === 'close')
    const spawn = w.ran.findIndex(r => r[0] === SPAWN)
    expect(w.ran[close]?.[3]).toBe('w1:t9')
    expect(close).toBeLessThan(spawn)
    expect(w.ran[spawn]?.slice(0, 5)).toEqual([SPAWN, 'mg677r', 'MG677R', 'worker', '/repo'])
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to a new worker - old worker\'s prompt cache cold for 10 min')
  })
  test('a pane gone since the tab list takes no text: the old tab is closed and a new worker started', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 20 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    w.paneGone = true
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-keys')).toEqual([])
    expect(ran(w, 'herdr', 'tab', 'close')[0]?.[3]).toBe('w1:t9')
    expect(ran(w, SPAWN)[0]?.slice(0, 5)).toEqual([SPAWN, 'mg677r', 'MG677R', 'worker', '/repo'])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to a new worker - tab alive, cache warm for 40 more min, context 80k, but its pane took no text')
  })
})

describe('give back, the edges', () => {
  test('a record a day old decides nothing and is dropped', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 999_000, tokens: 1, ttl: 3_600_000, updatedAt: 1_000_000 - 25 * 3_600_000,
    } } })
    await startSession($)
    expect(w.store.has('session:old')).toBe(false)
  })

  test('without Herdr it says the ticket is back in in-progress and nobody was started', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'back-MG677R' })
    await (await $.ui.mount(GIVE('Give back MG677R'))).input({ key: 'reason-MG677R', text: 'fehlt' })
    expect(await ui.find({ type: 'Text', text: "MG677R back in in-progress; no worker started (the ticket's note says why)" })).toBeDefined()
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to nobody (no worker started) - no Herdr in this session')
    expect(ran(w, SPAWN)).toEqual([])
  })
})

describe('give back with nothing to start a worker', () => {
  test('keeps the old tab and says nobody was started', async ($, on) => {
    const w = world(on, { herdr: true, spawn: false, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 70 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 70 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'back-MG677R' })
    await (await $.ui.mount(GIVE('Give back MG677R'))).input({ key: 'reason-MG677R', text: 'fehlt' })
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toStartWith('cockpit: given back to nobody (no worker started) - old worker\'s prompt cache cold for 10 min; no spawn.sh found')
  })
})

describe('the workers of the board', () => {
  const record = (session: string, ticket: string, extra: Record<string, unknown>) => ({
    session, pane: `w1:${session}`, cwd: '/x', repo: '/repo', ticket, role: 'worker', lane: '', ttl: 3_600_000, window: 1_000_000, ...extra,
  })

  test('lists every agent carrying a ticket here with lane, context, cache and whom a new round goes to', async ($, on) => {
    const w = world(on, { herdr: true, store: {
      // Idle in signoff, cache warm: the round would go back to it.
      'session:a': record('a', 'MG677R', { lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 80_000, updatedAt: 1_000_000 - 20 * 60_000 }),
      // At work: nothing is decided about it.
      'session:b': record('b', 'T8QS1Y', { lastAnswerAt: 1_000_000 - 90 * 60_000, tokens: 120_000, updatedAt: 1_000_000 - 90 * 60_000 }),
      // Another repository's worker, and one whose tab is gone: neither listed.
      'session:c': record('c', 'ZZZZZZ', { repo: '/other', lastAnswerAt: 999_000, tokens: 1, updatedAt: 999_000 }),
      'session:d': record('d', 'K7P2QX', { lastAnswerAt: 999_000, tokens: 1, updatedAt: 999_000 }),
      // A single lane's agent carries no ticket as a whole.
      'session:e': { ...record('e', 'MG677R', { lastAnswerAt: 999_000, tokens: 1, updatedAt: 999_999 }), role: 'lane', lane: 'review' },
    } })
    w.tabs = [
      { tab_id: 'w1:t1', label: 'MG677R/worker', agent_status: 'idle' },
      { tab_id: 'w1:t2', label: 'T8QS1Y/worker', agent_status: 'working' },
      { tab_id: 'w1:t3', label: 'ZZZZZZ/worker', agent_status: 'idle' },
    ]
    w.usage.tokens = 142_000
    await startSession($)
    const { result } = await $.tool.call({ tool: 'mcp__jaira-cockpit__workers' })
    expect(String(result).split('\n')).toEqual([
      'this session: context 142k of 1000k, cache warm while it works and 60 min after its last answer',
      'Agents carrying a ticket of this board:',
      'MG677R · signoff · worker · context 80k of 1000k · cache 40 min left · next round: same worker (tab alive, cache warm for 40 more min, context 80k)',
      'T8QS1Y · in-progress · worker · context 120k of 1000k · cache warm (working) · working: nothing to decide',
    ])
  })

  test('a cold worker waiting on a person is still listed, with a new worker for its next round', async ($, on) => {
    const w = world(on, { herdr: true, store: {
      'session:a': record('a', 'MG677R', { lastAnswerAt: 1_000_000 - 70 * 60_000, tokens: 80_000, updatedAt: 1_000_000 - 70 * 60_000 }),
    } })
    w.tabs = [{ tab_id: 'w1:t1', label: 'MG677R/worker', agent_status: 'idle' }]
    await startSession($)
    const text = (await $.command.run(typed('workers'))).text
    expect(text).toContain("MG677R · signoff · worker · context 80k of 1000k · cache cold · next round: new worker (old worker's prompt cache cold for 10 min)")
  })

  test('a worker writes the board it belongs to and its window into its record', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await $.prompt.submit({ text: '/jaira-worker MG677R --parent jaira-4c', origin: { kind: 'composer' }, wait: false })
    w.usage.tokens = 50_000
    await answer($, 'done')
    expect(w.store.get('session:abcdef12')).toMatchObject({ repo: '/repo', window: 1_000_000, tokens: 50_000 })
    // Asked from inside its own turn, it sees itself first.
    expect(String((await $.tool.call({ tool: 'mcp__jaira-cockpit__workers' })).result).split('\n')[0]).toBe(
      'this session: context 50k of 1000k, cache warm while it works and 60 min after its last answer',
    )
  })
})

describe("the cockpit's own words in a worker's pane", () => {
  test('are no person\'s prompt: never held, never filed as a thread', async ($, on) => {
    const w = world(on, { board: false })
    await startSession($)
    w.usage.tokens = 40_000
    w.usage.rateLimits = []
    await answer($)
    await w.clock.advance(7 * 60_000)
    w.answer = 'Cancel'
    const back = '[jaira-cockpit] Ticket MG677R kam aus signoff zurueck nach in-progress. Grund: fehlt.'
    expect((await $.prompt.submit({ text: back, origin: { kind: 'composer' }, wait: false })).text).toBe(back)
    await answer($, 'ok, mache weiter')
    await w.clock.settle()
    expect(w.asked).toEqual([])
    expect(w.modelCalls).toBe(0)
  })
})

describe('context and cache', () => {
  test('the band shows the context and how long the cache stays warm, and when to compact', async ($, on) => {
    const w = world(on, { board: false })
    await startSession($)
    w.usage.tokens = 142_000
    await answer($)
    const ui = await $.ui.mount(BAND)
    // What asks for action is drawn in the warning colour, a new session in red.
    expect(await sessionLine(ui)).toEqual({ text: 'context 142k | cache 60 min left | /compact advised', colors: ['', '', 'warning'] })
    w.usage.tokens = 320_000
    w.usage.rateLimits = []
    await answer($)
    const again = await $.ui.mount(BAND)
    expect(await sessionLine(again)).toEqual({ text: 'context 320k | cache 5 min left | start a new session', colors: ['', '', 'error'] })
  })
})

describe('the cache guard', () => {
  async function cold($: Test$, w: World, tokens: number, on: On) {
    await startSession($)
    w.usage.tokens = tokens
    await answer($)
    // An API key's five minutes, and seven minutes go by.
    w.usage.rateLimits = []
    await answer($)
    void on
  }

  test('holds a person\'s prompt onto a cold cache with a large context, and Cancel keeps it in the prompt', async ($, on) => {
    const w = world(on, { board: false })
    const clock = w.clock
    await cold($, w, 40_000, on)
    await clock.advance(7 * 60_000)
    w.answer = 'Cancel'
    const r = await $.prompt.submit({ text: 'and now the tests', origin: { kind: 'composer' }, wait: false })
    expect(w.asked).toHaveLength(1)
    expect(w.asked[0]).toContain('40k tokens')
    expect(r.drop).toBe('Cancel')
    expect(w.filled).toEqual(['and now the tests'])
  })

  test('never holds an agent, a slash command, a small context or a warm cache', async ($, on) => {
    const w = world(on, { board: false })
    const clock = w.clock
    await cold($, w, 40_000, on)
    await clock.advance(7 * 60_000)
    w.answer = 'Cancel'
    await $.prompt.submit({ text: 'Ticket X done', origin: { kind: 'peer' }, wait: false })
    await $.prompt.submit({ text: '/jaira-worker MG677R', origin: { kind: 'composer' }, wait: false })
    await $.prompt.submit({ text: 'scheduled', origin: { kind: 'scheduled-trigger' }, wait: false })
    expect(w.asked).toEqual([])
    w.usage.tokens = 20_000
    await answer($)
    await clock.advance(7 * 60_000)
    expect((await $.prompt.submit({ text: 'small', origin: { kind: 'composer' }, wait: false })).text).toBe('small')
    expect(w.asked).toEqual([])
  })

  test('Compact first and Start over keep the words in the prompt and run /compact or /clear', async ($, on) => {
    const w = world(on, { board: false })
    const ran: string[] = []
    on('session.compact', () => {
      ran.push('compact')
      return { skip: 'test' } as never
    })
    on('command.run', { command: 'clear' }, () => {
      ran.push('clear')
      return { text: '' }
    })
    await cold($, w, 40_000, on)
    await w.clock.advance(7 * 60_000)
    w.answer = 'Compact first'
    expect((await $.prompt.submit({ text: 'first', origin: { kind: 'composer' }, wait: false })).drop).toBe('Compact first')
    await w.clock.settle()
    w.answer = 'Start over (/clear)'
    expect((await $.prompt.submit({ text: 'second', origin: { kind: 'composer' }, wait: false })).drop).toBe('Start over (/clear)')
    await w.clock.settle()
    expect(ran).toEqual(['compact', 'clear'])
    expect(w.filled).toEqual(['first', 'second'])
  })

  test('Send lets it through', async ($, on) => {
    const w = world(on, { board: false })
    const clock = w.clock
    await cold($, w, 40_000, on)
    await clock.advance(7 * 60_000)
    w.answer = 'Send'
    expect((await $.prompt.submit({ text: 'go on', origin: { kind: 'composer' }, wait: false })).text).toBe('go on')
  })

  test('the threshold is a setting', { options: { cacheGuardTokens: 50_000 } }, async ($, on) => {
    const w = world(on, { board: false })
    await cold($, w, 40_000, on)
    await w.clock.advance(7 * 60_000)
    w.answer = 'Cancel'
    expect((await $.prompt.submit({ text: 'go on', origin: { kind: 'composer' }, wait: false })).text).toBe('go on')
    expect(w.asked).toEqual([])
  })
})

describe('the head line, in the band\'s frame', () => {
  test('counts the agents that run and the tickets waiting for a person, and what was used', async ($, on) => {
    const w = world(on, { herdr: true })
    // T8QS1Y sits in in-progress and MG677R in signoff: lanes count no agent.
    // A tab whose agent finished its turn stays open and counts for nothing.
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' },
      { tab_id: 'w1:t2', label: 'K7P2QX/dispatch', agent_status: 'working' },
      { tab_id: 'w1:t3', label: 'MG677R/review', agent_status: 'working' },
      { tab_id: 'w1:t4', label: 'MG677R/explain', agent_status: 'working' },
      { tab_id: 'w1:t5', label: 'ZZ99ZZ/worker', agent_status: 'working' },
      { tab_id: 'w1:t6', label: 'notes', agent_status: 'working' },
      { tab_id: 'w1:t7', label: 'MG677R/worker', agent_status: 'done' },
      { tab_id: 'w1:t8', label: 'T8QS1Y/review', agent_status: 'idle' },
      { tab_id: 'w1:t9', label: 'K7P2QX/review', agent_status: 'done' },
    ]
    await startSession($)
    expect(await headLine($)).toBe('2 worker | 1 review | 2 waiting for you | 5h 34%')
  })

  test('without Herdr the workers are this session\'s running subagents on the board\'s tickets', async ($, on) => {
    const w = world(on)
    w.agents = [
      { id: 'a1', description: 'NZZ T8QS1Y worker', type: 'claude', status: 'running' },
      { id: 'a2', description: 'carry a ticket', name: 'K7P2QX', type: 'claude', status: 'waiting' },
      { id: 'a3', description: 'MG677R review', type: 'claude', status: 'running' },
      { id: 'a4', description: 'MG677R critique', type: 'claude', status: 'completed' },
      { id: 'a5', description: 'search the code', type: 'Explore', status: 'running' },
    ]
    await startSession($)
    expect(await headLine($)).toBe('2 worker | 1 review | 2 waiting for you | 5h 34%')
  })

  test('a worker tab left open after its turn is no worker', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'done' },
      { tab_id: 'w1:t2', label: 'MG677R/review', agent_status: 'idle' },
    ]
    await startSession($)
    expect(await headLine($)).toBe('0 worker | 0 review | 2 waiting for you | 5h 34%')
  })

  test('a jaira built before human_exit still shows the tickets waiting in signoff and human', async ($, on) => {
    const w = world(on)
    w.oldJaira = true
    await startSession($)
    expect(await headLine($)).toBe('0 worker | 0 review | 2 waiting for you | 5h 34%')
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'accept-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'discuss-K7P2QX' })).toBeDefined()
  })

  test('shows the seven-day window beside the five-hour one when the subscription reports it', async ($, on) => {
    const w = world(on, { board: false })
    w.usage.rateLimits = [{ kind: 'five_hour', percentUsed: 34 }, { kind: 'seven_day', percentUsed: 61.6 }]
    await startSession($)
    expect(await headLine($)).toBe('5h 34% | 7d 62%')
  })

  test('without a board it shows what was used alone', async ($, on) => {
    const w = world(on, { board: false })
    await startSession($)
    expect(await headLine($)).toBe('5h 34%')
    // One frame for everything: no status line of the mod's own below the prompt.
    expect(w.status).toEqual([])
  })
})

describe('a project without a board', () => {
  test('shows context and cache, and nothing of tickets or lanes, without an error', async ($, on) => {
    const w = world(on, { board: false })
    await startSession($)
    w.usage.tokens = 50_000
    await answer($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ text: /Waiting for you/ })).toBeUndefined()
    expect(await ui.find({ type: 'Button' })).toBeUndefined()
    expect(await ui.find({ type: 'Text', text: /context 50k/ })).toBeDefined()
    expect(ran(w, 'jaira', 'list')).toEqual([])
  })
})

describe('conversation threads', () => {
  test('each turn of the person is filed under a thread; the band counts open questions', async ($, on) => {
    const w = world(on)
    const clock = w.clock
    await startSession($)
    w.model = '{"thread": null, "title": "Tags als Popup", "line": "Idee: Tags als Popup", "ticket": null, "question": "Popup oder eigene Seite?", "answered": []}'
    await $.prompt.submit({ text: 'Koennen die Tags ein Popup werden?', origin: { kind: 'composer' }, wait: false })
    await answer($, 'Popup oder eigene Seite?')
    await clock.settle()
    const list = await $.command.run(typed('threads'))
    expect(list.text).toContain('? #1 Tags als Popup ? Popup oder eigene Seite?')
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ type: 'Text', text: /threads: 1 open question\(s\), 0 left lying/ })).toBeDefined()
    w.model = '{"thread": 1, "title": "Tags als Popup", "line": "Ticket K7P2QX angelegt", "ticket": "K7P2QX", "question": null, "answered": [1]}'
    await $.prompt.submit({ text: 'Popup.', origin: { kind: 'composer' }, wait: false })
    await answer($, 'Angelegt: K7P2QX')
    await clock.settle()
    expect((await $.command.run(typed('threads'))).text).toContain('-> #1 Tags als Popup -> K7P2QX')
  })

  test('an agent\'s report costs no model call; a thread with no ticket and no answer is left lying after five turns', async ($, on) => {
    const w = world(on, { board: false })
    const clock = w.clock
    await startSession($)
    w.model = '{"thread": null, "title": "Logo", "line": "Logo groesser?", "ticket": "ABCDEF", "question": null, "answered": []}'
    await $.prompt.submit({ text: 'Logo groesser?', origin: { kind: 'composer' }, wait: false })
    await answer($, 'Kann man machen.')
    await clock.settle()
    await $.prompt.submit({ text: 'Ticket MG677R moved to signoff', origin: { kind: 'peer' }, wait: false })
    await answer($, 'noted')
    for (let i = 0; i < 4; i++) await answer($, 'more')
    await clock.settle()
    expect(w.modelCalls).toBe(1)
    // No board: no thread stands on a ticket.
    expect((await $.command.run(typed('threads'))).text).toContain('~ #1 Logo')
  })
})

describe('a session that carries a ticket', () => {
  const START = { text: '/jaira-worker MG677R --parent jaira-4c', origin: { kind: 'composer' }, wait: false } as const

  test('passes every agent message on as it was written', async ($, on) => {
    world(on, { herdr: true })
    await startSession($)
    await $.prompt.submit(START)
    const order = 'Ticket MG677R: review -> signoff, review-verdict="approve". Commit the fix with the ticket.'
    expect((await $.session.receive({ origin: { kind: 'peer' }, text: order })).text).toBe(order)
    expect((await $.prompt.submit({ text: order, origin: { kind: 'peer' }, wait: false })).text).toBe(order)
  })

  test('stays an agent\'s session after the mod is reloaded, though the line that started it never comes again', async ($, on) => {
    // A reload starts the module over: only the store remembers the start line.
    const w = world(on, { herdr: true, store: { 'carrying:abcdef12': { ticket: 'MG677R', role: 'worker', lane: '' } } })
    await startSession($)
    const order = 'Ticket MG677R: review -> signoff, review-verdict="approve". Commit the fix with the ticket.'
    expect((await $.session.receive({ origin: { kind: 'peer' }, text: order })).text).toBe(order)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ text: /Waiting for you/ })).toBeUndefined()
    w.usage.tokens = 50_000
    await answer($, 'done')
    expect(w.store.get('session:abcdef12')).toMatchObject({ ticket: 'MG677R', role: 'worker' })
    await $.session.end({ reason: 'other' } as never)
    expect(w.store.has('carrying:abcdef12')).toBe(false)
  })

  test('keeps what started it in the store for a reload to read back', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await $.prompt.submit(START)
    expect(w.store.get('carrying:abcdef12')).toEqual({ ticket: 'MG677R', role: 'worker', lane: '' })
  })

  test('draws no band and no status line, files no thread and guards no cache, and still writes its record', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await $.prompt.submit(START)
    expect(w.status.at(-1)).toBeUndefined()
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ text: /Waiting for you/ })).toBeUndefined()
    expect(await ui.find({ type: 'Button' })).toBeUndefined()
    w.usage.tokens = 200_000
    await answer($, 'done')
    const before = w.ran.length
    // An hour and more later a person types into the worker's tab: no guard.
    await w.clock.advance(2 * 3_600_000)
    await $.prompt.submit({ text: 'how far are you?', origin: { kind: 'composer' }, wait: false })
    await answer($, 'nearly')
    await w.clock.settle()
    expect(w.asked).toEqual([])
    expect(w.modelCalls).toBe(0)
    expect(w.status.at(-1)).toBeUndefined()
    expect(w.ran.slice(before).filter(r => r[0] === 'jaira')).toEqual([])
    expect(w.store.get('session:abcdef12')).toMatchObject({ ticket: 'MG677R', role: 'worker', pane: 'w1:p1', tokens: 200_000 })
  })
})

describe('what the store keeps', () => {
  test('threads of the 20 sessions heard from last, no more', async ($, on) => {
    const old = Object.fromEntries(Array.from({ length: 25 }, (_, i) => [`threads:s${i}`, { turn: 1, threads: [] }]))
    const w = world(on, { board: false, store: old })
    await startSession($)
    const left = [...w.store.keys()].filter(k => k.startsWith('threads:'))
    expect(left).toHaveLength(20)
    expect(left[0]).toBe('threads:s5')
  })
})

describe('metrics and the pane', () => {
  test('a worker session writes one row per session in the lane-eval columns', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await $.prompt.submit({ text: '/jaira-role-lane MG677R critique', origin: { kind: 'composer' }, wait: false })
    w.usage.tokens = 30_000
    await answer($, 'ok', 2000)
    await $.session.end({ reason: 'other' } as never)
    const csv = w.files.get('/home/a/.jaira/metrics.csv') ?? ''
    expect(csv.split('\n')[0]).toBe('arm,ticket,run,lane,input_tokens,cache_write_tokens,cache_read_tokens,output_tokens,cost_eq,peak_context,duration_ms,tool_uses,loops,verdict')
    expect(csv.split('\n')[1]).toMatch(/^live,MG677R,abcdef12,critique,100,2000,500,20,2650,2620,\d+,0,,$/)
    const record = w.store.get('session:abcdef12')
    expect(record).toMatchObject({ ticket: 'MG677R', role: 'lane', lane: 'critique', pane: 'w1:p1', tokens: 30_000, ttl: 3_600_000 })
  })

  test('/cockpit ticket opens a pane with the ticket\'s reports and usage', async ($, on) => {
    const w = world(on)
    await startSession($)
    w.files.set('/home/a/.jaira/metrics.csv', 'arm,ticket,run,lane,input_tokens,cache_write_tokens,cache_read_tokens,output_tokens,cost_eq,peak_context,duration_ms,tool_uses,loops,verdict\nworker,MG677R,s1,worker,10,20,30,40,38,90000,600000,12,,\n')
    await $.session.receive({ origin: { kind: 'peer' }, text: 'Ticket MG677R: review -> signoff, review-verdict="approve"' })
    await $.command.run(typed('ticket mg677r'))
    const pane = await $.ui.mount({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-ticket', props: { title: 'Ticket MG677R', isFocused: false, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } } as never)
    expect(await pane.find({ type: 'Text', text: /signoff · approve mit Vorbehalt/ })).toBeDefined()
    expect(await pane.find({ type: 'Text', text: 'Ticket MG677R: review -> signoff, review-verdict="approve"' })).toBeDefined()
    expect(await pane.find({ type: 'Text', text: 'worker · cost_eq 38 · peak 90k · 10 min · 12 tools' })).toBeDefined()
  })
})
