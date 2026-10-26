<script lang="ts">
  import { byBoardOrder, laneColour } from '@/entities/lane';
  import { emptyFilter, narrows, type Filter } from '../model/filter';

  let { filter = $bindable(), lanes, tags, people }: {
    filter: Filter;
    lanes: string[];
    tags: string[];
    people: string[];
  } = $props();
</script>

<form class="filters" role="search" onsubmit={(e) => e.preventDefault()}>
  <label class="q">Search <input type="search" bind:value={filter.q} placeholder="title, goal or id" /></label>
  <label>Lane
    <select bind:value={filter.lane} class="lane" style:border-left-color={filter.lane ? laneColour(filter.lane) : undefined}>
      <option value="">any</option>
      {#each byBoardOrder(lanes) as l (l)}<option value={l} style:color={laneColour(l)}>● {l}</option>{/each}
    </select>
  </label>
  <label>Tag
    <select bind:value={filter.tag}>
      <option value="">any</option>
      {#each tags as t (t)}<option value={t}>{t}</option>{/each}
    </select>
  </label>
  <label title="Created it, is assigned it, or changed it last">Person
    <select bind:value={filter.person}>
      <option value="">anyone</option>
      {#each people as p (p)}<option value={p}>{p}</option>{/each}
    </select>
  </label>
  <div class="range" role="group" aria-label="Date">
    <span>Date</span>
    <div class="row">
      <select bind:value={filter.date} aria-label="Which date">
        <option value="created">created</option>
        <option value="updated">last changed</option>
        <option value="filed">filed away</option>
      </select>
      <input type="date" bind:value={filter.from} aria-label="From" />
      <span aria-hidden="true">–</span>
      <input type="date" bind:value={filter.to} aria-label="To" />
    </div>
  </div>
  <label class="all"><input type="checkbox" bind:checked={filter.all} /> All tickets ever <small>logbook and archive too</small></label>
  {#if narrows(filter) || filter.all}
    <button type="button" onclick={() => (filter = emptyFilter())}>Clear</button>
  {/if}
</form>

<style>
  .filters { display: flex; flex-wrap: wrap; gap: 0.5rem 0.75rem; align-items: end; margin-top: 0.75rem; }
  label, .range { display: grid; gap: 0.15rem; font-size: 0.85rem; color: var(--muted); }
  .row { display: flex; gap: 0.3rem; align-items: center; }
  .q { flex: 1 1 14rem; }
  input, select, button { font: inherit; color: var(--text); background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.25rem 0.4rem; }
  .all { display: flex; gap: 0.35rem; align-items: center; color: var(--text); }
  .all small { color: var(--muted); }
  button { cursor: pointer; }
  .lane { border-left-width: 6px; }
</style>
