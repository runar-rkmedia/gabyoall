<script lang="ts">
  import { api, db } from 'api'
  import Icon from './Icon.svelte'

  import Tab from './Tab.svelte'
  export let value = 'request'
  const onClick = (s: keyof typeof api) => (e: MouseEvent) => {
    value = s
    console.log('clicking', s)
    // api[s]?.list?.()
  }
  $: scheduleCount = Object.keys($db.schedule).length
  $: requestCount = Object.keys($db.request).length
  $: endpointCount = Object.keys($db.endpoint).length
</script>

<nav>
  <Tab active={value === 'endpoint'} on:click={onClick('endpoint')}>
    <Icon icon={'gEndpoint'} />
    Endpoints ({endpointCount})
  </Tab>
  <Tab active={value === 'request'} on:click={onClick('request')}>
    <Icon icon={'gRequest'} />
    Requests ({requestCount})
  </Tab>
  <Tab
    active={value === 'schedule'}
    on:click={onClick('schedule')}
    disabled={!endpointCount && !requestCount}>
    <Icon icon={'gSchedula'} />
    Schedule ({scheduleCount})
  </Tab>
  <Tab active={value === 'stats'} on:click={onClick('stat')}>
    <Icon icon={'gStat'} />
    Stats
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
  .fas {
    color: var(--color-red-500);
  }
</style>
