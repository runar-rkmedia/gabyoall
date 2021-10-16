<script lang="ts">
  import Icon from './Icon.svelte'

  export let kind: 'error' | 'warning' | 'info' | 'success'
</script>

<div class={`alert ${kind}`}>
  <Icon icon={kind} class={kind} />
  <div class="content">
    <div class="title">
      <slot name="title" />
    </div>
    <slot />
  </div>
</div>

<style>
  .alert {
    border-radius: var(--radius);
    padding-block: var(--size-4);
    padding-inline: var(--size-6);
    box-shadow: var(--elevation-4);
    margin-block-start: var(--size-4);
    display: flex;
  }
  .title {
    display: block;
    font-size: 1.2rem;
    margin-block-end: var(--size-2);
    margin-block-start: -2px;
  }
  .alert .fas {
    min-width: var(--size-4);
    margin-inline-end: var(--size-3);
    margin-block-start: 7px;
  }

  .error {
    background-color: hsl(0, 88%, 96.1%);
    color: hsl(1, 49.6%, 24.9%);
  }
  .warning {
    background-color: hsl(34.6, 100%, 94.9%);
    color: hsl(35.3, 100%, 28%);
  }
  .info {
    background-color: hsl(197.5, 85.7%, 94.5%);
    color: hsl(198.8, 98%, 19.2%);
  }
  .success {
    background-color: hsl(120, 38.5%, 94.9%);
    color: hsl(123, 40%, 19.6%);
  }
</style>
