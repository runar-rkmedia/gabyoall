<script lang="ts">
  import reporter from '@felte/reporter-tippy'

  import { api, db } from '../api'
  import Button from './Button.svelte'
  import Collapse from './Collapse.svelte'
  import ConfigForm from './ConfigForm.svelte'
  import configStore from './configStore'
  import Spinner from './Spinner.svelte'
  import Tip from './Tip.svelte'
  import { createForm } from 'felte'
  import FormErrors from './FormErrors.svelte'

  const {
    form,
    errors: formErrors,
    isValid,
    data,
    validate,
    touched,
    setField,
    setFields,
  } = createForm<ApiDef.SchedulePayload>({
    extend: [reporter()],
    onSubmit: (_values) => {
      const values = {
        ..._values,
        start_date: deserializeInputDate(_values.start_date),
        end_date: deserializeInputDate(_values.end_date),
        monday: deserilizeInputDuration(_values.monday),
        tuesday: deserilizeInputDuration(_values.tuesday),
        wednesday: deserilizeInputDuration(_values.wednesday),
        thursday: deserilizeInputDuration(_values.thursday),
        friday: deserilizeInputDuration(_values.friday),
        saturday: deserilizeInputDuration(_values.saturday),
        sunday: deserilizeInputDuration(_values.sunday),
      }
      console.log('submit', values)
      scheduleCreate(values as any)
    },
    validate: (values) => {
      const errors: Record<string, string[]> = {}
      if ((values.label || '').length < 3)
        errors.label = ['A minimum of 3 characters must be supplied']
      if (!values.endpointID) errors.endpointID = ['Must be set']
      if (!values.requestID) errors.requestID = ['Must be set']
      if ((values.frequency || 0) < 0)
        errors.requestID = ['Must not be negative']

      return errors
    },
  })

  $: isTouched = Object.values($touched).some(Boolean)
  $: disabled =
    !isTouched || loading || !$isValid || !!$configStore.__validationMessage
  let createResponse: ReturnType<typeof api.schedule.create> | undefined
  let loading = false

  /** set if this should be an edit to an existing route*/
  export let editID = ''

  function deserilizeInputDuration(duration: string | undefined | number) {
    if (!duration) {
      return null
    }
    if (typeof duration === 'number') {
      throw new Error('What? why are you a number? Get FireFox!')
    }
    return duration.replace(':', 'h') + 'm'
  }
  function serializeInputDuration(duration: string | undefined | number) {
    console.log('dur', duration)
    if (!duration) {
      return ''
    }
    if (typeof duration === 'number') {
      throw new Error('What? why are you a number?')
    }
    const res = duration
      .split('m')[0]
      .split('h')
      .map((s) => s.padStart(2, '0'))
      .join(':')
    console.log('durOut', res)
    // TODO: type for weekdays are wrong
    return res as any
  }
  function serializeInputDate(date: Date | string | undefined) {
    if (!date) {
      return ''
    }
    const s = typeof date == 'string' ? date : date.toISOString()
    return s.slice(0, 16)
  }
  function deserializeInputDate(s: string | undefined) {
    if (!s) {
      return null
    }
    if ((s as any) instanceof Date) {
      return s
    }
    const d = new Date(s)
    if (isNaN(d.getTime())) {
      return null
    }
    var [y, m, ...rest] = s.split(/\D/).map(Number)
    const date = new Date(
      new Date(y, m - 1, ...rest).getTime() -
        new Date().getTimezoneOffset() * 60e3
    )
    return date
  }
  let lastEdit = ''

  // Change the form if the user chooses to edit a previous form
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
        setFields({
          ...s,
          start_date: serializeInputDate(s.start_date),
          end_date: serializeInputDate(s.end_date),
          monday: serializeInputDuration(s.monday),
          tuesday: serializeInputDuration(s.tuesday),
          wednesday: serializeInputDuration(s.wednesday),
          thursday: serializeInputDuration(s.thursday),
          friday: serializeInputDuration(s.friday),
          saturday: serializeInputDuration(s.saturday),
          sunday: serializeInputDuration(s.sunday),
        })
        validate()
      }
    }
  }

  async function scheduleCreate(payload: ApiDef.SchedulePayload) {
    loading = true
    // TODO: use form-library for config-store too
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
  const scheduleAsWeek = true
  const localTimeSone = Intl.DateTimeFormat().resolvedOptions().timeZone
  const weekDays: Record<
    keyof Pick<
      ApiDef.SchedulePayload,
      | 'monday'
      | 'tuesday'
      | 'wednesday'
      | 'thursday'
      | 'friday'
      | 'saturday'
      | 'sunday'
    >,
    string
  > = {
    monday: '',
    tuesday: '',
    wednesday: '',
    thursday: '',
    friday: '',
    saturday: '',
    sunday: '',
  }
  $: hasScheduleValue =
    !!$data.monday ||
    !!$data.tuesday ||
    !!$data.wednesday ||
    !!$data.thursday ||
    !!$data.friday ||
    !!$data.saturday ||
    $data.sunday
</script>

<h3>
  {!!editID
    ? `Editing schedule ${$db.schedule[editID]?.label}`
    : 'Creating schedule'}
</h3>

<form use:form>
  <label>
    <span class="required">Label:</span>
    <input
      type="text"
      name="label"
      required={true}
      minlength="3"
      maxLength="300" />
  </label>
  <div class="label-group">
    <label>
      <span class="required">Endpoint:</span>
      <select name="endpointID">
        {#if !$data.endpointID}
          <option value="" disabled={true}>None</option>
        {/if}
        {#each Object.values($db.endpoint).filter((e) => !e.deleted) as v}
          <option value={v.id}>{v.url}</option>
        {/each}
      </select>
    </label>
    <label>
      <span class="required">Request</span>
      <select name="requestID">
        {#if !$data.requestID}
          <option value="" disabled={true}>None</option>
        {/if}
        {#each Object.entries($db.request) as [id, v]}
          <option value={id}>{v.operationName}</option>
        {/each}
      </select>
    </label>
  </div>

  <div class="label-group">
    <label for="start_date">
      <span class="required">Start-time</span>
      <div class="input-button">
        <input id="start_date" type="datetime-local" name="start_date" />
        <Button
          icon="clock"
          color="secondary"
          on:click={() =>
            setField('start_date', serializeInputDate(new Date()))}>Now</Button>
      </div>
    </label>
    <label>
      End-time
      <input type="datetime-local" name="end_date" />
    </label>
  </div>

  {#if scheduleAsWeek}
    <paper class="weekdays">
      <Collapse key="weekdays">
        <h3 slot="title">Weekly schedule</h3>
        <Tip key="schedule-weekdays">
          <p>
            The schedule can be set to run on specific weekdays. You may
            additionally set a TimeZone of which these should run.
          </p>
        </Tip>
        {#each Object.keys(weekDays) as weekday}
          <label for={weekday} class="weekday">
            {weekday}
            <div>
              <Button
                color="secondary"
                icon="copy"
                disabled={!$data[weekday]}
                on:click={() => {
                  for (const w of Object.keys(weekDays)) {
                    if (w === weekday) {
                      continue
                    }
                    if (weekDays[w]) {
                      continue
                    }
                    setField(w, $data[weekday])
                  }
                }}>
                Copy to all unset
              </Button>
              <input
                id={weekday}
                name={weekday}
                placeholder="19:30m"
                type="time" />
            </div>
          </label>
        {/each}
        <hr />
        <label>
          TimeZone
          <input type="text" name={'location'} value={localTimeSone} />
        </label>
        <Button
          color="secondary"
          on:click={() => {
            setField(
              'location',
              $data.location !== localTimeSone ? localTimeSone : 'Local'
            )
          }}>
          {#if $data.location !== localTimeSone}
            Set to '{localTimeSone}'
          {:else}
            Set to servers local time
          {/if}
        </Button>

        {#if hasScheduleValue}
          <Button
            color="danger"
            on:click={() => {
              for (const w of Object.keys(weekDays)) {
                setField(w, '')
              }
            }}>
            Clear schedule
          </Button>
        {/if}
      </Collapse>
    </paper>
  {/if}

  <paper>
    <Collapse key="schedule-config">
      <h3 slot="title">Config</h3>
      <fieldset name="config">
        <ConfigForm />
      </fieldset>
    </Collapse>
  </paper>
  {#if !$isValid}
    <FormErrors formErrors={$formErrors} />
  {/if}
  <Button preventDefault={false} color="primary" {disabled} type="submit">
    {!!editID ? 'Update' : 'Create'}
  </Button>

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
  label.weekday {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    text-transform: capitalize;
  }
</style>
