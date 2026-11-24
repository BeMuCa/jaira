import type { On, SessionRateLimit } from 'claude-code'
import { describe, expect, mock, test } from 'claude-code/testing'
import { acceptTreeFor, carrierFor, ghRepoOf, parseRemote, classOf, clipLines, commitNaming, commitsOf, executionOf, followTitle, gapItems, notePaths, explanationNotes, type GapItem, languageOf, lightOf, reasonOf, reasonPrompt, releaseLines, rephrasedOf, seenToTick, snapshotOf, storyOf, translatePrompt, type Waiting } from '../hooks/board.ts'
import { reportLine } from '../hooks/digest.ts'
import { durationOf, stoneSize } from '../hooks/overview.ts'
import { answer as answerThread, answerLine, asksOf, file, parseFiling, titled, untitled, type Filing } from '../hooks/threads.ts'

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
  // What each `herdr pane get` answers as the pane's agent_status, in order;
  // once used up, the pane is idle.
  paneStates: string[]
  // When set, every `jaira move` waits here, in order, until the test lets it go.
  held: (() => void)[] | undefined
  // A forced move answers without moving the ticket: the row stays in the
  // band, for a test about the row and not about where the ticket went.
  stay?: boolean
  // This session's subagents, as $.agent.list() gives them.
  agents: { id: string; description: string; type: string; status: string; name?: string }[]
  // `jaira lanes --json` of a build before human_exit.
  oldJaira: boolean
  // The repository's boards besides today's, as `jaira boards --json` names
  // them; one without lanes is a board whose `lanes --board` read fails.
  boards: Record<string, { lanes?: unknown[]; tickets: unknown[] }>
  asked: string[]
  answer: string
  model: string
  modelCalls: number
  // The prompts the cheap model was handed, in order.
  prompts: string[]
  // The cheap model gives no answer: offline, or the call failed.
  modelDown?: boolean
  // Translate answers each text it is handed from this, a text not in it as
  // it is; undefined: it answers model.
  translations?: Record<string, string>
  // The Explain-on-one-note calls, and what each answers; undefined: no
  // answer.
  reasonPrompts: string[]
  reason?: string
  // What `jaira show --json` adds for a ticket beyond its body and lane, and
  // what `git log --format=%h --name-only --grep=<id>` prints.
  shows: Record<string, Record<string, unknown>>
  idLog: string
  // What \`git log -p\` prints; what a move to in-progress fails with.
  patches?: string
  moveFails?: string
  // When set, every cheap-model call waits here, in order, answering with
  // what w.model held when it was made.
  modelHeld?: (() => void)[]
  filled: string[]
  store: Map<string, unknown>
  status: (string | undefined)[]
  // What a dry run of a move refuses, by code: a person's lane always refuses
  // needs_human, the exit Accept stands for.
  refusals: { code: string; message: string }[]
  // A jaira older than the board: it warns on a line before what it prints.
  warns: boolean
  // What `jaira --version` names: before 0.3.7 a move back out of a
  // person's lane without --force is refused.
  version: string
  // `jaira note` with a give-back reason fails after the move went through.
  noteFails: boolean
  toasts: string[]
  // The definition-of-done item whose `jaira dod` fails (DKMG91).
  dodFails?: string
  // What `jaira whoami --json` says; undefined: it says nothing.
  me?: { identity: string; aliases: string[] }
  // The ticket bodies `jaira show --json` gives, by id.
  bodies: Record<string, string>
  // The bughunt key `jaira show --brief --json` gives, by id; an id without
  // one gets a brief without the key, as from a jaira older than it.
  hunts: Record<string, { at: string; state: string } | null>
  // The commits `git merge-base --is-ancestor <commit> HEAD` finds in the
  // checkout it runs in, by directory ('@/...'); any other answers 1.
  contains: Record<string, string[]>
  // What the plugin submitted to this session's conversation itself.
  submitted: string[]
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
  branchCommits: { sha: string; subject: string; body?: string }[]
  holders: string[]
  toasts: string[]
  // `git remote -v`.
  remotes: string
  // .jaira/pr as the base holds it (`git show <base>:.jaira/pr`); undefined:
  // no such file there.
  basePr?: string
  // origin's default branch where no forge names it (`git rev-parse
  // --abbrev-ref origin/HEAD`); undefined: origin/HEAD not set here.
  originHead?: string
  // The forge answers `repo view` but refuses `pr list`.
  prListFails: boolean
  nextFails: boolean
  // The latest release `jaira self upgrade --check` finds; the check writes
  // it into the release cache, as jaira does. `jaira self upgrade` installs
  // it, unless upgradeFails.
  latest: string
  upgradeFails: boolean
  // How `jaira update --roles --lanes --dry-run` exits: 0, 3 (an edited
  // file it would leave alone) or 2 (a jaira without --roles).
  dryRunExit: number
  // The terminal has no room for a pane: $.ui.open places nothing.
  noRoom?: boolean
  // When set, every `herdr tab list` waits here until the test lets it go.
  heldTabs?: (() => void)[]
  // Every `herdr tab list` fails with this on stderr (W80HNH).
  tabsFail?: string
  // Every `herdr tab list` prints the tabs but exits 3 with this on stderr,
  // as herdr-hc.sh does when a terminal gave no answer (8NWE4P).
  tabsIncomplete?: string
  // spawn.sh calls wait here until released, as one waiting for claude does.
  heldSpawn?: (() => void)[]
  // `git worktree list --porcelain`, WORKTREES unless set.
  worktrees?: string
  // `jaira create` of a follow-up fails.
  createFails?: boolean
  // Another session's move came first: jaira answers this one 'already in'.
  movedFirst?: boolean
  // .jaira/pr-when as `jaira roles list --json` reports it; done unless set.
  prWhen?: string
  // What each `jaira pr <branch> --json` answers, in turn, the last one
  // repeating; '!' + text fails with exit 3, '?' + text with exit 1.
  // { prs: [] } when empty.
  prReads: string[]
  // `jaira roles list --json` fails.
  rolesFails?: boolean
  // Where MG677R's file is not as its ref has it, by directory ('@/...'):
  // none there ('ref-only') or an older one ('behind'), until `jaira pull`
  // there brings it. Every write there is refused, as jaira refuses it.
  disk: Record<string, 'ref-only' | 'behind'>
  // The ticket `disk` speaks of; MG677R where unset.
  diskId?: string
  // What spawn.sh fails with; undefined: it opens its tab.
  spawnFails?: string
  // What `jaira pull` refuses with; undefined: it pulls.
  pullFails?: string
  // What `git clean` refuses with; undefined: it removes the file.
  cleanFails?: string
  // Parents spawn.sh refuses as no session listening there (exit 2).
  deadParents?: string[]
}

const LANES = {
  lanes: [
    { id: 'in-progress', agentic: true, context: '' },
    { id: 'critique', agentic: true, context: '' },
    { id: 'review', agentic: true, context: 'fresh' },
    { id: 'signoff', agentic: false, context: '', human_exit: true },
    { id: 'human', agentic: false, context: '', human_exit: false, requires_question: true },
    { id: 'done', agentic: false, context: '', terminal: true },
  ],
}

const SIGNOFF = { handle: 'MG677R', title: 'cache key', status: 'signoff', next_lane: 'done', question: '', review: { verdict: 'approve mit Vorbehalt: der Schluessel ignoriert den Branch' } }
const HUMAN = { handle: 'K7P2QX', title: 'tags popup', status: 'human', next_lane: 'review', question: 'Popup oder eigene Seite?', review: { verdict: '' } }
const WORKING = { handle: 'T8QS1Y', title: 'filter', status: 'in-progress', question: '', review: {} }

const OLDER = "jaira: this board was set up by an older version of jaira — run 'jaira update' to bring it up to date"
const EXPLANATION = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-explanation', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never
const ANSWER = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-answer', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never
const MESSAGE = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-message', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never

// Herdr's review worktree of feat/mg677r is listed first and ends in the id
// as well, without the ticket on its disk: every press must pass it by.
const WORKTREES = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n\nworktree /.herdr/worktrees/repo/review-feat-mg677r\nHEAD 2\nbranch refs/heads/review/feat/mg677r\n\nworktree /.worktrees/repo-mg677r\nHEAD 2\nbranch refs/heads/feat/mg677r\n'
const SPAWN = '/repo/.claude/skills/jaira-dispatcher/scripts/spawn.sh'
// The release cache jaira writes under ~/.jaira, HOME being /home/a.
const CACHE = '/home/a/.jaira/update-check.json'

// This session's messaging socket, as Herdr tests set it; socket: false in
// world() is a session whose socket file is gone.
const SOCKET = '/tmp/cc-socks-1000/77.sock'

