<script type="ts">
  import { state } from '../state'

  import Alert from './Alert.svelte'
  import Icon from './Icon.svelte'
  import Button from './Button.svelte'
  /**
   * Key should be a unique identifier for this tip, preferably human-readable
   *
   * @example
   * key="server-message"
   */
  export let key: string
  /** If the content is significanly changed, users who have already read the previous message,
   * will be presented again if the version is incremented.x
   */
  export let version: number = 0
  /** Indicated if the tip is read or not */
  $: read = ($state.seenHints[key]?.[0] ?? -1) >= version
</script>

<!-- 
	@component
	Tip is a stateful component. It will display the content in an alert of type info.

	If the user has achnolledged that the content is read by clicking _Got it_, the component will
	be collapsed from that point on.
	
	- The `key`-prop is required, and should be a unique key identifying the content.
	- Use the `version`-prop if the content is significantly changed
 -->
<Alert kind="info" collapse={read}>
  <slot />
  <Button
    color="secondary"
    on:click={() => {
      $state.seenHints[key] = [version, new Date()]
      read = false
    }}>
    <Icon icon={'success'} />
    Got it
  </Button>
</Alert>
