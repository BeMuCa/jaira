<script lang="ts">
  import type { SyncState } from '@/entities/sync';
  import { clock } from '@/shared/lib';

  let { sync, onsync }: { sync: SyncState; onsync: () => void } = $props();
</script>

<div class="sync">
  <span>Last sync: {clock(sync.last_fetch)}</span>
  {#if sync.fetch_error}<span class="err" title={sync.fetch_error}>fetch failed</span>{/if}
  {#if sync.push_error}<span class="err" title={sync.push_error}>push failed — queued, retried on the next change</span>{/if}
  <button type="button" onclick={onsync}>Sync now</button>
</div>

<style>
  .sync { display: flex; gap: 0.75rem; align-items: center; font-size: 0.85rem; color: var(--muted); }
  .err { color: var(--error); }
  button { font: inherit; padding: 0.2rem 0.6rem; border: 1px solid var(--line); border-radius: 4px; background: var(--card); color: var(--text); cursor: pointer; }
</style>
