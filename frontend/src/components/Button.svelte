<script type="ts">
  import type { Icon as IconType } from './icons'
  import Icon from './Icon.svelte'
  import type { Colors } from 'appTypes'

  export let icon: IconType | undefined = undefined
  // TODO: support all colors
  export let color: Colors | 'danger' = ''
  export let disabled: boolean = false
  export let type: string = ''
</script>

<button
  class:btn-reset={true}
  class={color}
  {type}
  class:icon-button={!!icon}
  {disabled}
  on:click>
  {#if icon}
    <Icon {icon} />
  {/if}
  <slot />
</button>
