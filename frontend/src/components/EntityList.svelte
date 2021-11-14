<script lang="ts">
  import Spinner from './Spinner.svelte'
  import ButtonShowDeleted from './buttons/ButtonShowDeleted.svelte'
  export let loading: boolean
  export let deletedCount: number
  export let error: string | undefined
</script>

<div class="spinner"><Spinner active={loading} /></div>
{#if error}
  {error}
{/if}
<ButtonShowDeleted {deletedCount} />
<ul>
  <slot />
</ul>

<style>
  .spinner {
    float: right;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    border-radius: var(--radius);
    box-shadow: var(--elevation-4);
  }
</style>
