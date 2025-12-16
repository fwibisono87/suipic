<script lang="ts">
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { onMount } from 'svelte';
	import { LoadingSpinner } from '$lib/components';
	import { EUserRole } from '$lib/types';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
		}
	});
</script>

<svelte:head>
	<title>Dashboard - Suipic</title>
</svelte:head>

{#if $isAuthenticated}
	<div class="space-y-6">
		<div class="hero bg-base-200 rounded-lg">
			<div class="hero-content text-center py-12">
				<div class="max-w-md">
					<h1 class="text-5xl font-bold">
						Welcome back, {$currentUser?.friendlyName || $currentUser?.username}!
					</h1>
					<p class="py-6">Manage your photos and albums with ease.</p>
					<div class="flex gap-4 justify-center">
						<a href="/albums" class="btn btn-primary">View Albums</a>
						{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
							<a href="/albums/new" class="btn btn-outline">Create Album</a>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="stats shadow">
				<div class="stat">
					<div class="stat-figure text-primary">
						<Icon icon="mdi:image-multiple" class="text-4xl" />
					</div>
					<div class="stat-title">Total Albums</div>
					<div class="stat-value text-primary">-</div>
					<div class="stat-desc">View all your albums</div>
				</div>
			</div>

			<div class="stats shadow">
				<div class="stat">
					<div class="stat-figure text-secondary">
						<Icon icon="mdi:camera" class="text-4xl" />
					</div>
					<div class="stat-title">Total Photos</div>
					<div class="stat-value text-secondary">-</div>
					<div class="stat-desc">Across all albums</div>
				</div>
			</div>

			<div class="stats shadow">
				<div class="stat">
					<div class="stat-figure text-accent">
						<Icon icon="mdi:clock-outline" class="text-4xl" />
					</div>
					<div class="stat-title">Recent Activity</div>
					<div class="stat-value text-accent">-</div>
					<div class="stat-desc">Last 7 days</div>
				</div>
			</div>
		</div>

		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Quick Actions</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
					<button class="btn btn-outline" on:click={() => goto('/albums')}>
						<Icon icon="mdi:image-multiple" class="text-xl" />
						Browse Albums
					</button>

					<button class="btn btn-outline" on:click={() => goto('/search')}>
						<Icon icon="mdi:magnify" class="text-xl" />
						Search Photos
					</button>

					{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
						<button class="btn btn-outline" on:click={() => goto('/albums/new')}>
							<Icon icon="mdi:plus" class="text-xl" />
							New Album
						</button>
					{/if}

					<button class="btn btn-outline" on:click={() => goto('/profile')}>
						<Icon icon="mdi:account" class="text-xl" />
						My Profile
					</button>

					<button class="btn btn-outline" on:click={() => goto('/settings')}>
						<Icon icon="mdi:cog" class="text-xl" />
						Settings
					</button>
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="flex items-center justify-center min-h-[60vh]">
		<LoadingSpinner size="lg" />
	</div>
{/if}
