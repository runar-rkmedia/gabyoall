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
    <div class="icon">
      {#if show}
        <Icon icon={'collapseUp'} class="toggle-icon" />
      {:else}
        <Icon icon={'collapseDown'} class="toggle-icon" />
      {/if}
    </div>
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

  .icon {
    font-size: 1.4rem;
  }
</style>
