<script lang="ts">
  import { blur, draw, fade, fly, scale, slide } from 'svelte/transition'

  import Icon from './Icon.svelte'

  export let kind: 'error' | 'warning' | 'info' | 'success'
  export let collapse = false
</script>

{#if collapse}
  <div class="wrapper">
    <button
      class={`btn-reset circle ${kind}`}
      on:click|preventDefault={() => (collapse = !collapse)}>
      <Icon icon={kind} class="kind" color="inherit" />
    </button>
  </div>
{:else}
  <div transition:slide class={`alert ${kind}`}>
    <div class="alert-icon">
      <Icon icon={kind} class="kind" />
    </div>
    <div class="content">
      <div class="title">
        <slot name="title" />
      </div>
      <slot />
    </div>
  </div>
{/if}

<style>
  .wrapper {
    display: flex;
    justify-content: flex-end;
  }
  .circle {
    width: 1em;
    height: 1em;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    border-radius: 50%;
    outline-width: 1px;
    outline-style: solid;
    outline-color: var(--color-black);
    background-color: var(--color-black);
    color: var(--color-grey-100);
  }
  .circle.info {
    background-color: var(--color-info-icon);
  }
  .circle.warning {
    background-color: var(--color-warning-icon);
  }
  .circle.error {
    background-color: var(--color-error-icon);
  }
  .circle.success {
    background-color: var(--color-success-icon);
  }
  .alert {
    border-radius: var(--radius);
    padding-block: var(--size-4);
    padding-inline: var(--size-6);
    box-shadow: var(--elevation-4);
    outline-width: 1px;
    outline-style: solid;
    margin-block-start: var(--size-4);
    display: flex;
  }
  .title {
    display: block;
    font-size: 1.2rem;
    margin-block-end: var(--size-2);
    margin-block-start: -2px;
  }
  .alert-icon {
    margin-inline-end: var(--size-3);
    margin-block-start: 7px;
  }

  .error:not(.circle) {
    outline-color: var(--color-warning-icon);
    background-color: hsl(0, 88%, 96.1%);
    color: hsl(1, 49.6%, 24.9%);
  }
  .warning:not(.circle) {
    outline-color: var(--color-warning-icon);
    background-color: hsl(34.6, 100%, 94.9%);
    color: hsl(35.3, 100%, 28%);
  }
  .info:not(.circle) {
    outline-color: var(--color-info-icon);
    background-color: hsl(197.5, 85.7%, 94.5%);
    color: hsl(198.8, 98%, 19.2%);
  }
  .success:not(.circle) {
    outline-color: var(--color-success-icon);
    background-color: hsl(120, 38.5%, 94.9%);
    color: hsl(123, 40%, 19.6%);
  }
</style>
