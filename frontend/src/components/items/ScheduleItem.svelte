<script type="ts">
  import { db } from '../../api'
  import Collapse from '../../components/Collapse.svelte'

  import Icon from '../../components/Icon.svelte'
  import ListItem from '../../components/ListItem.svelte'
  import formatDate from '../../dates'

  import ConfigItem from './ConfigItem.svelte'
  import EndpointItem from './EndpointItem.svelte'
  import RequestItem from './RequestItem.svelte'

  export let schedule: ApiDef.ScheduleEntity
  export let onEdit: ((ID: string) => void) | undefined = undefined
  export let onDelete: ((ID: string) => void) | undefined = undefined
  export let selectedID: string = ''
</script>

<ListItem
  deleteDisabled={selectedID === schedule.id}
  editDisabled={selectedID === schedule.id}
  {onEdit}
  {onDelete}
  ID={schedule.id}
  deleted={!!schedule.deleted}>
  <svelte:fragment slot="header">
    {schedule.label}
  </svelte:fragment>

  <svelte:fragment slot="error">
    {#if schedule.lastError}
      The last time this ran, an error occured:
      <p>
        {schedule.lastError}
      </p>
    {/if}
  </svelte:fragment>
  <svelte:fragment slot="description">
    {schedule.dates?.map((d) => new Date(d).toISOString()).join(', ') || ''}
    Created: {formatDate(schedule.start_date)}

    Last Run: {formatDate(schedule.lastRun)}
  </svelte:fragment>
  <svelte:fragment slot="details">
    <Collapse key="bob">
      <span slot="title">Details</span>
      {#if schedule.config}
        <div class="sub-item">
          <ConfigItem config={schedule.config} />
        </div>
      {/if}
      {#if schedule.endpointID && !!$db.endpoint[schedule.endpointID]}
        <Icon slot="title" icon={'gEndpoint'} />
        Endpoint:
        <div class="sub-item">
          <EndpointItem endpoint={$db.endpoint[schedule.endpointID]} />
        </div>
      {/if}
      {#if schedule.requestID && !!$db.request[schedule.requestID]}
        <Icon slot="title" icon={'gRequest'} />
        Request:
        <div class="sub-item">
          <RequestItem request={$db.request[schedule.requestID]} />
        </div>
      {/if}
    </Collapse>
  </svelte:fragment>
</ListItem>
