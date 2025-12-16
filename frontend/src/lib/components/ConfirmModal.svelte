<script lang="ts">
	import Icon from '@iconify/svelte';

	export let isOpen = false;
	export let title = 'Confirm Action';
	export let message = 'Are you sure you want to proceed?';
	export let confirmText = 'Confirm';
	export let cancelText = 'Cancel';
	export let confirmButtonClass = 'btn-error';
	export let onConfirm: () => void;
	export let onCancel: () => void;
</script>

{#if isOpen}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg flex items-center gap-2">
				<Icon icon="mdi:alert-circle" class="text-2xl text-warning" />
				{title}
			</h3>
			<p class="py-4">{message}</p>
			<div class="modal-action">
				<button class="btn btn-outline" on:click={onCancel}>
					{cancelText}
				</button>
				<button class="btn {confirmButtonClass}" on:click={onConfirm}>
					{confirmText}
				</button>
			</div>
		</div>
		<div class="modal-backdrop" on:click={onCancel} on:keydown={onCancel}></div>
	</div>
{/if}
