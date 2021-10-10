<script lang="ts">
  import createId from 'createId'

  import { api } from '../api'
  import Spinner from './Spinner.svelte'

  let createResponse: ReturnType<typeof api.endpoint.create> | undefined
  let loading = false
  let url = 'https://'
  let headers: Array<{
    key: string
    value: string
    /** Used internally wihtin the form to keep track of items*/
    id: string
  }> = []
  async function endpointCreate() {
    loading = true
    const payload: ApiDef.EndpointPayload = {
      url,
      headers: headers.reduce((r, h) => {
        if (!h.key) {
          return r
        }
        r[h.key] = r[h.key] ? [...r[h.key], h.value] : [h.value]
        return r
      }, {}),
    }
    if (!payload.url.includes('://')) {
      payload.url = 'https://' + payload.url
    }
    payload.headers
    createResponse = api.endpoint.create(payload)
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
    on:click|preventDefault={endpointCreate}
  >
    Create endpoint
  </button>
  {#each headers as h, i}
    <div class="flexer">
      <label>
        Header {h.key}
        <input type="text" bind:value={headers[i].key} />
      </label>
      <label>
        Value {h.value}
        <input type="text" bind:value={headers[i].value} />
      </label>
      <button
        on:click|preventDefault={() =>
          (headers = headers.filter((he) => he.id !== h.id))}>X</button
      >
    </div>
  {/each}
  <div class="spinner"><Spinner active={loading} /></div>
  {#if createResponse}
    {#await createResponse then [_, err]}
      {#if err}
        {err.error} ({err.code})
      {/if}
    {/await}
  {/if}
  <button
    disabled={loading}
    on:click|preventDefault={() =>
      (headers = [...headers, { key: '', value: '', id: createId() }])}
  >
    Add header
  </button>
</form>

<style>
  .flexer {
    display: flex;
  }
  .flexer label {
    min-width: 20em;
  }
</style>
