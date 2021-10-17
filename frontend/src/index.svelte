<script lang="ts">
  import { api, db } from 'api'

  import ScheduleForm from 'components/ScheduleForm.svelte'
  import ScheduleList from 'components/ScheduleList.svelte'
  import ServerInfo from 'components/ServerInfo.svelte'
  import Stats from 'components/Stats.svelte'
  import Tabs from 'components/Tabs.svelte'
  import formatDate from 'dates'
  import { state } from 'state'
  import EndpointForm from './components/EndpointForm.svelte'
  import EndpointList from './components/EndpointList.svelte'
  import RequestForm from './components/RequestForm.svelte'
  import RequestLIst from './components/RequestLIst.svelte'
  let scheduleID = ''
  api.endpoint.list()
  api.request.list()
  api.schedule.list()
  api.serverInfo()
  api.stat.list()
</script>

<div class="wrapper">
<header>
  <img src="/android-chrome-192x192.png" alt="Logo" />
  <h1>Gobyoall - Stress tester</h1>
  <Tabs bind:value={$state.tab} />
</header>
<div />

  <main>
    {#if $state.tab === 'schedule'}
      <h2>Schedules</h2>
      <div class="paper">
        <ScheduleForm endpointID="" requestID="" bind:editID={scheduleID} />
      </div>
      <ScheduleList bind:selectedID={scheduleID} />
    {:else if $state.tab === 'endpoint'}
      <h2>Endpoints</h2>
      <div class="paper">
        <EndpointForm />
      </div>
      <div class="paper">
        <EndpointList />
      </div>
    {:else if $state.tab === 'stat'}
      <h2>Statistics</h2>
      <div class="paper">
        <Stats />
      </div>
    {:else}
      <h2>Request</h2>
      <div class="paper">
        <RequestForm />
      </div>
      <div class="paper">
        <RequestLIst />
      </div>
    {/if}
  </main>

  <footer>
    <ServerInfo />
  </footer>
</div>

<style>
  main {
    margin-block-end: var(--size-12);
  }
  .wrapper {
    background-color: var(--color-blue-300);
    display: flex;
    flex-direction: column;
    min-height: 100%;
  }
  header {
    background-color: var(--color-black);
    color: hsl(240, 80%, 95%);
    display: flex;
    box-shadow: var(--elevation-4);
  }
  header h1 {
    margin-inline: var(--size-4);
    align-self: center;
  }
  main {
    margin-inline: var(--size-4);
  }
  img {
    height: 100px;
    max-width: 20vw;
  }
  footer {
    margin-top: auto;
    display: flex;
    width: 100%;
    justify-content: space-between;
    padding: var(--size-4);
    background-color: var(--color-black);
    color: var(--color-grey-100);
  }
</style>
