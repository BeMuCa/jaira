// The engine wraps what an agent sends in frames of its own: <agent-message>
// around a SendMessage, a [Subagent hand-back] disclaimer above a subagent's
// final report (every report line indented two spaces), <task-notification>
// XML around a background task's end. /cockpit report shows the report, so
// the frame is taken off here and kept as one header line. Pure, like digest.ts.

const HANDBACK = /\[Subagent hand-back\][\s\S]*?The report follows:\s*\n/
const AGENT_MESSAGE = /<agent-message\b[^>]*>\s*([\s\S]*?)\s*<\/agent-message>/
const NOTIFICATION = /<task-notification>([\s\S]*?)<\/task-notification>/

function tag(xml: string, name: string): string | undefined {
  return new RegExp(`<${name}>\\s*([\\s\\S]*?)\\s*</${name}>`).exec(xml)?.[1]
}

function dedent(text: string): string {
  return text.replace(/^ {2}/gm, '')
}

// What the agent itself said, with no frame of the engine's around it.
export function body(text: string): string {
  const handback = HANDBACK.exec(text)
  if (handback !== null) return dedent(text.slice(handback.index + handback[0].length)).trim()
  const message = AGENT_MESSAGE.exec(text)
  if (message !== null) return body(message[1]!)
  const notification = NOTIFICATION.exec(text)
  if (notification !== null) {
    const inner = notification[1]!
    const result = tag(inner, 'result')
    if (result !== undefined) return body(result)
    return tag(inner, 'summary') ?? inner.replace(/<[^>]+>/g, '').trim()
  }
  return text.trim()
}

// Who sent it, in a few words: the origin, and the agent where the frame
// names one ("peer · general-purpose", "task-notification · Agent \"review\"").
export function source(text: string, origin: string, agent?: string): string {
  const named = agent ?? tag(NOTIFICATION.exec(text)?.[1] ?? '', 'summary')
  return named === undefined || named === '' ? origin : `${origin} · ${named}`
}
