<script lang="ts">
  import type { EditorFromTextArea, EditorConfiguration } from 'codemirror'
  import { onMount } from 'svelte'
  import type EditorType from './EditorInner.svelte'
  import Spinner from './Spinner.svelte'
  let loaded = false
  let Editor: EditorType
  onMount(() =>
    import('./EditorInner.svelte').then((editor) => {
      Editor = editor.default as any
      loaded = true
    })
  )

  export let id: string | undefined = undefined
  export let name: string | undefined = undefined
  export let value: string = ''
  export let initialValue: string = ''
  export let noFormatSelector: boolean = false
  export let initialLanguage: 'json' | 'toml' | 'yaml' | 'graphql' = 'json'
  export let config: EditorConfiguration = {}
</script>

{#if loaded}
  <svelte:component
    this={Editor}
    bind:id
    bind:name
    bind:initialValue
    bind:value
    bind:noFormatSelector
    bind:initialLanguage
    bind:config />
{:else}
  <div class="spinner"><Spinner /></div>
{/if}
