<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from '@iconify/svelte';
	import type { TPhoto } from '$lib/types';
	import { photosApi } from '$lib/api/photos';
	import { formatDateTime } from '$lib/utils/format';
	import { CommentSection } from '$lib/components';

	export let photo: TPhoto;
	export let photographerId: number | null = null;

	const dispatch = createEventDispatcher<{
		update: TPhoto;
	}>();

	let isUpdating = false;

	$: pickRejectState = photo.pickRejectState || 'none';
	$: stars = photo.stars || 0;

	async function updatePickRejectState(newState: 'none' | 'pick' | 'reject') {
		if (isUpdating) return;
		isUpdating = true;
		
		try {
			const updated = await photosApi.update(photo.id, { 
				pickRejectState: newState === 'none' ? null : newState 
			});
			photo = updated;
			dispatch('update', updated);
		} catch (err) {
			console.error('Failed to update pick/reject state:', err);
		} finally {
			isUpdating = false;
		}
	}

	async function updateStars(newStars: number) {
		if (isUpdating || newStars === stars) return;
		isUpdating = true;
		
		try {
			const updated = await photosApi.update(photo.id, { stars: newStars });
			photo = updated;
			dispatch('update', updated);
		} catch (err) {
			console.error('Failed to update stars:', err);
		} finally {
			isUpdating = false;
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

<div class="client-photo-interaction-panel bg-base-200 rounded-lg p-4 space-y-4">
	<div class="space-y-3">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:flag" class="text-lg opacity-60" />
			<span class="text-sm font-semibold opacity-80">Select Your Favorites</span>
		</div>
		
		<div class="btn-group w-full">
			<button
				class="btn btn-sm flex-1"
				class:btn-success={pickRejectState === 'pick'}
				class:btn-outline={pickRejectState !== 'pick'}
				disabled={isUpdating}
				on:click={() => updatePickRejectState(pickRejectState === 'pick' ? 'none' : 'pick')}
			>
				<Icon icon="mdi:check-circle" class="text-lg" />
				Pick
			</button>
			<button
				class="btn btn-sm flex-1"
				class:btn-ghost={pickRejectState === 'none'}
				class:btn-outline={pickRejectState === 'none'}
				disabled={isUpdating || pickRejectState === 'none'}
				on:click={() => updatePickRejectState('none')}
			>
				<Icon icon="mdi:flag-outline" class="text-lg" />
				Unflagged
			</button>
			<button
				class="btn btn-sm flex-1"
				class:btn-error={pickRejectState === 'reject'}
				class:btn-outline={pickRejectState !== 'reject'}
				disabled={isUpdating}
				on:click={() => updatePickRejectState(pickRejectState === 'reject' ? 'none' : 'reject')}
			>
				<Icon icon="mdi:close-circle" class="text-lg" />
				Reject
			</button>
		</div>
		<p class="text-xs opacity-60">Mark your favorite photos or ones you'd like to exclude</p>
	</div>

	<div class="divider my-2"></div>

	<div class="space-y-3">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:star" class="text-lg opacity-60" />
			<span class="text-sm font-semibold opacity-80">Rate This Photo</span>
		</div>
		
		<div class="flex items-center gap-1">
			{#each [1, 2, 3, 4, 5] as rating}
				<button
					class="btn btn-ghost btn-sm p-1 min-h-0 h-8"
					disabled={isUpdating}
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
					disabled={isUpdating}
					on:click={() => updateStars(0)}
				>
					Clear
				</button>
			{/if}
		</div>
		<p class="text-xs opacity-60">Give a star rating from 1 to 5</p>
	</div>

	<div class="divider my-2"></div>

	<div class="space-y-3">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:information" class="text-lg opacity-60" />
			<span class="text-sm font-semibold opacity-80">Photo Details</span>
		</div>

		<div class="space-y-2 text-sm">
			<div class="flex items-center gap-2">
				<Icon icon="mdi:file" class="opacity-60" />
				<span class="font-medium">Filename:</span>
				<span class="flex-1 truncate opacity-80">{photo.filename}</span>
			</div>

			{#if photo.title}
				<div class="flex items-start gap-2">
					<Icon icon="mdi:text" class="opacity-60 mt-0.5" />
					<span class="font-medium">Title:</span>
					<span class="flex-1 opacity-80">{photo.title}</span>
				</div>
			{/if}

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
					<span class="font-medium">Camera Info</span>
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
	.client-photo-interaction-panel {
		max-height: calc(100vh - 200px);
		overflow-y: auto;
	}
</style>
