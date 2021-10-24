<script lang="ts">
  import { api, db } from '../api'
  import Spinner from './Spinner.svelte'
  import EndpointItem from './items/EndpointItem.svelte'
  import Button from './Button.svelte'
  import Alert from './Alert.svelte'
  import { state } from 'state'
  export let selectedID = ''
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
<Button
  icon="delete"
  on:click={() => ($state.showDeleted = !$state.showDeleted)}>
  {#if $state.showDeleted}
    Hide deleted
  {:else}
    Show deleted
  {/if}
</Button>
<ul>
  {#each Object.values($db.endpoint)
    .filter((e) => {
      if (!$state.showDeleted) {
        return !e.deleted
      }
      return true
    })
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
    {#if v.deleted}
      <Alert kind="warning">
        <EndpointItem endpoint={v} />
      </Alert>
    {:else}
      <EndpointItem endpoint={v} />
      <div class="item-actions">
        <Button icon="edit" on:click={() => (selectedID = v.id)}>Edit</Button>
      </div>
    {/if}
  {/each}
</ul>

<style>
  .spinner {
    float: right;
  }
</style>
