<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createEventDispatcher } from 'svelte';
	import { photosApi } from '$lib/api';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import Alert from './Alert.svelte';
	import type { TPhoto } from '$lib/types';
	import { extractExifData, type ExifDisplayData } from '$lib/utils/exif';
	import { generateImagePreview } from '$lib/utils/image';

	export let isOpen = false;
	export let albumId: number;

	const dispatch = createEventDispatcher<{
		close: void;
		uploadComplete: TPhoto[];
	}>();

	const MAX_CONCURRENT_UPLOADS = 3;
	const ACCEPTED_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'];

	interface FileWithPreview {
		file: File;
		preview: string;
		exif: ExifDisplayData | null;
		id: string;
		uploadProgress: number;
		status: 'pending' | 'uploading' | 'complete' | 'error';
		error?: string;
	}

	let files: FileWithPreview[] = [];
	let isDragging = false;
	let uploadQueue: FileWithPreview[] = [];
	let isUploading = false;
	let globalError = '';
	let completedUploads: TPhoto[] = [];

	$: pendingCount = files.filter((f) => f.status === 'pending').length;
	$: uploadingCount = files.filter((f) => f.status === 'uploading').length;
	$: completeCount = files.filter((f) => f.status === 'complete').length;
	$: errorCount = files.filter((f) => f.status === 'error').length;

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;

		const droppedFiles = Array.from(e.dataTransfer?.files || []);
		processFiles(droppedFiles);
	}

	function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		const selectedFiles = Array.from(input.files || []);
		processFiles(selectedFiles);
		input.value = '';
	}

	async function processFiles(newFiles: File[]) {
		const validFiles = newFiles.filter((file) => {
			if (!ACCEPTED_TYPES.includes(file.type)) {
				console.warn(`Skipping ${file.name}: unsupported type`);
				return false;
			}
			return true;
		});

		for (const file of validFiles) {
			const id = `${file.name}-${Date.now()}-${Math.random()}`;
			const preview = await generateImagePreview(file);
			const exif = await extractExifData(file);

			files = [
				...files,
				{
					file,
					preview,
					exif,
					id,
					uploadProgress: 0,
					status: 'pending'
				}
			];
		}
	}

	function removeFile(id: string) {
		const fileToRemove = files.find((f) => f.id === id);
		if (fileToRemove?.preview.startsWith('blob:')) {
			URL.revokeObjectURL(fileToRemove.preview);
		}
		files = files.filter((f) => f.id !== id);
	}

	async function startUpload() {
		if (isUploading) return;

		isUploading = true;
		globalError = '';
		uploadQueue = files.filter((f) => f.status === 'pending');

		while (uploadQueue.length > 0 || uploadingCount > 0) {
			const activeUploads = files.filter((f) => f.status === 'uploading').length;
			const slotsAvailable = MAX_CONCURRENT_UPLOADS - activeUploads;

			for (let i = 0; i < slotsAvailable && uploadQueue.length > 0; i++) {
				const fileItem = uploadQueue.shift();
				if (fileItem) {
					uploadFile(fileItem);
				}
			}

			await new Promise((resolve) => setTimeout(resolve, 100));
		}

		isUploading = false;

		if (errorCount === 0 && completeCount > 0) {
			dispatch('uploadComplete', completedUploads);
			handleClose();
		}
	}

	async function uploadFile(fileItem: FileWithPreview) {
		files = files.map((f) => (f.id === fileItem.id ? { ...f, status: 'uploading' } : f));

		try {
			const photo = await photosApi.create(albumId, fileItem.file);
			completedUploads = [...completedUploads, photo];

			files = files.map((f) =>
				f.id === fileItem.id ? { ...f, status: 'complete', uploadProgress: 100 } : f
			);
		} catch (err: unknown) {
			const errorMessage = (err as { message?: string }).message || 'Upload failed';
			files = files.map((f) =>
				f.id === fileItem.id ? { ...f, status: 'error', error: errorMessage } : f
			);
		}
	}

	function retryFailedUploads() {
		files = files.map((f) => (f.status === 'error' ? { ...f, status: 'pending', error: undefined } : f));
		startUpload();
	}

	function handleClose() {
		if (isUploading) {
			if (!confirm('Upload in progress. Are you sure you want to cancel?')) {
				return;
			}
		}

		files.forEach((f) => {
			if (f.preview.startsWith('blob:')) {
				URL.revokeObjectURL(f.preview);
			}
		});

		files = [];
		uploadQueue = [];
		isUploading = false;
		globalError = '';
		completedUploads = [];
		dispatch('close');
	}
</script>

