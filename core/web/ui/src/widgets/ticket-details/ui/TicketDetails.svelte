<script lang="ts">
  import { NoteItem, behindRef, type TicketDetail } from '@/entities/ticket';

  let { ticket, onclose }: { ticket: TicketDetail; onclose: () => void } = $props();

  // A ticket worked over several rounds carries dozens of notes; the newest
  // ones say where it stands, so the older ones fold away until asked for.
  const SHOWN = 5;
  let all = $state(false);
  const hidden = $derived(all ? 0 : Math.max(0, ticket.notes.length - SHOWN));
  const notes = $derived(ticket.notes.slice(hidden));
</script>

<article aria-labelledby="ticket-title">
  <header>
    <span class="handle">{ticket.handle}</span>
    <span class="lane">{ticket.lane}</span>
    <button class="close" onclick={onclose} aria-label="Close">×</button>
  </header>
  <h2 id="ticket-title">{ticket.title}</h2>
  {#if ticket.behind_ref}<p class="behind">{behindRef(ticket.handle)}</p>{/if}
  {#if ticket.goal}<h3>Goal</h3><p class="text">{ticket.goal}</p>{/if}
  {#if ticket.context}<h3>Context</h3><p class="text">{ticket.context}</p>{/if}
  {#if ticket.dod.length}
    <h3>Definition of done</h3>
    <ul class="dod">
      {#each ticket.dod as item, i (i)}
        <li class:done={item.done}><span aria-hidden="true">{item.done ? '☑' : '☐'}</span> {item.text}</li>
      {/each}
    </ul>
  {/if}
  <h3>Notes <small>{ticket.notes.length}</small></h3>
  {#if ticket.notes.length === 0}
    <p class="muted">No notes yet.</p>
  {:else}
    {#if hidden > 0}
      <button class="more" onclick={() => (all = true)}>Show all {ticket.notes.length} notes ({hidden} older)</button>
    {/if}
    <ol>
      {#each notes as note, i (hidden + i)}
        <NoteItem {note} />
      {/each}
    </ol>
  {/if}
</article>

<style>
  .behind { font-size: 0.85rem; color: var(--waiting); }
  article { display: grid; gap: 0.25rem; }
  header { display: flex; gap: 0.75rem; align-items: baseline; }
  .handle { font-family: ui-monospace, monospace; color: var(--muted); }
  .lane { font-size: 0.85rem; padding: 0 0.4rem; border: 1px solid var(--line); border-radius: 4px; }
  .close { margin-left: auto; font-size: 1.25rem; line-height: 1; background: none; border: none; color: inherit; cursor: pointer; }
  h2 { margin: 0.25rem 0; font-size: 1.15rem; }
  h3 { margin: 0.75rem 0 0.1rem; font-size: 0.85rem; text-transform: uppercase; color: var(--muted); }
  .text { margin: 0; white-space: pre-line; }
  .dod { padding: 0; margin: 0; display: grid; gap: 0.2rem; }
  .dod li { list-style: none; }
  .dod .done { color: var(--muted); }
  ol { padding: 0; margin: 0; }
  small { font-weight: normal; }
  .muted { color: var(--muted); }
  .more { justify-self: start; background: none; border: 1px solid var(--line); border-radius: 4px; padding: 0.2rem 0.6rem; color: inherit; cursor: pointer; }
</style>
