<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import CodeMirror, { fromTextArea } from 'codemirror'
  import 'codemirror/mode/javascript/javascript'
  import 'codemirror/mode/toml/toml'
  import 'codemirror/mode/yaml/yaml'
  import 'codemirror/theme/dracula.css'
  import 'codemirror/lib/codemirror.css'
  import type { EditorFromTextArea, EditorConfiguration } from 'codemirror'
  import { convertStringToCodeFormat } from 'util/codeFormat'
  import Button from './Button.svelte'
  import Alert from './Alert.svelte'
  import { state } from 'state'
  import Tip from './Tip.svelte'
  import CodeInner from './CodeInner.svelte'
  import Collapse from './Collapse.svelte'
  export let id: string | undefined = undefined
  export let name: string | undefined = undefined
  export let value: string = ''
  export let initialValue = ''
  export let noFormatSelector = false
  export let initialLanguage: 'json' | 'toml' | 'yaml' | 'graphql' = 'json'
  export let config: EditorConfiguration = {}
  // export let outFormat: 'json' | 'toml' | 'yaml' = initialLanguage

  let editor: EditorFromTextArea
  let errorMessage: string | null
  let textarea: HTMLTextAreaElement
  let _config: EditorConfiguration = {
    theme: 'dracula',
    mode: config?.mode || $state.codeLanguage,
    lineWrapping: true,
    lineNumbers: true,
    tabSize: 2,
    ...config,
  }
  if (_config.mode === 'json') {
    _config.mode = 'javascript'
  }
  $: {
    if (editor && config.readOnly) {
      editor.setValue(value)
      editor.setOption('mode', _config.mode)
    }
  }
  const setFormat = (format: typeof initialLanguage) => {
    editor.setOption('mode', format === 'json' ? 'javascript' : format)
    switch (format) {
      case 'json':
      case 'toml':
      case 'yaml':
        $state.codeLanguage = format
        break
      default:
        return
    }
    if ($state.editorRawFormat) {
      return
    }
    const [c, err] = convertStringToCodeFormat(editor.getValue(), format)
    if (!err) {
      editor.setValue(c as string)
      return
    }
  }
  function reset() {
    editor.setValue(initialValue || '')
    // if ($state.editorRawFormat) {
    //   return
    // }
    if (initialLanguage) {
      setFormat($state.codeLanguage || 'yaml')
      return
    }
  }
  onMount(() => {
    editor = fromTextArea(textarea, _config)
    editor.setSize('100%', '100%')
    editor.on('change', (e) => {
      if ($state.editorRawFormat) {
        value = editor.getValue()
        return
      }
      switch (_config.mode) {
        case 'javascript':
        case 'toml':
        case 'yaml':
          break
        default:
          return
      }
      const [c, err] = convertStringToCodeFormat(
        editor.getValue(),
        initialLanguage
      )
      errorMessage = err
      if (err) {
        return
      }
      value = c || ''
    })
    if (initialValue) {
      reset()
    }
  })

  onDestroy(() => {
    editor.toTextArea()
  })
  const langs: Array<typeof initialLanguage> = ['yaml', 'toml', 'json']
</script>

<slot name="title" />
{#if !config.readOnly}
  <div class="header">
    <Button color="danger" disabled={value === initialValue} on:click={reset}>
      Reset
    </Button>
  </div>
{/if}
<textarea bind:this={textarea} {name} {id} />
{#if !noFormatSelector}
  <div class="footer">
    <Button
      color="primary"
      toggle={$state.editorRawFormat}
      on:click={() => ($state.editorRawFormat = !$state.editorRawFormat)}>
      Raw
    </Button>
    {#each langs as l}
      <Button
        active={$state.codeLanguage === l}
        color="secondary"
        on:click={() => setFormat(l)}>{l}</Button>
    {/each}
  </div>
{/if}
{#if !config.readOnly}
  <Tip key="editor-format">
    <p>
      JSON is often used in API's, but is not always something that you want to
      edit manually.
    </p>
    <p>
      TOML and YAML are often considered better for human-readability and
      editing.
    </p>
    <p>
      You can therefore edit in a different language than what the backend
      supports
    </p>
    <p>Before sending, this value will be converted.</p>
    <p>
      If you do not want this behaviour, you can disable convertion by setting
      "Raw"
    </p>
  </Tip>
  {#if errorMessage}
    <Alert kind="error">{errorMessage}</Alert>
  {/if}
  {#if _config.mode !== 'graphql'}
    <paper>
      <Collapse key="editor-preview">
        <div slot="title">Preview</div>
        <CodeInner noFormatSelector={true} bind:code={value} />
      </Collapse>
    </paper>
  {/if}
{/if}

<style>
  .footer,
  .header {
    background-color: #2b2b2b;
    display: flex;
    justify-content: flex-end;
  }
</style>
