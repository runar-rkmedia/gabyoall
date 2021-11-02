<script type="ts">
  import type { Icon as IconType } from './icons'
  import Icon from './Icon.svelte'
  import type { Colors } from '../appTypes'
  import { createEventDispatcher } from 'svelte'
  export let preventDefault = true

  export let icon: IconType | undefined = undefined
  // TODO: support all colors
  export let color: Colors | 'danger' = ''
  export let disabled: boolean = false
  export let type: string = ''
  const dispatch = createEventDispatcher()
</script>

<button
  class:btn-reset={true}
  class={color}
  {type}
  class:icon-button={!!icon}
  {disabled}
  on:click={(e) => {
    if (preventDefault) {
      e.preventDefault()
    }
    dispatch('click', e)
  }}>
  {#if icon}
    <Icon {icon} />
  {/if}
  <slot />
</button>

<style>
  button.focus:not(:disabled),
  button:hover:not(:disabled) {
    /* TODO: move to psudo-element and transition opacity for perf. */
    box-shadow: 0 8px 16px -2px rgba(0, 32, 128, 0.25);
    transform: scale(1.05);
    transition: all 120ms ease-in-out;
  }
</style>
