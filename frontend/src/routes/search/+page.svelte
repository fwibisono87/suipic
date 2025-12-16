<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';
	import { isAuthenticated } from '$lib/stores';
	import { LoadingSpinner, Alert, PhotoGallery } from '$lib/components';
	import { searchApi, albumsApi } from '$lib/api';
	import type { TSearchParams, TSearchResponse } from '$lib/api/search';
	import type { TAlbum } from '$lib/types';
	
	let searchQuery = '';
	let selectedAlbum: number | undefined;
	let dateFrom = '';
	let dateTo = '';
	let minStars = 0;
	let maxStars = 5;
	let selectedStates: string[] = [];
	let currentPage = 1;
	let pageSize = 50;
	let showFilters = false;
	
	let searchResults: TSearchResponse | null = null;
	let albums: TAlbum[] = [];
	let isSearching = false;
	let searchError = '';
	let isLoadingAlbums = false;

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}
		
		loadAlbums();
		
		const params = $page.url.searchParams;
		if (params.get('q')) {
			searchQuery = params.get('q') || '';
			performSearch();
		}
	});

	async function loadAlbums() {
		isLoadingAlbums = true;
		try {
			albums = await albumsApi.list();
		} catch (err: unknown) {
			console.error('Failed to load albums:', err);
		} finally {
			isLoadingAlbums = false;
		}
	}

	async function performSearch() {
		isSearching = true;
		searchError = '';
		
		try {
			const params: TSearchParams = {
				limit: pageSize,
				offset: (currentPage - 1) * pageSize
			};

			if (searchQuery.trim()) {
				params.q = searchQuery.trim();
			}
			
			if (selectedAlbum) {
				params.album = selectedAlbum;
			}
			
			if (dateFrom) {
				params.dateFrom = new Date(dateFrom).toISOString();
			}
			
			if (dateTo) {
				const date = new Date(dateTo);
				date.setHours(23, 59, 59, 999);
				params.dateTo = date.toISOString();
			}
			
			if (minStars > 0) {
				params.minStars = minStars;
			}
			
			if (maxStars < 5) {
				params.maxStars = maxStars;
			}
			
			if (selectedStates.length > 0 && selectedStates.length < 3) {
				params.state = selectedStates.join(',');
			}

			searchResults = await searchApi.search(params);
		} catch (err: unknown) {
			searchError = (err as { message: string }).message || 'Search failed';
			searchResults = null;
		} finally {
			isSearching = false;
		}
	}

	function handleSearch() {
		currentPage = 1;
		performSearch();
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			handleSearch();
		}
	}

	function handleStateToggle(state: string) {
		if (selectedStates.includes(state)) {
			selectedStates = selectedStates.filter(s => s !== state);
		} else {
			selectedStates = [...selectedStates, state];
		}
	}

	function clearFilters() {
		selectedAlbum = undefined;
		dateFrom = '';
		dateTo = '';
		minStars = 0;
		maxStars = 5;
		selectedStates = [];
		currentPage = 1;
		performSearch();
	}

	function handlePreviousPage() {
		if (currentPage > 1) {
			currentPage--;
			performSearch();
			window.scrollTo({ top: 0, behavior: 'smooth' });
		}
	}

	function handleNextPage() {
		if (searchResults && currentPage * pageSize < searchResults.total) {
			currentPage++;
			performSearch();
			window.scrollTo({ top: 0, behavior: 'smooth' });
		}
	}

	$: totalPages = searchResults ? Math.ceil(searchResults.total / pageSize) : 0;
	$: hasActiveFilters = selectedAlbum || dateFrom || dateTo || minStars > 0 || maxStars < 5 || selectedStates.length > 0;
</script>

<svelte:head>
	<title>Search Photos - Suipic</title>
</svelte:head>

