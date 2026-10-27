import type { Ticket } from '@/entities/ticket';

// Which date the from-to range reads.
export type DateField = 'created' | 'updated' | 'filed';

// One value picked in a lane, tag or person filter: the ticket must carry it,
// or, with not, must not.
export interface Chip {
  value: string;
  not: boolean;
}

// The filters that take chips. Each name is its key in the address; the
// values it excludes go under not-<name>.
const CHIPPED = ['lane', 'tag', 'person'] as const;

// What the page filters by. Every field set narrows the list further: the
// filters are ANDed, as 'jaira list' ANDs its flags. Within one chipped filter
// the values a ticket must carry are ORed, and every excluded one is ANDed on.
export interface Filter {
  // Substring of the title, the goal or the id, any case.
  q: string;
  lane: Chip[];
  tag: Chip[];
  // Creator, assignee or whoever changed it last.
  person: Chip[];
  date: DateField;
  // yyyy-mm-dd, both ends included; empty is open.
  from: string;
  to: string;
  // Filed tickets too, not only the board.
  all: boolean;
}

export const emptyFilter = (): Filter => ({ q: '', lane: [], tag: [], person: [], date: 'created', from: '', to: '', all: false });

const DATES: DateField[] = ['created', 'updated', 'filed'];

// The filter lives in the address, so a link shows the same list.
export function fromQuery(search: string): Filter {
  const p = new URLSearchParams(search);
  const f = emptyFilter();
  f.q = p.get('q') ?? '';
  // ?lane=x from before the chips reads as one value the ticket must carry.
  for (const k of CHIPPED) {
    for (const v of p.getAll(k)) f[k] = addChip(f[k], v);
    for (const v of p.getAll(`not-${k}`)) f[k] = addChip(f[k], v, true);
  }
  const date = p.get('date') as DateField;
  if (DATES.includes(date)) f.date = date;
  f.from = p.get('from') ?? '';
  f.to = p.get('to') ?? '';
  f.all = p.get('all') === '1';
  return f;
}

// Only what is set goes into the address; an unfiltered page keeps a bare one.
export function toQuery(f: Filter): string {
  const p = new URLSearchParams();
  if (f.q) p.set('q', f.q);
  for (const k of CHIPPED) {
    for (const c of f[k]) p.append(c.not ? `not-${k}` : k, c.value);
  }
  for (const k of ['from', 'to'] as const) {
    if (f[k]) p.set(k, f[k]);
  }
  if ((f.from || f.to) && f.date !== 'created') p.set('date', f.date);
  if (f.all) p.set('all', '1');
  const s = p.toString();
  return s ? `?${s}` : '';
}

// Whether anything narrows the list; the switch only widens it.
export const narrows = (f: Filter) => !!(f.q || f.lane.length || f.tag.length || f.person.length || f.from || f.to);

// A value picked once more is not a second chip: it keeps the place it has.
export const addChip = (chips: Chip[], value: string, not = false): Chip[] =>
  chips.some((c) => c.value === value) ? chips : [...chips, { value, not }];

export const toggleChip = (chips: Chip[], value: string): Chip[] =>
  chips.map((c) => (c.value === value ? { ...c, not: !c.not } : c));

export const removeChip = (chips: Chip[], value: string): Chip[] => chips.filter((c) => c.value !== value);

// Whether a ticket carrying these values passes the chips: one of the wanted
// ones, when any are wanted, and none of the excluded.
function passes(chips: Chip[], carried: (string | undefined)[]): boolean {
  const wanted = chips.filter((c) => !c.not);
  if (wanted.length && !wanted.some((c) => carried.includes(c.value))) return false;
  return !chips.some((c) => c.not && carried.includes(c.value));
}

// The local day of an RFC 3339 time, comparable as a string with the
// range's yyyy-mm-dd.
function day(iso: string | undefined): string {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

const dateOf = (t: Ticket, f: DateField) => day(f === 'created' ? t.created_at : f === 'updated' ? t.updated_at : t.filed_at);

export function matches(t: Ticket, f: Filter): boolean {
  if (f.q) {
    const q = f.q.toLowerCase();
    if (![t.title, t.goal ?? '', t.handle, t.id].some((s) => s.toLowerCase().includes(q))) return false;
  }
  if (!passes(f.lane, [t.lane])) return false;
  if (!passes(f.tag, t.tags ?? [])) return false;
  if (!passes(f.person, [t.creator, t.assignee, t.updated_by])) return false;
  if (f.from || f.to) {
    // A ticket with no such date (never filed, or filed before it was
    // recorded) is outside every range.
    const d = dateOf(t, f.date);
    if (!d || (f.from && d < f.from) || (f.to && d > f.to)) return false;
  }
  return true;
}

// The values the dropdowns offer: what the tickets on the page carry.
export function choices(tickets: Ticket[]) {
  const lanes = new Set<string>();
  const tags = new Set<string>();
  const people = new Set<string>();
  for (const t of tickets) {
    if (t.lane) lanes.add(t.lane);
    for (const g of t.tags ?? []) tags.add(g);
    for (const p of [t.creator, t.assignee, t.updated_by]) if (p) people.add(p);
  }
  const sorted = (s: Set<string>) => [...s].sort((a, b) => a.localeCompare(b));
  return { lanes: sorted(lanes), tags: sorted(tags), people: sorted(people) };
}
