import { test } from 'node:test';
import assert from 'node:assert/strict';
import { badgeLane, fate } from './place.ts';

test('an archived ticket shows the lane it came from, a logbook one archive', () => {
  assert.equal(badgeLane({ lane: 'backlog', place: 'archive' }), 'backlog');
  assert.equal(badgeLane({ lane: 'done', place: 'logbook' }), 'archive');
  assert.equal(badgeLane({ lane: 'signoff', place: 'board' }), 'signoff');
  assert.equal(badgeLane({ lane: 'review', place: 'ref' }), 'review');
  assert.equal(badgeLane({ lane: 'todo' }), 'todo');
});

test('an archived ticket is finished, left behind or rejected, and each differs from an open one', () => {
  assert.equal(fate({ done: true, place: 'archive' }), 'done', 'archived from done');
  assert.equal(fate({ done: true, place: 'logbook' }), 'done');
  assert.equal(fate({ done: false, place: 'archive' }), 'left', 'archived from another lane');
  assert.equal(fate({ done: false, place: 'archive', rejected: 'wontfix' }), 'rejected');
  assert.equal(fate({ done: true, place: 'archive', rejected: 'duplicate' }), 'rejected', 'a rejection outranks the lane');
  assert.equal(fate({ done: false, place: 'board' }), 'open');
  assert.equal(fate({ done: false, place: 'ref' }), 'open');
  assert.equal(fate({ done: true, place: 'board' }), 'done');
});