function world(on: On, opts: { board?: boolean; herdr?: boolean; store?: Record<string, unknown>; spawn?: boolean; socket?: boolean; env?: Record<string, string> } = {}): World {
  const w: World = {
    ran: [],
    files: new Map(),
    usage: { tokens: 0, rateLimits: [{ kind: 'five_hour', percentUsed: 34 }] },
    tickets: [SIGNOFF, HUMAN, WORKING],
    tabs: [],
    paneGone: false,
    paneStates: [],
    held: undefined,
    agents: [],
    oldJaira: false,
    boards: {},
    asked: [],
    answer: '',
    model: '',
    modelCalls: 0,
    prompts: [],
    reasonPrompts: [],
    shows: {},
    idLog: '',
    filled: [],
    store: new Map(Object.entries(opts.store ?? {})),
    status: [],
    refusals: [{ code: 'needs_human', message: 'lane "signoff" is a human checkpoint' }],
    warns: false,
    version: '0.3.7',
    noteFails: false,
    toasts: [],
    disk: {},
    bodies: {},
    hunts: {},
    contains: {},
    submitted: [],
    prReads: [],
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
    nextFails: false,
    latest: '0.3.7',
    upgradeFails: false,
    dryRunExit: 0,
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
  mock.env(on, opts.env ?? { HOME: '/home/a', ...(opts.herdr ? { HERDR_ENV: '1', HERDR_BIN_PATH: 'herdr', HERDR_PANE_ID: 'w1:p1', CLAUDE_CODE_MESSAGING_SOCKET: '/tmp/cc-socks-1000/77.sock' } : {}) })
  on('session.start', ($, e) => ({ cwd: e.cwd }))
  on('command.register', () => ({ value: undefined }) as never)
  on('tool.register', () => ({ value: undefined }) as never)
  on('ui.render', () => ({ type: 'Box', props: {}, children: [] }) as never)
  on('ui.open', ($, e) => {
    const { id, title, focus } = e as unknown as { id: string; title?: string; focus?: boolean }
    w.opened.push({ id, ...(title !== undefined ? { title } : {}), ...(focus ? { focus } : {}) })
    return { value: { isPlaced: !w.noRoom } }
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
  on('prompt.submit', ($, e) => {
    if (e.origin.kind === 'plugin') w.submitted.push(e.text)
    return { text: e.text, origin: e.origin }
  })
  on('turn.complete', ($, e) => ({ text: e.answer, ...(e.usage ? { usage: e.usage } : {}) }))
  on('session.cwd', () => ({ value: '/repo' }))
  on('session.repo', () => ({ value: { root: '/repo', remote: null, internal: false, name: null } }))
  on('session.id', () => ({ value: 'abcdef123456' }))
  on('session.usage', () => ({ value: { startedAt: 0, context: { tokens: w.usage.tokens, window: 1_000_000 }, rateLimits: w.usage.rateLimits } }))
  on('fs.exists', ($, e) => ({ value: (hasBoard && e.path === '/repo/.jaira') || (e.path === SPAWN && opts.spawn !== false) || (e.path === SOCKET && opts.socket !== false) || w.files.has(e.path) }))
  on('fs.read', ($, e) => (w.files.has(e.path) ? { value: w.files.get(e.path)! } : { deny: 'no such file' }))
  on('fs.write', ($, e) => {
    w.files.set(e.path, e.text)
    return { value: undefined }
  })
  on('process.run', async ($, e) => {
    w.ran.push([...e.argv, `@${e.init?.cwd ?? ''}`])
    if (w.held && e.argv[0] === 'jaira' && e.argv[1] === 'move') await new Promise<void>(go => w.held!.push(go))
    if (w.heldTabs && e.argv[1] === 'tab' && e.argv[2] === 'list') await new Promise<void>(go => w.heldTabs!.push(go))
    if (w.heldSpawn && e.argv[0] === 'bash' && e.argv[1] === SPAWN) await new Promise<void>(go => w.heldSpawn!.push(go))
    const ok = (stdout = '') => ({ value: { exitCode: 0, stdout, stderr: '', isStdoutTruncated: false, isStderrTruncated: false } })
    const [cmd, ...args] = e.argv
    const here = `@${e.init?.cwd ?? ''}`
    const refuse = (stderr: string) => ({ value: { exitCode: 1, stdout: '', stderr, isStdoutTruncated: false, isStderrTruncated: false } })
    if (cmd === 'bash' && args[0] === SPAWN && w.spawnFails !== undefined) return refuse(w.spawnFails)
    if (cmd === 'jaira' && args[0] === 'pull') {
      if (w.pullFails !== undefined) return refuse(w.pullFails)
      delete w.disk[here]
      return ok()
    }
    if (cmd === 'git' && args[0] === 'clean' && w.cleanFails !== undefined) return refuse(w.cleanFails)
    if (args[0] === SPAWN && args.includes('--parent') && w.deadParents?.includes(args[args.indexOf('--parent') + 1]!)) {
      const p = args[args.indexOf('--parent') + 1]
      return { value: { exitCode: 2, stdout: '', stderr: `--parent ${p}: no session listens there. Leave --parent out`, isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (args[0] === SPAWN) return ok('w1:p42\n')
    const diskId = w.diskId ?? 'MG677R'
    if (cmd === 'jaira' && ['move', 'note', 'set', 'dod'].includes(args[0]!) && args[1] === diskId && w.disk[here] === 'ref-only') return refuse(`jaira: ${diskId}: ticket: on its ref and not on your disk; pull it first`)
    if (cmd === 'jaira' && ['move', 'note', 'set', 'dod'].includes(args[0]!) && args[1] === diskId && w.disk[here] === 'behind') return refuse(`jaira: ${diskId}: the file here is older than its ref (stale_file); jaira pull ${diskId}`)
    const on = args.includes('--board') ? w.boards[args[args.indexOf('--board') + 1] ?? ''] : undefined
    if (cmd === 'jaira' && args[0] === 'roles' && args[1] === 'list') {
      return w.rolesFails ? { value: { exitCode: 1, stdout: '', stderr: 'jaira: base unreadable', isStdoutTruncated: false, isStderrTruncated: false } } : ok(JSON.stringify({ pr_when: w.prWhen ?? 'done' }))
    }
    if (cmd === 'jaira' && args[0] === 'pr') {
      const a = (w.prReads.length > 1 ? w.prReads.shift() : w.prReads[0]) ?? JSON.stringify({ prs: [] })
      const code = a.startsWith('!') ? 3 : a.startsWith('?') ? 1 : 0
      return code ? { value: { exitCode: code, stdout: '', stderr: a.slice(1), isStdoutTruncated: false, isStderrTruncated: false } } : ok(a)
    }
    if (cmd === 'jaira' && args[0] === 'boards') return ok(JSON.stringify({ boards: ['default', ...Object.keys(w.boards)].map(name => ({ name })) }))
    if (cmd === 'jaira' && args.includes('--board') && (on === undefined || (args[0] === 'lanes' && on.lanes === undefined))) return failed('jaira: unknown flag: --board')
    if (cmd === 'jaira' && args[0] === 'lanes' && on !== undefined) return ok(JSON.stringify({ lanes: on.lanes }))
    if (cmd === 'jaira' && args[0] === 'list' && on !== undefined) return ok(JSON.stringify({ count: on.tickets.length, tickets: on.tickets }))
    if (cmd === 'jaira' && args[0] === 'lanes') return ok(JSON.stringify(w.oldJaira ? { lanes: LANES.lanes.map(({ human_exit, requires_question, ...l }) => l) } : LANES))
    if (cmd === 'jaira' && args[0] === 'list') return ok(JSON.stringify({ count: w.tickets.length, tickets: w.tickets }))
    if (cmd === 'jaira' && args[0] === 'move' && w.movedFirst && !args.includes('--dry-run')) return ok(`${args[1]} is already in ${args[3]}\n`)
    if (cmd === 'jaira' && args[0] === 'move' && args[3] === 'in-progress' && w.moveFails !== undefined) return { value: { exitCode: 3, stdout: '', stderr: w.moveFails, isStdoutTruncated: false, isStderrTruncated: false } }
    if (cmd === 'jaira' && args[0] === 'move' && args.includes('--dry-run')) {
      const error = (w.warns ? `${OLDER}\n` : '') + JSON.stringify({ error: { code: 3, reason: 'gate_refused', violations: w.refusals } })
      return { value: { exitCode: 3, stdout: '', stderr: error, isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === '--version') return ok(`jaira version ${w.version}\n`)
    if (cmd === 'jaira' && args[0] === 'self' && args.includes('--check')) {
      w.files.set(CACHE, JSON.stringify({ checked_at: '2026-10-16T10:00:00Z', latest: w.latest }))
      return ok(JSON.stringify({ current: w.version, latest: w.latest, up_to_date: w.latest === w.version, upgraded: false, target: '/home/a/.local/bin/jaira', install_method: 'self-managed', checked_at: '2026-10-16T10:00:00Z' }))
    }
    if (cmd === 'jaira' && args[0] === 'self' && w.upgradeFails) {
      return { value: { exitCode: 3, stdout: '', stderr: 'jaira: this jaira was installed with go install; upgrade it with go install github.com/BeMuCa/jaira/cmd/jaira@latest', isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'self') {
      const from = w.version
      w.version = w.latest
      // jaira names a release host override on stderr, after the upgrade.
      const stdout = `upgraded jaira ${from} -> ${w.latest} (/home/a/.local/bin/jaira)\nbringing your setup in step with ${w.latest}: jaira update --roles\n`
      return { value: { exitCode: 0, stdout, stderr: 'warning: release host overridden by JAIRA_RELEASE_API=http://127.0.0.1:8765/api\n', isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'update' && args.includes('--dry-run') && w.dryRunExit === 2) {
      return { value: { exitCode: 2, stdout: '', stderr: 'jaira: unknown flag: --roles\n', isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'update' && args.includes('--dry-run') && w.dryRunExit === 3) {
      const stdout = 'skipped /repo/.claude/skills/jaira-worker/SKILL.md (edited)\nnothing was written (--dry-run)\n'
      return { value: { exitCode: 3, stdout, stderr: "jaira: some files were edited here and left alone; 'jaira roles install --project --force' (or '--global --force') and 'jaira lanes use <id> --force' replace them\n", isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'update' && args.includes('--dry-run')) return ok(`${OLDER}\nskills in /repo/.claude/skills: 3 would be updated\nlanes: critique would be updated\nnothing was written (--dry-run)\n`)
    if (cmd === 'jaira' && args[0] === 'update') return ok('skills in /repo/.claude/skills: 3 updated\nlanes: critique updated\n')
    if (cmd === 'jaira' && args[0] === 'move' && !args.includes('--force') && /^0\.3\.[0-6]$/.test(w.version)) {
      return { value: { exitCode: 3, stdout: '', stderr: 'jaira: lane "signoff" is a human checkpoint; use --force', isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'note' && w.noteFails && args[2]?.startsWith('Given back')) {
      return { value: { exitCode: 1, stdout: '', stderr: 'jaira: ticket file is locked', isStdoutTruncated: false, isStderrTruncated: false } }
    }
    if (cmd === 'jaira' && args[0] === 'dod' && args[2] === w.dodFails) return { value: { exitCode: 3, stdout: '', stderr: 'jaira: the board is locked', isStdoutTruncated: false, isStderrTruncated: false } }
    if (cmd === 'jaira' && args[0] === 'whoami') return ok(w.me === undefined ? '' : JSON.stringify(w.me))
    if (cmd === 'jaira' && args[0] === 'show' && args.includes('--brief')) {
      const id = args[1]!
      return ok((w.warns ? `${OLDER}\n` : '') + JSON.stringify({ id, ...(id in w.hunts ? { bughunt: w.hunts[id] } : {}) }))
    }
    if (cmd === 'jaira' && args[0] === 'show') {
      const status = (w.tickets.find(x => (x as { handle: string }).handle === args[1]) as { status?: string } | undefined)?.status
      const disk = args[1] === (w.diskId ?? 'MG677R') ? w.disk[here] : undefined
      const file = disk === 'ref-only' ? {} : { path: `${e.init?.cwd ?? ''}/.jaira/tickets/01${args[1]}-cache-key.md` }
      return ok((w.warns ? `${OLDER}\n` : '') + JSON.stringify({ ...w.shows[args[1]!], handle: args[1], body: w.bodies[args[1]!] ?? '', ...(status !== undefined ? { status, board: '' } : {}), 'on-ref-only': disk === 'ref-only', 'behind-ref': disk === 'behind', ...file }))
    }
    if (cmd === 'git' && args[0] === 'rev-parse' && args[1] === '--abbrev-ref') return w.originHead === undefined ? { value: { exitCode: 1, stdout: '', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } } : ok(`${w.originHead}\n`)
    if (cmd === 'git' && args[0] === 'rev-parse') return ok('/repo/.git\n')
    if (cmd === 'git' && args[0] === '--exec-path') return ok('C:\\Program Files\\Git\\mingw64\\libexec\\git-core\n')
    if (cmd === 'git' && args[0] === 'worktree') return ok(w.worktrees ?? WORKTREES)
    if (cmd === 'git' && args[0] === 'merge-base' && args[1] === '--is-ancestor') {
      return (w.contains[`@${e.init?.cwd ?? ''}`] ?? []).includes(args[2]!) ? ok('') : { value: { exitCode: 1, stdout: '', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } }
    }
    const failed = (stderr: string) => ({ value: { exitCode: 1, stdout: '', stderr, isStdoutTruncated: false, isStderrTruncated: false } })
    // A follow-up gets the next handle, F00001 and on; createFails refuses it.
    if (cmd === 'jaira' && args[0] === 'create') {
      if (w.createFails) return failed('jaira: the board is locked')
      return ok(JSON.stringify({ handle: `F${String(w.ran.filter(r => r[0] === 'jaira' && r[1] === 'create').length).padStart(5, '0')}` }))
    }
    if (cmd === 'git' && args[0] === 'config') return w.forge ? ok(`${w.forge}\n`) : failed('')
    if (cmd === 'git' && args[0] === 'remote' && args[1] === 'get-url') return w.origin === null ? failed("error: No such remote 'origin'") : ok(`${w.origin}\n`)
    if ((cmd === 'gh' || cmd === 'glab') && w.prs === null) return failed(`${cmd}: command not found`)
    if ((cmd === 'gh' || cmd === 'glab') && args[1] === 'list' && w.prListFails) return failed("HTTP 401: Bad credentials")
    // origin is a fork of upstream/repo: pull requests are asked of the parent.
    // Its default branch is main; the fork's own lags behind as master.
    // On GitHub Enterprise gh is asked host/owner/repo (XSGWPX).
    if (cmd === 'gh' && args[0] === 'repo') {
      return ok(JSON.stringify(args[2]?.replace(/^ghe\.example\.com\//, '') === 'upstream/repo'
        ? { isFork: false, nameWithOwner: 'upstream/repo', parent: null, defaultBranchRef: { name: 'main' } }
        : { isFork: true, nameWithOwner: 'acme/repo', parent: { name: 'repo', owner: { login: 'upstream' } }, defaultBranchRef: { name: 'master' } }))
    }
    if (cmd === 'glab' && args[0] === 'repo') {
      return ok(JSON.stringify(args[2] === 'upstream/repo'
        ? { path_with_namespace: 'upstream/repo', default_branch: 'main' }
        : { path_with_namespace: 'acme/repo', forked_from_project: { path_with_namespace: 'upstream/repo' }, default_branch: 'master' }))
    }
    if (cmd === 'git' && args[0] === 'remote' && args[1] === '-v') return ok(w.remotes)
    if (cmd === 'git' && args[0] === 'show' && args[1]?.endsWith(':.jaira/pr')) return w.basePr === undefined ? failed(`fatal: path '.jaira/pr' does not exist in '${args[1].split(':')[0]}'`) : ok(w.basePr)
    if (cmd === 'git' && args[0] === 'branch' && args[1] === '--contains') return ok(w.holders.join('\n'))
    if (cmd === 'git' && args[0] === 'branch') return ok(w.branches.join('\n'))
    // Newest first, as git prints them. A --grep reads the whole message,
    // body included; --format=%H %s prints the subject alone.
    if (cmd === 'git' && args[0] === 'log' && args[1] === '--format=%h') return ok(w.idLog)
    if (cmd === 'git' && args[0] === 'log' && args[1] === '-p') return ok(w.patches ?? '')
    if (cmd === 'git' && args[0] === 'log' && args.includes('--branches')) {
      const grep = args.find(a => a.startsWith('--grep='))
      if (grep) {
        const pattern = new RegExp(grep.slice('--grep='.length), 'm')
        return ok(w.branchCommits.find(c => pattern.test(`${c.subject}\n\n${c.body ?? ''}`))?.sha ?? '')
      }
      return ok(w.branchCommits.map(c => `${c.sha} ${c.subject}`).join('\n'))
    }
    // As many subjects as commits ahead: the first `ahead` of them.
    if (cmd === 'git' && args[0] === 'log') return ok(w.subjects.concat(Array(Math.max(0, w.ahead - w.subjects.length)).fill('chore: more')).slice(0, w.ahead).join('\n'))
    // The move a press forces lands: jaira list reads the ticket in done after it.
    if (cmd === 'jaira' && args[0] === 'move' && args.includes('--force')) {
      if (w.stay) return ok()
      w.tickets = w.tickets.map(x => ((x as { handle: string }).handle === args[1] ? { ...(x as object), status: args[3] } : x))
      return ok()
    }
    if (cmd === 'gh' || cmd === 'glab') return ok(JSON.stringify(w.prs))
    if (cmd === 'jaira' && args[0] === 'next' && w.nextFails) return failed('jaira: board unreadable')
    if (cmd === 'jaira' && args[0] === 'next') return ok(JSON.stringify({ count: w.next.length, tickets: w.next.map(handle => ({ handle })) }))
    if (cmd === 'herdr' && args[0] === 'tab' && args[1] === 'list' && w.tabsFail !== undefined) return failed(w.tabsFail)
    if (cmd === 'herdr' && args[0] === 'tab' && args[1] === 'list' && w.tabsIncomplete !== undefined) return { value: { exitCode: 3, stdout: JSON.stringify({ result: { tabs: w.tabs } }), stderr: w.tabsIncomplete, isStdoutTruncated: false, isStderrTruncated: false } }
    if (cmd === 'herdr' && args[0] === 'tab' && args[1] === 'list') return ok(JSON.stringify({ result: { tabs: w.tabs } }))
    if (cmd === 'herdr' && args[0] === 'pane' && args[1] === 'get' && !w.paneGone) return ok(JSON.stringify({ result: { pane: { agent: 'claude', agent_status: w.paneStates.shift() ?? 'idle' } } }))
    if (cmd === 'herdr' && args[0] === 'pane' && w.paneGone) return { value: { exitCode: 1, stdout: '', stderr: 'no such pane', isStdoutTruncated: false, isStderrTruncated: false } }
    return ok()
  })
  on('tool.call', { tool: 'AskUserQuestion' }, ($, e) => {
    w.asked.push(JSON.stringify(e))
    const { questions } = e as { questions: { question: string }[] }
    return { result: { questions, answers: { [questions[0]!.question]: w.answer } }, text: w.answer }
  })
  on('agent.list', () => ({ value: w.agents }) as never)
  on('model.complete', async ($, e) => {
    const prompt = (e as { prompt: string }).prompt
    // Explain on one note is a call of its own, answered from its own field
    // and never held: the counts and the holds below are Translate's.
    if (prompt.startsWith('A person decides whether to accept')) {
      w.reasonPrompts.push(prompt)
      return (w.reason === undefined ? { value: { isAnswered: false } } : { value: { isAnswered: true, text: w.reason, usage: { input_tokens: 1, output_tokens: 1, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 } } }) as never
    }
    w.modelCalls++
    w.prompts.push(prompt)
    if (w.modelDown) return { value: { isAnswered: false } } as never
    const text = w.translations === undefined ? w.model : JSON.stringify((JSON.parse(prompt.slice(prompt.lastIndexOf('\n') + 1)) as string[]).map(x => w.translations![x] ?? x))
    if (w.modelHeld) await new Promise<void>(go => w.modelHeld!.push(go))
    return { value: { isAnswered: true, text, usage: { input_tokens: 1, output_tokens: 1, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 } } }
  })
  return w
}

// A Log line's start: its kind's sign, the clock and two spaces; `$1` keeps the sign.
const LOG_LINE = /^(\p{Emoji_Presentation}) \d\d:\d\d {2}/u
const BAND = { plugin: 'jaira-cockpit', surface: 'terminal', component: 'AbovePrompt', props: { hasSurvey: false, isWorking: false, maxRows: 20, bodyColumns: 120, scroll: { offset: 0, bodyRows: 20 }, view: {} } } as const

// The engine's `$` a test body is handed.
type Test$ = Parameters<Extract<Parameters<typeof test>[1], (...args: never[]) => unknown>>[0]

const typed = (args = '') =>
  ({ command: 'cockpit', args, origin: { kind: 'composer' }, presentation: { isFullscreen: false, columns: 120 } }) as const

// What the Acceptance pane says of a ticket, opened from the band: what the row leaves out.
async function verdictOf($: Test$, band: Pressable, id: string) {
  const pane = await acceptancePane($, band, id)
  const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
  await pane.unmount()
  return { texts, find: async (q: { type: 'Text'; text: string | RegExp }) => texts.find(x => (typeof q.text === 'string' ? x === q.text : q.text.test(x))) }
}

// 2V50Y5: the Acceptance pane replaced Accept, Give back, Explain and Verdict
// in the row: its button opens the pane, whose Execute and Explain all do
// what the row's buttons did.
const ACCEPT = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-acceptance', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 40 }, view: {} } }) as const
type Pressable = { press: (q: { key: string }) => Promise<unknown> }
// One pane mounted at a time, as the engine allows: a second press remounts it.
let mountedPane: { unmount: () => Promise<unknown> } | undefined
async function acceptancePane($: Test$, band: Pressable, id: string) {
  await band.press({ key: `acceptance-${id}` })
  await mountedPane?.unmount().catch(() => undefined)
  const pane = await $.ui.mount(ACCEPT(`Acceptance of ${id}`))
  mountedPane = pane
  return pane
}
// Explain sits in the Acceptance pane now; Discuss stays in the row.
async function pressChat($: Test$, band: Pressable, kind: string, id: string, w?: World) {
  if (kind === 'explain') return pressExplain($, band, id, w)
  return band.press({ key: `${kind}-${id}` })
}
// What opening the pane ran and opened - the history's git log, the pull
// request look - is no part of the press: dropped, so a test reads what the
// press itself did.
const opening = (r: string[]) =>
  r[0] === 'gh' || r[0] === 'glab' || (r[0] === 'jaira' && r[1] === 'show' && r[3] === '--json') ||
  (r[0] === 'git' && (['worktree', 'config', 'remote', 'fetch', 'rev-parse'].includes(r[1]!) || (r[1] === 'log' && r[2] === '--format=%h')))
async function quietly($: Test$, band: Pressable, id: string, w?: World) {
  const ran = w?.ran.length ?? 0
  const opened = w?.opened.length ?? 0
  const pane = await acceptancePane($, band, id)
  // The pane is up once its opening read the ticket: a person sees it, and
  // presses in it, only then.
  await w?.clock.settle()
  if (w !== undefined) {
    const added = w.ran.slice(ran)
    w.ran.splice(ran, added.length, ...added.filter(r => !opening(r)))
    const shown = w.opened.slice(opened)
    w.opened.splice(opened, shown.length, ...shown.filter(o => o.id !== 'cockpit-acceptance'))
  }
  return pane
}
async function pressAccept($: Test$, band: Pressable, id: string, w?: World) {
  const pane = await quietly($, band, id, w)
  await pane.press({ key: `execute-${id}` })
  return pane
}
async function pressExplain($: Test$, band: Pressable, id: string, w?: World) {
  const pane = await quietly($, band, id, w)
  await pane.press({ key: `explain-all-${id}` })
  return pane
}
// One step through a note's choices, as the old row's press was: fix now,
// follow-up, accept, fix now again.
async function cycle(pane: Pressable & { find: (q: { key: string }) => Promise<{ text: string } | undefined> }, id: string, i: number) {
  const order = ['now', 'later', 'accept']
  for (const c of order) {
    if ((await pane.find({ key: `${c}-${id}-${i}` }))?.text.startsWith('[•')) return pane.press({ key: `${order[(order.indexOf(c) + 1) % 3]}-${id}-${i}` })
  }
  throw new Error(`no choice drawn for note ${i} of ${id}`)
}

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

// The band as one string, in the order it is drawn: its frame's text.
async function drawn($: Test$) {
  return (await (await $.ui.mount(BAND)).find({ type: 'Box' }))!.text
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
    // Not the first Box: Claude Code 2.1.329 mounts the band inside one of its own.
    expect((await ui.findAll({ type: 'Box' })).some(x => x.props.borderStyle === 'round')).toBe(true)
    expect(await ui.find({ type: 'Text', text: 'jaira-cockpit' })).toBeDefined()
    expect((await ui.find({ text: 'Waiting for you' }))).toBeDefined()
    expect((await ui.find({ type: 'Text', text: /MG677R \| signoff \| DoD/ }))).toBeDefined()
    expect((await ui.find({ type: 'Text', text: /K7P2QX \| human/ }))).toBeDefined()
    // The verdict is a press on Verdict away, not in the row.
    expect(await ui.find({ text: /approve mit Vorbehalt/ })).toBeUndefined()
    expect(await ui.find({ text: /T8QS1Y/ })).toBeUndefined()
  })

  test('a ticket in the band is one row: its marks, no verdict, no notes, no question under it', async ($, on) => {
    const w = world(on)
    const verdict = 'approve, but the cache key ignores the branch'
    w.tickets = [
      { ...SIGNOFF, goal: 'Two branches never share a cache entry.', dod_complete: true, round: 3, review_round: 2, question: 'keep the old key?', review: { verdict, gaps: 'removed a helper', verdict_gaps: 'blocker: the key ignores the branch\nnote: a test is missing\nnote: a comment is stale' } },
      { ...SIGNOFF, handle: 'AB12CD', dod_complete: false, review: { verdict: 'approve', gaps: 'none' } },
      { ...SIGNOFF, handle: 'EF34GH', review: { verdict: 'approve' } },
    ]
    w.hunts = { MG677R: { at: 'abc1234', state: 'ok' }, AB12CD: null }
    await startSession($)
    // Read once the band is up: the hunts are read for the tickets it shows.
    await w.clock.advance(60_000)
    const band = await $.ui.mount(BAND)
    for (const id of ['MG677R', 'AB12CD', 'EF34GH']) {
      const row = await band.find({ key: `row-${id}` })
      expect(row?.children.length).toBe(1)
    }
    expect(await band.find({ text: /ignores the branch|a test is missing|keep the old key/ })).toBeUndefined()
    expect(await band.find({ type: 'Text', text: 'MG677R | signoff | DoD ✓ | bughunt ✓ | 1 blocker | 2 notes | review of round 2 | ? ' })).toBeDefined()
    // No hunt finished is a cross; a jaira that does not say, or a ticket
    // with no dod_complete, a '?'.
    expect(await band.find({ type: 'Text', text: 'AB12CD | signoff | DoD ✗ | bughunt ✗ | 0 blockers | 0 notes ' })).toBeDefined()
    expect(await band.find({ type: 'Text', text: 'EF34GH | signoff | DoD ? | bughunt ? | 0 blockers | 0 notes ' })).toBeDefined()
    // The hunt is read where the ticket's code is: its own worktree.
    expect(ran(w, 'jaira', 'show', 'MG677R', '--brief', '--json')[0]?.at(-1)).toBe('@/.worktrees/repo-mg677r')
    // The marks in the theme's colours.
    const row = await band.find({ key: 'row-MG677R' })
    const colours = JSON.stringify(row)
    expect(colours).toContain('"color":"success"')
    expect(colours).toContain('"color":"error"')
    // Everything left out of the row is in the Verdict pane, the goal on top.
    const { texts } = await verdictOf($, band, 'MG677R')
    // The light comes first, before the title: red for the blocker, the
    // other reasons counted.
    expect(texts[0]).toBe('● red: 1 blocker in the review (+4 more)')
    const title = texts.indexOf('cache key')
    expect(title).toBeGreaterThan(0)
    expect(texts.slice(title, title + 2)).toEqual(['cache key', '┃ Goal: Two branches never share a cache entry.'])
    for (const text of ['? keep the old key?', '🔴 blocker: the key ignores the branch', '🟡 note: a test is missing', '🟡 note: a comment is stale']) {
      expect(texts).toContain(text)
    }
    expect(texts.some(x => x.endsWith(`review: ${verdict}`))).toBe(true)
    expect(texts.some(x => /review of round 2, changed since/.test(x))).toBe(true)
  })

  test('a hunt\'s fail is a red cross only in a checkout that holds the hunt\'s commit', async ($, on) => {
    const w = world(on)
    w.tickets = [
      // In its own worktree, which holds the hunt: the fail stands.
      { ...SIGNOFF, review: { verdict: 'approve' } },
      // No worktree of its own: read in the main checkout, which does not.
      { ...SIGNOFF, handle: 'AB12CD', review: { verdict: 'approve' } },
    ]
    w.hunts = { MG677R: { at: 'abc1234', state: 'fail' }, AB12CD: { at: 'def5678', state: 'fail' } }
    w.contains = { '@/.worktrees/repo-mg677r': ['abc1234'] }
    await startSession($)
    await w.clock.advance(60_000)
    const band = await $.ui.mount(BAND)
    expect(await band.find({ type: 'Text', text: /^MG677R \| signoff \| DoD \? \| bughunt ✗/ })).toBeDefined()
    expect(await band.find({ type: 'Text', text: /^AB12CD \| signoff \| DoD \? \| bughunt \?/ })).toBeDefined()
    expect(ran(w, 'git', 'merge-base', '--is-ancestor', 'def5678', 'HEAD')[0]?.at(-1)).toBe('@/repo')
  })

  test('an earlier review whose notes were not kept counts nothing: unknown, not zero', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, round: 3, review_round: 2, review: { verdict: 'back to in-progress: gap 1', gaps: 'removed a helper' } }]
    await startSession($)
    const band = await $.ui.mount(BAND)
    expect(await band.find({ type: 'Text', text: /MG677R \| signoff \| .*\? blockers \| \? notes \| review of round 2/ })).toBeDefined()
    const { texts } = await verdictOf($, band, 'MG677R')
    expect(texts).toContain('The notes of that review were not kept: a later lane rewrote them.')
    expect(texts.some(x => /removed a helper/.test(x))).toBe(false)
  })

  test('the Explain pane starts with the ticket\'s title and goal', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, goal: 'Two branches never share a cache entry.' }]
    w.bodies.MG677R = '## Progress\n\n- **2026-10-12 15:00 · Alex** — explain: The key leaves out the branch.'
    await startSession($)
    const band = await $.ui.mount(BAND)
    await band.press({ key: 'read-MG677R' })
    const texts = (await (await $.ui.mount(EXPLANATION('Explanation of MG677R'))).findAll({ type: 'Text' })).map(x => x.text)
    expect(texts.slice(0, 2)).toEqual(['cache key', 'Goal: Two branches never share a cache entry.'])
    expect(texts.indexOf('Explanation of MG677R')).toBeGreaterThan(1)
  })

  test('tickets waiting on another board are listed too, each by its own board\'s lanes', async ($, on) => {
    const w = world(on)
    w.boards.video = {
      lanes: [
        { id: 'script', agentic: true, context: '' },
        { id: 'approve', agentic: false, context: '', human_exit: true },
        { id: 'ask', agentic: false, context: '', requires_question: true },
        { id: 'done', agentic: false, context: '', terminal: true },
      ],
      tickets: [
        { handle: 'S0EG0M', title: 'intro', status: 'approve', next_lane: 'done', question: '', review: {} },
        { handle: 'V1D3OQ', title: 'outro', status: 'ask', next_lane: 'script', question: 'Musik oder Stille?', review: {} },
        { handle: 'V2R3ND', title: 'credits', status: 'script', question: '', review: {} },
      ],
    }
    await startSession($)
    expect(await headLine($)).toBe('0 worker | 0 review | 4 waiting for you | 5h 34%')
    await (await $.ui.mount(BAND)).press({ key: 'more' })
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ type: 'Text', text: /S0EG0M \| approve/ })).toBeDefined()
    expect(await ui.find({ type: 'Text', text: /V1D3OQ \| ask/ })).toBeDefined()
    expect(await ui.find({ key: 'acceptance-S0EG0M' })).toBeDefined()
    expect(await ui.find({ key: 'acceptance-V1D3OQ' })).toBeDefined()
    expect(await ui.find({ text: /V2R3ND/ })).toBeUndefined()
    expect(ran(w, 'jaira', 'list', '--board', 'video').length).toBeGreaterThan(0)
  })

  test('a board whose lanes cannot be read is left out, not the band', async ($, on) => {
    const w = world(on)
    w.boards.video = { tickets: [{ handle: 'S0EG0M', title: 'intro', status: 'approve', next_lane: 'done', question: '', review: {} }] }
    await startSession($)
    expect(await headLine($)).toBe('0 worker | 0 review | 2 waiting for you | 5h 34%')
  })

  test('Accept moves the ticket to done in its own worktree', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'MG677R', w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toEqual([
      ['jaira', 'move', 'MG677R', '--to', 'done', '--dry-run', '--json', '@/.worktrees/repo-mg677r'],
      ['jaira', 'move', 'MG677R', '--to', 'done', '--force', '@/.worktrees/repo-mg677r'],
    ])
  })

  test('Accept finds the ticket\'s own worktree mid-rebase, detached, past the Herdr review worktree on a branch', async ($, on) => {
    const w = world(on)
    w.worktrees = WORKTREES.replace('branch refs/heads/feat/mg677r', 'detached')
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'MG677R', w)
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.at(-1)).toBe('@/.worktrees/repo-mg677r')
  })

  test('Acceptance is the one coloured button of the row; Read explanation stays plain, and Accept, Give back, Explain and Verdict are gone', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = '## Progress\n\n- **2026-10-12 15:00 · Alex** — explain: The key leaves out the branch.'
    await startSession($)
    const ui = await $.ui.mount(BAND)
    // The colour a Button's label is drawn in: that of the Text it holds.
    const colour = async (key: string) => {
      const button = await ui.find({ key })
      expect(button).toBeDefined()
      return button!.children.map(c => (c as { props?: { color?: string } }).props?.color).find(Boolean)
    }
    expect(await colour('acceptance-MG677R')).toBe('claude')
    expect(await colour('read-MG677R')).toBeUndefined()
    for (const gone of ['accept', 'back', 'explain', 'verdict']) expect(await ui.find({ key: `${gone}-MG677R` })).toBeUndefined()
  })

  test('on the terminal Acceptance draws its brackets in the label\'s colour; elsewhere the native button stays', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = '## Progress\n\n- **2026-10-12 15:00 · Alex** — explain: The key leaves out the branch.'
    await startSession($)
    for (const [surface, plain, text] of [['terminal', true, '[ Acceptance ]'], ['desktop', undefined, 'Acceptance']] as const) {
      const ui = await $.ui.mount({ ...BAND, surface })
      const button = await ui.find({ key: 'acceptance-MG677R' })
      expect(button?.props.plain).toBe(plain)
      expect(button?.props.label).toBe('Acceptance')
      expect(button?.children).toEqual([{ type: 'Text', props: { color: 'claude', dimColor: false }, children: [text] }])
      await ui.unmount()
    }
    const ui = await $.ui.mount(BAND)
    // Read explanation keeps the engine's own brackets.
    expect((await ui.find({ key: 'read-MG677R' }))?.props.plain).toBeUndefined()
  })

  // Lets the oldest `jaira move` still waiting answer, and what it set off run.
  const release = async (w: World) => {
    w.held!.shift()!()
    await w.clock.settle()
  }
  const dimmed = async (ui: { find: (q: { key: string }) => Promise<{ children: unknown[] } | undefined> }, key: string) =>
    ((await ui.find({ key }))?.children[0] as { props: { dimColor?: boolean } }).props.dimColor

  test('a pressed row locks, drawn dim, until its action is done; a second press does nothing', async ($, on) => {
    const w = world(on)
    w.held = []
    w.stay = true
    await startSession($)
    const ui = await $.ui.mount(BAND)
    const pane = await quietly($, ui, 'MG677R', w)
    void pane.press({ key: 'execute-MG677R' })
    await w.clock.settle()
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(true)
    // Another ticket's row is not this press's.
    expect(await dimmed(ui, 'acceptance-K7P2QX')).toBe(false)
    await ui.press({ key: 'acceptance-MG677R' })
    void pane.press({ key: 'execute-MG677R' })
    await w.clock.settle()
    expect(w.opened).toEqual([])
    await release(w)
    await release(w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(2)
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(false)
  })

  test('a row whose action has not ended is free again after 20 seconds, and the late end frees only its own lock', async ($, on) => {
    const w = world(on)
    w.held = []
    w.stay = true
    await startSession($)
    const ui = await $.ui.mount(BAND)
    const pane = await quietly($, ui, 'MG677R', w)
    void pane.press({ key: 'execute-MG677R' })
    await w.clock.advance(19_999)
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(true)
    await w.clock.advance(1)
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(false)
    void pane.press({ key: 'execute-MG677R' })
    await w.clock.settle()
    // The first press's dry run and move answer: it ends, the second waits on.
    await release(w)
    await release(w)
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(true)
    await release(w)
    await release(w)
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(false)
  })

  test('a reload frees a row whose press it cut off', async ($, on) => {
    const w = world(on)
    w.held = []
    await startSession($)
    const ui = await $.ui.mount(BAND)
    const pane = await quietly($, ui, 'MG677R', w)
    void pane.press({ key: 'execute-MG677R' })
    await w.clock.settle()
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(true)
    await startSession($)
    expect(await dimmed(ui, 'acceptance-MG677R')).toBe(false)
    // The cut-off press ends before the test does.
    await release(w)
    await release(w)
    await w.clock.advance(20_000)
  })

  test('Acceptance is offered on every waiting ticket; where the next lane does not end the board, Execute accepts nothing and says so', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'acceptance-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'discuss-K7P2QX' })).toBeDefined()
    const pane = await pressAccept($, ui, 'K7P2QX', w)
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain('K7P2QX waits in human, which does not end the work: nothing to accept there. Mark a note fix, or type what is missing, to give it back.')
  })

  // 2V50Y5 DoD 3: a gate that refuses done says so in the pane, changes
  // nothing, and the person's choices stay as they were.
  test('a gate refusing done is said in the Acceptance pane: nothing moves, nothing is written, the choices stay', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps: 'note: the size\nnote: pre-existing: the clock' } }]
    w.refusals = [...w.refusals, { code: 'needs_nonmodel_signal', message: 'the definition of done is not met: criterion 1 is still open' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    const pane = await acceptancePane($, ui, 'MG677R')
    // The clock: accept -> fix now -> follow-up.
    await cycle(pane, 'MG677R', 1)
    await cycle(pane, 'MG677R', 1)
    const opened = w.opened.length
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move').filter(r => r.includes('--force'))).toEqual([])
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(ran(w, 'jaira', 'create')).toEqual([])
    expect(w.opened.slice(opened)).toEqual([])
    expect(w.closed).not.toContain('cockpit-acceptance')
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    expect(texts).toContain('MG677R was not accepted, nothing was changed:\nthe definition of done is not met: criterion 1 is still open')
    expect((await pane.find({ key: 'later-MG677R-1' }))?.text).toBe('[•ticket]')
    // The next press, the gate satisfied, carries the same choices out.
    w.refusals = w.refusals.filter(r => r.code === 'needs_human')
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(ran(w, 'jaira', 'create')).toHaveLength(1)
  })

  test('Accept reads the dry run past the warning an older jaira prints before it', async ($, on) => {
    const w = world(on)
    w.warns = true
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'MG677R', w)
    expect(ran(w, 'jaira', 'move', 'MG677R').map(r => r[5])).toEqual(['--dry-run', '--force'])
    expect(w.opened.filter(o => o.title === 'Accept MG677R failed')).toEqual([])
  })

  test('a refusal is shown whole in the Acceptance pane, never cut to a line', async ($, on) => {
    const w = world(on)
    const long = `the definition of done is not met: ${'criterion still open, '.repeat(20)}and one more`
    w.refusals = [...w.refusals, { code: 'needs_nonmodel_signal', message: long }, { code: 'needs_commits', message: 'no commit names MG677R' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    const pane = await pressAccept($, ui, 'MG677R', w)
    expect(await pane.find({ type: 'Text', text: `MG677R was not accepted, nothing was changed:\n${long}\nno commit names MG677R` })).toBeDefined()
    expect(await ui.find({ type: 'Text', text: /criterion still open/ })).toBeUndefined()
  })

  test('Explain all in the Acceptance pane opens a Herdr tab of its own', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'explain-K7P2QX' })).toBeUndefined()
    await pressExplain($, ui, 'MG677R', w)
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'explain', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('an open question shows under its ticket with Discuss', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    const ui = await $.ui.mount(BAND)
    // The row marks the question; the question itself is in the Acceptance pane.
    expect(await ui.find({ text: /Popup oder eigene Seite/ })).toBeUndefined()
    expect(await ui.find({ type: 'Text', text: /K7P2QX \| human \| .* \| \?/ })).toBeDefined()
    expect(await (await acceptancePane($, ui, 'K7P2QX')).find({ type: 'Text', text: '? Popup oder eigene Seite?' })).toBeDefined()
    await ui.press({ key: 'discuss-K7P2QX' })
    // The chat sends its decision to this session, by its messaging socket.
    expect(ran(w, 'bash', SPAWN)[0]?.slice(0, 8)).toEqual(['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'k7p2qx', 'K7P2QX', 'discuss'])
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
    expect(await band.find({ type: 'Text', text: /K7P2QX \| decided/ })).toBeDefined()
    await band.press({ key: 'decision-K7P2QX' })
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: 'Decision on K7P2QX', focus: true })
    const pane = await $.ui.mount(MESSAGE('Decision on K7P2QX'))
    expect(await pane.find({ type: 'Text', text: 'a page of its own, back to in-progress' })).toBeDefined()
    expect(await (await $.ui.mount(BAND)).find({ key: 'decision-K7P2QX' })).toBeUndefined()
  })

  test('a person\'s session runs reap.sh from its first minute on, then every ten minutes', async ($, on) => {
    const REAP = '/repo/.claude/skills/jaira-dispatcher/scripts/reap.sh'
    const w = world(on, { herdr: true })
    await startSession($)
    await w.clock.settle()
    expect(ran(w, 'bash', REAP)).toEqual([])
    await w.clock.advance(60_000)
    expect(ran(w, 'bash', REAP)).toEqual([['bash', REAP, '/repo', '@/repo']])
    for (let i = 0; i < 9; i++) await w.clock.advance(60_000)
    expect(ran(w, 'bash', REAP)).toHaveLength(1)
    await w.clock.advance(60_000)
    expect(ran(w, 'bash', REAP)).toHaveLength(2)
  })

  test('a tab spawn.sh opened runs no reap.sh, though its role line comes only after the start', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await w.clock.advance(30_000)
    await $.prompt.submit({ text: '/jaira-worker MG677R --parent jaira-4c', origin: { kind: 'composer' }, wait: false })
    await w.clock.advance(11 * 60_000)
    expect(ran(w, 'bash', '/repo/.claude/skills/jaira-dispatcher/scripts/reap.sh')).toEqual([])
  })

  test('a session that carries a ticket runs no reap.sh', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'carrying:abcdef12': { ticket: 'MG677R', role: 'worker', lane: '' } } })
    await startSession($)
    await w.clock.advance(11 * 60_000)
    expect(ran(w, 'bash', '/repo/.claude/skills/jaira-dispatcher/scripts/reap.sh')).toEqual([])
  })

  test('without Herdr the cockpit runs no reap.sh: it would keep everything', async ($, on) => {
    const w = world(on)
    await startSession($)
    await w.clock.advance(60_000)
    expect(ran(w, 'bash', '/repo/.claude/skills/jaira-dispatcher/scripts/reap.sh')).toEqual([])
  })

  test('without Herdr, Explain asks this session for a subagent that writes the explain note', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressExplain($, ui, 'MG677R', w)
    expect(ran(w, 'bash')).toEqual([])
    expect(w.opened).toEqual([])
    expect(w.toasts).toContain('MG677R: a subagent of this session explains the finding; read it under Read explanation')
    expect(w.submitted).toHaveLength(1)
    const asked = w.submitted[0]!
    expect(asked).toContain('start one subagent (Agent tool)')
    // In the ticket's worktree, and in the shape the cockpit reads back.
    expect(asked).toContain("jaira -C '/.worktrees/repo-mg677r' note MG677R \"explain: ")
    expect(asked).toContain('Status: DoD ')
    expect(asked).toContain('Recommendation: ')
    // a ticket reviewed before jaira kept the copy has an empty verdict_gaps: its notes stay in review-gaps
    expect(asked).toContain('Where review.verdict_gaps is not empty and differs from review.gaps')
    // the hunt's note has scrolled out of the four newest notes by then: the brief says whether it covers the code
    expect(asked).toContain("bughunt on HEAD ends in the brief's bughunt.state word as it stands")
    // the format line and the status sentence name all four words, so unknown is never written as fail
    expect(asked).toContain('bughunt on HEAD ok|fail|unknown|skipped')
    expect(asked).not.toContain('each check ends in ok or fail')
    for (const cls of ['[blocker]', '[small edge]', '[known limit]', '[not this ticket]']) expect(asked).toContain(cls)
    expect(asked).toContain('[give back]')
    // Nothing reads as "this must be changed" by its shape alone.
    expect(asked).not.toContain('Warum:')
    expect(asked).toContain('can be accepted as it is')
  })

  test('without Herdr, Discuss talks the question through in this chat', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'discuss-K7P2QX' })
    expect(ran(w, 'bash')).toEqual([])
    expect(w.opened).toEqual([])
    expect(w.toasts).toContain('K7P2QX: the open question is talked through in this chat')
    expect(w.submitted).toHaveLength(1)
    const asked = w.submitted[0]!
    expect(asked).toContain('waits for a decision of mine')
    expect(asked).toContain("note K7P2QX \"decision: <decision>\"")
    expect(asked).toContain('set K7P2QX question=')
    expect(asked).not.toContain('subagent')
  })

  test('without Herdr, a double click on Explain or Discuss sends one request', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    // A person's double click: the second press a moment after the first
    // has long returned.
    await ui.press({ key: 'discuss-K7P2QX' })
    await w.clock.advance(300)
    await ui.press({ key: 'discuss-K7P2QX' })
    expect(w.submitted).toHaveLength(1)
    await pressExplain($, ui, 'MG677R', w)
    await w.clock.advance(300)
    await pressExplain($, ui, 'MG677R', w)
    expect(w.submitted).toHaveLength(2)
    // Pressed again after the lock, it is asked again.
    await w.clock.advance(20_000)
    await pressExplain($, ui, 'MG677R', w)
    expect(w.submitted).toHaveLength(3)
  })

  test('with Herdr but no spawn.sh, Explain says spawn.sh is missing', async ($, on) => {
    const w = world(on, { herdr: true, spawn: false })
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(w.submitted).toEqual([])
    const pane = await $.ui.mount(MESSAGE('Explain MG677R failed'))
    expect(await pane.find({ type: 'Text', text: /spawn\.sh.*jaira roles install/ })).toBeDefined()
  })

  test('on native Windows, spawn.sh is found under USERPROFILE and started by Git Bash', async ($, on) => {
    // Native Windows has no HOME; an empty one stands for it here. The test
    // engine takes absolute POSIX paths only, so USERPROFILE has no drive
    // letter — its backslashes are what is turned around.
    const home = '/Users/a/.claude/skills/jaira-dispatcher/scripts/spawn.sh'
    const env = { OS: 'Windows_NT', HOME: '', USERPROFILE: '\\Users\\a', HERDR_ENV: '1', HERDR_BIN_PATH: 'herdr' }
    const w = world(on, { herdr: true, spawn: false, env })
    w.files.set(home, '')
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(ran(w, 'C:/Program Files/Git/bin/bash.exe', home)).toEqual([['C:/Program Files/Git/bin/bash.exe', home, '--no-worktree', 'mg677r', 'MG677R', 'explain', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('on native Windows, CLAUDE_CODE_GIT_BASH_PATH names the bash', async ($, on) => {
    const env = { OS: 'Windows_NT', HOME: '/repo', HERDR_ENV: '1', CLAUDE_CODE_GIT_BASH_PATH: 'D:\\Git\\bin\\bash.exe' }
    const w = world(on, { herdr: true, env })
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(ran(w, 'D:\\Git\\bin\\bash.exe', SPAWN)).toHaveLength(1)
    expect(ran(w, 'git', '--exec-path')).toEqual([])
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
    expect(w.opened).toEqual([{ id: 'cockpit-explanation', title: 'Explanation of MG677R', focus: true }])
    // Written before the shape was asked for: shown whole, as it always was,
    // and still answerable in words of one's own.
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect(await pane.find({ type: 'Text', text: 'The key leaves out the branch.\nTwo branches share one cache entry.' })).toBeDefined()
    expect(await pane.find({ key: 'pick-MG677R-own' })).toBeDefined()
    expect((await pane.find({ key: 'pick-MG677R-own' }))?.props.autoFocus).toBe(true)
  })

  const SHAPED = [
    '## Progress',
    '',
    '- **2026-10-12 15:00 · Alex** — explain: The cache key leaves out the branch.',
    'Warum: two branches then share one entry and read each other\'s results.',
    'a) add the branch to the key',
    'b) accept it as it is',
    'c) clear the cache on checkout, and on every pull too, so that no branch ever reads an entry another branch wrote before it',
    'Empfehlung: a - smallest change, no stale reads',
  ].join('\n')

  // The same as the Explain chat now writes it: numbered, and with what
  // picking each option does.
  const MARKED = SHAPED
    .replace('a) add the branch to the key', '1) add the branch to the key [give back]')
    .replace('b) accept it as it is', '2) accept it as it is [accept]')
    .replace('c) clear the cache', '3) clear the cache')
    .replace('Empfehlung: a', 'Empfehlung: 1')

  const CLASSED = [
    '## Progress',
    '',
    '- **2026-11-20 18:00 · Alex** — explain: Status: DoD 7/7 ok | tests ok | blockers ok | bughunt on HEAD fail',
    '[small edge] the export test covers 1MB, not 40MB',
    '[known limit] Windows paths longer than 260 characters are not read',
    '[not this ticket] the stale comment in load.go belongs to K7Q2MX',
    'Status of a thing in prose, no checks here',
    '1) accept it as it is [accept]',
    '2) close the small edge first [give back]',
    'Recommendation: 1 - no note is a blocker, it can be accepted as it is',
  ].join('\n')

  test('a classed explanation shows its checks, coloured class marks and the recommendation apart', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = CLASSED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    // A tick or a cross per check, in the theme's colours.
    const all = await pane.findAll({ type: 'Text' })
    expect(all.filter(x => x.text === '✓ ' && x.props.color === 'success')).toHaveLength(3)
    expect(all.filter(x => x.text === '✗ ' && x.props.color === 'error')).toHaveLength(1)
    expect(await pane.find({ type: 'Text', text: 'bughunt on HEAD' })).toBeDefined()
    // Only the marks coloured; the note text plain beside them.
    expect(all.some(x => x.text === '[small edge]' && x.props.color === 'warning')).toBe(true)
    expect(all.some(x => x.text === '[known limit]' && x.props.dimColor === true)).toBe(true)
    expect(all.some(x => x.text === '[not this ticket]' && x.props.color === 'suggestion')).toBe(true)
    expect(await pane.find({ type: 'Text', text: /the export test covers 1MB, not 40MB$/ })).toBeDefined()
    // A line in neither shape stays as written.
    expect(await pane.find({ type: 'Text', text: 'Status of a thing in prose, no checks here' })).toBeDefined()
    // The recommendation bold, in a colour of its own; the option it names chosen.
    expect(all.some(x => x.text === 'Recommendation: 1 - no note is a blocker, it can be accepted as it is' && x.props.bold === true && x.props.color === 'claude')).toBe(true)
    expect((await pane.find({ key: 'pick-MG677R-1' }))?.props.autoFocus).toBe(true)
  })

  test('a hunt skipped on purpose is drawn without a cross', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = CLASSED.replace('bughunt on HEAD fail', 'bughunt on HEAD skipped')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect(await pane.find({ type: 'Text', text: /bughunt on HEAD skipped/ })).toBeDefined()
    const all = await pane.findAll({ type: 'Text' })
    expect(all.filter(x => x.text === '✗ ')).toHaveLength(0)
    expect(all.filter(x => x.text === '✓ ')).toHaveLength(3)
  })

  test('a hunt nobody can judge is a dim question mark, not a cross', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = CLASSED.replace('bughunt on HEAD fail', 'bughunt on HEAD unknown')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    const all = await pane.findAll({ type: 'Text' })
    expect(all.some(x => x.text === '? ' && x.props.dimColor === true)).toBe(true)
    expect(all.filter(x => x.text === '✗ ')).toHaveLength(0)
    expect(await pane.find({ type: 'Text', text: 'bughunt on HEAD' })).toBeDefined()
  })

  test('a blocker mark is red', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = CLASSED.replace('[small edge]', '[blocker]')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect((await pane.findAll({ type: 'Text' })).some(x => x.text === '[blocker]' && x.props.color === 'error')).toBe(true)
  })

  test('a class mark in markdown emphasis is still a class mark', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = CLASSED.replace('[small edge]', '**[blocker]**')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect((await pane.findAll({ type: 'Text' })).some(x => x.text === '[blocker]' && x.props.color === 'error')).toBe(true)
  })

  test('a shaped explanation reads as paragraphs and one line per option, the recommended one chosen', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = SHAPED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    expect(texts).toEqual(expect.arrayContaining([
      'The cache key leaves out the branch.',
      "Warum: two branches then share one entry and read each other's results.",
      'Empfehlung: a - smallest change, no stale reads',
    ]))
    // Every option is a row of its own to pick, carrying its whole text in
    // a wrapping Text — never a one-line label the surface cuts — with its
    // digit as the key that picks it; the recommended one has the focus.
    const LONG = 'clear the cache on checkout, and on every pull too, so that no branch ever reads an entry another branch wrote before it'
    for (const [k, text] of [['a', 'add the branch to the key (recommended)'], ['b', 'accept it as it is'], ['c', LONG]] as const) {
      const row = await pane.find({ key: `pick-MG677R-${k}` })
      // Its digit is its place, the key that picks it, drawn by the engine.
      expect(row?.props.hotkey).toBe(String('abc'.indexOf(k) + 1))
      expect(row?.text).toBe(text)
      // A surface that draws no children draws the label: the whole option too.
      // The engine puts the digit before the label itself: the label is the
      // option alone, or a surface without children would draw '1: 1: …'.
      expect(row?.props.label).toBe(text)
      expect((await pane.find({ type: 'Text', text }))?.props.wrap).toBe('wrap')
    }
    expect((await pane.find({ key: 'pick-MG677R-a' }))?.props.autoFocus).toBe(true)
    expect((await pane.find({ key: 'pick-MG677R-b' }))?.props.autoFocus).toBeUndefined()
    expect((await pane.find({ key: 'pick-MG677R-own' }))?.text).toBe('Own answer…')
    expect((await pane.find({ key: 'pick-MG677R-own' }))?.props.hotkey).toBe('4')
  })

  test('an option that is neither step goes onto the ticket as a decision and to the teamlead\'s conversation', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = SHAPED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    await (await $.ui.mount(EXPLANATION('Explanation of MG677R'))).press({ key: 'pick-MG677R-b' })
    expect(ran(w, 'jaira', 'note')).toEqual([['jaira', 'note', 'MG677R', 'decision: b) accept it as it is', '@/.worktrees/repo-mg677r']])
    expect(w.closed).toEqual(['cockpit-explanation'])
    // No mark: nothing moves, and the conversation the band sits in hears it.
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R: b) accept it as it is (Person, Cockpit) - decided in the explanation pane, noted on the ticket, not acted on'])
    const band = await $.ui.mount(BAND)
    expect(await band.find({ type: 'Text', text: /MG677R \| decided/ })).toBeDefined()
    await band.press({ key: 'decision-MG677R' })
    expect(await (await $.ui.mount(MESSAGE('Decision on MG677R'))).find({ type: 'Text', text: 'b) accept it as it is' })).toBeDefined()
  })

  test('an own answer is typed in a pane that takes the keyboard and gives the work back with it as the reason', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = SHAPED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    await (await $.ui.mount(EXPLANATION('Explanation of MG677R'))).press({ key: 'pick-MG677R-own' })
    expect(w.closed).toEqual(['cockpit-explanation'])
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-answer', title: 'Answer on MG677R', focus: true })
    const pane = await $.ui.mount(ANSWER('Answer on MG677R'))
    expect((await pane.find({ key: 'answer-MG677R' }))?.props.autoFocus).toBe(true)
    // Enter on an empty field notes nothing and keeps the pane open.
    await pane.input({ key: 'answer-MG677R', text: '  ' })
    expect(ran(w, 'jaira', 'note')).toEqual([])
    await pane.input({ key: 'answer-MG677R', text: 'key by branch and commit' })
    const notes = ran(w, 'jaira', 'note').map(r => r[3])
    expect(notes.slice(0, 2)).toEqual(['decision: key by branch and commit', 'Given back from signoff: key by branch and commit'])
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.slice(0, 5)).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress'])
    expect(w.closed).toEqual(['cockpit-explanation', 'cockpit-answer'])
  })

  test('an option marked [give back] gives the work back with it as the reason, without a second press', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = MARKED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    // The mark is read, and drawn in the colour of the band's button for that
    // step; the recommendation highlighted; the option's own words neutral.
    expect((await pane.find({ key: 'pick-MG677R-1' }))?.text).toBe('add the branch to the key [give back] (recommended)')
    const colours = new Map((await pane.findAll({ type: 'Text' })).map(x => [x.text, x.props.color]))
    expect(colours.get(' [give back]')).toBe('error')
    expect(colours.get(' [accept]')).toBe('success')
    expect(colours.get(' (recommended)')).toBe('suggestion')
    expect(colours.get('add the branch to the key [give back] (recommended)')).toBeUndefined()
    expect((await pane.find({ key: 'pick-MG677R-1' }))?.props.autoFocus).toBe(true)
    await pane.press({ key: 'pick-MG677R-1' })
    expect(ran(w, 'jaira', 'note').map(r => r[3]).slice(0, 2)).toEqual(['decision: 1) add the branch to the key', 'Given back from signoff: 1) add the branch to the key'])
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.slice(0, 5)).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress'])
    // The decision itself is not sent: only the line every Give back hands this session.
    expect(w.submitted).toEqual(["jaira-cockpit: MG677R given back: 1) add the branch to the key, worker none (the ticket's note says why) (Person, Cockpit; signoff -> in-progress)"])
  })

  test('an option marked [accept] accepts the ticket as the band\'s Accept does', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = MARKED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    await (await $.ui.mount(EXPLANATION('Explanation of MG677R'))).press({ key: 'pick-MG677R-2' })
    expect(ran(w, 'jaira', 'note').map(r => r[3])).toEqual(['decision: 2) accept it as it is', 'cockpit: accepted by the person, signoff -> done'])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R: accepted, signoff -> done (Person, Cockpit)'])
  })

  test('after a restart an answered explanation is read with its decision, not offered again; a newer one is open', async ($, on) => {
    const w = world(on)
    // Answered in an earlier session: the decision stands after the explanation.
    w.bodies.MG677R = SHAPED + '\n- **2026-10-12 15:02 · Alex** — decision: c) clear the cache on checkout\n'
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    let pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect(await pane.find({ key: 'pick-MG677R-a' })).toBeUndefined()
    expect(await pane.find({ type: 'Text', text: 'Decided: c) clear the cache on checkout' })).toBeDefined()
    // A newer explanation, a new round's finding, is open again.
    await pane.unmount()
    w.bodies.MG677R += '- **2026-10-12 16:00 · Alex** — explain: A new finding.\na) fix it [give back]\n'
    await w.clock.advance(60_000)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect(await pane.find({ key: 'pick-MG677R-a' })).toBeDefined()
    expect(await pane.find({ type: 'Text', text: /^Decided:/ })).toBeUndefined()
  })

  test('a choice made in this session is not offered again before the notes are next read', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = SHAPED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const first = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    await first.press({ key: 'pick-MG677R-b' })
    await first.unmount()
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect(await pane.find({ key: 'pick-MG677R-a' })).toBeUndefined()
    expect(await pane.find({ type: 'Text', text: 'Decided: b) accept it as it is' })).toBeDefined()
  })

  test('a verdict of an earlier round keeps its Explain and says the ticket changed since', async ($, on) => {
    const w = world(on)
    const reserved = 'approve, but the cache key ignores the branch'
    w.tickets = [
      { ...SIGNOFF, handle: 'OLD123', round: 2, review_round: 1, review: { verdict: reserved, gaps: '' } },
      { ...SIGNOFF, handle: 'NEW456', round: 2, review_round: 2, review: { verdict: reserved, gaps: '' } },
      // Review sent it back with a blocker, the teamlead moved it on past a
      // new review: the person must not get it silent.
      { ...SIGNOFF, handle: 'BLK789', round: 8, review_round: 7, review: { verdict: 'Zurueck an in-progress: ein Blocker', gaps: 'blocker: the export drops the last row' } },
    ]
    await startSession($)
    const band = await $.ui.mount(BAND)
    for (const id of ['OLD123', 'NEW456', 'BLK789']) expect(await band.find({ key: `acceptance-${id}` })).toBeDefined()
    // The row marks it; the Acceptance pane says it whole.
    expect(await band.find({ type: 'Text', text: /OLD123 \| signoff \| .*review of round 1/ })).toBeDefined()
    expect(await band.find({ type: 'Text', text: /BLK789 \| signoff \| .*review of round 7/ })).toBeDefined()
    expect(await band.find({ type: 'Text', text: /review of round 2/ })).toBeUndefined()
    expect(await (await verdictOf($, band, 'OLD123')).find({ type: 'Text', text: /review of round 1, changed since/ })).toBeDefined()
    expect(await (await verdictOf($, band, 'NEW456')).find({ type: 'Text', text: /changed since/ })).toBeUndefined()
  })

  test('a later round\'s optimize gaps are no review notes under the row', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...HUMAN, handle: 'OPT123', question: '', round: 2, review_round: 1, review: { verdict: 'approve', gaps: 'removed the unused helper\nnote: inlined the close element' } },
      { ...SIGNOFF, handle: 'CUR456', round: 2, review_round: 2, review: { verdict: 'approve', gaps: 'note: a test is missing' } },
    ]
    await startSession($)
    const band = await $.ui.mount(BAND)
    expect(await (await verdictOf($, band, 'OPT123')).find({ type: 'Text', text: /removed the unused helper|inlined the close element/ })).toBeUndefined()
    expect(await (await verdictOf($, band, 'CUR456')).find({ type: 'Text', text: /🟡 note: a test is missing/ })).toBeDefined()
  })

  test('a verdict of an earlier round lists the review\'s own notes, kept beside it, not optimize\'s', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...HUMAN, handle: 'OLD123', question: '', round: 3, review_round: 2, review: { verdict: 'approve', gaps: 'removed the unused helper', verdict_gaps: 'note: the export test covers 1MB, not 40MB' } },
      { ...SIGNOFF, handle: 'BLK456', round: 2, review_round: 1, review: { verdict: 'approve', gaps: 'none', verdict_gaps: 'blocker: DoD 2 - an empty export writes a header' } },
    ]
    await startSession($)
    const band = await $.ui.mount(BAND)
    const old = await verdictOf($, band, 'OLD123')
    expect(await old.find({ type: 'Text', text: /🟡 note: the export test covers 1MB/ })).toBeDefined()
    expect(await old.find({ type: 'Text', text: /removed the unused helper/ })).toBeUndefined()
    expect(await (await verdictOf($, band, 'BLK456')).find({ type: 'Text', text: /🔴 blocker: DoD 2 - an empty export writes a header/ })).toBeDefined()
  })

  test('an earlier-round verdict without a kept copy reads no blocker out of a later lane\'s gaps', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, handle: 'OPB123', round: 3, review_round: 2, review: { verdict: 'approve', gaps: 'blocker: removed the unused helper' } }]
    await startSession($)
    const band = await $.ui.mount(BAND)
    expect(await band.find({ key: 'explain-OPB123' })).toBeUndefined()
    expect(await band.find({ type: 'Text', text: /removed the unused helper/ })).toBeUndefined()
  })

  test('an explanation stays readable once its verdict is an earlier round\'s', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, round: 2, review_round: 1, review: { verdict: 'approve, but the cache key ignores the branch', gaps: 'note: a test is missing' } }]
    w.bodies.MG677R = SHAPED + '\n- **2026-10-12 15:02 · Alex** — decision: c) clear the cache on checkout\n'
    await startSession($)
    const band = await $.ui.mount(BAND)
    await band.press({ key: 'read-MG677R' })
    const pane = await $.ui.mount(EXPLANATION('Explanation of MG677R'))
    expect(await pane.find({ type: 'Text', text: 'Decided: c) clear the cache on checkout' })).toBeDefined()
  })

  test('the row counts blockers and notes; the Acceptance pane lists each with its mark', async ($, on) => {
    const w = world(on)
    const verdict = 'The diff meets all five DoD items and the tests prove each one (87 pass).'
    w.tickets = [
      // A blocker line, however the list is marked.
      { ...SIGNOFF, handle: 'BL0CK1', review: { verdict, gaps: 'blocker: DoD 2 - an empty export still writes a header\nnote: the export test covers 1MB, not 40MB' } },
      { ...SIGNOFF, handle: 'BL0CK2', review: { verdict, gaps: '1. Blocker: a crash on an empty file' } },
      { ...SIGNOFF, handle: 'BL0CK3', review: { verdict, gaps: '- blocker: data loss on CRLF\r\n- note: a comment is stale' } },
      // ... and however it is emphasised in markdown.
      { ...SIGNOFF, handle: 'BL0CK4', review: { verdict, gaps: '**blocker:** a crash on an empty file' } },
      { ...SIGNOFF, handle: 'BL0CK5', review: { verdict, gaps: '- **Blocker**: data loss on CRLF' } },
      { ...SIGNOFF, handle: 'BL0CK6', review: { verdict, gaps: '`blocker`: the gate is skipped' } },
      // A reservation in the verdict, notes beside it.
      { ...SIGNOFF, handle: 'RESRV1', review: { verdict: 'approve, but the cache key ignores the branch', gaps: 'note: a test is missing' } },
      // Notes only: listed, no Explain.
      { ...SIGNOFF, handle: 'NOTES1', review: { verdict, gaps: 'note: the export test covers 1MB, not 40MB\nnote: pre-existing: load.go:41 accepts an empty file' } },
      { ...SIGNOFF, handle: 'NOTES2', review: { verdict, gaps: '1. Minor: the note after a successful move is fire-and-forget.' } },
      { ...SIGNOFF, handle: 'NOTES3', review: { verdict, gaps: '**note:** the blocker: rule reads fine' } },
      // Nothing found, said in many ways.
      { ...SIGNOFF, handle: 'AB12CD', review: { verdict, gaps: 'none' } },
      { ...SIGNOFF, handle: 'EF34GH', review: { verdict: 'approve', gaps: '' } },
      { ...SIGNOFF, handle: 'GH56JK', review: { verdict, gaps: 'none — the docs could say more' } },
      { ...SIGNOFF, handle: 'MN12PQ', review: { verdict, gaps: 'keine Lücken' } },
      { ...SIGNOFF, handle: 'PQ34RS', review: { verdict, gaps: 'No gaps found.' } },
      // Not reviewed yet: a blocker-looking line of the optimize lane is no review's.
      { ...HUMAN, handle: 'LM90NP', question: '', review: { verdict: '', gaps: 'blocker: inlined the single-use close element' } },
      { ...HUMAN, handle: 'LM91NP', question: '', review: { verdict: '', gaps: 'inlined the single-use close element' } },
    ]
    await startSession($)
    const band = await $.ui.mount(BAND)
    await band.press({ key: 'more' })
    const open = await $.ui.mount(BAND)
    // The row counts them; the notes themselves stand in the Acceptance pane.
    expect(await open.find({ text: /export test covers/ })).toBeUndefined()
    expect(await open.find({ type: 'Text', text: /NOTES1 \| signoff \| .*0 blockers \| 2 notes/ })).toBeDefined()
    expect(await open.find({ type: 'Text', text: /BL0CK3 \| signoff \| .*1 blocker \| 1 note $/ })).toBeDefined()
    expect(await open.find({ type: 'Text', text: /AB12CD \| signoff \| .*0 blockers \| 0 notes/ })).toBeDefined()
    const notes1 = await verdictOf($, open, 'NOTES1')
    expect(await notes1.find({ type: 'Text', text: /🟡 note: the export test covers 1MB, not 40MB/ })).toBeDefined()
    expect(await notes1.find({ type: 'Text', text: /🔵 note: pre-existing: load\.go:41/ })).toBeDefined()
    expect(await (await verdictOf($, open, 'NOTES2')).find({ type: 'Text', text: /🟡 Minor: the note after a successful move/ })).toBeDefined()
    const block3 = await verdictOf($, open, 'BL0CK3')
    expect(await block3.find({ type: 'Text', text: /🟡 note: a comment is stale/ })).toBeDefined()
    // A blocker is listed as one, not as a note, and 'none' is no finding.
    expect(await block3.find({ type: 'Text', text: /🔴 blocker: data loss on CRLF/ })).toBeDefined()
    expect(await block3.find({ type: 'Text', text: /🟡 blocker:/i })).toBeUndefined()
    expect(await (await verdictOf($, open, 'AB12CD')).find({ type: 'Text', text: /none/ })).toBeUndefined()
    expect(await open.find({ key: 'verdict-LM91NP' })).toBeUndefined()
  })

  test('the whole review verdict is read in the Acceptance pane, with no explanation on the ticket', async ($, on) => {
    const w = world(on)
    const long = 'The diff meets all five DoD items; the one weak spot is the cache key, which leaves out the branch, so two branches share an entry.'
    w.tickets = [{ ...SIGNOFF, review: { verdict: long } }, { ...HUMAN, question: '' }]
    await startSession($)
    const band = await $.ui.mount(BAND)
    expect(await band.find({ key: 'read-MG677R' })).toBeUndefined()
    await band.press({ key: 'acceptance-MG677R' })
    expect(w.opened).toEqual([{ id: 'cockpit-acceptance', title: 'Acceptance of MG677R', focus: true }])
    const pane = await $.ui.mount(ACCEPT('Acceptance of MG677R'))
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    expect(texts.some(x => x.endsWith(`review: ${long}`))).toBe(true)
    // It starts with the light, then what the ticket is for. No hunt was
    // read for this ticket: nobody can tell, so yellow.
    expect(texts[0]).toBe('● yellow: bughunt unknown')
    expect(texts.filter(x => !x.startsWith('●'))[0]).toBe('cache key')
  })

  test('the Explain tab closes once its explanation is on the ticket, not on an older one', async ($, on) => {
    const w = world(on, { herdr: true })
    // 1_000_000 ms is 1970-01-01 00:16:40 UTC: jaira stamps notes to the minute.
    w.bodies.MG677R = '## Progress\n\n- **1970-01-01 00:10 · Alex** — explain: an older one\n'
    w.tabs = [{ tab_id: 'w1:t4', label: 'MG677R/explain', agent_status: 'working' }]
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressExplain($, ui, 'MG677R', w)
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
    expect(await ui.find({ key: 'acceptance-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'acceptance-AB12CD' })).toBeDefined()
    expect(await ui.find({ key: 'acceptance-EF34GH' })).toBeUndefined()
    expect(await ui.find({ key: 'discuss-K7P2QX' })).toBeUndefined()
    expect((await ui.find({ key: 'theirs' }))?.props.label).toBe('2 waiting for Berk, Carla')
    await ui.press({ key: 'theirs' })
    const open = await $.ui.mount(BAND)
    expect(await open.find({ type: 'Text', text: /EF34GH \| signoff \| Berk/ })).toBeDefined()
    // A teammate's ticket: talked through, never accepted or given back from here.
    expect(await open.find({ key: 'discuss-K7P2QX' })).toBeDefined()
    for (const id of ['EF34GH', 'K7P2QX']) {
      expect(await open.find({ key: `acceptance-${id}` })).toBeUndefined()
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
    expect(await ui.find({ key: 'acceptance-GH56JK' })).toBeUndefined()
    expect((await ui.find({ key: 'more' }))?.props.label).toBe('+2 more')
    await ui.press({ key: 'more' })
    const open = await $.ui.mount(BAND)
    expect(await open.find({ key: 'acceptance-EF34GH' })).toBeDefined()
    expect(await open.find({ key: 'acceptance-GH56JK' })).toBeDefined()
    await open.press({ key: 'fewer' })
    const closed = await $.ui.mount(BAND)
    expect(await closed.find({ key: 'acceptance-GH56JK' })).toBeUndefined()
    expect(await closed.find({ key: 'more' })).toBeDefined()
  })

  test('the band reads top down: head, notice, waiting for you, decisions, the teammates\' tickets, the session line', async ($, on) => {
    const w = world(on, { herdr: true })
    w.me = { identity: 'Alexander Sacharov', aliases: ['Alexander Sacharov'] }
    const BERK = { ...SIGNOFF, handle: 'EF34GH', assignee: 'Berk' }
    // Accept takes MG677R to done; AB12CD still waits in signoff below the notice.
    const MINE = { ...SIGNOFF, handle: 'AB12CD' }
    w.tickets = [SIGNOFF, MINE, HUMAN, WORKING, BERK]
    // Accept's notice stays on the band.
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'working' }]
    await startSession($)
    // A decision from a Discuss chat, as in the test of its own above.
    await (await $.ui.mount(BAND)).press({ key: 'discuss-K7P2QX' })
    w.bodies.K7P2QX = '## Progress\n\n- **1970-01-01 00:17 · Alex** — decision: a page of its own\n'
    w.tickets = [SIGNOFF, MINE, { ...HUMAN, status: 'in-progress', question: '' }, WORKING, BERK]
    w.tabs = [{ tab_id: 'w1:t5', label: 'K7P2QX/discuss', agent_status: 'idle' }]
    await w.clock.advance(60_000)
    // A notice: what Accept says.
    await pressAccept($, await $.ui.mount(BAND), 'MG677R', w)
    w.usage.tokens = 80_000
    await answer($)
    await (await $.ui.mount(BAND)).press({ key: 'theirs' })
    const band = await drawn($)
    const at = ['jaira-cockpit', 'MG677R accepted', 'Waiting for you', 'AB12CD | signoff', 'K7P2QX | decided', '1 waiting for Berk', 'EF34GH | signoff | Berk', 'context 80k'].map(part => band.indexOf(part))
    expect(at.every(i => i >= 0)).toBe(true)
    expect(at).toEqual([...at].sort((a, b) => a - b))
  })

  test('opened by +N more, the person\'s own tickets still sit above the teammates\'', async ($, on) => {
    const w = world(on)
    w.me = { identity: 'Alexander Sacharov', aliases: ['Alexander Sacharov'] }
    const own = ['MG677R', 'AB12CD', 'CD34EF', 'GH56JK']
    w.tickets = [...own.map(handle => ({ ...SIGNOFF, handle })), { ...SIGNOFF, handle: 'QR34ST', assignee: 'Berk' }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'more' })
    await (await $.ui.mount(BAND)).press({ key: 'theirs' })
    const band = await drawn($)
    const theirs = band.indexOf('1 waiting for Berk')
    for (const id of own) {
      expect(band.indexOf(id)).toBeGreaterThanOrEqual(0)
      expect(band.indexOf(id)).toBeLessThan(theirs)
    }
    expect(band.indexOf('show fewer')).toBeLessThan(theirs)
    expect(band.indexOf('QR34ST | signoff | Berk')).toBeGreaterThan(theirs)
  })

  test('a ticket nobody explained yet offers no explanation to read', async ($, on) => {
    world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ key: 'acceptance-MG677R' })).toBeDefined()
    expect(await ui.find({ key: 'read-MG677R' })).toBeUndefined()
  })
})

