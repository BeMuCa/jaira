import { test } from 'node:test';
import assert from 'node:assert/strict';
import { editOf, rebase } from './edit.ts';

const read = { title: 'Login', description: 'no password', body: '## Goal' };

test('a save sends only what the person changed, beside what they read', () => {
  assert.deepEqual(editOf(read, { ...read, title: 'Anmeldung' }), { base: read, title: 'Anmeldung' });
  assert.equal(editOf(read, { ...read }), null);
});

// After a 409 the typed text stays, the form knows the new state, and the
// second save is checked against that state, so it goes through. A field the
// person never touched takes the other side's value instead of reverting it.
test('a refused save keeps the typing and saves against the new state', () => {
  const mine = { ...read, title: 'Mein Titel', body: '## Goal\n\nmore' };
  const now = { ...read, title: 'Anders', description: 'per Link' };
  const r = rebase(read, mine, now);
  assert.deepEqual(r.base, now);
  assert.deepEqual(r.changed, ['title', 'description']);
  assert.deepEqual(r.mine, { title: 'Mein Titel', description: 'per Link', body: '## Goal\n\nmore' });
  assert.deepEqual(editOf(r.base, r.mine), { base: now, title: 'Mein Titel', body: '## Goal\n\nmore' });
});
