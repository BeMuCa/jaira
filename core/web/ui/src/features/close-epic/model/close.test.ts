import { test } from 'node:test';
import assert from 'node:assert/strict';
import { addQuestion, closeQuestion, visible } from './close.ts';

test('closing asks only while tickets are unfinished, and says how many', () => {
  assert.equal(closeQuestion({ title: 'Login', done: 3, total: 3 }), null);
  assert.equal(closeQuestion({ title: 'Login', done: 0, total: 0 }), null);
  assert.match(closeQuestion({ title: 'Login', done: 1, total: 3 }) ?? '', /2 of 3/);
});

test('a drop asks only on a closed epic, and says it stays closed', () => {
  assert.equal(addQuestion({ title: 'Login', done: 0, total: 1 }), null);
  assert.match(addQuestion({ title: 'Login', done: 1, total: 1, closed: true }) ?? '', /stays closed/);
});

test('closed epics are hidden unless asked for', () => {
  const epics = [{ name: 'a' }, { name: 'b', closed: true }];
  assert.deepEqual(visible(epics, false).map((e) => e.name), ['a']);
  assert.deepEqual(visible(epics, true).map((e) => e.name), ['a', 'b']);
});
