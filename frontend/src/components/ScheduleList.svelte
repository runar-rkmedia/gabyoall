<script lang="ts">
  import { api, db } from '../api'
  import ScheduleItem from './items/ScheduleItem.svelte'
  import EntityList from './EntityList.svelte'
  export let selectedID: string = ''
  let loading = true
  let error: undefined | string = undefined
  $: deletedCount = Object.values($db.schedule).filter((s) => s.deleted).length
</script>

<EntityList {loading} {error} {deletedCount}>
  {#each Object.values($db.schedule)
    .sort((a, b) => {
      const A = a.start_date || ''
      const B = b.start_date || ''
      if (A > B) {
        return 1
      }
      if (A < B) {
        return -1
      }

      return 0
    })
    .reverse() as schedule}
    <ScheduleItem
      {selectedID}
      {schedule}
      onEdit={(id) => (selectedID = id)}
      onDelete={(id) => api.schedule.delete(id)} />
  {/each}
</EntityList>
