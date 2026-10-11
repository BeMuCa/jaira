import type { On, SessionReceiveOrigin } from 'claude-code'
import { describe, expect, mock, test } from 'claude-code/testing'

// Beneath the plugin: a session that queues whatever reaches it, and a store.
function session(on: On) {
  mock.store(on)
  on('session.receive', ($, e) => ({ text: e.text }))
  on('prompt.submit', ($, e) => ({ text: e.text, origin: e.origin }))
}

// This session's own background task finishing: the engine wraps the
// notification in a warning and submits it as a prompt.
const LOCAL = (text: string) =>
  ({ text: `[SYSTEM NOTIFICATION - NOT USER INPUT]\n\n${text}`, origin: { kind: 'task-notification' }, wait: false }) as const

// A person typing `/flight <args>` at the prompt.
const typed = (args = '') =>
  ({ command: 'flight', args, origin: { kind: 'composer' }, presentation: { isFullscreen: false, columns: 120 } }) as const

const NOTIFICATION: SessionReceiveOrigin = { kind: 'task-notification' }

const REPORT = [
  '<task-notification>',
  '<task-id>a81f</task-id>',
  '<status>completed</status>',
  '<summary>Agent "review MG677R" completed</summary>',
  '<result>',
  'Review of ticket MG677R done. Moved review -> signoff.',
  'review-verdict="approve with reservations: the cache key ignores the branch"',
  'open question: should the cache survive a rebase?',
  '... forty more lines of findings ...',
  '</result>',
  '</task-notification>',
].join('\n')

describe('session.receive', () => {
  test('a report naming a ticket becomes one line and is kept', async ($, on) => {
    session(on)
    const { text } = await $.session.receive({ origin: NOTIFICATION, text: REPORT })
    expect(text).toBe(
      'MG677R · review → signoff · approve with reservations: the cache key ignores the branch · ' +
        'question: should the cache survive a rebase? — full report: /flight MG677R',
    )
    const shown = await $.command.run(typed('mg677r'))
    expect(shown.text).toContain('forty more lines of findings')
  })

  test("a worker's three-line report through SendMessage", async ($, on) => {
    session(on)
    const worker = [
      'T8QS1Y: noise filter mod under core/mod/builtin/jaira-cockpit, 6 tests green.',
      'Ticket T8QS1Y sits in signoff.',
      'Open: nothing; a person accepts it.',
    ].join('\n')
    const { text } = await $.session.receive({ origin: { kind: 'peer-send-message' }, text: worker })
    expect(text).toBe('T8QS1Y · → signoff — full report: /flight T8QS1Y')
  })

  test('a verdict keeps its reservation', async ($, on) => {
    session(on)
    const review = 'Ticket MG677R reviewed: approve mit Vorbehalt, moved to signoff.'
    const { text } = await $.session.receive({ origin: NOTIFICATION, text: review })
    expect(text).toBe('MG677R · → signoff · approve mit vorbehalt — full report: /flight MG677R')
  })

  test('a report without a ticket passes unchanged', async ($, on) => {
    session(on)
    const plain = 'Agent finished: refactored the parser, all tests pass.'
    const { text } = await $.session.receive({ origin: NOTIFICATION, text: plain })
    expect(text).toBe(plain)
  })

  test('a failed agent stays whole, ticket or not', async ($, on) => {
    session(on)
    const failed = REPORT.replace('<status>completed</status>', '<status>failed</status>')
    expect((await $.session.receive({ origin: NOTIFICATION, text: failed })).text).toBe(failed)
    const aborted = 'Worker on ticket MG677R aborted: usage limit reached in critique.'
    expect((await $.session.receive({ origin: { kind: 'peer' }, text: aborted })).text).toBe(aborted)
  })

  test('a completed agent that reports a failure stays whole, as a peer would', async ($, on) => {
    session(on)
    const failing = REPORT.replace('Moved review -> signoff.', 'MG677R: go test ./... exit status 1, worker stopped.')
    expect((await $.session.receive({ origin: NOTIFICATION, text: failing })).text).toBe(failing)
    expect((await $.session.receive({ origin: { kind: 'peer' }, text: failing })).text).toBe(failing)
  })

  test('failed tests stay whole, none failed is shortened', async ($, on) => {
    session(on)
    const some = 'MG677R: in-progress -> critique. go test ./...: 41 passed, 3 failed.'
    expect((await $.session.receive({ origin: NOTIFICATION, text: some })).text).toBe(some)
    const go = 'MG677R: in-progress -> critique.\n--- FAIL: TestDigest (0.00s)\nFAIL\tjaira/core/lane'
    expect((await $.session.receive({ origin: NOTIFICATION, text: go })).text).toBe(go)
    const none = 'MG677R: in-progress -> critique. go test ./...: 44 passed, 0 failed.'
    expect((await $.session.receive({ origin: NOTIFICATION, text: none })).text).toBe(
      'MG677R · in-progress → critique — full report: /flight MG677R',
    )
  })

  test("a person's message is never touched", async ($, on) => {
    session(on)
    const mine = 'Ticket MG677R: move it review -> signoff, verdict: approve'
    expect((await $.session.receive({ origin: { kind: 'bridge' }, text: mine })).text).toBe(mine)
  })
})

