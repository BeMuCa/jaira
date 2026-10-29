<script lang="ts">
  // signoff: finished and waiting for a person to accept it, counted apart
  // from done so it reads neither as finished nor as untouched.
  let { done, signoff = 0, total }: { done: number; signoff?: number; total: number } = $props();
  const pct = (n: number) => (total === 0 ? 0 : Math.round((n / total) * 100));
  const label = $derived(signoff > 0 ? `${done} done · ${signoff} in signoff of ${total}` : `${done} of ${total} done`);
</script>

<div class="progress" title={label}>
  <div class="bar" style:width="{pct(done)}%"></div>
  <div class="bar signoff" style:left="{pct(done)}%" style:width="{pct(signoff)}%"></div>
  <span>{label}</span>
</div>

<style>
  .progress { position: relative; height: 1.4rem; border-radius: 4px; background: var(--track); overflow: hidden; }
  .bar { position: absolute; inset: 0 auto 0 0; background: var(--accent); opacity: 0.35; }
  .signoff { background: var(--waiting); }
  span { position: relative; padding: 0 0.5rem; font-size: 0.85rem; line-height: 1.4rem; }
</style>
