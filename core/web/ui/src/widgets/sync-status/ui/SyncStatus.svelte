<script lang="ts">
  import type { SyncState } from '@/entities/sync';
  import { clock } from '@/shared/lib';

  let { sync }: { sync: SyncState } = $props();
</script>

<div class="sync">
  <span>Last sync: {clock(sync.last_fetch)}</span>
  {#if sync.fetch_error}<span class="err" title={sync.fetch_error}>fetch failed</span>{/if}
  {#if sync.push_error}<span class="err" title={sync.push_error}>push failed — queued, retried on the next change</span>{/if}
</div>

<style>
  .sync { display: flex; gap: 0.75rem; align-items: center; font-size: 0.85rem; color: var(--muted); }
  .err { color: var(--error); }
</style>
