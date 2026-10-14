<script lang="ts">
  let { oncreate }: {
    // Resolves when the epic is written; throws with the server's message.
    oncreate: (name: string, title: string, description: string) => Promise<void>;
  } = $props();

  let open = $state(false);
  let name = $state('');
  let title = $state('');
  let description = $state('');
  let error = $state('');
  let busy = $state(false);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    try {
      await oncreate(name.trim(), title.trim(), description.trim());
      name = title = description = error = '';
      open = false;
    } catch (err) {
      error = (err as Error).message;
    } finally {
      busy = false;
    }
  }
</script>

{#if !open}
  <button class="new" onclick={() => (open = true)}>+ New epic</button>
{:else}
  <form onsubmit={submit}>
    <label>Name <input bind:value={name} required placeholder="pay-later" /></label>
    <label>Title <input bind:value={title} placeholder="what a person reads" /></label>
    <label>Description <textarea bind:value={description} rows="2" placeholder="what the feature is for"></textarea></label>
    {#if error}<p class="error" role="alert">{error}</p>{/if}
    <div class="actions">
      <button type="submit" disabled={busy}>Create</button>
      <button type="button" onclick={() => (open = false)}>Cancel</button>
    </div>
  </form>
{/if}

<style>
  form { display: grid; gap: 0.4rem; padding: 0.75rem; border: 1px solid var(--line); border-radius: 8px; background: var(--panel); }
  label { display: grid; gap: 0.15rem; font-size: 0.85rem; color: var(--muted); }
  input, textarea { font: inherit; color: inherit; background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.3rem; }
  .actions { display: flex; gap: 0.5rem; }
  button { font: inherit; color: inherit; background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.25rem 0.75rem; cursor: pointer; }
  .new { justify-self: start; }
  .error { margin: 0; color: var(--error); }
</style>
