import { test } from 'node:test';
import assert from 'node:assert/strict';
import { addChip, emptyFilter, fromQuery, matches, narrows, removeChip, toggleChip, toQuery } from './filter.ts';

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

const filter = (q: string) => fromQuery(q);

test('a link from before the chips still filters by its one value', () => {
  const f = filter('?lane=review&tag=web&person=Berk');
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
  const f = filter('?lane=review&lane=testing&tag=web');
  assert.ok(matches(ticket('review', ['web']), f));
  assert.ok(matches(ticket('testing', ['web', 'cli']), f));
  assert.ok(!matches(ticket('done', ['web']), f), 'lane not picked');
  assert.ok(!matches(ticket('review', ['cli']), f), 'tag not picked');
});

test('an excluded value drops every ticket carrying it', () => {
  const allButDone = filter('?not-lane=done');
  assert.ok(matches(ticket('review', []), allButDone));
  assert.ok(!matches(ticket('done', []), allButDone));
  assert.ok(narrows(allButDone));

  const notCli = filter('?tag=web&not-tag=cli');
  assert.ok(matches(ticket('review', ['web']), notCli));
  assert.ok(!matches(ticket('review', ['web', 'cli']), notCli));

  // A person counts as on the ticket as creator, assignee or last changer,
  // so excluding them drops it in each of those roles.
  const notAlex = filter('?not-person=Alex');
  assert.ok(!matches(ticket('review', [], 'Berk'), notAlex), 'Alex created it');
  const notBerk = filter('?not-person=Berk');
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
