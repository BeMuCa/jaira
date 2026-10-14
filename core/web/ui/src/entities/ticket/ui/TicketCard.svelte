<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { Ticket } from '../model/types';

  let { ticket, draggable = false, ondragstart, children }: {
    ticket: Ticket;
    draggable?: boolean;
    ondragstart?: (e: DragEvent) => void;
    children?: Snippet;
  } = $props();
</script>

<li class="card" class:done={ticket.done} class:indirect={!ticket.direct} {draggable} {ondragstart}>
  <span class="handle">{ticket.handle}</span>
  <span class="title">{ticket.title}</span>
  <span class="lane">{ticket.lane}</span>
  {@render children?.()}
</li>

<style>
  .card {
    display: grid;
    grid-template-columns: auto 1fr auto;
    gap: 0.5rem;
    align-items: baseline;
    padding: 0.4rem 0.6rem;
    border: 1px solid var(--line);
    border-radius: 6px;
    background: var(--card);
    list-style: none;
  }
  .card[draggable='true'] { cursor: grab; }
  .handle { font-family: ui-monospace, monospace; color: var(--muted); }
  .lane { font-size: 0.8rem; color: var(--muted); }
  .done .title { text-decoration: line-through; color: var(--muted); }
  .indirect { border-style: dashed; }
</style>
