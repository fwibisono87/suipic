<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from '@iconify/svelte';

	export let rating: number = 0;
	export let maxStars: number = 5;
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let readonly: boolean = false;
	export let showClear: boolean = false;

	const dispatch = createEventDispatcher<{
		change: number;
	}>();

	let hoveredStar: number | null = null;

	$: displayRating = hoveredStar !== null ? hoveredStar : rating;
	$: sizeClass = {
		sm: 'text-base',
		md: 'text-2xl',
		lg: 'text-3xl'
	}[size];

	function handleStarClick(star: number) {
		if (readonly) return;
		const newRating = star === rating ? 0 : star;
		dispatch('change', newRating);
	}

	function handleMouseEnter(star: number) {
		if (!readonly) {
			hoveredStar = star;
		}
	}

	function handleMouseLeave() {
		hoveredStar = null;
	}

	function clearRating() {
		if (!readonly) {
			dispatch('change', 0);
		}
	}
</script>

<div class="star-rating flex items-center gap-1">
	{#each Array(maxStars) as _, index}
		{@const star = index + 1}
		<button
			class="btn btn-ghost p-1 min-h-0"
			class:cursor-pointer={!readonly}
			class:cursor-default={readonly}
			class:h-6={size === 'sm'}
			class:h-8={size === 'md'}
			class:h-10={size === 'lg'}
			disabled={readonly}
			on:click={() => handleStarClick(star)}
			on:mouseenter={() => handleMouseEnter(star)}
			on:mouseleave={handleMouseLeave}
			aria-label={`Rate ${star} stars`}
		>
			<Icon
				icon={star <= displayRating ? 'mdi:star' : 'mdi:star-outline'}
				class={`${sizeClass} ${star <= displayRating ? 'text-warning' : ''} ${
					star > displayRating ? 'opacity-30' : ''
				}`.trim()}
			/>
		</button>
	{/each}
	{#if showClear && rating > 0 && !readonly}
		<button
			class="btn btn-ghost btn-xs ml-1"
			on:click={clearRating}
		>
			Clear
		</button>
	{/if}
</div>

<style>
	.star-rating button:not(:disabled):hover {
		background-color: transparent;
	}
</style>
