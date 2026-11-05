import { test } from 'node:test';
import assert from 'node:assert/strict';
import { epicCardExit, handleOf } from './epic-card.ts';

test('every member card can be taken out, on the board or filed away', () => {
  assert.deepEqual(epicCardExit({ direct: true }), { removable: true });
  // place does not enter into it: a logbook or archive member has its line too.
  assert.deepEqual(epicCardExit({ direct: true, via: '' }), { removable: true });
});

test('a follow-up names the member it is counted through instead', () => {
  assert.deepEqual(epicCardExit({ direct: false, via: '01M6NNAXQ276DMK3A1GH3XXF5B' }), {
    removable: false,
    via: '01M6NNAXQ276DMK3A1GH3XXF5B',
  });
  assert.deepEqual(epicCardExit({ direct: false }), { removable: false, via: '' });
});

test('a handle is the last six characters of the id', () => {
  assert.equal(handleOf('01M6NNAXQ276DMK3A1GH3XXF5B'), '3XXF5B');
});
