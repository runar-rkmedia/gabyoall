<script lang="ts">
  export let language = 'json'
  export let code: string
  export let convert = false
  export let noFormatSelector = false
  import { onMount } from 'svelte'

  let component: any
  onMount(async () =>
    import('./CodeInner.svelte').then((d) => (component = d.default))
  )
</script>

{#if component}
  <svelte:component
    this={component}
    {noFormatSelector}
    {language}
    {code}
    {convert} />
{/if}
loading
