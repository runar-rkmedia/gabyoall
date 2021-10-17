<script lang="ts">
  import { api, db } from '../api'
  import { slide } from 'svelte/transition'
  import Spinner from './Spinner.svelte'
  import formatDate from 'dates'
  import Icon from './Icon.svelte'
  export let selectedID: string = ''
  let loading = true
</script>

<div class="spinner"><Spinner active={loading} /></div>
<ul>
  {#each Object.entries($db.schedule)
    .sort((a, b) => {
      const A = a[1].start_date || ''
      const B = b[1].start_date || ''
      if (A > B) {
        return 1
      }
      if (A < B) {
        return -1
      }

      return 0
    })
    .reverse() as [k, v]}
    <li id={k} transition:slide|local>
      <div class="item-content">
        <div class="item-header">
          {#if v.lastError}
            <Icon icon={'error'} />
          {/if}
          {v.label}
        </div>
        <div class="item-details">
          {v.dates?.map((d) => new Date(d).toISOString()).join(', ') || ''}
          {formatDate(v.start_date)}

          Last Run: {formatDate(v.lastRun)}
        </div>
      </div>
      <div class="item-actions">
        <button class="icon-button" on:click={() => (selectedID = k)}>
          <Icon icon="edit" />
          Edit
        </button>
      </div>
    </li>
  {/each}
</ul>

<style>
  .spinner {
    float: right;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    border-radius: var(--radius);
    box-shadow: var(--elevation-4);
  }
  li {
    background-color: var(--color-grey-100);
    display: flex;
    justify-content: space-between;
    width: 100%;
  }
  li:nth-child(even) {
    background-color: var(--color-grey-300);
  }
  .item-header {
    font-size: large;
  }
  .item-details {
    font-size: small;
  }
  .item-content {
    padding-inline: var(--size-4);
    margin-block-start: var(--size-3);
    margin-block-end: var(--size-2);
  }
</style>
