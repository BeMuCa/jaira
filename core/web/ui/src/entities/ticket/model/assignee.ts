// A card has room for a short name: at most this many characters, the
// ellipsis included. The full name goes into the tooltip.
export const ASSIGNEE_MAX = 20;

// Counted in code points, so a name with an umlaut or an emoji is not cut in
// the middle of a character.
export function shortAssignee(name: string): string {
  const chars = Array.from(name);
  if (chars.length <= ASSIGNEE_MAX) return name;
  return chars.slice(0, ASSIGNEE_MAX - 1).join('') + '…';
}
