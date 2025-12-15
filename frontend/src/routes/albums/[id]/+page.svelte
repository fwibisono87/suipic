<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { createQuery } from '@tanstack/svelte-query';
	import Icon from '@iconify/svelte';
	import { albumsApi, photosApi } from '$lib/api';
	import { LoadingSpinner, Alert } from '$lib/components';
	import { formatDate } from '$lib/utils';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { onMount } from 'svelte';

	$: albumId = parseInt($page.params.id);

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
		}
	});

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

	let uploadingFile: File | null = null;
	let uploadError = '';
	let isUploading = false;

	const handleFileSelect = async (e: Event) => {
		const target = e.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;

		uploadingFile = file;
		uploadError = '';
		isUploading = true;

		try {
			await photosApi.create(albumId, file);
			await photosQuery.refetch();
			uploadingFile = null;
			target.value = '';
		} catch (err: unknown) {
			uploadError = (err as { message: string }).message || 'Upload failed';
		} finally {
			isUploading = false;
		}
	};

	const canUpload = () => {
		if (!$albumQuery.data || !$currentUser) return false;
		return (
			$currentUser.role === EUserRole.ADMIN ||
			$albumQuery.data.photographerId === $currentUser.id
		);
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
			<div class="flex items-start justify-between">
				<div class="flex-1">
					<div class="breadcrumbs text-sm">
						<ul>
							<li><a href="/albums">Albums</a></li>
							<li>{$albumQuery.data.title}</li>
						</ul>
					</div>
					<h1 class="text-4xl font-bold mt-2">{$albumQuery.data.title}</h1>
					{#if $albumQuery.data.description}
						<p class="mt-2 text-lg opacity-70">{$albumQuery.data.description}</p>
					{/if}
					<div class="flex items-center gap-4 mt-4 text-sm opacity-60">
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
					</div>
				</div>
				{#if canUpload()}
					<div class="flex gap-2">
						<label class="btn btn-primary" for="photo-upload">
							<Icon icon="mdi:upload" class="text-xl" />
							Upload Photo
						</label>
						<input
							id="photo-upload"
							type="file"
							accept="image/*"
							class="hidden"
							on:change={handleFileSelect}
							disabled={isUploading}
						/>
					</div>
				{/if}
			</div>

			{#if uploadError}
				<Alert type="error" message={uploadError} dismissible onDismiss={() => (uploadError = '')} />
			{/if}

			{#if isUploading}
				<div class="alert alert-info">
					<LoadingSpinner size="sm" />
					<span>Uploading {uploadingFile?.name}...</span>
				</div>
			{/if}

			<div class="divider"></div>

			{#if $photosQuery.isLoading}
				<div class="flex justify-center py-20">
					<LoadingSpinner size="lg" />
				</div>
			{:else if $photosQuery.isError}
				<Alert type="error" message={$photosQuery.error?.message || 'Failed to load photos'} />
			{:else if $photosQuery.data && $photosQuery.data.length > 0}
				<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
					{#each $photosQuery.data as photo}
						<div class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
							<figure class="h-48 bg-base-300">
								<img
									src={`/api/photos/${photo.id}/thumbnail`}
									alt={photo.title || photo.filename}
									class="object-cover w-full h-full"
								/>
							</figure>
							<div class="card-body p-4">
								<h3 class="card-title text-sm">{photo.title || photo.filename}</h3>
								<div class="flex items-center gap-2">
									{#if photo.stars > 0}
										<div class="flex">
											{#each Array(photo.stars) as _}
												<Icon icon="mdi:star" class="text-warning text-base" />
											{/each}
										</div>
									{/if}
									{#if photo.pickRejectState === 'pick'}
										<span class="badge badge-success badge-sm">Pick</span>
									{:else if photo.pickRejectState === 'reject'}
										<span class="badge badge-error badge-sm">Reject</span>
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<div class="text-center py-20">
					<Icon icon="mdi:camera" class="text-8xl mx-auto opacity-30" />
					<p class="text-xl opacity-60 mt-4">No photos in this album</p>
					{#if canUpload()}
						<label class="btn btn-primary mt-4" for="photo-upload-empty">
							Upload Your First Photo
						</label>
						<input
							id="photo-upload-empty"
							type="file"
							accept="image/*"
							class="hidden"
							on:change={handleFileSelect}
							disabled={isUploading}
						/>
					{/if}
				</div>
			{/if}
		</div>
	{/if}
{/if}
