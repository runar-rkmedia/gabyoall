<script type="ts">
  import { state } from 'state'

  import Button from '../Button.svelte'
  export let deletedCount: number
</script>

{#if !!deletedCount}
  <Button
    icon="delete"
    color="secondary"
    active={$state.showDeleted}
    disabled={!$state.showDeleted && !deletedCount}
    on:click={() => ($state.showDeleted = !$state.showDeleted)}>
    Show deleted ({deletedCount})
  </Button>
{/if}