describe('after Accept', () => {
  const press = async ($: Test$, w?: World) => {
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'MG677R', w)
  }
  const prLane = (w: World) => ran(w, 'bash', SPAWN).filter(r => r.includes('pr'))
  const failure = async ($: Test$) => {
    const pane = await $.ui.mount(MESSAGE('Next step after MG677R failed'))
    return (await pane.findAll({ type: 'Text' })).map(x => x.text).join('\n')
  }

  test('a branch ahead of the target without a pull request: the press opens one, by a person, in the ticket\'s worktree', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await press($, w)
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'upstream/main..feat/mg677r', '@/repo']])
    expect(ran(w, 'gh')).toEqual([
      ['gh', 'repo', 'view', 'git@github.com:acme/repo.git', '--json', 'isFork,parent,nameWithOwner,defaultBranchRef', '@/repo'],
      ['gh', 'repo', 'view', 'upstream/repo', '--json', 'isFork,parent,nameWithOwner,defaultBranchRef', '@/repo'],
      ['gh', 'pr', 'list', '--repo', 'upstream/repo', '--head', 'feat/mg677r', '--state', 'all', '--json', 'number', '@/repo'],
    ])
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
    // After the move, never before it: the press accepted first.
    const order = w.ran.map(r => r[1] === SPAWN ? 'spawn' : r.includes('--force') ? 'move' : '').filter(Boolean)
    expect(order).toEqual(['move', 'spawn'])
    expect(w.opened).toEqual([])
  })

  test('a pull request in any state already there: nothing starts, not even with a ticket waiting next, and the press says accepted', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prs = [{ number: 12 }]
    w.next = ['T8QS1Y', 'AB12CD']
    await press($, w)
    // Which ticket runs next is the teamlead's call, never the button's (X6JHQX).
    expect(ran(w, 'jaira', 'next')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted: done')
  })

  test('nothing ahead of the target: no pull request is asked for, no tab opens, and the press says accepted', async ($, on) => {
    const w = world(on, { herdr: true })
    w.next = ['AB12CD']
    await press($, w)
    expect(ran(w, 'gh', 'pr')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted: done')
  })

  test('GitLab, by jaira.forge or by host, asks glab', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 1
    w.forge = 'gitlab'
    await press($, w)
    expect(ran(w, 'glab', 'mr')).toEqual([['glab', 'mr', 'list', '--repo', 'upstream/repo', '--source-branch', 'feat/mg677r', '--all', '-F', 'json', '@/repo']])
    expect(prLane(w)).toHaveLength(1)
  })

  test('GitHub Enterprise behind a user@host:owner/repo remote: gh is asked host/owner/repo every time (XSGWPX)', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 1
    w.forge = 'github'
    w.origin = 'alice@ghe.example.com:acme/repo.git'
    w.remotes = 'origin\talice@ghe.example.com:acme/repo.git (fetch)\nupstream\talice@ghe.example.com:upstream/repo.git (fetch)\n'
    await press($, w)
    expect(ran(w, 'gh')).toEqual([
      ['gh', 'repo', 'view', 'ghe.example.com/acme/repo', '--json', 'isFork,parent,nameWithOwner,defaultBranchRef', '@/repo'],
      ['gh', 'repo', 'view', 'ghe.example.com/upstream/repo', '--json', 'isFork,parent,nameWithOwner,defaultBranchRef', '@/repo'],
      ['gh', 'pr', 'list', '--repo', 'ghe.example.com/upstream/repo', '--head', 'feat/mg677r', '--state', 'all', '--json', 'number', '@/repo'],
    ])
    expect(ran(w, 'git', 'fetch')).toEqual([['git', 'fetch', 'upstream', 'main', '@/repo']])
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
    await pressAccept($, ui, 'ZZ99ZZ', w)
    // No commit of its own beyond the main checkout's HEAD: Accept looks for
    // no carrier and moves there, as before; only the pull request asks.
    expect(ran(w, 'git', 'branch')).toEqual([['git', 'branch', '--list', '--format=%(refname:short)', '*/zz99zz', '@/repo']])
    expect(ran(w, 'jaira', 'move', 'ZZ99ZZ', '--to', 'done', '--force').map(r => r.at(-1))).toEqual(['@/repo'])
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'upstream/main..feat/zz99zz', '@/repo']])
    expect(ran(w, 'git', 'worktree', 'add')).toEqual([['git', 'worktree', 'add', '/.worktrees/repo-zz99zz', 'feat/zz99zz', '@/repo']])
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-zz99zz', '@/.worktrees/repo-zz99zz']])
  })

  test('a branch left behind beside the Herdr review branch of it: the pull request opens from the ticket branch itself', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }]
    w.branches = ['review/task/zz99zz', 'task/zz99zz']
    w.subjects = ['fix(ZZ99ZZ): one']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'ZZ99ZZ', w)
    expect(ran(w, 'git', 'worktree', 'add')).toEqual([['git', 'worktree', 'add', '/.worktrees/repo-zz99zz', 'task/zz99zz', '@/repo']])
  })

  test('two tickets on one branch: no pull request while the other is not in done, nothing starts, and the press says which one is missing', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 4
    w.subjects = ['feat(MG677R): cache key', 'fix(T8QS1Y): filter', 'feat(T8QS1Y): filter', 'docs: no ticket']
    w.next = ['AB12CD']
    await press($, w)
    expect(ran(w, 'gh', 'pr')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; no pull request for feat/mg677r yet: T8QS1Y not in done')
  })

  test('two tickets on one branch, both in done: the pull request opens; a handle the board no longer lists is filed away and holds nothing up', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 4
    w.tickets = [SIGNOFF, { ...WORKING, status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'fix(T8QS1Y): filter', 'chore(OLD123): filed away long ago']
    await press($, w)
    expect(prLane(w)).toHaveLength(1)
  })

  test('a ticket worked on another ticket\'s branch finds that branch by its own commit, the one branch holding it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }, { ...WORKING, handle: 'MG677R', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'fix(ZZ99ZZ): on the same branch']
    w.branchCommits = [{ sha: 'abc123', subject: 'fix(ZZ99ZZ): on the same branch' }]
    // A main checkout detached mid-rebase lists first and is no branch.
    w.holders = ['(HEAD detached at abc123)', 'master', 'feat/mg677r']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'ZZ99ZZ', w)
    // Once against the main checkout's HEAD, where Accept moves it, and once
    // against the target, for the pull request.
    expect(ran(w, 'git', 'log', '--branches')).toEqual([
      ['git', 'log', '--branches', '--not', 'HEAD', '--format=%H %s', '@/repo'],
      ['git', 'log', '--branches', '--not', 'upstream/main', '--format=%H %s', '@/repo'],
    ])
    const contains = ['git', 'branch', '--contains', 'abc123', '--format=%(refname:short)', '@/repo']
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([contains, contains])
    // The move runs where its commits are, the other ticket's worktree, not
    // the main checkout, which reaches none of them (VHGWQF).
    expect(ran(w, 'jaira', 'move', 'ZZ99ZZ', '--to', 'done', '--force').map(r => r.at(-1))).toEqual(['@/.worktrees/repo-mg677r'])
    // Its worktree is the other ticket's: the role pushes from there.
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  // A wip or a branch stacked on the ticket's holds its commit as well: Accept
  // opens no pull request for a guess, and says between which (Alex, QTQW15).
  test('several branches hold its own commit: no pull request, the press names them', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }, { ...WORKING, handle: 'MG677R', status: 'done' }]
    w.branchCommits = [{ sha: 'abc123', subject: 'fix(ZZ99ZZ): on the same branch' }]
    w.holders = ['master', 'feat/mg677r', 'wip']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'ZZ99ZZ', w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    // Nothing moves, in the main checkout or anywhere: which worktree the
    // move belongs in is a guess, and the press names the candidates (VHGWQF).
    expect(ran(w, 'jaira', 'move', 'ZZ99ZZ', '--to', 'done')).toEqual([])
    expect(ran(w, 'jaira', 'pull', 'ZZ99ZZ')).toEqual([])
    const pane = await $.ui.mount(MESSAGE('Accept ZZ99ZZ failed'))
    expect(await pane.find({ type: 'Text', text: /feat\/mg677r and wip could each carry ZZ99ZZ's commits.*nothing moved - check out the right one and run 'jaira move ZZ99ZZ --to done --force' there/ })).toBeDefined()
  })

  // Branches named after it decide nothing: a merged feat/ and a fix/ kept
  // beside it, while its newest own commit sits on another ticket's branch.
  test('branches named after it do not pick the carrier: the branch holding its commit does', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }, { ...WORKING, handle: 'MG677R', status: 'done' }]
    w.branches = ['feat/zz99zz', 'fix/zz99zz']
    w.branchCommits = [{ sha: 'abc123', subject: 'fix(ZZ99ZZ): on the other branch' }]
    w.holders = ['feat/mg677r']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'ZZ99ZZ', w)
    expect(ran(w, 'jaira', 'move', 'ZZ99ZZ', '--to', 'done', '--force').map(r => r.at(-1))).toEqual(['@/.worktrees/repo-mg677r'])
  })

  // Two branches named after it, both merged: the main checkout reaches every
  // commit of its own, so no carrier is asked for and nothing is ambiguous.
  test('branches named after it but nothing of its own beyond HEAD: Accept moves in the main checkout', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }]
    w.branches = ['feat/zz99zz', 'fix/zz99zz']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'ZZ99ZZ', w)
    expect(ran(w, 'jaira', 'move', 'ZZ99ZZ', '--to', 'done', '--force').map(r => r.at(-1))).toEqual(['@/repo'])
  })

  test('a ticket named only beside others in a subject has no commit of its own: no branch, no pull request (Alex, QTQW15)', async ($, on) => {
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
    await pressAccept($, ui, 'ZZ99ZZ', w)
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
  })

  test('the branch search reads subjects only: a merge commit quoting the ticket in its body is passed over', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }, { ...WORKING, handle: 'MG677R', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'fix(ZZ99ZZ): on the same branch']
    w.branchCommits = [
      { sha: 'merge1', subject: "Merge branch 'feat/other' into feat/stack", body: 'fix(ZZ99ZZ): on the same branch' },
      { sha: 'abc123', subject: 'fix(ZZ99ZZ): on the same branch' },
    ]
    w.holders = ['feat/mg677r']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await pressAccept($, ui, 'ZZ99ZZ', w)
    const contains = ['git', 'branch', '--contains', 'abc123', '--format=%(refname:short)', '@/repo']
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([contains, contains])
  })

  test('in a fork, ahead counts against the parent\'s default branch, fetched first — never the lagging local master', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await press($, w)
    expect(ran(w, 'git', 'fetch')).toEqual([['git', 'fetch', 'upstream', 'main', '@/repo']])
    const order = w.ran.map(r => r[1] === 'fetch' ? 'fetch' : r[1] === 'log' ? 'log' : '').filter(Boolean)
    // Every log after the fetch: the branch search and the count alike.
    expect(order[0]).toBe('fetch')
    expect(order.slice(1).every(x => x === 'log')).toBe(true)
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'upstream/main..feat/mg677r', '@/repo']])
  })

  test('not a fork: the base is origin\'s own default branch', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 1
    w.origin = 'git@github.com:upstream/repo.git'
    w.remotes = 'origin\tgit@github.com:upstream/repo.git (fetch)\n'
    await press($, w)
    expect(ran(w, 'git', 'fetch')).toEqual([['git', 'fetch', 'origin', 'main', '@/repo']])
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'origin/main..feat/mg677r', '@/repo']])
  })

  test('a fork without its parent as a remote: nothing is counted, nothing starts, and the press says which remote to add', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.remotes = 'origin\tgit@github.com:acme/repo.git (fetch)\n'
    await press($, w)
    expect(ran(w, 'git', 'log')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toContain('Pull requests go to upstream/repo, and no remote here points at it')
  })

  test('no gh: the pull request goes unchecked, nothing starts, and the press says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prs = null
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(ran(w, 'jaira', 'move', 'MG677R').map(r => r[5])).toEqual(['--dry-run', '--force'])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; pull requests not checked (gh repo view: gh: command not found)')
  })

  test('a forge that answers the repository but not the pull requests: nothing starts, and the press says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prListFails = true
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toContain('pull requests not checked (gh: HTTP 401')
  })

  test('a subject naming several tickets holds the pull request for each of them', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [SIGNOFF, WORKING, { ...WORKING, handle: 'AB12CD', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'chore(AB12CD,T8QS1Y): both at once']
    await press($, w)
    expect(prLane(w)).toEqual([])
    expect(w.toasts.at(-1)).toContain('no pull request for feat/mg677r yet: T8QS1Y not in done')
  })

  test('no remote: nothing starts, and the pane says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('.jaira/pr says none and the host names no forge: the press lands on origin through the pr lane, asks no pull request, and says landing', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = 'git@git.example.de:acme/repo.git'
    w.basePr = '# solo\nnone\n'
    await press($, w)
    expect(ran(w, 'gh')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
    expect(w.toasts.at(-1)).toContain('MG677R: landing feat/mg677r on the default branch')
    expect(w.toasts.join('\n')).not.toContain('pull request for')
  })

  test('.jaira/pr says none and there is no origin: nothing lands, and the pane says there is nothing to land on', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    w.basePr = 'none\n'
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('no remote origin to land on')
  })

  test('.jaira/pr says none and the forge does not answer: nothing lands, since whether origin is a fork cannot be told, and the pane says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.prs = null
    w.basePr = '\uFEFFnone\r\n'
    await press($, w)
    expect(prLane(w)).toEqual([])
    const why = await failure($)
    expect(why).toContain('Nothing lands from here')
    expect(why).toContain('gh: command not found')
  })

  test('.jaira/pr says none and another ticket on the branch is not in done: nothing lands, and the press says which', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.origin = null
    w.subjects = ['feat(MG677R): cache key', 'fix(T8QS1Y): filter']
    w.basePr = 'none\n'
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; feat/mg677r does not land yet: T8QS1Y not in done')
  })

  test('.jaira/pr says none only in this checkout, not in the base: the base decides, so it is the pull request path', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    w.files.set('/repo/.jaira/pr', 'none\n')
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('.jaira/pr says none and no forge names the base: it is read from origin\'s default branch, as the gate reads it, not the local one', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.origin = 'git@git.example.de:acme/repo.git'
    w.originHead = 'origin/main'
    w.basePr = 'none\n'
    await press($, w)
    expect(ran(w, 'git', 'show')).toEqual([['git', 'show', 'origin/main:.jaira/pr', '@/repo']])
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'origin/main..feat/mg677r', '@/repo']])
    expect(prLane(w)).toHaveLength(1)
  })

  test('.jaira/pr with a typo is the pull request path: no forge still says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    w.basePr = 'nnone\n'
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('a host neither GitHub nor GitLab: nothing starts, and the pane names the setting that settles it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = 'git@git.example.de:acme/repo.git'
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('git config jaira.forge')
  })

  test('no pull request due and no Herdr: nothing is missing, so nothing fails, and the press says accepted', async ($, on) => {
    const w = world(on)
    w.next = ['AB12CD']
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted: done')
  })

  test('a socket whose file is gone is no parent: the press still opens the pull request, which asks in its own tab', async ($, on) => {
    const w = world(on, { herdr: true, socket: false })
    w.ahead = 3
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--no-parent', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('Herdr App (HC_TERMINAL_ID, no HERDR_ENV): the press opens the pull request through spawn.sh all the same', async ($, on) => {
    const w = world(on, { herdr: true, env: { HOME: '/home/a', HC_TERMINAL_ID: 't3' } })
    w.ahead = 3
    await press($, w)
    // No messaging socket in this session: no parent, the role asks in its tab.
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--no-parent', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('no Herdr: nothing starts, and the pane says the step needs it', async ($, on) => {
    const w = world(on)
    w.ahead = 3
    await press($, w)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs a pull request')
  })

  test('a refused Accept starts no next step at all', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.refusals = [...w.refusals, { code: 'needs_commits', message: 'no commit names MG677R' }]
    await press($, w)
    expect(ran(w, 'git', 'log')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.submitted).toEqual([])
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

describe('the head of every pane', () => {
  // A head row as the terminal lays it out: '[ label ]' per Button, one
  // column between the tabs, Close pushed to the right edge.
  type Node = { type: string; props: Record<string, unknown>; children?: Node[] }
  const buttonsOf = (n: Node | undefined): string[] =>
    n === undefined ? [] : n.type === 'Button' ? [n.props.label as string] : (n.children ?? []).flatMap(buttonsOf)
  async function headOf(pane: { find: (q: { key: string }) => Promise<unknown> }, key: string) {
    const row = (await pane.find({ key })) as Node | undefined
    expect(row).toBeDefined()
    // One line: a row that never wraps, its parts never shrunk to a second line.
    expect(row!.props).toMatchObject({ flexDirection: 'row', flexWrap: 'nowrap', justifyContent: 'space-between' })
    const [tabs, close] = row!.children
    expect(tabs!.props).toMatchObject({ flexWrap: 'nowrap', flexShrink: 0 })
    expect(close!.props).toMatchObject({ flexShrink: 0 })
    const left = buttonsOf(tabs!)
    const right = buttonsOf(close!)
    expect(right).toEqual(['Close'])
    // Marked as the pane's close: a desktop draws its native close for it.
    expect(close!.children![0]!.props.role).toBe('dismiss')
    const drawn = [...left, ...right].map(l => `[ ${l} ]`)
    // Every button is hit across at least three columns, brackets included.
    for (const d of drawn) expect(d.length).toBeGreaterThanOrEqual(3)
    const width = drawn.reduce((n, d) => n + d.length, 0) + Math.max(0, left.length - 1) + (left.length > 0 ? 1 : 0)
    return { left, width }
  }

  test('Overview: the tabs and Close on one line that fits 80 columns', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...WORKING, assignee: 'Alex', updated_at: ago(0.1) }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    const { left, width } = await headOf(pane, 'overview-top-row')
    expect(left).toEqual(['Queue', 'Log', 'Workers'])
    expect(width).toBeLessThanOrEqual(80)
  })

  test('Verdict, Give back, Explanation, Message and the ticket pane: Close at the right of a one-line head, first', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = ['## Progress', '', '- **2026-10-12 15:00 · Alex** — explain: The key leaves out the branch.'].join('\n')
    await startSession($)
    const band = await $.ui.mount(BAND)
    const first = async (pane: { find: (q: { type: 'Box' }) => Promise<unknown> }, key: string) =>
      expect(((await pane.find({ type: 'Box' })) as Node).children[0]?.props.key).toBe(key)
    const panes: [string, () => Promise<{ find: never; press: (q: { key: string }) => Promise<unknown> }>, string][] = [
      ['acceptance', async () => { await band.press({ key: 'acceptance-MG677R' }); return $.ui.mount(ACCEPT('Acceptance of MG677R')) as never }, 'cockpit-acceptance'],
      ['explanation', async () => { await band.press({ key: 'read-MG677R' }); return $.ui.mount(EXPLANATION('Explanation of MG677R')) as never }, 'cockpit-explanation'],
      ['message', async () => {
        // Another session moved it first: that is said in the message pane.
        w.movedFirst = true
        await (await pressAccept($, band, 'MG677R', w)).unmount()
        return $.ui.mount(MESSAGE('Accept MG677R failed')) as never
      }, 'cockpit-message'],
      ['ticket', async () => {
        await $.command.run(typed('ticket MG677R'))
        return $.ui.mount({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-ticket', props: { title: 'Ticket MG677R', isFocused: false, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } } as never) as never
      }, 'cockpit-ticket'],
    ]
    for (const [name, open, id] of panes) {
      const pane = await open()
      await first(pane as never, `${name}-head`)
      const { left, width } = await headOf(pane as never, `${name}-head`)
      expect(left).toEqual([])
      expect(width).toBeLessThanOrEqual(80)
      await pane.press({ key: `${name}-close` })
      expect(w.closed.at(-1)).toBe(id)
      await (pane as unknown as { unmount: () => Promise<unknown> }).unmount()
    }
  })

  test('a pane left with nothing to show keeps its head, and its Close closes it', async ($, on) => {
    const w = world(on)
    await startSession($)
    // Mounted with nothing opened: the state each reads is empty, as after
    // its ticket left the board.
    const panes: [string, unknown, string][] = [
      ['acceptance', ACCEPT('Acceptance of MG677R'), 'cockpit-acceptance'],
      ['explanation', EXPLANATION('Explanation of MG677R'), 'cockpit-explanation'],
      ['ticket', { plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-ticket', props: { title: 'Ticket MG677R', isFocused: false, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }, 'cockpit-ticket'],
    ]
    for (const [name, element, id] of panes) {
      const pane = await $.ui.mount(element as never)
      expect(await pane.find({ type: 'Text', text: 'nothing recorded yet' })).toBeDefined()
      const { left } = await headOf(pane as never, `${name}-head`)
      expect(left).toEqual([])
      await pane.press({ key: `${name}-close` })
      expect(w.closed.at(-1)).toBe(id)
      await pane.unmount()
    }
  })
})

describe('give back', () => {
  async function giveBack($: Test$, w: World) {
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'acceptance-MG677R' })
    const pane = await $.ui.mount(ACCEPT('Acceptance of MG677R'))
    const typed = pane.input({ key: 'typed-MG677R', text: 'der Branch fehlt im Schluessel' })
    // Herdr is given a second between the text and its Enter.
    await w.clock.advance(1000)
    await typed
    await pane.unmount()
    await ui.unmount()
    return w
  }

  test('asks what is missing in a pane that takes the keyboard, never in a band field the press leaves unfocused', async ($, on) => {
    const w = world(on)
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'acceptance-MG677R' })
    // Live, the band's own field never got the ring: the press keeps it on
    // the button. A pane opened with focus hands it to its autoFocus field.
    expect(await ui.find({ key: 'typed-MG677R' })).toBeUndefined()
    expect(w.opened).toEqual([{ id: 'cockpit-acceptance', title: 'Acceptance of MG677R', focus: true }])
    const pane = await $.ui.mount(ACCEPT('Acceptance of MG677R'))
    expect((await pane.find({ key: 'typed-MG677R' }))?.props.autoFocus).toBe(true)
    // Enter on an empty field gives nothing back, accepts nothing, and keeps the pane open.
    await pane.input({ key: 'typed-MG677R', text: '  ' })
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(w.closed).toEqual([])
    await pane.input({ key: 'typed-MG677R', text: 'fehlt' })
    expect(w.closed).toEqual(['cockpit-acceptance'])
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.slice(0, 5)).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress'])
  })

  test('notes why and moves the ticket back to in-progress, without --force, so any other refusal still stops it', async ($, on) => {
    const w = world(on)
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'note', 'MG677R')[0]?.[3]).toBe('Given back from signoff: der Branch fehlt im Schluessel')
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]).toEqual(['jaira', 'move', 'MG677R', '--to', 'in-progress', '@/.worktrees/repo-mg677r'])
    // The move first: the note is written only once the ticket went back.
    const at = (verb: string) => w.ran.findIndex(r => r[0] === 'jaira' && r[1] === verb && r[2] === 'MG677R')
    expect(at('move')).toBeLessThan(at('note'))
  })

  test('a refused move writes no note: the ticket stays where it is and the refusal is shown', async ($, on) => {
    const w = world(on)
    w.version = '0.3.6'
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(1)
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(w.submitted).toEqual([])
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: 'Give back MG677R failed', focus: true })
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
    expect(ran(w, 'bash', SPAWN)).toEqual([])
  })

  test('a note that fails after the move is shown with the reason typed, and the worker still gets the ticket', async ($, on) => {
    const w = world(on, { herdr: true, store: big(1_000_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    w.noteFails = true
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(1)
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: 'Note on MG677R failed', focus: true })
    const pane = await $.ui.mount(MESSAGE('Note on MG677R failed'))
    expect(await pane.find({ type: 'Text', text: /MG677R went back to in-progress, but the note with your reason was not written: jaira: ticket file is locked/ })).toBeDefined()
    expect(await pane.find({ type: 'Text', text: /jaira note MG677R 'Given back from signoff: der Branch fehlt im Schluessel'/ })).toBeDefined()
    expect(ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')).toHaveLength(1)
  })

  test('past half of its window a new worker takes the round', async ($, on) => {
    const w = world(on, { herdr: true, store: big(500_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toStartWith("cockpit: given back to a new worker - old worker's context 300k >= 250k (50% of its window); ")
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
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to the old worker - tab alive, cache warm for 40 more min, context 80k')
  })

  // QTQW15: the new worker reports to the ticket's teamlead, starts where the
  // ticket's commits are, and is never a second one beside a worker at work.
  const PARENT = '- **2026-11-20 10:00 · Alex** — parent session: uds:/run/lead.sock\n'
  const spawned = (w: World) => ran(w, 'bash', SPAWN).map(r => r.slice(2, -1))
  const lastNote = (w: World) => ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3] ?? ''

  test('the new worker reports to the parent session the ticket names, and the teamlead hears who took it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.bodies.MG677R = `## Progress\n- **2026-11-19 09:00 · Alex** — parent session: uds:/run/old.sock\n${PARENT}`
    await startSession($)
    await giveBack($, w)
    expect(spawned(w)).toEqual([['--no-worktree', '--parent', 'uds:/run/lead.sock', 'mg677r', 'MG677R', 'worker', '/.worktrees/repo-mg677r']])
    expect(lastNote(w)).toContain("reports to uds:/run/lead.sock, the ticket's parent session")
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R given back: der Branch fehlt im Schluessel, worker new MG677R/worker (pane w1:p42, on feat/mg677r) (Person, Cockpit; signoff -> in-progress)'])
  })

  test('a parent session that is gone: this session takes the reports, never no parent at once', async ($, on) => {
    const w = world(on, { herdr: true })
    w.bodies.MG677R = PARENT
    w.deadParents = ['uds:/run/lead.sock']
    await startSession($)
    await giveBack($, w)
    expect(spawned(w).map(a => a.slice(1, 3))).toEqual([['--parent', 'uds:/run/lead.sock'], ['--parent', 'uds:/tmp/cc-socks-1000/77.sock']])
    expect(lastNote(w)).toContain('reports to this session (uds:/run/lead.sock does not listen)')
  })

  test('no parent anywhere: the worker asks in its own tab, and the note says so', async ($, on) => {
    const w = world(on, { herdr: true, env: { HOME: '/home/a', HERDR_ENV: '1', HERDR_BIN_PATH: 'herdr', HERDR_PANE_ID: 'w1:p1' } })
    await startSession($)
    await giveBack($, w)
    expect(spawned(w).map(a => a[1])).toEqual(['--no-parent'])
    expect(lastNote(w)).toContain('no parent session reachable')
  })

  test('no worktree of its own: the one branch holding its commit, in the worktree that branch has', async ($, on) => {
    const w = world(on, { herdr: true })
    w.worktrees = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n\nworktree /.worktrees/repo-fix-login\nHEAD 3\nbranch refs/heads/fix-login\n'
    w.branchCommits = [{ sha: 'c2', subject: 'fix(MG677R): the key' }, { sha: 'c1', subject: 'chore(X1Y2Z3): other' }]
    w.holders = ['fix-login']
    await startSession($)
    await giveBack($, w)
    expect(spawned(w)).toEqual([['--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'worker', '/.worktrees/repo-fix-login']])
    expect(lastNote(w)).toContain('branch fix-login holds the newest commit naming the ticket, worktree /.worktrees/repo-fix-login')
  })

  test('several branches hold the commit and none is the ticket\'s own: nothing is started, the teamlead hears which', async ($, on) => {
    const w = world(on, { herdr: true })
    w.worktrees = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n\nworktree /.worktrees/repo-0gmyqb\nHEAD 4\nbranch refs/heads/feat/0gmyqb\n\nworktree /.worktrees/repo-fix-login\nHEAD 3\nbranch refs/heads/fix-login\n'
    w.branchCommits = [{ sha: 'c2', subject: 'fix(MG677R): the key' }]
    w.holders = ['feat/0gmyqb', 'fix-login']
    // No tip lies under the other: master was merged into fix-login after
    // feat/0gmyqb was stacked on it.
    await startSession($)
    await giveBack($, w)
    expect(spawned(w)).toEqual([])
    expect(lastNote(w)).toContain("feat/0gmyqb and fix-login could each be MG677R's branch and none is its worktree's")
    expect(w.submitted[0]).toContain('worker none: feat/0gmyqb and fix-login could each be')
    // The band says why, not only the note.
    expect(w.toasts.at(-1)).toContain("MG677R: feat/0gmyqb and fix-login could each be MG677R's branch")
  })

  test('a branch whose worktree was removed is checked out again', async ($, on) => {
    const w = world(on, { herdr: true })
    w.worktrees = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n'
    w.branchCommits = [{ sha: 'c2', subject: 'fix(MG677R): the key' }]
    w.holders = ['fix-login']
    await startSession($)
    await giveBack($, w)
    expect(spawned(w)).toEqual([['--branch', 'fix-login', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'worker', '/repo']])
  })

  test('two branches hold the commit: the ticket\'s own worktree branch, not the one stacked on it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.branchCommits = [{ sha: 'c2', subject: 'fix(MG677R): the key' }]
    w.holders = ['stacked-on-it', 'feat/mg677r']
    await startSession($)
    await giveBack($, w)
    expect(spawned(w)[0]?.at(-1)).toBe('/.worktrees/repo-mg677r')
  })

  test('no branch of the ticket\'s at all: spawn.sh makes one, as before', async ($, on) => {
    const w = world(on, { herdr: true })
    w.worktrees = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n'
    await startSession($)
    await giveBack($, w)
    expect(spawned(w)).toEqual([['--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'worker', '/repo']])
    expect(lastNote(w)).toContain('no branch carries the ticket yet')
  })

  test('a worker at work gets the decision, however full its context: its tab is never closed, no second one starts', async ($, on) => {
    const w = world(on, { herdr: true, store: big(500_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker', agent_status: 'working' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')).toHaveLength(1)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(spawned(w)).toEqual([])
    expect(lastNote(w)).toBe('cockpit: given back to the old worker - the old worker is at work in its tab')
  })

  test('a worker at work with nothing recorded to reach it: still no second one, and the note says where it is', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker', agent_status: 'working' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(spawned(w)).toEqual([])
    expect(lastNote(w)).toContain('tell it in tab MG677R/worker - no second worker started')
    // The teamlead hears there is one, so it starts none either.
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R given back: der Branch fehlt im Schluessel, worker MG677R/worker at work, not reached: type the decision into its tab, start no other (Person, Cockpit; signoff -> in-progress)'])
  })

  test('a worker at an approval dialog gets nothing typed and no Enter, is not replaced, and the person hears which tab waits', async ($, on) => {
    // Its pane is recorded: without the check the decision would be typed in.
    const w = world(on, { herdr: true, store: big(500_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker', agent_status: 'blocked' }]
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'herdr', 'pane', 'send-keys')).toEqual([])
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(spawned(w)).toEqual([])
    expect(lastNote(w)).toBe('cockpit: given back to nobody (no worker started) - the old worker waits at an approval dialog in tab MG677R/worker; no Enter pressed, no second worker started')
    expect(w.toasts.at(-1)).toBe('MG677R back in in-progress; its worker waits at an approval dialog in tab MG677R/worker - answer it there, then type the decision into it. No Enter was pressed for you.')
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R given back: der Branch fehlt im Schluessel, worker tab MG677R/worker waits at an approval dialog: answer it there, then type the decision into it (Person, Cockpit; signoff -> in-progress)'])
  })

  test('the recorded pane at a dialog while its tab looks idle: Herdr is asked about the pane itself, and nothing is typed', async ($, on) => {
    // A warm cache and a small context: the old worker would get the text.
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 20 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker', agent_status: 'idle' }]
    w.paneStates = ['blocked']
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'get', 'w1:p9')).toHaveLength(1)
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'herdr', 'pane', 'send-keys')).toEqual([])
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(spawned(w)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R back in in-progress; its worker waits at an approval dialog in pane w1:p9 - answer it there, then type the decision into it. No Enter was pressed for you.')
  })

  test('a dialog that opens in the second between the text and its Enter gets no Enter, and the decision counts as handed over', async ($, on) => {
    const w = world(on, { herdr: true, store: big(500_000) })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker', agent_status: 'working' }]
    w.paneStates = ['working', 'blocked']
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')).toHaveLength(1)
    expect(ran(w, 'herdr', 'pane', 'send-keys')).toEqual([])
    expect(spawned(w)).toEqual([])
    // The text arrived (Herdr App's hc send even pressed Enter with it): the
    // old worker has the decision, nobody is told to type it again.
    expect(lastNote(w)).toBe('cockpit: given back to the old worker - the old worker is at work in its tab; a dialog opened before the Enter, none pressed')
    expect(w.toasts.at(-1)).toBe('MG677R given back to the worker that carried it; an approval dialog opened in pane w1:p9 before the Enter - answer it there, then press Enter if the decision still waits in its input.')
    expect(w.submitted.at(-1)).toContain('worker that carried it MG677R/worker (pane w1:p9) at an approval dialog')
  })

  test('Give back twice in a row: one worker, the second press finds the ticket already back', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await giveBack($, w)
    w.movedFirst = true
    await w.clock.advance(60_000)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(2)
    expect(spawned(w)).toHaveLength(1)
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
    const spawn = w.ran.findIndex(r => r[1] === SPAWN)
    expect(w.ran[close]?.[3]).toBe('w1:t9')
    expect(close).toBeLessThan(spawn)
    expect(w.ran[spawn]?.slice(0, 8)).toEqual(['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'worker'])
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toStartWith('cockpit: given back to a new worker - old worker\'s prompt cache cold for 10 min; ')
  })
  test('a session without a messaging socket starts the new worker with no parent: it asks in its own tab', async ($, on) => {
    const w = world(on, { herdr: true, env: { HOME: '/home/a', HERDR_ENV: '1', HERDR_BIN_PATH: 'herdr', HERDR_PANE_ID: 'w1:p1' } })
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'bash', SPAWN)[0]?.slice(0, 7)).toEqual(['bash', SPAWN, '--no-worktree', '--no-parent', 'mg677r', 'MG677R', 'worker'])
  })
  test('a socket whose file is gone is no parent: the new worker still starts, asking in its own tab', async ($, on) => {
    const w = world(on, { herdr: true, socket: false })
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'bash', SPAWN)[0]?.slice(0, 7)).toEqual(['bash', SPAWN, '--no-worktree', '--no-parent', 'mg677r', 'MG677R', 'worker'])
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
    expect(ran(w, 'bash', SPAWN)[0]?.slice(0, 8)).toEqual(['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'worker'])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toStartWith('cockpit: given back to a new worker - tab alive, cache warm for 40 more min, context 80k, but its pane took no text; ')
  })
  // W80HNH: a tab list Herdr did not give is not "no tab open": the old worker
  // may still run, and a new one beside it is the second worker Give back must
  // never make.
  test('a failed herdr tab list starts no worker, types nothing and says the tab state was unreadable', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 70 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 70 * 60_000,
    } } })
    w.tabsFail = 'herdr: timed out'
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(1)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(ran(w, 'herdr', 'tab', 'close')).toEqual([])
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to nobody (no worker started) - herdr tab list failed: herdr: timed out; the tab state could not be read, so no worker started and nothing typed')
    expect(w.toasts.at(-1)).toBe("MG677R back in in-progress; Herdr's tab list could not be read, so it is unknown whether its worker still runs - no worker started. Check its tab, and start one only if none runs.")
  })
  // 8NWE4P: under Herdr App a terminal that does not answer `hc get` is left
  // out of the list, and herdr-hc.sh says so with exit 3. The worker's tab may
  // be the one missing: the tabs that did come back prove nothing about it.
  test('a tab list missing a terminal that gave no answer starts no worker either', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 70 * 60_000, tokens: 80_000, ttl: 3_600_000, updatedAt: 1_000_000 - 70 * 60_000,
    } } })
    w.tabs = [{ tab_id: 't1', label: 'AB12CD/worker', agent_status: 'idle' }]
    w.tabsIncomplete = 'herdr-hc.sh: tab list incomplete: hc get gave no answer for t7'
    await startSession($)
    await giveBack($, w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(1)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(ran(w, 'herdr', 'pane', 'send-text')).toEqual([])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to nobody (no worker started) - herdr tab list failed: herdr-hc.sh: tab list incomplete: hc get gave no answer for t7; the tab state could not be read, so no worker started and nothing typed')
    expect(w.toasts.at(-1)).toBe("MG677R back in in-progress; Herdr's tab list could not be read, so it is unknown whether its worker still runs - no worker started. Check its tab, and start one only if none runs.")
  })
})

