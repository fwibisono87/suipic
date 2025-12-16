<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from '@iconify/svelte';
	import type { TPhoto } from '$lib/types';
	import { formatDateTime } from '$lib/utils/format';
	import { CommentSection } from '$lib/components';
	import { useUpdatePhotoMutation } from '$lib/queries/photos';

	export let photo: TPhoto;
	export let photographerId: number | null = null;

	const dispatch = createEventDispatcher<{
		update: TPhoto;
	}>();

	const updatePhotoMutation = useUpdatePhotoMutation();

	let isEditingTitle = false;
	let titleInput = photo.title || '';

	$: pickRejectState = photo.pickRejectState || 'none';
	$: stars = photo.stars || 0;

	function updatePickRejectState(newState: 'none' | 'pick' | 'reject') {
		$updatePhotoMutation.mutate(
			{ 
				photoId: photo.id, 
				updates: { pickRejectState: newState === 'none' ? null : newState } 
			},
			{
				onSuccess: (updated) => {
					photo = updated;
					dispatch('update', updated);
				},
				onError: (err) => {
					console.error('Failed to update pick/reject state:', err);
				}
			}
		);
	}

	function updateStars(newStars: number) {
		if (newStars === stars) return;
		
		$updatePhotoMutation.mutate(
			{ 
				photoId: photo.id, 
				updates: { stars: newStars } 
			},
			{
				onSuccess: (updated) => {
					photo = updated;
					dispatch('update', updated);
				},
				onError: (err) => {
					console.error('Failed to update stars:', err);
				}
			}
		);
	}

	function startEditingTitle() {
		isEditingTitle = true;
		titleInput = photo.title || '';
		setTimeout(() => {
			const input = document.getElementById('title-input');
			if (input) (input as HTMLInputElement).focus();
		}, 0);
	}

	function saveTitle() {
		const newTitle = titleInput.trim();
		if (newTitle === (photo.title || '')) {
			isEditingTitle = false;
			return;
		}

		$updatePhotoMutation.mutate(
			{ 
				photoId: photo.id, 
				updates: { title: newTitle || null } 
			},
			{
				onSuccess: (updated) => {
					photo = updated;
					dispatch('update', updated);
					isEditingTitle = false;
				},
				onError: (err) => {
					console.error('Failed to update title:', err);
				}
			}
		);
	}

	function cancelEditingTitle() {
		isEditingTitle = false;
		titleInput = photo.title || '';
	}

	function handleTitleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			saveTitle();
		} else if (event.key === 'Escape') {
			event.preventDefault();
			cancelEditingTitle();
		}
	}

	function getExifDisplay(): Array<{ label: string; value: string }> {
		const items: Array<{ label: string; value: string }> = [];
		
		if (!photo.exifData) return items;

		const data = photo.exifData as Record<string, unknown>;

		if (data.Make && data.Model) {
			items.push({ label: 'Camera', value: `${data.Make} ${data.Model}` });
		}
		if (data.LensModel) {
			items.push({ label: 'Lens', value: String(data.LensModel) });
		}
		if (data.FocalLength) {
			items.push({ label: 'Focal Length', value: `${data.FocalLength}mm` });
		}
		if (data.FNumber) {
			items.push({ label: 'Aperture', value: `f/${data.FNumber}` });
		}
		if (data.ExposureTime) {
			const exp = Number(data.ExposureTime);
			const value = exp < 1 ? `1/${Math.round(1 / exp)}s` : `${exp}s`;
			items.push({ label: 'Shutter Speed', value });
		}
		if (data.ISO) {
			items.push({ label: 'ISO', value: `ISO ${data.ISO}` });
		}
		if (data.ImageWidth && data.ImageHeight) {
			items.push({ label: 'Dimensions', value: `${data.ImageWidth} × ${data.ImageHeight}` });
		}

		return items;
	}

	$: exifItems = getExifDisplay();
</script>

