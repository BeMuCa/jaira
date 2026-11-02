<script lang="ts">
  import { age, date } from '@/shared/lib';
  import { behindRef } from '../model/behind';
  import type { Ticket } from '../model/types';
  import { PLACE_LABELS, badgeLane, fate } from '../model/place';
  import { shortAssignee } from '../model/assignee';
  import { LaneBadge } from '@/entities/lane';

  let { ticket, draggable = false, ondragstart, onopen }: {
    ticket: Ticket;
    draggable?: boolean;
    ondragstart?: (e: DragEvent) => void;
    onopen?: (id: string) => void;
  } = $props();

  const how = $derived(fate(ticket));
</script>

<li
  class="card"
  class:done={how === 'done'}
  class:left={how === 'left'}
  class:rejected={how === 'rejected'}
  class:indirect={!ticket.direct}
  {draggable}
  {ondragstart}
  title={ticket.direct ? undefined : 'Counted through follows or parent: it has no line in the epic, so it cannot be dragged out'}
>
  <span class="handle">{ticket.handle}</span>
  <button class="title" onclick={() => onopen?.(ticket.id)}>{ticket.title}</button>
  <LaneBadge lane={badgeLane(ticket)} />
  {#if ticket.assignee}
    <span class="assignee" title={ticket.assignee}>{shortAssignee(ticket.assignee)}</span>
  {/if}
  {#if ticket.place && ticket.place !== 'board'}
    <span class="where">{PLACE_LABELS[ticket.place]}{#if how === 'left'} · left unfinished{/if}</span>
  {/if}
  {#if how === 'rejected'}
    <span class="why">rejected: {ticket.rejected}{#if ticket.rejected_why} – {ticket.rejected_why}{/if}</span>
  {/if}
  {#if ticket.created_at || ticket.updated_at || ticket.behind_ref}
    <span class="dates">
      {#if ticket.created_at}<time datetime={ticket.created_at}>created {date(ticket.created_at)}</time>{/if}
      {#if ticket.created_at && ticket.updated_at} · {/if}
      {#if ticket.updated_at}<time datetime={ticket.updated_at} title="Time since the ticket was last written">untouched {age(ticket.updated_at)}</time>{/if}
      {#if ticket.behind_ref} · <span class="behind" title={behindRef(ticket.handle)}>file behind ref</span>{/if}
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
  .via, .where, .why, .assignee, .dates { grid-column: 2 / -1; font-size: 0.75rem; color: var(--muted); }
  .behind { color: var(--waiting); }
  .done .title { text-decoration: line-through; color: var(--muted); }
  .left { opacity: 0.6; }
  .rejected { border-left: 3px solid var(--error); }
  .why { color: var(--error); overflow-wrap: anywhere; }
  .indirect { border-style: dashed; }
</style>
