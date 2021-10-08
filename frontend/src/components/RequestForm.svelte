<script lang="ts">
  import { api } from '../api'
  import Spinner from './Spinner.svelte'

  let query = 'query {galaxies}'
  let variables = '{}'
  let body = ''
  let method = '{}'
  let operationName = 'Name'
  let isGraphql = false
  let createResponse: ReturnType<typeof api.request.create>
  let loading = false
  async function endpointCreate() {
    loading = true
    createResponse = api.request.create({
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
    })
    await createResponse
    loading = false
  }
</script>

<form>
  <label>
    operationName
    <textarea type="text" bind:value={operationName} />
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
  {/if}
  <button
    disabled={loading}
    type="submit"
    on:click|preventDefault={endpointCreate}>Create request</button
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
