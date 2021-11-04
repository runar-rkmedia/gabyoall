<script lang="ts">
  import { isRedacted } from '../store'
  import Code from './Code.svelte'
  import Alert from './Alert.svelte'
  import Collapse from './Collapse.svelte'
  import decodeJWT from '../decodeJWT'
  import store from './configStore'
  import Icon from './Icon.svelte'
  import Tip from './Tip.svelte'
  import { api } from '../api'

  /** Clears the form. This should be done after submission.*/
  export const clear = store.reset
  export const set = store.set
  let dryDynamicPromise: ReturnType<typeof api.dryDynamic>
  // Be aware of nullability. We need to know if a user did set a setting or not.

  $: {
    if ($store._.showDecodedJWT && $store.auth.token) {
      const result = decodeJWT($store.auth.token)
      $store._.decodedToken = result
    }
  }
  const redactionFocusHandler: svelte.JSX.FocusEventHandler<HTMLInputElement> =
    (e) => {
      if (isRedacted(e?.currentTarget?.value)) {
        e.currentTarget.value = ''
      }
    }
</script>

<Tip key="config-merge-order" version={2}>
  TIP: Configurations are merged in the following order:
  <ol>
    <li>Endpoint</li>
    <li>Request</li>
    <li>Schedule</li>
  </ol>

  <p>
    This means that if <code>Schedule</code> has defined a value for Request-Count
    in the Config, it will take precedence. Null-values are ignored
  </p>
</Tip>
<button
  class="icon-button danger"
  disabled={!$store.__didChange}
  on:click|preventDefault={() => store.restore()}>
  <Icon icon={'delete'} />
  Reset</button>

<div class="label-group">
  <label>
    Request-count
    <input type="number" min="1" bind:value={$store.request_count} /></label>
  <label>
    Concurrency
    <input
      type="number"
      min="1"
      {...$store.request_count && { max: $store.request_count }}
      bind:value={$store.concurrency} /></label>
</div>
<label>
  Include response-data in output
  <select bind:value={$store.response_data}>
    <option value={null} />
    <option value={true}>Include</option>
    <option value={false}>Exclude</option>
  </select>
