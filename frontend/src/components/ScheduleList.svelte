<script lang="ts">
  import { db } from '../api'
  import Spinner from './Spinner.svelte'
  import ScheduleItem from './items/ScheduleItem.svelte'
  export let selectedID: string = ''
  let loading = true
</script>

<div class="spinner"><Spinner active={loading} /></div>
<ul>
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
    <ScheduleItem {schedule} on:click={() => (selectedID = schedule.id)} />
  {/each}
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
