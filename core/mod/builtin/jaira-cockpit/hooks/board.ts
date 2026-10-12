// What the cockpit knows about a jaira board it reads off the board itself,
// through the jaira binary: `jaira lanes --json` says whose lane is whose,
// `jaira list --json` where the tickets stand. Not from the agents' reports:
// so it is the same in both dispatch modes and with Herdr tabs. register.tsx
// runs the commands; this file only reads what they printed.

import { carrierTab, type Tab } from './sessions.ts'

export type Waiting = {
  id: string
  title: string
  lane: string
  verdict: string
  question: string
  // Where Accept moves it: the ticket's next lane, when that lane is the
  // board's terminal one (signoff -> done); null where the next lane is
  // still work to do (human -> review), so a press never skips a review.
  acceptTo: string | null
  // Whose it is; '' for nobody's.
  assignee: string
}

export type Counts = { worker: number; review: number; person: number }

type Lane = { id: string; agentic: boolean; terminal?: boolean; context: string; human_exit?: boolean }
type Ticket = {
  handle: string
  title: string
  status: string
  next_lane?: string
  question?: string
  assignee?: string
  review?: { verdict?: string }
}

function parse<T>(json: string): T | undefined {
  try {
    return JSON.parse(json) as T
  } catch {
    return undefined
  }
}

const oneLine = (s: string | undefined) => (s ?? '').replace(/\s+/g, ' ').trim()

// A jaira built before `lanes --json` carried human_exit says nothing about
// whose lane is whose: then the shipped board's person lanes are taken.
const PERSON_LANES = new Set(['human', 'signoff'])

function isPersons(lane: Lane): boolean {
  return lane.human_exit ?? PERSON_LANES.has(lane.id)
}

export type Agent = { description: string; name?: string; status: string }

// The chats a person opens from the band are theirs, not an agent at work.
const PERSONS_TABS = new Set(['explain', 'discuss'])
const WORKING = new Set(['pending', 'running', 'waiting'])

// The tickets waiting in a lane a person owns: those of the person this
// session is (`me`, their name and aliases from `jaira whoami`) or of nobody,
// counted as waiting for them, and the teammates' apart. Without `me` all
// count as the person's: rather a button too many than a ticket hidden. And
// the agents at work on this board's tickets, by kind: a Herdr tab spawn.sh
// labelled `<id>/<lane>`, and a subagent of this session whose description
// or name carries a ticket's id. One whose lane is the fresh one is a review,
// every other a worker. A tab or an agent naming no ticket of this board is
// not counted.
export function snapshotOf(
  lanesJson: string,
  listJson: string,
  tabs: readonly Tab[],
  agents: readonly Agent[],
  me?: readonly string[],
): { waiting: Waiting[]; theirs: Waiting[]; counts: Counts } {
  const mine = new Set((me ?? []).map(n => n.toLowerCase()))
  const lanes = parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? []
  const tickets = parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? []
  const byId = new Map(lanes.map(l => [l.id, l]))
  const waiting: Waiting[] = []
  const theirs: Waiting[] = []
  const counts: Counts = { worker: 0, review: 0, person: 0 }
  for (const t of tickets) {
    const lane = byId.get(t.status)
    if (lane === undefined || !isPersons(lane)) continue
    const assignee = (t.assignee ?? '').trim()
    const isMine = me === undefined || assignee === '' || mine.has(assignee.toLowerCase())
    if (isMine) counts.person++
    const next = t.next_lane === undefined ? undefined : byId.get(t.next_lane)
    const into = isMine ? waiting : theirs
    into.push({
      id: t.handle,
      title: t.title,
      lane: t.status,
      verdict: oneLine(t.review?.verdict),
      question: oneLine(t.question),
      acceptTo: next?.terminal === true ? next.id : null,
      assignee,
    })
  }
  const ids = tickets.map(t => t.handle)
  const fresh = lanes.filter(l => l.agentic && l.context === 'fresh').map(l => l.id)
  // A worker's tab stays open after its turn, for Give back: only one whose
  // agent is working counts.
  for (const { label, agent_status } of tabs) {
    const [id = '', lane = ''] = label.split('/')
    if (agent_status !== 'working' || !ids.includes(id) || lane === '' || PERSONS_TABS.has(lane)) continue
    if (fresh.includes(lane)) counts.review++
    else counts.worker++
  }
  for (const a of agents) {
    const said = `${a.description} ${a.name ?? ''}`
    if (!WORKING.has(a.status) || !ids.some(id => said.includes(id))) continue
    if (fresh.some(lane => new RegExp(`\\b${lane}\\b`, 'i').test(said))) counts.review++
    else counts.worker++
  }
  return { waiting, theirs, counts }
}

