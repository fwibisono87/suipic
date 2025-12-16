<script lang="ts">
	import { page } from '$app/stores';
	import { createQuery } from '@tanstack/svelte-query';
	import Icon from '@iconify/svelte';
	import { albumsApi, photosApi } from '$lib/api';
	import { LoadingSpinner, Alert } from '$lib/components';
	import ClientPhotoGallery from '$lib/components/ClientPhotoGallery.svelte';
	import { formatDate } from '$lib/utils';
	import { isAuthenticated } from '$lib/stores';
	import type { TPhoto } from '$lib/types';

	$: albumId = parseInt($page.params.id);

	let galleryLayout: 'grid' | 'masonry' = 'grid';
	let filterState: 'all' | 'pick' | 'reject' | 'unflagged' = 'all';

	const albumQuery = createQuery({
		queryKey: ['album', albumId],
		queryFn: () => albumsApi.get(albumId),
		enabled: $isAuthenticated && !!albumId
	});

	const photosQuery = createQuery({
		queryKey: ['photos', albumId],
		queryFn: () => photosApi.listByAlbum(albumId),
		enabled: $isAuthenticated && !!albumId
	});

	$: filteredPhotos = $photosQuery.data
		? $photosQuery.data.filter((photo: TPhoto) => {
				if (filterState === 'all') return true;
				if (filterState === 'unflagged') return !photo.pickRejectState || photo.pickRejectState === 'none';
				return photo.pickRejectState === filterState;
		  })
		: [];

	$: pickCount = $photosQuery.data?.filter((p: TPhoto) => p.pickRejectState === 'pick').length || 0;
	$: rejectCount = $photosQuery.data?.filter((p: TPhoto) => p.pickRejectState === 'reject').length || 0;
	$: unflaggedCount = $photosQuery.data?.filter((p: TPhoto) => !p.pickRejectState || p.pickRejectState === 'none').length || 0;

	const handlePhotoUpdate = () => {
		photosQuery.refetch();
	};
</script>

<svelte:head>
	<title>{$albumQuery.data?.title || 'Album'} - Suipic</title>
</svelte:head>

