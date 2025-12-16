<script lang="ts">
	import { goto } from '$app/navigation';
	import { albumsApi } from '$lib/api';
	import { Alert, AlbumForm } from '$lib/components';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { onMount } from 'svelte';

	let isLoading = false;
	let error = '';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if (
			$currentUser?.role !== EUserRole.PHOTOGRAPHER &&
			$currentUser?.role !== EUserRole.ADMIN
		) {
			goto('/albums');
		}
	});

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
				customFields: formData.customFields
			};

			const album = await albumsApi.create(albumData);

			if (formData.clientIds && formData.clientIds.length > 0) {
				await albumsApi.assignUsers(album.id, formData.clientIds);
			}

			goto(`/albums/${album.id}`);
		} catch (err: unknown) {
			error = (err as { message: string }).message || 'Failed to create album';
		} finally {
			isLoading = false;
		}
	};

	const handleCancel = () => {
		goto('/albums');
	};
</script>

<svelte:head>
	<title>New Album - Suipic</title>
</svelte:head>

<div class="max-w-4xl mx-auto">
	<div class="mb-6">
		<div class="breadcrumbs text-sm">
			<ul>
				<li><a href="/albums">Albums</a></li>
				<li>New Album</li>
			</ul>
		</div>
		<h1 class="text-4xl font-bold mt-2">Create New Album</h1>
		<p class="text-sm opacity-60 mt-2">Create a new album to organize and share your photos</p>
	</div>

	{#if error}
		<Alert type="error" message={error} dismissible onDismiss={() => (error = '')} />
	{/if}

	<div class="card bg-base-100 shadow-xl">
		<div class="card-body">
			<AlbumForm onSubmit={handleSubmit} onCancel={handleCancel} {isLoading} />
		</div>
	</div>
</div>
