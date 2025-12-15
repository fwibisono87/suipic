<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { isAuthenticated, currentUser, authToken } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { Card, Alert, LoadingSpinner } from '$lib/components';
	import { photographerApi, type TCreateClientRequest, type TClient } from '$lib/api';
	import { validateEmail, validateUsername } from '$lib/utils';

	let clients: TClient[] = [];
	let isLoadingList = true;
	let listError = '';

	let friendlyName = '';
	let username = '';
	let email = '';
	let password = '';
	let isCreating = false;
	let createError = '';
	let createSuccess = false;
	let showCreateForm = true;

	let searchQuery = '';
	let searchResults: TClient[] = [];
	let isSearching = false;
	let searchError = '';
	let showSearchResults = false;
	let searchTimeout: NodeJS.Timeout | null = null;

	let currentPage = 1;
	let itemsPerPage = 10;

	$: paginatedClients = clients.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
	$: totalPages = Math.ceil(clients.length / itemsPerPage);
	$: hasNextPage = currentPage < totalPages;
	$: hasPrevPage = currentPage > 1;

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if ($currentUser?.role !== EUserRole.PHOTOGRAPHER) {
			goto('/');
			return;
		}

		loadClients();
	});

	async function loadClients() {
		if (!$authToken) return;

		isLoadingList = true;
		listError = '';

		try {
			clients = await photographerApi.listClients($authToken);
		} catch (err: unknown) {
			listError = (err as { message: string }).message || 'Failed to load clients';
		} finally {
			isLoadingList = false;
		}
	}

	async function handleCreateClient(e: Event) {
		e.preventDefault();
		createError = '';
		createSuccess = false;

		if (!username || !validateUsername(username)) {
			createError = 'Username must be between 3 and 50 characters';
			return;
		}

		if (!email || !validateEmail(email)) {
			createError = 'Please enter a valid email address';
			return;
		}

		if (!password || password.length < 6) {
			createError = 'Password must be at least 6 characters';
			return;
		}

		if (!$authToken) {
			createError = 'Not authenticated';
			return;
		}

		isCreating = true;

		try {
			const data: TCreateClientRequest = {
				username,
				email,
				password,
				friendlyName: friendlyName || undefined
			};

			await photographerApi.createOrLinkClient(data, $authToken);
			createSuccess = true;
			friendlyName = '';
			username = '';
			email = '';
			password = '';

			await loadClients();
		} catch (err: unknown) {
			createError = (err as { message: string }).message || 'Failed to create client';
		} finally {
			isCreating = false;
		}
	}

	async function handleSearch() {
		if (!searchQuery.trim()) {
			showSearchResults = false;
			searchResults = [];
			return;
		}

		if (!$authToken) return;

		isSearching = true;
		searchError = '';

		try {
			searchResults = await photographerApi.searchClients(searchQuery, $authToken);
			showSearchResults = true;
		} catch (err: unknown) {
			searchError = (err as { message: string }).message || 'Failed to search clients';
		} finally {
			isSearching = false;
		}
	}

	function onSearchInput() {
		if (searchTimeout) {
			clearTimeout(searchTimeout);
		}

		searchTimeout = setTimeout(() => {
			handleSearch();
		}, 300);
	}

	async function handleLinkClient(client: TClient) {
		if (!$authToken) return;

		isCreating = true;
		createError = '';

		try {
			const data: TCreateClientRequest = {
				username: client.username
			};

			await photographerApi.createOrLinkClient(data, $authToken);
			showSearchResults = false;
			searchQuery = '';
			searchResults = [];

			await loadClients();
		} catch (err: unknown) {
			createError = (err as { message: string }).message || 'Failed to link client';
		} finally {
			isCreating = false;
		}
	}

	function toggleFormMode() {
		showCreateForm = !showCreateForm;
		createError = '';
		searchError = '';
		createSuccess = false;
		if (showCreateForm) {
			showSearchResults = false;
			searchQuery = '';
		}
	}

	function goToPage(page: number) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
		}
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Client Management - Suipic</title>
</svelte:head>

