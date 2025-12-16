<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import Icon from '@iconify/svelte';
	import { albumsApi } from '$lib/api';
	import { LoadingSpinner, Alert, Card } from '$lib/components';
	import { isAuthenticated } from '$lib/stores';
	import { formatDate } from '$lib/utils';
	import { config } from '$lib/config';

	const albumsQuery = createQuery({
		queryKey: ['albums'],
		queryFn: () => albumsApi.list(),
		enabled: $isAuthenticated
	});

	function getAlbumThumbnailUrl(album: any): string | null {
		if (!album.thumbnailPhotoId) return null;
		return `${config.apiUrl}/photos/${album.thumbnailPhotoId}?size=medium`;
	}
</script>

<svelte:head>
	<title>My Albums - Suipic</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-4xl font-bold">My Albums</h1>
			<p class="text-lg opacity-70 mt-2">View your assigned photo albums</p>
		</div>
	</div>

	<div class="divider"></div>

	{#if $albumsQuery.isLoading}
		<div class="flex justify-center py-20">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $albumsQuery.isError}
		<Alert type="error" message={$albumsQuery.error?.message || 'Failed to load albums'} />
	{:else if $albumsQuery.data && $albumsQuery.data.length === 0}
		<div class="text-center py-20">
			<Icon icon="mdi:folder-open" class="text-8xl mx-auto opacity-30" />
			<p class="text-xl opacity-60 mt-4">No albums assigned to you yet</p>
			<p class="text-sm opacity-50 mt-2">Your photographer will share albums with you soon</p>
		</div>
	{:else if $albumsQuery.data}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each $albumsQuery.data as album}
				<a href="/client/albums/{album.id}" class="block">
					<Card hoverable>
						<div class="aspect-video bg-base-300 rounded-lg overflow-hidden relative">
							{#if getAlbumThumbnailUrl(album)}
								<img
									src={getAlbumThumbnailUrl(album)}
									alt={album.title}
									class="w-full h-full object-cover"
								/>
							{:else}
								<div class="w-full h-full flex items-center justify-center">
									<Icon icon="mdi:image-multiple" class="text-6xl opacity-30" />
								</div>
							{/if}
						</div>
						<div class="p-4 space-y-2">
							<h3 class="text-xl font-semibold line-clamp-1">{album.title}</h3>
							{#if album.description}
								<p class="text-sm opacity-70 line-clamp-2">{album.description}</p>
							{/if}
							<div class="flex flex-wrap items-center gap-3 text-sm opacity-60">
								{#if album.dateTaken}
									<div class="flex items-center gap-1">
										<Icon icon="mdi:calendar" />
										<span>{formatDate(album.dateTaken)}</span>
									</div>
								{/if}
								{#if album.location}
									<div class="flex items-center gap-1">
										<Icon icon="mdi:map-marker" />
										<span class="line-clamp-1">{album.location}</span>
									</div>
								{/if}
							</div>
						</div>
					</Card>
				</a>
			{/each}
		</div>
	{/if}
</div>
