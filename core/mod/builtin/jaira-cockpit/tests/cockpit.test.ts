import type { On, SessionReceiveOrigin } from 'claude-code'
import { describe, expect, mock, test } from 'claude-code/testing'

// Beneath the plugin: a session that queues whatever reaches it, and a store.
function session(on: On) {
  mock.store(on)
  on('session.receive', ($, e) => ({ text: e.text }))
}

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

  test("a person's message is never touched", async ($, on) => {
    session(on)
    const mine = 'Ticket MG677R: move it review -> signoff, verdict: approve'
    expect((await $.session.receive({ origin: { kind: 'bridge' }, text: mine })).text).toBe(mine)
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
    const ran = await $.tool.call({ tool: 'mcp__jaira-cockpit__flight', input: { ticket: 'MG677R' } })
    expect(String(ran.result)).toContain('forty more lines of findings')
  })
})
