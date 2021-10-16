<script type="ts">
  import formatDate from 'dates'

  import { LinkedChart } from 'svelte-tiny-linked-charts'

  export let stat: ApiDef.CompactRequestStatisticsEntity
</script>

<span class="count">
  min: {stat.MinText}
</span>
<span class="count">
  avg: {stat.AverageText}
</span>
<span class="count">
  max: {stat.MaxText}
</span>
<span class="count">
  count: {Object.keys(stat.Requests || {}).length}
</span>
{formatDate(stat.StartTime)}
{#if stat.Requests}
  <ul>
    {#each Object.entries(stat.Requests).filter(([_, s]) => s.error) as [_, req]}
      <li>
        {req.error}
      </li>
    {/each}
  </ul>
  <LinkedChart
    showValue={true}
    width="500px"
    grow={true}
    barMinWidth="1"
    data={Object.values(stat.Requests).reduce((r, req, i) => {
      if (!req.offset) {
        return r
      }
      r[req.offset] = req.duration
      return r
    }, {})} />
{/if}
