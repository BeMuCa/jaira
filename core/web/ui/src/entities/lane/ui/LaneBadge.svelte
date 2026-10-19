<script lang="ts">
  import { laneColour } from '../model/colour.svelte';

  let { lane }: { lane: string } = $props();
</script>

<span class="lane" style:color={laneColour(lane)}>{lane}</span>

<style>
  .lane {
    font-size: 0.8rem;
    padding: 0 0.4rem;
    border: 1px solid currentColor;
    border-radius: 4px;
    background: color-mix(in srgb, currentColor 12%, transparent);
    white-space: nowrap;
  }
</style>
