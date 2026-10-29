import { test } from 'node:test';
import assert from 'node:assert/strict';
import { age } from './time.ts';

test('age counts days from a day on, hours below, and <1h under an hour', () => {
  const now = Date.parse('2026-10-29T20:00:00Z');
  assert.equal(age('2026-10-29T19:30:00Z', now), '<1h');
  assert.equal(age('2026-10-29T15:00:00Z', now), '5h');
  assert.equal(age('2026-10-28T20:00:00Z', now), '1d');
  assert.equal(age('2026-10-26T08:00:00Z', now), '3d');
  assert.equal(age(undefined, now), '');
});
