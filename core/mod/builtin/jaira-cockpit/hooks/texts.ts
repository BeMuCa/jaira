// Every text the cockpit shows a person, in English whatever language the
// person speaks with Claude: one set of words, the same on every machine.
// The ticket's own texts are the one exception: they come from the ticket,
// and the Acceptance pane's Translate shows them in the "language" the person
// set in their own settings.json (register.tsx translate).

import { type LightReason, oneLine } from './board.ts'

// A text a person typed, as part of one line: whitespace collapsed (board.ts
// oneLine), cut at max characters with an ellipsis. Counted in code points:
// a cut by UTF-16 index splits an emoji and leaves a lone surrogate in the
// prompt this session is handed.
function inLine(text: string, max: number): string {
  const chars = Array.from(oneLine(text))
  return chars.length <= max ? chars.join('') : `${chars.slice(0, max - 1).join('').trimEnd()}…`
}

// A path in single quotes for a shell line a person or a chat pastes.
function shellQuoted(path: string): string {
  return `'${path.replaceAll("'", `'\\''`)}'`
}

// The last step of a chat the cockpit pulled the ticket for (DF46B9): only the
// chat knows when it wrote its last line, so it removes the copy itself.
function removeCopy(copy: string | undefined): string {
  return copy === undefined ? '' : ` As the very last step, once you write nothing more onto the ticket, decided or not: git clean -fqx -- ${shellQuoted(copy)} - the copy of the ticket the cockpit pulled for this chat; left there, it hides the ticket's ref.`
}

