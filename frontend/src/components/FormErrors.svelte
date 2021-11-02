<script lang="ts">
  import getParentUntil from '../util/getParentUntil'

  type Obj = Record<string, unknown>
  type Errors<Data extends Obj | Obj[]> = {
    [key in keyof Data]?: Data[key] extends Obj
      ? Errors<Data[key]>
      : Data[key] extends Obj[]
      ? Errors<Data[key]>[]
      : string | string[] | null
  }
  export let formErrors: Errors<any>
  function isScrolledIntoView(el: Element) {
    var rect = el.getBoundingClientRect()
    var elemTop = rect.top
    var elemBottom = rect.bottom

    // Only completely visible elements return true:
    var isVisible = elemTop >= 0 && elemBottom <= window.innerHeight
    // Partially visible elements return true:
    //isVisible = elemTop < window.innerHeight && elemBottom >= 0;
    return isVisible
  }
  function clickEl(key: string) {
    const selector = `[name="${key}"`
    const el: HTMLInputElement | null = document.querySelector(selector)
    if (!el) {
      console.warn('Could not find element for selector: ', selector)
      return
    }
    const labelEl = getParentUntil(el, 'label')
    if (labelEl && !isScrolledIntoView(labelEl)) {
      labelEl.scrollIntoView()
    }
    el.focus()
  }
</script>

{#if Object.values(formErrors).some((err) => !!err && !!err.length)}
  Errors:
  <ul>
    {#each Object.entries(formErrors).filter(([_, val]) => val && val.length) as [key, value]}
      <li on:click={() => clickEl(key)}>
        <code>{key}: </code>
        {Array.isArray(value) ? value.join('. ') : value}
      </li>
    {/each}
  </ul>
{/if}
