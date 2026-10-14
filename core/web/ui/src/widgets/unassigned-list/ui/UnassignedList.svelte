<script lang="ts">
  import { TicketCard, type Ticket } from '@/entities/ticket';
  import { DropTarget, startTicketDrag } from '@/features/assign-to-epic';

  let { tickets, onremove, onopen }: {
    tickets: Ticket[];
    // A ticket dragged here out of an epic leaves that epic.
    onremove: (epic: string, id: string) => void;
    onopen: (id: string) => void;
  } = $props();
</script>

<DropTarget ondropticket={(t) => t.from && onremove(t.from, t.id)}>
  <section>
    <h2>Tickets in no epic <small>{tickets.length}</small></h2>
    {#if tickets.length === 0}
      <p class="empty">Every open ticket belongs to an epic.</p>
    {:else}
      <ul>
        {#each tickets as ticket (ticket.id)}
          <TicketCard {ticket} {onopen} draggable ondragstart={(e) => startTicketDrag(e, ticket.id)} />
        {/each}
      </ul>
    {/if}
  </section>
</DropTarget>

<style>
  section { min-height: 6rem; }
  ul { display: grid; gap: 0.4rem; padding: 0; margin: 0; }
  small { color: var(--muted); font-weight: normal; }
  .empty { color: var(--muted); }
</style>
