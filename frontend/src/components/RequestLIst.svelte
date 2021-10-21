<script lang="ts">
  import { api, db } from '../api'
  import { slide } from 'svelte/transition'
  import Spinner from './Spinner.svelte'
  import RequestItem from './items/RequestItem.svelte'
  let requests = api.request.list()
  let loading = true
  requests.then(() => (loading = false))
</script>

<div class="spinner"><Spinner active={loading} /></div>
{#await requests then [_, err]}
  {#if err}
    {err.error}
  {/if}
{/await}
<ul>
  {#each Object.values($db.request)
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
    <RequestItem request={v} />
  {/each}
</ul>

<style>
  .spinner {
    float: right;
  }
</style>
