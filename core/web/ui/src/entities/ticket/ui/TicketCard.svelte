<script lang="ts">
  import type { Ticket } from '../model/types';
  import { PLACE_LABELS } from '../model/place';

  let { ticket, draggable = false, ondragstart, onopen }: {
    ticket: Ticket;
    draggable?: boolean;
    ondragstart?: (e: DragEvent) => void;
    onopen?: (id: string) => void;
  } = $props();
</script>

<li
  class="card"
  class:done={ticket.done}
  class:indirect={!ticket.direct}
  {draggable}
  {ondragstart}
  title={ticket.direct ? undefined : 'Counted through follows or parent: it has no line in the epic, so it cannot be dragged out'}
>
  <span class="handle">{ticket.handle}</span>
  <button class="title" onclick={() => onopen?.(ticket.id)}>{ticket.title}</button>
  <span class="lane">{ticket.lane}</span>
  {#if ticket.place && ticket.place !== 'board'}
    <span class="where">{PLACE_LABELS[ticket.place]}{#if ticket.rejected} · rejected: {ticket.rejected}{/if}</span>
  {/if}
  {#if !ticket.direct}<span class="via">via follows/parent</span>{/if}
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
  .title { all: unset; cursor: pointer; overflow-wrap: anywhere; }
  .title:hover, .title:focus-visible { text-decoration: underline; }
  .lane { font-size: 0.8rem; color: var(--muted); }
  .via, .where { grid-column: 2 / -1; font-size: 0.75rem; color: var(--muted); }
  .done .title { text-decoration: line-through; color: var(--muted); }
  .indirect { border-style: dashed; }
</style>