// Of several branches that hold the ticket's newest commit, the ticket's own
// wins — its worktree's, else the one named after it — so a branch stacked on
// top of it never takes the work, whatever order git lists them in (QTQW15).
// Give back takes the ticket's own worktree, else its <prefix>/<slug>, else
// the one branch holding its own commit — and guesses at nothing else
// (Alex, QTQW15, after three rounds of guessing rules).
describe('the branch a ticket rides on', () => {
  const cases: [string, Parameters<typeof carrierFor>[0], string | undefined][] = [
    ['its worktree, whatever holds the commit', { own: 'feat/qtqw15', base: 'master', holders: 'feat/0gmyqb\nfix-login\n' }, 'feat/qtqw15'],
    ['its worktree on a branch of another name', { own: 'fix-login', base: 'master' }, 'fix-login'],
    ['no worktree, its branch left behind', { base: 'master', slugBranches: 'review/feat/qtqw15\nfeat/qtqw15\n', holders: 'feat/0gmyqb\n' }, 'feat/qtqw15'],
    ['one holder of another name', { base: 'master', holders: 'fix-login\n' }, 'fix-login'],
    ['one holder and its review wrapper', { base: 'master', holders: 'review/fix-login\nfix-login\n' }, 'fix-login'],
    ['the branch the main checkout stands on holds it alone', { base: 'master', holders: 'fix-login\nmaster\n' }, 'fix-login'],
    ['a detached HEAD is no holder', { base: 'master', holders: '(HEAD detached at c2)\nfix-login\n' }, 'fix-login'],
    ['nothing at all', { base: 'master', holders: '', slugBranches: '' }, undefined],
    ['only the default branch holds it', { base: 'master', holders: 'master\n' }, undefined],
  ]
  for (const [name, input, want] of cases) test(name, () => {
    const got = carrierFor(input)
    expect(got !== undefined && 'branch' in got ? got.branch : undefined).toBe(want)
  })
  test('two branches named after it, feat/ and fix/, nothing picked', () => {
    const got = carrierFor({ base: 'master', slugBranches: 'feat/qtqw15\nfix/qtqw15\nreview/feat/qtqw15\n', holders: 'feat/qtqw15\n' })
    expect(got !== undefined && 'ambiguous' in got ? got.ambiguous : undefined).toEqual(['feat/qtqw15', 'fix/qtqw15'])
  })
  for (const holders of ['feat/0gmyqb\nfix-login\n', 'fix-login\nfeat/0gmyqb\n', 'fix-login\nreview/fix-login\nwip\n']) {
    test(`several holders, none its own, nothing picked: ${holders.trim().replaceAll('\n', ', ')}`, () => {
      const got = carrierFor({ base: 'master', holders })
      expect(got !== undefined && 'ambiguous' in got ? got.ambiguous : undefined).toEqual(holders.trim().split('\n').filter(b => !b.startsWith('review/')))
    })
  }
})

describe('the commit that is a ticket\'s own', () => {
  const log = (...subjects: string[]) => subjects.map((s, i) => `c${i} ${s}`).join('\n')
  const cases: [string, string, string | undefined][] = [
    ['its own scope', log('fix(QTQW15): the key'), 'c0'],
    ['breaking mark', log('feat(QTQW15)!: the key'), 'c0'],
    ['a newer follow-up listing it is passed by', log('chore(0GMYQB,QTQW15): stack', 'fix(QTQW15): the key'), 'c1'],
    ['a list with spaces is passed by', log('fix(QTQW15, 0GMYQB): both', 'fix(QTQW15): own'), 'c1'],
    ['no commit of its own: a list naming it is still not its own', log('fix(0GMYQB): other', 'chore(0GMYQB,QTQW15): both'), undefined],
    ['the handle in the text only', log('fix(0GMYQB): builds on QTQW15', 'feat(QTQW15): own'), 'c1'],
    ['another handle that contains it', log('fix(XQTQW15): no'), undefined],
    ['no scope at all', log('QTQW15: no scope'), undefined],
    ['nothing names it', log('fix(0GMYQB): other'), undefined],
  ]
  for (const [name, input, want] of cases) test(name, () => expect(commitNaming(input, 'QTQW15')).toBe(want))
})

describe('the jaira it runs beside', () => {
  test('one older than the cockpit needs is said once, with the version installed and the one needed', async ($, on) => {
    const w = world(on)
    w.version = '0.3.6'
    await startSession($)
    await answer($)
    await answer($)
    expect(ran(w, 'jaira', '--version')).toHaveLength(1)
    expect(w.toasts.filter(x => x.includes('too old'))).toEqual([
      'jaira 0.3.6 is too old for this cockpit: it needs 0.3.7 or later. Give back will be refused until you update jaira.',
    ])
  })

  for (const version of ['0.3.7', '0.10.0', 'dev']) {
    test(`jaira ${version} is not warned about`, async ($, on) => {
      const w = world(on)
      w.version = version
      await startSession($)
      expect(ran(w, 'jaira', '--version')).toHaveLength(1)
      expect(w.toasts.filter(x => x.includes('too old'))).toEqual([])
    })
  }

  test('without a board jaira is not asked', async ($, on) => {
    const w = world(on, { board: false })
    w.version = '0.3.6'
    await startSession($)
    expect(ran(w, 'jaira', '--version')).toEqual([])
    expect(w.toasts).toEqual([])
  })
})

describe('a newer jaira release', () => {
  const UPGRADE = (title: string) =>
    ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-upgrade', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never
  const TITLE = 'Update jaira 0.5.1 -> 0.6.0?'
  // The release cache as the hourly check leaves it.
  const released = (w: World, latest: string, at = CACHE) => w.files.set(at, JSON.stringify({ checked_at: '2026-10-16T09:00:00Z', latest }))
  // Every jaira command that could write: the binary and the setup.
  const writes = (w: World) => w.ran.filter(r => r[0] === 'jaira' && (r[1] === 'self' || r[1] === 'update')).map(r => r.slice(1, -1).join(' '))
  const offered = async ($: Test$) => (await $.ui.mount(BAND)).find({ key: 'release-update' })

  test('notice shown: a release newer than the jaira installed gets a line with Update and Dismiss', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    released(w, '0.6.0')
    await startSession($)
    const ui = await $.ui.mount(BAND)
    expect(await ui.find({ type: 'Text', text: /^jaira 0\.6\.0 is released \(this is 0\.5\.1\)/ })).toBeDefined()
    expect((await ui.find({ key: 'release-update' }))?.text).toBe('[ Update ]')
    expect(await ui.find({ key: 'release-dismiss' })).toBeDefined()
    // The cache is read, nothing more: no check of its own, nothing written.
    expect(writes(w)).toEqual([])
  })

  for (const [name, version, cache] of [
    ['the same version', '0.6.0', '0.6.0'],
    ['a cache older than the jaira installed', '0.6.0', '0.5.1'],
    ['a cache never written', '0.5.1', undefined],
    ['a dev build', 'dev', '0.6.0'],
  ] as const) {
    test(`no notice: ${name}`, async ($, on) => {
      const w = world(on)
      w.version = version
      if (cache !== undefined) released(w, cache)
      await startSession($)
      expect(await offered($)).toBeUndefined()
      expect(await (await $.ui.mount(BAND)).find({ type: 'Text', text: /is released/ })).toBeUndefined()
    })
  }

  test('no notice: JAIRA_NO_UPDATE_CHECK=1, even with a cache naming a newer release', async ($, on) => {
    const w = world(on, { env: { HOME: '/home/a', JAIRA_NO_UPDATE_CHECK: '1' } })
    w.version = '0.5.1'
    released(w, '0.6.0')
    await startSession($)
    expect(await offered($)).toBeUndefined()
  })

  test('the cache is read where JAIRA_HOME puts it', async ($, on) => {
    const w = world(on, { env: { HOME: '/home/a', JAIRA_HOME: '/srv/jaira' } })
    w.version = '0.5.1'
    released(w, '0.6.0', '/srv/jaira/update-check.json')
    await startSession($)
    expect(await offered($)).toBeDefined()
  })

  test('Update path with preview: the check and the dry run first, shown in a pane; nothing is written before the yes', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    w.latest = '0.6.0'
    released(w, '0.6.0')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'release-update' })
    expect(writes(w)).toEqual(['self upgrade --check --json', 'update --roles --lanes --dry-run'])
    expect(w.opened).toEqual([{ id: 'cockpit-upgrade', title: TITLE, focus: true }])
    const pane = await $.ui.mount(UPGRADE(TITLE))
    expect(await pane.find({ type: 'Text', text: /^Nothing has been written yet\. Update runs 'jaira self upgrade' \(0\.5\.1 -> 0\.6\.0\), then 'jaira update --roles --lanes'/ })).toBeDefined()
    // The dry run's own words, the older-board warning taken off.
    expect(await pane.find({ type: 'Text', text: 'skills in /repo/.claude/skills: 3 would be updated\nlanes: critique would be updated\nnothing was written (--dry-run)' })).toBeDefined()
    expect(await pane.find({ type: 'Text', text: /older version/ })).toBeUndefined()
    expect((await pane.find({ key: 'upgrade-confirm' }))?.props.label).toBe('Update to 0.6.0')
    // An Enter pressed once too often lands on Cancel.
    expect((await pane.find({ key: 'upgrade-cancel' }))?.props.autoFocus).toBe(true)
    await pane.press({ key: 'upgrade-cancel' })
    expect(w.closed).toEqual(['cockpit-upgrade'])
    expect(writes(w)).toEqual(['self upgrade --check --json', 'update --roles --lanes --dry-run'])
    // Cancelled, the offer stands.
    expect(await offered($)).toBeDefined()
  })

  test('the yes in the pane upgrades the binary, then the skills and lane files, and the notice goes', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    w.latest = '0.6.0'
    released(w, '0.6.0')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'release-update' })
    await (await $.ui.mount(UPGRADE(TITLE))).press({ key: 'upgrade-confirm' })
    expect(writes(w)).toEqual(['self upgrade --check --json', 'update --roles --lanes --dry-run', 'self upgrade', 'update --roles --lanes'])
    expect(w.closed).toContain('cockpit-upgrade')
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: 'jaira updated to 0.6.0', focus: true })
    const result = await $.ui.mount(MESSAGE('jaira updated to 0.6.0'))
    expect(await result.find({ type: 'Text', text: /^upgraded jaira 0\.5\.1 -> 0\.6\.0[\s\S]*lanes: critique updated$/ })).toBeDefined()
    // What self upgrade says on stderr is shown too: a release host override
    // must never pass quietly.
    expect(await result.find({ type: 'Text', text: /\nwarning: release host overridden by JAIRA_RELEASE_API=http:\/\/127\.0\.0\.1:8765\/api\n/ })).toBeDefined()
    // jaira is asked again which it is: 0.6.0 now, nothing left to offer.
    expect(await offered($)).toBeUndefined()
  })

  test('a dry run that gives no preview is shown as the error it is, and no update is offered', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    w.latest = '0.6.0'
    w.dryRunExit = 2
    released(w, '0.6.0')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'release-update' })
    expect(writes(w)).toEqual(['self upgrade --check --json', 'update --roles --lanes --dry-run'])
    expect(w.opened).toEqual([{ id: 'cockpit-message', title: 'Update jaira failed', focus: true }])
    const pane = await $.ui.mount(MESSAGE('Update jaira failed'))
    expect(await pane.find({ type: 'Text', text: "Nothing was written and no update is offered: 'jaira update --roles --lanes --dry-run' gave no preview. It said:\njaira: unknown flag: --roles" })).toBeDefined()
    expect(await pane.find({ key: 'upgrade-confirm' })).toBeUndefined()
  })

  test('a dry run that would leave an edited file alone is a preview, and the update is offered', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    w.latest = '0.6.0'
    w.dryRunExit = 3
    released(w, '0.6.0')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'release-update' })
    expect(w.opened).toEqual([{ id: 'cockpit-upgrade', title: TITLE, focus: true }])
    const pane = await $.ui.mount(UPGRADE(TITLE))
    expect(await pane.find({ type: 'Text', text: /^skipped \/repo\/\.claude\/skills\/jaira-worker\/SKILL\.md \(edited\)[\s\S]*some files were edited here and left alone/ })).toBeDefined()
    expect(await pane.find({ key: 'upgrade-confirm' })).toBeDefined()
  })

  test('an upgrade jaira refuses is shown whole, and the setup is not touched', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    w.latest = '0.6.0'
    w.upgradeFails = true
    released(w, '0.6.0')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'release-update' })
    await (await $.ui.mount(UPGRADE(TITLE))).press({ key: 'upgrade-confirm' })
    expect(writes(w)).toEqual(['self upgrade --check --json', 'update --roles --lanes --dry-run', 'self upgrade'])
    const pane = await $.ui.mount(MESSAGE('Update jaira failed'))
    expect(await pane.find({ type: 'Text', text: /installed with go install/ })).toBeDefined()
    expect(await offered($)).toBeDefined()
  })

  test('a check that finds nothing newer opens no pane and takes the notice away', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    released(w, '0.6.0')
    await startSession($)
    // Upgraded in a terminal meanwhile.
    w.version = '0.6.0'
    w.latest = '0.6.0'
    await (await $.ui.mount(BAND)).press({ key: 'release-update' })
    expect(writes(w)).toEqual(['self upgrade --check --json'])
    expect(w.opened).toEqual([])
    expect(w.toasts).toContain('jaira 0.6.0 is already the latest release')
    expect(await offered($)).toBeUndefined()
  })

  test('Dismiss hides exactly that version; the next release is offered again', async ($, on) => {
    const w = world(on)
    w.version = '0.5.1'
    released(w, '0.6.0')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'release-dismiss' })
    expect(w.store.get('dismissed-release')).toBe('0.6.0')
    expect(w.toasts).toContain('jaira 0.6.0 dismissed: the band offers the next release again')
    expect(await offered($)).toBeUndefined()
    // A minute on the cache is read again: the same version stays hidden.
    await w.clock.advance(60_000)
    await w.clock.settle()
    expect(await offered($)).toBeUndefined()
    released(w, '0.6.1')
    await w.clock.advance(60_000)
    await w.clock.settle()
    expect(await (await $.ui.mount(BAND)).find({ type: 'Text', text: /^jaira 0\.6\.1 is released \(this is 0\.5\.1\)/ })).toBeDefined()
    expect(writes(w)).toEqual([])
  })

  test('a dismissed version stays hidden in the next session', async ($, on) => {
    const w = world(on, { store: { 'dismissed-release': '0.6.0' } })
    w.version = '0.5.1'
    released(w, '0.6.0')
    await startSession($)
    expect(await offered($)).toBeUndefined()
  })
})

describe('give back with the review notes', () => {
  const GAPS = 'blocker: DoD 2 - an empty export still writes a header\nnote: the export test covers 1MB, not 40MB\nnote: pre-existing: the log rotates at midnight UTC'
  function withGaps(on: On, gaps: string, extra: Record<string, unknown> = {}) {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps }, ...extra }, HUMAN, WORKING]
    return w
  }
  async function openPane($: Test$) {
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'acceptance-MG677R' })
    return $.ui.mount(ACCEPT('Acceptance of MG677R'))
  }
  type Mounted = Awaited<ReturnType<typeof openPane>>
  // Each note as '[choice] mark text': the choice its bold button names.
  const LABEL: Record<string, string> = { now: '[fix now]', later: '[follow-up]', accept: '[accept]' }
  async function rows(pane: Mounted) {
    const out: string[] = []
    // A note's own line opens with its mark; the marks inside it are its children's.
    const own = (await pane.findAll({ type: 'Text' })).map(x => x.text).filter(x => /^[🔴🟡⚪🔵] /u.test(x))
    for (let i = 0; ; i++) {
      if ((await pane.find({ key: `accept-MG677R-${i}` })) === undefined) return out
      let chosen = ''
      for (const c of ['now', 'later', 'accept']) if ((await pane.find({ key: `${c}-MG677R-${i}` }))?.text.startsWith('[•')) chosen = c
      out.push(`${LABEL[chosen] ?? '[?]'} ${own[i] ?? ''}`)
    }
  }
  const notes = (w: World) => ran(w, 'jaira', 'note', 'MG677R').map(r => r[3])

  test('every note of the review is a row; a blocker starts as fix now, every other note as accepted - no follow-up nobody chose', async ($, on) => {
    withGaps(on, GAPS)
    await startSession($)
    const pane = await openPane($)
    expect(await rows(pane)).toEqual([
      '[fix now] 🔴 blocker: DoD 2 - an empty export still writes a header',
      '[accept] 🟡 note: the export test covers 1MB, not 40MB',
      '[accept] 🔵 note: pre-existing: the log rotates at midnight UTC',
    ])
    // The first row takes the keyboard; the field below is for anything else.
    expect((await pane.find({ key: 'accept-MG677R-0' }))?.props.autoFocus).toBe(true)
    expect((await pane.find({ key: 'typed-MG677R' }))?.props.autoFocus).toBeUndefined()
  })

  test('the Explain note classes a known limit, which starts as accepted, when it classes every note', async ($, on) => {
    const w = withGaps(on, GAPS)
    w.bodies.MG677R = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok | tests ok\n[blocker] the header\n[small edge] the size\n[known limit] the clock\n1) fix it [give back]\n'
    await startSession($)
    expect((await rows(await openPane($)))[2]).toBe('[accept] ⚪ note: pre-existing: the log rotates at midnight UTC')
  })

  test('a blocker: line stays a blocker whatever class the Explain note gives its position', async ($, on) => {
    const w = withGaps(on, GAPS)
    // Every note classed, so position counts - but the first is a known limit there.
    w.bodies.MG677R = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n[known limit] the header\n[small edge] the size\n[known limit] the clock\n'
    await startSession($)
    expect(await rows(await openPane($))).toEqual([
      '[fix now] 🔴 blocker: DoD 2 - an empty export still writes a header',
      '[accept] 🟡 note: the export test covers 1MB, not 40MB',
      '[accept] ⚪ note: pre-existing: the log rotates at midnight UTC',
    ])
  })

  test('an Explain note that classes fewer notes than the review wrote is not matched by position', async ($, on) => {
    const w = withGaps(on, GAPS)
    w.bodies.MG677R = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n[known limit] the clock\n'
    await startSession($)
    expect((await rows(await openPane($)))[0]).toStartWith('[fix now] 🔴')
  })

  test('a press changes a row: accept, fix now, follow-up, accept again', async ($, on) => {
    withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    const seen: string[] = []
    for (let i = 0; i < 4; i++) {
      seen.push((await rows(pane))[0]!.split(']')[0]! + ']')
      await cycle(pane, 'MG677R', 0)
      }
    expect(seen).toEqual(['[accept]', '[fix now]', '[follow-up]', '[accept]'])
  })

  test('Give back writes exactly the notes to fix, the follow-ups with their tickets and the accepted ones, after the move', async ($, on) => {
    const w = withGaps(on, GAPS)
    await startSession($)
    const pane = await openPane($)
    // The clock note: accept -> fix now -> follow-up.
    await cycle(pane, 'MG677R', 2)
    await cycle(pane, 'MG677R', 2)
    await pane.input({ key: 'typed-MG677R', text: '' })
    expect(notes(w)[0]).toBe([
      'Given back from signoff: beheben:',
      '- blocker: DoD 2 - an empty export still writes a header',
      'Folge-Ticket:',
      '- F00001: note: pre-existing: the log rotates at midnight UTC',
      'angenommen:',
      '- note: the export test covers 1MB, not 40MB',
    ].join('\n'))
    expect(ran(w, 'jaira', 'create')).toEqual([[
      'jaira', 'create', '--context=note: pre-existing: the log rotates at midnight UTC', '--follows=MG677R', '--json', '--', 'the log rotates at midnight UTC', '@/.worktrees/repo-mg677r',
    ]])
    // The move first: a follow-up only for a ticket that really went back.
    const at = (verb: string) => w.ran.findIndex(r => r[0] === 'jaira' && r[1] === verb)
    expect(at('move')).toBeLessThan(at('create'))
    expect(w.closed).toContain('cockpit-acceptance')
  })

  test("a follow-up's title drops the review's fix-now mark as it drops pre-existing", async ($, on) => {
    const w = withGaps(on, 'blocker: the header\nnote: fix-now: the guard misses an empty name')
    await startSession($)
    const pane = await openPane($)
    await cycle(pane, 'MG677R', 1)
    await cycle(pane, 'MG677R', 1)
    await pane.input({ key: 'typed-MG677R', text: '' })
    const create = ran(w, 'jaira', 'create')[0]!
    expect(create[create.indexOf('--') + 1]).toBe('the guard misses an empty name')
  })

  test('what the person types goes along under the notes, on one line', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await cycle(pane, 'MG677R', 0)
    await pane.input({ key: 'typed-MG677R', text: 'and the docs' })
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size\ndazu: and the docs')
  })

  test('nothing marked fix now and nothing typed: Enter in the field accepts nothing, Execute accepts', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await pane.input({ key: 'typed-MG677R', text: ' ' })
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(w.closed).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R: no note is marked fix now and nothing is typed - mark one, type what is missing, or press Execute to accept')
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(w.closed).toContain('cockpit-acceptance')
  })

  test('a follow-up that could not be created is named so, and the note is still written', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    w.createFails = true
    await startSession($)
    const pane = await openPane($)
    await cycle(pane, 'MG677R', 0)
    await cycle(pane, 'MG677R', 1)
    await cycle(pane, 'MG677R', 1)
    await pane.input({ key: 'typed-MG677R', text: '' })
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size\nFolge-Ticket:\n- (not created: jaira: the board is locked): note: pre-existing: the clock')
    // Heard, not only noted: the toast after it says the ticket went back.
    expect(w.opened).toContainEqual({ id: 'cockpit-message', title: 'Follow-up tickets of MG677R failed', focus: true })
  })

  test('every value of the follow-up goes with =, so a finding starting with - is no flag', async ($, on) => {
    const w = withGaps(on, '-race fails in core/gitref')
    await startSession($)
    const pane = await openPane($)
    await cycle(pane, 'MG677R', 0)
    await cycle(pane, 'MG677R', 0)
    await pane.input({ key: 'typed-MG677R', text: 'and the docs' })
    // jaira create takes the title only as its argument: after '--', where
    // jaira itself (checked by hand: exit 0) reads '-race' as a title.
    expect(ran(w, 'jaira', 'create')[0]?.slice(2, 7)).toEqual(['--context=-race fails in core/gitref', '--follows=MG677R', '--json', '--', '-race fails in core/gitref'])
  })

  test('a note kept a follow-up at the next acceptance keeps its ticket: no second one is created', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    w.bodies.MG677R = '## Progress\n- **2026-11-20 10:00 · Alex** — Given back from signoff: beheben:\n- note: the size\nFolge-Ticket:\n- A1B2C3: note: pre-existing: the clock\n'
    await startSession($)
    await (await openPane($)).input({ key: 'typed-MG677R', text: '' })
    expect(ran(w, 'jaira', 'create')).toEqual([])
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size\nFolge-Ticket:\n- A1B2C3: note: pre-existing: the clock')
  })

  test('words typed with no note to fix still open the beheben: list the worker reads', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await pane.input({ key: 'typed-MG677R', text: 'the docs are missing' })
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\ndazu: the docs are missing\nangenommen:\n- note: the size')
  })

  test('the dazu: line does not end the decisions read back after it', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    w.bodies.MG677R = '## Progress\n- **2026-11-20 10:00 · Alex** — Given back from signoff: beheben:\ndazu: the docs\nangenommen:\n- note: the size\n'
    await startSession($)
    expect(await rows(await openPane($))).toEqual(['[accept] 🟡 note: the size (decided before)'])
  })

  test('a note marked fix now: Execute gives the ticket back, accepts nothing', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await cycle(pane, 'MG677R', 0)
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'in-progress')).toHaveLength(1)
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toEqual([])
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size')
  })

  test('a note that is no blocker, not marked fix now, does not hold Accept: the decisions go onto the ticket once it is accepted', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    await startSession($)
    const pane = await openPane($)
    // the clock: accept -> fix now -> follow-up
    await cycle(pane, 'MG677R', 1)
    await cycle(pane, 'MG677R', 1)
    await pane.press({ key: 'execute-MG677R' })
    expect(notes(w)[0]).toBe('Accepted from signoff: Folge-Ticket:\n- F00001: note: pre-existing: the clock\nangenommen:\n- note: the size')
    const forced = w.ran.findIndex(r => r[1] === 'move' && r.includes('--force'))
    const noted = w.ran.findIndex(r => r[1] === 'note')
    expect(forced).toBeLessThan(noted)
    expect(forced).toBeLessThan(w.ran.findIndex(r => r[1] === 'create'))
  })

  test('an Accept the gates refuse writes no decisions and creates no follow-up', async ($, on) => {
    const w = withGaps(on, 'note: pre-existing: the clock')
    w.refusals.push({ code: 'dod_open', message: 'DoD 2 is not ticked' })
    await startSession($)
    await (await openPane($)).press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(ran(w, 'jaira', 'create')).toEqual([])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toEqual([])
    expect(w.submitted).toEqual([])
  })

  test('a blocker not marked fix now is accepted only with a reason, which goes onto the ticket as a decision', async ($, on) => {
    const w = withGaps(on, 'blocker: the header\nnote: the size')
    await startSession($)
    const pane = await openPane($)
    // Fix now: the field asks for anything else, Execute gives back.
    expect((await pane.find({ key: 'typed-MG677R' }))?.props.label).toBe('Anything else?')
    await cycle(pane, 'MG677R', 0)
    await cycle(pane, 'MG677R', 0)
    // Not fix now: the field asks why, and Execute without a reason does nothing.
    expect((await pane.find({ key: 'typed-MG677R' }))?.props.label).toBe('Why accept it with the blocker?')
    for (const empty of ['', '   ', '\n']) {
      await pane.input({ key: 'typed-MG677R', text: empty })
      expect(ran(w, 'jaira', 'move')).toEqual([])
      expect(w.closed).toEqual([])
      expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain('MG677R: a blocker is not marked fix - type why it is accepted anyway, or mark it fix')
    }
    await pane.input({ key: 'typed-MG677R', text: 'the header is\nharmless for now' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(notes(w)).toEqual([
      'Accepted from signoff: angenommen:\n- blocker: the header\n- note: the size',
      'decision: accepted with a blocker not fixed: the header is harmless for now',
      'cockpit: accepted by the person, signoff -> done',
    ])
  })

  test('Execute accepts at once without a blocker, the pane closing behind it', async ($, on) => {
    const w = withGaps(on, 'note: the header')
    w.bodies.MG677R = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n[small edge] the header\n'
    await startSession($)
    await (await openPane($)).press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(w.closed).toContain('cockpit-acceptance')
  })

  // A blocker by its prefix, in markdown, or by the Explain note's class;
  // a 'blocker:' line the Explain note calls a small edge stays a blocker.
  for (const [gaps, explain] of [
    ['blocker: the header', ''],
    ['blocker: the header', '[small edge] the header'],
    ['- **Blocker:** the header\nnote: the size', ''],
    ['note: the header\nnote: the size', '[blocker] the header\n[small edge] the size'],
  ] as const) {
    test(`a blocker in the review starts as fix now: Execute gives back and accepts nothing: ${JSON.stringify(gaps)}${explain === '' ? '' : ` (Explain: ${JSON.stringify(explain)})`}`, async ($, on) => {
      const w = withGaps(on, gaps)
      if (explain !== '') w.bodies.MG677R = `## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n${explain}\n`
      await startSession($)
      await (await openPane($)).press({ key: 'execute-MG677R' })
      expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toEqual([])
      expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'in-progress')).toHaveLength(1)
    })
  }

  for (const press of ['back', 'accept'] as const) {
    test(`a move another session made first: the second press writes nothing (${press})`, async ($, on) => {
      const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
      w.movedFirst = true
      await startSession($)
      {
        const pane = await openPane($)
        await cycle(pane, 'MG677R', 0)
        await cycle(pane, 'MG677R', 1)
        await cycle(pane, 'MG677R', 1)
        if (press === 'back') await pane.input({ key: 'typed-MG677R', text: 'and the docs' })
        else {
          await cycle(pane, 'MG677R', 0)
          await cycle(pane, 'MG677R', 0)
          await pane.press({ key: 'execute-MG677R' })
        }
      }
      expect(ran(w, 'jaira', 'move').filter(r => !r.includes('--dry-run'))).toHaveLength(1)
      expect(ran(w, 'jaira', 'create')).toEqual([])
      expect(ran(w, 'jaira', 'note')).toEqual([])
      expect(w.ran.filter(r => r[0] === 'herdr' || r[1] === SPAWN)).toEqual([])
      expect(w.opened.at(-1)?.title).toBe(`${press === 'back' ? 'Give back' : 'Accept'} MG677R failed`)
      expect(w.submitted).toEqual([])
    })
  }

  test('the next acceptance starts every note decided before as it was decided, follow-ups included', async ($, on) => {
    const w = withGaps(on, GAPS)
    w.bodies.MG677R = [
      '## Progress',
      '- **2026-11-20 10:00 · Alex** — Given back from signoff: beheben:',
      '- blocker: DoD 2 - an empty export still writes a header',
      'Folge-Ticket:',
      '- A1B2C3: note: the export test covers 1MB, not 40MB',
      'angenommen:',
      '- note: pre-existing: the log rotates at midnight UTC',
      '- **2026-11-20 11:00 · Alex** — cockpit: given back to a new worker - x',
    ].join('\n')
    await startSession($)
    expect(await rows(await openPane($))).toEqual([
      '[fix now] 🔴 blocker: DoD 2 - an empty export still writes a header (decided before)',
      '[follow-up] 🟡 note: the export test covers 1MB, not 40MB (decided before)',
      '[accept] 🔵 note: pre-existing: the log rotates at midnight UTC (decided before)',
    ])
  })

  test('a newer decision wins over an older one; a note given back with words only decides nothing', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: the docs')
    w.bodies.MG677R = [
      '## Progress',
      '- **2026-11-19 10:00 · Alex** — Given back from signoff: beheben:',
      '- note: the size',
      '- **2026-11-19 12:00 · Alex** — Given back from signoff: note: the docs',
      '- **2026-11-20 10:00 · Alex** — Accepted from signoff: angenommen:',
      '- note: the size',
    ].join('\n')
    await startSession($)
    expect(await rows(await openPane($))).toEqual(['[accept] 🟡 note: the size (decided before)', '[accept] 🟡 note: the docs'])
  })

  // The input classes of the ticket, one row each.
  test('no notes: the pane asks what is missing, as before', async ($, on) => {
    withGaps(on, 'none')
    await startSession($)
    const pane = await openPane($)
    expect(await rows(pane)).toEqual([])
    expect((await pane.find({ key: 'typed-MG677R' }))?.props.autoFocus).toBe(true)
    expect((await pane.find({ key: 'typed-MG677R' }))?.props.label).toBe('What is missing?')
  })

  test('a note wrapped over lines is one row, a blank line and a none dropped, CRLF and emphasis read', async ($, on) => {
    withGaps(on, '- note: a long finding\n  that wraps here\n- **blocker:** the header\r\n\nnone')
    await startSession($)
    // What stops the ticket first: the blocker on top.
    expect(await rows(await openPane($))).toEqual(['[fix now] 🔴 **blocker:** the header', '[accept] 🟡 note: a long finding that wraps here'])
  })

  test('gaps without prefixes, as an older review wrote them, are a row per line, none merged', async ($, on) => {
    withGaps(on, 'first finding\nsecond finding')
    await startSession($)
    expect(await rows(await openPane($))).toEqual(['[accept] 🟡 first finding', '[accept] 🟡 second finding'])
  })

  test('a review of an earlier round lists its own notes, not what a later lane wrote into review-gaps', async ($, on) => {
    withGaps(on, 'note: optimize says this', { round: 2, review_round: 1, review: { verdict: SIGNOFF.review.verdict, gaps: 'note: optimize says this', verdict_gaps: 'blocker: the review said this' } })
    await startSession($)
    expect(await rows(await openPane($))).toEqual(['[fix now] 🔴 blocker: the review said this'])
  })

  test('very many notes are all rows', async ($, on) => {
    withGaps(on, Array.from({ length: 40 }, (_, i) => `note: finding ${i}`).join('\n'))
    await startSession($)
    const all = await rows(await openPane($))
    expect(all).toHaveLength(40)
    expect(all[39]).toBe('[accept] 🟡 note: finding 39')
  })

  test('decided meanwhile by someone else: nothing is moved, created or noted, and the pane says who was faster', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    await startSession($)
    const pane = await openPane($)
    w.tickets = w.tickets.map(x => ((x as { handle: string }).handle === 'MG677R' ? { ...(x as object), status: 'in-progress' } : x))
    await cycle(pane, 'MG677R', 0)
    await pane.input({ key: 'typed-MG677R', text: '' })
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(ran(w, 'jaira', 'create')).toEqual([])
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: 'Give back MG677R failed', focus: true })
  })

  test('the old worker gets the decisions on one line, the ticket keeps the lines', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', repo: '/repo', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 60_000, tokens: 80_000, window: 1_000_000, ttl: 3_600_000, updatedAt: 1_000_000 - 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps: 'note: the size\nnote: the docs' } }]
    await startSession($)
    const pane = await openPane($)
    await cycle(pane, 'MG677R', 0)
    await cycle(pane, 'MG677R', 1)
    const typed = pane.input({ key: 'typed-MG677R', text: '' })
    await w.clock.advance(1000)
    await typed
    const sent = ran(w, 'herdr', 'pane', 'send-text', 'w1:p9')[0]?.[4] ?? ''
    expect(sent).toContain('Reason: beheben: - note: the size - note: the docs - also a note')
    expect(sent).not.toContain('\n')
  })
})

