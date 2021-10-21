<script lang="ts">
  import createId from 'createId'

  import { api } from '../api'
  import Collapse from './Collapse.svelte'
  import ConfigForm from './ConfigForm.svelte'
  import configStore from './configStore'
  import Icon from './Icon.svelte'
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
    // TODO: use validation in store
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
    if ($configStore.__validationMessage) {
      console.error(
        'There was an error with validation of config',
        $configStore.__validationMessage
      )
      return
    }
    if ($configStore.__validationPayload) {
      payload.config = $configStore.__validationPayload
    }
    console.log('payload', payload)
    loading = true
    createResponse = api.endpoint.create(payload)
    await createResponse
    loading = false
  }
</script>

<form>
  <label>
    Url
    <input type="text" bind:value={url} /></label>
  <div class="spinner"><Spinner active={loading} /></div>
  {#if createResponse}
    {#await createResponse then [_, err]}
      {#if err}
        {err.error} ({err.code})
      {/if}
    {/await}
  {/if}
  <div class="paper">
    <Collapse>
      <h3 slot="title">Config</h3>
      <ConfigForm />
    </Collapse>
  </div>
  <div class="paper">
    <Collapse>
      <h3 slot="title">Headers</h3>

      {#each headers as h, i}
        <div class="label-group">
          <label>
            Header {h.key}
            <input type="text" bind:value={headers[i].key} />
          </label>
          <label>
            Value {h.value}
            <input type="text" bind:value={headers[i].value} />
          </label>
          <button
            class="icon-button label-button"
            on:click|preventDefault={() =>
              (headers = headers.filter((he) => he.id !== h.id))}>
            <Icon icon={'closeCross'} />
          </button>
        </div>
      {/each}
      <button
        class="secondary"
        disabled={loading}
        on:click|preventDefault={() =>
          (headers = [...headers, { key: '', value: '', id: createId() }])}>
        Add header
      </button>
    </Collapse>
  </div>

  <button
    class="primary"
    disabled={loading || !!$configStore.__validationMessage}
    type="submit"
    on:click|preventDefault={endpointCreate}>
    Create endpoint
  </button>
</form>

<style>
  .label-button {
    align-self: end;
    margin-bottom: var(--size-4);
  }
</style>
