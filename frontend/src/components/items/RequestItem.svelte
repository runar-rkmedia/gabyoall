<script type="ts">
  import { slide } from 'svelte/transition'
  import ConfigItem from './ConfigItem.svelte'

  export let request: ApiDef.RequestEntity
</script>

<li transition:slide|local>
  {request.method}
  {request.operationName}
  {new Date(request.createdAt).toLocaleTimeString()}
  {request.query || ''}
  {request.variables || ''}
  {request.body || ''}
  {#if request.config}
    <ConfigItem config={request.config} short={true} />
  {/if}
</li>
