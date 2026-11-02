import { test } from 'node:test';
import assert from 'node:assert/strict';
import { nameBeside } from './label.ts';

test('the name goes beside the title only when it says something else', () => {
  assert.equal(nameBeside({ name: 'web-0-2', title: 'Web 0.2' }), 'web-0-2');
  assert.equal(nameBeside({ name: 'login', title: 'login' }), '');
  assert.equal(nameBeside({ name: 'login', title: '' }), '');
});
