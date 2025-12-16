<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from '@iconify/svelte';

	export let state: 'none' | 'pick' | 'reject' = 'none';
	export let disabled: boolean = false;
	export let size: 'sm' | 'md' = 'md';

	const dispatch = createEventDispatcher<{
		change: 'none' | 'pick' | 'reject';
	}>();

	function handleStateChange(newState: 'none' | 'pick' | 'reject') {
		if (disabled) return;
		dispatch('change', newState);
	}

	$: buttonSize = size === 'sm' ? 'btn-sm' : 'btn-md';
</script>

<div class="pick-reject-buttons btn-group w-full">
	<button
		class="btn {buttonSize} flex-1"
		class:btn-success={state === 'pick'}
		class:btn-outline={state !== 'pick'}
		disabled={disabled}
		on:click={() => handleStateChange(state === 'pick' ? 'none' : 'pick')}
		aria-label="Mark as pick"
	>
		<Icon icon="mdi:check-circle" class="text-lg" />
		<span class="hidden sm:inline">Pick</span>
	</button>
	<button
		class="btn {buttonSize} flex-1"
		class:btn-ghost={state === 'none'}
		class:btn-outline={state === 'none'}
		disabled={disabled || state === 'none'}
		on:click={() => handleStateChange('none')}
		aria-label="Mark as unflagged"
	>
		<Icon icon="mdi:flag-outline" class="text-lg" />
		<span class="hidden sm:inline">Unflagged</span>
	</button>
	<button
		class="btn {buttonSize} flex-1"
		class:btn-error={state === 'reject'}
		class:btn-outline={state !== 'reject'}
		disabled={disabled}
		on:click={() => handleStateChange(state === 'reject' ? 'none' : 'reject')}
		aria-label="Mark as reject"
	>
		<Icon icon="mdi:close-circle" class="text-lg" />
		<span class="hidden sm:inline">Reject</span>
	</button>
</div>
