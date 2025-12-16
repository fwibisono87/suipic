<script lang="ts">
	import Icon from '@iconify/svelte';
	import { StarRating, PickRejectButtons, PhotoMetadata, PhotoInteractionPanel } from '$lib/components';
	import type { TPhoto } from '$lib/types';

	let demoRating = 3;
	let demoPickRejectState: 'none' | 'pick' | 'reject' = 'none';

	const demoPhoto: TPhoto = {
		id: 1,
		albumId: 1,
		filename: 'demo-photo.jpg',
		title: 'Beautiful Landscape',
		description: 'A stunning view of the mountains at sunset',
		dateTime: '2024-01-15T18:30:00Z',
		exifData: {
			Make: 'Canon',
			Model: 'EOS R5',
			LensModel: 'RF 24-70mm F2.8 L IS USM',
			FocalLength: 35,
			FNumber: 2.8,
			ExposureTime: 0.008,
			ISO: 100,
			ImageWidth: 8192,
			ImageHeight: 5464
		},
		stars: 4,
		pickRejectState: 'pick',
		createdAt: '2024-01-15T20:00:00Z',
		updatedAt: '2024-01-15T20:00:00Z'
	};

	let interactivePhoto = { ...demoPhoto };

	function handleRatingChange(event: CustomEvent<number>) {
		demoRating = event.detail;
	}

	function handlePickRejectChange(event: CustomEvent<'none' | 'pick' | 'reject'>) {
		demoPickRejectState = event.detail;
	}

	function handlePhotoUpdate(event: CustomEvent<TPhoto>) {
		interactivePhoto = event.detail;
	}

	function handleTitleUpdate(event: CustomEvent<string | null>) {
		console.log('Title updated to:', event.detail);
	}
</script>

<svelte:head>
	<title>Photo Interaction Components Demo - Suipic</title>
</svelte:head>