{#if $isAuthenticated && $currentUser?.role === EUserRole.PHOTOGRAPHER}
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<h1 class="text-4xl font-bold">Client Management</h1>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<div class="lg:col-span-1">
				<Card title="Manage Clients">
					<div class="flex gap-2 mb-4">
						<button
							class="btn btn-sm flex-1 {showCreateForm ? 'btn-primary' : ''}"
							on:click={toggleFormMode}
							disabled={showCreateForm}
						>
							<Icon icon="mdi:account-plus" class="text-lg" />
							Create New
						</button>
						<button
							class="btn btn-sm flex-1 {!showCreateForm ? 'btn-primary' : ''}"
							on:click={toggleFormMode}
							disabled={!showCreateForm}
						>
							<Icon icon="mdi:link" class="text-lg" />
							Link Existing
						</button>
					</div>

					{#if createError}
						<div class="mb-4">
							<Alert type="error" message={createError} dismissible onDismiss={() => (createError = '')} />
						</div>
					{/if}

					{#if searchError}
						<div class="mb-4">
							<Alert type="error" message={searchError} dismissible onDismiss={() => (searchError = '')} />
						</div>
					{/if}

					{#if createSuccess}
						<div class="mb-4">
							<Alert type="success" message="Client created successfully!" />
						</div>
					{/if}

					{#if showCreateForm}
						<form on:submit={handleCreateClient} class="space-y-4">
							<div class="form-control">
								<label class="label" for="friendlyName">
									<span class="label-text">Friendly Name (Optional)</span>
								</label>
								<input
									type="text"
									id="friendlyName"
									name="friendlyName"
									bind:value={friendlyName}
									placeholder="John Doe"
									class="input input-bordered w-full"
									disabled={isCreating}
								/>
							</div>

							<div class="form-control">
								<label class="label" for="username">
									<span class="label-text">Username</span>
								</label>
								<input
									type="text"
									id="username"
									name="username"
									bind:value={username}
									placeholder="client_username"
									class="input input-bordered w-full"
									disabled={isCreating}
									required
								/>
								<label class="label">
									<span class="label-text-alt">3-50 characters</span>
								</label>
							</div>

							<div class="form-control">
								<label class="label" for="email">
									<span class="label-text">Email</span>
								</label>
								<input
									type="email"
									id="email"
									name="email"
									bind:value={email}
									placeholder="client@example.com"
									class="input input-bordered w-full"
									disabled={isCreating}
									required
								/>
							</div>

							<div class="form-control">
								<label class="label" for="password">
									<span class="label-text">Password</span>
								</label>
								<input
									type="password"
									id="password"
									name="password"
									bind:value={password}
									placeholder="••••••••"
									class="input input-bordered w-full"
									disabled={isCreating}
									required
								/>
								<label class="label">
									<span class="label-text-alt">At least 6 characters</span>
								</label>
							</div>

							<button type="submit" class="btn btn-primary w-full" disabled={isCreating}>
								{#if isCreating}
									<LoadingSpinner size="sm" />
								{:else}
									<Icon icon="mdi:account-plus" class="text-lg" />
									Create Client
								{/if}
							</button>
						</form>
					{:else}
						<div class="space-y-4">
							<div class="form-control">
								<label class="label" for="search">
									<span class="label-text">Search by Username</span>
								</label>
								<input
									type="text"
									id="search"
									name="search"
									bind:value={searchQuery}
									on:input={onSearchInput}
									placeholder="Type to search..."
									class="input input-bordered w-full"
									disabled={isSearching}
								/>
							</div>

							{#if isSearching}
								<div class="flex justify-center py-4">
									<LoadingSpinner size="sm" />
								</div>
							{:else if showSearchResults}
								{#if searchResults.length === 0}
									<div class="text-center py-4 opacity-60">
										<Icon icon="mdi:account-search" class="text-4xl mb-2" />
										<p>No clients found</p>
									</div>
								{:else}
									<div class="space-y-2">
										{#each searchResults as client}
											<div class="card bg-base-200 shadow-sm">
												<div class="card-body p-4">
													<div class="flex items-center justify-between">
														<div>
															<div class="font-semibold">{client.username}</div>
															{#if client.friendlyName}
																<div class="text-sm opacity-70">{client.friendlyName}</div>
															{/if}
															<div class="text-xs opacity-60">{client.email}</div>
														</div>
														<button
															class="btn btn-sm btn-primary"
															on:click={() => handleLinkClient(client)}
															disabled={isCreating}
														>
															<Icon icon="mdi:link" />
															Link
														</button>
													</div>
												</div>
											</div>
										{/each}
									</div>
								{/if}
							{/if}
						</div>
					{/if}

					<div class="mt-4 p-4 bg-info/10 rounded-lg">
						<p class="text-sm">
							<Icon icon="mdi:information" class="inline text-info" />
							{#if showCreateForm}
								Create a new client account with email and password.
							{:else}
								Link an existing client by searching for their username.
							{/if}
						</p>
					</div>
				</Card>
			</div>

			<div class="lg:col-span-2">
				<Card title="My Clients">
					{#if listError}
						<div class="mb-4">
							<Alert type="error" message={listError} />
						</div>
					{/if}

					{#if isLoadingList}
						<div class="flex justify-center py-8">
							<LoadingSpinner />
						</div>
					{:else if clients.length === 0}
						<div class="text-center py-8 opacity-60">
							<Icon icon="mdi:account-group" class="text-6xl mb-2" />
							<p>No clients yet</p>
						</div>
					{:else}
						<div class="overflow-x-auto">
							<table class="table table-zebra w-full">
								<thead>
									<tr>
										<th>Client</th>
										<th>Email</th>
										<th>Status</th>
										<th>Created</th>
									</tr>
								</thead>
								<tbody>
									{#each paginatedClients as client}
										<tr>
											<td>
												<div class="flex items-center gap-2">
													<Icon icon="mdi:account" class="text-lg" />
													<div>
														<div class="font-semibold">{client.username}</div>
														{#if client.friendlyName}
															<div class="text-sm opacity-70">{client.friendlyName}</div>
														{/if}
													</div>
												</div>
											</td>
											<td>{client.email}</td>
											<td>
												{#if client.isShared}
													<div class="badge badge-warning gap-1">
														<Icon icon="mdi:account-multiple" />
														Shared
													</div>
												{:else}
													<div class="badge badge-success gap-1">
														<Icon icon="mdi:account-check" />
														Exclusive
													</div>
												{/if}
											</td>
											<td>{formatDate(client.createdAt)}</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>

						{#if totalPages > 1}
							<div class="flex justify-center items-center gap-2 mt-4">
								<button
									class="btn btn-sm"
									disabled={!hasPrevPage}
									on:click={() => goToPage(currentPage - 1)}
								>
									<Icon icon="mdi:chevron-left" />
								</button>

								{#if totalPages <= 7}
									{#each Array(totalPages) as _, i}
										<button
											class="btn btn-sm {currentPage === i + 1 ? 'btn-primary' : ''}"
											on:click={() => goToPage(i + 1)}
										>
											{i + 1}
										</button>
									{/each}
								{:else}
									<button
										class="btn btn-sm {currentPage === 1 ? 'btn-primary' : ''}"
										on:click={() => goToPage(1)}
									>
										1
									</button>

									{#if currentPage > 3}
										<span class="px-2">...</span>
									{/if}

									{#each Array(5) as _, i}
										{@const page = Math.max(2, Math.min(totalPages - 1, currentPage - 2 + i))}
										{#if page > 1 && page < totalPages}
											<button
												class="btn btn-sm {currentPage === page ? 'btn-primary' : ''}"
												on:click={() => goToPage(page)}
											>
												{page}
											</button>
										{/if}
									{/each}

									{#if currentPage < totalPages - 2}
										<span class="px-2">...</span>
									{/if}

									<button
										class="btn btn-sm {currentPage === totalPages ? 'btn-primary' : ''}"
										on:click={() => goToPage(totalPages)}
									>
										{totalPages}
									</button>
								{/if}

								<button
									class="btn btn-sm"
									disabled={!hasNextPage}
									on:click={() => goToPage(currentPage + 1)}
								>
									<Icon icon="mdi:chevron-right" />
								</button>
							</div>

							<div class="text-center text-sm opacity-60 mt-2">
								Showing {(currentPage - 1) * itemsPerPage + 1} - {Math.min(
									currentPage * itemsPerPage,
									clients.length
								)} of {clients.length} clients
							</div>
						{/if}
					{/if}
				</Card>
			</div>
		</div>
	</div>
{/if}
