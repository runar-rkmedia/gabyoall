<script lang="ts">
  import { api, db } from '../api'
  import Icon from './Icon.svelte'

  import Tab from './Tab.svelte'
  export let value = 'request'
  const onClick = (s: keyof typeof api) => (e: MouseEvent) => (value = s)

  $: scheduleCount = Object.keys($db.schedule).length
  $: requestCount = Object.keys($db.request).length
  $: endpointCount = Object.keys($db.endpoint).length
  $: stat = Object.values($db.stat).find(
    (s) => s.total_requests && s.CompletedRequests! !== s.total_requests
  )
</script>

<nav>
  <Tab active={value === 'endpoint'} on:click={onClick('endpoint')}>
    <Icon icon={'gEndpoint'} color="tertiary" />
    Endpoints ({endpointCount})
  </Tab>
  <Tab active={value === 'request'} on:click={onClick('request')}>
    <Icon icon={'gRequest'} color="tertiary" />
    Requests ({requestCount})
  </Tab>
  <Tab
    active={value === 'schedule'}
    on:click={onClick('schedule')}
    disabled={!endpointCount && !requestCount}>
    <Icon icon={'gSchedule'} color="tertiary" />
    Schedule ({scheduleCount})
  </Tab>
  <Tab active={value === 'stats'} on:click={onClick('stat')}>
    <Icon icon={'gStat'} color="tertiary" />
    Stats
    {#if stat && stat.total_requests}
      <Icon icon="play" />
      <span class="running-stat">
        {(
          (100 * (stat.CompletedRequests || 1)) /
          (stat.total_requests || 1)
        ).toFixed(1)}% ({stat.CompletedRequests || 0}
        / {stat.total_requests || 0})
      </span>
    {/if}
  </Tab>
</nav>

<style>
  nav {
    display: flex;
    position: relative;
  }
  nav::after {
    content: '';
    width: 1px;

    position: absolute;
    left: 0;
    top: var(--size-6);
    bottom: var(--size-6);

    background-color: var(--color-red-500);
  }
</style>