{#if isOpen}
	<div class="modal modal-open" role="dialog">
		<div class="modal-box max-w-5xl max-h-[90vh] flex flex-col p-0">
			<div class="flex items-center justify-between p-6 border-b sticky top-0 bg-base-100 z-10">
				<h3 class="font-bold text-lg">Upload Photos</h3>
				<button class="btn btn-sm btn-circle btn-ghost" on:click={handleClose} disabled={isUploading}>
					<Icon icon="mdi:close" class="text-xl" />
				</button>
			</div>

			<div class="flex-1 overflow-y-auto p-6 space-y-4">
				{#if globalError}
					<Alert type="error" message={globalError} dismissible onDismiss={() => (globalError = '')} />
				{/if}

				{#if files.length === 0}
					<div
						class="border-2 border-dashed rounded-lg p-12 text-center transition-colors"
						class:border-primary={isDragging}
						class:bg-primary={isDragging}
						class:bg-opacity-10={isDragging}
						on:dragover={handleDragOver}
						on:dragleave={handleDragLeave}
						on:drop={handleDrop}
						role="button"
						tabindex="0"
					>
						<Icon icon="mdi:cloud-upload" class="text-6xl mx-auto mb-4 opacity-50" />
						<p class="text-xl mb-2">Drag and drop photos here</p>
						<p class="text-sm opacity-60 mb-4">or</p>
						<label class="btn btn-primary">
							<Icon icon="mdi:folder-open" class="text-xl" />
							Browse Files
							<input
								type="file"
								accept="image/jpeg,image/jpg,image/png,image/webp,image/gif"
								multiple
								class="hidden"
								on:change={handleFileSelect}
							/>
						</label>
						<p class="text-xs opacity-50 mt-4">Supported formats: JPEG, PNG, WebP, GIF</p>
					</div>
				{:else}
					<div
						class="border-2 border-dashed rounded-lg p-4 text-center transition-colors"
						class:border-primary={isDragging}
						class:bg-primary={isDragging}
						class:bg-opacity-10={isDragging}
						on:dragover={handleDragOver}
						on:dragleave={handleDragLeave}
						on:drop={handleDrop}
						role="button"
						tabindex="0"
					>
						<Icon icon="mdi:plus" class="text-3xl inline-block mr-2 opacity-50" />
						<span class="text-sm opacity-70">Drop more files or</span>
						<label class="btn btn-sm btn-ghost">
							browse
							<input
								type="file"
								accept="image/jpeg,image/jpg,image/png,image/webp,image/gif"
								multiple
								class="hidden"
								on:change={handleFileSelect}
							/>
						</label>
					</div>

					<div class="stats shadow w-full">
						<div class="stat">
							<div class="stat-title">Total</div>
							<div class="stat-value text-2xl">{files.length}</div>
						</div>
						<div class="stat">
							<div class="stat-title">Pending</div>
							<div class="stat-value text-2xl">{pendingCount}</div>
						</div>
						<div class="stat">
							<div class="stat-title">Uploading</div>
							<div class="stat-value text-2xl text-info">{uploadingCount}</div>
						</div>
						<div class="stat">
							<div class="stat-title">Complete</div>
							<div class="stat-value text-2xl text-success">{completeCount}</div>
						</div>
						{#if errorCount > 0}
							<div class="stat">
								<div class="stat-title">Failed</div>
								<div class="stat-value text-2xl text-error">{errorCount}</div>
							</div>
						{/if}
					</div>

					<div class="space-y-3 max-h-96 overflow-y-auto">
						{#each files as fileItem (fileItem.id)}
							<div class="card bg-base-200 shadow-sm">
								<div class="card-body p-4">
									<div class="flex gap-4">
										<div class="flex-shrink-0">
											<img
												src={fileItem.preview}
												alt={fileItem.file.name}
												class="w-24 h-24 object-cover rounded"
											/>
										</div>

										<div class="flex-1 min-w-0">
											<div class="flex items-start justify-between gap-2">
												<div class="flex-1 min-w-0">
													<h4 class="font-semibold truncate">{fileItem.file.name}</h4>
													<p class="text-xs opacity-60">
														{(fileItem.file.size / 1024 / 1024).toFixed(2)} MB
													</p>
												</div>

												<div class="flex items-center gap-2">
													{#if fileItem.status === 'pending'}
														<span class="badge badge-ghost">Pending</span>
													{:else if fileItem.status === 'uploading'}
														<span class="badge badge-info gap-1">
															<LoadingSpinner size="xs" />
															Uploading
														</span>
													{:else if fileItem.status === 'complete'}
														<span class="badge badge-success gap-1">
															<Icon icon="mdi:check" />
															Complete
														</span>
													{:else if fileItem.status === 'error'}
														<span class="badge badge-error gap-1">
															<Icon icon="mdi:alert" />
															Failed
														</span>
													{/if}

													{#if fileItem.status === 'pending' && !isUploading}
														<button
															class="btn btn-xs btn-ghost"
															on:click={() => removeFile(fileItem.id)}
														>
															<Icon icon="mdi:close" class="text-base" />
														</button>
													{/if}
												</div>
											</div>

											{#if fileItem.status === 'uploading'}
												<progress
													class="progress progress-info w-full mt-2"
													value={fileItem.uploadProgress}
													max="100"
												></progress>
											{/if}

											{#if fileItem.status === 'error' && fileItem.error}
												<p class="text-xs text-error mt-1">{fileItem.error}</p>
											{/if}

											{#if fileItem.exif}
												<div class="mt-2 flex flex-wrap gap-2 text-xs">
													{#if fileItem.exif.make && fileItem.exif.model}
														<span class="badge badge-sm badge-outline">
															<Icon icon="mdi:camera" class="mr-1" />
															{fileItem.exif.make} {fileItem.exif.model}
														</span>
													{/if}
													{#if fileItem.exif.lens}
														<span class="badge badge-sm badge-outline">
															{fileItem.exif.lens}
														</span>
													{/if}
													{#if fileItem.exif.focalLength}
														<span class="badge badge-sm badge-outline">
															{fileItem.exif.focalLength}
														</span>
													{/if}
													{#if fileItem.exif.aperture}
														<span class="badge badge-sm badge-outline">
															{fileItem.exif.aperture}
														</span>
													{/if}
													{#if fileItem.exif.shutterSpeed}
														<span class="badge badge-sm badge-outline">
															{fileItem.exif.shutterSpeed}
														</span>
													{/if}
													{#if fileItem.exif.iso}
														<span class="badge badge-sm badge-outline">
															{fileItem.exif.iso}
														</span>
													{/if}
													{#if fileItem.exif.dimensions}
														<span class="badge badge-sm badge-outline">
															<Icon icon="mdi:image-size-select-large" class="mr-1" />
															{fileItem.exif.dimensions}
														</span>
													{/if}
													{#if fileItem.exif.dateTime}
														<span class="badge badge-sm badge-outline">
															<Icon icon="mdi:calendar" class="mr-1" />
															{fileItem.exif.dateTime}
														</span>
													{/if}
												</div>
											{/if}
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<div class="border-t p-6 sticky bottom-0 bg-base-100 flex justify-end gap-2">
				{#if files.length > 0}
					{#if !isUploading}
						<button class="btn btn-ghost" on:click={handleClose}>Cancel</button>
						{#if errorCount > 0}
							<button class="btn btn-warning" on:click={retryFailedUploads}>
								<Icon icon="mdi:refresh" class="text-xl" />
								Retry Failed ({errorCount})
							</button>
						{/if}
						<button class="btn btn-primary" on:click={startUpload} disabled={pendingCount === 0}>
							<Icon icon="mdi:upload" class="text-xl" />
							Upload {pendingCount} {pendingCount === 1 ? 'Photo' : 'Photos'}
						</button>
					{:else}
						<div class="flex items-center gap-2">
							<LoadingSpinner size="sm" />
							<span>Uploading {uploadingCount} of {files.length}...</span>
						</div>
					{/if}
				{/if}
			</div>
		</div>
	</div>
{/if}
			</span>
													{/if}
													{#if fileItem.exif.dateTime}
														<span class="badge badge-sm badge-outline">
															<Icon icon="mdi:calendar" class="mr-1" />
															{fileItem.exif.dateTime}
														</span>
													{/if}
												</div>
											{/if}
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<div class="border-t p-6 sticky bottom-0 bg-base-100 flex justify-end gap-2">
				{#if files.length > 0}
					{#if !isUploading}
						<button class="btn btn-ghost" on:click={handleClose}>Cancel</button>
						{#if errorCount > 0}
							<button class="btn btn-warning" on:click={retryFailedUploads}>
								<Icon icon="mdi:refresh" class="text-xl" />
								Retry Failed ({errorCount})
							</button>
						{/if}
						<button class="btn btn-primary" on:click={startUpload} disabled={pendingCount === 0}>
							<Icon icon="mdi:upload" class="text-xl" />
							Upload {pendingCount} {pendingCount === 1 ? 'Photo' : 'Photos'}
						</button>
					{:else}
						<div class="flex items-center gap-2">
							<LoadingSpinner size="sm" />
							<span>Uploading {uploadingCount} of {files.length}...</span>
						</div>
					{/if}
				{/if}
			</div>
		</div>
	</div>
{/if}
