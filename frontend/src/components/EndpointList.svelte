<script lang="ts">
  import { api, db } from '../api'
  import { slide } from 'svelte/transition'
  import Spinner from './Spinner.svelte'
  import formatDate from 'dates'
  import ConfigItem from './items/ConfigItem.svelte'
  import EndpointItem from './items/EndpointItem.svelte'
  let endpoints = api.endpoint.list()
  let loading = true
  endpoints.then(() => (loading = false))
</script>

<div class="spinner"><Spinner active={loading} /></div>
{#await endpoints then [_, err]}
  {#if err}
    {err.error}
  {/if}
{/await}
<ul>
  {#each Object.values($db.endpoint)
    .sort((a, b) => {
      const A = a.createdAt
      const B = b.createdAt
      if (A > B) {
        return 1
      }
      if (A < B) {
        return -1
      }

      return 0
    })
    .reverse() as v}
    <EndpointItem endpoint={v} />
  {/each}
</ul>

<style>
  .spinner {
    float: right;
  }
</style>
