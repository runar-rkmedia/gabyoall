<script lang="ts">
import { db } from "api";

import formatDate, { isValidDate } from "dates";

	$: serverStartedAt = isValidDate($db.serverInfo.ServerStartedAt)
	$: buildDate = isValidDate($db.serverInfo.BuildDate)
</script>

	<div>
		{#if serverStartedAt}
			
		{formatDate(serverStartedAt)}
		{/if}
		</div>
	<div>{$db.serverInfo.Version} ({$db.serverInfo.Version || 'development'})</div>
	<div>
		{#if buildDate}
			
		{formatDate(buildDate)}
		{/if}
	</div>
