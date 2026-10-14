<script lang="ts">
  let { done, total }: { done: number; total: number } = $props();
  const pct = $derived(total === 0 ? 0 : Math.round((done / total) * 100));
</script>

<div class="progress" title="{done} of {total} done">
  <div class="bar" style:width="{pct}%"></div>
  <span>{done} of {total} done</span>
</div>

<style>
  .progress { position: relative; height: 1.4rem; border-radius: 4px; background: var(--track); overflow: hidden; }
  .bar { position: absolute; inset: 0 auto 0 0; background: var(--accent); opacity: 0.35; }
  span { position: relative; padding: 0 0.5rem; font-size: 0.85rem; line-height: 1.4rem; }
</style>