// 409XCV: the teamlead's session hears of a press in the same minute, as one
// line, and only of a move that went through.
describe('the line a press hands this session', () => {
  test('Accept: one line, with the pull request the press started, so the teamlead starts no second one', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await startSession($)
    await pressAccept($, await $.ui.mount(BAND), 'MG677R', w)
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R: accepted, signoff -> done (Person, Cockpit); MG677R: opening the pull request for feat/mg677r in a Herdr tab'])
  })

  test('Accept whose next step fails: still one line, since the ticket is in done', async ($, on) => {
    const w = world(on)
    w.ahead = 3
    await startSession($)
    await pressAccept($, await $.ui.mount(BAND), 'MG677R', w)
    expect(w.submitted).toEqual(['jaira-cockpit: MG677R: accepted, signoff -> done (Person, Cockpit); the next step after MG677R failed, see the cockpit'])
  })

  test('Give back: one line naming who carries the work now, the reason on one line', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', repo: '/repo', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 300_000, window: 500_000, ttl: 3_600_000, updatedAt: 1_000_000 - 20 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'acceptance-MG677R' })
    await (await $.ui.mount(ACCEPT('Acceptance of MG677R'))).input({ key: 'typed-MG677R', text: 'der Branch\nfehlt   im Schluessel ' + 'x'.repeat(200) })
    expect(w.submitted).toHaveLength(1)
    const line = w.submitted[0]!
    expect(line.startsWith('jaira-cockpit: MG677R given back: der Branch fehlt im Schluessel xx')).toBe(true)
    expect(line).not.toContain('\n')
    expect(line.endsWith('…, worker new MG677R/worker (pane w1:p42, on feat/mg677r) (Person, Cockpit; signoff -> in-progress)')).toBe(true)
  })

  test('a reason cut at its limit never splits a character outside the BMP', async ($, on) => {
    const w = world(on)
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'acceptance-MG677R' })
    // 118 characters, then emoji: the cut at 119 falls inside a surrogate pair by UTF-16 index.
    await (await $.ui.mount(ACCEPT('Acceptance of MG677R'))).input({ key: 'typed-MG677R', text: 'x'.repeat(118) + '\u{1F44D}\u{1F44D}\u{1F44D} passt nicht' })
    const line = w.submitted[0]!
    expect(line).toContain('x\u{1F44D}…, worker')
    expect(/[\uD800-\uDBFF](?![\uDC00-\uDFFF])|(?<![\uD800-\uDBFF])[\uDC00-\uDFFF]/.test(line)).toBe(false)
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
    await ui.press({ key: 'acceptance-MG677R' })
    await (await $.ui.mount(ACCEPT('Acceptance of MG677R'))).input({ key: 'typed-MG677R', text: 'fehlt' })
    expect(await ui.find({ type: 'Text', text: "MG677R back in in-progress; no worker started (the ticket's note says why)" })).toBeDefined()
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to nobody (no worker started) - no Herdr in this session')
    expect(ran(w, 'bash', SPAWN)).toEqual([])
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
    await ui.press({ key: 'acceptance-MG677R' })
    await (await $.ui.mount(ACCEPT('Acceptance of MG677R'))).input({ key: 'typed-MG677R', text: 'fehlt' })
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

  test('a worker held at an approval dialog says so instead of a next round', async ($, on) => {
    const w = world(on, { herdr: true, store: {
      'session:a': record('a', 'T8QS1Y', { lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 120_000, updatedAt: 1_000_000 - 20 * 60_000 }),
    } })
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'blocked' }]
    await startSession($)
    const text = (await $.command.run(typed('workers'))).text
    expect(text).toContain('T8QS1Y · in-progress · worker · context 120k of 1000k · cache 40 min left · waits for approval: answer the dialog in its tab')
    expect(String((await $.tool.call({ tool: 'mcp__jaira-cockpit__workers' })).result)).toContain('waits for approval: answer the dialog in its tab')
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

  test('an agent held at an approval dialog is still in its turn and counts', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'blocked' },
      { tab_id: 'w1:t2', label: 'MG677R/review', agent_status: 'blocked' },
    ]
    await startSession($)
    expect(await headLine($)).toBe('1 worker | 1 review | 2 waiting for you | 5h 34%')
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
    expect(await ui.find({ key: 'acceptance-MG677R' })).toBeDefined()
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

const ASKS_PANE = { plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-asks', props: { title: 'Waiting for you', isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 30 }, view: {} } } as never

describe('open questions with their ticket and options (1M3X77)', () => {
  const filing = (o: Record<string, unknown>) => parseFiling(JSON.stringify({ thread: null, title: 'T', line: 'l', ticket: null, question: 'Q?', answered: [], ...o }))!
  // Every shape the cheap model may hand back for the two new fields.
  for (const [name, raw, want] of [
    ['named options', { options: ['ja', 'nein'] }, ['ja', 'nein']],
    ['no options key', {}, []],
    ['options not a list', { options: 'ja, nein' }, []],
    ['null in the list, blanks, numbers', { options: [null, ' ', 3, ' ja '] }, ['ja']],
    ['more than nine', { options: Array.from({ length: 12 }, (_, i) => `o${i}`) }, Array.from({ length: 9 }, (_, i) => `o${i}`)],
  ] as const) {
    test(`options: ${name}`, () => {
      expect(filing(raw).options).toEqual(want)
    })
  }
  for (const [name, about, want] of [
    ['an id', 'VTSF24', 'VTSF24'],
    ['lower case', 'vtsf24', 'VTSF24'],
    ['too long', 'VTSF245', null],
    ['prose', 'the ticket', null],
    ['null', null, null],
    ['a number', 123456, null],
  ] as const) {
    test(`about: ${name}`, () => {
      expect(filing({ about }).about).toBe(want)
    })
  }

  const f = (o: Partial<Filing>): Filing => ({ thread: null, title: 'Release', line: '', ticket: null, question: null, about: null, options: [], answered: [], ...o })
  test('the ticket and options ride with the question and go with it', () => {
    let ts = file([], f({ question: 'Antworte ja zu VTSF24?', about: 'VTSF24', options: ['ja', 'nein'] }), 1)
    expect(untitled(ts)).toEqual(['VTSF24'])
    ts = titled(ts, 'VTSF24', 'Release-Knopf')
    expect(untitled(ts)).toEqual([])
    expect(asksOf(ts)).toEqual([{ thread: 1, ticket: 'VTSF24', title: 'Release-Knopf', question: 'Antworte ja zu VTSF24?', options: ['ja', 'nein'] }])
    // A new question on the same thread replaces ticket, title and options.
    const again = file(ts, f({ thread: 1, question: 'Welche Farbe?' }), 2)
    expect(asksOf(again)).toEqual([{ thread: 1, title: 'Release', question: 'Welche Farbe?', options: [] }])
    // Answered by a filing, or by a pick: nothing of it is left.
    for (const done of [file(ts, f({ thread: 1, answered: [1] }), 2), answerThread(ts, 1, 2)]) {
      expect(asksOf(done)).toEqual([])
      expect(Object.keys(done[0]!).sort()).toEqual(['id', 'lastTurn', 'lines', 'title'])
    }
  })

  test('the thread\'s own ticket stands in where the question names none; asked last comes first', () => {
    const ts = file(file([], f({ title: 'A', ticket: 'AAAAAA', question: 'a?' }), 1), f({ title: 'B', question: 'b?' }), 2)
    expect(asksOf(ts).map(a => [a.ticket, a.title])).toEqual([[undefined, 'B'], ['AAAAAA', 'A']])
    expect(answerLine(asksOf(ts)[0]!, 'ja')).toBe('B: ja')
    expect(answerLine(asksOf(ts)[1]!, 'ja')).toBe('AAAAAA: ja')
  })

  test('the band opens the pane; a digit answers into the session and the question leaves', async ($, on) => {
    const w = world(on)
    w.shows['VTSF24'] = { title: 'Release-Knopf in Actions' }
    await startSession($)
    w.model = '{"thread": null, "title": "Release", "line": "Frage zu VTSF24", "ticket": null, "question": "Antworte ja zu VTSF24?", "about": "vtsf24", "options": ["ja", "nein, spaeter"], "answered": []}'
    await $.prompt.submit({ text: 'Was ist mit dem Release?', origin: { kind: 'composer' }, wait: false })
    await answer($, 'Naechster Schritt: antworte ja zu VTSF24')
    await w.clock.settle()
    const band = await $.ui.mount(BAND)
    const button = await band.find({ key: 'asks' })
    expect(button?.props.hotkey).toBe('q')
    await band.press({ key: 'asks' })
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-asks', title: 'Waiting for you', focus: true })
    const pane = await $.ui.mount(ASKS_PANE)
    const text = (await pane.find({ type: 'Box' }))!.text
    expect(text).toContain('VTSF24')
    expect(text).toContain('Release-Knopf in Actions')
    expect(text).toContain('Antworte ja zu VTSF24?')
    expect((await pane.find({ key: 'ask-1-0' }))?.props.hotkey).toBe('1')
    expect((await pane.find({ key: 'ask-1-1' }))?.props.hotkey).toBe('2')
    await pane.press({ key: 'ask-1-1' })
    await w.clock.settle()
    expect(w.submitted).toContain('VTSF24: nein, spaeter')
    expect(w.closed).toContain('cockpit-asks')
    expect(await (await $.ui.mount(BAND)).find({ key: 'asks' })).toBeUndefined()
    await pane.unmount()
    expect((await (await $.ui.mount(ASKS_PANE)).find({ type: 'Box' }))!.text).toContain('No question is waiting for you.')
  })

  test('a ticket jaira has no title for is asked about once, and the thread\'s title stands in', async ($, on) => {
    const w = world(on)
    await startSession($)
    w.model = '{"thread": null, "title": "Logo", "line": "Frage", "ticket": null, "question": "Wie gross?", "about": "ZZZZZZ", "answered": []}'
    for (let i = 0; i < 3; i++) {
      await $.prompt.submit({ text: `Logo ${i}?`, origin: { kind: 'composer' }, wait: false })
      await answer($, 'Wie gross?')
      await w.clock.settle()
      w.model = '{"thread": 1, "title": "Logo", "line": "noch offen", "ticket": null, "question": null, "answered": []}'
    }
    expect(ran(w, 'jaira', 'show', 'ZZZZZZ')).toHaveLength(1)
    expect((await (await $.ui.mount(ASKS_PANE)).find({ type: 'Box' }))!.text).toContain('ZZZZZZ | Logo')
  })

  test('one button holds the focus: the newest question\'s first', async ($, on) => {
    const w = world(on)
    await startSession($)
    w.model = '{"thread": null, "title": "A", "line": "a", "ticket": null, "question": "a?", "options": ["x", "y"], "answered": []}'
    await $.prompt.submit({ text: 'a', origin: { kind: 'composer' }, wait: false })
    await answer($, 'a?')
    await w.clock.settle()
    w.model = '{"thread": null, "title": "B", "line": "b", "ticket": null, "question": "b?", "answered": []}'
    await $.prompt.submit({ text: 'b', origin: { kind: 'composer' }, wait: false })
    await answer($, 'b?')
    await w.clock.settle()
    const pane = await $.ui.mount(ASKS_PANE)
    const focused = (await pane.findAll({ type: 'Button' })).filter(b => b.props.autoFocus === true).map(b => b.key)
    expect(focused).toEqual(['ask-2-chat'])
  })

  test('a question with no options is answered in the chat, the prompt started with its ticket', async ($, on) => {
    const w = world(on)
    await startSession($)
    w.model = '{"thread": null, "title": "Logo", "line": "Frage", "ticket": "ABCDEF", "question": "Wie gross?", "answered": []}'
    await $.prompt.submit({ text: 'Logo?', origin: { kind: 'composer' }, wait: false })
    await answer($, 'Wie gross?')
    await w.clock.settle()
    const pane = await $.ui.mount(ASKS_PANE)
    // The ticket is the thread's own, and jaira was not asked for a title.
    expect((await pane.find({ type: 'Box' }))!.text).toContain('ABCDEF')
    expect(ran(w, 'jaira', 'show', 'ABCDEF')).toEqual([])
    await pane.press({ key: 'ask-1-chat' })
    expect(w.filled).toEqual(['ABCDEF: '])
    expect(w.submitted).toEqual([])
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

  // Every way the start line arrives: typed into the pane, and relayed by
  // Herdr App from the teamlead's tab inside an <hc> wrapper.
  const starts = [
    ['typed', '/jaira-worker MG677R --parent jaira-4c', 'composer', ''],
    ['typed with leading space', '  /jaira-worker MG677R', 'composer', ''],
    ['relayed by Herdr App', '<hc from="t3">/jaira-worker MG677R --parent uds:/run/user/1000/cc-socks/1.sock</hc>', 'human', ''],
    ['relayed with attributes and a newline', '<hc from="t3" via="x">\n/jaira-dispatcher MG677R</hc>', 'human', ''],
    ['a lane relayed by Herdr App', '<hc from="t1">/jaira-role-lane MG677R critique</hc>', 'human', 'critique'],
    ['a lane whose id holds digits', '/jaira-role-lane MG677R qa2', 'composer', 'qa2'],
  ] as const
  for (const [name, text, kind, lane] of starts) {
    test(`recognises its start line ${name}, and a subagent's report reaches it whole`, async ($, on) => {
      const w = world(on, { herdr: true })
      await startSession($)
      await $.prompt.submit({ text, origin: { kind }, wait: false })
      expect(w.store.get('carrying:abcdef12')).toMatchObject({ ticket: 'MG677R', lane })
      const report = '<task-notification>Ticket MG677R: bughunt done at b5ec407 - 0 findings</task-notification>'
      expect((await $.prompt.submit({ text: report, origin: { kind: 'task-notification' }, wait: false })).text).toBe(report)
    })
  }

  test('an <hc> line that starts nothing leaves the session a person\'s', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await $.prompt.submit({ text: '<hc from="t3">how far is /jaira-worker MG677R?</hc>', origin: { kind: 'human' }, wait: false })
    expect(w.store.has('carrying:abcdef12')).toBe(false)
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

  test('/cockpit ticket reads the reports: the next report has nothing new, a later one is new', async ($, on) => {
    world(on)
    await startSession($)
    await $.session.receive({ origin: { kind: 'peer' }, text: 'Ticket MG677R: review -> signoff, review-verdict="approve"' })
    await $.command.run(typed('ticket MG677R'))
    const pane = await $.ui.mount({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-ticket', props: { title: 'Ticket MG677R', isFocused: false, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } } as never)
    expect(await pane.find({ type: 'Text', text: 'Ticket MG677R: review -> signoff, review-verdict="approve"' })).toBeDefined()
    expect((await $.command.run(typed('report MG677R'))).text).toStartWith('Nothing new for MG677R')
    await $.session.receive({ origin: { kind: 'peer' }, text: 'Ticket MG677R: signoff -> review' })
    expect((await $.command.run(typed('report MG677R'))).text).toContain('signoff -> review')
  })
})

const OVERVIEW = { plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-overview', props: { title: 'Board overview', isFocused: true, bodyColumns: 160, placement: 'dock', scroll: { offset: 0, bodyRows: 40 }, view: {} } } as never

// An instant `days` before the test clock's now, as jaira writes it.
const ago = (days: number) => new Date(1_000_000 - days * 86_400_000).toISOString()

describe('the overview', () => {
  test('the band has an Overview button only where there is a board', async ($, on) => {
    world(on)
    await startSession($)
    expect(await (await $.ui.mount(BAND)).find({ key: 'overview' })).toBeDefined()
  })

  test('no board, no Overview button', async ($, on) => {
    const w = world(on, { board: false })
    w.usage.tokens = 50_000
    await startSession($)
    await answer($)
    expect(await (await $.ui.mount(BAND)).find({ key: 'overview' })).toBeUndefined()
  })

  test('Overview reads the board afresh and shows in work, queue and put back, each row Handle | Lane | Titel', async ($, on) => {
    const w = world(on)
    w.tickets = [
      SIGNOFF,
      HUMAN,
      { ...WORKING, assignee: 'Alex', updated_at: ago(0.1) },
      // Not offered by next for an open blocked_by, and nobody on it: in nobody's work.
      { handle: 'W2W2W2', title: 'waits', status: 'in-progress', assignee: '', updated_at: ago(0.1) },
      { handle: 'Q1Q1Q1', title: 'queued', status: 'critique', assignee: '', updated_at: ago(0.2) },
      { handle: 'B1B1B1', title: 'later', status: 'backlog', creator: 'Berk', updated_by: 'Berk', updated_at: ago(20), created_at: ago(20) },
      { handle: 'X1X1X1', title: 'stuck', status: 'blocked', updated_by: 'Alex', blocked_reason: 'waits on Berk', updated_at: ago(15), created_at: ago(30) },
      { handle: 'D1D1D1', title: 'finished', status: 'done', updated_at: ago(0.1) },
    ]
    w.next = ['Q1Q1Q1', 'B1B1B1', 'K7P2QX']
    await startSession($)
    w.ran = []
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    expect(ran(w, 'jaira', 'next')).toEqual([['jaira', 'next', '--all', '--json', '@/repo']])
    expect(w.opened.at(-1)).toEqual({ id: 'cockpit-overview', title: 'Board overview', focus: true })
    const pane = await $.ui.mount(OVERVIEW)
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    // Queue is the open tab: one section beneath the other, each heading a fold button with its count.
    const labels = (await pane.findAll({ type: 'Button' } as never)).map(b => (b.props as { label: string }).label)
    expect(labels).toEqual(['Queue', 'Log', 'Workers', 'Close', '▾ In work (1)', '▾ Queue (1)', '▾ Put back (2)'])
    expect((await pane.findAll({ type: 'Box' })).filter(b => 'width' in b.props)).toEqual([])
    expect(texts).not.toContain('Created and put back, last 7 days')
    expect(texts).toContain('T8QS1Y | in-progress | filter')
    expect(texts).toContain('  Alex')
    expect(texts).toContain('Q1Q1Q1 | critique | queued')
    // The last touched first; a blocked ticket says why.
    expect(texts.indexOf('X1X1X1 | blocked | stuck')).toBeLessThan(texts.indexOf('B1B1B1 | backlog | later'))
    expect(texts).toContain('  Alex: waits on Berk')
    // A person's lanes are the band's, the terminal lane is nobody's.
    expect(texts.filter(x => /MG677R|K7P2QX|D1D1D1|W2W2W2/.test(x))).toEqual([])
    await pane.press({ key: 'overview-close' })
    expect(w.closed).toContain('cockpit-overview')
  })

  test('the history: created and put back over the last seven days, by day, newest first, with who', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new today', status: 'backlog', creator: 'Alex', updated_by: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'P1P1P1', title: 'put back', status: 'blocked', creator: 'Berk', updated_by: 'Alex', blocked_reason: 'x', created_at: ago(5), updated_at: ago(2) },
      { handle: 'O1O1O1', title: 'too old', status: 'backlog', creator: 'Berk', created_at: ago(10), updated_at: ago(10) },
      { handle: 'W1W1W1', title: 'started', status: 'in-progress', creator: 'Alex', assignee: 'Alex', created_at: ago(1), updated_at: ago(0.1) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    expect(texts).toContain('Created and put back, last 7 days')
    expect(texts.filter(x => /^In work|^Put back/.test(x))).toEqual([])
    const history = texts.slice(texts.indexOf('Created and put back, last 7 days') + 1)
    const days = history.filter(x => /^── \d{4}-\d\d-\d\d ──$/.test(x))
    expect(days).toHaveLength(4)
    expect([...days].sort().reverse()).toEqual(days)
    const events = history.filter(x => LOG_LINE.test(x)).map(x => x.replace(LOG_LINE, '$1 '))
    expect(events).toEqual([
      '🆕 N1N1N1 | backlog | new today | created by Alex',
      '🆕 W1W1W1 | in-progress | started | created by Alex',
      '🔙 P1P1P1 | blocked | put back | put back by Alex',
      '🆕 P1P1P1 | blocked | put back | created by Berk',
    ])
  })

  test('a title is never cut off: rows wrap', async ($, on) => {
    const w = world(on)
    const long = 'a title long enough that three columns of a third each would have cut it off well before its end'
    w.tickets = [{ ...WORKING, assignee: 'Alex', title: long, updated_at: ago(0.1) }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const row = await (await $.ui.mount(OVERVIEW)).find({ type: 'Text', text: `T8QS1Y | in-progress | ${long}` })
    expect(row?.props).toMatchObject({ wrap: 'wrap' })
  })

  test('a heading folds its section, and tab and folds outlive closing the pane', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...WORKING, assignee: 'Alex', updated_at: ago(0.1) }, { handle: 'Q1Q1Q1', title: 'queued', status: 'critique', assignee: '', updated_at: ago(0.2) }]
    w.next = ['Q1Q1Q1']
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    let pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-fold-working' })
    expect((await pane.find({ key: 'overview-fold-working' }))?.props.label).toBe('▸ In work (1)')
    expect(await pane.find({ type: 'Text', text: 'T8QS1Y | in-progress | filter' })).toBeUndefined()
    expect(await pane.find({ type: 'Text', text: 'Q1Q1Q1 | critique | queued' })).toBeDefined()
    await pane.press({ key: 'overview-tab-log' })
    await pane.press({ key: 'overview-close' })
    await pane.unmount()
    // Opened again: still on Log, and back on Queue the section is still folded.
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    pane = await $.ui.mount(OVERVIEW)
    expect(await pane.find({ type: 'Text', text: 'Created and put back, last 7 days' })).toBeDefined()
    await pane.press({ key: 'overview-tab-queue' })
    expect((await pane.find({ key: 'overview-fold-working' }))?.props.label).toBe('▸ In work (1)')
    await pane.press({ key: 'overview-fold-working' })
    expect(await pane.find({ type: 'Text', text: 'T8QS1Y | in-progress | filter' })).toBeDefined()
  })

  test('in work, queue and put back each carry their own colour, at the heading and at every row', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...WORKING, assignee: 'Alex', updated_at: ago(0.1) },
      { handle: 'Q1Q1Q1', title: 'queued', status: 'critique', assignee: '', updated_at: ago(0.2) },
      { handle: 'X1X1X1', title: 'stuck', status: 'blocked', updated_by: 'Alex', blocked_reason: 'waits', updated_at: ago(15), created_at: ago(30) },
    ]
    w.next = ['Q1Q1Q1']
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    const colorOf = async (text: string) => (await pane.find({ type: 'Text', text: new RegExp(`^${text.replace(/[()]/g, '\\$&')}$`) }))?.props.color
    const groups = [
      [await colorOf('▾ In work (1)'), await colorOf('T8QS1Y')],
      [await colorOf('▾ Queue (1)'), await colorOf('Q1Q1Q1')],
      [await colorOf('▾ Put back (1)'), await colorOf('X1X1X1')],
    ]
    for (const [heading, row] of groups) {
      expect(heading).toBeDefined()
      expect(row).toBe(heading)
    }
    expect(new Set(groups.map(g => g[0])).size).toBe(3)
  })

  test('a person has one colour, derived from the name: the same in both tabs and after reopening', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { ...WORKING, assignee: 'Alex', creator: 'Berk', created_at: ago(1), updated_at: ago(0.1) },
      { handle: 'X1X1X1', title: 'stuck', status: 'blocked', creator: 'Alex', updated_by: 'Berk', blocked_reason: 'waits', updated_at: ago(2), created_at: ago(3) },
    ]
    await startSession($)
    const colours = async (name: string) => {
      await (await $.ui.mount(BAND)).press({ key: 'overview' })
      const pane = await $.ui.mount(OVERVIEW)
      const found = [...await pane.findAll({ type: 'Text', text: new RegExp(`^${name}$`) })]
      await pane.press({ key: 'overview-tab-log' })
      found.push(...await pane.findAll({ type: 'Text', text: new RegExp(`^${name}$`) }))
      await pane.press({ key: 'overview-tab-queue' })
      await pane.press({ key: 'overview-close' })
      await pane.unmount()
      return found.map(x => x.props.color)
    }
    const alex = [...await colours('Alex'), ...await colours('Alex')]
    const berk = await colours('Berk')
    // Alex: in work, created in the log, twice over; Berk: put back, in both tabs.
    expect(alex.length).toBeGreaterThanOrEqual(4)
    expect(new Set(alex).size).toBe(1)
    expect(new Set(berk).size).toBe(1)
    expect(alex[0]).toBeDefined()
    expect(berk[0]).not.toBe(alex[0])
    expect(['claude', 'suggestion', 'warning']).not.toContain(alex[0])
  })

  test('the Log filters to the people chosen, several at once, All clears it, and the choice outlives closing the pane', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new today', status: 'backlog', creator: 'Alex', updated_by: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'B1B1B1', title: 'by Berk', status: 'backlog', creator: 'Berk', updated_by: 'Berk', created_at: ago(3), updated_at: ago(3) },
      { handle: 'C1C1C1', title: 'by Carla', status: 'backlog', creator: 'Carla', updated_by: 'Carla', created_at: ago(5), updated_at: ago(5) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    let pane = await $.ui.mount(OVERVIEW)
    // Queue has no filter: only the Log does.
    expect(await pane.find({ key: 'overview-people-all' })).toBeUndefined()
    await pane.press({ key: 'overview-tab-log' })
    const labels = (await pane.findAll({ type: 'Button' } as never)).map(b => (b.props as { label: string }).label)
    expect(labels).toEqual(['Queue', 'Log', 'Workers', 'Close', 'All', 'Alex', 'Berk', 'Carla', 'All', 'backlog'])
    // A name button is drawn in that person's colour.
    expect((await pane.find({ type: 'Text', text: /^Berk$/ }))?.props.color).toBeDefined()
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    const days = async () => (await pane.findAll({ type: 'Text', text: /^── / })).length
    expect(await handles()).toEqual(['N1N1N1', 'B1B1B1', 'C1C1C1'])
    await pane.press({ key: 'overview-person-Berk' })
    expect(await handles()).toEqual(['B1B1B1'])
    // A day with nobody chosen drops out.
    expect(await days()).toBe(1)
    expect((await pane.find({ key: 'overview-person-Berk' }))?.props.variant).toBe('primary')
    await pane.press({ key: 'overview-person-Carla' })
    expect(await handles()).toEqual(['B1B1B1', 'C1C1C1'])
    // A second press leaves Berk out: Carla's stays.
    await pane.press({ key: 'overview-person-Berk' })
    expect(await handles()).toEqual(['C1C1C1'])
    await pane.press({ key: 'overview-close' })
    await pane.unmount()
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    pane = await $.ui.mount(OVERVIEW)
    expect(await handles()).toEqual(['C1C1C1'])
    await pane.press({ key: 'overview-people-all' })
    expect(await handles()).toEqual(['N1N1N1', 'B1B1B1', 'C1C1C1'])
    expect((await pane.find({ key: 'overview-people-all' }))?.props.variant).toBe('primary')
  })

  test('the Log filters to the lanes chosen, in board order, several at once, together with the people filter', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new', status: 'backlog', creator: 'Alex', updated_by: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'B1B1B1', title: 'stuck', status: 'blocked', creator: 'Berk', updated_by: 'Berk', blocked_reason: 'x', created_at: ago(0.02), updated_at: ago(0.02) },
      { handle: 'S1S1S1', title: 'judged', status: 'signoff', creator: 'Alex', created_at: ago(0.03), updated_at: ago(0.03) },
      { handle: 'W1W1W1', title: 'started', status: 'in-progress', creator: 'Berk', assignee: 'Berk', created_at: ago(0.04), updated_at: ago(0.04) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const lanes = (await pane.findAll({ type: 'Button' } as never))
      .map(b => (b.props as { key?: string }).key ?? '')
      .filter(k => k.startsWith('overview-lane-'))
      .map(k => k.slice('overview-lane-'.length))
    // The board's lanes in its order; one it does not name goes last.
    expect(lanes).toEqual(['in-progress', 'signoff', 'backlog', 'blocked'])
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    expect(await handles()).toEqual(['N1N1N1', 'B1B1B1', 'S1S1S1', 'W1W1W1'])
    await pane.press({ key: 'overview-lane-blocked' })
    expect(await handles()).toEqual(['B1B1B1'])
    expect((await pane.find({ key: 'overview-lane-blocked' }))?.props.variant).toBe('primary')
    await pane.press({ key: 'overview-lane-signoff' })
    expect(await handles()).toEqual(['B1B1B1', 'S1S1S1'])
    // Both filters hold at once: blocked or signoff, and Alex's.
    await pane.press({ key: 'overview-person-Alex' })
    expect(await handles()).toEqual(['S1S1S1'])
    // signoff pressed again is left out: blocked is Berk's, so nothing stays.
    await pane.press({ key: 'overview-lane-signoff' })
    expect(await handles()).toEqual([])
    expect(await pane.find({ type: 'Text', text: 'none' })).toBeDefined()
    await pane.press({ key: 'overview-lanes-all' })
    expect(await handles()).toEqual(['N1N1N1', 'S1S1S1'])
    expect((await pane.find({ key: 'overview-lanes-all' }))?.props.variant).toBe('primary')
  })

  test('a lane or person chosen with no event left keeps its button, and pressing it on clears the filter', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'S1S1S1', title: 'judged', status: 'signoff', creator: 'Alex', created_at: ago(0.03), updated_at: ago(0.03) },
    ]
    await startSession($)
    const open = async () => {
      await (await $.ui.mount(BAND)).press({ key: 'overview' })
      const pane = await $.ui.mount(OVERVIEW)
      await pane.press({ key: 'overview-tab-log' })
      return pane
    }
    let pane = await open()
    await pane.press({ key: 'overview-lane-signoff' })
    await pane.press({ key: 'overview-person-Alex' })
    await pane.press({ key: 'overview-close' })
    await pane.unmount()
    // The ticket moves on to done and its events now belong to Berk.
    w.tickets = [{ handle: 'S1S1S1', title: 'judged', status: 'done', creator: 'Berk', created_at: ago(0.03), updated_at: ago(0.02) }]
    pane = await open()
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    expect(await handles()).toEqual([])
    expect((await pane.find({ key: 'overview-lane-signoff' }))?.props.variant).toBe('primary')
    expect((await pane.find({ key: 'overview-person-Alex' }))?.props.variant).toBe('primary')
    // Left out, the buttons stay, so the next press can still clear them.
    await pane.press({ key: 'overview-lane-signoff' })
    await pane.press({ key: 'overview-person-Alex' })
    expect((await pane.find({ key: 'overview-lane-signoff' }))?.props.label).toBe('✕ signoff')
    expect((await pane.find({ key: 'overview-person-Alex' }))?.props.label).toBe('✕ Alex')
    await pane.press({ key: 'overview-lane-signoff' })
    await pane.press({ key: 'overview-person-Alex' })
    expect(await handles()).toEqual(['S1S1S1'])
    // Unchosen and without an event, the buttons go again.
    expect(await pane.find({ key: 'overview-lane-signoff' })).toBeUndefined()
    expect(await pane.find({ key: 'overview-person-Alex' })).toBeUndefined()
    expect((await pane.find({ key: 'overview-lanes-all' }))?.props.variant).toBe('primary')
  })

  test('a chip cycles off, chosen, left out and off again, and shows when it is left out', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'A1A1A1', title: 'by Alex', status: 'backlog', creator: 'Alex', updated_by: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'B1B1B1', title: 'by Berk', status: 'backlog', creator: 'Berk', updated_by: 'Berk', created_at: ago(0.02), updated_at: ago(0.02) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    const chip = async () => {
      const p = (await pane.find({ key: 'overview-person-Berk' }))!.props as { label: string; variant?: string; dimColor?: boolean }
      return { label: p.label, variant: p.variant ?? null, dimColor: p.dimColor ?? false }
    }
    expect(await chip()).toEqual({ label: 'Berk', variant: null, dimColor: true })
    await pane.press({ key: 'overview-person-Berk' })
    expect(await chip()).toEqual({ label: 'Berk', variant: 'primary', dimColor: false })
    expect(await handles()).toEqual(['B1B1B1'])
    await pane.press({ key: 'overview-person-Berk' })
    expect(await chip()).toEqual({ label: '✕ Berk', variant: null, dimColor: false })
    expect((await pane.find({ type: 'Text', text: '✕ Berk' }))?.props.color).toBe('error')
    expect(await handles()).toEqual(['A1A1A1'])
    await pane.press({ key: 'overview-person-Berk' })
    expect(await chip()).toEqual({ label: 'Berk', variant: null, dimColor: true })
    expect(await handles()).toEqual(['A1A1A1', 'B1B1B1'])
  })

  test('left out alone, the Log shows everything else, several left out at once', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new', status: 'backlog', creator: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'B1B1B1', title: 'stuck', status: 'blocked', creator: 'Berk', blocked_reason: 'x', created_at: ago(0.02), updated_at: ago(0.02) },
      { handle: 'S1S1S1', title: 'judged', status: 'signoff', creator: 'Alex', created_at: ago(0.03), updated_at: ago(0.03) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    const out = async (key: string) => {
      await pane.press({ key })
      await pane.press({ key })
    }
    await out('overview-lane-signoff')
    expect(await handles()).toEqual(['N1N1N1', 'B1B1B1'])
    await out('overview-lane-blocked')
    expect(await handles()).toEqual(['N1N1N1'])
    // All is not lit while anything is left out.
    expect((await pane.find({ key: 'overview-lanes-all' }))?.props.variant).toBeUndefined()
  })

  test('chosen and left out together, across both rows: what is chosen and not left out', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new', status: 'backlog', creator: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'B1B1B1', title: 'stuck', status: 'blocked', creator: 'Berk', blocked_reason: 'x', created_at: ago(0.02), updated_at: ago(0.02) },
      { handle: 'S1S1S1', title: 'judged', status: 'signoff', creator: 'Alex', created_at: ago(0.03), updated_at: ago(0.03) },
      { handle: 'C1C1C1', title: 'by Carla', status: 'signoff', creator: 'Carla', created_at: ago(0.04), updated_at: ago(0.04) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    // Not Carla, but signoff and backlog: Alex's two.
    await pane.press({ key: 'overview-person-Carla' })
    await pane.press({ key: 'overview-person-Carla' })
    await pane.press({ key: 'overview-lane-signoff' })
    await pane.press({ key: 'overview-lane-backlog' })
    expect(await handles()).toEqual(['N1N1N1', 'S1S1S1'])
    // Alex chosen too, and backlog left out in the same row as signoff chosen.
    await pane.press({ key: 'overview-person-Alex' })
    await pane.press({ key: 'overview-lane-backlog' })
    expect(await handles()).toEqual(['S1S1S1'])
  })

  test('All clears its row of chosen and left out alike and is lit only then; the other row stays', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new', status: 'backlog', creator: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'B1B1B1', title: 'stuck', status: 'blocked', creator: 'Berk', blocked_reason: 'x', created_at: ago(0.02), updated_at: ago(0.02) },
      { handle: 'S1S1S1', title: 'judged', status: 'signoff', creator: 'Alex', created_at: ago(0.03), updated_at: ago(0.03) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const handles = async () => (await pane.findAll({ type: 'Text', text: LOG_LINE })).map(x => x.text.replace(LOG_LINE, '').slice(0, 6))
    const lit = async (key: string) => (await pane.find({ key }))?.props.variant === 'primary'
    await pane.press({ key: 'overview-person-Berk' })
    await pane.press({ key: 'overview-person-Berk' })
    await pane.press({ key: 'overview-lane-signoff' })
    await pane.press({ key: 'overview-lane-backlog' })
    await pane.press({ key: 'overview-lane-backlog' })
    expect(await handles()).toEqual(['S1S1S1'])
    expect(await lit('overview-lanes-all')).toBe(false)
    expect(await lit('overview-people-all')).toBe(false)
    await pane.press({ key: 'overview-lanes-all' })
    expect(await lit('overview-lanes-all')).toBe(true)
    expect((await pane.find({ key: 'overview-lane-backlog' }))?.props.label).toBe('backlog')
    // Berk stays left out: the people row was not cleared.
    expect(await handles()).toEqual(['N1N1N1', 'S1S1S1'])
    expect(await lit('overview-people-all')).toBe(false)
    await pane.press({ key: 'overview-people-all' })
    expect(await lit('overview-people-all')).toBe(true)
    expect(await handles()).toEqual(['N1N1N1', 'B1B1B1', 'S1S1S1'])
  })

  test('the top row is one line of buttons alike — the tabs on the left, Close at the right — and the open tab is the one highlighted', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...WORKING, assignee: 'Alex', updated_at: ago(0.1) }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    for (const tab of ['queue', 'log', 'workers']) {
      await pane.press({ key: `overview-tab-${tab}` })
      // No frame anywhere: every button is drawn as the terminal draws a
      // Button, '[ label ]', one line high.
      expect(await pane.find({ key: 'overview-close-frame' })).toBeUndefined()
      expect(await pane.find({ key: `overview-tab-${tab}-frame` })).toBeUndefined()
      const labels = (await pane.findAll({ type: 'Button' } as never)).map(b => (b.props as { label: string }).label)
      expect(labels.slice(0, 4)).toEqual(['Queue', 'Log', 'Workers', 'Close'])
      // The open tab, and only it, is the primary button; the others are
      // dimmed. Close is never a tab.
      for (const other of ['queue', 'log', 'workers']) {
        const open = other === tab
        const b = await pane.find({ key: `overview-tab-${other}` })
        expect(b?.props.variant).toBe(open ? 'primary' : undefined)
        expect(b?.props.dimColor).toBe(!open)
      }
      expect((await pane.find({ key: 'overview-close' }))?.props.variant).toBeUndefined()
    }
    await pane.press({ key: 'overview-close' })
    await pane.unmount()
    // An unread board keeps the same head, Close alone at its right.
    w.nextFails = true
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    expect(await (await $.ui.mount(OVERVIEW)).find({ key: 'overview-top-row-close' })).toBeDefined()
  })

  test('every Log line starts with its kind\'s sign, one wide code point, so no column shifts', async ($, on) => {
    const w = world(on)
    w.tickets = [
      { handle: 'N1N1N1', title: 'new', status: 'backlog', creator: 'Alex', created_at: ago(0.01), updated_at: ago(0.01) },
      { handle: 'P1P1P1', title: 'back', status: 'blocked', creator: 'Berk', updated_by: 'Alex', blocked_reason: 'x', created_at: ago(5), updated_at: ago(2) },
    ]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    await pane.press({ key: 'overview-tab-log' })
    const lines = (await pane.findAll({ type: 'Text', text: /\d\d:\d\d {2}[0-9A-Z]{6} / })).map(x => x.text)
    expect(lines).toHaveLength(3)
    const signs = lines.map(l => [...l][0]!)
    // Created and put back read apart, and the same kind always the same sign.
    expect(new Set(signs).size).toBe(2)
    for (const [i, line] of lines.entries()) {
      // A single code point a terminal draws two cells wide by default: no
      // variation selector, no joiner, nothing a terminal may draw one cell.
      expect(signs[i]).toMatch(/^\p{Emoji_Presentation}$/u)
      expect(line).not.toMatch(/[\uFE0E\uFE0F\u200D]/)
      // The clock starts at the same column on every line.
      expect([...line].indexOf(line.match(/\d\d:\d\d/)![0][0]!)).toBe(2)
    }
  })

  test('a board that cannot be read says so', async ($, on) => {
    const w = world(on)
    w.nextFails = true
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    expect(await (await $.ui.mount(OVERVIEW)).find({ type: 'Text', text: /could not be read/ })).toBeDefined()
  })
})

