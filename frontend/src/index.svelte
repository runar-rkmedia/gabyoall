<svelte:options immutable={true} />

<script lang="ts">
  import 'tippy.js/dist/tippy.css' // Tooltips popover
  import { api, db } from './api'

  import ScheduleForm from './components/ScheduleForm.svelte'
  import ScheduleList from './components/ScheduleList.svelte'
  import ServerInfo from './components/ServerInfo.svelte'
  import Stats from './components/Stats.svelte'
  import Tabs from './components/Tabs.svelte'
  import Alert from './components/Alert.svelte'
  import { state } from './state'
  import EndpointForm from './components/EndpointForm.svelte'
  import EndpointList from './components/EndpointList.svelte'
  import RequestForm from './components/RequestForm.svelte'
  import RequestLIst from './components/RequestLIst.svelte'
  import Button from './components/Button.svelte'
  let scheduleID = ''
  let endpointID = ''
  api.endpoint.list()
  api.request.list()
  api.schedule.list()
  api.serverInfo()
  api.stat.list()
  const dbWarnSizeGB = 0.5
  const dbWarnSize = dbWarnSizeGB * 1e9
</script>

<div class="wrapper">
  <header>
    <img src="/android-chrome-192x192.png" alt="Logo" />
    <h1>Gobyoall - Stress tester</h1>
    <Tabs bind:value={$state.tab} />
  </header>
  <div />

  <main>
    {#if ($db.serverInfo?.DatabaseSize || 0) > dbWarnSize}
      <Alert kind="warning">
        <div slot="title">The servers database has grown a bit big.</div>

        <p>It is currently {$db.serverInfo.DatabaseSizeStr}</p>
        <p>This may affect performance.</p>
        <p>Some functionality may have been disabled.</p>
        <p>It is adviced to clean the database</p>
        <Button icon="delete" color="danger" on:click={() => api.stat.clean()}
          >Click to remove Stats</Button>
      </Alert>
    {/if}
    {#if $state.tab === 'schedule'}
      <h2>Schedules</h2>
      <paper>
        {#if !Object.keys($db.endpoint).length}
          <p>You must first create Endpoint to be able to create a Schedule.</p>
        {:else if !Object.keys($db.request).length}
          <p>You must first create Request to be able to create a Schedule.</p>
        {:else}
          <ScheduleForm bind:editID={scheduleID} />
        {/if}
      </paper>
      <ScheduleList bind:selectedID={scheduleID} />
    {:else if $state.tab === 'endpoint'}
      <h2>Endpoints</h2>
      <paper>
        <EndpointForm bind:editID={endpointID} />
      </paper>
      <EndpointList bind:selectedID={endpointID} />
    {:else if $state.tab === 'stat'}
      <h2>Statistics</h2>
      <paper>
        <Stats />
      </paper>
    {:else}
      <h2>Request</h2>
      <paper>
        <RequestForm />
      </paper>
      <RequestLIst />
    {/if}
  </main>
  {#if $state.serverStats}
    <iframe
      title="Server Statistics"
      id="statsviz"
      height="600"
      width="100%"
      src="https://localhost/debug/statsviz/" />

    <a href="https://localhost/debug/statsviz/">Statwiz statistics</a>
  {/if}
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
