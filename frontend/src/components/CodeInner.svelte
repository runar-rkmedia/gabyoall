<script lang="ts">
  export let language = 'json'
  export let code: string
  export let convert = false
  import formatterYaml from 'yaml'
  import formatterToml from 'toml-js'
  import Highlight from 'svelte-highlight'
  import yaml from 'svelte-highlight/src/languages/yaml'
  import json from 'svelte-highlight/src/languages/json'
  import toml from 'svelte-highlight/src/languages/ini'
  import atomOneDark from 'svelte-highlight/src/styles/atom-one-dark'
  import { state } from '../state'
  import Alert from './Alert.svelte'
  import Icon from './Icon.svelte'
  let errorMsg = ''
  function convertFormat(code: string | {}, format: string) {
    errorMsg = ''
    if (!code) {
      return ''
    }
    if (!format) {
      return code
    }
    try {
      const obj = typeof code === 'string' ? JSON.parse(code) : code
      switch (format) {
        case 'yaml':
        case 'yml':
          return formatterYaml.stringify(obj, { sortMapEntries: true })
        case 'toml':
          return formatterToml.dump(obj)
        case 'json':
          return JSON.stringify(obj, null, 2)
        default:
          console.error('Unsupported format', format)
          break
      }
    } catch (error) {
      console.error('failed to convert to format', { code, format, error })
      errorMsg = error
    }
    return code
  }
</script>

<svelte:head>
  {@html atomOneDark}
</svelte:head>

<div class:no-code={!code}>
  <button
    class="icon-button primary"
    on:click|preventDefault={() => {
      switch ($state.codeLanguage) {
        case 'json':
          $state.codeLanguage = 'yaml'
          break
        case 'yaml':
          $state.codeLanguage = 'toml'
          break
        case 'toml':
          $state.codeLanguage = 'json'
          break
        default:
          $state.codeLanguage = 'toml'
          break
      }
    }}>
    <Icon icon="code" />
    {$state.codeLanguage}</button>
  {#if errorMsg}
    <Alert kind="error">
      {errorMsg}
    </Alert>
  {/if}
  <Highlight
    language={{
      yaml: yaml,
      yml: yaml,
      json: json,
      toml: toml,
    }[language]}
    code={convert
      ? convertFormat(code, $state.codeLanguage || language)
      : code} />
</div>

<style>
  div {
    position: relative;
  }
  button {
    position: absolute;
    right: 0;
    top: calc(-1 * var(--size-3));
  }
  .no-code {
    opacity: 0.5;
  }
</style>
