<script lang="ts">
  import { api, db } from '../api'
  import { slide } from 'svelte/transition'
  import Spinner from './Spinner.svelte'
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
  {#each Object.entries($db.request)
    .sort((a, b) => {
      const A = a[1].createdAt
      const B = b[1].createdAt
      if (A > B) {
        return 1
      }
      if (A < B) {
        return -1
      }

      return 0
    })
    .reverse() as [k, v]}
    <li id={k} transition:slide|local>
      {v.method}
      {v.operationName}
      {new Date(v.createdAt).toLocaleTimeString()}
      {v.query}
      {v.variables}
      {v.body}
    </li>
  {/each}
</ul>

<style>
  .spinner {
    float: right;
  }
</style>
