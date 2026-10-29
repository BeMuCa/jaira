import { test } from 'node:test';
import assert from 'node:assert/strict';
import { badgeLane } from './place.ts';

test('a filed ticket shows archive, every other its lane', () => {
  assert.equal(badgeLane({ lane: 'backlog', place: 'archive' }), 'archive');
  assert.equal(badgeLane({ lane: 'done', place: 'logbook' }), 'archive');
  assert.equal(badgeLane({ lane: 'signoff', place: 'board' }), 'signoff');
  assert.equal(badgeLane({ lane: 'review', place: 'ref' }), 'review');
  assert.equal(badgeLane({ lane: 'todo' }), 'todo');
});
