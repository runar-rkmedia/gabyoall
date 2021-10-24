<script lang="ts">
  import createId from 'createId'

  import { api, db } from '../api'
  import Code from './Code.svelte'
  import Collapse from './Collapse.svelte'
  import ConfigForm from './ConfigForm.svelte'
  import configStore from './configStore'
  import Icon from './Icon.svelte'
  import Spinner from './Spinner.svelte'

  export let editID = ''
  let createResponse: ReturnType<typeof api.endpoint.create> | undefined
  let loading = false
  let url = 'https://'
  let headers: Array<{
    key: string
    value: string
    /** Used internally wihtin the form to keep track of items*/
    id: string
  }> = []
  let lastEdit = ''
  $: {
    if (editID && lastEdit !== editID) {
      const e = $db.endpoint[editID]
      if (e) {
        if (e.config) {
          configStore.restore(e.config)
          console.debug(e.config)
        } else {
          configStore.reset()
        }

        lastEdit = editID
        url = e.url
        if (e.headers) {
          headers = Object.entries(e.headers).reduce((r, [k, v]) => {
            if (Array.isArray(v)) {
              for (const val of v) {
                // TODO: we should probably support multiple headers. (altough, where is this really needed in practice?)
                r[k] = v
              }
            } else {
              r[k] = v
            }
            return r
          }, {} as typeof headers)
        }
      }
    }
  }

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
      console.debug('payload-config', payload.config)
    }
    console.log('payload', payload, payload.config?.auth)
    loading = true
    createResponse = !!editID
      ? api.endpoint.update(editID, payload)
      : api.endpoint.create(payload)
    // createResponse = api.endpoint.create(payload)
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
  <paper>
    <Collapse key="endpoint-config">
      <h3 slot="title">Config</h3>
      <ConfigForm />
    </Collapse>
  </paper>
  <paper>
    <Collapse key="endpoint-headers">
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
  </paper>

  <button
    class="primary"
    disabled={loading || !!$configStore.__validationMessage}
    type="submit"
    on:click|preventDefault={endpointCreate}>
    {!!editID ? 'Update' : 'Create'}{' '}
    endpoint
  </button>
  <Code
    code={JSON.stringify($configStore.__validationPayload)}
    convert={true} />
  <Code
    code={JSON.stringify($configStore.__validationMessage)}
    convert={true} />
</form>

<style>
  .label-button {
    align-self: end;
    margin-bottom: var(--size-4);
  }
</style>
