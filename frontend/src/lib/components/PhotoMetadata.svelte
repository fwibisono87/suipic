<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from '@iconify/svelte';
	import type { TPhoto } from '$lib/types';
	import { formatDateTime } from '$lib/utils/format';

	export let photo: TPhoto;
	export let editable: boolean = false;

	const dispatch = createEventDispatcher<{
		updateTitle: string | null;
	}>();

	let isEditingTitle = false;
	let titleInput = photo.title || '';
	let isUpdating = false;

	function startEditingTitle() {
		if (!editable) return;
		isEditingTitle = true;
		titleInput = photo.title || '';
		setTimeout(() => {
			const input = document.getElementById('metadata-title-input');
			if (input) (input as HTMLInputElement).focus();
		}, 0);
	}

	function saveTitle() {
		const newTitle = titleInput.trim();
		if (newTitle === (photo.title || '')) {
			isEditingTitle = false;
			return;
		}
		dispatch('updateTitle', newTitle || null);
		isEditingTitle = false;
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

<div class="photo-metadata bg-base-100 rounded-lg p-4 space-y-3">
	<div class="flex items-center gap-2 mb-4">
		<Icon icon="mdi:information" class="text-lg opacity-60" />
		<span class="text-sm font-semibold opacity-80">Photo Information</span>
	</div>

	<div class="space-y-2 text-sm">
		<div class="flex items-center gap-2">
			<Icon icon="mdi:file" class="opacity-60 flex-shrink-0" />
			<span class="font-medium">Filename:</span>
			<span class="flex-1 truncate opacity-80">{photo.filename}</span>
		</div>

		<div class="flex items-start gap-2">
			<Icon icon="mdi:text" class="opacity-60 mt-0.5 flex-shrink-0" />
			<span class="font-medium">Title:</span>
			{#if isEditingTitle}
				<div class="flex-1 flex gap-1">
					<input
						id="metadata-title-input"
						type="text"
						class="input input-xs input-bordered flex-1"
						bind:value={titleInput}
						on:keydown={handleTitleKeydown}
						disabled={isUpdating}
						placeholder="Enter title..."
					/>
					<button
						class="btn btn-xs btn-success"
						on:click={saveTitle}
						disabled={isUpdating}
					>
						<Icon icon="mdi:check" />
					</button>
					<button
						class="btn btn-xs btn-ghost"
						on:click={cancelEditingTitle}
						disabled={isUpdating}
					>
						<Icon icon="mdi:close" />
					</button>
				</div>
			{:else if editable}
				<button
					class="flex-1 text-left truncate opacity-80 hover:opacity-100 hover:underline"
					on:click={startEditingTitle}
				>
					{photo.title || 'Click to add title...'}
				</button>
			{:else}
				<span class="flex-1 truncate opacity-80">{photo.title || 'No title'}</span>
			{/if}
		</div>

		{#if photo.dateTime}
			<div class="flex items-center gap-2">
				<Icon icon="mdi:calendar" class="opacity-60 flex-shrink-0" />
				<span class="font-medium">Date Taken:</span>
				<span class="flex-1 opacity-80">{formatDateTime(photo.dateTime)}</span>
			</div>
		{/if}

		<div class="flex items-center gap-2">
			<Icon icon="mdi:clock" class="opacity-60 flex-shrink-0" />
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
					<span class="flex-1 opacity-80 break-words">{item.value}</span>
				</div>
			{/each}
		</div>
	{/if}
</div>
