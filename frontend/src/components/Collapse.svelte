<script lang="ts">
  import Icon from './Icon.svelte'

  export let show = false
</script>

<div class="collapse">
  <button
    class="btn-reset toggle"
    aria-label="Collapse"
    on:click|preventDefault={() => (show = !show)}>
    <slot name="title" />
    {#if show}
      <Icon icon={'collapseUp'} />
    {:else}
      <Icon icon={'collapseDown'} />
    {/if}
  </button>
  {#if show}
    <slot />
  {/if}
</div>

<style>
  button.toggle {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }
  button.toggle .fas {
    font-size: 1.4rem;
    /* position: absolute;
    right: var(--size-4);
    /* top: var(--size-2); */
  }
</style>
