<script lang="ts">
  import { age, date } from '@/shared/lib';
  import { behindRef } from '../model/behind';
  import type { Ticket } from '../model/types';
  import { PLACE_LABELS, badgeLane } from '../model/place';
  import { shortAssignee } from '../model/assignee';
  import { handleOf } from '../model/epic-card';
  import { LaneBadge } from '@/entities/lane';

  let { ticket, draggable = false, ondragstart, onopen, onremove, removeLabel = 'Remove' }: {
    ticket: Ticket;
    draggable?: boolean;
    ondragstart?: (e: DragEvent) => void;
    onopen?: (id: string) => void;
    // Set, the card carries a button that takes the ticket out of where it
    // is listed: a tap does what a drag does, where there is no drag.
    onremove?: () => void;
    // What the button says to a screen reader.
    removeLabel?: string;
  } = $props();
</script>

<li
  class="card"
  class:done={ticket.done}
  class:indirect={!ticket.direct}
  {draggable}
  {ondragstart}
>
  <span class="handle">{ticket.handle}</span>
  <button class="title" onclick={() => onopen?.(ticket.id)}>{ticket.title}</button>
  <LaneBadge lane={badgeLane(ticket)} />
  {#if onremove}
    <button class="remove" aria-label={removeLabel} title={removeLabel} onclick={onremove}>Remove</button>
  {/if}
  {#if ticket.assignee}
    <span class="assignee" title={ticket.assignee}>{shortAssignee(ticket.assignee)}</span>
  {/if}
  {#if ticket.place && ticket.place !== 'board'}
    <span class="where">{PLACE_LABELS[ticket.place]}{#if ticket.rejected} · rejected: {ticket.rejected}{/if}</span>
  {/if}
  {#if ticket.created_at || ticket.updated_at || ticket.behind_ref}
    <span class="dates">
      {#if ticket.created_at}<time datetime={ticket.created_at}>created {date(ticket.created_at)}</time>{/if}
      {#if ticket.created_at && ticket.updated_at} · {/if}
      {#if ticket.updated_at}<time datetime={ticket.updated_at} title="Time since the ticket was last written">untouched {age(ticket.updated_at)}</time>{/if}
      {#if ticket.behind_ref} · <span class="behind" title={behindRef(ticket.handle)}>file behind ref</span>{/if}
    </span>
  {/if}
  {#if !ticket.direct}
    <!-- It has no line in the epic to take out: the member it hangs off has. -->
    <span class="via">
      {#if ticket.via}
        counted through
        <button class="link" title="Open the member whose line brings this ticket in" onclick={() => onopen?.(ticket.via!)}>{handleOf(ticket.via)}</button>
        (follows/parent) — take that one out to remove this
      {:else}
        via follows/parent
      {/if}
    </span>
  {/if}
</li>

<style>
  .card {
    display: grid;
    grid-template-columns: auto 1fr auto auto;
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
  .via, .where, .assignee, .dates { grid-column: 2 / -1; font-size: 0.75rem; color: var(--muted); }
  .behind { color: var(--waiting); }
  .done .title { text-decoration: line-through; color: var(--muted); }
  .indirect { border-style: dashed; }
  .remove { font: inherit; font-size: 0.75rem; color: var(--muted); background: none; border: 1px solid var(--line); border-radius: 4px; padding: 0.2rem 0.5rem; min-height: 1.75rem; cursor: pointer; }
  .remove:hover, .remove:focus-visible { color: inherit; border-color: currentColor; }
  .link { all: unset; cursor: pointer; font-family: ui-monospace, monospace; text-decoration: underline; }
</style>
