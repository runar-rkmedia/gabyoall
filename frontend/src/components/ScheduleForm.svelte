<script lang="ts">
  import { serializeDate } from 'apiFetcher'

  import { objectKeys } from 'simplytyped'

  import { api, db } from '../api'
  import Collapse from './Collapse.svelte'
  import ConfigForm from './ConfigForm.svelte'
  import configStore from './configStore'
  import Spinner from './Spinner.svelte'

  let createResponse: ReturnType<typeof api.schedule.create> | undefined
  let loading = false
  let frequency = 0
  let label = ''
  // let maxInterJobConcurrency = 0
  let multiplier = 1
  // let offsets = 0
  let start_date_str = ''
  // onMount(() => {
  //   start_date_str = new Date().
  // })
  /** set if this should be an edit to an existing route*/
  export let editID = ''
  export let endpointID
  export let requestID

  function serializeInputDate(date: Date | string | undefined) {
    if (!date) {
      return ''
    }
    const s = typeof date == 'string' ? date : date.toISOString()
    return s.split('Z')[0]
  }
  function deserializeInputDate(s: string) {
    if (!s) {
      return null
    }
    const d = new Date(s)
    if (isNaN(d.getTime())) {
      return null
    }
    return new Date(
      new Date(s).getTime() - new Date().getTimezoneOffset() * 60e3
    )
  }
  let lastEdit = ''

  $: {
    if (editID && lastEdit !== editID) {
      const s = $db.schedule[editID]
      if (s) {
        if (s.config) {
          configStore.restore(s.config)
        } else {
          configStore.reset()
        }

        lastEdit = editID
        endpointID = s.endpointID
        requestID = s.requestID
        frequency = s.frequency || 0
        multiplier = s.multiplier || 0

        start_date_str = serializeInputDate(s.start_date)
        label = s.label || ''
      }
    }
  }
  function validate() {
    const now = new Date()
    const errors: Partial<
      Record<keyof ApiDef.SchedulePayload, string | null | false | undefined>
    > = {
      endpointID: !endpointID && 'No endpoint-id set',
      requestID: !requestID && 'No request-id set',
      frequency: frequency < 0 && 'Frequency cannot be negative',
      // maxInterJobConcurrency:
      //   maxInterJobConcurrency < 0 && 'maxInterJobConcurrency cannot be negative',
      multiplier: multiplier <= 0 && 'multiplier must be positive',
      // offsets: offsets < 0 && 'offsets cannot be negative',
      start_date: (() => {
        const t = deserializeInputDate(start_date_str)?.getTime()
        if (!t) {
          return 'Start-date must be set'
        }
        const diff = now.getTime() - t
        if (diff > 60 * 60e3) {
          return 'Start-time cannot be in the past'
        }
        return null
        // !start_date_str
        // ? 'Start-date must be set'
        // : now.getTime() - new Date(start_date_str).getTime() > 60 * 60e3
        // ? 'start_date cannot be in the past'
        // : '',
      })(),
      label: !label && 'Label must be set',
    }
    const errs = objectKeys(errors).reduce((r, k) => {
      if (!errors[k]) {
        return r
      }
      r[k] = errors[k]
      return r
    }, {} as typeof errors)
    if (!Object.keys(errs).length) {
      return null
    }
    return errs
  }
  $: errors = validate()
  $: valid = !errors
  $: disabled = loading || !valid || !!$configStore.__validationMessage
  async function scheduleCreate() {
    loading = true
    const d = deserializeInputDate(start_date_str)
    const payload: ApiDef.SchedulePayload = {
      endpointID,
      requestID,
      frequency,
      // maxInterJobConcurrency,
      multiplier,
      // offsets,
      ...(!!d && {
        start_date: serializeDate(d),
      }),
      label,
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

    createResponse = !!editID
      ? api.schedule.update(editID, payload)
      : api.schedule.create(payload)
    await createResponse
    loading = false
  }
  // FIXME: I'm sure there is a better way to do this in svelte...
  setInterval(() => {
    errors = validate()
  }, 100)
</script>

<h3>
  {!!editID
    ? `Editing schedule ${$db.schedule[editID]?.label}`
    : 'Creating schedule'}
</h3>
<form>
  <label>
    Label
    <input type="text" name="label" bind:value={label} />
  </label>
  <label>
    Endpint:
    <select bind:value={endpointID}>
      {#each Object.values($db.endpoint).filter((e) => !e.deleted) as v}
        <option value={v.id}>{v.url}</option>
      {/each}
    </select>
  </label>
  <label>
    Request
    <select bind:value={requestID}>
      {#each Object.entries($db.request) as [id, v]}
        <option value={id}>{v.operationName}</option>
      {/each}
    </select>
  </label>
  <label>
    Start-time {start_date_str}
    {start_date_str ? new Date(start_date_str) : 'no'}
    <input
      type="datetime-local"
      name="start_date"
      bind:value={start_date_str} />
  </label>
  <button
    on:click|preventDefault={() =>
      (start_date_str = serializeInputDate(new Date()))}>Now</button>
  {#if errors}
    {#each Object.values(errors) as err}
      <div class="error">{err}</div>
    {/each}
  {/if}
  <paper>
    <Collapse key="schedule-config">
      <h3 slot="title">Config</h3>
      <ConfigForm />
    </Collapse>
  </paper>
  <button {disabled} type="submit" on:click|preventDefault={scheduleCreate}>
    {!!editID ? 'Update' : 'Create'}
  </button>

  <div class="spinner"><Spinner active={loading} /></div>
  {#if createResponse}
    {#await createResponse then [_, err]}
      {#if err}
        {err.error} ({err.code})
      {/if}
    {/await}
  {/if}
</form>

<style>
  .error::before {
    content: 'ERR: ';
    color: var(--color-red-700);
  }
</style>
