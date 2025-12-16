<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { createQuery } from '@tanstack/svelte-query';
	import { albumsApi } from '$lib/api';
	import { Alert, AlbumForm, LoadingSpinner } from '$lib/components';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { onMount } from 'svelte';

	$: albumId = parseInt($page.params.id);

	let isLoading = false;
	let error = '';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}
	});

	const albumQuery = createQuery({
		queryKey: ['album', albumId],
		queryFn: () => albumsApi.get(albumId),
		enabled: $isAuthenticated && !!albumId
	});

	$: if (
		$albumQuery.data &&
		$currentUser &&
		$currentUser.role !== EUserRole.ADMIN &&
		$albumQuery.data.photographerId !== $currentUser.id
	) {
		goto(`/albums/${albumId}`);
	}

	interface AlbumFormData {
		title: string;
		description: string | null;
		location: string | null;
		dateTaken: string | null;
		customFields: Record<string, string> | null;
		thumbnailPhotoId: number | null;
		clientIds: number[];
	}

	const handleSubmit = async (formData: AlbumFormData) => {
		error = '';
		isLoading = true;

		try {
			const albumData = {
				title: formData.title,
				description: formData.description,
				location: formData.location,
				dateTaken: formData.dateTaken,
				customFields: formData.customFields,
				thumbnailPhotoId: formData.thumbnailPhotoId
			};

			await albumsApi.update(albumId, albumData);

			if (formData.clientIds && formData.clientIds.length > 0) {
				await albumsApi.assignUsers(albumId, formData.clientIds);
			}

			goto(`/albums/${albumId}`);
		} catch (err: unknown) {
			error = (err as { message: string }).message || 'Failed to update album';
		} finally {
			isLoading = false;
		}
	};

	const handleCancel = () => {
		goto(`/albums/${albumId}`);
	};
</script>

<svelte:head>
	<title>Edit Album - Suipic</title>
</svelte:head>

{#if $isAuthenticated}
	{#if $albumQuery.isLoading}
		<div class="flex justify-center py-20">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $albumQuery.isError}
		<Alert type="error" message={$albumQuery.error?.message || 'Failed to load album'} />
	{:else if $albumQuery.data}
		<div class="max-w-4xl mx-auto">
			<div class="mb-6">
				<div class="breadcrumbs text-sm">
					<ul>
						<li><a href="/albums">Albums</a></li>
						<li><a href="/albums/{albumId}">{$albumQuery.data.title}</a></li>
						<li>Edit</li>
					</ul>
				</div>
				<h1 class="text-4xl font-bold mt-2">Edit Album</h1>
				<p class="text-sm opacity-60 mt-2">Update album details and settings</p>
			</div>

			{#if error}
				<Alert type="error" message={error} dismissible onDismiss={() => (error = '')} />
			{/if}

			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<AlbumForm album={$albumQuery.data} onSubmit={handleSubmit} onCancel={handleCancel} {isLoading} />
				</div>
			</div>
		</div>
	{/if}
{/if}
