export const meta = {
  name: 'lane-eval',
  description: 'Run the same jaira tickets in two ways of working, judge the diffs blind, record what every agent did',
  whenToUse: 'Comparing lane-per-agent against one worker per ticket on golden tickets; see docs/LANE-EVAL.md',
  phases: [
    { title: 'Setup', detail: 'one clone without a remote per arm, ticket and run' },
    { title: 'Arms', detail: 'lane-per-agent: a fresh agent per lane; worker: one agent up to review' },
    { title: 'Review', detail: 'a fresh agent works the review lane in every arm' },
    { title: 'Judge', detail: 'a blind judge per ticket and run, twice with X and Y swapped' },
    { title: 'Record', detail: 'results.json, read by scripts/lane-eval/metrics.sh' },
  ],
}

// args: {
//   repo:    absolute path of the repository the tickets and scripts/lane-eval come from
//   out:     absolute, not yet used directory for clones, judge folders, KEY.txt, results.json
//   tickets: ['BG5QJ6', {id: 'GDB5YJ', base: '43080ed'}]  (a golden ticket needs no base)
//   base:    base commit for every ticket that names none (optional)
//   arms:    ['lane-per-agent', 'worker'] (default: both)
//   n:       runs per arm and ticket (default 3)
//   maxSteps: lane agents per run before lane-per-agent gives up (default 20)
// }
const A = args || {}
if (!A.repo || !A.out || !Array.isArray(A.tickets) || A.tickets.length === 0) {
  throw new Error('lane-eval needs args {repo, out, tickets: [...]}; see docs/LANE-EVAL.md')
}
const KNOWN_ARMS = ['lane-per-agent', 'worker']
const ARMS = A.arms || KNOWN_ARMS
for (const arm of ARMS) {
  if (!KNOWN_ARMS.includes(arm)) throw new Error(`unknown arm '${arm}'; known: ${KNOWN_ARMS.join(', ')}`)
}
const N = A.n || 3
const MAX_STEPS = A.maxSteps || 20
const OUT = A.out.replace(/\/+$/, '')
const EVAL_ID = OUT.split('/').pop()
const TOOLS = `${A.repo}/scripts/lane-eval`
const LOOP_LANES = ['critique', 'optimize', 'testing']
// Lanes with no agent step that a fresh ticket is walked through rather than stopped at.
const PASS_THROUGH = ['backlog', 'todo']

const TICKETS = A.tickets.map(t => (typeof t === 'string' ? { id: t, base: A.base || '' } : { id: t.id, base: t.base || A.base || '' }))
const PAIRS = []
TICKETS.forEach((t, ti) => {
  for (let run = 1; run <= N; run++) PAIRS.push({ ticket: t.id, base: t.base, run, ti })
})
log(`${TICKETS.length} ticket(s) x ${ARMS.length} arm(s) x ${N} run(s) = ${PAIRS.length * ARMS.length} clones under ${OUT}`)

// metrics.sh finds an arm's agents by this first prompt line; nothing else carries it.
const mark = (arm, ticket, run, step, lane) =>
  `[lane-eval id=${EVAL_ID} arm=${arm} ticket=${ticket} run=${run} step=${step} lane=${lane}]`

const rules = (dir) => `Work ONLY in ${dir} — an isolated git clone with no remote. \`cd\` there in every shell command. Never push, never touch ${A.repo} or any path outside this clone, never run \`jaira update\`, do not create git worktrees or branches: commit locally on the current branch. Conversational mode is off. A lane that only takes a question (human) is entered only with a genuine open question; without one, follow \`next_lane\` from \`jaira show <id> --json\` past it.
Go tests: set HOME, USERPROFILE and JAIRA_HOME to one ABSOLUTE empty temp directory before \`go test\` — otherwise the tests read this machine's ~/.jaira and about 16 of them fail on the untouched base as well. A relative HOME breaks the Go module cache. Judge a run by its exit code.`

