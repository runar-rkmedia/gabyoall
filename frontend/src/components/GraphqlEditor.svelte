<script type="ts">
  import Editor from './Editor.svelte'
  import type { EditorConfiguration } from 'codemirror'

  export let id: string | undefined = undefined
  export let name: string | undefined = undefined
  export let value: string = ''
  export let initialValue = ''
  export let noFormatSelector = false
  export let initialLanguage: 'json' | 'toml' | 'yaml' | 'graphql' = 'json'
  export let config: EditorConfiguration = {}
  // TODO: add a schema-fetcher and highlighting for graphql
</script>

<Editor
  config={{
    mode: 'graphql',
    ...config,
  }}
  bind:id
  bind:name
  bind:value
  bind:initialValue
  bind:noFormatSelector
  bind:initialLanguage />
