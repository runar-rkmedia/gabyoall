<script lang="ts">
  import { onMount } from 'svelte'

  let App
  let loaded = false
  onMount(() =>
    import('./App.svelte').then(({ default: d }) => {
      App = d
      loaded = true
    })
  )
</script>

<svelte:component this={App} />
