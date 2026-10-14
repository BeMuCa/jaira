<script lang="ts">
  import type { Snippet } from 'svelte';
  import { carriesTicket, droppedTicket, type DraggedTicket } from '../model/drag';

  let { ondropticket, children }: {
    ondropticket: (t: DraggedTicket) => void;
    children: Snippet;
  } = $props();

  let over = $state(false);

  function dragover(e: DragEvent) {
    if (!carriesTicket(e)) return;
    e.preventDefault();
    over = true;
  }

  function drop(e: DragEvent) {
    over = false;
    const t = droppedTicket(e);
    if (!t) return;
    e.preventDefault();
    ondropticket(t);
  }
</script>

<div class="target" class:over role="group" ondragover={dragover} ondragleave={() => (over = false)} ondrop={drop}>
  {@render children()}
</div>

<style>
  .target { border-radius: 8px; outline: 2px dashed transparent; outline-offset: 2px; }
  .over { outline-color: var(--accent); }
</style>
