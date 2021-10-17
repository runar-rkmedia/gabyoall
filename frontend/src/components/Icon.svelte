<script type="ts">
  import { Icon, iconMap } from './icons'

  export let icon: Icon
  export let kind = 'fas'
  export let color:
    | 'success'
    | 'warning'
    | 'error'
    | 'primary'
    | 'info'
    | 'secondary'
    | 'tertiary'
    | '' = ''
  const colorMap = {
    error: 'color-error',
    warning: 'color-warning',
    info: 'color-info',
  }
</script>

<i
  class={[
    color ? `color-${color}` : colorMap[icon],
    kind,
    iconMap[icon],
    $$props.class,
  ]
    .filter(Boolean)
    .join(' ')} />
