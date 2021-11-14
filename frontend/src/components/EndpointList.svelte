<script lang="ts">
  import { api, db } from '../api'
  import EndpointItem from './items/EndpointItem.svelte'
  import { state } from '../state'
  import EntityList from './EntityList.svelte'
  export let selectedID = ''
  let endpoints = api.endpoint.list()
  let loading = true
  endpoints.then(() => (loading = false))
  $: deletedCount = Object.values($db.endpoint).filter((e) => e.deleted).length
  let error: string | undefined = undefined
  $: {
    endpoints.then((response) => (error = response[1]?.error))
  }
</script>

<EntityList {error} {loading} {deletedCount}>
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
    <EndpointItem
      endpoint={v}
      onEdit={(id) => (selectedID = id)}
      onDelete={(id) => api.endpoint.delete(id)} />
  {/each}
</EntityList>
