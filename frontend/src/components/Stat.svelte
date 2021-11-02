<script type="ts">
  import formatDate from '../dates'

  export let stat: ApiDef.StatEntity
  function formatDuration(duration: number | undefined | null) {
    if (!duration) {
      return ''
    }
    const ms = Math.floor(duration / 1e6)

    if (ms > 1500) {
      const s = Math.floor(ms / 1e3)
      const msRest = ms - s * 1e3
      return `${s}s ${msRest}ms`
    }
    return ms + 'ms'
  }
</script>

<span class="count">
  min: {formatDuration(stat.Min)}
</span>
<span class="count">
  avg: {formatDuration(stat.Average)}
</span>
<span class="count">
  max: {formatDuration(stat.Max)}
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
{/if}
