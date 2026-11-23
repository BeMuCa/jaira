import { test } from 'node:test';
import assert from 'node:assert/strict';
import { handleOf } from './handle.ts';

test('a handle is the last six characters of the id', () => {
  assert.equal(handleOf('01M6NNAXQ276DMK3A1GH3XXF5B'), '3XXF5B');
});
