<script lang="ts">
	import { PhotoUploadModal, PhotoUploadButton } from '$lib/components';
	import type { TPhoto } from '$lib/types';

	let showModal = false;
	let uploadedPhotos: TPhoto[] = [];
	let selectedFiles: File[] = [];
	const demoAlbumId = 1;

	function handleUploadComplete(event: CustomEvent<TPhoto[]>) {
		uploadedPhotos = [...uploadedPhotos, ...event.detail];
		console.log('Upload complete:', event.detail);
	}

	function handleFilesSelected(event: CustomEvent<File[]>) {
		selectedFiles = event.detail;
		console.log('Files selected:', selectedFiles);
	}
</script>

<svelte:head>
	<title>Photo Upload Demo - Suipic</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<h1 class="text-4xl font-bold mb-2">Photo Upload Demo</h1>
		<p class="text-lg opacity-70 mb-8">
			Test the photo upload functionality with drag-and-drop, multiple file selection, EXIF display,
			and concurrent uploads.
		</p>

		<div class="space-y-6">
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">PhotoUploadModal Component</h2>
					<p class="opacity-70">
						Full-featured upload modal with drag-and-drop, previews, EXIF data, and progress tracking.
					</p>
					<div class="card-actions">
						<button class="btn btn-primary" on:click={() => (showModal = true)}>
							Open Upload Modal
						</button>
					</div>
				</div>
			</div>

			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">PhotoUploadButton Component</h2>
					<p class="opacity-70">
						Simple button component for file selection without the modal UI.
					</p>
					<div class="card-actions gap-2">
						<PhotoUploadButton
							label="Upload Photos"
							variant="primary"
							size="md"
							on:filesSelected={handleFilesSelected}
						/>
						<PhotoUploadButton
							label="Single File"
							variant="secondary"
							size="md"
							multiple={false}
							on:filesSelected={handleFilesSelected}
						/>
						<PhotoUploadButton
							label="Small"
							variant="outline"
							size="sm"
							on:filesSelected={handleFilesSelected}
						/>
					</div>
					{#if selectedFiles.length > 0}
						<div class="mt-4">
							<p class="font-semibold">Selected files:</p>
							<ul class="list-disc list-inside">
								{#each selectedFiles as file}
									<li class="text-sm">{file.name} ({(file.size / 1024 / 1024).toFixed(2)} MB)</li>
								{/each}
							</ul>
						</div>
					{/if}
				</div>
			</div>

			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">Uploaded Photos</h2>
					{#if uploadedPhotos.length === 0}
						<p class="opacity-70">No photos uploaded yet. Use the upload modal to test.</p>
					{:else}
						<div class="overflow-x-auto">
							<table class="table table-zebra">
								<thead>
									<tr>
										<th>ID</th>
										<th>Filename</th>
										<th>Date/Time</th>
										<th>Stars</th>
									</tr>
								</thead>
								<tbody>
									{#each uploadedPhotos as photo}
										<tr>
											<td>{photo.id}</td>
											<td>{photo.filename}</td>
											<td>{photo.dateTime || 'N/A'}</td>
											<td>{photo.stars}</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</div>
			</div>

			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">Features</h2>
					<ul class="list-disc list-inside space-y-2">
						<li>Drag-and-drop file upload</li>
						<li>Multiple file selection</li>
						<li>Client-side preview generation (300x300 thumbnails)</li>
						<li>EXIF metadata extraction and display</li>
						<li>Concurrent uploads with configurable limit (3-5 concurrent)</li>
						<li>Per-file progress indicators</li>
						<li>Upload status tracking (pending, uploading, complete, error)</li>
						<li>Retry failed uploads</li>
						<li>Supported formats: JPEG, PNG, WebP, GIF</li>
					</ul>
				</div>
			</div>

			<div class="alert alert-info">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					class="stroke-current shrink-0 w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					></path>
				</svg>
				<div>
					<p class="font-semibold">Note</p>
					<p class="text-sm">
						This demo page uses a hardcoded album ID ({demoAlbumId}). Make sure you're authenticated
						and have access to an album to test uploads.
					</p>
				</div>
			</div>
		</div>
	</div>
</div>

<PhotoUploadModal
	isOpen={showModal}
	albumId={demoAlbumId}
	on:close={() => (showModal = false)}
	on:uploadComplete={handleUploadComplete}
/>
