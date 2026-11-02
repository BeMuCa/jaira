import type { External } from '@/entities/epic';
import type { EpicEdit, EpicFields } from '../api/set';

const FIELDS = ['title', 'description', 'body'] as const;

// editOf is what a save sends: the fields the person changed against what
// the form was opened on, and that base. Null when nothing changed.
export function editOf(base: EpicFields, mine: EpicFields): EpicEdit | null {
  const edit: EpicEdit = { base };
  for (const f of FIELDS) if (mine[f] !== base[f]) edit[f] = mine[f];
  return Object.keys(edit).length === 1 ? null : edit;
}

// keysEdit is what a save sends for the keys: the rows the person added or
// relinked, and the keys they took out, against the keys the form opened on.
// Keys compare in any case, as jaira compares them. An empty key row is left
// out; nothing changed is an empty object.
export function keysEdit(read: External[], rows: { key: string; url: string }[]): Pick<EpicEdit, 'add_external' | 'drop_external'> {
  const norm = (k: string) => k.trim().toUpperCase();
  const kept = rows.filter((r) => r.key.trim());
  const add = kept
    .filter((r) => !read.some((x) => x.key === r.key.trim() && x.url === r.url.trim()))
    .map((r) => `${r.key.trim()} ${r.url.trim()}`.trim());
  const drop = read.filter((x) => !kept.some((r) => norm(r.key) === norm(x.key))).map((x) => x.key);
  return { ...(add.length ? { add_external: add } : {}), ...(drop.length ? { drop_external: drop } : {}) };
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
