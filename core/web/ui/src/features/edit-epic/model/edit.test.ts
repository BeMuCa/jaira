import { test } from 'node:test';
import assert from 'node:assert/strict';
import { editOf, keysEdit, linkable, rebase } from './edit.ts';

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

// The key rows are sent as changes: a new row and a new link are added, a
// row taken out is dropped, an untouched one is not sent at all.
test('the keys go as what was added and what was taken out', () => {
  const read = [{ key: 'JIRA-1', url: 'https://j.example/1' }, { key: 'YT-7', url: '' }, { key: 'OLD-2', url: '' }];
  const rows = [{ key: 'JIRA-1', url: 'https://j.example/one' }, { key: 'YT-7', url: '' }, { key: ' NEW-3 ', url: '' }, { key: '', url: 'x' }];
  assert.deepEqual(keysEdit(read, rows), { add_external: ['JIRA-1 https://j.example/one', 'NEW-3'], drop_external: ['OLD-2'] });
  assert.deepEqual(keysEdit(read, read.map((x) => ({ ...x }))), {});
});

// A row whose link the server would refuse is left out of the save and not
// dropped either; every other change still goes.
test('a bad link holds back its own row only', () => {
  const read = [{ key: 'C-3', url: 'not a link' }];
  const rows = [{ key: 'C-3', url: 'not a link' }, { key: 'NEW-1', url: 'javascript:alert(1)' }, { key: 'OK-2', url: 'https://ok.example/2' }];
  assert.deepEqual(keysEdit(read, rows), { add_external: ['OK-2 https://ok.example/2'] });
  assert.ok(linkable('') && linkable('https://x.example/a') && linkable('http://x.example'));
  assert.ok(!linkable('not a link') && !linkable('ftp://x.example') && !linkable('https://'));
});