describe('the workers tab', () => {
  const WIDTH = 160
  const workersTab = async ($: Test$, columns = WIDTH, rows = 40) => {
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const base = OVERVIEW as { props: { scroll: object } }
    const pane = await $.ui.mount({ ...base, props: { ...base.props, bodyColumns: columns, scroll: { ...base.props.scroll, bodyRows: rows } } } as never)
    await pane.press({ key: 'overview-tab-workers' })
    return pane
  }
  type Pane = Awaited<ReturnType<typeof workersTab>>
  const texts = async (pane: Pane) => (await pane.findAll({ type: 'Text' })).map(x => x.text)
  // The picture's rows: whole lines of the pane's width. A row of one style
  // is its own single run too: findAll gives it twice in a row.
  const pictureOf = async (pane: Pane, columns = WIDTH) =>
    (await texts(pane)).filter((t, i, all) => t.length === columns && !/^[A-Z0-9]{6}  /.test(t) && all[i - 1] !== t)
  // One line per agent beneath the picture, its trailing room trimmed.
  const linesOf = async (pane: Pane) => (await texts(pane)).filter(t => /^[A-Z0-9]{6}  /.test(t)).map(t => t.replace(/\s+$/, ''))
  const dod = (done: number, of: number) => Array.from({ length: of }, (_, i) => ({ done: i < done, state: i < done ? 'done' : 'todo' }))

  test('the stone: a pebble for nothing open, never smaller for more, never past its box', () => {
    expect(stoneSize(0, 30, 10)).toEqual({ w: 3, h: 1 })
    expect(stoneSize(-2, 30, 10)).toEqual({ w: 3, h: 1 })
    let last = stoneSize(0, 30, 10)
    for (const open of [1, 2, 3, 5, 9, 16, 40, 100, 1000]) {
      const size = stoneSize(open, 30, 10)
      expect(size.w).toBeGreaterThanOrEqual(last.w)
      expect(size.h).toBeGreaterThanOrEqual(last.h)
      expect(size.w).toBeLessThanOrEqual(30)
      expect(size.h).toBeLessThanOrEqual(10)
      last = size
    }
    // Four times the work reads larger, not four times as wide.
    expect(stoneSize(4, 99, 99)).toEqual({ w: 12, h: 3 })
    expect(stoneSize(16, 99, 99)).toEqual({ w: 20, h: 5 })
    // A box smaller than a pebble still holds one cell.
    expect(stoneSize(9, 0, 0)).toEqual({ w: 1, h: 1 })
  })

  test('how long an agent has been seen reads in minutes, then hours', () => {
    expect([0, 59_999, 60_000, 59 * 60_000, 65 * 60_000, 25 * 3_600_000, -5].map(durationOf)).toEqual(['<1m', '<1m', '1m', '59m', '1h05', '25h00', '<1m'])
  })

  test('a report\'s first line: a notification\'s result without tags and blank lines, else the text', () => {
    expect(reportLine('<task-notification>\n<summary>x</summary>\n<result>\n\n  done: moved.\nmore\n</result>')).toBe('done: moved.')
    expect(reportLine('<task-notification>\n<task-id>a1</task-id> <status>done</status>\nplain message\nsecond')).toBe('plain message')
    expect(reportLine('')).toBe('')
    expect(reportLine('<result>\n</result>')).toBe('')
    // Angle brackets inside a line of content are content, never stripped.
    for (const line of ['keyed by Map<string, number> and a < b check > fixed', 'if n<3 then x>2', 'moved <id> to critique', '<b>bold</b> start'])
      expect(reportLine(`<task-notification>\n<result>\n${line}\n</result>`)).toBe(line)
    // A line an agent wrote to make a pattern backtrack is read in linear time.
    for (const line of ['<a></a>'.repeat(5000) + 'x', '<a>'.repeat(5000) + 'x', '<a '.repeat(5000) + '>x', `<${'a'.repeat(20000)}>x</b>`]) {
      const start = performance.now()
      expect(reportLine(`<result>\n${line}\n</result>`)).toBe(line)
      expect(performance.now() - start).toBeLessThan(200)
    }
  })

  test('a ticket with an open question and no agent is no figure; one with an agent waits for the answer', async ($, on) => {
    const w = world(on, { herdr: true })
    const asking = { handle: 'Q1Q1Q1', title: 'asks', status: 'human', next_lane: 'review', question: 'which one?' }
    w.tickets = [HUMAN, asking, { ...WORKING, assignee: 'Alex' }]
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' },
      { tab_id: 'w1:t2', label: 'Q1Q1Q1/worker', agent_status: 'idle' },
    ]
    await startSession($)
    const lines = await linesOf(await workersTab($))
    expect(lines.map(l => l.slice(0, 6))).toEqual(['T8QS1Y', 'Q1Q1Q1'])
    expect(lines[1]).toContain('waits for an answer')
    expect(lines.join('\n')).not.toContain('K7P2QX')
  })

  // Herdr's `blocked` is an approval dialog open in the agent's tab: it holds
  // the agent whatever else is true, an open question included.
  for (const question of ['', 'which one?']) {
    test(`an agent held at an approval dialog waits for approval, its mark a '!'${question ? ', over an open question' : ''}`, async ($, on) => {
      const w = world(on, { herdr: true })
      w.tickets = [{ ...WORKING, assignee: 'Alex', question }]
      w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'blocked' }]
      await startSession($)
      const pane = await workersTab($)
      const lines = await linesOf(pane)
      expect(lines).toHaveLength(1)
      expect(lines[0]).toMatch(/^T8QS1Y {2}in-progress {2}worker {2}waits for approval /)
      const picture = await pictureOf(pane)
      expect(picture.some(l => l.includes('!'))).toBe(true)
      expect(picture.some(l => /[?z]/.test(l.replace(/\d+ open/, '')))).toBe(false)
    })
  }

  test('only tickets with an open question and nobody on them: nobody on the board', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [HUMAN, { handle: 'Z1Z1Z1', title: 'old', status: 'signoff', next_lane: 'done', question: 'accept?' }]
    await startSession($)
    expect(await (await workersTab($)).find({ type: 'Text', text: 'No agent on the board right now.' })).toBeDefined()
  })

  test('one line per agent: ticket, lane, kind, state, how long it has been seen, and the first line of its last report', async ($, on) => {
    const report = '<task-notification>\n<task-id>a1</task-id>\n<result>\n\nT8QS1Y: tests green, moving to critique.\nsecond line\n</result>\n</task-notification>'
    const w = world(on, { herdr: true, store: { 'report:/repo:T8QS1Y': [{ at: '2026-11-21T10:00:00Z', origin: 'peer', session: 's', text: report }] } })
    w.tickets = [{ ...WORKING, assignee: 'Alex' }, { handle: 'R1R1R1', title: 'judged', status: 'review', next_lane: 'signoff', question: '' }]
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' },
      { tab_id: 'w1:t2', label: 'R1R1R1/review', agent_status: 'idle' },
      // A person's own chat is not an agent.
      { tab_id: 'w1:t3', label: 'T8QS1Y/explain', agent_status: 'working' },
    ]
    w.agents = [{ id: 'a1', description: 'bughunt T8QS1Y', type: 'claude', status: 'running' }]
    await startSession($)
    const pane = await workersTab($)
    expect(await linesOf(pane)).toEqual([
      'T8QS1Y  in-progress  subagent  working    <1m  T8QS1Y: tests green, moving to critique.',
      'T8QS1Y  in-progress  worker    working    <1m  T8QS1Y: tests green, moving to critique.',
      'R1R1R1  review       review    resting    <1m  judged',
    ])
    // Read again after another tab: the time it was first seen stays.
    await pane.press({ key: 'overview-tab-queue' })
    await w.clock.advance(6 * 60_000)
    await pane.press({ key: 'overview-tab-workers' })
    expect((await linesOf(pane))[1]).toBe('T8QS1Y  in-progress  worker    working     6m  T8QS1Y: tests green, moving to critique.')
  })

  test('every agent stands at one stone, its count the open definition-of-done items of the tickets with an agent', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [
      { ...WORKING, assignee: 'Alex', dod_items: [...dod(1, 4), { done: false, state: 'superseded' }] },
      { handle: 'C1C1C1', title: 'judged', status: 'critique', question: '', dod_items: dod(0, 2) },
      // No agent on it: its work is not at the stone.
      { handle: 'N1N1N1', title: 'idle', status: 'in-progress', question: '', dod_items: dod(0, 9) },
    ]
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' },
      { tab_id: 'w1:t2', label: 'C1C1C1/worker', agent_status: 'working' },
      { tab_id: 'w1:t3', label: 'C1C1C1/bughunt', agent_status: 'working' },
    ]
    await startSession($)
    const picture = await pictureOf(await workersTab($))
    // Three open on T8QS1Y, two on C1C1C1 counted once however many work it.
    expect(picture.some(l => l.includes('5 open'))).toBe(true)
    // Three crabs on one ground line, none on a track of its own.
    expect(picture.filter(l => l.includes('▐▛███▜▌')).join('').split('▐▛███▜▌').length - 1).toBe(3)
    expect(picture.filter(l => /^─+$/.test(l))).toHaveLength(1)
  })

  test('a definition-of-done item ticked makes the stone smaller, and a chunk falls from it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [{ ...WORKING, assignee: 'Alex', dod_items: dod(0, 9) }]
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' }]
    await startSession($)
    const pane = await workersTab($)
    // The stone's cells: its runs of rock, which no crab's run matches.
    const stone = async () => (await texts(pane)).filter(t => /^(▟█*▙|█+|▄+)$/.test(t)).join('').length
    const big = await stone()
    const before = await pictureOf(pane)
    expect(before.some(l => l.includes('9 open'))).toBe(true)
    expect(before.some(l => l.includes('▞'))).toBe(false)
    w.tickets = [{ ...WORKING, assignee: 'Alex', dod_items: dod(8, 9) }]
    await w.clock.advance(5000)
    const after = await pictureOf(pane)
    expect(after.some(l => l.includes('1 open'))).toBe(true)
    expect(after.some(l => l.includes('▞'))).toBe(true)
    expect(await stone()).toBeLessThan(big)
    // It has fallen after a few frames.
    await w.clock.advance(250 * 8)
    expect((await pictureOf(pane)).some(l => l.includes('▞'))).toBe(false)
  })

  for (const n of [1, 4, 8, 16, 30]) {
    test(`${n} agents and all their lines fit 40 rows × 100 columns without scrolling`, async ($, on) => {
      const w = world(on, { herdr: true })
      const ids = Array.from({ length: n }, (_, i) => `A${String(i).padStart(2, '0')}AAA`)
      w.tickets = ids.map(handle => ({ handle, title: handle, status: 'in-progress', question: '', dod_items: dod(0, 3) }))
      w.tabs = ids.map((id, i) => ({ tab_id: `w1:t${i}`, label: `${id}/worker`, agent_status: 'working' }))
      await startSession($)
      const pane = await workersTab($, 100, 40)
      const picture = await pictureOf(pane, 100)
      const lines = await linesOf(pane)
      expect(lines).toHaveLength(n)
      // 40 rows less the head: the tabs and Close, one line.
      expect(picture.length + lines.length).toBeLessThanOrEqual(39)
      const crabs = picture.join('\n')
      const big = crabs.split('▐▛███▜▌').length - 1
      const small = crabs.split('▐▛▜▌').length - 1
      // Every agent stands in the picture: the big crab while they all fit, else the small one.
      expect(big + small).toBe(n)
      expect(big === n || small === n).toBe(true)
      if (n === 30) expect(small).toBe(n)
    })
  }

  test('working hammers in place; resting sits still', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [{ ...WORKING, assignee: 'Alex' }]
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' }]
    await startSession($)
    const pane = await workersTab($)
    const seen: string[] = []
    for (let i = 0; i < 4; i++) {
      seen.push((await pictureOf(pane)).join('\n').replace(/[✻✳✶✢]/g, ''))
      await w.clock.advance(250)
    }
    expect(new Set(seen).size).toBe(2)
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'idle' }]
    await w.clock.advance(5000)
    const resting: string[] = []
    for (let i = 0; i < 4; i++) {
      resting.push((await pictureOf(pane)).join('\n'))
      await w.clock.advance(250)
    }
    expect(new Set(resting).size).toBe(1)
    expect(resting[0]).toContain('▀▀▀▀▀')
  })

  test('two tickets never wear the same colour or hat, and each keeps both on every read', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [{ ...WORKING, assignee: 'Alex' }, { handle: 'C1C1C1', title: 'judged', status: 'critique', question: '' }]
    w.tabs = [
      { tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' },
      { tab_id: 'w1:t2', label: 'C1C1C1/worker', agent_status: 'working' },
    ]
    await startSession($)
    const pane = await workersTab($)
    // Each ticket's colour as its line's handle wears it.
    const colors = async () => Object.fromEntries((await pane.findAll({ type: 'Text' }))
      .filter(x => /^[A-Z0-9]{6}$/.test(x.text) && (x.props as { bold?: boolean }).bold === true)
      .map(x => [x.text, (x.props as { color?: string }).color]))
    const first = await colors()
    expect(first.T8QS1Y).toBeDefined()
    expect(first.T8QS1Y).not.toBe(first.C1C1C1)
    const more = ['A1A1A1', 'A2A2A2', 'A3A3A3', 'A4A4A4']
    w.tickets = [...w.tickets, ...more.map(handle => ({ handle, title: handle, status: 'in-progress', question: '' }))]
    w.tabs = [...w.tabs, ...more.map((id, i) => ({ tab_id: `w1:x${i}`, label: `${id}/worker`, agent_status: 'working' }))]
    await w.clock.advance(5000)
    const after = await colors()
    expect(after.T8QS1Y).toBe(first.T8QS1Y)
    expect(after.C1C1C1).toBe(first.C1C1C1)
    expect(new Set(Object.values(after)).size).toBe(6)
  })

  test('once a ticket leaves, two that had to share a colour never keep sharing it', async ($, on) => {
    const w = world(on, { herdr: true })
    const ids = ['AAAAAA', 'BBBBBB', 'CCCCCC', 'DDDDDD', 'EEEEEE', 'FFFFFF', 'GGGGGG']
    const board = (on: string[]) => {
      w.tickets = on.map(handle => ({ handle, title: handle, status: 'in-progress', question: '' }))
      w.tabs = on.map((id, i) => ({ tab_id: `w1:t${i}`, label: `${id}/worker`, agent_status: 'working' }))
    }
    board(ids)
    await startSession($)
    const pane = await workersTab($)
    const colors = async () => (await pane.findAll({ type: 'Text' })).filter(x => /^[A-Z]{6}$/.test(x.text) && (x.props as { bold?: boolean }).bold === true).map(x => (x.props as { color?: string }).color)
    expect(new Set(await colors()).size).toBe(6)
    board(ids.filter(id => id !== 'BBBBBB'))
    await w.clock.advance(5000)
    expect(new Set(await colors()).size).toBe(6)
  })

  test('nobody on the board says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [WORKING]
    await startSession($)
    expect(await (await workersTab($)).find({ type: 'Text', text: 'No agent on the board right now.' })).toBeDefined()
  })

  test('an Overview with no room to show it starts no frames, even on the Workers tab', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' }]
    await startSession($)
    const pane = await workersTab($)
    await pane.press({ key: 'overview-close' })
    w.noRoom = true
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const before = ran(w, 'herdr', 'tab', 'list').length
    await w.clock.advance(20_000)
    expect(ran(w, 'herdr', 'tab', 'list').length).toBe(before)
  })

  test('a read of the board still running is not started a second time', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' }]
    await startSession($)
    await workersTab($)
    const reads = () => ran(w, 'herdr', 'tab', 'list').length
    const before = reads()
    w.heldTabs = []
    // Two reads due, the first one hanging: only one asked.
    await w.clock.advance(10_000)
    expect(reads()).toBe(before + 1)
    // Once it ends, the next one is asked when it is due.
    w.heldTabs.shift()!()
    w.heldTabs = undefined
    await w.clock.advance(5000)
    expect(reads()).toBe(before + 2)
  })

  test('a read still running when the tab changes or the pane closes starts no frames: nothing is read after Close', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' }]
    await startSession($)
    const reads = () => ran(w, 'herdr', 'tab', 'list').length
    // Lets the plugin run until its read of Herdr's tabs is held.
    const held = async () => {
      for (let i = 0; i < 50 && w.heldTabs!.length === 0; i++) await w.clock.advance(0)
      expect(w.heldTabs!.length).toBe(1)
    }
    // Lets the held read go, and the plugin finish what it did with it.
    const release = async () => {
      w.heldTabs!.shift()!()
      for (let i = 0; i < 50; i++) await w.clock.advance(0)
    }
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    w.heldTabs = []
    // Workers pressed, its read hangs; Queue is shown before it ends.
    const workers = pane.press({ key: 'overview-tab-workers' }).then(() => undefined)
    await held()
    await pane.press({ key: 'overview-tab-queue' })
    await release()
    await workers
    // Workers again, and the pane closed while its read hangs.
    const again = pane.press({ key: 'overview-tab-workers' }).then(() => undefined)
    await held()
    await pane.press({ key: 'overview-close' })
    await release()
    await again
    w.heldTabs = undefined
    const before = reads()
    await w.clock.advance(20_000)
    expect(reads()).toBe(before)
  })

  test('the frames stop when another tab is shown and when the pane closes', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tabs = [{ tab_id: 'w1:t1', label: 'T8QS1Y/worker', agent_status: 'working' }]
    await startSession($)
    const reads = () => ran(w, 'herdr', 'tab', 'list').length
    const pane = await workersTab($)
    let before = reads()
    await w.clock.advance(5000)
    expect(reads()).toBe(before + 1)
    await pane.press({ key: 'overview-tab-queue' })
    before = reads()
    await w.clock.advance(20_000)
    expect(reads()).toBe(before)
    await pane.press({ key: 'overview-tab-workers' })
    await pane.press({ key: 'overview-close' })
    before = reads()
    await w.clock.advance(20_000)
    expect(reads()).toBe(before)
  })
})

describe('following a pull request (CRASZS)', () => {
  const press = async ($: Test$) => {
    await startSession($)
    await pressAccept($, await $.ui.mount(BAND), 'MG677R')
  }
  const open = (n: number, extra: Record<string, unknown> = {}) =>
    JSON.stringify({ prs: [{ number: n, state: 'open', tickets: ['MG677R'], settled: true, unsettled: [], ...extra }] })
  const prReads = (w: World) => ran(w, 'jaira', 'pr')
  const follow = (n: number) =>
    `jaira-cockpit: pull request #${n} for MG677R is open - follow it: run 'jaira pr --watch --number ${n}' under a Monitor (one line per change, it ends when the pull request is merged or closed)`

  test('the pr lane Accept started: the pull request is looked for each minute and handed to this session once it is there, once', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prReads = [JSON.stringify({ prs: [] }), open(12)]
    await press($)
    // Nothing is asked at the press itself: the lane has only just started.
    expect(prReads(w)).toEqual([])
    await w.clock.advance(60_000)
    expect(prReads(w)).toEqual([['jaira', 'pr', 'feat/mg677r', '--json', '@/repo']])
    expect(w.submitted.filter(x => x.includes('pull request #'))).toEqual([])
    await w.clock.advance(60_000)
    expect(w.submitted.filter(x => x.includes('pull request #'))).toEqual([follow(12)])
    await w.clock.advance(10 * 60_000)
    expect(prReads(w)).toHaveLength(2)
  })

  test('a pull request that never appears is looked for twenty times, then no more', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    await press($)
    for (let i = 0; i < 40; i++) await w.clock.advance(60_000)
    expect(prReads(w)).toHaveLength(20)
    expect(w.submitted.filter(x => x.includes('pull request #'))).toEqual([])
  })

  test('jaira pr refusing (no gh, no sign-in) ends the looking at once', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prReads = ['!jaira: pr: not signed in']
    await press($)
    for (let i = 0; i < 5; i++) await w.clock.advance(60_000)
    expect(prReads(w)).toHaveLength(1)
  })

  test('a forge failing for a minute (exit 1) is asked again; the pull request is handed over once it answers', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prReads = ['?jaira: pr: gh: HTTP 502', open(12)]
    await press($)
    await w.clock.advance(60_000)
    await w.clock.advance(60_000)
    expect(prReads(w)).toHaveLength(2)
    expect(w.submitted.filter(x => x.includes('pull request #'))).toEqual([follow(12)])
  })

  test('a pull request handed over before is not handed over again, also by a later session', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'pr-followed:/repo:12': true } })
    w.ahead = 3
    w.prReads = [open(12)]
    await press($)
    await w.clock.advance(60_000)
    expect(prReads(w)).toHaveLength(1)
    expect(w.submitted.filter(x => x.includes('pull request #'))).toEqual([])
  })

  describe('.jaira/pr-when says signoff', () => {
    const paneText = async ($: Test$) => {
      const pane = await $.ui.mount(MESSAGE('Accept MG677R failed'))
      return (await pane.findAll({ type: 'Text' })).map(x => x.text).join('\n')
    }
    const moved = (w: World) => ran(w, 'jaira', 'move').filter(r => r.includes('--force'))

    test('no pull request yet: Accept opens it through the pr lane and the ticket stays in signoff', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      await press($)
      expect(moved(w)).toEqual([])
      expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', '--parent', 'uds:/tmp/cc-socks-1000/77.sock', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
      expect(w.toasts.at(-1)).toContain('MG677R stays in signoff: opening the pull request for feat/mg677r')
      expect(w.submitted).toEqual(['jaira-cockpit: the person pressed Accept on MG677R; this board opens the pull request in signoff, so the pr lane is opening it for feat/mg677r and MG677R stays in signoff until the review round is settled'])
    })

    test('a pull request open and not settled: nothing moves, and the pane says what is in the way', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      w.prs = [{ number: 12 }]
      w.prReads = [open(12, { settled: false, unsettled: ['2 open threads', 'review pending: copilot'] })]
      await press($)
      expect(moved(w)).toEqual([])
      expect(ran(w, 'bash', SPAWN)).toEqual([])
      const text = await paneText($)
      expect(text).toContain('#12: 2 open threads')
      expect(text).toContain('#12: review pending: copilot')
      expect(text).toContain('MG677R stays in signoff; nothing was moved.')
    })

    test('a pull request settled: Accept moves the ticket to done as everywhere', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      w.prs = [{ number: 12 }]
      w.prReads = [open(12)]
      await press($)
      expect(moved(w)).toHaveLength(1)
    })

    test('the forge lists a pull request jaira pr does not show open from origin: held, a person tells whose it is', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      w.prs = [{ number: 12 }]
      w.prReads = [JSON.stringify({ prs: [] })]
      await press($)
      expect(moved(w)).toEqual([])
      expect(ran(w, 'bash', SPAWN)).toEqual([])
      expect(await paneText($)).toContain('none is open from this repository')
    })

    test('another ticket on the branch not in signoff yet: nothing opens and nothing moves', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 2
      w.subjects = ['feat(MG677R): cache key', 'fix(AB12CD): other']
      w.tickets = [SIGNOFF, { ...SIGNOFF, handle: 'AB12CD', status: 'critique' }]
      await press($)
      expect(moved(w)).toEqual([])
      expect(ran(w, 'bash', SPAWN)).toEqual([])
      expect(await paneText($)).toContain('the pull request for feat/mg677r waits: AB12CD not in signoff yet')
    })

    test('no lock against a second press: each press starts the pr lane, no tab state is read, nothing moves', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      await press($)
      w.tabs = [{ tab_id: 'w1:t7', label: 'MG677R/pr', agent_status: 'blocked' }]
      await w.clock.advance(3_000)
      await pressAccept($, await $.ui.mount(BAND), 'MG677R')
      expect(ran(w, 'bash', SPAWN)).toHaveLength(2)
      expect(moved(w)).toEqual([])
    })

    test('an open definition-of-done item refuses the press before any pull request is opened', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      w.refusals = [...w.refusals, { code: 'needs_nonmodel_signal', message: 'the definition of done is not met: criterion 1 is still open' }]
      await startSession($)
      const pane = await pressAccept($, await $.ui.mount(BAND), 'MG677R')
      expect(ran(w, 'bash', SPAWN)).toEqual([])
      expect(ran(w, 'jaira', 'roles')).toEqual([])
      expect(moved(w)).toEqual([])
      // A gate's refusal is said in the Acceptance pane (2V50Y5), not the message pane.
      expect((await pane.findAll({ type: 'Text' })).map(x => x.text).join('\n')).toContain('criterion 1 is still open')
    })

    test('pr-when that cannot be read counts as done: Accept moves the ticket as on any board', async ($, on) => {
      const w = world(on, { herdr: true })
      w.rolesFails = true
      w.ahead = 3
      await press($)
      expect(moved(w)).toHaveLength(1)
      expect(ran(w, 'jaira', 'pr')).toEqual([])
    })

    test('a forge that cannot be asked holds the ticket too, and says why', async ($, on) => {
      const w = world(on, { herdr: true })
      w.prWhen = 'signoff'
      w.ahead = 3
      w.prs = null
      await press($)
      expect(moved(w)).toEqual([])
      expect(await paneText($)).toContain('pull requests not checked')
    })
  })
})

