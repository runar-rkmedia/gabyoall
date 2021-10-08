<script lang="ts">
  import { api } from '../api'
  import Spinner from './Spinner.svelte'

  let url = 'https://'
  let createResponse: ReturnType<typeof api.endpoint.create>
  let loading = false
  async function endpointCreate() {
    loading = true
    let uri = url
    if (!uri.includes('://')) {
      uri = 'https://' + uri
    }
    createResponse = api.endpoint.create({ url: uri })
    await createResponse
    loading = false
  }
</script>

<form>
  <label>
    Url
    <input type="text" bind:value={url} /></label
  >
  <button
    disabled={loading}
    type="submit"
    on:click|preventDefault={endpointCreate}>Create endpoint</button
  >
  <div class="spinner"><Spinner active={loading} /></div>
  {#if createResponse}
    {#await createResponse then [_, err]}
      {#if err}
        {err.error} ({err.code})
      {/if}
    {/await}
  {/if}
</form>