const SETUP = {
  type: 'object',
  properties: {
    dir: { type: 'string' }, base: { type: 'string' }, handle: { type: 'string' },
    lane: { type: 'string' }, next_lane: { type: 'string' },
    agent_lanes: { type: 'array', items: { type: 'string' } },
  },
  required: ['dir', 'base', 'handle', 'lane', 'next_lane', 'agent_lanes'],
}
const LANE = {
  type: 'object',
  properties: {
    status: { type: 'string', description: 'the lane the ticket sits in after your move (status in jaira show --json)' },
    next_lane: { type: 'string', description: 'next_lane in jaira show --json after your move' },
    question: { type: 'string', description: 'the open question, if you moved the ticket to a human lane; else empty' },
  },
  required: ['status', 'next_lane'],
}
const WORKER = {
  type: 'object',
  properties: {
    status: { type: 'string', description: 'the lane the ticket sits in when you stop' },
    loops: { type: 'integer', description: 'how many times critique, optimize or testing sent the ticket back to in-progress' },
    lanes: { type: 'string', description: 'the lanes you worked, in order, e.g. in-progress,critique,in-progress,critique,optimize,testing' },
    question: { type: 'string' },
  },
  required: ['status', 'loops', 'lanes'],
}
const REVIEW = {
  type: 'object',
  properties: {
    status: { type: 'string' },
    review_verdict: { type: 'string', description: 'the review-verdict you wrote, first sentence' },
  },
  required: ['status', 'review_verdict'],
}
const JUDGE = {
  type: 'object',
  properties: {
    better: { type: 'string', enum: ['X', 'Y', 'tie'] },
    x_meets_dod: { type: 'boolean' },
    y_meets_dod: { type: 'boolean' },
    reasons: { type: 'string' },
  },
  required: ['better', 'x_meets_dod', 'y_meets_dod', 'reasons'],
}

async function setup(arm, pair) {
  return agent(
    `Run exactly this one command and nothing else, then return the JSON it prints:\nbash '${TOOLS}/clone.sh' '${A.repo}' '${OUT}' '${arm}' '${pair.ticket}' ${pair.run} '${pair.base}'\nIf it fails, return its error output as your reply instead.`,
    { label: `setup ${arm} ${pair.ticket} r${pair.run}`, phase: 'Setup', schema: SETUP, effort: 'low' })
}

// lane-per-agent: the script decides the lane, a fresh agent works exactly that one.
async function laneByLane(arm, pair, s) {
  const agentic = new Set(s.agent_lanes)
  let status = s.lane, next = s.next_lane, step = 0, loops = 0, question = ''
  const lanes = []
  while (step < MAX_STEPS) {
    if (!agentic.has(status) && !PASS_THROUGH.includes(status)) break
    const lane = agentic.has(status) ? status : next
    if (lane === 'review' || !agentic.has(lane)) break
    step++
    const r = await agent(`${mark(arm, pair.ticket, pair.run, step, lane)}
${rules(s.dir)}

Ticket ${s.handle}, lane \`${lane}\`. If the ticket does not sit in \`${lane}\` yet, \`jaira move ${s.handle} --to ${lane}\` first. Then invoke the Skill \`jaira-role-lane\` with args \`${s.handle} ${lane}\` and follow it for exactly that one lane, including the \`jaira move\` that ends it. Then stop and return the ticket's status and next_lane as \`jaira show ${s.handle} --json\` reports them.`,
      { label: `${arm} ${pair.ticket} r${pair.run} ${lane}`, phase: 'Arms', schema: LANE })
    if (!r) break
    lanes.push(lane)
    if (LOOP_LANES.includes(lane) && r.status === 'in-progress') loops++
    status = r.status; next = r.next_lane; question = r.question || ''
  }
  if (step >= MAX_STEPS) log(`${arm} ${pair.ticket} r${pair.run}: stopped after ${MAX_STEPS} lane agents in ${status}`)
  return { status, loops, lanes: lanes.join(','), question, steps: step }
}

// worker: one agent carries the ticket through every agent lane before review.
async function worker(arm, pair, s) {
  const r = await agent(`${mark(arm, pair.ticket, pair.run, 1, 'in-progress..testing')}
${rules(s.dir)}

You are the single worker for ticket ${s.handle}. Carry it yourself, in this one session, through every agent lane up to but NOT including \`review\`. \`jaira claim ${s.handle}\`, then repeatedly: \`jaira show ${s.handle} --json\` gives \`next_lane\`; for an agent lane run \`jaira show ${s.handle} --for-lane <lane> --json\`, do that lane's work as its prompt says, produce its required outputs, and \`jaira move\` onward. Critique, optimize and testing send work back to in-progress until they have nothing left to say — run those loops honestly, and when you critique your own work, reread the diff as a stranger would. Stop when the ticket sits in \`review\` (or in a human lane, moved there with a genuine question).`,
    { label: `${arm} ${pair.ticket} r${pair.run}`, phase: 'Arms', schema: WORKER })
  if (!r) return { status: 'unknown', loops: 0, lanes: '', question: '', steps: 1 }
  return { status: r.status, loops: r.loops, lanes: r.lanes, question: r.question || '', steps: 1 }
}