<div class="max-w-7xl mx-auto p-4 sm:p-6 lg:p-8">
	<div class="mb-8">
		<h1 class="text-4xl font-bold mb-2">Photo Interaction Components</h1>
		<p class="text-lg opacity-70">
			Demo of the new photo interaction UI components including pick/reject buttons, star rating, and metadata display.
		</p>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
		<div class="space-y-8">
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title flex items-center gap-2">
						<Icon icon="mdi:star" class="text-2xl" />
						Star Rating Component
					</h2>
					<p class="text-sm opacity-70 mb-4">
						Interactive 5-star rating system with hover effects
					</p>
					
					<div class="space-y-4">
						<div>
							<p class="text-sm font-semibold mb-2">Default Size</p>
							<StarRating rating={demoRating} on:change={handleRatingChange} showClear={true} />
							<p class="text-xs opacity-60 mt-2">Current rating: {demoRating} stars</p>
						</div>

						<div class="divider"></div>

						<div>
							<p class="text-sm font-semibold mb-2">Small Size</p>
							<StarRating rating={3} size="sm" readonly={false} />
						</div>

						<div class="divider"></div>

						<div>
							<p class="text-sm font-semibold mb-2">Large Size (Read-only)</p>
							<StarRating rating={5} size="lg" readonly={true} />
						</div>
					</div>
				</div>
			</div>

			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title flex items-center gap-2">
						<Icon icon="mdi:flag" class="text-2xl" />
						Pick/Reject Buttons
					</h2>
					<p class="text-sm opacity-70 mb-4">
						Color-coded toggle buttons for photo selection workflow
					</p>
					
					<div class="space-y-4">
						<div>
							<p class="text-sm font-semibold mb-2">Default Size</p>
							<PickRejectButtons 
								state={demoPickRejectState} 
								on:change={handlePickRejectChange} 
							/>
							<p class="text-xs opacity-60 mt-2">Current state: {demoPickRejectState}</p>
						</div>

						<div class="divider"></div>

						<div>
							<p class="text-sm font-semibold mb-2">Small Size</p>
							<PickRejectButtons state="pick" size="sm" disabled={false} />
						</div>

						<div class="divider"></div>

						<div>
							<p class="text-sm font-semibold mb-2">Disabled State</p>
							<PickRejectButtons state="reject" disabled={true} />
						</div>
					</div>
				</div>
			</div>

			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title flex items-center gap-2">
						<Icon icon="mdi:information" class="text-2xl" />
						Photo Metadata Component
					</h2>
					<p class="text-sm opacity-70 mb-4">
						Display photo information with inline title editing
					</p>
					
					<PhotoMetadata 
						photo={demoPhoto} 
						editable={true} 
						on:updateTitle={handleTitleUpdate}
					/>
				</div>
			</div>
		</div>

		<div class="space-y-8">
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h2 class="card-title flex items-center gap-2">
						<Icon icon="mdi:image-edit" class="text-2xl" />
						Complete Photo Interaction Panel
					</h2>
					<p class="text-sm opacity-70 mb-4">
						All-in-one panel combining pick/reject, star rating, and metadata
					</p>
					
					<PhotoInteractionPanel 
						photo={interactivePhoto} 
						on:update={handlePhotoUpdate}
					/>
				</div>
			</div>

			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h3 class="card-title text-lg">
						<Icon icon="mdi:code-tags" class="text-xl" />
						Features
					</h3>
					<ul class="space-y-2 text-sm">
						<li class="flex items-start gap-2">
							<Icon icon="mdi:check-circle" class="text-success mt-0.5 flex-shrink-0" />
							<span><strong>Pick/Reject/Unflagged:</strong> Color-coded buttons (green/gray/red) for quick photo selection</span>
						</li>
						<li class="flex items-start gap-2">
							<Icon icon="mdi:check-circle" class="text-success mt-0.5 flex-shrink-0" />
							<span><strong>Star Rating:</strong> 0-5 stars with hover preview and toggle on/off</span>
						</li>
						<li class="flex items-start gap-2">
							<Icon icon="mdi:check-circle" class="text-success mt-0.5 flex-shrink-0" />
							<span><strong>Metadata Display:</strong> Filename, title, dates, and EXIF data</span>
						</li>
						<li class="flex items-start gap-2">
							<Icon icon="mdi:check-circle" class="text-success mt-0.5 flex-shrink-0" />
							<span><strong>Inline Title Editing:</strong> Click to edit, Enter to save, Escape to cancel</span>
						</li>
						<li class="flex items-start gap-2">
							<Icon icon="mdi:check-circle" class="text-success mt-0.5 flex-shrink-0" />
							<span><strong>EXIF Parsing:</strong> Camera, lens, aperture, shutter speed, ISO, dimensions</span>
						</li>
						<li class="flex items-start gap-2">
							<Icon icon="mdi:check-circle" class="text-success mt-0.5 flex-shrink-0" />
							<span><strong>Lightbox Integration:</strong> Press 'I' key to toggle info panel in lightbox view</span>
						</li>
					</ul>
				</div>
			</div>

			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h3 class="card-title text-lg">
						<Icon icon="mdi:keyboard" class="text-xl" />
						Keyboard Shortcuts (Lightbox)
					</h3>
					<div class="overflow-x-auto">
						<table class="table table-sm">
							<tbody>
								<tr>
									<td class="font-mono font-semibold">I</td>
									<td>Toggle info panel</td>
								</tr>
								<tr>
									<td class="font-mono font-semibold">←</td>
									<td>Previous photo</td>
								</tr>
								<tr>
									<td class="font-mono font-semibold">→</td>
									<td>Next photo</td>
								</tr>
								<tr>
									<td class="font-mono font-semibold">ESC</td>
									<td>Close panel or lightbox</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</div>
	</div>

	<div class="mt-8 text-center">
		<a href="/albums" class="btn btn-primary">
			<Icon icon="mdi:arrow-left" class="text-xl" />
			Back to Albums
		</a>
	</div>
</div>
