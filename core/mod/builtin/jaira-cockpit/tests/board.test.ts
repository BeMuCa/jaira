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
  // The ticket bodies `jaira show --json` gives, by id.
  bodies: Record<string, string>
  // The panes the plugin opened, and closed, by id and title.
  opened: { id: string; title?: string; focus?: boolean }[]
  closed: string[]
  clock: ReturnType<typeof mock.clock>
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
  on('ui.toast', () => ({ value: undefined }))
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
    if (cmd === 'jaira' && args[0] === 'show') return ok((w.warns ? `${OLDER}\n` : '') + JSON.stringify({ handle: args[1], body: w.bodies[args[1]!] ?? '' }))
    if (cmd === 'git' && args[0] === 'rev-parse') return ok('/repo/.git\n')
    if (cmd === 'git' && args[0] === 'worktree') return ok(WORKTREES)
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
    expect(w.opened).toEqual([])
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
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.slice(0, 6)).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress', '--force'])
  })

  test('notes why and moves the ticket back to in-progress', async ($, on) => {
    const w = world(on)
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'note', 'MG677R')[0]?.[3]).toBe('Given back from signoff: der Branch fehlt im Schluessel')
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress', '--force', '@/.worktrees/repo-mg677r'])
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
