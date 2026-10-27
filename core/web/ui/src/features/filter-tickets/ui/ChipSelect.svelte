<script lang="ts">
  import { addChip, removeChip, toggleChip, type Chip } from '../model/filter';

  // A dropdown that adds what is picked as a chip beside it. A click on the
  // chip turns 'X' into 'not X' and back; × takes it away.
  let { label, title, any, options, chips = $bindable(), colour }: {
    label: string;
    title?: string;
    any: string;
    options: string[];
    chips: Chip[];
    colour?: (value: string) => string;
  } = $props();

  const left = $derived(options.filter((o) => !chips.some((c) => c.value === o)));

  function pick(e: Event) {
    const select = e.currentTarget as HTMLSelectElement;
    if (select.value) chips = addChip(chips, select.value);
    select.value = '';
  }
</script>

<div class="chips" role="group" aria-label={label} {title}>
  <span>{label}</span>
  <div class="row">
    <select onchange={pick} aria-label="Add a {label.toLowerCase()}">
      <option value="">{chips.length ? '+ add' : any}</option>
      {#each left as o (o)}<option value={o} style:color={colour?.(o)}>{colour ? '● ' : ''}{o}</option>{/each}
    </select>
    {#each chips as c (c.value)}
      <span class="chip" class:not={c.not} style:border-left-color={colour?.(c.value)}>
        <button type="button" class="toggle" title={c.not ? 'Click to keep only this' : 'Click to leave this out'} onclick={() => (chips = toggleChip(chips, c.value))}>{c.not ? `not ${c.value}` : c.value}</button>
        <button type="button" class="remove" aria-label="Remove {c.value}" onclick={() => (chips = removeChip(chips, c.value))}>×</button>
      </span>
    {/each}
  </div>
</div>

<style>
  .chips { display: grid; gap: 0.15rem; font-size: 0.85rem; color: var(--muted); }
  .row { display: flex; flex-wrap: wrap; gap: 0.3rem; align-items: center; }
  .chip { display: inline-flex; align-items: center; border: 1px solid var(--line); border-left-width: 4px; border-radius: 4px; background: var(--card); color: var(--text); }
  .chip button { all: unset; cursor: pointer; padding: 0.2rem 0.35rem; }
  .chip button:focus-visible { outline: 2px solid var(--line); }
  select { font: inherit; color: var(--text); background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.25rem 0.4rem; }
  .remove { color: var(--muted); }
  .not { border-style: dashed; }
  .not .toggle { color: var(--muted); }
</style>
