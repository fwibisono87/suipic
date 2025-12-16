<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Icon from '@iconify/svelte';
	import type { TPhoto } from '$lib/types';
	import { config } from '$lib/config';
	import { LoadingSpinner } from '$lib/components';
	import ClientPhotoInteractionPanel from './ClientPhotoInteractionPanel.svelte';
	import { imageProtectionEnabled } from '$lib/stores';

	export let photos: TPhoto[];
	export let currentIndex: number = 0;
	export let isOpen: boolean = false;
	export let onClose: () => void = () => {};
	export let onPhotoUpdate: ((photo: TPhoto) => void) | null = null;
	export let photographerId: number | null = null;

	let isLoading = true;
	let isError = false;
	let imageElement: HTMLImageElement;
	let touchStartX = 0;
	let touchEndX = 0;
	let preloadedImages = new Map<number, HTMLImageElement>();
	let showInteractionPanel = false;

	$: currentPhoto = photos[currentIndex];
	$: photoUrl = currentPhoto ? `${config.apiUrl}/photos/${currentPhoto.id}` : '';
	$: canGoPrev = currentIndex > 0;
	$: canGoNext = currentIndex < photos.length - 1;

	$: if (isOpen && currentPhoto) {
		isLoading = true;
		isError = false;
		if (imageElement) {
			imageElement.src = photoUrl;
		}
		preloadAdjacentImages();
	}

	function preloadAdjacentImages() {
		if (canGoNext && !preloadedImages.has(currentIndex + 1)) {
			const img = new Image();
			img.src = `${config.apiUrl}/photos/${photos[currentIndex + 1].id}`;
			preloadedImages.set(currentIndex + 1, img);
		}
		if (canGoPrev && !preloadedImages.has(currentIndex - 1)) {
			const img = new Image();
			img.src = `${config.apiUrl}/photos/${photos[currentIndex - 1].id}`;
			preloadedImages.set(currentIndex - 1, img);
		}
	}

	function goToPrevious() {
		if (canGoPrev) {
			currentIndex--;
		}
	}

	function goToNext() {
		if (canGoNext) {
			currentIndex++;
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (!isOpen) return;

		switch (event.key) {
			case 'Escape':
				if (showInteractionPanel) {
					showInteractionPanel = false;
				} else {
					onClose();
				}
				break;
			case 'ArrowLeft':
				goToPrevious();
				break;
			case 'ArrowRight':
				goToNext();
				break;
			case 'i':
			case 'I':
				showInteractionPanel = !showInteractionPanel;
				break;
		}
	}

	function handlePhotoUpdate(event: CustomEvent<TPhoto>) {
		const updatedPhoto = event.detail;
		photos[currentIndex] = updatedPhoto;
		if (onPhotoUpdate) {
			onPhotoUpdate(updatedPhoto);
		}
	}

	function toggleInteractionPanel() {
		showInteractionPanel = !showInteractionPanel;
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	function handleImageLoad() {
		isLoading = false;
	}

	function handleImageError() {
		isLoading = false;
		isError = true;
	}

	function handleTouchStart(event: TouchEvent) {
		touchStartX = event.touches[0].clientX;
	}

	function handleTouchMove(event: TouchEvent) {
		touchEndX = event.touches[0].clientX;
	}

	function handleTouchEnd() {
		const swipeThreshold = 50;
		const diff = touchStartX - touchEndX;

		if (Math.abs(diff) > swipeThreshold) {
			if (diff > 0) {
				goToNext();
			} else {
				goToPrevious();
			}
		}

		touchStartX = 0;
		touchEndX = 0;
	}

	onMount(() => {
		document.addEventListener('keydown', handleKeyDown);
	});

	onDestroy(() => {
		document.removeEventListener('keydown', handleKeyDown);
		preloadedImages.clear();
	});
</script>

{#if isOpen}
	<div
		class="lightbox fixed inset-0 z-50 bg-black/95 flex items-center justify-center"
		on:click={handleBackdropClick}
		role="dialog"
		aria-modal="true"
		aria-label="Photo viewer"
		on:touchstart={handleTouchStart}
		on:touchmove={handleTouchMove}
		on:touchend={handleTouchEnd}
	>
		<div class="absolute top-2 sm:top-4 right-2 sm:right-4 z-10 flex gap-2">
			<button
				class="btn btn-circle btn-ghost text-white btn-sm sm:btn-md"
				class:btn-active={showInteractionPanel}
				on:click={toggleInteractionPanel}
				aria-label="Toggle actions panel"
			>
				<Icon icon="mdi:gesture-tap" class="text-xl sm:text-2xl" />
			</button>
			<button
				class="btn btn-circle btn-ghost text-white btn-sm sm:btn-md"
				on:click={onClose}
				aria-label="Close"
			>
				<Icon icon="mdi:close" class="text-xl sm:text-2xl" />
			</button>
		</div>

		<button
			class="btn btn-circle btn-ghost absolute left-2 sm:left-4 top-1/2 -translate-y-1/2 text-white z-10 btn-sm sm:btn-md"
			class:opacity-50={!canGoPrev}
			class:cursor-not-allowed={!canGoPrev}
			on:click={goToPrevious}
			disabled={!canGoPrev}
			aria-label="Previous photo"
		>
			<Icon icon="mdi:chevron-left" class="text-2xl sm:text-3xl" />
		</button>

		<button
			class="btn btn-circle btn-ghost absolute right-2 sm:right-4 top-1/2 -translate-y-1/2 text-white z-10 btn-sm sm:btn-md"
			class:opacity-50={!canGoNext}
			class:cursor-not-allowed={!canGoNext}
			on:click={goToNext}
			disabled={!canGoNext}
			aria-label="Next photo"
		>
			<Icon icon="mdi:chevron-right" class="text-2xl sm:text-3xl" />
		</button>

		<div class="max-w-[95vw] max-h-[90vh] flex items-center justify-center relative">
			{#if isLoading}
				<div class="absolute inset-0 flex items-center justify-center">
					<LoadingSpinner size="lg" />
				</div>
			{/if}

			{#if isError}
				<div class="text-center text-white">
					<Icon icon="mdi:image-broken" class="text-6xl mx-auto opacity-50" />
					<p class="mt-4 text-lg">Failed to load image</p>
				</div>
			{:else}
				<img
					bind:this={imageElement}
					src={photoUrl}
					alt={currentPhoto.title || currentPhoto.filename}
					class="max-w-full max-h-[90vh] object-contain transition-opacity duration-300 protected-image"
					class:opacity-0={isLoading}
					class:opacity-100={!isLoading}
					class:user-select-none={$imageProtectionEnabled}
					class:pointer-events-none={$imageProtectionEnabled}
					draggable={!$imageProtectionEnabled}
					on:load={handleImageLoad}
					on:error={handleImageError}
					on:contextmenu={handleContextMenu}
				/>
			{/if}
		</div>

		<div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-3 sm:p-6">
			<div class="max-w-4xl mx-auto">
				<div class="flex items-start justify-between gap-2 sm:gap-4 text-white">
					<div class="flex-1 min-w-0">
						{#if currentPhoto.title}
							<h2 class="text-base sm:text-xl font-semibold truncate">{currentPhoto.title}</h2>
						{/if}
						<p class="text-xs sm:text-sm opacity-80 truncate">{currentPhoto.filename}</p>
					</div>
					<div class="flex flex-col items-end gap-1 sm:gap-2 flex-shrink-0">
						<div class="text-xs sm:text-sm opacity-70">
							{currentIndex + 1} / {photos.length}
						</div>
						<div class="flex items-center gap-1 sm:gap-2">
							{#if currentPhoto.stars > 0}
								<div class="flex items-center gap-0.5">
									{#each Array(currentPhoto.stars) as _}
										<Icon icon="mdi:star" class="text-warning text-base sm:text-lg" />
									{/each}
								</div>
							{/if}
							{#if currentPhoto.pickRejectState === 'pick'}
								<span class="badge badge-success badge-xs sm:badge-sm">Pick</span>
							{:else if currentPhoto.pickRejectState === 'reject'}
								<span class="badge badge-error badge-xs sm:badge-sm">Reject</span>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="absolute top-12 sm:top-20 left-1/2 -translate-x-1/2 text-center text-white text-xs opacity-50 hidden sm:block">
			<p>Use arrow keys or swipe to navigate • I for actions • ESC to close</p>
		</div>

		{#if showInteractionPanel}
			<div 
				class="absolute right-0 top-0 h-full w-full sm:w-96 bg-base-100 shadow-2xl overflow-y-auto z-20 transition-transform"
				on:click|stopPropagation
			>
				<div class="p-4">
					<div class="flex items-center justify-between mb-4">
						<h2 class="text-lg font-semibold">Photo Actions</h2>
						<button
							class="btn btn-ghost btn-sm btn-circle"
							on:click={toggleInteractionPanel}
							aria-label="Close panel"
						>
							<Icon icon="mdi:close" class="text-xl" />
						</button>
					</div>
					<ClientPhotoInteractionPanel photo={currentPhoto} {photographerId} on:update={handlePhotoUpdate} />
				</div>
			</div>
		{/if}
	</div>
{/if}

<style>
	.lightbox {
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.btn-circle {
		transition: all 0.2s ease;
	}

	.btn-circle:hover:not(:disabled) {
		transform: scale(1.1);
		background-color: rgba(255, 255, 255, 0.1);
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
