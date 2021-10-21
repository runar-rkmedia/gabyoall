<script type="ts">
  import { db } from 'api'
  import Collapse from 'components/Collapse.svelte'

  import Icon from 'components/Icon.svelte'
  import formatDate from 'dates'

  import { slide } from 'svelte/transition'
  import ConfigItem from './ConfigItem.svelte'
  import EndpointItem from './EndpointItem.svelte'
  import RequestItem from './RequestItem.svelte'

  export let schedule: ApiDef.ScheduleEntity
</script>

<li transition:slide|local>
  <div class="item-content">
    <div class="item-header">
      {#if schedule.lastError}
        <Icon icon={'error'} />
      {/if}
      {schedule.label}
    </div>
    <div class="item-details">
      {schedule.dates?.map((d) => new Date(d).toISOString()).join(', ') || ''}
      {formatDate(schedule.start_date)}

      Last Run: {formatDate(schedule.lastRun)}
    </div>
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
  </div>
  <div class="item-actions">
    <button class="icon-button" on:click>
      <Icon icon="edit" />
      Edit
    </button>
  </div>
</li>

<style>
  li {
    background-color: var(--color-grey-100);
    display: flex;
    justify-content: space-between;
    width: 100%;
  }
  li:nth-child(even) {
    background-color: var(--color-grey-300);
  }
  .item-header {
    font-size: large;
  }
  .item-details {
    font-size: small;
  }
  .item-content {
    padding-inline: var(--size-4);
    margin-block-start: var(--size-3);
    margin-block-end: var(--size-2);
  }
  .sub-item {
    margin-inline-start: 16px;
  }
</style>
