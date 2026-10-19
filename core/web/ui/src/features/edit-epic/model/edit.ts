import type { EpicEdit, EpicFields } from '../api/set';

const FIELDS = ['title', 'description', 'body'] as const;

// editOf is what a save sends: the fields the person changed against what
// the form was opened on, and that base. Null when nothing changed.
export function editOf(base: EpicFields, mine: EpicFields): EpicEdit | null {
  const edit: EpicEdit = { base };
  for (const f of FIELDS) if (mine[f] !== base[f]) edit[f] = mine[f];
  return Object.keys(edit).length === 1 ? null : edit;
}

// rebase answers a refused save: somebody changed the epic since the form
// was opened. The new state becomes the base, so a second save, made knowing
// it, goes through. A field the person typed in keeps what they typed; a field
// they left alone takes the new value, or that second save would quietly put
// the old one back. changed names the fields the other side changed.
export function rebase(base: EpicFields, mine: EpicFields, now: EpicFields): { base: EpicFields; mine: EpicFields; changed: (keyof EpicFields)[] } {
  const next = { ...mine };
  for (const f of FIELDS) if (mine[f] === base[f]) next[f] = now[f];
  return { base: { ...now }, mine: next, changed: FIELDS.filter((f) => now[f] !== base[f]) };
}
