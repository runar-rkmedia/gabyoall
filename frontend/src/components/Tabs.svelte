<script lang="ts">
  import { api, db } from 'api'

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
  <Tab active={value === 'request'} on:click={onClick('request')}>
    <i class="fas fa-network-wired" /> Requests ({requestCount})
  </Tab>
  <Tab active={value === 'endpoint'} on:click={onClick('endpoint')}>
    <i class="fas fa-ethernet" /> Endpoints ({endpointCount})
  </Tab>
  <Tab
    active={value === 'schedule'}
    on:click={onClick('schedule')}
    disabled={!scheduleCount}
  >
    <i class="fas fa-calendar" /> Schedule ({scheduleCount})
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