// 76ERHH: a press writes from a disk that holds the ticket as its ref has it.
// A file on the ref alone, or one older than the ref, refused every Accept and
// Give back ("on its ref and not on your disk"); the cockpit pulls it first.
describe('a ticket not on the disk the press writes from', () => {
  const OWN = '@/.worktrees/repo-mg677r'
  const MAIN = '@/repo'
  // The main checkout alone: the ticket has no worktree of its own.
  const ALONE = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n'
  const press = async ($: Test$, w: World, how: 'accept' | 'give back') => {
    const ui = await $.ui.mount(BAND)
    if (how === 'accept') return void (await pressAccept($, ui, 'MG677R', w))
    await ui.press({ key: 'acceptance-MG677R' })
    const typed = (await $.ui.mount(ACCEPT('Acceptance of MG677R'))).input({ key: 'typed-MG677R', text: 'fehlt' })
    await w.clock.advance(1000)
    await typed
  }
  const moved = (w: World, how: 'accept' | 'give back') =>
    ran(w, 'jaira', 'move', 'MG677R').filter(r => !r.includes('--dry-run') && r[4] === (how === 'accept' ? 'done' : 'in-progress'))

  const cases = [
    { tree: 'own worktree', where: OWN, state: 'ref-only', pulls: true, cleans: false },
    { tree: 'own worktree', where: OWN, state: 'behind', pulls: true, cleans: false },
    { tree: 'own worktree', where: OWN, state: 'current', pulls: false, cleans: false },
    { tree: 'no worktree', where: MAIN, state: 'ref-only', pulls: true, cleans: true },
    { tree: 'no worktree', where: MAIN, state: 'behind', pulls: true, cleans: false },
    { tree: 'no worktree', where: MAIN, state: 'current', pulls: false, cleans: false },
  ] as const
  for (const how of ['accept', 'give back'] as const) {
    for (const c of cases) {
      test(`${how}, ${c.tree}, file ${c.state}: ${c.pulls ? 'pulled first, then decided' : 'decided as it is'}${c.cleans ? ', the pulled copy taken away after' : ''}`, async ($, on) => {
        const w = world(on)
        if (c.tree === 'no worktree') w.worktrees = ALONE
        if (c.state !== 'current') w.disk[c.where] = c.state
        await startSession($)
        await press($, w, how)
        expect(moved(w, how)).toHaveLength(1)
        expect(moved(w, how)[0]!.at(-1)).toBe(c.where)
        expect(w.opened.filter(o => o.id === 'cockpit-message')).toEqual([])
        expect(ran(w, 'jaira', 'pull', 'MG677R').map(r => r.at(-1))).toEqual(c.pulls ? [c.where] : [])
        // The pull before the first write, the dry run included.
        if (c.pulls) expect(w.ran.findIndex(r => r[1] === 'pull')).toBeLessThan(w.ran.findIndex(r => r[1] === 'move'))
        const cleaned = ran(w, 'git', 'clean')
        expect(cleaned).toEqual(c.cleans ? [['git', 'clean', '-fqx', '--', '/repo/.jaira/tickets/01MG677R-cache-key.md', MAIN]] : [])
        // Taken away only once the press wrote all it writes.
        if (c.cleans) expect(w.ran.findIndex(r => r[0] === 'git' && r[1] === 'clean')).toBeGreaterThan(w.ran.findLastIndex(r => r[0] === 'jaira' && ['move', 'note', 'create'].includes(r[1]!)))
      })
    }
  }

  test('a pulled copy is taken away from the main checkout when the press is refused too', async ($, on) => {
    const w = world(on)
    w.worktrees = ALONE
    w.disk[MAIN] = 'ref-only'
    w.refusals = [{ code: 'dod_open', message: 'definition of done item 2 is open' }]
    await startSession($)
    await press($, w, 'accept')
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toEqual([])
    // The refusal is said in the Acceptance pane, which stays open.
    expect(w.closed).not.toContain('cockpit-acceptance')
    expect(ran(w, 'git', 'clean')).toHaveLength(1)
  })

  test('a copy that cannot be taken away is said beside the band, leaving the message pane to what the press said', async ($, on) => {
    const w = world(on)
    w.worktrees = ALONE
    w.disk[MAIN] = 'ref-only'
    w.cleanFails = 'warning: failed to remove .jaira/tickets/01MG677R-cache-key.md: Permission denied'
    w.refusals = [{ code: 'dod_open', message: 'definition of done item 2 is open' }]
    await startSession($)
    await press($, w, 'accept')
    // The refusal is said in the Acceptance pane; the copy beside the band.
    expect(w.closed).not.toContain('cockpit-acceptance')
    const said = w.toasts.at(-1) ?? ''
    expect(said).toContain('Permission denied')
    expect(said).toContain("git clean -fqx -- '/repo/.jaira/tickets/01MG677R-cache-key.md'")
  })

  test('a lost give-back note stays in the message pane when the copy cannot be taken away either', async ($, on) => {
    const w = world(on)
    w.worktrees = ALONE
    w.disk[MAIN] = 'ref-only'
    w.cleanFails = 'Permission denied'
    w.noteFails = true
    await startSession($)
    await press($, w, 'give back')
    expect(w.opened.filter(o => o.id === 'cockpit-message').at(-1)?.title).toBe('Note on MG677R failed')
    expect(w.toasts.some(x => x.includes('git clean -fqx'))).toBe(true)
  })

  for (const how of ['accept', 'give back'] as const) {
    test(`${how}: a refused pull stops the press with its reason, writes nothing and never steals`, async ($, on) => {
      const w = world(on)
      w.worktrees = ALONE
      w.disk[MAIN] = 'ref-only'
      w.pullFails = 'jaira: MG677R is taken by Berk; --steal takes it anyway'
      await startSession($)
      await press($, w, how)
      expect(ran(w, 'jaira', 'pull', 'MG677R')).toEqual([['jaira', 'pull', 'MG677R', MAIN]])
      expect(ran(w, 'jaira', 'move', 'MG677R')).toEqual([])
      expect(ran(w, 'jaira', 'note')).toEqual([])
      expect(ran(w, 'git', 'clean')).toEqual([])
      expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: `${how === 'accept' ? 'Accept' : 'Give back'} MG677R failed`, focus: true })
      expect(w.submitted).toEqual([])
    })
  }

  test('a decision picked in the explanation pulls once for its note and the Give back it carries out', async ($, on) => {
    const w = world(on)
    w.worktrees = ALONE
    w.disk[MAIN] = 'ref-only'
    w.bodies.MG677R = [
      '## Progress',
      '',
      '- **2026-10-12 15:00 · Alex** — explain: The cache key leaves out the branch.',
      '1) add the branch to the key [give back]',
      '2) accept it as it is [accept]',
      'Empfehlung: 1',
    ].join('\n')
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    await (await $.ui.mount(EXPLANATION('Explanation of MG677R'))).press({ key: 'pick-MG677R-1' })
    expect(ran(w, 'jaira', 'note', 'MG677R').map(r => r[3]).slice(0, 2)).toEqual(['decision: 1) add the branch to the key', 'Given back from signoff: 1) add the branch to the key'])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'in-progress')).toHaveLength(1)
    expect(ran(w, 'jaira', 'pull', 'MG677R')).toHaveLength(1)
    expect(ran(w, 'git', 'clean')).toHaveLength(1)
    expect(w.ran.findIndex(r => r[0] === 'git' && r[1] === 'clean')).toBeGreaterThan(w.ran.findLastIndex(r => r[0] === 'jaira' && ['move', 'note', 'create'].includes(r[1]!)))
  })
})

// The Verdict pane's light, worked out from the ticket alone: one rule set per
// colour, the reason it names first.
describe('the verdict light', () => {
  const ready: Waiting = {
    id: 'MG677R', title: 'cache key', goal: '', lane: 'signoff', verdict: 'approve', question: '', dodComplete: true,
    dodOpen: { person: 0, other: 0 }, acceptTo: 'done', assignee: '', blockers: [], notes: [], findingsKnown: true, gaps: '',
  }
  const first = (w: Waiting, hunt: Parameters<typeof lightOf>[1]) => {
    const l = lightOf(w, hunt)
    return [l.color, l.reasons[0]?.kind ?? '-']
  }

  test('green: an accepting review, a hunt on HEAD, the definition of done ticked, nothing open', () => {
    expect(first(ready, { state: 'ok', at: 'abc1234' })).toEqual(['green', '-'])
    // A ticket that runs no hunt after a fix, by choice, is no reason to look.
    expect(first(ready, { state: 'skipped', at: 'abc1234' })).toEqual(['green', '-'])
    expect(first({ ...ready, verdict: 'approve, no changes needed' }, { state: 'ok', at: 'abc1234' })).toEqual(['green', '-'])
  })

  test('yellow: accepted, but something wants a look first', () => {
    expect(first({ ...ready, dodComplete: false, dodOpen: { person: 1, other: 0 } }, { state: 'ok', at: 'abc1234' })).toEqual(['yellow', 'byAPerson'])
    expect(first({ ...ready, notes: ['note: a test is missing'] }, { state: 'ok', at: 'abc1234' })).toEqual(['yellow', 'notes'])
    expect(first({ ...ready, verdict: 'approve, but the docs are thin' }, { state: 'ok', at: 'abc1234' })).toEqual(['yellow', 'reservation'])
    expect(first({ ...ready, question: 'keep the old key?' }, { state: 'ok', at: 'abc1234' })).toEqual(['yellow', 'question'])
    expect(first({ ...ready, reviewedIn: 2 }, { state: 'ok', at: 'abc1234' })).toEqual(['yellow', 'oldReview'])
    expect(first({ ...ready, verdict: '' }, { state: 'ok', at: 'abc1234' })).toEqual(['yellow', 'notReviewed'])
    // Nobody can tell: a hunted file changed and not committed, or a jaira that does not say.
    expect(first(ready, { state: 'unknown', at: 'abc1234' })).toEqual(['yellow', 'huntUnknown'])
    expect(first(ready, undefined)).toEqual(['yellow', 'huntUnknown'])
  })

  test('red: the review gives it back, a blocker, a hunt not on HEAD, or a gate that refuses', () => {
    expect(first({ ...ready, verdict: 'request changes: the key ignores the branch' }, { state: 'ok', at: 'abc1234' })).toEqual(['red', 'refused'])
    expect(first({ ...ready, blockers: ['blocker: the key ignores the branch'] }, { state: 'ok', at: 'abc1234' })).toEqual(['red', 'blockers'])
    expect(first(ready, { state: 'fail', at: 'abc1234' })).toEqual(['red', 'staleHunt'])
    // An open item other than a person's own: the terminal lane's gate refuses it.
    expect(first({ ...ready, dodComplete: false, dodOpen: { person: 1, other: 1 } }, { state: 'ok', at: 'abc1234' })).toEqual(['red', 'dodOpen'])
    // Every red reason counts before any yellow one.
    expect(lightOf({ ...ready, question: 'why?', notes: ['note: x'] }, { state: 'fail', at: 'abc1234' }).reasons.map(r => r.kind)).toEqual(['staleHunt', 'notes', 'question'])
  })

  // A ticket in human waits for an answer: its open items are no reason to
  // give it back, so the question decides the light.
  test('a question in a lane that does not end the work makes the light yellow, the question first', () => {
    const asking = { ...ready, lane: 'human', acceptTo: null, question: 'keep the old key?', dodComplete: false, dodOpen: { person: 1, other: 2 } }
    const l = lightOf(asking, { state: 'fail', at: 'abc1234' })
    expect([l.color, l.reasons.map(r => r.kind)]).toEqual(['yellow', ['question', 'staleHunt', 'dodOpen', 'byAPerson']])
    // In signoff, where Accept ends the work, the same open items stay red.
    expect(first({ ...asking, lane: 'signoff', acceptTo: 'done' }, { state: 'ok', at: 'abc1234' })).toEqual(['red', 'dodOpen'])
  })

  // GY8T52: the review states each note's class after 'note:'; the light
  // and the Acceptance pane read that word, so they show what the teamlead says.
  test('a note colours the light only when its class is edge, or when it states none', () => {
    const ok = { state: 'ok', at: 'abc1234' } as const
    const cases: [string, string][] = [
      ['note: limit: windows paths over 260 characters', 'green'],
      ['note: elsewhere: load.go:41 accepts an empty file', 'green'],
      ['note: pre-existing: load.go:41 accepts an empty file', 'green'],
      ['note: edge: the size check is missing', 'yellow'],
      ['note: edge: fix-now: the guard misses an empty name', 'yellow'],
      ['note: fix-now: the guard misses an empty name', 'yellow'],
      ['note: the size check is missing', 'yellow'],
    ]
    for (const [note, color] of cases) expect([note, lightOf({ ...ready, notes: [note] }, ok).color]).toEqual([note, color])
  })

  // Every way a review writes a finding, and the class the pane gives it: the
  // stated word first; without one as before - the Explain note when it
  // classes every finding, then pre-existing, then edge.
  test('gapItems reads the class each finding states', () => {
    const cases: [string, string | undefined, string][] = [
      ['blocker: the header is written twice', 'blocker', 'blocker'],
      ['**Blocker**: the header is written twice', 'blocker', 'blocker'],
      ['note: edge: the size check is missing', 'edge', 'edge'],
      ['note: limit: windows paths over 260 characters', 'limit', 'limit'],
      ['note: elsewhere: load.go:41 accepts an empty file', 'elsewhere', 'elsewhere'],
      ['- **note:** `limit`: windows paths', 'limit', 'limit'],
      ['note: Limit: windows paths', 'limit', 'limit'],
      ['note: pre-existing: load.go:41', 'elsewhere', 'elsewhere'],
      ['note: preexisting: load.go:41', 'elsewhere', 'elsewhere'],
      ['note: pre-existing load.go:41 accepts an empty file', 'elsewhere', 'elsewhere'],
      ['note: fix-now: the guard', 'edge', 'edge'],
      ['note: edge: fix-now: the guard', 'edge', 'edge'],
      // A class word that is only the sentence's first word is no class.
      ['note: edge cases are untested', undefined, 'edge'],
      ['note: limits are not checked', undefined, 'edge'],
      ['note: the size check is missing', undefined, 'edge'],
    ]
    for (const [line, stated, cls] of cases) {
      expect([line, classOf(line.replace(/^- /, ''))]).toEqual([line, stated])
      expect([line, gapItems(line)[0]?.cls]).toEqual([line, cls])
    }
  })

  test('a stated class outranks the Explain note; an older mark or none still follows it', () => {
    const gaps = 'note: limit: the paths\nnote: pre-existing: load.go\nnote: the size'
    const explain = '[small edge] the paths\n[known limit] load.go\n[not this ticket] the size\n'
    expect(gapItems(gaps, explain).map(i => i.cls)).toEqual(['limit', 'limit', 'elsewhere'])
    expect(gapItems(gaps).map(i => i.cls)).toEqual(['limit', 'elsewhere', 'edge'])
  })

  test("a follow-up's title drops the class word and an older mark, once each", () => {
    expect(followTitle('note: elsewhere: load.go accepts an empty file')).toBe('load.go accepts an empty file')
    expect(followTitle('note: edge: fix-now: the guard misses an empty name')).toBe('the guard misses an empty name')
    // Once each: a finding whose own words open with a class word keeps them.
    expect(followTitle('note: edge: Limit: 0 means unlimited in jaira list')).toBe('Limit: 0 means unlimited in jaira list')
    expect(followTitle('note: limit: Edge: lanes with one ticket')).toBe('Edge: lanes with one ticket')
    expect(followTitle('note: pre-existing: load.go accepts an empty file')).toBe('load.go accepts an empty file')
  })

  // No hunt finished (jaira's null) is a look, never red: a ticket that
  // changed no code has none, and the brief cannot tell it from one owed.
  test('no hunt at all is yellow; a hunt that no longer covers HEAD is red', () => {
    expect(first(ready, { state: 'fail', at: '' })).toEqual(['yellow', 'noHunt'])
    expect(first(ready, { state: 'fail', at: 'abc1234' })).toEqual(['red', 'staleHunt'])
  })

  // Which way a verdict goes is read off the words the review lane declares,
  // as jaira's own gate reads them, never off its wording.
  test('the review lane\'s declared words decide which way a verdict goes', () => {
    const lanes = JSON.stringify({ lanes: [
      { id: 'review', agentic: true, context: 'fresh', verdicts: { accept: ['accept'], send_back: ['send-back'] } },
      { id: 'signoff', agentic: false, context: '', human_exit: true },
    ] })
    const light = (verdict: string, gaps = 'none') => {
      const list = JSON.stringify({ tickets: [{ handle: 'MG677R', title: 't', status: 'signoff', dod_complete: true, review: { verdict, gaps } }] })
      const w = snapshotOf(lanes, list, [], []).waiting[0]!
      return first(w, { state: 'ok', at: 'abc1234' })
    }
    // A send-back with nothing in its gaps is still a send-back.
    expect(light('send-back: DoD 2 unmet')).toEqual(['red', 'refused'])
    expect(light('Send back, DoD 1')).toEqual(['red', 'refused'])
    expect(light('accept')).toEqual(['green', '-'])
    expect(light('Accept.')).toEqual(['green', '-'])
    // Free text on a board with words: neither way, so a look — never red for a word like 'not'.
    expect(light('Not a single defect found; accept.')).toEqual(['yellow', 'freeVerdict'])
    expect(light('send-backs')).toEqual(['yellow', 'freeVerdict'])
    // Without declared words (an older jaira) the wording is all there is.
    const old = JSON.stringify({ lanes: [{ id: 'signoff', agentic: false, context: '', human_exit: true }] })
    const list = JSON.stringify({ tickets: [{ handle: 'MG677R', title: 't', status: 'signoff', dod_complete: true, review: { verdict: 'request changes', gaps: 'none' } }] })
    expect(first(snapshotOf(old, list, [], []).waiting[0]!, { state: 'ok', at: 'abc1234' })).toEqual(['red', 'refused'])
  })
})

// Every review note opens with the mark of its class, and reads in the
// language the person set in their own settings.json - shortened by the
// cheap model, kept per ticket and notes, the original where there is none.
describe('review notes as a person reads them', () => {
  const NOTES = 'blocker: the header is written twice\nnote: the size check is missing\nnote: windows paths over 260 characters\nnote: pre-existing: load.go:41 accepts an empty file'
  // The Explain note classes every note, so the third is a known limit.
  const CLASSES = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n[blocker] the header\n[small edge] the size\n[known limit] the paths\n[not this ticket] load.go\n'
  const SETTINGS = '/home/a/.jaira/settings.json'
  const ru = ['Заголовок пишется дважды.', 'Нет проверки размера.', 'Пути Windows длиннее 260 символов.', 'load.go:41 принимает пустой файл.']
  const de = ['Der Header wird doppelt geschrieben.', 'Die Größenprüfung fehlt.', 'Windows-Pfade über 260 Zeichen.', 'load.go:41 nimmt eine leere Datei an.']
  const noteList = NOTES.split('\n')
  const dictionary = (answer: readonly string[]) => Object.fromEntries(noteList.map((n, i) => [n, answer[i]!]))

  function withNotes(on: On, language?: string, answer: readonly string[] = ru) {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps: NOTES } }, HUMAN, WORKING]
    w.bodies.MG677R = CLASSES
    if (language !== undefined) w.files.set(SETTINGS, JSON.stringify({ dispatch: 'worker', language }))
    w.translations = dictionary(answer)
    return w
  }
  // The pane as it opens, or after one press on Translate.
  async function verdictNotes($: Test$, translated = false) {
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    if (translated) await pane.press({ key: 'translate-MG677R' })
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    await pane.unmount()
    return texts.filter(x => /^\p{Extended_Pictographic}/u.test(x))
  }
  async function giveBackRows($: Test$, translated = false) {
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    if (translated) await pane.press({ key: 'translate-MG677R' })
    const marks = (await pane.findAll({ type: 'Text' })).map(x => x.text).filter(x => /^[🔴🟡⚪🔵] /u.test(x))
    const out: string[] = []
    for (let i = 0; i < marks.length; i++) {
      let label = ''
      for (const [c, l] of [['now', '[fix now]'], ['later', '[follow-up]'], ['accept', '[accept]']] as const) if ((await pane.find({ key: `${c}-MG677R-${i}` }))?.text.startsWith('[•')) label = l
      out.push(`${label} ${marks[i]}`)
    }
    return { out, pane }
  }
  const ORIGINAL = ['🔴 blocker: the header is written twice', '🟡 note: the size check is missing', '⚪ note: windows paths over 260 characters', '🔵 note: pre-existing: load.go:41 accepts an empty file']

  test('each class has its mark in the Acceptance pane; without a language there is no Translate and nothing is sent', async ($, on) => {
    const w = withNotes(on)
    await startSession($)
    expect(await verdictNotes($)).toEqual(ORIGINAL)
    expect((await giveBackRows($)).out).toEqual([
      '[fix now] 🔴 blocker: the header is written twice',
      '[accept] 🟡 note: the size check is missing',
      '[accept] ⚪ note: windows paths over 260 characters',
      '[accept] 🔵 note: pre-existing: load.go:41 accepts an empty file',
    ])
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    expect(await pane.find({ key: 'translate-MG677R' })).toBeUndefined()
    expect(w.modelCalls).toBe(0)
  })

  // FMWBA8: a language set is no call - the pane opens on the original, as
  // often as it opens, until Translate is pressed.
  test('with a language set the pane opens on the original and asks nothing until Translate is pressed', async ($, on) => {
    const w = withNotes(on, 'Русский')
    await startSession($)
    expect(await verdictNotes($)).toEqual(ORIGINAL)
    expect(await verdictNotes($)).toEqual(ORIGINAL)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    const button = await pane.find({ key: 'translate-MG677R' })
    expect(button?.props.hotkey).toBe('t')
    expect(button?.text).toContain('Translate')
    expect(w.modelCalls).toBe(0)
  })

  // Two people, two machines: one test each, as each has its own settings.
  for (const [language, answer] of [['Русский', ru], ['Deutsch', de]] as const) {
    test(`a person who set ${language} presses Translate and reads the notes in it, the ticket keeps its words`, async ($, on) => {
      const w = withNotes(on, language, answer)
      await startSession($)
      expect(await verdictNotes($, true)).toEqual([`🔴 ${answer[0]}`, `🟡 ${answer[1]}`, `⚪ ${answer[2]}`, `🔵 ${answer[3]}`])
      expect(w.modelCalls).toBe(1)
      expect(w.prompts.at(-1)).toContain(`into ${language}`)
      expect(w.prompts.at(-1)).toContain('blocker: the header is written twice')
      // Nothing of it reaches the board: no jaira command writes a word of it.
      expect(w.ran.filter(r => r[0] === 'jaira' && ['note', 'set', 'move'].includes(r[1]!))).toEqual([])
      expect(JSON.stringify(w.ran)).not.toContain(answer[0])
    })
  }

  test('one press translates the head and every note in one call; a second press goes back to the original', async ($, on) => {
    const w = withNotes(on, 'Русский')
    w.shows.MG677R = { title: 'cache key', goal: 'a key per branch', outcome: { what: 'the key names the branch', why: 'two branches shared one' } }
    w.translations = { ...dictionary(ru), 'cache key': 'ключ кэша', 'a key per branch': 'ключ на ветку', 'the key names the branch': 'ключ называет ветку', 'two branches shared one': 'две ветки делили один' }
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await w.clock.settle()
    await pane.press({ key: 'translate-MG677R' })
    const all = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    expect(w.modelCalls).toBe(1)
    const asked = JSON.parse(w.prompts[0]!.slice(w.prompts[0]!.lastIndexOf('\n') + 1)) as string[]
    expect(asked.slice(0, 4)).toEqual(['cache key', 'a key per branch', 'the key names the branch', 'two branches shared one'])
    expect(asked.slice(-4)).toEqual(noteList)
    for (const text of ['ключ кэша', '┃ Goal: ключ на ветку', '┃ Done: ключ называет ветку', '┃ Why: две ветки делили один', `🔴 ${ru[0]}`]) expect(all).toContain(text)
    expect((await pane.find({ key: 'translate-MG677R' }))?.text).toContain('Original')
    await pane.press({ key: 'translate-MG677R' })
    const back = (await pane.findAll({ type: 'Text' })).map(x => x.text)
    expect(back).toContain('cache key')
    expect(back).toContain(ORIGINAL[0])
    expect(w.modelCalls).toBe(1)
  })

  test('a head the pane draws from the board row - no goal in show - is translated as drawn', async ($, on) => {
    const w = withNotes(on, 'Русский')
    w.tickets = [{ ...SIGNOFF, goal: 'a key per branch', review: { verdict: SIGNOFF.review.verdict, gaps: NOTES } }, HUMAN, WORKING]
    w.shows.MG677R = { title: 'cache key', goal: '' }
    w.translations = { ...dictionary(ru), 'a key per branch': 'ключ на ветку' }
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain('┃ Goal: a key per branch')
    await pane.press({ key: 'translate-MG677R' })
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain('┃ Goal: ключ на ветку')
  })

  test('Give back shows the notes translated and writes the originals onto the ticket', async ($, on) => {
    const w = withNotes(on, 'Русский')
    await startSession($)
    const { out, pane } = await giveBackRows($, true)
    expect(out).toEqual([`[fix now] 🔴 ${ru[0]}`, `[accept] 🟡 ${ru[1]}`, `[accept] ⚪ ${ru[2]}`, `[accept] 🔵 ${ru[3]}`])
    await pane.input({ key: 'typed-MG677R', text: '' })
    const written = ran(w, 'jaira', 'note', 'MG677R').map(r => r[3]).join('\n')
    expect(written).toContain('- blocker: the header is written twice')
    expect(written).not.toContain(ru[0])
  })

  test('kept in the store per ticket and texts: a press on a ticket translated before asks nothing, other notes or another language ask again', async ($, on) => {
    const w = withNotes(on, 'Русский')
    await startSession($)
    await verdictNotes($, true)
    expect(await verdictNotes($)).toEqual(ORIGINAL)
    expect((await verdictNotes($, true))[0]).toBe(`🔴 ${ru[0]}`)
    await giveBackRows($, true)
    expect(w.modelCalls).toBe(1)
    w.files.set(SETTINGS, JSON.stringify({ language: 'Deutsch' }))
    w.translations = dictionary(de)
    expect((await verdictNotes($, true))[0]).toBe(`🔴 ${de[0]}`)
    expect(w.modelCalls).toBe(2)
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps: 'note: only one now' } }, HUMAN, WORKING]
    w.translations = { 'note: only one now': 'Nur noch eine.' }
    await w.clock.advance(60_000)
    expect(await verdictNotes($, true)).toEqual(['🟡 Nur noch eine.'])
    expect(w.modelCalls).toBe(3)
  })

  for (const answer of ['not json at all', JSON.stringify(ru.slice(0, 3)), JSON.stringify(['', ...ru.slice(1)]), JSON.stringify([1, 2, 3, 4]), '{"a": 1}']) {
    test(`an answer that does not fit leaves the original text with its mark: ${answer.slice(0, 30)}`, async ($, on) => {
      const w = withNotes(on, 'Русский')
      w.translations = undefined
      w.model = answer
      await startSession($)
      expect((await verdictNotes($, true))[0]).toBe('🔴 blocker: the header is written twice')
      expect(w.modelCalls).toBe(1)
    })
  }

  test('a model that does not answer leaves the original, and the next press asks again', async ($, on) => {
    const w = withNotes(on, 'Русский')
    w.modelDown = true
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'translate-MG677R' })
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain(ORIGINAL[1])
    expect((await pane.find({ key: 'translate-MG677R' }))?.text).toContain('Translate')
    w.modelDown = false
    await pane.press({ key: 'translate-MG677R' })
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain(`🟡 ${ru[1]}`)
    expect(w.modelCalls).toBe(2)
  })

  test('a call answered after a second press, or after the pane opened again, shows nothing', async ($, on) => {
    const w = withNotes(on, 'Deutsch', de)
    w.modelHeld = []
    await startSession($)
    const band = await $.ui.mount(BAND)
    const pane = await acceptancePane($, band, 'MG677R')
    void pane.press({ key: 'translate-MG677R' })
    await w.clock.settle()
    expect(w.modelHeld).toHaveLength(1)
    expect((await pane.find({ key: 'translate-MG677R' }))?.text).toContain('translating…')
    // The second press while it waits: back to the original at once.
    await pane.press({ key: 'translate-MG677R' })
    w.modelHeld.shift()!()
    await w.clock.settle()
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain(ORIGINAL[0])
    // The answer that came back is kept all the same: the next press shows
    // it at once, without a call.
    await pane.press({ key: 'translate-MG677R' })
    expect((await pane.findAll({ type: 'Text' })).map(x => x.text)).toContain(`🔴 ${de[0]}`)
    expect(w.modelCalls).toBe(1)
    // Another language, pressed, and the pane opened anew while the call waits.
    w.translations = dictionary(ru)
    w.files.set(SETTINGS, JSON.stringify({ language: 'Русский' }))
    const ru_ = await acceptancePane($, band, 'MG677R')
    void ru_.press({ key: 'translate-MG677R' })
    await w.clock.settle()
    expect(w.modelHeld).toHaveLength(1)
    const again = await acceptancePane($, band, 'MG677R')
    w.modelHeld.shift()!()
    await w.clock.settle()
    w.modelHeld = undefined
    expect((await again.findAll({ type: 'Text' })).map(x => x.text)).toContain(ORIGINAL[0])
  })
})

describe('the language and the rephrased answer, read alone', () => {
  test('the language is the settings file\'s, trimmed; none, an empty one or a broken file is none', () => {
    expect(languageOf('{"language": " Русский "}')).toBe('Русский')
    expect(languageOf('{"language": "Deutsch", "dispatch": "worker"}')).toBe('Deutsch')
    for (const json of [undefined, '', '{}', '{"language": ""}', '{"language": "  "}', '{"language": 7}', 'null', '{broken']) expect(languageOf(json)).toBeUndefined()
  })

  test('an answer is taken only as an array of exactly as many non-empty strings, prose and fences around it allowed', () => {
    expect(rephrasedOf('["a", "b"]', 2)).toEqual(['a', 'b'])
    expect(rephrasedOf('```json\n["a",\n "b  c"]\n```', 2)).toEqual(['a', 'b c'])
    expect(rephrasedOf('Here: ["a", "b"] - done', 2)).toEqual(['a', 'b'])
    for (const answer of ['', 'a, b', '["a"]', '["a", "b", "c"]', '["a", " "]', '["a", null]', '["a", "b"']) expect(rephrasedOf(answer, 2)).toBeUndefined()
  })

  test('the prompt names the language and carries the head and every note word for word, in one array', () => {
    const p = translatePrompt('Русский', ['cache key'], ['blocker: `x` in a.go:3', 'note: "quoted"'])
    expect(p).toContain('into Русский')
    expect(p).toContain('The first 1 are the title')
    expect(p).toContain('The last 2 are review notes')
    expect(p).toContain('exactly 3 strings')
    expect(p).toContain(JSON.stringify(['cache key', 'blocker: `x` in a.go:3', 'note: "quoted"']))
    expect(translatePrompt('Deutsch', ['cache key'], [])).not.toContain('review notes')
  })
})

// DF46B9: Explain and Discuss open a chat that writes its note onto the
// ticket, which a file on the ref alone or older than the ref refuses. They
// pull it first through the same ticketHere as Accept (76ERHH). A copy pulled
// into the main checkout is the chat's to remove, as its last step: the
// cockpit removes it only when no chat started, and names it a minute after
// the chat is done if it is still there.
describe('Explain and Discuss on a ticket not on the disk the chat writes from', () => {
  const MAIN = '@/repo'
  const PATH = (id: string) => `/repo/.jaira/tickets/01${id}-cache-key.md`
  const ALONE = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n'
  const KINDS = [
    { kind: 'explain', id: 'MG677R', title: 'Explain MG677R failed' },
    { kind: 'discuss', id: 'K7P2QX', title: 'Discuss K7P2QX failed' },
  ] as const
  const setUp = (w: World, id: string, state: 'ref-only' | 'behind') => {
    w.worktrees = ALONE
    w.diskId = id
    w.disk[MAIN] = state
    // The file jaira pull writes there.
    w.files.set(PATH(id), '')
  }
  // The chat has done its part: its note on the ticket (1_000_000 ms is
  // 00:16 UTC, the minute it was opened), its tab no longer working.
  const finish = (w: World, kind: 'explain' | 'discuss', id: string) => {
    w.bodies[id] = `## Progress\n\n- **1970-01-01 00:16 · Alex** — ${kind === 'explain' ? 'explain' : 'decision'}: done\n`
    w.tabs = [{ tab_id: 'w1:t4', label: `${id}/${kind}`, agent_status: 'idle' }]
  }
  const copyLine = (w: World, id: string) => w.toasts.filter(x => x.startsWith(`${id}: the copy pulled into the main checkout for a chat`))

  for (const { kind, id } of KINDS) {
    test(`${kind}, no worktree, file ref-only: pulled before the chat opens, which is handed the copy to remove itself`, async ($, on) => {
      const w = world(on, { herdr: true })
      setUp(w, id, 'ref-only')
      w.tabs = [{ tab_id: 'w1:t4', label: `${id}/${kind}`, agent_status: 'working' }]
      await startSession($)
      await pressChat($, await $.ui.mount(BAND), kind, id, w)
      expect(ran(w, 'jaira', 'pull', id)).toEqual([['jaira', 'pull', id, MAIN]])
      const spawned = w.ran.findIndex(r => r[0] === 'bash' && r[1] === SPAWN)
      expect(spawned).toBeGreaterThan(w.ran.findIndex(r => r[1] === 'pull'))
      const argv = w.ran[spawned]!
      expect(argv.at(-1)).toBe(MAIN)
      expect(argv[argv.indexOf('--copy') + 1]).toBe(PATH(id))
      expect(argv.indexOf('--copy')).toBeLessThan(argv.indexOf(id.toLowerCase()))
      expect(w.opened.filter(o => o.id === 'cockpit-message')).toEqual([])
      // The cockpit never removes it under a chat.
      finish(w, kind, id)
      w.files.delete(PATH(id))
      for (let i = 0; i < 4; i++) await w.clock.advance(60_000)
      expect(ran(w, 'git', 'clean')).toEqual([])
      expect(copyLine(w, id)).toEqual([])
    })

    test(`${kind}, no worktree, file behind: pulled before the chat opens, no copy to hand over`, async ($, on) => {
      const w = world(on, { herdr: true })
      setUp(w, id, 'behind')
      await startSession($)
      await pressChat($, await $.ui.mount(BAND), kind, id, w)
      expect(ran(w, 'jaira', 'pull', id)).toEqual([['jaira', 'pull', id, MAIN]])
      expect(ran(w, 'bash', SPAWN)[0]).not.toContain('--copy')
      finish(w, kind, id)
      for (let i = 0; i < 4; i++) await w.clock.advance(60_000)
      expect(copyLine(w, id)).toEqual([])
    })

    test(`${kind}: a copy the chat left behind is named beside the band a minute after the chat is done, once`, async ($, on) => {
      const w = world(on, { herdr: true })
      setUp(w, id, 'ref-only')
      w.tabs = [{ tab_id: 'w1:t4', label: `${id}/${kind}`, agent_status: 'working' }]
      await startSession($)
      await pressChat($, await $.ui.mount(BAND), kind, id, w)
      await w.clock.advance(60_000)
      finish(w, kind, id)
      await w.clock.advance(60_000)
      // Done this minute: the chat may still be writing its last line.
      expect(copyLine(w, id)).toEqual([])
      await w.clock.advance(60_000)
      await w.clock.advance(60_000)
      expect(copyLine(w, id)).toEqual([`${id}: the copy pulled into the main checkout for a chat is still there and hides its ref; once no chat on ${id} runs any more, remove it: git clean -fqx -- '${PATH(id)}'`])
      expect(ran(w, 'git', 'clean')).toEqual([])
    })

    test(`${kind}: a refused pull stops the press with its reason, opens no chat and never steals`, async ($, on) => {
      const w = world(on, { herdr: true })
      setUp(w, id, 'ref-only')
      w.pullFails = `jaira: ${id} is taken by Berk; --steal takes it anyway`
      await startSession($)
      await pressChat($, await $.ui.mount(BAND), kind, id, w)
      expect(ran(w, 'jaira', 'pull', id)).toEqual([['jaira', 'pull', id, MAIN]])
      expect(ran(w, 'bash', SPAWN)).toEqual([])
      expect(ran(w, 'git', 'clean')).toEqual([])
      expect(w.opened.at(-1)).toEqual({ id: 'cockpit-message', title: KINDS.find(k => k.kind === kind)!.title, focus: true })
      const pane = await $.ui.mount(MESSAGE(KINDS.find(k => k.kind === kind)!.title))
      expect(await pane.find({ type: 'Text', text: /--steal/ })).toBeDefined()
      expect(w.submitted).toEqual([])
    })
  }

  for (const { kind, id } of KINDS) {
    test(`without Herdr, ${kind} pulls and asks this session to remove the copy as its last step`, async ($, on) => {
      const w = world(on)
      setUp(w, id, 'ref-only')
      await startSession($)
      await pressChat($, await $.ui.mount(BAND), kind, id, w)
      expect(ran(w, 'jaira', 'pull', id)).toEqual([['jaira', 'pull', id, MAIN]])
      expect(w.submitted).toHaveLength(1)
      expect(w.submitted[0]).toEndWith(`As the very last step, once you write nothing more onto the ticket, decided or not: git clean -fqx -- '${PATH(id)}' - the copy of the ticket the cockpit pulled for this chat; left there, it hides the ticket's ref.`)
      finish(w, kind, id)
      w.files.delete(PATH(id))
      for (let i = 0; i < 4; i++) await w.clock.advance(60_000)
      expect(ran(w, 'git', 'clean')).toEqual([])
      expect(copyLine(w, id)).toEqual([])
    })
  }

  test('a ticket already on the disk: no copy is mentioned to the chat', async ($, on) => {
    const w = world(on)
    w.worktrees = ALONE
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(ran(w, 'jaira', 'pull')).toEqual([])
    expect(w.submitted[0]).not.toContain('git clean')
  })

  test('a chat that cannot be opened takes the pulled copy away at once', async ($, on) => {
    const w = world(on, { herdr: true, spawn: false })
    setUp(w, 'MG677R', 'ref-only')
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(w.opened.at(-1)?.title).toBe('Explain MG677R failed')
    expect(ran(w, 'git', 'clean')).toEqual([['git', 'clean', '-fqx', '--', PATH('MG677R'), MAIN]])
  })

  test('a spawn.sh that failed before its tab came up takes the pulled copy away at once', async ($, on) => {
    const w = world(on, { herdr: true })
    setUp(w, 'MG677R', 'ref-only')
    w.spawnFails = 'unknown flag: --copy'
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(w.opened.at(-1)?.title).toBe('Explain MG677R failed')
    expect(ran(w, 'git', 'clean')).toEqual([['git', 'clean', '-fqx', '--', PATH('MG677R'), MAIN]])
  })

  test('a spawn.sh that failed after its tab came up leaves the copy to that chat', async ($, on) => {
    const w = world(on, { herdr: true })
    setUp(w, 'MG677R', 'ref-only')
    w.spawnFails = 'claude did not come up in p1: none'
    w.tabs = [{ tab_id: 'w1:t4', label: 'MG677R/explain', agent_status: 'working' }]
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(w.opened.at(-1)?.title).toBe('Explain MG677R failed')
    expect(ran(w, 'git', 'clean')).toEqual([])
    // Still running: nothing is named under it.
    for (let i = 0; i < 3; i++) await w.clock.advance(60_000)
    expect(copyLine(w, 'MG677R')).toEqual([])
  })

  test('a second press while spawn.sh still waits for claude is refused: the chat holds the copy from before spawn.sh runs', async ($, on) => {
    const w = world(on, { herdr: true })
    setUp(w, 'MG677R', 'ref-only')
    w.tickets = [{ ...SIGNOFF, question: 'Branch in den Schluessel?' }, HUMAN, WORKING]
    const held: (() => void)[] = []
    w.heldSpawn = held
    await startSession($)
    const pane = await quietly($, await $.ui.mount(BAND), 'MG677R', w)
    const first = pane.press({ key: 'explain-all-MG677R' })
    await w.clock.settle()
    expect(held).toHaveLength(1)
    // The row unlocks after BUSY_MS while spawn.sh still waits.
    await w.clock.advance(21_000)
    await (await $.ui.mount(BAND)).press({ key: 'discuss-MG677R' })
    expect(ran(w, 'bash', SPAWN)).toHaveLength(1)
    expect(w.opened.at(-1)?.title).toBe('Discuss MG677R failed')
    held[0]!()
    await first
    await w.clock.settle()
    expect(ran(w, 'git', 'clean')).toEqual([])
  })

  test('a failed spawn.sh frees the ticket for the next chat', async ($, on) => {
    const w = world(on, { herdr: true })
    setUp(w, 'MG677R', 'ref-only')
    w.spawnFails = 'unknown flag: --copy'
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    w.spawnFails = undefined
    w.disk['@/repo'] = 'ref-only'
    await w.clock.advance(21_000)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(ran(w, 'bash', SPAWN)).toHaveLength(2)
  })

  for (const second of ['explain', 'discuss'] as const) {
    test(`a second chat (${second}) on a ticket whose copy an Explain chat holds is refused until that copy was looked at`, async ($, on) => {
      const w = world(on, { herdr: true })
      setUp(w, 'MG677R', 'ref-only')
      // MG677R with an open question too, so it offers Discuss beside Explain.
      w.tickets = [{ ...SIGNOFF, question: 'Branch in den Schluessel?' }, HUMAN, WORKING]
      w.tabs = [{ tab_id: 'w1:t4', label: 'MG677R/explain', agent_status: 'working' }]
      await startSession($)
      await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
      expect(ran(w, 'bash', SPAWN)).toHaveLength(1)
      const before = w.ran.length
      await pressChat($, await $.ui.mount(BAND), second, 'MG677R', w)
      expect(ran(w, 'bash', SPAWN)).toHaveLength(1)
      expect(w.ran.slice(before).filter(r => r[1] === 'pull')).toEqual([])
      expect(w.opened.at(-1)?.title).toBe(`${second === 'explain' ? 'Explain' : 'Discuss'} MG677R failed`)
      const pane = await $.ui.mount(MESSAGE(`${second === 'explain' ? 'Explain' : 'Discuss'} MG677R failed`))
      expect(await pane.find({ type: 'Text', text: /The Explain chat on MG677R is still open/ })).toBeDefined()
    })
  }
})

