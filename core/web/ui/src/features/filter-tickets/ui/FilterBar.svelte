<script lang="ts">
  import { byBoardOrder, laneColour } from '@/entities/lane';
  import { WHERES, emptyFilter, narrows, type Filter } from '../model/filter';
  import ChipSelect from './ChipSelect.svelte';

  let { filter = $bindable(), lanes, tags, people }: {
    filter: Filter;
    lanes: string[];
    tags: string[];
    people: string[];
  } = $props();
</script>

<form class="filters" role="search" onsubmit={(e) => e.preventDefault()}>
  <label class="q">Search <input type="search" bind:value={filter.q} placeholder="title, goal or id" /></label>
  <ChipSelect label="Lane" any="any" options={byBoardOrder(lanes)} bind:chips={filter.lane} colour={laneColour} />
  <ChipSelect label="Tag" any="any" options={tags} bind:chips={filter.tag} />
  <ChipSelect label="Where" title="Board, logbook, archive (put away unfinished or finished) or rejected" any="board" options={WHERES} bind:chips={filter.where} />
  <ChipSelect label="Person" title="Created it, is assigned it, or changed it last" any="anyone" options={people} bind:chips={filter.person} />
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
  {#if narrows(filter)}
    <button type="button" onclick={() => (filter = emptyFilter())}>Clear</button>
  {/if}
</form>

<style>
  .filters { display: flex; flex-wrap: wrap; gap: 0.5rem 0.75rem; align-items: end; margin-top: 0.75rem; }
  label, .range { display: grid; gap: 0.15rem; font-size: 0.85rem; color: var(--muted); }
  .row { display: flex; gap: 0.3rem; align-items: center; }
  .q { flex: 1 1 14rem; }
  input, select, button { font: inherit; color: var(--text); background: var(--card); border: 1px solid var(--line); border-radius: 4px; padding: 0.25rem 0.4rem; }
  button { cursor: pointer; }
</style>
