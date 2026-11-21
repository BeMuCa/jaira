// Every text the cockpit shows a person, in English whatever language the
// person speaks with Claude: one set of words, the same on every machine.

export const t = {
  waiting: 'Waiting for you',
  accept: 'Accept',
  giveBack: 'Give back',
  explain: 'Explain',
  discuss: 'Discuss',
  more: (n: number) => `+${n} more`,
  fewer: 'show fewer',
  theirs: (n: number, names: string[]) => `${n} waiting for ${names.join(', ')}`,
  reasonLabel: 'What is missing?',
  reasonPlaceholder: 'goes onto the ticket as a note',
  // The Give back pane with the review's notes, one row each.
  gapsHeading: 'Review notes - press a row to change what happens to it:',
  choice: { now: '[fix now]', later: '[follow-up]', accept: '[accept]' } as Record<string, string>,
  decidedBefore: '(decided before)',
  moreLabel: 'Anything else?',
  morePlaceholder: 'optional, goes onto the ticket with the notes',
  acceptWithNow: (id: string, n: number) => `${id}: ${n} note${n === 1 ? ' is' : 's are'} marked fix now - press Give back, or mark ${n === 1 ? 'it' : 'them'} follow-up or accept`,
  blockerNeedsWhy: (id: string) => `${id}: a blocker is not marked fix now - type why it is accepted anyway, or mark it fix now and press Give back`,
  whyLabel: 'Why accept it with the blocker?',
  whyPlaceholder: 'required - goes onto the ticket as a decision',
  acceptedDespite: 'accepted with a blocker not fixed',
  alreadyIn: (id: string, lane: string) => `${id} is already in ${lane}: someone else decided it meanwhile. Nothing was written.`,
  nothingToFix: (id: string) => `${id}: no note is marked fix now and nothing is typed - mark one, type what is missing, or press Accept`,
  followUpsOf: (id: string) => `Follow-up tickets of ${id}`,
  followUpsLost: (id: string, texts: string[]) => `These notes of ${id} were marked follow-up, but jaira created no ticket for them; the note on ${id} names each as not created:\n${texts.map(x => `- ${x}`).join('\n')}`,
  movedMeanwhile: (id: string, lane: string, now: string) => `${id} is no longer in ${lane} but in ${now}: someone else decided it meanwhile. Nothing was written.`,
  accepted: (id: string, lane: string) => `${id} accepted: ${lane}`,
  returnedTo: (id: string, who: 'old' | 'new' | 'none') =>
    who === 'none' ? `${id} back in in-progress; no worker started (the ticket's note says why)` : `${id} given back to ${who === 'old' ? 'the worker that carried it' : 'a new worker'}`,
  failed: (what: string) => `${what} failed`,
  afterAccept: (id: string) => `Next step after ${id}`,
  staysDone: (id: string) => `${id} stays in done; start the next step by hand.`,
  prStarted: (id: string, branch: string) => `${id}: opening the pull request for ${branch} in a Herdr tab`,
  acceptedNoPr: (id: string, hint: string) => `${id} accepted; ${hint}`,
  prWaits: (branch: string, pending: string[]) => `no pull request for ${branch} yet: ${pending.join(', ')} not in done`,
  prUnchecked: (why: string) => `pull requests not checked (${why})`,
  noForge: (branch: string) => `${branch} has commits to propose, but no pull request can be opened from here: this repository needs an origin on GitHub or GitLab, or git config jaira.forge github (or gitlab).`,
  noRemoteFor: (repo: string) => `Pull requests go to ${repo}, and no remote here points at it: add it with git remote add upstream <its URL>.`,
  landStarted: (id: string, branch: string) => `${id}: landing ${branch} on the default branch in a Herdr tab (no pull requests on this board)`,
  landWaits: (branch: string, pending: string[]) => `${branch} does not land yet: ${pending.join(', ')} not in done`,
  cannotLand: (why: string) => `Nothing lands from here: ${why}`,
  noOrigin: 'this repository has no remote origin to land on.',
  landNeedsHerdr: (branch: string) => `${branch} lands on the default branch from a Herdr tab of its own (HERDR_ENV=1 or Herdr App, and spawn.sh from jaira roles install).`,
  prNeedsHerdr: (branch: string) => `${branch} needs a pull request, which opens in a Herdr tab of its own (HERDR_ENV=1 or Herdr App, and spawn.sh from jaira roles install).`,
  noteOn: (id: string) => `Note on ${id}`,
  noteLost: (id: string, why: string, note: string) => `${id} went back to in-progress, but the note with your reason was not written: ${why}\nAdd it by hand: jaira note ${id} '${note.replaceAll("'", `'\\''`)}'`,
  tooOld: (installed: string, needed: string) => `jaira ${installed} is too old for this cockpit: it needs ${needed} or later. Give back will be refused until you update jaira.`,
  released: (latest: string, installed: string) => `jaira ${latest} is released (this is ${installed})`,
  update: 'Update',
  dismiss: 'Dismiss',
  dismissed: (version: string) => `jaira ${version} dismissed: the band offers the next release again`,
  updateJaira: 'Update jaira',
  upToDate: (version: string) => `jaira ${version} is already the latest release`,
  upgradeTitle: (from: string, to: string) => `Update jaira ${from} -> ${to}?`,
  upgradeIntro: (from: string, to: string) =>
    `Nothing has been written yet. Update runs 'jaira self upgrade' (${from} -> ${to}), then 'jaira update --roles --lanes' in this repository. What that would change here, as the jaira installed now sees it (--dry-run):`,
  confirmUpgrade: (to: string) => `Update to ${to}`,
  cancel: 'Cancel',
  upgrading: (to: string) => `updating jaira to ${to}...`,
  upgraded: (to: string) => `jaira updated to ${to}`,
  noRoomUpdate: 'No room for a pane here: widen the terminal and press Update again.',
  dryRunFailed: "Nothing was written and no update is offered: 'jaira update --roles --lanes --dry-run' gave no preview. It said:",
  close: 'Close',
  bandTitle: 'jaira-cockpit',
  noRoom: 'No room for a pane here: widen the terminal and press Give back again.',
  readExplanation: 'Read explanation',
  reviewedIn: (round: number) => `review of round ${round}, changed since`,
  explanationOf: (id: string) => `Explanation of ${id}`,
  decided: 'decided',
  decidedAs: 'Decided:',
  ownAnswer: 'Own answer…',
  recommended: '(recommended)',
  ownLabel: 'Your answer',
  ownPlaceholder: 'goes onto the ticket as a decision',
  decide: 'Decide',
  decidedOn: (id: string) => `${id}: decision noted on the ticket`,
  actionGiveBack: '[give back]',
  actionAccept: '[accept]',
  decisionFor: (id: string, text: string) =>
    `Decision on ${id} in the cockpit's explanation pane, noted on the ticket and not acted on: ${text}`,
  answerOn: (id: string) => `Answer on ${id}`,
  verdict: 'Verdict',
  verdictOf: (id: string) => `Review verdict of ${id}`,
  readDecision: 'Read decision',
  decisionOf: (id: string) => `Decision on ${id}`,
  noSpawn: 'A chat of its own opens in a Herdr tab through spawn.sh, and none was found: run jaira roles install.',
  noBoard: 'No jaira board here.',
  explainAsked: (id: string) => `${id}: a subagent of this session explains the finding; read it under Read explanation`,
  discussHere: (id: string) => `${id}: the open question is talked through in this chat`,
  // Without Herdr the session itself is asked: the same requests spawn.sh
  // types into a tab of its own (lane explain and discuss there) — a change to
  // one belongs in the other. jaira runs with -C: a ticket in a worktree is
  // "on its ref and not on your disk" anywhere else.
  explainBySubagent: (id: string, dir: string) =>
    `The cockpit's Explain on ${id}, without Herdr: start one subagent (Agent tool) with exactly the prompt below and do nothing on ${id} yourself. When it is done, tell me in one line that the explanation of ${id} is on the ticket, under Read explanation in the cockpit.\n\n` +
    `Read ticket ${id} with jaira -C '${dir}' show ${id} --brief (title, goal, question, review, the newest bug hunt and the latest notes). Explain the review on ${id} in plain words: where it stands, what each note in review-gaps means in one everyday sentence, which options there are and which you recommend. Where review.verdict_gaps is not empty and differs from review.gaps, the verdict is from an earlier round and a later lane rewrote review-gaps: the review's notes are those in review.verdict_gaps - read them wherever this says review-gaps, and say that the ticket changed since. ` +
    `Then write the same explanation onto the ticket as one note with exactly these lines: jaira -C '${dir}' note ${id} "explain: Status: DoD <done>/<all> ok|fail | tests ok|fail | blockers ok|fail | bughunt on HEAD ok|fail|unknown|skipped\n[<class>] <one note, one everyday sentence>\n1) <option, one line> [accept]\n2) <option, one line> [give back]\nRecommendation: <number> - <why, one sentence, naming the note it is about>" - with real line breaks in the text. ` +
    `The status line: DoD, tests and blockers each end in ok or fail, bughunt on HEAD in one of ok, fail, unknown or skipped; blockers is ok when review-gaps has no blocker: line; bughunt on HEAD ends in the brief's bughunt.state word as it stands - ok, fail, unknown (a hunted file is changed and not committed, so nobody can tell; never write fail for it) or skipped (the ticket runs no hunt after a fix, by choice) - and in fail when bughunt is null (no hunt finished); only a brief without a bughunt key, from an older jaira, falls back to: ok when the newest 'bughunt: done at <commit>' note names the commit git -C '${dir}' rev-parse --short HEAD prints. ` +
    `One [<class>] line per note in review-gaps, the class exactly one of: [blocker] - a wrong result, data loss, a crash, a security hole; [small edge] - formally fine, one round of about ten minutes closes it; [known limit] - cannot be solved here and is documented; [not this ticket] - belongs elsewhere, say where. Classify first; nothing else per note, no line on why it matters. ` +
    `Two to four options with 1) 2) 3) 4) at the start of the line - the numbers the person picks by in the cockpit - in the order you recommend them, no fixed order. Each option ends in [give back] when it means giving the work back to change something, [accept] when it means accepting the ticket as it is, and no mark when it is neither. When no note is a blocker, the recommendation says plainly that the ticket can be accepted as it is. ` +
    `Change nothing else, commit nothing, do not move the ticket. Answer in the language the person speaks with this session; keep the words Status:, Recommendation:, ok, fail and the class marks in English, the cockpit reads them.`,
  discussInThisChat: (id: string, dir: string) =>
    `The cockpit's Discuss on ${id}, without Herdr: ticket ${id} waits for a decision of mine. Read it with jaira -C '${dir}' show ${id} --brief (the question field and the latest notes), explain the question and the options to me in plain words and ask me. ` +
    `When I have decided: write the decision in one sentence onto the ticket with jaira -C '${dir}' note ${id} "decision: <decision>" and clear the question with jaira -C '${dir}' set ${id} question=. If the decision means the work goes on, move the ticket back yourself, exactly so: jaira -C '${dir}' move ${id} --to in-progress (without --what/--why/--resolves: the ticket has its outcome already, and the decision is in the note). ` +
    `Accept it into the next lane only when I say yes to that explicitly. Change no code. Answer in the language I speak with you.`,
  opened: (id: string) => `${id}: opened in a Herdr tab of its own`,
  context: (k: number) => `context ${k}k`,
  cacheLeft: (min: number) => `cache ${min} min left`,
  cacheCold: 'cache cold',
  hintCompact: '/compact advised',
  hintNew: 'start a new session',
  guardQuestion: (k: number, min: number, usd: string) => `The prompt cache went cold ${min} min ago. Sending writes ${k}k tokens to it again (about ${usd}). What now?`,
  guardSend: 'Send',
  guardCompact: 'Compact first',
  guardClear: 'Start over (/clear)',
  guardCancel: 'Cancel',
  guardHeader: 'Cold cache',
  threadsOpen: (q: number, s: number) => `threads: ${q} open question(s), ${s} left lying`,
  threadsNone: 'No conversation threads yet.',
  threadsTitle: 'Conversation threads (? open question to you, ~ left lying, -> ticket):',
  workers: (n: number) => `${n} worker`,
  reviews: (n: number) => `${n} review`,
  waitingFor: (n: number) => `${n} waiting for you`,
  usage: 'session',
  workersTitle: 'Agents carrying a ticket of this board:',
  workersNone: 'No agent carries a ticket of this board.',
  workersNoBoard: 'No jaira board here.',
  paneTitle: (id: string) => `Ticket ${id}`,
  paneReports: 'Reports',
  paneUsage: 'Usage per agent',
  paneNothing: 'nothing recorded yet',
  usageText:
    'Usage: /cockpit <subcommand>\n\n  report [id] [all]  full text of the agent reports shortened for a ticket that were not read yet, newest first; all: the older ones too; without an id, the tickets that have some\n  threads            the conversation threads of this session and where each stands\n  ticket <id>        a pane with the ticket\'s reports and what its agents used\n  workers            every agent carrying a ticket of this board: lane, context, cache, and whether a new round goes to it',
}

