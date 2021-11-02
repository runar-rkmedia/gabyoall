<script type="ts">
  import { db } from '../api'
  import Stat from './Stat.svelte'
</script>

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
  }) as [key, stat]}
    <li>
      <Stat {stat} />
    </li>
  {/each}
</ol>
