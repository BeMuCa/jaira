import { test } from 'node:test';
import assert from 'node:assert/strict';
import { ASSIGNEE_MAX, shortAssignee } from './assignee.ts';

test('a name that fits stays as it is', () => {
  assert.equal(shortAssignee('Alexander Sacharov'), 'Alexander Sacharov');
  assert.equal(shortAssignee('x'.repeat(ASSIGNEE_MAX)), 'x'.repeat(ASSIGNEE_MAX));
});

test('a longer name is cut to the limit, the ellipsis counted', () => {
  const short = shortAssignee('Maximiliane Müller-Lüdenscheidt');
  assert.equal(short, 'Maximiliane Müller-…');
  assert.equal(Array.from(short).length, ASSIGNEE_MAX);
});

test('characters outside the basic plane are not split', () => {
  const short = shortAssignee('🦊'.repeat(25));
  assert.equal(short, '🦊'.repeat(ASSIGNEE_MAX - 1) + '…');
});