{#if $isAuthenticated}
	<div class="space-y-6">
		<div class="flex items-start justify-between gap-4">
			<div>
				<h1 class="text-4xl font-bold">Search Photos</h1>
				<p class="text-lg opacity-70 mt-2">Search across all photos using ElasticSearch</p>
			</div>
		</div>

		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<div class="flex gap-2">
					<div class="flex-1">
						<div class="form-control">
							<div class="input-group">
								<input
									type="text"
									placeholder="Search photos by title, album, location, comments, or EXIF data..."
									class="input input-bordered w-full"
									bind:value={searchQuery}
									on:keydown={handleKeyDown}
								/>
								<button class="btn btn-primary" on:click={handleSearch} disabled={isSearching}>
									{#if isSearching}
										<LoadingSpinner size="sm" />
									{:else}
										<Icon icon="mdi:magnify" class="text-xl" />
									{/if}
									Search
								</button>
							</div>
						</div>
					</div>
					<button
						class="btn btn-outline"
						class:btn-active={showFilters}
						on:click={() => (showFilters = !showFilters)}
					>
						<Icon icon="mdi:filter-variant" class="text-xl" />
						<span class="hidden sm:inline">Filters</span>
						{#if hasActiveFilters}
							<span class="badge badge-sm badge-primary">Active</span>
						{/if}
					</button>
				</div>

				{#if showFilters}
					<div class="divider"></div>
					<div class="space-y-4">
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div class="form-control">
								<label class="label" for="album-filter">
									<span class="label-text font-semibold">Album</span>
								</label>
								<select
									id="album-filter"
									class="select select-bordered w-full"
									bind:value={selectedAlbum}
									disabled={isLoadingAlbums}
								>
									<option value={undefined}>All Albums</option>
									{#each albums as album}
										<option value={album.id}>{album.title}</option>
									{/each}
								</select>
							</div>

							<div class="form-control">
								<label class="label">
									<span class="label-text font-semibold">Star Rating</span>
									<span class="label-text-alt">{minStars} - {maxStars} stars</span>
								</label>
								<div class="flex items-center gap-4">
									<div class="flex-1">
										<label class="label-text text-xs">Min</label>
										<input
											type="range"
											min="0"
											max="5"
											bind:value={minStars}
											class="range range-xs range-primary"
										/>
										<div class="flex justify-between text-xs opacity-60 mt-1">
											<span>0</span>
											<span>5</span>
										</div>
									</div>
									<div class="flex-1">
										<label class="label-text text-xs">Max</label>
										<input
											type="range"
											min="0"
											max="5"
											bind:value={maxStars}
											class="range range-xs range-primary"
										/>
										<div class="flex justify-between text-xs opacity-60 mt-1">
											<span>0</span>
											<span>5</span>
										</div>
									</div>
								</div>
							</div>
						</div>

						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div class="form-control">
								<label class="label" for="date-from">
									<span class="label-text font-semibold">Date From</span>
								</label>
								<input
									id="date-from"
									type="date"
									class="input input-bordered w-full"
									bind:value={dateFrom}
								/>
							</div>

							<div class="form-control">
								<label class="label" for="date-to">
									<span class="label-text font-semibold">Date To</span>
								</label>
								<input
									id="date-to"
									type="date"
									class="input input-bordered w-full"
									bind:value={dateTo}
								/>
							</div>
						</div>

						<div class="form-control">
							<label class="label">
								<span class="label-text font-semibold">Pick/Reject State</span>
							</label>
							<div class="flex flex-wrap gap-2">
								<button
									class="btn btn-sm"
									class:btn-outline={!selectedStates.includes('none')}
									class:btn-neutral={selectedStates.includes('none')}
									on:click={() => handleStateToggle('none')}
								>
									<Icon icon="mdi:minus-circle-outline" class="text-lg" />
									None
								</button>
								<button
									class="btn btn-sm"
									class:btn-outline={!selectedStates.includes('pick')}
									class:btn-success={selectedStates.includes('pick')}
									on:click={() => handleStateToggle('pick')}
								>
									<Icon icon="mdi:check-circle" class="text-lg" />
									Pick
								</button>
								<button
									class="btn btn-sm"
									class:btn-outline={!selectedStates.includes('reject')}
									class:btn-error={selectedStates.includes('reject')}
									on:click={() => handleStateToggle('reject')}
								>
									<Icon icon="mdi:close-circle" class="text-lg" />
									Reject
								</button>
							</div>
						</div>

						<div class="flex justify-between items-center pt-2">
							<button class="btn btn-ghost btn-sm" on:click={clearFilters} disabled={!hasActiveFilters}>
								<Icon icon="mdi:filter-remove" class="text-lg" />
								Clear Filters
							</button>
							<button class="btn btn-primary btn-sm" on:click={handleSearch} disabled={isSearching}>
								Apply Filters
							</button>
						</div>
					</div>
				{/if}
			</div>
		</div>

		{#if searchError}
			<Alert type="error" message={searchError} dismissible onDismiss={() => (searchError = '')} />
		{/if}

		{#if isSearching && !searchResults}
			<div class="flex justify-center py-20">
				<LoadingSpinner size="lg" />
			</div>
		{:else if searchResults}
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div>
						<h2 class="text-2xl font-bold">
							{searchResults.total} {searchResults.total === 1 ? 'Result' : 'Results'}
						</h2>
						{#if searchQuery}
							<p class="text-sm opacity-60 mt-1">
								Searching for "{searchQuery}"
								{#if hasActiveFilters}
									<span class="badge badge-sm badge-primary ml-2">Filtered</span>
								{/if}
							</p>
						{/if}
					</div>
					
					{#if totalPages > 1}
						<div class="join">
							<button
								class="join-item btn btn-sm"
								on:click={handlePreviousPage}
								disabled={currentPage === 1 || isSearching}
							>
								<Icon icon="mdi:chevron-left" class="text-lg" />
							</button>
							<button class="join-item btn btn-sm no-animation">
								Page {currentPage} of {totalPages}
							</button>
							<button
								class="join-item btn btn-sm"
								on:click={handleNextPage}
								disabled={currentPage >= totalPages || isSearching}
							>
								<Icon icon="mdi:chevron-right" class="text-lg" />
							</button>
						</div>
					{/if}
				</div>

				{#if searchResults.photos.length > 0}
					<PhotoGallery photos={searchResults.photos} layout="grid" onPhotoUpdate={() => performSearch()} photographerId={null} />
					
					{#if totalPages > 1}
						<div class="flex justify-center pt-4">
							<div class="join">
								<button
									class="join-item btn"
									on:click={handlePreviousPage}
									disabled={currentPage === 1 || isSearching}
								>
									<Icon icon="mdi:chevron-left" class="text-lg" />
									Previous
								</button>
								<button class="join-item btn no-animation">
									Page {currentPage} of {totalPages}
								</button>
								<button
									class="join-item btn"
									on:click={handleNextPage}
									disabled={currentPage >= totalPages || isSearching}
								>
									Next
									<Icon icon="mdi:chevron-right" class="text-lg" />
								</button>
							</div>
						</div>
					{/if}
				{:else}
					<div class="text-center py-20">
						<Icon icon="mdi:image-search" class="text-8xl mx-auto opacity-30" />
						<p class="text-xl opacity-60 mt-4">No photos found</p>
						<p class="text-sm opacity-50 mt-2">Try adjusting your search or filters</p>
					</div>
				{/if}
			</div>
		{:else}
			<div class="text-center py-20">
				<Icon icon="mdi:image-search-outline" class="text-8xl mx-auto opacity-30" />
				<p class="text-xl opacity-60 mt-4">Start searching to find photos</p>
				<p class="text-sm opacity-50 mt-2">Use the search box above or apply filters to begin</p>
			</div>
		{/if}
	</div>
{:else}
	<div class="flex items-center justify-center min-h-[60vh]">
		<LoadingSpinner size="lg" />
	</div>
{/if}