// Each ticket's lane, by handle, from `jaira list --json`.
export function lanesOf(listJson: string): Map<string, string> {
  return new Map((parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? []).map(t => [t.handle, t.status]))
}

// Who `jaira whoami --json` says the person is: the identity and its aliases.
export function whoOf(whoamiJson: string): string[] | undefined {
  const who = parse<{ identity?: string; aliases?: string[] }>(withoutWarnings(whoamiJson))
  const names = [who?.identity ?? '', ...(who?.aliases ?? [])].map(n => n.trim()).filter(n => n !== '')
  return names.length === 0 ? undefined : names
}

// A review verdict a person should have explained: one that names a
// reservation, or does not accept the work at all. A verdict is a sentence
// ("Approve. The diff meets the DoD ...", "the diff matches the criteria; no
// defects found"); a plain yes, however long, gets no Explain.
const RESERVATION = /vorbehalt|reservation|caveat|changes|änderungen|\bbut\b|\baber\b|however|jedoch|except|ausser|außer|allerdings|offen(e|er)? (punkt|frage)|open (point|question)/i
const REFUSAL = /^\s*(request|reject|refuse|fail|zur(ü|ue)ck|ablehn|nicht|not)\b/i

// A reservation word under a negation is none: "approve, no changes needed",
// "no open questions", "ohne Vorbehalt". The negation reaches two words at
// most and never past a comma or a full stop, and a connective is never
// negated: "no blockers, but the docs are missing" keeps its reservation.
const CONNECTIVE = /\b(but|aber|however|jedoch)\b/i
const NEGATED = new RegExp(`\\b(no|without|kein(e|en|er)?|ohne)\\s+([^\\s,;:.]+\\s+){0,2}?(${RESERVATION.source})`, 'gi')

export function hasReservation(verdict: string): boolean {
  const unnegated = verdict.replace(NEGATED, m => (CONNECTIVE.test(m) ? m : ''))
  return RESERVATION.test(unnegated) || REFUSAL.test(verdict)
}

// A jaira older than the board it reads says so on a line of its own before
// anything else it prints; that line is no part of an answer or an error.
const OLDER_BUILD = /^jaira: this board was set up by an older version.*$\n?/gm

export function withoutWarnings(text: string): string {
  return text.replace(OLDER_BUILD, '').trim()
}

