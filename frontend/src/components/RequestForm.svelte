<script lang="ts">
  import { api, db } from '../api'
  import Collapse from './Collapse.svelte'
  import ConfigForm from './ConfigForm.svelte'
  import configStore from './configStore'
  import Spinner from './Spinner.svelte'
  import Editor from './Editor.svelte'
  import Button from './Button.svelte'
  import GraphqlEditor from './GraphqlEditor.svelte'

  let query = 'query {galaxies}'
  let variables = '{}'
  let body = ''
  let method = 'POST'
  let operationName = 'Name'
  let isGraphql = false
  export let editID = ''
  let lastEdit = ''

  let createResponse: ReturnType<typeof api.request.create> | undefined
  let loading = false
  $: {
    if (editID && lastEdit !== editID) {
      const r = $db.request[editID]
      if (r) {
        if (r.config) {
          configStore.restore(r.config)
          console.debug(r.config)
        } else {
          configStore.reset()
        }
        console.log(r)

        lastEdit = editID
        query = r.query || ''
        variables = r.variables ? JSON.stringify(r.variables) : '{}'
        body =
          typeof r.body === 'string' ? r.body : JSON.stringify(r.body) || ''
        method = r.method || ''
        operationName = r.operationName || ''
        isGraphql = !!r.query
      }
    }
  }

  async function endpointCreate() {
    loading = true
    const payload: ApiDef.RequestPayload = {
      method,
      operationName,
      ...(isGraphql
        ? {
            query,
            variables: JSON.parse(variables || '{}'),
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
    console.log('isF', isGraphql, payload)
    createResponse = !!editID
      ? api.request.update(editID, payload)
      : api.request.create(payload)
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
    <label for="query-editor">
      Query
      <GraphqlEditor
        name="query"
        id="query-editor"
        initialValue={query}
        initialLanguage="graphql"
        bind:value={query}
        noFormatSelector={true} />
    </label>
    <label for="variables-editor">
      Variables
      <Editor
        name="variables"
        id="variables-editor"
        initialValue={variables}
        bind:value={variables} />
    </label>
  {:else}
    <paper>
      <Collapse key="request-body">
        <label for="body-editor" slot="title"> Body </label>
        <Editor
          name="body"
          id="body-editor"
          initialValue={body}
          bind:value={body} />
      </Collapse>
    </paper>
  {/if}
  <div class="spinner"><Spinner active={loading} /></div>
  <paper>
    <Collapse key="request-config">
      <h3 slot="title">Config</h3>
      <ConfigForm />
    </Collapse>
  </paper>
  <Button
    icon="gRequest"
    color="primary"
    disabled={loading || !!$configStore.__validationMessage}
    type="submit"
    on:click={endpointCreate}>
    {editID ? 'Update' : 'Create'} request
  </Button>
  {#if createResponse}
    {#await createResponse then [_, err]}
      {#if err}
        {err.error} ({err.code})
      {/if}
    {/await}
  {/if}
</form>
