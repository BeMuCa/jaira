<script lang="ts">
  import { EpicProgress, type Epic } from '@/entities/epic';
  import { TicketCard } from '@/entities/ticket';
  import { DropTarget } from '@/features/assign-to-epic';

  let { epics, onassign }: {
    epics: Epic[];
    onassign: (epic: string, id: string) => void;
  } = $props();
</script>

<section>
  <h2>Epics <small>{epics.length}</small></h2>
  {#if epics.length === 0}
    <p class="empty">No epics yet. Start one with <code>jaira epic create</code>.</p>
  {/if}
  {#each epics as epic (epic.name)}
    <DropTarget ondropticket={(id) => onassign(epic.name, id)}>
      <article>
        <header>
          <h3>{epic.title}</h3>
          <EpicProgress done={epic.done} total={epic.total} />
        </header>
        {#if epic.description}<p class="desc">{epic.description}</p>{/if}
        <details>
          <summary>{epic.tickets.length} tickets</summary>
          <ul>
            {#each epic.tickets as ticket (ticket.id)}
              <TicketCard {ticket} />
            {/each}
          </ul>
        </details>
      </article>
    </DropTarget>
  {/each}
</section>

<style>
  section { display: grid; gap: 0.75rem; align-content: start; }
  article { padding: 0.75rem; border: 1px solid var(--line); border-radius: 8px; background: var(--panel); }
  header { display: grid; gap: 0.4rem; }
  h3 { margin: 0; font-size: 1rem; }
  .desc { margin: 0.5rem 0; white-space: pre-line; }
  ul { display: grid; gap: 0.3rem; padding: 0; margin: 0.4rem 0 0; }
  summary { cursor: pointer; color: var(--muted); font-size: 0.85rem; }
  small { color: var(--muted); font-weight: normal; }
  .empty { color: var(--muted); }
</style>