// The noise filter's own words: the digest line and /cockpit report.
export const filterWords = {
  // Leads every digest, so the model never takes the line for the agent's
  // own words.
  marker: '[shortened by jaira-cockpit]',
  question: 'question',
  fullReport: 'full report',
  none: 'No agent reports kept for this project yet.',
  noneFor: (id: string) => `No reports kept for ${id} in this project.`,
  nothingNew: (id: string, older: number) =>
    `Nothing new for ${id} since it was last read; ${older} older report(s) kept: /cockpit report ${id} all`,
  list: 'Tickets with kept reports in this project, most recent first:',
  hint: '/cockpit report <id> shows one.',
  count: (n: number) => `${n} report(s)`,
  unread: (n: number) => `${n} new`,
  last: 'last',
  of: 'of',
}

// The Overview pane's words: the board at a glance and its last days.
export const overviewWords = {
  button: 'Overview',
  title: 'Board overview',
  working: 'In work',
  queue: 'Queue',
  parked: 'Put back',
  tabQueue: 'Queue',
  tabLog: 'Log',
  tabWorkers: 'Workers',
  // What a figure on the Workers tab is doing, beneath it in so many words.
  states: { working: 'working', idle: 'resting', question: 'waits for an answer' },
  noWorkers: 'No agent on the board right now.',
  everybody: 'All',
  none: 'none',
  history: (days: number) => `Created and put back, last ${days} days`,
  created: 'created',
  putBack: 'put back',
  by: ' by ',
  unread: 'The board could not be read (jaira lanes, list or next failed).',
  noRoom: 'No room for a pane here: widen the terminal and press Overview again.',
  close: 'Close',
}
