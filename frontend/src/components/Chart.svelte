<script type="ts">
  import Plotly from 'plotly.js-basic-dist-min'
  import { afterUpdate } from 'svelte'
  let div: HTMLDivElement
  export let data: any = {}
  export let layout: any = {}

  function createChart() {
    Plotly.newPlot(div, data, layout)
  }
  afterUpdate(() => {
    createChart()
  })
</script>

<div bind:this={div} />
