<script type="ts">
  import formatDate from '../dates'
  import Button from './Button.svelte'
  import type ChartType from './Chart.svelte'
  let Chart: ChartType
  import { onMount } from 'svelte'
  onMount(async () => {
    const { default: def } = await import('./Chart.svelte')
    Chart = def as any
  })

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
  export let showChart = true
</script>

<div class="wrapper">
  <div class="header">
    <div class="title">
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
      {#if stat.total_requests && stat.CompletedRequests && stat.total_requests !== stat.CompletedRequests}
        <span class="count">
          Completed: {(
            (100 * stat.CompletedRequests) /
            (stat.total_requests || 0)
          ).toFixed(1)}% ({stat.CompletedRequests}
          / {stat.total_requests})
        </span>
      {/if}
      {formatDate(stat.StartTime)}
    </div>
    <div class="actions">
      <Button color="secondary" on:click={() => (showChart = !showChart)}
        >Toggle Chart</Button>
    </div>
  </div>
  {#if stat.TimeSeries && showChart}
    <div class="chart-container">
      {#if Chart}
        <svelte:component
          this={Chart}
          data={Object.entries(stat.TimeSeries).map(([k, s], i) => {
            const count = s.Series.length
            return s.Series.reduce(
              (r, [x, y]) => {
                r.x.push(x + y)
                r.y.push(y)
                /* r.marker.color.push(y) */

                return r
              },
              {
                x: [],
                y: [],
                type: 'scatter',
                name: `${k || 'ok'} (${count})`,
                mode: 'markers',
                marker: { size: 5, color: i },
              }
            )
          })}
          layout={{
            title: 'Request-duration.',
            yaxis: {
              type: 'log',
              autorange: true,
            },
          }} />
      {/if}
    </div>
  {/if}
</div>

<style>
  .header {
    display: flex;
    justify-content: space-between;
  }
  .chart-container {
    border: 1px solid red;
  }
</style>