<div class="photo-interaction-panel bg-base-200 rounded-lg p-4 space-y-4">
	<div class="space-y-3">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:flag" class="text-lg opacity-60" />
			<span class="text-sm font-semibold opacity-80">Pick / Reject</span>
		</div>
		
		<div class="btn-group w-full">
			<button
				class="btn btn-sm flex-1"
				class:btn-success={pickRejectState === 'pick'}
				class:btn-outline={pickRejectState !== 'pick'}
				on:click={() => updatePickRejectState(pickRejectState === 'pick' ? 'none' : 'pick')}
			>
				<Icon icon="mdi:check-circle" class="text-lg" />
				Pick
			</button>
			<button
				class="btn btn-sm flex-1"
				class:btn-ghost={pickRejectState === 'none'}
				class:btn-outline={pickRejectState === 'none'}
				disabled={pickRejectState === 'none'}
				on:click={() => updatePickRejectState('none')}
			>
				<Icon icon="mdi:flag-outline" class="text-lg" />
				Unflagged
			</button>
			<button
				class="btn btn-sm flex-1"
				class:btn-error={pickRejectState === 'reject'}
				class:btn-outline={pickRejectState !== 'reject'}
				on:click={() => updatePickRejectState(pickRejectState === 'reject' ? 'none' : 'reject')}
			>
				<Icon icon="mdi:close-circle" class="text-lg" />
				Reject
			</button>
		</div>
	</div>

	<div class="divider my-2"></div>

	<div class="space-y-3">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:star" class="text-lg opacity-60" />
			<span class="text-sm font-semibold opacity-80">Rating</span>
		</div>
		
		<div class="flex items-center gap-1">
			{#each [1, 2, 3, 4, 5] as rating}
				<button
					class="btn btn-ghost btn-sm p-1 min-h-0 h-8"
					on:click={() => updateStars(rating === stars ? 0 : rating)}
				>
					<Icon
						icon={rating <= stars ? 'mdi:star' : 'mdi:star-outline'}
						class={`text-2xl ${rating <= stars ? 'text-warning' : ''} ${
							rating > stars ? 'opacity-30' : ''
						}`.trim()}
					/>
				</button>
			{/each}
			{#if stars > 0}
				<button
					class="btn btn-ghost btn-xs ml-2"
					on:click={() => updateStars(0)}
				>
					Clear
				</button>
			{/if}
		</div>
	</div>

	<div class="divider my-2"></div>

	<div class="space-y-3">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:information" class="text-lg opacity-60" />
			<span class="text-sm font-semibold opacity-80">Metadata</span>
		</div>

		<div class="space-y-2 text-sm">
			<div class="flex items-center gap-2">
				<Icon icon="mdi:file" class="opacity-60" />
				<span class="font-medium">Filename:</span>
				<span class="flex-1 truncate opacity-80">{photo.filename}</span>
			</div>

			<div class="flex items-start gap-2">
				<Icon icon="mdi:text" class="opacity-60 mt-0.5" />
				<span class="font-medium">Title:</span>
				{#if isEditingTitle}
					<div class="flex-1 flex gap-1">
						<input
							id="title-input"
							type="text"
							class="input input-xs input-bordered flex-1"
							bind:value={titleInput}
							on:keydown={handleTitleKeydown}
							disabled={$updatePhotoMutation.isPending}
							placeholder="Enter title..."
						/>
						<button
							class="btn btn-xs btn-success"
							on:click={saveTitle}
							disabled={$updatePhotoMutation.isPending}
						>
							<Icon icon="mdi:check" />
						</button>
						<button
							class="btn btn-xs btn-ghost"
							on:click={cancelEditingTitle}
							disabled={$updatePhotoMutation.isPending}
						>
							<Icon icon="mdi:close" />
						</button>
					</div>
				{:else}
					<button
						class="flex-1 text-left truncate opacity-80 hover:opacity-100 hover:underline"
						on:click={startEditingTitle}
					>
						{photo.title || 'Click to add title...'}
					</button>
				{/if}
			</div>

			{#if photo.dateTime}
				<div class="flex items-center gap-2">
					<Icon icon="mdi:calendar" class="opacity-60" />
					<span class="font-medium">Date Taken:</span>
					<span class="flex-1 opacity-80">{formatDateTime(photo.dateTime)}</span>
				</div>
			{/if}

			<div class="flex items-center gap-2">
				<Icon icon="mdi:clock" class="opacity-60" />
				<span class="font-medium">Uploaded:</span>
				<span class="flex-1 opacity-80">{formatDateTime(photo.createdAt)}</span>
			</div>
		</div>

		{#if exifItems.length > 0}
			<div class="divider my-2"></div>
			
			<div class="space-y-2 text-sm">
				<div class="flex items-center gap-2 mb-2">
					<Icon icon="mdi:camera" class="opacity-60" />
					<span class="font-medium">EXIF Data</span>
				</div>
				
				{#each exifItems as item}
					<div class="flex items-start gap-2 pl-6">
						<span class="font-medium min-w-[100px]">{item.label}:</span>
						<span class="flex-1 opacity-80">{item.value}</span>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	{#if photographerId !== null}
		<div class="divider my-2"></div>
		
		<div class="mt-4">
			<CommentSection {photo} {photographerId} />
		</div>
	{/if}
</div>

<style>
	.photo-interaction-panel {
		max-height: calc(100vh - 200px);
		overflow-y: auto;
	}
</style>
