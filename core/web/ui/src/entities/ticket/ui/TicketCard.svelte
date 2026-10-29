<script lang="ts">
  import { age, date } from '@/shared/lib';
  import type { Ticket } from '../model/types';

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
  {#if ticket.created_at || ticket.updated_at}
    <span class="dates">
      {#if ticket.created_at}<time datetime={ticket.created_at}>created {date(ticket.created_at)}</time>{/if}
      {#if ticket.created_at && ticket.updated_at} · {/if}
      {#if ticket.updated_at}<time datetime={ticket.updated_at} title="Time since the ticket was last written">untouched {age(ticket.updated_at)}</time>{/if}
    </span>
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
  .dates, .via { grid-column: 2 / -1; font-size: 0.75rem; color: var(--muted); }
  .done .title { text-decoration: line-through; color: var(--muted); }
  .indirect { border-style: dashed; }
</style>
