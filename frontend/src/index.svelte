<script lang="ts">
  import EndpointForm from './components/EndpointForm.svelte'
  import EndpointList from './components/EndpointList.svelte'
  import RequestForm from './components/RequestForm.svelte'
  import RequestLIst from './components/RequestLIst.svelte'
  let showView = ''

  const toggleView = (view: 'request' | 'endpoint') => {
    return () => {
      if (showView) {
        showView = ''
        return
      }
      showView = view
    }
  }
</script>

<h1>Gobyoall - Stress tester</h1>

<button on:click={toggleView('endpoint')}>
  {showView !== 'endpoint' ? 'Create endpoint' : 'Hide form'}
</button>
<button on:click={toggleView('request')}>
  {showView !== 'request' ? 'Create request' : 'Hide form'}
</button>
<div>
  {#if showView === 'endpoint'}
    <EndpointForm />
  {:else if showView === 'request'}
    <RequestForm />
  {/if}
</div>
<h2>Endpoints</h2>
<EndpointList />
<h2>requests</h2>
<RequestLIst />

<style>
  button {
    float: right;
  }
</style>
