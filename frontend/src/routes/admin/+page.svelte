<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { Card, Alert, LoadingSpinner } from '$lib/components';
	import { useListPhotographers, useCreatePhotographer } from '$lib/queries/admin';
	import type { TCreatePhotographerRequest, TCreatePhotographerResponse } from '$lib/api';
	import { validateEmail, validateUsername } from '$lib/utils';

	let email = '';
	let username = '';
	let generatedCredentials: TCreatePhotographerResponse | null = null;
	let showCredentialsModal = false;

	let currentPage = 1;
	let itemsPerPage = 10;

	const photographersQuery = useListPhotographers();
	const createPhotographerMutation = useCreatePhotographer();

	$: photographers = $photographersQuery.data?.photographers || [];
	$: paginatedPhotographers = photographers.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
	$: totalPages = Math.ceil(photographers.length / itemsPerPage);
	$: hasNextPage = currentPage < totalPages;
	$: hasPrevPage = currentPage > 1;

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if ($currentUser?.role !== EUserRole.ADMIN) {
			goto('/');
			return;
		}
	});

	async function handleCreatePhotographer(e: Event) {
		e.preventDefault();

		if (!email || !validateEmail(email)) {
			return;
		}

		if (!username || !validateUsername(username)) {
			return;
		}

		const data: TCreatePhotographerRequest = {
			email,
			username
		};

		$createPhotographerMutation.mutate(data, {
			onSuccess: (response) => {
				generatedCredentials = response;
				showCredentialsModal = true;
				email = '';
				username = '';
			}
		});
	}

	function closeCredentialsModal() {
		showCredentialsModal = false;
		generatedCredentials = null;
	}

	function copyToClipboard(text: string) {
		if (typeof navigator !== 'undefined' && navigator.clipboard) {
			navigator.clipboard.writeText(text);
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
	<title>Admin Dashboard - Suipic</title>
</svelte:head>

{#if $isAuthenticated && $currentUser?.role === EUserRole.ADMIN}
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<h1 class="text-4xl font-bold">Admin Dashboard</h1>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<div class="lg:col-span-1">
				<Card title="Create Photographer Account">
					{#if $createPhotographerMutation.isError}
						<div class="mb-4">
							<Alert 
								type="error" 
								message={$createPhotographerMutation.error?.message || 'Failed to create photographer'} 
								dismissible 
								onDismiss={() => $createPhotographerMutation.reset()} 
							/>
						</div>
					{/if}

					{#if $createPhotographerMutation.isSuccess && !showCredentialsModal}
						<div class="mb-4">
							<Alert type="success" message="Photographer account created successfully!" />
						</div>
					{/if}

					<form on:submit={handleCreatePhotographer} class="space-y-4">
						<div class="form-control">
							<label class="label" for="email">
								<span class="label-text">Email</span>
							</label>
							<input
								type="email"
								id="email"
								name="email"
								bind:value={email}
								placeholder="photographer@example.com"
								class="input input-bordered w-full"
								disabled={$createPhotographerMutation.isPending}
								required
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
								placeholder="photographer_username"
								class="input input-bordered w-full"
								disabled={$createPhotographerMutation.isPending}
								required
							/>
							<label class="label">
								<span class="label-text-alt">3-50 characters</span>
							</label>
						</div>

						<button type="submit" class="btn btn-primary w-full" disabled={$createPhotographerMutation.isPending}>
							{#if $createPhotographerMutation.isPending}
								<LoadingSpinner size="sm" />
							{:else}
								<Icon icon="mdi:account-plus" class="text-lg" />
								Create Photographer
							{/if}
						</button>
					</form>

					<div class="mt-4 p-4 bg-info/10 rounded-lg">
						<p class="text-sm">
							<Icon icon="mdi:information" class="inline text-info" />
							A random password will be generated. Make sure to save the credentials displayed after creation.
						</p>
					</div>
				</Card>
			</div>

			<div class="lg:col-span-2">
				<Card title="Photographer Accounts">
					{#if $photographersQuery.isError}
						<div class="mb-4">
							<Alert type="error" message={$photographersQuery.error?.message || 'Failed to load photographers'} />
						</div>
					{/if}

					{#if $photographersQuery.isLoading}
						<div class="flex justify-center py-8">
							<LoadingSpinner />
						</div>
					{:else if photographers.length === 0}
						<div class="text-center py-8 opacity-60">
							<Icon icon="mdi:account-group" class="text-6xl mb-2" />
							<p>No photographer accounts yet</p>
						</div>
					{:else}
						<div class="overflow-x-auto">
							<table class="table table-zebra w-full">
								<thead>
									<tr>
										<th>Username</th>
										<th>Email</th>
										<th>Created</th>
									</tr>
								</thead>
								<tbody>
									{#each paginatedPhotographers as photographer}
										<tr>
											<td>
												<div class="flex items-center gap-2">
													<Icon icon="mdi:account" class="text-lg" />
													<span class="font-semibold">{photographer.username}</span>
												</div>
											</td>
											<td>{photographer.email}</td>
											<td>{formatDate(photographer.createdAt)}</td>
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
									photographers.length
								)} of {photographers.length} photographers
							</div>
						{/if}
					{/if}
				</Card>
			</div>
		</div>
	</div>
{/if}

{#if showCredentialsModal && generatedCredentials}
	<div class="modal modal-open">
		<div class="modal-box max-w-2xl">
			<h3 class="font-bold text-lg mb-4">
				<Icon icon="mdi:check-circle" class="inline text-success text-2xl" />
				Photographer Account Created
			</h3>

			<div class="alert alert-warning mb-4">
				<Icon icon="mdi:alert" />
				<span>Save these credentials now! The password will not be shown again.</span>
			</div>

			<div class="space-y-4">
				<div class="form-control">
					<label class="label">
						<span class="label-text font-semibold">Username</span>
					</label>
					<div class="join w-full">
						<input
							type="text"
							value={generatedCredentials.user.username}
							class="input input-bordered join-item flex-1"
							readonly
						/>
						<button
							class="btn join-item"
							on:click={() => copyToClipboard(generatedCredentials?.user.username || '')}
						>
							<Icon icon="mdi:content-copy" />
						</button>
					</div>
				</div>

				<div class="form-control">
					<label class="label">
						<span class="label-text font-semibold">Email</span>
					</label>
					<div class="join w-full">
						<input
							type="text"
							value={generatedCredentials.user.email}
							class="input input-bordered join-item flex-1"
							readonly
						/>
						<button
							class="btn join-item"
							on:click={() => copyToClipboard(generatedCredentials?.user.email || '')}
						>
							<Icon icon="mdi:content-copy" />
						</button>
					</div>
				</div>

				<div class="form-control">
					<label class="label">
						<span class="label-text font-semibold">Generated Password</span>
					</label>
					<div class="join w-full">
						<input
							type="text"
							value={generatedCredentials.password}
							class="input input-bordered join-item flex-1 font-mono"
							readonly
						/>
						<button
							class="btn join-item"
							on:click={() => copyToClipboard(generatedCredentials?.password || '')}
						>
							<Icon icon="mdi:content-copy" />
						</button>
					</div>
				</div>
			</div>

			<div class="modal-action">
				<button class="btn btn-primary" on:click={closeCredentialsModal}>
					I've Saved the Credentials
				</button>
			</div>
		</div>
	</div>
{/if}