</label>
{#if $store.__validationMessage}
  {#each Object.entries($store.__validationMessage) as [k, v]}
    {#if !k.startsWith('auth')}
      <Alert kind="error">
        {v}
      </Alert>
    {/if}
  {/each}
{/if}
<paper>
  <Collapse key="config-authentication">
    <h3 slot="title">
      {#if $store.__validationMessage && Object.keys($store.__validationMessage).some( (key) => key.includes('auth') )}
        <Icon icon={'warning'} class="error" />
      {/if}
      Authentication
    </h3>
    <label>
      Kind: Type {$store.auth.kind}
      <select name="authentication.kind" bind:value={$store.auth.kind}>
        <option value={null}>None</option>
        <option value="bearer">Bearer-token (raw)</option>
        <option value="impersonation">Impersonation (keycloak)</option>
        <option value="dynamic">Dynamic</option>
      </select>
    </label>
    {#if $store.auth.kind === 'bearer'}
      <label>
        Token (JWT only):
        <input type="text" bind:value={$store.auth.token} />
      </label>
      <label class="checkbox">
        Show decoded <input
          type="checkbox"
          disabled={!$store.auth.token}
          bind:checked={$store._.showDecodedJWT} /></label>
      {#if $store._.showDecodedJWT}
        <paper>
          {#each $store._.decodedToken as [code, err]}
            {#if err}
              <Alert kind="error">
                {err}
              </Alert>
            {/if}
            <Code {code} convert={true} />
          {/each}
        </paper>
      {/if}
    {:else if $store.auth.kind === 'dynamic'}
      <Tip key="auth-dynamic">
        Dynamic authentication is used where the authentication can be done in a
        series of requests, which ultimatly results in a token of some kind.

        <p>
          Each request is passed to eachother and the results can be further
          nested to any final result
        </p>
      </Tip>
      <Alert kind="warning">
        This functionality is experimental, and <em>will change</em>.
      </Alert>
      <label>
        Header-key. The header-key to use for authentciation. Defaults to <code
          >authorization</code>
        <input type="text" bind:value={$store.auth.dynamic.headerKey} /></label>
      {#each $store.auth.dynamic.requests as r}
        <paper>
          <label>
            Uri
            <input type="text" name="" bind:value={r.uri} />
          </label>
          <label>
            Method
            <input type="text" name="" bind:value={r.method} />
          </label>
          <label>
            Body
            <textarea rows="4" name="" bind:value={r.body} />
          </label>
          <label class="checkbox">
            Json Request (short for setting <code
              >Content-Type: Application/json</code
            >)
            <input type="checkbox" name="" bind:checked={r.json_request} />
          </label>
          <label class="checkbox">
            If set, will marshal the response as json
            <input type="checkbox" name="" bind:checked={r.json_response} />
          </label>
          <label>
            <span>
              Optional GoLang-templating string. For instance:

              <code>
                {'{{.Response.Headers.Authorization | first}}'}
              </code>
              would extract the first <em>Authorization</em>-header
            </span>
            <Tip key="golang-templating">
              <p>Go has a very effective and powerful templating-language.</p>
              <p>
                The template-string is to be wrapped in a set of double curly
                braces like <code>{'{{content}}'}</code>
              </p>
              <p>
                You refer to variables by a <code>.</code>-prefix like
                <code>{'{{ .Response.BodyJson.Foo }}'}</code>
              </p>
              <p>
                Functions are available, and are used without a prefix. For
                instance <code>{' {{ .Response.BodyJson.Foo | title }}'}</code> would
                convert to title-case.
              </p>
              <p>
                Functions from <a href="http://masterminds.github.io/sprig/"
                  >Sprig</a> should all be available. (although a few are ommited
                for security)
              </p>
              <p>In addition, there are a few custom-functions:</p>
              <table>
                <tr>
                  <th>Name</th>
                  <th>Arg</th>
                  <th>Description</th>
                </tr>
                <tr>
                  <td>jmes</td>
                  <td>path-string</td>
                  <td
                    ><a href="https://jmespath.org/">JMES-path</a>-string to
                    extract</td>
                </tr>
              </table>
            </Tip>
            <input type="text" bind:value={r.result_jmes_path} /></label>
        </paper>
      {/each}
      {#await dryDynamicPromise then res}
        <Code code={JSON.stringify(res)} convert={true} />
      {/await}
      <Code code={JSON.stringify($store.auth.dynamic)} convert={true} />
      <button
        class="primary"
        on:click|preventDefault={() => {
          const d = $store.auth.dynamic
          dryDynamicPromise = api.dryDynamic({
            ...d,
            requests: d.requests.map(({ headers, ...req }) => {
              return {
                ...req,
                // headers: req.headers?.reduce((r, h) => {
                //   // TODO: map
                //   return r
                // }, {}),
              }
            }),
          })
        }}>
        Dry-run
      </button>
      <button
        class="secondary"
        on:click|preventDefault={() =>
          ($store.auth.dynamic.requests = [
            ...$store.auth.dynamic.requests,
            {
              headers: [],
              method: 'GET',
              uri: 'https://',
              body: '',
              json_request: true,
              json_response: true,
              result_jmes_path: '',
            },
          ])}>
        Add request
      </button>
    {:else if $store.auth.kind === 'impersonation'}
      <div class="label-group">
        <label>
          Client-ID
          <input type="text" bind:value={$store.auth.client_id} />
        </label>
        <label>
          Client-Secret
          <input
            type="text"
            on:focus={redactionFocusHandler}
            bind:value={$store.auth.client_secret} />
        </label>
      </div>
      <label>
        Redirection-url
        <input type="text" bind:value={$store.auth.redirect_uri} />
      </label>
      <label>
        Header-key. (defaults to 'Authorization')
        <input type="text" bind:value={$store.auth.header_key} />
      </label>
      <div class="label-group">
        <label>
          Impersonator-password
          <input
            type="text"
            bind:value={$store.auth.impersionation_credentials.password} />
        </label>
        <label>
          Impersonator-username
          <input
            type="text"
            bind:value={$store.auth.impersionation_credentials.username} />
        </label>
        <label>
          <select bind:value={$store.auth.endpoint_type}>
            <option value="keycloak">Keycloak</option>
          </select>
        </label>
      </div>
      {#if $store._.impersonateUserName}
        <label
          class:has-err={$store.__validationMessage?.[
            'auth.impersionation_credentials.userID/userName'
          ]}>
          User-name to impersonate
          <input
            type="text"
            bind:value={$store.auth.impersionation_credentials
              .user_name_to_impersonate} />
        </label>
      {:else}
        <label>
          User-id to impersonate
          <input
            type="text"
            bind:value={$store.auth.impersionation_credentials
              .user_id_to_impersonate} />
        </label>
      {/if}
      <label class="checkbox">
        Use Username to impersonate (requires <code
          >Rolemapping: manage-users</code
        >)
        <input
          type="checkbox"
          bind:checked={$store._.impersonateUserName} /></label>
    {/if}
    {#if $store.__validationMessage}
      {#each Object.entries($store.__validationMessage) as [k, v]}
        {#if k.startsWith('auth')}
          <Alert kind="error">
            {v}
          </Alert>
        {/if}
      {/each}
    {/if}
  </Collapse>
</paper>

<style>
  paper {
    margin-top: var(--size-8);
  }
  label.has-err {
    color: brown;
  }
  table {
    margin-block-end: var(--size-8);
  }
</style>
