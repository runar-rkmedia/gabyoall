<script type="ts">
  import { api, db } from '../api'
  import Button from './Button.svelte'
  import Stat from './Stat.svelte'
</script>

<Button icon="delete" color="danger" on:click={() => api.stat.clean()}
  >Click to remove Stats
</Button>
<Button icon="refresh" color="primary" on:click={() => api.stat.list()}>
  Refresh
</Button>
<ol>
  {#each Object.entries($db.stat).sort(([_, a], [__, b]) => {
    const A = a.StartTime
    const B = b.StartTime
    if (A > B) {
      return -1
    }
    if (A < B) {
      return 1
    }

    return 0
  }) as [key, stat], i}
    <li>
      {#key stat.updatedAt}
        <Stat stat={$db.stat[key]} showChart={!i} />
      {/key}
    </li>
  {/each}
</ol>
