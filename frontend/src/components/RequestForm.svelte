<script lang="ts">
  import { api } from '../api'
  import Collapse from './Collapse.svelte'
  import ConfigForm from './ConfigForm.svelte'
  import configStore from './configStore'
  import Spinner from './Spinner.svelte'

  let query = 'query {galaxies}'
  let variables = '{}'
  let body = ''
  let method = 'POST'
  let operationName = 'Name'
  let isGraphql = false
  let createResponse: ReturnType<typeof api.request.create> | undefined
  let loading = false
  async function endpointCreate() {
    loading = true
    const payload: ApiDef.RequestPayload = {
      method,
      operationName,
      ...(isGraphql
        ? {
            query,
            variables: JSON.parse(variables),
          }
        : {
            body,
          }),
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
    createResponse = api.request.create(payload)
    await createResponse
    loading = false
  }
</script>

<form>
  <label>
    {isGraphql ? 'Operation Name' : 'Label'}
    <input type="text" bind:value={operationName} />
  </label>
  <label>
    Method
    <input type="text" bind:value={method} />
  </label>
  <label>
    As GraphQL
    <input type="checkbox" bind:checked={isGraphql} />
  </label>
  {#if isGraphql}
    <label>
      query
      <textarea type="text" bind:value={query} />
    </label>
    <label>
      variables
      <textarea type="text" bind:value={variables} />
    </label>
  {:else}
    <label>
      body
      <textarea type="text" bind:value={body} />
    </label>
  {/if}
  <button
    disabled={loading || !!$configStore.__validationMessage}
    type="submit"
    on:click|preventDefault={endpointCreate}>Create request</button>
  <div class="spinner"><Spinner active={loading} /></div>
  <div class="paper">
    <Collapse>
      <h3 slot="title">Config</h3>
      <ConfigForm />
    </Collapse>
  </div>
  {#if createResponse}
    {#await createResponse then [_, err]}
      {#if err}
        {err.error} ({err.code})
      {/if}
    {/await}
  {/if}
</form>
