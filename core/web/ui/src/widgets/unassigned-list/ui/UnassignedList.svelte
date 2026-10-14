<script lang="ts">
  import { TicketCard, type Ticket } from '@/entities/ticket';
  import { startTicketDrag } from '@/features/assign-to-epic';

  let { tickets }: { tickets: Ticket[] } = $props();
</script>

<section>
  <h2>Tickets in no epic <small>{tickets.length}</small></h2>
  {#if tickets.length === 0}
    <p class="empty">Every open ticket belongs to an epic.</p>
  {:else}
    <ul>
      {#each tickets as ticket (ticket.id)}
        <TicketCard {ticket} draggable ondragstart={(e) => startTicketDrag(e, ticket.id)} />
      {/each}
    </ul>
  {/if}
</section>

<style>
  ul { display: grid; gap: 0.4rem; padding: 0; margin: 0; }
  small { color: var(--muted); font-weight: normal; }
  .empty { color: var(--muted); }
</style>
