<script lang="ts">
  export let language: 'json' | 'yaml' | 'toml' = 'json'
  export let code: string
  export let noFormatSelector = false
  import Editor from './Editor.svelte'
</script>

<Editor
  initialValue={code}
  bind:value={code}
  {noFormatSelector}
  config={{ readOnly: true, mode: language }}
  initialLanguage={language}>
  <slot name="title" slot="title" />
</Editor>
