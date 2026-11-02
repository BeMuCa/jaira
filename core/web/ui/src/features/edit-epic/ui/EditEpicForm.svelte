<script lang="ts">
  import { untrack } from 'svelte';
  import type { Epic } from '@/entities/epic';
  import type { EpicEdit, EpicFields } from '../api/set';
  import { editOf, keysEdit, linkable, rebase } from '../model/edit';

  let { epic, onsave, oncancel }: {
    epic: Epic;
    // Resolves when the epic is written; throws with the server's message,
    // after the page has reloaded the epics.
    onsave: (edit: EpicEdit) => Promise<void>;
    oncancel: () => void;
  } = $props();

  const fieldsOf = (e: Epic): EpicFields => ({ title: e.title, description: e.description, body: e.body });

  // What the editor opened on. The page reloads the epics every minute; the
  // form keeps this, so a save is checked against what the person read.
  const opened = untrack(() => fieldsOf(epic));
  let base = $state(opened);
  let title = $state(opened.title);
  let description = $state(opened.description);
  let body = $state(opened.body);
  // The keys as the form opened on them, and one row per key to edit, a blank
  // one at the end to add a key.
  const keysRead = untrack(() => epic.external ?? []);
  let keys = $state([...keysRead.map((x) => ({ key: x.key, url: x.url })), { key: '', url: '' }]);
  let error = $state('');
  let busy = $state(false);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    const k = keysEdit(keysRead, keys);
    const fields = editOf(base, { title, description, body });
    if (!fields && !k.add_external && !k.drop_external) return oncancel();
    const edit = { ...(fields ?? { base }), ...k };
    busy = true;
    try {
      await onsave(edit);
    } catch (err) {
      error = (err as Error).message;
      // The page has reloaded by now. When somebody else changed the epic,
      // the form takes that state as its base and keeps what was typed: a
      // second Save is then a decision made knowing it, and goes through.
      const r = rebase(base, { title, description, body }, fieldsOf(epic));
      if (r.changed.length) {
        base = r.base;
        ({ title, description, body } = r.mine);
        const now = r.changed.map((f) => (f === 'body' ? 'the text' : `the ${f} to "${r.base[f]}"`));
        error = `Somebody changed ${now.join(' and ')} meanwhile. What you typed is kept; Save again to write it over that.`;
      }
    } finally {
      busy = false;
    }
  }
</script>

<form onsubmit={submit}>
  <label>Title <input bind:value={title} required /></label>
  <label>Description <input bind:value={description} placeholder="one line on what the feature is for" /></label>
  <label>Text <textarea bind:value={body} rows="12" placeholder="markdown"></textarea></label>
  <fieldset class="keys">
    <legend title="The key of the task in another system this epic stands for, with its link. It stays on the epic, never on a ticket.">External keys</legend>
    {#each keys as row, i (i)}
      <!-- The link is plain text, not type="url": the browser would refuse the
           whole form over one hand-written line. A bad link marks its row
           and is left out of the save; everything else is written. -->
      <div class="key" class:bad={!linkable(row.url)}>
        <input bind:value={row.key} placeholder="JIRA-1234" aria-label="Key" size="12" />
        <input bind:value={row.url} placeholder="https://…" aria-label="Link" aria-invalid={!linkable(row.url)} />
        {#if i < keys.length - 1}
          <button type="button" aria-label="Remove {row.key}" onclick={() => keys.splice(i, 1)}>×</button>
        {:else}
          <button type="button" onclick={() => keys.push({ key: '', url: '' })}>+</button>
        {/if}
      </div>
      {#if !linkable(row.url)}<p class="bad-hint">Not an http(s) link — this line is not saved; the rest of the form is.</p>{/if}
    {/each}
  </fieldset>
  {#if error}<p class="error" role="alert">{error}</p>{/if}
  <div class="actions">
    <button type="submit" disabled={busy}>Save</button>
    <button type="button" onclick={oncancel}>Cancel</button>
  </div>
</form>

<style>
  form { display: grid; gap: 0.4rem; margin-top: 0.5rem; }
  label { display: grid; gap: 0.15rem; font-size: 0.85rem; color: var(--muted); }
  input, textarea { font: inherit; color: inherit; background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.3rem; }
  textarea { font-family: ui-monospace, monospace; font-size: 0.85rem; resize: vertical; }
  .actions { display: flex; gap: 0.5rem; }
  button { font: inherit; color: inherit; background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.25rem 0.75rem; cursor: pointer; }
  .error { margin: 0; color: var(--error); }
  .keys { display: grid; gap: 0.3rem; border: none; padding: 0; margin: 0; font-size: 0.85rem; color: var(--muted); }
  .keys legend { padding: 0; margin-bottom: 0.15rem; }
  .key { display: flex; gap: 0.3rem; }
  .key input + input { flex: 1; min-width: 0; }
  .bad input + input { border-color: var(--error); }
  .bad-hint { margin: 0; color: var(--error); font-size: 0.8rem; }
  .key button { padding: 0.1rem 0.5rem; }
</style>
