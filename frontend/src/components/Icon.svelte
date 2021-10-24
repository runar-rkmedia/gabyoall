<script type="ts">
  import type { Colors } from 'appTypes'

  import { Icon, iconMap } from './icons'

  export let icon: Icon
  export let kind = 'fas'
  export let color: Colors = ''
  const colorMap = {
    error: 'color-error',
    warning: 'color-warning',
    info: 'color-info',
  }
</script>

<i
  class={[
    color !== 'inherit' && (color ? `color-${color}` : colorMap[icon]),
    kind,
    iconMap[icon],
    $$props.class,
  ]
    .filter(Boolean)
    .join(' ')} />
