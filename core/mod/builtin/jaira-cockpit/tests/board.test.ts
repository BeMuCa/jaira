import type { On, SessionRateLimit } from 'claude-code'
import { describe, expect, mock, test } from 'claude-code/testing'
import { languageOf, lightOf, rephrasedOf, rephrasePrompt, snapshotOf, type Waiting } from '../hooks/board.ts'
import { reportLine } from '../hooks/digest.ts'
import { durationOf, stoneSize } from '../hooks/overview.ts'

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
  // `git worktree list --porcelain`, WORKTREES unless set.
  worktrees?: string
  // `jaira create` of a follow-up fails.
  createFails?: boolean
  // Another session's move came first: jaira answers this one 'already in'.
  movedFirst?: boolean
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
const GIVE = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-give-back', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as const
const EXPLANATION = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-explanation', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never
const ANSWER = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-answer', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never
const VERDICT = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-verdict', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never
const MESSAGE = (title: string) =>
  ({ plugin: 'jaira-cockpit', surface: 'terminal', component: 'Pane', requestId: 'cockpit-message', props: { title, isFocused: true, bodyColumns: 80, placement: 'dock', scroll: { offset: 0, bodyRows: 20 }, view: {} } }) as never

// Herdr's review worktree of feat/mg677r is listed first and ends in the id
// as well, without the ticket on its disk: every press must pass it by.
const WORKTREES = 'worktree /repo\nHEAD 1\nbranch refs/heads/master\n\nworktree /.herdr/worktrees/repo/review-feat-mg677r\nHEAD 2\nbranch refs/heads/review/feat/mg677r\n\nworktree /.worktrees/repo-mg677r\nHEAD 2\nbranch refs/heads/feat/mg677r\n'
const SPAWN = '/repo/.claude/skills/jaira-dispatcher/scripts/spawn.sh'
// The release cache jaira writes under ~/.jaira, HOME being /home/a.
const CACHE = '/home/a/.jaira/update-check.json'

