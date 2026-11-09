// What the page asks before it writes past a refusal the server would make:
// the server refuses both without force, so a page that sent force without
// asking would make the refusal mean nothing. null means no question — the
// write goes as it is, without force.

interface Closable {
  title: string;
  done: number;
  total: number;
  closed?: boolean;
}

// Closing an epic that still counts unfinished tickets.
export function closeQuestion(epic: Closable): string | null {
  const open = epic.total - epic.done;
  if (open <= 0) return null;
  return `${open} of ${epic.total} ticket(s) of "${epic.title}" are not finished. Close it anyway?`;
}

// Putting a ticket into a closed epic. It stays closed.
export function addQuestion(epic: Closable): string | null {
  if (!epic.closed) return null;
  return `"${epic.title}" is closed. Put the ticket in anyway? The epic stays closed.`;
}

// The epics the list shows: closed ones only when asked for.
export const visible = <T extends { closed?: boolean }>(epics: T[], showClosed: boolean): T[] =>
  showClosed ? epics : epics.filter((e) => !e.closed);
