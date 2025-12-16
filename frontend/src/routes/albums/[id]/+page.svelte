<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Icon from '@iconify/svelte';
	import { albumsApi, photosApi, photographerApi } from '$lib/api';
	import { LoadingSpinner, Alert, ConfirmModal } from '$lib/components';
	import { formatDate } from '$lib/utils';
	import { isAuthenticated, currentUser, authToken } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { onMount } from 'svelte';

	$: albumId = parseInt($page.params.id);

	const queryClient = useQueryClient();

	let uploadingFile: File | null = null;
	let uploadError = '';
	let isUploading = false;
	let showDeleteModal = false;
	let isDeleting = false;
	let deleteError = '';

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

	const canManage = () => {
		if (!$albumQuery.data || !$currentUser) return false;
		return (
			$currentUser.role === EUserRole.ADMIN ||
			$albumQuery.data.photographerId === $currentUser.id
		);
	};

	const albumUsersQuery = createQuery({
		queryKey: ['albumUsers', albumId],
		queryFn: () => albumsApi.getUsers(albumId),
		enabled: $isAuthenticated && !!albumId && canManage()
	});

	const clientsQuery = createQuery({
		queryKey: ['clients'],
		queryFn: () => photographerApi.listClients($authToken || ''),
		enabled: $isAuthenticated && !!$authToken && canManage()
	});

	$: assignedClients =
		$albumUsersQuery.data && $clientsQuery.data
			? $clientsQuery.data.filter((client) => $albumUsersQuery.data?.includes(client.id))
			: [];

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

	const handleEdit = () => {
		goto(`/albums/${albumId}/edit`);
	};

	const handleDelete = async () => {
		deleteError = '';
		isDeleting = true;

		try {
			await albumsApi.delete(albumId);
			queryClient.invalidateQueries({ queryKey: ['albums'] });
			goto('/albums');
		} catch (err: unknown) {
			deleteError = (err as { message: string }).message || 'Failed to delete album';
			showDeleteModal = false;
		} finally {
			isDeleting = false;
		}
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
							<li><a href="/albums">Albums</a></li>
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
					</div>
					{#if $albumQuery.data.customFields && Object.keys($albumQuery.data.customFields).length > 0}
						<div class="mt-4 space-y-2">
							<h3 class="font-semibold text-sm opacity-80">Custom Fields</h3>
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
					{#if canManage() && assignedClients.length > 0}
						<div class="mt-4 space-y-2">
							<h3 class="font-semibold text-sm opacity-80">Assigned Clients</h3>
							<div class="flex flex-wrap gap-2">
								{#each assignedClients as client}
									<div class="badge badge-lg badge-primary">
										<Icon icon="mdi:account" class="text-base mr-1" />
										{client.friendlyName || client.username}
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>
				<div class="flex gap-2 flex-shrink-0">
					{#if canManage()}
						<button class="btn btn-outline btn-sm" on:click={handleEdit}>
							<Icon icon="mdi:pencil" class="text-xl" />
							Edit
						</button>
						<button class="btn btn-error btn-outline btn-sm" on:click={() => (showDeleteModal = true)}>
							<Icon icon="mdi:delete" class="text-xl" />
							Delete
						</button>
					{/if}
				</div>
			</div>

			{#if deleteError}
				<Alert type="error" message={deleteError} dismissible onDismiss={() => (deleteError = '')} />
			{/if}

			{#if canManage()}
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
				<div>
					<h2 class="text-2xl font-bold mb-4">Photos</h2>
					<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
						{#each $photosQuery.data as photo}
							<div class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
								<figure class="h-48 bg-base-300">
									<img
										src={`/api/thumbnails/${photo.id}`}
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
				</div>
			{:else}
				<div class="text-center py-20">
					<Icon icon="mdi:camera" class="text-8xl mx-auto opacity-30" />
					<p class="text-xl opacity-60 mt-4">No photos in this album</p>
					{#if canManage()}
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

<ConfirmModal
	isOpen={showDeleteModal}
	title="Delete Album"
	message="Are you sure you want to delete this album? This action cannot be undone and will delete all photos in the album."
	confirmText="Delete"
	cancelText="Cancel"
	confirmButtonClass="btn-error"
	onConfirm={handleDelete}
	onCancel={() => (showDeleteModal = false)}
/>