function world(on: On, opts: { board?: boolean; herdr?: boolean; store?: Record<string, unknown>; spawn?: boolean; env?: Record<string, string> } = {}): World {
  const w: World = {
    ran: [],
    files: new Map(),
    usage: { tokens: 0, rateLimits: [{ kind: 'five_hour', percentUsed: 34 }] },
    tickets: [SIGNOFF, HUMAN, WORKING],
    tabs: [],
    paneGone: false,
    held: undefined,
    agents: [],
    oldJaira: false,
    boards: {},
    asked: [],
    answer: '',
    model: '',
    modelCalls: 0,
    prompts: [],
    filled: [],
    store: new Map(Object.entries(opts.store ?? {})),
    status: [],
    refusals: [{ code: 'needs_human', message: 'lane "signoff" is a human checkpoint' }],
    warns: false,
    version: '0.3.7',
    noteFails: false,
    toasts: [],
    bodies: {},
    hunts: {},
    contains: {},
    submitted: [],
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
  on('fs.exists', ($, e) => ({ value: (hasBoard && e.path === '/repo/.jaira') || (e.path === SPAWN && opts.spawn !== false) || w.files.has(e.path) }))
  on('fs.read', ($, e) => (w.files.has(e.path) ? { value: w.files.get(e.path)! } : { deny: 'no such file' }))
  on('fs.write', ($, e) => {
    w.files.set(e.path, e.text)
    return { value: undefined }
  })
  on('process.run', async ($, e) => {
    w.ran.push([...e.argv, `@${e.init?.cwd ?? ''}`])
    if (w.held && e.argv[0] === 'jaira' && e.argv[1] === 'move') await new Promise<void>(go => w.held!.push(go))
    if (w.heldTabs && e.argv[1] === 'tab' && e.argv[2] === 'list') await new Promise<void>(go => w.heldTabs!.push(go))
    const ok = (stdout = '') => ({ value: { exitCode: 0, stdout, stderr: '', isStdoutTruncated: false, isStderrTruncated: false } })
    const [cmd, ...args] = e.argv
    const on = args.includes('--board') ? w.boards[args[args.indexOf('--board') + 1] ?? ''] : undefined
    if (cmd === 'jaira' && args[0] === 'boards') return ok(JSON.stringify({ boards: ['default', ...Object.keys(w.boards)].map(name => ({ name })) }))
    if (cmd === 'jaira' && args.includes('--board') && (on === undefined || (args[0] === 'lanes' && on.lanes === undefined))) return failed('jaira: unknown flag: --board')
    if (cmd === 'jaira' && args[0] === 'lanes' && on !== undefined) return ok(JSON.stringify({ lanes: on.lanes }))
    if (cmd === 'jaira' && args[0] === 'list' && on !== undefined) return ok(JSON.stringify({ count: on.tickets.length, tickets: on.tickets }))
    if (cmd === 'jaira' && args[0] === 'lanes') return ok(JSON.stringify(w.oldJaira ? { lanes: LANES.lanes.map(({ human_exit, requires_question, ...l }) => l) } : LANES))
    if (cmd === 'jaira' && args[0] === 'list') return ok(JSON.stringify({ count: w.tickets.length, tickets: w.tickets }))
    if (cmd === 'jaira' && args[0] === 'move' && w.movedFirst && !args.includes('--dry-run')) return ok(`${args[1]} is already in ${args[3]}\n`)
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
    if (cmd === 'jaira' && args[0] === 'whoami') return ok(w.me === undefined ? '' : JSON.stringify(w.me))
    if (cmd === 'jaira' && args[0] === 'show' && args.includes('--brief')) {
      const id = args[1]!
      return ok((w.warns ? `${OLDER}\n` : '') + JSON.stringify({ id, ...(id in w.hunts ? { bughunt: w.hunts[id] } : {}) }))
    }
    if (cmd === 'jaira' && args[0] === 'show') {
      const status = (w.tickets.find(x => (x as { handle: string }).handle === args[1]) as { status?: string } | undefined)?.status
      return ok((w.warns ? `${OLDER}\n` : '') + JSON.stringify({ handle: args[1], body: w.bodies[args[1]!] ?? '', ...(status !== undefined ? { status, board: '' } : {}) }))
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
    if (cmd === 'git' && args[0] === 'show' && args[1]?.endsWith(':.jaira/pr')) return w.basePr === undefined ? failed(`fatal: path '.jaira/pr' does not exist in '${args[1].split(':')[0]}'`) : ok(w.basePr)
    if (cmd === 'git' && args[0] === 'branch' && args[1] === '--contains') return ok(w.holders.join('\n'))
    if (cmd === 'git' && args[0] === 'branch') return ok(w.branches.join('\n'))
    // Newest first, as git prints them. A --grep reads the whole message,
    // body included; --format=%H %s prints the subject alone.
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
  on('model.complete', ($, e) => {
    w.modelCalls++
    w.prompts.push((e as { prompt: string }).prompt)
    if (w.modelDown) return { value: { isAnswered: false } } as never
    return { value: { isAnswered: true, text: w.model, usage: { input_tokens: 1, output_tokens: 1, cache_creation_input_tokens: 0, cache_read_input_tokens: 0 } } }
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

// The Verdict pane of a ticket, opened from the band: what the row leaves out.
async function verdictOf($: Test$, band: { press: (q: { key: string }) => Promise<unknown> }, id: string) {
  await band.press({ key: `verdict-${id}` })
  const pane = await $.ui.mount(VERDICT(`Review verdict of ${id}`))
  const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
  await pane.unmount()
  return { texts, find: async (q: { type: 'Text'; text: string | RegExp }) => texts.find(x => (typeof q.text === 'string' ? x === q.text : q.text.test(x))) }
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
    expect(texts.slice(title, title + 2)).toEqual(['cache key', 'Goal: Two branches never share a cache entry.'])
    for (const text of ['? keep the old key?', `Verdict: ${verdict}`, '🔴 blocker: the key ignores the branch', '🟡 note: a test is missing', '🟡 note: a comment is stale']) {
      expect(texts).toContain(text)
    }
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
    expect(await ui.find({ key: 'accept-S0EG0M' })).toBeDefined()
    expect(await ui.find({ key: 'accept-V1D3OQ' })).toBeUndefined()
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
    await ui.press({ key: 'accept-MG677R' })
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
    await ui.press({ key: 'accept-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R')[0]?.at(-1)).toBe('@/.worktrees/repo-mg677r')
  })

  test('Accept, Give back and Explain are told apart by the theme\'s colours; Read explanation and Verdict stay plain', async ($, on) => {
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
    expect(await colour('accept-MG677R')).toBe('success')
    // Not the accent colour of a primary Button around it: green is the mark.
    expect((await ui.find({ key: 'accept-MG677R' }))?.props.variant).toBeUndefined()
    expect(await colour('back-MG677R')).toBe('error')
    expect(await colour('explain-MG677R')).toBe('warning')
    expect(await colour('read-MG677R')).toBeUndefined()
    expect(await colour('verdict-MG677R')).toBeUndefined()
  })

  test('on the terminal the coloured buttons draw their brackets in the label\'s colour; elsewhere the native button stays', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = '## Progress\n\n- **2026-10-12 15:00 · Alex** — explain: The key leaves out the branch.'
    await startSession($)
    for (const [surface, plain, text] of [['terminal', true, '[ Accept ]'], ['desktop', undefined, 'Accept']] as const) {
      const ui = await $.ui.mount({ ...BAND, surface })
      const accept = await ui.find({ key: 'accept-MG677R' })
      expect(accept?.props.plain).toBe(plain)
      expect(accept?.props.label).toBe('Accept')
      expect(accept?.children).toEqual([{ type: 'Text', props: { color: 'success', dimColor: false }, children: [text] }])
      await ui.unmount()
    }
    const ui = await $.ui.mount(BAND)
    expect((await ui.find({ key: 'back-MG677R' }))?.text).toBe('[ Give back ]')
    expect((await ui.find({ key: 'explain-MG677R' }))?.text).toBe('[ Explain ]')
    // Read explanation and Verdict keep the engine's own brackets.
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
    void ui.press({ key: 'accept-MG677R' })
    await w.clock.settle()
    for (const key of ['accept-MG677R', 'back-MG677R', 'explain-MG677R']) expect(await dimmed(ui, key)).toBe(true)
    // Another ticket's row is not this press's.
    expect(await dimmed(ui, 'back-K7P2QX')).toBe(false)
    await ui.press({ key: 'back-MG677R' })
    void ui.press({ key: 'accept-MG677R' })
    await w.clock.settle()
    expect(w.opened).toEqual([])
    await release(w)
    await release(w)
    expect(ran(w, 'jaira', 'move', 'MG677R')).toHaveLength(2)
    expect(await dimmed(ui, 'accept-MG677R')).toBe(false)
  })

  test('a row whose action has not ended is free again after 20 seconds, and the late end frees only its own lock', async ($, on) => {
    const w = world(on)
    w.held = []
    w.stay = true
    await startSession($)
    const ui = await $.ui.mount(BAND)
    void ui.press({ key: 'accept-MG677R' })
    await w.clock.advance(19_999)
    expect(await dimmed(ui, 'accept-MG677R')).toBe(true)
    await w.clock.advance(1)
    expect(await dimmed(ui, 'accept-MG677R')).toBe(false)
    void ui.press({ key: 'accept-MG677R' })
    await w.clock.settle()
    // The first press's dry run and move answer: it ends, the second waits on.
    await release(w)
    await release(w)
    expect(await dimmed(ui, 'accept-MG677R')).toBe(true)
    await release(w)
    await release(w)
    expect(await dimmed(ui, 'accept-MG677R')).toBe(false)
  })

  test('a reload frees a row whose press it cut off', async ($, on) => {
    const w = world(on)
    w.held = []
    await startSession($)
    const ui = await $.ui.mount(BAND)
    void ui.press({ key: 'accept-MG677R' })
    await w.clock.settle()
    expect(await dimmed(ui, 'accept-MG677R')).toBe(true)
    await startSession($)
    expect(await dimmed(ui, 'accept-MG677R')).toBe(false)
    // The cut-off press ends before the test does.
    await release(w)
    await release(w)
    await w.clock.advance(20_000)
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
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'explain', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('an open question shows under its ticket with Discuss', async ($, on) => {
    const w = world(on, { herdr: true })
    await startSession($)
    const ui = await $.ui.mount(BAND)
    // The row marks the question; the question itself is in the Verdict pane.
    expect(await ui.find({ text: /Popup oder eigene Seite/ })).toBeUndefined()
    expect(await ui.find({ type: 'Text', text: /K7P2QX \| human \| .* \| \?/ })).toBeDefined()
    await ui.press({ key: 'verdict-K7P2QX' })
    expect(await (await $.ui.mount(VERDICT('Review verdict of K7P2QX'))).find({ type: 'Text', text: '? Popup oder eigene Seite?' })).toBeDefined()
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
    await ui.press({ key: 'explain-MG677R' })
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
    await ui.press({ key: 'explain-MG677R' })
    await w.clock.advance(300)
    await ui.press({ key: 'explain-MG677R' })
    expect(w.submitted).toHaveLength(2)
    // Pressed again after the lock, it is asked again.
    await w.clock.advance(20_000)
    await ui.press({ key: 'explain-MG677R' })
    expect(w.submitted).toHaveLength(3)
  })

  test('with Herdr but no spawn.sh, Explain says spawn.sh is missing', async ($, on) => {
    const w = world(on, { herdr: true, spawn: false })
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'explain-MG677R' })
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
    await (await $.ui.mount(BAND)).press({ key: 'explain-MG677R' })
    expect(ran(w, 'C:/Program Files/Git/bin/bash.exe', home)).toEqual([['C:/Program Files/Git/bin/bash.exe', home, '--no-worktree', 'mg677r', 'MG677R', 'explain', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('on native Windows, CLAUDE_CODE_GIT_BASH_PATH names the bash', async ($, on) => {
    const env = { OS: 'Windows_NT', HOME: '/repo', HERDR_ENV: '1', CLAUDE_CODE_GIT_BASH_PATH: 'D:\\Git\\bin\\bash.exe' }
    const w = world(on, { herdr: true, env })
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'explain-MG677R' })
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
    expect(w.submitted).toEqual(["Decision on MG677R in the cockpit's explanation pane, noted on the ticket and not acted on: b) accept it as it is"])
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
    expect(w.submitted).toEqual(["jaira-cockpit: the person gave MG677R back in the cockpit (signoff -> in-progress), no worker started (the ticket's note says why): 1) add the branch to the key"])
  })

  test('an option marked [accept] accepts the ticket as the band\'s Accept does', async ($, on) => {
    const w = world(on)
    w.bodies.MG677R = MARKED
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'read-MG677R' })
    await (await $.ui.mount(EXPLANATION('Explanation of MG677R'))).press({ key: 'pick-MG677R-2' })
    expect(ran(w, 'jaira', 'note').map(r => r[3])).toEqual(['decision: 2) accept it as it is'])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(w.submitted).toEqual(['jaira-cockpit: the person accepted MG677R in the cockpit (signoff -> done)'])
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
    for (const id of ['OLD123', 'NEW456', 'BLK789']) expect(await band.find({ key: `explain-${id}` })).toBeDefined()
    // The row marks it; the Verdict pane says it whole.
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
    // the blocker kept beside an earlier verdict offers Explain; notes alone do not
    expect(await band.find({ key: 'explain-BLK456' })).toBeDefined()
    expect(await band.find({ key: 'explain-OLD123' })).toBeUndefined()
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

  test('Explain is offered for a blocker or a reservation; notes alone are listed without a button', async ($, on) => {
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
    for (const id of ['BL0CK1', 'BL0CK2', 'BL0CK3', 'BL0CK4', 'BL0CK5', 'BL0CK6', 'RESRV1']) expect(await open.find({ key: `explain-${id}` })).toBeDefined()
    for (const id of ['NOTES1', 'NOTES2', 'NOTES3', 'AB12CD', 'EF34GH', 'GH56JK', 'MN12PQ', 'PQ34RS', 'LM90NP']) expect(await open.find({ key: `explain-${id}` })).toBeUndefined()
    // The row counts them; the notes themselves stand in the Verdict pane.
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

  test('the whole review verdict is read from the band, with no explanation on the ticket', async ($, on) => {
    const w = world(on)
    const long = 'The diff meets all five DoD items; the one weak spot is the cache key, which leaves out the branch, so two branches share an entry.'
    w.tickets = [{ ...SIGNOFF, review: { verdict: long } }, { ...HUMAN, question: '' }]
    await startSession($)
    const band = await $.ui.mount(BAND)
    expect(await band.find({ key: 'read-MG677R' })).toBeUndefined()
    // No verdict, no notes, no question: no button.
    expect(await band.find({ key: 'verdict-K7P2QX' })).toBeUndefined()
    await band.press({ key: 'verdict-MG677R' })
    expect(w.opened).toEqual([{ id: 'cockpit-verdict', title: 'Review verdict of MG677R', focus: true }])
    const pane = await $.ui.mount(VERDICT('Review verdict of MG677R'))
    expect(await pane.find({ type: 'Text', text: `Verdict: ${long}` })).toBeDefined()
    // It starts with the light, then what the ticket is for. No hunt was
    // read for this ticket: nobody can tell, so yellow.
    const texts = (await pane.findAll({ type: 'Text' })).map(x => x.text)
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
    await (await $.ui.mount(BAND)).press({ key: 'accept-MG677R' })
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
  const prLane = (w: World) => ran(w, 'bash', SPAWN).filter(r => r.includes('pr'))
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
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
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
    await press($)
    // Which ticket runs next is the teamlead's call, never the button's (X6JHQX).
    expect(ran(w, 'jaira', 'next')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted: done')
  })

  test('nothing ahead of the target: no pull request is asked for, no tab opens, and the press says accepted', async ($, on) => {
    const w = world(on, { herdr: true })
    w.next = ['AB12CD']
    await press($)
    expect(ran(w, 'gh', 'pr')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted: done')
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
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-zz99zz', '@/.worktrees/repo-zz99zz']])
  })

  test('a branch left behind beside the Herdr review branch of it: the pull request opens from the ticket branch itself', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [{ ...SIGNOFF, handle: 'ZZ99ZZ' }]
    w.branches = ['review/task/zz99zz', 'task/zz99zz']
    w.subjects = ['fix(ZZ99ZZ): one']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-ZZ99ZZ' })
    expect(ran(w, 'git', 'worktree', 'add')).toEqual([['git', 'worktree', 'add', '/.worktrees/repo-zz99zz', 'task/zz99zz', '@/repo']])
  })

  test('two tickets on one branch: no pull request while the other is not in done, nothing starts, and the press says which one is missing', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 4
    w.subjects = ['feat(MG677R): cache key', 'fix(T8QS1Y): filter', 'feat(T8QS1Y): filter', 'docs: no ticket']
    w.next = ['AB12CD']
    await press($)
    expect(ran(w, 'gh', 'pr')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; no pull request for feat/mg677r yet: T8QS1Y not in done')
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
    // A wip holding the same commit has fewer parts and is no reason to leave
    // it; a main checkout detached mid-rebase lists first and is no branch.
    w.holders = ['(HEAD detached at abc123)', 'master', 'feat/mg677r', 'wip']
    await startSession($)
    const ui = await $.ui.mount(BAND)
    await ui.press({ key: 'accept-ZZ99ZZ' })
    expect(ran(w, 'git', 'log', '--branches')).toEqual([['git', 'log', '--branches', '--not', 'upstream/main', '--format=%H %s', '@/repo']])
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([['git', 'branch', '--contains', 'abc123', '--format=%(refname:short)', '@/repo']])
    // Its worktree is the other ticket's: the role pushes from there.
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
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
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'zz99zz', 'ZZ99ZZ', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
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
    await ui.press({ key: 'accept-ZZ99ZZ' })
    expect(ran(w, 'git', 'branch', '--contains')).toEqual([['git', 'branch', '--contains', 'abc123', '--format=%(refname:short)', '@/repo']])
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

  test('a fork without its parent as a remote: nothing is counted, nothing starts, and the press says which remote to add', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.remotes = 'origin\tgit@github.com:acme/repo.git (fetch)\n'
    await press($)
    expect(ran(w, 'git', 'log')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toContain('Pull requests go to upstream/repo, and no remote here points at it')
  })

  test('no gh: the pull request goes unchecked, nothing starts, and the press says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prs = null
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(ran(w, 'jaira', 'move', 'MG677R').map(r => r[5])).toEqual(['--dry-run', '--force'])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; pull requests not checked (gh repo view: gh: command not found)')
  })

  test('a forge that answers the repository but not the pull requests: nothing starts, and the press says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.prListFails = true
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toContain('pull requests not checked (gh: HTTP 401')
  })

  test('a subject naming several tickets holds the pull request for each of them', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.tickets = [SIGNOFF, WORKING, { ...WORKING, handle: 'AB12CD', status: 'done' }]
    w.subjects = ['feat(MG677R): cache key', 'chore(AB12CD,T8QS1Y): both at once']
    await press($)
    expect(prLane(w)).toEqual([])
    expect(w.toasts.at(-1)).toContain('no pull request for feat/mg677r yet: T8QS1Y not in done')
  })

  test('no remote: nothing starts, and the pane says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('.jaira/pr says none and the host names no forge: the press lands on origin through the pr lane, asks no pull request, and says landing', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = 'git@git.example.de:acme/repo.git'
    w.basePr = '# solo\nnone\n'
    await press($)
    expect(ran(w, 'gh')).toEqual([])
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
    expect(w.toasts.at(-1)).toContain('MG677R: landing feat/mg677r on the default branch')
    expect(w.toasts.join('\n')).not.toContain('pull request for')
  })

  test('.jaira/pr says none and there is no origin: nothing lands, and the pane says there is nothing to land on', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    w.basePr = 'none\n'
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('no remote origin to land on')
  })

  test('.jaira/pr says none and the forge does not answer: nothing lands, since whether origin is a fork cannot be told, and the pane says why', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.prs = null
    w.basePr = '\uFEFFnone\r\n'
    await press($)
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
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted; feat/mg677r does not land yet: T8QS1Y not in done')
  })

  test('.jaira/pr says none only in this checkout, not in the base: the base decides, so it is the pull request path', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    w.files.set('/repo/.jaira/pr', 'none\n')
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('.jaira/pr says none and no forge names the base: it is read from origin\'s default branch, as the gate reads it, not the local one', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 2
    w.origin = 'git@git.example.de:acme/repo.git'
    w.originHead = 'origin/main'
    w.basePr = 'none\n'
    await press($)
    expect(ran(w, 'git', 'show')).toEqual([['git', 'show', 'origin/main:.jaira/pr', '@/repo']])
    expect(ran(w, 'git', 'log', '--format=%s')).toEqual([['git', 'log', '--format=%s', 'origin/main..feat/mg677r', '@/repo']])
    expect(prLane(w)).toHaveLength(1)
  })

  test('.jaira/pr with a typo is the pull request path: no forge still says so', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = null
    w.basePr = 'nnone\n'
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs an origin on GitHub or GitLab')
  })

  test('a host neither GitHub nor GitLab: nothing starts, and the pane names the setting that settles it', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.origin = 'git@git.example.de:acme/repo.git'
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('git config jaira.forge')
  })

  test('no pull request due and no Herdr: nothing is missing, so nothing fails, and the press says accepted', async ($, on) => {
    const w = world(on)
    w.next = ['AB12CD']
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(w.opened).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R accepted: done')
  })

  test('Herdr App (HC_TERMINAL_ID, no HERDR_ENV): the press opens the pull request through spawn.sh all the same', async ($, on) => {
    const w = world(on, { herdr: true, env: { HOME: '/home/a', HC_TERMINAL_ID: 't3' } })
    w.ahead = 3
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([['bash', SPAWN, '--no-worktree', 'mg677r', 'MG677R', 'pr', '/.worktrees/repo-mg677r', '@/.worktrees/repo-mg677r']])
  })

  test('no Herdr: nothing starts, and the pane says the step needs it', async ($, on) => {
    const w = world(on)
    w.ahead = 3
    await press($)
    expect(ran(w, 'bash', SPAWN)).toEqual([])
    expect(await failure($)).toContain('needs a pull request')
  })

  test('a refused Accept starts no next step at all', async ($, on) => {
    const w = world(on, { herdr: true })
    w.ahead = 3
    w.refusals = [...w.refusals, { code: 'needs_commits', message: 'no commit names MG677R' }]
    await press($)
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
    expect(ran(w, 'bash', SPAWN)).toEqual([])
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
    const spawn = w.ran.findIndex(r => r[1] === SPAWN)
    expect(w.ran[close]?.[3]).toBe('w1:t9')
    expect(close).toBeLessThan(spawn)
    expect(w.ran[spawn]?.slice(0, 7)).toEqual(['bash', SPAWN, '--no-parent', 'mg677r', 'MG677R', 'worker', '/repo'])
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
    expect(ran(w, 'bash', SPAWN)[0]?.slice(0, 7)).toEqual(['bash', SPAWN, '--no-parent', 'mg677r', 'MG677R', 'worker', '/repo'])
    expect(ran(w, 'jaira', 'note', 'MG677R').at(-1)?.[3]).toBe('cockpit: given back to a new worker - tab alive, cache warm for 40 more min, context 80k, but its pane took no text')
  })
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
    await ui.press({ key: 'back-MG677R' })
    return $.ui.mount(GIVE('Give back MG677R'))
  }
  type Mounted = Awaited<ReturnType<typeof openPane>>
  async function rows(pane: Mounted) {
    const out: string[] = []
    for (let i = 0; ; i++) {
      const row = await pane.find({ key: `gap-MG677R-${i}` })
      if (row === undefined) return out
      out.push(row.text)
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
    expect((await pane.find({ key: 'gap-MG677R-0' }))?.props.autoFocus).toBe(true)
    expect((await pane.find({ key: 'reason-MG677R' }))?.props.autoFocus).toBeUndefined()
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
      await pane.press({ key: 'gap-MG677R-0' })
      }
    expect(seen).toEqual(['[accept]', '[fix now]', '[follow-up]', '[accept]'])
  })

  test('Give back writes exactly the notes to fix, the follow-ups with their tickets and the accepted ones, after the move', async ($, on) => {
    const w = withGaps(on, GAPS)
    await startSession($)
    const pane = await openPane($)
    // The clock note: accept -> fix now -> follow-up.
    await pane.press({ key: 'gap-MG677R-2' })
    await pane.press({ key: 'gap-MG677R-2' })
    await pane.input({ key: 'reason-MG677R', text: '' })
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
    expect(w.closed).toContain('cockpit-give-back')
  })

  test('what the person types goes along under the notes, on one line', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.input({ key: 'reason-MG677R', text: 'and the docs' })
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size\ndazu: and the docs')
  })

  test('nothing marked fix now and nothing typed gives nothing back and keeps the pane open', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await pane.input({ key: 'reason-MG677R', text: ' ' })
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(w.closed).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R: no note is marked fix now and nothing is typed - mark one, type what is missing, or press Accept')
  })

  test('a follow-up that could not be created is named so, and the note is still written', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    w.createFails = true
    await startSession($)
    const pane = await openPane($)
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.press({ key: 'gap-MG677R-1' })
    await pane.press({ key: 'gap-MG677R-1' })
    await pane.input({ key: 'reason-MG677R', text: '' })
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size\nFolge-Ticket:\n- (not created: jaira: the board is locked): note: pre-existing: the clock')
    // Heard, not only noted: the toast after it says the ticket went back.
    expect(w.opened).toContainEqual({ id: 'cockpit-message', title: 'Follow-up tickets of MG677R failed', focus: true })
  })

  test('every value of the follow-up goes with =, so a finding starting with - is no flag', async ($, on) => {
    const w = withGaps(on, '-race fails in core/gitref')
    await startSession($)
    const pane = await openPane($)
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.input({ key: 'reason-MG677R', text: 'and the docs' })
    // jaira create takes the title only as its argument: after '--', where
    // jaira itself (checked by hand: exit 0) reads '-race' as a title.
    expect(ran(w, 'jaira', 'create')[0]?.slice(2, 7)).toEqual(['--context=-race fails in core/gitref', '--follows=MG677R', '--json', '--', '-race fails in core/gitref'])
  })

  test('a note kept a follow-up at the next acceptance keeps its ticket: no second one is created', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    w.bodies.MG677R = '## Progress\n- **2026-11-20 10:00 · Alex** — Given back from signoff: beheben:\n- note: the size\nFolge-Ticket:\n- A1B2C3: note: pre-existing: the clock\n'
    await startSession($)
    await (await openPane($)).input({ key: 'reason-MG677R', text: '' })
    expect(ran(w, 'jaira', 'create')).toEqual([])
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\n- note: the size\nFolge-Ticket:\n- A1B2C3: note: pre-existing: the clock')
  })

  test('words typed with no note to fix still open the beheben: list the worker reads', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await pane.input({ key: 'reason-MG677R', text: 'the docs are missing' })
    expect(notes(w)[0]).toBe('Given back from signoff: beheben:\ndazu: the docs are missing\nangenommen:\n- note: the size')
  })

  test('the dazu: line does not end the decisions read back after it', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    w.bodies.MG677R = '## Progress\n- **2026-11-20 10:00 · Alex** — Given back from signoff: beheben:\ndazu: the docs\nangenommen:\n- note: the size\n'
    await startSession($)
    expect(await rows(await openPane($))).toEqual(['[accept] 🟡 note: the size (decided before)'])
  })

  test('Accept with notes still marked fix now is refused and says so', async ($, on) => {
    const w = withGaps(on, 'note: the size')
    await startSession($)
    const pane = await openPane($)
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.press({ key: 'gaps-accept-MG677R' })
    expect(ran(w, 'jaira', 'move')).toEqual([])
    expect(w.toasts.at(-1)).toBe('MG677R: 1 note is marked fix now - press Give back, or mark it follow-up or accept')
  })

  test('a note that is no blocker, not marked fix now, does not hold Accept: the decisions go onto the ticket once it is accepted', async ($, on) => {
    const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
    await startSession($)
    const pane = await openPane($)
    // the clock: accept -> fix now -> follow-up
    await pane.press({ key: 'gap-MG677R-1' })
    await pane.press({ key: 'gap-MG677R-1' })
    await pane.press({ key: 'gaps-accept-MG677R' })
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
    await (await openPane($)).press({ key: 'gaps-accept-MG677R' })
    expect(ran(w, 'jaira', 'note')).toEqual([])
    expect(ran(w, 'jaira', 'create')).toEqual([])
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toEqual([])
    expect(w.submitted).toEqual([])
  })

  test('a blocker not marked fix now is accepted only with a reason, which goes onto the ticket as a decision', async ($, on) => {
    const w = withGaps(on, 'blocker: the header\nnote: the size')
    await startSession($)
    const pane = await openPane($)
    // Fix now: no field for a reason, the plain Accept.
    expect(await pane.find({ key: 'why-MG677R' })).toBeUndefined()
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.press({ key: 'gap-MG677R-0' })
    // Not fix now: the Accept button gives way to the required field.
    expect(await pane.find({ key: 'gaps-accept-MG677R' })).toBeUndefined()
    for (const empty of ['', '   ', '\n']) {
      await pane.input({ key: 'why-MG677R', text: empty })
      expect(ran(w, 'jaira', 'move')).toEqual([])
      expect(w.closed).toEqual([])
      expect(w.toasts.at(-1)).toBe('MG677R: a blocker is not marked fix now - type why it is accepted anyway, or mark it fix now and press Give back')
    }
    await pane.input({ key: 'why-MG677R', text: 'the header is\nharmless for now' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(notes(w)).toEqual([
      'Accepted from signoff: angenommen:\n- blocker: the header\n- note: the size',
      'decision: accepted with a blocker not fixed: the header is harmless for now',
    ])
  })

  test('the band\'s Accept accepts at once without a blocker', async ($, on) => {
    const w = withGaps(on, 'note: the header')
    w.bodies.MG677R = '## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n[small edge] the header\n'
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'accept-MG677R' })
    expect(ran(w, 'jaira', 'move', 'MG677R', '--to', 'done', '--force')).toHaveLength(1)
    expect(w.opened.map(o => o.id)).not.toContain('cockpit-give-back')
  })

  // A blocker by its prefix, in markdown, or by the Explain note's class;
  // a 'blocker:' line the Explain note calls a small edge stays a blocker.
  for (const [gaps, explain] of [
    ['blocker: the header', ''],
    ['blocker: the header', '[small edge] the header'],
    ['- **Blocker:** the header\nnote: the size', ''],
    ['note: the header\nnote: the size', '[blocker] the header\n[small edge] the size'],
  ] as const) {
    test(`the band's Accept with a blocker in the review opens the notes and moves nothing: ${JSON.stringify(gaps)}${explain === '' ? '' : ` (Explain: ${JSON.stringify(explain)})`}`, async ($, on) => {
      const w = withGaps(on, gaps)
      if (explain !== '') w.bodies.MG677R = `## Progress\n- **2026-11-21 10:00 · Alex** — explain: Status: DoD 3/3 ok\n${explain}\n`
      await startSession($)
      await (await $.ui.mount(BAND)).press({ key: 'accept-MG677R' })
      expect(ran(w, 'jaira', 'move')).toEqual([])
      expect(w.opened.at(-1)).toEqual({ id: 'cockpit-give-back', title: 'Give back MG677R', focus: true })
    })
  }

  for (const press of ['back', 'accept', 'band-accept'] as const) {
    test(`a move another session made first: the second press writes nothing (${press})`, async ($, on) => {
      const w = withGaps(on, 'note: the size\nnote: pre-existing: the clock')
      w.movedFirst = true
      await startSession($)
      if (press === 'band-accept') await (await $.ui.mount(BAND)).press({ key: 'accept-MG677R' })
      else {
        const pane = await openPane($)
        await pane.press({ key: 'gap-MG677R-0' })
        await pane.press({ key: 'gap-MG677R-1' })
        await pane.press({ key: 'gap-MG677R-1' })
        if (press === 'back') await pane.input({ key: 'reason-MG677R', text: 'and the docs' })
        else {
          await pane.press({ key: 'gap-MG677R-0' })
          await pane.press({ key: 'gap-MG677R-0' })
          await pane.press({ key: 'gaps-accept-MG677R' })
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
    expect((await pane.find({ key: 'reason-MG677R' }))?.props.autoFocus).toBe(true)
    expect(await pane.find({ key: 'gaps-accept-MG677R' })).toBeUndefined()
  })

  test('a note wrapped over lines is one row, a blank line and a none dropped, CRLF and emphasis read', async ($, on) => {
    withGaps(on, '- note: a long finding\n  that wraps here\n- **blocker:** the header\r\n\nnone')
    await startSession($)
    expect(await rows(await openPane($))).toEqual(['[accept] 🟡 note: a long finding that wraps here', '[fix now] 🔴 **blocker:** the header'])
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
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.input({ key: 'reason-MG677R', text: '' })
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
    await pane.press({ key: 'gap-MG677R-0' })
    await pane.press({ key: 'gap-MG677R-1' })
    const typed = pane.input({ key: 'reason-MG677R', text: '' })
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
    await (await $.ui.mount(BAND)).press({ key: 'accept-MG677R' })
    expect(w.submitted).toEqual(['jaira-cockpit: the person accepted MG677R in the cockpit (signoff -> done); MG677R: opening the pull request for feat/mg677r in a Herdr tab'])
  })

  test('Accept whose next step fails: still one line, since the ticket is in done', async ($, on) => {
    const w = world(on)
    w.ahead = 3
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'accept-MG677R' })
    expect(w.submitted).toEqual(['jaira-cockpit: the person accepted MG677R in the cockpit (signoff -> done); the next step after MG677R failed, see the cockpit'])
  })

  test('Give back: one line naming who carries the work now, the reason on one line', async ($, on) => {
    const w = world(on, { herdr: true, store: { 'session:old': {
      session: 'old', pane: 'w1:p9', cwd: '/x', repo: '/repo', ticket: 'MG677R', role: 'worker', lane: '',
      lastAnswerAt: 1_000_000 - 20 * 60_000, tokens: 300_000, window: 500_000, ttl: 3_600_000, updatedAt: 1_000_000 - 20 * 60_000,
    } } })
    w.tabs = [{ tab_id: 'w1:t9', label: 'MG677R/worker' }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'back-MG677R' })
    await (await $.ui.mount(GIVE('Give back MG677R'))).input({ key: 'reason-MG677R', text: 'der Branch\nfehlt   im Schluessel ' + 'x'.repeat(200) })
    expect(w.submitted).toHaveLength(1)
    const line = w.submitted[0]!
    expect(line.startsWith('jaira-cockpit: the person gave MG677R back in the cockpit (signoff -> in-progress), a new worker started: der Branch fehlt im Schluessel xx')).toBe(true)
    expect(line).not.toContain('\n')
    expect(line.endsWith('…')).toBe(true)
  })

  test('a reason cut at its limit never splits a character outside the BMP', async ($, on) => {
    const w = world(on)
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'back-MG677R' })
    // 118 characters, then emoji: the cut at 119 falls inside a surrogate pair by UTF-16 index.
    await (await $.ui.mount(GIVE('Give back MG677R'))).input({ key: 'reason-MG677R', text: 'x'.repeat(118) + '\u{1F44D}\u{1F44D}\u{1F44D} passt nicht' })
    const line = w.submitted[0]!
    expect(line.endsWith('x\u{1F44D}…')).toBe(true)
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
    await ui.press({ key: 'back-MG677R' })
    await (await $.ui.mount(GIVE('Give back MG677R'))).input({ key: 'reason-MG677R', text: 'fehlt' })
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
    await pane.press({ key: 'overview-lane-signoff' })
    expect(await handles()).toEqual([])
    expect(await pane.find({ type: 'Text', text: 'none' })).toBeDefined()
    await pane.press({ key: 'overview-lanes-all' })
    expect(await handles()).toEqual(['N1N1N1', 'S1S1S1'])
    expect((await pane.find({ key: 'overview-lanes-all' }))?.props.variant).toBe('primary')
  })

  test('a lane or person chosen with no event left keeps its button, and pressing it clears the filter', async ($, on) => {
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
    await pane.press({ key: 'overview-lane-signoff' })
    await pane.press({ key: 'overview-person-Alex' })
    expect(await handles()).toEqual(['S1S1S1'])
    // Unchosen and without an event, the buttons go again.
    expect(await pane.find({ key: 'overview-lane-signoff' })).toBeUndefined()
    expect(await pane.find({ key: 'overview-person-Alex' })).toBeUndefined()
    expect((await pane.find({ key: 'overview-lanes-all' }))?.props.variant).toBe('primary')
  })

  test('Close stands in the tab row, framed and wider than a tab, on every tab and on an unread board', async ($, on) => {
    const w = world(on)
    w.tickets = [{ ...WORKING, assignee: 'Alex', updated_at: ago(0.1) }]
    await startSession($)
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    const pane = await $.ui.mount(OVERVIEW)
    for (const tab of ['queue', 'log', 'workers']) {
      await pane.press({ key: `overview-tab-${tab}` })
      const frame = await pane.find({ key: 'overview-close-frame' })
      expect(frame?.props).toMatchObject({ borderStyle: 'round', paddingX: 1 })
      const labels = (await pane.findAll({ type: 'Button' } as never)).map(b => (b.props as { label: string }).label)
      expect(labels.slice(0, 4)).toEqual(['Queue', 'Log', 'Workers', 'Close'])
    }
    // What Close draws is longer than the longest tab's label.
    const closeText = (await pane.find({ type: 'Text', text: /^ +Close +$/ }))?.text ?? ''
    expect(closeText.length).toBeGreaterThan('Workers'.length)
    await pane.press({ key: 'overview-close' })
    await pane.unmount()
    w.nextFails = true
    await (await $.ui.mount(BAND)).press({ key: 'overview' })
    expect(await (await $.ui.mount(OVERVIEW)).find({ key: 'overview-close-frame' })).toBeDefined()
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
      // 40 rows less the tab row, three high with the framed Close in it.
      expect(picture.length + lines.length).toBeLessThanOrEqual(37)
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

// The Verdict pane's light, worked out from the ticket alone: one rule set per
// colour, the reason it names first.
describe('the verdict light', () => {
  const ready: Waiting = {
    id: 'MG677R', title: 'cache key', goal: '', lane: 'signoff', verdict: 'approve', question: '', dodComplete: true,
    dodOpen: { person: 0, other: 0 }, acceptTo: 'done', assignee: '', explainable: false, blockers: [], notes: [], findingsKnown: true, gaps: '',
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

  function withNotes(on: On, language?: string, answer: string[] = ru) {
    const w = world(on)
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps: NOTES } }, HUMAN, WORKING]
    w.bodies.MG677R = CLASSES
    if (language !== undefined) w.files.set(SETTINGS, JSON.stringify({ dispatch: 'worker', language }))
    w.model = JSON.stringify(answer)
    return w
  }
  async function verdictNotes($: Test$) {
    const { texts } = await verdictOf($, await $.ui.mount(BAND), 'MG677R')
    return texts.filter(x => /^\p{Extended_Pictographic}/u.test(x))
  }
  async function giveBackRows($: Test$) {
    await (await $.ui.mount(BAND)).press({ key: 'back-MG677R' })
    const pane = await $.ui.mount(GIVE('Give back MG677R'))
    const out: string[] = []
    for (let i = 0; ; i++) {
      const row = await pane.find({ key: `gap-MG677R-${i}` })
      if (row === undefined) return { out, pane }
      out.push(row.text)
    }
  }

  test('each class has its mark, in the Verdict pane and in the Give back pane; without a language nothing is sent', async ($, on) => {
    const w = withNotes(on)
    await startSession($)
    expect(await verdictNotes($)).toEqual([
      '🔴 blocker: the header is written twice',
      '🟡 note: the size check is missing',
      '⚪ note: windows paths over 260 characters',
      '🔵 note: pre-existing: load.go:41 accepts an empty file',
    ])
    expect((await giveBackRows($)).out).toEqual([
      '[fix now] 🔴 blocker: the header is written twice',
      '[accept] 🟡 note: the size check is missing',
      '[accept] ⚪ note: windows paths over 260 characters',
      '[accept] 🔵 note: pre-existing: load.go:41 accepts an empty file',
    ])
    expect(w.modelCalls).toBe(0)
  })

  // Two people, two machines: one test each, as each has its own settings.
  for (const [language, answer] of [['Русский', ru], ['Deutsch', de]] as const) {
    test(`a person who set ${language} reads the notes in it, the ticket keeps its words`, async ($, on) => {
      const w = withNotes(on, language, answer)
      await startSession($)
      expect(await verdictNotes($)).toEqual([`🔴 ${answer[0]}`, `🟡 ${answer[1]}`, `⚪ ${answer[2]}`, `🔵 ${answer[3]}`])
      expect(w.prompts.at(-1)).toContain(`in ${language}`)
      expect(w.prompts.at(-1)).toContain('blocker: the header is written twice')
      // Nothing of it reaches the board: no jaira command writes a word of it.
      expect(w.ran.filter(r => r[0] === 'jaira' && ['note', 'set', 'move'].includes(r[1]!))).toEqual([])
      expect(JSON.stringify(w.ran)).not.toContain(answer[0])
    })
  }

  test('Give back shows the notes in the language and writes the originals onto the ticket', async ($, on) => {
    const w = withNotes(on, 'Русский')
    await startSession($)
    const { out, pane } = await giveBackRows($)
    expect(out).toEqual([`[fix now] 🔴 ${ru[0]}`, `[accept] 🟡 ${ru[1]}`, `[accept] ⚪ ${ru[2]}`, `[accept] 🔵 ${ru[3]}`])
    await pane.input({ key: 'reason-MG677R', text: '' })
    const written = ran(w, 'jaira', 'note', 'MG677R').map(r => r[3]).join('\n')
    expect(written).toContain('- blocker: the header is written twice')
    expect(written).not.toContain(ru[0])
  })

  test('kept per ticket and notes: a second look asks nothing, other notes or another language ask again', async ($, on) => {
    const w = withNotes(on, 'Русский')
    await startSession($)
    await verdictNotes($)
    await verdictNotes($)
    await giveBackRows($)
    expect(w.modelCalls).toBe(1)
    w.files.set(SETTINGS, JSON.stringify({ language: 'Deutsch' }))
    w.model = JSON.stringify(de)
    expect((await verdictNotes($))[0]).toBe(`🔴 ${de[0]}`)
    expect(w.modelCalls).toBe(2)
    w.tickets = [{ ...SIGNOFF, review: { verdict: SIGNOFF.review.verdict, gaps: 'note: only one now' } }, HUMAN, WORKING]
    w.model = JSON.stringify(['Nur noch eine.'])
    await w.clock.advance(60_000)
    expect(await verdictNotes($)).toEqual(['🟡 Nur noch eine.'])
    expect(w.modelCalls).toBe(3)
  })

  for (const answer of ['not json at all', JSON.stringify(ru.slice(0, 3)), JSON.stringify(['', ...ru.slice(1)]), JSON.stringify([1, 2, 3, 4]), '{"a": 1}']) {
    test(`an answer that does not fit leaves the original text with its mark: ${answer.slice(0, 30)}`, async ($, on) => {
      const w = withNotes(on, 'Русский')
      w.model = answer
      await startSession($)
      expect((await verdictNotes($))[0]).toBe('🔴 blocker: the header is written twice')
      expect(w.modelCalls).toBe(1)
    })
  }

  test('a model that does not answer leaves the original text with its mark', async ($, on) => {
    const w = withNotes(on, 'Русский')
    w.modelDown = true
    await startSession($)
    expect((await verdictNotes($))[1]).toBe('🟡 note: the size check is missing')
    expect(w.modelCalls).toBe(1)
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

  test('the prompt names the language and carries every note word for word', () => {
    const p = rephrasePrompt('Русский', ['blocker: `x` in a.go:3', 'note: "quoted"'])
    expect(p).toContain('in Русский')
    expect(p).toContain('exactly 2 strings')
    expect(p).toContain(JSON.stringify(['blocker: `x` in a.go:3', 'note: "quoted"']))
  })
})