// 2V50Y5: the Acceptance pane's rules, read alone.
describe('the Acceptance pane, read alone', () => {
  const item = (cls: GapItem['cls'], choice: GapItem['choice'], follow?: string): GapItem => ({ text: `${cls}-${choice}`, cls, choice, ...(follow ? { follow } : {}) })
  // Every combination of what the notes say, words typed or not, and a lane
  // that ends the work or not: what Execute does.
  const TABLE: [string, GapItem[], string, string | null, string, number, number][] = [
    ['no notes, nothing typed', [], '', 'done', 'accept', 0, 0],
    ['no notes, words typed', [], 'fehlt', 'done', 'giveBack', 0, 0],
    ['all ok', [item('edge', 'accept'), item('limit', 'accept')], '', 'done', 'accept', 0, 0],
    ['all ok, words typed', [item('edge', 'accept')], 'noch was', 'done', 'giveBack', 0, 0],
    ['one to fix', [item('edge', 'now'), item('limit', 'accept')], '', 'done', 'giveBack', 1, 0],
    ['one to fix and a follow-up', [item('edge', 'now'), item('elsewhere', 'later')], '', 'done', 'giveBack', 1, 1],
    ['a follow-up with ok', [item('edge', 'later'), item('limit', 'accept')], '', 'done', 'accept', 0, 1],
    ['a follow-up made before', [item('edge', 'later', 'AB12CD')], '', 'done', 'accept', 0, 0],
    ['blocker to fix', [item('blocker', 'now')], '', 'done', 'giveBack', 1, 0],
    ['blocker kept, no reason', [item('blocker', 'accept')], '', 'done', 'needWhy', 0, 0],
    ['blocker kept with a reason', [item('blocker', 'accept')], 'ist bekannt', 'done', 'accept', 0, 0],
    ['blocker to a follow-up with a reason', [item('blocker', 'later')], 'später', 'done', 'accept', 0, 1],
    ['blocker kept, another to fix', [item('blocker', 'accept'), item('edge', 'now')], '', 'done', 'giveBack', 1, 0],
    ['all ok in a lane that does not end', [item('edge', 'accept')], '', null, 'notTerminal', 0, 0],
    ['to fix in a lane that does not end', [item('edge', 'now')], '', null, 'giveBack', 1, 0],
    ['words in a lane that does not end', [], 'fehlt', null, 'giveBack', 0, 0],
    ['blocker kept, words typed, in a lane that does not end', [item('blocker', 'accept')], 'fehlt', null, 'giveBack', 0, 0],
    ['blocker kept, nothing typed, in a lane that does not end', [item('blocker', 'later')], '', null, 'notTerminal', 0, 1],
  ]
  for (const [name, items, typed, to, step, fix, followUps] of TABLE) {
    test(`Execute: ${name} -> ${step}`, () => {
      const e = executionOf(items, typed, to)
      expect([e.step, e.fix, e.followUps]).toEqual([step, fix, followUps])
    })
  }

  test('the field asks why exactly while a blocker is kept and nothing is to fix', () => {
    expect(executionOf([item('blocker', 'accept')], '', 'done').asksWhy).toBe(true)
    expect(executionOf([item('blocker', 'later')], '', 'done').asksWhy).toBe(true)
    expect(executionOf([item('blocker', 'now')], '', 'done').asksWhy).toBe(false)
    expect(executionOf([item('blocker', 'accept'), item('edge', 'now')], '', 'done').asksWhy).toBe(false)
    expect(executionOf([item('edge', 'accept')], '', 'done').asksWhy).toBe(false)
    // Where nothing can be accepted, nothing asks why.
    expect(executionOf([item('blocker', 'accept')], '', null).asksWhy).toBe(false)
  })

  test('the story is read off one show: head, rounds, the DoD with superseded items, the verdict, when it was created', () => {
    const story = storyOf(JSON.stringify({
      title: 'cache  key', goal: 'a key per branch', outcome: { what: 'the key names the branch', why: 'two branches shared one' }, round: 3,
      back: { by: 'person', count: 2, from: 'signoff' }, created_at: '2026-11-20T10:00:00Z', review: { verdict: 'accept' },
      dod_items: [{ done: true, state: 'done' }, { done: false, state: 'todo' }, { done: false, state: 'superseded' }],
    }))
    expect(story).toEqual({
      title: 'cache key', goal: 'a key per branch', what: 'the key names the branch', why: 'two branches shared one', round: 3,
      back: { count: 2, from: 'signoff' }, dod: { done: 1, superseded: 1, total: 3 }, verdict: 'accept', created: Date.parse('2026-11-20T10:00:00Z'),
    })
    // Never back, no rounds counted, no date: nothing made up.
    expect(storyOf(`${OLDER}\n{"title": "t", "back": {"count": 0, "from": ""}}`)).toEqual({ title: 't', goal: '', what: '', why: '', round: 1, dod: { done: 0, superseded: 0, total: 0 }, verdict: '' })
    expect(storyOf('{broken')).toBeUndefined()
  })

  test('a release note gives its change lines, the mark to act on kept; commits and their files are counted once', () => {
    expect(releaseLines('<!-- x -->\n- one change\n\nprose\n- ! act on this\n-not a change\n')).toEqual(['one change', '! act on this'])
    expect(commitsOf('abc1234\na.ts\nb.ts\n\ndef5678\na.ts\n')).toEqual({ commits: 2, files: 2 })
    expect(commitsOf('')).toEqual({ commits: 0, files: 0 })
  })

  test('Explain on one note asks one sentence in the language, the note word for word; the answer is one line', () => {
    const ctx = { title: 'cache key', goal: '', what: '', diff: '' }
    const p = reasonPrompt('Русский', ctx, 'note: `x` in a.go:3')
    expect(p).toContain('In Русский')
    expect(p).toContain('note: `x` in a.go:3')
    expect(reasonPrompt(undefined, { ...ctx, title: 't' }, 'n')).toContain('the language the note is written in')
    // DoD 14: the goal, what was done and the diff go along, the diff cut short.
    const full = reasonPrompt('Deutsch', { title: 't', goal: 'a key per branch', what: 'the key names the branch', diff: '+const key = branch' }, 'note: register.tsx:12 x')
    for (const part of ["The ticket's goal: a key per branch", 'What was done: the key names the branch', '+const key = branch', 'name concretely what the code does']) expect(full).toContain(part)
    expect(reasonPrompt('Deutsch', ctx, 'n')).not.toContain('The change to the files')
    // A placeholder names no file; a file named twice is one.
    expect(notePaths('note: register.tsx:1214 reads core/release/notes/<ID>.md and hooks/board.ts, not register.tsx')).toEqual(['register.tsx', 'hooks/board.ts'])
    expect(notePaths('nothing named here')).toEqual([])
    expect(notePaths('see /home/a/repo/core/mod/hooks/register.tsx:9')).toEqual(['register.tsx'])
    expect(clipLines('a\nb\nc\n', 2)).toBe('a\nb\n... 1 more lines')
    expect(clipLines('a\nb', 2)).toBe('a\nb')
    expect(reasonOf('  Это\nважно. ')).toBe('Это важно.')
    expect(reasonOf('  ')).toBeUndefined()
  })

  test('an explanation a give-back came after is none, until a new one is written', () => {
    const note = (at: string, text: string) => `- **2026-11-22 ${at} · Alex** — ${text}`
    const show = (...notes: string[]) => JSON.stringify({ body: `# t\n\n## Progress\n${notes.join('\n')}\n` })
    expect(explanationNotes(show(note('10:00', 'explain: Status: ok')))).toEqual({ text: 'Status: ok' })
    expect(explanationNotes(show(note('10:00', 'explain: Status: ok'), note('10:05', 'Given back from signoff: beheben:')))).toBeUndefined()
    expect(explanationNotes(show(note('10:00', 'explain: old'), note('10:05', 'Given back from signoff: x'), note('11:00', 'explain: new')))).toEqual({ text: 'new' })
    // A note that only mentions a give-back is no give-back.
    expect(explanationNotes(show(note('10:00', 'explain: Status: ok'), note('10:05', 'testing: Given back from signoff earlier')))).toEqual({ text: 'Status: ok' })
  })
})

// 2V50Y5: what the Acceptance pane shows of a ticket, and what its buttons do
// beside Execute.
describe('the Acceptance pane', () => {
  const SETTINGS = '/home/a/.jaira/settings.json'
  const texts = async (pane: { findAll: (q: { type: 'Text' }) => Promise<{ text: string }[]> }) => (await pane.findAll({ type: 'Text' })).map(x => x.text)

  test('the head says what the ticket is for, what was done and why; the history how it got here', async ($, on) => {
    const w = world(on, { herdr: true })
    w.tickets = [{ ...SIGNOFF, goal: 'Two branches never share a cache entry.', review: { verdict: 'accept', gaps: 'note: a test is missing' } }, HUMAN, WORKING]
    w.shows.MG677R = {
      title: 'cache key', goal: 'Two branches never share a cache entry.', outcome: { what: 'the key names the branch', why: 'two branches shared one entry' },
      round: 2, back: { by: 'person', count: 1, from: 'signoff' }, created_at: new Date(1_000_000 - 3 * 86_400_000).toISOString(), review: { verdict: 'accept' },
      dod_items: [{ done: true, state: 'done' }, { done: true, state: 'done' }, { done: false, state: 'superseded' }],
    }
    w.hunts = { MG677R: { at: 'abc1234', state: 'ok' } }
    w.contains = { '@/.worktrees/repo-mg677r': ['abc1234'] }
    w.idLog = 'abc1234\nhooks/a.ts\nhooks/b.ts\n\ndef5678\nhooks/a.ts\n'
    w.files.set('/.worktrees/repo-mg677r/core/release/notes/MG677R.md', '<!-- x -->\n- Press Acceptance to decide a ticket.\n')
    w.prs = [{ number: 12 }]
    await startSession($)
    await w.clock.advance(60_000)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await w.clock.settle()
    const all = await texts(pane)
    expect(all).toContain('cache key')
    expect(all).toContain('┃ Goal: Two branches never share a cache entry.')
    expect(all).toContain('┃ Done: the key names the branch')
    expect(all).toContain('┃ Why: two branches shared one entry')
    expect(all).toContain('round 2 | back 1x from signoff | 3 days')
    expect(all).toContain('DoD 2/2 (1 superseded) | bughunt ✓ on HEAD | review: accept')
    expect(all).toContain('2 commits, 2 files | pull request exists')
    expect(all).toContain('Press Acceptance to decide a ticket.')
    // MK8HD2: each block opens with its title and a rule; each head field
    // has its own colour; what the ticket went through is dim, where it
    // stands is not.
    for (const title of ['History', 'Release note', 'Review notes - choose for each']) expect(all.some(x => x.startsWith(title) && x.endsWith('───'))).toBe(true)
    const colour = async (key: string) => JSON.stringify(await pane.find({ key }))
    expect(await colour('goal')).toContain('"color":"suggestion"')
    expect(await colour('what')).toContain('"color":"success"')
    expect(await colour('why')).toContain('"color":"warning"')
    const dim = async (text: string) => (await pane.find({ type: 'Text', text }))?.props.dimColor
    expect(await dim('round 2 | back 1x from signoff | 3 days')).toBe(true)
    // Where it stands wraps on a line of its own under the title, never cut
    // with the rule.
    expect(await dim('MG677R | signoff')).toBe(true)
    expect(await dim('2 commits, 2 files | pull request exists')).toBe(true)
    expect(await dim('DoD 2/2 (1 superseded) | bughunt ✓ on HEAD | review: accept')).toBeFalsy()
    // The light follows the marks: one 🟡 note makes it yellow.
    expect(all[0]).toBe('● yellow: 1 review note to weigh')
    // Looked up from the ticket's own worktree, without a fetch.
    expect(ran(w, 'git', 'log', '--format=%h')[0]?.at(-1)).toBe('@/.worktrees/repo-mg677r')
    expect(ran(w, 'git', 'fetch')).toEqual([])
  })

  test('Explain on one note asks the cheap model once for a sentence, shown under the note', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: a test is missing\nnote: pre-existing: the clock' } }, HUMAN, WORKING]
    w.files.set(SETTINGS, JSON.stringify({ language: 'Deutsch' }))
    w.reason = 'Unwichtig fürs Annehmen: es fehlt nur ein Test.'
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    expect(await pane.find({ key: 'why-MG677R-0' })).toBeDefined()
    expect(await pane.find({ key: 'why-MG677R-1' })).toBeDefined()
    await pane.press({ key: 'why-MG677R-0' })
    expect(w.reasonPrompts).toHaveLength(1)
    expect(w.reasonPrompts[0]).toContain('In Deutsch')
    expect(w.reasonPrompts[0]).toContain('note: a test is missing')
    expect(await texts(pane)).toContain('  Unwichtig fürs Annehmen: es fehlt nur ein Test.')
    // Asked once: a second press, or the pane opened again, asks nothing.
    await pane.press({ key: 'why-MG677R-0' })
    await (await acceptancePane($, await $.ui.mount(BAND), 'MG677R')).unmount()
    expect(w.reasonPrompts).toHaveLength(1)
    // Nothing goes onto the ticket.
    expect(ran(w, 'jaira', 'note')).toEqual([])
  })

  test('a model that gives no sentence says so under the note, and a later press asks again', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: a test is missing' } }, HUMAN, WORKING]
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'why-MG677R-0' })
    expect(await texts(pane)).toContain('  no explanation came back')
    // Without a language set, the sentence is asked in the note's own.
    expect(w.reasonPrompts[0]).toContain('the language the note is written in')
    // A failure is not kept: the next press asks again.
    w.reason = 'Only a test is missing.'
    await pane.press({ key: 'why-MG677R-0' })
    expect(w.reasonPrompts).toHaveLength(2)
    expect(await texts(pane)).toContain('  Only a test is missing.')
  })

  test('a sentence kept in one language is not shown once the language changed: the pane opened again asks in the new one', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: a test is missing' } }, HUMAN, WORKING]
    w.files.set(SETTINGS, JSON.stringify({ language: 'Deutsch' }))
    w.reason = 'Nur ein Test fehlt.'
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'why-MG677R-0' })
    w.files.set(SETTINGS, JSON.stringify({ language: 'Русский' }))
    w.reason = 'Не хватает только теста.'
    const again = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    expect((await texts(again)).some(x => x.includes('Nur ein Test fehlt.'))).toBe(false)
    await again.press({ key: 'why-MG677R-0' })
    expect(w.reasonPrompts.at(-1)).toContain('In Русский')
    expect(await texts(again)).toContain('  Не хватает только теста.')
  })

  test('Execute while an Explain chat still opens is said, not dropped, and writes nothing', async ($, on) => {
    const w = world(on, { herdr: true })
    w.heldSpawn = []
    await startSession($)
    const pane = await quietly($, await $.ui.mount(BAND), 'MG677R', w)
    void pane.press({ key: 'explain-all-MG677R' })
    await w.clock.settle()
    expect(w.heldSpawn).toHaveLength(1)
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R: an earlier press or an Explain chat is still at work on it - press Execute again once it is done')
    w.heldSpawn.shift()!()
    await w.clock.settle()
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
  })

  test('a known limit or a note for another ticket colours nothing: the light follows the marks', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, dod_complete: true, review: { verdict: 'accept', gaps: 'note: windows paths\nnote: pre-existing: load.go' } }, HUMAN, WORKING]
    w.bodies.MG677R = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: ok\n[known limit] the paths\n[not this ticket] load.go\n'
    w.hunts = { MG677R: { at: 'abc1234', state: 'ok' } }
    w.contains = { '@/.worktrees/repo-mg677r': ['abc1234'] }
    await startSession($)
    await w.clock.advance(60_000)
    const all = await texts(await acceptancePane($, await $.ui.mount(BAND), 'MG677R'))
    expect(all[0]).toBe('● green: accept')
    expect(all).toContain('⚪ note: windows paths')
    expect(all).toContain('🔵 note: pre-existing: load.go')
  })

  test('Explain on a note hands the model the goal, what was done and the diff of the file the note names', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: register.tsx:12 closes the pane early' } }, HUMAN, WORKING]
    w.shows.MG677R = { title: 'cache key', goal: 'a key per branch', outcome: { what: 'the key names the branch', why: 'x' } }
    w.patches = 'abc1234 fix(MG677R): x\n+  closeAcceptance($)\n'
    w.reason = 'Das Fenster schließt vor dem Schritt.'
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await w.clock.settle()
    await pane.press({ key: 'why-MG677R-0' })
    const prompt = w.reasonPrompts.at(-1) ?? ''
    expect(prompt).toContain("The ticket's goal: a key per branch")
    expect(prompt).toContain('What was done: the key names the branch')
    expect(prompt).toContain('+  closeAcceptance($)')
    // Only the ticket's own commits, the name matched at any depth.
    expect(ran(w, 'git', 'log', '-p')[0]?.slice(0, 9)).toEqual(['git', 'log', '-p', '--format=%h %s', '-F', '--grep=MG677R', 'HEAD', '--', ':(glob)**/register.tsx'])
  })

  test('a give-back the move refuses keeps the pane open with the person\'s marks', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: the size\nnote: the docs' } }, HUMAN, WORKING]
    w.moveFails = 'jaira: lane "in-progress" refuses: the board is locked'
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await cycle(pane, 'MG677R', 0)
    await pane.press({ key: 'execute-MG677R' })
    expect(w.closed).not.toContain('cockpit-acceptance')
    expect((await pane.find({ key: 'now-MG677R-0' }))?.text).toBe('[•fix]')
    expect(ran(w, 'jaira', 'note')).toEqual([])
    // Once the move goes through, the pane closes.
    w.moveFails = undefined
    await pane.press({ key: 'execute-MG677R' })
    expect(w.closed).toContain('cockpit-acceptance')
  })

  test('a repository that keeps no release notes shows no release note line', async ($, on) => {
    const w = world(on)
    await startSession($)
    const all = await texts(await acceptancePane($, await $.ui.mount(BAND), 'MG677R'))
    await w.clock.settle()
    expect(all.some(x => x.startsWith('Release note'))).toBe(false)
  })

  test('the pull request is asked once in five minutes, however often the pane opens', async ($, on) => {
    const w = world(on)
    await startSession($)
    for (let i = 0; i < 3; i++) {
      await (await acceptancePane($, await $.ui.mount(BAND), 'MG677R')).unmount()
      await w.clock.settle()
    }
    expect(ran(w, 'gh', 'pr', 'list')).toHaveLength(1)
    await w.clock.advance(5 * 60_000)
    await (await acceptancePane($, await $.ui.mount(BAND), 'MG677R')).unmount()
    await w.clock.settle()
    expect(ran(w, 'gh', 'pr', 'list')).toHaveLength(2)
  })

  test('a forge that does not answer is asked again at the next opening, not taken for no pull request', async ($, on) => {
    const w = world(on)
    w.prListFails = true
    await startSession($)
    for (let i = 0; i < 2; i++) {
      await (await acceptancePane($, await $.ui.mount(BAND), 'MG677R')).unmount()
      await w.clock.settle()
    }
    expect(ran(w, 'gh', 'pr', 'list')).toHaveLength(2)
  })

  test('two notes explained one after the other both keep their sentence in the store', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: a\nnote: b' } }, HUMAN, WORKING]
    w.reason = 'one sentence'
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await Promise.all([pane.press({ key: 'why-MG677R-0' }), pane.press({ key: 'why-MG677R-1' })])
    await w.clock.settle()
    expect(Object.keys((w.store.get('reason:MG677R') ?? {}) as object)).toHaveLength(2)
  })

  test('Explain all opens the explaining chat from the pane', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    await pressExplain($, await $.ui.mount(BAND), 'MG677R', w)
    expect(ran(w, 'bash', SPAWN)[0]?.slice(0, 6)).toEqual(['bash', SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'explain'])
  })

  // DoD 10: an explanation of an earlier round is gone once the ticket was
  // given back, and back with a new one.
  test('Read explanation goes once the ticket was given back after it, and comes back with a new explanation', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = '## Progress\n- **2026-11-22 10:00 · Alex** — explain: Status: DoD 3/3 ok\n[small edge] the size\n'
    await startSession($)
    expect(await (await $.ui.mount(BAND)).find({ key: 'read-MG677R' })).toBeDefined()
    w.bodies.MG677R += '- **2026-11-22 10:05 · Alex** — Given back from signoff: beheben:\n- note: the size\n'
    await w.clock.advance(60_000)
    expect(await (await $.ui.mount(BAND)).find({ key: 'read-MG677R' })).toBeUndefined()
    w.bodies.MG677R += '- **2026-11-22 11:00 · Alex** — explain: Status: DoD 3/3 ok\n[known limit] the size\n'
    await w.clock.advance(60_000)
    expect(await (await $.ui.mount(BAND)).find({ key: 'read-MG677R' })).toBeDefined()
  })

  test('the classes of an explanation a give-back came after no longer class the notes', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: the size' } }, HUMAN, WORKING]
    w.bodies.MG677R = '## Progress\n- **2026-11-22 10:00 · Alex** — explain: Status: ok\n[known limit] the size\n- **2026-11-22 10:05 · Alex** — Given back from signoff: dazu: more\n'
    await startSession($)
    const all = await texts(await acceptancePane($, await $.ui.mount(BAND), 'MG677R'))
    expect(all).toContain('🟡 note: the size')
  })
})

// Where Accept moves a ticket (VHGWQF), by input class: its own worktree
// first, then the worktree of the branch carrying its commits, never a guess.
describe('where Accept moves a ticket', () => {
  const OTHER = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n\nworktree /.worktrees/repo-mg677r\nHEAD 2\nbranch refs/heads/feat/mg677r\n'
  const OWN = OTHER + '\nworktree /.worktrees/repo-zz99zz\nHEAD 3\nbranch refs/heads/feat/zz99zz\n'
  const cases: [string, string, Parameters<typeof acceptTreeFor>[2], ReturnType<typeof acceptTreeFor>][] = [
    ['its own worktree, whatever carries the commits', OWN, { branch: 'feat/mg677r', by: 'commit' }, { dir: '/.worktrees/repo-zz99zz' }],
    ['its own worktree, even where the carrier is ambiguous', OWN, { ambiguous: ['feat/mg677r', 'wip'] }, { dir: '/.worktrees/repo-zz99zz' }],
    ["another ticket's branch holds its commits: that worktree", OTHER, { branch: 'feat/mg677r', by: 'commit' }, { dir: '/.worktrees/repo-mg677r' }],
    ['the carrier has no worktree: the main checkout', OTHER, { branch: 'feat/zz99zz', by: 'name' }, {}],
    ['no carrier at all: the main checkout', OTHER, undefined, {}],
    ['several could carry it: named, nothing picked', OTHER, { ambiguous: ['feat/mg677r', 'wip'] }, { ambiguous: ['feat/mg677r', 'wip'] }],
  ]
  for (const [name, porcelain, carrier, want] of cases) {
    test(name, () => {
      expect(acceptTreeFor(porcelain, 'ZZ99ZZ', carrier)).toEqual(want)
    })
  }
})

// DKMG91: a 'seen by a person' item of the definition of done stands in the
// Acceptance pane as a section of its own, above the review notes; the person
// ticks it there, and Accept ticks it on the ticket before it moves it.
describe('what the person checks themselves', () => {
  const texts = async (pane: { findAll: (q: { type: 'Text' }) => Promise<{ text: string }[]> }) => (await pane.findAll({ type: 'Text' })).map(x => x.text)
  const SEEN = [
    { n: 1, done: true, state: 'done', text: 'the key names the branch' },
    { n: 2, done: false, state: 'todo', text: 'seen by a person: the band shows the branch' },
    { n: 3, done: true, state: 'done', text: 'seen by a person: already ticked' },
    { n: 4, done: false, state: 'superseded', text: 'seen by a person: superseded' },
  ]
  const setUp = (on: On) => {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: 'accept', gaps: 'note: a test is missing' } }, HUMAN, WORKING]
    w.shows.MG677R = { title: 'cache key', dod_items: SEEN }
    w.me = { identity: 'Alex', aliases: [] }
    return w
  }

  test('open items stand above the review notes, each with its text; ticked and superseded ones do not', async ($, on) => {
    setUp(on)
    await startSession($)
    const all = await texts(await acceptancePane($, await $.ui.mount(BAND), 'MG677R'))
    const heading = all.indexOf('Check yourself - tick what you have seen')
    expect(heading).toBeGreaterThan(-1)
    // MK8HD2's section heading: the title, then a dim rule as a text of its own.
    const item = all.indexOf('[ ] 2. seen by a person: the band shows the branch')
    expect(item).toBeGreaterThan(heading)
    expect(all.slice(heading + 1, item).every(x => /^\s*─+$/.test(x))).toBe(true)
    expect(all.indexOf('Review notes - choose for each')).toBeGreaterThan(item)
    expect(all.some(x => x.includes('already ticked') || x.includes('person: superseded'))).toBe(false)
    expect(all).toContain('Execute: tick what you checked above first, then done')
  })

  test('a ticket without such items shows no section, and Execute accepts as before', async ($, on) => {
    const w = setUp(on)
    w.shows.MG677R = { title: 'cache key', dod_items: [SEEN[0], SEEN[2]] }
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    const all = await texts(pane)
    expect(all.some(x => x.startsWith('Check yourself'))).toBe(false)
    expect(all).toContain('Execute: done')
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'dod')).toEqual([])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
  })

  test('Execute with an item not ticked accepts nothing and names the item', async ($, on) => {
    const w = setUp(on)
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'dod')).toEqual([])
    expect(ran(w, 'jaira', 'move').filter(r => r[1] === 'MG677R')).toEqual([])
    expect(await texts(pane)).toContain('MG677R: not accepted - one item is still to check yourself and tick:\n- seen by a person: the band shows the branch')
    expect(w.closed).not.toContain('cockpit-acceptance')
  })

  test('ticked, one Execute ticks it on the ticket with the person named, then accepts', async ($, on) => {
    const w = setUp(on)
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'seen-MG677R-2' })
    expect(await texts(pane)).toContain('[✓] 2. seen by a person: the band shows the branch')
    expect(await texts(pane)).toContain('Execute: done')
    await pane.press({ key: 'execute-MG677R' })
    expect(ran(w, 'jaira', 'dod').map(r => r.slice(0, -1))).toEqual([['jaira', 'dod', 'MG677R', '2', '--done', '--proof', 'seen by Alex in the cockpit']])
    const order = w.ran.filter(r => r[0] === 'jaira' && ['dod', 'move'].includes(r[1]!) && r[2] === 'MG677R').map(r => (r[1] === 'dod' ? 'dod' : r.includes('--dry-run') ? 'dry-run' : 'move'))
    expect(order).toEqual(['dod', 'dry-run', 'move'])
    expect(w.closed).toContain('cockpit-acceptance')
  })

  test('a tick pressed twice is taken back', async ($, on) => {
    setUp(on)
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'seen-MG677R-2' })
    await pane.press({ key: 'seen-MG677R-2' })
    expect(await texts(pane)).toContain('[ ] 2. seen by a person: the band shows the branch')
  })

  // The pane's list is read when it opens: what changed on the ticket since
  // is read again before anything is written (bughunt 23.11., two findings).
  const TWO = [SEEN[1], { n: 5, done: false, state: 'todo', text: 'seen by a person: the docs read well' }]
  const ticked = async ($: Test$, w: World, items: unknown[]) => {
    w.shows.MG677R = { title: 'cache key', dod_items: TWO }
    await startSession($)
    const pane = await acceptancePane($, await $.ui.mount(BAND), 'MG677R')
    await pane.press({ key: 'seen-MG677R-2' })
    await pane.press({ key: 'seen-MG677R-5' })
    w.shows.MG677R = { title: 'cache key', dod_items: items }
    await pane.press({ key: 'execute-MG677R' })
    return pane
  }

  test('a refusal after some items were ticked names them instead of saying nothing changed', async ($, on) => {
    const w = setUp(on)
    w.dodFails = '5'
    const pane = await ticked($, w, TWO)
    expect(await texts(pane)).toContain('MG677R was not accepted; only definition-of-done item 2 was ticked as seen by you:\njaira: the board is locked')
    expect(ran(w, 'jaira', 'move').filter(r => r[1] === 'MG677R')).toEqual([])
  })

  test('a gate refusing after the ticks names them too', async ($, on) => {
    const w = setUp(on)
    w.refusals = [...w.refusals, { code: 'needs_nonmodel_signal', message: 'criterion 1 is still open' }]
    const pane = await ticked($, w, TWO)
    expect(await texts(pane)).toContain('MG677R was not accepted; only definition-of-done items 2, 5 were ticked as seen by you:\ncriterion 1 is still open')
  })

  for (const [how, changed] of [['reworded', { ...TWO[1], text: 'seen by a person: the docs and the help read well' }], ['superseded', { ...TWO[1], state: 'superseded' }]] as const) {
    test(`an item ${how} since the pane opened ticks nothing and accepts nothing`, async ($, on) => {
      const w = setUp(on)
      const pane = await ticked($, w, [TWO[0], changed])
      expect(ran(w, 'jaira', 'dod')).toEqual([])
      expect(ran(w, 'jaira', 'move').filter(r => r[1] === 'MG677R')).toEqual([])
      expect(await texts(pane)).toContain('MG677R was not accepted, nothing was changed:\nthis item changed on the ticket since the pane opened - close it, open it again and look once more:\n- seen by a person: the docs read well')
    })
  }

  test('an item ticked meanwhile keeps its own proof', async ($, on) => {
    const w = setUp(on)
    await ticked($, w, [TWO[0], { ...TWO[1], done: true, state: 'done', proof: 'seen by Berk' }])
    expect(ran(w, 'jaira', 'dod').map(r => r[3])).toEqual(['2'])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
  })

  test('seenToTick: open and unchanged ticks, done meanwhile skips, gone or changed stops', () => {
    const shown = (items: unknown[]) => JSON.stringify({ dod_items: items })
    const a = { n: 1, text: 'seen by a person: a', checked: true }
    expect(seenToTick([a], shown([{ n: 1, done: false, state: 'todo', text: 'seen by a person:  a' }]))).toEqual({ tick: [1], changed: [] })
    expect(seenToTick([a], shown([{ n: 1, done: true, state: 'done', text: 'seen by a person: a' }]))).toEqual({ tick: [], changed: [] })
    expect(seenToTick([a], shown([]))).toEqual({ tick: [], changed: ['seen by a person: a'] })
  })

  test('executionOf waits for the items; a note to fix still gives back', () => {
    expect(executionOf([], '', 'done', 1).step).toBe('needSeen')
    expect(executionOf([], '', 'done', 0).step).toBe('accept')
    expect(executionOf([], 'missing', 'done', 1).step).toBe('giveBack')
    expect(executionOf([], '', null, 1).step).toBe('notTerminal')
  })
})

// Every form git takes for a remote, taken apart as core/pr's parseRemote
// does (XSGWPX): the user before the @ and the port are no part of the host.
describe('parseRemote', () => {
  const cases: [string, { host: string; repo: string } | undefined][] = [
    ['https://github.com/me/jaira.git', { host: 'github.com', repo: 'me/jaira' }],
    ['https://token@github.com/me/jaira.git/', { host: 'github.com', repo: 'me/jaira' }],
    ['https://ghe.example.com:8443/team/app.git', { host: 'ghe.example.com', repo: 'team/app' }],
    ['ssh://alice@ghe.example.com:2222/team/app.git', { host: 'ghe.example.com', repo: 'team/app' }],
    ['git+ssh://git@GHE.Example.com/team/app', { host: 'ghe.example.com', repo: 'team/app' }],
    ['git@github.com:me/jaira', { host: 'github.com', repo: 'me/jaira' }],
    ['alice@ghe.example.com:team/app.git', { host: 'ghe.example.com', repo: 'team/app' }],
    ['ghe.example.com:team/app.git', { host: 'ghe.example.com', repo: 'team/app' }],
    ['git@gitlab.example.com:group/sub/project.git', { host: 'gitlab.example.com', repo: 'group/sub/project' }],
    ['/srv/git/app.git', undefined],
    ['C:/repos/app', undefined],
    ['C:\\repos\\app', undefined],
    ['https://ghe.example.com/app', undefined],
    ['git@host:app.git', undefined],
    ['file:///srv/git/app.git', undefined],
  ]
  for (const [url, want] of cases) {
    test(url, () => {
      expect(parseRemote(url)).toEqual(want)
    })
  }
})

// A host gh takes for github.com goes to gh as it is, with no host of its own.
describe('ghRepoOf', () => {
  const cases: [string, { repo: string; host?: string }][] = [
    ['alice@ghe.example.com:team/app.git', { repo: 'ghe.example.com/team/app', host: 'ghe.example.com' }],
    ['git@github.com:me/jaira.git', { repo: 'git@github.com:me/jaira.git' }],
    ['ssh://git@ssh.github.com:443/me/jaira.git', { repo: 'ssh://git@ssh.github.com:443/me/jaira.git' }],
    ['https://www.GitHub.com/me/jaira', { repo: 'https://www.GitHub.com/me/jaira' }],
    ['/srv/git/app.git', { repo: '/srv/git/app.git' }],
  ]
  for (const [url, want] of cases) {
    test(url, () => {
      expect(ghRepoOf(url)).toEqual(want)
    })
  }
})