// What stops a move besides the person's own exit, read off the error
// `jaira move --dry-run --json` prints when a gate refuses: every other
// refusal, or undefined when the human exit (needs_human) is the only one.
// That one Accept stands for; every other gate holds, as for a move made by
// hand in the board. An error that lists no refusals is a refusal of its own.
export function otherRefusal(errorJson: string): string | undefined {
  type Refused = { error?: { message?: string; violations?: { code?: string; message?: string }[] } }
  const text = withoutWarnings(errorJson)
  let error: Refused['error']
  try {
    error = (JSON.parse(text.slice(Math.max(0, text.search(/^\{/m)))) as Refused).error
  } catch {
    return text || 'jaira refused the move'
  }
  if (error?.violations === undefined) return error?.message || 'jaira refused the move'
  const others = error.violations.filter(v => v.code !== 'needs_human')
  return others.length === 0 ? undefined : others.map(v => v.message ?? v.code).join('\n')
}

// What a chat the band opened wrote back onto the ticket, the newest of its
// kind, read off `jaira show --json`: an Explain chat's `explain:` note, a
// Discuss chat's `decision:` note, up to the next note or section, and when
// it was written (jaira stamps notes to the minute, in UTC). undefined when
// there is none yet.
export function noteOf(showJson: string, kind: 'explain' | 'decision'): { text: string; at: number } | undefined {
  const lead = new RegExp(`^- \\*\\*(\\d{4}-\\d\\d-\\d\\d \\d\\d:\\d\\d) · [^*]*\\*\\* — ${kind}:\\s*`)
  const body = parse<{ body?: string }>(withoutWarnings(showJson))?.body ?? ''
  const lines = body.split('\n')
  let at = -1
  for (let i = 0; i < lines.length; i++) if (lead.test(lines[i]!)) at = i
  if (at < 0) return undefined
  const stamp = lead.exec(lines[at]!)![1]!
  const rest = lines.slice(at + 1)
  const end = rest.findIndex(l => l.startsWith('- **') || l.startsWith('#'))
  const text = [lines[at]!.replace(lead, ''), ...(end < 0 ? rest : rest.slice(0, end))].join('\n').trim()
  return { text, at: Date.parse(`${stamp.replace(' ', 'T')}:00Z`) }
}

// The main working tree, from `git rev-parse --git-common-dir`: spawn.sh puts
// a ticket's worktree beside it, never beside another worktree.
export function mainRepoOf(gitCommonDir: string, fallback: string): string {
  const dir = gitCommonDir.trim().replace(/\/+$/, '')
  if (!dir.endsWith('/.git')) return fallback
  return dir.slice(0, -'/.git'.length)
}

// Where the ticket's work lives, from `git worktree list --porcelain`: the
// worktree on a branch ending in its lower-case id, else the main working
// tree. Running `jaira move` anywhere else answers "on its ref and not on
// your disk".
export function worktreeFor(porcelain: string, id: string, fallback: string): string {
  return worktreesOf(porcelain).find(w => isTicketsTree(w, id))?.path ?? fallback
}

type Worktree = { path: string; branch: string }

function worktreesOf(porcelain: string): Worktree[] {
  const out: Worktree[] = []
  for (const block of porcelain.split(/\n\n+/)) {
    const path = /^worktree (.+)$/m.exec(block)?.[1]
    const branch = (/^branch (.+)$/m.exec(block)?.[1] ?? '').replace(/^refs\/heads\//, '')
    if (path !== undefined) out.push({ path, branch })
  }
  return out
}

function isTicketsTree(w: Worktree, id: string): boolean {
  const slug = id.toLowerCase()
  return w.branch.endsWith(`/${slug}`) || w.path.endsWith(`-${slug}`)
}

// The branch a ticket's work rides on and the main working tree's, from the
// same porcelain; undefined where there is none. The main tree's branch is
// the base only where no forge can name one: in a fork the local master
// lags the upstream by hundreds of commits, every one of them "ahead".
export function branchesFor(porcelain: string, id: string): { ticket?: string; base?: string } {
  const trees = worktreesOf(porcelain)
  const ticket = trees.find(w => isTicketsTree(w, id))?.branch
  const base = trees[0]?.branch
  return { ...(ticket ? { ticket } : {}), ...(base ? { base } : {}) }
}

// Where a branch is checked out, off the same porcelain; undefined where it
// is in no worktree. Without one, newWorktreeFor is where spawn.sh would have
// put it: beside the repository, never inside it.
export function worktreeOfBranch(porcelain: string, branch: string): string | undefined {
  return worktreesOf(porcelain).find(w => w.branch === branch)?.path
}

export function newWorktreeFor(repo: string, branch: string): string {
  const cut = repo.lastIndexOf('/')
  return `${repo.slice(0, cut)}/.worktrees/${repo.slice(cut + 1)}-${branch.slice(branch.lastIndexOf('/') + 1)}`
}

// What an Accept starts once the ticket is in done — the one rule, decided
// by Alex on GPFA7B. A branch with commits master does not have and no pull
// request yet is finished work nobody can see: the press that accepted it is
// a person's, which is the only thing allowed to open a pull request, so it
// opens one. But a branch can carry several tickets, and a pull request
// opened while one of them is still being worked proposes half the work: so
// only once every ticket its commits name is in done (`pending` empty). Any
// other case — a ticket still pending, the branch already proposed (a pull
// request in any state, merged and closed included), or nothing ahead of
// master — has no further step of its own, and the press moves the board on
// to the next ticket instead. One button, so the rule decides, never the
// person per press. The caller asks the forge only where the answer can
// still change the step, and passes false otherwise. A forge that could not
// be asked — offline, logged out, its repository no remote here — is
// undefined: whether a pull request exists is unknown, so none is opened,
// and the board moves on rather than wait for the network (Alex, GPFA7B
// DoD 8); the press says why.
export type StepAfterAccept = 'pr' | 'next'

export function stepAfterAccept(ahead: number, pending: readonly string[], hasPullRequest: boolean | undefined): StepAfterAccept {
  return ahead > 0 && pending.length === 0 && hasPullRequest === false ? 'pr' : 'next'
}

// The tickets a branch carries, off `git log --format=%s <base>..<branch>`:
// every handle a subject names the way this board's commits do,
// `type(A3K9QP): ...`, or several at once, `chore(A3K9QP,B7XQ2M): ...`,
// once each.
export function ticketsNamed(subjects: string): string[] {
  const out: string[] = []
  for (const line of subjects.split('\n')) {
    const scope = /^[a-z]+\(([^)]*)\)!?:/.exec(line.trim())?.[1] ?? ''
    for (const id of scope.split(/[\s,]+/)) if (/^[0-9A-Z]{6}$/.test(id) && !out.includes(id)) out.push(id)
  }
  return out
}

// Which of those are not in the terminal lane yet, off `jaira list --json`.
// A handle the board does not list is filed away in the logbook — past done
// — or no ticket of this board, and holds nothing up.
export function pendingOf(listJson: string, named: readonly string[], terminal: string): string[] {
  const tickets = parse<{ tickets?: { handle?: string; status?: string }[] }>(withoutWarnings(listJson))?.tickets ?? []
  return named.filter(id => {
    const status = tickets.find(t => t.handle === id)?.status
    return status !== undefined && status !== terminal
  })
}

// A ticket's branch without a worktree of its own, off a branch listing
// (`--format=%(refname:short)`): the first that is not the base. Read from
// `git branch --list '*/<slug>'` — spawn.sh names it <prefix>/<slug>, and a
// worktree removed after done leaves it behind — and from `git branch
// --contains <a commit naming it>`, where it was worked on another ticket's
// branch.
export function carrierBranch(listing: string, base: string): string | undefined {
  return listing.split('\n').map(l => l.trim()).find(b => b && b !== base)
}

// The forge jaira-role-pr would open a pull request on: `git config
// jaira.forge` where it is set, else the host of origin. undefined where
// neither settles it.
export type Forge = 'github' | 'gitlab'

export function forgeOf(forge: string, originUrl: string): Forge | undefined {
  const on = forge.trim() || (/github\.com/i.test(originUrl) ? 'github' : /gitlab/i.test(originUrl) ? 'gitlab' : '')
  return on === 'github' || on === 'gitlab' ? on : undefined
}

// Where a pull request of this branch would be, as jaira-role-pr settles it:
// in a fork, the parent — asking origin there answers for the fork, finds
// nothing, and a second pull request gets opened onto a merged branch. The
// command describes a repository (origin's URL, or owner/repo);
// targetRepoOf reads the target off origin's answer, defaultBranchOf the
// branch pull requests go into off the target's.
export function repoViewCommand(on: Forge, repo: string): string[] {
  return on === 'github'
    ? ['gh', 'repo', 'view', repo, '--json', 'isFork,parent,nameWithOwner,defaultBranchRef']
    : ['glab', 'repo', 'view', repo, '-F', 'json']
}

export function targetRepoOf(on: Forge, viewJson: string): string | undefined {
  if (on === 'github') {
    const v = parse<{ nameWithOwner?: string; parent?: { name?: string; owner?: { login?: string } } | null }>(withoutWarnings(viewJson))
    const parent = v?.parent?.owner?.login && v.parent.name ? `${v.parent.owner.login}/${v.parent.name}` : undefined
    return parent ?? v?.nameWithOwner
  }
  const v = parse<{ path_with_namespace?: string; forked_from_project?: { path_with_namespace?: string } | null }>(withoutWarnings(viewJson))
  return v?.forked_from_project?.path_with_namespace ?? v?.path_with_namespace
}

export function defaultBranchOf(on: Forge, viewJson: string): string | undefined {
  const v = parse<{ defaultBranchRef?: { name?: string }; default_branch?: string }>(withoutWarnings(viewJson))
  return (on === 'github' ? v?.defaultBranchRef?.name : v?.default_branch) || undefined
}

// The local remote that points at a repository, off `git remote -v`: the
// one whose URL ends in owner/repo, SSH or HTTPS. undefined where none does —
// a fork cloned without its upstream added.
export function remoteFor(remotes: string, repo: string): string | undefined {
  const tail = new RegExp(`[:/]${repo.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}(\\.git)?$`, 'i')
  for (const line of remotes.split('\n')) {
    const [name, url] = line.trim().split(/\s+/)
    if (name && url && tail.test(url)) return name
  }
  return undefined
}

// The pull requests of a branch in that repository, in every state. Both
// tools print a JSON array, empty where there is none.
export function pullRequestsCommand(on: Forge, repo: string, branch: string): string[] {
  return on === 'github'
    ? ['gh', 'pr', 'list', '--repo', repo, '--head', branch, '--state', 'all', '--json', 'number']
    : ['glab', 'mr', 'list', '--repo', repo, '--source-branch', branch, '--all', '-F', 'json']
}

export function hasPullRequest(listJson: string): boolean | undefined {
  const list = parse<unknown>(withoutWarnings(listJson))
  return Array.isArray(list) ? list.length > 0 : undefined
}

// The lane spawn.sh is given for the next ticket: the role a teamlead starts
// under the person's dispatch setting (~/.jaira/settings.json), "worker" only
// where it says so, as settings.DispatchMode reads it.
export function dispatchLane(settingsJson: string | undefined): 'worker' | 'dispatch' {
  const dispatch = parse<{ dispatch?: unknown }>(settingsJson ?? '')?.dispatch
  return typeof dispatch === 'string' && dispatch.trim() === 'worker' ? 'worker' : 'dispatch'
}

// The next ticket to start, off `jaira next --all --json`: the first one no
// worker or dispatcher tab carries already, so an Accept never puts a second
// agent on a ticket that is being worked. undefined when there is none.
export function nextToStart(nextJson: string, tabs: readonly Tab[]): string | undefined {
  const tickets = parse<{ tickets?: { handle?: string }[] }>(withoutWarnings(nextJson))?.tickets ?? []
  return tickets.map(t => t.handle).find((h): h is string => !!h && carrierTab(tabs, h) === undefined)
}

// The directory holding .jaira/, from `cwd` upwards: the candidates in the
// order they are tried.
export function upwards(cwd: string): string[] {
  const dirs: string[] = []
  let dir = cwd.replace(/\/+$/, '') || '/'
  for (let depth = 0; depth < 12; depth++) {
    dirs.push(dir)
    if (dir === '/') break
    const cut = dir.lastIndexOf('/')
    dir = cut <= 0 ? '/' : dir.slice(0, cut)
  }
  return dirs
}
