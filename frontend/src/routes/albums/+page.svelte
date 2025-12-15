<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import { albumsApi } from '$lib/api';
	import { LoadingSpinner, Alert } from '$lib/components';
	import { formatDate } from '$lib/utils';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { onMount } from 'svelte';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
		}
	});

	const albumsQuery = createQuery({
		queryKey: ['albums'],
		queryFn: () => albumsApi.list(),
		enabled: $isAuthenticated
	});

	const handleAlbumClick = (albumId: number) => {
		goto(`/albums/${albumId}`);
	};

	const handleCreateAlbum = () => {
		goto('/albums/new');
	};
</script>

<svelte:head>
	<title>Albums - Suipic</title>
</svelte:head>

{#if $isAuthenticated}
	<div class="space-y-6">
		<div class="flex justify-between items-center">
			<h1 class="text-4xl font-bold">Albums</h1>
			{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
				<button class="btn btn-primary" on:click={handleCreateAlbum}>
					<Icon icon="mdi:plus" class="text-xl" />
					New Album
				</button>
			{/if}
		</div>

		{#if $albumsQuery.isLoading}
			<div class="flex justify-center py-20">
				<LoadingSpinner size="lg" />
			</div>
		{:else if $albumsQuery.isError}
			<Alert type="error" message={$albumsQuery.error?.message || 'Failed to load albums'} />
		{:else if $albumsQuery.data && $albumsQuery.data.length > 0}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each $albumsQuery.data as album}
					<div
						class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow cursor-pointer"
						on:click={() => handleAlbumClick(album.id)}
						on:keydown={(e) => e.key === 'Enter' && handleAlbumClick(album.id)}
						role="button"
						tabindex="0"
					>
						<figure class="h-48 bg-base-300">
							{#if album.thumbnailPhotoId}
								<img
									src={`/api/photos/${album.thumbnailPhotoId}/thumbnail`}
									alt={album.title}
									class="object-cover w-full h-full"
								/>
							{:else}
								<div class="flex items-center justify-center h-full">
									<Icon icon="mdi:image-multiple" class="text-6xl opacity-30" />
								</div>
							{/if}
						</figure>
						<div class="card-body">
							<h2 class="card-title">{album.title}</h2>
							{#if album.description}
								<p class="text-sm opacity-70 line-clamp-2">{album.description}</p>
							{/if}
							<div class="flex items-center gap-2 text-sm opacity-60 mt-2">
								{#if album.dateTaken}
									<div class="flex items-center gap-1">
										<Icon icon="mdi:calendar" class="text-base" />
										<span>{formatDate(album.dateTaken)}</span>
									</div>
								{/if}
								{#if album.location}
									<div class="flex items-center gap-1">
										<Icon icon="mdi:map-marker" class="text-base" />
										<span class="truncate">{album.location}</span>
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<div class="text-center py-20">
				<div class="flex flex-col items-center gap-4">
					<Icon icon="mdi:image-multiple" class="text-8xl opacity-30" />
					<p class="text-xl opacity-60">No albums found</p>
					{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
						<button class="btn btn-primary" on:click={handleCreateAlbum}>
							Create Your First Album
						</button>
					{/if}
				</div>
			</div>
		{/if}
	</div>
{/if}
