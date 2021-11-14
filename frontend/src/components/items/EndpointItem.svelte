<script lang="ts">
  import formatDate from '../../dates'

  import ConfigItem from './ConfigItem.svelte'
  import ListItem from '../ListItem.svelte'
  import Collapse from 'components/Collapse.svelte'

  export let endpoint: ApiDef.EndpointEntity
  export let onEdit: ((ID: string) => void) | undefined = undefined
  export let onDelete: ((ID: string) => void) | undefined = undefined
  export let selectedID: string = ''
</script>

<ListItem
  deleteDisabled={selectedID === endpoint.id}
  editDisabled={selectedID === endpoint.id}
  {onEdit}
  {onDelete}
  ID={endpoint.id}
  deleted={!!endpoint.deleted}>
  <svelte:fragment slot="header">
    {endpoint.url}
  </svelte:fragment>
  <svelte:fragment slot="description">
    Created: {formatDate(endpoint.createdAt)}

    {#if endpoint.updatedAt}
      Updated: {formatDate(endpoint.updatedAt)}
    {/if}
  </svelte:fragment>
  <svelte:fragment slot="details">
    <Collapse key="item-config">
      <span slot="title">Details</span>
      {#if endpoint.config}
        <div class="sub-item">
          <ConfigItem config={endpoint.config} />
        </div>
      {/if}
    </Collapse>
  </svelte:fragment>
</ListItem>
