<script lang="ts">
  import { api, db } from 'api'

  import ScheduleForm from 'components/ScheduleForm.svelte'
  import ScheduleList from 'components/ScheduleList.svelte'
  import Tabs from 'components/Tabs.svelte'
  import { state } from 'state'
  import EndpointForm from './components/EndpointForm.svelte'
  import EndpointList from './components/EndpointList.svelte'
  import RequestForm from './components/RequestForm.svelte'
  import RequestLIst from './components/RequestLIst.svelte'
  let scheduleID = ''
  api.endpoint.list()
  api.request.list()
  api.schedule.list()
</script>

<header>
  <img src="/android-chrome-192x192.png" alt="Logo" />
  <h1>Gobyoall - Stress tester</h1>
  <Tabs bind:value={$state.tab} />
</header>

<main>
  {#if $state.tab === 'endpoint'}
    <h2>Endpoints</h2>
    <div class="paper">
      <EndpointForm />
    </div>
    <div class="paper">
      <EndpointList />
    </div>
  {:else if $state.tab === 'request'}
    <h2>Request</h2>
    <div class="paper">
      <RequestForm />
    </div>
    <div class="paper">
      <RequestLIst />
    </div>
  {:else}
    <h2>Schedules</h2>
    <div class="paper">
      <ScheduleForm endpointID="" requestID="" bind:editID={scheduleID} />
    </div>
    <ScheduleList bind:selectedID={scheduleID} />
  {/if}
</main>

<style>
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
</style>
