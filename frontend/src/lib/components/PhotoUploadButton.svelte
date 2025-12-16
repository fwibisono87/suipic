<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createEventDispatcher } from 'svelte';

	export let accept = 'image/jpeg,image/jpg,image/png,image/webp,image/gif';
	export let multiple = true;
	export let disabled = false;
	export let size: 'xs' | 'sm' | 'md' | 'lg' = 'md';
	export let variant: 'primary' | 'secondary' | 'ghost' | 'outline' = 'primary';
	export let label = 'Upload Photos';

	const dispatch = createEventDispatcher<{
		filesSelected: File[];
	}>();

	function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		const files = Array.from(input.files || []);
		
		if (files.length > 0) {
			dispatch('filesSelected', files);
		}
		
		input.value = '';
	}

	$: buttonClass = `btn btn-${variant} btn-${size}`;
</script>

<label class={buttonClass} class:btn-disabled={disabled}>
	<Icon icon="mdi:upload" class="text-lg" />
	{label}
	<input
		type="file"
		{accept}
		{multiple}
		{disabled}
		class="hidden"
		on:change={handleFileSelect}
	/>
</label>
