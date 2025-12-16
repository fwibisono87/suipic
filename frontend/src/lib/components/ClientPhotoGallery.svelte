<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { TPhoto } from '$lib/types';
	import PhotoThumbnail from './PhotoThumbnail.svelte';
	import ClientLightbox from './ClientLightbox.svelte';
	import Icon from '@iconify/svelte';

	export let photos: TPhoto[];
	export let layout: 'grid' | 'masonry' = 'grid';
	export let onPhotoUpdate: ((photo: TPhoto) => void) | null = null;
	export let photographerId: number | null = null;

	let lightboxOpen = false;
	let currentPhotoIndex = 0;
	let containerElement: HTMLElement;
	let visibleRange = { start: 0, end: 50 };
	let windowWidth = 0;
	let itemHeight = 300;

	$: gridCols = getGridColumns(windowWidth);
	$: visiblePhotos = layout === 'grid' ? photos.slice(visibleRange.start, visibleRange.end) : photos;

	function getGridColumns(width: number): number {
		if (width >= 1280) return 5;
		if (width >= 1024) return 4;
		if (width >= 768) return 3;
		if (width >= 640) return 3;
		return 2;
	}

	function openLightbox(index: number) {
		currentPhotoIndex = index;
		lightboxOpen = true;
		document.body.style.overflow = 'hidden';
	}

	function closeLightbox() {
		lightboxOpen = false;
		document.body.style.overflow = '';
	}

	function handlePhotoUpdate(photo: TPhoto) {
		photos[currentPhotoIndex] = photo;
		photos = [...photos];
		if (onPhotoUpdate) {
			onPhotoUpdate(photo);
		}
	}

	function handleScroll() {
		if (!containerElement || layout === 'masonry') return;

		const rect = containerElement.getBoundingClientRect();
		const scrollTop = window.scrollY;
		const containerTop = containerElement.offsetTop;

		const viewportHeight = window.innerHeight;
		const buffer = viewportHeight * 2;

		const relativeScroll = scrollTop - containerTop;
		const startRow = Math.max(0, Math.floor((relativeScroll - buffer) / itemHeight));
		const endRow = Math.ceil((relativeScroll + viewportHeight + buffer) / itemHeight);

		const start = Math.max(0, startRow * gridCols);
		const end = Math.min(photos.length, (endRow + 1) * gridCols);

		if (start !== visibleRange.start || end !== visibleRange.end) {
			visibleRange = { start, end };
		}
	}

	function handleResize() {
		windowWidth = window.innerWidth;
		if (layout === 'grid') {
			handleScroll();
		}
	}

	$: if (layout === 'grid' && containerElement) {
		handleScroll();
	}

	onMount(() => {
		windowWidth = window.innerWidth;
		handleScroll();

		window.addEventListener('scroll', handleScroll, { passive: true });
		window.addEventListener('resize', handleResize, { passive: true });
	});

	onDestroy(() => {
		window.removeEventListener('scroll', handleScroll);
		window.removeEventListener('resize', handleResize);
		document.body.style.overflow = '';
	});
</script>

<div class="photo-gallery">
	{#if photos.length === 0}
		<div class="text-center py-20">
			<Icon icon="mdi:camera" class="text-8xl mx-auto opacity-30" />
			<p class="text-xl opacity-60 mt-4">No photos to display</p>
		</div>
	{:else}
		<div
			bind:this={containerElement}
			class="gallery-container"
			class:grid-layout={layout === 'grid'}
			class:masonry-layout={layout === 'masonry'}
		>
			{#if layout === 'grid'}
				<div
					class="grid gap-2 sm:gap-3 md:gap-4 grid-cols-2 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
				>
					{#each visiblePhotos as photo, index}
						<PhotoThumbnail {photo} onClick={() => openLightbox(visibleRange.start + index)} />
					{/each}
				</div>
			{:else}
				<div
					class="masonry-grid columns-2 sm:columns-3 md:columns-3 lg:columns-4 xl:columns-5 gap-2 sm:gap-3 md:gap-4"
				>
					{#each visiblePhotos as photo, index}
						<div class="break-inside-avoid mb-2 sm:mb-3 md:mb-4">
							<PhotoThumbnail {photo} onClick={() => openLightbox(index)} />
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<ClientLightbox {photos} {photographerId} bind:currentIndex={currentPhotoIndex} bind:isOpen={lightboxOpen} onClose={closeLightbox} onPhotoUpdate={handlePhotoUpdate} />
	{/if}
</div>

<style>
	.gallery-container {
		position: relative;
		width: 100%;
	}

	.masonry-grid {
		column-gap: 1rem;
	}

	@media (min-width: 640px) {
		.masonry-grid {
			column-gap: 1rem;
		}
	}
</style>
