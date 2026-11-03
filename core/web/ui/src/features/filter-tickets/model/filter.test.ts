import { test } from 'node:test';
import assert from 'node:assert/strict';
import { addChip, emptyFilter, epicPasses, fromQuery, matches, narrows, removeChip, toggleChip, toQuery, wantsFiled, whereOf } from './filter.ts';

const ticket = (lane: string, tags: string[], assignee?: string) => ({
  id: lane + tags.join(),
  handle: 'AAAAAA',
  title: 't',
  lane,
  done: false,
  direct: true,
  place: 'board' as const,
  tags,
  creator: 'Alex',
  assignee,
});

test('a link from before the chips still filters by its one value', () => {
  const f = fromQuery('?lane=review&tag=web&person=Berk');
  assert.deepEqual(f.lane, [{ value: 'review', not: false }]);
  assert.deepEqual(f.tag, [{ value: 'web', not: false }]);
  assert.deepEqual(f.person, [{ value: 'Berk', not: false }]);
});

test('several values and excluded ones go into the link and come back', () => {
  let f = emptyFilter();
  f.lane = addChip(addChip(f.lane, 'review'), 'testing');
  f.lane = addChip(f.lane, 'done', true);
  f.person = addChip(f.person, 'Alex Sacharov', true);
  const q = toQuery(f);
  assert.equal(q, '?lane=review&lane=testing&not-lane=done&not-person=Alex+Sacharov');
  assert.deepEqual(fromQuery(q), f);
  assert.equal(toQuery(emptyFilter()), '');
});

test('values within a filter are ORed, filters are ANDed', () => {
  const f = fromQuery('?lane=review&lane=testing&tag=web');
  assert.ok(matches(ticket('review', ['web']), f));
  assert.ok(matches(ticket('testing', ['web', 'cli']), f));
  assert.ok(!matches(ticket('done', ['web']), f), 'lane not picked');
  assert.ok(!matches(ticket('review', ['cli']), f), 'tag not picked');
});

test('an excluded value drops every ticket carrying it', () => {
  const allButDone = fromQuery('?not-lane=done');
  assert.ok(matches(ticket('review', []), allButDone));
  assert.ok(!matches(ticket('done', []), allButDone));
  assert.ok(narrows(allButDone));

  const notCli = fromQuery('?tag=web&not-tag=cli');
  assert.ok(matches(ticket('review', ['web']), notCli));
  assert.ok(!matches(ticket('review', ['web', 'cli']), notCli));

  // A person counts as on the ticket as creator, assignee or last changer,
  // so excluding them drops it in each of those roles.
  const notAlex = fromQuery('?not-person=Alex');
  assert.ok(!matches(ticket('review', [], 'Berk'), notAlex), 'Alex created it');
  const notBerk = fromQuery('?not-person=Berk');
  assert.ok(!matches(ticket('review', [], 'Berk'), notBerk));
  assert.ok(matches(ticket('review', []), notBerk));
});

test('a chip is picked once, toggles between X and not X, and goes', () => {
  let chips = addChip([], 'web');
  assert.deepEqual(addChip(chips, 'web', true), chips, 'picked again: no second chip');
  chips = toggleChip(chips, 'web');
  assert.deepEqual(chips, [{ value: 'web', not: true }]);
  chips = toggleChip(chips, 'web');
  assert.deepEqual(chips, [{ value: 'web', not: false }]);
  assert.deepEqual(removeChip(chips, 'web'), []);
});

const at = (place: 'board' | 'ref' | 'logbook' | 'archive', rejected?: string) => ({ ...ticket('done', []), place, rejected });

test('the place filter keeps the board, the logbook, the archive and the rejected apart', () => {
  assert.equal(whereOf(at('board')), 'board');
  assert.equal(whereOf(at('ref')), 'board', 'only on its ref is still on the board');
  assert.equal(whereOf(at('logbook')), 'logbook');
  assert.equal(whereOf(at('archive')), 'archive');
  assert.equal(whereOf(at('archive', 'wontfix')), 'rejected');

  const tickets = [at('board'), at('logbook'), at('archive'), at('archive', 'wontfix')];
  const shown = (q: string) => tickets.filter((t) => matches(t, fromQuery(q))).map(whereOf);
  assert.deepEqual(shown('?where=board'), ['board']);
  assert.deepEqual(shown('?where=logbook'), ['logbook']);
  assert.deepEqual(shown('?where=archive'), ['archive'], 'the archive without the rejected');
  assert.deepEqual(shown('?where=rejected'), ['rejected']);
  assert.deepEqual(shown('?where=logbook&where=archive'), ['logbook', 'archive']);
  assert.deepEqual(shown('?not-where=board'), ['logbook', 'archive', 'rejected']);
  assert.deepEqual(shown(''), ['board', 'logbook', 'archive', 'rejected'], 'no place picked filters nothing');
});

test('the server is asked for filed tickets only when a place off the board is picked', () => {
  assert.ok(!wantsFiled(emptyFilter()));
  assert.ok(!wantsFiled(fromQuery('?where=board')));
  assert.ok(wantsFiled(fromQuery('?where=archive')));
  assert.ok(wantsFiled(fromQuery('?not-where=board')));
  assert.ok(narrows(fromQuery('?where=rejected')));
});

test('a link from before the place filter, ?all=1, picks every place', () => {
  const f = fromQuery('?all=1');
  assert.deepEqual(f.where.map((c) => c.value), ['board', 'logbook', 'archive', 'rejected']);
  assert.ok(wantsFiled(f));
  assert.equal(toQuery(f), '?where=board&where=logbook&where=archive&where=rejected');
  assert.deepEqual(fromQuery(toQuery(f)), f);
});

// A key filter keeps the tickets of the epics carrying the key, in any case,
// and goes into the link like every chip.
test('the key filter matches the keys of the epics a ticket counts towards', () => {
  const f = fromQuery('?key=jira-1');
  assert.deepEqual(f.key, [{ value: 'jira-1', not: false }]);
  assert.equal(toQuery(f), '?key=jira-1');
  assert.ok(narrows(f));
  const t = ticket('review', []);
  assert.ok(matches(t, f, ['JIRA-1', 'YT-7']));
  assert.ok(!matches(t, f, ['YT-7']));
  assert.ok(!matches(t, f));
  assert.ok(epicPasses(['Jira-1'], f));
  assert.ok(!epicPasses([], f));
  assert.ok(epicPasses([], emptyFilter()));
});
