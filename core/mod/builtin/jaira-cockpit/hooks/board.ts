// What the cockpit knows about a jaira board it reads off the board itself,
// through the jaira binary: `jaira lanes --json` says whose lane is whose,
// `jaira list --json` where the tickets stand. Not from the agents' reports:
// so it is the same in both dispatch modes and with Herdr tabs. register.tsx
// runs the commands; this file only reads what they printed.

export type Waiting = {
  id: string
  title: string
  lane: string
  verdict: string
  question: string
}

export type Counts = { worker: number; review: number; person: number }

type Lane = { id: string; agentic: boolean; context: string; human_exit?: boolean }
type Ticket = {
  handle: string
  title: string
  status: string
  question?: string
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

// The tickets waiting in a lane a person owns, and how many tickets each
// kind of lane holds: an agent's, the fresh review's, a person's.
export function snapshotOf(lanesJson: string, listJson: string): { waiting: Waiting[]; counts: Counts } {
  const lanes = parse<{ lanes: Lane[] }>(lanesJson)?.lanes ?? []
  const tickets = parse<{ tickets: Ticket[] }>(listJson)?.tickets ?? []
  const byId = new Map(lanes.map(l => [l.id, l]))
  const waiting: Waiting[] = []
  const counts: Counts = { worker: 0, review: 0, person: 0 }
  for (const t of tickets) {
    const lane = byId.get(t.status)
    if (lane === undefined) continue
    if (lane.human_exit === true) {
      counts.person++
      waiting.push({ id: t.handle, title: t.title, lane: t.status, verdict: oneLine(t.review?.verdict), question: oneLine(t.question) })
    } else if (lane.agentic && lane.context === 'fresh') counts.review++
    else if (lane.agentic) counts.worker++
  }
  return { waiting, counts }
}

// A review verdict a person should have explained: anything but a plain yes.
export function hasReservation(verdict: string): boolean {
  if (verdict === '') return false
  return !/^(approve|approved|accept|accepted|pass|lgtm|ok)\.?$/i.test(verdict.trim())
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
  const slug = id.toLowerCase()
  for (const block of porcelain.split(/\n\n+/)) {
    const path = /^worktree (.+)$/m.exec(block)?.[1]
    const branch = /^branch (.+)$/m.exec(block)?.[1] ?? ''
    if (path !== undefined && (branch.endsWith(`/${slug}`) || path.endsWith(`-${slug}`))) return path
  }
  return fallback
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
