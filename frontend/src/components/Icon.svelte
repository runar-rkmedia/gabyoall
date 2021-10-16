<script type="ts">
  import { Icon, iconMap } from './icons'

  export let icon: Icon
  export let kind = 'fas'
</script>

<i class={[kind, iconMap[icon], $$props.class].filter(Boolean).join(' ')} />