{#if $isAuthenticated}
	{#if $albumQuery.isLoading}
		<div class="flex justify-center py-20">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $albumQuery.isError}
		<Alert type="error" message={$albumQuery.error?.message || 'Failed to load album'} />
	{:else if $albumQuery.data}
		<div class="space-y-6">
			<div class="flex items-start justify-between gap-4">
				<div class="flex-1">
					<div class="breadcrumbs text-sm">
						<ul>
							<li><a href="/client/albums">My Albums</a></li>
							<li>{$albumQuery.data.title}</li>
						</ul>
					</div>
					<h1 class="text-4xl font-bold mt-2">{$albumQuery.data.title}</h1>
					{#if $albumQuery.data.description}
						<p class="mt-2 text-lg opacity-70">{$albumQuery.data.description}</p>
					{/if}
					<div class="flex flex-wrap items-center gap-4 mt-4 text-sm opacity-60">
						{#if $albumQuery.data.dateTaken}
							<div class="flex items-center gap-1">
								<Icon icon="mdi:calendar" class="text-base" />
								<span>{formatDate($albumQuery.data.dateTaken)}</span>
							</div>
						{/if}
						{#if $albumQuery.data.location}
							<div class="flex items-center gap-1">
								<Icon icon="mdi:map-marker" class="text-base" />
								<span>{$albumQuery.data.location}</span>
							</div>
						{/if}
						{#if $photosQuery.data}
							<div class="flex items-center gap-1">
								<Icon icon="mdi:image-multiple" class="text-base" />
								<span>{$photosQuery.data.length} {$photosQuery.data.length === 1 ? 'photo' : 'photos'}</span>
							</div>
						{/if}
					</div>
					{#if $albumQuery.data.customFields && Object.keys($albumQuery.data.customFields).length > 0}
						<div class="mt-4 space-y-2">
							<h3 class="font-semibold text-sm opacity-80">Album Details</h3>
							<div class="flex flex-wrap gap-2">
								{#each Object.entries($albumQuery.data.customFields) as [key, value]}
									<div class="badge badge-lg badge-outline">
										<span class="font-semibold">{key}:</span>
										<span class="ml-1">{value}</span>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			</div>

			{#if $photosQuery.data && $photosQuery.data.length > 0}
				<div class="divider"></div>

				<div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
					<div class="flex flex-wrap gap-2">
						<button
							class="btn btn-sm"
							class:btn-primary={filterState === 'all'}
							class:btn-outline={filterState !== 'all'}
							on:click={() => (filterState = 'all')}
						>
							<Icon icon="mdi:image-multiple" class="text-lg" />
							All ({$photosQuery.data.length})
						</button>
						<button
							class="btn btn-sm"
							class:btn-success={filterState === 'pick'}
							class:btn-outline={filterState !== 'pick'}
							on:click={() => (filterState = 'pick')}
						>
							<Icon icon="mdi:check-circle" class="text-lg" />
							Picks ({pickCount})
						</button>
						<button
							class="btn btn-sm"
							class:btn-ghost={filterState === 'unflagged'}
							class:btn-outline={filterState !== 'unflagged'}
							on:click={() => (filterState = 'unflagged')}
						>
							<Icon icon="mdi:flag-outline" class="text-lg" />
							Unflagged ({unflaggedCount})
						</button>
						<button
							class="btn btn-sm"
							class:btn-error={filterState === 'reject'}
							class:btn-outline={filterState !== 'reject'}
							on:click={() => (filterState = 'reject')}
						>
							<Icon icon="mdi:close-circle" class="text-lg" />
							Rejects ({rejectCount})
						</button>
					</div>

					<div class="flex gap-2">
						<button
							class="btn btn-xs sm:btn-sm btn-outline"
							class:btn-active={galleryLayout === 'grid'}
							on:click={() => (galleryLayout = 'grid')}
							aria-label="Grid layout"
						>
							<Icon icon="mdi:view-grid" class="text-base sm:text-xl" />
							<span class="hidden sm:inline">Grid</span>
						</button>
						<button
							class="btn btn-xs sm:btn-sm btn-outline"
							class:btn-active={galleryLayout === 'masonry'}
							on:click={() => (galleryLayout = 'masonry')}
							aria-label="Masonry layout"
						>
							<Icon icon="mdi:view-quilt" class="text-base sm:text-xl" />
							<span class="hidden sm:inline">Masonry</span>
						</button>
					</div>
				</div>
			{/if}

			<div class="divider"></div>

			{#if $photosQuery.isLoading}
				<div class="flex justify-center py-20">
					<LoadingSpinner size="lg" />
				</div>
			{:else if $photosQuery.isError}
				<Alert type="error" message={$photosQuery.error?.message || 'Failed to load photos'} />
			{:else if filteredPhotos.length > 0}
				<div>
					<ClientPhotoGallery 
						photos={filteredPhotos} 
						layout={galleryLayout} 
						onPhotoUpdate={handlePhotoUpdate} 
						photographerId={$albumQuery.data?.photographerId || null} 
					/>
				</div>
			{:else if $photosQuery.data && $photosQuery.data.length > 0}
				<div class="text-center py-20">
					<Icon icon="mdi:filter-off" class="text-8xl mx-auto opacity-30" />
					<p class="text-xl opacity-60 mt-4">No photos match the current filter</p>
					<button class="btn btn-primary mt-4" on:click={() => (filterState = 'all')}>
						Show All Photos
					</button>
				</div>
			{:else}
				<div class="text-center py-20">
					<Icon icon="mdi:camera" class="text-8xl mx-auto opacity-30" />
					<p class="text-xl opacity-60 mt-4">No photos in this album yet</p>
					<p class="text-sm opacity-50 mt-2">Your photographer will upload photos soon</p>
				</div>
			{/if}
		</div>
	{/if}
{/if}
