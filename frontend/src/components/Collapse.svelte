<script lang="ts">
  import { state } from '../state'

  import { onMount } from 'svelte'
  import Icon from './Icon.svelte'

  export let show = false
  export let key: string
  onMount(() => {
    if (!key) {
      return
    }
    show = $state.collapse[key]
  })
</script>

<div class="collapse">
  <button
    class="btn-reset toggle"
    aria-label="Collapse"
    on:click|preventDefault={() => {
      show = !show
      if (!key) {
        return
      }
      $state.collapse[key] = show
    }}>
    <slot name="title" class="title" />
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
