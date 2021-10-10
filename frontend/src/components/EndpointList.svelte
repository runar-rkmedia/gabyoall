<script lang="ts">
  import { api, db } from '../api'
  import { slide } from 'svelte/transition'
  import Spinner from './Spinner.svelte'
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
  {#each Object.entries($db.endpoint)
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
      {v.url}
      {new Date(v.createdAt).toLocaleTimeString()}
      {#if v.headers}
        <strong>Headers</strong>
        <ul>
          {#each Object.entries(v.headers) as [hKey, hVal]}
            <li>
              {hKey} - {hVal.join('; ')}
            </li>
          {/each}
        </ul>
      {/if}
    </li>
  {/each}
</ul>

<style>
  .spinner {
    float: right;
  }
</style>
