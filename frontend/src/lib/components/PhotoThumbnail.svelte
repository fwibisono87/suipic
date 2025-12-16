<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import type { TPhoto } from '$lib/types';
	import { config } from '$lib/config';
	import { imageProtectionEnabled } from '$lib/stores';

	export let photo: TPhoto;
	export let onClick: () => void = () => {};
	export let lazyLoad: boolean = true;

	let isLoaded = false;
	let isError = false;
	let imgElement: HTMLImageElement;
	let containerElement: HTMLElement;
	let observer: IntersectionObserver | null = null;

	const thumbnailUrl = `${config.apiUrl}/thumbnails/${photo.id}`;

	onMount(() => {
		if (!lazyLoad) {
			loadImage();
			return;
		}

		if ('IntersectionObserver' in window) {
			observer = new IntersectionObserver(
				(entries) => {
					entries.forEach((entry) => {
						if (entry.isIntersecting && !isLoaded && !isError) {
							loadImage();
							observer?.unobserve(entry.target);
						}
					});
				},
				{
					rootMargin: '100px'
				}
			);

			if (containerElement) {
				observer.observe(containerElement);
			}
		} else {
			loadImage();
		}

		return () => {
			if (observer) {
				observer.disconnect();
			}
		};
	});

	function loadImage() {
		if (imgElement) {
			imgElement.src = thumbnailUrl;
		}
	}

	function handleLoad() {
		isLoaded = true;
	}

	function handleError() {
		isError = true;
	}

	function handleContextMenu(event: MouseEvent) {
		// Prevent saving/downloading when image protection is enabled
		if ($imageProtectionEnabled) {
			event.preventDefault();
		}
	}

	function handleClick() {
		onClick();
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			onClick();
		}
	}
</script>

<div
	bind:this={containerElement}
	class="photo-thumbnail relative group cursor-pointer overflow-hidden rounded-md sm:rounded-lg bg-base-300 transition-transform active:scale-95 sm:hover:scale-[1.02]"
	on:click={handleClick}
	on:keydown={handleKeyDown}
	role="button"
	tabindex="0"
>
	<div class="aspect-square">
		{#if !isLoaded && !isError}
			<div class="absolute inset-0 flex items-center justify-center bg-base-300 animate-pulse">
				<Icon icon="mdi:image-outline" class="text-2xl sm:text-4xl opacity-30" />
			</div>
		{/if}

		{#if isError}
			<div class="absolute inset-0 flex flex-col items-center justify-center bg-error/10">
				<Icon icon="mdi:image-broken" class="text-2xl sm:text-4xl text-error opacity-50" />
				<span class="text-xs mt-2 opacity-50">Failed to load</span>
			</div>
		{:else}
			<img
				bind:this={imgElement}
				alt={photo.title || photo.filename}
				class="w-full h-full object-cover transition-opacity duration-300 protected-image"
				class:opacity-0={!isLoaded}
				class:opacity-100={isLoaded}
				class:user-select-none={$imageProtectionEnabled}
				class:pointer-events-none={$imageProtectionEnabled}
				draggable={!$imageProtectionEnabled}
				on:load={handleLoad}
				on:error={handleError}
				on:contextmenu={handleContextMenu}
			/>
		{/if}
	</div>

	{#if isLoaded}
		<div
			class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 group-active:opacity-100 transition-opacity"
		>
			<div class="absolute bottom-0 left-0 right-0 p-2 sm:p-3 text-white">
				<div class="flex items-center justify-between">
					<div class="flex-1 min-w-0">
						{#if photo.title}
							<p class="text-xs sm:text-sm font-semibold truncate">{photo.title}</p>
						{/if}
						<p class="text-[10px] sm:text-xs opacity-80 truncate">{photo.filename}</p>
					</div>
					<div class="flex items-center gap-1 sm:gap-2 ml-2 flex-shrink-0">
						{#if photo.stars > 0}
							<div class="flex items-center gap-0.5">
								<Icon icon="mdi:star" class="text-warning text-xs sm:text-sm" />
								<span class="text-[10px] sm:text-xs">{photo.stars}</span>
							</div>
						{/if}
						{#if photo.pickRejectState === 'pick'}
							<div class="badge badge-success badge-xs">Pick</div>
						{:else if photo.pickRejectState === 'reject'}
							<div class="badge badge-error badge-xs">Reject</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.photo-thumbnail {
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
		animation: fadeIn 0.3s ease-out;
	}

	.photo-thumbnail:hover {
		box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.protected-image::selection {
		background: transparent;
	}

	.user-select-none {
		user-select: none;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
	}

	.pointer-events-none {
		pointer-events: none;
	}
</style>
