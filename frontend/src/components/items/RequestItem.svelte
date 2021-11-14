<script type="ts">
  import ConfigItem from './ConfigItem.svelte'
  import ListItem from '../ListItem.svelte'
  import formatDate from 'dates'
  import Collapse from 'components/Collapse.svelte'

  export let request: ApiDef.RequestEntity
  export let onEdit: ((ID: string) => void) | undefined = undefined
  export let onDelete: ((ID: string) => void) | undefined = undefined
  export let selectedID: string = ''
</script>

<ListItem
  deleteDisabled={selectedID === request.id}
  editDisabled={selectedID === request.id}
  {onEdit}
  {onDelete}
  ID={request.id}
  deleted={!!request.deleted}>
  <svelte:fragment slot="header">
    {request.operationName}
  </svelte:fragment>
  <svelte:fragment slot="description">
    Created: {formatDate(request.createdAt)}

    {#if request.updatedAt}
      Updated: {formatDate(request.updatedAt)}
    {/if}
  </svelte:fragment>
  <svelte:fragment slot="details">
    <Collapse key="item-config">
      <span slot="title">Details</span>
      {#if request.config}
        <div class="sub-item">
          <ConfigItem config={request.config} />
        </div>
      {/if}
    </Collapse>
  </svelte:fragment>
</ListItem>