export const t = {
  waiting: 'Waiting for you',
  accept: 'Accept',
  giveBack: 'Give back',
  // The Acceptance pane (2V50Y5): one window for Accept, Give back and the
  // verdict - the ticket, how it got here, and what happens to each note.
  acceptance: 'Acceptance',
  acceptanceOf: (id: string) => `Acceptance of ${id}`,
  done: 'Done:',
  why: 'Why:',
  historyHeading: 'History',
  round: (n: number) => `round ${n}`,
  backFrom: (count: number, from: string) => `back ${count}x${from ? ` from ${from}` : ''}`,
  age: (ms: number) => {
    const min = Math.max(0, Math.floor(ms / 60_000))
    return min < 60 ? `${min} min` : min < 48 * 60 ? `${Math.floor(min / 60)} h` : `${Math.floor(min / 1440)} days`
  },
  dodCount: (done: number, total: number, superseded: number) => `DoD ${done}/${total - superseded}${superseded > 0 ? ` (${superseded} superseded)` : ''}`,
  huntOnHead: (state: string) => `bughunt ${state === 'ok' ? '✓ on HEAD' : state === 'fail' ? '✗ not on HEAD' : state === 'skipped' ? '– skipped' : '?'}`,
  review: (verdict: string) => `review: ${verdict || 'none yet'}`,
  commitsFiles: (commits: number, files: number) => `${commits} commit${commits === 1 ? '' : 's'}, ${files} file${files === 1 ? '' : 's'}`,
  prOpen: 'pull request exists',
  releaseNote: 'Release note',
  notesChoose: 'Review notes - choose for each',
  choiceOk: 'ok',
  choiceFix: 'fix',
  choiceTicket: 'ticket',
  explainNote: '?',
  explaining: 'explaining…',
  noReason: 'no explanation came back',
  explainAll: 'Explain all',
  translate: 'Translate',
  translating: 'translating…',
  original: 'Original',
  execute: 'Execute',
  // N5ZZP6: Execute while its row is locked - by its own press or an
  // Explain chat - drawn dim; a press on it does nothing.
  running: 'running…',
  // Execute's lock ran out while its step is still at work: it takes a
  // press again, and a second one may write beside the first.
  stillRunning: (id: string, seconds: number) => `${id}: Execute still at work after ${seconds} s - it takes a press again; look at the ticket before you press once more`,
  // What Execute will do, in the pane's last line before the button.
  // ends: whether to is the terminal lane; elsewhere the line names the lane
  // the ticket goes on to, so nobody reads it as done (M5JSE9).
  summary: (step: 'giveBack' | 'accept' | 'needWhy' | 'noNext' | 'needAnswer' | 'needSeen', fix: number, followUps: number, to: string, ends: boolean) => {
    const into = ends ? to : `accept → ${to}`
    const tickets = followUps > 0 ? ` + ${followUps} follow-up ticket${followUps === 1 ? '' : 's'}` : ''
    switch (step) {
      case 'giveBack': return `Execute: give back${fix > 0 ? ` with ${fix} note${fix === 1 ? '' : 's'} to fix` : ' with what you typed'}${tickets}`
      case 'accept': return `Execute: ${into}${tickets}`
      case 'needWhy': return `Execute: type why the blocker is accepted${tickets}`
      case 'noNext': return 'Execute: nothing to accept in this lane - mark a note fix, or type what is missing'
      case 'needAnswer': return 'Execute: answer the question first - or mark a note fix, or type what is missing'
      case 'needSeen': return `Execute: tick what you checked above first, then ${into}${tickets}`
    }
  },
  // DKMG91: the 'seen by a person' items of the definition of done, which
  // only the person ticks - listed above the review notes, ticked by Accept.
  checkYourself: 'Check yourself - tick what you have seen',
  seenMissing: (id: string, texts: string[]) => `${id}: not accepted - ${texts.length === 1 ? 'one item is' : `${texts.length} items are`} still to check yourself and tick:\n${texts.map(x => `- ${x}`).join('\n')}`,
  noNext: (id: string, lane: string) => `${id} waits in ${lane}, and jaira names no lane after it: nothing to accept there. Mark a note fix, or type what is missing, to give it back.`,
  needAnswer: (id: string, lane: string, to: string) => `${id} waits in ${lane} for the answer to its question: answer it first - Discuss, or an option marked [accept] under Explain - and Execute then moves it on to ${to}. Mark a note fix, or type what is missing, to give it back.`,
  refused: (id: string) => `${id} was not accepted, nothing was changed:`,
  refusedTicked: (id: string, ns: readonly number[]) => `${id} was not accepted; only definition-of-done item${ns.length === 1 ? '' : 's'} ${ns.join(', ')} ${ns.length === 1 ? 'was' : 'were'} ticked as seen by you:`,
  seenChanged: (texts: string[]) => `${texts.length === 1 ? 'this item' : 'these items'} changed on the ticket since the pane opened - close it, open it again and look once more:\n${texts.map(x => `- ${x}`).join('\n')}`,
  explain: 'Explain',
  discuss: 'Discuss',
  more: (n: number) => `+${n} more`,
  fewer: 'show fewer',
  theirs: (n: number, names: string[]) => `${n} waiting for ${names.join(', ')}`,
  reasonLabel: 'What is missing?',
  reasonPlaceholder: 'goes onto the ticket as a note',
  decidedBefore: '(decided before)',
  moreLabel: 'Anything else?',
  morePlaceholder: 'optional, goes onto the ticket with the notes',
  blockerNeedsWhy: (id: string) => `${id}: a blocker is not marked fix - type why it is accepted anyway, or mark it fix`,
  blockerNeedsTyping: (id: string) => `${id}: a blocker is not marked fix - accepting it needs a reason typed in a terminal or on the desktop; nothing was accepted`,
  whyLabel: 'Why accept it with the blocker?',
  whyPlaceholder: 'required - goes onto the ticket as a decision',
  acceptedDespite: 'accepted with a blocker not fixed',
  alreadyIn: (id: string, lane: string) => `${id} is already in ${lane}: someone else decided it meanwhile. Nothing was written.`,
  nothingToFix: (id: string) => `${id}: no note is marked fix now and nothing is typed - mark one, type what is missing, or press Execute to accept`,
  followUpsOf: (id: string) => `Follow-up tickets of ${id}`,
  followUpsLost: (id: string, texts: string[]) => `These notes of ${id} were marked follow-up, but jaira created no ticket for them; the note on ${id} names each as not created:\n${texts.map(x => `- ${x}`).join('\n')}`,
  movedMeanwhile: (id: string, lane: string, now: string) => `${id} is no longer in ${lane} but in ${now}: someone else decided it meanwhile. Nothing was written.`,
  accepted: (id: string, lane: string) => `${id} accepted: ${lane}`,
  returnedTo: (id: string, who: 'old' | 'new' | 'none') =>
    who === 'none' ? `${id} back in in-progress; no worker started (the ticket's note says why)` : `${id} given back to ${who === 'old' ? 'the worker that carried it' : 'a new worker'}`,
  tabsUnreadable: (id: string) => `${id} back in in-progress; Herdr's tab list could not be read, so it is unknown whether its worker still runs - no worker started. Check its tab, and start one only if none runs.`,
  workerBlocked: (id: string, where: string) => `${id} back in in-progress; its worker waits at an approval dialog in ${where} - answer it there, then type the decision into it. No Enter was pressed for you.`,
  dialogBeforeEnter: (id: string, pane: string) => `${id} given back to the worker that carried it; an approval dialog opened in pane ${pane} before the Enter - answer it there, then press Enter if the decision still waits in its input.`,
  // The one line Accept and Give back hand this session, the teamlead's,
  // after a move that went through: what happened, and what already runs.
  // Each says who decided — '(Person, Cockpit)' — and leads with the
  // ticket, so a teamlead stops listing it as open the moment it arrives.
  toSessionAccepted: (id: string, from: string, to: string, next: string) =>
    `jaira-cockpit: ${id}: accepted, ${from} -> ${to} (Person, Cockpit)${next ? `; ${inLine(next, 200)}` : ''}`,
  acceptedNote: (from: string, to: string) => `accepted by the person, ${from} -> ${to}`,
  // Two named after the ticket, or several holding its newest commit: both
  // are 'could each be its branch', and the line says only that.
  branchAmbiguous: (id: string, branches: string[]) => `${branches.join(' and ')} could each be ${id}'s branch and none is its worktree's: no branch picked for it - start its worker or pull request on the right one by hand`,
  acceptAmbiguous: (id: string, branches: string[], to: string) => `${branches.join(' and ')} could each carry ${id}'s commits and none is its worktree's: nothing moved - check out the right one and run 'jaira move ${id} --to ${to} --force' there`,
  acceptNoteLost: (id: string, why: string) => `${id} accepted; the note saying so was not written: ${oneLine(why).slice(0, 120)}`,
  toSessionGivenBack: (id: string, from: string, who: 'old' | 'new' | 'none', reason: string, worker: string) =>
    `jaira-cockpit: ${id} given back: ${inLine(reason, 120)}, worker ${who === 'none' ? worker || "none (the ticket's note says why)" : `${who === 'old' ? 'that carried it' : 'new'} ${worker}`} (Person, Cockpit; ${from} -> in-progress)`,
  afterAcceptFailed: (id: string) => `the next step after ${id} failed, see the cockpit`,
  failed: (what: string) => `${what} failed`,
  afterAccept: (id: string) => `Next step after ${id}`,
  staysDone: (id: string) => `${id} stays in done; start the next step by hand.`,
  // A pull request handed to the teamlead to follow (CRASZS): once each.
  followPr: (n: number, tickets: string[]) =>
    `jaira-cockpit: pull request #${n} for ${tickets.join(', ')} is open - follow it: run 'jaira pr --watch --number ${n}' under a Monitor (one line per change, it ends when the pull request is merged or closed)`,
  // .jaira/pr-when says signoff: Accept opens the pull request first, and
  // done waits for its review round.
  prOpensInSignoff: (id: string, branch: string) => `${id} stays in signoff: opening the pull request for ${branch} in a Herdr tab; press Accept again once its review round is settled`,
  toSessionPrInSignoff: (id: string, branch: string) =>
    `jaira-cockpit: the person pressed Accept on ${id}; this board opens the pull request in signoff, so the pr lane is opening it for ${branch} and ${id} stays in signoff until the review round is settled`,
  prElsewhere: (branch: string) => `the forge lists a pull request for ${branch}, but none is open from this repository's origin - merged, closed, or another fork's branch of the same name; check it and move the ticket by hand`,
  prWaitsSignoff: (branch: string, pending: string[]) => `the pull request for ${branch} waits: ${pending.join(', ')} not in signoff yet`,
  notSettled: (branch: string, unsettled: string[]) => `the review round of ${branch} is not settled:\n${unsettled.map(u => `- ${u}`).join('\n')}`,
  staysInLane: (id: string, lane: string) => `${id} stays in ${lane}; nothing was moved.`,
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
  copyNotRemoved: (id: string, path: string, why: string) => `${id}: the copy pulled into the main checkout is still there and hides its ref (${oneLine(why)}); remove it: git clean -fqx -- ${shellQuoted(path)}`,
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
  noRoom: 'No room for a pane here: widen the terminal and press Acceptance again.',
  readExplanation: 'Read explanation',
  reviewedIn: (round: number) => `review of round ${round}, changed since`,
  // The band's row: one mark per check, counts for the review's findings.
  markDod: 'DoD',
  markHunt: 'bughunt',
  markBlockers: (n: number) => `${n} blocker${n === 1 ? '' : 's'}`,
  markNotes: (n: number) => `${n} note${n === 1 ? '' : 's'}`,
  markOldReview: (round: number) => `review of round ${round}`,
  markBlockersUnknown: '? blockers',
  markNotesUnknown: '? notes',
  // The Acceptance pane: what the ticket is for, then how it got here.
  goal: 'Goal:',
  // The light on top of the Acceptance pane: its colour, the reason it names
  // first, and how many more there are.
  light: { green: 'green', yellow: 'yellow', red: 'red' } as Record<'green' | 'yellow' | 'red', string>,
  lightAccept: 'accept',
  lightMore: (n: number) => ` (+${n} more)`,
  lightReason: (r: LightReason): string => {
    const s = (n: number) => (n === 1 ? '' : 's')
    switch (r.kind) {
      case 'refused': return 'the review gives it back'
      case 'blockers': return `${r.n} blocker${s(r.n)} in the review`
      case 'staleHunt': return 'bughunt not on HEAD'
      case 'noHunt': return 'no bughunt yet'
      case 'dodOpen': return `${r.n} DoD item${s(r.n)} open`
      case 'dodIncomplete': return 'DoD not complete'
      case 'notReviewed': return t.noVerdict
      case 'freeVerdict': return "the verdict uses none of the lane's words"
      case 'reservation': return 'the review has a reservation'
      case 'oldReview': return t.reviewedIn(r.round)
      case 'notesUnknown': return 'review notes not kept'
      case 'notes': return `${r.n} review note${s(r.n)} to weigh`
      case 'byAPerson': return `${r.n} 'seen by a person' item${s(r.n)} open`
      case 'question': return 'an open question'
      case 'huntUnknown': return 'bughunt unknown'
    }
  },
  noVerdict: 'not reviewed yet',
  notesNotKept: 'The notes of that review were not kept: a later lane rewrote them.',
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
    `jaira-cockpit: ${id}: ${text} (Person, Cockpit) - decided in the explanation pane, noted on the ticket, not acted on`,
  answerOn: (id: string) => `Answer on ${id}`,
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
  explainBySubagent: (id: string, dir: string, copy?: string) =>
    `The cockpit's Explain on ${id}, without Herdr: start one subagent (Agent tool) with exactly the prompt below and do nothing on ${id} yourself. When it is done, tell me in one line that the explanation of ${id} is on the ticket, under Read explanation in the cockpit.\n\n` +
    `Read ticket ${id} with jaira -C '${dir}' show ${id} --brief (title, goal, question, review, the newest bug hunt and the latest notes). Explain the review on ${id} in plain words: where it stands, what each note in review-gaps means in one everyday sentence, which options there are and which you recommend. Where review.verdict_gaps is not empty and differs from review.gaps, the verdict is from an earlier round and a later lane rewrote review-gaps: the review's notes are those in review.verdict_gaps - read them wherever this says review-gaps, and say that the ticket changed since. ` +
    `Then write the same explanation onto the ticket as one note with exactly these lines: jaira -C '${dir}' note ${id} "explain: Status: DoD <done>/<all> ok|fail | tests ok|fail | blockers ok|fail | bughunt on HEAD ok|fail|unknown|skipped\n[<class>] <one note, one everyday sentence>\n1) <option, one line> [accept]\n2) <option, one line> [give back]\nRecommendation: <number> - <why, one sentence, naming the note it is about>" - with real line breaks in the text. ` +
    `The status line: DoD, tests and blockers each end in ok or fail, bughunt on HEAD in one of ok, fail, unknown or skipped; blockers is ok when review-gaps has no blocker: line; bughunt on HEAD ends in the brief's bughunt.state word as it stands - ok, fail, unknown (a hunted file is changed and not committed, so nobody can tell; never write fail for it) or skipped (the ticket runs no hunt after a fix, by choice) - and in fail when bughunt is null (no hunt finished); only a brief without a bughunt key, from an older jaira, falls back to: ok when the newest 'bughunt: done at <commit>' note names the commit git -C '${dir}' rev-parse --short HEAD prints. ` +
    `One [<class>] line per note in review-gaps, the class exactly one of: [blocker] - a wrong result, data loss, a crash, a security hole; [small edge] - formally fine, one round of about ten minutes closes it; [known limit] - cannot be solved here and is documented; [not this ticket] - belongs elsewhere, say where. A note that states its class keeps it - blocker: is [blocker], note: edge: [small edge], note: limit: [known limit], note: elsewhere: [not this ticket]; the cockpit shows that class beside yours. Classify first; nothing else per note, no line on why it matters. ` +
    `Two to four options with 1) 2) 3) 4) at the start of the line - the numbers the person picks by in the cockpit - in the order you recommend them, no fixed order. Each option ends in [give back] when it means giving the work back to change something, [accept] when it means accepting the ticket as it is, and no mark when it is neither. When no note is a blocker, the recommendation says plainly that the ticket can be accepted as it is. ` +
    `Change nothing else, commit nothing, do not move the ticket. Answer in the language the person speaks with this session; keep the words Status:, Recommendation:, ok, fail and the class marks in English, the cockpit reads them.` +
    removeCopy(copy),
  discussInThisChat: (id: string, dir: string, copy?: string) =>
    `The cockpit's Discuss on ${id}, without Herdr: ticket ${id} waits for a decision of mine. Read it with jaira -C '${dir}' show ${id} --brief (the question field and the latest notes), explain the question and the options to me in plain words and ask me. ` +
    `When I have decided: write the decision in one sentence onto the ticket with jaira -C '${dir}' note ${id} "decision: <decision>" and clear the question with jaira -C '${dir}' set ${id} question=. If the decision means the work goes on, move the ticket back yourself, exactly so: jaira -C '${dir}' move ${id} --to in-progress (without --what/--why/--resolves: the ticket has its outcome already, and the decision is in the note). ` +
    `Accept it into the next lane only when I say yes to that explicitly. Change no code. Answer in the language I speak with you.` +
    removeCopy(copy),
  copyLeft: (id: string, path: string) => `${id}: the copy pulled into the main checkout for a chat is still there and hides its ref; once no chat on ${id} runs any more, remove it: git clean -fqx -- ${shellQuoted(path)}`,
  chatHoldsCopy: (id: string, kind: string) => `The ${kind} chat on ${id} is still open and removes the copy of the ticket pulled for it when it ends. Finish that chat first.`,
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
  asksTitle: 'Waiting for you',
  asksNone: 'No question is waiting for you.',
  asksNoTicket: 'no ticket',
  answerInChat: 'Answer in the chat…',
  asksDone: 'x done',
  asksAnswered: 'Answered',
  askedAt: (when: string) => `asked ${when}`,
  answeredAt: (when: string) => `answered ${when}`,
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
  needsYou: 'needs you',
  fullReport: 'full report',
  // The fields of a worker's status line, as the digest names them.
  commits: 'commits',
  gates: 'gates',
  notes: 'notes',
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
  // What an agent on the Workers tab is doing, on its line in so many words.
  states: { working: 'working', idle: 'resting', question: 'waits for an answer', blocked: 'waits for approval' },
  noWorkers: 'No agent on the board right now.',
  // Above the stone: the definition-of-done items still open on the tickets
  // with an agent at them.
  stoneOpen: (n: number) => `${n} open`,
  // Agents the picture has no room for; their lines are beneath it all the same.
  moreWorkers: (n: number) => `+${n} more`,
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
