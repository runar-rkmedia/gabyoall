<script lang="ts">
  export let active: boolean
  export let disabled: boolean = false
</script>

<li class:active><button {disabled} class="" on:click><slot /></button></li>

<style>
  .active button {
    background-color: hsla(240, 60%, 80%, 20%);
  }
  button:hover:not(:disabled) {
    background-color: hsla(240, 60%, 80%, 30%);
  }
  .active button:hover:not(:disabled) {
    background-color: hsla(240, 60%, 80%, 40%);
  }
  li {
    all: initial;
  }
  button:disabled {
    cursor: unset;
    opacity: 0.5;
  }
  button {
    font-weight: bold;
    cursor: pointer;
    font-family: revert;
    position: relative;
    transition: transform 60ms ease-in-out;
    color: revert;
    color: var(--color-grey-100);
    /* background-color: var(--color-grey-700); */
    padding-inline: 12px;
    height: 100%;
  }
</style>