describe('prompt.submit', () => {
  test("this session's own task notification is shortened, its wrapper kept", async ($, on) => {
    session(on)
    const { text } = await $.prompt.submit(LOCAL(REPORT))
    expect(text).toBe(
      '[SYSTEM NOTIFICATION - NOT USER INPUT]\n\nMG677R · review → signoff · approve with reservations: the cache key ignores the branch · ' +
        'question: should the cache survive a rebase? — full report: /flight MG677R',
    )
    expect((await $.command.run(typed('MG677R'))).text).not.toContain('SYSTEM NOTIFICATION')
  })

  test('a line already shortened is not kept a second time', async ($, on) => {
    session(on)
    const { text = '' } = await $.session.receive({ origin: { kind: 'peer' }, text: 'ticket K7P2QX moved to testing' })
    expect((await $.prompt.submit({ text, origin: { kind: 'peer' }, wait: false })).text).toBe(text)
    expect((await $.command.run(typed('K7P2QX'))).text).toContain('report 1 of 1')
  })

  test('a comment on an artifact stays whole on both hooks', async ($, on) => {
    session(on)
    const comment = '<task-notification>\nA comment on the page: ticket MG677R looks wrong, move review -> signoff?\n</task-notification>'
    const event = { source: 'artifact-room', kind: 'events', data: {}, untrustedKeys: [] }
    const { text = '' } = await $.session.receive({ origin: NOTIFICATION, text: comment, event })
    expect(text).toBe(comment)
    expect((await $.prompt.submit({ text, origin: NOTIFICATION, wait: false })).text).toBe(comment)
  })

  test('what a person types is never touched', async ($, on) => {
    session(on)
    const mine = 'Ticket MG677R: move it review -> signoff, verdict: approve'
    expect((await $.prompt.submit({ text: mine, origin: { kind: 'composer' }, wait: false })).text).toBe(mine)
  })
})

describe('/flight', () => {
  test('lists reports newest first', async ($, on) => {
    session(on)
    await $.session.receive({ origin: { kind: 'peer' }, text: 'first: ticket K7P2QX moved to critique' })
    await $.session.receive({ origin: { kind: 'peer' }, text: 'second: ticket K7P2QX moved to testing' })
    const { text = '' } = await $.command.run(typed('K7P2QX'))
    expect(text.indexOf('second')).toBeLessThan(text.indexOf('first'))
    expect((await $.command.run(typed())).text).toContain('K7P2QX')
  })
})

describe('the flight tool', () => {
  test('gives the model the same full text', async ($, on) => {
    session(on)
    await $.session.receive({ origin: NOTIFICATION, text: REPORT })
    const ran = await $.tool.call({ tool: 'mcp__jaira-cockpit__flight', ticket: 'MG677R' })
    expect(String(ran.result)).toContain('forty more lines of findings')
  })
})
