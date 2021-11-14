<script lang="ts">
  import { api, db } from '../api'
  import EntityList from './EntityList.svelte'
  import RequestItem from './items/RequestItem.svelte'

  let loading = true
  let error: undefined | string = undefined
  export let selectedID: string = ''

  let requests = api.request.list()
  requests.then((response) => {
    loading = false
    error = response[1]?.error
  })

  $: deletedCount = Object.values($db.request).filter((s) => s.deleted).length
</script>

<EntityList {loading} {error} {deletedCount}>
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
    <RequestItem
      {selectedID}
      onEdit={(id) => (selectedID = id)}
      onDelete={(id) => api.request.delete(id)}
      request={v} />
  {/each}
</EntityList>
