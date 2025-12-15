<script lang="ts">
	export let type: 'info' | 'success' | 'warning' | 'error' = 'info';
	export let message: string;
	export let dismissible = false;
	export let onDismiss: (() => void) | undefined = undefined;

	const alertClasses = {
		info: 'alert-info',
		success: 'alert-success',
		warning: 'alert-warning',
		error: 'alert-error'
	};
</script>

<div class="alert {alertClasses[type]}">
	<span>{message}</span>
	{#if dismissible && onDismiss}
		<button class="btn btn-sm btn-ghost" on:click={onDismiss}>✕</button>
	{/if}
</div>
