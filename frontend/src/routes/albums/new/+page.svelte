<script lang="ts">
	import { goto } from '$app/navigation';
	import { albumsApi } from '$lib/api';
	import { Alert, LoadingSpinner } from '$lib/components';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole, type TCreateAlbumRequest } from '$lib/types';
	import { onMount } from 'svelte';

	let title = '';
	let description = '';
	let location = '';
	let dateTaken = '';
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

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		error = '';

		if (!title.trim()) {
			error = 'Title is required';
			return;
		}

		isLoading = true;

		try {
			const albumData: TCreateAlbumRequest = {
				title: title.trim(),
				description: description.trim() || null,
				location: location.trim() || null,
				dateTaken: dateTaken || null
			};

			const album = await albumsApi.create(albumData);
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

<div class="max-w-2xl mx-auto">
	<div class="mb-6">
		<h1 class="text-4xl font-bold">Create New Album</h1>
		<p class="text-sm opacity-60 mt-2">Create a new album to organize your photos</p>
	</div>

	{#if error}
		<Alert type="error" message={error} dismissible onDismiss={() => (error = '')} />
	{/if}

	<div class="card bg-base-100 shadow-xl">
		<div class="card-body">
			<form on:submit={handleSubmit} class="space-y-4">
				<div class="form-control">
					<label class="label" for="title">
						<span class="label-text">Title *</span>
					</label>
					<input
						type="text"
						id="title"
						bind:value={title}
						placeholder="Album title"
						class="input input-bordered w-full"
						disabled={isLoading}
						required
					/>
				</div>

				<div class="form-control">
					<label class="label" for="description">
						<span class="label-text">Description</span>
					</label>
					<textarea
						id="description"
						bind:value={description}
						placeholder="Album description"
						class="textarea textarea-bordered h-24"
						disabled={isLoading}
					></textarea>
				</div>

				<div class="form-control">
					<label class="label" for="location">
						<span class="label-text">Location</span>
					</label>
					<input
						type="text"
						id="location"
						bind:value={location}
						placeholder="Where the photos were taken"
						class="input input-bordered w-full"
						disabled={isLoading}
					/>
				</div>

				<div class="form-control">
					<label class="label" for="dateTaken">
						<span class="label-text">Date Taken</span>
					</label>
					<input
						type="date"
						id="dateTaken"
						bind:value={dateTaken}
						class="input input-bordered w-full"
						disabled={isLoading}
					/>
				</div>

				<div class="flex gap-4 mt-6">
					<button type="submit" class="btn btn-primary flex-1" disabled={isLoading}>
						{#if isLoading}
							<LoadingSpinner size="sm" />
						{:else}
							Create Album
						{/if}
					</button>
					<button
						type="button"
						class="btn btn-outline flex-1"
						on:click={handleCancel}
						disabled={isLoading}
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
</div>
