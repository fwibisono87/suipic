<script lang="ts">
	import Icon from '@iconify/svelte';
	import { LoadingSpinner } from '$lib/components';
	import type { TAlbum } from '$lib/types';
	import { useListClients } from '$lib/queries/photographer';
	import { usePhotosQuery } from '$lib/queries/photos';
	import { useAlbumUsersQuery } from '$lib/queries/albums';

	export let album: TAlbum | null = null;
	export let onSubmit: (data: AlbumFormData) => Promise<void>;
	export let onCancel: () => void;
	export let isLoading = false;

	interface AlbumFormData {
		title: string;
		description: string | null;
		location: string | null;
		dateTaken: string | null;
		customFields: Record<string, string> | null;
		thumbnailPhotoId: number | null;
		clientIds: number[];
	}

	let title = album?.title || '';
	let description = album?.description || '';
	let location = album?.location || '';
	let dateTaken = album?.dateTaken ? album.dateTaken.split('T')[0] : '';
	let customFields: Array<{ key: string; value: string }> = album?.customFields
		? Object.entries(album.customFields).map(([key, value]) => ({
				key,
				value: String(value)
		  }))
		: [];
	let selectedClients: number[] = [];
	let thumbnailPhotoId: number | null = album?.thumbnailPhotoId || null;

	const clientsQuery = useListClients();
	const photosQuery = usePhotosQuery(album?.id || 0, !!album?.id);
	const albumUsersQuery = useAlbumUsersQuery(album?.id || 0, !!album?.id);

	$: if ($albumUsersQuery.data) {
		selectedClients = $albumUsersQuery.data;
	}

	const addCustomField = () => {
		customFields = [...customFields, { key: '', value: '' }];
	};

	const removeCustomField = (index: number) => {
		customFields = customFields.filter((_, i) => i !== index);
	};

	const toggleClient = (clientId: number) => {
		if (selectedClients.includes(clientId)) {
			selectedClients = selectedClients.filter((id) => id !== clientId);
		} else {
			selectedClients = [...selectedClients, clientId];
		}
	};

	const handleSubmit = async (e: Event) => {
		e.preventDefault();

		const formData: AlbumFormData = {
			title: title.trim(),
			description: description.trim() || null,
			location: location.trim() || null,
			dateTaken: dateTaken || null,
			customFields:
				customFields.length > 0
					? customFields.reduce(
							(acc, field) => {
								if (field.key.trim()) {
									acc[field.key.trim()] = field.value;
								}
								return acc;
							},
							{} as Record<string, string>
					  )
					: null,
			thumbnailPhotoId,
			clientIds: selectedClients
		};

		await onSubmit(formData);
	};

	$: isEditMode = !!album;
</script>

<form on:submit={handleSubmit} class="space-y-6">
	<div class="form-control">
		<label class="label" for="title">
			<span class="label-text font-semibold">Title *</span>
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
			<span class="label-text font-semibold">Description</span>
		</label>
		<textarea
			id="description"
			bind:value={description}
			placeholder="Album description"
			class="textarea textarea-bordered h-24"
			disabled={isLoading}
		></textarea>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
		<div class="form-control">
			<label class="label" for="location">
				<span class="label-text font-semibold">Location</span>
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
				<span class="label-text font-semibold">Date Taken</span>
			</label>
			<input
				type="date"
				id="dateTaken"
				bind:value={dateTaken}
				class="input input-bordered w-full"
				disabled={isLoading}
			/>
		</div>
	</div>

	<div class="form-control">
		<label class="label">
			<span class="label-text font-semibold">Custom Fields</span>
			<button
				type="button"
				class="btn btn-sm btn-outline"
				on:click={addCustomField}
				disabled={isLoading}
			>
				<Icon icon="mdi:plus" />
				Add Field
			</button>
		</label>
		<div class="space-y-2">
			{#each customFields as field, index}
				<div class="flex gap-2">
					<input
						type="text"
						bind:value={field.key}
						placeholder="Field name"
						class="input input-bordered flex-1"
						disabled={isLoading}
					/>
					<input
						type="text"
						bind:value={field.value}
						placeholder="Field value"
						class="input input-bordered flex-1"
						disabled={isLoading}
					/>
					<button
						type="button"
						class="btn btn-square btn-outline btn-error"
						on:click={() => removeCustomField(index)}
						disabled={isLoading}
					>
						<Icon icon="mdi:delete" />
					</button>
				</div>
			{/each}
		</div>
	</div>

	<div class="form-control">
		<label class="label">
			<span class="label-text font-semibold">Assign Clients</span>
		</label>
		{#if $clientsQuery.isLoading}
			<div class="flex justify-center p-4">
				<LoadingSpinner size="sm" />
			</div>
		{:else if $clientsQuery.isError}
			<p class="text-error text-sm">Failed to load clients</p>
		{:else if $clientsQuery.data && $clientsQuery.data.length > 0}
			<div class="border rounded-lg p-4 max-h-48 overflow-y-auto space-y-2">
				{#each $clientsQuery.data as client}
					<label class="flex items-center gap-2 cursor-pointer hover:bg-base-200 p-2 rounded">
						<input
							type="checkbox"
							class="checkbox checkbox-primary"
							checked={selectedClients.includes(client.id)}
							on:change={() => toggleClient(client.id)}
							disabled={isLoading}
						/>
						<span class="flex-1">
							{client.friendlyName || client.username}
							<span class="text-xs opacity-60">({client.email})</span>
						</span>
					</label>
				{/each}
			</div>
		{:else}
			<p class="text-sm opacity-60">No clients available</p>
		{/if}
	</div>

	{#if isEditMode && $photosQuery.data && $photosQuery.data.length > 0}
		<div class="form-control">
			<label class="label">
				<span class="label-text font-semibold">Thumbnail</span>
			</label>
			<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
				{#each $photosQuery.data as photo}
					<button
						type="button"
						class="relative aspect-square rounded-lg overflow-hidden border-2 transition-all {thumbnailPhotoId ===
						photo.id
							? 'border-primary ring-2 ring-primary'
							: 'border-base-300 hover:border-primary'}"
						on:click={() => (thumbnailPhotoId = photo.id)}
						disabled={isLoading}
					>
						<img
							src={`/api/thumbnails/${photo.id}`}
							alt={photo.title || photo.filename}
							class="w-full h-full object-cover"
						/>
						{#if thumbnailPhotoId === photo.id}
							<div
								class="absolute inset-0 bg-primary bg-opacity-20 flex items-center justify-center"
							>
								<Icon icon="mdi:check-circle" class="text-4xl text-primary" />
							</div>
						{/if}
					</button>
				{/each}
			</div>
		</div>
	{/if}

	<div class="flex gap-4 pt-4">
		<button type="submit" class="btn btn-primary flex-1" disabled={isLoading || !title.trim()}>
			{#if isLoading}
				<LoadingSpinner size="sm" />
			{:else}
				<Icon icon="mdi:check" class="text-xl" />
				{isEditMode ? 'Update Album' : 'Create Album'}
			{/if}
		</button>
		<button type="button" class="btn btn-outline flex-1" on:click={onCancel} disabled={isLoading}>
			<Icon icon="mdi:close" class="text-xl" />
			Cancel
		</button>
	</div>
</form>
