<script lang="ts">
  import { api, db } from 'api'

  import Button from 'components/Button.svelte'

  import formatDate from 'dates'

  import { slide } from 'svelte/transition'
  import ConfigItem from './ConfigItem.svelte'

  export let endpoint: ApiDef.EndpointEntity
  $: {
    console.log($db.endpoint[endpoint.id].deleted)
  }
</script>

<li transition:slide|local class:deleted={!!endpoint.deleted}>
  {endpoint.url}
  {formatDate(endpoint.createdAt)}
  {#if endpoint.config}
    <ConfigItem config={endpoint.config} short={true} />
  {/if}
  {#if endpoint.headers}
    <strong>Headers</strong>
    <ul>
      {#each Object.entries(endpoint.headers) as [hKey, hVal]}
        <li>
          {hKey} - {hVal.join('; ')}
        </li>
      {/each}
    </ul>
  {/if}
  <Button icon="delete" on:click={() => api.endpoint.delete(endpoint.id)}>
    Delete
  </Button>
</li>

<style>
  .deleted {
    text-decoration: line-through;
  }
</style>
