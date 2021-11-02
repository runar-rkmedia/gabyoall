<script lang="ts">
  export let config: ApiDef.Config
  export let short = true
</script>

<div>
  <div class="">
    Request-count: <strong>{config.request_count}</strong>
  </div>
  <div class="">
    Concurrency: <strong>{config.concurrency}</strong>
  </div>

  {#if !short}
    {JSON.stringify(config)}
  {/if}
</div>