async function runUnit(arm, pair) {
  const s = await setup(arm, pair)
  if (!s) return null
  const done = arm === 'worker' ? await worker(arm, pair, s) : await laneByLane(arm, pair, s)
  let review = null
  if (done.status === 'review') {
    review = await agent(`${mark(arm, pair.ticket, pair.run, done.steps + 1, 'review')}
${rules(s.dir)}

Ticket ${s.handle}, lane \`review\`. Invoke the Skill \`jaira-role-lane\` with args \`${s.handle} review\` and follow it for exactly that one lane, including the \`jaira move\` that ends it. Then stop.`,
      { label: `${arm} ${pair.ticket} r${pair.run} review`, phase: 'Review', schema: REVIEW })
  } else {
    log(`${arm} ${pair.ticket} r${pair.run}: no review, the ticket stopped in ${done.status}`)
  }
  return {
    arm, ticket: pair.ticket, run: pair.run, handle: s.handle, dir: s.dir, base: s.base,
    loops: done.loops, lanes: done.lanes, question: done.question,
    ended_in: review ? review.status : done.status,
    review_verdict: review ? review.review_verdict : '',
  }
}

const winnerOf = (j, x, y) => (!j ? null : j.better === 'X' ? x : j.better === 'Y' ? y : 'tie')

async function judgePair(units, pair) {
  const ok = (units || []).filter(Boolean)
  if (ARMS.length !== 2 || ok.length !== 2) return { units: ok, judged: null }
  // Which arm is X in the first pass alternates over tickets and runs, so X is not always the same arm.
  const [ax, ay] = (pair.ti + pair.run) % 2 === 0 ? ARMS : [ARMS[1], ARMS[0]]
  const prep = await agent(
    `Run exactly this one command and nothing else, then return its output verbatim:\nbash '${TOOLS}/judge-prep.sh' '${OUT}' '${pair.ticket}' ${pair.run} '${ax}' '${ay}'`,
    { label: `prep ${pair.ticket} r${pair.run}`, phase: 'Judge', effort: 'low' })
  if (!prep) return { units: ok, judged: null }
  const judge = (p) => agent(`You are a blind judge. Read ONLY the files in ${OUT}/judge/${pair.ticket}-r${pair.run}-${p}/: ticket.md (what was demanded), X.diff and Y.diff (two independent attempts at it) and, if present, reference.diff (the change that was accepted back then — a reference, not a target to copy). Open no other path, run no git, and do not try to find out how either diff was produced.

Decide which diff is the better change to merge for this ticket: does it meet every Definition of Done item, is it correct, does it stay in scope, is it tested, would a reviewer understand it. An empty or broken diff fails. Prefer neither position; answer tie only if you cannot separate them.`,
    { label: `judge ${pair.ticket} r${pair.run} ${p}`, phase: 'Judge', schema: JUDGE })
  const [j1, j2] = await parallel([() => judge('p1'), () => judge('p2')])
  const w1 = winnerOf(j1, ax, ay), w2 = winnerOf(j2, ay, ax)
  let winner
  if (!w1 || !w2) winner = null
  else if (w1 === w2) winner = w1 // an arm, or 'tie'
  else winner = 'split'           // the two orders disagree: position, not quality, decided
  return { units: ok, judged: { winner, p1: { x: ax, y: ay, ...(j1 || {}) }, p2: { x: ay, y: ax, ...(j2 || {}) } } }
}

phase('Setup')
const pairs = await pipeline(PAIRS,
  pair => parallel(ARMS.map(arm => () => runUnit(arm, pair))),
  (units, pair) => judgePair(units, pair))

const runs = []
const judgements = []
pairs.forEach((p, i) => {
  if (!p) return
  const pair = PAIRS[i]
  const w = p.judged ? p.judged.winner : null
  for (const u of p.units) {
    const verdict = w === null ? 'unjudged' : w === 'tie' || w === 'split' ? w : w === u.arm ? 'win' : 'lose'
    runs.push({ ...u, verdict })
  }
  if (p.judged) judgements.push({ ticket: pair.ticket, run: pair.run, ...p.judged })
})
const results = { eval: EVAL_ID, arms: ARMS, n: N, tickets: TICKETS, runs, judgements }

phase('Record')
await agent(`Write the JSON below to ${OUT}/results.json exactly as given — one heredoc with a quoted delimiter, nothing else — and reply "ok".

${JSON.stringify(results, null, 1)}`, { label: 'record results.json', phase: 'Record', effort: 'low' })
log(`next: bash ${TOOLS}/metrics.sh ${OUT} <transcript dir of this run>`)
return results
